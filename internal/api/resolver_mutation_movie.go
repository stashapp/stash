package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
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
	// generate checksum from movie name rather than image
	checksum := md5.FromString(input.Name)

	var frontimageData []byte
	var backimageData []byte
	var err error

	// HACK: if back image is being set, set the front image to the default.
	// This is because we can't have a null front image with a non-null back image.
	if input.FrontImage == nil && input.BackImage != nil {
		input.FrontImage = &models.DefaultMovieImage
	}

	// Process the base 64 encoded image string
	if input.FrontImage != nil {
		frontimageData, err = utils.ProcessImageInput(ctx, *input.FrontImage)
		if err != nil {
			return nil, err
		}
	}

	// Process the base 64 encoded image string
	if input.BackImage != nil {
		backimageData, err = utils.ProcessImageInput(ctx, *input.BackImage)
		if err != nil {
			return nil, err
		}
	}

	// Populate a new movie from the input
	currentTime := time.Now()
	newMovie := models.Movie{
		Checksum:  checksum,
		Name:      sql.NullString{String: input.Name, Valid: true},
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}

	if input.Aliases != nil {
		newMovie.Aliases = sql.NullString{String: *input.Aliases, Valid: true}
	}
	if input.Duration != nil {
		duration := int64(*input.Duration)
		newMovie.Duration = sql.NullInt64{Int64: duration, Valid: true}
	}

	if input.Date != nil {
		newMovie.Date = models.SQLiteDate{String: *input.Date, Valid: true}
	}

	if input.Rating != nil {
		rating := int64(*input.Rating)
		newMovie.Rating = sql.NullInt64{Int64: rating, Valid: true}
	}

	if input.StudioID != nil {
		studioID, _ := strconv.ParseInt(*input.StudioID, 10, 64)
		newMovie.StudioID = sql.NullInt64{Int64: studioID, Valid: true}
	}

	if input.Director != nil {
		newMovie.Director = sql.NullString{String: *input.Director, Valid: true}
	}

	if input.Synopsis != nil {
		newMovie.Synopsis = sql.NullString{String: *input.Synopsis, Valid: true}
	}

	if input.URL != nil {
		newMovie.URL = sql.NullString{String: *input.URL, Valid: true}
	}

	// Start the transaction and save the movie
	var movie *models.Movie
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie
		movie, err = qb.Create(ctx, newMovie)
		if err != nil {
			return err
		}

		// update image table
		if len(frontimageData) > 0 {
			if err := qb.UpdateImages(ctx, movie.ID, frontimageData, backimageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, movie.ID, plugin.MovieCreatePost, input, nil)
	return r.getMovie(ctx, movie.ID)
}

func (r *mutationResolver) MovieUpdate(ctx context.Context, input MovieUpdateInput) (*models.Movie, error) {
	// Populate movie from the input
	movieID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	updatedMovie := models.MoviePartial{
		ID:        movieID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	var frontimageData []byte
	frontImageIncluded := translator.hasField("front_image")
	if input.FrontImage != nil {
		frontimageData, err = utils.ProcessImageInput(ctx, *input.FrontImage)
		if err != nil {
			return nil, err
		}
	}
	backImageIncluded := translator.hasField("back_image")
	var backimageData []byte
	if input.BackImage != nil {
		backimageData, err = utils.ProcessImageInput(ctx, *input.BackImage)
		if err != nil {
			return nil, err
		}
	}

	if input.Name != nil {
		// generate checksum from movie name rather than image
		checksum := md5.FromString(*input.Name)
		updatedMovie.Name = &sql.NullString{String: *input.Name, Valid: true}
		updatedMovie.Checksum = &checksum
	}

	updatedMovie.Aliases = translator.nullString(input.Aliases, "aliases")
	updatedMovie.Duration = translator.nullInt64(input.Duration, "duration")
	updatedMovie.Date = translator.sqliteDate(input.Date, "date")
	updatedMovie.Rating = translator.nullInt64(input.Rating, "rating")
	updatedMovie.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedMovie.Director = translator.nullString(input.Director, "director")
	updatedMovie.Synopsis = translator.nullString(input.Synopsis, "synopsis")
	updatedMovie.URL = translator.nullString(input.URL, "url")

	// Start the transaction and save the movie
	var movie *models.Movie
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie
		movie, err = qb.Update(ctx, updatedMovie)
		if err != nil {
			return err
		}

		// update image table
		if frontImageIncluded || backImageIncluded {
			if !frontImageIncluded {
				frontimageData, err = qb.GetFrontImage(ctx, updatedMovie.ID)
				if err != nil {
					return err
				}
			}
			if !backImageIncluded {
				backimageData, err = qb.GetBackImage(ctx, updatedMovie.ID)
				if err != nil {
					return err
				}
			}

			if len(frontimageData) == 0 && len(backimageData) == 0 {
				// both images are being nulled. Destroy them.
				if err := qb.DestroyImages(ctx, movie.ID); err != nil {
					return err
				}
			} else {
				// HACK - if front image is null and back image is not null, then set the front image
				// to the default image since we can't have a null front image and a non-null back image
				if frontimageData == nil && backimageData != nil {
					frontimageData, _ = utils.ProcessImageInput(ctx, models.DefaultMovieImage)
				}

				if err := qb.UpdateImages(ctx, movie.ID, frontimageData, backimageData); err != nil {
					return err
				}
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

	updatedTime := time.Now()

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedMovie := models.MoviePartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedMovie.Rating = translator.nullInt64(input.Rating, "rating")
	updatedMovie.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedMovie.Director = translator.nullString(input.Director, "director")

	ret := []*models.Movie{}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie

		for _, movieID := range movieIDs {
			updatedMovie.ID = movieID

			existing, err := qb.Find(ctx, movieID)
			if err != nil {
				return err
			}

			if existing == nil {
				return fmt.Errorf("movie with id %d not found", movieID)
			}

			movie, err := qb.Update(ctx, updatedMovie)
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
