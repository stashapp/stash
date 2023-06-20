package models

import (
	"fmt"
	"io"
	"strconv"
)

var audioCodecOptions = map[AudioCodecEnum]string{
	AudioCodecEnum("AAC"):           "aac",
	AudioCodecEnum("AC3"):           "ac3",
	AudioCodecEnum("ADPCM_IMA_WAV"): "adpcm_ima_wav",
	AudioCodecEnum("MP3"):           "mp3",
	AudioCodecEnum("PCM_S16LE"):     "pcm_s16le",
	AudioCodecEnum("WMAPRO"):        "wmapro",
	AudioCodecEnum("WMAV2"):         "wmav2",
}

type AudioCodecEnum string

const (
	// aac
	AudioCodecEnumAac AudioCodecEnum = "AAC"
	// ac3
	AudioCodecEnumAc3 AudioCodecEnum = "AC3"
	// adpcm_ima_wav
	AudioCodecEnumAdpcmImaWav AudioCodecEnum = "ADPCM_IMA_WAV"
	// mp3
	AudioCodecEnumMp3 AudioCodecEnum = "MP3"
	// pcm_s16le
	AudioCodecEnumPcmS16le AudioCodecEnum = "PCM_S16LE"
	// wmapro
	AudioCodecEnumWmaPro AudioCodecEnum = "WMAPRO"
	// wmav2
	AudioCodecEnumWmav2 AudioCodecEnum = "WMAV2"
)

var AllAudioCodecEnum = []AudioCodecEnum{
	AudioCodecEnumAac,
	AudioCodecEnumAc3,
	AudioCodecEnumAdpcmImaWav,
	AudioCodecEnumMp3,
	AudioCodecEnumPcmS16le,
	AudioCodecEnumWmaPro,
	AudioCodecEnumWmav2,
}

func (e AudioCodecEnum) IsValid() bool {
	switch e {
	case AudioCodecEnumAac,
		AudioCodecEnumAc3,
		AudioCodecEnumAdpcmImaWav,
		AudioCodecEnumMp3,
		AudioCodecEnumPcmS16le,
		AudioCodecEnumWmaPro,
		AudioCodecEnumWmav2:
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
