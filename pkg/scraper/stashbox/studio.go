package stashbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
)

func (c Client) FindStashBoxStudio(ctx context.Context, query string) (*models.ScrapedStudio, error) {
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
		r := c.repository
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			ret = studioFragmentToScrapedStudio(*studio.FindStudio)

			err = match.ScrapedStudio(ctx, r.Studio, ret, &c.box.Endpoint)
			if err != nil {
				return err
			}

			if studio.FindStudio.Parent != nil {
				parentStudio, err := c.client.FindStudio(ctx, &studio.FindStudio.Parent.ID, nil)
				if err != nil {
					return err
				}

				if parentStudio.FindStudio != nil {
					ret.Parent = studioFragmentToScrapedStudio(*parentStudio.FindStudio)

					err = match.ScrapedStudio(ctx, r.Studio, ret.Parent, &c.box.Endpoint)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func studioFragmentToScrapedStudio(s graphql.StudioFragment) *models.ScrapedStudio {
	images := []string{}
	for _, image := range s.Images {
		images = append(images, image.URL)
	}

	st := &models.ScrapedStudio{
		Name:         s.Name,
		URL:          findURL(s.Urls, "HOME"),
		Images:       images,
		RemoteSiteID: &s.ID,
	}

	if len(st.Images) > 0 {
		st.Image = &st.Images[0]
	}

	return st
}
