package models

import (
	"fmt"
	"io"
	"strconv"
)

type GenerateMetadataOptions struct {
	Sprites                   *bool                   `json:"sprites"`
	Previews                  *bool                   `json:"previews"`
	ImagePreviews             *bool                   `json:"imagePreviews"`
	PreviewOptions            *GeneratePreviewOptions `json:"previewOptions"`
	Markers                   *bool                   `json:"markers"`
	MarkerImagePreviews       *bool                   `json:"markerImagePreviews"`
	MarkerScreenshots         *bool                   `json:"markerScreenshots"`
	Transcodes                *bool                   `json:"transcodes"`
	Phashes                   *bool                   `json:"phashes"`
	InteractiveHeatmapsSpeeds *bool                   `json:"interactiveHeatmapsSpeeds"`
}

type GeneratePreviewOptions struct {
	// Number of segments in a preview file
	PreviewSegments *int `json:"previewSegments"`
	// Preview segment duration, in seconds
	PreviewSegmentDuration *float64 `json:"previewSegmentDuration"`
	// Duration of start of video to exclude when generating previews
	PreviewExcludeStart *string `json:"previewExcludeStart"`
	// Duration of end of video to exclude when generating previews
	PreviewExcludeEnd *string `json:"previewExcludeEnd"`
	// Preset when generating preview
	PreviewPreset *PreviewPreset `json:"previewPreset"`
}

type PreviewPreset string

const (
	// X264_ULTRAFAST
	PreviewPresetUltrafast PreviewPreset = "ultrafast"
	// X264_VERYFAST
	PreviewPresetVeryfast PreviewPreset = "veryfast"
	// X264_FAST
	PreviewPresetFast PreviewPreset = "fast"
	// X264_MEDIUM
	PreviewPresetMedium PreviewPreset = "medium"
	// X264_SLOW
	PreviewPresetSlow PreviewPreset = "slow"
	// X264_SLOWER
	PreviewPresetSlower PreviewPreset = "slower"
	// X264_VERYSLOW
	PreviewPresetVeryslow PreviewPreset = "veryslow"
)

var AllPreviewPreset = []PreviewPreset{
	PreviewPresetUltrafast,
	PreviewPresetVeryfast,
	PreviewPresetFast,
	PreviewPresetMedium,
	PreviewPresetSlow,
	PreviewPresetSlower,
	PreviewPresetVeryslow,
}

func (e PreviewPreset) IsValid() bool {
	switch e {
	case PreviewPresetUltrafast, PreviewPresetVeryfast, PreviewPresetFast, PreviewPresetMedium, PreviewPresetSlow, PreviewPresetSlower, PreviewPresetVeryslow:
		return true
	}
	return false
}

func (e PreviewPreset) String() string {
	return string(e)
}

func (e *PreviewPreset) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PreviewPreset(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PreviewPreset", str)
	}
	return nil
}

func (e PreviewPreset) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
