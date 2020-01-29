package models

// import (
// 	"database/sql"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/stashapp/stash/pkg/database"
// )

// NO NECESARIO COMPROBAR
// type MovieSceneQueryBuilder struct{}

// func NewMovieSceneQueryBuilder() MovieSceneQueryBuilder {
// 	return MovieSceneQueryBuilder{}
// }


// func (qb *MovieSceneQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*SceneNumberForSceneResultType, error) {
// 	query := "SELECT movie_id, scene_idx FROM movies_scenes WHERE scene_id = ?"
// 	args := []interface{}{sceneID}
// 	return qb.queryMoviesScenes(query, args, tx)
// }


// func (qb *MovieSceneQueryBuilder) queryMoviesScenes(query string, args []interface{}, tx *sqlx.Tx) ([]*SceneNumberForSceneResultType, error) {
// 	var rows *sqlx.Rows
// 	var err error
// 	if tx != nil {
// 		rows, err = tx.Queryx(query, args...)
// 	} else {
// 		rows, err = database.DB.Queryx(query, args...)
// 	}

// 	if err != nil && err != sql.ErrNoRows {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	moviesscenes := make([]*SceneNumberForSceneResultType, 0)
// 	for rows.Next() {
// 		moviescene := SceneNumberForSceneResultType{}
// 		if err := rows.StructScan(&moviescene); err != nil {
// 			return nil, err
// 		}
// 		moviesscenes = append(moviesscenes, &moviescene)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return moviesscenes, nil
// }
