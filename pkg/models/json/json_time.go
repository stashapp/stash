// Package json provides generic JSON types.
package json

import (
	"fmt"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/utils"
)

var currentLocation = time.Now().Location()

type JSONTime struct {
	time.Time
}

func (jt *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		jt.Time = time.Time{}
		return nil
	}

	// #731 - returning an error here causes the entire JSON parse to fail for ffprobe.
	jt.Time, _ = utils.ParseDateStringAsTime(s)
	return nil
}

func (jt *JSONTime) MarshalJSON() ([]byte, error) {
	if jt.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", jt.Time.Format(time.RFC3339))), nil
}

func (jt JSONTime) GetTime() time.Time {
	if currentLocation != nil {
		if jt.IsZero() {
			return time.Now().In(currentLocation)
		} else {
			return jt.Time.In(currentLocation)
		}
	} else {
		if jt.IsZero() {
			return time.Now()
		} else {
			return jt.Time
		}
	}
}
