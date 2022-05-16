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

// AppendUnique appends a fingerprint to the list if a Fingerprint of the same type does not already exist in the list. If one does, then it is updated with o's Fingerprint value.
func (f Fingerprints) AppendUnique(o Fingerprint) Fingerprints {
	ret := f
	for i, fp := range ret {
		if fp.Type == o.Type {
			ret[i] = o
			return ret
		}
	}

	return append(f, o)
}

// FingerprintCalculator calculates a fingerprint for the provided file.
type FingerprintCalculator interface {
	CalculateFingerprints(f *BaseFile, o Opener) ([]Fingerprint, error)
}
