package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var imageTable = "images"
var imageTableMgr = &table{table: goqu.T(imageTable)}

const imageIDColumn = "image_id"
const performersImagesTable = "performers_images"
const imagesTagsTable = "images_tags"

type imageRow struct {
	ID          int        `db:"id" goqu:"skipinsert"`
	Checksum    string     `db:"checksum"`
	Path        string     `db:"path"`
	Title       nullString `db:"title"`
	Rating      nullInt    `db:"rating"`
	Organized   bool       `db:"organized"`
	OCounter    int        `db:"o_counter"`
	Size        nullInt64  `db:"size"`
	Width       nullInt    `db:"width"`
	Height      nullInt    `db:"height"`
	StudioID    nullInt    `db:"studio_id,omitempty"`
	FileModTime nullTime   `db:"file_mod_time"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (r *imageRow) fromImage(i models.Image) {
	r.ID = i.ID
	r.Checksum = i.Checksum
	r.Path = i.Path
	r.Title = newNullStringPtr(i.Title)
	r.Rating = newNullIntPtr(i.Rating)
	r.Organized = i.Organized
	r.OCounter = i.OCounter
	r.Size = newNullInt64Ptr(i.Size)
	r.Width = newNullIntPtr(i.Width)
	r.Height = newNullIntPtr(i.Height)
	r.StudioID = newNullIntPtr(i.StudioID)
	r.FileModTime = newNullTime(i.FileModTime)
	r.CreatedAt = i.CreatedAt
	r.UpdatedAt = i.UpdatedAt
}

type imageRowRecord struct {
	updateRecord
}

func (r *imageRowRecord) fromPartial(i models.ImagePartial) {
	r.setString("checksum", i.Checksum)
	r.setString("path", i.Path)
	r.setNullStringPtr("title", i.Title)
	r.setNullIntPtr("rating", i.Rating)
	r.setBool("organized", i.Organized)
	r.setInt("o_counter", i.OCounter)
	r.setNullInt64Ptr("size", i.Size)
	r.setNullIntPtr("width", i.Width)
	r.setNullIntPtr("height", i.Height)
	r.setNullIntPtr("studio_id", i.StudioID)
	r.setNullTimePtr("file_mod_time", i.FileModTime)
	r.setTime("created_at", i.CreatedAt)
	r.setTime("updated_at", i.UpdatedAt)
}

type imageQueryRow struct {
	imageRow
}

func (r *imageQueryRow) resolve() *models.Image {
	return &models.Image{
		ID:          r.ID,
		Checksum:    r.Checksum,
		Path:        r.Path,
		Title:       r.Title.stringPtr(),
		Rating:      r.Rating.intPtr(),
		Organized:   r.Organized,
		OCounter:    r.OCounter,
		Size:        r.Size.int64Ptr(),
		Width:       r.Width.intPtr(),
		Height:      r.Height.intPtr(),
		StudioID:    r.StudioID.intPtr(),
		FileModTime: r.FileModTime.timePtr(),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type imageQueryBuilder struct {
	repository
}

var ImageReaderWriter = &imageQueryBuilder{
	repository{
		tableName: imageTable,
		idColumn:  idColumn,
	},
}

func (qb *imageQueryBuilder) table() exp.IdentifierExpression {
	return imageTableMgr.table
}

func (qb *imageQueryBuilder) Create(ctx context.Context, newObject *models.Image) error {
	var r imageRow
	r.fromImage(*newObject)

	id, err := imageTableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	newObject.ID = id

	return nil
}

func (qb *imageQueryBuilder) UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error) {
	r := imageRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(partial)

	if err := imageTableMgr.updateByID(ctx, id, r.Record); err != nil {
		return nil, err
	}

	return qb.find(ctx, id)
}

func (qb *imageQueryBuilder) Update(ctx context.Context, updatedObject *models.Image) error {
	var r imageRow
	r.fromImage(*updatedObject)

	return imageTableMgr.updateByID(ctx, updatedObject.ID, r)
}

func (qb *imageQueryBuilder) getOCounter(ctx context.Context, id int) (int, error) {
	q := goqu.From(qb.table()).Select("o_counter").Where(goqu.Ex{"id": id})

	const single = true
	var ret int
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		if err := rows.Scan(&ret); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *imageQueryBuilder) IncrementOCounter(ctx context.Context, id int) (int, error) {
	if err := imageTableMgr.updateByID(ctx, id, goqu.Record{
		"o_counter": goqu.L("o_counter + 1"),
	}); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *imageQueryBuilder) DecrementOCounter(ctx context.Context, id int) (int, error) {
	if err := imageTableMgr.updateByID(ctx, id, goqu.Record{
		"o_counter": goqu.L("o_counter - 1"),
	}); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *imageQueryBuilder) ResetOCounter(ctx context.Context, id int) (int, error) {
	if err := imageTableMgr.updateByID(ctx, id, goqu.Record{
		"o_counter": 0,
	}); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *imageQueryBuilder) Destroy(ctx context.Context, id int) error {
	return imageTableMgr.destroyExisting(ctx, []int{id})
}

func (qb *imageQueryBuilder) Find(ctx context.Context, id int) (*models.Image, error) {
	return qb.find(ctx, id)
}

func (qb *imageQueryBuilder) FindMany(ctx context.Context, ids []int) ([]*models.Image, error) {
	var images []*models.Image
	for _, id := range ids {
		image, err := qb.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if image == nil {
			return nil, fmt.Errorf("image with id %d not found", id)
		}

		images = append(images, image)
	}

	return images, nil
}

func (qb *imageQueryBuilder) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	return goqu.From(table).Select(table.All())
}

func (qb *imageQueryBuilder) get(ctx context.Context, q *goqu.SelectDataset) (*models.Image, error) {
	var row imageQueryRow

	if err := imageTableMgr.get(ctx, q, &row); err != nil {
		return nil, err
	}

	return row.resolve(), nil
}

func (qb *imageQueryBuilder) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Image, error) {
	var ret []*models.Image

	const single = false
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var f imageQueryRow
		if err := rows.StructScan(&f); err != nil {
			return err
		}

		ret = append(ret, f.resolve())
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *imageQueryBuilder) find(ctx context.Context, id int) (*models.Image, error) {
	q := qb.selectDataset().Where(imageTableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting image by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *imageQueryBuilder) FindByChecksum(ctx context.Context, checksum string) (*models.Image, error) {
	q := qb.selectDataset().Where(qb.table().Col("checksum").Eq(checksum))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting image by checksum %s: %w", checksum, err)
	}

	return ret, nil
}

func (qb *imageQueryBuilder) FindByPath(ctx context.Context, path string) (*models.Image, error) {
	q := qb.selectDataset().Where(qb.table().Col("path").Eq(path))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting image by path %s: %w", path, err)
	}

	return ret, nil
}

func (qb *imageQueryBuilder) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Image, error) {
	table := qb.table()
	joinTable := goqu.T(galleriesImagesTable)

	q := qb.selectDataset().InnerJoin(
		joinTable,
		goqu.On(joinTable.Col("image_id").Eq(table.Col(idColumn))),
	).Where(
		joinTable.Col("gallery_id").Eq(galleryID),
	).GroupBy(table.Col(idColumn)).Order(table.Col("path").Asc())

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting images for gallery %d: %w", galleryID, err)
	}

	return ret, nil
}

func (qb *imageQueryBuilder) CountByGalleryID(ctx context.Context, galleryID int) (int, error) {
	joinTable := goqu.T(galleriesImagesTable)

	q := goqu.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col("gallery_id").Eq(galleryID))
	return count(ctx, q)
}

func (qb *imageQueryBuilder) Count(ctx context.Context) (int, error) {
	q := goqu.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *imageQueryBuilder) Size(ctx context.Context) (float64, error) {
	q := goqu.Select(goqu.SUM(qb.table().Col("size").Cast("double"))).From(qb.table())
	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *imageQueryBuilder) All(ctx context.Context) ([]*models.Image, error) {
	return qb.getMany(ctx, qb.selectDataset())
}

func (qb *imageQueryBuilder) validateFilter(imageFilter *models.ImageFilterType) error {
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

func (qb *imageQueryBuilder) makeFilter(ctx context.Context, imageFilter *models.ImageFilterType) *filterBuilder {
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

	query.handleCriterion(ctx, stringCriterionHandler(imageFilter.Checksum, "images.checksum"))
	query.handleCriterion(ctx, stringCriterionHandler(imageFilter.Title, "images.title"))
	query.handleCriterion(ctx, stringCriterionHandler(imageFilter.Path, "images.path"))
	query.handleCriterion(ctx, intCriterionHandler(imageFilter.Rating, "images.rating"))
	query.handleCriterion(ctx, intCriterionHandler(imageFilter.OCounter, "images.o_counter"))
	query.handleCriterion(ctx, boolCriterionHandler(imageFilter.Organized, "images.organized"))
	query.handleCriterion(ctx, resolutionCriterionHandler(imageFilter.Resolution, "images.height", "images.width"))
	query.handleCriterion(ctx, imageIsMissingCriterionHandler(qb, imageFilter.IsMissing))

	query.handleCriterion(ctx, imageTagsCriterionHandler(qb, imageFilter.Tags))
	query.handleCriterion(ctx, imageTagCountCriterionHandler(qb, imageFilter.TagCount))
	query.handleCriterion(ctx, imageGalleriesCriterionHandler(qb, imageFilter.Galleries))
	query.handleCriterion(ctx, imagePerformersCriterionHandler(qb, imageFilter.Performers))
	query.handleCriterion(ctx, imagePerformerCountCriterionHandler(qb, imageFilter.PerformerCount))
	query.handleCriterion(ctx, imageStudioCriterionHandler(qb, imageFilter.Studios))
	query.handleCriterion(ctx, imagePerformerTagsCriterionHandler(qb, imageFilter.PerformerTags))
	query.handleCriterion(ctx, imagePerformerFavoriteCriterionHandler(imageFilter.PerformerFavorite))

	return query
}

func (qb *imageQueryBuilder) makeQuery(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if imageFilter == nil {
		imageFilter = &models.ImageFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, imageTable)

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"images.title", "images.path", "images.checksum"}
		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(imageFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, imageFilter)

	query.addFilter(filter)

	query.sortAndPagination = qb.getImageSort(findFilter) + getPagination(findFilter)

	return &query, nil
}

func (qb *imageQueryBuilder) Query(ctx context.Context, options models.ImageQueryOptions) (*models.ImageQueryResult, error) {
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

func (qb *imageQueryBuilder) queryGroupedFields(ctx context.Context, options models.ImageQueryOptions, query queryBuilder) (*models.ImageQueryResult, error) {
	if !options.Count && !options.Megapixels && !options.TotalSize {
		// nothing to do - return empty result
		return models.NewImageQueryResult(qb), nil
	}

	aggregateQuery := qb.newQuery()

	if options.Count {
		aggregateQuery.addColumn("COUNT(temp.id) as total")
	}

	if options.Megapixels {
		query.addColumn("COALESCE(images.width, 0) * COALESCE(images.height, 0) / 1000000 as megapixels")
		aggregateQuery.addColumn("COALESCE(SUM(temp.megapixels), 0) as megapixels")
	}

	if options.TotalSize {
		query.addColumn("COALESCE(images.size, 0) as size")
		aggregateQuery.addColumn("COALESCE(SUM(temp.size), 0) as size")
	}

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

func (qb *imageQueryBuilder) QueryCount(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, imageFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

func imageIsMissingCriterionHandler(qb *imageQueryBuilder, isMissing *string) criterionHandlerFunc {
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

func (qb *imageQueryBuilder) getMultiCriterionHandlerBuilder(foreignTable, joinTable, foreignFK string, addJoinsFunc func(f *filterBuilder)) multiCriterionHandlerBuilder {
	return multiCriterionHandlerBuilder{
		primaryTable: imageTable,
		foreignTable: foreignTable,
		joinTable:    joinTable,
		primaryFK:    imageIDColumn,
		foreignFK:    foreignFK,
		addJoinsFunc: addJoinsFunc,
	}
}

func imageTagsCriterionHandler(qb *imageQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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

func imageTagCountCriterionHandler(qb *imageQueryBuilder, tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: imageTable,
		joinTable:    imagesTagsTable,
		primaryFK:    imageIDColumn,
	}

	return h.handler(tagCount)
}

func imageGalleriesCriterionHandler(qb *imageQueryBuilder, galleries *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		qb.galleriesRepository().join(f, "", "images.id")
		f.addLeftJoin(galleryTable, "", "galleries_images.gallery_id = galleries.id")
	}
	h := qb.getMultiCriterionHandlerBuilder(galleryTable, galleriesImagesTable, galleryIDColumn, addJoinsFunc)

	return h.handler(galleries)
}

func imagePerformersCriterionHandler(qb *imageQueryBuilder, performers *models.MultiCriterionInput) criterionHandlerFunc {
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

func imagePerformerCountCriterionHandler(qb *imageQueryBuilder, performerCount *models.IntCriterionInput) criterionHandlerFunc {
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

func imageStudioCriterionHandler(qb *imageQueryBuilder, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := hierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: imageTable,
		foreignTable: studioTable,
		foreignFK:    studioIDColumn,
		derivedTable: "studio",
		parentFK:     "parent_id",
	}

	return h.handler(studios)
}

func imagePerformerTagsCriterionHandler(qb *imageQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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

func (qb *imageQueryBuilder) getImageSort(findFilter *models.FindFilterType) string {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return ""
	}
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()

	switch sort {
	case "tag_count":
		return getCountSort(imageTable, imagesTagsTable, imageIDColumn, direction)
	case "performer_count":
		return getCountSort(imageTable, performersImagesTable, imageIDColumn, direction)
	default:
		return getSort(sort, direction, "images")
	}
}

func (qb *imageQueryBuilder) galleriesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesImagesTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: galleryIDColumn,
	}
}

func (qb *imageQueryBuilder) GetGalleryIDs(ctx context.Context, imageID int) ([]int, error) {
	return qb.galleriesRepository().getIDs(ctx, imageID)
}

func (qb *imageQueryBuilder) UpdateGalleries(ctx context.Context, imageID int, galleryIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.galleriesRepository().replace(ctx, imageID, galleryIDs)
}

func (qb *imageQueryBuilder) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersImagesTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: performerIDColumn,
	}
}

func (qb *imageQueryBuilder) GetPerformerIDs(ctx context.Context, imageID int) ([]int, error) {
	return qb.performersRepository().getIDs(ctx, imageID)
}

func (qb *imageQueryBuilder) UpdatePerformers(ctx context.Context, imageID int, performerIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.performersRepository().replace(ctx, imageID, performerIDs)
}

func (qb *imageQueryBuilder) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: imagesTagsTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: tagIDColumn,
	}
}

func (qb *imageQueryBuilder) GetTagIDs(ctx context.Context, imageID int) ([]int, error) {
	return qb.tagsRepository().getIDs(ctx, imageID)
}

func (qb *imageQueryBuilder) UpdateTags(ctx context.Context, imageID int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.tagsRepository().replace(ctx, imageID, tagIDs)
}
