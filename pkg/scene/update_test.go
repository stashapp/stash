package scene

import (
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"

	"github.com/stretchr/testify/mock"
)

func TestUpdater_IsEmpty(t *testing.T) {
	organized := true
	ids := []int{1}
	stashIDs := []models.StashID{
		{},
	}
	cover := []byte{1}

	tests := []struct {
		name string
		u    *Updater
		want bool
	}{
		{
			"empty",
			&Updater{},
			true,
		},
		{
			"id only",
			&Updater{
				Partial: models.ScenePartial{
					ID: 1,
				},
			},
			true,
		},
		{
			"partial set",
			&Updater{
				Partial: models.ScenePartial{
					Organized: &organized,
				},
			},
			false,
		},
		{
			"performer set",
			&Updater{
				PerformerIDs: ids,
			},
			false,
		},
		{
			"tags set",
			&Updater{
				TagIDs: ids,
			},
			false,
		},
		{
			"performer set",
			&Updater{
				StashIDs: stashIDs,
			},
			false,
		},
		{
			"cover set",
			&Updater{
				CoverImage: cover,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsEmpty(); got != tt.want {
				t.Errorf("Updater.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockScreenshotSetter struct{}

func (s *mockScreenshotSetter) SetScreenshot(scene *models.Scene, imageData []byte) error {
	return nil
}

func TestUpdater_Update(t *testing.T) {
	const (
		sceneID = iota + 1
		badUpdateID
		badPerformersID
		badTagsID
		badStashIDsID
		badCoverID
		performerID
		tagID
	)

	performerIDs := []int{performerID}
	tagIDs := []int{tagID}
	stashID := "stashID"
	endpoint := "endpoint"
	stashIDs := []models.StashID{
		{
			StashID:  stashID,
			Endpoint: endpoint,
		},
	}

	title := "title"
	cover := []byte("cover")

	validScene := &models.Scene{}

	updateErr := errors.New("error updating")

	qb := mocks.SceneReaderWriter{}
	qb.On("Update", mock.MatchedBy(func(s models.ScenePartial) bool {
		return s.ID != badUpdateID
	})).Return(validScene, nil)
	qb.On("Update", mock.MatchedBy(func(s models.ScenePartial) bool {
		return s.ID == badUpdateID
	})).Return(nil, updateErr)

	qb.On("UpdatePerformers", sceneID, performerIDs).Return(nil).Once()
	qb.On("UpdateTags", sceneID, tagIDs).Return(nil).Once()
	qb.On("UpdateStashIDs", sceneID, stashIDs).Return(nil).Once()
	qb.On("UpdateCover", sceneID, cover).Return(nil).Once()

	qb.On("UpdatePerformers", badPerformersID, performerIDs).Return(updateErr).Once()
	qb.On("UpdateTags", badTagsID, tagIDs).Return(updateErr).Once()
	qb.On("UpdateStashIDs", badStashIDsID, stashIDs).Return(updateErr).Once()
	qb.On("UpdateCover", badCoverID, cover).Return(updateErr).Once()

	tests := []struct {
		name    string
		u       *Updater
		wantNil bool
		wantErr bool
	}{
		{
			"empty",
			&Updater{
				ID: sceneID,
			},
			true,
			true,
		},
		{
			"update all",
			&Updater{
				ID:           sceneID,
				PerformerIDs: performerIDs,
				TagIDs:       tagIDs,
				StashIDs: []models.StashID{
					{
						StashID:  stashID,
						Endpoint: endpoint,
					},
				},
				CoverImage: cover,
			},
			false,
			false,
		},
		{
			"update fields only",
			&Updater{
				ID: sceneID,
				Partial: models.ScenePartial{
					Title: models.NullStringPtr(title),
				},
			},
			false,
			false,
		},
		{
			"error updating scene",
			&Updater{
				ID: badUpdateID,
				Partial: models.ScenePartial{
					Title: models.NullStringPtr(title),
				},
			},
			true,
			true,
		},
		{
			"error updating performers",
			&Updater{
				ID:           badPerformersID,
				PerformerIDs: performerIDs,
			},
			true,
			true,
		},
		{
			"error updating tags",
			&Updater{
				ID:     badTagsID,
				TagIDs: tagIDs,
			},
			true,
			true,
		},
		{
			"error updating stash IDs",
			&Updater{
				ID:       badStashIDsID,
				StashIDs: stashIDs,
			},
			true,
			true,
		},
		{
			"error updating cover",
			&Updater{
				ID:         badCoverID,
				CoverImage: cover,
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.Update(&qb, &mockScreenshotSetter{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Updater.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("Updater.Update() = %v, want %v", got, tt.wantNil)
			}
		})
	}

	qb.AssertExpectations(t)
}
