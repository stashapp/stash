package scene

import (
	"database/sql"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

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
	performerIDs = utils.IntAppendUnique(performerIDs, performerID)

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
	tagIDs = utils.IntAppendUnique(tagIDs, tagID)

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
	galleryIDs = utils.IntAppendUnique(galleryIDs, galleryID)

	if len(galleryIDs) != oldLen {
		if err := qb.UpdateGalleries(id, galleryIDs); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
