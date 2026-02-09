// Package sha1 provides utility functions for generating SHA1 hashes.
package sha1

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

// FromBytes returns a SHA1 checksum string from data.
func FromBytes(data []byte) string {
	result := sha1.Sum(data)
	return fmt.Sprintf("%x", result)
}

// FromString returns a SHA1 checksum string from str.
func FromString(str string) string {
	data := []byte(str)
	return FromBytes(data)
}

// FromFilePath returns a SHA1 checksum string for the file at filePath.
// It returns an empty string and an error if an error occurs opening the file.
func FromFilePath(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return FromReader(f)
}

// FromReader returns a SHA1 checksum string from data read from src.
// It returns an empty string and an error if an error occurs reading from src.
func FromReader(src io.Reader) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}
	checksum := h.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}
