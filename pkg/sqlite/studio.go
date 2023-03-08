package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/studio"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

const (
	studioTable          = "studios"
	studioIDColumn       = "studio_id"
	studioAliasesTable   = "studio_aliases"
	studioAliasColumn    = "alias"
	studioParentIDColumn = "parent_id"
	studioNameColumn     = "name"
)

type studioRow struct {
	ID        int                    `db:"id" goqu:"skipinsert"`
	Checksum  string                 `db:"checksum"`
	Name      zero.String            `db:"name"`
	URL       zero.String            `db:"url"`
	ParentID  null.Int               `db:"parent_id,omitempty"`
	CreatedAt models.SQLiteTimestamp `db:"created_at"`
	UpdatedAt models.SQLiteTimestamp `db:"updated_at"`
	Details   zero.String            `db:"details"`
	// expressed as 1-100
	Rating        null.Int `db:"rating"`
	IgnoreAutoTag bool     `db:"ignore_auto_tag"`
}

func (r *studioRow) fromStudio(o models.Studio) {
	r.ID = o.ID
	r.Checksum = md5.FromString(o.Name)
	r.Name = zero.StringFrom(o.Name)
	r.URL = zero.StringFrom(o.URL)
	r.ParentID = intFromPtr(o.ParentID)
	r.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}
	r.Details = zero.StringFrom(o.Details)
	r.Rating = intFromPtr(o.Rating)
	r.IgnoreAutoTag = o.IgnoreAutoTag
}

func (r *studioRow) resolve() *models.Studio {
	ret := &models.Studio{
		ID:        r.ID,
		Checksum:  r.Checksum,
		Name:      r.Name.String,
		URL:       r.URL.String,
		ParentID:  nullIntPtr(r.ParentID),
		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,
		Details:   r.Details.String,
		// expressed as 1-100
		Rating:        nullIntPtr(r.Rating),
		IgnoreAutoTag: r.IgnoreAutoTag,
	}

	return ret
}

type studioRowRecord struct {
	updateRecord
}

func (r *studioRowRecord) fromPartial(o models.StudioPartial) {
	if !o.Name.Null && o.Name.Value != "" {
		r.setString("checksum", models.NewOptionalString(md5.FromString(o.Name.Value)))
	}
	r.setNullString("name", o.Name)
	r.setNullString("url", o.URL)
	r.setNullInt("parent_id", o.ParentID)
	r.setSQLiteTimestamp("updated_at", models.NewOptionalTime(time.Now()))
	r.setNullString("details", o.Details)
	r.setNullInt("rating", o.Rating)
	r.setBool("ignore_auto_tag", o.IgnoreAutoTag)
}

type StudioStore struct {
	repository

	tableMgr *table
}

func NewStudioStore() *StudioStore {
	return &StudioStore{
		repository: repository{
			tableName: studioTable,
			idColumn:  idColumn,
		},
		tableMgr: studioTableMgr,
	}
}

func (qb *StudioStore) Create(ctx context.Context, input models.StudioDBInput) (*int, error) {
	var err error
	var parentID *int
	parentID, err = qb.handleParentStudio(ctx, input)
	if err != nil {
		return nil, err
	}
	if parentID != nil {
		input.StudioCreate.ParentID = parentID
	}

	// Create the main studio
	var r studioRow
	r.fromStudio(*input.StudioCreate)
	time := time.Now()
	r.CreatedAt = models.SQLiteTimestamp{Timestamp: time}
	r.UpdatedAt = models.SQLiteTimestamp{Timestamp: time}

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return nil, err
	}

	// Update image table
	if len(input.StudioCreate.ImageBytes) > 0 {
		if err := qb.UpdateImage(ctx, id, input.StudioCreate.ImageBytes); err != nil {
			return nil, err
		}
	}

	if input.StudioCreate.Aliases.Loaded() {
		if err := studio.EnsureAliasesUnique(ctx, id, input.StudioCreate.Aliases.List(), qb); err != nil {
			return nil, err
		}

		if err := studiosAliasesTableMgr.insertJoins(ctx, id, input.StudioCreate.Aliases.List()); err != nil {
			return nil, err
		}
	}

	if input.StudioCreate.StashIDs.Loaded() {
		if err := studiosStashIDsTableMgr.insertJoins(ctx, id, input.StudioCreate.StashIDs.List()); err != nil {
			return nil, err
		}
	}

	return &id, nil
}

func (qb *StudioStore) UpdatePartial(ctx context.Context, input models.StudioDBInput) (*models.Studio, error) {
	var err error
	var parentID *int
	parentID, err = qb.handleParentStudio(ctx, input)
	if err != nil {
		return nil, err
	} else if parentID != nil {
		input.StudioUpdate.ParentID = models.NewOptionalIntPtr(parentID)
	}

	r := studioRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(*input.StudioUpdate)

	if len(r.Record) > 0 {
		if err := qb.tableMgr.updateByID(ctx, input.StudioUpdate.ID, r.Record); err != nil {
			return nil, err
		}
	}

	// Update image table
	if len(input.StudioUpdate.ImageBytes) > 0 {
		if err := qb.UpdateImage(ctx, input.StudioUpdate.ID, input.StudioUpdate.ImageBytes); err != nil {
			return nil, err
		}
	} else if input.StudioUpdate.ImageIncluded {
		// must be unsetting
		if err := qb.DestroyImage(ctx, input.StudioUpdate.ID); err != nil {
			return nil, err
		}
	}

	if input.StudioUpdate.Aliases != nil {
		if err := studio.EnsureAliasesUnique(ctx, input.StudioUpdate.ID, input.StudioUpdate.Aliases.Values, qb); err != nil {
			return nil, err
		}

		if err := studiosAliasesTableMgr.modifyJoins(ctx, input.StudioUpdate.ID, input.StudioUpdate.Aliases.Values, input.StudioUpdate.Aliases.Mode); err != nil {
			return nil, err
		}
	}

	if input.StudioUpdate.StashIDs != nil {
		if err := studiosStashIDsTableMgr.modifyJoins(ctx, input.StudioUpdate.ID, input.StudioUpdate.StashIDs.StashIDs, input.StudioUpdate.StashIDs.Mode); err != nil {
			return nil, err
		}
	}

	return qb.Find(ctx, input.StudioUpdate.ID)
}

// Returns a studio ID if a new one was created
func (qb *StudioStore) handleParentStudio(ctx context.Context, input models.StudioDBInput) (*int, error) {
	var err error
	var id *int
	var parentDBInput models.StudioDBInput

	if input.ParentCreate != nil {
		parentDBInput.StudioCreate = input.ParentCreate
		id, err = qb.Create(ctx, parentDBInput)
		if err != nil {
			return nil, err
		}
	} else if input.ParentUpdate != nil {
		parentDBInput.StudioUpdate = input.ParentUpdate
		_, err := qb.UpdatePartial(ctx, parentDBInput)
		if err != nil {
			return nil, err
		}
	}
	return id, nil
}

// This is only used by the Import/Export functionality, which already handles parent/child studios
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
	// TODO - set null on foreign key in scraped items
	// remove studio from scraped items
	_, err := qb.tx.Exec(ctx, "UPDATE scraped_items SET studio_id = null WHERE studio_id = ?", id)
	if err != nil {
		return err
	}

	return qb.destroyExisting(ctx, []int{id})
}

func (qb *StudioStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *StudioStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *StudioStore) Find(ctx context.Context, id int) (*models.Studio, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting studio by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *StudioStore) FindMany(ctx context.Context, ids []int) ([]*models.Studio, error) {
	tableMgr := studioTableMgr
	q := goqu.Select("*").From(tableMgr.table).Where(tableMgr.byIDInts(ids...))
	unsorted, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.Studio, len(ids))

	for _, s := range unsorted {
		i := intslice.IntIndex(ids, s.ID)
		ret[i] = s
	}

	for i := range ret {
		if ret[i] == nil {
			return nil, fmt.Errorf("studio with id %d not found", ids[i])
		}
	}

	return ret, nil
}

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
	table := qb.table()

	sq := dialect.From(table).
		Select(table.Col(idColumn)).Where(
		table.Col(studioParentIDColumn).Eq(id),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting child studios for studio %d: %w", id, err)
	}

	return ret, nil
}

func (qb *StudioStore) FindBySceneID(ctx context.Context, sceneID int) (*models.Studio, error) {
	table := qb.table()
	sceneTable := sceneTableMgr.table

	sq := dialect.From(table).
		InnerJoin(sceneTable,
			goqu.On(table.Col(idColumn).Eq(sceneTable.Col(studioIDColumn))),
		).
		Select(table.Col(idColumn)).Where(
		sceneTable.Col(idColumn).Eq(sceneID),
	).Limit(1)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting studio for scene %d: %w", sceneID, err)
	} else if len(ret) < 1 {
		return nil, err
	}

	return ret[0], nil
}

func (qb *StudioStore) FindByName(ctx context.Context, name string, nocase bool) (*models.Studio, error) {
	clause := "name "
	if nocase {
		clause += "COLLATE NOCASE "
	}
	clause += "= ?"

	sq := qb.selectDataset().Prepared(true).Where(
		goqu.L(clause, name),
	).Limit(1)
	ret, err := qb.getMany(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting studio by name: %s: %w", name, err)
	} else if len(ret) < 1 {
		return nil, err
	}

	return ret[0], nil
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
	sq := dialect.From(table).Select(table.Col(idColumn))
	// TODO - disabled alias matching until we get finer control over it
	// .LeftJoin(
	// 	studiosAliasesJoinTable,
	// 	goqu.On(studiosAliasesJoinTable.Col(studioIDColumn).Eq(table.Col(idColumn))),
	// )

	var whereClauses []exp.Expression

	for _, w := range words {
		whereClauses = append(whereClauses, table.Col(studioNameColumn).Like(w+"%"))
		// TODO - see above
		// whereClauses = append(whereClauses, studiosAliasesJoinTable.Col("alias").Like(w+"%"))
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
				f.addLeftJoin("studios_image", "", "studios_image.studio_id = studios.id")
				f.addWhere("studios_image.studio_id IS NULL")
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

	switch sort {
	case "scenes_count":
		return getCountSort(studioTable, sceneTable, studioIDColumn, direction)
	case "images_count":
		return getCountSort(studioTable, imageTable, studioIDColumn, direction)
	case "galleries_count":
		return getCountSort(studioTable, galleryTable, studioIDColumn, direction)
	default:
		return getSort(sort, direction, "studios")
	}
}

func (qb *StudioStore) imageRepository() *imageRepository {
	return &imageRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: "studios_image",
			idColumn:  studioIDColumn,
		},
		imageColumn: "image",
	}
}

func (qb *StudioStore) GetImage(ctx context.Context, studioID int) ([]byte, error) {
	return qb.imageRepository().get(ctx, studioID)
}

func (qb *StudioStore) HasImage(ctx context.Context, studioID int) (bool, error) {
	return qb.imageRepository().exists(ctx, studioID)
}

func (qb *StudioStore) UpdateImage(ctx context.Context, studioID int, image []byte) error {
	return qb.imageRepository().replace(ctx, studioID, image)
}

func (qb *StudioStore) DestroyImage(ctx context.Context, studioID int) error {
	return qb.imageRepository().destroy(ctx, []int{studioID})
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
