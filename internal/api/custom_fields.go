package api

import "github.com/stashapp/stash/pkg/models"

func handleUpdateCustomFields(input models.CustomFieldsInput) models.CustomFieldsInput {
	ret := input
	// convert json.Numbers to int/float
	ret.Full = convertMapJSONNumbers(ret.Full)
	ret.Partial = convertMapJSONNumbers(ret.Partial)

	return ret
}
