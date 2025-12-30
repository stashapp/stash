package config

import (
	"fmt"
	"io"
	"strconv"
)

type ConfigImageLightboxResult struct {
	SlideshowDelay             *int                      `json:"slideshowDelay"`
	DisplayMode                *ImageLightboxDisplayMode `json:"displayMode"`
	ScaleUp                    *bool                     `json:"scaleUp"`
	ResetZoomOnNav             *bool                     `json:"resetZoomOnNav"`
	ScrollMode                 *ImageLightboxScrollMode  `json:"scrollMode"`
	ScrollAttemptsBeforeChange int                       `json:"scrollAttemptsBeforeChange"`
	DisableAnimation           *bool                     `json:"disableAnimation"`
}

type ImageLightboxDisplayMode string

const (
	ImageLightboxDisplayModeOriginal ImageLightboxDisplayMode = "ORIGINAL"
	ImageLightboxDisplayModeFitXy    ImageLightboxDisplayMode = "FIT_XY"
	ImageLightboxDisplayModeFitX     ImageLightboxDisplayMode = "FIT_X"
)

var AllImageLightboxDisplayMode = []ImageLightboxDisplayMode{
	ImageLightboxDisplayModeOriginal,
	ImageLightboxDisplayModeFitXy,
	ImageLightboxDisplayModeFitX,
}

func (e ImageLightboxDisplayMode) IsValid() bool {
	switch e {
	case ImageLightboxDisplayModeOriginal, ImageLightboxDisplayModeFitXy, ImageLightboxDisplayModeFitX:
		return true
	}
	return false
}

func (e ImageLightboxDisplayMode) String() string {
	return string(e)
}

func (e *ImageLightboxDisplayMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ImageLightboxDisplayMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ImageLightboxDisplayMode", str)
	}
	return nil
}

func (e ImageLightboxDisplayMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ImageLightboxScrollMode string

const (
	ImageLightboxScrollModeZoom ImageLightboxScrollMode = "ZOOM"
	ImageLightboxScrollModePanY ImageLightboxScrollMode = "PAN_Y"
)

var AllImageLightboxScrollMode = []ImageLightboxScrollMode{
	ImageLightboxScrollModeZoom,
	ImageLightboxScrollModePanY,
}

func (e ImageLightboxScrollMode) IsValid() bool {
	switch e {
	case ImageLightboxScrollModeZoom, ImageLightboxScrollModePanY:
		return true
	}
	return false
}

func (e ImageLightboxScrollMode) String() string {
	return string(e)
}

func (e *ImageLightboxScrollMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ImageLightboxScrollMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ImageLightboxScrollMode", str)
	}
	return nil
}

func (e ImageLightboxScrollMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ConfigDisableDropdownCreate struct {
	Performer bool `json:"performer"`
	Tag       bool `json:"tag"`
	Studio    bool `json:"studio"`
	Movie     bool `json:"movie"`
	Gallery   bool `json:"gallery"`
}
