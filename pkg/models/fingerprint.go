package models

import (
	"fmt"
	"strconv"
)

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

func (f *Fingerprint) Value() string {
	switch v := f.Fingerprint.(type) {
	case int64:
		return strconv.FormatUint(uint64(v), 16)
	default:
		return fmt.Sprintf("%v", f.Fingerprint)
	}
}

// String returns the string representation of the Fingerprint.
// It will return an empty string if the Fingerprint is not a string.
func (f Fingerprint) String() string {
	s, _ := f.Fingerprint.(string)
	return s
}

// Int64 returns the int64 representation of the Fingerprint.
// It will return 0 if the Fingerprint is not an int64.
func (f Fingerprint) Int64() int64 {
	v, _ := f.Fingerprint.(int64)
	return v
}

type Fingerprints []Fingerprint

func (f Fingerprints) Remove(type_ string) Fingerprints {
	var ret Fingerprints

	for _, ff := range f {
		if ff.Type != type_ {
			ret = append(ret, ff)
		}
	}

	return ret
}

func (f Fingerprints) Filter(types ...string) Fingerprints {
	var ret Fingerprints

	for _, ff := range f {
		for _, t := range types {
			if ff.Type == t {
				ret = append(ret, ff)
				break
			}
		}
	}

	return ret
}

// Equals returns true if the contents of this slice are equal to those in the other slice.
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

// ContentsChanged returns true if this Fingerprints slice contains any Fingerprints that different Fingerprint values for the matching type in other, or if this slice contains any Fingerprint types that are not in other.
func (f Fingerprints) ContentsChanged(other Fingerprints) bool {
	for _, ff := range f {
		oo := other.For(ff.Type)
		if oo == nil || oo.Fingerprint != ff.Fingerprint {
			return true
		}
	}

	return false
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
	fp := f.For(type_)
	if fp == nil {
		return nil
	}
	return fp.Fingerprint
}

func (f Fingerprints) GetString(type_ string) string {
	fp := f.For(type_)
	if fp == nil {
		return ""
	}
	return fp.String()
}

func (f Fingerprints) GetInt64(type_ string) int64 {
	fp := f.For(type_)
	if fp != nil {
		return 0
	}
	return fp.Int64()
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
