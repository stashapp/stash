package identify

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stretchr/testify/mock"
)

type mockSceneScraper struct {
	errIDs  []int
	results map[int]*models.ScrapedScene
}

func (s mockSceneScraper) ScrapeScene(ctx context.Context, sceneID int) (*models.ScrapedScene, error) {
	if intslice.IntInclude(s.errIDs, sceneID) {
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
		errUpdateID
	)

	var scrapedTitle = "scrapedTitle"

	defaultOptions := &models.IdentifyMetadataOptionsInput{}
	sources := []ScraperSource{
		{
			Scraper: mockSceneScraper{
				errIDs: []int{errID1},
				results: map[int]*models.ScrapedScene{
					found1ID: {
						Title: &scrapedTitle,
					},
				},
			},
		},
		{
			Scraper: mockSceneScraper{
				errIDs: []int{errID2},
				results: map[int]*models.ScrapedScene{
					found2ID: {
						Title: &scrapedTitle,
					},
					errUpdateID: {
						Title: &scrapedTitle,
					},
				},
			},
		},
	}

	repo := mocks.NewTransactionManager()
	repo.Scene().(*mocks.SceneReaderWriter).On("Update", mock.MatchedBy(func(partial models.ScenePartial) bool {
		return partial.ID != errUpdateID
	})).Return(nil, nil)
	repo.Scene().(*mocks.SceneReaderWriter).On("Update", mock.MatchedBy(func(partial models.ScenePartial) bool {
		return partial.ID == errUpdateID
	})).Return(nil, errors.New("update error"))

	tests := []struct {
		name    string
		sceneID int
		wantErr bool
	}{
		{
			"error scraping",
			errID1,
			false,
		},
		{
			"error scraping from second",
			errID2,
			false,
		},
		{
			"found in first scraper",
			found1ID,
			false,
		},
		{
			"found in second scraper",
			found2ID,
			false,
		},
		{
			"not found",
			missingID,
			false,
		},
		{
			"error modifying",
			errUpdateID,
			true,
		},
	}

	identifier := SceneIdentifier{
		DefaultOptions:              defaultOptions,
		Sources:                     sources,
		SceneUpdatePostHookExecutor: mockHookExecutor{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene := &models.Scene{
				ID: tt.sceneID,
			}
			if err := identifier.Identify(context.TODO(), repo, scene); (err != nil) != tt.wantErr {
				t.Errorf("SceneIdentifier.Identify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSceneIdentifier_modifyScene(t *testing.T) {
	repo := mocks.NewTransactionManager()
	tr := &SceneIdentifier{}

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
				&models.Scene{},
				&scrapeResult{
					result: &models.ScrapedScene{},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tr.modifyScene(context.TODO(), repo, tt.args.scene, tt.args.result); (err != nil) != tt.wantErr {
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
		options []models.IdentifyMetadataOptionsInput
	}
	tests := []struct {
		name string
		args args
		want map[string]*models.IdentifyFieldOptionsInput
	}{
		{
			"simple",
			args{
				[]models.IdentifyMetadataOptionsInput{
					{
						FieldOptions: []*models.IdentifyFieldOptionsInput{
							{
								Field:    inFirst,
								Strategy: models.IdentifyFieldStrategyIgnore,
							},
							{
								Field:    inBoth,
								Strategy: models.IdentifyFieldStrategyIgnore,
							},
						},
					},
					{
						FieldOptions: []*models.IdentifyFieldOptionsInput{
							{
								Field:    inSecond,
								Strategy: models.IdentifyFieldStrategyMerge,
							},
							{
								Field:    inBoth,
								Strategy: models.IdentifyFieldStrategyMerge,
							},
						},
					},
				},
			},
			map[string]*models.IdentifyFieldOptionsInput{
				inFirst: {
					Field:    inFirst,
					Strategy: models.IdentifyFieldStrategyIgnore,
				},
				inSecond: {
					Field:    inSecond,
					Strategy: models.IdentifyFieldStrategyMerge,
				},
				inBoth: {
					Field:    inBoth,
					Strategy: models.IdentifyFieldStrategyIgnore,
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
		originalDate    = "originalDate"
		originalDetails = "originalDetails"
		originalURL     = "originalURL"
	)

	var (
		scrapedTitle   = "scrapedTitle"
		scrapedDate    = "scrapedDate"
		scrapedDetails = "scrapedDetails"
		scrapedURL     = "scrapedURL"
	)

	originalScene := &models.Scene{
		Title: models.NullString(originalTitle),
		Date: models.SQLiteDate{
			String: originalDate,
			Valid:  true,
		},
		Details: models.NullString(originalDetails),
		URL:     models.NullString(originalURL),
	}

	organisedScene := *originalScene
	organisedScene.Organized = true

	emptyScene := &models.Scene{}

	postPartial := models.ScenePartial{
		Title: models.NullStringPtr(scrapedTitle),
		Date: &models.SQLiteDate{
			String: scrapedDate,
			Valid:  true,
		},
		Details: models.NullStringPtr(scrapedDetails),
		URL:     models.NullStringPtr(scrapedURL),
	}

	scrapedScene := &models.ScrapedScene{
		Title:   &scrapedTitle,
		Date:    &scrapedDate,
		Details: &scrapedDetails,
		URL:     &scrapedURL,
	}

	scrapedUnchangedScene := &models.ScrapedScene{
		Title:   &originalTitle,
		Date:    &originalDate,
		Details: &originalDetails,
		URL:     &originalURL,
	}

	makeFieldOptions := func(input *models.IdentifyFieldOptionsInput) map[string]*models.IdentifyFieldOptionsInput {
		return map[string]*models.IdentifyFieldOptionsInput{
			"title":   input,
			"date":    input,
			"details": input,
			"url":     input,
		}
	}

	overwriteAll := makeFieldOptions(&models.IdentifyFieldOptionsInput{
		Strategy: models.IdentifyFieldStrategyOverwrite,
	})
	ignoreAll := makeFieldOptions(&models.IdentifyFieldOptionsInput{
		Strategy: models.IdentifyFieldStrategyIgnore,
	})
	mergeAll := makeFieldOptions(&models.IdentifyFieldOptionsInput{
		Strategy: models.IdentifyFieldStrategyMerge,
	})

	setOrganised := true

	type args struct {
		scene        *models.Scene
		scraped      *models.ScrapedScene
		fieldOptions map[string]*models.IdentifyFieldOptionsInput
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
			models.ScenePartial{},
		},
		{
			"merge (empty values)",
			args{
				emptyScene,
				scrapedScene,
				mergeAll,
				false,
			},
			postPartial,
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
				Organized: &setOrganised,
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
			if got := getScenePartial(tt.args.scene, tt.args.scraped, tt.args.fieldOptions, tt.args.setOrganized); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScenePartial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shouldSetSingleValueField(t *testing.T) {
	const invalid = "invalid"

	type args struct {
		strategy         *models.IdentifyFieldOptionsInput
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
				&models.IdentifyFieldOptionsInput{
					Strategy: models.IdentifyFieldStrategyIgnore,
				},
				false,
			},
			false,
		},
		{
			"merge existing",
			args{
				&models.IdentifyFieldOptionsInput{
					Strategy: models.IdentifyFieldStrategyMerge,
				},
				true,
			},
			false,
		},
		{
			"merge absent",
			args{
				&models.IdentifyFieldOptionsInput{
					Strategy: models.IdentifyFieldStrategyMerge,
				},
				false,
			},
			true,
		},
		{
			"overwrite",
			args{
				&models.IdentifyFieldOptionsInput{
					Strategy: models.IdentifyFieldStrategyOverwrite,
				},
				true,
			},
			true,
		},
		{
			"nil (merge) existing",
			args{
				&models.IdentifyFieldOptionsInput{},
				true,
			},
			false,
		},
		{
			"nil (merge) absent",
			args{
				&models.IdentifyFieldOptionsInput{},
				false,
			},
			true,
		},
		{
			"invalid (merge) existing",
			args{
				&models.IdentifyFieldOptionsInput{
					Strategy: invalid,
				},
				true,
			},
			false,
		},
		{
			"invalid (merge) absent",
			args{
				&models.IdentifyFieldOptionsInput{
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
