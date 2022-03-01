package api

import (
	"context"
	"database/sql"
	"strconv"

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
	args := getArgumentMap(ctx)

	input := args[updateInputField]
	var ret map[string]interface{}
	if input != nil {
		ret, _ = input.(map[string]interface{})
	}

	if ret == nil {
		ret = make(map[string]interface{})
	}

	return ret
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

func (t changesetTranslator) nullString(value *string, field string) *sql.NullString {
	if !t.hasField(field) {
		return nil
	}

	ret := &sql.NullString{}

	if value != nil {
		ret.String = *value
		ret.Valid = true
	}

	return ret
}

func (t changesetTranslator) sqliteDate(value *string, field string) *models.SQLiteDate {
	if !t.hasField(field) {
		return nil
	}

	ret := &models.SQLiteDate{}

	if value != nil {
		ret.String = *value
		ret.Valid = true
	}

	return ret
}

func (t changesetTranslator) nullInt64(value *int, field string) *sql.NullInt64 {
	if !t.hasField(field) {
		return nil
	}

	ret := &sql.NullInt64{}

	if value != nil {
		ret.Int64 = int64(*value)
		ret.Valid = true
	}

	return ret
}

func (t changesetTranslator) nullInt64FromString(value *string, field string) *sql.NullInt64 {
	if !t.hasField(field) {
		return nil
	}

	ret := &sql.NullInt64{}

	if value != nil {
		ret.Int64, _ = strconv.ParseInt(*value, 10, 64)
		ret.Valid = true
	}

	return ret
}

func (t changesetTranslator) nullBool(value *bool, field string) *sql.NullBool {
	if !t.hasField(field) {
		return nil
	}

	ret := &sql.NullBool{}

	if value != nil {
		ret.Bool = *value
		ret.Valid = true
	}

	return ret
}
