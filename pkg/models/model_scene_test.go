package models

import (
	"reflect"
	"testing"
)

func TestScenePartial_UpdateInput(t *testing.T) {
	const (
		id    = 1
		idStr = "1"
	)

	var (
		title       = "title"
		titlePtr    = &title
		details     = "details"
		detailsPtr  = &details
		url         = "url"
		urlPtr      = &url
		date        = "2001-02-03"
		rating      = 4
		ratingPtr   = &rating
		organized   = true
		studioID    = 2
		studioIDPtr = &studioID
		studioIDStr = "2"
	)

	dateObj := NewDate(date)
	dateObjPtr := &dateObj

	tests := []struct {
		name string
		id   int
		s    ScenePartial
		want SceneUpdateInput
	}{
		{
			"full",
			id,
			ScenePartial{
				Title:     &titlePtr,
				Details:   &detailsPtr,
				URL:       &urlPtr,
				Date:      &dateObjPtr,
				Rating:    &ratingPtr,
				Organized: &organized,
				StudioID:  &studioIDPtr,
			},
			SceneUpdateInput{
				ID:        idStr,
				Title:     &title,
				Details:   &details,
				URL:       &url,
				Date:      &date,
				Rating:    &rating,
				Organized: &organized,
				StudioID:  &studioIDStr,
			},
		},
		{
			"empty",
			id,
			ScenePartial{},
			SceneUpdateInput{
				ID: idStr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.UpdateInput(tt.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScenePartial.UpdateInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
