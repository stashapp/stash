package astits

import (
	"fmt"

	"github.com/asticode/go-astikit"
)

const (
	patSectionEntryBytesSize = 4 // 16 bits + 3 reserved + 13 bits = 32 bits
)

// PATData represents a PAT data
// https://en.wikipedia.org/wiki/Program-specific_information
type PATData struct {
	Programs          []*PATProgram
	TransportStreamID uint16
}

// PATProgram represents a PAT program
type PATProgram struct {
	ProgramMapID  uint16 // The packet identifier that contains the associated PMT
	ProgramNumber uint16 // Relates to the Table ID extension in the associated PMT. A value of 0 is reserved for a NIT packet identifier.
}

// parsePATSection parses a PAT section
func parsePATSection(i *astikit.BytesIterator, offsetSectionsEnd int, tableIDExtension uint16) (d *PATData, err error) {
	// Create data
	d = &PATData{TransportStreamID: tableIDExtension}

	// Loop until end of section data is reached
	for i.Offset() < offsetSectionsEnd {
		// Get next bytes
		var bs []byte
		if bs, err = i.NextBytesNoCopy(4); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Append program
		d.Programs = append(d.Programs, &PATProgram{
			ProgramMapID:  uint16(bs[2]&0x1f)<<8 | uint16(bs[3]),
			ProgramNumber: uint16(bs[0])<<8 | uint16(bs[1]),
		})
	}
	return
}

func calcPATSectionLength(d *PATData) uint16 {
	return uint16(4 * len(d.Programs))
}

func writePATSection(w *astikit.BitsWriter, d *PATData) (int, error) {
	b := astikit.NewBitsWriterBatch(w)

	for _, p := range d.Programs {
		b.Write(p.ProgramNumber)
		b.WriteN(uint8(0xff), 3)
		b.WriteN(p.ProgramMapID, 13)
	}

	return len(d.Programs) * patSectionEntryBytesSize, b.Err()
}
