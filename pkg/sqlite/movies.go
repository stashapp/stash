package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

const movieTable = "movies"
const movieIDColumn = "movie_id"

type movieQueryBuilder struct {
	repository
}

func NewMovieReaderWriter(tx dbi) *movieQueryBuilder {
	return &movieQueryBuilder{
		repository{
			tx:        tx,
			tableName: movieTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *movieQueryBuilder) Create(newObject models.Movie) (*models.Movie, error) {
	var ret models.Movie
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *movieQueryBuilder) Update(updatedObject models.MoviePartial) (*models.Movie, error) {
	const partial = true
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *movieQueryBuilder) UpdateFull(updatedObject models.Movie) (*models.Movie, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *movieQueryBuilder) Destroy(id int) error {
	return qb.destroyExisting([]int{id})
}

func (qb *movieQueryBuilder) Find(id int) (*models.Movie, error) {
	var ret models.Movie
	if err := qb.get(id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *movieQueryBuilder) FindMany(ids []int) ([]*models.Movie, error) {
	var movies []*models.Movie
	for _, id := range ids {
		movie, err := qb.Find(id)
		if err != nil {
			return nil, err
		}

		if movie == nil {
			return nil, fmt.Errorf("movie with id %d not found", id)
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (qb *movieQueryBuilder) FindByName(name string, nocase bool) (*models.Movie, error) {
	query := "SELECT * FROM movies WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryMovie(query, args)
}

func (qb *movieQueryBuilder) FindByNames(names []string, nocase bool) ([]*models.Movie, error) {
	query := "SELECT * FROM movies WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryMovies(query, args)
}

func (qb *movieQueryBuilder) Count() (int, error) {
	return qb.runCountQuery(qb.buildCountQuery("SELECT movies.id FROM movies"), nil)
}

func (qb *movieQueryBuilder) All() ([]*models.Movie, error) {
	return qb.queryMovies(selectAll("movies")+qb.getMovieSort(nil), nil)
}

func (qb *movieQueryBuilder) makeFilter(movieFilter *models.MovieFilterType) *filterBuilder {
	query := &filterBuilder{}

	query.handleCriterion(stringCriterionHandler(movieFilter.Name, "movies.name"))
	query.handleCriterion(stringCriterionHandler(movieFilter.Director, "movies.director"))
	query.handleCriterion(stringCriterionHandler(movieFilter.Synopsis, "movies.synopsis"))
	query.handleCriterion(intCriterionHandler(movieFilter.Rating, "movies.rating"))
	query.handleCriterion(durationCriterionHandler(movieFilter.Duration, "movies.duration"))
	query.handleCriterion(movieIsMissingCriterionHandler(qb, movieFilter.IsMissing))
	query.handleCriterion(stringCriterionHandler(movieFilter.URL, "movies.url"))
	query.handleCriterion(movieStudioCriterionHandler(qb, movieFilter.Studios))
	query.handleCriterion(moviePerformersCriterionHandler(qb, movieFilter.Performers))

	return query
}

func (qb *movieQueryBuilder) Query(movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) ([]*models.Movie, int, error) {
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

	filter := qb.makeFilter(movieFilter)

	query.addFilter(filter)

	query.sortAndPagination = qb.getMovieSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind()
	if err != nil {
		return nil, 0, err
	}

	var movies []*models.Movie
	for _, id := range idsResult {
		movie, err := qb.Find(id)
		if err != nil {
			return nil, 0, err
		}

		movies = append(movies, movie)
	}

	return movies, countResult, nil
}

func movieIsMissingCriterionHandler(qb *movieQueryBuilder, isMissing *string) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "front_image":
				f.addLeftJoin("movies_images", "", "movies_images.movie_id = movies.id")
				f.addWhere("movies_images.front_image IS NULL")
			case "back_image":
				f.addLeftJoin("movies_images", "", "movies_images.movie_id = movies.id")
				f.addWhere("movies_images.back_image IS NULL")
			case "scenes":
				f.addLeftJoin("movies_scenes", "", "movies_scenes.movie_id = movies.id")
				f.addWhere("movies_scenes.scene_id IS NULL")
			default:
				f.addWhere("(movies." + *isMissing + " IS NULL OR TRIM(movies." + *isMissing + ") = '')")
			}
		}
	}
}

func movieStudioCriterionHandler(qb *movieQueryBuilder, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := hierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: movieTable,
		foreignTable: studioTable,
		foreignFK:    studioIDColumn,
		derivedTable: "studio",
		parentFK:     "parent_id",
	}

	return h.handler(studios)
}

func moviePerformersCriterionHandler(qb *movieQueryBuilder, performers *models.MultiCriterionInput) criterionHandlerFunc {
	return func(f *filterBuilder) {
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

func (qb *movieQueryBuilder) getMovieSort(findFilter *models.FindFilterType) string {
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
	case "name": // #943 - override name sorting to use natural sort
		return " ORDER BY " + getColumn("movies", sort) + " COLLATE NATURAL_CS " + direction
	case "scenes_count": // generic getSort won't work for this
		return getCountSort(movieTable, moviesScenesTable, movieIDColumn, direction)
	default:
		return getSort(sort, direction, "movies")
	}
}

func (qb *movieQueryBuilder) queryMovie(query string, args []interface{}) (*models.Movie, error) {
	results, err := qb.queryMovies(query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *movieQueryBuilder) queryMovies(query string, args []interface{}) ([]*models.Movie, error) {
	var ret models.Movies
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Movie(ret), nil
}

func (qb *movieQueryBuilder) UpdateImages(movieID int, frontImage []byte, backImage []byte) error {
	// Delete the existing cover and then create new
	if err := qb.DestroyImages(movieID); err != nil {
		return err
	}

	_, err := qb.tx.Exec(
		`INSERT INTO movies_images (movie_id, front_image, back_image) VALUES (?, ?, ?)`,
		movieID,
		frontImage,
		backImage,
	)

	return err
}

func (qb *movieQueryBuilder) DestroyImages(movieID int) error {
	// Delete the existing joins
	_, err := qb.tx.Exec("DELETE FROM movies_images WHERE movie_id = ?", movieID)
	if err != nil {
		return err
	}
	return err
}

func (qb *movieQueryBuilder) GetFrontImage(movieID int) ([]byte, error) {
	query := `SELECT front_image from movies_images WHERE movie_id = ?`
	return getImage(qb.tx, query, movieID)
}

func (qb *movieQueryBuilder) GetBackImage(movieID int) ([]byte, error) {
	query := `SELECT back_image from movies_images WHERE movie_id = ?`
	return getImage(qb.tx, query, movieID)
}

func (qb *movieQueryBuilder) FindByPerformerID(performerID int) ([]*models.Movie, error) {
	query := `SELECT DISTINCT movies.*
FROM movies
INNER JOIN movies_scenes ON movies.id = movies_scenes.movie_id
INNER JOIN performers_scenes ON performers_scenes.scene_id = movies_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return qb.queryMovies(query, args)
}

func (qb *movieQueryBuilder) CountByPerformerID(performerID int) (int, error) {
	query := `SELECT COUNT(DISTINCT movies_scenes.movie_id) AS count
FROM movies_scenes
INNER JOIN performers_scenes ON performers_scenes.scene_id = movies_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return qb.runCountQuery(query, args)
}

func (qb *movieQueryBuilder) FindByStudioID(studioID int) ([]*models.Movie, error) {
	query := `SELECT movies.*
FROM movies
WHERE movies.studio_id = ?
`
	args := []interface{}{studioID}
	return qb.queryMovies(query, args)
}

func (qb *movieQueryBuilder) CountByStudioID(studioID int) (int, error) {
	query := `SELECT COUNT(1) AS count
FROM movies
WHERE movies.studio_id = ?
`
	args := []interface{}{studioID}
	return qb.runCountQuery(query, args)
}
