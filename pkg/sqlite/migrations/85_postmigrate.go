package migrations

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

func post85(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 85")

	ffprobePath, _ := exec.LookPath("ffprobe")

	mm := schema85PostMigrator{
		migrator: migrator{
			db: db,
		},
		ffprobe: ffmpeg.NewFFProbe(ffprobePath),
	}

	return mm.migrate(ctx)
}

type schema85PostMigrator struct {
	migrator
	ffprobe *ffmpeg.FFProbe
}

func (m *schema85PostMigrator) migrate(ctx context.Context) error {
	const (
		limit    = 1000
		logEvery = 10000
	)

	result := struct {
		Count int `db:"count"`
	}{0}

	if err := m.db.Get(&result, "SELECT COUNT(*) AS count FROM `video_files` WHERE `frames` IS NULL"); err != nil {
		return err
	}

	if result.Count == 0 {
		return nil
	}

	logger.Infof("Backfilling frames for %d video files...", result.Count)

	lastID := 0
	count := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := `
				SELECT f.id, folders.path, f.basename
				FROM video_files vf
				JOIN files f ON f.id = vf.file_id
				JOIN folders ON folders.id = f.parent_folder_id
				WHERE vf.frames IS NULL
			`
			if lastID != 0 {
				query += fmt.Sprintf(" AND f.id > %d", lastID)
			}
			query += fmt.Sprintf(" ORDER BY f.id LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var fileID int
				var dir string
				var basename string

				if err := rows.Scan(&fileID, &dir, &basename); err != nil {
					return err
				}

				gotSome = true
				lastID = fileID
				count++

				path := filepath.Join(dir, basename)

				frames, err := m.ffprobe.GetReadFrameCount(path)
				if err != nil || frames <= 0 {
					continue
				}

				if _, err := tx.Exec("UPDATE `video_files` SET `frames` = ? WHERE `file_id` = ?", frames, fileID); err != nil {
					return err
				}
			}

			return rows.Err()
		}); err != nil {
			return err
		}

		if !gotSome {
			break
		}

		if count%logEvery == 0 {
			logger.Infof("Checked %d video files", count)
		}
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(85, post85)
}
