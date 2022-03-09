package ffmpeg

import "fmt"

type VideoFilter string

func (f VideoFilter) Args() []string {
	if f == "" {
		return nil
	}

	return []string{"-vf", string(f)}
}

func (f VideoFilter) ScaleWidth(w int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=%v:-2", w))
}

func (f VideoFilter) ScaleHeight(w int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=-2:%v", w))
}

// ScaleDimesions scales using w and h. Use -n to maintain aspect ratio and maintain as multiple of n.
func (f VideoFilter) ScaleDimensions(w, h int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=%v:%v", w, h))
}

func (f VideoFilter) ScaleMaxSize(maxDimensions int) VideoFilter {
	return f.Append(fmt.Sprintf("scale=%v:%v:force_original_aspect_ratio=decrease", maxDimensions, maxDimensions))
}

// ScaleMax scales to maxSize. It will scale width if it is larger than height, otherwise it will scale height.
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

func (f VideoFilter) Fps(fps int) VideoFilter {
	return f.Append(fmt.Sprintf("fps=%v", fps))
}

func (f VideoFilter) Select(frame int) VideoFilter {
	return f.Append(fmt.Sprintf("select=eq(n\\,%d)", frame))
}

func (f VideoFilter) Append(s string) VideoFilter {
	// if filter is empty, then just set
	if f == "" {
		return VideoFilter(s)
	}

	return VideoFilter(fmt.Sprintf(",%s", s))
}
