package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func movieFromGroupCreateInput(ctx context.Context, input GroupCreateInput) (*models.Movie, error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate a new movie from the input
	newMovie := models.NewMovie()

	newMovie.Name = input.Name
	newMovie.Aliases = translator.string(input.Aliases)
	newMovie.Duration = input.Duration
	newMovie.Rating = input.Rating100
	newMovie.Director = translator.string(input.Director)
	newMovie.Synopsis = translator.string(input.Synopsis)

	var err error

	newMovie.Date, err = translator.datePtr(input.Date)
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	newMovie.StudioID, err = translator.intPtrFromString(input.StudioID)
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	newMovie.TagIDs, err = translator.relatedIds(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	if input.Urls != nil {
		newMovie.URLs = models.NewRelatedStrings(input.Urls)
	}

	return &newMovie, nil
}

func (r *mutationResolver) GroupCreate(ctx context.Context, input GroupCreateInput) (*models.Movie, error) {
	newMovie, err := movieFromGroupCreateInput(ctx, input)
	if err != nil {
		return nil, err
	}

	// Process the base 64 encoded image string
	var frontimageData []byte
	if input.FrontImage != nil {
		frontimageData, err = utils.ProcessImageInput(ctx, *input.FrontImage)
		if err != nil {
			return nil, fmt.Errorf("processing front image: %w", err)
		}
	}

	// Process the base 64 encoded image string
	var backimageData []byte
	if input.BackImage != nil {
		backimageData, err = utils.ProcessImageInput(ctx, *input.BackImage)
		if err != nil {
			return nil, fmt.Errorf("processing back image: %w", err)
		}
	}

	// HACK: if back image is being set, set the front image to the default.
	// This is because we can't have a null front image with a non-null back image.
	if len(frontimageData) == 0 && len(backimageData) != 0 {
		frontimageData = static.ReadAll(static.DefaultMovieImage)
	}

	// Start the transaction and save the movie
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Movie

		err = qb.Create(ctx, newMovie)
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

	// for backwards compatibility - run both movie and group hooks
	r.hookExecutor.ExecutePostHooks(ctx, newMovie.ID, hook.GroupCreatePost, input, nil)
	r.hookExecutor.ExecutePostHooks(ctx, newMovie.ID, hook.MovieCreatePost, input, nil)
	return r.getMovie(ctx, newMovie.ID)
}

func moviePartialFromGroupUpdateInput(translator changesetTranslator, input GroupUpdateInput) (ret models.MoviePartial, err error) {
	// Populate movie from the input
	updatedMovie := models.NewMoviePartial()

	updatedMovie.Name = translator.optionalString(input.Name, "name")
	updatedMovie.Aliases = translator.optionalString(input.Aliases, "aliases")
	updatedMovie.Duration = translator.optionalInt(input.Duration, "duration")
	updatedMovie.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedMovie.Director = translator.optionalString(input.Director, "director")
	updatedMovie.Synopsis = translator.optionalString(input.Synopsis, "synopsis")

	updatedMovie.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		err = fmt.Errorf("converting date: %w", err)
		return
	}
	updatedMovie.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		err = fmt.Errorf("converting studio id: %w", err)
		return
	}

	updatedMovie.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		err = fmt.Errorf("converting tag ids: %w", err)
		return
	}

	updatedMovie.URLs = translator.updateStrings(input.Urls, "urls")

	return updatedMovie, nil
}

func (r *mutationResolver) GroupUpdate(ctx context.Context, input GroupUpdateInput) (*models.Movie, error) {
	movieID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedMovie, err := moviePartialFromGroupUpdateInput(translator, input)
	if err != nil {
		return nil, err
	}

	var frontimageData []byte
	frontImageIncluded := translator.hasField("front_image")
	if input.FrontImage != nil {
		frontimageData, err = utils.ProcessImageInput(ctx, *input.FrontImage)
		if err != nil {
			return nil, fmt.Errorf("processing front image: %w", err)
		}
	}

	var backimageData []byte
	backImageIncluded := translator.hasField("back_image")
	if input.BackImage != nil {
		backimageData, err = utils.ProcessImageInput(ctx, *input.BackImage)
		if err != nil {
			return nil, fmt.Errorf("processing back image: %w", err)
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

	// for backwards compatibility - run both movie and group hooks
	r.hookExecutor.ExecutePostHooks(ctx, movie.ID, hook.GroupUpdatePost, input, translator.getFields())
	r.hookExecutor.ExecutePostHooks(ctx, movie.ID, hook.MovieUpdatePost, input, translator.getFields())
	return r.getMovie(ctx, movie.ID)
}

func moviePartialFromBulkGroupUpdateInput(translator changesetTranslator, input BulkGroupUpdateInput) (ret models.MoviePartial, err error) {
	updatedMovie := models.NewMoviePartial()

	updatedMovie.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedMovie.Director = translator.optionalString(input.Director, "director")

	updatedMovie.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		err = fmt.Errorf("converting studio id: %w", err)
		return
	}

	updatedMovie.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		err = fmt.Errorf("converting tag ids: %w", err)
		return
	}

	updatedMovie.URLs = translator.optionalURLsBulk(input.Urls, nil)

	return updatedMovie, nil
}

func (r *mutationResolver) BulkGroupUpdate(ctx context.Context, input BulkGroupUpdateInput) ([]*models.Movie, error) {
	movieIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate movie from the input
	updatedMovie, err := moviePartialFromBulkGroupUpdateInput(translator, input)
	if err != nil {
		return nil, err
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
		// for backwards compatibility - run both movie and group hooks
		r.hookExecutor.ExecutePostHooks(ctx, movie.ID, hook.GroupUpdatePost, input, translator.getFields())
		r.hookExecutor.ExecutePostHooks(ctx, movie.ID, hook.MovieUpdatePost, input, translator.getFields())

		movie, err = r.getMovie(ctx, movie.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, movie)
	}

	return newRet, nil
}

func (r *mutationResolver) GroupDestroy(ctx context.Context, input GroupDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Movie.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	// for backwards compatibility - run both movie and group hooks
	r.hookExecutor.ExecutePostHooks(ctx, id, hook.GroupDestroyPost, input, nil)
	r.hookExecutor.ExecutePostHooks(ctx, id, hook.MovieDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) GroupsDestroy(ctx context.Context, groupIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(groupIDs)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
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
		// for backwards compatibility - run both movie and group hooks
		r.hookExecutor.ExecutePostHooks(ctx, id, hook.GroupDestroyPost, groupIDs, nil)
		r.hookExecutor.ExecutePostHooks(ctx, id, hook.MovieDestroyPost, groupIDs, nil)
	}

	return true, nil
}
