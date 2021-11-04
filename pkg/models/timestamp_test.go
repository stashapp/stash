package models

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/plugin/common/log"
)

func TestTimestampSymmetry(t *testing.T) {
	initialTime := time.Now()
	buf := bytes.NewBuffer([]byte{})
	MarshalTimestamp(initialTime).MarshalGQL(buf)

	str, err := strconv.Unquote(buf.String())
	if err != nil {
		t.Fatal("could not unquote string")
	}
	newTime, err := UnmarshalTimestamp(str)
	if err != nil {
		t.Fatalf("could not unmarshal time: %v", err)
	}

	if !initialTime.Equal(newTime) {
		t.Fatalf("have %v, want %v", newTime, initialTime)
	}
}

func TestTimestamp(t *testing.T) {
	n := time.Now()
	testCases := []struct {
		name string
		have string
		want string
	}{
		{"reflexivity", n.Format(time.RFC3339Nano), n.Format(time.RFC3339Nano)},
		{"date-only", "2021-11-04T00:00:00Z", "2021-11-04T00:00:00Z"},
		{"unix", "@1636035887", "2021-11-04T15:24:47+01:00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := UnmarshalTimestamp(tc.have)
			if err != nil {
				t.Fatalf("could not unmarshal time: %v", err)
			}

			buf := bytes.NewBuffer([]byte{})
			MarshalTimestamp(p).MarshalGQL(buf)

			got, err := strconv.Unquote(buf.String())
			if err != nil {
				t.Fatalf("count not unquote string")
			}
			if got != tc.want {
				t.Errorf("got %s; want %s", got, tc.want)
			}
		})
	}
}

const epsilon = 10 * time.Second

func TestTimestampRelative(t *testing.T) {
	n := time.Now()
	testCases := []struct {
		name string
		have string
		want time.Time
	}{
		{"past", "<4h", n.Add(-4 * time.Hour)},
		{"future", ">5m", n.Add(5 * time.Minute)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := UnmarshalTimestamp(tc.have)
			if err != nil {
				t.Fatalf("could not unmarshal time: %v", err)
			}

			log.Infof("got: %v", got)
			if got.Sub(tc.want) > epsilon {
				t.Errorf("not within bound of %v; got %s; want %s", epsilon, got, tc.want)
			}
		})
	}

}
