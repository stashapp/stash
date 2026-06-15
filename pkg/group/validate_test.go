package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAliases(t *testing.T) {
	tests := []struct {
		name    string
		aliases []string
		want    error
	}{
		{"nil aliases", nil, nil},
		{"empty aliases", []string{}, nil},
		{"valid single alias", []string{"alias one"}, nil},
		{"valid multiple aliases", []string{"alias one", "alias two"}, nil},
		{"internal spaces valid", []string{"foo bar baz"}, nil},
		{"empty string alias", []string{""}, ErrEmptyAlias},
		{"whitespace-only alias", []string{"   "}, ErrEmptyAlias},
		{"leading space", []string{" alias"}, ErrAliasNotTrimmed},
		{"trailing space", []string{"alias "}, ErrAliasNotTrimmed},
		{"valid then untrimmed", []string{"good", " bad"}, ErrAliasNotTrimmed},
		{"duplicate aliases", []string{"alias", "alias"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAliases(testCtx, tt.aliases)
			assert.Equal(t, tt.want, got)
		})
	}
}
