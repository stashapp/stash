package utils

import (
	"fmt"
	"math"
)

// from stdlib's time.go
func norm(hi, lo, base int) (nhi, nlo int) {
	if lo < 0 {
		n := (-lo-1)/base + 1
		hi -= n
		lo += n * base
	}
	if lo >= base {
		n := lo / base
		hi += n
		lo -= n * base
	}
	return hi, lo
}

// GetVTTTime returns a timestamp appropriate for VTT files (hh:mm:ss.mmm)
func GetVTTTime(fracSeconds float64) string {
	if fracSeconds < 0 || math.IsNaN(fracSeconds) || math.IsInf(fracSeconds, 0) {
		return "00:00:00.000"
	}

	var msec, sec, mnt, hour int
	msec = int(fracSeconds * 1000)
	sec, msec = norm(sec, msec, 1000)
	mnt, sec = norm(mnt, sec, 60)
	hour, mnt = norm(hour, mnt, 60)

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hour, mnt, sec, msec)

}
