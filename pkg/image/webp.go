package image

import (
	"bytes"
)

const (
	formatWebP = "webp"
	formatGif  = "gif"
)

// https://developers.google.com/speed/webp/docs/riff_container
func isWebPAnimated(buf []byte) bool {
	const (
		webPHeaderStart = 8
		webPHeaderEnd   = 12
		webPHeader      = "WEBP"

		animationHeaderLoc    = 16
		minAnimSignatureIndex = 20

		maxSize = 48
	)

	// truncate the buffer to the max size
	if len(buf) > maxSize {
		buf = buf[:maxSize]
	}

	isWebp := len(buf) >= webPHeaderEnd && string(buf[webPHeaderStart:webPHeaderEnd]) == "WEBP" // is WEBP

	if isWebp {
		const animBit byte = 1 << 1
		if len(buf) > minAnimSignatureIndex {
			// Animation Bit is set and ANIM header is present
			return (buf[animationHeaderLoc]&animBit == animBit) && containsAnimSignature(buf[minAnimSignatureIndex:])
		}
	}
	return false
}

// https://developers.google.com/speed/webp/docs/riff_container#animation
func containsAnimSignature(buf []byte) bool {
	index := bytes.Index(buf, []byte("ANIM"))
	return index != -1
}
