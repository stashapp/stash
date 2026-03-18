package migrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

func post86(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 86")

	m := schema86Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrateSavedFilters(ctx)
}

type schema86Migrator struct {
	migrator
}

func (m *schema86Migrator) migrateSavedFilters(ctx context.Context) error {
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		rows, err := tx.Query("SELECT id, object_filter FROM saved_filters ORDER BY id")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id           int
				objectFilter []byte
			)

			err := rows.Scan(&id, &objectFilter)
			if err != nil {
				return err
			}

			if len(objectFilter) == 0 {
				continue
			}

			newObjectFilter, err := m.convertObjectFilter(objectFilter)
			if err != nil {
				return fmt.Errorf("failed to convert object filter for saved filter %d: %w", id, err)
			}

			if newObjectFilter != nil {
				_, err = tx.Exec("UPDATE saved_filters SET object_filter = ? WHERE id = ?", newObjectFilter, id)
				if err != nil {
					return fmt.Errorf("failed to update saved filter %d: %w", id, err)
				}
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	return nil
}

func (m *schema86Migrator) convertObjectFilter(data []byte) ([]byte, error) {
	var filter map[string]interface{}
	if err := json.Unmarshal(data, &filter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object filter: %w", err)
	}

	for _, value := range filter {
		criterion, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		v, hasValue := criterion["value"]
		if !hasValue || v == nil {
			continue
		}

		if valueObj, isObj := v.(map[string]interface{}); isObj {
			_, hasItems := valueObj["items"]
			_, hasExcluded := valueObj["excluded"]

			if hasItems || hasExcluded {
				var values []string
				if items, ok := valueObj["items"].([]interface{}); ok {
					for _, item := range items {
						if itemMap, isMap := item.(map[string]interface{}); isMap {
							if idStr, ok := itemMap["id"].(string); ok {
								values = append(values, idStr)
							} else if idFloat, ok := itemMap["id"].(float64); ok {
								values = append(values, fmt.Sprintf("%d", int(idFloat)))
							}
						}
					}
				}

				var excludes []string
				if excluded, ok := valueObj["excluded"].([]interface{}); ok {
					for _, item := range excluded {
						if itemMap, isMap := item.(map[string]interface{}); isMap {
							if idStr, ok := itemMap["id"].(string); ok {
								excludes = append(excludes, idStr)
							} else if idFloat, ok := itemMap["id"].(float64); ok {
								excludes = append(excludes, fmt.Sprintf("%d", int(idFloat)))
							}
						}
					}
				}

				var depth interface{}
				if d, ok := valueObj["depth"]; ok {
					depth = d
				} else if d, ok := valueObj["Depth"]; ok {
					depth = d
				}

				if len(values) > 0 {
					criterion["value"] = values
				} else {
					criterion["value"] = []string{}
				}

				if len(excludes) > 0 || hasExcluded {
					if excludes == nil {
						criterion["excludes"] = []string{}
					} else {
						criterion["excludes"] = excludes
					}
				}

				if depth != nil {
					criterion["depth"] = depth
				}
			}
		}
	}

	return json.Marshal(filter)
}

func init() {
	sqlite.RegisterPostMigration(86, post86)
}
