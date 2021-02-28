// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goimagehash

import (
	"image"
	"image/jpeg"
	"os"
	"testing"
)

func TestHashCompute(t *testing.T) {
	for _, tt := range []struct {
		img1     string
		img2     string
		method   func(img image.Image) (*ImageHash, error)
		name     string
		distance int
	}{
		{"_examples/sample1.jpg", "_examples/sample1.jpg", AverageHash, "AverageHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", AverageHash, "AverageHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", AverageHash, "AverageHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", AverageHash, "AverageHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", AverageHash, "AverageHash", 42},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", AverageHash, "AverageHash", 4},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", AverageHash, "AverageHash", 38},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", AverageHash, "AverageHash", 40},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", AverageHash, "AverageHash", 6},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", DifferenceHash, "DifferenceHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", DifferenceHash, "DifferenceHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", DifferenceHash, "DifferenceHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", DifferenceHash, "DifferenceHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", DifferenceHash, "DifferenceHash", 43},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", DifferenceHash, "DifferenceHash", 0},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", DifferenceHash, "DifferenceHash", 37},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", DifferenceHash, "DifferenceHash", 43},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", DifferenceHash, "DifferenceHash", 16},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", PerceptionHash, "PerceptionHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", PerceptionHash, "PerceptionHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", PerceptionHash, "PerceptionHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", PerceptionHash, "PerceptionHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", PerceptionHash, "PerceptionHash", 32},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", PerceptionHash, "PerceptionHash", 2},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", PerceptionHash, "PerceptionHash", 30},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", PerceptionHash, "PerceptionHash", 34},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", PerceptionHash, "PerceptionHash", 20},
	} {
		file1, err := os.Open(tt.img1)
		if err != nil {

		}
		defer file1.Close()

		file2, err := os.Open(tt.img2)
		if err != nil {
			t.Errorf("%s", err)
		}
		defer file2.Close()

		img1, err := jpeg.Decode(file1)
		if err != nil {
			t.Errorf("%s", err)
		}

		img2, err := jpeg.Decode(file2)
		if err != nil {
			t.Errorf("%s", err)
		}

		hash1, err := tt.method(img1)
		if err != nil {
			t.Errorf("%s", err)
		}
		hash2, err := tt.method(img2)
		if err != nil {
			t.Errorf("%s", err)
		}

		dis1, err := hash1.Distance(hash2)
		if err != nil {
			t.Errorf("%s", err)
		}

		dis2, err := hash2.Distance(hash1)
		if err != nil {
			t.Errorf("%s", err)
		}

		if dis1 != dis2 {
			t.Errorf("Distance should be identical %v vs %v", dis1, dis2)
		}

		if dis1 != tt.distance {
			t.Errorf("%s: Distance between %v and %v is expected %v but got %v", tt.name, tt.img1, tt.img2, tt.distance, dis1)
		}
	}
}

func TestNilHashCompute(t *testing.T) {
	hash, err := AverageHash(nil)
	if err == nil {
		t.Errorf("Error should be got.")
	}
	if hash != nil {
		t.Errorf("Nil hash should be got. but got %v", hash)
	}

	hash, err = DifferenceHash(nil)
	if err == nil {
		t.Errorf("Error should be got.")
	}
	if hash != nil {
		t.Errorf("Nil hash should be got. but got %v", hash)
	}

	hash, err = PerceptionHash(nil)
	if err == nil {
		t.Errorf("Error should be got.")
	}
	if hash != nil {
		t.Errorf("Nil hash should be got. but got %v", hash)
	}
}

func TestNilExtendHashCompute(t *testing.T) {
	hash, err := ExtAverageHash(nil, 8, 8)
	if err == nil {
		t.Errorf("Error should be got.")
	}
	if hash != nil {
		t.Errorf("Nil hash should be got. but got %v", hash)
	}

	hash, err = ExtDifferenceHash(nil, 8, 8)
	if err == nil {
		t.Errorf("Error should be got.")
	}
	if hash != nil {
		t.Errorf("Nil hash should be got. but got %v", hash)
	}

	hash, err = ExtPerceptionHash(nil, 8, 8)
	if err == nil {
		t.Errorf("Error should be got.")
	}
	if hash != nil {
		t.Errorf("Nil hash should be got. but got %v", hash)
	}

	hash, err = ExtPerceptionHash(nil, 8, 9)
	if err == nil {
		t.Errorf("Error should be got.")
	}
	if hash != nil {
		t.Errorf("Nil hash should be got. but got %v", hash)
	}
}

func BenchmarkDistanceIdentical(b *testing.B) {
	h1 := &ImageHash{hash: 0xe48ae53c05e502f7}
	h2 := &ImageHash{hash: 0xe48ae53c05e502f7}

	for i := 0; i < b.N; i++ {
		h1.Distance(h2)
	}
}

func BenchmarkDistanceDifferent(b *testing.B) {
	h1 := &ImageHash{hash: 0xe48ae53c05e502f7}
	h2 := &ImageHash{hash: 0x678be53815e510f7} // 8 bits flipped

	for i := 0; i < b.N; i++ {
		h1.Distance(h2)
	}
}

func TestExtImageHashCompute(t *testing.T) {
	for _, tt := range []struct {
		img1     string
		img2     string
		width    int
		height   int
		method   func(img image.Image, width, height int) (*ExtImageHash, error)
		name     string
		distance int
	}{
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 42},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 4},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 38},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 40},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", 8, 8, ExtAverageHash, "ExtAverageHash", 6},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 149},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 8},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 152},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 155},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", 16, 16, ExtAverageHash, "ExtAverageHash", 27},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 17, 17, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 17, 17, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 17, 17, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 17, 17, ExtAverageHash, "ExtAverageHash", 0},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 32},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 2},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 30},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 34},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", 8, 8, ExtPerceptionHash, "ExtPerceptionHash", 20},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 122},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 12},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 122},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 118},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", 16, 16, ExtPerceptionHash, "ExtPerceptionHash", 104},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 43},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 37},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 43},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", 8, 8, ExtDifferenceHash, "ExtDifferenceHash", 16},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample1.jpg", "_examples/sample2.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 139},
		{"_examples/sample1.jpg", "_examples/sample3.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 14},
		{"_examples/sample1.jpg", "_examples/sample4.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 130},
		{"_examples/sample2.jpg", "_examples/sample3.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 147},
		{"_examples/sample2.jpg", "_examples/sample4.jpg", 16, 16, ExtDifferenceHash, "ExtDifferenceHash", 89},
		{"_examples/sample1.jpg", "_examples/sample1.jpg", 17, 17, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample2.jpg", "_examples/sample2.jpg", 17, 17, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample3.jpg", "_examples/sample3.jpg", 17, 17, ExtDifferenceHash, "ExtDifferenceHash", 0},
		{"_examples/sample4.jpg", "_examples/sample4.jpg", 17, 17, ExtDifferenceHash, "ExtDifferenceHash", 0},
	} {
		file1, err := os.Open(tt.img1)
		if err != nil {
			t.Errorf("%s", err)
		}
		defer file1.Close()

		file2, err := os.Open(tt.img2)
		if err != nil {
			t.Errorf("%s", err)
		}
		defer file2.Close()

		img1, err := jpeg.Decode(file1)
		if err != nil {
			t.Errorf("%s", err)
		}

		img2, err := jpeg.Decode(file2)
		if err != nil {
			t.Errorf("%s", err)
		}

		hash1, err := tt.method(img1, tt.width, tt.height)
		if err != nil {
			t.Errorf("%s", err)
		}
		hash2, err := tt.method(img2, tt.width, tt.height)
		if err != nil {
			t.Errorf("%s", err)
		}

		dis1, err := hash1.Distance(hash2)
		if err != nil {
			t.Errorf("%s", err)
		}

		dis2, err := hash2.Distance(hash1)
		if err != nil {
			t.Errorf("%s", err)
		}

		if dis1 != dis2 {
			t.Errorf("Distance should be identical %v vs %v", dis1, dis2)
		}

		if dis1 != tt.distance {
			t.Errorf("%s: Distance between %v and %v is expected %v but got %v", tt.name, tt.img1, tt.img2, tt.distance, dis1)
		}
	}
}

func BenchmarkExtImageHashDistanceDifferent(b *testing.B) {
	h1 := &ExtImageHash{hash: []uint64{0xe48ae53c05e502f7}}
	h2 := &ExtImageHash{hash: []uint64{0x678be53815e510f7}} // 8 bits flipped

	for i := 0; i < b.N; i++ {
		_, err := h1.Distance(h2)
		if err != nil {
			b.Errorf("%s", err)
		}
	}
}

func BenchmarkPerceptionHash(b *testing.B) {
	file1, err := os.Open("_examples/sample3.jpg")
	if err != nil {
		b.Errorf("%s", err)
	}
	defer file1.Close()
	img1, err := jpeg.Decode(file1)
	if err != nil {
		b.Errorf("%s", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := ExtPerceptionHash(img1, 8, 8)
		if err != nil {
			b.Errorf("%s", err)
		}
	}
}

func BenchmarkAverageHash(b *testing.B) {
	file1, err := os.Open("_examples/sample3.jpg")
	if err != nil {
		b.Errorf("%s", err)
	}
	defer file1.Close()
	img1, err := jpeg.Decode(file1)
	if err != nil {
		b.Errorf("%s", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := ExtAverageHash(img1, 8, 8)
		if err != nil {
			b.Errorf("%s", err)
		}
	}
}

func BenchmarkDiffrenceHash(b *testing.B) {
	file1, err := os.Open("_examples/sample3.jpg")
	if err != nil {
		b.Errorf("%s", err)
	}
	defer file1.Close()
	img1, err := jpeg.Decode(file1)
	if err != nil {
		b.Errorf("%s", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := ExtDifferenceHash(img1, 8, 8)
		if err != nil {
			b.Errorf("%s", err)
		}
	}
}
