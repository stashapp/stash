package utils

import (
	"fmt"
	"math/rand"
	"strings"
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

// Cut cuts s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
// TODO: This function will be present in go 1.18. When it
// appears, replace calls to utils.Cut with strings.Cut
// replace
func Cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
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

type StrFormatMap map[string]interface{}

// StrFormat formats the provided format string, replacing placeholders
// in the form of "{fieldName}" with the values in the provided
// StrFormatMap.
//
// For example,
// StrFormat("{foo} bar {baz}", StrFormatMap{
//     "foo": "bar",
//     "baz": "abc",
// })
//
// would return: "bar bar abc"
func StrFormat(format string, m StrFormatMap) string {
	args := make([]string, len(m)*2)
	i := 0

	for k, v := range m {
		args[i] = fmt.Sprintf("{%s}", k)
		args[i+1] = fmt.Sprint(v)
		i += 2
	}

	return strings.NewReplacer(args...).Replace(format)
}
