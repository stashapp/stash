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
	cap := false
	for i, c := range v {
		switch {
		case c == '_' && i > 0:
			cap = true
		case cap:
			buf.WriteRune(unicode.ToUpper(c))
			cap = false
		default:
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
		case map[interface{}]interface{}:
			nm[adjKey] = toSnakeCaseMap(cast.ToStringMap(v))
		case map[string]interface{}:
			nm[adjKey] = toSnakeCaseMap(v)
		default:
			nm[adjKey] = v
		}
	}

	return nm
}

func fromSnakeCaseMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		adjKey := fromSnakeCase(key)
		switch v := val.(type) {
		case map[interface{}]interface{}:
			nm[adjKey] = fromSnakeCaseMap(cast.ToStringMap(v))
		case map[string]interface{}:
			nm[adjKey] = fromSnakeCaseMap(v)
		default:
			nm[adjKey] = v
		}
	}

	return nm
}
