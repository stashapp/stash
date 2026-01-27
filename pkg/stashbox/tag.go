package stashbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/stashbox/graphql"
)

// QueryTag searches for tags by name or ID.
// If query is a valid UUID, it searches by ID (returns single result).
// Otherwise, it searches by name (returns multiple results).
func (c Client) QueryTag(ctx context.Context, query string) ([]*models.ScrapedTag, error) {
	_, err := uuid.Parse(query)
	if err == nil {
		// Query is a UUID, use findTag for exact match
		return c.findTagByID(ctx, query)
	}
	// Otherwise search by name
	return c.queryTagsByName(ctx, query)
}

func (c Client) findTagByID(ctx context.Context, id string) ([]*models.ScrapedTag, error) {
	tag, err := c.client.FindTag(ctx, &id, nil)
	if err != nil {
		return nil, err
	}

	if tag.FindTag == nil {
		return nil, nil
	}

	return []*models.ScrapedTag{{
		Name:         tag.FindTag.Name,
		RemoteSiteID: &tag.FindTag.ID,
	}}, nil
}

func (c Client) queryTagsByName(ctx context.Context, name string) ([]*models.ScrapedTag, error) {
	input := graphql.TagQueryInput{
		Name:      &name,
		Page:      1,
		PerPage:   25,
		Direction: graphql.SortDirectionEnumAsc,
		Sort:      graphql.TagSortEnumName,
	}

	result, err := c.client.QueryTags(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.QueryTags.Tags == nil {
		return nil, nil
	}

	var ret []*models.ScrapedTag
	for _, t := range result.QueryTags.Tags {
		ret = append(ret, &models.ScrapedTag{
			Name:         t.Name,
			RemoteSiteID: &t.ID,
		})
	}

	return ret, nil
}
