package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

// used to refetch group after hooks run
func (r *mutationResolver) getGroup(ctx context.Context, id int) (ret *models.Group, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Group.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) MovieCreate(ctx context.Context, input MovieCreateInput) (*models.Group, error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate a new group from the input
	newGroup := models.NewGroup()

	newGroup.Name = strings.TrimSpace(input.Name)
	newGroup.Aliases = translator.string(input.Aliases)
	newGroup.Duration = input.Duration
	newGroup.Rating = input.Rating100
	newGroup.Director = translator.string(input.Director)
	newGroup.Synopsis = translator.string(input.Synopsis)

	var err error

	newGroup.Date, err = translator.datePtr(input.Date)
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	newGroup.StudioID, err = translator.intPtrFromString(input.StudioID)
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	newGroup.TagIDs, err = translator.relatedIds(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	if input.Urls != nil {
		newGroup.URLs = models.NewRelatedStrings(stringslice.TrimSpace(input.Urls))
	} else if input.URL != nil {
		newGroup.URLs = models.NewRelatedStrings([]string{strings.TrimSpace(*input.URL)})
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
		frontimageData = static.ReadAll(static.DefaultGroupImage)
	}

	// Start the transaction and save the group
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Group

		err = qb.Create(ctx, &newGroup)
		if err != nil {
			return err
		}

		// update image table
		if len(frontimageData) > 0 {
			if err := qb.UpdateFrontImage(ctx, newGroup.ID, frontimageData); err != nil {
				return err
			}
		}

		if len(backimageData) > 0 {
			if err := qb.UpdateBackImage(ctx, newGroup.ID, backimageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// for backwards compatibility - run both movie and group hooks
	r.hookExecutor.ExecutePostHooks(ctx, newGroup.ID, hook.GroupCreatePost, input, nil)
	r.hookExecutor.ExecutePostHooks(ctx, newGroup.ID, hook.MovieCreatePost, input, nil)
	return r.getGroup(ctx, newGroup.ID)
}

func (r *mutationResolver) MovieUpdate(ctx context.Context, input MovieUpdateInput) (*models.Group, error) {
	groupID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate group from the input
	updatedGroup := models.NewGroupPartial()

	updatedGroup.Name = translator.optionalString(input.Name, "name")
	updatedGroup.Aliases = translator.optionalString(input.Aliases, "aliases")
	updatedGroup.Duration = translator.optionalInt(input.Duration, "duration")
	updatedGroup.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedGroup.Director = translator.optionalString(input.Director, "director")
	updatedGroup.Synopsis = translator.optionalString(input.Synopsis, "synopsis")

	updatedGroup.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedGroup.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedGroup.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	updatedGroup.URLs = translator.optionalURLs(input.Urls, input.URL)

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

	// Start the transaction and save the group
	var group *models.Group
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Group
		group, err = qb.UpdatePartial(ctx, groupID, updatedGroup)
		if err != nil {
			return err
		}

		// update image table
		if frontImageIncluded {
			if err := qb.UpdateFrontImage(ctx, group.ID, frontimageData); err != nil {
				return err
			}
		}

		if backImageIncluded {
			if err := qb.UpdateBackImage(ctx, group.ID, backimageData); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// for backwards compatibility - run both movie and group hooks
	r.hookExecutor.ExecutePostHooks(ctx, group.ID, hook.GroupUpdatePost, input, translator.getFields())
	r.hookExecutor.ExecutePostHooks(ctx, group.ID, hook.MovieUpdatePost, input, translator.getFields())
	return r.getGroup(ctx, group.ID)
}

func (r *mutationResolver) BulkMovieUpdate(ctx context.Context, input BulkMovieUpdateInput) ([]*models.Group, error) {
	groupIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate group from the input
	updatedGroup := models.NewGroupPartial()

	updatedGroup.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedGroup.Director = translator.optionalString(input.Director, "director")

	updatedGroup.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedGroup.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	updatedGroup.URLs = translator.optionalURLsBulk(input.Urls, nil)

	ret := []*models.Group{}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Group

		for _, groupID := range groupIDs {
			group, err := qb.UpdatePartial(ctx, groupID, updatedGroup)
			if err != nil {
				return err
			}

			ret = append(ret, group)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	var newRet []*models.Group
	for _, group := range ret {
		// for backwards compatibility - run both movie and group hooks
		r.hookExecutor.ExecutePostHooks(ctx, group.ID, hook.GroupUpdatePost, input, translator.getFields())
		r.hookExecutor.ExecutePostHooks(ctx, group.ID, hook.MovieUpdatePost, input, translator.getFields())

		group, err = r.getGroup(ctx, group.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, group)
	}

	return newRet, nil
}

func (r *mutationResolver) MovieDestroy(ctx context.Context, input MovieDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Group.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	// for backwards compatibility - run both movie and group hooks
	r.hookExecutor.ExecutePostHooks(ctx, id, hook.GroupDestroyPost, input, nil)
	r.hookExecutor.ExecutePostHooks(ctx, id, hook.MovieDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) MoviesDestroy(ctx context.Context, groupIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(groupIDs)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Group
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
