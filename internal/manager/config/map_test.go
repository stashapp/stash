package config

import (
	"testing"
)

func Test_toSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		v    string
		want string
	}{
		{
			"basic",
			"basic",
			"basic",
		},
		{
			"two words",
			"twoWords",
			"two_words",
		},
		{
			"three word value",
			"threeWordValue",
			"three_word_value",
		},
		{
			"snake case",
			"snake_case",
			"snake_case",
		},
		{
			"double capital",
			"doubleCApital",
			"double_capital",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toSnakeCase(tt.v); got != tt.want {
				t.Errorf("toSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		v    string
		want string
	}{
		{
			"basic",
			"basic",
			"basic",
		},
		{
			"two words",
			"two_words",
			"twoWords",
		},
		{
			"three word value",
			"three_word_value",
			"threeWordValue",
		},
		{
			"camel case",
			"camelCase",
			"camelCase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromSnakeCase(tt.v); got != tt.want {
				t.Errorf("fromSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
