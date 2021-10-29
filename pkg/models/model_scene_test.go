package models

import (
	"database/sql"
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
		details     = "details"
		url         = "url"
		date        = "2001-02-03"
		rating      = 4
		organized   = true
		studioID    = 2
		studioIDStr = "2"
	)

	tests := []struct {
		name string
		s    ScenePartial
		want SceneUpdateInput
	}{
		{
			"full",
			ScenePartial{
				ID:      id,
				Title:   NullStringPtr(title),
				Details: NullStringPtr(details),
				URL:     NullStringPtr(url),
				Date: &SQLiteDate{
					String: date,
					Valid:  true,
				},
				Rating: &sql.NullInt64{
					Int64: int64(rating),
					Valid: true,
				},
				Organized: &organized,
				StudioID: &sql.NullInt64{
					Int64: int64(studioID),
					Valid: true,
				},
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
			ScenePartial{
				ID: id,
			},
			SceneUpdateInput{
				ID: idStr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.UpdateInput(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScenePartial.UpdateInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
