// based on code from golang src/time/time.go

package mpd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

var (
	rStart   = "^P"          // Must start with a 'P'
	rDays    = "(\\d+D)?"    // We only allow Days for durations, not Months or Years
	rTime    = "(?:T"        // If there's any 'time' units then they must be preceded by a 'T'
	rHours   = "(\\d+H)?"    // Hours
	rMinutes = "(\\d+M)?"    // Minutes
	rSeconds = "([\\d.]+S)?" // Seconds (Potentially decimal)
	rEnd     = ")?$"         // end of regex must close "T" capture group
)

var xmlDurationRegex = regexp.MustCompile(rStart + rDays + rTime + rHours + rMinutes + rSeconds + rEnd)

func (d Duration) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: d.String()}, nil
}

func (d *Duration) UnmarshalXMLAttr(attr xml.Attr) error {
	dur, err := ParseDuration(attr.Value)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

// String renders a Duration in XML Duration Data Type format
func (d *Duration) String() string {
	// Largest time is 2540400h10m10.000000000s
	var buf [32]byte
	w := len(buf)

	u := uint64(*d)
	neg := *d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 'S'
		w--
		if u == 0 {
			return "PT0S"
		}
		/*
			switch {
			case u < uint64(Millisecond):
				// print microseconds
				prec = 3
				// U+00B5 'µ' micro sign == 0xC2 0xB5
				w-- // Need room for two bytes.
				copy(buf[w:], "µ")
			default:
				// print milliseconds
				prec = 6
				buf[w] = 'm'
			}
		*/
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u)
	} else {
		w--
		buf[w] = 'S'

		w, u = fmtFrac(buf[:w], u, 9)

		// u is now integer seconds
		w = fmtInt(buf[:w], u%60)
		u /= 60

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'M'
			w = fmtInt(buf[:w], u%60)
			u /= 60

			// u is now integer hours
			// Stop at hours because days can be different lengths.
			if u > 0 {
				w--
				buf[w] = 'H'
				w = fmtInt(buf[:w], u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return "PT" + string(buf[w:])
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros.  it omits the decimal
// point too when the fraction is 0.  It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

func ParseDuration(str string) (time.Duration, error) {
	if len(str) < 3 {
		return 0, errors.New("At least one number and designator are required")
	}

	if strings.Contains(str, "-") {
		return 0, errors.New("Duration cannot be negative")
	}

	// Check that only the parts we expect exist and that everything's in the correct order
	if !xmlDurationRegex.Match([]byte(str)) {
		return 0, errors.New("Duration must be in the format: P[nD][T[nH][nM][nS]]")
	}

	var parts = xmlDurationRegex.FindStringSubmatch(str)
	var total time.Duration

	if parts[1] != "" {
		days, err := strconv.Atoi(strings.TrimRight(parts[1], "D"))
		if err != nil {
			return 0, fmt.Errorf("Error parsing Days: %s", err)
		}
		total += time.Duration(days) * time.Hour * 24
	}

	if parts[2] != "" {
		hours, err := strconv.Atoi(strings.TrimRight(parts[2], "H"))
		if err != nil {
			return 0, fmt.Errorf("Error parsing Hours: %s", err)
		}
		total += time.Duration(hours) * time.Hour
	}

	if parts[3] != "" {
		mins, err := strconv.Atoi(strings.TrimRight(parts[3], "M"))
		if err != nil {
			return 0, fmt.Errorf("Error parsing Minutes: %s", err)
		}
		total += time.Duration(mins) * time.Minute
	}

	if parts[4] != "" {
		secs, err := strconv.ParseFloat(strings.TrimRight(parts[4], "S"), 64)
		if err != nil {
			return 0, fmt.Errorf("Error parsing Seconds: %s", err)
		}
		total += time.Duration(secs * float64(time.Second))
	}

	return total, nil
}
