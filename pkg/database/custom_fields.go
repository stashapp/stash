package database

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

type customFieldsStore interface {
	GetCustomFields(ctx context.Context, id int) (map[string]interface{}, error)
	GetCustomFieldsBulk(ctx context.Context, ids []int) ([]models.CustomFieldMap, error)
	SetCustomFields(ctx context.Context, id int, values models.CustomFieldsInput) error
}
