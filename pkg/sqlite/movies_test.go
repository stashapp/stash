//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func loadMovieRelationships(ctx context.Context, expected models.Movie, actual *models.Movie) error {
	if expected.URLs.Loaded() {
		if err := actual.LoadURLs(ctx, db.Movie); err != nil {
			return err
		}
	}
	if expected.TagIDs.Loaded() {
		if err := actual.LoadTagIDs(ctx, db.Movie); err != nil {
			return err
		}
	}

	return nil
}

func Test_MovieStore_Create(t *testing.T) {
	var (
		name      = "name"
		url       = "url"
		aliases   = "alias1, alias2"
		director  = "director"
		rating    = 60
		duration  = 34
		synopsis  = "synopsis"
		date, _   = models.ParseDate("2003-02-01")
		createdAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		newObject models.Movie
		wantErr   bool
	}{
		{
			"full",
			models.Movie{
				Name:      name,
				Duration:  &duration,
				Date:      &date,
				Rating:    &rating,
				StudioID:  &studioIDs[studioIdxWithMovie],
				Director:  director,
				Synopsis:  synopsis,
				URLs:      models.NewRelatedStrings([]string{url}),
				TagIDs:    models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithMovie]}),
				Aliases:   aliases,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"invalid tag id",
			models.Movie{
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
	}

	qb := db.Movie

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			p := tt.newObject
			if err := qb.Create(ctx, &p); (err != nil) != tt.wantErr {
				t.Errorf("MovieStore.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(p.ID)
				return
			}

			assert.NotZero(p.ID)

			copy := tt.newObject
			copy.ID = p.ID

			// load relationships
			if err := loadMovieRelationships(ctx, copy, &p); err != nil {
				t.Errorf("loadMovieRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, p)

			// ensure can find the movie
			found, err := qb.Find(ctx, p.ID)
			if err != nil {
				t.Errorf("MovieStore.Find() error = %v", err)
			}

			if !assert.NotNil(found) {
				return
			}

			// load relationships
			if err := loadMovieRelationships(ctx, copy, found); err != nil {
				t.Errorf("loadMovieRelationships() error = %v", err)
				return
			}
			assert.Equal(copy, *found)

			return
		})
	}
}

func Test_movieQueryBuilder_Update(t *testing.T) {
	var (
		name      = "name"
		url       = "url"
		aliases   = "alias1, alias2"
		director  = "director"
		rating    = 60
		duration  = 34
		synopsis  = "synopsis"
		date, _   = models.ParseDate("2003-02-01")
		createdAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name          string
		updatedObject *models.Movie
		wantErr       bool
	}{
		{
			"full",
			&models.Movie{
				ID:        movieIDs[movieIdxWithTag],
				Name:      name,
				Duration:  &duration,
				Date:      &date,
				Rating:    &rating,
				StudioID:  &studioIDs[studioIdxWithMovie],
				Director:  director,
				Synopsis:  synopsis,
				URLs:      models.NewRelatedStrings([]string{url}),
				TagIDs:    models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithMovie]}),
				Aliases:   aliases,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear tag ids",
			&models.Movie{
				ID:     movieIDs[movieIdxWithTag],
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"invalid studio id",
			&models.Movie{
				ID:       movieIDs[movieIdxWithScene],
				Name:     name,
				StudioID: &invalidID,
			},
			true,
		},
		{
			"invalid tag id",
			&models.Movie{
				ID:     movieIDs[movieIdxWithScene],
				Name:   name,
				TagIDs: models.NewRelatedIDs([]int{invalidID}),
			},
			true,
		},
	}

	qb := db.Movie
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("movieQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.ID)
			if err != nil {
				t.Errorf("movieQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadMovieRelationships(ctx, copy, s); err != nil {
				t.Errorf("loadMovieRelationships() error = %v", err)
				return
			}

			assert.Equal(copy, *s)
		})
	}
}

func clearMoviePartial() models.MoviePartial {
	// leave mandatory fields
	return models.MoviePartial{
		Aliases:  models.OptionalString{Set: true, Null: true},
		Synopsis: models.OptionalString{Set: true, Null: true},
		Director: models.OptionalString{Set: true, Null: true},
		Duration: models.OptionalInt{Set: true, Null: true},
		URLs:     &models.UpdateStrings{Mode: models.RelationshipUpdateModeSet},
		Date:     models.OptionalDate{Set: true, Null: true},
		Rating:   models.OptionalInt{Set: true, Null: true},
		StudioID: models.OptionalInt{Set: true, Null: true},
		TagIDs:   &models.UpdateIDs{Mode: models.RelationshipUpdateModeSet},
	}
}

func Test_movieQueryBuilder_UpdatePartial(t *testing.T) {
	var (
		name      = "name"
		url       = "url"
		aliases   = "alias1, alias2"
		director  = "director"
		rating    = 60
		duration  = 34
		synopsis  = "synopsis"
		date, _   = models.ParseDate("2003-02-01")
		createdAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name    string
		id      int
		partial models.MoviePartial
		want    models.Movie
		wantErr bool
	}{
		{
			"full",
			movieIDs[movieIdxWithScene],
			models.MoviePartial{
				Name:     models.NewOptionalString(name),
				Director: models.NewOptionalString(director),
				Synopsis: models.NewOptionalString(synopsis),
				Aliases:  models.NewOptionalString(aliases),
				URLs: &models.UpdateStrings{
					Values: []string{url},
					Mode:   models.RelationshipUpdateModeSet,
				},
				Date:      models.NewOptionalDate(date),
				Duration:  models.NewOptionalInt(duration),
				Rating:    models.NewOptionalInt(rating),
				StudioID:  models.NewOptionalInt(studioIDs[studioIdxWithMovie]),
				CreatedAt: models.NewOptionalTime(createdAt),
				UpdatedAt: models.NewOptionalTime(updatedAt),
				TagIDs: &models.UpdateIDs{
					IDs:  []int{tagIDs[tagIdx1WithMovie], tagIDs[tagIdx1WithDupName]},
					Mode: models.RelationshipUpdateModeSet,
				},
			},
			models.Movie{
				ID:        movieIDs[movieIdxWithScene],
				Name:      name,
				Director:  director,
				Synopsis:  synopsis,
				Aliases:   aliases,
				URLs:      models.NewRelatedStrings([]string{url}),
				Date:      &date,
				Duration:  &duration,
				Rating:    &rating,
				StudioID:  &studioIDs[studioIdxWithMovie],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				TagIDs:    models.NewRelatedIDs([]int{tagIDs[tagIdx1WithDupName], tagIDs[tagIdx1WithMovie]}),
			},
			false,
		},
		{
			"clear all",
			movieIDs[movieIdxWithScene],
			clearMoviePartial(),
			models.Movie{
				ID:     movieIDs[movieIdxWithScene],
				Name:   movieNames[movieIdxWithScene],
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			false,
		},
		{
			"invalid id",
			invalidID,
			models.MoviePartial{},
			models.Movie{},
			true,
		},
	}
	for _, tt := range tests {
		qb := db.Movie

		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			got, err := qb.UpdatePartial(ctx, tt.id, tt.partial)
			if (err != nil) != tt.wantErr {
				t.Errorf("movieQueryBuilder.UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// load relationships
			if err := loadMovieRelationships(ctx, tt.want, got); err != nil {
				t.Errorf("loadMovieRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *got)

			s, err := qb.Find(ctx, tt.id)
			if err != nil {
				t.Errorf("movieQueryBuilder.Find() error = %v", err)
			}

			// load relationships
			if err := loadMovieRelationships(ctx, tt.want, s); err != nil {
				t.Errorf("loadMovieRelationships() error = %v", err)
				return
			}

			assert.Equal(tt.want, *s)
		})
	}
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

		movies := queryMovies(ctx, t, &filter, nil)
		assert.Len(t, movies, 0, "Expected no movies to be found")

		// query for movies that exclude the URL "ccc"
		urlCriterion.Value = "ccc"
		movies = queryMovies(ctx, t, &filter, nil)

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

		movies := queryMovies(ctx, t, &filter, nil)

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

func queryMovies(ctx context.Context, t *testing.T, movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) []*models.Movie {
	sqb := db.Movie
	movies, _, err := sqb.Query(ctx, movieFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying movie: %s", err.Error())
	}

	return movies
}

func TestMovieQueryTags(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		tagCriterion := models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithMovie]),
				strconv.Itoa(tagIDs[tagIdx1WithMovie]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		movieFilter := models.MovieFilterType{
			Tags: &tagCriterion,
		}

		// ensure ids are correct
		movies := queryMovies(ctx, t, &movieFilter, nil)
		assert.Len(t, movies, 3)
		for _, movie := range movies {
			assert.True(t, movie.ID == movieIDs[movieIdxWithTag] || movie.ID == movieIDs[movieIdxWithTwoTags] || movie.ID == movieIDs[movieIdxWithThreeTags])
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithMovie]),
				strconv.Itoa(tagIDs[tagIdx2WithMovie]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		movies = queryMovies(ctx, t, &movieFilter, nil)

		if assert.Len(t, movies, 2) {
			assert.Equal(t, sceneIDs[movieIdxWithTwoTags], movies[0].ID)
			assert.Equal(t, sceneIDs[movieIdxWithThreeTags], movies[1].ID)
		}

		tagCriterion = models.HierarchicalMultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithMovie]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(movieIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		movies = queryMovies(ctx, t, &movieFilter, &findFilter)
		assert.Len(t, movies, 0)

		return nil
	})
}

func TestMovieQueryTagCount(t *testing.T) {
	const tagCount = 1
	tagCountCriterion := models.IntCriterionInput{
		Value:    tagCount,
		Modifier: models.CriterionModifierEquals,
	}

	verifyMoviesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierNotEquals
	verifyMoviesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyMoviesTagCount(t, tagCountCriterion)

	tagCountCriterion.Modifier = models.CriterionModifierLessThan
	verifyMoviesTagCount(t, tagCountCriterion)
}

func verifyMoviesTagCount(t *testing.T, tagCountCriterion models.IntCriterionInput) {
	withTxn(func(ctx context.Context) error {
		sqb := db.Movie
		movieFilter := models.MovieFilterType{
			TagCount: &tagCountCriterion,
		}

		movies := queryMovies(ctx, t, &movieFilter, nil)
		assert.Greater(t, len(movies), 0)

		for _, movie := range movies {
			ids, err := sqb.GetTagIDs(ctx, movie.ID)
			if err != nil {
				return err
			}
			verifyInt(t, len(ids), tagCountCriterion)
		}

		return nil
	})
}

func TestMovieQuerySorting(t *testing.T) {
	sort := "scenes_count"
	direction := models.SortDirectionEnumDesc
	findFilter := models.FindFilterType{
		Sort:      &sort,
		Direction: &direction,
	}

	withTxn(func(ctx context.Context) error {
		movies := queryMovies(ctx, t, nil, &findFilter)

		// scenes should be in same order as indexes
		firstMovie := movies[0]

		assert.Equal(t, movieIDs[movieIdxWithScene], firstMovie.ID)

		// sort in descending order
		direction = models.SortDirectionEnumAsc

		movies = queryMovies(ctx, t, nil, &findFilter)
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
