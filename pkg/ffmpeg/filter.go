package ffmpeg

import (
	"fmt"
)

// VideoFilter represents video filter parameters to be passed to ffmpeg.
type VideoFilter string

// Args converts the video filter parameters to a slice of arguments to be passed to ffmpeg.
// Returns an empty slice if the filter is empty.
func (f VideoFilter) Args() []string {
	if f == "" {
		return nil
	}

	return []string{"-vf", string(f)}
}

// ScaleWidth returns a VideoFilter scaling the width to the given width, maintaining aspect ratio and a height as a multiple of 2.
func (f VideoFilter) ScaleWidth(w int) VideoFilter {
	return f.ScaleDimensions(w, -2)
}

func (f VideoFilter) ScaleHeight(h int) VideoFilter {
	return f.ScaleDimensions(-2, h)
}

// ScaleDimesions returns a VideoFilter scaling using w and h. Use -n to maintain aspect ratio and maintain as multiple of n.
func (f VideoFilter) ScaleDimensions(w, h int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=%v:%v", w, h))
}

// ScaleMaxSize returns a VideoFilter scaling to maxDimensions, maintaining aspect ratio using force_original_aspect_ratio=decrease.
func (f VideoFilter) ScaleMaxSize(maxDimensions int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=%v:%v:force_original_aspect_ratio=decrease", maxDimensions, maxDimensions))
}

// ScaleMax returns a VideoFilter scaling to maxSize. It will scale width if it is larger than height, otherwise it will scale height.
func (f VideoFilter) ScaleMax(inputWidth, inputHeight, maxSize int) VideoFilter {
	// get the smaller dimension of the input
	videoSize := inputHeight
	if inputWidth < videoSize {
		videoSize = inputWidth
	}

	// if maxSize is larger than the video dimension, then no-op
	if maxSize >= videoSize || maxSize == 0 {
		return f
	}

	// we're setting either the width or height
	// we'll set the smaller dimesion
	if inputWidth > inputHeight {
		// set the height
		return f.ScaleDimensions(-2, maxSize)
	}

	return f.ScaleDimensions(maxSize, -2)
}

// ScaleMaxLM returns a VideoFilter scaling to maxSize with respect to a max size.
func (f VideoFilter) ScaleMaxLM(width int, height int, reqHeight int, maxWidth int, maxHeight int) VideoFilter {
	// calculate the aspect ratio of the current resolution
	aspectRatio := width / height

	// find the max height
	desiredHeight := reqHeight
	if desiredHeight == 0 {
		desiredHeight = height
	}

	// calculate the desired width based on the desired height and the aspect ratio
	desiredWidth := int(desiredHeight * aspectRatio)

	// check which dimension to scale based on the maximum resolution
	if desiredHeight > maxHeight || desiredWidth > maxWidth {
		if desiredHeight-maxHeight > desiredWidth-maxWidth {
			// scale the height down to the maximum height
			return f.ScaleDimensions(-2, maxHeight)
		} else {
			// scale the width down to the maximum width
			return f.ScaleDimensions(maxWidth, -2)
		}
	}

	// the current resolution can be scaled to the desired height without exceeding the maximum resolution
	return f.ScaleMax(width, height, reqHeight)
}

// Fps returns a VideoFilter setting the frames per second.
func (f VideoFilter) Fps(fps int) VideoFilter {
	return f.Append(fmt.Sprintf("fps=%v", fps))
}

// Select returns a VideoFilter to select the given frame.
func (f VideoFilter) Select(frame int) VideoFilter {
	return f.Append(fmt.Sprintf("select=eq(n\\,%d)", frame))
}

// Append returns a VideoFilter appending the given string.
func (f VideoFilter) Append(s string) VideoFilter {
	// if filter is empty, then just set
	if f == "" {
		return VideoFilter(s)
	}

	return VideoFilter(fmt.Sprintf("%s,%s", f, s))
}
