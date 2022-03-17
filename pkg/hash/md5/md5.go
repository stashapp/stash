// Package md5 provides utility functions for generating MD5 hashes.
package md5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// FromBytes returns an MD5 checksum string from data.
func FromBytes(data []byte) string {
	result := md5.Sum(data)
	return fmt.Sprintf("%x", result)
}

// FromString returns an MD5 checksum string from str.
func FromString(str string) string {
	data := []byte(str)
	return FromBytes(data)
}

// FromFilePath returns an MD5 checksum string for the file at filePath.
// It returns an empty string and an error if an error occurs opening the file.
func FromFilePath(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return FromReader(f)
}

// FromReader returns an MD5 checksum string from data read from src.
// It returns an empty string and an error if an error occurs reading from src.
func FromReader(src io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}
	checksum := h.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}
