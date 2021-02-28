// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goimagehash

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestNewImageHash(t *testing.T) {
	for _, tt := range []struct {
		datas    [][]uint8
		hash1    Kind
		hash2    Kind
		distance int
		err      error
	}{
		{[][]uint8{{1, 0, 1, 1}, {0, 0, 0, 0}}, Unknown, Unknown, 3, nil},
		{[][]uint8{{0, 0, 0, 0}, {0, 0, 0, 0}}, Unknown, Unknown, 0, nil},
		{[][]uint8{{0, 0, 0, 0}, {0, 0, 0, 1}}, Unknown, Unknown, 1, nil},
		{[][]uint8{{0, 0, 0, 0}, {0, 0, 0, 1}}, Unknown, AHash, -1, errors.New("Image hashes's kind should be identical")},
	} {
		data1 := tt.datas[0]
		data2 := tt.datas[1]
		hash1 := NewImageHash(0, tt.hash1)
		hash2 := NewImageHash(0, tt.hash2)

		for i := 0; i < len(data1); i++ {
			if data1[i] == 1 {
				hash1.leftShiftSet(i)
			}
		}

		for i := 0; i < len(data2); i++ {
			if data2[i] == 1 {
				hash2.leftShiftSet(i)
			}
		}

		dis, err := hash1.Distance(hash2)
		if dis != tt.distance {
			t.Errorf("Distance between %v and %v expected as %d but got %d", data1, data2, tt.distance, dis)
		}
		if err != nil && err.Error() != tt.err.Error() {
			t.Errorf("Expected err %s, actual %s", tt.err, err)
		}
	}
}

func TestSerialization(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Errorf("%v", err)
		}
	}

	methods := []func(img image.Image) (*ImageHash, error){
		AverageHash, PerceptionHash, DifferenceHash,
	}
	examples := []string{
		"_examples/sample1.jpg", "_examples/sample2.jpg", "_examples/sample3.jpg", "_examples/sample4.jpg",
	}

	for _, ex := range examples {
		file, err := os.Open(ex)
		checkErr(err)

		defer file.Close()

		img, _, err := image.Decode(file)
		checkErr(err)

		for _, method := range methods {
			methodStr := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()

			hash, err := method(img)
			checkErr(err)

			hex := hash.ToString()
			// len(kind) == 1, len(":") == 1, len(hash) == 16
			if len(hex) != 18 {
				t.Errorf("Got invalid hex string '%v'; %v of '%v'", hex, methodStr, ex)
			}

			reHash, err := ImageHashFromString(hex)
			checkErr(err)

			distance, err := hash.Distance(reHash)
			checkErr(err)

			if distance != 0 {
				t.Errorf("Original and unserialized objects should be identical, got distance=%v; %v of '%v'", distance, methodStr, ex)
			}
		}

		// test for ExtIExtImageHash
		sizeList := []int{8, 16}
		for _, size := range sizeList {
			hash, err := ExtPerceptionHash(img, size, size)
			checkErr(err)

			hex := hash.ToString()
			// len(kind) == 1, len(":") == 1
			if len(hex) != size*size/4+2 {
				t.Errorf("Got invalid hex string '%v'; %v of '%v'", hex, "ExtPerceptionHash", ex)
			}

			reHash, err := ExtImageHashFromString(hex)
			checkErr(err)

			distance, err := hash.Distance(reHash)
			checkErr(err)

			if distance != 0 {
				t.Errorf("Original and unserialized objects should be identical, got distance=%v; %v of '%v'", distance, "ExtPerceptionHash", ex)
			}
		}
	}
}

func TestDifferentBitSizeHash(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Errorf("%v", err)
		}
	}
	file, err := os.Open("_examples/sample1.jpg")
	checkErr(err)
	defer file.Close()

	img, _, err := image.Decode(file)
	checkErr(err)

	hash1, _ := ExtAverageHash(img, 32, 32)
	hash2, _ := ExtDifferenceHash(img, 32, 32)
	_, err = hash1.Distance(hash2)
	if err == nil {
		t.Errorf("Should got error with different kinds of hashes")
	}
	hash3, _ := ExtAverageHash(img, 31, 31)
	_, err = hash1.Distance(hash3)
	if err == nil {
		t.Errorf("Should got error with different bits of hashes")
	}
}
func TestDumpAndLoad(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Errorf("%v", err)
		}
	}

	methods := []func(img image.Image) (*ImageHash, error){
		AverageHash, PerceptionHash, DifferenceHash,
	}
	examples := []string{
		"_examples/sample1.jpg", "_examples/sample2.jpg", "_examples/sample3.jpg", "_examples/sample4.jpg",
	}

	for _, ex := range examples {
		file, err := os.Open(ex)
		checkErr(err)

		defer file.Close()

		img, _, err := image.Decode(file)
		checkErr(err)

		for _, method := range methods {
			hash, err := method(img)
			checkErr(err)
			var b bytes.Buffer
			foo := bufio.NewWriter(&b)
			err = hash.Dump(foo)
			checkErr(err)
			foo.Flush()
			bar := bufio.NewReader(&b)
			reHash, err := LoadImageHash(bar)
			checkErr(err)

			distance, err := hash.Distance(reHash)
			checkErr(err)

			if distance != 0 {
				t.Errorf("Original and unserialized objects should be identical, got distance=%v", distance)
			}

			if hash.Bits() != 64 || reHash.Bits() != 64 {
				t.Errorf("Hash bits should be 64 but got, %v, %v", hash.Bits(), reHash.Bits())
			}
		}

		// test for ExtIExtImageHash
		extMethods := []func(img image.Image, width, height int) (*ExtImageHash, error){
			ExtAverageHash, ExtPerceptionHash, ExtDifferenceHash,
		}

		sizeList := []int{8, 16}
		for _, size := range sizeList {
			for _, method := range extMethods {
				hash, err := method(img, size, size)
				checkErr(err)
				var b bytes.Buffer
				foo := bufio.NewWriter(&b)
				err = hash.Dump(foo)
				checkErr(err)
				foo.Flush()
				bar := bufio.NewReader(&b)
				reHash, err := LoadExtImageHash(bar)
				checkErr(err)

				distance, err := hash.Distance(reHash)
				checkErr(err)

				if distance != 0 {
					t.Errorf("Original and unserialized objects should be identical, got distance=%v", distance)
				}

				if hash.Bits() != size*size || reHash.Bits() != size*size {
					t.Errorf("Hash bits should be 64 but got, %v, %v", hash.Bits(), reHash.Bits())
				}
			}
		}
	}
}
