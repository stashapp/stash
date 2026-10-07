package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string { return &s }

func TestDeprecatedAliasList(t *testing.T) {
	assert.Nil(t, deprecatedAliasList(nil))
	assert.Equal(t, []string{}, deprecatedAliasList(strPtr("")))
	assert.Equal(t, []string{}, deprecatedAliasList(strPtr("   ")))
	assert.Equal(t, []string{"a"}, deprecatedAliasList(strPtr(" a ")))
	assert.Equal(t, []string{"a, b"}, deprecatedAliasList(strPtr("a, b")))
}
