package config

import (
	"bytes"
	"unicode"

	"github.com/spf13/cast"
)

// HACK: viper changes map keys to case insensitive values, so the workaround is to
// convert the map to use snake-case keys

// toSnakeCase converts a string from camelCase to snake_case
// NOTE: a double capital will be converted in a way that will yield a different result
// when converted back to camel case.
// For example: someIDs => some_ids => someIds
func toSnakeCase(v string) string {
	var buf bytes.Buffer
	underscored := false
	for i, c := range v {
		if !underscored && unicode.IsUpper(c) && i > 0 {
			buf.WriteByte('_')
			underscored = true
		} else {
			underscored = false
		}

		buf.WriteRune(unicode.ToLower(c))
	}
	return buf.String()
}

// fromSnakeCase converts a string from snake_case to camelCase
func fromSnakeCase(v string) string {
	var buf bytes.Buffer
	leadingUnderscore := true
	capvar := false
	for i, c := range v {
		switch {
		case c == '_' && !leadingUnderscore && i > 0:
			capvar = true
		case c == '_' && leadingUnderscore:
			buf.WriteRune(c)
		case capvar:
			buf.WriteRune(unicode.ToUpper(c))
			capvar = false
		default:
			leadingUnderscore = false
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

// fromSnakeCaseMap recursively converts a map using snake_case keys to camelCase keys
func fromSnakeCaseMap(m map[string]interface{}) map[string]interface{} {
	return fromSnakeCaseValue(m).(map[string]interface{})
}

func fromSnakeCaseValue(val interface{}) interface{} {
	switch v := val.(type) {
	case map[interface{}]interface{}:
		ret := cast.ToStringMap(v)
		for k, vv := range ret {
			adjKey := fromSnakeCase(k)
			ret[adjKey] = fromSnakeCaseValue(vv)
		}
		return ret
	case map[string]interface{}:
		ret := make(map[string]interface{})
		for k, vv := range v {
			adjKey := fromSnakeCase(k)
			ret[adjKey] = fromSnakeCaseValue(vv)
		}
		return ret
	case []interface{}:
		ret := make([]interface{}, len(v))
		for i, vv := range v {
			ret[i] = fromSnakeCaseValue(vv)
		}
		return ret
	default:
		return v
	}
}

// toSnakeCaseMap recursively converts a map using camelCase keys to snake_case keys
func toSnakeCaseMap(m map[string]interface{}) map[string]interface{} {
	return toSnakeCaseValue(m).(map[string]interface{})
}

func toSnakeCaseValue(val interface{}) interface{} {
	switch v := val.(type) {
	case map[interface{}]interface{}:
		ret := cast.ToStringMap(v)
		for k, vv := range ret {
			adjKey := toSnakeCase(k)
			ret[adjKey] = toSnakeCaseValue(vv)
		}
		return ret
	case map[string]interface{}:
		ret := make(map[string]interface{})
		for k, vv := range v {
			adjKey := toSnakeCase(k)
			ret[adjKey] = toSnakeCaseValue(vv)
		}
		return ret
	case []interface{}:
		ret := make([]interface{}, len(v))
		for i, vv := range v {
			ret[i] = toSnakeCaseValue(vv)
		}
		return ret
	default:
		return v
	}
}
