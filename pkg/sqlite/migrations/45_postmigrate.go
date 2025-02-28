package migrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

type schema45Migrator struct {
	migrator
	hasBlobs bool
}

func post45(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 45")

	m := schema45Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable: "tags_image",
		joinIDCol: "tag_id",
		destTable: "tags",
		cols: []migrateImageToBlobOptions{
			{
				joinImageCol: "image",
				destCol:      "image_blob",
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to migrate images table for tags: %w", err)
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable: "studios_image",
		joinIDCol: "studio_id",
		destTable: "studios",
		cols: []migrateImageToBlobOptions{
			{
				joinImageCol: "image",
				destCol:      "image_blob",
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to migrate images table for studios: %w", err)
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable: "performers_image",
		joinIDCol: "performer_id",
		destTable: "performers",
		cols: []migrateImageToBlobOptions{
			{
				joinImageCol: "image",
				destCol:      "image_blob",
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to migrate images table for performers: %w", err)
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable: "scenes_cover",
		joinIDCol: "scene_id",
		destTable: "scenes",
		cols: []migrateImageToBlobOptions{
			{
				joinImageCol: "cover",
				destCol:      "cover_blob",
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to migrate images table for scenes: %w", err)
	}

	if err := m.migrateImagesTable(ctx, migrateImagesTableOptions{
		joinTable: "movies_images",
		joinIDCol: "movie_id",
		destTable: "movies",
		cols: []migrateImageToBlobOptions{
			{
				joinImageCol: "front_image",
				destCol:      "front_image_blob",
			},
			{
				joinImageCol: "back_image",
				destCol:      "back_image_blob",
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to migrate images table for movies: %w", err)
	}

	tablesToDrop := []string{
		"tags_image",
		"studios_image",
		"performers_image",
		"scenes_cover",
		"movies_images",
	}

	for _, table := range tablesToDrop {
		if err := m.dropTable(ctx, table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	if err := m.migrateConfig(ctx); err != nil {
		return fmt.Errorf("failed to migrate config: %w", err)
	}

	return nil
}

type migrateImageToBlobOptions struct {
	joinImageCol string
	destCol      string
}

type migrateImagesTableOptions struct {
	joinTable string
	joinIDCol string
	destTable string
	cols      []migrateImageToBlobOptions
}

func (o migrateImagesTableOptions) selectColumns() string {
	var cols []string
	for _, c := range o.cols {
		cols = append(cols, "`"+c.joinImageCol+"`")
	}

	return strings.Join(cols, ", ")
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
			query := fmt.Sprintf("SELECT %s, %s FROM `%s`", options.joinIDCol, options.selectColumns(), options.joinTable)

			query += fmt.Sprintf(" LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				m.hasBlobs = true

				var id int

				result := make([]interface{}, len(options.cols)+1)
				result[0] = &id
				for i := range options.cols {
					v := []byte{}
					result[i+1] = &v
				}

				err := rows.Scan(result...)
				if err != nil {
					return err
				}

				gotSome = true
				count++

				for i, col := range options.cols {
					image := result[i+1].(*[]byte)

					if len(*image) > 0 {
						if err := m.insertImage(tx, *image, id, options.destTable, col.destCol); err != nil {
							return err
						}
					}
				}

				// delete the row from the join table so we don't process it again
				deleteSQL := utils.StrFormat("DELETE FROM `{joinTable}` WHERE `{joinIDCol}` = ?", utils.StrFormatMap{
					"joinTable": options.joinTable,
					"joinIDCol": options.joinIDCol,
				})
				if _, err := tx.Exec(deleteSQL, id); err != nil {
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

	return nil
}

func (m *schema45Migrator) insertImage(tx *sqlx.Tx, data []byte, id int, destTable string, destCol string) error {
	// calculate checksum and insert into blobs table
	checksum := md5.FromBytes(data)

	if _, err := tx.Exec("INSERT INTO `blobs` (`checksum`, `blob`) VALUES (?, ?) ON CONFLICT DO NOTHING", checksum, data); err != nil {
		return err
	}

	// set the tag image checksum
	updateSQL := utils.StrFormat("UPDATE `{destTable}` SET `{destCol}` = ? WHERE `id` = ?", utils.StrFormatMap{
		"destTable": destTable,
		"destCol":   destCol,
	})
	if _, err := tx.Exec(updateSQL, checksum, id); err != nil {
		return err
	}

	return nil
}

func (m *schema45Migrator) dropTable(ctx context.Context, table string) error {
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		logger.Debugf("Dropping %s", table)
		_, err := tx.Exec(fmt.Sprintf("DROP TABLE `%s`", table))
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (m *schema45Migrator) migrateConfig(ctx context.Context) error {
	c := config.GetInstance()

	// if we don't have blobs, and storage is already set, then don't overwrite
	if !m.hasBlobs && c.GetBlobsStorage().IsValid() {
		logger.Infof("Blobs storage already set, not overwriting")
		return nil
	}

	// if we have blobs in the database, then default to database storage
	// otherwise default to filesystem storage
	defaultStorage := config.BlobStorageTypeFilesystem
	if m.hasBlobs || c.GetBlobsPath() == "" {
		defaultStorage = config.BlobStorageTypeDatabase
	}

	logger.Infof("Setting blobs storage to %s", defaultStorage.String())
	c.SetInterface(config.BlobsStorage, defaultStorage)
	if err := c.Write(); err != nil {
		logger.Errorf("Error while writing configuration file: %s", err.Error())
	}

	// if default scan settings are set, then set to generate scene covers by default
	scanDefaults := c.GetDefaultScanSettings()
	if scanDefaults != nil {
		scanDefaults.ScanGenerateCovers = true
		c.SetInterface(config.DefaultScanSettings, scanDefaults)
		if err := c.Write(); err != nil {
			logger.Errorf("Error while writing configuration file: %s", err.Error())
		}
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(45, post45)
}
