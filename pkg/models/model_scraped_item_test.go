package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_scrapedToStudioInput(t *testing.T) {
	const name = "name"
	url := "url"
	remoteSiteID := "remoteSiteID"

	tests := []struct {
		name   string
		studio *ScrapedStudio
		want   *Studio
	}{
		{
			"set all",
			&ScrapedStudio{
				Name:         name,
				URL:          &url,
				RemoteSiteID: &remoteSiteID,
			},
			&Studio{
				Name: name,
				URL:  url,
				StashIDs: NewRelatedStashIDs([]StashID{
					{
						StashID: remoteSiteID,
					},
				}),
			},
		},
		{
			"set none",
			&ScrapedStudio{
				Name:         name,
				RemoteSiteID: &remoteSiteID,
			},
			&Studio{
				Name: name,
				StashIDs: NewRelatedStashIDs([]StashID{
					{
						StashID: remoteSiteID,
					},
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.studio.ToStudio("", nil)

			assert.NotEqual(t, time.Time{}, got.CreatedAt)
			assert.NotEqual(t, time.Time{}, got.UpdatedAt)

			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
			assert.Equal(t, tt.want, got)
		})
	}
}
