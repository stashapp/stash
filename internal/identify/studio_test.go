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
	createdID64 := int64(createdID)

	repo := mocks.NewTxnRepository()
	mockStudioReaderWriter := repo.Studio.(*mocks.StudioReaderWriter)
	mockStudioReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p models.Studio) bool {
		return p.Name.String == validName
	})).Return(&models.Studio{
		ID: createdID,
	}, nil)
	mockStudioReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p models.Studio) bool {
		return p.Name.String == invalidName
	})).Return(nil, errors.New("error creating performer"))

	mockStudioReaderWriter.On("UpdateStashIDs", testCtx, createdID, []models.StashID{
		{
			Endpoint: invalidEndpoint,
			StashID:  remoteSiteID,
		},
	}).Return(errors.New("error updating stash ids"))
	mockStudioReaderWriter.On("UpdateStashIDs", testCtx, createdID, []models.StashID{
		{
			Endpoint: validEndpoint,
			StashID:  remoteSiteID,
		},
	}).Return(nil)

	type args struct {
		endpoint string
		studio   *models.ScrapedStudio
	}
	tests := []struct {
		name    string
		args    args
		want    *int64
		wantErr bool
	}{
		{
			"simple",
			args{
				emptyEndpoint,
				&models.ScrapedStudio{
					Name: validName,
				},
			},
			&createdID64,
			false,
		},
		{
			"error creating",
			args{
				emptyEndpoint,
				&models.ScrapedStudio{
					Name: invalidName,
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
			&createdID64,
			false,
		},
		{
			"invalid stash id",
			args{
				invalidEndpoint,
				&models.ScrapedStudio{
					Name:         validName,
					RemoteSiteID: &remoteSiteID,
				},
			},
			nil,
			true,
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
				t.Errorf("createMissingStudio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scrapedToStudioInput(t *testing.T) {
	const name = "name"
	const md5 = "b068931cc450442b63f5b3d276ea4297"
	url := "url"

	tests := []struct {
		name   string
		studio *models.ScrapedStudio
		want   models.Studio
	}{
		{
			"set all",
			&models.ScrapedStudio{
				Name: name,
				URL:  &url,
			},
			models.Studio{
				Name:     models.NullString(name),
				Checksum: md5,
				URL:      models.NullString(url),
			},
		},
		{
			"set none",
			&models.ScrapedStudio{
				Name: name,
			},
			models.Studio{
				Name:     models.NullString(name),
				Checksum: md5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scrapedToStudioInput(tt.studio)

			// clear created/updated dates
			got.CreatedAt = models.SQLiteTimestamp{}
			got.UpdatedAt = got.CreatedAt

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scrapedToStudioInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
