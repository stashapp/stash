package identify

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_createMissingStudio(t *testing.T) {
	emptyEndpoint := ""
	validEndpoint := "validEndpoint"
	invalidEndpoint := "invalidEndpoint"
	remoteSiteID := "remoteSiteID"
	validName := "validName"
	invalidName := "invalidName"
	createdID := 1

	db := mocks.NewDatabase()

	db.Studio.On("Create", testCtx, mock.MatchedBy(func(p *models.Studio) bool {
		return p.Name == validName
	})).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Studio)
		s.ID = createdID
	}).Return(nil)
	db.Studio.On("Create", testCtx, mock.MatchedBy(func(p *models.Studio) bool {
		return p.Name == invalidName
	})).Return(errors.New("error creating studio"))

	db.Studio.On("UpdatePartial", testCtx, models.StudioPartial{
		ID: createdID,
		StashIDs: &models.UpdateStashIDs{
			StashIDs: []models.StashID{
				{
					Endpoint: invalidEndpoint,
					StashID:  remoteSiteID,
				},
			},
			Mode: models.RelationshipUpdateModeSet,
		},
	}).Return(nil, errors.New("error updating stash ids"))
	db.Studio.On("UpdatePartial", testCtx, models.StudioPartial{
		ID: createdID,
		StashIDs: &models.UpdateStashIDs{
			StashIDs: []models.StashID{
				{
					Endpoint: validEndpoint,
					StashID:  remoteSiteID,
				},
			},
			Mode: models.RelationshipUpdateModeSet,
		},
	}).Return(models.Studio{
		ID: createdID,
	}, nil)

	type args struct {
		endpoint string
		studio   *models.ScrapedStudio
	}
	tests := []struct {
		name    string
		args    args
		want    *int
		wantErr bool
	}{
		{
			"simple",
			args{
				emptyEndpoint,
				&models.ScrapedStudio{
					Name:         validName,
					RemoteSiteID: &remoteSiteID,
				},
			},
			&createdID,
			false,
		},
		{
			"error creating",
			args{
				emptyEndpoint,
				&models.ScrapedStudio{
					Name:         invalidName,
					RemoteSiteID: &remoteSiteID,
				},
			},
			nil,
			true,
		},
		{
			"valid stash id",
			args{
				validEndpoint,
				&models.ScrapedStudio{
					Name:         validName,
					RemoteSiteID: &remoteSiteID,
				},
			},
			&createdID,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createMissingStudio(testCtx, tt.args.endpoint, db.Studio, tt.args.studio)
			if (err != nil) != tt.wantErr {
				t.Errorf("createMissingStudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMissingStudio() = %d, want %d", got, tt.want)
			}
		})
	}
}
