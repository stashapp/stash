package tag

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type FinderAliasImageGetter interface {
	GetAliases(ctx context.Context, studioID int) ([]string, error)
	GetImage(ctx context.Context, tagID int) ([]byte, error)
	FindByChildTagID(ctx context.Context, childID int) ([]*models.Tag, error)
	models.StashIDLoader
}

// ToJSON converts a Tag object into its JSON equivalent.
func ToJSON(ctx context.Context, reader FinderAliasImageGetter, tag *models.Tag) (*jsonschema.Tag, error) {
	newTagJSON := jsonschema.Tag{
		Name:          tag.Name,
		SortName:      tag.SortName,
		Description:   tag.Description,
		Favorite:      tag.Favorite,
		IgnoreAutoTag: tag.IgnoreAutoTag,
		CreatedAt:     json.JSONTime{Time: tag.CreatedAt},
		UpdatedAt:     json.JSONTime{Time: tag.UpdatedAt},
	}

	aliases, err := reader.GetAliases(ctx, tag.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting tag aliases: %v", err)
	}

	newTagJSON.Aliases = aliases

	if err := tag.LoadStashIDs(ctx, reader); err != nil {
		return nil, fmt.Errorf("loading tag stash ids: %w", err)
	}

	stashIDs := tag.StashIDs.List()
	if len(stashIDs) > 0 {
		newTagJSON.StashIDs = stashIDs
	}

	image, err := reader.GetImage(ctx, tag.ID)
	if err != nil {
		logger.Errorf("Error getting tag image: %v", err)
	}

	if len(image) > 0 {
		newTagJSON.Image = utils.GetBase64StringFromData(image)
	}

	parents, err := reader.FindByChildTagID(ctx, tag.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting parents: %v", err)
	}

	newTagJSON.Parents = GetNames(parents)

	return &newTagJSON, nil
}

// GetDependentTagIDs returns a slice of unique tag IDs that this tag references.
func GetDependentTagIDs(ctx context.Context, reader FinderAliasImageGetter, tag *models.Tag) ([]int, error) {
	var ret []int

	parents, err := reader.FindByChildTagID(ctx, tag.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting parents: %v", err)
	}

	for _, tt := range parents {
		toAdd, err := GetDependentTagIDs(ctx, reader, tt)
		if err != nil {
			return nil, fmt.Errorf("error getting dependent tag IDs: %v", err)
		}

		ret = sliceutil.AppendUniques(ret, toAdd)
		ret = sliceutil.AppendUnique(ret, tt.ID)
	}

	return ret, nil
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
