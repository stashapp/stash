package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

const (
	galleryTable = "galleries"

	galleriesFilesTable      = "galleries_files"
	performersGalleriesTable = "performers_galleries"
	galleriesTagsTable       = "galleries_tags"
	galleriesImagesTable     = "galleries_images"
	galleriesScenesTable     = "scenes_galleries"
	galleryIDColumn          = "gallery_id"
	galleriesURLsTable       = "gallery_urls"
	galleriesURLColumn       = "url"
)

type galleryRow struct {
	ID           int         `db:"id" goqu:"skipinsert"`
	Title        zero.String `db:"title"`
	Code         zero.String `db:"code"`
	Date         NullDate    `db:"date"`
	Details      zero.String `db:"details"`
	Photographer zero.String `db:"photographer"`
	// expressed as 1-100
	Rating    null.Int  `db:"rating"`
	Organized bool      `db:"organized"`
	StudioID  null.Int  `db:"studio_id,omitempty"`
	FolderID  null.Int  `db:"folder_id,omitempty"`
	CreatedAt Timestamp `db:"created_at"`
	UpdatedAt Timestamp `db:"updated_at"`
}

func (r *galleryRow) fromGallery(o models.Gallery) {
	r.ID = o.ID
	r.Title = zero.StringFrom(o.Title)
	r.Code = zero.StringFrom(o.Code)
	r.Date = NullDateFromDatePtr(o.Date)
	r.Details = zero.StringFrom(o.Details)
	r.Photographer = zero.StringFrom(o.Photographer)
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.StudioID = intFromPtr(o.StudioID)
	r.FolderID = nullIntFromFolderIDPtr(o.FolderID)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

type galleryQueryRow struct {
	galleryRow
	FolderPath            zero.String `db:"folder_path"`
	PrimaryFileID         null.Int    `db:"primary_file_id"`
	PrimaryFileFolderPath zero.String `db:"primary_file_folder_path"`
	PrimaryFileBasename   zero.String `db:"primary_file_basename"`
	PrimaryFileChecksum   zero.String `db:"primary_file_checksum"`
}

func (r *galleryQueryRow) resolve() *models.Gallery {
	ret := &models.Gallery{
		ID:            r.ID,
		Title:         r.Title.String,
		Code:          r.Code.String,
		Date:          r.Date.DatePtr(),
		Details:       r.Details.String,
		Photographer:  r.Photographer.String,
		Rating:        nullIntPtr(r.Rating),
		Organized:     r.Organized,
		StudioID:      nullIntPtr(r.StudioID),
		FolderID:      nullIntFolderIDPtr(r.FolderID),
		PrimaryFileID: nullIntFileIDPtr(r.PrimaryFileID),
		CreatedAt:     r.CreatedAt.Timestamp,
		UpdatedAt:     r.UpdatedAt.Timestamp,
	}

	if r.PrimaryFileFolderPath.Valid && r.PrimaryFileBasename.Valid {
		ret.Path = filepath.Join(r.PrimaryFileFolderPath.String, r.PrimaryFileBasename.String)
	} else if r.FolderPath.Valid {
		ret.Path = r.FolderPath.String
	}

	return ret
}

type galleryRowRecord struct {
	updateRecord
}

func (r *galleryRowRecord) fromPartial(o models.GalleryPartial) {
	r.setNullString("title", o.Title)
	r.setNullString("code", o.Code)
	r.setNullDate("date", o.Date)
	r.setNullString("details", o.Details)
	r.setNullString("photographer", o.Photographer)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setNullInt("studio_id", o.StudioID)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
}

type galleryRepositoryType struct {
	repository
	performers joinRepository
	images     joinRepository
	tags       joinRepository
	scenes     joinRepository
	files      filesRepository
}

func (r *galleryRepositoryType) addGalleriesFilesTable(f *filterBuilder) {
	f.addLeftJoin(galleriesFilesTable, "", "galleries_files.gallery_id = galleries.id")
}

func (r *galleryRepositoryType) addFilesTable(f *filterBuilder) {
	r.addGalleriesFilesTable(f)
	f.addLeftJoin(fileTable, "", "galleries_files.file_id = files.id")
}

func (r *galleryRepositoryType) addFoldersTable(f *filterBuilder) {
	r.addFilesTable(f)
	f.addLeftJoin(folderTable, "", "files.parent_folder_id = folders.id")
}

var (
	galleryRepository = galleryRepositoryType{
		repository: repository{
			tableName: galleryTable,
			idColumn:  idColumn,
		},
		performers: joinRepository{
			repository: repository{
				tableName: performersGalleriesTable,
				idColumn:  galleryIDColumn,
			},
			fkColumn: "performer_id",
		},
		tags: joinRepository{
			repository: repository{
				tableName: galleriesTagsTable,
				idColumn:  galleryIDColumn,
			},
			fkColumn:     "tag_id",
			foreignTable: tagTable,
			orderBy:      "tags.name ASC",
		},
		images: joinRepository{
			repository: repository{
				tableName: galleriesImagesTable,
				idColumn:  galleryIDColumn,
			},
			fkColumn: "image_id",
		},
		scenes: joinRepository{
			repository: repository{
				tableName: galleriesScenesTable,
				idColumn:  galleryIDColumn,
			},
			fkColumn: sceneIDColumn,
		},
		files: filesRepository{
			repository: repository{
				tableName: galleriesFilesTable,
				idColumn:  galleryIDColumn,
			},
		},
	}
)

type GalleryStore struct {
	tableMgr *table

	fileStore   *FileStore
	folderStore *FolderStore
}

func NewGalleryStore(fileStore *FileStore, folderStore *FolderStore) *GalleryStore {
	return &GalleryStore{
		tableMgr:    galleryTableMgr,
		fileStore:   fileStore,
		folderStore: folderStore,
	}
}

func (qb *GalleryStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *GalleryStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	files := fileTableMgr.table
	folders := folderTableMgr.table
	galleryFolder := folderTableMgr.table.As("gallery_folder")

	return dialect.From(table).LeftJoin(
		galleriesFilesJoinTable,
		goqu.On(
			galleriesFilesJoinTable.Col(galleryIDColumn).Eq(table.Col(idColumn)),
			galleriesFilesJoinTable.Col("primary").Eq(1),
		),
	).LeftJoin(
		files,
		goqu.On(files.Col(idColumn).Eq(galleriesFilesJoinTable.Col(fileIDColumn))),
	).LeftJoin(
		folders,
		goqu.On(folders.Col(idColumn).Eq(files.Col("parent_folder_id"))),
	).LeftJoin(
		galleryFolder,
		goqu.On(galleryFolder.Col(idColumn).Eq(table.Col("folder_id"))),
	).Select(
		qb.table().All(),
		galleriesFilesJoinTable.Col(fileIDColumn).As("primary_file_id"),
		folders.Col("path").As("primary_file_folder_path"),
		files.Col("basename").As("primary_file_basename"),
		galleryFolder.Col("path").As("folder_path"),
	)
}

func (qb *GalleryStore) Create(ctx context.Context, newObject *models.Gallery, fileIDs []models.FileID) error {
	var r galleryRow
	r.fromGallery(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if len(fileIDs) > 0 {
		const firstPrimary = true
		if err := galleriesFilesTableMgr.insertJoins(ctx, id, firstPrimary, fileIDs); err != nil {
			return err
		}
	}

	if newObject.URLs.Loaded() {
		const startPos = 0
		if err := galleriesURLsTableMgr.insertJoins(ctx, id, startPos, newObject.URLs.List()); err != nil {
			return err
		}
	}
	if newObject.PerformerIDs.Loaded() {
		if err := galleriesPerformersTableMgr.insertJoins(ctx, id, newObject.PerformerIDs.List()); err != nil {
			return err
		}
	}
	if newObject.TagIDs.Loaded() {
		if err := galleriesTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs.List()); err != nil {
			return err
		}
	}
	if newObject.SceneIDs.Loaded() {
		if err := galleriesScenesTableMgr.insertJoins(ctx, id, newObject.SceneIDs.List()); err != nil {
			return err
		}
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated

	return nil
}

func (qb *GalleryStore) Update(ctx context.Context, updatedObject *models.Gallery) error {
	var r galleryRow
	r.fromGallery(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.URLs.Loaded() {
		if err := galleriesURLsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.URLs.List()); err != nil {
			return err
		}
	}
	if updatedObject.PerformerIDs.Loaded() {
		if err := galleriesPerformersTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.PerformerIDs.List()); err != nil {
			return err
		}
	}
	if updatedObject.TagIDs.Loaded() {
		if err := galleriesTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs.List()); err != nil {
			return err
		}
	}
	if updatedObject.SceneIDs.Loaded() {
		if err := galleriesScenesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.SceneIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.Files.Loaded() {
		fileIDs := make([]models.FileID, len(updatedObject.Files.List()))
		for i, f := range updatedObject.Files.List() {
			fileIDs[i] = f.Base().ID
		}

		if err := galleriesFilesTableMgr.replaceJoins(ctx, updatedObject.ID, fileIDs); err != nil {
			return err
		}
	}

	return nil
}

func (qb *GalleryStore) UpdatePartial(ctx context.Context, id int, partial models.GalleryPartial) (*models.Gallery, error) {
	r := galleryRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(partial)

	if len(r.Record) > 0 {
		if err := qb.tableMgr.updateByID(ctx, id, r.Record); err != nil {
			return nil, err
		}
	}

	if partial.URLs != nil {
		if err := galleriesURLsTableMgr.modifyJoins(ctx, id, partial.URLs.Values, partial.URLs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.PerformerIDs != nil {
		if err := galleriesPerformersTableMgr.modifyJoins(ctx, id, partial.PerformerIDs.IDs, partial.PerformerIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.TagIDs != nil {
		if err := galleriesTagsTableMgr.modifyJoins(ctx, id, partial.TagIDs.IDs, partial.TagIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.SceneIDs != nil {
		if err := galleriesScenesTableMgr.modifyJoins(ctx, id, partial.SceneIDs.IDs, partial.SceneIDs.Mode); err != nil {
			return nil, err
		}
	}

	if partial.PrimaryFileID != nil {
		if err := galleriesFilesTableMgr.setPrimary(ctx, id, *partial.PrimaryFileID); err != nil {
			return nil, err
		}
	}

	return qb.find(ctx, id)
}

func (qb *GalleryStore) Destroy(ctx context.Context, id int) error {
	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

func (qb *GalleryStore) GetFiles(ctx context.Context, id int) ([]models.File, error) {
	fileIDs, err := galleryRepository.files.get(ctx, id)
	if err != nil {
		return nil, err
	}

	// use fileStore to load files
	files, err := qb.fileStore.Find(ctx, fileIDs...)
	if err != nil {
		return nil, err
	}

	ret := make([]models.File, len(files))
	copy(ret, files)

	return ret, nil
}

func (qb *GalleryStore) GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error) {
	const primaryOnly = false
	return galleryRepository.files.getMany(ctx, ids, primaryOnly)
}

// returns nil, nil if not found
func (qb *GalleryStore) Find(ctx context.Context, id int) (*models.Gallery, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *GalleryStore) FindMany(ctx context.Context, ids []int) ([]*models.Gallery, error) {
	galleries := make([]*models.Gallery, len(ids))

	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := qb.selectDataset().Prepared(true).Where(qb.table().Col(idColumn).In(batch))
		unsorted, err := qb.getMany(ctx, q)
		if err != nil {
			return err
		}

		for _, s := range unsorted {
			i := slices.Index(ids, s.ID)
			galleries[i] = s
		}

		return nil
	}); err != nil {
		return nil, err
	}

	for i := range galleries {
		if galleries[i] == nil {
			return nil, fmt.Errorf("gallery with id %d not found", ids[i])
		}
	}

	return galleries, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *GalleryStore) find(ctx context.Context, id int) (*models.Gallery, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *GalleryStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Gallery, error) {
	table := qb.table()

	q := qb.selectDataset().Prepared(true).Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

// returns nil, sql.ErrNoRows if not found
func (qb *GalleryStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Gallery, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *GalleryStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Gallery, error) {
	const single = false
	var ret []*models.Gallery
	var lastID int
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f galleryQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()

		if s.ID == lastID {
			return fmt.Errorf("internal error: multiple rows returned for single gallery id %d", s.ID)
		}
		lastID = s.ID

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *GalleryStore) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Gallery, error) {
	sq := dialect.From(galleriesFilesJoinTable).Select(galleriesFilesJoinTable.Col(galleryIDColumn)).Where(
		galleriesFilesJoinTable.Col(fileIDColumn).Eq(fileID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting gallery by file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *GalleryStore) CountByFileID(ctx context.Context, fileID models.FileID) (int, error) {
	joinTable := galleriesFilesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(fileIDColumn).Eq(fileID))
	return count(ctx, q)
}

func (qb *GalleryStore) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Gallery, error) {
	fingerprintTable := fingerprintTableMgr.table

	var ex []exp.Expression

	for _, v := range fp {
		ex = append(ex, goqu.And(
			fingerprintTable.Col("type").Eq(v.Type),
			fingerprintTable.Col("fingerprint").Eq(v.Fingerprint),
		))
	}

	sq := dialect.From(galleriesFilesJoinTable).
		InnerJoin(
			fingerprintTable,
			goqu.On(fingerprintTable.Col(fileIDColumn).Eq(galleriesFilesJoinTable.Col(fileIDColumn))),
		).
		Select(galleriesFilesJoinTable.Col(galleryIDColumn)).Where(goqu.Or(ex...))

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting gallery by fingerprints: %w", err)
	}

	return ret, nil
}

func (qb *GalleryStore) FindByChecksum(ctx context.Context, checksum string) ([]*models.Gallery, error) {
	return qb.FindByFingerprints(ctx, []models.Fingerprint{
		{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: checksum,
		},
	})
}

func (qb *GalleryStore) FindByChecksums(ctx context.Context, checksums []string) ([]*models.Gallery, error) {
	fingerprints := make([]models.Fingerprint, len(checksums))

	for i, c := range checksums {
		fingerprints[i] = models.Fingerprint{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: c,
		}
	}
	return qb.FindByFingerprints(ctx, fingerprints)
}

func (qb *GalleryStore) FindByPath(ctx context.Context, p string) ([]*models.Gallery, error) {
	table := qb.table()
	filesTable := fileTableMgr.table
	fileFoldersTable := folderTableMgr.table.As("file_folders")
	foldersTable := folderTableMgr.table

	basename := filepath.Base(p)
	dir := filepath.Dir(p)

	sq := dialect.From(table).LeftJoin(
		galleriesFilesJoinTable,
		goqu.On(galleriesFilesJoinTable.Col(galleryIDColumn).Eq(table.Col(idColumn))),
	).LeftJoin(
		filesTable,
		goqu.On(filesTable.Col(idColumn).Eq(galleriesFilesJoinTable.Col(fileIDColumn))),
	).LeftJoin(
		fileFoldersTable,
		goqu.On(fileFoldersTable.Col(idColumn).Eq(filesTable.Col("parent_folder_id"))),
	).LeftJoin(
		foldersTable,
		goqu.On(foldersTable.Col(idColumn).Eq(table.Col("folder_id"))),
	).Select(table.Col(idColumn)).Where(
		goqu.Or(
			goqu.And(
				fileFoldersTable.Col("path").Eq(dir),
				filesTable.Col("basename").Eq(basename),
			),
			foldersTable.Col("path").Eq(p),
		),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting gallery by path %s: %w", p, err)
	}

	return ret, nil
}

func (qb *GalleryStore) FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Gallery, error) {
	table := qb.table()

	sq := dialect.From(table).Select(table.Col(idColumn)).Where(
		table.Col("folder_id").Eq(folderID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting galleries for folder %d: %w", folderID, err)
	}

	return ret, nil
}

func (qb *GalleryStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.Gallery, error) {
	sq := dialect.From(galleriesScenesJoinTable).Select(galleriesScenesJoinTable.Col(galleryIDColumn)).Where(
		galleriesScenesJoinTable.Col(sceneIDColumn).Eq(sceneID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting galleries for scene %d: %w", sceneID, err)
	}

	return ret, nil
}

func (qb *GalleryStore) FindByImageID(ctx context.Context, imageID int) ([]*models.Gallery, error) {
	sq := dialect.From(galleriesImagesJoinTable).Select(galleriesImagesJoinTable.Col(galleryIDColumn)).Where(
		galleriesImagesJoinTable.Col(imageIDColumn).Eq(imageID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting galleries for image %d: %w", imageID, err)
	}

	return ret, nil
}

func (qb *GalleryStore) CountByImageID(ctx context.Context, imageID int) (int, error) {
	joinTable := galleriesImagesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(imageIDColumn).Eq(imageID))
	return count(ctx, q)
}

func (qb *GalleryStore) FindUserGalleryByTitle(ctx context.Context, title string) ([]*models.Gallery, error) {
	table := qb.table()

	sq := dialect.From(table).LeftJoin(
		galleriesFilesJoinTable,
		goqu.On(galleriesFilesJoinTable.Col(galleryIDColumn).Eq(table.Col(idColumn))),
	).Select(table.Col(idColumn)).Where(
		table.Col("folder_id").IsNull(),
		galleriesFilesJoinTable.Col("file_id").IsNull(),
		table.Col("title").Eq(title),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting user galleries for title %s: %w", title, err)
	}

	return ret, nil
}

func (qb *GalleryStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *GalleryStore) All(ctx context.Context) ([]*models.Gallery, error) {
	return qb.getMany(ctx, qb.selectDataset())
}

func (qb *GalleryStore) makeQuery(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if galleryFilter == nil {
		galleryFilter = &models.GalleryFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := galleryRepository.newQuery()
	distinctIDs(&query, galleryTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.addJoins(
			join{
				table:    galleriesFilesTable,
				onClause: "galleries_files.gallery_id = galleries.id",
			},
			join{
				table:    fileTable,
				onClause: "galleries_files.file_id = files.id",
			},
			join{
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			},
			join{
				table:    fingerprintTable,
				onClause: "files_fingerprints.file_id = galleries_files.file_id",
			},
			join{
				table:    folderTable,
				as:       "gallery_folder",
				onClause: "galleries.folder_id = gallery_folder.id",
			},
			join{
				table:    galleriesChaptersTable,
				onClause: "galleries_chapters.gallery_id = galleries.id",
			},
		)

		// add joins for files and checksum
		filepathColumn := "folders.path || '" + string(filepath.Separator) + "' || files.basename"
		searchColumns := []string{"galleries.title", "gallery_folder.path", filepathColumn, "files_fingerprints.fingerprint", "galleries_chapters.title"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := filterBuilderFromHandler(ctx, &galleryFilterHandler{
		galleryFilter: galleryFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setGallerySort(&query, findFilter); err != nil {
		return nil, err
	}
	query.sortAndPagination += getPagination(findFilter)

	return &query, nil
}

func (qb *GalleryStore) Query(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) ([]*models.Gallery, int, error) {
	query, err := qb.makeQuery(ctx, galleryFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	galleries, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return galleries, countResult, nil
}

func (qb *GalleryStore) QueryCount(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, galleryFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

var gallerySortOptions = sortOptions{
	"created_at",
	"date",
	"file_count",
	"file_mod_time",
	"id",
	"images_count",
	"path",
	"performer_count",
	"random",
	"rating",
	"tag_count",
	"title",
	"updated_at",
}

func (qb *GalleryStore) setGallerySort(query *queryBuilder, findFilter *models.FindFilterType) error {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return nil
	}

	sort := findFilter.GetSort("path")
	direction := findFilter.GetDirection()

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := gallerySortOptions.validateSort(sort); err != nil {
		return err
	}

	addFileTable := func() {
		query.addJoins(
			join{
				table:    galleriesFilesTable,
				onClause: "galleries_files.gallery_id = galleries.id",
			},
			join{
				table:    fileTable,
				onClause: "galleries_files.file_id = files.id",
			},
		)
	}

	addFolderTable := func() {
		query.addJoins(
			join{
				table:    folderTable,
				onClause: "folders.id = galleries.folder_id",
			},
			join{
				table:    folderTable,
				as:       "file_folder",
				onClause: "files.parent_folder_id = file_folder.id",
			},
		)
	}

	switch sort {
	case "file_count":
		query.sortAndPagination += getCountSort(galleryTable, galleriesFilesTable, galleryIDColumn, direction)
	case "images_count":
		query.sortAndPagination += getCountSort(galleryTable, galleriesImagesTable, galleryIDColumn, direction)
	case "tag_count":
		query.sortAndPagination += getCountSort(galleryTable, galleriesTagsTable, galleryIDColumn, direction)
	case "performer_count":
		query.sortAndPagination += getCountSort(galleryTable, performersGalleriesTable, galleryIDColumn, direction)
	case "path":
		// special handling for path
		addFileTable()
		addFolderTable()
		query.sortAndPagination += fmt.Sprintf(" ORDER BY COALESCE(folders.path, '') || COALESCE(file_folder.path, '') || COALESCE(files.basename, '') COLLATE NATURAL_CI %s", direction)
	case "file_mod_time":
		sort = "mod_time"
		addFileTable()
		query.sortAndPagination += getSort(sort, direction, fileTable)
	case "title":
		addFileTable()
		addFolderTable()
		query.sortAndPagination += " ORDER BY COALESCE(galleries.title, files.basename, basename(COALESCE(folders.path, ''))) COLLATE NATURAL_CI " + direction + ", file_folder.path COLLATE NATURAL_CI " + direction
	default:
		query.sortAndPagination += getSort(sort, direction, "galleries")
	}

	// Whatever the sorting, always use title/id as a final sort
	query.sortAndPagination += ", COALESCE(galleries.title, galleries.id) COLLATE NATURAL_CI ASC"

	return nil
}

func (qb *GalleryStore) GetURLs(ctx context.Context, galleryID int) ([]string, error) {
	return galleriesURLsTableMgr.get(ctx, galleryID)
}

func (qb *GalleryStore) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	const firstPrimary = false
	return galleriesFilesTableMgr.insertJoins(ctx, id, firstPrimary, []models.FileID{fileID})
}

func (qb *GalleryStore) GetPerformerIDs(ctx context.Context, id int) ([]int, error) {
	return galleryRepository.performers.getIDs(ctx, id)
}

func (qb *GalleryStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return galleryRepository.tags.getIDs(ctx, id)
}

func (qb *GalleryStore) GetImageIDs(ctx context.Context, galleryID int) ([]int, error) {
	return galleryRepository.images.getIDs(ctx, galleryID)
}

func (qb *GalleryStore) AddImages(ctx context.Context, galleryID int, imageIDs ...int) error {
	return galleryRepository.images.insertOrIgnore(ctx, galleryID, imageIDs...)
}

func (qb *GalleryStore) RemoveImages(ctx context.Context, galleryID int, imageIDs ...int) error {
	return galleryRepository.images.destroyJoins(ctx, galleryID, imageIDs...)
}

func (qb *GalleryStore) UpdateImages(ctx context.Context, galleryID int, imageIDs []int) error {
	// Delete the existing joins and then create new ones
	return galleryRepository.images.replace(ctx, galleryID, imageIDs)
}

func (qb *GalleryStore) SetCover(ctx context.Context, galleryID int, coverImageID int) error {
	return imageGalleriesTableMgr.setCover(ctx, coverImageID, galleryID)
}

func (qb *GalleryStore) ResetCover(ctx context.Context, galleryID int) error {
	return imageGalleriesTableMgr.resetCover(ctx, galleryID)
}

func (qb *GalleryStore) GetSceneIDs(ctx context.Context, id int) ([]int, error) {
	return galleryRepository.scenes.getIDs(ctx, id)
}
