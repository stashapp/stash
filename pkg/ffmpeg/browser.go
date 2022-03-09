package ffmpeg

// only support H264 by default, since Safari does not support VP8/VP9
var defaultSupportedCodecs = []string{H264, H265}

var validForH264Mkv = []Container{Mp4, Matroska}
var validForH264 = []Container{Mp4}
var validForH265Mkv = []Container{Mp4, Matroska}
var validForH265 = []Container{Mp4}
var validForVp8 = []Container{Webm}
var validForVp9Mkv = []Container{Webm, Matroska}
var validForVp9 = []Container{Webm}
var validForHevcMkv = []Container{Mp4, Matroska}
var validForHevc = []Container{Mp4}

var validAudioForMkv = []AudioCodec{Aac, Mp3, Vorbis, Opus}
var validAudioForWebm = []AudioCodec{Vorbis, Opus}
var validAudioForMp4 = []AudioCodec{Aac, Mp3}

func IsStreamable(videoCodec string, audioCodec AudioCodec, container Container) bool {
	supportedVideoCodecs := defaultSupportedCodecs

	// check if the video codec matches the supported codecs
	return isValidCodec(videoCodec, supportedVideoCodecs) && isValidCombo(videoCodec, container, supportedVideoCodecs) && IsValidAudioForContainer(audioCodec, container)
}

func isValidCodec(codecName string, supportedCodecs []string) bool {
	for _, c := range supportedCodecs {
		if c == codecName {
			return true
		}
	}
	return false
}

func isValidAudio(audio AudioCodec, validCodecs []AudioCodec) bool {
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

func IsValidAudioForContainer(audio AudioCodec, format Container) bool {
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
func isValidCombo(codecName string, format Container, supportedVideoCodecs []string) bool {
	supportMKV := isValidCodec(Mkv, supportedVideoCodecs)
	supportHEVC := isValidCodec(Hevc, supportedVideoCodecs)

	switch codecName {
	case H264:
		if supportMKV {
			return isValidForContainer(format, validForH264Mkv)
		}
		return isValidForContainer(format, validForH264)
	case H265:
		if supportMKV {
			return isValidForContainer(format, validForH265Mkv)
		}
		return isValidForContainer(format, validForH265)
	case Vp8:
		return isValidForContainer(format, validForVp8)
	case Vp9:
		if supportMKV {
			return isValidForContainer(format, validForVp9Mkv)
		}
		return isValidForContainer(format, validForVp9)
	case Hevc:
		if supportHEVC {
			if supportMKV {
				return isValidForContainer(format, validForHevcMkv)
			}
			return isValidForContainer(format, validForHevc)
		}
	}
	return false
}

func isValidForContainer(format Container, validContainers []Container) bool {
	for _, fmt := range validContainers {
		if fmt == format {
			return true
		}
	}
	return false
}
