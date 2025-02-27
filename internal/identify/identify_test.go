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
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testCtx = context.Background()

type mockSceneScraper struct {
	errIDs  []int
	results map[int][]*scraper.ScrapedScene
}

func (s mockSceneScraper) ScrapeScenes(ctx context.Context, sceneID int) ([]*scraper.ScrapedScene, error) {
	if slices.Contains(s.errIDs, sceneID) {
		return nil, errors.New("scrape scene error")
	}
	return s.results[sceneID], nil
}

type mockHookExecutor struct {
}

func (s mockHookExecutor) ExecuteSceneUpdatePostHooks(ctx context.Context, input models.SceneUpdateInput, inputFields []string) {
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
		errUpdateID
	)

	var (
		skipMultipleTagID    = 1
		skipMultipleTagIDStr = strconv.Itoa(skipMultipleTagID)
	)

	var (
		scrapedTitle  = "scrapedTitle"
		scrapedTitle2 = "scrapedTitle2"

		boolFalse = false
		boolTrue  = true
	)

	defaultOptions := &MetadataOptions{
		SetOrganized:             &boolFalse,
		SetCoverImage:            &boolFalse,
		IncludeMalePerformers:    &boolFalse,
		SkipSingleNamePerformers: &boolFalse,
	}
	sources := []ScraperSource{
		{
			Scraper: mockSceneScraper{
				errIDs: []int{errID1},
				results: map[int][]*scraper.ScrapedScene{
					found1ID: {{
						Title: &scrapedTitle,
					}},
				},
			},
		},
		{
			Scraper: mockSceneScraper{
				errIDs: []int{errID2},
				results: map[int][]*scraper.ScrapedScene{
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
		SetOrganized:             &boolFalse,
		SetCoverImage:            &boolFalse,
		IncludeMalePerformers:    &boolFalse,
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
					result: &scraper.ScrapedScene{},
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

	scrapedScene := &scraper.ScrapedScene{
		Title:   &scrapedTitle,
		Date:    &scrapedDate,
		Details: &scrapedDetails,
		URLs:    []string{scrapedURL},
	}

	scrapedUnchangedScene := &scraper.ScrapedScene{
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
		scraped      *scraper.ScrapedScene
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
