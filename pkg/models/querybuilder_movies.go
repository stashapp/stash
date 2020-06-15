package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

type MovieQueryBuilder struct{}

func NewMovieQueryBuilder() MovieQueryBuilder {
	return MovieQueryBuilder{}
}

func (qb *MovieQueryBuilder) Create(newMovie Movie, tx *sqlx.Tx) (*Movie, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO movies (front_image, back_image, checksum, name, aliases, duration, date, rating, studio_id, director, synopsis, url, created_at, updated_at)
				VALUES (:front_image, :back_image, :checksum, :name, :aliases, :duration, :date, :rating, :studio_id, :director, :synopsis, :url, :created_at, :updated_at)
		`,
		newMovie,
	)
	if err != nil {
		return nil, err
	}
	movieID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newMovie, `SELECT * FROM movies WHERE id = ? LIMIT 1`, movieID); err != nil {
		return nil, err
	}
	return &newMovie, nil
}

func (qb *MovieQueryBuilder) Update(updatedMovie MoviePartial, tx *sqlx.Tx) (*Movie, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE movies SET `+SQLGenKeysPartial(updatedMovie)+` WHERE movies.id = :id`,
		updatedMovie,
	)
	if err != nil {
		return nil, err
	}

	return qb.Find(updatedMovie.ID, tx)
}

func (qb *MovieQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	// delete movie from movies_scenes

	_, err := tx.Exec("DELETE FROM movies_scenes WHERE movie_id = ?", id)
	if err != nil {
		return err
	}

	// // remove movie from scraped items
	// _, err = tx.Exec("UPDATE scraped_items SET movie_id = null WHERE movie_id = ?", id)
	// if err != nil {
	// 	return err
	// }

	return executeDeleteQuery("movies", id, tx)
}

func (qb *MovieQueryBuilder) Find(id int, tx *sqlx.Tx) (*Movie, error) {
	query := "SELECT * FROM movies WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryMovie(query, args, tx)
}

func (qb *MovieQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*Movie, error) {
	query := `
		SELECT movies.* FROM movies
		LEFT JOIN movies_scenes as scenes_join on scenes_join.movie_id = movies.id
		WHERE scenes_join.scene_id = ?
		GROUP BY movies.id
	`
	args := []interface{}{sceneID}
	return qb.queryMovies(query, args, tx)
}

func (qb *MovieQueryBuilder) FindByName(name string, tx *sqlx.Tx, nocase bool) (*Movie, error) {
	query := "SELECT * FROM movies WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryMovie(query, args, tx)
}

func (qb *MovieQueryBuilder) FindByNames(names []string, tx *sqlx.Tx, nocase bool) ([]*Movie, error) {
	query := "SELECT * FROM movies WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryMovies(query, args, tx)
}

func (qb *MovieQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT movies.id FROM movies"), nil)
}

func (qb *MovieQueryBuilder) All() ([]*Movie, error) {
	return qb.queryMovies(selectAll("movies")+qb.getMovieSort(nil), nil, nil)
}

func (qb *MovieQueryBuilder) AllSlim() ([]*Movie, error) {
	return qb.queryMovies("SELECT movies.id, movies.name FROM movies "+qb.getMovieSort(nil), nil, nil)
}

func (qb *MovieQueryBuilder) Query(movieFilter *MovieFilterType, findFilter *FindFilterType) ([]*Movie, int) {
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}
	if movieFilter == nil {
		movieFilter = &MovieFilterType{}
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

	sortAndPagination := qb.getMovieSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("movies", body, args, sortAndPagination, whereClauses, havingClauses)

	var movies []*Movie
	for _, id := range idsResult {
		movie, _ := qb.Find(id, nil)
		movies = append(movies, movie)
	}

	return movies, countResult
}

func (qb *MovieQueryBuilder) getMovieSort(findFilter *FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "movies")
}

func (qb *MovieQueryBuilder) queryMovie(query string, args []interface{}, tx *sqlx.Tx) (*Movie, error) {
	results, err := qb.queryMovies(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *MovieQueryBuilder) queryMovies(query string, args []interface{}, tx *sqlx.Tx) ([]*Movie, error) {
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, args...)
	} else {
		rows, err = database.DB.Queryx(query, args...)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	movies := make([]*Movie, 0)
	for rows.Next() {
		movie := Movie{}
		if err := rows.StructScan(&movie); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
