// +build integration

package models_test

import (
	"strconv"
	"strings"
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

func TestMovieFindByName(t *testing.T) {

	mqb := models.NewMovieQueryBuilder()

	name := movieNames[movieIdxWithScene] // find a movie by name

	movie, err := mqb.FindByName(name, nil, false)

	if err != nil {
		t.Fatalf("Error finding movies: %s", err.Error())
	}

	assert.Equal(t, movieNames[movieIdxWithScene], movie.Name.String)

	name = movieNames[movieIdxWithDupName] // find a movie by name nocase

	movie, err = mqb.FindByName(name, nil, true)

	if err != nil {
		t.Fatalf("Error finding movies: %s", err.Error())
	}
	// movieIdxWithDupName and movieIdxWithScene should have similar names ( only diff should be Name vs NaMe)
	//movie.Name should match with movieIdxWithScene since its ID is before moveIdxWithDupName
	assert.Equal(t, movieNames[movieIdxWithScene], movie.Name.String)
	//movie.Name should match with movieIdxWithDupName if the check is not case sensitive
	assert.Equal(t, strings.ToLower(movieNames[movieIdxWithDupName]), strings.ToLower(movie.Name.String))
}

func TestMovieFindByNames(t *testing.T) {
	var names []string

	mqb := models.NewMovieQueryBuilder()

	names = append(names, movieNames[movieIdxWithScene]) // find movies by names

	movies, err := mqb.FindByNames(names, nil, false)
	if err != nil {
		t.Fatalf("Error finding movies: %s", err.Error())
	}
	assert.Len(t, movies, 1)
	assert.Equal(t, movieNames[movieIdxWithScene], movies[0].Name.String)

	movies, err = mqb.FindByNames(names, nil, true) // find movies by names nocase
	if err != nil {
		t.Fatalf("Error finding movies: %s", err.Error())
	}
	assert.Len(t, movies, 2) // movieIdxWithScene and movieIdxWithDupName
	assert.Equal(t, strings.ToLower(movieNames[movieIdxWithScene]), strings.ToLower(movies[0].Name.String))
	assert.Equal(t, strings.ToLower(movieNames[movieIdxWithScene]), strings.ToLower(movies[1].Name.String))
}

func TestMovieQueryStudio(t *testing.T) {
	mqb := models.NewMovieQueryBuilder()
	studioCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithMovie]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	movieFilter := models.MovieFilterType{
		Studios: &studioCriterion,
	}

	movies, _ := mqb.Query(&movieFilter, nil)

	assert.Len(t, movies, 1)

	// ensure id is correct
	assert.Equal(t, movieIDs[movieIdxWithStudio], movies[0].ID)

	studioCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithMovie]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getMovieStringValue(movieIdxWithStudio, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	movies, _ = mqb.Query(&movieFilter, &findFilter)
	assert.Len(t, movies, 0)
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO Count
// TODO All
// TODO Query
