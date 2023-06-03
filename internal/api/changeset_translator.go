package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/models"
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

func (t changesetTranslator) string(value *string, field string) string {
	if value == nil {
		return ""
	}

	return *value
}

func (t changesetTranslator) optionalString(value *string, field string) models.OptionalString {
	if !t.hasField(field) {
		return models.OptionalString{}
	}

	return models.NewOptionalStringPtr(value)
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

func (t changesetTranslator) datePtr(value *string, field string) (*models.Date, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	date, err := models.ParseDate(*value)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func (t changesetTranslator) intPtrFromString(value *string, field string) (*int, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	vv, err := strconv.Atoi(*value)
	if err != nil {
		return nil, fmt.Errorf("converting %v to int: %w", *value, err)
	}
	return &vv, nil
}

func (t changesetTranslator) ratingConversionInt(legacyValue *int, rating100Value *int) *int {
	const (
		legacyField    = "rating"
		rating100Field = "rating100"
	)

	legacyRating := t.optionalInt(legacyValue, legacyField)
	if legacyRating.Set && !(legacyRating.Null) {
		ret := int(models.Rating5To100(int(legacyRating.Value)))
		return &ret
	}

	o := t.optionalInt(rating100Value, rating100Field)
	if o.Set && !(o.Null) {
		return &o.Value
	}

	return nil
}

func (t changesetTranslator) ratingConversionOptional(legacyValue *int, rating100Value *int) models.OptionalInt {
	const (
		legacyField    = "rating"
		rating100Field = "rating100"
	)

	legacyRating := t.optionalInt(legacyValue, legacyField)
	if legacyRating.Set && !(legacyRating.Null) {
		legacyRating.Value = int(models.Rating5To100(int(legacyRating.Value)))
		return legacyRating
	}
	return t.optionalInt(rating100Value, rating100Field)
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

func (t changesetTranslator) bool(value *bool, field string) bool {
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
