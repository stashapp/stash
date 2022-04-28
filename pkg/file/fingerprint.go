package file

var (
	FingerprintTypeOshash = "oshash"
	FingerprintTypeMD5    = "md5"
	FingerprintTypePhash  = "phash"
)

// Fingerprint represents a fingerprint of a file.
type Fingerprint struct {
	Type        string
	Fingerprint interface{}
}

type Fingerprints []Fingerprint

func (f Fingerprints) Get(type_ string) interface{} {
	for _, fp := range f {
		if fp.Type == type_ {
			return fp.Fingerprint
		}
	}

	return nil
}

// FingerprintCalculator calculates a fingerprint for the provided file.
type FingerprintCalculator interface {
	CalculateFingerprints(f *BaseFile, o Opener) ([]Fingerprint, error)
}
