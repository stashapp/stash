package models

import (
	"fmt"
	"io"
	"strconv"
)

type ResolutionRange struct {
	min, max int
}

var resolutionRanges = map[ResolutionEnum]ResolutionRange{
	ResolutionEnum("VERY_LOW"):    {144, 239},
	ResolutionEnum("LOW"):         {240, 359},
	ResolutionEnum("R360P"):       {360, 479},
	ResolutionEnum("STANDARD"):    {480, 539},
	ResolutionEnum("WEB_HD"):      {540, 719},
	ResolutionEnum("STANDARD_HD"): {720, 1079},
	ResolutionEnum("FULL_HD"):     {1080, 1439},
	ResolutionEnum("QUAD_HD"):     {1440, 1919},
	ResolutionEnum("VR_HD"):       {1920, 2159},
	ResolutionEnum("FOUR_K"):      {1920, 2559},
	ResolutionEnum("FIVE_K"):      {2560, 2999},
	ResolutionEnum("SIX_K"):       {3000, 3583},
	ResolutionEnum("SEVEN_K"):     {3584, 3839},
	ResolutionEnum("EIGHT_K"):     {3840, 6143},
	ResolutionEnum("HUGE"):        {6144, 9999},
}

type ResolutionEnum string

const (
	// 144p
	ResolutionEnumVeryLow ResolutionEnum = "VERY_LOW"
	// 240p
	ResolutionEnumLow ResolutionEnum = "LOW"
	// 360p
	ResolutionEnumR360p ResolutionEnum = "R360P"
	// 480p
	ResolutionEnumStandard ResolutionEnum = "STANDARD"
	// 540p
	ResolutionEnumWebHd ResolutionEnum = "WEB_HD"
	// 720p
	ResolutionEnumStandardHd ResolutionEnum = "STANDARD_HD"
	// 1080p
	ResolutionEnumFullHd ResolutionEnum = "FULL_HD"
	// 1440p
	ResolutionEnumQuadHd ResolutionEnum = "QUAD_HD"
	// 1920p - deprecated
	ResolutionEnumVrHd ResolutionEnum = "VR_HD"
	// 4k
	ResolutionEnumFourK ResolutionEnum = "FOUR_K"
	// 5k
	ResolutionEnumFiveK ResolutionEnum = "FIVE_K"
	// 6k
	ResolutionEnumSixK ResolutionEnum = "SIX_K"
	// 7k
	ResolutionEnumSevenK ResolutionEnum = "SEVEN_K"
	// 8k
	ResolutionEnumEightK ResolutionEnum = "EIGHT_K"
	// 8K+
	ResolutionEnumHuge ResolutionEnum = "HUGE"
)

var AllResolutionEnum = []ResolutionEnum{
	ResolutionEnumVeryLow,
	ResolutionEnumLow,
	ResolutionEnumR360p,
	ResolutionEnumStandard,
	ResolutionEnumWebHd,
	ResolutionEnumStandardHd,
	ResolutionEnumFullHd,
	ResolutionEnumQuadHd,
	ResolutionEnumVrHd,
	ResolutionEnumFourK,
	ResolutionEnumFiveK,
	ResolutionEnumSixK,
	ResolutionEnumSevenK,
	ResolutionEnumEightK,
	ResolutionEnumHuge,
}

func (e ResolutionEnum) IsValid() bool {
	switch e {
	case ResolutionEnumVeryLow, ResolutionEnumLow, ResolutionEnumR360p, ResolutionEnumStandard, ResolutionEnumWebHd, ResolutionEnumStandardHd, ResolutionEnumFullHd, ResolutionEnumQuadHd, ResolutionEnumVrHd, ResolutionEnumFourK, ResolutionEnumFiveK, ResolutionEnumSixK, ResolutionEnumSevenK, ResolutionEnumEightK, ResolutionEnumHuge:
		return true
	}
	return false
}

func (e ResolutionEnum) String() string {
	return string(e)
}

func (e *ResolutionEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ResolutionEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ResolutionEnum", str)
	}
	return nil
}

func (e ResolutionEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// GetMaxResolution returns the maximum width or height that media must be
// to qualify as this resolution.
func (e *ResolutionEnum) GetMaxResolution() int {
	return resolutionRanges[*e].max
}

// GetMinResolution returns the minimum width or height that media must be
// to qualify as this resolution.
func (e *ResolutionEnum) GetMinResolution() int {
	return resolutionRanges[*e].min
}

type StreamingResolutionEnum string

const (
	// 240p
	StreamingResolutionEnumLow StreamingResolutionEnum = "LOW"
	// 480p
	StreamingResolutionEnumStandard StreamingResolutionEnum = "STANDARD"
	// 720p
	StreamingResolutionEnumStandardHd StreamingResolutionEnum = "STANDARD_HD"
	// 1080p
	StreamingResolutionEnumFullHd StreamingResolutionEnum = "FULL_HD"
	// 4k
	StreamingResolutionEnumFourK StreamingResolutionEnum = "FOUR_K"
	// Original
	StreamingResolutionEnumOriginal StreamingResolutionEnum = "ORIGINAL"
)

var AllStreamingResolutionEnum = []StreamingResolutionEnum{
	StreamingResolutionEnumLow,
	StreamingResolutionEnumStandard,
	StreamingResolutionEnumStandardHd,
	StreamingResolutionEnumFullHd,
	StreamingResolutionEnumFourK,
	StreamingResolutionEnumOriginal,
}

func (e StreamingResolutionEnum) IsValid() bool {
	switch e {
	case StreamingResolutionEnumLow, StreamingResolutionEnumStandard, StreamingResolutionEnumStandardHd, StreamingResolutionEnumFullHd, StreamingResolutionEnumFourK, StreamingResolutionEnumOriginal:
		return true
	}
	return false
}

func (e StreamingResolutionEnum) String() string {
	return string(e)
}

func (e *StreamingResolutionEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StreamingResolutionEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StreamingResolutionEnum", str)
	}
	return nil
}

func (e StreamingResolutionEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

var streamingResolutionMax = map[StreamingResolutionEnum]int{
	StreamingResolutionEnumLow:        resolutionRanges[ResolutionEnumLow].min,
	StreamingResolutionEnumStandard:   resolutionRanges[ResolutionEnumStandard].min,
	StreamingResolutionEnumStandardHd: resolutionRanges[ResolutionEnumStandardHd].min,
	StreamingResolutionEnumFullHd:     resolutionRanges[ResolutionEnumFullHd].min,
	StreamingResolutionEnumFourK:      resolutionRanges[ResolutionEnumFourK].min,
	StreamingResolutionEnumOriginal:   0,
}

func (e StreamingResolutionEnum) GetMaxResolution() int {
	return streamingResolutionMax[e]
}
