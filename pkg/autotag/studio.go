package autotag

import (
	"database/sql"

	"github.com/stashapp/stash/pkg/models"
)

func getMatchingStudios(path string, reader models.StudioReader) ([]*models.Studio, error) {
	words := getPathWords(path)
	candidates, err := reader.QueryForAutoTag(words)

	if err != nil {
		return nil, err
	}

	var ret []*models.Studio
	for _, c := range candidates {
		if nameMatchesPath(c.Name.String, path) {
			ret = append(ret, c)
		}
	}

	return ret, nil
}

func addSceneStudio(sceneWriter models.SceneReaderWriter, sceneID, studioID int) (bool, error) {
	// don't set if already set
	scene, err := sceneWriter.Find(sceneID)
	if err != nil {
		return false, err
	}

	if scene.StudioID.Valid {
		return false, nil
	}

	// set the studio id
	s := sql.NullInt64{Int64: int64(studioID), Valid: true}
	scenePartial := models.ScenePartial{
		ID:       sceneID,
		StudioID: &s,
	}

	if _, err := sceneWriter.Update(scenePartial); err != nil {
		return false, err
	}
	return true, nil
}

func getStudioTagger(p *models.Studio) tagger {
	return tagger{
		ID:   p.ID,
		Type: "studio",
		Name: p.Name.String,
	}
}

func StudioScenes(p *models.Studio, paths []string, rw models.SceneReaderWriter) error {
	t := getStudioTagger(p)

	return t.tagScenes(paths, rw, func(subjectID, otherID int) (bool, error) {
		return addSceneStudio(rw, otherID, subjectID)
	})
}
