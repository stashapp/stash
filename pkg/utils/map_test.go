package utils

import (
	"reflect"
	"testing"
)

// func TestNestedMap_Get(t *testing.T) {
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name  string
// 		m     NestedMap
// 		args  args
// 		want  interface{}
// 		want1 bool
// 	}{

// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := tt.m.Get(tt.args.key)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NestedMap.Get() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("NestedMap.Get() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

func TestNestedMapGet(t *testing.T) {
	m := NestedMap{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": "qux",
			},
		},
	}

	tests := []struct {
		name  string
		key   string
		want  interface{}
		found bool
	}{
		{
			name:  "Get a value from a nested map",
			key:   "foo.bar.baz",
			want:  "qux",
			found: true,
		},
		{
			name:  "Get a value from a nested map with a missing key",
			key:   "foo.bar.quux",
			want:  nil,
			found: false,
		},
		{
			name:  "Get a value from a nested map with a missing key",
			key:   "foo.quux.baz",
			want:  nil,
			found: false,
		},
		{
			name:  "Get a value from a nested map with a missing key",
			key:   "quux.bar.baz",
			want:  nil,
			found: false,
		},
		{
			name:  "Get a value from a nested map with a missing key",
			key:   "foo.bar",
			want:  map[string]interface{}{"baz": "qux"},
			found: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := m.Get(tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NestedMap.Get() got = %v, want %v", got, tt.want)
			}
			if found != tt.found {
				t.Errorf("NestedMap.Get() found = %v, want %v", found, tt.found)
			}
		})
	}
}

func TestNestedMapSet(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		existing NestedMap
		want     NestedMap
	}{
		{
			name:     "Set a value in a nested map",
			key:      "foo.bar.baz",
			existing: NestedMap{},
			want: NestedMap{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"baz": "qux",
					},
				},
			},
		},
		{
			name: "Overwrite existing value",
			key:  "foo.bar",
			existing: NestedMap{
				"foo": map[string]interface{}{
					"bar": "old",
				},
			},
			want: NestedMap{
				"foo": map[string]interface{}{
					"bar": "qux",
				},
			},
		},
		{
			name: "Set a value overwriting a primitive with a nested map",
			key:  "foo.bar",
			existing: NestedMap{
				"foo": "bar",
			},
			want: NestedMap{
				"foo": map[string]interface{}{
					"bar": "qux",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.existing.Set(tt.key, "qux")
			if !reflect.DeepEqual(tt.existing, tt.want) {
				t.Errorf("NestedMap.Set() got = %v, want %v", tt.existing, tt.want)
			}
		})
	}
}
