package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

const (
	performerTable         = "performers"
	performerIDColumn      = "performer_id"
	performersAliasesTable = "performer_aliases"
	performerAliasColumn   = "alias"
	performersTagsTable    = "performers_tags"

	performerImageBlobColumn = "image_blob"
)

type performerRow struct {
	ID            int         `db:"id" goqu:"skipinsert"`
	Name          null.String `db:"name"` // TODO: make schema non-nullable
	Disambigation zero.String `db:"disambiguation"`
	Gender        zero.String `db:"gender"`
	URL           zero.String `db:"url"`
	Twitter       zero.String `db:"twitter"`
	Instagram     zero.String `db:"instagram"`
	Birthdate     NullDate    `db:"birthdate"`
	Ethnicity     zero.String `db:"ethnicity"`
	Country       zero.String `db:"country"`
	EyeColor      zero.String `db:"eye_color"`
	Height        null.Int    `db:"height"`
	Measurements  zero.String `db:"measurements"`
	FakeTits      zero.String `db:"fake_tits"`
	PenisLength   null.Float  `db:"penis_length"`
	Circumcised   zero.String `db:"circumcised"`
	CareerLength  zero.String `db:"career_length"`
	Tattoos       zero.String `db:"tattoos"`
	Piercings     zero.String `db:"piercings"`
	Favorite      bool        `db:"favorite"`
	CreatedAt     Timestamp   `db:"created_at"`
	UpdatedAt     Timestamp   `db:"updated_at"`
	// expressed as 1-100
	Rating        null.Int    `db:"rating"`
	Details       zero.String `db:"details"`
	DeathDate     NullDate    `db:"death_date"`
	HairColor     zero.String `db:"hair_color"`
	Weight        null.Int    `db:"weight"`
	IgnoreAutoTag bool        `db:"ignore_auto_tag"`

	// not used in resolution or updates
	ImageBlob zero.String `db:"image_blob"`
}

func (r *performerRow) fromPerformer(o models.Performer) {
	r.ID = o.ID
	r.Name = null.StringFrom(o.Name)
	r.Disambigation = zero.StringFrom(o.Disambiguation)
	if o.Gender != nil && o.Gender.IsValid() {
		r.Gender = zero.StringFrom(o.Gender.String())
	}
	r.URL = zero.StringFrom(o.URL)
	r.Twitter = zero.StringFrom(o.Twitter)
	r.Instagram = zero.StringFrom(o.Instagram)
	r.Birthdate = NullDateFromDatePtr(o.Birthdate)
	r.Ethnicity = zero.StringFrom(o.Ethnicity)
	r.Country = zero.StringFrom(o.Country)
	r.EyeColor = zero.StringFrom(o.EyeColor)
	r.Height = intFromPtr(o.Height)
	r.Measurements = zero.StringFrom(o.Measurements)
	r.FakeTits = zero.StringFrom(o.FakeTits)
	r.PenisLength = null.FloatFromPtr(o.PenisLength)
	if o.Circumcised != nil && o.Circumcised.IsValid() {
		r.Circumcised = zero.StringFrom(o.Circumcised.String())
	}
	r.CareerLength = zero.StringFrom(o.CareerLength)
	r.Tattoos = zero.StringFrom(o.Tattoos)
	r.Piercings = zero.StringFrom(o.Piercings)
	r.Favorite = o.Favorite
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
	r.Rating = intFromPtr(o.Rating)
	r.Details = zero.StringFrom(o.Details)
	r.DeathDate = NullDateFromDatePtr(o.DeathDate)
	r.HairColor = zero.StringFrom(o.HairColor)
	r.Weight = intFromPtr(o.Weight)
	r.IgnoreAutoTag = o.IgnoreAutoTag
}

func (r *performerRow) resolve() *models.Performer {
	ret := &models.Performer{
		ID:             r.ID,
		Name:           r.Name.String,
		Disambiguation: r.Disambigation.String,
		URL:            r.URL.String,
		Twitter:        r.Twitter.String,
		Instagram:      r.Instagram.String,
		Birthdate:      r.Birthdate.DatePtr(),
		Ethnicity:      r.Ethnicity.String,
		Country:        r.Country.String,
		EyeColor:       r.EyeColor.String,
		Height:         nullIntPtr(r.Height),
		Measurements:   r.Measurements.String,
		FakeTits:       r.FakeTits.String,
		PenisLength:    nullFloatPtr(r.PenisLength),
		CareerLength:   r.CareerLength.String,
		Tattoos:        r.Tattoos.String,
		Piercings:      r.Piercings.String,
		Favorite:       r.Favorite,
		CreatedAt:      r.CreatedAt.Timestamp,
		UpdatedAt:      r.UpdatedAt.Timestamp,
		// expressed as 1-100
		Rating:        nullIntPtr(r.Rating),
		Details:       r.Details.String,
		DeathDate:     r.DeathDate.DatePtr(),
		HairColor:     r.HairColor.String,
		Weight:        nullIntPtr(r.Weight),
		IgnoreAutoTag: r.IgnoreAutoTag,
	}

	if r.Gender.ValueOrZero() != "" {
		v := models.GenderEnum(r.Gender.String)
		ret.Gender = &v
	}

	if r.Circumcised.ValueOrZero() != "" {
		v := models.CircumisedEnum(r.Circumcised.String)
		ret.Circumcised = &v
	}

	return ret
}

type performerRowRecord struct {
	updateRecord
}

func (r *performerRowRecord) fromPartial(o models.PerformerPartial) {
	r.setString("name", o.Name)
	r.setNullString("disambiguation", o.Disambiguation)
	r.setNullString("gender", o.Gender)
	r.setNullString("url", o.URL)
	r.setNullString("twitter", o.Twitter)
	r.setNullString("instagram", o.Instagram)
	r.setNullDate("birthdate", o.Birthdate)
	r.setNullString("ethnicity", o.Ethnicity)
	r.setNullString("country", o.Country)
	r.setNullString("eye_color", o.EyeColor)
	r.setNullInt("height", o.Height)
	r.setNullString("measurements", o.Measurements)
	r.setNullString("fake_tits", o.FakeTits)
	r.setNullFloat64("penis_length", o.PenisLength)
	r.setNullString("circumcised", o.Circumcised)
	r.setNullString("career_length", o.CareerLength)
	r.setNullString("tattoos", o.Tattoos)
	r.setNullString("piercings", o.Piercings)
	r.setBool("favorite", o.Favorite)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
	r.setNullInt("rating", o.Rating)
	r.setNullString("details", o.Details)
	r.setNullDate("death_date", o.DeathDate)
	r.setNullString("hair_color", o.HairColor)
	r.setNullInt("weight", o.Weight)
	r.setBool("ignore_auto_tag", o.IgnoreAutoTag)
}

type PerformerStore struct {
	repository
	blobJoinQueryBuilder

	tableMgr *table
}

func NewPerformerStore(blobStore *BlobStore) *PerformerStore {
	return &PerformerStore{
		repository: repository{
			tableName: performerTable,
			idColumn:  idColumn,
		},
		blobJoinQueryBuilder: blobJoinQueryBuilder{
			blobStore: blobStore,
			joinTable: performerTable,
		},
		tableMgr: performerTableMgr,
	}
}

func (qb *PerformerStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *PerformerStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *PerformerStore) Create(ctx context.Context, newObject *models.Performer) error {
	var r performerRow
	r.fromPerformer(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if newObject.Aliases.Loaded() {
		if err := performersAliasesTableMgr.insertJoins(ctx, id, newObject.Aliases.List()); err != nil {
			return err
		}
	}

	if newObject.TagIDs.Loaded() {
		if err := performersTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if newObject.StashIDs.Loaded() {
		if err := performersStashIDsTableMgr.insertJoins(ctx, id, newObject.StashIDs.List()); err != nil {
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

func (qb *PerformerStore) UpdatePartial(ctx context.Context, id int, partial models.PerformerPartial) (*models.Performer, error) {
	r := performerRowRecord{
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

	if partial.Aliases != nil {
		if err := performersAliasesTableMgr.modifyJoins(ctx, id, partial.Aliases.Values, partial.Aliases.Mode); err != nil {
			return nil, err
		}
	}

	if partial.TagIDs != nil {
		if err := performersTagsTableMgr.modifyJoins(ctx, id, partial.TagIDs.IDs, partial.TagIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.StashIDs != nil {
		if err := performersStashIDsTableMgr.modifyJoins(ctx, id, partial.StashIDs.StashIDs, partial.StashIDs.Mode); err != nil {
			return nil, err
		}
	}

	return qb.find(ctx, id)
}

func (qb *PerformerStore) Update(ctx context.Context, updatedObject *models.Performer) error {
	var r performerRow
	r.fromPerformer(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.Aliases.Loaded() {
		if err := performersAliasesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.Aliases.List()); err != nil {
			return err
		}
	}

	if updatedObject.TagIDs.Loaded() {
		if err := performersTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.StashIDs.Loaded() {
		if err := performersStashIDsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.StashIDs.List()); err != nil {
			return err
		}
	}

	return nil
}

func (qb *PerformerStore) Destroy(ctx context.Context, id int) error {
	// must handle image checksums manually
	if err := qb.destroyImage(ctx, id); err != nil {
		return err
	}

	return qb.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *PerformerStore) Find(ctx context.Context, id int) (*models.Performer, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *PerformerStore) FindMany(ctx context.Context, ids []int) ([]*models.Performer, error) {
	tableMgr := performerTableMgr
	ret := make([]*models.Performer, len(ids))

	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := goqu.Select("*").From(tableMgr.table).Where(tableMgr.byIDInts(batch...))
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
			return nil, fmt.Errorf("performer with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *PerformerStore) find(ctx context.Context, id int) (*models.Performer, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *PerformerStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Performer, error) {
	table := qb.table()

	q := qb.selectDataset().Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

// returns nil, sql.ErrNoRows if not found
func (qb *PerformerStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Performer, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *PerformerStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Performer, error) {
	const single = false
	var ret []*models.Performer
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f performerRow
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

func (qb *PerformerStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.Performer, error) {
	sq := dialect.From(scenesPerformersJoinTable).Select(scenesPerformersJoinTable.Col(performerIDColumn)).Where(
		scenesPerformersJoinTable.Col(sceneIDColumn).Eq(sceneID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting performers for scene %d: %w", sceneID, err)
	}

	return ret, nil
}

func (qb *PerformerStore) FindByImageID(ctx context.Context, imageID int) ([]*models.Performer, error) {
	sq := dialect.From(performersImagesJoinTable).Select(performersImagesJoinTable.Col(performerIDColumn)).Where(
		performersImagesJoinTable.Col(imageIDColumn).Eq(imageID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting performers for image %d: %w", imageID, err)
	}

	return ret, nil
}

func (qb *PerformerStore) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Performer, error) {
	sq := dialect.From(performersGalleriesJoinTable).Select(performersGalleriesJoinTable.Col(performerIDColumn)).Where(
		performersGalleriesJoinTable.Col(galleryIDColumn).Eq(galleryID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting performers for gallery %d: %w", galleryID, err)
	}

	return ret, nil
}

func (qb *PerformerStore) FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Performer, error) {
	clause := "name "
	if nocase {
		clause += "COLLATE NOCASE "
	}
	clause += "IN " + getInBinding(len(names))

	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}

	sq := qb.selectDataset().Prepared(true).Where(
		goqu.L(clause, args...),
	)
	ret, err := qb.getMany(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting performers by names: %w", err)
	}

	return ret, nil
}

func (qb *PerformerStore) CountByTagID(ctx context.Context, tagID int) (int, error) {
	joinTable := performersTagsJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(tagIDColumn).Eq(tagID))
	return count(ctx, q)
}

func (qb *PerformerStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *PerformerStore) All(ctx context.Context) ([]*models.Performer, error) {
	table := qb.table()
	return qb.getMany(ctx, qb.selectDataset().Order(table.Col("name").Asc()))
}

func (qb *PerformerStore) QueryForAutoTag(ctx context.Context, words []string) ([]*models.Performer, error) {
	// TODO - Query needs to be changed to support queries of this type, and
	// this method should be removed
	table := qb.table()
	sq := dialect.From(table).Select(table.Col(idColumn))
	// TODO - disabled alias matching until we get finer control over it
	// .LeftJoin(
	// 	performersAliasesJoinTable,
	// 	goqu.On(performersAliasesJoinTable.Col(performerIDColumn).Eq(table.Col(idColumn))),
	// )

	var whereClauses []exp.Expression

	for _, w := range words {
		whereClauses = append(whereClauses, table.Col("name").Like(w+"%"))
		// TODO - see above
		// whereClauses = append(whereClauses, performersAliasesJoinTable.Col("alias").Like(w+"%"))
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

func (qb *PerformerStore) validateFilter(filter *models.PerformerFilterType) error {
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

	// if legacy height filter used, ensure only supported modifiers are used
	if filter.Height != nil {
		// treat as an int filter
		intCrit := &models.IntCriterionInput{
			Modifier: filter.Height.Modifier,
		}
		if !intCrit.ValidModifier() {
			return fmt.Errorf("invalid height modifier: %s", filter.Height.Modifier)
		}

		// ensure value is a valid number
		if _, err := strconv.Atoi(filter.Height.Value); err != nil {
			return fmt.Errorf("invalid height value: %s", filter.Height.Value)
		}
	}

	return nil
}

func (qb *PerformerStore) makeFilter(ctx context.Context, filter *models.PerformerFilterType) *filterBuilder {
	query := &filterBuilder{}

	if filter.And != nil {
		query.and(qb.makeFilter(ctx, filter.And))
	}
	if filter.Or != nil {
		query.or(qb.makeFilter(ctx, filter.Or))
	}
	if filter.Not != nil {
		query.not(qb.makeFilter(ctx, filter.Not))
	}

	const tableName = performerTable
	query.handleCriterion(ctx, stringCriterionHandler(filter.Name, tableName+".name"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.Disambiguation, tableName+".disambiguation"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.Details, tableName+".details"))

	query.handleCriterion(ctx, boolCriterionHandler(filter.FilterFavorites, tableName+".favorite", nil))
	query.handleCriterion(ctx, boolCriterionHandler(filter.IgnoreAutoTag, tableName+".ignore_auto_tag", nil))

	query.handleCriterion(ctx, yearFilterCriterionHandler(filter.BirthYear, tableName+".birthdate"))
	query.handleCriterion(ctx, yearFilterCriterionHandler(filter.DeathYear, tableName+".death_date"))

	query.handleCriterion(ctx, performerAgeFilterCriterionHandler(filter.Age))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if gender := filter.Gender; gender != nil {
			genderCopy := *gender
			if genderCopy.Value.IsValid() && len(genderCopy.ValueList) == 0 {
				genderCopy.ValueList = []models.GenderEnum{genderCopy.Value}
			}

			v := utils.StringerSliceToStringSlice(genderCopy.ValueList)
			enumCriterionHandler(genderCopy.Modifier, v, tableName+".gender")(ctx, f)
		}
	}))

	query.handleCriterion(ctx, performerIsMissingCriterionHandler(qb, filter.IsMissing))
	query.handleCriterion(ctx, stringCriterionHandler(filter.Ethnicity, tableName+".ethnicity"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.Country, tableName+".country"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.EyeColor, tableName+".eye_color"))

	// special handler for legacy height filter
	heightCmCrit := filter.HeightCm
	if heightCmCrit == nil && filter.Height != nil {
		heightCm, _ := strconv.Atoi(filter.Height.Value) // already validated
		heightCmCrit = &models.IntCriterionInput{
			Value:    heightCm,
			Modifier: filter.Height.Modifier,
		}
	}

	query.handleCriterion(ctx, intCriterionHandler(heightCmCrit, tableName+".height", nil))

	query.handleCriterion(ctx, stringCriterionHandler(filter.Measurements, tableName+".measurements"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.FakeTits, tableName+".fake_tits"))
	query.handleCriterion(ctx, floatCriterionHandler(filter.PenisLength, tableName+".penis_length", nil))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if circumcised := filter.Circumcised; circumcised != nil {
			v := utils.StringerSliceToStringSlice(circumcised.Value)
			enumCriterionHandler(circumcised.Modifier, v, tableName+".circumcised")(ctx, f)
		}
	}))

	query.handleCriterion(ctx, stringCriterionHandler(filter.CareerLength, tableName+".career_length"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.Tattoos, tableName+".tattoos"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.Piercings, tableName+".piercings"))
	query.handleCriterion(ctx, intCriterionHandler(filter.Rating100, tableName+".rating", nil))
	query.handleCriterion(ctx, stringCriterionHandler(filter.HairColor, tableName+".hair_color"))
	query.handleCriterion(ctx, stringCriterionHandler(filter.URL, tableName+".url"))
	query.handleCriterion(ctx, intCriterionHandler(filter.Weight, tableName+".weight", nil))
	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if filter.StashID != nil {
			qb.stashIDRepository().join(f, "performer_stash_ids", "performers.id")
			stringCriterionHandler(filter.StashID, "performer_stash_ids.stash_id")(ctx, f)
		}
	}))
	query.handleCriterion(ctx, &stashIDCriterionHandler{
		c:                 filter.StashIDEndpoint,
		stashIDRepository: qb.stashIDRepository(),
		stashIDTableAs:    "performer_stash_ids",
		parentIDCol:       "performers.id",
	})

	query.handleCriterion(ctx, performerAliasCriterionHandler(qb, filter.Aliases))

	query.handleCriterion(ctx, performerTagsCriterionHandler(qb, filter.Tags))

	query.handleCriterion(ctx, performerStudiosCriterionHandler(qb, filter.Studios))

	query.handleCriterion(ctx, performerAppearsWithCriterionHandler(qb, filter.Performers))

	query.handleCriterion(ctx, performerTagCountCriterionHandler(qb, filter.TagCount))
	query.handleCriterion(ctx, performerSceneCountCriterionHandler(qb, filter.SceneCount))
	query.handleCriterion(ctx, performerImageCountCriterionHandler(qb, filter.ImageCount))
	query.handleCriterion(ctx, performerGalleryCountCriterionHandler(qb, filter.GalleryCount))
	query.handleCriterion(ctx, performerPlayCounterCriterionHandler(qb, filter.PlayCount))
	query.handleCriterion(ctx, performerOCounterCriterionHandler(qb, filter.OCounter))
	query.handleCriterion(ctx, dateCriterionHandler(filter.Birthdate, tableName+".birthdate"))
	query.handleCriterion(ctx, dateCriterionHandler(filter.DeathDate, tableName+".death_date"))
	query.handleCriterion(ctx, timestampCriterionHandler(filter.CreatedAt, tableName+".created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(filter.UpdatedAt, tableName+".updated_at"))

	return query
}

func (qb *PerformerStore) makeQuery(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if performerFilter == nil {
		performerFilter = &models.PerformerFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, performerTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.join(performersAliasesTable, "", "performer_aliases.performer_id = performers.id")
		searchColumns := []string{"performers.name", "performer_aliases.alias"}
		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(performerFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, performerFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	var err error
	query.sortAndPagination, err = qb.getPerformerSort(findFilter)
	if err != nil {
		return nil, err
	}
	query.sortAndPagination += getPagination(findFilter)

	return &query, nil
}

func (qb *PerformerStore) Query(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, int, error) {
	query, err := qb.makeQuery(ctx, performerFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	performers, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return performers, countResult, nil
}

func (qb *PerformerStore) QueryCount(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, performerFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

// TODO - we need to provide a whitelist of possible values
func performerIsMissingCriterionHandler(qb *PerformerStore, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "scenes": // Deprecated: use `scene_count == 0` filter instead
				f.addLeftJoin(performersScenesTable, "scenes_join", "scenes_join.performer_id = performers.id")
				f.addWhere("scenes_join.scene_id IS NULL")
			case "image":
				f.addWhere("performers.image_blob IS NULL")
			case "stash_id":
				performersStashIDsTableMgr.join(f, "performer_stash_ids", "performers.id")
				f.addWhere("performer_stash_ids.performer_id IS NULL")
			case "aliases":
				performersAliasesTableMgr.join(f, "", "performers.id")
				f.addWhere("performer_aliases.alias IS NULL")
			default:
				f.addWhere("(performers." + *isMissing + " IS NULL OR TRIM(performers." + *isMissing + ") = '')")
			}
		}
	}
}

func yearFilterCriterionHandler(year *models.IntCriterionInput, col string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if year != nil && year.Modifier.IsValid() {
			clause, args := getIntCriterionWhereClause("cast(strftime('%Y', "+col+") as int)", *year)
			f.addWhere(clause, args...)
		}
	}
}

func performerAgeFilterCriterionHandler(age *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if age != nil && age.Modifier.IsValid() {
			clause, args := getIntCriterionWhereClause(
				"cast(IFNULL(strftime('%Y.%m%d', performers.death_date), strftime('%Y.%m%d', 'now')) - strftime('%Y.%m%d', performers.birthdate) as int)",
				*age,
			)
			f.addWhere(clause, args...)
		}
	}
}

func performerAliasCriterionHandler(qb *PerformerStore, alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    performersAliasesTable,
		stringColumn: performerAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			performersAliasesTableMgr.join(f, "", "performers.id")
		},
	}

	return h.handler(alias)
}

func performerTagsCriterionHandler(qb *PerformerStore, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: performerTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "image_tag",
		joinTable:      performersTagsTable,
		primaryFK:      performerIDColumn,
	}

	return h.handler(tags)
}

func performerTagCountCriterionHandler(qb *PerformerStore, count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersTagsTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func performerSceneCountCriterionHandler(qb *PerformerStore, count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersScenesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func performerImageCountCriterionHandler(qb *PerformerStore, count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersImagesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

func performerGalleryCountCriterionHandler(qb *PerformerStore, count *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: performerTable,
		joinTable:    performersGalleriesTable,
		primaryFK:    performerIDColumn,
	}

	return h.handler(count)
}

// used for sorting and filtering on performer o-count
var selectPerformerOCountSQL = utils.StrFormat(
	"SELECT SUM(o_counter) "+
		"FROM ("+
		"SELECT SUM(o_counter) as o_counter from {performers_images} s "+
		"LEFT JOIN {images} ON {images}.id = s.{images_id} "+
		"WHERE s.{performer_id} = {performers}.id "+
		"UNION ALL "+
		"SELECT COUNT({scenes_o_dates}.{o_date}) as o_counter from {performers_scenes} s "+
		"LEFT JOIN {scenes} ON {scenes}.id = s.{scene_id} "+
		"LEFT JOIN {scenes_o_dates} ON {scenes_o_dates}.{scene_id} = {scenes}.id "+
		"WHERE s.{performer_id} = {performers}.id "+
		")",
	map[string]interface{}{
		"performers_images": performersImagesTable,
		"images":            imageTable,
		"performer_id":      performerIDColumn,
		"images_id":         imageIDColumn,
		"performers":        performerTable,
		"performers_scenes": performersScenesTable,
		"scenes":            sceneTable,
		"scene_id":          sceneIDColumn,
		"scenes_o_dates":    scenesODatesTable,
		"o_date":            sceneODateColumn,
	},
)

// used for sorting and filtering play count on performer view count
var selectPerformerPlayCountSQL = utils.StrFormat(
	"SELECT COUNT(DISTINCT {view_date}) FROM ("+
		"SELECT {view_date} FROM {performers_scenes} s "+
		"LEFT JOIN {scenes} ON {scenes}.id = s.{scene_id} "+
		"LEFT JOIN {scenes_view_dates} ON {scenes_view_dates}.{scene_id} = {scenes}.id "+
		"WHERE s.{performer_id} = {performers}.id"+
		")",
	map[string]interface{}{
		"performer_id":      performerIDColumn,
		"performers":        performerTable,
		"performers_scenes": performersScenesTable,
		"scenes":            sceneTable,
		"scene_id":          sceneIDColumn,
		"scenes_view_dates": scenesViewDatesTable,
		"view_date":         sceneViewDateColumn,
	},
)

// used for sorting on performer last o_date
var selectPerformerLastOAtSQL = utils.StrFormat(
	"SELECT MAX(o_date) FROM ("+
		"SELECT {o_date} FROM {performers_scenes} s "+
		"LEFT JOIN {scenes} ON {scenes}.id = s.{scene_id} "+
		"LEFT JOIN {scenes_o_dates} ON {scenes_o_dates}.{scene_id} = {scenes}.id "+
		"WHERE s.{performer_id} = {performers}.id"+
		")",
	map[string]interface{}{
		"performer_id":      performerIDColumn,
		"performers":        performerTable,
		"performers_scenes": performersScenesTable,
		"scenes":            sceneTable,
		"scene_id":          sceneIDColumn,
		"scenes_o_dates":    scenesODatesTable,
		"o_date":            sceneODateColumn,
	},
)

// used for sorting on performer last view_date
var selectPerformerLastPlayedAtSQL = utils.StrFormat(
	"SELECT MAX(view_date) FROM ("+
		"SELECT {view_date} FROM {performers_scenes} s "+
		"LEFT JOIN {scenes} ON {scenes}.id = s.{scene_id} "+
		"LEFT JOIN {scenes_view_dates} ON {scenes_view_dates}.{scene_id} = {scenes}.id "+
		"WHERE s.{performer_id} = {performers}.id"+
		")",
	map[string]interface{}{
		"performer_id":      performerIDColumn,
		"performers":        performerTable,
		"performers_scenes": performersScenesTable,
		"scenes":            sceneTable,
		"scene_id":          sceneIDColumn,
		"scenes_view_dates": scenesViewDatesTable,
		"view_date":         sceneViewDateColumn,
	},
)

func performerOCounterCriterionHandler(qb *PerformerStore, count *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if count == nil {
			return
		}

		lhs := "(" + selectPerformerOCountSQL + ")"
		clause, args := getIntCriterionWhereClause(lhs, *count)

		f.addWhere(clause, args...)
	}
}

func performerPlayCounterCriterionHandler(qb *PerformerStore, count *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if count == nil {
			return
		}

		lhs := "(" + selectPerformerPlayCountSQL + ")"
		clause, args := getIntCriterionWhereClause(lhs, *count)

		f.addWhere(clause, args...)
	}
}

func performerStudiosCriterionHandler(qb *PerformerStore, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if studios != nil {
			formatMaps := []utils.StrFormatMap{
				{
					"primaryTable": sceneTable,
					"joinTable":    performersScenesTable,
					"primaryFK":    sceneIDColumn,
				},
				{
					"primaryTable": imageTable,
					"joinTable":    performersImagesTable,
					"primaryFK":    imageIDColumn,
				},
				{
					"primaryTable": galleryTable,
					"joinTable":    performersGalleriesTable,
					"primaryFK":    galleryIDColumn,
				},
			}

			if studios.Modifier == models.CriterionModifierIsNull || studios.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if studios.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				var conditions []string
				for _, c := range formatMaps {
					f.addLeftJoin(c["joinTable"].(string), "", fmt.Sprintf("%s.performer_id = performers.id", c["joinTable"]))
					f.addLeftJoin(c["primaryTable"].(string), "", fmt.Sprintf("%s.%s = %s.id", c["joinTable"], c["primaryFK"], c["primaryTable"]))

					conditions = append(conditions, fmt.Sprintf("%s.studio_id IS NULL", c["primaryTable"]))
				}

				f.addWhere(fmt.Sprintf("%s (%s)", notClause, strings.Join(conditions, " AND ")))
				return
			}

			if len(studios.Value) == 0 {
				return
			}

			var clauseCondition string

			switch studios.Modifier {
			case models.CriterionModifierIncludes:
				// return performers who appear in scenes/images/galleries with any of the given studios
				clauseCondition = "NOT"
			case models.CriterionModifierExcludes:
				// exclude performers who appear in scenes/images/galleries with any of the given studios
				clauseCondition = ""
			default:
				return
			}

			const derivedPerformerStudioTable = "performer_studio"
			valuesClause, err := getHierarchicalValues(ctx, qb.tx, studios.Value, studioTable, "", "parent_id", "child_id", studios.Depth)
			if err != nil {
				f.setError(err)
				return
			}
			f.addWith("studio(root_id, item_id) AS (" + valuesClause + ")")

			templStr := `SELECT performer_id FROM {primaryTable}
	INNER JOIN {joinTable} ON {primaryTable}.id = {joinTable}.{primaryFK}
	INNER JOIN studio ON {primaryTable}.studio_id = studio.item_id`

			var unions []string
			for _, c := range formatMaps {
				unions = append(unions, utils.StrFormat(templStr, c))
			}

			f.addWith(fmt.Sprintf("%s AS (%s)", derivedPerformerStudioTable, strings.Join(unions, " UNION ")))

			f.addLeftJoin(derivedPerformerStudioTable, "", fmt.Sprintf("performers.id = %s.performer_id", derivedPerformerStudioTable))
			f.addWhere(fmt.Sprintf("%s.performer_id IS %s NULL", derivedPerformerStudioTable, clauseCondition))
		}
	}
}

func performerAppearsWithCriterionHandler(qb *PerformerStore, performers *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performers != nil {
			formatMaps := []utils.StrFormatMap{
				{
					"primaryTable": performersScenesTable,
					"joinTable":    performersScenesTable,
					"primaryFK":    sceneIDColumn,
				},
				{
					"primaryTable": performersImagesTable,
					"joinTable":    performersImagesTable,
					"primaryFK":    imageIDColumn,
				},
				{
					"primaryTable": performersGalleriesTable,
					"joinTable":    performersGalleriesTable,
					"primaryFK":    galleryIDColumn,
				},
			}

			if len(performers.Value) == '0' {
				return
			}

			const derivedPerformerPerformersTable = "performer_performers"

			valuesClause := strings.Join(performers.Value, "),(")

			f.addWith("performer(id) AS (VALUES(" + valuesClause + "))")

			templStr := `SELECT {primaryTable}2.performer_id FROM {primaryTable}
			INNER JOIN {primaryTable} AS {primaryTable}2 ON {primaryTable}.{primaryFK} = {primaryTable}2.{primaryFK}
			INNER JOIN performer ON {primaryTable}.performer_id = performer.id
			WHERE {primaryTable}2.performer_id != performer.id`

			if performers.Modifier == models.CriterionModifierIncludesAll && len(performers.Value) > 1 {
				templStr += `
							GROUP BY {primaryTable}2.performer_id
							HAVING(count(distinct {primaryTable}.performer_id) IS ` + strconv.Itoa(len(performers.Value)) + `)`
			}

			var unions []string
			for _, c := range formatMaps {
				unions = append(unions, utils.StrFormat(templStr, c))
			}

			f.addWith(fmt.Sprintf("%s AS (%s)", derivedPerformerPerformersTable, strings.Join(unions, " UNION ")))

			f.addInnerJoin(derivedPerformerPerformersTable, "", fmt.Sprintf("performers.id = %s.performer_id", derivedPerformerPerformersTable))
		}
	}
}

func (qb *PerformerStore) sortByOCounter(direction string) string {
	// need to sum the o_counter from scenes and images
	return " ORDER BY (" + selectPerformerOCountSQL + ") " + direction
}

func (qb *PerformerStore) sortByPlayCount(direction string) string {
	// need to sum the o_counter from scenes and images
	return " ORDER BY (" + selectPerformerPlayCountSQL + ") " + direction
}

func (qb *PerformerStore) sortByLastOAt(direction string) string {
	// need to get the o_dates from scenes
	return " ORDER BY (" + selectPerformerLastOAtSQL + ") " + direction
}

func (qb *PerformerStore) sortByLastPlayedAt(direction string) string {
	// need to get the view_dates from scenes
	return " ORDER BY (" + selectPerformerLastPlayedAtSQL + ") " + direction
}

var performerSortOptions = sortOptions{
	"birthdate",
	"created_at",
	"galleries_count",
	"height",
	"id",
	"images_count",
	"last_o_at",
	"last_played_at",
	"name",
	"o_counter",
	"penis_length",
	"play_count",
	"random",
	"rating",
	"scenes_count",
	"tag_count",
	"updated_at",
}

func (qb *PerformerStore) getPerformerSort(findFilter *models.FindFilterType) (string, error) {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := performerSortOptions.validateSort(sort); err != nil {
		return "", err
	}

	sortQuery := ""
	switch sort {
	case "tag_count":
		sortQuery += getCountSort(performerTable, performersTagsTable, performerIDColumn, direction)
	case "scenes_count":
		sortQuery += getCountSort(performerTable, performersScenesTable, performerIDColumn, direction)
	case "images_count":
		sortQuery += getCountSort(performerTable, performersImagesTable, performerIDColumn, direction)
	case "galleries_count":
		sortQuery += getCountSort(performerTable, performersGalleriesTable, performerIDColumn, direction)
	case "play_count":
		sortQuery += qb.sortByPlayCount(direction)
	case "o_counter":
		sortQuery += qb.sortByOCounter(direction)
	case "last_played_at":
		sortQuery += qb.sortByLastPlayedAt(direction)
	case "last_o_at":
		sortQuery += qb.sortByLastOAt(direction)
	default:
		sortQuery += getSort(sort, direction, "performers")
	}

	// Whatever the sorting, always use name/id as a final sort
	sortQuery += ", COALESCE(performers.name, performers.id) COLLATE NATURAL_CI ASC"
	return sortQuery, nil
}

func (qb *PerformerStore) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersTagsTable,
			idColumn:  performerIDColumn,
		},
		fkColumn:     tagIDColumn,
		foreignTable: tagTable,
		orderBy:      "tags.name ASC",
	}
}

func (qb *PerformerStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return qb.tagsRepository().getIDs(ctx, id)
}

func (qb *PerformerStore) GetImage(ctx context.Context, performerID int) ([]byte, error) {
	return qb.blobJoinQueryBuilder.GetImage(ctx, performerID, performerImageBlobColumn)
}

func (qb *PerformerStore) HasImage(ctx context.Context, performerID int) (bool, error) {
	return qb.blobJoinQueryBuilder.HasImage(ctx, performerID, performerImageBlobColumn)
}

func (qb *PerformerStore) UpdateImage(ctx context.Context, performerID int, image []byte) error {
	return qb.blobJoinQueryBuilder.UpdateImage(ctx, performerID, performerImageBlobColumn, image)
}

func (qb *PerformerStore) destroyImage(ctx context.Context, performerID int) error {
	return qb.blobJoinQueryBuilder.DestroyImage(ctx, performerID, performerImageBlobColumn)
}

func (qb *PerformerStore) stashIDRepository() *stashIDRepository {
	return &stashIDRepository{
		repository{
			tx:        qb.tx,
			tableName: "performer_stash_ids",
			idColumn:  performerIDColumn,
		},
	}
}

func (qb *PerformerStore) GetAliases(ctx context.Context, performerID int) ([]string, error) {
	return performersAliasesTableMgr.get(ctx, performerID)
}

func (qb *PerformerStore) GetStashIDs(ctx context.Context, performerID int) ([]models.StashID, error) {
	return performersStashIDsTableMgr.get(ctx, performerID)
}

func (qb *PerformerStore) FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Performer, error) {
	sq := dialect.From(performersStashIDsJoinTable).Select(performersStashIDsJoinTable.Col(performerIDColumn)).Where(
		performersStashIDsJoinTable.Col("stash_id").Eq(stashID.StashID),
		performersStashIDsJoinTable.Col("endpoint").Eq(stashID.Endpoint),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting performers for stash ID %s: %w", stashID.StashID, err)
	}

	return ret, nil
}

func (qb *PerformerStore) FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*models.Performer, error) {
	table := qb.table()
	sq := dialect.From(table).LeftJoin(
		performersStashIDsJoinTable,
		goqu.On(table.Col(idColumn).Eq(performersStashIDsJoinTable.Col(performerIDColumn))),
	).Select(table.Col(idColumn))

	if hasStashID {
		sq = sq.Where(
			performersStashIDsJoinTable.Col("stash_id").IsNotNull(),
			performersStashIDsJoinTable.Col("endpoint").Eq(stashboxEndpoint),
		)
	} else {
		sq = sq.Where(
			performersStashIDsJoinTable.Col("stash_id").IsNull(),
		)
	}

	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting performers for stash-box endpoint %s: %w", stashboxEndpoint, err)
	}

	return ret, nil
}
