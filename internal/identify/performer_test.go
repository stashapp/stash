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
		endpoint       string
		p              *models.ScrapedPerformer
		createMissing  bool
		skipSingleName bool
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
				false,
			},
			nil,
			false,
		},
		{
			"single name no disambig creating",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{
					Name: &name,
				},
				true,
				true,
			},
			nil,
			true,
		},
		{
			"valid name creating",
			args{
				emptyEndpoint,
				&models.ScrapedPerformer{
					Name: &name,
				},
				true,
				false,
			},
			&validStoredID,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPerformerID(testCtx, tt.args.endpoint, &mockPerformerReaderWriter, tt.args.p, tt.args.createMissing, tt.args.skipSingleName)
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

	dateFromInt := func(i int) *models.Date {
		t := time.Date(2001, 1, i, 0, 0, 0, 0, time.UTC)
		d := models.Date{Time: t}
		return &d
	}
	dateStrFromInt := func(i int) *string {
		s := dateFromInt(i).String()
		return &s
	}

	genderFromInt := func(i int) *models.GenderEnum {
		g := models.AllGenderEnum[i%len(models.AllGenderEnum)]
		return &g
	}
	genderStrFromInt := func(i int) *string {
		s := genderFromInt(i).String()
		return &s
	}

	tests := []struct {
		name      string
		performer *models.ScrapedPerformer
		want      models.Performer
	}{
		{
			"set all",
			&models.ScrapedPerformer{
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
			},
			models.Performer{
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
				Aliases:        models.NewRelatedStrings([]string{*nextVal()}),
				Twitter:        *nextVal(),
				Instagram:      *nextVal(),
				URL:            *nextVal(),
				Details:        *nextVal(),
			},
		},
		{
			"set none",
			&models.ScrapedPerformer{
				Name: &name,
			},
			models.Performer{
				Name: name,
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
