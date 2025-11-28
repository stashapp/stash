package stashbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/stashapp/stash/pkg/models"
)

func (c Client) FindTag(ctx context.Context, query string) (*models.ScrapedTag, error) {
	var id *string
	var name *string

	_, err := uuid.Parse(query)
	if err == nil {
		// Confirmed the user passed in a Stash ID
		id = &query
	} else {
		// Otherwise assume they're searching on a name
		name = &query
	}

	tag, err := c.client.FindTag(ctx, id, name)
	if err != nil {
		return nil, err
	}

	if tag.FindTag == nil {
		return nil, nil
	}

	return &models.ScrapedTag{
		Name:         tag.FindTag.Name,
		RemoteSiteID: &tag.FindTag.ID,
	}, nil
}
