package models

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

type RatingSystem string

const (
	FiveStar             = "FiveStar"
	FivePointFiveStar    = "FivePointFiveStar"
	FivePointTwoFiveStar = "FivePointTwoFiveStar"
	// TenStar              = "TenStar"
	// TenPointFiveStar     = "TenPointFiveStar"
	// TenPointTwoFiveStar  = "TenPointTwoFiveStar"
	TenPointDecimal = "TenPointDecimal"
)

func (e RatingSystem) IsValid() bool {
	switch e {
	// case FiveStar, FivePointFiveStar, FivePointTwoFiveStar, TenStar, TenPointFiveStar, TenPointTwoFiveStar, TenPointDecimal:
	case FiveStar, FivePointFiveStar, FivePointTwoFiveStar, TenPointDecimal:
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

const (
	maxRating100 = 100
	maxRating5   = 5
	minRating5   = 1
	minRating100 = 20
)

// Rating100To5 converts a 1-100 rating to a 1-5 rating.
// Values <= 30 are converted to 1. Otherwise, rating is divided by 20 and rounded to the nearest integer.
func Rating100To5(rating100 int) int {
	val := math.Round((float64(rating100) / 20))
	return int(math.Max(minRating5, math.Min(maxRating5, val)))
}

// Rating5To100 converts a 1-5 rating to a 1-100 rating
func Rating5To100(rating5 int) int {
	return int(math.Max(minRating100, math.Min(maxRating100, float64(rating5*20))))
}
