package models

import (
	"fmt"
	"github.com/stashapp/stash/pkg/utils"
	"strings"
	"time"
)

type JSONTime struct {
	time.Time
}

func (jt *JSONTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		jt.Time = time.Time{}
		return
	}

	jt.Time, err = utils.ParseDateStringAsTime(s)
	return
}

func (jt *JSONTime) MarshalJSON() ([]byte, error) {
	if jt.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", jt.Time.Format(time.RFC3339))), nil
}
