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
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

const (
	movieTable    = "movies"
	movieIDColumn = "movie_id"

	movieFrontImageBlobColumn = "front_image_blob"
	movieBackImageBlobColumn  = "back_image_blob"
)

type movieRow struct {
	ID       int         `db:"id" goqu:"skipinsert"`
	Name     zero.String `db:"name"`
	Aliases  zero.String `db:"aliases"`
	Duration null.Int    `db:"duration"`
	Date     NullDate    `db:"date"`
	// expressed as 1-100
	Rating    null.Int    `db:"rating"`
	StudioID  null.Int    `db:"studio_id,omitempty"`
	Director  zero.String `db:"director"`
	Synopsis  zero.String `db:"synopsis"`
	URL       zero.String `db:"url"`
	CreatedAt Timestamp   `db:"created_at"`
	UpdatedAt Timestamp   `db:"updated_at"`

	// not used in resolutions or updates
	FrontImageBlob zero.String `db:"front_image_blob"`
	BackImageBlob  zero.String `db:"back_image_blob"`
}

func (r *movieRow) fromMovie(o models.Movie) {
	r.ID = o.ID
	r.Name = zero.StringFrom(o.Name)
	r.Aliases = zero.StringFrom(o.Aliases)
	r.Duration = intFromPtr(o.Duration)
	r.Date = NullDateFromDatePtr(o.Date)
	r.Rating = intFromPtr(o.Rating)
	r.StudioID = intFromPtr(o.StudioID)
	r.Director = zero.StringFrom(o.Director)
	r.Synopsis = zero.StringFrom(o.Synopsis)
	r.URL = zero.StringFrom(o.URL)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *movieRow) resolve() *models.Movie {
	ret := &models.Movie{
		ID:        r.ID,
		Name:      r.Name.String,
		Aliases:   r.Aliases.String,
		Duration:  nullIntPtr(r.Duration),
		Date:      r.Date.DatePtr(),
		Rating:    nullIntPtr(r.Rating),
		StudioID:  nullIntPtr(r.StudioID),
		Director:  r.Director.String,
		Synopsis:  r.Synopsis.String,
		URL:       r.URL.String,
		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,
	}

	return ret
}

type movieRowRecord struct {
	updateRecord
}

func (r *movieRowRecord) fromPartial(o models.MoviePartial) {
	r.setNullString("name", o.Name)
	r.setNullString("aliases", o.Aliases)
	r.setNullInt("duration", o.Duration)
	r.setNullDate("date", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setNullInt("studio_id", o.StudioID)
	r.setNullString("director", o.Director)
	r.setNullString("synopsis", o.Synopsis)
	r.setNullString("url", o.URL)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
}

type MovieStore struct {
	repository
	blobJoinQueryBuilder

	tableMgr *table
}

func NewMovieStore(blobStore *BlobStore) *MovieStore {
	return &MovieStore{
		repository: repository{
			tableName: movieTable,
			idColumn:  idColumn,
		},
		blobJoinQueryBuilder: blobJoinQueryBuilder{
			blobStore: blobStore,
			joinTable: movieTable,
		},

		tableMgr: movieTableMgr,
	}
}

func (qb *MovieStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *MovieStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *MovieStore) Create(ctx context.Context, newObject *models.Movie) error {
	var r movieRow
	r.fromMovie(*newObject)

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

func (qb *MovieStore) UpdatePartial(ctx context.Context, id int, partial models.MoviePartial) (*models.Movie, error) {
	r := movieRowRecord{
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

func (qb *MovieStore) Update(ctx context.Context, updatedObject *models.Movie) error {
	var r movieRow
	r.fromMovie(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *MovieStore) Destroy(ctx context.Context, id int) error {
	// must handle image checksums manually
	if err := qb.destroyImages(ctx, id); err != nil {
		return err
	}

	return qb.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *MovieStore) Find(ctx context.Context, id int) (*models.Movie, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *MovieStore) FindMany(ctx context.Context, ids []int) ([]*models.Movie, error) {
	ret := make([]*models.Movie, len(ids))

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
			return nil, fmt.Errorf("movie with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *MovieStore) find(ctx context.Context, id int) (*models.Movie, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *MovieStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Movie, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *MovieStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Movie, error) {
	const single = false
	var ret []*models.Movie
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f movieRow
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

func (qb *MovieStore) FindByName(ctx context.Context, name string, nocase bool) (*models.Movie, error) {
	// query := "SELECT * FROM movies WHERE name = ?"
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

func (qb *MovieStore) FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Movie, error) {
	// query := "SELECT * FROM movies WHERE name"
	// if nocase {
	// 	query += " COLLATE NOCASE"
	// }
	// query += " IN " + getInBinding(len(names))
	where := "name"
	if nocase {
		where += " COLLATE NOCASE"
	}
	where += " IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	sq := qb.selectDataset().Prepared(true).Where(goqu.L(where, args...))
	ret, err := qb.getMany(ctx, sq)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *MovieStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *MovieStore) All(ctx context.Context) ([]*models.Movie, error) {
	table := qb.table()

	return qb.getMany(ctx, qb.selectDataset().Order(
		table.Col("name").Asc(),
		table.Col(idColumn).Asc(),
	))
}

func (qb *MovieStore) makeFilter(ctx context.Context, movieFilter *models.MovieFilterType) *filterBuilder {
	query := &filterBuilder{}

	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.Name, "movies.name"))
	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.Director, "movies.director"))
	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.Synopsis, "movies.synopsis"))
	query.handleCriterion(ctx, intCriterionHandler(movieFilter.Rating100, "movies.rating", nil))
	// legacy rating handler
	query.handleCriterion(ctx, rating5CriterionHandler(movieFilter.Rating, "movies.rating", nil))
	query.handleCriterion(ctx, floatIntCriterionHandler(movieFilter.Duration, "movies.duration", nil))
	query.handleCriterion(ctx, movieIsMissingCriterionHandler(qb, movieFilter.IsMissing))
	query.handleCriterion(ctx, stringCriterionHandler(movieFilter.URL, "movies.url"))
	query.handleCriterion(ctx, studioCriterionHandler(movieTable, movieFilter.Studios))
	query.handleCriterion(ctx, moviePerformersCriterionHandler(qb, movieFilter.Performers))
	query.handleCriterion(ctx, dateCriterionHandler(movieFilter.Date, "movies.date"))
	query.handleCriterion(ctx, timestampCriterionHandler(movieFilter.CreatedAt, "movies.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(movieFilter.UpdatedAt, "movies.updated_at"))

	return query
}

func (qb *MovieStore) makeQuery(ctx context.Context, movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}
	if movieFilter == nil {
		movieFilter = &models.MovieFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, movieTable)

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"movies.name"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := qb.makeFilter(ctx, movieFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	query.sortAndPagination = qb.getMovieSort(findFilter) + getPagination(findFilter)

	return &query, nil
}

func (qb *MovieStore) Query(ctx context.Context, movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) ([]*models.Movie, int, error) {
	query, err := qb.makeQuery(ctx, movieFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	movies, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return movies, countResult, nil
}

func (qb *MovieStore) QueryCount(ctx context.Context, movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, movieFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

func movieIsMissingCriterionHandler(qb *MovieStore, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "front_image":
				f.addWhere("movies.front_image_blob IS NULL")
			case "back_image":
				f.addWhere("movies.back_image_blob IS NULL")
			case "scenes":
				f.addLeftJoin("movies_scenes", "", "movies_scenes.movie_id = movies.id")
				f.addWhere("movies_scenes.scene_id IS NULL")
			default:
				f.addWhere("(movies." + *isMissing + " IS NULL OR TRIM(movies." + *isMissing + ") = '')")
			}
		}
	}
}

func moviePerformersCriterionHandler(qb *MovieStore, performers *models.MultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performers != nil {
			if performers.Modifier == models.CriterionModifierIsNull || performers.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if performers.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("movies_scenes", "", "movies.id = movies_scenes.movie_id")
				f.addLeftJoin("performers_scenes", "", "movies_scenes.scene_id = performers_scenes.scene_id")

				f.addWhere(fmt.Sprintf("performers_scenes.performer_id IS %s NULL", notClause))
				return
			}

			if len(performers.Value) == 0 {
				return
			}

			var args []interface{}
			for _, arg := range performers.Value {
				args = append(args, arg)
			}

			// Hack, can't apply args to join, nor inner join on a left join, so use CTE instead
			f.addWith(`movies_performers AS (
				SELECT movies_scenes.movie_id, performers_scenes.performer_id
				FROM movies_scenes
				INNER JOIN performers_scenes ON movies_scenes.scene_id = performers_scenes.scene_id
				WHERE performers_scenes.performer_id IN`+getInBinding(len(performers.Value))+`
			)`, args...)
			f.addLeftJoin("movies_performers", "", "movies.id = movies_performers.movie_id")

			switch performers.Modifier {
			case models.CriterionModifierIncludes:
				f.addWhere("movies_performers.performer_id IS NOT NULL")
			case models.CriterionModifierIncludesAll:
				f.addWhere("movies_performers.performer_id IS NOT NULL")
				f.addHaving("COUNT(DISTINCT movies_performers.performer_id) = ?", len(performers.Value))
			case models.CriterionModifierExcludes:
				f.addWhere("movies_performers.performer_id IS NULL")
			}
		}
	}
}

func (qb *MovieStore) getMovieSort(findFilter *models.FindFilterType) string {
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
	case "scenes_count": // generic getSort won't work for this
		sortQuery += getCountSort(movieTable, moviesScenesTable, movieIDColumn, direction)
	default:
		sortQuery += getSort(sort, direction, "movies")
	}

	// Whatever the sorting, always use name/id as a final sort
	sortQuery += ", COALESCE(movies.name, movies.id) COLLATE NATURAL_CI ASC"
	return sortQuery
}

func (qb *MovieStore) queryMovies(ctx context.Context, query string, args []interface{}) ([]*models.Movie, error) {
	const single = false
	var ret []*models.Movie
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f movieRow
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

func (qb *MovieStore) UpdateFrontImage(ctx context.Context, movieID int, frontImage []byte) error {
	return qb.UpdateImage(ctx, movieID, movieFrontImageBlobColumn, frontImage)
}

func (qb *MovieStore) UpdateBackImage(ctx context.Context, movieID int, backImage []byte) error {
	return qb.UpdateImage(ctx, movieID, movieBackImageBlobColumn, backImage)
}

func (qb *MovieStore) destroyImages(ctx context.Context, movieID int) error {
	if err := qb.DestroyImage(ctx, movieID, movieFrontImageBlobColumn); err != nil {
		return err
	}
	if err := qb.DestroyImage(ctx, movieID, movieBackImageBlobColumn); err != nil {
		return err
	}

	return nil
}

func (qb *MovieStore) GetFrontImage(ctx context.Context, movieID int) ([]byte, error) {
	return qb.GetImage(ctx, movieID, movieFrontImageBlobColumn)
}

func (qb *MovieStore) HasFrontImage(ctx context.Context, movieID int) (bool, error) {
	return qb.HasImage(ctx, movieID, movieFrontImageBlobColumn)
}

func (qb *MovieStore) GetBackImage(ctx context.Context, movieID int) ([]byte, error) {
	return qb.GetImage(ctx, movieID, movieBackImageBlobColumn)
}

func (qb *MovieStore) HasBackImage(ctx context.Context, movieID int) (bool, error) {
	return qb.HasImage(ctx, movieID, movieBackImageBlobColumn)
}

func (qb *MovieStore) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Movie, error) {
	query := `SELECT DISTINCT movies.*
FROM movies
INNER JOIN movies_scenes ON movies.id = movies_scenes.movie_id
INNER JOIN performers_scenes ON performers_scenes.scene_id = movies_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return qb.queryMovies(ctx, query, args)
}

func (qb *MovieStore) CountByPerformerID(ctx context.Context, performerID int) (int, error) {
	query := `SELECT COUNT(DISTINCT movies_scenes.movie_id) AS count
FROM movies_scenes
INNER JOIN performers_scenes ON performers_scenes.scene_id = movies_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return qb.runCountQuery(ctx, query, args)
}

func (qb *MovieStore) FindByStudioID(ctx context.Context, studioID int) ([]*models.Movie, error) {
	query := `SELECT movies.*
FROM movies
WHERE movies.studio_id = ?
`
	args := []interface{}{studioID}
	return qb.queryMovies(ctx, query, args)
}

func (qb *MovieStore) CountByStudioID(ctx context.Context, studioID int) (int, error) {
	query := `SELECT COUNT(1) AS count
FROM movies
WHERE movies.studio_id = ?
`
	args := []interface{}{studioID}
	return qb.runCountQuery(ctx, query, args)
}
