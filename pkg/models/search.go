package models

import "strings"

const (
	or         = "OR"
	orSymbol   = "|"
	notPrefix  = '-'
	phraseChar = '"'
)

// SearchSpecs provides the specifications for text-based searches.
type SearchSpecs struct {
	// MustHave specifies all of the terms that must appear in the results.
	MustHave []string

	// AnySets specifies sets of terms where one of each set must appear in the results.
	AnySets [][]string

	// MustNot specifies all terms that must not appear in the results.
	MustNot []string
}

// combinePhrases detects quote characters at the start and end of
// words and combines the contents into a single word.
func combinePhrases(words []string) []string {
	var ret []string
	startIndex := -1
	for i, w := range words {
		if startIndex == -1 {
			// looking for start of phrase
			// this could either be " or -"
			ww := w
			if len(w) > 0 && w[0] == notPrefix {
				ww = w[1:]
			}
			if len(ww) > 0 && ww[0] == phraseChar && (len(ww) < 2 || ww[len(ww)-1] != phraseChar) {
				startIndex = i
				continue
			}

			ret = append(ret, w)
		} else if len(w) > 0 && w[len(w)-1] == phraseChar { // looking for end of phrase
			// combine words
			phrase := strings.Join(words[startIndex:i+1], " ")

			// add to return value
			ret = append(ret, phrase)
			startIndex = -1
		}
	}

	if startIndex != -1 {
		ret = append(ret, words[startIndex:]...)
	}

	return ret
}

func extractOrConditions(words []string, searchSpec *SearchSpecs) []string {
	for foundOr := true; foundOr; {
		foundOr = false
		for i, w := range words {
			if i > 0 && i < len(words)-1 && (strings.EqualFold(w, or) || w == orSymbol) {
				// found an OR keyword
				// first operand will be the last word
				startIndex := i - 1

				// find the last operand
				// this will be the last word not preceded by OR
				lastIndex := len(words) - 1
				for ii := i + 2; ii < len(words); ii += 2 {
					if !strings.EqualFold(words[ii], or) {
						lastIndex = ii - 1
						break
					}
				}

				foundOr = true

				// combine the words into an any set
				var set []string
				for ii := startIndex; ii <= lastIndex; ii += 2 {
					word := extractPhrase(words[ii])
					if word == "" {
						continue
					}
					set = append(set, word)
				}

				searchSpec.AnySets = append(searchSpec.AnySets, set)

				// take out the OR'd words
				words = append(words[0:startIndex], words[lastIndex+1:]...)

				// break and reparse
				break
			}
		}
	}

	return words
}

func extractNotConditions(words []string, searchSpec *SearchSpecs) []string {
	var ret []string

	for _, w := range words {
		if len(w) > 1 && w[0] == notPrefix {
			word := extractPhrase(w[1:])
			if word == "" {
				continue
			}
			searchSpec.MustNot = append(searchSpec.MustNot, word)
		} else {
			ret = append(ret, w)
		}
	}

	return ret
}

func extractPhrase(w string) string {
	if len(w) > 1 && w[0] == phraseChar && w[len(w)-1] == phraseChar {
		return w[1 : len(w)-1]
	}

	return w
}

// ParseSearchString parses the Q value and returns a SearchSpecs object.
//
// By default, any words in the search value must appear in the results.
// Words encompassed by quotes (") as treated as a single term.
// Where keyword "OR" (case-insensitive) appears (and is not part of a quoted phrase), one of the
// OR'd terms must appear in the results.
// Where a keyword is prefixed with "-", that keyword must not appear in the results.
// Where OR appears as the first or last term, or where one of the OR operands has a
// not prefix, then the OR is treated literally.
func ParseSearchString(s string) SearchSpecs {
	s = strings.TrimSpace(s)

	if s == "" {
		return SearchSpecs{}
	}

	// break into words
	words := strings.Split(s, " ")

	// combine phrases first, then extract OR conditions, then extract NOT conditions
	// and the leftovers will be AND'd
	ret := SearchSpecs{}
	words = combinePhrases(words)
	words = extractOrConditions(words, &ret)
	words = extractNotConditions(words, &ret)

	for _, w := range words {
		// ignore empty quotes
		word := extractPhrase(w)
		if word == "" {
			continue
		}
		ret.MustHave = append(ret.MustHave, word)
	}

	return ret
}
