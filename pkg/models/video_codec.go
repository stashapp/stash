package models

import (
	"fmt"
	"io"
	"strconv"
)

var videoCodecOptions = map[VideoCodecEnum]string{
	VideoCodecEnum("AV1"):         "av1",
	VideoCodecEnum("H264"):        "h264",
	VideoCodecEnum("HEVC"):        "hevc",
	VideoCodecEnum("MPEG2_VIDEO"): "mpeg2video",
	VideoCodecEnum("MPEG4"):       "mpeg4",
	VideoCodecEnum("VC1"):         "vc1",
	VideoCodecEnum("VP6F"):        "vp6f",
	VideoCodecEnum("WMV1"):        "wmv1",
	VideoCodecEnum("WMV2"):        "wmv2",
	VideoCodecEnum("WMV3"):        "wmv3",
}

type VideoCodecEnum string

const (
	// av1
	VideoCodecEnumAv1 VideoCodecEnum = "AV1"
	// h264
	VideoCodecEnumH264 VideoCodecEnum = "H264"
	// hevc
	VideoCodecEnumHevc VideoCodecEnum = "HEVC"
	// mpeg2video
	VideoCodecEnumMpeg2Video VideoCodecEnum = "MPEG2_VIDEO"
	// mpeg4
	VideoCodecEnumMpeg4 VideoCodecEnum = "MPEG4"
	// vc1
	VideoCodecEnumVc1 VideoCodecEnum = "VC1"
	// vp6f
	VideoCodecEnumVp6f VideoCodecEnum = "VP6F"
	// wmv1
	VideoCodecEnumWmv1 VideoCodecEnum = "WMV1"
	// wmv2
	VideoCodecEnumWmv2 VideoCodecEnum = "WMV2"
	// wmv3
	VideoCodecEnumWmv3 VideoCodecEnum = "WMV3"
)

var AllVideoCodecEnum = []VideoCodecEnum{
	VideoCodecEnumAv1,
	VideoCodecEnumH264,
	VideoCodecEnumHevc,
	VideoCodecEnumMpeg2Video,
	VideoCodecEnumMpeg4,
	VideoCodecEnumVc1,
	VideoCodecEnumVp6f,
	VideoCodecEnumWmv1,
	VideoCodecEnumWmv2,
	VideoCodecEnumWmv3,
}

func (e VideoCodecEnum) IsValid() bool {
	switch e {
	case VideoCodecEnumAv1,
		VideoCodecEnumH264,
		VideoCodecEnumHevc,
		VideoCodecEnumMpeg2Video,
		VideoCodecEnumMpeg4,
		VideoCodecEnumVc1,
		VideoCodecEnumVp6f,
		VideoCodecEnumWmv1,
		VideoCodecEnumWmv2,
		VideoCodecEnumWmv3:
		return true
	}
	return false
}

func (e *VideoCodecEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VideoCodecEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VideoCodecEnum", str)
	}
	return nil
}

func (e VideoCodecEnum) String() string {
	return string(e)
}

func (e VideoCodecEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *VideoCodecEnum) GetCodecValue() string {
	return videoCodecOptions[*e]
}
