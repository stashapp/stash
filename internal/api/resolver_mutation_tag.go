package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getTag(ctx context.Context, id int) (ret *models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) TagCreate(ctx context.Context, input TagCreateInput) (*models.Tag, error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate a new tag from the input
	newTag := models.NewTag()

	newTag.Name = input.Name
	newTag.Description = translator.string(input.Description)
	newTag.IgnoreAutoTag = translator.bool(input.IgnoreAutoTag)

	var err error

	var parentIDs []int
	if len(input.ParentIds) > 0 {
		parentIDs, err = stringslice.StringSliceToIntSlice(input.ParentIds)
		if err != nil {
			return nil, fmt.Errorf("converting parent ids: %w", err)
		}
	}

	var childIDs []int
	if len(input.ChildIds) > 0 {
		childIDs, err = stringslice.StringSliceToIntSlice(input.ChildIds)
		if err != nil {
			return nil, fmt.Errorf("converting child ids: %w", err)
		}
	}

	// Process the base 64 encoded image string
	var imageData []byte
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, fmt.Errorf("processing image: %w", err)
		}
	}

	// Start the transaction and save the tag
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Tag

		// ensure name is unique
		if err := tag.EnsureTagNameUnique(ctx, 0, newTag.Name, qb); err != nil {
			return err
		}

		err = qb.Create(ctx, &newTag)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, newTag.ID, imageData); err != nil {
				return err
			}
		}

		if len(input.Aliases) > 0 {
			if err := tag.EnsureAliasesUnique(ctx, newTag.ID, input.Aliases, qb); err != nil {
				return err
			}

			if err := qb.UpdateAliases(ctx, newTag.ID, input.Aliases); err != nil {
				return err
			}
		}

		if len(parentIDs) > 0 {
			if err := qb.UpdateParentTags(ctx, newTag.ID, parentIDs); err != nil {
				return err
			}
		}

		if len(childIDs) > 0 {
			if err := qb.UpdateChildTags(ctx, newTag.ID, childIDs); err != nil {
				return err
			}
		}

		// FIXME: This should be called before any changes are made, but
		// requires a rewrite of ValidateHierarchy.
		if len(parentIDs) > 0 || len(childIDs) > 0 {
			if err := tag.ValidateHierarchy(ctx, &newTag, parentIDs, childIDs, qb); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newTag.ID, hook.TagCreatePost, input, nil)
	return r.getTag(ctx, newTag.ID)
}

func (r *mutationResolver) TagUpdate(ctx context.Context, input TagUpdateInput) (*models.Tag, error) {
	tagID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate tag from the input
	updatedTag := models.NewTagPartial()

	updatedTag.IgnoreAutoTag = translator.optionalBool(input.IgnoreAutoTag, "ignore_auto_tag")
	updatedTag.Description = translator.optionalString(input.Description, "description")

	var parentIDs []int
	if translator.hasField("parent_ids") {
		parentIDs, err = stringslice.StringSliceToIntSlice(input.ParentIds)
		if err != nil {
			return nil, fmt.Errorf("converting parent ids: %w", err)
		}
	}

	var childIDs []int
	if translator.hasField("child_ids") {
		childIDs, err = stringslice.StringSliceToIntSlice(input.ChildIds)
		if err != nil {
			return nil, fmt.Errorf("converting child ids: %w", err)
		}
	}

	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, fmt.Errorf("processing image: %w", err)
		}
	}

	// Start the transaction and save the tag
	var t *models.Tag
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Tag

		// ensure name is unique
		t, err = qb.Find(ctx, tagID)
		if err != nil {
			return err
		}

		if t == nil {
			return fmt.Errorf("tag with id %d not found", tagID)
		}

		if input.Name != nil && t.Name != *input.Name {
			if err := tag.EnsureTagNameUnique(ctx, tagID, *input.Name, qb); err != nil {
				return err
			}

			updatedTag.Name = models.NewOptionalString(*input.Name)
		}

		t, err = qb.UpdatePartial(ctx, tagID, updatedTag)
		if err != nil {
			return err
		}

		// update image table
		if imageIncluded {
			if err := qb.UpdateImage(ctx, tagID, imageData); err != nil {
				return err
			}
		}

		if translator.hasField("aliases") {
			if err := tag.EnsureAliasesUnique(ctx, tagID, input.Aliases, qb); err != nil {
				return err
			}

			if err := qb.UpdateAliases(ctx, tagID, input.Aliases); err != nil {
				return err
			}
		}

		if parentIDs != nil {
			if err := qb.UpdateParentTags(ctx, tagID, parentIDs); err != nil {
				return err
			}
		}

		if childIDs != nil {
			if err := qb.UpdateChildTags(ctx, tagID, childIDs); err != nil {
				return err
			}
		}

		// FIXME: This should be called before any changes are made, but
		// requires a rewrite of ValidateHierarchy.
		if parentIDs != nil || childIDs != nil {
			if err := tag.ValidateHierarchy(ctx, t, parentIDs, childIDs, qb); err != nil {
				logger.Errorf("Error saving tag: %s", err)
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, t.ID, hook.TagUpdatePost, input, translator.getFields())
	return r.getTag(ctx, t.ID)
}

func (r *mutationResolver) TagDestroy(ctx context.Context, input TagDestroyInput) (bool, error) {
	tagID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Tag.Destroy(ctx, tagID)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, tagID, hook.TagDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) TagsDestroy(ctx context.Context, tagIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(tagIDs)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Tag
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
		r.hookExecutor.ExecutePostHooks(ctx, id, hook.TagDestroyPost, tagIDs, nil)
	}

	return true, nil
}

func (r *mutationResolver) TagsMerge(ctx context.Context, input TagsMergeInput) (*models.Tag, error) {
	source, err := stringslice.StringSliceToIntSlice(input.Source)
	if err != nil {
		return nil, fmt.Errorf("converting source ids: %w", err)
	}

	destination, err := strconv.Atoi(input.Destination)
	if err != nil {
		return nil, fmt.Errorf("converting destination id: %w", err)
	}

	if len(source) == 0 {
		return nil, nil
	}

	var t *models.Tag
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Tag

		var err error
		t, err = qb.Find(ctx, destination)
		if err != nil {
			return err
		}

		if t == nil {
			return fmt.Errorf("tag with id %d not found", destination)
		}

		parents, children, err := tag.MergeHierarchy(ctx, destination, source, qb)
		if err != nil {
			return err
		}

		if err = qb.Merge(ctx, source, destination); err != nil {
			return err
		}

		err = qb.UpdateParentTags(ctx, destination, parents)
		if err != nil {
			return err
		}
		err = qb.UpdateChildTags(ctx, destination, children)
		if err != nil {
			return err
		}

		err = tag.ValidateHierarchy(ctx, t, parents, children, qb)
		if err != nil {
			logger.Errorf("Error merging tag: %s", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, t.ID, hook.TagMergePost, input, nil)

	return t, nil
}
