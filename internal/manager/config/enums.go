package config

import (
	"fmt"
	"io"
	"strconv"
)

type BlobsStorageType string

const (
	// Database
	BlobStorageTypeDatabase BlobsStorageType = "DATABASE"
	// Filesystem
	BlobStorageTypeFilesystem BlobsStorageType = "FILESYSTEM"
)

var AllBlobStorageType = []BlobsStorageType{
	BlobStorageTypeDatabase,
	BlobStorageTypeFilesystem,
}

func (e BlobsStorageType) IsValid() bool {
	switch e {
	case BlobStorageTypeDatabase, BlobStorageTypeFilesystem:
		return true
	}
	return false
}

func (e BlobsStorageType) String() string {
	return string(e)
}

func (e *BlobsStorageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BlobsStorageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BlobStorageType", str)
	}
	return nil
}

func (e BlobsStorageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
