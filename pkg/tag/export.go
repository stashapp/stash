package tag

import (
	"fmt"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// ToJSON converts a Tag object into its JSON equivalent.
func ToJSON(reader models.TagReader, tag *models.Tag) (*jsonschema.Tag, error) {
	newTagJSON := jsonschema.Tag{
		Name:      tag.Name,
		CreatedAt: models.JSONTime{Time: tag.CreatedAt.Timestamp},
		UpdatedAt: models.JSONTime{Time: tag.UpdatedAt.Timestamp},
	}

	image, err := reader.GetImage(tag.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting tag image: %s", err.Error())
	}

	if len(image) > 0 {
		newTagJSON.Image = utils.GetBase64StringFromData(image)
	}

	return &newTagJSON, nil
}

func GetIDs(tags []*models.Tag) []int {
	var results []int
	for _, tag := range tags {
		results = append(results, tag.ID)
	}

	return results
}

func GetNames(tags []*models.Tag) []string {
	var results []string
	for _, tag := range tags {
		results = append(results, tag.Name)
	}

	return results
}
