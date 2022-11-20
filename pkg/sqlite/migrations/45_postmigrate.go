package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
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

	if err := m.migrateTagImages(ctx); err != nil {
		return err
	}

	return nil
}

func (m *schema45Migrator) migrateTagImages(ctx context.Context) error {
	logger.Infof("Moving tag images to new table")

	const (
		limit    = 1000
		logEvery = 10000
	)

	count := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := "SELECT `tags_image`.`tag_id`, `tags_image`.`image` FROM `tags_image`"

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
				if _, err := m.db.Exec("UPDATE `tags` SET `image_checksum` = ? WHERE `id` = ?", checksum, id); err != nil {
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
			logger.Infof("Migrated %d folders", count)
		}
	}

	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		// drop the tags_image table
		logger.Debugf("Dropping tags_image")
		_, err := m.db.Exec("DROP TABLE `tags_image`")
		return err
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(45, post45)
}
