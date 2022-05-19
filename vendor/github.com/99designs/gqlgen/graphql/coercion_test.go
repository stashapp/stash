package graphql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoerceList(t *testing.T) {

	mapInput := map[string]interface{}{
		"test": "value",
		"nested": map[string]interface{}{
			"nested": true,
		},
	}

	jsonNumber := json.Number("12")

	assert.Equal(t, []interface{}{"test", "values"}, CoerceList([]interface{}{"test", "values"}))
	assert.Equal(t, []interface{}{"test"}, CoerceList("test"))
	assert.Equal(t, []interface{}{"test"}, CoerceList([]string{"test"}))
	assert.Equal(t, []interface{}{3}, CoerceList([]int{3}))
	assert.Equal(t, []interface{}{3}, CoerceList(3))
	assert.Equal(t, []interface{}{int32(3)}, CoerceList([]int32{3}))
	assert.Equal(t, []interface{}{int64(2)}, CoerceList([]int64{2}))
	assert.Equal(t, []interface{}{float32(3.14)}, CoerceList([]float32{3.14}))
	assert.Equal(t, []interface{}{3.14}, CoerceList([]float64{3.14}))
	assert.Equal(t, []interface{}{jsonNumber}, CoerceList([]json.Number{jsonNumber}))
	assert.Equal(t, []interface{}{jsonNumber}, CoerceList(jsonNumber))
	assert.Equal(t, []interface{}{true}, CoerceList([]bool{true}))
	assert.Equal(t, []interface{}{mapInput}, CoerceList(mapInput))
	assert.Equal(t, []interface{}{mapInput}, CoerceList([]interface{}{mapInput}))
	assert.Equal(t, []interface{}{mapInput}, CoerceList([]map[string]interface{}{mapInput}))
	assert.Empty(t, CoerceList(nil))

}
