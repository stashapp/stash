package sqlite

import (
	"regexp"

	lru "github.com/hashicorp/golang-lru/v2"
)

// size of the regex LRU cache in elements.
// A small number number was chosen because it's most likely use is for a
// single query - this function gets called for every row in the (filtered)
// results. It's likely to only need no more than 1 or 2 in any given query.
// After that point, it's just sitting in the cache and is unlikely to be used
// again.
const regexCacheSize = 10

var regexCache *lru.Cache[string, *regexp.Regexp]

func init() {
	regexCache, _ = lru.New[string, *regexp.Regexp](regexCacheSize)
}

// regexFn is registered as an SQLite function as "regexp"
// It uses an LRU cache to cache recent regex patterns to reduce CPU load over
// identical patterns.
func regexFn(re, s string) (bool, error) {
	compiled, ok := regexCache.Get(re)
	if !ok {
		var err error
		compiled, err = regexp.Compile(re)
		if err != nil {
			return false, err
		}
		regexCache.Add(re, compiled)
	}

	return compiled.MatchString(s), nil
}
