//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSetCustomFields(t *testing.T) {
	performerIdx := performerIdx1WithScene

	mergeCustomFields := func(i map[string]interface{}) map[string]interface{} {
		m := getPerformerCustomFields(performerIdx)
		for k, v := range i {
			m[k] = v
		}
		return m
	}

	tests := []struct {
		name     string
		input    models.CustomFieldsInput
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			"valid full",
			models.CustomFieldsInput{
				Full: map[string]interface{}{
					"key": "value",
				},
			},
			map[string]interface{}{
				"key": "value",
			},
			false,
		},
		{
			"valid partial",
			models.CustomFieldsInput{
				Partial: map[string]interface{}{
					"key": "value",
				},
			},
			mergeCustomFields(map[string]interface{}{
				"key": "value",
			}),
			false,
		},
		{
			"valid partial overwrite",
			models.CustomFieldsInput{
				Partial: map[string]interface{}{
					"real": float64(4.56),
				},
			},
			mergeCustomFields(map[string]interface{}{
				"real": float64(4.56),
			}),
			false,
		},
		{
			"valid remove",
			models.CustomFieldsInput{
				Remove: []string{"real"},
			},
			func() map[string]interface{} {
				m := getPerformerCustomFields(performerIdx)
				delete(m, "real")
				return m
			}(),
			false,
		},
		{
			"leading space full",
			models.CustomFieldsInput{
				Full: map[string]interface{}{
					" key": "value",
				},
			},
			nil,
			true,
		},
		{
			"trailing space full",
			models.CustomFieldsInput{
				Full: map[string]interface{}{
					"key ": "value",
				},
			},
			nil,
			true,
		},
		{
			"leading space partial",
			models.CustomFieldsInput{
				Partial: map[string]interface{}{
					" key": "value",
				},
			},
			nil,
			true,
		},
		{
			"trailing space partial",
			models.CustomFieldsInput{
				Partial: map[string]interface{}{
					"key ": "value",
				},
			},
			nil,
			true,
		},
		{
			"big key full",
			models.CustomFieldsInput{
				Full: map[string]interface{}{
					"12345678901234567890123456789012345678901234567890123456789012345": "value",
				},
			},
			nil,
			true,
		},
		{
			"big key partial",
			models.CustomFieldsInput{
				Partial: map[string]interface{}{
					"12345678901234567890123456789012345678901234567890123456789012345": "value",
				},
			},
			nil,
			true,
		},
		{
			"empty key full",
			models.CustomFieldsInput{
				Full: map[string]interface{}{
					"": "value",
				},
			},
			nil,
			true,
		},
		{
			"empty key partial",
			models.CustomFieldsInput{
				Partial: map[string]interface{}{
					"": "value",
				},
			},
			nil,
			true,
		},
		{
			"invalid remove full",
			models.CustomFieldsInput{
				Full: map[string]interface{}{
					"key": "value",
				},
				Remove: []string{"key"},
			},
			nil,
			true,
		},
		{
			"invalid remove partial",
			models.CustomFieldsInput{
				Partial: map[string]interface{}{
					"real": float64(4.56),
				},
				Remove: []string{"real"},
			},
			nil,
			true,
		},
	}

	// use performer custom fields store
	store := db.Performer
	id := performerIDs[performerIdx]

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			err := store.SetCustomFields(ctx, id, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCustomFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			actual, err := store.GetCustomFields(ctx, id)
			if err != nil {
				t.Errorf("GetCustomFields() error = %v", err)
				return
			}

			assert.Equal(tt.expected, actual)
		})
	}
}
