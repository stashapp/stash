package models

import "context"

type CustomFieldMap map[string]interface{}

type CustomFieldsInput struct {
	// If populated, the entire custom fields map will be replaced with this value
	Full map[string]interface{} `json:"full"`
	// If populated, only the keys in this map will be updated
	Partial map[string]interface{} `json:"partial"`
}

type CustomFieldsReader interface {
	GetCustomFieldsBulk(ctx context.Context, ids []int) ([]CustomFieldMap, error)
}
