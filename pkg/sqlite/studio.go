package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

const (
	studioTable        = "studios"
	studioIDColumn     = "studio_id"
	studioAliasesTable = "studio_aliases"
	studioAliasColumn  = "alias"

	studioImageBlobColumn = "image_blob"
)

type studioRow struct {
	ID        int                    `db:"id" goqu:"skipinsert"`
	Checksum  string                 `db:"checksum"`
	Name      zero.String            `db:"name"`
	URL       zero.String            `db:"url"`
	ParentID  null.Int               `db:"parent_id,omitempty"`
	CreatedAt models.SQLiteTimestamp `db:"created_at"`
	UpdatedAt models.SQLiteTimestamp `db:"updated_at"`
	// expressed as 1-100
	Rating        null.Int    `db:"rating"`
	Details       zero.String `db:"details"`
	IgnoreAutoTag bool        `db:"ignore_auto_tag"`

	// not used in resolutions or updates
	CoverBlob zero.String `db:"image_blob"`
}

func (r *studioRow) fromStudio(o models.Studio) {
	r.ID = o.ID
	r.Checksum = o.Checksum
	r.Name = zero.StringFrom(o.Name)
	r.URL = zero.StringFrom(o.URL)
	r.ParentID = intFromPtr(o.ParentID)
	r.CreatedAt = models.SQLiteTimestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = models.SQLiteTimestamp{Timestamp: o.UpdatedAt}
	r.Rating = intFromPtr(o.Rating)
	r.Details = zero.StringFrom(o.Details)
	r.IgnoreAutoTag = o.IgnoreAutoTag
}

func (r *studioRow) resolve() *models.Studio {
	ret := &models.Studio{
		ID:            r.ID,
		Checksum:      r.Checksum,
		Name:          r.Name.String,
		URL:           r.URL.String,
		ParentID:      nullIntPtr(r.ParentID),
		CreatedAt:     r.CreatedAt.Timestamp,
		UpdatedAt:     r.UpdatedAt.Timestamp,
		Rating:        nullIntPtr(r.Rating),
		Details:       r.Details.String,
		IgnoreAutoTag: r.IgnoreAutoTag,
	}

	return ret
}

type studioRowRecord struct {
	updateRecord
}

func (r *studioRowRecord) fromPartial(o models.StudioPartial) {
	r.setString("checksum", o.Checksum)
	r.setNullString("name", o.Name)
	r.setNullString("url", o.URL)
	r.setNullInt("parent_id", o.ParentID)
	r.setSQLiteTimestamp("created_at", o.CreatedAt)
	r.setSQLiteTimestamp("updated_at", o.UpdatedAt)
	r.setNullInt("rating", o.Rating)
	r.setNullString("details", o.Details)
	r.setBool("ignore_auto_tag", o.IgnoreAutoTag)
}

type StudioStore struct {
	repository
	blobJoinQueryBuilder

	tableMgr *table
}

func NewStudioStore(blobStore *BlobStore) *StudioStore {
	return &StudioStore{
		repository: repository{
			tableName: studioTable,
			idColumn:  idColumn,
		},
		blobJoinQueryBuilder: blobJoinQueryBuilder{
			blobStore: blobStore,
			joinTable: studioTable,
		},

		tableMgr: studioTableMgr,
	}
}

func (qb *StudioStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *StudioStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *StudioStore) Create(ctx context.Context, newObject *models.Studio) error {
	var r studioRow
	r.fromStudio(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated

	return nil
}

func (qb *StudioStore) UpdatePartial(ctx context.Context, id int, partial models.StudioPartial) (*models.Studio, error) {
	r := studioRowRecord{
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

	return qb.find(ctx, id)
}

func (qb *StudioStore) Update(ctx context.Context, updatedObject *models.Studio) error {
	var r studioRow
	r.fromStudio(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *StudioStore) Destroy(ctx context.Context, id int) error {
	// must handle image checksums manually
	if err := qb.destroyImage(ctx, id); err != nil {
		return err
	}

	// TODO - set null on foreign key in scraped items
	// remove studio from scraped items
	_, err := qb.tx.Exec(ctx, "UPDATE scraped_items SET studio_id = null WHERE studio_id = ?", id)
	if err != nil {
		return err
	}

	return qb.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *StudioStore) Find(ctx context.Context, id int) (*models.Studio, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *StudioStore) FindMany(ctx context.Context, ids []int) ([]*models.Studio, error) {
	ret := make([]*models.Studio, len(ids))

	table := qb.table()
	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(batch))
		unsorted, err := qb.getMany(ctx, q)
		if err != nil {
			return err
		}

		for _, s := range unsorted {
			i := intslice.IntIndex(ids, s.ID)
			ret[i] = s
		}

		return nil
	}); err != nil {
		return nil, err
	}

	for i := range ret {
		if ret[i] == nil {
			return nil, fmt.Errorf("studio with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *StudioStore) find(ctx context.Context, id int) (*models.Studio, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *StudioStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Studio, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *StudioStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Studio, error) {
	const single = false
	var ret []*models.Studio
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f studioRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *StudioStore) FindChildren(ctx context.Context, id int) ([]*models.Studio, error) {
	// SELECT studios.* FROM studios WHERE studios.parent_id = ?
	table := qb.table()
	sq := qb.selectDataset().Where(table.Col("parent_id").Eq(id))
	ret, err := qb.getMany(ctx, sq)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *StudioStore) FindBySceneID(ctx context.Context, sceneID int) (*models.Studio, error) {
	// SELECT studios.* FROM studios JOIN scenes ON studios.id = scenes.studio_id WHERE scenes.id = ? LIMIT 1
	table := qb.table()
	scenes := sceneTableMgr.table
	sq := qb.selectDataset().Join(
		scenes, goqu.On(table.Col(idColumn), scenes.Col(studioIDColumn)),
	).Where(
		scenes.Col(idColumn),
	).Limit(1)
	ret, err := qb.get(ctx, sq)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return ret, nil
}

func (qb *StudioStore) FindByName(ctx context.Context, name string, nocase bool) (*models.Studio, error) {
	// query := "SELECT * FROM studios WHERE name = ?"
	// if nocase {
	// 	query += " COLLATE NOCASE"
	// }
	// query += " LIMIT 1"
	where := "name = ?"
	if nocase {
		where += " COLLATE NOCASE"
	}
	sq := qb.selectDataset().Prepared(true).Where(goqu.L(where, name)).Limit(1)
	ret, err := qb.get(ctx, sq)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return ret, nil
}

func (qb *StudioStore) FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Studio, error) {
	query := selectAll("studios") + `
		LEFT JOIN studio_stash_ids on studio_stash_ids.studio_id = studios.id
		WHERE studio_stash_ids.stash_id = ?
		AND studio_stash_ids.endpoint = ?
	`
	args := []interface{}{stashID.StashID, stashID.Endpoint}
	return qb.queryStudios(ctx, query, args)
}

func (qb *StudioStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *StudioStore) All(ctx context.Context) ([]*models.Studio, error) {
	table := qb.table()

	return qb.getMany(ctx, qb.selectDataset().Order(
		table.Col("name").Asc(),
		table.Col(idColumn).Asc(),
	))
}

func (qb *StudioStore) QueryForAutoTag(ctx context.Context, words []string) ([]*models.Studio, error) {
	// TODO - Query needs to be changed to support queries of this type, and
	// this method should be removed
	query := selectAll(studioTable)
	query += " LEFT JOIN studio_aliases ON studio_aliases.studio_id = studios.id"

	var whereClauses []string
	var args []interface{}

	for _, w := range words {
		ww := w + "%"
		whereClauses = append(whereClauses, "studios.name like ?")
		args = append(args, ww)

		// include aliases
		whereClauses = append(whereClauses, "studio_aliases.alias like ?")
		args = append(args, ww)
	}

	whereOr := "(" + strings.Join(whereClauses, " OR ") + ")"
	where := strings.Join([]string{
		"studios.ignore_auto_tag = 0",
		whereOr,
	}, " AND ")
	return qb.queryStudios(ctx, query+" WHERE "+where, args)
}

func (qb *StudioStore) validateFilter(filter *models.StudioFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if filter.And != nil {
		if filter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if filter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(filter.And)
	}

	if filter.Or != nil {
		if filter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(filter.Or)
	}

	if filter.Not != nil {
		return qb.validateFilter(filter.Not)
	}

	return nil
}

func (qb *StudioStore) makeFilter(ctx context.Context, studioFilter *models.StudioFilterType) *filterBuilder {
	query := &filterBuilder{}

	if studioFilter.And != nil {
		query.and(qb.makeFilter(ctx, studioFilter.And))
	}
	if studioFilter.Or != nil {
		query.or(qb.makeFilter(ctx, studioFilter.Or))
	}
	if studioFilter.Not != nil {
		query.not(qb.makeFilter(ctx, studioFilter.Not))
	}

	query.handleCriterion(ctx, stringCriterionHandler(studioFilter.Name, studioTable+".name"))
	query.handleCriterion(ctx, stringCriterionHandler(studioFilter.Details, studioTable+".details"))
	query.handleCriterion(ctx, stringCriterionHandler(studioFilter.URL, studioTable+".url"))
	query.handleCriterion(ctx, intCriterionHandler(studioFilter.Rating100, studioTable+".rating", nil))
	// legacy rating handler
	query.handleCriterion(ctx, rating5CriterionHandler(studioFilter.Rating, studioTable+".rating", nil))
	query.handleCriterion(ctx, boolCriterionHandler(studioFilter.IgnoreAutoTag, studioTable+".ignore_auto_tag", nil))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if studioFilter.StashID != nil {
			qb.stashIDRepository().join(f, "studio_stash_ids", "studios.id")
			stringCriterionHandler(studioFilter.StashID, "studio_stash_ids.stash_id")(ctx, f)
		}
	}))
	query.handleCriterion(ctx, &stashIDCriterionHandler{
		c:                 studioFilter.StashIDEndpoint,
		stashIDRepository: qb.stashIDRepository(),
		stashIDTableAs:    "studio_stash_ids",
		parentIDCol:       "studios.id",
	})

	query.handleCriterion(ctx, studioIsMissingCriterionHandler(qb, studioFilter.IsMissing))
	query.handleCriterion(ctx, studioSceneCountCriterionHandler(qb, studioFilter.SceneCount))
	query.handleCriterion(ctx, studioImageCountCriterionHandler(qb, studioFilter.ImageCount))
	query.handleCriterion(ctx, studioGalleryCountCriterionHandler(qb, studioFilter.GalleryCount))
	query.handleCriterion(ctx, studioParentCriterionHandler(qb, studioFilter.Parents))
	query.handleCriterion(ctx, studioAliasCriterionHandler(qb, studioFilter.Aliases))
	query.handleCriterion(ctx, timestampCriterionHandler(studioFilter.CreatedAt, "studios.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(studioFilter.UpdatedAt, "studios.updated_at"))

	return query
}

func (qb *StudioStore) Query(ctx context.Context, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int, error) {
	if studioFilter == nil {
		studioFilter = &models.StudioFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, studioTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.join(studioAliasesTable, "", "studio_aliases.studio_id = studios.id")
		searchColumns := []string{"studios.name", "studio_aliases.alias"}

		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(studioFilter); err != nil {
		return nil, 0, err
	}
	filter := qb.makeFilter(ctx, studioFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, 0, err
	}

	query.sortAndPagination = qb.getStudioSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	studios, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return studios, countResult, nil
}

func studioIsMissingCriterionHandler(qb *StudioStore, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "image":
				f.addWhere("studios.image_blob IS NULL")
			case "stash_id":
				qb.stashIDRepository().join(f, "studio_stash_ids", "studios.id")
				f.addWhere("studio_stash_ids.studio_id IS NULL")
			default:
				f.addWhere("(studios." + *isMissing + " IS NULL OR TRIM(studios." + *isMissing + ") = '')")
			}
		}
	}
}

func studioSceneCountCriterionHandler(qb *StudioStore, sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes", "", "scenes.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes.id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioImageCountCriterionHandler(qb *StudioStore, imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images", "", "images.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct images.id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioGalleryCountCriterionHandler(qb *StudioStore, galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries", "", "galleries.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries.id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioParentCriterionHandler(qb *StudioStore, parents *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		f.addLeftJoin("studios", "parent_studio", "parent_studio.id = studios.parent_id")
	}
	h := multiCriterionHandlerBuilder{
		primaryTable: studioTable,
		foreignTable: "parent_studio",
		joinTable:    "",
		primaryFK:    studioIDColumn,
		foreignFK:    "parent_id",
		addJoinsFunc: addJoinsFunc,
	}
	return h.handler(parents)
}

func studioAliasCriterionHandler(qb *StudioStore, alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    studioAliasesTable,
		stringColumn: studioAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			qb.aliasRepository().join(f, "", "studios.id")
		},
	}

	return h.handler(alias)
}

func (qb *StudioStore) getStudioSort(findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}

	sortQuery := ""
	switch sort {
	case "scenes_count":
		sortQuery += getCountSort(studioTable, sceneTable, studioIDColumn, direction)
	case "images_count":
		sortQuery += getCountSort(studioTable, imageTable, studioIDColumn, direction)
	case "galleries_count":
		sortQuery += getCountSort(studioTable, galleryTable, studioIDColumn, direction)
	default:
		sortQuery += getSort(sort, direction, "studios")
	}

	// Whatever the sorting, always use name/id as a final sort
	sortQuery += ", COALESCE(studios.name, studios.id) COLLATE NATURAL_CI ASC"
	return sortQuery
}

func (qb *StudioStore) queryStudios(ctx context.Context, query string, args []interface{}) ([]*models.Studio, error) {
	const single = false
	var ret []*models.Studio
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f studioRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *StudioStore) GetImage(ctx context.Context, studioID int) ([]byte, error) {
	return qb.blobJoinQueryBuilder.GetImage(ctx, studioID, studioImageBlobColumn)
}

func (qb *StudioStore) HasImage(ctx context.Context, studioID int) (bool, error) {
	return qb.blobJoinQueryBuilder.HasImage(ctx, studioID, studioImageBlobColumn)
}

func (qb *StudioStore) UpdateImage(ctx context.Context, studioID int, image []byte) error {
	return qb.blobJoinQueryBuilder.UpdateImage(ctx, studioID, studioImageBlobColumn, image)
}

func (qb *StudioStore) destroyImage(ctx context.Context, studioID int) error {
	return qb.blobJoinQueryBuilder.DestroyImage(ctx, studioID, studioImageBlobColumn)
}

func (qb *StudioStore) stashIDRepository() *stashIDRepository {
	return &stashIDRepository{
		repository{
			tx:        qb.tx,
			tableName: "studio_stash_ids",
			idColumn:  studioIDColumn,
		},
	}
}

func (qb *StudioStore) GetStashIDs(ctx context.Context, studioID int) ([]models.StashID, error) {
	return qb.stashIDRepository().get(ctx, studioID)
}

func (qb *StudioStore) UpdateStashIDs(ctx context.Context, studioID int, stashIDs []models.StashID) error {
	return qb.stashIDRepository().replace(ctx, studioID, stashIDs)
}

func (qb *StudioStore) aliasRepository() *stringRepository {
	return &stringRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: studioAliasesTable,
			idColumn:  studioIDColumn,
		},
		stringColumn: studioAliasColumn,
	}
}

func (qb *StudioStore) GetAliases(ctx context.Context, studioID int) ([]string, error) {
	return qb.aliasRepository().get(ctx, studioID)
}

func (qb *StudioStore) UpdateAliases(ctx context.Context, studioID int, aliases []string) error {
	return qb.aliasRepository().replace(ctx, studioID, aliases)
}
