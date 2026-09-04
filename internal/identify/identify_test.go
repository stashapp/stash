package identify

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testCtx = context.Background()

type mockSceneScraper struct {
	errIDs  []int
	results map[int][]*models.ScrapedScene
}

func (s mockSceneScraper) ScrapeScenes(ctx context.Context, sceneID int) ([]*models.ScrapedScene, error) {
	if slices.Contains(s.errIDs, sceneID) {
		return nil, errors.New("scrape scene error")
	}
	return s.results[sceneID], nil
}

type mockHookExecutor struct {
}

func (s mockHookExecutor) ExecuteSceneUpdatePostHooks(ctx context.Context, input models.SceneUpdateInput, inputFields []string) {
}

func (s mockHookExecutor) ExecuteGalleryUpdatePostHooks(ctx context.Context, input models.GalleryUpdateInput, inputFields []string) {
}

func TestSceneIdentifier_Identify(t *testing.T) {
	const (
		errID1 = iota
		errID2
		missingID
		found1ID
		found2ID
		multiFoundID
		multiFound2ID
		multiInvalidTagID
		errUpdateID
	)

	var (
		skipMultipleTagID    = 1
		skipMultipleTagIDStr = strconv.Itoa(skipMultipleTagID)
	)

	var (
		scrapedTitle  = "scrapedTitle"
		scrapedTitle2 = "scrapedTitle2"
		invalidTagID  = "invalid"

		boolFalse = false
		boolTrue  = true
	)

	defaultOptions := &MetadataOptions{
		SetOrganized:  &boolFalse,
		SetCoverImage: &boolFalse,
		PerformerGenders: []models.GenderEnum{
			models.GenderEnumFemale,
			models.GenderEnumTransgenderFemale,
			models.GenderEnumTransgenderMale,
			models.GenderEnumIntersex,
			models.GenderEnumNonBinary,
		},
		SkipSingleNamePerformers: &boolFalse,
	}
	sources := []ScraperSource{
		{
			Scraper: mockSceneScraper{
				errIDs: []int{errID1},
				results: map[int][]*models.ScrapedScene{
					found1ID: {{
						Title: &scrapedTitle,
					}},
				},
			},
		},
		{
			Scraper: mockSceneScraper{
				errIDs: []int{errID2},
				results: map[int][]*models.ScrapedScene{
					found2ID: {{
						Title: &scrapedTitle,
					}},
					errUpdateID: {{
						Title: &scrapedTitle,
					}},
					multiFoundID: {
						{
							Title: &scrapedTitle,
						},
						{
							Title: &scrapedTitle2,
						},
					},
					multiFound2ID: {
						{
							Title: &scrapedTitle,
						},
						{
							Title: &scrapedTitle2,
						},
					},
					multiInvalidTagID: {
						{
							Title: &scrapedTitle,
						},
						{
							Title: &scrapedTitle2,
						},
					},
				},
			},
		},
	}

	db := mocks.NewDatabase()

	db.Scene.On("GetURLs", mock.Anything, mock.Anything).Return(nil, nil)
	db.Scene.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(id int) bool {
		return id == errUpdateID
	}), mock.Anything).Return(nil, errors.New("update error"))
	db.Scene.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(id int) bool {
		return id != errUpdateID
	}), mock.Anything).Return(nil, nil)

	db.Tag.On("Find", mock.Anything, skipMultipleTagID).Return(&models.Tag{
		ID:   skipMultipleTagID,
		Name: skipMultipleTagIDStr,
	}, nil)

	tests := []struct {
		name    string
		sceneID int
		options *MetadataOptions
		wantErr bool
	}{
		{
			"error scraping",
			errID1,
			nil,
			false,
		},
		{
			"error scraping from second",
			errID2,
			nil,
			false,
		},
		{
			"found in first scraper",
			found1ID,
			nil,
			false,
		},
		{
			"found in second scraper",
			found2ID,
			nil,
			false,
		},
		{
			"not found",
			missingID,
			nil,
			false,
		},
		{
			"error modifying",
			errUpdateID,
			nil,
			true,
		},
		{
			"multiple found",
			multiFoundID,
			nil,
			false,
		},
		{
			"multiple found - set tag",
			multiFound2ID,
			&MetadataOptions{
				SkipMultipleMatches:  &boolTrue,
				SkipMultipleMatchTag: &skipMultipleTagIDStr,
			},
			false,
		},
		{
			"multiple found - invalid tag",
			multiInvalidTagID,
			&MetadataOptions{
				SkipMultipleMatches:  &boolTrue,
				SkipMultipleMatchTag: &invalidTagID,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identifier := SceneIdentifier{
				TxnManager:                  db,
				SceneReaderUpdater:          db.Scene,
				StudioReaderWriter:          db.Studio,
				PerformerCreator:            db.Performer,
				TagFinderCreator:            db.Tag,
				DefaultOptions:              defaultOptions,
				Sources:                     sources,
				SceneUpdatePostHookExecutor: mockHookExecutor{},
			}

			if tt.options != nil {
				identifier.DefaultOptions = tt.options
			}

			scene := &models.Scene{
				ID:           tt.sceneID,
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			}
			if err := identifier.Identify(testCtx, scene); (err != nil) != tt.wantErr {
				t.Errorf("SceneIdentifier.Identify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSceneIdentifier_modifyScene(t *testing.T) {
	db := mocks.NewDatabase()

	boolFalse := false
	defaultOptions := &MetadataOptions{
		SetOrganized:  &boolFalse,
		SetCoverImage: &boolFalse,
		PerformerGenders: []models.GenderEnum{
			models.GenderEnumFemale,
			models.GenderEnumTransgenderFemale,
			models.GenderEnumTransgenderMale,
			models.GenderEnumIntersex,
			models.GenderEnumNonBinary,
		},
		SkipSingleNamePerformers: &boolFalse,
	}
	tr := &SceneIdentifier{
		TxnManager:         db,
		SceneReaderUpdater: db.Scene,
		StudioReaderWriter: db.Studio,
		PerformerCreator:   db.Performer,
		TagFinderCreator:   db.Tag,
		DefaultOptions:     defaultOptions,
	}

	type args struct {
		scene  *models.Scene
		result *scrapeResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty update",
			args{
				&models.Scene{
					URLs:         models.NewRelatedStrings([]string{}),
					PerformerIDs: models.NewRelatedIDs([]int{}),
					TagIDs:       models.NewRelatedIDs([]int{}),
					StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
				},
				&scrapeResult{
					result: &models.ScrapedScene{},
					source: ScraperSource{
						Options: defaultOptions,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tr.modifyScene(testCtx, tt.args.scene, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("SceneIdentifier.modifyScene() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSceneIdentifier_modifyScene_paths(t *testing.T) {
	const sceneID = 1
	scrapedTitle := "scrapedTitle"
	scrapedCode := "scrapedCode"
	invalidID := "invalid"
	boolFalse := false
	sourceOptions := &MetadataOptions{SetCoverImage: &boolFalse}

	tests := []struct {
		name     string
		scene    *models.Scene
		scraped  *models.ScrapedScene
		executor SceneUpdatePostHookExecutor
		setup    func(db *mocks.Database)
		wantErr  bool
	}{
		{
			name: "load URLs error",
			scene: &models.Scene{
				ID:           sceneID,
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			},
			setup: func(db *mocks.Database) {
				db.Scene.On("GetURLs", mock.Anything, sceneID).Return(nil, errors.New("load URLs error"))
			},
			wantErr: true,
		},
		{
			name: "load performers error",
			scene: &models.Scene{
				ID:       sceneID,
				URLs:     models.NewRelatedStrings([]string{}),
				TagIDs:   models.NewRelatedIDs([]int{}),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{}),
			},
			setup: func(db *mocks.Database) {
				db.Scene.On("GetPerformerIDs", mock.Anything, sceneID).Return(nil, errors.New("load performers error"))
			},
			wantErr: true,
		},
		{
			name: "load tags error",
			scene: &models.Scene{
				ID:           sceneID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			},
			setup: func(db *mocks.Database) {
				db.Scene.On("GetTagIDs", mock.Anything, sceneID).Return(nil, errors.New("load tags error"))
			},
			wantErr: true,
		},
		{
			name: "load stash IDs error",
			scene: &models.Scene{
				ID:           sceneID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Scene.On("GetStashIDs", mock.Anything, sceneID).Return(nil, errors.New("load stash IDs error"))
			},
			wantErr: true,
		},
		{
			name: "updater error",
			scene: &models.Scene{
				ID:           sceneID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			},
			scraped: &models.ScrapedScene{
				Studio: &models.ScrapedStudio{StoredID: &invalidID},
			},
			wantErr: true,
		},
		{
			name: "update error",
			scene: &models.Scene{
				ID:           sceneID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			},
			scraped: &models.ScrapedScene{Title: &scrapedTitle},
			setup: func(db *mocks.Database) {
				db.Scene.On("UpdatePartial", mock.Anything, sceneID, mock.AnythingOfType("models.ScenePartial")).Return(nil, errors.New("update error"))
			},
			wantErr: true,
		},
		{
			name: "success with title",
			scene: &models.Scene{
				ID:           sceneID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			},
			scraped:  &models.ScrapedScene{Title: &scrapedTitle},
			executor: mockHookExecutor{},
			setup: func(db *mocks.Database) {
				db.Scene.On("UpdatePartial", mock.Anything, sceneID, mock.AnythingOfType("models.ScenePartial")).Return(&models.Scene{ID: sceneID}, nil)
			},
		},
		{
			name: "success without title",
			scene: &models.Scene{
				ID:           sceneID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			},
			scraped:  &models.ScrapedScene{Code: &scrapedCode},
			executor: mockHookExecutor{},
			setup: func(db *mocks.Database) {
				db.Scene.On("UpdatePartial", mock.Anything, sceneID, mock.AnythingOfType("models.ScenePartial")).Return(&models.Scene{ID: sceneID}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()
			if tt.setup != nil {
				tt.setup(db)
			}

			identifier := SceneIdentifier{
				TxnManager:                  db,
				SceneReaderUpdater:          db.Scene,
				StudioReaderWriter:          db.Studio,
				PerformerCreator:            db.Performer,
				TagFinderCreator:            db.Tag,
				SceneUpdatePostHookExecutor: tt.executor,
			}
			result := &scrapeResult{
				source: ScraperSource{Options: sourceOptions},
				result: tt.scraped,
			}

			err := identifier.modifyScene(testCtx, tt.scene, result)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_getFieldOptions(t *testing.T) {
	const (
		inFirst  = "inFirst"
		inSecond = "inSecond"
		inBoth   = "inBoth"
	)

	type args struct {
		options []MetadataOptions
	}
	tests := []struct {
		name string
		args args
		want map[string]*FieldOptions
	}{
		{
			"simple",
			args{
				[]MetadataOptions{
					{
						FieldOptions: []*FieldOptions{
							{
								Field:    inFirst,
								Strategy: FieldStrategyIgnore,
							},
							{
								Field:    inBoth,
								Strategy: FieldStrategyIgnore,
							},
						},
					},
					{
						FieldOptions: []*FieldOptions{
							{
								Field:    inSecond,
								Strategy: FieldStrategyMerge,
							},
							{
								Field:    inBoth,
								Strategy: FieldStrategyMerge,
							},
						},
					},
				},
			},
			map[string]*FieldOptions{
				inFirst: {
					Field:    inFirst,
					Strategy: FieldStrategyIgnore,
				},
				inSecond: {
					Field:    inSecond,
					Strategy: FieldStrategyMerge,
				},
				inBoth: {
					Field:    inBoth,
					Strategy: FieldStrategyIgnore,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldOptions(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultipleMatchesFoundError_Error(t *testing.T) {
	err := &MultipleMatchesFoundError{
		Source: ScraperSource{Name: "Test Source"},
	}

	assert.Equal(t, "multiple matches found for Test Source", err.Error())
}

func Test_mergeSharedMetadataOptions(t *testing.T) {
	boolFalse := false
	emptyTag := ""
	tag := "tag"
	genders := []models.GenderEnum{models.GenderEnumFemale}

	base := MetadataOptions{
		SetOrganized:               &boolFalse,
		IncludeMalePerformers:      &boolFalse,
		PerformerGenders:           genders,
		SkipMultipleMatches:        &boolFalse,
		SkipMultipleMatchTag:       &emptyTag,
		SkipSingleNamePerformers:   &boolFalse,
		SkipSingleNamePerformerTag: &emptyTag,
	}

	assert.Equal(t, base, mergeSharedMetadataOptions(base, nil))

	override := &MetadataOptions{
		SetOrganized:               &boolFalse,
		IncludeMalePerformers:      &boolFalse,
		PerformerGenders:           genders,
		SkipMultipleMatches:        &boolFalse,
		SkipMultipleMatchTag:       &tag,
		SkipSingleNamePerformers:   &boolFalse,
		SkipSingleNamePerformerTag: &tag,
	}

	assert.Equal(t, *override, mergeSharedMetadataOptions(base, override))
}

func Test_getSceneUpdater_errorsAndSourceOptions(t *testing.T) {
	const sceneID = 1

	var (
		boolFalse     = false
		boolTrue      = true
		createMissing = true
		singleName    = "Single"
		studioID      = "1"
		tagID         = "2"
		skipTagID     = 3
		invalidID     = "invalid"
		remoteSiteID  = "remote"
		emptyImage    = ""
		invalidImage  = "data:image/png;base64,%%%%"
	)

	skipTagIDStr := strconv.Itoa(skipTagID)

	tests := []struct {
		name     string
		options  *MetadataOptions
		scraped  *models.ScrapedScene
		setup    func(db *mocks.Database)
		wantErr  bool
		validate func(t *testing.T, updater *scene.UpdateSet)
	}{
		{
			name: "source options with skipped performer tag and stash IDs",
			options: &MetadataOptions{
				SetCoverImage: &boolFalse,
				FieldOptions: []*FieldOptions{
					{
						Field:         "performers",
						Strategy:      FieldStrategyMerge,
						CreateMissing: &createMissing,
					},
				},
				SkipSingleNamePerformers:   &boolTrue,
				SkipSingleNamePerformerTag: &skipTagIDStr,
			},
			scraped: &models.ScrapedScene{
				Studio:       &models.ScrapedStudio{StoredID: &studioID},
				Performers:   []*models.ScrapedPerformer{{Name: &singleName}},
				Tags:         []*models.ScrapedTag{{StoredID: &tagID}},
				RemoteSiteID: &remoteSiteID,
			},
			validate: func(t *testing.T, updater *scene.UpdateSet) {
				assert.True(t, updater.Partial.StudioID.Set)
				assert.Equal(t, []int{2, skipTagID}, updater.Partial.TagIDs.IDs)
				assert.Len(t, updater.Partial.StashIDs.StashIDs, 1)
				assert.Equal(t, "endpoint", updater.Partial.StashIDs.StashIDs[0].Endpoint)
				assert.Equal(t, remoteSiteID, updater.Partial.StashIDs.StashIDs[0].StashID)
			},
		},
		{
			name: "invalid skipped performer tag",
			options: &MetadataOptions{
				SetCoverImage: &boolFalse,
				FieldOptions: []*FieldOptions{
					{
						Field:         "performers",
						Strategy:      FieldStrategyMerge,
						CreateMissing: &createMissing,
					},
				},
				SkipSingleNamePerformers:   &boolTrue,
				SkipSingleNamePerformerTag: &invalidID,
			},
			scraped: &models.ScrapedScene{
				Performers: []*models.ScrapedPerformer{{Name: &singleName}},
			},
			wantErr: true,
		},
		{
			name: "skipped performer without tag",
			options: &MetadataOptions{
				SetCoverImage: &boolFalse,
				FieldOptions: []*FieldOptions{
					{
						Field:         "performers",
						Strategy:      FieldStrategyMerge,
						CreateMissing: &createMissing,
					},
				},
				SkipSingleNamePerformers: &boolTrue,
			},
			scraped: &models.ScrapedScene{
				Performers: []*models.ScrapedPerformer{{Name: &singleName}},
			},
			validate: func(t *testing.T, updater *scene.UpdateSet) {
				assert.Nil(t, updater.Partial.PerformerIDs)
				assert.Nil(t, updater.Partial.TagIDs)
			},
		},
		{
			name:    "studio error",
			options: &MetadataOptions{SetCoverImage: &boolFalse},
			scraped: &models.ScrapedScene{Studio: &models.ScrapedStudio{StoredID: &invalidID}},
			wantErr: true,
		},
		{
			name:    "performer error",
			options: &MetadataOptions{SetCoverImage: &boolFalse},
			scraped: &models.ScrapedScene{Performers: []*models.ScrapedPerformer{{StoredID: &invalidID}}},
			wantErr: true,
		},
		{
			name:    "tag error",
			options: &MetadataOptions{SetCoverImage: &boolFalse},
			scraped: &models.ScrapedScene{Tags: []*models.ScrapedTag{{StoredID: &invalidID}}},
			wantErr: true,
		},
		{
			name:    "cover called with empty image",
			options: &MetadataOptions{SetCoverImage: &boolTrue},
			scraped: &models.ScrapedScene{Image: &emptyImage},
			validate: func(t *testing.T, updater *scene.UpdateSet) {
				assert.Nil(t, updater.CoverImage)
			},
		},
		{
			name:    "cover error",
			options: &MetadataOptions{SetCoverImage: &boolTrue},
			scraped: &models.ScrapedScene{Image: &invalidImage},
			setup: func(db *mocks.Database) {
				db.Scene.On("GetCover", testCtx, sceneID).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()
			if tt.setup != nil {
				tt.setup(db)
			}

			identifier := SceneIdentifier{
				SceneReaderUpdater: db.Scene,
				StudioReaderWriter: db.Studio,
				PerformerCreator:   db.Performer,
				TagFinderCreator:   db.Tag,
			}
			sceneObj := &models.Scene{
				ID:           sceneID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
				StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
			}
			result := &scrapeResult{
				source: ScraperSource{
					Options:    tt.options,
					RemoteSite: "endpoint",
				},
				result: tt.scraped,
			}

			updater, err := identifier.getSceneUpdater(testCtx, sceneObj, result)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, updater)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updater)
				if tt.validate != nil {
					tt.validate(t, updater)
				}
			}
		})
	}
}

func Test_getScenePartial(t *testing.T) {
	var (
		originalTitle   = "originalTitle"
		originalDate    = "2001-01-01"
		originalDetails = "originalDetails"
		originalURL     = "originalURL"
	)

	var (
		scrapedTitle   = "scrapedTitle"
		scrapedDate    = "2002-02-02"
		scrapedDetails = "scrapedDetails"
		scrapedURL     = "scrapedURL"
	)

	originalDateObj, _ := models.ParseDate(originalDate)
	scrapedDateObj, _ := models.ParseDate(scrapedDate)

	originalScene := &models.Scene{
		Title:   originalTitle,
		Date:    &originalDateObj,
		Details: originalDetails,
		URLs:    models.NewRelatedStrings([]string{originalURL}),
	}

	organisedScene := *originalScene
	organisedScene.Organized = true

	emptyScene := &models.Scene{
		URLs: models.NewRelatedStrings([]string{}),
	}

	postPartial := models.ScenePartial{
		Title:   models.NewOptionalString(scrapedTitle),
		Date:    models.NewOptionalDate(scrapedDateObj),
		Details: models.NewOptionalString(scrapedDetails),
		URLs: &models.UpdateStrings{
			Values: []string{scrapedURL},
			Mode:   models.RelationshipUpdateModeSet,
		},
	}

	postPartialMerge := postPartial
	postPartialMerge.URLs = &models.UpdateStrings{
		Values: []string{scrapedURL},
		Mode:   models.RelationshipUpdateModeSet,
	}

	scrapedScene := &models.ScrapedScene{
		Title:   &scrapedTitle,
		Date:    &scrapedDate,
		Details: &scrapedDetails,
		URLs:    []string{scrapedURL},
	}

	scrapedUnchangedScene := &models.ScrapedScene{
		Title:   &originalTitle,
		Date:    &originalDate,
		Details: &originalDetails,
		URLs:    []string{originalURL},
	}

	makeFieldOptions := func(input *FieldOptions) map[string]*FieldOptions {
		return map[string]*FieldOptions{
			"title":   input,
			"date":    input,
			"details": input,
			"url":     input,
		}
	}

	overwriteAll := makeFieldOptions(&FieldOptions{
		Strategy: FieldStrategyOverwrite,
	})
	ignoreAll := makeFieldOptions(&FieldOptions{
		Strategy: FieldStrategyIgnore,
	})
	mergeAll := makeFieldOptions(&FieldOptions{
		Strategy: FieldStrategyMerge,
	})

	setOrganised := true

	type args struct {
		scene        *models.Scene
		scraped      *models.ScrapedScene
		fieldOptions map[string]*FieldOptions
		setOrganized bool
	}
	tests := []struct {
		name string
		args args
		want models.ScenePartial
	}{
		{
			"overwrite all",
			args{
				originalScene,
				scrapedScene,
				overwriteAll,
				false,
			},
			postPartial,
		},
		{
			"ignore all",
			args{
				originalScene,
				scrapedScene,
				ignoreAll,
				false,
			},
			models.ScenePartial{},
		},
		{
			"merge (existing values)",
			args{
				originalScene,
				scrapedScene,
				mergeAll,
				false,
			},
			models.ScenePartial{
				URLs: &models.UpdateStrings{
					Values: []string{originalURL, scrapedURL},
					Mode:   models.RelationshipUpdateModeSet,
				},
			},
		},
		{
			"merge (empty values)",
			args{
				emptyScene,
				scrapedScene,
				mergeAll,
				false,
			},
			postPartialMerge,
		},
		{
			"unchanged",
			args{
				originalScene,
				scrapedUnchangedScene,
				overwriteAll,
				false,
			},
			models.ScenePartial{},
		},
		{
			"overwrite url removal",
			args{
				originalScene,
				&models.ScrapedScene{
					URLs: []string{scrapedURL},
				},
				overwriteAll,
				false,
			},
			models.ScenePartial{
				URLs: &models.UpdateStrings{
					Values: []string{scrapedURL},
					Mode:   models.RelationshipUpdateModeSet,
				},
			},
		},
		{
			"set organized",
			args{
				originalScene,
				scrapedUnchangedScene,
				overwriteAll,
				true,
			},
			models.ScenePartial{
				Organized: models.NewOptionalBool(setOrganised),
			},
		},
		{
			"set organized unchanged",
			args{
				&organisedScene,
				scrapedUnchangedScene,
				overwriteAll,
				true,
			},
			models.ScenePartial{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getScenePartial(tt.args.scene, tt.args.scraped, tt.args.fieldOptions, tt.args.setOrganized)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getOptions(t *testing.T) {
	boolTrue := true
	boolFalse := false

	tests := []struct {
		name           string
		defaultOptions *MetadataOptions
		sourceOptions  *MetadataOptions
		want           MetadataOptions
	}{
		{
			name: "nil source options returns defaults unchanged",
			defaultOptions: &MetadataOptions{
				SetOrganized:     &boolTrue,
				PerformerGenders: []models.GenderEnum{models.GenderEnumFemale},
			},
			sourceOptions: nil,
			want: MetadataOptions{
				SetOrganized:     &boolTrue,
				PerformerGenders: []models.GenderEnum{models.GenderEnumFemale},
			},
		},
		{
			name: "nil PerformerGenders in source does not override default",
			defaultOptions: &MetadataOptions{
				PerformerGenders: []models.GenderEnum{models.GenderEnumFemale},
			},
			sourceOptions: &MetadataOptions{
				PerformerGenders: nil,
			},
			want: MetadataOptions{
				PerformerGenders: []models.GenderEnum{models.GenderEnumFemale},
			},
		},
		{
			// When the UI sends an empty performerGenders array (all genders allowed),
			// it must not be treated as a filter that blocks all performers.
			name: "empty PerformerGenders in source does not override default",
			defaultOptions: &MetadataOptions{
				PerformerGenders: []models.GenderEnum{models.GenderEnumFemale},
			},
			sourceOptions: &MetadataOptions{
				PerformerGenders: []models.GenderEnum{},
			},
			want: MetadataOptions{
				PerformerGenders: []models.GenderEnum{models.GenderEnumFemale},
			},
		},
		{
			name: "non-empty PerformerGenders in source overrides default",
			defaultOptions: &MetadataOptions{
				PerformerGenders: []models.GenderEnum{models.GenderEnumFemale},
			},
			sourceOptions: &MetadataOptions{
				PerformerGenders: []models.GenderEnum{models.GenderEnumMale},
			},
			want: MetadataOptions{
				PerformerGenders: []models.GenderEnum{models.GenderEnumMale},
			},
		},
		{
			// Empty source PerformerGenders with nil default means no filter (include all).
			name:           "empty PerformerGenders with nil default yields nil (no filter)",
			defaultOptions: nil,
			sourceOptions: &MetadataOptions{
				PerformerGenders: []models.GenderEnum{},
			},
			want: MetadataOptions{},
		},
		{
			name: "source overrides scalar fields",
			defaultOptions: &MetadataOptions{
				SetCoverImage: &boolTrue,
				SetOrganized:  &boolTrue,
			},
			sourceOptions: &MetadataOptions{
				SetCoverImage: &boolFalse,
			},
			want: MetadataOptions{
				SetCoverImage: &boolFalse,
				SetOrganized:  &boolTrue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identifier := SceneIdentifier{
				DefaultOptions: tt.defaultOptions,
			}
			source := ScraperSource{Options: tt.sourceOptions}
			got := identifier.getOptions(source)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getSceneUpdater_performerGenders(t *testing.T) {
	const (
		femalePerformerID = 1
		malePerformerID   = 2
	)
	femaleIDStr := strconv.Itoa(femalePerformerID)
	maleIDStr := strconv.Itoa(malePerformerID)
	female := models.GenderEnumFemale.String()
	male := models.GenderEnumMale.String()
	boolFalse := false

	db := mocks.NewDatabase()

	// Scene with no existing performers; all relationship fields pre-loaded so
	// getSceneUpdater does not need to hit the database.
	scene := &models.Scene{
		ID:           1,
		URLs:         models.NewRelatedStrings([]string{}),
		PerformerIDs: models.NewRelatedIDs([]int{}),
		TagIDs:       models.NewRelatedIDs([]int{}),
		StashIDs:     models.NewRelatedStashIDs([]models.StashID{}),
	}

	// Scraped scene with one female and one male performer, both already in the
	// database (StoredID set). No studio, tags, cover or stash IDs, so no DB
	// calls are needed beyond resolving the performer IDs.
	scrapedScene := &models.ScrapedScene{
		Performers: []*models.ScrapedPerformer{
			{StoredID: &femaleIDStr, Gender: &female},
			{StoredID: &maleIDStr, Gender: &male},
		},
	}

	tests := []struct {
		name                  string
		performerGenders      []models.GenderEnum
		includeMalePerformers *bool
		wantPerformerIDs      []int
	}{
		{
			// nil means "no filter configured" — all performers pass through.
			name:             "nil PerformerGenders includes all genders",
			performerGenders: nil,
			wantPerformerIDs: []int{femalePerformerID, malePerformerID},
		},
		{
			// An empty slice sent by the UI when no gender restriction is set must
			// also mean "no filter". This was the root cause of the identify bug.
			name:             "empty PerformerGenders includes all genders",
			performerGenders: []models.GenderEnum{},
			wantPerformerIDs: []int{femalePerformerID, malePerformerID},
		},
		{
			name:             "female-only filter excludes male performer",
			performerGenders: []models.GenderEnum{models.GenderEnumFemale},
			wantPerformerIDs: []int{femalePerformerID},
		},
		{
			name:             "male-only filter excludes female performer",
			performerGenders: []models.GenderEnum{models.GenderEnumMale},
			wantPerformerIDs: []int{malePerformerID},
		},
		{
			// Legacy field: empty PerformerGenders falls back to IncludeMalePerformers.
			name:                  "empty PerformerGenders falls back to IncludeMalePerformers=false",
			performerGenders:      []models.GenderEnum{},
			includeMalePerformers: &boolFalse,
			wantPerformerIDs:      []int{femalePerformerID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := MetadataOptions{
				SetCoverImage:         &boolFalse,
				PerformerGenders:      tt.performerGenders,
				IncludeMalePerformers: tt.includeMalePerformers,
			}

			identifier := SceneIdentifier{
				TxnManager:         db,
				SceneReaderUpdater: db.Scene,
				StudioReaderWriter: db.Studio,
				PerformerCreator:   db.Performer,
				TagFinderCreator:   db.Tag,
				DefaultOptions:     &opts,
			}

			result := &scrapeResult{
				source: ScraperSource{},
				result: scrapedScene,
			}

			updater, err := identifier.getSceneUpdater(testCtx, scene, result)
			assert.NoError(t, err)
			assert.NotNil(t, updater.Partial.PerformerIDs, "expected PerformerIDs to be set")
			assert.ElementsMatch(t, tt.wantPerformerIDs, updater.Partial.PerformerIDs.IDs)
		})
	}
}

func Test_shouldSetSingleValueField(t *testing.T) {
	const invalid = "invalid"

	type args struct {
		strategy         *FieldOptions
		hasExistingValue bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"ignore",
			args{
				&FieldOptions{
					Strategy: FieldStrategyIgnore,
				},
				false,
			},
			false,
		},
		{
			"merge existing",
			args{
				&FieldOptions{
					Strategy: FieldStrategyMerge,
				},
				true,
			},
			false,
		},
		{
			"merge absent",
			args{
				&FieldOptions{
					Strategy: FieldStrategyMerge,
				},
				false,
			},
			true,
		},
		{
			"overwrite",
			args{
				&FieldOptions{
					Strategy: FieldStrategyOverwrite,
				},
				true,
			},
			true,
		},
		{
			"nil (merge) existing",
			args{
				&FieldOptions{},
				true,
			},
			false,
		},
		{
			"nil (merge) absent",
			args{
				&FieldOptions{},
				false,
			},
			true,
		},
		{
			"invalid (merge) existing",
			args{
				&FieldOptions{
					Strategy: invalid,
				},
				true,
			},
			false,
		},
		{
			"invalid (merge) absent",
			args{
				&FieldOptions{
					Strategy: invalid,
				},
				false,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSetSingleValueField(tt.args.strategy, tt.args.hasExistingValue); got != tt.want {
				t.Errorf("shouldSetSingleValueField() = %v, want %v", got, tt.want)
			}
		})
	}
}
