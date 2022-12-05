package jsonschema

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadPerformer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Performer
		wantErr bool
	}{
		{
			name: "alias list",
			input: `
{
	"aliases": ["alias1", "alias2"]
}`,
			want: Performer{
				Aliases: []string{"alias1", "alias2"},
			},
			wantErr: false,
		},
		{
			name: "alias string list",
			input: `
{
	"aliases": "alias1, alias2"
}`,
			want: Performer{
				Aliases: []string{"alias1", "alias2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := loadPerformer(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadPerformer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, &tt.want, got)
		})
	}
}
