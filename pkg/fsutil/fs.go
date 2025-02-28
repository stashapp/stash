// Package fsutil provides filesystem utility functions for the application.
package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// IsFsPathCaseSensitive checks the fs of the given path to see if it is case sensitive
// if the case sensitivity can not be determined false and an error != nil are returned
func IsFsPathCaseSensitive(path string) (bool, error) {
	// The case sensitivity of the fs of "path" is determined by case flipping
	// the first letter rune from the base string of the path
	// If the resulting flipped path exists then the fs should not be case sensitive
	// ( we check the file mod time to avoid matching an existing path )

	fi, err := os.Stat(path)
	if err != nil { // path cannot be stat'd
		return false, err
	}

	base := filepath.Base(path)
	fBase, err := flipCaseSingle(base)
	if err != nil { // cannot be case flipped
		return false, err
	}
	i := strings.LastIndex(path, base)
	if i < 0 { // shouldn't happen
		return false, fmt.Errorf("could not case flip path %s", path)
	}

	flipped := []byte(path)
	for _, c := range []byte(fBase) { // replace base of path with the flipped one ( we need to flip the base or last dir part )
		flipped[i] = c
		i++
	}

	fiCase, err := os.Stat(string(flipped))
	if err != nil { // cannot stat the case flipped path
		return true, nil // fs of path should be case sensitive
	}

	if fiCase.ModTime() == fi.ModTime() { // file path exists and is the same
		return false, nil // fs of path is not case sensitive
	}
	return false, fmt.Errorf("can not determine case sensitivity of path %s", path)
}

// flipCaseSingle flips the case ( lower<->upper ) of a single char from the string s
// If the string cannot be flipped, the original string value and an error are returned
func flipCaseSingle(s string) (string, error) {
	rr := []rune(s)
	for i, r := range rr {
		if unicode.IsLetter(r) { // look for a letter  to flip
			if unicode.IsUpper(r) {
				rr[i] = unicode.ToLower(r)
				return string(rr), nil
			}
			rr[i] = unicode.ToUpper(r)
			return string(rr), nil
		}

	}
	return s, fmt.Errorf("could not case flip string %s", s)
}
