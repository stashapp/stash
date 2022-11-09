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
		code        = "1337"
		details     = "details"
		director    = "director"
		url         = "url"
		date        = "2001-02-03"
		rating      = 4
		organized   = true
		studioID    = 2
		studioIDStr = "2"
	)

	dateObj := NewDate(date)

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
				Title:     NewOptionalString(title),
				Code:      NewOptionalString(code),
				Details:   NewOptionalString(details),
				Director:  NewOptionalString(director),
				URL:       NewOptionalString(url),
				Date:      NewOptionalDate(dateObj),
				Rating:    NewOptionalInt(rating),
				Organized: NewOptionalBool(organized),
				StudioID:  NewOptionalInt(studioID),
			},
			SceneUpdateInput{
				ID:        idStr,
				Title:     &title,
				Code:      &code,
				Details:   &details,
				Director:  &director,
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
