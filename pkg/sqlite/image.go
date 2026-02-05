package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

const imageTable = "images"

const (
	imageIDColumn         = "image_id"
	performersImagesTable = "performers_images"
	imagesTagsTable       = "images_tags"
	imagesFilesTable      = "images_files"
	imagesURLsTable       = "image_urls"
	imageURLColumn        = "url"
)

type imageRow struct {
	ID    int         `db:"id" goqu:"skipinsert"`
	Title zero.String `db:"title"`
	Code  zero.String `db:"code"`
	// expressed as 1-100
	Rating        null.Int    `db:"rating"`
	Date          NullDate    `db:"date"`
	DatePrecision null.Int    `db:"date_precision"`
	Details       zero.String `db:"details"`
	Photographer  zero.String `db:"photographer"`
	Organized     bool        `db:"organized"`
	OCounter      int         `db:"o_counter"`
	StudioID      null.Int    `db:"studio_id,omitempty"`
	CreatedAt     Timestamp   `db:"created_at"`
	UpdatedAt     Timestamp   `db:"updated_at"`
}

func (r *imageRow) fromImage(i models.Image) {
	r.ID = i.ID
	r.Title = zero.StringFrom(i.Title)
	r.Code = zero.StringFrom(i.Code)
	r.Rating = intFromPtr(i.Rating)
	r.Date = NullDateFromDatePtr(i.Date)
	r.DatePrecision = datePrecisionFromDatePtr(i.Date)
	r.Details = zero.StringFrom(i.Details)
	r.Photographer = zero.StringFrom(i.Photographer)
	r.Organized = i.Organized
	r.OCounter = i.OCounter
	r.StudioID = intFromPtr(i.StudioID)
	r.CreatedAt = Timestamp{Timestamp: i.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: i.UpdatedAt}
}

type imageQueryRow struct {
	imageRow
	PrimaryFileID         null.Int    `db:"primary_file_id"`
	PrimaryFileFolderPath zero.String `db:"primary_file_folder_path"`
	PrimaryFileBasename   zero.String `db:"primary_file_basename"`
	PrimaryFileChecksum   zero.String `db:"primary_file_checksum"`
}

func (r *imageQueryRow) resolve() *models.Image {
	ret := &models.Image{
		ID:           r.ID,
		Title:        r.Title.String,
		Code:         r.Code.String,
		Rating:       nullIntPtr(r.Rating),
		Date:         r.Date.DatePtr(r.DatePrecision),
		Details:      r.Details.String,
		Photographer: r.Photographer.String,
		Organized:    r.Organized,
		OCounter:     r.OCounter,
		StudioID:     nullIntPtr(r.StudioID),

		PrimaryFileID: nullIntFileIDPtr(r.PrimaryFileID),
		Checksum:      r.PrimaryFileChecksum.String,

		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,
	}

	if r.PrimaryFileFolderPath.Valid && r.PrimaryFileBasename.Valid {
		ret.Path = filepath.Join(r.PrimaryFileFolderPath.String, r.PrimaryFileBasename.String)
	}

	return ret
}

type imageRowRecord struct {
	updateRecord
}

func (r *imageRowRecord) fromPartial(i models.ImagePartial) {
	r.setNullString("title", i.Title)
	r.setNullString("code", i.Code)
	r.setNullInt("rating", i.Rating)
	r.setNullDate("date", "date_precision", i.Date)
	r.setNullString("details", i.Details)
	r.setNullString("photographer", i.Photographer)
	r.setBool("organized", i.Organized)
	r.setInt("o_counter", i.OCounter)
	r.setNullInt("studio_id", i.StudioID)
	r.setTimestamp("created_at", i.CreatedAt)
	r.setTimestamp("updated_at", i.UpdatedAt)
}

type imageRepositoryType struct {
	repository
	performers joinRepository
	galleries  joinRepository
	tags       joinRepository
	files      filesRepository
}

func (r *imageRepositoryType) addImagesFilesTable(f *filterBuilder) {
	f.addLeftJoin(imagesFilesTable, "", "images_files.image_id = images.id")
}

func (r *imageRepositoryType) addFilesTable(f *filterBuilder) {
	r.addImagesFilesTable(f)
	f.addLeftJoin(fileTable, "", "images_files.file_id = files.id")
}

func (r *imageRepositoryType) addFoldersTable(f *filterBuilder) {
	r.addFilesTable(f)
	f.addLeftJoin(folderTable, "", "files.parent_folder_id = folders.id")
}

func (r *imageRepositoryType) addImageFilesTable(f *filterBuilder) {
	r.addImagesFilesTable(f)
	f.addLeftJoin(imageFileTable, "", "image_files.file_id = images_files.file_id")
}

var (
	imageRepository = imageRepositoryType{
		repository: repository{
			tableName: imageTable,
			idColumn:  idColumn,
		},

		performers: joinRepository{
			repository: repository{
				tableName: performersImagesTable,
				idColumn:  imageIDColumn,
			},
			fkColumn: performerIDColumn,
		},

		galleries: joinRepository{
			repository: repository{
				tableName: galleriesImagesTable,
				idColumn:  imageIDColumn,
			},
			fkColumn: galleryIDColumn,
		},

		files: filesRepository{
			repository: repository{
				tableName: imagesFilesTable,
				idColumn:  imageIDColumn,
			},
		},

		tags: joinRepository{
			repository: repository{
				tableName: imagesTagsTable,
				idColumn:  imageIDColumn,
			},
			fkColumn:     tagIDColumn,
			foreignTable: tagTable,
			orderBy:      tagTableSortSQL,
		},
	}
)

type ImageStore struct {
	tableMgr *table
	oCounterManager

	repo *storeRepository
}

func NewImageStore(r *storeRepository) *ImageStore {
	return &ImageStore{
		tableMgr:        imageTableMgr,
		oCounterManager: oCounterManager{imageTableMgr},
		repo:            r,
	}
}

func (qb *ImageStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *ImageStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	files := fileTableMgr.table
	folders := folderTableMgr.table
	checksum := fingerprintTableMgr.table

	return dialect.From(table).LeftJoin(
		imagesFilesJoinTable,
		goqu.On(
			imagesFilesJoinTable.Col(imageIDColumn).Eq(table.Col(idColumn)),
			imagesFilesJoinTable.Col("primary").Eq(1),
		),
	).LeftJoin(
		files,
		goqu.On(files.Col(idColumn).Eq(imagesFilesJoinTable.Col(fileIDColumn))),
	).LeftJoin(
		folders,
		goqu.On(folders.Col(idColumn).Eq(files.Col("parent_folder_id"))),
	).LeftJoin(
		checksum,
		goqu.On(
			checksum.Col(fileIDColumn).Eq(imagesFilesJoinTable.Col(fileIDColumn)),
			checksum.Col("type").Eq(models.FingerprintTypeMD5),
		),
	).Select(
		qb.table().All(),
		imagesFilesJoinTable.Col(fileIDColumn).As("primary_file_id"),
		folders.Col("path").As("primary_file_folder_path"),
		files.Col("basename").As("primary_file_basename"),
		checksum.Col("fingerprint").As("primary_file_checksum"),
	)
}

func (qb *ImageStore) Create(ctx context.Context, newObject *models.Image, fileIDs []models.FileID) error {
	var r imageRow
	r.fromImage(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if len(fileIDs) > 0 {
		const firstPrimary = true
		if err := imagesFilesTableMgr.insertJoins(ctx, id, firstPrimary, fileIDs); err != nil {
			return err
		}
	}

	if newObject.URLs.Loaded() {
		const startPos = 0
		if err := imagesURLsTableMgr.insertJoins(ctx, id, startPos, newObject.URLs.List()); err != nil {
			return err
		}
	}

	if newObject.PerformerIDs.Loaded() {
		if err := imagesPerformersTableMgr.insertJoins(ctx, id, newObject.PerformerIDs.List()); err != nil {
			return err
		}
	}
	if newObject.TagIDs.Loaded() {
		if err := imagesTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if newObject.GalleryIDs.Loaded() {
		if err := imageGalleriesTableMgr.insertJoins(ctx, id, newObject.GalleryIDs.List()); err != nil {
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

func (qb *ImageStore) UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error) {
	r := imageRowRecord{
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

	if partial.GalleryIDs != nil {
		if err := imageGalleriesTableMgr.modifyJoins(ctx, id, partial.GalleryIDs.IDs, partial.GalleryIDs.Mode); err != nil {
			return nil, err
		}
	}

	if partial.URLs != nil {
		if err := imagesURLsTableMgr.modifyJoins(ctx, id, partial.URLs.Values, partial.URLs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.PerformerIDs != nil {
		if err := imagesPerformersTableMgr.modifyJoins(ctx, id, partial.PerformerIDs.IDs, partial.PerformerIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.TagIDs != nil {
		if err := imagesTagsTableMgr.modifyJoins(ctx, id, partial.TagIDs.IDs, partial.TagIDs.Mode); err != nil {
			return nil, err
		}
	}

	if partial.PrimaryFileID != nil {
		if err := imagesFilesTableMgr.setPrimary(ctx, id, *partial.PrimaryFileID); err != nil {
			return nil, err
		}
	}

	return qb.find(ctx, id)
}

func (qb *ImageStore) Update(ctx context.Context, updatedObject *models.Image) error {
	var r imageRow
	r.fromImage(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.URLs.Loaded() {
		if err := imagesURLsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.URLs.List()); err != nil {
			return err
		}
	}

	if updatedObject.PerformerIDs.Loaded() {
		if err := imagesPerformersTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.PerformerIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.TagIDs.Loaded() {
		if err := imagesTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.GalleryIDs.Loaded() {
		if err := imageGalleriesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.GalleryIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.Files.Loaded() {
		fileIDs := make([]models.FileID, len(updatedObject.Files.List()))
		for i, f := range updatedObject.Files.List() {
			fileIDs[i] = f.Base().ID
		}

		if err := imagesFilesTableMgr.replaceJoins(ctx, updatedObject.ID, fileIDs); err != nil {
			return err
		}
	}
	return nil
}

func (qb *ImageStore) Destroy(ctx context.Context, id int) error {
	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *ImageStore) Find(ctx context.Context, id int) (*models.Image, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *ImageStore) FindMany(ctx context.Context, ids []int) ([]*models.Image, error) {
	images := make([]*models.Image, len(ids))

	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := qb.selectDataset().Prepared(true).Where(qb.table().Col(idColumn).In(batch))
		unsorted, err := qb.getMany(ctx, q)
		if err != nil {
			return err
		}

		for _, s := range unsorted {
			i := slices.Index(ids, s.ID)
			images[i] = s
		}

		return nil
	}); err != nil {
		return nil, err
	}

	for i := range images {
		if images[i] == nil {
			return nil, fmt.Errorf("image with id %d not found", ids[i])
		}
	}

	return images, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *ImageStore) find(ctx context.Context, id int) (*models.Image, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *ImageStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Image, error) {
	table := qb.table()

	q := qb.selectDataset().Prepared(true).Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

// returns nil, sql.ErrNoRows if not found
func (qb *ImageStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Image, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *ImageStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Image, error) {
	const single = false
	var ret []*models.Image
	var lastID int
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f imageQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		i := f.resolve()

		if i.ID == lastID {
			return fmt.Errorf("internal error: multiple rows returned for single image id %d", i.ID)
		}
		lastID = i.ID

		ret = append(ret, i)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

// Returns the custom cover for the gallery, if one has been set.
func (qb *ImageStore) CoverByGalleryID(ctx context.Context, galleryID int) (*models.Image, error) {
	table := qb.table()

	sq := dialect.From(table).
		InnerJoin(
			galleriesImagesJoinTable,
			goqu.On(table.Col(idColumn).Eq(galleriesImagesJoinTable.Col(imageIDColumn))),
		).
		Select(table.Col(idColumn)).
		Where(goqu.And(
			galleriesImagesJoinTable.Col("gallery_id").Eq(galleryID),
			galleriesImagesJoinTable.Col("cover").Eq(true),
		))

	q := qb.selectDataset().Prepared(true).Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting cover for gallery %d: %w", galleryID, err)
	}

	switch {
	case len(ret) > 1:
		return nil, fmt.Errorf("internal error: multiple covers returned for gallery %d", galleryID)
	case len(ret) == 1:
		return ret[0], nil
	default:
		return nil, nil
	}
}

func (qb *ImageStore) GetFiles(ctx context.Context, id int) ([]models.File, error) {
	fileIDs, err := imageRepository.files.get(ctx, id)
	if err != nil {
		return nil, err
	}

	// use fileStore to load files
	files, err := qb.repo.File.Find(ctx, fileIDs...)
	if err != nil {
		return nil, err
	}

	ret := make([]models.File, len(files))
	copy(ret, files)

	return ret, nil
}

func (qb *ImageStore) GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error) {
	const primaryOnly = false
	return imageRepository.files.getMany(ctx, ids, primaryOnly)
}

func (qb *ImageStore) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Image, error) {
	table := qb.table()

	sq := dialect.From(table).
		InnerJoin(
			imagesFilesJoinTable,
			goqu.On(table.Col(idColumn).Eq(imagesFilesJoinTable.Col(imageIDColumn))),
		).
		Select(table.Col(idColumn)).Where(imagesFilesJoinTable.Col(fileIDColumn).Eq(fileID))

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting image by file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *ImageStore) CountByFileID(ctx context.Context, fileID models.FileID) (int, error) {
	joinTable := imagesFilesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(fileIDColumn).Eq(fileID))
	return count(ctx, q)
}

func (qb *ImageStore) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Image, error) {
	table := qb.table()
	fingerprintTable := fingerprintTableMgr.table

	var ex []exp.Expression

	for _, v := range fp {
		ex = append(ex, goqu.And(
			fingerprintTable.Col("type").Eq(v.Type),
			fingerprintTable.Col("fingerprint").Eq(v.Fingerprint),
		))
	}

	sq := dialect.From(table).
		InnerJoin(
			imagesFilesJoinTable,
			goqu.On(table.Col(idColumn).Eq(imagesFilesJoinTable.Col(imageIDColumn))),
		).
		InnerJoin(
			fingerprintTable,
			goqu.On(fingerprintTable.Col(fileIDColumn).Eq(imagesFilesJoinTable.Col(fileIDColumn))),
		).
		Select(table.Col(idColumn)).Where(goqu.Or(ex...))

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting image by fingerprints: %w", err)
	}

	return ret, nil
}

func (qb *ImageStore) FindByChecksum(ctx context.Context, checksum string) ([]*models.Image, error) {
	return qb.FindByFingerprints(ctx, []models.Fingerprint{
		{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: checksum,
		},
	})
}

var defaultGalleryOrder = []exp.OrderedExpression{
	goqu.L("COALESCE(folders.path, '') || COALESCE(files.basename, '') COLLATE NATURAL_CI").Asc(),
	goqu.L("COALESCE(images.title, images.id) COLLATE NATURAL_CI").Asc(),
}

func (qb *ImageStore) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Image, error) {
	table := qb.table()

	sq := dialect.From(table).
		InnerJoin(
			galleriesImagesJoinTable,
			goqu.On(table.Col(idColumn).Eq(galleriesImagesJoinTable.Col(imageIDColumn))),
		).
		Select(table.Col(idColumn)).Where(
		galleriesImagesJoinTable.Col("gallery_id").Eq(galleryID),
	)

	q := qb.selectDataset().Prepared(true).Where(
		table.Col(idColumn).Eq(
			sq,
		),
	).Order(defaultGalleryOrder...)

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting images for gallery %d: %w", galleryID, err)
	}

	return ret, nil
}

func (qb *ImageStore) FindByGalleryIDIndex(ctx context.Context, galleryID int, index uint) (*models.Image, error) {
	table := qb.table()

	q := qb.selectDataset().
		InnerJoin(
			galleriesImagesJoinTable,
			goqu.On(table.Col(idColumn).Eq(galleriesImagesJoinTable.Col(imageIDColumn))),
		).
		Where(galleriesImagesJoinTable.Col(galleryIDColumn).Eq(galleryID)).
		Prepared(true).
		Order(defaultGalleryOrder...).
		Limit(1).Offset(index)

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting images for gallery %d: %w", galleryID, err)
	}

	if len(ret) == 0 {
		return nil, nil
	}

	return ret[0], nil
}

func (qb *ImageStore) CountByGalleryID(ctx context.Context, galleryID int) (int, error) {
	joinTable := goqu.T(galleriesImagesTable)

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col("gallery_id").Eq(galleryID))
	return count(ctx, q)
}

func (qb *ImageStore) OCountByPerformerID(ctx context.Context, performerID int) (int, error) {
	table := qb.table()
	joinTable := performersImagesJoinTable
	q := dialect.Select(goqu.COALESCE(goqu.SUM("o_counter"), 0)).From(table).InnerJoin(joinTable, goqu.On(table.Col(idColumn).Eq(joinTable.Col(imageIDColumn)))).Where(joinTable.Col(performerIDColumn).Eq(performerID))

	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *ImageStore) OCountByStudioID(ctx context.Context, studioID int) (int, error) {
	table := qb.table()
	q := dialect.Select(goqu.COALESCE(goqu.SUM("o_counter"), 0)).From(table).Where(
		table.Col(studioIDColumn).Eq(studioID),
	)

	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *ImageStore) OCount(ctx context.Context) (int, error) {
	table := qb.table()

	q := dialect.Select(goqu.COALESCE(goqu.SUM("o_counter"), 0)).From(table)
	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *ImageStore) FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Image, error) {
	table := qb.table()
	fileTable := goqu.T(fileTable)

	sq := dialect.From(table).
		InnerJoin(
			imagesFilesJoinTable,
			goqu.On(table.Col(idColumn).Eq(imagesFilesJoinTable.Col(imageIDColumn))),
		).
		InnerJoin(
			fileTable,
			goqu.On(imagesFilesJoinTable.Col(fileIDColumn).Eq(fileTable.Col(idColumn))),
		).
		Select(table.Col(idColumn)).Where(
		fileTable.Col("parent_folder_id").Eq(folderID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting image by folder: %w", err)
	}

	return ret, nil
}

func (qb *ImageStore) FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Image, error) {
	table := qb.table()
	fileTable := goqu.T(fileTable)

	sq := dialect.From(table).
		InnerJoin(
			imagesFilesJoinTable,
			goqu.On(table.Col(idColumn).Eq(imagesFilesJoinTable.Col(imageIDColumn))),
		).
		InnerJoin(
			fileTable,
			goqu.On(imagesFilesJoinTable.Col(fileIDColumn).Eq(fileTable.Col(idColumn))),
		).
		Select(table.Col(idColumn)).Where(
		fileTable.Col("zip_file_id").Eq(zipFileID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting image by zip file: %w", err)
	}

	return ret, nil
}

func (qb *ImageStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *ImageStore) Size(ctx context.Context) (float64, error) {
	table := qb.table()
	fileTable := fileTableMgr.table
	q := dialect.Select(
		goqu.COALESCE(goqu.SUM(fileTableMgr.table.Col("size")), 0),
	).From(table).InnerJoin(
		imagesFilesJoinTable,
		goqu.On(table.Col(idColumn).Eq(imagesFilesJoinTable.Col(imageIDColumn))),
	).InnerJoin(
		fileTable,
		goqu.On(imagesFilesJoinTable.Col(fileIDColumn).Eq(fileTable.Col(idColumn))),
	)
	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *ImageStore) All(ctx context.Context) ([]*models.Image, error) {
	return qb.getMany(ctx, qb.selectDataset())
}

func (qb *ImageStore) makeQuery(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if imageFilter == nil {
		imageFilter = &models.ImageFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := imageRepository.newQuery()
	distinctIDs(&query, imageTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.addJoins(
			join{
				table:    imagesFilesTable,
				onClause: "images_files.image_id = images.id",
			},
			join{
				table:    fileTable,
				onClause: "images_files.file_id = files.id",
			},
			join{
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			},
			join{
				table:    fingerprintTable,
				onClause: "files_fingerprints.file_id = images_files.file_id",
			},
		)

		filepathColumn := "folders.path || '" + string(filepath.Separator) + "' || files.basename"
		searchColumns := []string{"images.title", filepathColumn, "files_fingerprints.fingerprint"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := filterBuilderFromHandler(ctx, &imageFilterHandler{
		imageFilter: imageFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setImageSortAndPagination(&query, findFilter); err != nil {
		return nil, err
	}

	return &query, nil
}

func (qb *ImageStore) Query(ctx context.Context, options models.ImageQueryOptions) (*models.ImageQueryResult, error) {
	query, err := qb.makeQuery(ctx, options.ImageFilter, options.FindFilter)
	if err != nil {
		return nil, err
	}

	result, err := qb.queryGroupedFields(ctx, options, *query)
	if err != nil {
		return nil, fmt.Errorf("error querying aggregate fields: %w", err)
	}

	idsResult, err := query.findIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding IDs: %w", err)
	}

	result.IDs = idsResult
	return result, nil
}

func (qb *ImageStore) queryGroupedFields(ctx context.Context, options models.ImageQueryOptions, query queryBuilder) (*models.ImageQueryResult, error) {
	if !options.Count && !options.Megapixels && !options.TotalSize {
		// nothing to do - return empty result
		return models.NewImageQueryResult(qb), nil
	}

	aggregateQuery := imageRepository.newQuery()

	if options.Count {
		aggregateQuery.addColumn("COUNT(DISTINCT temp.id) as total")
	}

	if options.Megapixels {
		query.addJoins(
			join{
				table:    imagesFilesTable,
				onClause: "images_files.image_id = images.id",
			},
			join{
				table:    imageFileTable,
				onClause: "images_files.file_id = image_files.file_id",
			},
		)
		query.addColumn("COALESCE(image_files.width, 0) * COALESCE(image_files.height, 0) as megapixels")
		aggregateQuery.addColumn("COALESCE(SUM(temp.megapixels), 0) / 1000000 as megapixels")
	}

	if options.TotalSize {
		query.addJoins(
			join{
				table:    imagesFilesTable,
				onClause: "images_files.image_id = images.id",
			},
			join{
				table:    fileTable,
				onClause: "images_files.file_id = files.id",
			},
		)
		query.addColumn("COALESCE(files.size, 0) as size")
		aggregateQuery.addColumn("SUM(temp.size) as size")
	}

	const includeSortPagination = false
	aggregateQuery.from = fmt.Sprintf("(%s) as temp", query.toSQL(includeSortPagination))

	out := struct {
		Total      int
		Megapixels null.Float
		Size       null.Float
	}{}
	if err := imageRepository.queryStruct(ctx, aggregateQuery.toSQL(includeSortPagination), query.args, &out); err != nil {
		return nil, err
	}

	ret := models.NewImageQueryResult(qb)
	ret.Count = out.Total
	ret.Megapixels = out.Megapixels.Float64
	ret.TotalSize = out.Size.Float64
	return ret, nil
}

func (qb *ImageStore) QueryCount(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, imageFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

var imageSortOptions = sortOptions{
	"created_at",
	"date",
	"file_count",
	"file_mod_time",
	"filesize",
	"id",
	"o_counter",
	"path",
	"performer_count",
	"random",
	"rating",
	"resolution",
	"tag_count",
	"title",
	"updated_at",
}

func (qb *ImageStore) setImageSortAndPagination(q *queryBuilder, findFilter *models.FindFilterType) error {
	sortClause := ""

	if findFilter != nil && findFilter.Sort != nil && *findFilter.Sort != "" {
		sort := findFilter.GetSort("title")
		direction := findFilter.GetDirection()

		// CVE-2024-32231 - ensure sort is in the list of allowed sorts
		if err := imageSortOptions.validateSort(sort); err != nil {
			return err
		}

		// translate sort field
		if sort == "file_mod_time" {
			sort = "mod_time"
		}

		addFilesJoin := func() {
			q.addJoins(
				join{
					sort:     true,
					table:    imagesFilesTable,
					onClause: "images_files.image_id = images.id",
				},
				join{
					sort:     true,
					table:    fileTable,
					onClause: "images_files.file_id = files.id",
				},
			)
		}

		addFolderJoin := func() {
			q.addJoins(join{
				sort:     true,
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			})
		}

		switch sort {
		case "path":
			addFilesJoin()
			addFolderJoin()
			sortClause = " ORDER BY COALESCE(folders.path, '') || COALESCE(files.basename, '') COLLATE NATURAL_CI " + direction
		case "file_count":
			sortClause = getCountSort(imageTable, imagesFilesTable, imageIDColumn, direction)
		case "tag_count":
			sortClause = getCountSort(imageTable, imagesTagsTable, imageIDColumn, direction)
		case "performer_count":
			sortClause = getCountSort(imageTable, performersImagesTable, imageIDColumn, direction)
		case "mod_time", "filesize":
			addFilesJoin()
			sortClause = getSort(sort, direction, "files")
		case "resolution":
			addFilesJoin()
			q.addJoins(join{
				sort:     true,
				table:    imageFileTable,
				onClause: "images_files.file_id = image_files.file_id",
			})
			sortClause = " ORDER BY MIN(image_files.width, image_files.height) " + direction
		case "title":
			addFilesJoin()
			addFolderJoin()
			sortClause = " ORDER BY COALESCE(images.title, files.basename) COLLATE NATURAL_CI " + direction + ", folders.path COLLATE NATURAL_CI " + direction
		default:
			sortClause = getSort(sort, direction, "images")
		}

		// Whatever the sorting, always use title/id as a final sort
		sortClause += ", COALESCE(images.title, images.id) COLLATE NATURAL_CI ASC"
	}

	q.sortAndPagination = sortClause + getPagination(findFilter)

	return nil
}

func (qb *ImageStore) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	const firstPrimary = false
	return imagesFilesTableMgr.insertJoins(ctx, id, firstPrimary, []models.FileID{fileID})
}

// RemoveFileID removes the file ID from the image.
// If the file ID is the primary file, then the next file in the list is set as the primary file.
func (qb *ImageStore) RemoveFileID(ctx context.Context, id int, fileID models.FileID) error {
	fileIDs, err := imagesFilesTableMgr.get(ctx, id)
	if err != nil {
		return fmt.Errorf("getting file IDs for image %d: %w", id, err)
	}

	fileIDs = sliceutil.Filter(fileIDs, func(f models.FileID) bool {
		return f != fileID
	})

	return imagesFilesTableMgr.replaceJoins(ctx, id, fileIDs)
}

func (qb *ImageStore) GetGalleryIDs(ctx context.Context, imageID int) ([]int, error) {
	return imageRepository.galleries.getIDs(ctx, imageID)
}

// func (qb *imageQueryBuilder) UpdateGalleries(ctx context.Context, imageID int, galleryIDs []int) error {
// 	// Delete the existing joins and then create new ones
// 	return qb.galleriesRepository().replace(ctx, imageID, galleryIDs)
// }

func (qb *ImageStore) GetPerformerIDs(ctx context.Context, imageID int) ([]int, error) {
	return imageRepository.performers.getIDs(ctx, imageID)
}

func (qb *ImageStore) UpdatePerformers(ctx context.Context, imageID int, performerIDs []int) error {
	// Delete the existing joins and then create new ones
	return imageRepository.performers.replace(ctx, imageID, performerIDs)
}

func (qb *ImageStore) GetTagIDs(ctx context.Context, imageID int) ([]int, error) {
	return imageRepository.tags.getIDs(ctx, imageID)
}

func (qb *ImageStore) UpdateTags(ctx context.Context, imageID int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return imageRepository.tags.replace(ctx, imageID, tagIDs)
}

func (qb *ImageStore) GetURLs(ctx context.Context, imageID int) ([]string, error) {
	return imagesURLsTableMgr.get(ctx, imageID)
}
