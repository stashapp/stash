package gallery

import "github.com/stashapp/stash/pkg/models"

func UpdateFileModTime(qb models.GalleryWriter, id int, modTime models.NullSQLiteTimestamp) (*models.Gallery, error) {
	return qb.UpdatePartial(models.GalleryPartial{
		ID:          id,
		FileModTime: &modTime,
	})
}
