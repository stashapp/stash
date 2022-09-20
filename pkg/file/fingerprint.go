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

func (f Fingerprints) Equals(other Fingerprints) bool {
	if len(f) != len(other) {
		return false
	}

	for _, ff := range f {
		found := false
		for _, oo := range other {
			if ff == oo {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// For returns a pointer to the first Fingerprint element matching the provided type.
func (f Fingerprints) For(type_ string) *Fingerprint {
	for _, fp := range f {
		if fp.Type == type_ {
			return &fp
		}
	}

	return nil
}

func (f Fingerprints) Get(type_ string) interface{} {
	for _, fp := range f {
		if fp.Type == type_ {
			return fp.Fingerprint
		}
	}

	return nil
}

func (f Fingerprints) GetString(type_ string) string {
	fp := f.Get(type_)
	if fp != nil {
		s, _ := fp.(string)
		return s
	}

	return ""
}

func (f Fingerprints) GetInt64(type_ string) int64 {
	fp := f.Get(type_)
	if fp != nil {
		v, _ := fp.(int64)
		return v
	}

	return 0
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
	CalculateFingerprints(f *BaseFile, o Opener, useExisting bool) ([]Fingerprint, error)
}
