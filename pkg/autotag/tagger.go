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
	// handle path separators
	const separator = `[` + separatorChars + `]`

	ret := strings.Replace(name, " ", separator+"*", -1)
	ret = `(?:^|_|[^\w\d])` + ret + `(?:$|_|[^\w\d])`
	return ret
}

func nameMatchesPath(name, path string) bool {
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
			ret = append(ret, w)
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
