package models

import (
	"fmt"
	"io"
	"strconv"
)

var audioCodecOptions = map[AudioCodecEnum]string{
	AudioCodecEnum("AAC"): "aac",
}

type AudioCodecEnum string

const (
	// aac
	AudioCodecEnumAac AudioCodecEnum = "AAC"
)

var AllAudioCodecEnum = []AudioCodecEnum{
	AudioCodecEnumAac,
}

func (e AudioCodecEnum) IsValid() bool {
	switch e {
	case AudioCodecEnumAac:
		return true
	}
	return false
}

func (e *AudioCodecEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AudioCodecEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AudioCodecEnum", str)
	}
	return nil
}

func (e AudioCodecEnum) String() string {
	return string(e)
}

func (e AudioCodecEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *AudioCodecEnum) GetCodecValue() string {
	return audioCodecOptions[*e]
}
