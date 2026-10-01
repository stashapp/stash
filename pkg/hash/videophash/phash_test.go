package videophash

import (
	"image"
	"image/color"
	"testing"

	"github.com/disintegration/imaging"
)

func TestSpriteColumns(t *testing.T) {
	tests := []struct {
		name     string
		duration float64
		want     int
	}{
		{"unknown duration", 0, 5},
		{"one second clip", 1, 2},
		{"45 seconds", 45, 2},
		{"just over 45 seconds", 45.1, 3},
		{"90 seconds", 90, 3},
		{"just over 90 seconds", 90.1, 4},
		{"150 seconds", 150, 4},
		{"just over 150 seconds", 150.1, 5},
		{"feature length", 3600, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := spriteColumns(tt.duration); got != tt.want {
				t.Errorf("spriteColumns(%v) = %d, want %d", tt.duration, got, tt.want)
			}
		})
	}
}

func TestCombineImages(t *testing.T) {
	const frameSize = 16

	for _, columns := range []int{2, 3, 4, 5} {
		chunkCount := columns * columns

		// each frame is a solid colour, keyed off its index, so that the
		// position of every frame in the montage can be verified
		images := make([]image.Image, chunkCount)
		for i := range images {
			images[i] = imaging.New(frameSize, frameSize, color.NRGBA{R: uint8(i + 1), A: 255})
		}

		montage := combineImages(images, columns)

		want := frameSize * columns
		if size := montage.Bounds().Size(); size.X != want || size.Y != want {
			t.Errorf("columns %d: montage size = %dx%d, want %dx%d", columns, size.X, size.Y, want, want)
		}

		for i := range images {
			x := frameSize * (i % columns)
			y := frameSize * (i / columns)

			r, _, _, _ := montage.At(x, y).RGBA()
			if got := uint8(r >> 8); got != uint8(i+1) {
				t.Errorf("columns %d: frame %d at (%d,%d) = %d, want %d", columns, i, x, y, got, uint8(i+1))
			}
		}
	}
}
