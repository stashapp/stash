package models

import (
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/utils"
)

type DatePrecision int

const (
	// default precision is day
	DatePrecisionDay DatePrecision = iota
	DatePrecisionMonth
	DatePrecisionYear
)

// Date wraps a time.Time with a format of "YYYY-MM-DD"
type Date struct {
	time.Time
	Precision DatePrecision
}

var dateFormatPrecision = []string{
	"2006-01-02",
	"2006-01",
	"2006",
}

func (d Date) String() string {
	return d.Format(dateFormatPrecision[d.Precision])
}

func (d Date) After(o Date) bool {
	return d.Time.After(o.Time)
}

// ParseDate tries to parse the input string into a date using utils.ParseDateStringAsTime.
// If that fails, it attempts to parse the string with decreasing precision (month, then year).
// It returns a Date struct with the appropriate precision set, or an error if all parsing attempts fail.
func ParseDate(s string) (Date, error) {
	var errs []error

	// default parse to day precision
	ret, err := utils.ParseDateStringAsTime(s)
	if err == nil {
		return Date{Time: ret, Precision: DatePrecisionDay}, nil
	}

	errs = append(errs, err)

	// try month and year precision
	for i, format := range dateFormatPrecision[1:] {
		ret, err := time.Parse(format, s)
		if err == nil {
			return Date{Time: ret, Precision: DatePrecision(i + 1)}, nil
		}
		errs = append(errs, err)
	}

	return Date{}, fmt.Errorf("failed to parse date %q: %v", s, errs)
}
