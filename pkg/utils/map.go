package utils

import (
	"encoding/json"
	"strings"
)

// NestedMap is a map that supports nested keys.
// It is expected that the nested maps are of type map[string]interface{}
type NestedMap map[string]interface{}

func (m NestedMap) Get(key string) (interface{}, bool) {
	fields := strings.Split(key, ".")

	current := m

	for _, f := range fields[:len(fields)-1] {
		v, found := current[f]
		if !found {
			return nil, false
		}

		current, _ = v.(map[string]interface{})
		if current == nil {
			return nil, false
		}
	}

	ret, found := current[fields[len(fields)-1]]
	return ret, found
}

func (m NestedMap) Set(key string, value interface{}) {
	fields := strings.Split(key, ".")

	current := m

	for _, f := range fields[:len(fields)-1] {
		v, ok := current[f].(map[string]interface{})
		if !ok {
			v = make(map[string]interface{})
			current[f] = v
		}

		current = v
	}

	current[fields[len(fields)-1]] = value
}

func (m NestedMap) Delete(key string) {
	fields := strings.Split(key, ".")

	current := m

	for _, f := range fields[:len(fields)-1] {
		v, ok := current[f].(map[string]interface{})
		if !ok {
			return
		}

		current = v
	}

	delete(current, fields[len(fields)-1])
}

// MergeMaps merges src into dest. If a key exists in both maps, the value from src is used.
func MergeMaps(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		if _, ok := dest[k]; ok {
			if srcMap, ok := v.(map[string]interface{}); ok {
				if destMap, ok := dest[k].(map[string]interface{}); ok {
					MergeMaps(destMap, srcMap)
					continue
				}
			}
		}

		dest[k] = v
	}
}

// ConvertMapJSONNumbers converts all JSON numbers in a map to either float64 or int64.
func ConvertMapJSONNumbers(m map[string]interface{}) (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	for k, v := range m {
		if n, ok := v.(json.Number); ok {
			ret[k] = JSONNumberToNumber(n)
		} else if mm, ok := v.(map[string]interface{}); ok {
			ret[k] = ConvertMapJSONNumbers(mm)
		} else {
			ret[k] = v
		}
	}

	return ret
}
