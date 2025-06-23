package identify

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/utils"
	"github.com/stretchr/testify/mock"
)

func Test_sceneRelationships_studio(t *testing.T) {
	validStoredID := "1"
	remoteSiteID := "2"
	var validStoredIDInt = 1
	invalidStoredID := "invalidStoredID"
	createMissing := true

	defaultOptions := &FieldOptions{
		Strategy: FieldStrategyMerge,
	}

	db := mocks.NewDatabase()

	db.Studio.On("Create", testCtx, mock.Anything).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Studio)
		s.ID = validStoredIDInt
	}).Return(nil)

	tr := sceneRelationships{
		studioReaderWriter: db.Studio,
		fieldOptions:       make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		scene        *models.Scene
		fieldOptions *FieldOptions
		result       *models.ScrapedStudio
		want         *int
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
				StudioID: &validStoredIDInt,
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
			&models.ScrapedStudio{RemoteSiteID: &remoteSiteID},
			&validStoredIDInt,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.scene = tt.scene
			tr.fieldOptions["studio"] = tt.fieldOptions
			tr.result = &scrapeResult{
				source: ScraperSource{
					RemoteSite: "endpoint",
				},
				result: &models.ScrapedScene{
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

	emptyScene := &models.Scene{
		ID:           sceneID,
		PerformerIDs: models.NewRelatedIDs([]int{}),
		TagIDs:       models.NewRelatedIDs([]int{}),
		StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
	}

	sceneWithPerformer := &models.Scene{
		ID: sceneWithPerformerID,
		PerformerIDs: models.NewRelatedIDs([]int{
			existingPerformerID,
		}),
	}

	db := mocks.NewDatabase()

	tr := sceneRelationships{
		sceneReader:  db.Scene,
		fieldOptions: make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		scene        *models.Scene
		fieldOptions *FieldOptions
		scraped      []*models.ScrapedPerformer
		ignoreMale   bool
		want         []int
		wantErr      bool
	}{
		{
			"ignore",
			emptyScene,
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
			emptyScene,
			defaultOptions,
			[]*models.ScrapedPerformer{},
			false,
			nil,
			false,
		},
		{
			"merge existing",
			sceneWithPerformer,
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
			sceneWithPerformer,
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
			emptyScene,
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
			sceneWithPerformer,
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
			sceneWithPerformer,
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
			emptyScene,
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
			tr.scene = tt.scene
			tr.fieldOptions["performers"] = tt.fieldOptions
			tr.result = &scrapeResult{
				result: &models.ScrapedScene{
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

	emptyScene := &models.Scene{
		ID:           sceneID,
		TagIDs:       models.NewRelatedIDs([]int{}),
		PerformerIDs: models.NewRelatedIDs([]int{}),
		StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
	}

	sceneWithTag := &models.Scene{
		ID: sceneWithTagID,
		TagIDs: models.NewRelatedIDs([]int{
			existingID,
		}),
		PerformerIDs: models.NewRelatedIDs([]int{}),
		StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
	}

	db := mocks.NewDatabase()

	db.Tag.On("Create", testCtx, mock.MatchedBy(func(p *models.Tag) bool {
		return p.Name == validName
	})).Run(func(args mock.Arguments) {
		t := args.Get(1).(*models.Tag)
		t.ID = validStoredIDInt
	}).Return(nil)
	db.Tag.On("Create", testCtx, mock.MatchedBy(func(p *models.Tag) bool {
		return p.Name == invalidName
	})).Return(errors.New("error creating tag"))

	tr := sceneRelationships{
		sceneReader:  db.Scene,
		tagCreator:   db.Tag,
		fieldOptions: make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		scene        *models.Scene
		fieldOptions *FieldOptions
		scraped      []*models.ScrapedTag
		want         []int
		wantErr      bool
	}{
		{
			"ignore",
			emptyScene,
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
			emptyScene,
			defaultOptions,
			[]*models.ScrapedTag{},
			nil,
			false,
		},
		{
			"merge existing",
			sceneWithTag,
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
			sceneWithTag,
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
			sceneWithTag,
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
			emptyScene,
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
			emptyScene,
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
			emptyScene,
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
			tr.scene = tt.scene
			tr.fieldOptions["tags"] = tt.fieldOptions
			tr.result = &scrapeResult{
				result: &models.ScrapedScene{
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

	emptyScene := &models.Scene{
		ID: sceneID,
	}

	sceneWithStashIDs := &models.Scene{
		ID: sceneWithStashID,
		StashIDs: models.NewRelatedStashIDs([]models.StashID{
			{
				StashID:   remoteSiteID,
				Endpoint:  existingEndpoint,
				UpdatedAt: time.Time{},
			},
		}),
	}

	db := mocks.NewDatabase()

	tr := sceneRelationships{
		sceneReader:  db.Scene,
		fieldOptions: make(map[string]*FieldOptions),
	}

	setTime := time.Now()

	tests := []struct {
		name          string
		scene         *models.Scene
		fieldOptions  *FieldOptions
		endpoint      string
		remoteSiteID  *string
		setUpdateTime bool
		want          []models.StashID
		wantErr       bool
	}{
		{
			"ignore",
			emptyScene,
			&FieldOptions{
				Strategy: FieldStrategyIgnore,
			},
			newEndpoint,
			&remoteSiteID,
			false,
			nil,
			false,
		},
		{
			"no endpoint",
			emptyScene,
			defaultOptions,
			"",
			&remoteSiteID,
			false,
			nil,
			false,
		},
		{
			"no site id",
			emptyScene,
			defaultOptions,
			newEndpoint,
			nil,
			false,
			nil,
			false,
		},
		{
			"merge existing",
			sceneWithStashIDs,
			defaultOptions,
			existingEndpoint,
			&remoteSiteID,
			false,
			nil,
			false,
		},
		{
			"merge existing set update time",
			sceneWithStashIDs,
			defaultOptions,
			existingEndpoint,
			&remoteSiteID,
			true,
			[]models.StashID{
				{
					StashID:   remoteSiteID,
					Endpoint:  existingEndpoint,
					UpdatedAt: setTime,
				},
			},
			false,
		},
		{
			"merge existing new value",
			sceneWithStashIDs,
			defaultOptions,
			existingEndpoint,
			&newRemoteSiteID,
			false,
			[]models.StashID{
				{
					StashID:   newRemoteSiteID,
					Endpoint:  existingEndpoint,
					UpdatedAt: setTime,
				},
			},
			false,
		},
		{
			"merge add",
			sceneWithStashIDs,
			defaultOptions,
			newEndpoint,
			&newRemoteSiteID,
			false,
			[]models.StashID{
				{
					StashID:   remoteSiteID,
					Endpoint:  existingEndpoint,
					UpdatedAt: time.Time{},
				},
				{
					StashID:   newRemoteSiteID,
					Endpoint:  newEndpoint,
					UpdatedAt: setTime,
				},
			},
			false,
		},
		{
			"overwrite",
			sceneWithStashIDs,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			newEndpoint,
			&newRemoteSiteID,
			false,
			[]models.StashID{
				{
					StashID:   newRemoteSiteID,
					Endpoint:  newEndpoint,
					UpdatedAt: setTime,
				},
			},
			false,
		},
		{
			"overwrite same",
			sceneWithStashIDs,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			existingEndpoint,
			&remoteSiteID,
			false,
			nil,
			false,
		},
		{
			"overwrite same set update time",
			sceneWithStashIDs,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			existingEndpoint,
			&remoteSiteID,
			true,
			[]models.StashID{
				{
					StashID:   remoteSiteID,
					Endpoint:  existingEndpoint,
					UpdatedAt: setTime,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.scene = tt.scene
			tr.fieldOptions["stash_ids"] = tt.fieldOptions
			tr.result = &scrapeResult{
				source: ScraperSource{
					RemoteSite: tt.endpoint,
				},
				result: &models.ScrapedScene{
					RemoteSiteID: tt.remoteSiteID,
				},
			}

			got, err := tr.stashIDs(testCtx, tt.setUpdateTime)

			if (err != nil) != tt.wantErr {
				t.Errorf("sceneRelationships.stashIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// massage updatedAt times to be consistent for comparison
			for i := range got {
				if !got[i].UpdatedAt.IsZero() {
					got[i].UpdatedAt = setTime
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sceneRelationships.stashIDs() = %+v, want %+v", got, tt.want)
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

	db := mocks.NewDatabase()

	db.Scene.On("GetCover", testCtx, sceneID).Return(existingData, nil)
	db.Scene.On("GetCover", testCtx, errSceneID).Return(nil, errors.New("error getting cover"))

	tr := sceneRelationships{
		sceneReader:  db.Scene,
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
			newData,
			false,
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
				result: &models.ScrapedScene{
					Image: tt.image,
				},
			}

			got, err := tr.cover(testCtx)
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
