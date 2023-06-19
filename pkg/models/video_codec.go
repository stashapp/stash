package models

import (
	"fmt"
	"io"
	"strconv"
)

var videoCodecOptions = map[VideoCodecEnum]string{
	VideoCodecEnum("H264"): "h264",
	VideoCodecEnum("HEVC"): "hevc",
	VideoCodecEnum("AV1"):  "av1",
	VideoCodecEnum("WMV3"): "wmv3",
	VideoCodecEnum("VC1"):  "vc1",
}

type VideoCodecEnum string

const (
	// h264
	VideoCodecEnumH264 VideoCodecEnum = "H264"
	// hevc
	VideoCodecEnumHevc VideoCodecEnum = "HEVC"
	// av1
	VideoCodecEnumAv1 VideoCodecEnum = "AV1"
	// wmv3
	VideoCodecEnumWmv3 VideoCodecEnum = "WMV3"
	// vc1
	VideoCodecEnumVc1 VideoCodecEnum = "VC1"
)

var AllVideoCodecEnum = []VideoCodecEnum{
	VideoCodecEnumH264,
	VideoCodecEnumHevc,
	VideoCodecEnumAv1,
	VideoCodecEnumWmv3,
	VideoCodecEnumVc1,
}

func (e VideoCodecEnum) IsValid() bool {
	switch e {
	case VideoCodecEnumH264, VideoCodecEnumHevc, VideoCodecEnumAv1, VideoCodecEnumWmv3, VideoCodecEnumVc1:
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
