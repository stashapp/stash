package astits

import (
	"fmt"
	"github.com/asticode/go-astikit"
)

// Scrambling Controls
const (
	ScramblingControlNotScrambled         = 0
	ScramblingControlReservedForFutureUse = 1
	ScramblingControlScrambledWithEvenKey = 2
	ScramblingControlScrambledWithOddKey  = 3
)

const (
	MpegTsPacketSize       = 188
	mpegTsPacketHeaderSize = 3
	pcrBytesSize           = 6
)

// Packet represents a packet
// https://en.wikipedia.org/wiki/MPEG_transport_stream
type Packet struct {
	AdaptationField *PacketAdaptationField
	Header          *PacketHeader
	Payload         []byte // This is only the payload content
}

// PacketHeader represents a packet header
type PacketHeader struct {
	ContinuityCounter          uint8 // Sequence number of payload packets (0x00 to 0x0F) within each stream (except PID 8191)
	HasAdaptationField         bool
	HasPayload                 bool
	PayloadUnitStartIndicator  bool   // Set when a PES, PSI, or DVB-MIP packet begins immediately following the header.
	PID                        uint16 // Packet Identifier, describing the payload data.
	TransportErrorIndicator    bool   // Set when a demodulator can't correct errors from FEC data; indicating the packet is corrupt.
	TransportPriority          bool   // Set when the current packet has a higher priority than other packets with the same PID.
	TransportScramblingControl uint8
}

// PacketAdaptationField represents a packet adaptation field
type PacketAdaptationField struct {
	AdaptationExtensionField          *PacketAdaptationExtensionField
	DiscontinuityIndicator            bool // Set if current TS packet is in a discontinuity state with respect to either the continuity counter or the program clock reference
	ElementaryStreamPriorityIndicator bool // Set when this stream should be considered "high priority"
	HasAdaptationExtensionField       bool
	HasOPCR                           bool
	HasPCR                            bool
	HasTransportPrivateData           bool
	HasSplicingCountdown              bool
	Length                            int
	IsOneByteStuffing                 bool            // Only used for one byte stuffing - if true, adaptation field will be written as one uint8(0). Not part of TS format
	StuffingLength                    int             // Only used in writePacketAdaptationField to request stuffing
	OPCR                              *ClockReference // Original Program clock reference. Helps when one TS is copied into another
	PCR                               *ClockReference // Program clock reference
	RandomAccessIndicator             bool            // Set when the stream may be decoded without errors from this point
	SpliceCountdown                   int             // Indicates how many TS packets from this one a splicing point occurs (Two's complement signed; may be negative)
	TransportPrivateDataLength        int
	TransportPrivateData              []byte
}

// PacketAdaptationExtensionField represents a packet adaptation extension field
type PacketAdaptationExtensionField struct {
	DTSNextAccessUnit      *ClockReference // The PES DTS of the splice point. Split up as 3 bits, 1 marker bit (0x1), 15 bits, 1 marker bit, 15 bits, and 1 marker bit, for 33 data bits total.
	HasLegalTimeWindow     bool
	HasPiecewiseRate       bool
	HasSeamlessSplice      bool
	LegalTimeWindowIsValid bool
	LegalTimeWindowOffset  uint16 // Extra information for rebroadcasters to determine the state of buffers when packets may be missing.
	Length                 int
	PiecewiseRate          uint32 // The rate of the stream, measured in 188-byte packets, to define the end-time of the LTW.
	SpliceType             uint8  // Indicates the parameters of the H.262 splice.
}

// parsePacket parses a packet
func parsePacket(i *astikit.BytesIterator) (p *Packet, err error) {
	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: getting next byte failed: %w", err)
		return
	}

	// Packet must start with a sync byte
	if b != syncByte {
		err = ErrPacketMustStartWithASyncByte
		return
	}

	// Create packet
	p = &Packet{}

	// In case packet size is bigger than 188 bytes, we don't care for the first bytes
	i.Seek(i.Len() - MpegTsPacketSize + 1)
	offsetStart := i.Offset()

	// Parse header
	if p.Header, err = parsePacketHeader(i); err != nil {
		err = fmt.Errorf("astits: parsing packet header failed: %w", err)
		return
	}

	// Parse adaptation field
	if p.Header.HasAdaptationField {
		if p.AdaptationField, err = parsePacketAdaptationField(i); err != nil {
			err = fmt.Errorf("astits: parsing packet adaptation field failed: %w", err)
			return
		}
	}

	// Build payload
	if p.Header.HasPayload {
		i.Seek(payloadOffset(offsetStart, p.Header, p.AdaptationField))
		p.Payload = i.Dump()
	}
	return
}

// payloadOffset returns the payload offset
func payloadOffset(offsetStart int, h *PacketHeader, a *PacketAdaptationField) (offset int) {
	offset = offsetStart + 3
	if h.HasAdaptationField {
		offset += 1 + a.Length
	}
	return
}

// parsePacketHeader parses the packet header
func parsePacketHeader(i *astikit.BytesIterator) (h *PacketHeader, err error) {
	// Get next bytes
	var bs []byte
	if bs, err = i.NextBytesNoCopy(3); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}

	// Create header
	h = &PacketHeader{
		ContinuityCounter:          uint8(bs[2] & 0xf),
		HasAdaptationField:         bs[2]&0x20 > 0,
		HasPayload:                 bs[2]&0x10 > 0,
		PayloadUnitStartIndicator:  bs[0]&0x40 > 0,
		PID:                        uint16(bs[0]&0x1f)<<8 | uint16(bs[1]),
		TransportErrorIndicator:    bs[0]&0x80 > 0,
		TransportPriority:          bs[0]&0x20 > 0,
		TransportScramblingControl: uint8(bs[2]) >> 6 & 0x3,
	}
	return
}

// parsePacketAdaptationField parses the packet adaptation field
func parsePacketAdaptationField(i *astikit.BytesIterator) (a *PacketAdaptationField, err error) {
	// Create adaptation field
	a = &PacketAdaptationField{}

	// Get next byte
	var b byte
	if b, err = i.NextByte(); err != nil {
		err = fmt.Errorf("astits: fetching next byte failed: %w", err)
		return
	}

	// Length
	a.Length = int(b)

	afStartOffset := i.Offset()

	// Valid length
	if a.Length > 0 {
		// Get next byte
		if b, err = i.NextByte(); err != nil {
			err = fmt.Errorf("astits: fetching next byte failed: %w", err)
			return
		}

		// Flags
		a.DiscontinuityIndicator = b&0x80 > 0
		a.RandomAccessIndicator = b&0x40 > 0
		a.ElementaryStreamPriorityIndicator = b&0x20 > 0
		a.HasPCR = b&0x10 > 0
		a.HasOPCR = b&0x08 > 0
		a.HasSplicingCountdown = b&0x04 > 0
		a.HasTransportPrivateData = b&0x02 > 0
		a.HasAdaptationExtensionField = b&0x01 > 0

		// PCR
		if a.HasPCR {
			if a.PCR, err = parsePCR(i); err != nil {
				err = fmt.Errorf("astits: parsing PCR failed: %w", err)
				return
			}
		}

		// OPCR
		if a.HasOPCR {
			if a.OPCR, err = parsePCR(i); err != nil {
				err = fmt.Errorf("astits: parsing PCR failed: %w", err)
				return
			}
		}

		// Splicing countdown
		if a.HasSplicingCountdown {
			if b, err = i.NextByte(); err != nil {
				err = fmt.Errorf("astits: fetching next byte failed: %w", err)
				return
			}
			a.SpliceCountdown = int(b)
		}

		// Transport private data
		if a.HasTransportPrivateData {
			// Length
			if b, err = i.NextByte(); err != nil {
				err = fmt.Errorf("astits: fetching next byte failed: %w", err)
				return
			}
			a.TransportPrivateDataLength = int(b)

			// Data
			if a.TransportPrivateDataLength > 0 {
				if a.TransportPrivateData, err = i.NextBytes(a.TransportPrivateDataLength); err != nil {
					err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
					return
				}
			}
		}

		// Adaptation extension
		if a.HasAdaptationExtensionField {
			// Create extension field
			a.AdaptationExtensionField = &PacketAdaptationExtensionField{}

			// Get next byte
			if b, err = i.NextByte(); err != nil {
				err = fmt.Errorf("astits: fetching next byte failed: %w", err)
				return
			}

			// Length
			a.AdaptationExtensionField.Length = int(b)
			if a.AdaptationExtensionField.Length > 0 {
				// Get next byte
				if b, err = i.NextByte(); err != nil {
					err = fmt.Errorf("astits: fetching next byte failed: %w", err)
					return
				}

				// Basic
				a.AdaptationExtensionField.HasLegalTimeWindow = b&0x80 > 0
				a.AdaptationExtensionField.HasPiecewiseRate = b&0x40 > 0
				a.AdaptationExtensionField.HasSeamlessSplice = b&0x20 > 0

				// Legal time window
				if a.AdaptationExtensionField.HasLegalTimeWindow {
					var bs []byte
					if bs, err = i.NextBytesNoCopy(2); err != nil {
						err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
						return
					}
					a.AdaptationExtensionField.LegalTimeWindowIsValid = bs[0]&0x80 > 0
					a.AdaptationExtensionField.LegalTimeWindowOffset = uint16(bs[0]&0x7f)<<8 | uint16(bs[1])
				}

				// Piecewise rate
				if a.AdaptationExtensionField.HasPiecewiseRate {
					var bs []byte
					if bs, err = i.NextBytesNoCopy(3); err != nil {
						err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
						return
					}
					a.AdaptationExtensionField.PiecewiseRate = uint32(bs[0]&0x3f)<<16 | uint32(bs[1])<<8 | uint32(bs[2])
				}

				// Seamless splice
				if a.AdaptationExtensionField.HasSeamlessSplice {
					// Get next byte
					if b, err = i.NextByte(); err != nil {
						err = fmt.Errorf("astits: fetching next byte failed: %w", err)
						return
					}

					// Splice type
					a.AdaptationExtensionField.SpliceType = uint8(b&0xf0) >> 4

					// We need to rewind since the current byte is used by the DTS next access unit as well
					i.Skip(-1)

					// DTS Next access unit
					if a.AdaptationExtensionField.DTSNextAccessUnit, err = parsePTSOrDTS(i); err != nil {
						err = fmt.Errorf("astits: parsing DTS failed: %w", err)
						return
					}
				}
			}
		}
	}

	a.StuffingLength = a.Length - (i.Offset() - afStartOffset)

	return
}

// parsePCR parses a Program Clock Reference
// Program clock reference, stored as 33 bits base, 6 bits reserved, 9 bits extension.
func parsePCR(i *astikit.BytesIterator) (cr *ClockReference, err error) {
	var bs []byte
	if bs, err = i.NextBytesNoCopy(6); err != nil {
		err = fmt.Errorf("astits: fetching next bytes failed: %w", err)
		return
	}
	pcr := uint64(bs[0])<<40 | uint64(bs[1])<<32 | uint64(bs[2])<<24 | uint64(bs[3])<<16 | uint64(bs[4])<<8 | uint64(bs[5])
	cr = newClockReference(int64(pcr>>15), int64(pcr&0x1ff))
	return
}

func writePacket(w *astikit.BitsWriter, p *Packet, targetPacketSize int) (written int, retErr error) {
	if retErr = w.Write(uint8(syncByte)); retErr != nil {
		return
	}
	written += 1

	n, retErr := writePacketHeader(w, p.Header)
	if retErr != nil {
		return
	}
	written += n

	if p.Header.HasAdaptationField {
		n, retErr = writePacketAdaptationField(w, p.AdaptationField)
		if retErr != nil {
			return
		}
		written += n
	}

	if targetPacketSize-written < len(p.Payload) {
		return 0, fmt.Errorf(
			"writePacket: can't write %d bytes of payload: only %d is available",
			len(p.Payload),
			targetPacketSize-written,
		)
	}

	if p.Header.HasPayload {
		retErr = w.Write(p.Payload)
		if retErr != nil {
			return
		}
		written += len(p.Payload)
	}

	for written < targetPacketSize {
		if retErr = w.Write(uint8(0xff)); retErr != nil {
			return
		}
		written++
	}

	return written, nil
}

func writePacketHeader(w *astikit.BitsWriter, h *PacketHeader) (written int, retErr error) {
	b := astikit.NewBitsWriterBatch(w)

	b.Write(h.TransportErrorIndicator)
	b.Write(h.PayloadUnitStartIndicator)
	b.Write(h.TransportPriority)
	b.WriteN(h.PID, 13)
	b.WriteN(h.TransportScramblingControl, 2)
	b.Write(h.HasAdaptationField) // adaptation_field_control higher bit
	b.Write(h.HasPayload)         // adaptation_field_control lower bit
	b.WriteN(h.ContinuityCounter, 4)

	return mpegTsPacketHeaderSize, b.Err()
}

func writePCR(w *astikit.BitsWriter, cr *ClockReference) (int, error) {
	b := astikit.NewBitsWriterBatch(w)

	b.WriteN(uint64(cr.Base), 33)
	b.WriteN(uint8(0xff), 6)
	b.WriteN(uint64(cr.Extension), 9)
	return pcrBytesSize, b.Err()
}

func calcPacketAdaptationFieldLength(af *PacketAdaptationField) (length uint8) {
	length++
	if af.HasPCR {
		length += pcrBytesSize
	}
	if af.HasOPCR {
		length += pcrBytesSize
	}
	if af.HasSplicingCountdown {
		length++
	}
	if af.HasTransportPrivateData {
		length += 1 + uint8(len(af.TransportPrivateData))
	}
	if af.HasAdaptationExtensionField {
		length += 1 + calcPacketAdaptationFieldExtensionLength(af.AdaptationExtensionField)
	}
	length += uint8(af.StuffingLength)
	return
}

func writePacketAdaptationField(w *astikit.BitsWriter, af *PacketAdaptationField) (bytesWritten int, retErr error) {
	b := astikit.NewBitsWriterBatch(w)

	if af.IsOneByteStuffing {
		b.Write(uint8(0))
		return 1, nil
	}

	length := calcPacketAdaptationFieldLength(af)
	b.Write(length)
	bytesWritten++

	b.Write(af.DiscontinuityIndicator)
	b.Write(af.RandomAccessIndicator)
	b.Write(af.ElementaryStreamPriorityIndicator)
	b.Write(af.HasPCR)
	b.Write(af.HasOPCR)
	b.Write(af.HasSplicingCountdown)
	b.Write(af.HasTransportPrivateData)
	b.Write(af.HasAdaptationExtensionField)

	bytesWritten++

	if af.HasPCR {
		n, err := writePCR(w, af.PCR)
		if err != nil {
			return 0, err
		}
		bytesWritten += n
	}

	if af.HasOPCR {
		n, err := writePCR(w, af.OPCR)
		if err != nil {
			return 0, err
		}
		bytesWritten += n
	}

	if af.HasSplicingCountdown {
		b.Write(uint8(af.SpliceCountdown))
		bytesWritten++
	}

	if af.HasTransportPrivateData {
		// we can get length from TransportPrivateData itself, why do we need separate field?
		b.Write(uint8(af.TransportPrivateDataLength))
		bytesWritten++
		if af.TransportPrivateDataLength > 0 {
			b.Write(af.TransportPrivateData)
		}
		bytesWritten += len(af.TransportPrivateData)
	}

	if af.HasAdaptationExtensionField {
		n, err := writePacketAdaptationFieldExtension(w, af.AdaptationExtensionField)
		if err != nil {
			return 0, err
		}
		bytesWritten += n
	}

	// stuffing
	for i := 0; i < af.StuffingLength; i++ {
		b.Write(uint8(0xff))
		bytesWritten++
	}

	retErr = b.Err()
	return
}

func calcPacketAdaptationFieldExtensionLength(afe *PacketAdaptationExtensionField) (length uint8) {
	length++
	if afe.HasLegalTimeWindow {
		length += 2
	}
	if afe.HasPiecewiseRate {
		length += 3
	}
	if afe.HasSeamlessSplice {
		length += ptsOrDTSByteLength
	}
	return length
}

func writePacketAdaptationFieldExtension(w *astikit.BitsWriter, afe *PacketAdaptationExtensionField) (bytesWritten int, retErr error) {
	b := astikit.NewBitsWriterBatch(w)

	length := calcPacketAdaptationFieldExtensionLength(afe)
	b.Write(length)
	bytesWritten++

	b.Write(afe.HasLegalTimeWindow)
	b.Write(afe.HasPiecewiseRate)
	b.Write(afe.HasSeamlessSplice)
	b.WriteN(uint8(0xff), 5) // reserved
	bytesWritten++

	if afe.HasLegalTimeWindow {
		b.Write(afe.LegalTimeWindowIsValid)
		b.WriteN(afe.LegalTimeWindowOffset, 15)
		bytesWritten += 2
	}

	if afe.HasPiecewiseRate {
		b.WriteN(uint8(0xff), 2)
		b.WriteN(afe.PiecewiseRate, 22)
		bytesWritten += 3
	}

	if afe.HasSeamlessSplice {
		n, err := writePTSOrDTS(w, afe.SpliceType, afe.DTSNextAccessUnit)
		if err != nil {
			return 0, err
		}
		bytesWritten += n
	}

	retErr = b.Err()
	return
}

func newStuffingAdaptationField(bytesToStuff int) *PacketAdaptationField {
	if bytesToStuff == 1 {
		return &PacketAdaptationField{
			IsOneByteStuffing: true,
		}
	}

	return &PacketAdaptationField{
		// one byte for length and one for flags
		StuffingLength: bytesToStuff - 2,
	}
}
