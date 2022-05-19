package graphql

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestJsonWriter(t *testing.T) {
	obj := NewFieldSet([]CollectedField{
		{Field: &ast.Field{Alias: "test"}},
		{Field: &ast.Field{Alias: "array"}},
		{Field: &ast.Field{Alias: "emptyArray"}},
		{Field: &ast.Field{Alias: "child"}},
	})
	obj.Values[0] = MarshalInt(10)

	obj.Values[1] = &Array{
		MarshalInt(1),
		MarshalString("2"),
		MarshalBoolean(true),
		False,
		Null,
		MarshalFloat(1.3),
		True,
	}

	obj.Values[2] = &Array{}

	child2 := NewFieldSet([]CollectedField{
		{Field: &ast.Field{Alias: "child"}},
	})
	child2.Values[0] = Null

	child1 := NewFieldSet([]CollectedField{
		{Field: &ast.Field{Alias: "child"}},
	})
	child1.Values[0] = child2

	obj.Values[3] = child1

	b := &bytes.Buffer{}
	obj.MarshalGQL(b)

	require.Equal(t, `{"test":10,"array":[1,"2",true,false,null,1.3,true],"emptyArray":[],"child":{"child":{"child":null}}}`, b.String())
}
