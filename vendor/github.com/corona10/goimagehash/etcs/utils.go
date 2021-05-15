// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etcs

// MeanOfPixels function returns a mean of pixels.
func MeanOfPixels(pixels []float64) float64 {
	m := 0.0
	lens := len(pixels)
	if lens == 0 {
		return 0
	}

	for _, p := range pixels {
		m += p
	}

	return m / float64(lens)
}

// MedianOfPixels function returns a median value of pixels.
// It uses quick selection algorithm.
func MedianOfPixels(pixels []float64) float64 {
	tmp := make([]float64, len(pixels))
	copy(tmp, pixels)
	l := len(tmp)
	pos := l / 2
	v := quickSelectMedian(tmp, 0, l-1, pos)
	return v
}

func quickSelectMedian(sequence []float64, low int, hi int, k int) float64 {
	if low == hi {
		return sequence[k]
	}

	for low < hi {
		pivot := low/2 + hi/2
		pivotValue := sequence[pivot]
		storeIdx := low
		sequence[pivot], sequence[hi] = sequence[hi], sequence[pivot]
		for i := low; i < hi; i++ {
			if sequence[i] < pivotValue {
				sequence[storeIdx], sequence[i] = sequence[i], sequence[storeIdx]
				storeIdx++
			}
		}
		sequence[hi], sequence[storeIdx] = sequence[storeIdx], sequence[hi]
		if k <= storeIdx {
			hi = storeIdx
		} else {
			low = storeIdx + 1
		}
	}

	if len(sequence)%2 == 0 {
		return sequence[k-1]/2 + sequence[k]/2
	}
	return sequence[k]
}
