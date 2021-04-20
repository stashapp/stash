package autotag

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func getMatchingTags(path string, tagReader models.TagReader) ([]*models.Tag, error) {
	words := getPathWords(path)
	tags, err := tagReader.QueryForAutoTag(words)

	if err != nil {
		return nil, err
	}

	var ret []*models.Tag
	for _, p := range tags {
		if nameMatchesPath(p.Name, path) {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func getTagTagger(p *models.Tag) tagger {
	return tagger{
		ID:   p.ID,
		Type: "tag",
		Name: p.Name,
	}
}

func TagScenes(p *models.Tag, paths []string, rw models.SceneReaderWriter) error {
	t := getTagTagger(p)

	return t.tagScenes(paths, rw, func(subjectID, otherID int) (bool, error) {
		return scene.AddTag(rw, otherID, subjectID)
	})
}
