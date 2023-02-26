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

	repo := mocks.NewTxnRepository()
	mockStudioReaderWriter := repo.Studio.(*mocks.StudioReaderWriter)
	mockStudioReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p models.StudioDBInput) bool {
		return p.StudioCreate.Name == validName
	})).Return(&createdID, nil)
	mockStudioReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p models.StudioDBInput) bool {
		return p.StudioCreate.Name == invalidName
	})).Return(nil, errors.New("error creating studio"))

	mockStudioReaderWriter.On("UpdatePartial", testCtx, models.StudioDBInput{
		StudioUpdate: &models.StudioPartial{
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
		},
	}).Return(nil, errors.New("error updating stash ids"))
	mockStudioReaderWriter.On("UpdatePartial", testCtx, models.StudioDBInput{
		StudioUpdate: &models.StudioPartial{
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
			got, err := createMissingStudio(testCtx, tt.args.endpoint, mockStudioReaderWriter, tt.args.studio)
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

func Test_scrapedToStudioInput(t *testing.T) {
	const name = "name"
	url := "url"
	remoteSiteID := "remoteSiteID"

	tests := []struct {
		name   string
		studio *models.ScrapedStudio
		want   *models.Studio
	}{
		{
			"set all",
			&models.ScrapedStudio{
				Name:         name,
				URL:          &url,
				RemoteSiteID: &remoteSiteID,
			},
			&models.Studio{
				Name: name,
				URL:  url,
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						StashID: remoteSiteID,
					},
				}),
			},
		},
		{
			"set none",
			&models.ScrapedStudio{
				Name:         name,
				RemoteSiteID: &remoteSiteID,
			},
			&models.Studio{
				Name: name,
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						StashID: remoteSiteID,
					},
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := studioFromScrapedStudio(testCtx, tt.studio, "")

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s, scrapedToStudioInput() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
