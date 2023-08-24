package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

var migrate49TypeResolution = map[string][]string{
	"Boolean": {
		/*
			"organized",
			"interactive",
			"ignore_auto_tag",
			"performer_favorite",
			"filter_favorites",
		*/
	},
	"Int": {
		"id",
		"rating",
		"rating100",
		"o_counter",
		"duration",
		"tag_count",
		"age",
		"height",
		"height_cm",
		"weight",
		"scene_count",
		"marker_count",
		"image_count",
		"gallery_count",
		"performer_count",
		"interactive_speed",
		"resume_time",
		"play_count",
		"play_duration",
		"parent_count",
		"child_count",
		"performer_age",
		"file_count",
	},
	"Float": {
		"penis_length",
	},
	"Object": {
		"tags",
		"performers",
		"studios",
		"movies",
		"galleries",
		"parents",
		"children",
		"scene_tags",
		"performer_tags",
	},
}
var migrate49NameChanges = map[string]string{
	"rating":             "rating100",
	"parent_studios":     "parents",
	"child_studios":      "children",
	"parent_tags":        "parents",
	"child_tags":         "children",
	"child_tag_count":    "child_count",
	"parent_tag_count":   "parent_count",
	"height":             "height_cm",
	"imageIsMissing":     "is_missing",
	"sceneIsMissing":     "is_missing",
	"galleryIsMissing":   "is_missing",
	"performerIsMissing": "is_missing",
	"tagIsMissing":       "is_missing",
	"studioIsMissing":    "is_missing",
	"movieIsMissing":     "is_missing",
	"favorite":           "filter_favorites",
	"hasMarkers":         "has_markers",
	"parentTags":         "parents",
	"childTags":          "children",
	"phash":              "phash_distance",
	"scene_code":         "code",
	"hasChapters":        "has_chapters",
	"sceneChecksum":      "checksum",
	"galleryChecksum":    "checksum",
	"sceneTags":          "scene_tags",
	"performerTags":      "performer_tags",
}

func post49(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 49")

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
		Page      *int    `json:"page"`
		PerPage   *int    `json:"perPage"`
		Sort      *string `json:"sortby"`
		Direction *string `json:"sortdir"`
	}

	ppDefault := 40
	pageDefault := 1
	qDefault := ""
	sortDefault := "date"
	asc := "asc"
	ff := findFilterJson{Q: &qDefault, Page: &pageDefault, PerPage: &ppDefault, Sort: &sortDefault, Direction: &asc}
	if err := json.Unmarshal(data, &ff); err != nil {
		return nil, fmt.Errorf("failed to unmarshal find filter: %w", err)
	}

	newDir := strings.ToUpper(*ff.Direction)
	ff.Direction = &newDir

	type findFilterRewrite struct {
		Q         *string `json:"q"`
		Page      *int    `json:"page"`
		PerPage   *int    `json:"per_page"`
		Sort      *string `json:"sort"`
		Direction *string `json:"direction"`
	}

	fr := findFilterRewrite(ff)

	return json.Marshal(fr)
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

	field := ret["type"].(string)
	// Some names are depracated
	if newFieldName, ok := migrate49NameChanges[field]; ok {
		field = newFieldName
	}
	delete(ret, "type")

	// Find out whether the object needs some adjustment/has non-string content attached
	if arrayContains(migrate49TypeResolution["Boolean"], field) {
		ret["value"] = adjustCriterionValue(ret["value"], "bool")
	}
	if arrayContains(migrate49TypeResolution["Int"], field) {
		ret["value"] = adjustCriterionValue(ret["value"], "int")
	}
	if arrayContains(migrate49TypeResolution["Float"], field) {
		ret["value"] = adjustCriterionValue(ret["value"], "float64")
	}
	if arrayContains(migrate49TypeResolution["Object"], field) {
		ret["value"] = adjustCriterionValue(ret["value"], "object")
	}
	out[field] = ret

	return nil
}

func arrayContains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

// General Function for converting the types inside a criterion
func adjustCriterionValue(value interface{}, t string) interface{} {
	if mapvalue, ok := value.(map[string]interface{}); ok {
		// Primitive values and lists of them
		for _, next := range []string{"value", "value2"} {
			if valmap, ok := mapvalue[next].([]string); ok {
				var valNewMap []interface{}
				for index, v := range valmap {
					valNewMap[index] = convertString(interface{}(v), t)
				}
				mapvalue[next] = interface{}(valNewMap)
			} else if _, ok := mapvalue[next]; ok {
				mapvalue[next] = convertString(mapvalue[next], t)
			}
		}
		// Items
		for _, next := range []string{"items", "excluded"} {
			if _, ok := mapvalue[next]; ok {
				mapvalue[next] = adjustCriterionItem(mapvalue[next])
			}
		}

		// Those Values are always Int
		for _, next := range []string{"Distance", "Depth"} {
			if _, ok := mapvalue[next]; ok {
				if formattedOut, ok := strconv.ParseInt(mapvalue[next].(string), 10, 64); ok == nil {
					mapvalue[next] = interface{}(formattedOut)
				}
			}
		}
		return mapvalue
	} else if _, ok := value.(string); ok {
		// Singular Primitive Values
		return convertString(value, t)
	} else if listvalue, ok := value.([]interface{}); ok {
		// Items as a singular value, as well as singular lists
		if t == "object" {
			value = adjustCriterionItem(value)
		} else {
			for index, val := range listvalue {
				listvalue[index] = convertString(val, t)
			}
			value = interface{}(listvalue)
		}
		return value
	} else if _, ok := value.(int); ok {
		return value
	}
	fmt.Printf("Could not recognize format of value %v\n", value)
	return value
}

// Converts values inside a criterion that represent some objects, like performer or studio.
func adjustCriterionItem(value interface{}) interface{} {
	// Basically, this first converts step by step the value, after that it adjusts id and Depth (of parent/child studios) to int
	if itemlist, ok := value.([]interface{}); ok {
		var itemNewList []interface{}
		for _, val := range itemlist {
			if val, ok := val.(map[string]interface{}); ok {
				newItem := make(map[string]interface{})
				for index, v := range val {
					if v, ok := v.(string); ok {
						switch index {
						case "id":
							if formattedOut, ok := strconv.ParseInt(v, 10, 64); ok == nil {
								newItem["id"] = formattedOut
							}
						case "Depth":
							if formattedOut, ok := strconv.ParseInt(v, 10, 64); ok == nil {
								newItem["Depth"] = formattedOut
							}
						default:
							newItem[index] = v
						}
					}
				}
				itemNewList = append(itemNewList, interface{}(newItem))
			}
		}
		return interface{}(itemNewList)
	}
	fmt.Printf("Could not recognize %v as an item list \n", value)
	return value
}

// Converts a value of type string to its according type, given by string
func convertString(value interface{}, t string) interface{} {
	if val, ok := value.(string); ok {
		switch t {
		case "float64":
			if formattedOut, ok := strconv.ParseFloat(val, 64); ok == nil {
				return interface{}(formattedOut)
			}
		case "int":
			if formattedOut, ok := strconv.ParseInt(val, 10, 64); ok == nil {
				return interface{}(formattedOut)
			}
		case "bool":
			if formattedOut, ok := strconv.ParseBool(val); ok == nil {
				return interface{}(formattedOut)
			}
		default:
			fmt.Printf("No valid conversiontype, need bool, int or float64\n")
			return value
		}
	}
	if reflect.TypeOf(value).Name() != t && !(t == "int" && reflect.TypeOf(value).Name() == "float64") && !(t == "float64" && reflect.TypeOf(value).Name() == "int") {
		fmt.Printf("Failed to convert %v to String, leaving unmodified.\n", value)
	}
	return value
}

func init() {
	sqlite.RegisterPostMigration(49, post49)
}
