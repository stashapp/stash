package scene

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/utils"
)

type Updater interface {
	PartialUpdater
	UpdatePerformers(ctx context.Context, sceneID int, performerIDs []int) error
	UpdateTags(ctx context.Context, sceneID int, tagIDs []int) error
	UpdateStashIDs(ctx context.Context, sceneID int, stashIDs []models.StashID) error
	UpdateCover(ctx context.Context, sceneID int, cover []byte) error
}

type PartialUpdater interface {
	Update(ctx context.Context, updatedScene models.ScenePartial) (*models.Scene, error)
}

type PerformerUpdater interface {
	GetPerformerIDs(ctx context.Context, sceneID int) ([]int, error)
	UpdatePerformers(ctx context.Context, sceneID int, performerIDs []int) error
}

type TagUpdater interface {
	GetTagIDs(ctx context.Context, sceneID int) ([]int, error)
	UpdateTags(ctx context.Context, sceneID int, tagIDs []int) error
}

type GalleryUpdater interface {
	GetGalleryIDs(ctx context.Context, sceneID int) ([]int, error)
	UpdateGalleries(ctx context.Context, sceneID int, galleryIDs []int) error
}

var ErrEmptyUpdater = errors.New("no fields have been set")

// UpdateSet is used to update a scene and its relationships.
type UpdateSet struct {
	ID int

	Partial models.ScenePartial

	// in future these could be moved into a separate struct and reused
	// for a Creator struct

	// Not set if nil. Set to []int{} to clear existing
	PerformerIDs []int
	// Not set if nil. Set to []int{} to clear existing
	TagIDs []int
	// Not set if nil. Set to []int{} to clear existing
	StashIDs []models.StashID
	// Not set if nil. Set to []byte{} to clear existing
	CoverImage []byte
}

// IsEmpty returns true if there is nothing to update.
func (u *UpdateSet) IsEmpty() bool {
	withoutID := u.Partial
	withoutID.ID = 0

	return withoutID == models.ScenePartial{} &&
		u.PerformerIDs == nil &&
		u.TagIDs == nil &&
		u.StashIDs == nil &&
		u.CoverImage == nil
}

// Update updates a scene by updating the fields in the Partial field, then
// updates non-nil relationships. Returns an error if there is no work to
// be done.
func (u *UpdateSet) Update(ctx context.Context, qb Updater, screenshotSetter ScreenshotSetter) (*models.Scene, error) {
	if u.IsEmpty() {
		return nil, ErrEmptyUpdater
	}

	partial := u.Partial
	partial.ID = u.ID
	partial.UpdatedAt = &models.SQLiteTimestamp{
		Timestamp: time.Now(),
	}

	ret, err := qb.Update(ctx, partial)
	if err != nil {
		return nil, fmt.Errorf("error updating scene: %w", err)
	}

	if u.PerformerIDs != nil {
		if err := qb.UpdatePerformers(ctx, u.ID, u.PerformerIDs); err != nil {
			return nil, fmt.Errorf("error updating scene performers: %w", err)
		}
	}

	if u.TagIDs != nil {
		if err := qb.UpdateTags(ctx, u.ID, u.TagIDs); err != nil {
			return nil, fmt.Errorf("error updating scene tags: %w", err)
		}
	}

	if u.StashIDs != nil {
		if err := qb.UpdateStashIDs(ctx, u.ID, u.StashIDs); err != nil {
			return nil, fmt.Errorf("error updating scene stash_ids: %w", err)
		}
	}

	if u.CoverImage != nil {
		if err := qb.UpdateCover(ctx, u.ID, u.CoverImage); err != nil {
			return nil, fmt.Errorf("error updating scene cover: %w", err)
		}

		if err := screenshotSetter.SetScreenshot(ret, u.CoverImage); err != nil {
			return nil, fmt.Errorf("error setting scene screenshot: %w", err)
		}
	}

	return ret, nil
}

// UpdateInput converts the UpdateSet into SceneUpdateInput for hook firing purposes.
func (u UpdateSet) UpdateInput() models.SceneUpdateInput {
	// ensure the partial ID is set
	u.Partial.ID = u.ID
	ret := u.Partial.UpdateInput()

	if u.PerformerIDs != nil {
		ret.PerformerIds = intslice.IntSliceToStringSlice(u.PerformerIDs)
	}

	if u.TagIDs != nil {
		ret.TagIds = intslice.IntSliceToStringSlice(u.TagIDs)
	}

	if u.StashIDs != nil {
		for _, s := range u.StashIDs {
			ss := s.StashIDInput()
			ret.StashIds = append(ret.StashIds, &ss)
		}
	}

	if u.CoverImage != nil {
		// convert back to base64
		data := utils.GetBase64StringFromData(u.CoverImage)
		ret.CoverImage = &data
	}

	return ret
}

func UpdateFormat(ctx context.Context, qb PartialUpdater, id int, format string) (*models.Scene, error) {
	return qb.Update(ctx, models.ScenePartial{
		ID: id,
		Format: &sql.NullString{
			String: format,
			Valid:  true,
		},
	})
}

func UpdateOSHash(ctx context.Context, qb PartialUpdater, id int, oshash string) (*models.Scene, error) {
	return qb.Update(ctx, models.ScenePartial{
		ID: id,
		OSHash: &sql.NullString{
			String: oshash,
			Valid:  true,
		},
	})
}

func UpdateChecksum(ctx context.Context, qb PartialUpdater, id int, checksum string) (*models.Scene, error) {
	return qb.Update(ctx, models.ScenePartial{
		ID: id,
		Checksum: &sql.NullString{
			String: checksum,
			Valid:  true,
		},
	})
}

func UpdateFileModTime(ctx context.Context, qb PartialUpdater, id int, modTime models.NullSQLiteTimestamp) (*models.Scene, error) {
	return qb.Update(ctx, models.ScenePartial{
		ID:          id,
		FileModTime: &modTime,
	})
}

func AddPerformer(ctx context.Context, qb PerformerUpdater, id int, performerID int) (bool, error) {
	performerIDs, err := qb.GetPerformerIDs(ctx, id)
	if err != nil {
		return false, err
	}

	oldLen := len(performerIDs)
	performerIDs = intslice.IntAppendUnique(performerIDs, performerID)

	if len(performerIDs) != oldLen {
		if err := qb.UpdatePerformers(ctx, id, performerIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddTag(ctx context.Context, qb TagUpdater, id int, tagID int) (bool, error) {
	tagIDs, err := qb.GetTagIDs(ctx, id)
	if err != nil {
		return false, err
	}

	oldLen := len(tagIDs)
	tagIDs = intslice.IntAppendUnique(tagIDs, tagID)

	if len(tagIDs) != oldLen {
		if err := qb.UpdateTags(ctx, id, tagIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddGallery(ctx context.Context, qb GalleryUpdater, id int, galleryID int) (bool, error) {
	galleryIDs, err := qb.GetGalleryIDs(ctx, id)
	if err != nil {
		return false, err
	}

	oldLen := len(galleryIDs)
	galleryIDs = intslice.IntAppendUnique(galleryIDs, galleryID)

	if len(galleryIDs) != oldLen {
		if err := qb.UpdateGalleries(ctx, id, galleryIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
