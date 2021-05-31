package utils

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"
)

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandomSequence(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

// FlipCaseSingle flips the case ( lower<->upper ) of a single char from the string s
// If the string cannot be flipped, the original string value and an error are returned
func FlipCaseSingle(s string) (string, error) {
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
