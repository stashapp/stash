package autotag

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func getMatchingPerformers(path string, performerReader models.PerformerReader) ([]*models.Performer, error) {
	performers, _, err := performerReader.Query(nil, nil)

	if err != nil {
		return nil, err
	}

	var ret []*models.Performer
	for _, p := range performers {
		if nameMatchesPath(p.Name.String, path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func getPerformerFilter(path string) models.PerformerFilterType {
	//words := getPathWords(path)

	// TODO - fix this with ors
	return models.PerformerFilterType{
		Aliases: &models.StringCriterionInput{
			Modifier: models.CriterionModifierEquals,
			Value:    path,
		},
	}
}

func getPerformerTagger(p *models.Performer) tagger {
	return tagger{
		ID:   p.ID,
		Type: "performer",
		Name: p.Name.String,
	}
}

func PerformerScenes(p *models.Performer, rw models.SceneReaderWriter) error {
	t := getPerformerTagger(p)

	return t.tagScenes(rw, func(subjectID, otherID int) (bool, error) {
		return scene.AddPerformer(rw, otherID, subjectID)
	})
}
