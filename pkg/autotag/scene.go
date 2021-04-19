package autotag

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func getMatchingScenes(name string, sceneReader models.SceneReader) ([]*models.Scene, error) {
	regex := getPathQueryRegex(name)
	organized := false
	filter := models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Value:    "(?i)" + regex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
		Organized: &organized,
	}

	pp := 0
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

func getSceneFileTagger(s *models.Scene) tagger {
	return tagger{
		ID:   s.ID,
		Type: "scene",
		Name: s.GetTitle(),
		Path: s.Path,
	}
}

func ScenePerformers(s *models.Scene, rw models.SceneReaderWriter, performerReader models.PerformerReader) error {
	t := getSceneFileTagger(s)

	return t.tagPerformers(performerReader, func(subjectID, otherID int) (bool, error) {
		return scene.AddPerformer(rw, subjectID, otherID)
	})
}

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

func SceneTags(s *models.Scene, rw models.SceneReaderWriter, tagReader models.TagReader) error {
	t := getSceneFileTagger(s)

	return t.tagTags(tagReader, func(subjectID, otherID int) (bool, error) {
		return scene.AddTag(rw, subjectID, otherID)
	})
}
