package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

type schema45Migrator struct {
	migrator
}

func post45(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 45")

	m := schema45Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable:    "tags_image",
		joinIDCol:    "tag_id",
		joinImageCol: "image",
		destTable:    "tags",
		destCol:      "image_blob",
	}); err != nil {
		return err
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable:    "studios_image",
		joinIDCol:    "studio_id",
		joinImageCol: "image",
		destTable:    "studios",
		destCol:      "image_blob",
	}); err != nil {
		return err
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable:    "performers_image",
		joinIDCol:    "performer_id",
		joinImageCol: "image",
		destTable:    "performers",
		destCol:      "image_blob",
	}); err != nil {
		return err
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable:    "scenes_cover",
		joinIDCol:    "scene_id",
		joinImageCol: "cover",
		destTable:    "scenes",
		destCol:      "cover_blob",
	}); err != nil {
		return err
	}

	return nil
}

type migrateImagesTableOptions struct {
	joinTable    string
	joinIDCol    string
	joinImageCol string
	destTable    string
	destCol      string
}

func (m *schema45Migrator) migrateImagesTable(ctx context.Context, options migrateImagesTableOptions) error {
	logger.Infof("Moving %s to blobs table", options.joinTable)

	const (
		limit    = 1000
		logEvery = 10000
	)

	count := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := utils.StrFormat("SELECT `{joinTable}`.`{joinIDCol}`, `{joinTable}`.`{joinImageCol}` FROM `{joinTable}`", utils.StrFormatMap{
				"joinTable":    options.joinTable,
				"joinIDCol":    options.joinIDCol,
				"joinImageCol": options.joinImageCol,
			})

			query += fmt.Sprintf(" LIMIT %d", limit)

			rows, err := m.db.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				var data []byte

				err := rows.Scan(&id, &data)
				if err != nil {
					return err
				}

				gotSome = true
				count++

				// calculate checksum and insert into blobs table
				checksum := md5.FromBytes(data)

				if _, err := m.db.Exec("INSERT INTO `blobs` (`checksum`, `blob`) VALUES (?, ?)", checksum, data); err != nil {
					return err
				}

				// set the tag image checksum
				updateSQL := utils.StrFormat("UPDATE `{destTable}` SET `{destCol}` = ? WHERE `id` = ?", utils.StrFormatMap{
					"destTable": options.destTable,
					"destCol":   options.destCol,
				})
				if _, err := m.db.Exec(updateSQL, checksum, id); err != nil {
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
			logger.Infof("Migrated %d images", count)
		}
	}

	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		// drop the tags_image table
		logger.Debugf("Dropping %s", options.joinTable)
		_, err := m.db.Exec(fmt.Sprintf("DROP TABLE `%s`", options.joinTable))
		return err
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(45, post45)
}
