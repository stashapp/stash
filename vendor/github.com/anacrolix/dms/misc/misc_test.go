package misc

import (
	"testing"
	"time"
)

func TestFormatDurationSexagesimal(t *testing.T) {
	expected := "0:22:57.628452"
	actual := FormatDurationSexagesimal(time.Duration(1377628452000))
	if actual != expected {
		t.Fatalf("got %q, expected %q", actual, expected)
	}
}
