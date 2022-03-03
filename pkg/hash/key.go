// Package hash provides utility functions for generating hashes from strings and random keys.
package hash

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
)

// GenerateRandomKey generates a random string of length l.
// It returns an empty string and an error if an error occurs while generating a random number.
func GenerateRandomKey(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// IntFromString generates a uint64 from a string.
// Values returned by this function are guaranteed to be the same for equal strings.
// They are not guaranteed to be unique for different strings.
func IntFromString(str string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(str))
	return h.Sum64()
}
