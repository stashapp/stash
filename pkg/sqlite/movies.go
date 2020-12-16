package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const movieTable = "movies"

type MovieQueryBuilder struct {
	repository
}

func NewMovieReaderWriter(tx *sqlx.Tx) *MovieQueryBuilder {
	return &MovieQueryBuilder{
		repository{
			tx:        tx,
			tableName: movieTable,
			idColumn:  idColumn,
			constructor: func() interface{} {
				return &models.Movie{}
			},
		},
	}
}

func (qb *MovieQueryBuilder) Create(newObject models.Movie) (*models.Movie, error) {
	var ret models.Movie
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *MovieQueryBuilder) Update(updatedObject models.MoviePartial) (*models.Movie, error) {
	const partial = true
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *MovieQueryBuilder) UpdateFull(updatedObject models.Movie) (*models.Movie, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *MovieQueryBuilder) Destroy(id int) error {
	return qb.destroyExisting([]int{id})
}

func (qb *MovieQueryBuilder) Find(id int) (*models.Movie, error) {
	var ret models.Movie
	if err := qb.get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *MovieQueryBuilder) FindMany(ids []int) ([]*models.Movie, error) {
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

func (qb *MovieQueryBuilder) FindByName(name string, nocase bool) (*models.Movie, error) {
	query := "SELECT * FROM movies WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryMovie(query, args)
}

func (qb *MovieQueryBuilder) FindByNames(names []string, nocase bool) ([]*models.Movie, error) {
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

func (qb *MovieQueryBuilder) Count() (int, error) {
	return qb.runCountQuery(qb.buildCountQuery("SELECT movies.id FROM movies"), nil)
}

func (qb *MovieQueryBuilder) All() ([]*models.Movie, error) {
	return qb.queryMovies(selectAll("movies")+qb.getMovieSort(nil), nil)
}

func (qb *MovieQueryBuilder) AllSlim() ([]*models.Movie, error) {
	return qb.queryMovies("SELECT movies.id, movies.name FROM movies "+qb.getMovieSort(nil), nil)
}

func (qb *MovieQueryBuilder) Query(movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) ([]*models.Movie, int, error) {
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}
	if movieFilter == nil {
		movieFilter = &models.MovieFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("movies")
	body += `
	left join movies_scenes as scenes_join on scenes_join.movie_id = movies.id
	left join scenes on scenes_join.scene_id = scenes.id
	left join studios as studio on studio.id = movies.studio_id
`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"movies.name"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		whereClauses = append(whereClauses, clause)
		args = append(args, thisArgs...)
	}

	if studiosFilter := movieFilter.Studios; studiosFilter != nil && len(studiosFilter.Value) > 0 {
		for _, studioID := range studiosFilter.Value {
			args = append(args, studioID)
		}

		whereClause, havingClause := getMultiCriterionClause("movies", "studio", "", "", "studio_id", studiosFilter)
		whereClauses = appendClause(whereClauses, whereClause)
		havingClauses = appendClause(havingClauses, havingClause)
	}

	if isMissingFilter := movieFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "front_image":
			body += `left join movies_images on movies_images.movie_id = movies.id
			`
			whereClauses = appendClause(whereClauses, "movies_images.front_image IS NULL")
		case "back_image":
			body += `left join movies_images on movies_images.movie_id = movies.id
			`
			whereClauses = appendClause(whereClauses, "movies_images.back_image IS NULL")
		case "scenes":
			body += `left join movies_scenes on movies_scenes.movie_id = movies.id
			`
			whereClauses = appendClause(whereClauses, "movies_scenes.scene_id IS NULL")
		default:
			whereClauses = appendClause(whereClauses, "movies."+*isMissingFilter+" IS NULL")
		}
	}

	sortAndPagination := qb.getMovieSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := qb.executeFindQuery(body, args, sortAndPagination, whereClauses, havingClauses)
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

func (qb *MovieQueryBuilder) getMovieSort(findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}

	// #943 - override name sorting to use natural sort
	if sort == "name" {
		return " ORDER BY " + getColumn("movies", sort) + " COLLATE NATURAL_CS " + direction
	}

	return getSort(sort, direction, "movies")
}

func (qb *MovieQueryBuilder) queryMovie(query string, args []interface{}) (*models.Movie, error) {
	results, err := qb.queryMovies(query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *MovieQueryBuilder) queryMovies(query string, args []interface{}) ([]*models.Movie, error) {
	var ret models.Movies
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Movie(ret), nil
}

func (qb *MovieQueryBuilder) UpdateImages(movieID int, frontImage []byte, backImage []byte) error {
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

func (qb *MovieQueryBuilder) DestroyImages(movieID int) error {
	// Delete the existing joins
	_, err := qb.tx.Exec("DELETE FROM movies_images WHERE movie_id = ?", movieID)
	if err != nil {
		return err
	}
	return err
}

func (qb *MovieQueryBuilder) GetFrontImage(movieID int) ([]byte, error) {
	query := `SELECT front_image from movies_images WHERE movie_id = ?`
	return getImage(qb.tx, query, movieID)
}

func (qb *MovieQueryBuilder) GetBackImage(movieID int) ([]byte, error) {
	query := `SELECT back_image from movies_images WHERE movie_id = ?`
	return getImage(qb.tx, query, movieID)
}
