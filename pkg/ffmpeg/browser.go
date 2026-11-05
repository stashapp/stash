package ffmpeg

import (
	"errors"
	"fmt"
)

// DefaultSupportedCodecs is used when no codecs are configured.
// Only H264 is safe across all browsers.
var DefaultSupportedCodecs = []string{H264}

// mkv is never listed: not all browsers (notably Safari) can play it, so mkv files are always transcoded.
var validForH264 = []Container{Mp4}
var validForVp8 = []Container{Webm}
var validForVp9 = []Container{Webm}
var validForHevc = []Container{Mp4}
var validForOtherCodecs = []Container{Mp4, Webm}

var validAudioForMkv = []ProbeAudioCodec{Aac, Mp3, Vorbis, Opus}
var validAudioForWebm = []ProbeAudioCodec{Vorbis, Opus}
var validAudioForMp4 = []ProbeAudioCodec{Aac, Mp3, Opus}

var (
	// ErrUnsupportedVideoCodecForBrowser is returned when the video codec is not supported for browser streaming.
	ErrUnsupportedVideoCodecForBrowser = errors.New("unsupported video codec for browser")

	// ErrUnsupportedVideoCodecContainer is returned when the video codec/container combination is not supported for browser streaming.
	ErrUnsupportedVideoCodecContainer = errors.New("video codec/container combination is unsupported for browser streaming")

	// ErrUnsupportedAudioCodecContainer is returned when the audio codec/container combination is not supported for browser streaming.
	ErrUnsupportedAudioCodecContainer = errors.New("audio codec/container combination is unsupported for browser streaming")
)

// IsStreamable returns nil if the file is streamable, or an error if it is not.
// supportedVideoCodecs are ffprobe codec names (eg h264, hevc, av1).
func IsStreamable(videoCodec string, audioCodec ProbeAudioCodec, container Container, supportedVideoCodecs []string) error {
	// check if the video codec matches the supported codecs
	if !isValidCodec(videoCodec, supportedVideoCodecs) {
		return fmt.Errorf("%w: %s", ErrUnsupportedVideoCodecForBrowser, videoCodec)
	}

	if !isValidCombo(videoCodec, container) {
		return fmt.Errorf("%w: %s/%s", ErrUnsupportedVideoCodecContainer, videoCodec, container)
	}

	if !IsValidAudioForContainer(audioCodec, container) {
		return fmt.Errorf("%w: %s/%s", ErrUnsupportedAudioCodecContainer, audioCodec, container)
	}

	return nil
}

func isValidCodec(codecName string, supportedCodecs []string) bool {
	for _, c := range supportedCodecs {
		if c == codecName {
			return true
		}
	}
	return false
}

func isValidAudio(audio ProbeAudioCodec, validCodecs []ProbeAudioCodec) bool {
	// if audio codec is missing or unsupported by ffmpeg we can't do anything about it
	// report it as valid so that the file can at least be streamed directly if the video codec is supported
	if audio == MissingUnsupported {
		return true
	}

	for _, c := range validCodecs {
		if c == audio {
			return true
		}
	}

	return false
}

// IsValidAudioForContainer returns true if the audio codec is valid for the container.
func IsValidAudioForContainer(audio ProbeAudioCodec, format Container) bool {
	switch format {
	case Matroska:
		return isValidAudio(audio, validAudioForMkv)
	case Webm:
		return isValidAudio(audio, validAudioForWebm)
	case Mp4:
		return isValidAudio(audio, validAudioForMp4)
	}
	return false
}

// isValidCombo checks if a codec/container combination is valid.
// Returns true on validity, false otherwise
func isValidCombo(codecName string, format Container) bool {
	switch codecName {
	case H264:
		return isValidForContainer(format, validForH264)
	case H265, Hevc:
		return isValidForContainer(format, validForHevc)
	case Vp8:
		return isValidForContainer(format, validForVp8)
	case Vp9:
		return isValidForContainer(format, validForVp9)
	default:
		// user-whitelisted codecs without specific rules (eg av1, vvc)
		return isValidForContainer(format, validForOtherCodecs)
	}
}

func isValidForContainer(format Container, validContainers []Container) bool {
	for _, fmt := range validContainers {
		if fmt == format {
			return true
		}
	}
	return false
}
