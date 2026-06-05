package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) StoryCreate(ctx context.Context, input models.StoryCreateInput) (*models.Story, error) {
	newStory := models.NewStory()
	newStory.Title = input.Title
	if input.Author != nil { newStory.Author = *input.Author }
	if input.URL != nil { newStory.URLs.Add(*input.URL) }
	if input.Language != nil { newStory.Language = *input.Language }
	if input.TagLine != nil { newStory.TagLine = *input.TagLine }
	if input.Details != nil { newStory.Details = *input.Details }
	if input.Rating100 != nil { newStory.Rating = input.Rating100 }

	var err error
	newStory.Date, err = utils.ParseDate(input.Date)
	if err != nil { return nil, fmt.Errorf("parsing date: %w", err) }

	imageData := make(map[string][]byte)
	for _, field := range []struct{ key string; val *string }{
		{"front_image", input.FrontImage},
		{"back_image", input.BackImage},
	} {
		if field.val != nil {
			data, err := utils.ProcessImageInput(ctx, *field.val)
			if err != nil { return nil, fmt.Errorf("processing %s: %w", field.key, err) }
			imageData[field.key] = data
		}
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Story
		if err := qb.Create(ctx, &models.CreateStoryInput{Story: &newStory}); err != nil {
			return err
		}
		if data, ok := imageData["front_image"]; ok && len(data) > 0 {
			if err := qb.UpdateFrontImage(ctx, newStory.ID, data); err != nil { return err }
		}
		if data, ok := imageData["back_image"]; ok && len(data) > 0 {
			if err := qb.UpdateBackImage(ctx, newStory.ID, data); err != nil { return err }
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return r.getStory(ctx, newStory.ID)
}

func (r *mutationResolver) StoryUpdate(ctx context.Context, input models.StoryUpdateInput) (*models.Story, error) {
	storyID, err := strconv.Atoi(input.ID)
	if err != nil { return nil, fmt.Errorf("converting id: %w", err) }

	var imageData []byte
	if input.FrontImage != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.FrontImage)
		if err != nil { return nil, fmt.Errorf("processing front image: %w", err) }
	}

	var backImageData []byte
	if input.BackImage != nil {
		backImageData, err = utils.ProcessImageInput(ctx, *input.BackImage)
		if err != nil { return nil, fmt.Errorf("processing back image: %w", err) }
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Story
		partial := models.NewStoryPartial()
		if input.Title != nil { partial.Title = models.NewOptionalString(*input.Title) }
		if input.Author != nil { partial.Author = models.NewOptionalString(*input.Author) }
		if input.URL != nil { partial.URLs = &models.UpdateStrings{Values: []string{*input.URL}, Mode: models.RelationshipUpdateModeSet} }
		if input.Language != nil { partial.Language = models.NewOptionalString(*input.Language) }
		if input.TagLine != nil { partial.TagLine = models.NewOptionalString(*input.TagLine) }
		if input.Details != nil { partial.Details = models.NewOptionalString(*input.Details) }
		if input.Rating100 != nil { partial.Rating = models.NewOptionalInt(*input.Rating100) }
		if input.StudioID != nil {
			sid, err := strconv.Atoi(*input.StudioID)
			if err != nil { return fmt.Errorf("converting studio id: %w", err) }
			partial.StudioID = models.NewOptionalInt(sid)
		}

		_, err = qb.UpdatePartial(ctx, storyID, partial)
		if err != nil { return err }

		if len(imageData) > 0 {
			if err := qb.UpdateFrontImage(ctx, storyID, imageData); err != nil { return err }
		}
		if len(backImageData) > 0 {
			if err := qb.UpdateBackImage(ctx, storyID, backImageData); err != nil { return err }
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return r.getStory(ctx, storyID)
}

func (r *mutationResolver) BulkStoryUpdate(ctx context.Context, input BulkStoryUpdateInput) ([]*models.Story, error) {
	storyIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil { return nil, fmt.Errorf("converting ids: %w", err) }

	var stories []*models.Story
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Story
		for _, id := range storyIDs {
			partial := models.NewStoryPartial()
			if input.URL != nil { partial.URLs = &models.UpdateStrings{Values: []string{*input.URL}, Mode: models.RelationshipUpdateModeSet} }
			if input.Language != nil { partial.Language = models.NewOptionalString(*input.Language) }
			if input.TagLine != nil { partial.TagLine = models.NewOptionalString(*input.TagLine) }
			if input.Details != nil { partial.Details = models.NewOptionalString(*input.Details) }
			if input.Rating100 != nil { partial.Rating = models.NewOptionalInt(*input.Rating100) }
			updated, err := qb.UpdatePartial(ctx, id, partial)
			if err != nil { return err }
			stories = append(stories, updated)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return stories, nil
}

func (r *mutationResolver) StoryDestroy(ctx context.Context, input StoryDestroyInput) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil { return false, fmt.Errorf("converting ids: %w", err) }

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for _, id := range ids {
			if err := r.repository.Story.Destroy(ctx, id); err != nil { return err }
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) getStory(ctx context.Context, id int) (*models.Story, error) {
	var ret *models.Story
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Story.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}
