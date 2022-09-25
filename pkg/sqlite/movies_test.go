//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func TestMovieFindByName(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := sqlite.MovieReaderWriter

		name := movieNames[movieIdxWithScene] // find a movie by name

		movie, err := mqb.FindByName(ctx, name, false)

		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}

		assert.Equal(t, movieNames[movieIdxWithScene], movie.Name.String)

		name = movieNames[movieIdxWithDupName] // find a movie by name nocase

		movie, err = mqb.FindByName(ctx, name, true)

		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}
		// movieIdxWithDupName and movieIdxWithScene should have similar names ( only diff should be Name vs NaMe)
		//movie.Name should match with movieIdxWithScene since its ID is before moveIdxWithDupName
		assert.Equal(t, movieNames[movieIdxWithScene], movie.Name.String)
		//movie.Name should match with movieIdxWithDupName if the check is not case sensitive
		assert.Equal(t, strings.ToLower(movieNames[movieIdxWithDupName]), strings.ToLower(movie.Name.String))

		return nil
	})
}

func TestMovieFindByNames(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		var names []string

		mqb := sqlite.MovieReaderWriter

		names = append(names, movieNames[movieIdxWithScene]) // find movies by names

		movies, err := mqb.FindByNames(ctx, names, false)
		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}
		assert.Len(t, movies, 1)
		assert.Equal(t, movieNames[movieIdxWithScene], movies[0].Name.String)

		movies, err = mqb.FindByNames(ctx, names, true) // find movies by names nocase
		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}
		assert.Len(t, movies, 2) // movieIdxWithScene and movieIdxWithDupName
		assert.Equal(t, strings.ToLower(movieNames[movieIdxWithScene]), strings.ToLower(movies[0].Name.String))
		assert.Equal(t, strings.ToLower(movieNames[movieIdxWithScene]), strings.ToLower(movies[1].Name.String))

		return nil
	})
}

func TestMovieQueryStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := sqlite.MovieReaderWriter
		studioCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithMovie]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		movieFilter := models.MovieFilterType{
			Studios: &studioCriterion,
		}

		movies, _, err := mqb.Query(ctx, &movieFilter, nil)
		if err != nil {
			t.Errorf("Error querying movie: %s", err.Error())
		}

		assert.Len(t, movies, 1)

		// ensure id is correct
		assert.Equal(t, movieIDs[movieIdxWithStudio], movies[0].ID)

		studioCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithMovie]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getMovieStringValue(movieIdxWithStudio, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		movies, _, err = mqb.Query(ctx, &movieFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying movie: %s", err.Error())
		}
		assert.Len(t, movies, 0)

		return nil
	})
}

func TestMovieQueryURL(t *testing.T) {
	const sceneIdx = 1
	movieURL := getMovieStringValue(sceneIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    movieURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.MovieFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(n *models.Movie) {
		t.Helper()
		verifyNullString(t, n.URL, urlCriterion)
	}

	verifyMovieQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyMovieQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "movie_.*1_URL"
	verifyMovieQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyMovieQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyMovieQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyMovieQuery(t, filter, verifyFn)
}

func verifyMovieQuery(t *testing.T, filter models.MovieFilterType, verifyFn func(s *models.Movie)) {
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := sqlite.MovieReaderWriter

		movies := queryMovie(ctx, t, sqb, &filter, nil)

		// assume it should find at least one
		assert.Greater(t, len(movies), 0)

		for _, m := range movies {
			verifyFn(m)
		}

		return nil
	})
}

func queryMovie(ctx context.Context, t *testing.T, sqb models.MovieReader, movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) []*models.Movie {
	movies, _, err := sqb.Query(ctx, movieFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying movie: %s", err.Error())
	}

	return movies
}

func TestMovieQuerySorting(t *testing.T) {
	sort := "scenes_count"
	direction := models.SortDirectionEnumDesc
	findFilter := models.FindFilterType{
		Sort:      &sort,
		Direction: &direction,
	}

	withTxn(func(ctx context.Context) error {
		sqb := sqlite.MovieReaderWriter
		movies := queryMovie(ctx, t, sqb, nil, &findFilter)

		// scenes should be in same order as indexes
		firstMovie := movies[0]

		assert.Equal(t, movieIDs[movieIdxWithScene], firstMovie.ID)

		// sort in descending order
		direction = models.SortDirectionEnumAsc

		movies = queryMovie(ctx, t, sqb, nil, &findFilter)
		lastMovie := movies[len(movies)-1]

		assert.Equal(t, movieIDs[movieIdxWithScene], lastMovie.ID)

		return nil
	})
}

func TestMovieUpdateMovieImages(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		mqb := sqlite.MovieReaderWriter

		// create movie to test against
		const name = "TestMovieUpdateMovieImages"
		movie := models.Movie{
			Name:     sql.NullString{String: name, Valid: true},
			Checksum: md5.FromString(name),
		}
		created, err := mqb.Create(ctx, movie)
		if err != nil {
			return fmt.Errorf("Error creating movie: %s", err.Error())
		}

		frontImage := []byte("frontImage")
		backImage := []byte("backImage")
		err = mqb.UpdateImages(ctx, created.ID, frontImage, backImage)
		if err != nil {
			return fmt.Errorf("Error updating movie images: %s", err.Error())
		}

		// ensure images are set
		storedFront, err := mqb.GetFrontImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting front image: %s", err.Error())
		}
		assert.Equal(t, storedFront, frontImage)

		storedBack, err := mqb.GetBackImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting back image: %s", err.Error())
		}
		assert.Equal(t, storedBack, backImage)

		// set front image only
		newImage := []byte("newImage")
		err = mqb.UpdateImages(ctx, created.ID, newImage, nil)
		if err != nil {
			return fmt.Errorf("Error updating movie images: %s", err.Error())
		}

		storedFront, err = mqb.GetFrontImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting front image: %s", err.Error())
		}
		assert.Equal(t, storedFront, newImage)

		// back image should be nil
		storedBack, err = mqb.GetBackImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting back image: %s", err.Error())
		}
		assert.Nil(t, nil)

		// set back image only
		err = mqb.UpdateImages(ctx, created.ID, nil, newImage)
		if err == nil {
			return fmt.Errorf("Expected error setting nil front image")
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestMovieDestroyMovieImages(t *testing.T) {
	if err := withTxn(func(ctx context.Context) error {
		mqb := sqlite.MovieReaderWriter

		// create movie to test against
		const name = "TestMovieDestroyMovieImages"
		movie := models.Movie{
			Name:     sql.NullString{String: name, Valid: true},
			Checksum: md5.FromString(name),
		}
		created, err := mqb.Create(ctx, movie)
		if err != nil {
			return fmt.Errorf("Error creating movie: %s", err.Error())
		}

		frontImage := []byte("frontImage")
		backImage := []byte("backImage")
		err = mqb.UpdateImages(ctx, created.ID, frontImage, backImage)
		if err != nil {
			return fmt.Errorf("Error updating movie images: %s", err.Error())
		}

		err = mqb.DestroyImages(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error destroying movie images: %s", err.Error())
		}

		// front image should be nil
		storedFront, err := mqb.GetFrontImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting front image: %s", err.Error())
		}
		assert.Nil(t, storedFront)

		// back image should be nil
		storedBack, err := mqb.GetBackImage(ctx, created.ID)
		if err != nil {
			return fmt.Errorf("Error getting back image: %s", err.Error())
		}
		assert.Nil(t, storedBack)

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO Count
// TODO All
// TODO Query
