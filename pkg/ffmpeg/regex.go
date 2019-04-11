package ffmpeg

import (
	"regexp"
	"strconv"
)

var TimeRegex = regexp.MustCompile(`time=\s*(\d+):(\d+):(\d+.\d+)`)
var FrameRegex = regexp.MustCompile(`frame=\s*([0-9]+)`)

func GetTimeFromRegex(str string) float64 {
	regexResult := TimeRegex.FindStringSubmatch(str)

	// Bail early if we don't have the results we expect
	if len(regexResult) != 4 {
		return 0
	}

	h, _ := strconv.ParseFloat(regexResult[1], 64)
	m, _ := strconv.ParseFloat(regexResult[2], 64)
	s, _ := strconv.ParseFloat(regexResult[3], 64)
	hours := h * 3600
	minutes := m * 60
	seconds := s
	return hours + minutes + seconds
}

func GetFrameFromRegex(str string) int {
	regexResult := FrameRegex.FindStringSubmatch(str)

	// Bail early if we don't have the results we expect
	if len(regexResult) < 2 {
		return 0
	}

	result, _ := strconv.Atoi(regexResult[1])
	return result
}
