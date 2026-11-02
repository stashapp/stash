package stashbox

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/stashbox/graphql"
)

func (c Client) FindStudio(ctx context.Context, query string) (*models.ScrapedStudio, error) {
	var studio *graphql.FindStudio

	_, err := uuid.Parse(query)
	if err == nil {
		// Confirmed the user passed in a Stash ID
		studio, err = c.client.FindStudio(ctx, &query, nil)
	} else {
		// Otherwise assume they're searching on a name
		studio, err = c.client.FindStudio(ctx, nil, &query)
	}

	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedStudio
	if studio.FindStudio != nil {
		ret = studioFragmentToScrapedStudio(*studio.FindStudio)
	}

	return ret, nil
}

func studioFragmentToScrapedStudio(s graphql.StudioFragment) *models.ScrapedStudio {
	images := []string{}
	for _, image := range s.Images {
		images = append(images, image.URL)
	}

	aliases := strings.Join(s.Aliases, ", ")

	st := &models.ScrapedStudio{
		Name:         s.Name,
		Aliases:      &aliases,
		Images:       images,
		RemoteSiteID: &s.ID,
	}

	for _, u := range s.Urls {
		st.URLs = append(st.URLs, u.URL)
	}

	if len(st.Images) > 0 {
		st.Image = &st.Images[0]
	}

	if s.Parent != nil {
		st.Parent = studioFragmentToScrapedStudio(graphql.StudioFragment{
			Name:    s.Parent.Name,
			ID:      s.Parent.ID,
			Aliases: s.Parent.Aliases,
			Urls:    s.Parent.Urls,
			Images:  s.Parent.Images,
		})
	}

	return st
}
