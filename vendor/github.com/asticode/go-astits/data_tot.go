package astits

import (
	"fmt"
	"time"

	"github.com/asticode/go-astikit"
)

// TOTData represents a TOT data
// Page: 39 | Chapter: 5.2.6 | Link: https://www.dvb.org/resources/public/standards/a38_dvb-si_specification.pdf
// (barbashov) the link above can be broken, alternative: https://dvb.org/wp-content/uploads/2019/12/a038_tm1217r37_en300468v1_17_1_-_rev-134_-_si_specification.pdf
type TOTData struct {
	Descriptors []*Descriptor
	UTCTime     time.Time
}

// parseTOTSection parses a TOT section
func parseTOTSection(i *astikit.BytesIterator) (d *TOTData, err error) {
	// Create data
	d = &TOTData{}

	// UTC time
	if d.UTCTime, err = parseDVBTime(i); err != nil {
		err = fmt.Errorf("astits: parsing DVB time failed: %w", err)
		return
	}

	// Descriptors
	if d.Descriptors, err = parseDescriptors(i); err != nil {
		err = fmt.Errorf("astits: parsing descriptors failed: %w", err)
		return
	}
	return
}
