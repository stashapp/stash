package utils

import (
	"strconv"
	"time"
)

// GetVTTTime returns a timestamp appropriate for VTT files (hh:mm:ss)
func GetVTTTime(totalSeconds float64) (s string) {
	totalSecondsString := strconv.FormatFloat(totalSeconds, 'f', -1, 64)
	secondsDuration, _ := time.ParseDuration(totalSecondsString + "s")

	// Hours
	var hours = int(secondsDuration / time.Hour)
	var n = secondsDuration % time.Hour
	if hours < 10 {
		s += "0"
	}
	s += strconv.Itoa(hours) + ":"

	// Minutes
	var minutes = int(n / time.Minute)
	n = secondsDuration % time.Minute
	if minutes < 10 {
		s += "0"
	}
	s += strconv.Itoa(minutes) + ":"

	// Seconds
	var seconds = int(n / time.Second)
	n = secondsDuration % time.Second
	if seconds < 10 {
		s += "0"
	}
	s += strconv.Itoa(seconds)

	return
}
