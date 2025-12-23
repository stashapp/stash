package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/group"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func groupFromGroupCreateInput(ctx context.Context, input GroupCreateInput) (*models.Group, error) {
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

	newGroup.ContainingGroups, err = translator.groupIDDescriptions(input.ContainingGroups)
	if err != nil {
		return nil, fmt.Errorf("converting containing group ids: %w", err)
	}

	newGroup.SubGroups, err = translator.groupIDDescriptions(input.SubGroups)
	if err != nil {
		return nil, fmt.Errorf("converting containing group ids: %w", err)
	}

	if input.Urls != nil {
		newGroup.URLs = models.NewRelatedStrings(stringslice.TrimSpace(input.Urls))
	}

	return &newGroup, nil
}

func (r *mutationResolver) GroupCreate(ctx context.Context, input GroupCreateInput) (*models.Group, error) {
	newGroup, err := groupFromGroupCreateInput(ctx, input)
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
		frontimageData = static.ReadAll(static.DefaultGroupImage)
	}

	// Start the transaction and save the group
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		if err = r.groupService.Create(ctx, newGroup, frontimageData, backimageData); err != nil {
			return err
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

func groupPartialFromGroupUpdateInput(translator changesetTranslator, input GroupUpdateInput) (ret models.GroupPartial, err error) {
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
		err = fmt.Errorf("converting date: %w", err)
		return
	}
	updatedGroup.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		err = fmt.Errorf("converting studio id: %w", err)
		return
	}

	updatedGroup.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		err = fmt.Errorf("converting tag ids: %w", err)
		return
	}

	updatedGroup.ContainingGroups, err = translator.updateGroupIDDescriptions(input.ContainingGroups, "containing_groups")
	if err != nil {
		err = fmt.Errorf("converting containing group ids: %w", err)
		return
	}

	updatedGroup.SubGroups, err = translator.updateGroupIDDescriptions(input.SubGroups, "sub_groups")
	if err != nil {
		err = fmt.Errorf("converting containing group ids: %w", err)
		return
	}

	updatedGroup.URLs = translator.updateStrings(input.Urls, "urls")

	return updatedGroup, nil
}

func (r *mutationResolver) GroupUpdate(ctx context.Context, input GroupUpdateInput) (*models.Group, error) {
	groupID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedGroup, err := groupPartialFromGroupUpdateInput(translator, input)
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

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		frontImage := group.ImageInput{
			Image: frontimageData,
			Set:   frontImageIncluded,
		}

		backImage := group.ImageInput{
			Image: backimageData,
			Set:   backImageIncluded,
		}

		_, err = r.groupService.UpdatePartial(ctx, groupID, updatedGroup, frontImage, backImage)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// for backwards compatibility - run both movie and group hooks
	r.hookExecutor.ExecutePostHooks(ctx, groupID, hook.GroupUpdatePost, input, translator.getFields())
	r.hookExecutor.ExecutePostHooks(ctx, groupID, hook.MovieUpdatePost, input, translator.getFields())
	return r.getGroup(ctx, groupID)
}

func groupPartialFromBulkGroupUpdateInput(translator changesetTranslator, input BulkGroupUpdateInput) (ret models.GroupPartial, err error) {
	updatedGroup := models.NewGroupPartial()

	updatedGroup.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedGroup.Director = translator.optionalString(input.Director, "director")

	updatedGroup.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		err = fmt.Errorf("converting studio id: %w", err)
		return
	}

	updatedGroup.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		err = fmt.Errorf("converting tag ids: %w", err)
		return
	}

	updatedGroup.ContainingGroups, err = translator.updateGroupIDDescriptionsBulk(input.ContainingGroups, "containing_groups")
	if err != nil {
		err = fmt.Errorf("converting containing group ids: %w", err)
		return
	}

	updatedGroup.SubGroups, err = translator.updateGroupIDDescriptionsBulk(input.SubGroups, "sub_groups")
	if err != nil {
		err = fmt.Errorf("converting containing group ids: %w", err)
		return
	}

	updatedGroup.URLs = translator.optionalURLsBulk(input.Urls, nil)

	return updatedGroup, nil
}

func (r *mutationResolver) BulkGroupUpdate(ctx context.Context, input BulkGroupUpdateInput) ([]*models.Group, error) {
	groupIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate group from the input
	updatedGroup, err := groupPartialFromBulkGroupUpdateInput(translator, input)
	if err != nil {
		return nil, err
	}

	ret := []*models.Group{}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for _, groupID := range groupIDs {
			group, err := r.groupService.UpdatePartial(ctx, groupID, updatedGroup, group.ImageInput{}, group.ImageInput{})
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

func (r *mutationResolver) GroupDestroy(ctx context.Context, input GroupDestroyInput) (bool, error) {
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

func (r *mutationResolver) GroupsDestroy(ctx context.Context, groupIDs []string) (bool, error) {
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

func (r *mutationResolver) AddGroupSubGroups(ctx context.Context, input GroupSubGroupAddInput) (bool, error) {
	groupID, err := strconv.Atoi(input.ContainingGroupID)
	if err != nil {
		return false, fmt.Errorf("converting group id: %w", err)
	}

	subGroups, err := groupsDescriptionsFromGroupInput(input.SubGroups)
	if err != nil {
		return false, fmt.Errorf("converting sub group ids: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.groupService.AddSubGroups(ctx, groupID, subGroups, input.InsertIndex)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveGroupSubGroups(ctx context.Context, input GroupSubGroupRemoveInput) (bool, error) {
	groupID, err := strconv.Atoi(input.ContainingGroupID)
	if err != nil {
		return false, fmt.Errorf("converting group id: %w", err)
	}

	subGroupIDs, err := stringslice.StringSliceToIntSlice(input.SubGroupIds)
	if err != nil {
		return false, fmt.Errorf("converting sub group ids: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.groupService.RemoveSubGroups(ctx, groupID, subGroupIDs)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ReorderSubGroups(ctx context.Context, input ReorderSubGroupsInput) (bool, error) {
	groupID, err := strconv.Atoi(input.GroupID)
	if err != nil {
		return false, fmt.Errorf("converting group id: %w", err)
	}

	subGroupIDs, err := stringslice.StringSliceToIntSlice(input.SubGroupIds)
	if err != nil {
		return false, fmt.Errorf("converting sub group ids: %w", err)
	}

	insertPointID, err := strconv.Atoi(input.InsertAtID)
	if err != nil {
		return false, fmt.Errorf("converting insert at id: %w", err)
	}

	insertAfter := utils.IsTrue(input.InsertAfter)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.groupService.ReorderSubGroups(ctx, groupID, subGroupIDs, insertPointID, insertAfter)
	}); err != nil {
		return false, err
	}

	return true, nil
}
