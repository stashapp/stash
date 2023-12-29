package models

type OrientationEnum string

const (
	OrientationLandscape OrientationEnum = "Landscape"
	OrientationPortrait  OrientationEnum = "Portrait"
	OrientationSquare    OrientationEnum = "Square"
)

func (e OrientationEnum) IsValid() bool {
	switch e {
	case OrientationLandscape, OrientationPortrait, OrientationSquare:
		return true
	}
	return false
}
