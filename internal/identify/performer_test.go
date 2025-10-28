package identify

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"

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
	remoteSiteID := "2"
	name := "name"

	db := mocks.NewDatabase()

	db.Performer.On("Create", testCtx, mock.AnythingOfType("*models.CreatePerformerInput")).Run(func(args mock.Arguments) {
		p := args.Get(1).(*models.CreatePerformerInput)
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
					Name:         &name,
					RemoteSiteID: &remoteSiteID,
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
			got, err := getPerformerID(testCtx, tt.args.endpoint, db.Performer, tt.args.p, tt.args.createMissing, tt.args.skipSingleName)
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

	db := mocks.NewDatabase()

	db.Performer.On("Create", testCtx, mock.MatchedBy(func(p *models.CreatePerformerInput) bool {
		return p.Name == validName
	})).Run(func(args mock.Arguments) {
		p := args.Get(1).(*models.CreatePerformerInput)
		p.ID = performerID
	}).Return(nil)

	db.Performer.On("Create", testCtx, mock.MatchedBy(func(p *models.CreatePerformerInput) bool {
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
					Name:         &validName,
					RemoteSiteID: &remoteSiteID,
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
					Name:         &invalidName,
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
			got, err := createMissingPerformer(testCtx, tt.args.endpoint, db.Performer, tt.args.p)
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
