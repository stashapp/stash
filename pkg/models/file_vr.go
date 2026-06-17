package models

import (
	"fmt"
	"io"
	"strconv"
)

// ProjectionEnum describes how a video frame is mapped onto a sphere
// or other projection surface. Generic naming; agnostic of any
// particular player or device.
type ProjectionEnum string

const (
	ProjectionEnumFlat            ProjectionEnum = "FLAT"
	ProjectionEnumEquirectangular ProjectionEnum = "EQUIRECTANGULAR"
	ProjectionEnumFisheye         ProjectionEnum = "FISHEYE"
	ProjectionEnumMKX200          ProjectionEnum = "MKX200"
	ProjectionEnumRF52            ProjectionEnum = "RF52"
	ProjectionEnumDome            ProjectionEnum = "DOME"
	ProjectionEnumCubemap         ProjectionEnum = "CUBEMAP"
	ProjectionEnumRectilinear     ProjectionEnum = "RECTILINEAR"
)

var AllProjectionEnum = []ProjectionEnum{
	ProjectionEnumFlat,
	ProjectionEnumEquirectangular,
	ProjectionEnumFisheye,
	ProjectionEnumMKX200,
	ProjectionEnumRF52,
	ProjectionEnumDome,
	ProjectionEnumCubemap,
	ProjectionEnumRectilinear,
}

func (e ProjectionEnum) IsValid() bool {
	switch e {
	case ProjectionEnumFlat, ProjectionEnumEquirectangular, ProjectionEnumFisheye,
		ProjectionEnumMKX200, ProjectionEnumRF52, ProjectionEnumDome,
		ProjectionEnumCubemap, ProjectionEnumRectilinear:
		return true
	}
	return false
}

func (e ProjectionEnum) String() string {
	return string(e)
}

func (e *ProjectionEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}
	*e = ProjectionEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProjectionEnum", str)
	}
	return nil
}

func (e ProjectionEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// StereoModeEnum describes how a stereoscopic video is laid out
// across the frame. Generic naming.
type StereoModeEnum string

const (
	StereoModeEnumMono              StereoModeEnum = "MONO"
	StereoModeEnumSBS               StereoModeEnum = "SBS"
	StereoModeEnumTB                StereoModeEnum = "TB"
	StereoModeEnumCUV               StereoModeEnum = "CUV"
	StereoModeEnumAlternatingFrames StereoModeEnum = "AF"
	StereoModeEnumInterleavedRows   StereoModeEnum = "INTERLEAVED_ROWS"
)

var AllStereoModeEnum = []StereoModeEnum{
	StereoModeEnumMono,
	StereoModeEnumSBS,
	StereoModeEnumTB,
	StereoModeEnumCUV,
	StereoModeEnumAlternatingFrames,
	StereoModeEnumInterleavedRows,
}

func (e StereoModeEnum) IsValid() bool {
	switch e {
	case StereoModeEnumMono, StereoModeEnumSBS, StereoModeEnumTB, StereoModeEnumCUV,
		StereoModeEnumAlternatingFrames, StereoModeEnumInterleavedRows:
		return true
	}
	return false
}

func (e StereoModeEnum) String() string {
	return string(e)
}

func (e *StereoModeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}
	*e = StereoModeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StereoModeEnum", str)
	}
	return nil
}

func (e StereoModeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// VRCorrections holds optional per-video corrections commonly
// applied when projecting stereoscopic or 360-degree content.
// All fields are optional; nil means "unset".
// AlphaMode follows DeoVR's alpha channel specification.
type VRCorrections struct {
	HorizontalOffset *float64 `json:"horizontal_offset"`
	VerticalOffset   *float64 `json:"vertical_offset"`
	AlphaMode        *string  `json:"alpha_mode"`
}
