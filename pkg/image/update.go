package image

import "github.com/stashapp/stash/pkg/models"

func UpdateFileModTime(qb models.ImageWriter, id int, modTime models.NullSQLiteTimestamp) (*models.Image, error) {
	return qb.Update(models.ImagePartial{
		ID:          id,
		FileModTime: &modTime,
	})
}
