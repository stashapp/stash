package database

import (
	"regexp"

	lru "github.com/hashicorp/golang-lru"
)

const regexCacheSize = 10

var regexCache *lru.Cache

func init() {
	regexCache, _ = lru.New(regexCacheSize)
}

// regexFn is registered as an SQLite function as "regexp"
// It uses an LRU cache to cache recent regex patterns to reduce CPU load over
// identical patterns.
func regexFn(re, s string) (bool, error) {
	entry, ok := regexCache.Get(re)
	var compiled *regexp.Regexp

	if !ok {
		var err error
		compiled, err = regexp.Compile(re)
		if err != nil {
			return false, err
		}
		regexCache.Add(re, compiled)
	} else {
		compiled = entry.(*regexp.Regexp)
	}

	return compiled.MatchString(s), nil
}
