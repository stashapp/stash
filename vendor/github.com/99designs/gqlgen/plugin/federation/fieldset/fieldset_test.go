package fieldset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnnestedWithoutPrefix(t *testing.T) {
	fieldSet := New("foo bar", nil)

	require.Len(t, fieldSet, 2)

	require.Len(t, fieldSet[0], 1)
	require.Equal(t, "foo", fieldSet[0][0])

	require.Len(t, fieldSet[1], 1)
	require.Equal(t, "bar", fieldSet[1][0])
}

func TestNestedWithoutPrefix(t *testing.T) {
	fieldSet := New("foo bar { baz} a b {c{d}}e", nil)

	require.Len(t, fieldSet, 5)

	require.Len(t, fieldSet[0], 1)
	require.Equal(t, "foo", fieldSet[0][0])

	require.Len(t, fieldSet[1], 2)
	require.Equal(t, "bar", fieldSet[1][0])
	require.Equal(t, "baz", fieldSet[1][1])

	require.Len(t, fieldSet[2], 1)
	require.Equal(t, "a", fieldSet[2][0])

	require.Len(t, fieldSet[3], 3)
	require.Equal(t, "b", fieldSet[3][0])
	require.Equal(t, "c", fieldSet[3][1])
	require.Equal(t, "d", fieldSet[3][2])

	require.Len(t, fieldSet[4], 1)
	require.Equal(t, "e", fieldSet[4][0])
}

func TestWithPrefix(t *testing.T) {
	fieldSet := New("foo bar{id}", []string{"prefix"})

	require.Len(t, fieldSet, 2)

	require.Len(t, fieldSet[0], 2)
	require.Equal(t, "prefix", fieldSet[0][0])
	require.Equal(t, "foo", fieldSet[0][1])

	require.Len(t, fieldSet[1], 3)
	require.Equal(t, "prefix", fieldSet[1][0])
	require.Equal(t, "bar", fieldSet[1][1])
	require.Equal(t, "id", fieldSet[1][2])
}

func TestInvalid(t *testing.T) {
	require.Panics(t, func() {
		New("foo bar{baz", nil)
	})
}

func TestToGo(t *testing.T) {
	require.Equal(t, Field{"foo"}.ToGo(), "Foo")
	require.Equal(t, Field{"foo", "bar"}.ToGo(), "FooBar")
	require.Equal(t, Field{"bar", "id"}.ToGo(), "BarID")
}

func TestToGoPrivate(t *testing.T) {
	require.Equal(t, Field{"foo"}.ToGoPrivate(), "foo")
	require.Equal(t, Field{"foo", "bar"}.ToGoPrivate(), "fooBar")
	require.Equal(t, Field{"bar", "id"}.ToGoPrivate(), "barID")
}
