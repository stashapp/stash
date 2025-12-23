package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

const updateInputField = "input"

func getArgumentMap(ctx context.Context) map[string]interface{} {
	rctx := graphql.GetFieldContext(ctx)
	reqCtx := graphql.GetOperationContext(ctx)
	return rctx.Field.ArgumentMap(reqCtx.Variables)
}

func getUpdateInputMap(ctx context.Context) map[string]interface{} {
	return getNamedUpdateInputMap(ctx, updateInputField)
}

func getNamedUpdateInputMap(ctx context.Context, field string) map[string]interface{} {
	args := getArgumentMap(ctx)

	// field can be qualified
	fields := strings.Split(field, ".")

	currArgs := args

	for _, f := range fields {
		v, found := currArgs[f]
		if !found {
			currArgs = nil
			break
		}

		currArgs, _ = v.(map[string]interface{})
		if currArgs == nil {
			break
		}
	}

	if currArgs != nil {
		return currArgs
	}

	return make(map[string]interface{})
}

func getUpdateInputMaps(ctx context.Context) []map[string]interface{} {
	args := getArgumentMap(ctx)

	input := args[updateInputField]
	var ret []map[string]interface{}
	if input != nil {
		// convert []interface{} into []map[string]interface{}
		iSlice, _ := input.([]interface{})
		for _, i := range iSlice {
			m, _ := i.(map[string]interface{})
			if m != nil {
				ret = append(ret, m)
			}
		}
	}

	return ret
}

type changesetTranslator struct {
	inputMap map[string]interface{}
}

func (t changesetTranslator) hasField(field string) bool {
	if t.inputMap == nil {
		return false
	}

	_, found := t.inputMap[field]
	return found
}

func (t changesetTranslator) getFields() []string {
	var ret []string
	for k := range t.inputMap {
		ret = append(ret, k)
	}

	return ret
}

func (t changesetTranslator) string(value *string) string {
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}

func (t changesetTranslator) optionalString(value *string, field string) models.OptionalString {
	if !t.hasField(field) {
		return models.OptionalString{}
	}

	if value == nil {
		return models.NewOptionalStringPtr(nil)
	}

	trimmed := strings.TrimSpace(*value)
	return models.NewOptionalString(trimmed)
}

func (t changesetTranslator) optionalDate(value *string, field string) (models.OptionalDate, error) {
	if !t.hasField(field) {
		return models.OptionalDate{}, nil
	}

	if value == nil || *value == "" {
		return models.OptionalDate{
			Set:  true,
			Null: true,
		}, nil
	}

	date, err := models.ParseDate(*value)
	if err != nil {
		return models.OptionalDate{}, err
	}

	return models.NewOptionalDate(date), nil
}

func (t changesetTranslator) datePtr(value *string) (*models.Date, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	date, err := models.ParseDate(*value)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func (t changesetTranslator) intPtrFromString(value *string) (*int, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	vv, err := strconv.Atoi(*value)
	if err != nil {
		return nil, fmt.Errorf("converting %v to int: %w", *value, err)
	}
	return &vv, nil
}

func (t changesetTranslator) optionalInt(value *int, field string) models.OptionalInt {
	if !t.hasField(field) {
		return models.OptionalInt{}
	}

	return models.NewOptionalIntPtr(value)
}

func (t changesetTranslator) optionalIntFromString(value *string, field string) (models.OptionalInt, error) {
	if !t.hasField(field) {
		return models.OptionalInt{}, nil
	}

	if value == nil {
		return models.OptionalInt{
			Set:  true,
			Null: true,
		}, nil
	}

	vv, err := strconv.Atoi(*value)
	if err != nil {
		return models.OptionalInt{}, fmt.Errorf("converting %v to int: %w", *value, err)
	}
	return models.NewOptionalInt(vv), nil
}

func (t changesetTranslator) bool(value *bool) bool {
	if value == nil {
		return false
	}

	return *value
}

func (t changesetTranslator) optionalBool(value *bool, field string) models.OptionalBool {
	if !t.hasField(field) {
		return models.OptionalBool{}
	}

	return models.NewOptionalBoolPtr(value)
}

func (t changesetTranslator) optionalFloat64(value *float64, field string) models.OptionalFloat64 {
	if !t.hasField(field) {
		return models.OptionalFloat64{}
	}

	return models.NewOptionalFloat64Ptr(value)
}

func (t changesetTranslator) fileIDPtrFromString(value *string) (*models.FileID, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	vv, err := strconv.Atoi(*value)
	if err != nil {
		return nil, fmt.Errorf("converting %v to int: %w", *value, err)
	}

	id := models.FileID(vv)
	return &id, nil
}

func (t changesetTranslator) fileIDSliceFromStringSlice(value []string) ([]models.FileID, error) {
	ints, err := stringslice.StringSliceToIntSlice(value)
	if err != nil {
		return nil, err
	}

	fileIDs := make([]models.FileID, len(ints))
	for i, v := range ints {
		fileIDs[i] = models.FileID(v)
	}

	return fileIDs, nil
}

func (t changesetTranslator) relatedIds(value []string) (models.RelatedIDs, error) {
	ids, err := stringslice.StringSliceToIntSlice(value)
	if err != nil {
		return models.RelatedIDs{}, err
	}

	return models.NewRelatedIDs(ids), nil
}

func (t changesetTranslator) updateIds(value []string, field string) (*models.UpdateIDs, error) {
	if !t.hasField(field) {
		return nil, nil
	}

	ids, err := stringslice.StringSliceToIntSlice(value)
	if err != nil {
		return nil, err
	}

	return &models.UpdateIDs{
		IDs:  ids,
		Mode: models.RelationshipUpdateModeSet,
	}, nil
}

func (t changesetTranslator) updateIdsBulk(value *BulkUpdateIds, field string) (*models.UpdateIDs, error) {
	if !t.hasField(field) || value == nil {
		return nil, nil
	}

	ids, err := stringslice.StringSliceToIntSlice(value.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids [%v]: %w", value.Ids, err)
	}

	return &models.UpdateIDs{
		IDs:  ids,
		Mode: value.Mode,
	}, nil
}

func (t changesetTranslator) optionalURLs(value []string, legacyValue *string) *models.UpdateStrings {
	const (
		legacyField = "url"
		field       = "urls"
	)

	// prefer urls over url
	if t.hasField(field) {
		return t.updateStrings(value, field)
	} else if t.hasField(legacyField) {
		var valueSlice []string
		if legacyValue != nil {
			valueSlice = []string{*legacyValue}
		}
		return t.updateStrings(valueSlice, legacyField)
	}

	return nil
}

func (t changesetTranslator) optionalURLsBulk(value *BulkUpdateStrings, legacyValue *string) *models.UpdateStrings {
	const (
		legacyField = "url"
		field       = "urls"
	)

	// prefer urls over url
	if t.hasField("urls") {
		return t.updateStringsBulk(value, field)
	} else if t.hasField(legacyField) {
		var valueSlice []string
		if legacyValue != nil {
			valueSlice = []string{*legacyValue}
		}
		return t.updateStrings(valueSlice, legacyField)
	}

	return nil
}

func (t changesetTranslator) updateStrings(value []string, field string) *models.UpdateStrings {
	if !t.hasField(field) {
		return nil
	}

	// Trim whitespace from each string
	trimmedValues := make([]string, len(value))
	for i, v := range value {
		trimmedValues[i] = strings.TrimSpace(v)
	}

	return &models.UpdateStrings{
		Values: trimmedValues,
		Mode:   models.RelationshipUpdateModeSet,
	}
}

func (t changesetTranslator) updateStringsBulk(value *BulkUpdateStrings, field string) *models.UpdateStrings {
	if !t.hasField(field) || value == nil {
		return nil
	}

	// Trim whitespace from each string
	trimmedValues := make([]string, len(value.Values))
	for i, v := range value.Values {
		trimmedValues[i] = strings.TrimSpace(v)
	}

	return &models.UpdateStrings{
		Values: trimmedValues,
		Mode:   value.Mode,
	}
}

func (t changesetTranslator) updateStashIDs(value models.StashIDInputs, field string) *models.UpdateStashIDs {
	if !t.hasField(field) {
		return nil
	}

	return &models.UpdateStashIDs{
		StashIDs: value.ToStashIDs(),
		Mode:     models.RelationshipUpdateModeSet,
	}
}

func (t changesetTranslator) relatedGroupsFromMovies(value []models.SceneMovieInput) (models.RelatedGroups, error) {
	groupsScenes, err := models.GroupsScenesFromInput(value)
	if err != nil {
		return models.RelatedGroups{}, err
	}

	return models.NewRelatedGroups(groupsScenes), nil
}

func groupsScenesFromGroupInput(input []models.SceneGroupInput) ([]models.GroupsScenes, error) {
	ret := make([]models.GroupsScenes, len(input))

	for i, v := range input {
		mID, err := strconv.Atoi(v.GroupID)
		if err != nil {
			return nil, fmt.Errorf("invalid group ID: %s", v.GroupID)
		}

		ret[i] = models.GroupsScenes{
			GroupID:    mID,
			SceneIndex: v.SceneIndex,
		}
	}

	return ret, nil
}

func (t changesetTranslator) relatedGroups(value []models.SceneGroupInput) (models.RelatedGroups, error) {
	groupsScenes, err := groupsScenesFromGroupInput(value)
	if err != nil {
		return models.RelatedGroups{}, err
	}

	return models.NewRelatedGroups(groupsScenes), nil
}

func (t changesetTranslator) updateGroupIDsFromMovies(value []models.SceneMovieInput, field string) (*models.UpdateGroupIDs, error) {
	if !t.hasField(field) {
		return nil, nil
	}

	groupsScenes, err := models.GroupsScenesFromInput(value)
	if err != nil {
		return nil, err
	}

	return &models.UpdateGroupIDs{
		Groups: groupsScenes,
		Mode:   models.RelationshipUpdateModeSet,
	}, nil
}

func (t changesetTranslator) updateGroupIDs(value []models.SceneGroupInput, field string) (*models.UpdateGroupIDs, error) {
	if !t.hasField(field) {
		return nil, nil
	}

	groupsScenes, err := groupsScenesFromGroupInput(value)
	if err != nil {
		return nil, err
	}

	return &models.UpdateGroupIDs{
		Groups: groupsScenes,
		Mode:   models.RelationshipUpdateModeSet,
	}, nil
}

func (t changesetTranslator) updateGroupIDsBulk(value *BulkUpdateIds, field string) (*models.UpdateGroupIDs, error) {
	if !t.hasField(field) || value == nil {
		return nil, nil
	}

	ids, err := stringslice.StringSliceToIntSlice(value.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids [%v]: %w", value.Ids, err)
	}

	groups := make([]models.GroupsScenes, len(ids))
	for i, id := range ids {
		groups[i] = models.GroupsScenes{GroupID: id}
	}

	return &models.UpdateGroupIDs{
		Groups: groups,
		Mode:   value.Mode,
	}, nil
}

func groupsDescriptionsFromGroupInput(input []*GroupDescriptionInput) ([]models.GroupIDDescription, error) {
	ret := make([]models.GroupIDDescription, len(input))

	for i, v := range input {
		gID, err := strconv.Atoi(v.GroupID)
		if err != nil {
			return nil, fmt.Errorf("invalid group ID: %s", v.GroupID)
		}

		ret[i] = models.GroupIDDescription{
			GroupID: gID,
		}
		if v.Description != nil {
			ret[i].Description = strings.TrimSpace(*v.Description)
		}
	}

	return ret, nil
}

func (t changesetTranslator) groupIDDescriptions(value []*GroupDescriptionInput) (models.RelatedGroupDescriptions, error) {
	groupsScenes, err := groupsDescriptionsFromGroupInput(value)
	if err != nil {
		return models.RelatedGroupDescriptions{}, err
	}

	return models.NewRelatedGroupDescriptions(groupsScenes), nil
}

func (t changesetTranslator) updateGroupIDDescriptions(value []*GroupDescriptionInput, field string) (*models.UpdateGroupDescriptions, error) {
	if !t.hasField(field) {
		return nil, nil
	}

	groupsScenes, err := groupsDescriptionsFromGroupInput(value)
	if err != nil {
		return nil, err
	}

	return &models.UpdateGroupDescriptions{
		Groups: groupsScenes,
		Mode:   models.RelationshipUpdateModeSet,
	}, nil
}

func (t changesetTranslator) updateGroupIDDescriptionsBulk(value *BulkUpdateGroupDescriptionsInput, field string) (*models.UpdateGroupDescriptions, error) {
	if !t.hasField(field) || value == nil {
		return nil, nil
	}

	groups, err := groupsDescriptionsFromGroupInput(value.Groups)
	if err != nil {
		return nil, err
	}

	return &models.UpdateGroupDescriptions{
		Groups: groups,
		Mode:   value.Mode,
	}, nil
}
