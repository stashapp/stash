package graphql

import (
	"encoding/json"
)

// CoerceList applies coercion from a single value to a list.
func CoerceList(v interface{}) []interface{} {
	var vSlice []interface{}
	if v != nil {
		switch v := v.(type) {
		case []interface{}:
			// already a slice no coercion required
			vSlice = v
		case []string:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []json.Number:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []bool:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []map[string]interface{}:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []float64:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []float32:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []int:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []int32:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		case []int64:
			if len(v) > 0 {
				vSlice = []interface{}{v[0]}
			}
		default:
			vSlice = []interface{}{v}
		}
	}
	return vSlice
}
