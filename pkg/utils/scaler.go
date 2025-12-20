package utils

// ScaleMode defines how to scale dimensions
type ScaleMode int

const (
	// ScaleToWidth scales to the specified width, maintaining aspect ratio
	ScaleToWidth ScaleMode = iota
	// ScaleToHeight scales to the specified height, maintaining aspect ratio
	ScaleToHeight
	// ScaleToMinSize scales so that the smaller dimension becomes the target size
	ScaleToMinSize
)

// ScaleDimensions scales the given width and height according to the specified mode and target.
// Returns the new width and height. If width or height is 0, returns 0,0.
// For ScaleToMinSize, if target is larger than or equal to the relevant dimension, returns the original dimensions.
func ScaleDimensions(width, height int, target int, mode ScaleMode) (int, int) {
	if width <= 0 || height <= 0 {
		return 0, 0
	}

	switch mode {
	case ScaleToWidth:
		newHeight := int(float64(target) * float64(height) / float64(width))
		return target, newHeight
	case ScaleToHeight:
		newWidth := int(float64(target) * float64(width) / float64(height))
		return newWidth, target
	case ScaleToMinSize:
		minDim := height
		if width < height {
			minDim = width
		}
		if target >= minDim {
			return width, height
		}
		if width >= height {
			return ScaleDimensions(width, height, target, ScaleToHeight)
		} else {
			return ScaleDimensions(width, height, target, ScaleToWidth)
		}
	}

	return width, height
}

// GetFFmpegScaleArgs returns the width and height arguments for ffmpeg ScaleDimensions filter.
// Uses -2 to maintain aspect ratio where appropriate.
// Returns 0,0 if no scaling is needed.
func GetFFmpegScaleArgs(width, height int, target int, mode ScaleMode) (int, int) {
	if width <= 0 || height <= 0 {
		return 0, 0
	}

	switch mode {
	case ScaleToWidth:
		return target, -2
	case ScaleToHeight:
		return -2, target
	case ScaleToMinSize:
		minDim := height
		if width < height {
			minDim = width
		}
		if target >= minDim || target == 0 {
			return 0, 0
		}
		if width >= height {
			return -2, target
		} else {
			return target, -2
		}
	}

	return 0, 0
}

// GetFFmpegScaleArgsForRect returns ffmpeg scale args to fit within maxWidth x maxHeight,
// starting from reqHeight as the desired height.
func GetFFmpegScaleArgsForRect(width, height, reqHeight, maxWidth, maxHeight int) (int, int) {
	if width <= 0 || height <= 0 || maxWidth <= 0 || maxHeight <= 0 {
		return 0, 0
	}

	aspectRatio := float64(width) / float64(height)
	desiredHeight := reqHeight
	if desiredHeight == 0 {
		desiredHeight = height
	}
	desiredWidth := int(float64(desiredHeight) * aspectRatio)

	if desiredHeight <= maxHeight && desiredWidth <= maxWidth {
		// Fits, so use ScaleToMinSize with reqHeight
		return GetFFmpegScaleArgs(width, height, reqHeight, ScaleToMinSize)
	}

	// Doesn't fit, scale to the limiting dimension
	if float64(desiredHeight-maxHeight) > float64(desiredWidth-maxWidth) {
		return -2, maxHeight
	} else {
		return maxWidth, -2
	}
}
