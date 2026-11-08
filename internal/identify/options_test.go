package identify

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldStrategy_String(t *testing.T) {
	assert.Equal(t, "MERGE", FieldStrategyMerge.String())
}

func TestFieldStrategy_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    FieldStrategy
		wantErr bool
	}{
		{"valid", "OVERWRITE", FieldStrategyOverwrite, false},
		{"invalid type", 1, "", true},
		{"invalid value", "INVALID", "INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got FieldStrategy
			err := got.UnmarshalGQL(tt.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFieldStrategy_MarshalGQL(t *testing.T) {
	var out strings.Builder

	FieldStrategyIgnore.MarshalGQL(&out)

	assert.Equal(t, `"IGNORE"`, out.String())
}
