package scene

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/utils"
)

type Updater interface {
	PartialUpdater
	UpdateCover(ctx context.Context, sceneID int, cover []byte) error
}

type PartialUpdater interface {
	UpdatePartial(ctx context.Context, id int, updatedScene models.ScenePartial) (*models.Scene, error)
}

var ErrEmptyUpdater = errors.New("no fields have been set")

// UpdateSet is used to update a scene and its relationships.
type UpdateSet struct {
	ID int

	Partial models.ScenePartial

	// in future these could be moved into a separate struct and reused
	// for a Creator struct

	// Not set if nil. Set to []byte{} to clear existing
	CoverImage []byte
}

// IsEmpty returns true if there is nothing to update.
func (u *UpdateSet) IsEmpty() bool {
	withoutID := u.Partial

	return withoutID == models.ScenePartial{} &&
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
	updatedAt := time.Now()
	partial.UpdatedAt = models.NewOptionalTime(updatedAt)

	ret, err := qb.UpdatePartial(ctx, u.ID, partial)
	if err != nil {
		return nil, fmt.Errorf("error updating scene: %w", err)
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
	ret := u.Partial.UpdateInput(u.ID)

	if u.CoverImage != nil {
		// convert back to base64
		data := utils.GetBase64StringFromData(u.CoverImage)
		ret.CoverImage = &data
	}

	return ret
}

func UpdateFormat(ctx context.Context, qb PartialUpdater, id int, format string) (*models.Scene, error) {
	v := format
	return qb.UpdatePartial(ctx, id, models.ScenePartial{
		Format: models.NewOptionalString(v),
	})
}

func UpdateOSHash(ctx context.Context, qb PartialUpdater, id int, oshash string) (*models.Scene, error) {
	v := oshash

	return qb.UpdatePartial(ctx, id, models.ScenePartial{
		OSHash: models.NewOptionalString(v),
	})
}

func UpdateChecksum(ctx context.Context, qb PartialUpdater, id int, checksum string) (*models.Scene, error) {
	v := checksum

	return qb.UpdatePartial(ctx, id, models.ScenePartial{
		Checksum: models.NewOptionalString(v),
	})
}

func UpdateFileModTime(ctx context.Context, qb PartialUpdater, id int, modTime time.Time) (*models.Scene, error) {
	v := modTime
	return qb.UpdatePartial(ctx, id, models.ScenePartial{
		FileModTime: models.NewOptionalTime(v),
	})
}

func AddPerformer(ctx context.Context, qb PartialUpdater, o *models.Scene, performerID int) (bool, error) {
	if !intslice.IntInclude(o.PerformerIDs, performerID) {
		if _, err := qb.UpdatePartial(ctx, o.ID, models.ScenePartial{
			PerformerIDs: &models.UpdateIDs{
				IDs:  []int{performerID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddTag(ctx context.Context, qb PartialUpdater, o *models.Scene, tagID int) (bool, error) {
	if !intslice.IntInclude(o.TagIDs, tagID) {
		if _, err := qb.UpdatePartial(ctx, o.ID, models.ScenePartial{
			TagIDs: &models.UpdateIDs{
				IDs:  []int{tagID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddGallery(ctx context.Context, qb PartialUpdater, o *models.Scene, galleryID int) (bool, error) {
	if !intslice.IntInclude(o.GalleryIDs, galleryID) {
		if _, err := qb.UpdatePartial(ctx, o.ID, models.ScenePartial{
			TagIDs: &models.UpdateIDs{
				IDs:  []int{galleryID},
				Mode: models.RelationshipUpdateModeAdd,
			},
		}); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
