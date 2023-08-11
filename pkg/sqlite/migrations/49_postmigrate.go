package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

// var migrate49CriterionMap = map[models.FilterMode]map[string]string{
// 	models.FilterModeScenes: map[string]string{
// 		"id":                 "IntCriterionInput",
// 		"title":              "StringCriterionInput",
// 		"code":               "StringCriterionInput",
// 		"details":            "StringCriterionInput",
// 		"director":           "StringCriterionInput",
// 		"oshash":             "StringCriterionInput",
// 		"checksum":           "StringCriterionInput",
// 		"phash":              "StringCriterionInput",
// 		"phash_distance":     "PhashDistanceCriterionInput",
// 		"path":               "StringCriterionInput",
// 		"file_count":         "IntCriterionInput",
// 		"rating":             "IntCriterionInput",
// 		"rating100":          "IntCriterionInput",
// 		"organized":          "Boolean",
// 		"o_counter":          "IntCriterionInput",
// 		"duplicated":         "PHashDuplicationCriterionInput",
// 		"resolution":         "ResolutionCriterionInput",
// 		"duration":           "IntCriterionInput",
// 		"has_markers":        "String",
// 		"is_missing":         "String",
// 		"studios":            "HierarchicalMultiCriterionInput",
// 		"movies":             "MultiCriterionInput",
// 		"tags":               "HierarchicalMultiCriterionInput",
// 		"tag_count":          "IntCriterionInput",
// 		"performer_tags":     "HierarchicalMultiCriterionInput",
// 		"performer_favorite": "Boolean",
// 		"performer_age":      "IntCriterionInput",
// 		"performers":         "MultiCriterionInput",
// 		"performer_count":    "IntCriterionInput",
// 		"stash_id":           "StringCriterionInput",
// 		"stash_id_endpoint":  "StashIDCriterionInput",
// 		"url":                "StringCriterionInput",
// 		"interactive":        "Boolean",
// 		"interactive_speed":  "IntCriterionInput",
// 		"captions":           "StringCriterionInput",
// 		"resume_time":        "IntCriterionInput",
// 		"play_count":         "IntCriterionInput",
// 		"play_duration":      "IntCriterionInput",
// 		"date":               "DateCriterionInput",
// 		"created_at":         "TimestampCriterionInput",
// 		"updated_at":         "TimestampCriterionInput",
// 	},
// }

func post49(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 49")

	// translate the existing saved filters into the new format

	// existing format:
	/* const result = {
	   perPage: this.itemsPerPage,
	   sortby: this.getSortBy(),
	   sortdir:
	     this.sortBy === "date"
	       ? this.sortDirection === SortDirectionEnum.Asc
	         ? "asc"
	         : undefined
	       : this.sortDirection === SortDirectionEnum.Desc
	       ? "desc"
	       : undefined,
	   disp: this.displayMode,
	   q: this.searchTerm || undefined,
	   z: this.zoomIndex,
	   c: encodedCriteria,
	 }; */

	m := schema49Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrateSavedFilters(ctx)
}

type schema49Migrator struct {
	migrator
}

func (m *schema49Migrator) migrateSavedFilters(ctx context.Context) error {
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		rows, err := m.db.Query("SELECT id, mode, find_filter FROM saved_filters ORDER BY id")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id         int
				mode       models.FilterMode
				findFilter string
			)

			err := rows.Scan(&id, &mode, &findFilter)
			if err != nil {
				return err
			}

			asRawMessage := json.RawMessage(findFilter)

			newFindFilter, err := m.getFindFilter(asRawMessage)
			if err != nil {
				return fmt.Errorf("failed to get find filter for saved filter %d: %w", id, err)
			}

			objectFilter, err := m.getObjectFilter(mode, asRawMessage)
			if err != nil {
				return fmt.Errorf("failed to get object filter for saved filter %d: %w", id, err)
			}

			uiOptions, err := m.getDisplayOptions(asRawMessage)
			if err != nil {
				return fmt.Errorf("failed to get display options for saved filter %d: %w", id, err)
			}

			_, err = m.db.Exec("UPDATE saved_filters SET find_filter = ?, object_filter = ?, ui_options = ? WHERE id = ?", newFindFilter, objectFilter, uiOptions, id)
			if err != nil {
				return fmt.Errorf("failed to update saved filter %d: %w", id, err)
			}
		}

		return rows.Err()
	}); err != nil {
		return err
	}

	return nil
}

func (m *schema49Migrator) getDisplayOptions(data json.RawMessage) (json.RawMessage, error) {
	type displayOptions struct {
		DisplayMode *int `json:"disp"`
		ZoomIndex   *int `json:"z"`
	}

	var opts displayOptions
	if err := json.Unmarshal(data, &opts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal display options: %w", err)
	}

	ret := make(map[string]interface{})
	if opts.DisplayMode != nil {
		ret["display_mode"] = *opts.DisplayMode
	}
	if opts.ZoomIndex != nil {
		ret["zoom_index"] = *opts.ZoomIndex
	}

	return json.Marshal(ret)
}

func (m *schema49Migrator) getFindFilter(data json.RawMessage) (json.RawMessage, error) {
	type findFilterJson struct {
		Q         *string `json:"q"`
		PerPage   *int    `json:"perPage"`
		Sort      *string `json:"sortby"`
		Direction *string `json:"sortdir"`
	}

	var ff findFilterJson
	if err := json.Unmarshal(data, &ff); err != nil {
		return nil, fmt.Errorf("failed to unmarshal find filter: %w", err)
	}

	if ff.Direction != nil {
		newDir := strings.ToUpper(*ff.Direction)
		ff.Direction = &newDir
	}

	// remarshal back to json
	return json.Marshal(ff)
}

func (m *schema49Migrator) getObjectFilter(mode models.FilterMode, data json.RawMessage) (json.RawMessage, error) {
	type criteriaJson struct {
		Criteria []string `json:"c"`
	}

	var c criteriaJson
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object filter: %w", err)
	}

	ret := make(map[string]interface{})
	for _, raw := range c.Criteria {
		if err := m.convertCriterion(mode, ret, raw); err != nil {
			return nil, err
		}
	}

	return json.Marshal(ret)
}

func (m *schema49Migrator) convertCriterion(mode models.FilterMode, out map[string]interface{}, criterion string) error {
	// convert to a map
	ret := make(map[string]interface{})

	if err := json.Unmarshal([]byte(criterion), &ret); err != nil {
		return fmt.Errorf("failed to unmarshal criterion: %w", err)
	}

	// TODO - convert UI type to parameter name
	field := ret["type"].(string)
	delete(ret, "type")

	out[field] = ret

	return nil
}

func init() {
	sqlite.RegisterPostMigration(49, post49)
}
