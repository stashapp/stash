package sqlgen

import "time"

var timeLocation = time.UTC

// Set the location to use when interpolating time.Time instances. See https://golang.org/pkg/time/#LoadLocation
// NOTE: This has no effect when using prepared statements.
func SetTimeLocation(loc *time.Location) {
	timeLocation = loc
}

func GetTimeLocation() *time.Location {
	return timeLocation
}
