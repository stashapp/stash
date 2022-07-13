package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getStudio(ctx context.Context, id int) (ret *models.Studio, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) StudioCreate(ctx context.Context, input StudioCreateInput) (*models.Studio, error) {
	// generate checksum from studio name rather than image
	checksum := md5.FromString(input.Name)

	var imageData []byte
	var err error

	// Process the base 64 encoded image string
	if input.Image != nil {
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	// Populate a new studio from the input
	currentTime := time.Now()
	newStudio := models.Studio{
		Checksum:  checksum,
		Name:      sql.NullString{String: input.Name, Valid: true},
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}
	if input.URL != nil {
		newStudio.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.ParentID != nil {
		parentID, _ := strconv.ParseInt(*input.ParentID, 10, 64)
		newStudio.ParentID = sql.NullInt64{Int64: parentID, Valid: true}
	}

	if input.Rating != nil {
		newStudio.Rating = sql.NullInt64{Int64: int64(*input.Rating), Valid: true}
	} else {
		newStudio.Rating = sql.NullInt64{Valid: false}
	}
	if input.Details != nil {
		newStudio.Details = sql.NullString{String: *input.Details, Valid: true}
	}
	if input.IgnoreAutoTag != nil {
		newStudio.IgnoreAutoTag = *input.IgnoreAutoTag
	}

	// Start the transaction and save the studio
	var s *models.Studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		var err error
		s, err = qb.Create(ctx, newStudio)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, s.ID, imageData); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if input.StashIds != nil {
			stashIDJoins := input.StashIds
			if err := qb.UpdateStashIDs(ctx, s.ID, stashIDJoins); err != nil {
				return err
			}
		}

		if len(input.Aliases) > 0 {
			if err := studio.EnsureAliasesUnique(ctx, s.ID, input.Aliases, qb); err != nil {
				return err
			}

			if err := qb.UpdateAliases(ctx, s.ID, input.Aliases); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, s.ID, plugin.StudioCreatePost, input, nil)
	return r.getStudio(ctx, s.ID)
}

func (r *mutationResolver) StudioUpdate(ctx context.Context, input StudioUpdateInput) (*models.Studio, error) {
	// Populate studio from the input
	studioID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedStudio := models.StudioPartial{
		ID:        studioID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	var imageData []byte
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		var err error
		imageData, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}
	if input.Name != nil {
		// generate checksum from studio name rather than image
		checksum := md5.FromString(*input.Name)
		updatedStudio.Name = &sql.NullString{String: *input.Name, Valid: true}
		updatedStudio.Checksum = &checksum
	}

	updatedStudio.URL = translator.nullString(input.URL, "url")
	updatedStudio.Details = translator.nullString(input.Details, "details")
	updatedStudio.ParentID = translator.nullInt64FromString(input.ParentID, "parent_id")
	updatedStudio.Rating = translator.nullInt64(input.Rating, "rating")
	updatedStudio.IgnoreAutoTag = input.IgnoreAutoTag

	// Start the transaction and save the studio
	var s *models.Studio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio

		if err := manager.ValidateModifyStudio(ctx, updatedStudio, qb); err != nil {
			return err
		}

		var err error
		s, err = qb.Update(ctx, updatedStudio)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(ctx, s.ID, imageData); err != nil {
				return err
			}
		} else if imageIncluded {
			// must be unsetting
			if err := qb.DestroyImage(ctx, s.ID); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if translator.hasField("stash_ids") {
			stashIDJoins := input.StashIds
			if err := qb.UpdateStashIDs(ctx, studioID, stashIDJoins); err != nil {
				return err
			}
		}

		if translator.hasField("aliases") {
			if err := studio.EnsureAliasesUnique(ctx, studioID, input.Aliases, qb); err != nil {
				return err
			}

			if err := qb.UpdateAliases(ctx, studioID, input.Aliases); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, s.ID, plugin.StudioUpdatePost, input, translator.getFields())
	return r.getStudio(ctx, s.ID)
}

func (r *mutationResolver) StudioDestroy(ctx context.Context, input StudioDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.repository.Studio.Destroy(ctx, id)
	}); err != nil {
		return false, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, id, plugin.StudioDestroyPost, input, nil)

	return true, nil
}

func (r *mutationResolver) StudiosDestroy(ctx context.Context, studioIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(studioIDs)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Studio
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
		r.hookExecutor.ExecutePostHooks(ctx, id, plugin.StudioDestroyPost, studioIDs, nil)
	}

	return true, nil
}
