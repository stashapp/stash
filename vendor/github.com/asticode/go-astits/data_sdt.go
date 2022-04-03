package astits

import (
	"fmt"

	"github.com/asticode/go-astikit"
)

// Running statuses
const (
	RunningStatusNotRunning          = 1
	RunningStatusPausing             = 3
	RunningStatusRunning             = 4
	RunningStatusServiceOffAir       = 5
	RunningStatusStartsInAFewSeconds = 2
	RunningStatusUndefined           = 0
)

// SDTData represents an SDT data
// Page: 33 | Chapter: 5.2.3 | Link: https://www.dvb.org/resources/public/standards/a38_dvb-si_specification.pdf
// (barbashov) the link above can be broken, alternative: https://dvb.org/wp-content/uploads/2019/12/a038_tm1217r37_en300468v1_17_1_-_rev-134_-_si_specification.pdf
type SDTData struct {
	OriginalNetworkID uint16
	Services          []*SDTDataService
	TransportStreamID uint16
}

// SDTDataService represents an SDT data service
type SDTDataService struct {
	Descriptors            []*Descriptor
	HasEITPresentFollowing bool // When true indicates that EIT present/following information for the service is present in the current TS
	HasEITSchedule         bool // When true indicates that EIT schedule information for the service is present in the current TS
	HasFreeCSAMode         bool // When true indicates that access to one or more streams may be controlled by a CA system.
	RunningStatus          uint8
	ServiceID              uint16
}

// parseSDTSection parses an SDT section
func parseSDTSection(i *astikit.BytesIterator, offsetSectionsEnd int, tableIDExtension uint16) (d *SDTData, err error) {
	// Create data
	d = &SDTData{TransportStreamID: tableIDExtension}

	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(2); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Original network ID
	d.OriginalNetworkID = uint16(bs[0])<<8 | uint16(bs[1])

	// Reserved for future use
	i.Skip(1)

	// Loop until end of section data is reached
	for i.Offset() < offsetSectionsEnd {
		// Create service
		s := &SDTDataService{}

		// Get next bytes
		if bs, err = i.NextBytesNoCopy(2); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Service ID
		s.ServiceID = uint16(bs[0])<<8 | uint16(bs[1])

		// Get next byte
		var b byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// EIT schedule flag
		s.HasEITSchedule = uint8(b&0x2) > 0

		// EIT present/following flag
		s.HasEITPresentFollowing = uint8(b&0x1) > 0

		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Running status
		s.RunningStatus = uint8(b) >> 5

		// Free CA mode
		s.HasFreeCSAMode = uint8(b&0x10) > 0

		// We need to rewind since the current byte is used by the descriptor as well
		i.Skip(-1)

		// Descriptors
		if s.Descriptors, err = parseDescriptors(i); err != nil {
			err = fmt.Errorf("astits: parsing descriptors failed: %w", err)
			return
		}

		// Append service
		d.Services = append(d.Services, s)
	}
	return
}
