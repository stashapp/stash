package utils

import (
	"reflect"
	"testing"
)

func TestNotNilFields(t *testing.T) {
	v := "value"
	var zeroStr string

	type testObject struct {
		ptrField      *string `tag:"ptrField"`
		noTagField    *string
		otherTagField *string  `otherTag:"otherTagField"`
		sliceField    []string `tag:"sliceField"`
	}

	type args struct {
		subject interface{}
		tag     string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"basic",
			args{
				testObject{
					ptrField:      &v,
					noTagField:    &v,
					otherTagField: &v,
					sliceField:    []string{v},
				},
				"tag",
			},
			[]string{"ptrField", "sliceField"},
		},
		{
			"empty",
			args{
				testObject{},
				"tag",
			},
			nil,
		},
		{
			"zero values",
			args{
				testObject{
					ptrField:      &zeroStr,
					noTagField:    &zeroStr,
					otherTagField: &zeroStr,
					sliceField:    []string{},
				},
				"tag",
			},
			[]string{"ptrField", "sliceField"},
		},
		{
			"other tag",
			args{
				testObject{
					ptrField:      &v,
					noTagField:    &v,
					otherTagField: &v,
					sliceField:    []string{v},
				},
				"otherTag",
			},
			[]string{"otherTagField"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotNilFields(tt.args.subject, tt.args.tag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotNilFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
