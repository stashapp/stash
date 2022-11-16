package identify

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_getPerformerID(t *testing.T) {
	const (
		emptyEndpoint = ""
		endpoint      = "endpoint"
	)
	invalidStoredID := "invalidStoredID"
	validStoredIDStr := "1"
	validStoredID := 1
	name := "name"

	mockPerformerReaderWriter := mocks.PerformerReaderWriter{}
	mockPerformerReaderWriter.On("Create", testCtx, mock.Anything).Run(func(args mock.Arguments) {
		p := args.Get(1).(*models.Performer)
		p.ID = validStoredID
	}).Return(nil)

	type args struct {
		endpoint      string
		p             *models.ScrapedPerformer
		createMissing bool
	}
	tests := []struct {
		name    string
		args    args
		want    *int
		wantErr bool
	}{
		{
			"no performer",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{},
				false,
			},
			nil,
			false,
		},
		{
			"invalid stored id",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{
					StoredID: &invalidStoredID,
				},
				false,
			},
			nil,
			true,
		},
		{
			"valid stored id",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{
					StoredID: &validStoredIDStr,
				},
				false,
			},
			&validStoredID,
			false,
		},
		{
			"nil stored not creating",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{
					Name: &name,
				},
				false,
			},
			nil,
			false,
		},
		{
			"nil name creating",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{},
				true,
			},
			nil,
			false,
		},
		{
			"valid name creating",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{
					Name: &name,
				},
				true,
			},
			&validStoredID,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPerformerID(testCtx, tt.args.endpoint, &mockPerformerReaderWriter, tt.args.p, tt.args.createMissing)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPerformerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPerformerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createMissingPerformer(t *testing.T) {
	emptyEndpoint := ""
	validEndpoint := "validEndpoint"
	invalidEndpoint := "invalidEndpoint"
	remoteSiteID := "remoteSiteID"
	validName := "validName"
	invalidName := "invalidName"
	performerID := 1

	mockPerformerReaderWriter := mocks.PerformerReaderWriter{}
	mockPerformerReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p *models.Performer) bool {
		return p.Name == validName
	})).Run(func(args mock.Arguments) {
		p := args.Get(1).(*models.Performer)
		p.ID = performerID
	}).Return(nil)

	mockPerformerReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p *models.Performer) bool {
		return p.Name == invalidName
	})).Return(errors.New("error creating performer"))

	mockPerformerReaderWriter.On("UpdateStashIDs", testCtx, performerID, []models.StashID{
		{
			Endpoint: invalidEndpoint,
			StashID:  remoteSiteID,
		},
	}).Return(errors.New("error updating stash ids"))
	mockPerformerReaderWriter.On("UpdateStashIDs", testCtx, performerID, []models.StashID{
		{
			Endpoint: validEndpoint,
			StashID:  remoteSiteID,
		},
	}).Return(nil)

	type args struct {
		endpoint string
		p        *models.ScrapedPerformer
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
				&models.ScrapedPerformer{
					Name: &validName,
				},
			},
			&performerID,
			false,
		},
		{
			"error creating",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{
					Name: &invalidName,
				},
			},
			nil,
			true,
		},
		{
			"valid stash id",
			args{
				validEndpoint,
				&models.ScrapedPerformer{
					Name:         &validName,
					RemoteSiteID: &remoteSiteID,
				},
			},
			&performerID,
			false,
		},
		{
			"invalid stash id",
			args{
				invalidEndpoint,
				&models.ScrapedPerformer{
					Name:         &validName,
					RemoteSiteID: &remoteSiteID,
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createMissingPerformer(testCtx, tt.args.endpoint, &mockPerformerReaderWriter, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("createMissingPerformer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMissingPerformer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scrapedToPerformerInput(t *testing.T) {
	name := "name"
	md5 := "b068931cc450442b63f5b3d276ea4297"

	var stringValues []string
	for i := 0; i < 17; i++ {
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

	dateToDatePtr := func(d models.Date) *models.Date {
		return &d
	}

	tests := []struct {
		name      string
		performer *models.ScrapedPerformer
		want      models.Performer
	}{
		{
			"set all",
			&models.ScrapedPerformer{
				Name:         &name,
				Birthdate:    nextVal(),
				DeathDate:    nextVal(),
				Gender:       nextVal(),
				Ethnicity:    nextVal(),
				Country:      nextVal(),
				EyeColor:     nextVal(),
				HairColor:    nextVal(),
				Height:       nextVal(),
				Weight:       nextVal(),
				Measurements: nextVal(),
				FakeTits:     nextVal(),
				CareerLength: nextVal(),
				Tattoos:      nextVal(),
				Piercings:    nextVal(),
				Aliases:      nextVal(),
				Twitter:      nextVal(),
				Instagram:    nextVal(),
			},
			models.Performer{
				Name:         name,
				Checksum:     md5,
				Birthdate:    dateToDatePtr(models.NewDate(*nextVal())),
				DeathDate:    dateToDatePtr(models.NewDate(*nextVal())),
				Gender:       models.GenderEnum(*nextVal()),
				Ethnicity:    *nextVal(),
				Country:      *nextVal(),
				EyeColor:     *nextVal(),
				HairColor:    *nextVal(),
				Height:       nextIntVal(),
				Weight:       nextIntVal(),
				Measurements: *nextVal(),
				FakeTits:     *nextVal(),
				CareerLength: *nextVal(),
				Tattoos:      *nextVal(),
				Piercings:    *nextVal(),
				Aliases:      *nextVal(),
				Twitter:      *nextVal(),
				Instagram:    *nextVal(),
			},
		},
		{
			"set none",
			&models.ScrapedPerformer{
				Name: &name,
			},
			models.Performer{
				Name:     name,
				Checksum: md5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scrapedToPerformerInput(tt.performer)

			// clear created/updated dates
			got.CreatedAt = time.Time{}
			got.UpdatedAt = got.CreatedAt

			assert.Equal(t, tt.want, got)
		})
	}
}
