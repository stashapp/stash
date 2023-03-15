package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var imageTable = "images"

const (
	imageIDColumn         = "image_id"
	performersImagesTable = "performers_images"
	imagesTagsTable       = "images_tags"
	imagesFilesTable      = "images_files"
)

type imageRow struct {
	ID    int         `db:"id" goqu:"skipinsert"`
	Title zero.String `db:"title"`
	// expressed as 1-100
	Rating    null.Int               `db:"rating"`
	URL       zero.String            `db:"url"`
	Date      models.SQLiteDate      `db:"date"`
	Organized bool                   `db:"organized"`
	OCounter  int                    `db:"o_counter"`
	StudioID  null.Int               `db:"studio_id,omitempty"`
	CreatedAt models.SQLiteTimestamp `db:"created_at"`
	UpdatedAt models.SQLiteTimestamp `db:"updated_at"`
}

func (r *imageRow) fromImage(i models.Image) {
	r.ID = i.ID
	r.Title = zero.StringFrom(i.Title)
	r.Rating = intFromPtr(i.Rating)
	r.URL = zero.StringFrom(i.URL)
	if i.Date != nil {
		_ = r.Date.Scan(i.Date.Time)
	}
	r.Organized = i.Organized
	r.OCounter = i.OCounter
	r.StudioID = intFromPtr(i.StudioID)
	r.CreatedAt = models.SQLiteTimestamp{Timestamp: i.CreatedAt}
	r.UpdatedAt = models.SQLiteTimestamp{Timestamp: i.UpdatedAt}
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
		ID:        r.ID,
		Title:     r.Title.String,
		Rating:    nullIntPtr(r.Rating),
		URL:       r.URL.String,
		Date:      r.Date.DatePtr(),
		Organized: r.Organized,
		OCounter:  r.OCounter,
		StudioID:  nullIntPtr(r.StudioID),

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
	r.setNullInt("rating", i.Rating)
	r.setNullString("url", i.URL)
	r.setSQLiteDate("date", i.Date)
	r.setBool("organized", i.Organized)
	r.setInt("o_counter", i.OCounter)
	r.setNullInt("studio_id", i.StudioID)
	r.setSQLiteTimestamp("created_at", i.CreatedAt)
	r.setSQLiteTimestamp("updated_at", i.UpdatedAt)
}

type ImageStore struct {
	repository

	tableMgr *table
	oCounterManager

	fileStore *FileStore
}

func NewImageStore(fileStore *FileStore) *ImageStore {
	return &ImageStore{
		repository: repository{
			tableName: imageTable,
			idColumn:  idColumn,
		},
		tableMgr:        imageTableMgr,
		oCounterManager: oCounterManager{imageTableMgr},
		fileStore:       fileStore,
	}
}

func (qb *ImageStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *ImageStore) Create(ctx context.Context, newObject *models.ImageCreateInput) error {
	var r imageRow
	r.fromImage(*newObject.Image)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if len(newObject.FileIDs) > 0 {
		const firstPrimary = true
		if err := imagesFilesTableMgr.insertJoins(ctx, id, firstPrimary, newObject.FileIDs); err != nil {
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

	updated, err := qb.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject.Image = *updated

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
		fileIDs := make([]file.ID, len(updatedObject.Files.List()))
		for i, f := range updatedObject.Files.List() {
			fileIDs[i] = f.ID
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

func (qb *ImageStore) Find(ctx context.Context, id int) (*models.Image, error) {
	return qb.find(ctx, id)
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
			i := intslice.IntIndex(ids, s.ID)
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
			checksum.Col("type").Eq(file.FingerprintTypeMD5),
		),
	).Select(
		qb.table().All(),
		imagesFilesJoinTable.Col(fileIDColumn).As("primary_file_id"),
		folders.Col("path").As("primary_file_folder_path"),
		files.Col("basename").As("primary_file_basename"),
		checksum.Col("fingerprint").As("primary_file_checksum"),
	)
}

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

func (qb *ImageStore) GetFiles(ctx context.Context, id int) ([]*file.ImageFile, error) {
	fileIDs, err := qb.filesRepository().get(ctx, id)
	if err != nil {
		return nil, err
	}

	// use fileStore to load files
	files, err := qb.fileStore.Find(ctx, fileIDs...)
	if err != nil {
		return nil, err
	}

	ret := make([]*file.ImageFile, len(files))
	for i, f := range files {
		var ok bool
		ret[i], ok = f.(*file.ImageFile)
		if !ok {
			return nil, fmt.Errorf("expected file to be *file.ImageFile not %T", f)
		}
	}

	return ret, nil
}

func (qb *ImageStore) GetManyFileIDs(ctx context.Context, ids []int) ([][]file.ID, error) {
	const primaryOnly = false
	return qb.filesRepository().getMany(ctx, ids, primaryOnly)
}

func (qb *ImageStore) find(ctx context.Context, id int) (*models.Image, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting image by id %d: %w", id, err)
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

func (qb *ImageStore) FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Image, error) {
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

func (qb *ImageStore) CountByFileID(ctx context.Context, fileID file.ID) (int, error) {
	joinTable := imagesFilesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(fileIDColumn).Eq(fileID))
	return count(ctx, q)
}

func (qb *ImageStore) FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Image, error) {
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
	return qb.FindByFingerprints(ctx, []file.Fingerprint{
		{
			Type:        file.FingerprintTypeMD5,
			Fingerprint: checksum,
		},
	})
}

func (qb *ImageStore) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Image, error) {
	table := qb.table()
	fileTable := fileTableMgr.table
	folderTable := folderTableMgr.table

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
	).Order(folderTable.Col("path").Asc(), fileTable.Col("basename").Asc())

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting images for gallery %d: %w", galleryID, err)
	}

	return ret, nil
}

func (qb *ImageStore) CountByGalleryID(ctx context.Context, galleryID int) (int, error) {
	joinTable := goqu.T(galleriesImagesTable)

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col("gallery_id").Eq(galleryID))
	return count(ctx, q)
}

func (qb *ImageStore) FindByFolderID(ctx context.Context, folderID file.FolderID) ([]*models.Image, error) {
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

func (qb *ImageStore) FindByZipFileID(ctx context.Context, zipFileID file.ID) ([]*models.Image, error) {
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
		goqu.SUM(fileTableMgr.table.Col("size")),
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

func (qb *ImageStore) validateFilter(imageFilter *models.ImageFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if imageFilter.And != nil {
		if imageFilter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if imageFilter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(imageFilter.And)
	}

	if imageFilter.Or != nil {
		if imageFilter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(imageFilter.Or)
	}

	if imageFilter.Not != nil {
		return qb.validateFilter(imageFilter.Not)
	}

	return nil
}

func (qb *ImageStore) makeFilter(ctx context.Context, imageFilter *models.ImageFilterType) *filterBuilder {
	query := &filterBuilder{}

	if imageFilter.And != nil {
		query.and(qb.makeFilter(ctx, imageFilter.And))
	}
	if imageFilter.Or != nil {
		query.or(qb.makeFilter(ctx, imageFilter.Or))
	}
	if imageFilter.Not != nil {
		query.not(qb.makeFilter(ctx, imageFilter.Not))
	}

	query.handleCriterion(ctx, intCriterionHandler(imageFilter.ID, "images.id", nil))
	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if imageFilter.Checksum != nil {
			qb.addImagesFilesTable(f)
			f.addInnerJoin(fingerprintTable, "fingerprints_md5", "images_files.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
		}

		stringCriterionHandler(imageFilter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
	}))
	query.handleCriterion(ctx, stringCriterionHandler(imageFilter.Title, "images.title"))

	query.handleCriterion(ctx, pathCriterionHandler(imageFilter.Path, "folders.path", "files.basename", qb.addFoldersTable))
	query.handleCriterion(ctx, imageFileCountCriterionHandler(qb, imageFilter.FileCount))
	query.handleCriterion(ctx, intCriterionHandler(imageFilter.Rating100, "images.rating", nil))
	// legacy rating handler
	query.handleCriterion(ctx, rating5CriterionHandler(imageFilter.Rating, "images.rating", nil))
	query.handleCriterion(ctx, intCriterionHandler(imageFilter.OCounter, "images.o_counter", nil))
	query.handleCriterion(ctx, boolCriterionHandler(imageFilter.Organized, "images.organized", nil))
	query.handleCriterion(ctx, dateCriterionHandler(imageFilter.Date, "images.date"))
	query.handleCriterion(ctx, stringCriterionHandler(imageFilter.URL, "images.url"))

	query.handleCriterion(ctx, resolutionCriterionHandler(imageFilter.Resolution, "image_files.height", "image_files.width", qb.addImageFilesTable))
	query.handleCriterion(ctx, imageIsMissingCriterionHandler(qb, imageFilter.IsMissing))

	query.handleCriterion(ctx, imageTagsCriterionHandler(qb, imageFilter.Tags))
	query.handleCriterion(ctx, imageTagCountCriterionHandler(qb, imageFilter.TagCount))
	query.handleCriterion(ctx, imageGalleriesCriterionHandler(qb, imageFilter.Galleries))
	query.handleCriterion(ctx, imagePerformersCriterionHandler(qb, imageFilter.Performers))
	query.handleCriterion(ctx, imagePerformerCountCriterionHandler(qb, imageFilter.PerformerCount))
	query.handleCriterion(ctx, imageStudioCriterionHandler(qb, imageFilter.Studios))
	query.handleCriterion(ctx, imagePerformerTagsCriterionHandler(qb, imageFilter.PerformerTags))
	query.handleCriterion(ctx, imagePerformerFavoriteCriterionHandler(imageFilter.PerformerFavorite))
	query.handleCriterion(ctx, timestampCriterionHandler(imageFilter.CreatedAt, "images.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(imageFilter.UpdatedAt, "images.updated_at"))

	return query
}

func (qb *ImageStore) addImagesFilesTable(f *filterBuilder) {
	f.addLeftJoin(imagesFilesTable, "", "images_files.image_id = images.id")
}

func (qb *ImageStore) addFilesTable(f *filterBuilder) {
	qb.addImagesFilesTable(f)
	f.addLeftJoin(fileTable, "", "images_files.file_id = files.id")
}

func (qb *ImageStore) addFoldersTable(f *filterBuilder) {
	qb.addFilesTable(f)
	f.addLeftJoin(folderTable, "", "files.parent_folder_id = folders.id")
}

func (qb *ImageStore) addImageFilesTable(f *filterBuilder) {
	qb.addImagesFilesTable(f)
	f.addLeftJoin(imageFileTable, "", "image_files.file_id = images_files.file_id")
}

func (qb *ImageStore) makeQuery(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if imageFilter == nil {
		imageFilter = &models.ImageFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
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

	if err := qb.validateFilter(imageFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, imageFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	qb.setImageSortAndPagination(&query, findFilter)

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

	aggregateQuery := qb.newQuery()

	if options.Count {
		aggregateQuery.addColumn("COUNT(DISTINCT temp.id) as total")
	}

	// TODO - this doesn't work yet
	// if options.Megapixels {
	// 	query.addColumn("COALESCE(images.width, 0) * COALESCE(images.height, 0) / 1000000 as megapixels")
	// 	aggregateQuery.addColumn("COALESCE(SUM(temp.megapixels), 0) as megapixels")
	// }

	// if options.TotalSize {
	// 	query.addColumn("COALESCE(images.size, 0) as size")
	// 	aggregateQuery.addColumn("COALESCE(SUM(temp.size), 0) as size")
	// }

	const includeSortPagination = false
	aggregateQuery.from = fmt.Sprintf("(%s) as temp", query.toSQL(includeSortPagination))

	out := struct {
		Total      int
		Megapixels float64
		Size       float64
	}{}
	if err := qb.repository.queryStruct(ctx, aggregateQuery.toSQL(includeSortPagination), query.args, &out); err != nil {
		return nil, err
	}

	ret := models.NewImageQueryResult(qb)
	ret.Count = out.Total
	ret.Megapixels = out.Megapixels
	ret.TotalSize = out.Size
	return ret, nil
}

func (qb *ImageStore) QueryCount(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, imageFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

func imageFileCountCriterionHandler(qb *ImageStore, fileCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    imagesFilesTable,
		primaryFK:    imageIDColumn,
	}

	return h.handler(fileCount)
}

func imageIsMissingCriterionHandler(qb *ImageStore, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "studio":
				f.addWhere("images.studio_id IS NULL")
			case "performers":
				qb.performersRepository().join(f, "performers_join", "images.id")
				f.addWhere("performers_join.image_id IS NULL")
			case "galleries":
				qb.galleriesRepository().join(f, "galleries_join", "images.id")
				f.addWhere("galleries_join.image_id IS NULL")
			case "tags":
				qb.tagsRepository().join(f, "tags_join", "images.id")
				f.addWhere("tags_join.image_id IS NULL")
			default:
				f.addWhere("(images." + *isMissing + " IS NULL OR TRIM(images." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *ImageStore) getMultiCriterionHandlerBuilder(foreignTable, joinTable, foreignFK string, addJoinsFunc func(f *filterBuilder)) multiCriterionHandlerBuilder {
	return multiCriterionHandlerBuilder{
		primaryTable: imageTable,
		foreignTable: foreignTable,
		joinTable:    joinTable,
		primaryFK:    imageIDColumn,
		foreignFK:    foreignFK,
		addJoinsFunc: addJoinsFunc,
	}
}

func imageTagsCriterionHandler(qb *ImageStore, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: imageTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "image_tag",
		joinTable:      imagesTagsTable,
		primaryFK:      imageIDColumn,
	}

	return h.handler(tags)
}

func imageTagCountCriterionHandler(qb *ImageStore, tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    imagesTagsTable,
		primaryFK:    imageIDColumn,
	}

	return h.handler(tagCount)
}

func imageGalleriesCriterionHandler(qb *ImageStore, galleries *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		if galleries.Modifier == models.CriterionModifierIncludes || galleries.Modifier == models.CriterionModifierIncludesAll {
			f.addInnerJoin(galleriesImagesTable, "", "galleries_images.image_id = images.id")
			f.addInnerJoin(galleryTable, "", "galleries_images.gallery_id = galleries.id")
		}
	}
	h := qb.getMultiCriterionHandlerBuilder(galleryTable, galleriesImagesTable, galleryIDColumn, addJoinsFunc)

	return h.handler(galleries)
}

func imagePerformersCriterionHandler(qb *ImageStore, performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    performersImagesTable,
		joinAs:       "performers_join",
		primaryFK:    imageIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			qb.performersRepository().join(f, "performers_join", "images.id")
		},
	}

	return h.handler(performers)
}

func imagePerformerCountCriterionHandler(qb *ImageStore, performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    performersImagesTable,
		primaryFK:    imageIDColumn,
	}

	return h.handler(performerCount)
}

func imagePerformerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_images", "", "images.id = performers_images.image_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_images.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_images.image_id as id FROM performers_images
JOIN performers ON performers.id = performers_images.performer_id
GROUP BY performers_images.image_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "images.id = nofaves.id")
				f.addWhere("performers_images.image_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func imageStudioCriterionHandler(qb *ImageStore, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := hierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: imageTable,
		foreignTable: studioTable,
		foreignFK:    studioIDColumn,
		parentFK:     "parent_id",
	}

	return h.handler(studios)
}

func imagePerformerTagsCriterionHandler(qb *ImageStore, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("performers_images", "", "images.id = performers_images.image_id")
				f.addLeftJoin("performers_tags", "", "performers_images.performer_id = performers_tags.performer_id")

				f.addWhere(fmt.Sprintf("performers_tags.tag_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 {
				return
			}

			valuesClause := getHierarchicalValues(ctx, qb.tx, tags.Value, tagTable, "tags_relations", "", tags.Depth)

			f.addWith(`performer_tags AS (
SELECT pi.image_id, t.column1 AS root_tag_id FROM performers_images pi
INNER JOIN performers_tags pt ON pt.performer_id = pi.performer_id
INNER JOIN (` + valuesClause + `) t ON t.column2 = pt.tag_id
)`)

			f.addLeftJoin("performer_tags", "", "performer_tags.image_id = images.id")

			addHierarchicalConditionClauses(f, tags, "performer_tags", "root_tag_id")
		}
	}
}

func (qb *ImageStore) setImageSortAndPagination(q *queryBuilder, findFilter *models.FindFilterType) {
	sortClause := ""

	if findFilter != nil && findFilter.Sort != nil && *findFilter.Sort != "" {
		sort := findFilter.GetSort("title")
		direction := findFilter.GetDirection()

		// translate sort field
		if sort == "file_mod_time" {
			sort = "mod_time"
		}

		addFilesJoin := func() {
			q.addJoins(
				join{
					table:    imagesFilesTable,
					onClause: "images_files.image_id = images.id",
				},
				join{
					table:    fileTable,
					onClause: "images_files.file_id = files.id",
				},
			)
		}

		addFolderJoin := func() {
			q.addJoins(join{
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			})
		}

		switch sort {
		case "path":
			addFilesJoin()
			addFolderJoin()
			sortClause = " ORDER BY folders.path " + direction + ", files.basename " + direction
		case "file_count":
			sortClause = getCountSort(imageTable, imagesFilesTable, imageIDColumn, direction)
		case "tag_count":
			sortClause = getCountSort(imageTable, imagesTagsTable, imageIDColumn, direction)
		case "performer_count":
			sortClause = getCountSort(imageTable, performersImagesTable, imageIDColumn, direction)
		case "mod_time", "filesize":
			addFilesJoin()
			sortClause = getSort(sort, direction, "files")
		case "title":
			addFilesJoin()
			addFolderJoin()
			sortClause = " ORDER BY COALESCE(images.title, files.basename) COLLATE NATURAL_CS " + direction + ", folders.path " + direction
		default:
			sortClause = getSort(sort, direction, "images")
		}
	}

	q.sortAndPagination = sortClause + getPagination(findFilter)
}

func (qb *ImageStore) galleriesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesImagesTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: galleryIDColumn,
	}
}

func (qb *ImageStore) filesRepository() *filesRepository {
	return &filesRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: imagesFilesTable,
			idColumn:  imageIDColumn,
		},
	}
}

func (qb *ImageStore) AddFileID(ctx context.Context, id int, fileID file.ID) error {
	const firstPrimary = false
	return imagesFilesTableMgr.insertJoins(ctx, id, firstPrimary, []file.ID{fileID})
}

func (qb *ImageStore) GetGalleryIDs(ctx context.Context, imageID int) ([]int, error) {
	return qb.galleriesRepository().getIDs(ctx, imageID)
}

// func (qb *imageQueryBuilder) UpdateGalleries(ctx context.Context, imageID int, galleryIDs []int) error {
// 	// Delete the existing joins and then create new ones
// 	return qb.galleriesRepository().replace(ctx, imageID, galleryIDs)
// }

func (qb *ImageStore) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersImagesTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: performerIDColumn,
	}
}

func (qb *ImageStore) GetPerformerIDs(ctx context.Context, imageID int) ([]int, error) {
	return qb.performersRepository().getIDs(ctx, imageID)
}

func (qb *ImageStore) UpdatePerformers(ctx context.Context, imageID int, performerIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.performersRepository().replace(ctx, imageID, performerIDs)
}

func (qb *ImageStore) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: imagesTagsTable,
			idColumn:  imageIDColumn,
		},
		fkColumn:     tagIDColumn,
		foreignTable: tagTable,
		orderBy:      "tags.name ASC",
	}
}

func (qb *ImageStore) GetTagIDs(ctx context.Context, imageID int) ([]int, error) {
	return qb.tagsRepository().getIDs(ctx, imageID)
}

func (qb *ImageStore) UpdateTags(ctx context.Context, imageID int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.tagsRepository().replace(ctx, imageID, tagIDs)
}
