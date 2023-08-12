package models

import (
	"strconv"
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

func Test_scrapedToPerformerInput(t *testing.T) {
	name := "name"
	remoteSiteID := "remoteSiteID"

	var stringValues []string
	for i := 0; i < 20; i++ {
		stringValues = append(stringValues, strconv.Itoa(i))
	}

	upTo := 0
	nextVal := func() *string {
		ret := stringValues[upTo]
		upTo = (upTo + 1) % len(stringValues)
		return &ret
	}

	nextIntVal := func() *int {
		ret := upTo
		upTo = (upTo + 1) % len(stringValues)
		return &ret
	}

	dateFromInt := func(i int) *Date {
		t := time.Date(2001, 1, i, 0, 0, 0, 0, time.UTC)
		d := Date{Time: t}
		return &d
	}
	dateStrFromInt := func(i int) *string {
		s := dateFromInt(i).String()
		return &s
	}

	genderFromInt := func(i int) *GenderEnum {
		g := AllGenderEnum[i%len(AllGenderEnum)]
		return &g
	}
	genderStrFromInt := func(i int) *string {
		s := genderFromInt(i).String()
		return &s
	}

	tests := []struct {
		name      string
		performer *ScrapedPerformer
		want      *Performer
	}{
		{
			"set all",
			&ScrapedPerformer{
				Name:           &name,
				Disambiguation: nextVal(),
				Birthdate:      dateStrFromInt(*nextIntVal()),
				DeathDate:      dateStrFromInt(*nextIntVal()),
				Gender:         genderStrFromInt(*nextIntVal()),
				Ethnicity:      nextVal(),
				Country:        nextVal(),
				EyeColor:       nextVal(),
				HairColor:      nextVal(),
				Height:         nextVal(),
				Weight:         nextVal(),
				Measurements:   nextVal(),
				FakeTits:       nextVal(),
				CareerLength:   nextVal(),
				Tattoos:        nextVal(),
				Piercings:      nextVal(),
				Aliases:        nextVal(),
				Twitter:        nextVal(),
				Instagram:      nextVal(),
				URL:            nextVal(),
				Details:        nextVal(),
				RemoteSiteID:   &remoteSiteID,
			},
			&Performer{
				Name:           name,
				Disambiguation: *nextVal(),
				Birthdate:      dateFromInt(*nextIntVal()),
				DeathDate:      dateFromInt(*nextIntVal()),
				Gender:         genderFromInt(*nextIntVal()),
				Ethnicity:      *nextVal(),
				Country:        *nextVal(),
				EyeColor:       *nextVal(),
				HairColor:      *nextVal(),
				Height:         nextIntVal(),
				Weight:         nextIntVal(),
				Measurements:   *nextVal(),
				FakeTits:       *nextVal(),
				CareerLength:   *nextVal(),
				Tattoos:        *nextVal(),
				Piercings:      *nextVal(),
				Aliases:        NewRelatedStrings([]string{*nextVal()}),
				Twitter:        *nextVal(),
				Instagram:      *nextVal(),
				URL:            *nextVal(),
				Details:        *nextVal(),
				StashIDs: NewRelatedStashIDs([]StashID{
					{
						StashID: remoteSiteID,
					},
				}),
			},
		},
		{
			"set none",
			&ScrapedPerformer{
				Name:         &name,
				RemoteSiteID: &remoteSiteID,
			},
			&Performer{
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
			got := tt.performer.ToPerformer("", nil)

			assert.NotEqual(t, time.Time{}, got.CreatedAt)
			assert.NotEqual(t, time.Time{}, got.UpdatedAt)

			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
			assert.Equal(t, tt.want, got)
		})
	}
}
