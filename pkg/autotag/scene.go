package autotag

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

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
			p = p + sep
		}

		or.Path = &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
}

func getMatchingScenes(name string, paths []string, sceneReader models.SceneReader) ([]*models.Scene, error) {
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
	scenes, _, _, _, err := sceneReader.Query(&filter, &models.FindFilterType{
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

func getSceneFileTagger(s *models.Scene) tagger {
	return tagger{
		ID:   s.ID,
		Type: "scene",
		Name: s.GetTitle(),
		Path: s.Path,
	}
}

// ScenePerformers tags the provided scene with performers whose name matches the scene's path.
func ScenePerformers(s *models.Scene, rw models.SceneReaderWriter, performerReader models.PerformerReader) error {
	t := getSceneFileTagger(s)

	return t.tagPerformers(performerReader, func(subjectID, otherID int) (bool, error) {
		return scene.AddPerformer(rw, subjectID, otherID)
	})
}

// SceneStudios tags the provided scene with the first studio whose name matches the scene's path.
//
// Scenes will not be tagged if studio is already set.
func SceneStudios(s *models.Scene, rw models.SceneReaderWriter, studioReader models.StudioReader) error {
	if s.StudioID.Valid {
		// don't modify
		return nil
	}

	t := getSceneFileTagger(s)

	return t.tagStudios(studioReader, func(subjectID, otherID int) (bool, error) {
		return addSceneStudio(rw, subjectID, otherID)
	})
}

// SceneTags tags the provided scene with tags whose name matches the scene's path.
func SceneTags(s *models.Scene, rw models.SceneReaderWriter, tagReader models.TagReader) error {
	t := getSceneFileTagger(s)

	return t.tagTags(tagReader, func(subjectID, otherID int) (bool, error) {
		return scene.AddTag(rw, subjectID, otherID)
	})
}
