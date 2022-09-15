package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type HashAlgorithm string

const (
	HashAlgorithmMd5 HashAlgorithm = "MD5"
	// oshash
	HashAlgorithmOshash HashAlgorithm = "OSHASH"
)

var AllHashAlgorithm = []HashAlgorithm{
	HashAlgorithmMd5,
	HashAlgorithmOshash,
}

func (e HashAlgorithm) IsValid() bool {
	switch e {
	case HashAlgorithmMd5, HashAlgorithmOshash:
		return true
	}
	return false
}

func (e HashAlgorithm) String() string {
	return string(e)
}

func (e *HashAlgorithm) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HashAlgorithm(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HashAlgorithm", str)
	}
	return nil
}

func (e HashAlgorithm) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type File struct {
	Checksum    string    `db:"checksum" json:"checksum"`
	OSHash      string    `db:"oshash" json:"oshash"`
	Path        string    `db:"path" json:"path"`
	Size        string    `db:"size" json:"size"`
	FileModTime time.Time `db:"file_mod_time" json:"file_mod_time"`
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s File) GetHash(hashAlgorithm HashAlgorithm) string {
	switch hashAlgorithm {
	case HashAlgorithmMd5:
		return s.Checksum
	case HashAlgorithmOshash:
		return s.OSHash
	default:
		panic("unknown hash algorithm")
	}
}

func (s File) Equal(o File) bool {
	return s.Path == o.Path && s.Checksum == o.Checksum && s.OSHash == o.OSHash && s.Size == o.Size && s.FileModTime.Equal(o.FileModTime)
}
