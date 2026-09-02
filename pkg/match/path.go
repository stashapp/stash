// Package match provides functions for matching paths to models.
package match

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil"
)

const (
	separatorChars   = `.\-_ `
	separatorPattern = `(?:_|[^\p{L}\w\d])+`

	reNotLetterWordUnicode = `[^\p{L}\w\d]`
	reNotLetterWord        = `[^\w\d]`
)

var separatorRE = regexp.MustCompile(separatorPattern)

func getPathQueryRegex(name string) string {
	// escape specific regex characters
	name = regexp.QuoteMeta(name)

	// handle path separators
	const separator = `[` + separatorChars + `]`

	ret := strings.ReplaceAll(name, " ", separator+"*")

	ret = `(?:^|_|[^\p{L}\d])` + ret + `(?:$|_|[^\p{L}\d])`
	return ret
}

func getPathWords(path string, trimExt bool) []string {
	retStr := path

	if trimExt {
		// remove the extension
		ext := filepath.Ext(retStr)
		if ext != "" {
			retStr = strings.TrimSuffix(retStr, ext)
		}
	}

	// handle path separators
	retStr = separatorRE.ReplaceAllString(retStr, " ")

	words := strings.Split(retStr, " ")

	// remove any single letter words
	var ret []string
	for _, w := range words {
		if utf8.RuneCountInString(w) > 1 {
			// #1450 - we need to open up the criteria for matching so that we
			// can match where path has no space between subject names -
			// ie name = "foo bar" - path = "foobar"
			// we post-match afterwards, so we can afford to be a little loose
			// with the query
			// just use the first two characters
			// #2293 - need to convert to unicode runes for the substring, otherwise
			// the resulting string is corrupted.
			ret = sliceutil.AppendUnique(ret, string([]rune(w)[0:2]))
		}
	}

	return ret
}

// https://stackoverflow.com/a/53069799
func allASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// nameMatchesPath returns the index in the path for the right-most match.
// Returns -1 if not found.
func nameMatchesPath(name, path string) int {
	// #2363 - optimisation: only use unicode character regexp if path contains
	// unicode characters
	re := nameToRegexp(name, !allASCII(path))
	return regexpMatchesPath(re, path)
}

// pathMatcher holds per-path precomputed values so they aren't recomputed
// for every candidate name. `allASCII` and `strings.ToLower(path)` were
// running once per (candidate, file) pair before; under a worker pool with
// thousands of candidates per file that was the dominant allocation.
type pathMatcher struct {
	loweredPath string
	useUnicode  bool
	cache       *Cache
}

func newPathMatcher(path string, cache *Cache) pathMatcher {
	return pathMatcher{
		loweredPath: strings.ToLower(path),
		useUnicode:  !allASCII(path),
		cache:       cache,
	}
}

// match returns the right-most index where name matches the path, or -1.
// Uses the cache's compiled-regexp table so each name is compiled once per
// autotag run instead of once per file.
func (m *pathMatcher) match(name string) int {
	re := m.cache.nameRegexp(name, m.useUnicode)
	found := re.FindAllStringIndex(m.loweredPath, -1)
	if found == nil {
		return -1
	}
	return found[len(found)-1][0]
}

// nameToRegexp compiles a regexp pattern to match paths from the given name.
// Set useUnicode to true if this regexp is to be used on any strings with unicode characters.
func nameToRegexp(name string, useUnicode bool) *regexp.Regexp {
	// escape specific regex characters
	name = regexp.QuoteMeta(name)

	name = strings.ToLower(name)

	// handle path separators
	const separator = `[` + separatorChars + `]`

	// performance optimisation: only use \p{L} is useUnicode is true
	notWord := reNotLetterWord
	if useUnicode {
		notWord = reNotLetterWordUnicode
	}

	reStr := strings.ReplaceAll(name, " ", separator+"*")
	reStr = `(?:^|_|` + notWord + `)` + reStr + `(?:$|_|` + notWord + `)`

	re := regexp.MustCompile(reStr)
	return re
}

func regexpMatchesPath(r *regexp.Regexp, path string) int {
	path = strings.ToLower(path)
	found := r.FindAllStringIndex(path, -1)
	if found == nil {
		return -1
	}
	return found[len(found)-1][0]
}

func getPerformers(ctx context.Context, words []string, performerReader models.PerformerAutoTagQueryer, cache *Cache) ([]*models.Performer, error) {
	performers, err := performerReader.QueryForAutoTag(ctx, words)
	if err != nil {
		return nil, err
	}

	swPerformers, err := getSingleLetterPerformers(ctx, cache, performerReader)
	if err != nil {
		return nil, err
	}

	return append(performers, swPerformers...), nil
}

// PathToPerformers returns performers whose name matches the given path.
//
// When the cache has been preloaded via Cache.PreloadPerformers, the full
// non-ignored performer set is already in memory and a 2-rune prefix index
// narrows candidates before regex-matching — this is the path the bulk
// file-based auto-tag job takes. Otherwise (e.g., the built-in scraper,
// which runs on a single scene per request) falls back to a per-call SQL
// prefilter via reader.QueryForAutoTag.
func PathToPerformers(ctx context.Context, path string, reader models.PerformerAutoTagQueryer, cache *Cache, trimExt bool) ([]*models.Performer, error) {
	var performers []*models.Performer
	if cache != nil && cache.allPerformers != nil {
		performers = cache.performerCandidates(getPathWords(path, trimExt))
	} else {
		words := getPathWords(path, trimExt)
		var err error
		performers, err = getPerformers(ctx, words, reader, cache)
		if err != nil {
			return nil, err
		}
	}

	pm := newPathMatcher(path, cache)
	var ret []*models.Performer
	for _, p := range performers {
		matches := false
		if pm.match(p.Name) != -1 {
			matches = true
		}

		// TODO - disabled alias matching until we can get finer
		// control over the matching. To re-enable:
		//   - uncomment this block (fallback path)
		//   - have Cache.PreloadPerformers load aliases (e.g. via
		//     loadAllAliases, as PreloadStudios/PreloadTags do) and
		//     iterate them here in the preloaded path too
		// if !matches {
		// 	if err := p.LoadAliases(ctx, reader); err != nil {
		// 		return nil, err
		// 	}
		// 	for _, alias := range p.Aliases.List() {
		// 		if pm.match(alias) != -1 {
		// 			matches = true
		// 			break
		// 		}
		// 	}
		// }

		if matches {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func getStudios(ctx context.Context, words []string, reader models.StudioAutoTagQueryer, cache *Cache) ([]*models.Studio, error) {
	studios, err := reader.QueryForAutoTag(ctx, words)
	if err != nil {
		return nil, err
	}

	swStudios, err := getSingleLetterStudios(ctx, cache, reader)
	if err != nil {
		return nil, err
	}

	return append(studios, swStudios...), nil
}

// PathToStudio returns the studio whose name or alias matches the given
// path. Where multiple match, the one matching the latest position wins.
//
// See PathToPerformers for the preloaded-vs-fallback behavior.
func PathToStudio(ctx context.Context, path string, reader models.StudioAutoTagQueryer, cache *Cache, trimExt bool) (*models.Studio, error) {
	pm := newPathMatcher(path, cache)

	if cache != nil && cache.allStudios != nil {
		candidates := cache.studioCandidates(getPathWords(path, trimExt))
		var ret *models.Studio
		index := -1
		for _, c := range candidates {
			if matchIndex := pm.match(c.Studio.Name); matchIndex != -1 && matchIndex > index {
				ret = c.Studio
				index = matchIndex
			}
			for _, alias := range c.Aliases {
				if matchIndex := pm.match(alias); matchIndex != -1 && matchIndex > index {
					ret = c.Studio
					index = matchIndex
				}
			}
		}
		return ret, nil
	}

	words := getPathWords(path, trimExt)
	candidates, err := getStudios(ctx, words, reader, cache)
	if err != nil {
		return nil, err
	}

	var ret *models.Studio
	index := -1
	for _, c := range candidates {
		if matchIndex := pm.match(c.Name); matchIndex != -1 && matchIndex > index {
			ret = c
			index = matchIndex
		}

		aliases, err := reader.GetAliases(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		for _, alias := range aliases {
			if matchIndex := pm.match(alias); matchIndex != -1 && matchIndex > index {
				ret = c
				index = matchIndex
			}
		}
	}

	return ret, nil
}

func getTags(ctx context.Context, words []string, reader models.TagAutoTagQueryer, cache *Cache) ([]*models.Tag, error) {
	tags, err := reader.QueryForAutoTag(ctx, words)
	if err != nil {
		return nil, err
	}

	swTags, err := getSingleLetterTags(ctx, cache, reader)
	if err != nil {
		return nil, err
	}

	return append(tags, swTags...), nil
}

// PathToTags returns tags whose name or alias matches the given path.
//
// See PathToPerformers for the preloaded-vs-fallback behavior.
func PathToTags(ctx context.Context, path string, reader models.TagAutoTagQueryer, cache *Cache, trimExt bool) ([]*models.Tag, error) {
	pm := newPathMatcher(path, cache)

	if cache != nil && cache.allTags != nil {
		candidates := cache.tagCandidates(getPathWords(path, trimExt))
		var ret []*models.Tag
		for _, c := range candidates {
			if pm.match(c.Tag.Name) != -1 {
				ret = append(ret, c.Tag)
				continue
			}
			for _, alias := range c.Aliases {
				if pm.match(alias) != -1 {
					ret = append(ret, c.Tag)
					break
				}
			}
		}
		return ret, nil
	}

	words := getPathWords(path, trimExt)
	tags, err := getTags(ctx, words, reader, cache)
	if err != nil {
		return nil, err
	}

	var ret []*models.Tag
	for _, t := range tags {
		matches := false
		if pm.match(t.Name) != -1 {
			matches = true
		}
		if !matches {
			aliases, err := reader.GetAliases(ctx, t.ID)
			if err != nil {
				return nil, err
			}
			for _, alias := range aliases {
				if pm.match(alias) != -1 {
					matches = true
					break
				}
			}
		}
		if matches {
			ret = append(ret, t)
		}
	}

	return ret, nil
}

func PathToScenesFn(ctx context.Context, name string, paths []string, sceneReader models.SceneQueryer, fn func(ctx context.Context, scene *models.Scene) error) error {
	regex := getPathQueryRegex(name)
	organized := false
	filter := models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Value:    "(?i)" + regex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
		Organized: &organized,
	}

	filter.And = scene.PathsFilter(paths)

	// do in batches
	pp := 1000
	sort := "id"
	sortDir := models.SortDirectionEnumAsc
	lastID := 0

	for {
		if lastID != 0 {
			filter.ID = &models.IntCriterionInput{
				Value:    lastID,
				Modifier: models.CriterionModifierGreaterThan,
			}
		}

		scenes, err := scene.Query(ctx, sceneReader, &filter, &models.FindFilterType{
			PerPage:   &pp,
			Sort:      &sort,
			Direction: &sortDir,
		})

		if err != nil {
			return fmt.Errorf("error querying scenes with regex '%s': %s", regex, err.Error())
		}

		// paths may have unicode characters
		const useUnicode = true

		r := nameToRegexp(name, useUnicode)
		for _, p := range scenes {
			if regexpMatchesPath(r, p.Path) != -1 {
				if err := fn(ctx, p); err != nil {
					return fmt.Errorf("processing scene %s: %w", p.GetTitle(), err)
				}
			}
		}

		if len(scenes) < pp {
			break
		}

		lastID = scenes[len(scenes)-1].ID
	}

	return nil
}

func PathToImagesFn(ctx context.Context, name string, paths []string, imageReader models.ImageQueryer, fn func(ctx context.Context, scene *models.Image) error) error {
	regex := getPathQueryRegex(name)
	organized := false
	filter := models.ImageFilterType{
		Path: &models.StringCriterionInput{
			Value:    "(?i)" + regex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
		Organized: &organized,
	}

	filter.And = image.PathsFilter(paths)

	// do in batches
	pp := 1000
	sort := "id"
	sortDir := models.SortDirectionEnumAsc
	lastID := 0

	for {
		if lastID != 0 {
			filter.ID = &models.IntCriterionInput{
				Value:    lastID,
				Modifier: models.CriterionModifierGreaterThan,
			}
		}

		images, err := image.Query(ctx, imageReader, &filter, &models.FindFilterType{
			PerPage:   &pp,
			Sort:      &sort,
			Direction: &sortDir,
		})

		if err != nil {
			return fmt.Errorf("error querying images with regex '%s': %s", regex, err.Error())
		}

		// paths may have unicode characters
		const useUnicode = true

		r := nameToRegexp(name, useUnicode)
		for _, p := range images {
			if regexpMatchesPath(r, p.Path) != -1 {
				if err := fn(ctx, p); err != nil {
					return fmt.Errorf("processing image %s: %w", p.GetTitle(), err)
				}
			}
		}

		if len(images) < pp {
			break
		}

		lastID = images[len(images)-1].ID
	}

	return nil
}

func PathToGalleriesFn(ctx context.Context, name string, paths []string, galleryReader models.GalleryQueryer, fn func(ctx context.Context, scene *models.Gallery) error) error {
	regex := getPathQueryRegex(name)
	organized := false
	filter := models.GalleryFilterType{
		Path: &models.StringCriterionInput{
			Value:    "(?i)" + regex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
		Organized: &organized,
	}

	filter.And = gallery.PathsFilter(paths)

	// do in batches
	pp := 1000
	sort := "id"
	sortDir := models.SortDirectionEnumAsc
	lastID := 0

	for {
		if lastID != 0 {
			filter.ID = &models.IntCriterionInput{
				Value:    lastID,
				Modifier: models.CriterionModifierGreaterThan,
			}
		}

		galleries, _, err := galleryReader.Query(ctx, &filter, &models.FindFilterType{
			PerPage:   &pp,
			Sort:      &sort,
			Direction: &sortDir,
		})

		if err != nil {
			return fmt.Errorf("error querying galleries with regex '%s': %s", regex, err.Error())
		}

		// paths may have unicode characters
		const useUnicode = true

		r := nameToRegexp(name, useUnicode)
		for _, p := range galleries {
			path := p.Path
			if path != "" && regexpMatchesPath(r, path) != -1 {
				if err := fn(ctx, p); err != nil {
					return fmt.Errorf("processing gallery %s: %w", p.GetTitle(), err)
				}
			}
		}

		if len(galleries) < pp {
			break
		}

		lastID = galleries[len(galleries)-1].ID
	}

	return nil
}
