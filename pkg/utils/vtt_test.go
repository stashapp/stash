package utils

import (
	"math"
	"testing"
)

func TestZeroTimestamp(t *testing.T) {
	if want, got := "00:00:00.000", GetVTTTime(0); want != got {
		t.Errorf("TestZeroTimestamp: GetVTTTime(0) = %v; want %v", got, want)
	}
}

func TestValidTimestamp(t *testing.T) {
	s := 0.1
	if want, got := "00:00:00.100", GetVTTTime(s); want != got {
		t.Errorf("TestValidTimestamp: GetVTTTime(%v) = %v; want %v", s, got, want)
	}
	s = ((24+1)*60+1)*60 + 1 + 0.1
	if want, got := "25:01:01.100", GetVTTTime(s); want != got {
		t.Errorf("TestValidTimestamp: GetVTTTime(%v) = %v; want %v", s, got, want)
	}
}

// Negative timestamps are not defined by WebVTT.
func TestNegativeTimestamp(t *testing.T) {
	if want, got := "00:00:00.000", GetVTTTime(-1); want != got {
		t.Errorf("TestNegativeTimestamp: GetVTTTime(-1) = %v; want %v", got, want)
	}
}

func TestInvalidTimestamp(t *testing.T) {
	if want, got := "00:00:00.000", GetVTTTime(math.NaN()); want != got {
		t.Errorf("TestInvalidTimestamp: GetVTTTime(NaN) = %v; want %v", got, want)
	}
	if want, got := "00:00:00.000", GetVTTTime(math.Inf(1)); want != got {
		t.Errorf("TestInvalidTimestamp: GetVTTTime(Inf) = %v; want %v", got, want)
	}
	if want, got := "00:00:00.000", GetVTTTime(math.Inf(-1)); want != got {
		t.Errorf("TestInvalidTimestamp: GetVTTTime(-Inf) = %v; want %v", got, want)
	}
}
