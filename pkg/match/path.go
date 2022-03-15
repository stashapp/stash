package match

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/utils"
)

const (
	separatorChars = `.\-_ `

	reNotLetterWordUnicode = `[^\p{L}\w\d]`
	reNotLetterWord        = `[^\w\d]`
)

func getPathQueryRegex(name string) string {
	// escape specific regex characters
	name = regexp.QuoteMeta(name)

	// handle path separators
	const separator = `[` + separatorChars + `]`

	ret := strings.ReplaceAll(name, " ", separator+"*")

	// \p{L} is specifically omitted here because of the performance hit when
	// including it. It does mean that paths where the name is bounded by
	// unicode letters will be returned. However, the results should be tested
	// by nameMatchesPath which does include \p{L}. The improvement in query
	// performance should be outweigh the performance hit of testing any extra
	// results.
	ret = `(?:^|_|[^\w\d])` + ret + `(?:$|_|[^\w\d])`
	return ret
}

func getPathWords(path string) []string {
	retStr := path

	// remove the extension
	ext := filepath.Ext(retStr)
	if ext != "" {
		retStr = strings.TrimSuffix(retStr, ext)
	}

	// handle path separators
	const separator = `(?:_|[^\p{L}\w\d])+`
	re := regexp.MustCompile(separator)
	retStr = re.ReplaceAllString(retStr, " ")

	words := strings.Split(retStr, " ")

	// remove any single letter words
	var ret []string
	for _, w := range words {
		if len(w) > 1 {
			// #1450 - we need to open up the criteria for matching so that we
			// can match where path has no space between subject names -
			// ie name = "foo bar" - path = "foobar"
			// we post-match afterwards, so we can afford to be a little loose
			// with the query
			// just use the first two characters
			// #2293 - need to convert to unicode runes for the substring, otherwise
			// the resulting string is corrupted.
			ret = utils.StrAppendUnique(ret, string([]rune(w)[0:2]))
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

func PathToPerformers(path string, performerReader models.PerformerReader) ([]*models.Performer, error) {
	words := getPathWords(path)
	performers, err := performerReader.QueryForAutoTag(words)

	if err != nil {
		return nil, err
	}

	var ret []*models.Performer
	for _, p := range performers {
		// TODO - commenting out alias handling until both sides work correctly
		if nameMatchesPath(p.Name.String, path) != -1 { // || nameMatchesPath(p.Aliases.String, path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

// PathToStudio returns the Studio that matches the given path.
// Where multiple matching studios are found, the one that matches the latest
// position in the path is returned.
func PathToStudio(path string, reader models.StudioReader) (*models.Studio, error) {
	words := getPathWords(path)
	candidates, err := reader.QueryForAutoTag(words)

	if err != nil {
		return nil, err
	}

	var ret *models.Studio
	index := -1
	for _, c := range candidates {
		matchIndex := nameMatchesPath(c.Name.String, path)
		if matchIndex != -1 && matchIndex > index {
			ret = c
			index = matchIndex
		}

		aliases, err := reader.GetAliases(c.ID)
		if err != nil {
			return nil, err
		}

		for _, alias := range aliases {
			matchIndex = nameMatchesPath(alias, path)
			if matchIndex != -1 && matchIndex > index {
				ret = c
				index = matchIndex
			}
		}
	}

	return ret, nil
}

func PathToTags(path string, tagReader models.TagReader) ([]*models.Tag, error) {
	words := getPathWords(path)
	tags, err := tagReader.QueryForAutoTag(words)

	if err != nil {
		return nil, err
	}

	var ret []*models.Tag
	for _, t := range tags {
		matches := false
		if nameMatchesPath(t.Name, path) != -1 {
			matches = true
		}

		if !matches {
			aliases, err := tagReader.GetAliases(t.ID)
			if err != nil {
				return nil, err
			}
			for _, alias := range aliases {
				if nameMatchesPath(alias, path) != -1 {
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

func PathToScenes(name string, paths []string, sceneReader models.SceneReader) ([]*models.Scene, error) {
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

	pp := models.PerPageAll
	scenes, err := scene.Query(sceneReader, &filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying scenes with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Scene

	// paths may have unicode characters
	const useUnicode = true

	r := nameToRegexp(name, useUnicode)
	for _, p := range scenes {
		if regexpMatchesPath(r, p.Path) != -1 {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func PathToImages(name string, paths []string, imageReader models.ImageReader) ([]*models.Image, error) {
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

	pp := models.PerPageAll
	images, err := image.Query(imageReader, &filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying images with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Image

	// paths may have unicode characters
	const useUnicode = true

	r := nameToRegexp(name, useUnicode)
	for _, p := range images {
		if regexpMatchesPath(r, p.Path) != -1 {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func PathToGalleries(name string, paths []string, galleryReader models.GalleryReader) ([]*models.Gallery, error) {
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

	pp := models.PerPageAll
	gallerys, _, err := galleryReader.Query(&filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying gallerys with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Gallery

	// paths may have unicode characters
	const useUnicode = true

	r := nameToRegexp(name, useUnicode)
	for _, p := range gallerys {
		if regexpMatchesPath(r, p.Path.String) != -1 {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
