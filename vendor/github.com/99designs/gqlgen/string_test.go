package graphql

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, `"hello"`, doStrMarshal("hello"))
	assert.Equal(t, `"he\tllo"`, doStrMarshal("he\tllo"))
	assert.Equal(t, `"he\tllo"`, doStrMarshal("he	llo"))
	assert.Equal(t, `"he\nllo"`, doStrMarshal("he\nllo"))
	assert.Equal(t, `"he\r\nllo"`, doStrMarshal("he\r\nllo"))
	assert.Equal(t, `"he\\llo"`, doStrMarshal(`he\llo`))
	assert.Equal(t, `"quotes\"nested\"in\"quotes\""`, doStrMarshal(`quotes"nested"in"quotes"`))
	assert.Equal(t, `"\u0000"`, doStrMarshal("\u0000"))
	assert.Equal(t, `"\u0000"`, doStrMarshal("\u0000"))
	assert.Equal(t, "\"\U000fe4ed\"", doStrMarshal("\U000fe4ed"))
}

func doStrMarshal(s string) string {
	var buf bytes.Buffer
	MarshalString(s).MarshalGQL(&buf)
	return buf.String()
}
