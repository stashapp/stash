// +build integration
package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func TestMovieFindBySceneID(t *testing.T) {
	mqb := models.NewMovieQueryBuilder()
	sceneID := sceneIDs[sceneIdxWithMovie]

	movies, err := mqb.FindBySceneID(sceneID, nil)

	if err != nil {
		t.Fatalf("Error finding movie: %s", err.Error())
	}

	assert.Equal(t, 1, len(movies), "expect 1 movie")

	movie := movies[0]
	assert.Equal(t, getMovieStringValue(movieIdxWithScene, "Name"), movie.Name.String)

	movies, err = mqb.FindBySceneID(0, nil)

	if err != nil {
		t.Fatalf("Error finding movie: %s", err.Error())
	}

	assert.Equal(t, 0, len(movies))
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO FindByName
// TODO FindByNames
// TODO Count
// TODO All
// TODO Query
