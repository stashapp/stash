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
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGalleryScraper struct {
	errIDs  []int
	results map[int][]*scraper.ScrapedGallery
}

func (s mockGalleryScraper) ScrapeGalleries(ctx context.Context, galleryID int) ([]*scraper.ScrapedGallery, error) {
	if sliceutil.Contains(s.errIDs, galleryID) {
		return nil, errors.New("scrape gallery error")
	}
	return s.results[galleryID], nil
}

func (s mockHookExecutor) ExecuteGalleryUpdatePostHooks(ctx context.Context, input models.GalleryUpdateInput, inputFields []string) {
}

func TestGalleryIdentifier_Identify(t *testing.T) {
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

	defaultOptions := &GalleryMetadataOptions{
		SetOrganized:             &boolFalse,
		IncludeMalePerformers:    &boolFalse,
		SkipSingleNamePerformers: &boolFalse,
	}
	sources := []GalleryScraperSource{
		{
			Scraper: mockGalleryScraper{
				errIDs: []int{errID1},
				results: map[int][]*scraper.ScrapedGallery{
					found1ID: {{
						Title: &scrapedTitle,
					}},
				},
			},
		},
		{
			Scraper: mockGalleryScraper{
				errIDs: []int{errID2},
				results: map[int][]*scraper.ScrapedGallery{
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

	db.Gallery.On("GetURLs", mock.Anything, mock.Anything).Return(nil, nil)
	db.Gallery.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(id int) bool {
		return id == errUpdateID
	}), mock.Anything).Return(nil, errors.New("update error"))
	db.Gallery.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(id int) bool {
		return id != errUpdateID
	}), mock.Anything).Return(nil, nil)

	db.Tag.On("Find", mock.Anything, skipMultipleTagID).Return(&models.Tag{
		ID:   skipMultipleTagID,
		Name: skipMultipleTagIDStr,
	}, nil)

	tests := []struct {
		name      string
		galleryID int
		options   *GalleryMetadataOptions
		wantErr   bool
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
			&GalleryMetadataOptions{
				SkipMultipleMatches:  &boolTrue,
				SkipMultipleMatchTag: &skipMultipleTagIDStr,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identifier := GalleryIdentifier{
				TxnManager:                    db,
				GalleryReaderUpdater:          db.Gallery,
				StudioReaderWriter:            db.Studio,
				PerformerCreator:              db.Performer,
				TagFinderCreator:              db.Tag,
				DefaultOptions:                defaultOptions,
				Sources:                       sources,
				GalleryUpdatePostHookExecutor: mockHookExecutor{},
			}

			if tt.options != nil {
				identifier.DefaultOptions = tt.options
			}

			gallery := &models.Gallery{
				ID:           tt.galleryID,
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			}
			if err := identifier.Identify(testCtx, gallery); (err != nil) != tt.wantErr {
				t.Errorf("GalleryIdentifier.Identify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGalleryIdentifier_modifyGallery(t *testing.T) {
	db := mocks.NewDatabase()

	boolFalse := false
	defaultOptions := &GalleryMetadataOptions{
		SetOrganized:             &boolFalse,
		IncludeMalePerformers:    &boolFalse,
		SkipSingleNamePerformers: &boolFalse,
	}
	tr := &GalleryIdentifier{
		TxnManager:           db,
		GalleryReaderUpdater: db.Gallery,
		StudioReaderWriter:   db.Studio,
		PerformerCreator:     db.Performer,
		TagFinderCreator:     db.Tag,
		DefaultOptions:       defaultOptions,
	}

	type args struct {
		gallery *models.Gallery
		result  *galleryScrapeResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty update",
			args{
				&models.Gallery{
					URLs:         models.NewRelatedStrings([]string{}),
					PerformerIDs: models.NewRelatedIDs([]int{}),
					TagIDs:       models.NewRelatedIDs([]int{}),
				},
				&galleryScrapeResult{
					result: &scraper.ScrapedGallery{},
					source: GalleryScraperSource{
						Options: defaultOptions,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tr.modifyGallery(testCtx, tt.args.gallery, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("GalleryIdentifier.modifyGallery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getGalleryFieldOptions(t *testing.T) {
	const (
		inFirst  = "inFirst"
		inSecond = "inSecond"
		inBoth   = "inBoth"
	)

	type args struct {
		options []GalleryMetadataOptions
	}
	tests := []struct {
		name string
		args args
		want map[string]*FieldOptions
	}{
		{
			"simple",
			args{
				[]GalleryMetadataOptions{
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
			if got := getFieldOptionsGallery(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGalleryPartial(t *testing.T) {
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

	originalGallery := &models.Gallery{
		Title:   originalTitle,
		Date:    &originalDateObj,
		Details: originalDetails,
		URLs:    models.NewRelatedStrings([]string{originalURL}),
	}

	organisedGallery := *originalGallery
	organisedGallery.Organized = true

	emptyGallery := &models.Gallery{
		URLs: models.NewRelatedStrings([]string{}),
	}

	postPartial := models.GalleryPartial{
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

	scrapedGallery := &scraper.ScrapedGallery{
		Title:   &scrapedTitle,
		Date:    &scrapedDate,
		Details: &scrapedDetails,
		URLs:    []string{scrapedURL},
	}

	scrapedUnchangedGallery := &scraper.ScrapedGallery{
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
		gallery      *models.Gallery
		scraped      *scraper.ScrapedGallery
		fieldOptions map[string]*FieldOptions
		setOrganized bool
	}
	tests := []struct {
		name string
		args args
		want models.GalleryPartial
	}{
		{
			"overwrite all",
			args{
				originalGallery,
				scrapedGallery,
				overwriteAll,
				false,
			},
			postPartial,
		},
		{
			"ignore all",
			args{
				originalGallery,
				scrapedGallery,
				ignoreAll,
				false,
			},
			models.GalleryPartial{},
		},
		{
			"merge (existing values)",
			args{
				originalGallery,
				scrapedGallery,
				mergeAll,
				false,
			},
			models.GalleryPartial{
				URLs: &models.UpdateStrings{
					Values: []string{originalURL, scrapedURL},
					Mode:   models.RelationshipUpdateModeSet,
				},
			},
		},
		{
			"merge (empty values)",
			args{
				emptyGallery,
				scrapedGallery,
				mergeAll,
				false,
			},
			postPartialMerge,
		},
		{
			"unchanged",
			args{
				originalGallery,
				scrapedUnchangedGallery,
				overwriteAll,
				false,
			},
			models.GalleryPartial{},
		},
		{
			"set organized",
			args{
				originalGallery,
				scrapedUnchangedGallery,
				overwriteAll,
				true,
			},
			models.GalleryPartial{
				Organized: models.NewOptionalBool(setOrganised),
			},
		},
		{
			"set organized unchanged",
			args{
				&organisedGallery,
				scrapedUnchangedGallery,
				overwriteAll,
				true,
			},
			models.GalleryPartial{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getGalleryPartial(tt.args.gallery, tt.args.scraped, tt.args.fieldOptions, tt.args.setOrganized)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_shouldGallerySetSingleValueField(t *testing.T) {
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
