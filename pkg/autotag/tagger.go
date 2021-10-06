// Package autotag provides methods to auto-tag scenes with performers,
// studios and tags.
//
// The autotag engine tags scenes with performers/studios/tags if the scene's
// path matches the performer/studio/tag name. A scene's path is considered
// a match if it contains the performer/studio/tag's full name, ignoring any
// '.', '-', '_' characters in the path.
//
// For example, for a performer "foo bar", the following paths would be
// considered a match: "foo bar.mp4", "foobar.mp4", "foo.bar.mp4",
// "foo-bar.mp4", "aaa.foo bar.bbb.mp4".
// The following would not be considered a match:
// "aafoo bar.mp4", "foo barbb.mp4", "foo/bar.mp4"
package autotag

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const separatorChars = `.\-_ `

func getPathQueryRegex(name string) string {
	// escape specific regex characters
	name = regexp.QuoteMeta(name)

	// handle path separators
	const separator = `[` + separatorChars + `]`

	ret := strings.Replace(name, " ", separator+"*", -1)
	ret = `(?:^|_|[^\w\d])` + ret + `(?:$|_|[^\w\d])`
	return ret
}

func nameMatchesPath(name, path string) bool {
	// escape specific regex characters
	name = regexp.QuoteMeta(name)

	name = strings.ToLower(name)
	path = strings.ToLower(path)

	// handle path separators
	const separator = `[` + separatorChars + `]`

	reStr := strings.Replace(name, " ", separator+"*", -1)
	reStr = `(?:^|_|[^\w\d])` + reStr + `(?:$|_|[^\w\d])`

	re := regexp.MustCompile(reStr)
	return re.MatchString(path)
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

type tagger struct {
	ID   int
	Type string
	Name string
	Path string
}

type addLinkFunc func(subjectID, otherID int) (bool, error)

func (t *tagger) addError(otherType, otherName string, err error) error {
	return fmt.Errorf("error adding %s '%s' to %s '%s': %s", otherType, otherName, t.Type, t.Name, err.Error())
}

func (t *tagger) addLog(otherType, otherName string) {
	logger.Infof("Added %s '%s' to %s '%s'", otherType, otherName, t.Type, t.Name)
}

func (t *tagger) tagPerformers(performerReader models.PerformerReader, addFunc addLinkFunc) error {
	others, err := getMatchingPerformers(t.Path, performerReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("performer", p.Name.String, err)
		}

		if added {
			t.addLog("performer", p.Name.String)
		}
	}

	return nil
}

func (t *tagger) tagStudios(studioReader models.StudioReader, addFunc addLinkFunc) error {
	others, err := getMatchingStudios(t.Path, studioReader)
	if err != nil {
		return err
	}

	// only add first studio
	if len(others) > 0 {
		studio := others[0]
		added, err := addFunc(t.ID, studio.ID)

		if err != nil {
			return t.addError("studio", studio.Name.String, err)
		}

		if added {
			t.addLog("studio", studio.Name.String)
		}
	}

	return nil
}

func (t *tagger) tagTags(tagReader models.TagReader, addFunc addLinkFunc) error {
	others, err := getMatchingTags(t.Path, tagReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("tag", p.Name, err)
		}

		if added {
			t.addLog("tag", p.Name)
		}
	}

	return nil
}

func (t *tagger) tagScenes(paths []string, sceneReader models.SceneReader, addFunc addLinkFunc) error {
	others, err := getMatchingScenes(t.Name, paths, sceneReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("scene", p.GetTitle(), err)
		}

		if added {
			t.addLog("scene", p.GetTitle())
		}
	}

	return nil
}

func (t *tagger) tagImages(paths []string, imageReader models.ImageReader, addFunc addLinkFunc) error {
	others, err := getMatchingImages(t.Name, paths, imageReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("image", p.GetTitle(), err)
		}

		if added {
			t.addLog("image", p.GetTitle())
		}
	}

	return nil
}

func (t *tagger) tagGalleries(paths []string, galleryReader models.GalleryReader, addFunc addLinkFunc) error {
	others, err := getMatchingGalleries(t.Name, paths, galleryReader)
	if err != nil {
		return err
	}

	for _, p := range others {
		added, err := addFunc(t.ID, p.ID)

		if err != nil {
			return t.addError("gallery", p.GetTitle(), err)
		}

		if added {
			t.addLog("gallery", p.GetTitle())
		}
	}

	return nil
}
