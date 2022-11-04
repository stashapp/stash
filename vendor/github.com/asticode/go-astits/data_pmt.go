package astits

import (
	"fmt"

	"github.com/asticode/go-astikit"
)

type StreamType uint8

// Stream types
const (
	StreamTypeMPEG1Video                 StreamType = 0x01
	StreamTypeMPEG2Video                 StreamType = 0x02
	StreamTypeMPEG1Audio                 StreamType = 0x03 // ISO/IEC 11172-3
	StreamTypeMPEG2HalvedSampleRateAudio StreamType = 0x04 // ISO/IEC 13818-3
	StreamTypeMPEG2Audio                 StreamType = 0x04
	StreamTypePrivateSection             StreamType = 0x05
	StreamTypePrivateData                StreamType = 0x06
	StreamTypeMPEG2PacketizedData        StreamType = 0x06 // Rec. ITU-T H.222 | ISO/IEC 13818-1 i.e., DVB subtitles/VBI and AC-3
	StreamTypeADTS                       StreamType = 0x0F // ISO/IEC 13818-7 Audio with ADTS transport syntax
	StreamTypeAACAudio                   StreamType = 0x0f
	StreamTypeMPEG4Video                 StreamType = 0x10
	StreamTypeAACLATMAudio               StreamType = 0x11
	StreamTypeMetadata                   StreamType = 0x15
	StreamTypeH264Video                  StreamType = 0x1B // Rec. ITU-T H.264 | ISO/IEC 14496-10
	StreamTypeH265Video                  StreamType = 0x24 // Rec. ITU-T H.265 | ISO/IEC 23008-2
	StreamTypeHEVCVideo                  StreamType = 0x24
	StreamTypeCAVSVideo                  StreamType = 0x42
	StreamTypeVC1Video                   StreamType = 0xea
	StreamTypeDIRACVideo                 StreamType = 0xd1
	StreamTypeAC3Audio                   StreamType = 0x81
	StreamTypeDTSAudio                   StreamType = 0x82
	StreamTypeTRUEHDAudio                StreamType = 0x83
	StreamTypeSCTE35                     StreamType = 0x86
	StreamTypeEAC3Audio                  StreamType = 0x87
)

// PMTData represents a PMT data
// https://en.wikipedia.org/wiki/Program-specific_information
type PMTData struct {
	ElementaryStreams  []*PMTElementaryStream
	PCRPID             uint16        // The packet identifier that contains the program clock reference used to improve the random access accuracy of the stream's timing that is derived from the program timestamp. If this is unused. then it is set to 0x1FFF (all bits on).
	ProgramDescriptors []*Descriptor // Program descriptors
	ProgramNumber      uint16
}

// PMTElementaryStream represents a PMT elementary stream
type PMTElementaryStream struct {
	ElementaryPID               uint16        // The packet identifier that contains the stream type data.
	ElementaryStreamDescriptors []*Descriptor // Elementary stream descriptors
	StreamType                  StreamType    // This defines the structure of the data contained within the elementary packet identifier.
}

// parsePMTSection parses a PMT section
func parsePMTSection(i *astikit.BytesIterator, offsetSectionsEnd int, tableIDExtension uint16) (d *PMTData, err error) {
	// Create data
	d = &PMTData{ProgramNumber: tableIDExtension}

	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(2); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// PCR PID
	d.PCRPID = uint16(bs[0]&0x1f)<<8 | uint16(bs[1])

	// Program descriptors
	if d.ProgramDescriptors, err = parseDescriptors(i); err != nil {
		err = fmt.Errorf("astits: parsing descriptors failed: %w", err)
		return
	}

	// Loop until end of section data is reached
	for i.Offset() < offsetSectionsEnd {
		// Create stream
		e := &PMTElementaryStream{}

		// Get next byte
		var b byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Stream type
		e.StreamType = StreamType(b)

		// Get next bytes
		if bs, err = i.NextBytesNoCopy(2); err != nil {
			err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
			return
		}

		// Elementary PID
		e.ElementaryPID = uint16(bs[0]&0x1f)<<8 | uint16(bs[1])

		// Elementary descriptors
		if e.ElementaryStreamDescriptors, err = parseDescriptors(i); err != nil {
			err = fmt.Errorf("astits: parsing descriptors failed: %w", err)
			return
		}

		// Add elementary stream
		d.ElementaryStreams = append(d.ElementaryStreams, e)
	}
	return
}

func calcPMTProgramInfoLength(d *PMTData) uint16 {
	ret := uint16(2) // program_info_length
	ret += calcDescriptorsLength(d.ProgramDescriptors)

	for _, es := range d.ElementaryStreams {
		ret += 5 // stream_type, elementary_pid, es_info_length
		ret += calcDescriptorsLength(es.ElementaryStreamDescriptors)
	}

	return ret
}

func calcPMTSectionLength(d *PMTData) uint16 {
	ret := uint16(4)
	ret += calcDescriptorsLength(d.ProgramDescriptors)

	for _, es := range d.ElementaryStreams {
		ret += 5
		ret += calcDescriptorsLength(es.ElementaryStreamDescriptors)
	}

	return ret
}

func writePMTSection(w *astikit.BitsWriter, d *PMTData) (int, error) {
	b := astikit.NewBitsWriterBatch(w)

	// TODO split into sections

	b.WriteN(uint8(0xff), 3)
	b.WriteN(d.PCRPID, 13)
	bytesWritten := 2

	n, err := writeDescriptorsWithLength(w, d.ProgramDescriptors)
	if err != nil {
		return 0, err
	}
	bytesWritten += n

	for _, es := range d.ElementaryStreams {
		b.Write(uint8(es.StreamType))
		b.WriteN(uint8(0xff), 3)
		b.WriteN(es.ElementaryPID, 13)
		bytesWritten += 3

		n, err = writeDescriptorsWithLength(w, es.ElementaryStreamDescriptors)
		if err != nil {
			return 0, err
		}
		bytesWritten += n
	}

	return bytesWritten, b.Err()
}

func (t StreamType) IsVideo() bool {
	switch t {
	case StreamTypeMPEG1Video,
		StreamTypeMPEG2Video,
		StreamTypeMPEG4Video,
		StreamTypeH264Video,
		StreamTypeH265Video,
		StreamTypeCAVSVideo,
		StreamTypeVC1Video,
		StreamTypeDIRACVideo:
		return true
	}
	return false
}

func (t StreamType) IsAudio() bool {
	switch t {
	case StreamTypeMPEG1Audio,
		StreamTypeMPEG2Audio,
		StreamTypeAACAudio,
		StreamTypeAACLATMAudio,
		StreamTypeAC3Audio,
		StreamTypeDTSAudio,
		StreamTypeTRUEHDAudio,
		StreamTypeEAC3Audio:
		return true
	}
	return false
}

func (t StreamType) String() string {
	switch t {
	case StreamTypeMPEG1Video:
		return "MPEG1 Video"
	case StreamTypeMPEG2Video:
		return "MPEG2 Video"
	case StreamTypeMPEG1Audio:
		return "MPEG1 Audio"
	case StreamTypeMPEG2Audio:
		return "MPEG2 Audio"
	case StreamTypePrivateSection:
		return "Private Section"
	case StreamTypePrivateData:
		return "Private Data"
	case StreamTypeAACAudio:
		return "AAC Audio"
	case StreamTypeMPEG4Video:
		return "MPEG4 Video"
	case StreamTypeAACLATMAudio:
		return "AAC LATM Audio"
	case StreamTypeMetadata:
		return "Metadata"
	case StreamTypeH264Video:
		return "H264 Video"
	case StreamTypeH265Video:
		return "H265 Video"
	case StreamTypeCAVSVideo:
		return "CAVS Video"
	case StreamTypeVC1Video:
		return "VC1 Video"
	case StreamTypeDIRACVideo:
		return "DIRAC Video"
	case StreamTypeAC3Audio:
		return "AC3 Audio"
	case StreamTypeDTSAudio:
		return "DTS Audio"
	case StreamTypeTRUEHDAudio:
		return "TRUEHD Audio"
	case StreamTypeSCTE35:
		return "SCTE 35"
	case StreamTypeEAC3Audio:
		return "EAC3 Audio"
	}
	return "Unknown"
}

func (t StreamType) ToPESStreamID() uint8 {
	switch t {
	case StreamTypeMPEG1Video, StreamTypeMPEG2Video, StreamTypeMPEG4Video, StreamTypeH264Video,
		StreamTypeH265Video, StreamTypeCAVSVideo, StreamTypeVC1Video:
		return 0xe0
	case StreamTypeDIRACVideo:
		return 0xfd
	case StreamTypeMPEG2Audio, StreamTypeAACAudio, StreamTypeAACLATMAudio:
		return 0xc0
	case StreamTypeAC3Audio, StreamTypeEAC3Audio: // m2ts_mode???
		return 0xfd
	case StreamTypePrivateSection, StreamTypePrivateData, StreamTypeMetadata:
		return 0xfc
	default:
		return 0xbd
	}
}
