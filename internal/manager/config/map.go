package config

import (
	"bytes"
	"unicode"

	"github.com/spf13/cast"
)

// HACK: viper changes map keys to case insensitive values, so the workaround is to
// convert the map to use snake-case keys

// toSnakeCase converts a string to snake_case
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

// copyAndInsensitiviseMap behaves like insensitiviseMap, but creates a copy of
// any map it makes case insensitive.
func toSnakeCaseMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		adjKey := toSnakeCase(key)

		switch v := val.(type) {
		case map[string]interface{}:
			nm[adjKey] = toSnakeCaseMap(v)
		default:
			nm[adjKey] = val
		}
	}

	return nm
}

// convertMapValue converts values into something that can be marshalled in JSON
// This means converting map[interface{}]interface{} to map[string]interface{} where ever
// encountered.
func convertMapValue(val interface{}) interface{} {
	switch v := val.(type) {
	case map[interface{}]interface{}:
		ret := cast.ToStringMap(v)
		for k, vv := range ret {
			adjKey := fromSnakeCase(k)
			ret[adjKey] = convertMapValue(vv)
		}
		return ret
	case map[string]interface{}:
		ret := make(map[string]interface{})
		for k, vv := range v {
			adjKey := fromSnakeCase(k)
			ret[adjKey] = convertMapValue(vv)
		}
		return ret
	case []interface{}:
		ret := make([]interface{}, len(v))
		for i, vv := range v {
			ret[i] = convertMapValue(vv)
		}
		return ret
	default:
		return v
	}
}

func fromSnakeCaseMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		adjKey := fromSnakeCase(key)
		nm[adjKey] = convertMapValue(val)
	}

	return nm
}
