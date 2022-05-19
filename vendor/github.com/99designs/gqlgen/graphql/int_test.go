package graphql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		assert.Equal(t, "123", m2s(MarshalInt(123)))
	})

	t.Run("unmarshal", func(t *testing.T) {
		assert.Equal(t, 123, mustUnmarshalInt(123))
		assert.Equal(t, 123, mustUnmarshalInt(int64(123)))
		assert.Equal(t, 123, mustUnmarshalInt(json.Number("123")))
		assert.Equal(t, 123, mustUnmarshalInt("123"))
	})
}

func mustUnmarshalInt(v interface{}) int {
	res, err := UnmarshalInt(v)
	if err != nil {
		panic(err)
	}
	return res
}

func TestInt32(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		assert.Equal(t, "123", m2s(MarshalInt32(123)))
	})

	t.Run("unmarshal", func(t *testing.T) {
		assert.Equal(t, int32(123), mustUnmarshalInt32(123))
		assert.Equal(t, int32(123), mustUnmarshalInt32(int64(123)))
		assert.Equal(t, int32(123), mustUnmarshalInt32(json.Number("123")))
		assert.Equal(t, int32(123), mustUnmarshalInt32("123"))
	})
}

func mustUnmarshalInt32(v interface{}) int32 {
	res, err := UnmarshalInt32(v)
	if err != nil {
		panic(err)
	}
	return res
}

func TestInt64(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		assert.Equal(t, "123", m2s(MarshalInt32(123)))
	})

	t.Run("unmarshal", func(t *testing.T) {
		assert.Equal(t, int64(123), mustUnmarshalInt64(123))
		assert.Equal(t, int64(123), mustUnmarshalInt64(int64(123)))
		assert.Equal(t, int64(123), mustUnmarshalInt64(json.Number("123")))
		assert.Equal(t, int64(123), mustUnmarshalInt64("123"))
	})
}

func mustUnmarshalInt64(v interface{}) int64 {
	res, err := UnmarshalInt64(v)
	if err != nil {
		panic(err)
	}
	return res
}
