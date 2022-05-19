package identify

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/utils"
	"github.com/stretchr/testify/mock"
)

func Test_sceneRelationships_studio(t *testing.T) {
	validStoredID := "1"
	var validStoredIDInt int64 = 1
	invalidStoredID := "invalidStoredID"
	createMissing := true

	defaultOptions := &FieldOptions{
		Strategy: FieldStrategyMerge,
	}

	mockStudioReaderWriter := &mocks.StudioReaderWriter{}
	mockStudioReaderWriter.On("Create", testCtx, mock.Anything).Return(&models.Studio{
		ID: int(validStoredIDInt),
	}, nil)

	tr := sceneRelationships{
		studioCreator: mockStudioReaderWriter,
		fieldOptions:  make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		scene        *models.Scene
		fieldOptions *FieldOptions
		result       *models.ScrapedStudio
		want         *int64
		wantErr      bool
	}{
		{
			"nil studio",
			&models.Scene{},
			defaultOptions,
			nil,
			nil,
			false,
		},
		{
			"ignore",
			&models.Scene{},
			&FieldOptions{
				Strategy: FieldStrategyIgnore,
			},
			&models.ScrapedStudio{
				StoredID: &validStoredID,
			},
			nil,
			false,
		},
		{
			"invalid stored id",
			&models.Scene{},
			defaultOptions,
			&models.ScrapedStudio{
				StoredID: &invalidStoredID,
			},
			nil,
			true,
		},
		{
			"same stored id",
			&models.Scene{
				StudioID: models.NullInt64(validStoredIDInt),
			},
			defaultOptions,
			&models.ScrapedStudio{
				StoredID: &validStoredID,
			},
			nil,
			false,
		},
		{
			"different stored id",
			&models.Scene{},
			defaultOptions,
			&models.ScrapedStudio{
				StoredID: &validStoredID,
			},
			&validStoredIDInt,
			false,
		},
		{
			"no create missing",
			&models.Scene{},
			defaultOptions,
			&models.ScrapedStudio{},
			nil,
			false,
		},
		{
			"create missing",
			&models.Scene{},
			&FieldOptions{
				Strategy:      FieldStrategyMerge,
				CreateMissing: &createMissing,
			},
			&models.ScrapedStudio{},
			&validStoredIDInt,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.scene = tt.scene
			tr.fieldOptions["studio"] = tt.fieldOptions
			tr.result = &scrapeResult{
				result: &scraper.ScrapedScene{
					Studio: tt.result,
				},
			}

			got, err := tr.studio(testCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneRelationships.studio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sceneRelationships.studio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneRelationships_performers(t *testing.T) {
	const (
		sceneID = iota
		sceneWithPerformerID
		errSceneID
		existingPerformerID
		validStoredIDInt
	)
	validStoredID := strconv.Itoa(validStoredIDInt)
	invalidStoredID := "invalidStoredID"
	createMissing := true
	existingPerformerStr := strconv.Itoa(existingPerformerID)
	validName := "validName"
	female := models.GenderEnumFemale.String()
	male := models.GenderEnumMale.String()

	defaultOptions := &FieldOptions{
		Strategy: FieldStrategyMerge,
	}

	mockSceneReaderWriter := &mocks.SceneReaderWriter{}
	mockSceneReaderWriter.On("GetPerformerIDs", testCtx, sceneID).Return(nil, nil)
	mockSceneReaderWriter.On("GetPerformerIDs", testCtx, sceneWithPerformerID).Return([]int{existingPerformerID}, nil)
	mockSceneReaderWriter.On("GetPerformerIDs", testCtx, errSceneID).Return(nil, errors.New("error getting IDs"))

	tr := sceneRelationships{
		sceneReader:  mockSceneReaderWriter,
		fieldOptions: make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		sceneID      int
		fieldOptions *FieldOptions
		scraped      []*models.ScrapedPerformer
		ignoreMale   bool
		want         []int
		wantErr      bool
	}{
		{
			"ignore",
			sceneID,
			&FieldOptions{
				Strategy: FieldStrategyIgnore,
			},
			[]*models.ScrapedPerformer{
				{
					StoredID: &validStoredID,
				},
			},
			false,
			nil,
			false,
		},
		{
			"none",
			sceneID,
			defaultOptions,
			[]*models.ScrapedPerformer{},
			false,
			nil,
			false,
		},
		{
			"error getting ids",
			errSceneID,
			defaultOptions,
			[]*models.ScrapedPerformer{
				{},
			},
			false,
			nil,
			true,
		},
		{
			"merge existing",
			sceneWithPerformerID,
			defaultOptions,
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &existingPerformerStr,
				},
			},
			false,
			nil,
			false,
		},
		{
			"merge add",
			sceneWithPerformerID,
			defaultOptions,
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &validStoredID,
				},
			},
			false,
			[]int{existingPerformerID, validStoredIDInt},
			false,
		},
		{
			"ignore male",
			sceneID,
			defaultOptions,
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &validStoredID,
					Gender:   &male,
				},
			},
			true,
			nil,
			false,
		},
		{
			"overwrite",
			sceneWithPerformerID,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &validStoredID,
				},
			},
			false,
			[]int{validStoredIDInt},
			false,
		},
		{
			"ignore male (not male)",
			sceneWithPerformerID,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &validStoredID,
					Gender:   &female,
				},
			},
			true,
			[]int{validStoredIDInt},
			false,
		},
		{
			"error getting tag ID",
			sceneID,
			&FieldOptions{
				Strategy:      FieldStrategyOverwrite,
				CreateMissing: &createMissing,
			},
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &invalidStoredID,
				},
			},
			false,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.scene = &models.Scene{
				ID: tt.sceneID,
			}
			tr.fieldOptions["performers"] = tt.fieldOptions
			tr.result = &scrapeResult{
				result: &scraper.ScrapedScene{
					Performers: tt.scraped,
				},
			}

			got, err := tr.performers(testCtx, tt.ignoreMale)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneRelationships.performers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sceneRelationships.performers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneRelationships_tags(t *testing.T) {
	const (
		sceneID = iota
		sceneWithTagID
		errSceneID
		existingID
		validStoredIDInt
	)
	validStoredID := strconv.Itoa(validStoredIDInt)
	invalidStoredID := "invalidStoredID"
	createMissing := true
	existingIDStr := strconv.Itoa(existingID)
	validName := "validName"
	invalidName := "invalidName"

	defaultOptions := &FieldOptions{
		Strategy: FieldStrategyMerge,
	}

	mockSceneReaderWriter := &mocks.SceneReaderWriter{}
	mockTagReaderWriter := &mocks.TagReaderWriter{}
	mockSceneReaderWriter.On("GetTagIDs", testCtx, sceneID).Return(nil, nil)
	mockSceneReaderWriter.On("GetTagIDs", testCtx, sceneWithTagID).Return([]int{existingID}, nil)
	mockSceneReaderWriter.On("GetTagIDs", testCtx, errSceneID).Return(nil, errors.New("error getting IDs"))

	mockTagReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p models.Tag) bool {
		return p.Name == validName
	})).Return(&models.Tag{
		ID: validStoredIDInt,
	}, nil)
	mockTagReaderWriter.On("Create", testCtx, mock.MatchedBy(func(p models.Tag) bool {
		return p.Name == invalidName
	})).Return(nil, errors.New("error creating tag"))

	tr := sceneRelationships{
		sceneReader:  mockSceneReaderWriter,
		tagCreator:   mockTagReaderWriter,
		fieldOptions: make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		sceneID      int
		fieldOptions *FieldOptions
		scraped      []*models.ScrapedTag
		want         []int
		wantErr      bool
	}{
		{
			"ignore",
			sceneID,
			&FieldOptions{
				Strategy: FieldStrategyIgnore,
			},
			[]*models.ScrapedTag{
				{
					StoredID: &validStoredID,
				},
			},
			nil,
			false,
		},
		{
			"none",
			sceneID,
			defaultOptions,
			[]*models.ScrapedTag{},
			nil,
			false,
		},
		{
			"error getting ids",
			errSceneID,
			defaultOptions,
			[]*models.ScrapedTag{
				{},
			},
			nil,
			true,
		},
		{
			"merge existing",
			sceneWithTagID,
			defaultOptions,
			[]*models.ScrapedTag{
				{
					Name:     validName,
					StoredID: &existingIDStr,
				},
			},
			nil,
			false,
		},
		{
			"merge add",
			sceneWithTagID,
			defaultOptions,
			[]*models.ScrapedTag{
				{
					Name:     validName,
					StoredID: &validStoredID,
				},
			},
			[]int{existingID, validStoredIDInt},
			false,
		},
		{
			"overwrite",
			sceneWithTagID,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			[]*models.ScrapedTag{
				{
					Name:     validName,
					StoredID: &validStoredID,
				},
			},
			[]int{validStoredIDInt},
			false,
		},
		{
			"error getting tag ID",
			sceneID,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			[]*models.ScrapedTag{
				{
					Name:     validName,
					StoredID: &invalidStoredID,
				},
			},
			nil,
			true,
		},
		{
			"create missing",
			sceneID,
			&FieldOptions{
				Strategy:      FieldStrategyOverwrite,
				CreateMissing: &createMissing,
			},
			[]*models.ScrapedTag{
				{
					Name: validName,
				},
			},
			[]int{validStoredIDInt},
			false,
		},
		{
			"error creating",
			sceneID,
			&FieldOptions{
				Strategy:      FieldStrategyOverwrite,
				CreateMissing: &createMissing,
			},
			[]*models.ScrapedTag{
				{
					Name: invalidName,
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.scene = &models.Scene{
				ID: tt.sceneID,
			}
			tr.fieldOptions["tags"] = tt.fieldOptions
			tr.result = &scrapeResult{
				result: &scraper.ScrapedScene{
					Tags: tt.scraped,
				},
			}

			got, err := tr.tags(testCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneRelationships.tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sceneRelationships.tags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneRelationships_stashIDs(t *testing.T) {
	const (
		sceneID = iota
		sceneWithStashID
		errSceneID
		existingID
		validStoredIDInt
	)
	existingEndpoint := "existingEndpoint"
	newEndpoint := "newEndpoint"
	remoteSiteID := "remoteSiteID"
	newRemoteSiteID := "newRemoteSiteID"

	defaultOptions := &FieldOptions{
		Strategy: FieldStrategyMerge,
	}

	mockSceneReaderWriter := &mocks.SceneReaderWriter{}
	mockSceneReaderWriter.On("GetStashIDs", testCtx, sceneID).Return(nil, nil)
	mockSceneReaderWriter.On("GetStashIDs", testCtx, sceneWithStashID).Return([]*models.StashID{
		{
			StashID:  remoteSiteID,
			Endpoint: existingEndpoint,
		},
	}, nil)
	mockSceneReaderWriter.On("GetStashIDs", testCtx, errSceneID).Return(nil, errors.New("error getting IDs"))

	tr := sceneRelationships{
		sceneReader:  mockSceneReaderWriter,
		fieldOptions: make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		sceneID      int
		fieldOptions *FieldOptions
		endpoint     string
		remoteSiteID *string
		want         []models.StashID
		wantErr      bool
	}{
		{
			"ignore",
			sceneID,
			&FieldOptions{
				Strategy: FieldStrategyIgnore,
			},
			newEndpoint,
			&remoteSiteID,
			nil,
			false,
		},
		{
			"no endpoint",
			sceneID,
			defaultOptions,
			"",
			&remoteSiteID,
			nil,
			false,
		},
		{
			"no site id",
			sceneID,
			defaultOptions,
			newEndpoint,
			nil,
			nil,
			false,
		},
		{
			"error getting ids",
			errSceneID,
			defaultOptions,
			newEndpoint,
			&remoteSiteID,
			nil,
			true,
		},
		{
			"merge existing",
			sceneWithStashID,
			defaultOptions,
			existingEndpoint,
			&remoteSiteID,
			nil,
			false,
		},
		{
			"merge existing new value",
			sceneWithStashID,
			defaultOptions,
			existingEndpoint,
			&newRemoteSiteID,
			[]models.StashID{
				{
					StashID:  newRemoteSiteID,
					Endpoint: existingEndpoint,
				},
			},
			false,
		},
		{
			"merge add",
			sceneWithStashID,
			defaultOptions,
			newEndpoint,
			&newRemoteSiteID,
			[]models.StashID{
				{
					StashID:  remoteSiteID,
					Endpoint: existingEndpoint,
				},
				{
					StashID:  newRemoteSiteID,
					Endpoint: newEndpoint,
				},
			},
			false,
		},
		{
			"overwrite",
			sceneWithStashID,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			newEndpoint,
			&newRemoteSiteID,
			[]models.StashID{
				{
					StashID:  newRemoteSiteID,
					Endpoint: newEndpoint,
				},
			},
			false,
		},
		{
			"overwrite same",
			sceneWithStashID,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			existingEndpoint,
			&remoteSiteID,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.scene = &models.Scene{
				ID: tt.sceneID,
			}
			tr.fieldOptions["stash_ids"] = tt.fieldOptions
			tr.result = &scrapeResult{
				source: ScraperSource{
					RemoteSite: tt.endpoint,
				},
				result: &scraper.ScrapedScene{
					RemoteSiteID: tt.remoteSiteID,
				},
			}

			got, err := tr.stashIDs(testCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneRelationships.stashIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sceneRelationships.stashIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneRelationships_cover(t *testing.T) {
	const (
		sceneID = iota
		sceneWithStashID
		errSceneID
		existingID
		validStoredIDInt
	)
	existingData := []byte("existingData")
	newData := []byte("newData")
	const base64Prefix = "data:image/png;base64,"
	existingDataEncoded := base64Prefix + utils.GetBase64StringFromData(existingData)
	newDataEncoded := base64Prefix + utils.GetBase64StringFromData(newData)
	invalidData := newDataEncoded + "!!!"

	mockSceneReaderWriter := &mocks.SceneReaderWriter{}
	mockSceneReaderWriter.On("GetCover", testCtx, sceneID).Return(existingData, nil)
	mockSceneReaderWriter.On("GetCover", testCtx, errSceneID).Return(nil, errors.New("error getting cover"))

	tr := sceneRelationships{
		sceneReader:  mockSceneReaderWriter,
		fieldOptions: make(map[string]*FieldOptions),
	}

	tests := []struct {
		name    string
		sceneID int
		image   *string
		want    []byte
		wantErr bool
	}{
		{
			"nil image",
			sceneID,
			nil,
			nil,
			false,
		},
		{
			"different image",
			sceneID,
			&newDataEncoded,
			newData,
			false,
		},
		{
			"same image",
			sceneID,
			&existingDataEncoded,
			nil,
			false,
		},
		{
			"error getting scene cover",
			errSceneID,
			&newDataEncoded,
			nil,
			true,
		},
		{
			"invalid data",
			sceneID,
			&invalidData,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.scene = &models.Scene{
				ID: tt.sceneID,
			}
			tr.result = &scrapeResult{
				result: &scraper.ScrapedScene{
					Image: tt.image,
				},
			}

			got, err := tr.cover(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("sceneRelationships.cover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sceneRelationships.cover() = %v, want %v", got, tt.want)
			}
		})
	}
}
