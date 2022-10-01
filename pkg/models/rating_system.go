package models

import (
	"fmt"
	"io"
	"strconv"
)

type RatingSystem string

const (
	FiveStar = "FiveStar"
	FivePointFiveStar = "FivePointFiveStar"
	FivePointTwoFiveStar = "FivePointTwoFiveStar"
	TenStar = "TenStar"
	TenPointFiveStar = "TenPointFiveStar"
	TenPointTwoFiveStar = "TenPointTwoFiveStar"
	TenPointDecimal = "TenPointDecimal"
	None = "None"
)

func (e RatingSystem) IsValid() bool {
	switch e {
	case FiveStar, FivePointFiveStar, FivePointTwoFiveStar, TenStar, TenPointFiveStar, TenPointTwoFiveStar, TenPointDecimal, None:
		return true
	}
	return false
}


func (e RatingSystem) String() string {
	return string(e)
}

func (e *RatingSystem) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RatingSystem(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RatingSystem", str)
	}
	return nil
}

func (e RatingSystem) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
