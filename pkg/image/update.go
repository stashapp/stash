package image

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

func UpdateFileModTime(qb models.ImageWriter, id int, modTime models.NullSQLiteTimestamp) (*models.Image, error) {
	return qb.Update(models.ImagePartial{
		ID:          id,
		FileModTime: &modTime,
	})
}

func AddPerformer(qb models.ImageReaderWriter, id int, performerID int) (bool, error) {
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

func AddTag(qb models.ImageReaderWriter, id int, tagID int) (bool, error) {
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
