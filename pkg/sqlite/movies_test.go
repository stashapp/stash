//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func loadMovieRelationships(ctx context.Context, expected models.Movie, actual *models.Movie) error {
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Gallery); err != nil {
			return err
		}
	}

	return nil
}

func TestMovieFindByName(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.Movie

		name := movieNames[movieIdxWithScene] // find a movie by name

		movie, err := mqb.FindByName(ctx, name, false)

		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}

		assert.Equal(t, movieNames[movieIdxWithScene], movie.Name)

		name = movieNames[movieIdxWithDupName] // find a movie by name nocase

		movie, err = mqb.FindByName(ctx, name, true)

		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}
		// movieIdxWithDupName and movieIdxWithScene should have similar names ( only diff should be Name vs NaMe)
		//movie.Name should match with movieIdxWithScene since its ID is before moveIdxWithDupName
		assert.Equal(t, movieNames[movieIdxWithScene], movie.Name)
		//movie.Name should match with movieIdxWithDupName if the check is not case sensitive
		assert.Equal(t, strings.ToLower(movieNames[movieIdxWithDupName]), strings.ToLower(movie.Name))

		return nil
	})
}

func TestMovieFindByNames(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		var names []string

		mqb := db.Movie

		names = append(names, movieNames[movieIdxWithScene]) // find movies by names

		movies, err := mqb.FindByNames(ctx, names, false)
		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}
		assert.Len(t, movies, 1)
		assert.Equal(t, movieNames[movieIdxWithScene], movies[0].Name)

		movies, err = mqb.FindByNames(ctx, names, true) // find movies by names nocase
		if err != nil {
			t.Errorf("Error finding movies: %s", err.Error())
		}
		assert.Len(t, movies, 2) // movieIdxWithScene and movieIdxWithDupName
		assert.Equal(t, strings.ToLower(movieNames[movieIdxWithScene]), strings.ToLower(movies[0].Name))
		assert.Equal(t, strings.ToLower(movieNames[movieIdxWithScene]), strings.ToLower(movies[1].Name))

		return nil
	})
}

func moviesToIDs(i []*models.Movie) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func TestMovieQuery(t *testing.T) {
	var (
		frontImage = "front_image"
		backImage  = "back_image"
	)

	tests := []struct {
		name        string
		findFilter  *models.FindFilterType
		filter      *models.MovieFilterType
		includeIdxs []int
		excludeIdxs []int
		wantErr     bool
	}{
		{
			"is missing front image",
			nil,
			&models.MovieFilterType{
				IsMissing: &frontImage,
			},
			// just ensure that it doesn't error
			nil,
			nil,
			false,
		},
		{
			"is missing back image",
			nil,
			&models.MovieFilterType{
				IsMissing: &backImage,
			},
			// just ensure that it doesn't error
			nil,
			nil,
			false,
		},
	}

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			results, _, err := db.Movie.Query(ctx, tt.filter, tt.findFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("MovieQueryBuilder.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ids := moviesToIDs(results)
			include := indexesToIDs(performerIDs, tt.includeIdxs)
			exclude := indexesToIDs(performerIDs, tt.excludeIdxs)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}
}

func TestMovieQueryStudio(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.Movie
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

		urls := n.URLs.List()
		var url string
		if len(urls) > 0 {
			url = urls[0]
		}

		verifyString(t, url, urlCriterion)
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

func TestMovieQueryURLExcludes(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		mqb := db.Movie

		// create movie with two URLs
		movie := models.Movie{
			Name: "TestMovieQueryURLExcludes",
			URLs: models.NewRelatedStrings([]string{
				"aaa",
				"bbb",
			}),
		}

		err := mqb.Create(ctx, &movie)

		if err != nil {
			return fmt.Errorf("Error creating movie: %w", err)
		}

		// query for movies that exclude the URL "aaa"
		urlCriterion := models.StringCriterionInput{
			Value:    "aaa",
			Modifier: models.CriterionModifierExcludes,
		}

		nameCriterion := models.StringCriterionInput{
			Value:    movie.Name,
			Modifier: models.CriterionModifierEquals,
		}

		filter := models.MovieFilterType{
			URL:  &urlCriterion,
			Name: &nameCriterion,
		}

		movies := queryMovie(ctx, t, mqb, &filter, nil)
		assert.Len(t, movies, 0, "Expected no movies to be found")

		// query for movies that exclude the URL "ccc"
		urlCriterion.Value = "ccc"
		movies = queryMovie(ctx, t, mqb, &filter, nil)

		if assert.Len(t, movies, 1, "Expected one movie to be found") {
			assert.Equal(t, movie.Name, movies[0].Name)
		}

		return nil
	})
}

func verifyMovieQuery(t *testing.T, filter models.MovieFilterType, verifyFn func(s *models.Movie)) {
	withTxn(func(ctx context.Context) error {
		t.Helper()
		sqb := db.Movie

		movies := queryMovie(ctx, t, sqb, &filter, nil)

		for _, movie := range movies {
			if err := movie.LoadURLs(ctx, sqb); err != nil {
				t.Errorf("Error loading movie relationships: %v", err)
			}
		}

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
		sqb := db.Movie
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

func TestMovieUpdateFrontImage(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Movie

		// create movie to test against
		const name = "TestMovieUpdateMovieImages"
		movie := models.Movie{
			Name: name,
		}
		err := qb.Create(ctx, &movie)
		if err != nil {
			return fmt.Errorf("Error creating movie: %s", err.Error())
		}

		return testUpdateImage(t, ctx, movie.ID, qb.UpdateFrontImage, qb.GetFrontImage)
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestMovieUpdateBackImage(t *testing.T) {
	if err := withRollbackTxn(func(ctx context.Context) error {
		qb := db.Movie

		// create movie to test against
		const name = "TestMovieUpdateMovieImages"
		movie := models.Movie{
			Name: name,
		}
		err := qb.Create(ctx, &movie)
		if err != nil {
			return fmt.Errorf("Error creating movie: %s", err.Error())
		}

		return testUpdateImage(t, ctx, movie.ID, qb.UpdateBackImage, qb.GetBackImage)
	}); err != nil {
		t.Error(err.Error())
	}
}

// TODO Update
// TODO Destroy - ensure image is destroyed
// TODO Find
// TODO Count
// TODO All
// TODO Query
