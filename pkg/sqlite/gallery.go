package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
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
	galleriesChaptersTable   = "galleries_chapters"
	galleryIDColumn          = "gallery_id"
)

type galleryRow struct {
	ID      int               `db:"id" goqu:"skipinsert"`
	Title   zero.String       `db:"title"`
	URL     zero.String       `db:"url"`
	Date    models.SQLiteDate `db:"date"`
	Details zero.String       `db:"details"`
	// expressed as 1-100
	Rating    null.Int               `db:"rating"`
	Organized bool                   `db:"organized"`
	StudioID  null.Int               `db:"studio_id,omitempty"`
	FolderID  null.Int               `db:"folder_id,omitempty"`
	CreatedAt models.SQLiteTimestamp `db:"created_at"`
	UpdatedAt models.SQLiteTimestamp `db:"updated_at"`
}

func (r *galleryRow) fromGallery(o models.Gallery) {
	r.ID = o.ID
	r.Title = zero.StringFrom(o.Title)
	r.URL = zero.StringFrom(o.URL)
	if o.Date != nil {
		_ = r.Date.Scan(o.Date.Time)
	}
	r.Details = zero.StringFrom(o.Details)
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.StudioID = intFromPtr(o.StudioID)
	r.FolderID = nullIntFromFolderIDPtr(o.FolderID)
	r.CreatedAt = models.SQLiteTimestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = models.SQLiteTimestamp{Timestamp: o.UpdatedAt}
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
		URL:           r.URL.String,
		Date:          r.Date.DatePtr(),
		Details:       r.Details.String,
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
	r.setNullString("url", o.URL)
	r.setSQLiteDate("date", o.Date)
	r.setNullString("details", o.Details)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setNullInt("studio_id", o.StudioID)
	r.setSQLiteTimestamp("created_at", o.CreatedAt)
	r.setSQLiteTimestamp("updated_at", o.UpdatedAt)
}

type GalleryStore struct {
	repository

	tableMgr *table

	fileStore   *FileStore
	folderStore *FolderStore
}

func NewGalleryStore(fileStore *FileStore, folderStore *FolderStore) *GalleryStore {
	return &GalleryStore{
		repository: repository{
			tableName: galleryTable,
			idColumn:  idColumn,
		},
		tableMgr:    galleryTableMgr,
		fileStore:   fileStore,
		folderStore: folderStore,
	}
}

func (qb *GalleryStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *GalleryStore) Create(ctx context.Context, newObject *models.Gallery, fileIDs []file.ID) error {
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

	updated, err := qb.Find(ctx, id)
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
		fileIDs := make([]file.ID, len(updatedObject.Files.List()))
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

	return qb.Find(ctx, id)
}

func (qb *GalleryStore) Destroy(ctx context.Context, id int) error {
	return qb.tableMgr.destroyExisting(ctx, []int{id})
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

func (qb *GalleryStore) GetFiles(ctx context.Context, id int) ([]file.File, error) {
	fileIDs, err := qb.filesRepository().get(ctx, id)
	if err != nil {
		return nil, err
	}

	// use fileStore to load files
	files, err := qb.fileStore.Find(ctx, fileIDs...)
	if err != nil {
		return nil, err
	}

	ret := make([]file.File, len(files))
	copy(ret, files)

	return ret, nil
}

func (qb *GalleryStore) GetManyFileIDs(ctx context.Context, ids []int) ([][]file.ID, error) {
	const primaryOnly = false
	return qb.filesRepository().getMany(ctx, ids, primaryOnly)
}

func (qb *GalleryStore) Find(ctx context.Context, id int) (*models.Gallery, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting gallery by id %d: %w", id, err)
	}

	return ret, nil
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
			i := intslice.IntIndex(ids, s.ID)
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

func (qb *GalleryStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Gallery, error) {
	table := qb.table()

	q := qb.selectDataset().Prepared(true).Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

func (qb *GalleryStore) FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Gallery, error) {
	sq := dialect.From(galleriesFilesJoinTable).Select(galleriesFilesJoinTable.Col(galleryIDColumn)).Where(
		galleriesFilesJoinTable.Col(fileIDColumn).Eq(fileID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting gallery by file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *GalleryStore) CountByFileID(ctx context.Context, fileID file.ID) (int, error) {
	joinTable := galleriesFilesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(fileIDColumn).Eq(fileID))
	return count(ctx, q)
}

func (qb *GalleryStore) FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Gallery, error) {
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
	return qb.FindByFingerprints(ctx, []file.Fingerprint{
		{
			Type:        file.FingerprintTypeMD5,
			Fingerprint: checksum,
		},
	})
}

func (qb *GalleryStore) FindByChecksums(ctx context.Context, checksums []string) ([]*models.Gallery, error) {
	fingerprints := make([]file.Fingerprint, len(checksums))

	for i, c := range checksums {
		fingerprints[i] = file.Fingerprint{
			Type:        file.FingerprintTypeMD5,
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

func (qb *GalleryStore) FindByFolderID(ctx context.Context, folderID file.FolderID) ([]*models.Gallery, error) {
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

func (qb *GalleryStore) validateFilter(galleryFilter *models.GalleryFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if galleryFilter.And != nil {
		if galleryFilter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if galleryFilter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(galleryFilter.And)
	}

	if galleryFilter.Or != nil {
		if galleryFilter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(galleryFilter.Or)
	}

	if galleryFilter.Not != nil {
		return qb.validateFilter(galleryFilter.Not)
	}

	return nil
}

func (qb *GalleryStore) makeFilter(ctx context.Context, galleryFilter *models.GalleryFilterType) *filterBuilder {
	query := &filterBuilder{}

	if galleryFilter.And != nil {
		query.and(qb.makeFilter(ctx, galleryFilter.And))
	}
	if galleryFilter.Or != nil {
		query.or(qb.makeFilter(ctx, galleryFilter.Or))
	}
	if galleryFilter.Not != nil {
		query.not(qb.makeFilter(ctx, galleryFilter.Not))
	}

	query.handleCriterion(ctx, intCriterionHandler(galleryFilter.ID, "galleries.id", nil))
	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.Title, "galleries.title"))
	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.Details, "galleries.details"))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if galleryFilter.Checksum != nil {
			qb.addGalleriesFilesTable(f)
			f.addLeftJoin(fingerprintTable, "fingerprints_md5", "galleries_files.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
		}

		stringCriterionHandler(galleryFilter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
	}))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if galleryFilter.IsZip != nil {
			qb.addGalleriesFilesTable(f)
			if *galleryFilter.IsZip {

				f.addWhere("galleries_files.file_id IS NOT NULL")
			} else {
				f.addWhere("galleries_files.file_id IS NULL")
			}
		}
	}))

	query.handleCriterion(ctx, qb.galleryPathCriterionHandler(galleryFilter.Path))
	query.handleCriterion(ctx, galleryFileCountCriterionHandler(qb, galleryFilter.FileCount))
	query.handleCriterion(ctx, intCriterionHandler(galleryFilter.Rating100, "galleries.rating", nil))
	// legacy rating handler
	query.handleCriterion(ctx, rating5CriterionHandler(galleryFilter.Rating, "galleries.rating", nil))
	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.URL, "galleries.url"))
	query.handleCriterion(ctx, boolCriterionHandler(galleryFilter.Organized, "galleries.organized", nil))
	query.handleCriterion(ctx, galleryIsMissingCriterionHandler(qb, galleryFilter.IsMissing))
	query.handleCriterion(ctx, galleryTagsCriterionHandler(qb, galleryFilter.Tags))
	query.handleCriterion(ctx, galleryTagCountCriterionHandler(qb, galleryFilter.TagCount))
	query.handleCriterion(ctx, galleryPerformersCriterionHandler(qb, galleryFilter.Performers))
	query.handleCriterion(ctx, galleryPerformerCountCriterionHandler(qb, galleryFilter.PerformerCount))
	query.handleCriterion(ctx, hasChaptersCriterionHandler(galleryFilter.HasChapters))
	query.handleCriterion(ctx, galleryStudioCriterionHandler(qb, galleryFilter.Studios))
	query.handleCriterion(ctx, galleryPerformerTagsCriterionHandler(qb, galleryFilter.PerformerTags))
	query.handleCriterion(ctx, galleryAverageResolutionCriterionHandler(qb, galleryFilter.AverageResolution))
	query.handleCriterion(ctx, galleryImageCountCriterionHandler(qb, galleryFilter.ImageCount))
	query.handleCriterion(ctx, galleryPerformerFavoriteCriterionHandler(galleryFilter.PerformerFavorite))
	query.handleCriterion(ctx, galleryPerformerAgeCriterionHandler(galleryFilter.PerformerAge))
	query.handleCriterion(ctx, dateCriterionHandler(galleryFilter.Date, "galleries.date"))
	query.handleCriterion(ctx, timestampCriterionHandler(galleryFilter.CreatedAt, "galleries.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(galleryFilter.UpdatedAt, "galleries.updated_at"))

	return query
}

func (qb *GalleryStore) addGalleriesFilesTable(f *filterBuilder) {
	f.addLeftJoin(galleriesFilesTable, "", "galleries_files.gallery_id = galleries.id")
}

func (qb *GalleryStore) addFilesTable(f *filterBuilder) {
	qb.addGalleriesFilesTable(f)
	f.addLeftJoin(fileTable, "", "galleries_files.file_id = files.id")
}

func (qb *GalleryStore) addFoldersTable(f *filterBuilder) {
	qb.addFilesTable(f)
	f.addLeftJoin(folderTable, "", "files.parent_folder_id = folders.id")
}

func (qb *GalleryStore) makeQuery(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if galleryFilter == nil {
		galleryFilter = &models.GalleryFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
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

	if err := qb.validateFilter(galleryFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, galleryFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	qb.setGallerySort(&query, findFilter)
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

	var galleries []*models.Gallery
	for _, id := range idsResult {
		gallery, err := qb.Find(ctx, id)
		if err != nil {
			return nil, 0, err
		}

		galleries = append(galleries, gallery)
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

func (qb *GalleryStore) galleryPathCriterionHandler(c *models.StringCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if c != nil {
			qb.addFoldersTable(f)
			f.addLeftJoin(folderTable, "gallery_folder", "galleries.folder_id = gallery_folder.id")

			const pathColumn = "folders.path"
			const basenameColumn = "files.basename"
			const folderPathColumn = "gallery_folder.path"

			addWildcards := true
			not := false

			if modifier := c.Modifier; c.Modifier.IsValid() {
				switch modifier {
				case models.CriterionModifierIncludes:
					clause := getPathSearchClauseMany(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := getStringSearchClause([]string{folderPathColumn}, c.Value, false)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierExcludes:
					not = true
					clause := getPathSearchClauseMany(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := getStringSearchClause([]string{folderPathColumn}, c.Value, true)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierEquals:
					addWildcards = false
					clause := getPathSearchClause(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := makeClause(folderPathColumn+" LIKE ?", c.Value)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierNotEquals:
					addWildcards = false
					not = true
					clause := getPathSearchClause(pathColumn, basenameColumn, c.Value, addWildcards, not)
					clause2 := makeClause(folderPathColumn+" NOT LIKE ?", c.Value)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					filepathColumn := fmt.Sprintf("%s || '%s' || %s", pathColumn, string(filepath.Separator), basenameColumn)
					clause := makeClause(fmt.Sprintf("%s IS NOT NULL AND %s IS NOT NULL AND %s regexp ?", pathColumn, basenameColumn, filepathColumn), c.Value)
					clause2 := makeClause(fmt.Sprintf("%s IS NOT NULL AND %[1]s regexp ?", folderPathColumn), c.Value)
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				case models.CriterionModifierNotMatchesRegex:
					if _, err := regexp.Compile(c.Value); err != nil {
						f.setError(err)
						return
					}
					filepathColumn := fmt.Sprintf("%s || '%s' || %s", pathColumn, string(filepath.Separator), basenameColumn)
					f.addWhere(fmt.Sprintf("%s IS NULL OR %s IS NULL OR %s NOT regexp ?", pathColumn, basenameColumn, filepathColumn), c.Value)
					f.addWhere(fmt.Sprintf("%s IS NULL OR %[1]s NOT regexp ?", folderPathColumn), c.Value)
				case models.CriterionModifierIsNull:
					f.addWhere(fmt.Sprintf("%s IS NULL OR TRIM(%[1]s) = '' OR %s IS NULL OR TRIM(%[2]s) = ''", pathColumn, basenameColumn))
					f.addWhere(fmt.Sprintf("%s IS NULL OR TRIM(%[1]s) = ''", folderPathColumn))
				case models.CriterionModifierNotNull:
					clause := makeClause(fmt.Sprintf("%s IS NOT NULL AND TRIM(%[1]s) != '' AND %s IS NOT NULL AND TRIM(%[2]s) != ''", pathColumn, basenameColumn))
					clause2 := makeClause(fmt.Sprintf("%s IS NOT NULL AND TRIM(%[1]s) != ''", folderPathColumn))
					f.whereClauses = append(f.whereClauses, orClauses(clause, clause2))
				default:
					panic("unsupported string filter modifier")
				}
			}
		}
	}
}

func galleryFileCountCriterionHandler(qb *GalleryStore, fileCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesFilesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(fileCount)
}

func galleryIsMissingCriterionHandler(qb *GalleryStore, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "scenes":
				f.addLeftJoin("scenes_galleries", "scenes_join", "scenes_join.gallery_id = galleries.id")
				f.addWhere("scenes_join.gallery_id IS NULL")
			case "studio":
				f.addWhere("galleries.studio_id IS NULL")
			case "performers":
				qb.performersRepository().join(f, "performers_join", "galleries.id")
				f.addWhere("performers_join.gallery_id IS NULL")
			case "date":
				f.addWhere("galleries.date IS NULL OR galleries.date IS \"\" OR galleries.date IS \"0001-01-01\"")
			case "tags":
				qb.tagsRepository().join(f, "tags_join", "galleries.id")
				f.addWhere("tags_join.gallery_id IS NULL")
			default:
				f.addWhere("(galleries." + *isMissing + " IS NULL OR TRIM(galleries." + *isMissing + ") = '')")
			}
		}
	}
}

func galleryTagsCriterionHandler(qb *GalleryStore, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: galleryTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "image_tag",
		joinTable:      galleriesTagsTable,
		primaryFK:      galleryIDColumn,
	}

	return h.handler(tags)
}

func galleryTagCountCriterionHandler(qb *GalleryStore, tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesTagsTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(tagCount)
}

func galleryPerformersCriterionHandler(qb *GalleryStore, performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    performersGalleriesTable,
		joinAs:       "performers_join",
		primaryFK:    galleryIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			qb.performersRepository().join(f, "performers_join", "galleries.id")
		},
	}

	return h.handler(performers)
}

func galleryPerformerCountCriterionHandler(qb *GalleryStore, performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    performersGalleriesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(performerCount)
}

func galleryImageCountCriterionHandler(qb *GalleryStore, imageCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesImagesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(imageCount)
}

func hasChaptersCriterionHandler(hasChapters *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if hasChapters != nil {
			f.addLeftJoin("galleries_chapters", "", "galleries_chapters.gallery_id = galleries.id")
			if *hasChapters == "true" {
				f.addHaving("count(galleries_chapters.gallery_id) > 0")
			} else {
				f.addWhere("galleries_chapters.id IS NULL")
			}
		}
	}
}

func galleryStudioCriterionHandler(qb *GalleryStore, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := hierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: galleryTable,
		foreignTable: studioTable,
		foreignFK:    studioIDColumn,
		parentFK:     "parent_id",
	}

	return h.handler(studios)
}

func galleryPerformerTagsCriterionHandler(qb *GalleryStore, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("performers_galleries", "", "galleries.id = performers_galleries.gallery_id")
				f.addLeftJoin("performers_tags", "", "performers_galleries.performer_id = performers_tags.performer_id")

				f.addWhere(fmt.Sprintf("performers_tags.tag_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 {
				return
			}

			valuesClause := getHierarchicalValues(ctx, qb.tx, tags.Value, tagTable, "tags_relations", "", tags.Depth)

			f.addWith(`performer_tags AS (
SELECT pg.gallery_id, t.column1 AS root_tag_id FROM performers_galleries pg
INNER JOIN performers_tags pt ON pt.performer_id = pg.performer_id
INNER JOIN (` + valuesClause + `) t ON t.column2 = pt.tag_id
)`)

			f.addLeftJoin("performer_tags", "", "performer_tags.gallery_id = galleries.id")

			addHierarchicalConditionClauses(f, tags, "performer_tags", "root_tag_id")
		}
	}
}

func galleryPerformerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_galleries", "", "galleries.id = performers_galleries.gallery_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_galleries.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_galleries.gallery_id as id FROM performers_galleries 
JOIN performers ON performers.id = performers_galleries.performer_id
GROUP BY performers_galleries.gallery_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "galleries.id = nofaves.id")
				f.addWhere("performers_galleries.gallery_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func galleryPerformerAgeCriterionHandler(performerAge *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerAge != nil {
			f.addInnerJoin("performers_galleries", "", "galleries.id = performers_galleries.gallery_id")
			f.addInnerJoin("performers", "", "performers_galleries.performer_id = performers.id")

			f.addWhere("galleries.date != '' AND performers.birthdate != ''")
			f.addWhere("galleries.date IS NOT NULL AND performers.birthdate IS NOT NULL")
			f.addWhere("galleries.date != '0001-01-01' AND performers.birthdate != '0001-01-01'")

			ageCalc := "cast(strftime('%Y.%m%d', galleries.date) - strftime('%Y.%m%d', performers.birthdate) as int)"
			whereClause, args := getIntWhereClause(ageCalc, performerAge.Modifier, performerAge.Value, performerAge.Value2)
			f.addWhere(whereClause, args...)
		}
	}
}

func galleryAverageResolutionCriterionHandler(qb *GalleryStore, resolution *models.ResolutionCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if resolution != nil && resolution.Value.IsValid() {
			qb.imagesRepository().join(f, "images_join", "galleries.id")
			f.addLeftJoin("images", "", "images_join.image_id = images.id")
			f.addLeftJoin("images_files", "", "images.id = images_files.image_id")
			f.addLeftJoin("image_files", "", "images_files.file_id = image_files.file_id")

			min := resolution.Value.GetMinResolution()
			max := resolution.Value.GetMaxResolution()

			const widthHeight = "avg(MIN(image_files.width, image_files.height))"

			switch resolution.Modifier {
			case models.CriterionModifierEquals:
				f.addHaving(fmt.Sprintf("%s BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierNotEquals:
				f.addHaving(fmt.Sprintf("%s NOT BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierLessThan:
				f.addHaving(fmt.Sprintf("%s < %d", widthHeight, min))
			case models.CriterionModifierGreaterThan:
				f.addHaving(fmt.Sprintf("%s > %d", widthHeight, max))
			}
		}
	}
}

func (qb *GalleryStore) setGallerySort(query *queryBuilder, findFilter *models.FindFilterType) {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return
	}

	sort := findFilter.GetSort("path")
	direction := findFilter.GetDirection()

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
		query.sortAndPagination += fmt.Sprintf(" ORDER BY folders.path %s, file_folder.path %[1]s, files.basename %[1]s", direction)
	case "file_mod_time":
		sort = "mod_time"
		addFileTable()
		query.sortAndPagination += getSort(sort, direction, fileTable)
	case "title":
		addFileTable()
		addFolderTable()
		query.sortAndPagination += " ORDER BY COALESCE(galleries.title, files.basename, basename(COALESCE(folders.path, ''))) COLLATE NATURAL_CS " + direction + ", file_folder.path " + direction
	default:
		query.sortAndPagination += getSort(sort, direction, "galleries")
	}
}

func (qb *GalleryStore) filesRepository() *filesRepository {
	return &filesRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesFilesTable,
			idColumn:  galleryIDColumn,
		},
	}
}

func (qb *GalleryStore) AddFileID(ctx context.Context, id int, fileID file.ID) error {
	const firstPrimary = false
	return galleriesFilesTableMgr.insertJoins(ctx, id, firstPrimary, []file.ID{fileID})
}

func (qb *GalleryStore) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersGalleriesTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "performer_id",
	}
}

func (qb *GalleryStore) GetPerformerIDs(ctx context.Context, id int) ([]int, error) {
	return qb.performersRepository().getIDs(ctx, id)
}

func (qb *GalleryStore) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesTagsTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn:     "tag_id",
		foreignTable: tagTable,
		orderBy:      "tags.name ASC",
	}
}

func (qb *GalleryStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return qb.tagsRepository().getIDs(ctx, id)
}

func (qb *GalleryStore) imagesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesImagesTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "image_id",
	}
}

func (qb *GalleryStore) GetImageIDs(ctx context.Context, galleryID int) ([]int, error) {
	return qb.imagesRepository().getIDs(ctx, galleryID)
}

func (qb *GalleryStore) AddImages(ctx context.Context, galleryID int, imageIDs ...int) error {
	return qb.imagesRepository().insertOrIgnore(ctx, galleryID, imageIDs...)
}

func (qb *GalleryStore) RemoveImages(ctx context.Context, galleryID int, imageIDs ...int) error {
	return qb.imagesRepository().destroyJoins(ctx, galleryID, imageIDs...)
}

func (qb *GalleryStore) UpdateImages(ctx context.Context, galleryID int, imageIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.imagesRepository().replace(ctx, galleryID, imageIDs)
}

func (qb *GalleryStore) scenesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesScenesTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: sceneIDColumn,
	}
}

func (qb *GalleryStore) GetSceneIDs(ctx context.Context, id int) ([]int, error) {
	return qb.scenesRepository().getIDs(ctx, id)
}
