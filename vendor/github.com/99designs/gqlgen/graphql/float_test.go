package graphql

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat(t *testing.T) {
	assert.Equal(t, "123", m2s(MarshalFloat(123)))
	assert.Equal(t, "1.2345678901", m2s(MarshalFloat(1.2345678901)))
	assert.Equal(t, "1.2345678901234567", m2s(MarshalFloat(1.234567890123456789)))
	assert.Equal(t, "1.2e+20", m2s(MarshalFloat(1.2e+20)))
	assert.Equal(t, "1.2e-20", m2s(MarshalFloat(1.2e-20)))
}

func m2s(m Marshaler) string {
	var b bytes.Buffer
	m.MarshalGQL(&b)
	return b.String()
}
