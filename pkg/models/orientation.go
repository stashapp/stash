package models

type OrientationEnum string

const (
	OrientationLandscape OrientationEnum = "LANDSCAPE"
	OrientationPortrait  OrientationEnum = "PORTRAIT"
	OrientationSquare    OrientationEnum = "SQUARE"
)

func (e OrientationEnum) IsValid() bool {
	switch e {
	case OrientationLandscape, OrientationPortrait, OrientationSquare:
		return true
	}
	return false
}
