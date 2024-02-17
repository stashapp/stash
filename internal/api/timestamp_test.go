package api

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

func TestTimestampSymmetry(t *testing.T) {
	n := time.Now()
	buf := bytes.NewBuffer([]byte{})
	MarshalTimestamp(n).MarshalGQL(buf)

	str, err := strconv.Unquote(buf.String())
	if err != nil {
		t.Fatal("could not unquote string")
	}
	got, err := UnmarshalTimestamp(str)
	if err != nil {
		t.Fatalf("could not unmarshal time: %v", err)
	}

	if !n.Equal(got) {
		t.Fatalf("have %v, want %v", got, n)
	}
}

func TestTimestamp(t *testing.T) {
	n := time.Now().In(time.UTC)
	testCases := []struct {
		name string
		have string
		want string
	}{
		{"reflexivity", n.Format(time.RFC3339Nano), n.Format(time.RFC3339Nano)},
		{"rfc3339", "2021-11-04T01:02:03Z", "2021-11-04T01:02:03Z"},
		{"date", "2021-04-05", "2021-04-05T00:00:00Z"},
		{"datetime", "2021-04-05 14:45:36", "2021-04-05T14:45:36Z"},
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

			if got.Sub(tc.want) > epsilon {
				t.Errorf("not within bound of %v; got %s; want %s", epsilon, got, tc.want)
			}
		})
	}

}
