package transport

import (
	"errors"
	"io"
)

type bytesReader struct {
	s *[]byte
	i int64 // current reading index
}

func (r *bytesReader) Read(b []byte) (n int, err error) {
	if r.s == nil {
		return 0, errors.New("byte slice pointer is nil")
	}
	if r.i >= int64(len(*r.s)) {
		return 0, io.EOF
	}
	n = copy(b, (*r.s)[r.i:])
	r.i += int64(n)
	return
}

func (r *bytesReader) Seek(offset int64, whence int) (int64, error) {
	if r.s == nil {
		return 0, errors.New("byte slice pointer is nil")
	}
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(*r.s)) + offset
	default:
		return 0, errors.New("invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("negative position")
	}
	r.i = abs
	return abs, nil
}
