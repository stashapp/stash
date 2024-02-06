package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/studio"
)

const (
	studioTable           = "studios"
	studioIDColumn        = "studio_id"
	studioAliasesTable    = "studio_aliases"
	studioAliasColumn     = "alias"
	studioParentIDColumn  = "parent_id"
	studioNameColumn      = "name"
	studioImageBlobColumn = "image_blob"
)

type studioRow struct {
	ID        int         `db:"id" goqu:"skipinsert"`
	Name      zero.String `db:"name"`
	URL       zero.String `db:"url"`
	ParentID  null.Int    `db:"parent_id,omitempty"`
	CreatedAt Timestamp   `db:"created_at"`
	UpdatedAt Timestamp   `db:"updated_at"`
	// expressed as 1-100
	Rating        null.Int    `db:"rating"`
	Details       zero.String `db:"details"`
	IgnoreAutoTag bool        `db:"ignore_auto_tag"`

	// not used in resolutions or updates
	ImageBlob zero.String `db:"image_blob"`
}

func (r *studioRow) fromStudio(o models.Studio) {
	r.ID = o.ID
	r.Name = zero.StringFrom(o.Name)
	r.URL = zero.StringFrom(o.URL)
	r.ParentID = intFromPtr(o.ParentID)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
	r.Rating = intFromPtr(o.Rating)
	r.Details = zero.StringFrom(o.Details)
	r.IgnoreAutoTag = o.IgnoreAutoTag
}

func (r *studioRow) resolve() *models.Studio {
	ret := &models.Studio{
		ID:            r.ID,
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
	r.setNullString("name", o.Name)
	r.setNullString("url", o.URL)
	r.setNullInt("parent_id", o.ParentID)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
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
	var err error

	var r studioRow
	r.fromStudio(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if newObject.Aliases.Loaded() {
		if err := studio.EnsureAliasesUnique(ctx, id, newObject.Aliases.List(), qb); err != nil {
			return err
		}

		if err := studiosAliasesTableMgr.insertJoins(ctx, id, newObject.Aliases.List()); err != nil {
			return err
		}
	}

	if newObject.StashIDs.Loaded() {
		if err := studiosStashIDsTableMgr.insertJoins(ctx, id, newObject.StashIDs.List()); err != nil {
			return err
		}
	}

	updated, _ := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated
	return nil
}

func (qb *StudioStore) UpdatePartial(ctx context.Context, input models.StudioPartial) (*models.Studio, error) {
	r := studioRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(input)

	if len(r.Record) > 0 {
		if err := qb.tableMgr.updateByID(ctx, input.ID, r.Record); err != nil {
			return nil, err
		}
	}

	if input.Aliases != nil {
		if err := studio.EnsureAliasesUnique(ctx, input.ID, input.Aliases.Values, qb); err != nil {
			return nil, err
		}

		if err := studiosAliasesTableMgr.modifyJoins(ctx, input.ID, input.Aliases.Values, input.Aliases.Mode); err != nil {
			return nil, err
		}
	}

	if input.StashIDs != nil {
		if err := studiosStashIDsTableMgr.modifyJoins(ctx, input.ID, input.StashIDs.StashIDs, input.StashIDs.Mode); err != nil {
			return nil, err
		}
	}

	return qb.Find(ctx, input.ID)
}

// This is only used by the Import/Export functionality
func (qb *StudioStore) Update(ctx context.Context, updatedObject *models.Studio) error {
	var r studioRow
	r.fromStudio(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.Aliases.Loaded() {
		if err := studiosAliasesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.Aliases.List()); err != nil {
			return err
		}
	}

	if updatedObject.StashIDs.Loaded() {
		if err := studiosStashIDsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.StashIDs.List()); err != nil {
			return err
		}
	}

	return nil
}

func (qb *StudioStore) Destroy(ctx context.Context, id int) error {
	// must handle image checksums manually
	if err := qb.destroyImage(ctx, id); err != nil {
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
			i := sliceutil.Index(ids, s.ID)
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

func (qb *StudioStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Studio, error) {
	table := qb.table()

	q := qb.selectDataset().Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

func (qb *StudioStore) FindChildren(ctx context.Context, id int) ([]*models.Studio, error) {
	// SELECT studios.* FROM studios WHERE studios.parent_id = ?
	table := qb.table()
	sq := qb.selectDataset().Where(table.Col(studioParentIDColumn).Eq(id))
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
	sq := dialect.From(studiosStashIDsJoinTable).Select(studiosStashIDsJoinTable.Col(studioIDColumn)).Where(
		studiosStashIDsJoinTable.Col("stash_id").Eq(stashID.StashID),
		studiosStashIDsJoinTable.Col("endpoint").Eq(stashID.Endpoint),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting studios for stash ID %s: %w", stashID.StashID, err)
	}

	return ret, nil
}

func (qb *StudioStore) FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*models.Studio, error) {
	table := qb.table()
	sq := dialect.From(table).LeftJoin(
		studiosStashIDsJoinTable,
		goqu.On(table.Col(idColumn).Eq(studiosStashIDsJoinTable.Col(studioIDColumn))),
	).Select(table.Col(idColumn))

	if hasStashID {
		sq = sq.Where(
			studiosStashIDsJoinTable.Col("stash_id").IsNotNull(),
			studiosStashIDsJoinTable.Col("endpoint").Eq(stashboxEndpoint),
		)
	} else {
		sq = sq.Where(
			studiosStashIDsJoinTable.Col("stash_id").IsNull(),
		)
	}

	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting studios for stash-box endpoint %s: %w", stashboxEndpoint, err)
	}

	return ret, nil
}

func (qb *StudioStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *StudioStore) All(ctx context.Context) ([]*models.Studio, error) {
	table := qb.table()
	return qb.getMany(ctx, qb.selectDataset().Order(table.Col(studioNameColumn).Asc()))
}

func (qb *StudioStore) QueryForAutoTag(ctx context.Context, words []string) ([]*models.Studio, error) {
	// TODO - Query needs to be changed to support queries of this type, and
	// this method should be removed
	table := qb.table()
	sq := dialect.From(table).Select(table.Col(idColumn)).LeftJoin(
		studiosAliasesJoinTable,
		goqu.On(studiosAliasesJoinTable.Col(studioIDColumn).Eq(table.Col(idColumn))),
	)

	var whereClauses []exp.Expression

	for _, w := range words {
		whereClauses = append(whereClauses, table.Col(studioNameColumn).Like(w+"%"))
		whereClauses = append(whereClauses, studiosAliasesJoinTable.Col("alias").Like(w+"%"))
	}

	sq = sq.Where(
		goqu.Or(whereClauses...),
		table.Col("ignore_auto_tag").Eq(0),
	)

	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting performers for autotag: %w", err)
	}

	return ret, nil
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
	query.handleCriterion(ctx, studioChildCountCriterionHandler(qb, studioFilter.ChildCount))
	query.handleCriterion(ctx, timestampCriterionHandler(studioFilter.CreatedAt, studioTable+".created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(studioFilter.UpdatedAt, studioTable+".updated_at"))

	return query
}

func (qb *StudioStore) makeQuery(ctx context.Context, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
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
		return nil, err
	}
	filter := qb.makeFilter(ctx, studioFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	query.sortAndPagination = qb.getStudioSort(findFilter) + getPagination(findFilter)

	return &query, nil
}

func (qb *StudioStore) Query(ctx context.Context, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int, error) {
	query, err := qb.makeQuery(ctx, studioFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

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
			studiosAliasesTableMgr.join(f, "", "studios.id")
		},
	}

	return h.handler(alias)
}

func studioChildCountCriterionHandler(qb *StudioStore, childCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if childCount != nil {
			f.addLeftJoin("studios", "children_count", "children_count.parent_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct children_count.id)", *childCount)

			f.addHaving(clause, args...)
		}
	}
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
	case "child_count":
		sortQuery += getCountSort(studioTable, studioTable, studioParentIDColumn, direction)
	default:
		sortQuery += getSort(sort, direction, "studios")
	}

	// Whatever the sorting, always use name/id as a final sort
	sortQuery += ", COALESCE(studios.name, studios.id) COLLATE NATURAL_CI ASC"
	return sortQuery
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
	return studiosStashIDsTableMgr.get(ctx, studioID)
}

func (qb *StudioStore) GetAliases(ctx context.Context, studioID int) ([]string, error) {
	return studiosAliasesTableMgr.get(ctx, studioID)
}
