package jsonschema

import (
	"fmt"
	"github.com/stashapp/stash/pkg/utils"
	"strings"
	"time"
)

type RailsTime struct {
	time.Time
}

func (ct *RailsTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	t, err := utils.ParseDateStringAsTime(s)
	if t != nil {
		ct.Time = *t
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
