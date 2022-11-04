// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transforms

import (
	"image"
)

// Rgb2Gray function converts RGB to a gray scale array.
func Rgb2Gray(colorImg image.Image) [][]float64 {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	pixels := make([][]float64, h)

	for i := range pixels {
		pixels[i] = make([]float64, w)
		for j := range pixels[i] {
			color := colorImg.At(j, i)
			r, g, b, _ := color.RGBA()
			lum := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
			pixels[i][j] = lum
		}
	}

	return pixels
}

// Rgb2GrayFast function converts RGB to a gray scale array.
func Rgb2GrayFast(colorImg image.Image, pixels *[]float64) {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	if w != h {
		return
	}
	switch c := colorImg.(type) {
	case *image.YCbCr:
		rgb2GrayYCbCR(c, *pixels, w)
	case *image.RGBA:
		rgb2GrayRGBA(c, *pixels, w)
	default:
		rgb2GrayDefault(c, *pixels, w)
	}
}

// pixel2Gray converts a pixel to grayscale value base on luminosity
func pixel2Gray(r, g, b, a uint32) float64 {
	return 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
}

// rgb2GrayDefault uses the image.Image interface
func rgb2GrayDefault(colorImg image.Image, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[j+(i*s)] = pixel2Gray(colorImg.At(j, i).RGBA())
		}
	}
}

// rgb2GrayYCbCR uses *image.YCbCr which is signifiantly faster than the image.Image interface.
func rgb2GrayYCbCR(colorImg *image.YCbCr, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[j+(i*s)] = pixel2Gray(colorImg.YCbCrAt(j, i).RGBA())
		}
	}
}

// rgb2GrayYCbCR uses *image.RGBA which is signifiantly faster than the image.Image interface.
func rgb2GrayRGBA(colorImg *image.RGBA, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = pixel2Gray(colorImg.At(j, i).RGBA())
		}
	}
}

// FlattenPixels function flattens 2d array into 1d array.
func FlattenPixels(pixels [][]float64, x int, y int) []float64 {
	flattens := make([]float64, x*y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattens[y*i+j] = pixels[i][j]
		}
	}
	return flattens
}

// FlattenPixelsFast64 function flattens 2d array into 1d array.
func FlattenPixelsFast64(pixels []float64, x int, y int) []float64 {
	flattens := [64]float64{}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattens[y*i+j] = pixels[(i*64)+j]
		}
	}
	return flattens[:]
}
