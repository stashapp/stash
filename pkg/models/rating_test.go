package models

import (
	"testing"
)

func TestRating100To5(t *testing.T) {
	tests := []struct {
		name      string
		rating100 int
		want      int
	}{
		{"20", 20, 1},
		{"100", 100, 5},
		{"1", 1, 1},
		{"10", 10, 1},
		{"11", 11, 1},
		{"21", 21, 1},
		{"31", 31, 2},
		{"0", 0, 1},
		{"-100", -100, 1},
		{"120", 120, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rating100To5(tt.rating100); got != tt.want {
				t.Errorf("Rating100To5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRating5To100(t *testing.T) {
	tests := []struct {
		name    string
		rating5 int
		want    int
	}{
		{"1", 1, 20},
		{"5", 5, 100},
		{"2", 2, 40},
		{"3", 3, 60},
		{"4", 4, 80},
		{"6", 6, 100},
		{"0", 0, 20},
		{"-1", -1, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rating5To100(tt.rating5); got != tt.want {
				t.Errorf("Rating5To100() = %v, want %v", got, tt.want)
			}
		})
	}
}
