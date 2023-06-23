package casefolded

import (
	"unicode"
	"unicode/utf8"
)

// Natural implements sort.Interface to sort strings in natural order. This
// means that e.g. "abc2" < "abc12".
//
// This is the simple case-folded version,
// which means that letters are considered equal if strings.SimpleFold says they are.
// For example, "abc2" < "ABC12" < "abc100" and 'k' == '\u212a' (the Kelvin symbol).
//
// Non-digit sequences and numbers are compared separately.
// The former are compared rune-by-rune using the lowest equivalent runes,
// while digits are compared numerically
// (except that the number of leading zeros is used as a tie-breaker, so e.g. "2" < "02")
//
// Limitations:
//   - only ASCII digits (0-9) are considered.
//   - comparisons are done on a rune-by-rune basis,
//     so some special case equivalences like 'ß' == 'SS" are not supported.
//   - Special cases like Turkish 'i' == 'İ' (and not regular dotless 'I')
//     are not supported either.
type Natural []string

func (n Natural) Len() int           { return len(n) }
func (n Natural) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Natural) Less(i, j int) bool { return NaturalLess(n[i], n[j]) }

func isDigit(b rune) bool { return '0' <= b && b <= '9' }

// caseFold returns the lowest-numbered rune equivalent to the parameter.
func caseFold(r rune) rune {
	// Iterate until SimpleFold returns a lower value.
	// This will be the lowest-numbered equivalent rune.
	var prev rune = -1
	for r > prev {
		prev, r = r, unicode.SimpleFold(r)
	}
	return r
}

// NaturalLess compares two strings using natural ordering. This means that e.g.
// "abc2" < "abc12".
//
// This is the simple case-folded version,
// which means that letters are considered equal if strings.SimpleFold says they are.
// For example, "abc2" < "ABC12" < "abc100" and 'k' == '\u212a' (the Kelvin symbol).
//
// Non-digit sequences and numbers are compared separately.
// The former are compared rune-by-rune using the lowest equivalent runes,
// while digits are compared numerically
// (except that the number of leading zeros is used as a tie-breaker, so e.g. "2" < "02")
//
// Limitations:
//   - only ASCII digits (0-9) are considered.
//   - comparisons are done on a rune-by-rune basis,
//     so some special case equivalences like 'ß' == 'SS" are not supported.
//   - Special cases like Turkish 'i' == 'İ' (and not regular dotless 'I')
//     are not supported either.
func NaturalLess(str1, str2 string) bool {
	// ASCII fast path.
	idx1, idx2 := 0, 0
	for idx1 < len(str1) && idx2 < len(str2) {
		c1, c2 := rune(str1[idx1]), rune(str2[idx2])

		// Bail out to full Unicode support?
		if c1|c2 >= utf8.RuneSelf {
			goto hasUnicode
		}

		dig1, dig2 := isDigit(c1), isDigit(c2)
		switch {
		case !dig1 || !dig2:
			// For ASCII it suffices to normalize letters to upper-case,
			// because upper-cased ASCII compares lexicographically.
			// Note: this does not account for regional special cases
			// like Turkish dotted capital 'İ'.

			// Canonicalize to upper-case.
			c1 = unicode.ToUpper(c1)
			c2 = unicode.ToUpper(c2)
			// Identical upper-cased ASCII runes are equal.
			if c1 == c2 {
				idx1++
				idx2++
				continue
			}
			return c1 < c2
		default: // Digits
			// Eat zeros.
			for ; idx1 < len(str1) && str1[idx1] == '0'; idx1++ {
			}
			for ; idx2 < len(str2) && str2[idx2] == '0'; idx2++ {
			}
			// Eat all digits.
			nonZero1, nonZero2 := idx1, idx2
			for ; idx1 < len(str1) && isDigit(rune(str1[idx1])); idx1++ {
			}
			for ; idx2 < len(str2) && isDigit(rune(str2[idx2])); idx2++ {
			}
			// If lengths of numbers with non-zero prefix differ, the shorter
			// one is less.
			if len1, len2 := idx1-nonZero1, idx2-nonZero2; len1 != len2 {
				return len1 < len2
			}
			// If they're equally long, string comparison is correct.
			if nr1, nr2 := str1[nonZero1:idx1], str2[nonZero2:idx2]; nr1 != nr2 {
				return nr1 < nr2
			}
			// Otherwise, the one with less zeros is less.
			// Because everything up to the number is equal, comparing the index
			// after the zeros is sufficient.
			if nonZero1 != nonZero2 {
				return nonZero1 < nonZero2
			}
		}
		// They're identical so far, so continue comparing.
	}
	// So far they are identical. At least one is ended. If the other continues,
	// it sorts last.
	return len(str1) < len(str2)

hasUnicode:
	for idx1 < len(str1) && idx2 < len(str2) {
		c1, delta1 := utf8.DecodeRuneInString(str1[idx1:])
		c2, delta2 := utf8.DecodeRuneInString(str2[idx2:])

		dig1, dig2 := isDigit(c1), isDigit(c2)
		switch {
		case !dig1 || !dig2:
			idx1 += delta1
			idx2 += delta2
			// Fast path: identical runes are equal.
			if c1 == c2 {
				continue
			}
			// ASCII fast path: ASCII characters compare by their upper-case equivalent (if any)
			// because 'A' < 'a', so upper-case them.
			if c1 <= unicode.MaxASCII && c2 <= unicode.MaxASCII {
				c1 = unicode.ToUpper(c1)
				c2 = unicode.ToUpper(c2)
				if c1 != c2 {
					return c1 < c2
				}
				continue
			}
			// Compare lowest equivalent characters.
			c1 = caseFold(c1)
			c2 = caseFold(c2)
			if c1 == c2 {
				continue
			}
			return c1 < c2
		default: // Digits
			// Eat zeros.
			for ; idx1 < len(str1) && str1[idx1] == '0'; idx1++ {
			}
			for ; idx2 < len(str2) && str2[idx2] == '0'; idx2++ {
			}
			// Eat all digits.
			nonZero1, nonZero2 := idx1, idx2
			for ; idx1 < len(str1) && isDigit(rune(str1[idx1])); idx1++ {
			}
			for ; idx2 < len(str2) && isDigit(rune(str2[idx2])); idx2++ {
			}
			// If lengths of numbers with non-zero prefix differ, the shorter
			// one is less.
			if len1, len2 := idx1-nonZero1, idx2-nonZero2; len1 != len2 {
				return len1 < len2
			}
			// If they're equally long, string comparison is correct.
			if nr1, nr2 := str1[nonZero1:idx1], str2[nonZero2:idx2]; nr1 != nr2 {
				return nr1 < nr2
			}
			// Otherwise, the one with less zeros is less.
			// Because everything up to the number is equal, comparing the index
			// after the zeros is sufficient.
			if nonZero1 != nonZero2 {
				return nonZero1 < nonZero2
			}
		}
		// They're identical so far, so continue comparing.
	}
	// So far they are identical. At least one is ended. If the other continues,
	// it sorts last.
	return len(str1[idx1:]) < len(str2[idx2:])
}
