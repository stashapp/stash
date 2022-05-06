package astikit

import "fmt"

// BytesIterator represents an object capable of iterating sequentially and safely
// through a slice of bytes. This is particularly useful when you need to iterate
// through a slice of bytes and don't want to check for "index out of range" errors
// manually.
type BytesIterator struct {
	bs     []byte
	offset int
}

// NewBytesIterator creates a new BytesIterator
func NewBytesIterator(bs []byte) *BytesIterator {
	return &BytesIterator{bs: bs}
}

// NextByte returns the next byte
func (i *BytesIterator) NextByte() (b byte, err error) {
	if len(i.bs) < i.offset+1 {
		err = fmt.Errorf("astikit: slice length is %d, offset %d is invalid", len(i.bs), i.offset)
		return
	}
	b = i.bs[i.offset]
	i.offset++
	return
}

// NextBytes returns the n next bytes
func (i *BytesIterator) NextBytes(n int) (bs []byte, err error) {
	if len(i.bs) < i.offset+n {
		err = fmt.Errorf("astikit: slice length is %d, offset %d is invalid", len(i.bs), i.offset+n)
		return
	}
	bs = make([]byte, n)
	copy(bs, i.bs[i.offset:i.offset+n])
	i.offset += n
	return
}

// NextBytesNoCopy returns the n next bytes
// Be careful with this function as it doesn't make a copy of returned data.
// bs will point to internal BytesIterator buffer.
// If you need to modify returned bytes or store it for some time, use NextBytes instead
func (i *BytesIterator) NextBytesNoCopy(n int) (bs []byte, err error) {
	if len(i.bs) < i.offset+n {
		err = fmt.Errorf("astikit: slice length is %d, offset %d is invalid", len(i.bs), i.offset+n)
		return
	}
	bs = i.bs[i.offset : i.offset+n]
	i.offset += n
	return
}

// Seek seeks to the nth byte
func (i *BytesIterator) Seek(n int) {
	i.offset = n
}

// Skip skips the n previous/next bytes
func (i *BytesIterator) Skip(n int) {
	i.offset += n
}

// HasBytesLeft checks whether there are bytes left
func (i *BytesIterator) HasBytesLeft() bool {
	return i.offset < len(i.bs)
}

// Offset returns the offset
func (i *BytesIterator) Offset() int {
	return i.offset
}

// Dump dumps the rest of the slice
func (i *BytesIterator) Dump() (bs []byte) {
	if !i.HasBytesLeft() {
		return
	}
	bs = make([]byte, len(i.bs)-i.offset)
	copy(bs, i.bs[i.offset:len(i.bs)])
	i.offset = len(i.bs)
	return
}

// Len returns the slice length
func (i *BytesIterator) Len() int {
	return len(i.bs)
}

const (
	padRight = "right"
	padLeft  = "left"
)

type bytesPadder struct {
	cut       bool
	direction string
	length    int
	repeat    byte
}

func newBytesPadder(repeat byte, length int) *bytesPadder {
	return &bytesPadder{
		direction: padLeft,
		length:    length,
		repeat:    repeat,
	}
}

func (p *bytesPadder) pad(i []byte) []byte {
	if len(i) == p.length {
		return i
	} else if len(i) > p.length {
		if p.cut {
			return i[:p.length]
		}
		return i
	} else {
		o := make([]byte, len(i))
		copy(o, i)
		for idx := 0; idx < p.length-len(i); idx++ {
			if p.direction == padRight {
				o = append(o, p.repeat)
			} else {
				o = append([]byte{p.repeat}, o...)
			}
			o = append(o, p.repeat)
		}
		o = o[:p.length]
		return o
	}
}

// PadOption represents a Pad option
type PadOption func(p *bytesPadder)

// PadCut is a PadOption
// It indicates to the padder it must cut the input to the provided length
// if its original length is bigger
func PadCut(p *bytesPadder) { p.cut = true }

// PadLeft is a PadOption
// It indicates additionnal bytes have to be added to the left
func PadLeft(p *bytesPadder) { p.direction = padLeft }

// PadRight is a PadOption
// It indicates additionnal bytes have to be added to the right
func PadRight(p *bytesPadder) { p.direction = padRight }

// BytesPad pads the slice of bytes with additionnal options
func BytesPad(i []byte, repeat byte, length int, options ...PadOption) []byte {
	p := newBytesPadder(repeat, length)
	for _, o := range options {
		o(p)
	}
	return p.pad(i)
}

// StrPad pads the string with additionnal options
func StrPad(i string, repeat rune, length int, options ...PadOption) string {
	return string(BytesPad([]byte(i), byte(repeat), length, options...))
}
