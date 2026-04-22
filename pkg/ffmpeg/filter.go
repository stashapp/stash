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

// ScaleMax scales to reqHeight (0 = source height), optionally clamped to a
// maxWidth x maxHeight rect. Aspect ratio is preserved.
func (f VideoFilter) ScaleMax(width, height, reqHeight, maxWidth, maxHeight int) VideoFilter {
	// if a rect is given, clamp to whichever edge overshoots it by more
	if maxWidth > 0 && maxHeight > 0 {
		// projected dimensions at reqHeight (or source height if 0)
		target := reqHeight
		if target == 0 {
			target = height
		}
		projectedWidth := target * width / height

		if target > maxHeight || projectedWidth > maxWidth {
			// cap the edge that exceeds its limit by more
			if target-maxHeight > projectedWidth-maxWidth {
				return f.ScaleDimensions(-2, maxHeight)
			}
			return f.ScaleDimensions(maxWidth, -2)
		}
	}

	// no-op if reqHeight is larger than the smaller dimension
	if reqHeight == 0 || reqHeight >= min(width, height) {
		return f
	}

	// scale the smaller dimension to reqHeight
	if width > height {
		// set the height
		return f.ScaleDimensions(-2, reqHeight)
	}
	return f.ScaleDimensions(reqHeight, -2)
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
