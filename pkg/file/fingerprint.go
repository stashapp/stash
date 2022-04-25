package file

import (
	"io"
)

type Fingerprint struct {
	Type        string
	Fingerprint interface{}
}

type FingerprintHandler interface {
	CalculateFingerprint(f *BasicFile, r io.Reader) (*Fingerprint, error)
}
