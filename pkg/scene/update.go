package scene

import (
	"database/sql"

	"github.com/stashapp/stash/pkg/models"
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
