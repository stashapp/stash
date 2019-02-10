package jsonschema

import (
	"fmt"
	"strings"
	"time"
)

type RailsTime struct {
	time.Time
}

const railsTimeLayout = "2006-01-02 15:04:05 MST"

func (ct *RailsTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(railsTimeLayout, s)
	if err != nil {
		ct.Time, err = time.Parse(time.RFC3339, s)
	}
	return
}

func (ct *RailsTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(time.RFC3339))), nil
}

func (ct *RailsTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}