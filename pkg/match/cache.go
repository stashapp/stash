package match

import (
	"context"
	"regexp"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stashapp/stash/pkg/models"
)

// regexpCacheSize bounds the compiled-regexp LRU. Sized generously so that
// for realistic libraries (up to ~100 k performers/studios/tags combined,
// each optionally with a unicode and ASCII variant) the cache never evicts
// during one auto-tag job. LRU is used for consistency with
// pkg/sqlite/regex.go; eviction only kicks in for libraries far past that.
const regexpCacheSize = 200_000

const singleFirstCharacterRegex = `^[\p{L}][.\-_ ]`

var singleFirstCharacterRE = regexp.MustCompile(singleFirstCharacterRegex)

// firstTwoRunesLower returns the first two runes of s, lowercased. Returns
// "" if s has fewer than two runes. Mirrors what getPathWords produces for
// path words, so the two can be compared as index keys.
func firstTwoRunesLower(s string) string {
	lower := strings.ToLower(s)
	runes := []rune(lower)
	if len(runes) < 2 {
		return ""
	}
	return string(runes[0:2])
}

// performerCandidates returns the set of preloaded performers that should
// be regex-checked for the given path words. Mirrors the SQL
// `name LIKE 'xx%' OR name LIKE 'yy%'` prefilter, plus always-check
// performers whose name begins with a single-letter word (which the 2-rune
// prefix lookup can't reach).
func (c *Cache) performerCandidates(pathWords []string) []*models.Performer {
	if len(c.performerByPrefix) == 0 && len(c.performerAlwaysCheck) == 0 {
		return nil
	}
	seen := make(map[int]bool, len(pathWords)*2)
	out := make([]*models.Performer, 0, len(pathWords)*2)
	for _, w := range pathWords {
		key := strings.ToLower(w)
		for _, p := range c.performerByPrefix[key] {
			if !seen[p.ID] {
				seen[p.ID] = true
				out = append(out, p)
			}
		}
	}
	for _, p := range c.performerAlwaysCheck {
		if !seen[p.ID] {
			seen[p.ID] = true
			out = append(out, p)
		}
	}
	return out
}

func (c *Cache) studioCandidates(pathWords []string) []cachedStudio {
	if len(c.studioByPrefix) == 0 && len(c.studioAlwaysCheck) == 0 {
		return nil
	}
	seen := make(map[int]bool, len(pathWords)*2)
	out := make([]cachedStudio, 0, len(pathWords)*2)
	for _, w := range pathWords {
		key := strings.ToLower(w)
		for _, s := range c.studioByPrefix[key] {
			if !seen[s.Studio.ID] {
				seen[s.Studio.ID] = true
				out = append(out, s)
			}
		}
	}
	for _, s := range c.studioAlwaysCheck {
		if !seen[s.Studio.ID] {
			seen[s.Studio.ID] = true
			out = append(out, s)
		}
	}
	return out
}

func (c *Cache) tagCandidates(pathWords []string) []cachedTag {
	if len(c.tagByPrefix) == 0 && len(c.tagAlwaysCheck) == 0 {
		return nil
	}
	seen := make(map[int]bool, len(pathWords)*2)
	out := make([]cachedTag, 0, len(pathWords)*2)
	for _, w := range pathWords {
		key := strings.ToLower(w)
		for _, t := range c.tagByPrefix[key] {
			if !seen[t.Tag.ID] {
				seen[t.Tag.ID] = true
				out = append(out, t)
			}
		}
	}
	for _, t := range c.tagAlwaysCheck {
		if !seen[t.Tag.ID] {
			seen[t.Tag.ID] = true
			out = append(out, t)
		}
	}
	return out
}

// Cache is used to cache queries that should not change across an autotag
// process. Safe for concurrent use by multiple goroutines.
type Cache struct {
	performersOnce sync.Once
	performersErr  error
	studiosOnce    sync.Once
	studiosErr     error
	tagsOnce       sync.Once
	tagsErr        error

	singleCharPerformers []*models.Performer
	singleCharStudios    []*models.Studio
	singleCharTags       []*models.Tag

	// Preloaded candidate sets. When populated (via PreloadX), the
	// PathTo* functions skip the per-path QueryForAutoTag DB roundtrip
	// and consult the in-memory prefix index instead. Nil means
	// "not preloaded, fall back to the old SQL-prefilter path".
	allPerformers []*models.Performer
	allStudios    []cachedStudio
	allTags       []cachedTag

	// Prefix indexes built at preload time. Map key is the first two
	// lowercased runes of name (or alias, for studios/tags). The
	// alwaysCandidate slice holds entries whose first "word" is a
	// single letter — they wouldn't be reached by 2-rune path word
	// lookup, so they must always be checked (mirroring the existing
	// single-letter regex query).
	performerByPrefix    map[string][]*models.Performer
	performerAlwaysCheck []*models.Performer
	studioByPrefix       map[string][]cachedStudio
	studioAlwaysCheck    []cachedStudio
	tagByPrefix          map[string][]cachedTag
	tagAlwaysCheck       []cachedTag

	regexpCacheOnce sync.Once
	regexpCache     *lru.Cache[regexpCacheKey, *regexp.Regexp]
}

// cachedStudio bundles a studio with its aliases so PathToStudio can match
// against both without an N+1 GetAliases query.
type cachedStudio struct {
	Studio  *models.Studio
	Aliases []string
}

// cachedTag bundles a tag with its aliases so PathToTags can match against
// both without an N+1 GetAliases query.
type cachedTag struct {
	Tag     *models.Tag
	Aliases []string
}

// PreloadPerformers loads all non-ignored performers into the cache and
// builds a 2-rune prefix index so subsequent PathToPerformers calls can
// skip both the per-path QueryForAutoTag and the per-candidate regex
// when no prefix matches.
func (c *Cache) PreloadPerformers(ctx context.Context, reader models.PerformerAutoTagQueryer) error {
	if c.allPerformers != nil {
		return nil
	}
	ignoreAutoTag := false
	perPage := -1
	perfs, _, err := reader.Query(ctx, &models.PerformerFilterType{
		IgnoreAutoTag: &ignoreAutoTag,
	}, &models.FindFilterType{PerPage: &perPage})
	if err != nil {
		return err
	}
	if perfs == nil {
		perfs = []*models.Performer{}
	}
	c.allPerformers = perfs

	c.performerByPrefix = make(map[string][]*models.Performer, len(perfs))
	for _, p := range perfs {
		if prefix := firstTwoRunesLower(p.Name); prefix != "" {
			c.performerByPrefix[prefix] = append(c.performerByPrefix[prefix], p)
		}
		if singleFirstCharacterRE.MatchString(p.Name) {
			c.performerAlwaysCheck = append(c.performerAlwaysCheck, p)
		}
	}
	return nil
}

// loadAllAliases loads aliases for the given ids. Uses the reader's bulk
// GetAllAliases method when available (avoiding the N+1 per-id roundtrip);
// otherwise falls back to per-id GetAliases.
func loadAllAliases(ctx context.Context, reader models.AliasLoader, ids []int) (map[int][]string, error) {
	if bulk, ok := reader.(models.AllAliasLoader); ok {
		return bulk.GetAllAliases(ctx)
	}
	ret := make(map[int][]string, len(ids))
	for _, id := range ids {
		a, err := reader.GetAliases(ctx, id)
		if err != nil {
			return nil, err
		}
		if len(a) > 0 {
			ret[id] = a
		}
	}
	return ret, nil
}

// PreloadStudios loads all non-ignored studios plus their aliases into the
// cache and builds a 2-rune prefix index (over names AND aliases, mirroring
// the SQL LEFT JOIN on studio_aliases).
func (c *Cache) PreloadStudios(ctx context.Context, reader models.StudioAutoTagQueryer) error {
	if c.allStudios != nil {
		return nil
	}
	ignoreAutoTag := false
	perPage := -1
	studios, _, err := reader.Query(ctx, &models.StudioFilterType{
		IgnoreAutoTag: &ignoreAutoTag,
	}, &models.FindFilterType{PerPage: &perPage})
	if err != nil {
		return err
	}
	ids := make([]int, len(studios))
	for i, s := range studios {
		ids[i] = s.ID
	}
	aliasesByID, err := loadAllAliases(ctx, reader, ids)
	if err != nil {
		return err
	}
	out := make([]cachedStudio, len(studios))
	c.studioByPrefix = make(map[string][]cachedStudio, len(studios))
	seenPerPrefix := make(map[string]map[int]bool)
	for i, s := range studios {
		aliases := aliasesByID[s.ID]
		cs := cachedStudio{Studio: s, Aliases: aliases}
		out[i] = cs

		c.indexByPrefix(s.ID, s.Name, aliases, seenPerPrefix, func(prefix string) {
			c.studioByPrefix[prefix] = append(c.studioByPrefix[prefix], cs)
		})
		if hasSingleFirstChar(s.Name, aliases) {
			c.studioAlwaysCheck = append(c.studioAlwaysCheck, cs)
		}
	}
	c.allStudios = out
	return nil
}

// PreloadTags loads all non-ignored tags plus their aliases into the cache
// and builds a 2-rune prefix index (over names AND aliases).
func (c *Cache) PreloadTags(ctx context.Context, reader models.TagAutoTagQueryer) error {
	if c.allTags != nil {
		return nil
	}
	ignoreAutoTag := false
	perPage := -1
	tags, _, err := reader.Query(ctx, &models.TagFilterType{
		IgnoreAutoTag: &ignoreAutoTag,
	}, &models.FindFilterType{PerPage: &perPage})
	if err != nil {
		return err
	}
	ids := make([]int, len(tags))
	for i, t := range tags {
		ids[i] = t.ID
	}
	aliasesByID, err := loadAllAliases(ctx, reader, ids)
	if err != nil {
		return err
	}
	out := make([]cachedTag, len(tags))
	c.tagByPrefix = make(map[string][]cachedTag, len(tags))
	seenPerPrefix := make(map[string]map[int]bool)
	for i, t := range tags {
		aliases := aliasesByID[t.ID]
		ct := cachedTag{Tag: t, Aliases: aliases}
		out[i] = ct

		c.indexByPrefix(t.ID, t.Name, aliases, seenPerPrefix, func(prefix string) {
			c.tagByPrefix[prefix] = append(c.tagByPrefix[prefix], ct)
		})
		if hasSingleFirstChar(t.Name, aliases) {
			c.tagAlwaysCheck = append(c.tagAlwaysCheck, ct)
		}
	}
	c.allTags = out
	return nil
}

// indexByPrefix records the entity under every distinct 2-rune prefix of
// its name/aliases (deduping so a name+alias that share a prefix bucket
// only add the entity once).
func (c *Cache) indexByPrefix(id int, name string, aliases []string, seen map[string]map[int]bool, add func(prefix string)) {
	emit := func(s string) {
		prefix := firstTwoRunesLower(s)
		if prefix == "" {
			return
		}
		if seen[prefix] == nil {
			seen[prefix] = make(map[int]bool)
		}
		if !seen[prefix][id] {
			seen[prefix][id] = true
			add(prefix)
		}
	}
	emit(name)
	for _, a := range aliases {
		emit(a)
	}
}

func hasSingleFirstChar(name string, aliases []string) bool {
	if singleFirstCharacterRE.MatchString(name) {
		return true
	}
	for _, a := range aliases {
		if singleFirstCharacterRE.MatchString(a) {
			return true
		}
	}
	return false
}

type regexpCacheKey struct {
	name       string
	useUnicode bool
}

// nameRegexp returns a compiled regexp for the given name, caching the
// result so repeated autotag calls across many files don't pay the
// compile cost each time.
func (c *Cache) nameRegexp(name string, useUnicode bool) *regexp.Regexp {
	if c == nil {
		return nameToRegexp(name, useUnicode)
	}

	c.regexpCacheOnce.Do(func() {
		c.regexpCache, _ = lru.New[regexpCacheKey, *regexp.Regexp](regexpCacheSize)
	})

	key := regexpCacheKey{name: name, useUnicode: useUnicode}
	if r, ok := c.regexpCache.Get(key); ok {
		return r
	}
	r := nameToRegexp(name, useUnicode)
	c.regexpCache.Add(key, r)
	return r
}

// getSingleLetterPerformers returns all performers with names that start with single character words.
// The autotag query splits the words into two-character words to query
// against. This means that performers with single-letter words in their names could potentially
// be missed.
// This query is expensive, so it's queried once and cached, if the cache if provided.
func getSingleLetterPerformers(ctx context.Context, c *Cache, reader models.PerformerAutoTagQueryer) ([]*models.Performer, error) {
	if c == nil {
		c = &Cache{}
	}

	c.performersOnce.Do(func() {
		pp := -1
		performers, _, err := reader.Query(ctx, &models.PerformerFilterType{
			Name: &models.StringCriterionInput{
				Value:    singleFirstCharacterRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}, &models.FindFilterType{
			PerPage: &pp,
		})

		if err != nil {
			c.performersErr = err
			return
		}

		if len(performers) == 0 {
			c.singleCharPerformers = make([]*models.Performer, 0)
		} else {
			c.singleCharPerformers = performers
		}
	})

	return c.singleCharPerformers, c.performersErr
}

// getSingleLetterStudios returns all studios with names that start with single character words.
// See getSingleLetterPerformers for details.
func getSingleLetterStudios(ctx context.Context, c *Cache, reader models.StudioAutoTagQueryer) ([]*models.Studio, error) {
	if c == nil {
		c = &Cache{}
	}

	c.studiosOnce.Do(func() {
		pp := -1
		studios, _, err := reader.Query(ctx, &models.StudioFilterType{
			Name: &models.StringCriterionInput{
				Value:    singleFirstCharacterRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
		}, &models.FindFilterType{
			PerPage: &pp,
		})

		if err != nil {
			c.studiosErr = err
			return
		}

		if len(studios) == 0 {
			c.singleCharStudios = make([]*models.Studio, 0)
		} else {
			c.singleCharStudios = studios
		}
	})

	return c.singleCharStudios, c.studiosErr
}

// getSingleLetterTags returns all tags with names that start with single character words.
// See getSingleLetterPerformers for details.
func getSingleLetterTags(ctx context.Context, c *Cache, reader models.TagAutoTagQueryer) ([]*models.Tag, error) {
	if c == nil {
		c = &Cache{}
	}

	c.tagsOnce.Do(func() {
		pp := -1
		tags, _, err := reader.Query(ctx, &models.TagFilterType{
			Name: &models.StringCriterionInput{
				Value:    singleFirstCharacterRegex,
				Modifier: models.CriterionModifierMatchesRegex,
			},
			OperatorFilter: models.OperatorFilter[models.TagFilterType]{
				Or: &models.TagFilterType{
					Aliases: &models.StringCriterionInput{
						Value:    singleFirstCharacterRegex,
						Modifier: models.CriterionModifierMatchesRegex,
					},
				},
			},
		}, &models.FindFilterType{
			PerPage: &pp,
		})

		if err != nil {
			c.tagsErr = err
			return
		}

		if len(tags) == 0 {
			c.singleCharTags = make([]*models.Tag, 0)
		} else {
			c.singleCharTags = tags
		}
	})

	return c.singleCharTags, c.tagsErr
}
