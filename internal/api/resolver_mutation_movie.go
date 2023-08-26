package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getMovie(ctx context.Context, id int) (ret *models.Movie, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) MovieCreate(ctx context.Context, input MovieCreateInput) (*models.Movie, error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate a new movie from the input
	currentTime := time.Now()
	newMovie := models.Movie{
		Name:      input.Name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Aliases:   translator.string(input.Aliases, "aliases"),
		Duration:  input.Duration,
		Rating:    translator.ratingConversionInt(input.Rating, input.Rating100),
		Director:  translator.string(input.Director, "director"),
		Synopsis:  translator.string(input.Synopsis, "synopsis"),
		URL:       translator.string(input.URL, "url"),
	}

	var err error

	newMovie.Date, err = translator.datePtr(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	newMovie.StudioID, err = translator.intPtrFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	// HACK: if back image is being set, set the front image to the default.
	// This is because we can't have a null front image with a non-null back image.
	if input.FrontImage == nil && input.BackImage != nil {
		input.FrontImage = &models.DefaultMovieImage
	}

	// Process the base 64 encoded image string
	var frontimageData []byte
	if input.FrontImage != nil {
		frontimageData, err = utils.ProcessImageInput(ctx, *input.FrontImage)
		if err != nil {
			return nil, err
		}
	}

	// Process the base 64 encoded image string
	var backimageData []byte
	if input.BackImage != nil {
		backimageData, err = utils.ProcessImageInput(ctx, *input.BackImage)
		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the movie
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie

		err = qb.Create(ctx, &newMovie)
		if err != nil {
			return err
		}

		// update image table
		if len(frontimageData) > 0 {
			if err := qb.UpdateFrontImage(ctx, newMovie.ID, frontimageData); err != nil {
				return err
			}
		}

		if len(backimageData) > 0 {
			if err := qb.UpdateBackImage(ctx, newMovie.ID, backimageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newMovie.ID, plugin.MovieCreatePost, input, nil)
	return r.getMovie(ctx, newMovie.ID)
}

func (r *mutationResolver) MovieUpdate(ctx context.Context, input MovieUpdateInput) (*models.Movie, error) {
	movieID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate movie from the input
	updatedMovie := models.NewMoviePartial()

	updatedMovie.Name = translator.optionalString(input.Name, "name")
	updatedMovie.Aliases = translator.optionalString(input.Aliases, "aliases")
	updatedMovie.Duration = translator.optionalInt(input.Duration, "duration")
	updatedMovie.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedMovie.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedMovie.Director = translator.optionalString(input.Director, "director")
	updatedMovie.Synopsis = translator.optionalString(input.Synopsis, "synopsis")
	updatedMovie.URL = translator.optionalString(input.URL, "url")
	updatedMovie.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	var frontimageData []byte
	frontImageIncluded := translator.hasField("front_image")
	if input.FrontImage != nil {
		frontimageData, err = utils.ProcessImageInput(ctx, *input.FrontImage)
		if err != nil {
			return nil, err
		}
	}

	var backimageData []byte
	backImageIncluded := translator.hasField("back_image")
	if input.BackImage != nil {
		backimageData, err = utils.ProcessImageInput(ctx, *input.BackImage)
		if err != nil {
			return nil, err
		}
	}

	// Start the transaction and save the movie
	var movie *models.Movie
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie
		movie, err = qb.UpdatePartial(ctx, movieID, updatedMovie)
		if err != nil {
			return err
		}

		// update image table
		if frontImageIncluded {
			if err := qb.UpdateFrontImage(ctx, movie.ID, frontimageData); err != nil {
				return err
			}
		}

		if backImageIncluded {
			if err := qb.UpdateBackImage(ctx, movie.ID, backimageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, movie.ID, plugin.MovieUpdatePost, input, translator.getFields())
	return r.getMovie(ctx, movie.ID)
}

func (r *mutationResolver) BulkMovieUpdate(ctx context.Context, input BulkMovieUpdateInput) ([]*models.Movie, error) {
	movieIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// populate movie from the input
	updatedMovie := models.NewMoviePartial()

	updatedMovie.Rating = translator.ratingConversionOptional(input.Rating, input.Rating100)
	updatedMovie.Director = translator.optionalString(input.Director, "director")
	updatedMovie.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	ret := []*models.Movie{}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie

		for _, movieID := range movieIDs {
			movie, err := qb.UpdatePartial(ctx, movieID, updatedMovie)
			if err != nil {
				return err
			}

			ret = append(ret, movie)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	var newRet []*models.Movie
	for _, movie := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, movie.ID, plugin.MovieUpdatePost, input, translator.getFields())

		movie, err = r.getMovie(ctx, movie.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, movie)
	}

	return newRet, nil
}

func (r *mutationResolver) MovieDestroy(ctx context.Context, input MovieDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Movie.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, id, plugin.MovieDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) MoviesDestroy(ctx context.Context, movieIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(movieIDs)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie
		for _, id := range ids {
			if err := qb.Destroy(ctx, id); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	for _, id := range ids {
		r.hookExecutor.ExecutePostHooks(ctx, id, plugin.MovieDestroyPost, movieIDs, nil)
	}

	return true, nil
}
