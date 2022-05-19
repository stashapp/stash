package graphql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		assert.Equal(t, "123", m2s(MarshalUint(123)))
	})

	t.Run("unmarshal", func(t *testing.T) {
		assert.Equal(t, uint(123), mustUnmarshalUint(123))
		assert.Equal(t, uint(123), mustUnmarshalUint(int64(123)))
		assert.Equal(t, uint(123), mustUnmarshalUint(json.Number("123")))
		assert.Equal(t, uint(123), mustUnmarshalUint("123"))
	})
}

func mustUnmarshalUint(v interface{}) uint {
	res, err := UnmarshalUint(v)
	if err != nil {
		panic(err)
	}
	return res
}

func TestUint32(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		assert.Equal(t, "123", m2s(MarshalUint32(123)))
	})

	t.Run("unmarshal", func(t *testing.T) {
		assert.Equal(t, uint32(123), mustUnmarshalUint32(123))
		assert.Equal(t, uint32(123), mustUnmarshalUint32(int64(123)))
		assert.Equal(t, uint32(123), mustUnmarshalUint32(json.Number("123")))
		assert.Equal(t, uint32(123), mustUnmarshalUint32("123"))
	})
}

func mustUnmarshalUint32(v interface{}) uint32 {
	res, err := UnmarshalUint32(v)
	if err != nil {
		panic(err)
	}
	return res
}

func TestUint64(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		assert.Equal(t, "123", m2s(MarshalUint64(123)))
	})

	t.Run("unmarshal", func(t *testing.T) {
		assert.Equal(t, uint64(123), mustUnmarshalUint64(123))
		assert.Equal(t, uint64(123), mustUnmarshalUint64(int64(123)))
		assert.Equal(t, uint64(123), mustUnmarshalUint64(json.Number("123")))
		assert.Equal(t, uint64(123), mustUnmarshalUint64("123"))
	})
}

func mustUnmarshalUint64(v interface{}) uint64 {
	res, err := UnmarshalUint64(v)
	if err != nil {
		panic(err)
	}
	return res
}
