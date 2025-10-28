package api

import (
	"encoding/json"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

// jsonNumberToNumber converts a JSON number to either a float64 or int64.
func jsonNumberToNumber(n json.Number) interface{} {
	if strings.Contains(string(n), ".") {
		f, _ := n.Float64()
		return f
	}
	ret, _ := n.Int64()
	return ret
}

// anyJSONNumberToNumber converts a JSON number using jsonNumberToNumber, otherwise it returns the existing value
func anyJSONNumberToNumber(v any) any {
	if n, ok := v.(json.Number); ok {
		return jsonNumberToNumber(n)
	}

	return v
}

// ConvertMapJSONNumbers converts all JSON numbers in a map to either float64 or int64.
func convertMapJSONNumbers(m map[string]interface{}) (ret map[string]interface{}) {
	if m == nil {
		return nil
	}

	ret = make(map[string]interface{})
	for k, v := range m {
		if n, ok := v.(json.Number); ok {
			ret[k] = jsonNumberToNumber(n)
		} else if mm, ok := v.(map[string]interface{}); ok {
			ret[k] = convertMapJSONNumbers(mm)
		} else {
			ret[k] = v
		}
	}

	return ret
}

func convertCustomFieldCriterionValues(c models.CustomFieldCriterionInput) models.CustomFieldCriterionInput {
	nv := make([]any, len(c.Value))
	for i, v := range c.Value {
		nv[i] = anyJSONNumberToNumber(v)
	}
	c.Value = nv
	return c
}

func convertCustomFieldCriterionInputJSONNumbers(c []models.CustomFieldCriterionInput) []models.CustomFieldCriterionInput {
	ret := make([]models.CustomFieldCriterionInput, len(c))
	for i, v := range c {
		ret[i] = convertCustomFieldCriterionValues(v)
	}

	return ret
}
