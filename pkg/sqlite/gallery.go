package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

const galleryTable = "galleries"

const performersGalleriesTable = "performers_galleries"
const galleriesTagsTable = "galleries_tags"
const galleriesImagesTable = "galleries_images"
const galleriesScenesTable = "scenes_galleries"
const galleryIDColumn = "gallery_id"

type galleryRow struct {
	ID          int               `db:"id" goqu:"skipinsert"`
	Path        zero.String       `db:"path"`
	Checksum    string            `db:"checksum"`
	Zip         bool              `db:"zip"`
	Title       zero.String       `db:"title"`
	URL         zero.String       `db:"url"`
	Date        models.SQLiteDate `db:"date"`
	Details     zero.String       `db:"details"`
	Rating      null.Int          `db:"rating"`
	Organized   bool              `db:"organized"`
	StudioID    null.Int          `db:"studio_id,omitempty"`
	FileModTime null.Time         `db:"file_mod_time"`
	CreatedAt   time.Time         `db:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at"`
}

func (r *galleryRow) fromGallery(o models.Gallery) {
	r.ID = o.ID
	r.Path = zero.StringFromPtr(o.Path)
	r.Checksum = o.Checksum
	r.Zip = o.Zip
	r.Title = zero.StringFrom(o.Title)
	r.URL = zero.StringFrom(o.URL)
	if o.Date != nil {
		_ = r.Date.Scan(o.Date.Time)
	}
	r.Details = zero.StringFrom(o.Details)
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.StudioID = intFromPtr(o.StudioID)
	r.FileModTime = null.TimeFromPtr(o.FileModTime)
	r.CreatedAt = o.CreatedAt
	r.UpdatedAt = o.UpdatedAt
}

type galleryRowRecord struct {
	updateRecord
}

func (r *galleryRowRecord) fromPartial(o models.GalleryPartial) {
	r.setNullString("path", o.Path)
	r.setString("checksum", o.Checksum)
	r.setBool("zip", o.Zip)
	r.setNullString("title", o.Title)
	r.setNullString("url", o.URL)
	r.setSQLiteDate("date", o.Date)
	r.setNullString("details", o.Details)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setNullInt("studio_id", o.StudioID)
	r.setNullTime("file_mod_time", o.FileModTime)
	r.setTime("created_at", o.CreatedAt)
	r.setTime("updated_at", o.UpdatedAt)
}

type galleryQueryRow struct {
	galleryRow

	SceneID     null.Int `db:"scene_id"`
	TagID       null.Int `db:"tag_id"`
	PerformerID null.Int `db:"performer_id"`
}

func (r *galleryQueryRow) resolve() *models.Gallery {
	ret := &models.Gallery{
		ID:          r.ID,
		Path:        r.Path.Ptr(),
		Checksum:    r.Checksum,
		Zip:         r.Zip,
		Title:       r.Title.String,
		URL:         r.URL.String,
		Date:        r.Date.DatePtr(),
		Details:     r.Details.String,
		Rating:      nullIntPtr(r.Rating),
		Organized:   r.Organized,
		StudioID:    nullIntPtr(r.StudioID),
		FileModTime: r.FileModTime.Ptr(),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	r.appendRelationships(ret)

	return ret
}

func (r *galleryQueryRow) appendRelationships(i *models.Gallery) {
	if r.TagID.Valid {
		i.TagIDs = intslice.IntAppendUnique(i.TagIDs, int(r.TagID.Int64))
	}
	if r.PerformerID.Valid {
		i.PerformerIDs = intslice.IntAppendUnique(i.PerformerIDs, int(r.PerformerID.Int64))
	}
	if r.SceneID.Valid {
		i.SceneIDs = intslice.IntAppendUnique(i.SceneIDs, int(r.SceneID.Int64))
	}
}

type galleryQueryRows []galleryQueryRow

func (r galleryQueryRows) resolve() []*models.Gallery {
	var ret []*models.Gallery
	var last *models.Gallery
	var lastID int

	for _, row := range r {
		if last == nil || lastID != row.ID {
			f := row.resolve()
			last = f
			lastID = row.ID
			ret = append(ret, last)
			continue
		}

		// must be merging with previous row
		row.appendRelationships(last)
	}

	return ret
}

type galleryQueryBuilder struct {
	repository

	tableMgr *table
}

var GalleryReaderWriter = &galleryQueryBuilder{
	repository: repository{
		tableName: galleryTable,
		idColumn:  idColumn,
	},
	tableMgr: galleryTableMgr,
}

func (qb *galleryQueryBuilder) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *galleryQueryBuilder) Create(ctx context.Context, newObject *models.Gallery) error {
	var r galleryRow
	r.fromGallery(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if err := galleriesPerformersTableMgr.insertJoins(ctx, id, newObject.PerformerIDs); err != nil {
		return err
	}
	if err := galleriesTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs); err != nil {
		return err
	}
	if err := galleriesScenesTableMgr.insertJoins(ctx, id, newObject.SceneIDs); err != nil {
		return err
	}

	// only assign id once we are successful
	newObject.ID = id

	return nil
}

func (qb *galleryQueryBuilder) Update(ctx context.Context, updatedObject *models.Gallery) error {
	var r galleryRow
	r.fromGallery(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if err := galleriesPerformersTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.PerformerIDs); err != nil {
		return err
	}
	if err := galleriesTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs); err != nil {
		return err
	}
	if err := galleriesScenesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.SceneIDs); err != nil {
		return err
	}

	return nil
}

func (qb *galleryQueryBuilder) UpdatePartial(ctx context.Context, id int, partial models.GalleryPartial) (*models.Gallery, error) {
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

	return qb.Find(ctx, id)
}

func (qb *galleryQueryBuilder) Destroy(ctx context.Context, id int) error {
	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

func (qb *galleryQueryBuilder) selectDataset() *goqu.SelectDataset {
	table := qb.table()

	return dialect.From(table).Select(
		table.All(),
		galleriesTagsJoinTable.Col("tag_id"),
		performersGalleriesJoinTable.Col("performer_id"),
		galleriesScenesJoinTable.Col("scene_id"),
	).LeftJoin(
		galleriesTagsJoinTable,
		goqu.On(table.Col(idColumn).Eq(galleriesTagsJoinTable.Col(galleryIDColumn))),
	).LeftJoin(
		performersGalleriesJoinTable,
		goqu.On(table.Col(idColumn).Eq(performersGalleriesJoinTable.Col(galleryIDColumn))),
	).LeftJoin(
		galleriesScenesJoinTable,
		goqu.On(table.Col(idColumn).Eq(galleriesScenesJoinTable.Col(galleryIDColumn))),
	)
}

func (qb *galleryQueryBuilder) get(ctx context.Context, q *goqu.SelectDataset) (*models.Gallery, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *galleryQueryBuilder) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Gallery, error) {
	const single = false
	var rows galleryQueryRows
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f galleryQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		rows = append(rows, f)
		return nil
	}); err != nil {
		return nil, err
	}

	return rows.resolve(), nil
}

func (qb *galleryQueryBuilder) Find(ctx context.Context, id int) (*models.Gallery, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting gallery by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *galleryQueryBuilder) FindMany(ctx context.Context, ids []int) ([]*models.Gallery, error) {
	var galleries []*models.Gallery
	for _, id := range ids {
		gallery, err := qb.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if gallery == nil {
			return nil, fmt.Errorf("gallery with id %d not found", id)
		}

		galleries = append(galleries, gallery)
	}

	return galleries, nil
}

func (qb *galleryQueryBuilder) FindByChecksum(ctx context.Context, checksum string) (*models.Gallery, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("checksum").Eq(checksum))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting gallery by checksum %s: %w", checksum, err)
	}

	return ret, nil
}

func (qb *galleryQueryBuilder) FindByChecksums(ctx context.Context, checksums []string) ([]*models.Gallery, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("checksum").In(checksums))

	ret, err := qb.getMany(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting gallery by checksums: %w", err)
	}

	return ret, nil
}

func (qb *galleryQueryBuilder) FindByPath(ctx context.Context, path string) (*models.Gallery, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("path").Eq(path))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting gallery by path %s: %w", path, err)
	}

	return ret, nil
}

func (qb *galleryQueryBuilder) FindBySceneID(ctx context.Context, sceneID int) ([]*models.Gallery, error) {
	table := qb.table()

	q := qb.selectDataset().Where(
		galleriesScenesJoinTable.Col("scene_id").Eq(sceneID),
	).GroupBy(table.Col(idColumn))

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting galleries for scene %d: %w", sceneID, err)
	}

	return ret, nil
}

func (qb *galleryQueryBuilder) FindByImageID(ctx context.Context, imageID int) ([]*models.Gallery, error) {
	table := qb.table()

	q := qb.selectDataset().InnerJoin(
		galleriesImagesJoinTable,
		goqu.On(table.Col(idColumn).Eq(galleriesImagesJoinTable.Col(galleryIDColumn))),
	).Where(
		galleriesImagesJoinTable.Col("image_id").Eq(imageID),
	).GroupBy(table.Col(idColumn))

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting galleries for image %d: %w", imageID, err)
	}

	return ret, nil
}

func (qb *galleryQueryBuilder) CountByImageID(ctx context.Context, imageID int) (int, error) {
	joinTable := galleriesImagesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(imageIDColumn).Eq(imageID))
	return count(ctx, q)
}

func (qb *galleryQueryBuilder) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *galleryQueryBuilder) All(ctx context.Context) ([]*models.Gallery, error) {
	return qb.getMany(ctx, qb.selectDataset())
}

func (qb *galleryQueryBuilder) validateFilter(galleryFilter *models.GalleryFilterType) error {
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

func (qb *galleryQueryBuilder) makeFilter(ctx context.Context, galleryFilter *models.GalleryFilterType) *filterBuilder {
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

	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.Title, "galleries.title"))
	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.Details, "galleries.details"))
	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.Checksum, "galleries.checksum"))
	query.handleCriterion(ctx, boolCriterionHandler(galleryFilter.IsZip, "galleries.zip"))
	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.Path, "galleries.path"))
	query.handleCriterion(ctx, intCriterionHandler(galleryFilter.Rating, "galleries.rating"))
	query.handleCriterion(ctx, stringCriterionHandler(galleryFilter.URL, "galleries.url"))
	query.handleCriterion(ctx, boolCriterionHandler(galleryFilter.Organized, "galleries.organized"))
	query.handleCriterion(ctx, galleryIsMissingCriterionHandler(qb, galleryFilter.IsMissing))
	query.handleCriterion(ctx, galleryTagsCriterionHandler(qb, galleryFilter.Tags))
	query.handleCriterion(ctx, galleryTagCountCriterionHandler(qb, galleryFilter.TagCount))
	query.handleCriterion(ctx, galleryPerformersCriterionHandler(qb, galleryFilter.Performers))
	query.handleCriterion(ctx, galleryPerformerCountCriterionHandler(qb, galleryFilter.PerformerCount))
	query.handleCriterion(ctx, galleryStudioCriterionHandler(qb, galleryFilter.Studios))
	query.handleCriterion(ctx, galleryPerformerTagsCriterionHandler(qb, galleryFilter.PerformerTags))
	query.handleCriterion(ctx, galleryAverageResolutionCriterionHandler(qb, galleryFilter.AverageResolution))
	query.handleCriterion(ctx, galleryImageCountCriterionHandler(qb, galleryFilter.ImageCount))
	query.handleCriterion(ctx, galleryPerformerFavoriteCriterionHandler(galleryFilter.PerformerFavorite))
	query.handleCriterion(ctx, galleryPerformerAgeCriterionHandler(galleryFilter.PerformerAge))

	return query
}

func (qb *galleryQueryBuilder) makeQuery(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if galleryFilter == nil {
		galleryFilter = &models.GalleryFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, galleryTable)

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"galleries.title", "galleries.path", "galleries.checksum"}
		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(galleryFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, galleryFilter)

	query.addFilter(filter)

	query.sortAndPagination = qb.getGallerySort(findFilter) + getPagination(findFilter)

	return &query, nil
}

func (qb *galleryQueryBuilder) Query(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) ([]*models.Gallery, int, error) {
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

func (qb *galleryQueryBuilder) QueryCount(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, galleryFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

func galleryIsMissingCriterionHandler(qb *galleryQueryBuilder, isMissing *string) criterionHandlerFunc {
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

func galleryTagsCriterionHandler(qb *galleryQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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

func galleryTagCountCriterionHandler(qb *galleryQueryBuilder, tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesTagsTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(tagCount)
}

func galleryPerformersCriterionHandler(qb *galleryQueryBuilder, performers *models.MultiCriterionInput) criterionHandlerFunc {
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

func galleryPerformerCountCriterionHandler(qb *galleryQueryBuilder, performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    performersGalleriesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(performerCount)
}

func galleryImageCountCriterionHandler(qb *galleryQueryBuilder, imageCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: galleryTable,
		joinTable:    galleriesImagesTable,
		primaryFK:    galleryIDColumn,
	}

	return h.handler(imageCount)
}

func galleryStudioCriterionHandler(qb *galleryQueryBuilder, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := hierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: galleryTable,
		foreignTable: studioTable,
		foreignFK:    studioIDColumn,
		derivedTable: "studio",
		parentFK:     "parent_id",
	}

	return h.handler(studios)
}

func galleryPerformerTagsCriterionHandler(qb *galleryQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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

func galleryAverageResolutionCriterionHandler(qb *galleryQueryBuilder, resolution *models.ResolutionCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if resolution != nil && resolution.Value.IsValid() {
			qb.imagesRepository().join(f, "images_join", "galleries.id")
			f.addLeftJoin("images", "", "images_join.image_id = images.id")

			min := resolution.Value.GetMinResolution()
			max := resolution.Value.GetMaxResolution()

			const widthHeight = "avg(MIN(images.width, images.height))"

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

func (qb *galleryQueryBuilder) getGallerySort(findFilter *models.FindFilterType) string {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return ""
	}

	sort := findFilter.GetSort("path")
	direction := findFilter.GetDirection()

	switch sort {
	case "images_count":
		return getCountSort(galleryTable, galleriesImagesTable, galleryIDColumn, direction)
	case "tag_count":
		return getCountSort(galleryTable, galleriesTagsTable, galleryIDColumn, direction)
	case "performer_count":
		return getCountSort(galleryTable, performersGalleriesTable, galleryIDColumn, direction)
	default:
		return getSort(sort, direction, "galleries")
	}
}

func (qb *galleryQueryBuilder) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersGalleriesTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "performer_id",
	}
}

func (qb *galleryQueryBuilder) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesTagsTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "tag_id",
	}
}

func (qb *galleryQueryBuilder) imagesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesImagesTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "image_id",
	}
}

func (qb *galleryQueryBuilder) GetImageIDs(ctx context.Context, galleryID int) ([]int, error) {
	return qb.imagesRepository().getIDs(ctx, galleryID)
}

func (qb *galleryQueryBuilder) UpdateImages(ctx context.Context, galleryID int, imageIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.imagesRepository().replace(ctx, galleryID, imageIDs)
}
