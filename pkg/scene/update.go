package scene

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/utils"
)

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
func (u *UpdateSet) Update(qb models.SceneWriter, screenshotSetter ScreenshotSetter) (*models.Scene, error) {
	if u.IsEmpty() {
		return nil, ErrEmptyUpdater
	}

	partial := u.Partial
	partial.ID = u.ID
	partial.UpdatedAt = &models.SQLiteTimestamp{
		Timestamp: time.Now(),
	}

	ret, err := qb.Update(partial)
	if err != nil {
		return nil, fmt.Errorf("error updating scene: %w", err)
	}

	if u.PerformerIDs != nil {
		if err := qb.UpdatePerformers(u.ID, u.PerformerIDs); err != nil {
			return nil, fmt.Errorf("error updating scene performers: %w", err)
		}
	}

	if u.TagIDs != nil {
		if err := qb.UpdateTags(u.ID, u.TagIDs); err != nil {
			return nil, fmt.Errorf("error updating scene tags: %w", err)
		}
	}

	if u.StashIDs != nil {
		if err := qb.UpdateStashIDs(u.ID, u.StashIDs); err != nil {
			return nil, fmt.Errorf("error updating scene stash_ids: %w", err)
		}
	}

	if u.CoverImage != nil {
		if err := qb.UpdateCover(u.ID, u.CoverImage); err != nil {
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

func UpdateFormat(qb models.SceneWriter, id int, format string) (*models.Scene, error) {
	return qb.Update(models.ScenePartial{
		ID: id,
		Format: &sql.NullString{
			String: format,
			Valid:  true,
		},
	})
}

func UpdateOSHash(qb models.SceneWriter, id int, oshash string) (*models.Scene, error) {
	return qb.Update(models.ScenePartial{
		ID: id,
		OSHash: &sql.NullString{
			String: oshash,
			Valid:  true,
		},
	})
}

func UpdateChecksum(qb models.SceneWriter, id int, checksum string) (*models.Scene, error) {
	return qb.Update(models.ScenePartial{
		ID: id,
		Checksum: &sql.NullString{
			String: checksum,
			Valid:  true,
		},
	})
}

func UpdateFileModTime(qb models.SceneWriter, id int, modTime models.NullSQLiteTimestamp) (*models.Scene, error) {
	return qb.Update(models.ScenePartial{
		ID:          id,
		FileModTime: &modTime,
	})
}

func AddPerformer(qb models.SceneReaderWriter, id int, performerID int) (bool, error) {
	performerIDs, err := qb.GetPerformerIDs(id)
	if err != nil {
		return false, err
	}

	oldLen := len(performerIDs)
	performerIDs = intslice.IntAppendUnique(performerIDs, performerID)

	if len(performerIDs) != oldLen {
		if err := qb.UpdatePerformers(id, performerIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddTag(qb models.SceneReaderWriter, id int, tagID int) (bool, error) {
	tagIDs, err := qb.GetTagIDs(id)
	if err != nil {
		return false, err
	}

	oldLen := len(tagIDs)
	tagIDs = intslice.IntAppendUnique(tagIDs, tagID)

	if len(tagIDs) != oldLen {
		if err := qb.UpdateTags(id, tagIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func AddGallery(qb models.SceneReaderWriter, id int, galleryID int) (bool, error) {
	galleryIDs, err := qb.GetGalleryIDs(id)
	if err != nil {
		return false, err
	}

	oldLen := len(galleryIDs)
	galleryIDs = intslice.IntAppendUnique(galleryIDs, galleryID)

	if len(galleryIDs) != oldLen {
		if err := qb.UpdateGalleries(id, galleryIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
