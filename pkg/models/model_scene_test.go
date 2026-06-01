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
		rating100   = 80
		organized   = true
		studioID    = 2
		studioIDStr = "2"
	)

	dateObj, _ := ParseDate(date)

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
				Title:    NewOptionalString(title),
				Code:     NewOptionalString(code),
				Details:  NewOptionalString(details),
				Director: NewOptionalString(director),
				URLs: &UpdateStrings{
					Values: []string{url},
					Mode:   RelationshipUpdateModeSet,
				},
				Date:      NewOptionalDate(dateObj),
				Rating:    NewOptionalInt(rating100),
				Organized: NewOptionalBool(organized),
				StudioID:  NewOptionalInt(studioID),
			},
			SceneUpdateInput{
				ID:        idStr,
				Title:     &title,
				Code:      &code,
				Details:   &details,
				Director:  &director,
				Urls:      []string{url},
				Date:      &date,
				Rating100: &rating100,
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

func TestScene_RangeDuration(t *testing.T) {
	startTime := 10.0
	endTime := 25.0

	tests := []struct {
		name         string
		scene        Scene
		fileDuration float64
		want         float64
	}{
		{
			name:         "full file",
			scene:        Scene{},
			fileDuration: 120,
			want:         120,
		},
		{
			name: "start only",
			scene: Scene{
				StartTime: &startTime,
			},
			fileDuration: 120,
			want:         110,
		},
		{
			name: "end only",
			scene: Scene{
				EndTime: &endTime,
			},
			fileDuration: 120,
			want:         25,
		},
		{
			name: "bounded range",
			scene: Scene{
				StartTime: &startTime,
				EndTime:   &endTime,
			},
			fileDuration: 120,
			want:         15,
		},
		{
			name: "invalid range is clamped to zero duration",
			scene: Scene{
				StartTime: &endTime,
				EndTime:   &startTime,
			},
			fileDuration: 120,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.scene.RangeDuration(tt.fileDuration)
			if got == nil || *got != tt.want {
				t.Errorf("Scene.RangeDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateSceneFileRangeWithDuration(t *testing.T) {
	negative := -1.0
	zero := 0.0
	start := 10.0
	end := 25.0
	tooLong := 130.0

	tests := []struct {
		name      string
		startTime *float64
		endTime   *float64
		wantErr   bool
	}{
		{
			name:      "unbounded range",
			startTime: nil,
			endTime:   nil,
		},
		{
			name:      "valid bounded range",
			startTime: &start,
			endTime:   &end,
		},
		{
			name:      "negative start",
			startTime: &negative,
			wantErr:   true,
		},
		{
			name:    "negative end",
			endTime: &negative,
			wantErr: true,
		},
		{
			name:      "end equals start",
			startTime: &zero,
			endTime:   &zero,
			wantErr:   true,
		},
		{
			name:      "start after file duration",
			startTime: &tooLong,
			wantErr:   true,
		},
		{
			name:    "end after file duration",
			endTime: &tooLong,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSceneFileRangeWithDuration(tt.startTime, tt.endTime, 120)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSceneFileRangeWithDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
