package models

import (
	"fmt"
	"io"
	"strconv"
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
