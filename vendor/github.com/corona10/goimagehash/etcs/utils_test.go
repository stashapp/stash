// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package etcs

import (
	"testing"
)

func TestMeanPixels(t *testing.T) {
	for _, tt := range []struct {
		pixels   []float64
		expected float64
	}{
		{[]float64{0, 0, 0, 0}, 0},
		{[]float64{1, 2, 3, 4}, 2.5},
	} {
		pixels := tt.pixels
		result := MeanOfPixels(pixels)
		if result != tt.expected {
			t.Errorf("Mean of %v is expected as %v but got %v.", pixels, tt.expected, result)
		}
	}
}

func TestMedianPixels(t *testing.T) {
	for _, tt := range []struct {
		pixels   []float64
		expected float64
	}{
		{[]float64{0, 0, 0, 0}, 0},
		{[]float64{1}, 1},
		{[]float64{1, 2, 3, 4}, 2.5},
		{[]float64{5, 3, 1, 7, 9}, 5},
		{[]float64{98.3, 33.4, 105.44, 1500.4, 22.5, 66.6}, 82.44999999999999},
	} {
		pixels := tt.pixels
		result := MedianOfPixels(pixels)
		if result != tt.expected {
			t.Errorf("Median of %v is expected as %v but got %v.", pixels, tt.expected, result)
		}
	}
}
