package match

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

const separatorChars = `.\-_ `

func getPathQueryRegex(name string) string {
	// escape specific regex characters
	name = regexp.QuoteMeta(name)

	// handle path separators
	const separator = `[` + separatorChars + `]`

	ret := strings.ReplaceAll(name, " ", separator+"*")
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
	const separator = `(?:_|[^\w\d])+`
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
			ret = append(ret, w[0:2])
		}
	}

	return ret
}

func nameMatchesPath(name, path string) bool {
	// escape specific regex characters
	name = regexp.QuoteMeta(name)

	name = strings.ToLower(name)
	path = strings.ToLower(path)

	// handle path separators
	const separator = `[` + separatorChars + `]`

	reStr := strings.ReplaceAll(name, " ", separator+"*")
	reStr = `(?:^|_|[^\w\d])` + reStr + `(?:$|_|[^\w\d])`

	re := regexp.MustCompile(reStr)
	return re.MatchString(path)
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
		if nameMatchesPath(p.Name.String, path) { // || nameMatchesPath(p.Aliases.String, path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func PathToStudios(path string, reader models.StudioReader) ([]*models.Studio, error) {
	words := getPathWords(path)
	candidates, err := reader.QueryForAutoTag(words)

	if err != nil {
		return nil, err
	}

	var ret []*models.Studio
	for _, c := range candidates {
		matches := false
		if nameMatchesPath(c.Name.String, path) {
			matches = true
		}

		if !matches {
			aliases, err := reader.GetAliases(c.ID)
			if err != nil {
				return nil, err
			}

			for _, alias := range aliases {
				if nameMatchesPath(alias, path) {
					matches = true
					break
				}
			}
		}

		if matches {
			ret = append(ret, c)
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
		if nameMatchesPath(t.Name, path) {
			matches = true
		}

		if !matches {
			aliases, err := tagReader.GetAliases(t.ID)
			if err != nil {
				return nil, err
			}
			for _, alias := range aliases {
				if nameMatchesPath(alias, path) {
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

func scenePathsFilter(paths []string) *models.SceneFilterType {
	if paths == nil {
		return nil
	}

	sep := string(filepath.Separator)

	var ret *models.SceneFilterType
	var or *models.SceneFilterType
	for _, p := range paths {
		newOr := &models.SceneFilterType{}
		if or != nil {
			or.Or = newOr
		} else {
			ret = newOr
		}

		or = newOr

		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
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

	filter.And = scenePathsFilter(paths)

	pp := models.PerPageAll
	scenes, _, err := sceneReader.Query(&filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying scenes with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Scene
	for _, p := range scenes {
		if nameMatchesPath(name, p.Path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func imagePathsFilter(paths []string) *models.ImageFilterType {
	if paths == nil {
		return nil
	}

	sep := string(filepath.Separator)

	var ret *models.ImageFilterType
	var or *models.ImageFilterType
	for _, p := range paths {
		newOr := &models.ImageFilterType{}
		if or != nil {
			or.Or = newOr
		} else {
			ret = newOr
		}

		or = newOr

		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
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

	filter.And = imagePathsFilter(paths)

	pp := models.PerPageAll
	images, _, err := imageReader.Query(&filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying images with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Image
	for _, p := range images {
		if nameMatchesPath(name, p.Path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func galleryPathsFilter(paths []string) *models.GalleryFilterType {
	if paths == nil {
		return nil
	}

	sep := string(filepath.Separator)

	var ret *models.GalleryFilterType
	var or *models.GalleryFilterType
	for _, p := range paths {
		newOr := &models.GalleryFilterType{}
		if or != nil {
			or.Or = newOr
		} else {
			ret = newOr
		}

		or = newOr

		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
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

	filter.And = galleryPathsFilter(paths)

	pp := models.PerPageAll
	gallerys, _, err := galleryReader.Query(&filter, &models.FindFilterType{
		PerPage: &pp,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying gallerys with regex '%s': %s", regex, err.Error())
	}

	var ret []*models.Gallery
	for _, p := range gallerys {
		if nameMatchesPath(name, p.Path.String) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
