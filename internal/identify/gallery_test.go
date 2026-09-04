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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_galleryRelationships_studio(t *testing.T) {
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
		s := args.Get(1).(*models.CreateStudioInput)
		s.ID = validStoredIDInt
	}).Return(nil)

	tr := galleryRelationships{
		studioReaderWriter: db.Studio,
		fieldOptions:       make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		gallery      *models.Gallery
		fieldOptions *FieldOptions
		result       *models.ScrapedStudio
		want         *int
		wantErr      bool
	}{
		{
			"nil studio",
			&models.Gallery{},
			defaultOptions,
			nil,
			nil,
			false,
		},
		{
			"ignore",
			&models.Gallery{},
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
			&models.Gallery{},
			defaultOptions,
			&models.ScrapedStudio{
				StoredID: &invalidStoredID,
			},
			nil,
			true,
		},
		{
			"same stored id",
			&models.Gallery{
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
			&models.Gallery{},
			defaultOptions,
			&models.ScrapedStudio{
				StoredID: &validStoredID,
			},
			&validStoredIDInt,
			false,
		},
		{
			"no create missing",
			&models.Gallery{},
			defaultOptions,
			&models.ScrapedStudio{},
			nil,
			false,
		},
		{
			"create missing",
			&models.Gallery{},
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
			tr.gallery = tt.gallery
			tr.fieldOptions["studio"] = tt.fieldOptions
			tr.scraped = &models.ScrapedGallery{
				Studio: tt.result,
			}
			tr.remoteSite = "endpoint"

			got, err := tr.studio(testCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryRelationships.studio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryRelationships.studio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_galleryRelationships_performers(t *testing.T) {
	const (
		galleryID = iota
		galleryWithPerformerID
		errGalleryID
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

	emptyGallery := &models.Gallery{
		ID:           galleryID,
		PerformerIDs: models.NewRelatedIDs([]int{}),
		TagIDs:       models.NewRelatedIDs([]int{}),
	}

	galleryWithPerformer := &models.Gallery{
		ID: galleryWithPerformerID,
		PerformerIDs: models.NewRelatedIDs([]int{
			existingPerformerID,
		}),
	}

	db := mocks.NewDatabase()

	tr := galleryRelationships{
		galleryReader: db.Gallery,
		fieldOptions:  make(map[string]*FieldOptions),
	}

	tests := []struct {
		name           string
		gallery        *models.Gallery
		fieldOptions   *FieldOptions
		scraped        []*models.ScrapedPerformer
		allowedGenders []models.GenderEnum
		want           []int
		wantErr        bool
	}{
		{
			"ignore",
			emptyGallery,
			&FieldOptions{
				Strategy: FieldStrategyIgnore,
			},
			[]*models.ScrapedPerformer{
				{
					StoredID: &validStoredID,
				},
			},
			nil,
			nil,
			false,
		},
		{
			"none",
			emptyGallery,
			defaultOptions,
			[]*models.ScrapedPerformer{},
			nil,
			nil,
			false,
		},
		{
			"merge existing",
			galleryWithPerformer,
			defaultOptions,
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &existingPerformerStr,
				},
			},
			nil,
			nil,
			false,
		},
		{
			"merge add",
			galleryWithPerformer,
			defaultOptions,
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &validStoredID,
				},
			},
			nil,
			[]int{existingPerformerID, validStoredIDInt},
			false,
		},
		{
			"ignore male",
			emptyGallery,
			defaultOptions,
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &validStoredID,
					Gender:   &male,
				},
			},
			[]models.GenderEnum{models.GenderEnumFemale, models.GenderEnumTransgenderMale, models.GenderEnumTransgenderFemale, models.GenderEnumIntersex, models.GenderEnumNonBinary},
			nil,
			false,
		},
		{
			"overwrite",
			galleryWithPerformer,
			&FieldOptions{
				Strategy: FieldStrategyOverwrite,
			},
			[]*models.ScrapedPerformer{
				{
					Name:     &validName,
					StoredID: &validStoredID,
				},
			},
			nil,
			[]int{validStoredIDInt},
			false,
		},
		{
			"ignore male (not male)",
			galleryWithPerformer,
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
			[]models.GenderEnum{models.GenderEnumFemale, models.GenderEnumTransgenderMale, models.GenderEnumTransgenderFemale, models.GenderEnumIntersex, models.GenderEnumNonBinary},
			[]int{validStoredIDInt},
			false,
		},
		{
			"error getting performer ID",
			emptyGallery,
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
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.gallery = tt.gallery
			tr.fieldOptions["performers"] = tt.fieldOptions
			tr.scraped = &models.ScrapedGallery{
				Performers: tt.scraped,
			}

			got, err := tr.performers(testCtx, tt.allowedGenders)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryRelationships.performers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryRelationships.performers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_galleryRelationships_tags(t *testing.T) {
	const (
		galleryID = iota
		galleryWithTagID
		errGalleryID
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

	emptyGallery := &models.Gallery{
		ID:           galleryID,
		TagIDs:       models.NewRelatedIDs([]int{}),
		PerformerIDs: models.NewRelatedIDs([]int{}),
	}

	galleryWithTag := &models.Gallery{
		ID: galleryWithTagID,
		TagIDs: models.NewRelatedIDs([]int{
			existingID,
		}),
		PerformerIDs: models.NewRelatedIDs([]int{}),
	}

	db := mocks.NewDatabase()

	db.Tag.On("Create", testCtx, mock.MatchedBy(func(p *models.CreateTagInput) bool {
		return p.Tag.Name == validName
	})).Run(func(args mock.Arguments) {
		t := args.Get(1).(*models.CreateTagInput)
		t.Tag.ID = validStoredIDInt
	}).Return(nil)
	db.Tag.On("Create", testCtx, mock.MatchedBy(func(p *models.CreateTagInput) bool {
		return p.Tag.Name == invalidName
	})).Return(errors.New("error creating tag"))

	tr := galleryRelationships{
		galleryReader: db.Gallery,
		tagCreator:    db.Tag,
		fieldOptions:  make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		gallery      *models.Gallery
		fieldOptions *FieldOptions
		scraped      []*models.ScrapedTag
		want         []int
		wantErr      bool
	}{
		{
			"ignore",
			emptyGallery,
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
			emptyGallery,
			defaultOptions,
			[]*models.ScrapedTag{},
			nil,
			false,
		},
		{
			"merge existing",
			galleryWithTag,
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
			galleryWithTag,
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
			galleryWithTag,
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
			emptyGallery,
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
			emptyGallery,
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
			emptyGallery,
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
			tr.gallery = tt.gallery
			tr.fieldOptions["tags"] = tt.fieldOptions
			tr.scraped = &models.ScrapedGallery{
				Tags: tt.scraped,
			}

			got, err := tr.tags(testCtx)
			if (err != nil) != tt.wantErr {
				t.Errorf("galleryRelationships.tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("galleryRelationships.tags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGalleryPartial(t *testing.T) {
	var (
		originalTitle        = "originalTitle"
		originalCode         = "originalCode"
		originalDate         = "2001-01-01"
		originalDetails      = "originalDetails"
		originalPhotographer = "originalPhotographer"
		originalURL          = "originalURL"
	)

	var (
		scrapedTitle        = "scrapedTitle"
		scrapedCode         = "scrapedCode"
		scrapedDate         = "2002-02-02"
		scrapedDetails      = "scrapedDetails"
		scrapedPhotographer = "scrapedPhotographer"
		scrapedURL          = "scrapedURL"
	)

	originalDateObj, _ := models.ParseDate(originalDate)
	scrapedDateObj, _ := models.ParseDate(scrapedDate)

	originalGallery := &models.Gallery{
		Title:        originalTitle,
		Code:         originalCode,
		Date:         &originalDateObj,
		Details:      originalDetails,
		Photographer: originalPhotographer,
		URLs:         models.NewRelatedStrings([]string{originalURL}),
	}

	organisedGallery := *originalGallery
	organisedGallery.Organized = true

	emptyGallery := &models.Gallery{
		URLs: models.NewRelatedStrings([]string{}),
	}

	postPartial := models.GalleryPartial{
		Title:        models.NewOptionalString(scrapedTitle),
		Code:         models.NewOptionalString(scrapedCode),
		Date:         models.NewOptionalDate(scrapedDateObj),
		Details:      models.NewOptionalString(scrapedDetails),
		Photographer: models.NewOptionalString(scrapedPhotographer),
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

	scrapedGallery := &models.ScrapedGallery{
		Title:        &scrapedTitle,
		Code:         &scrapedCode,
		Date:         &scrapedDate,
		Details:      &scrapedDetails,
		Photographer: &scrapedPhotographer,
		URLs:         []string{scrapedURL},
	}

	scrapedUnchangedGallery := &models.ScrapedGallery{
		Title:        &originalTitle,
		Code:         &originalCode,
		Date:         &originalDate,
		Details:      &originalDetails,
		Photographer: &originalPhotographer,
		URLs:         []string{originalURL},
	}

	makeFieldOptions := func(input *FieldOptions) map[string]*FieldOptions {
		return map[string]*FieldOptions{
			"title":        input,
			"code":         input,
			"date":         input,
			"details":      input,
			"photographer": input,
			"url":          input,
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
		scraped      *models.ScrapedGallery
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
			"overwrite url removal",
			args{
				originalGallery,
				&models.ScrapedGallery{
					URLs: []string{scrapedURL},
				},
				overwriteAll,
				false,
			},
			models.GalleryPartial{
				URLs: &models.UpdateStrings{
					Values: []string{scrapedURL},
					Mode:   models.RelationshipUpdateModeSet,
				},
			},
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

func Test_getGalleryUpdater_performerGenders(t *testing.T) {
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

	gallery := &models.Gallery{
		ID:           1,
		URLs:         models.NewRelatedStrings([]string{}),
		PerformerIDs: models.NewRelatedIDs([]int{}),
		TagIDs:       models.NewRelatedIDs([]int{}),
	}

	scrapedGallery := &models.ScrapedGallery{
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
			name:             "nil PerformerGenders includes all genders",
			performerGenders: nil,
			wantPerformerIDs: []int{femalePerformerID, malePerformerID},
		},
		{
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
			name:                  "empty PerformerGenders falls back to IncludeMalePerformers=false",
			performerGenders:      []models.GenderEnum{},
			includeMalePerformers: &boolFalse,
			wantPerformerIDs:      []int{femalePerformerID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := MetadataOptions{
				PerformerGenders:      tt.performerGenders,
				IncludeMalePerformers: tt.includeMalePerformers,
			}

			identifier := GalleryIdentifier{
				TxnManager:           db,
				GalleryReaderUpdater: db.Gallery,
				StudioReaderWriter:   db.Studio,
				PerformerCreator:     db.Performer,
				TagFinderCreator:     db.Tag,
				DefaultOptions:       &opts,
			}

			result := &galleryScrapeResult{
				source: GalleryScraperSource{},
				result: scrapedGallery,
			}

			updater, err := identifier.getGalleryUpdater(testCtx, gallery, result)
			assert.NoError(t, err)
			assert.NotNil(t, updater.Partial.PerformerIDs, "expected PerformerIDs to be set")
			assert.ElementsMatch(t, tt.wantPerformerIDs, updater.Partial.PerformerIDs.IDs)
		})
	}
}

func Test_GalleryIdentifier_Identify(t *testing.T) {
	const (
		galleryID         = 1
		storedStudioID    = 1
		storedPerformerID = 1
		storedTagID       = 1
	)
	validStoredStudioID := strconv.Itoa(storedStudioID)
	validStoredPerformerID := strconv.Itoa(storedPerformerID)
	validStoredTagID := strconv.Itoa(storedTagID)
	scrapedTitle := "scrapedTitle"
	scrapedCode := "scrapedCode"
	scrapedDate := "2002-02-02"
	scrapedURL := "https://example.com"

	boolTrue := true
	defaultOptions := &MetadataOptions{
		SetOrganized: &boolTrue,
	}

	source := GalleryScraperSource{
		Name: "test-gallery-scraper",
		Scraper: mockGalleryScraper{results: map[int][]*models.ScrapedGallery{
			galleryID: {{
				Title:      &scrapedTitle,
				Code:       &scrapedCode,
				Date:       &scrapedDate,
				URLs:       []string{scrapedURL},
				Studio:     &models.ScrapedStudio{StoredID: &validStoredStudioID},
				Performers: []*models.ScrapedPerformer{{StoredID: &validStoredPerformerID}},
				Tags:       []*models.ScrapedTag{{StoredID: &validStoredTagID}},
			}},
		}},
		RemoteSite: "test-endpoint",
	}

	db := mocks.NewDatabase()

	gallery := &models.Gallery{
		ID:       galleryID,
		Title:    "originalTitle",
		Code:     "originalCode",
		StudioID: nil,
	}

	db.Gallery.On("Find", mock.Anything, galleryID).Return(gallery, nil)
	db.Gallery.On("GetURLs", mock.Anything, galleryID).Return([]string{}, nil)
	db.Gallery.On("GetPerformerIDs", mock.Anything, galleryID).Return([]int{}, nil)
	db.Gallery.On("GetTagIDs", mock.Anything, galleryID).Return([]int{}, nil)
	db.Gallery.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(id int) bool { return id == galleryID }), mock.Anything).Return(&models.Gallery{ID: galleryID, Title: scrapedTitle}, nil)

	identifier := GalleryIdentifier{
		TxnManager:           db,
		GalleryReaderUpdater: db.Gallery,
		StudioReaderWriter:   db.Studio,
		PerformerCreator:     db.Performer,
		TagFinderCreator:     db.Tag,
		DefaultOptions:       defaultOptions,
		Sources:              []GalleryScraperSource{source},
		PostHookExecutor:     mockHookExecutor{},
	}

	err := identifier.Identify(testCtx, gallery)
	assert.NoError(t, err)

	// Verify UpdatePartial was called (gallery was modified)
	db.Gallery.AssertNumberOfCalls(t, "UpdatePartial", 1)
	// Verify scraped data was processed (studio, performer, tag all resolved via StoredID)
	db.Gallery.AssertNumberOfCalls(t, "GetURLs", 1)
	db.Gallery.AssertNumberOfCalls(t, "GetPerformerIDs", 1)
	db.Gallery.AssertNumberOfCalls(t, "GetTagIDs", 1)
}

func Test_getGalleryUpdater_errorsAndSourceOptions(t *testing.T) {
	boolTrue := true
	createMissing := true
	singleName := "Single"
	validTagID := "2"
	skipTagID := 3
	skipTagIDStr := strconv.Itoa(skipTagID)
	invalidID := "invalid"

	tests := []struct {
		name    string
		options *MetadataOptions
		scraped *models.ScrapedGallery
		wantIDs []int
		wantErr bool
	}{
		{
			name: "source options with skipped performer tag",
			options: &MetadataOptions{
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
			scraped: &models.ScrapedGallery{
				Performers: []*models.ScrapedPerformer{{Name: &singleName}},
				Tags:       []*models.ScrapedTag{{StoredID: &validTagID}},
			},
			wantIDs: []int{2, skipTagID},
		},
		{
			name: "invalid skipped performer tag",
			options: &MetadataOptions{
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
			scraped: &models.ScrapedGallery{
				Performers: []*models.ScrapedPerformer{{Name: &singleName}},
			},
			wantErr: true,
		},
		{
			name: "skipped performer without tag",
			options: &MetadataOptions{
				FieldOptions: []*FieldOptions{
					{
						Field:         "performers",
						Strategy:      FieldStrategyMerge,
						CreateMissing: &createMissing,
					},
				},
				SkipSingleNamePerformers: &boolTrue,
			},
			scraped: &models.ScrapedGallery{
				Performers: []*models.ScrapedPerformer{{Name: &singleName}},
			},
		},
		{
			name:    "studio error",
			scraped: &models.ScrapedGallery{Studio: &models.ScrapedStudio{StoredID: &invalidID}},
			wantErr: true,
		},
		{
			name:    "performer error",
			scraped: &models.ScrapedGallery{Performers: []*models.ScrapedPerformer{{StoredID: &invalidID}}},
			wantErr: true,
		},
		{
			name:    "tag error",
			scraped: &models.ScrapedGallery{Tags: []*models.ScrapedTag{{StoredID: &invalidID}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()
			identifier := GalleryIdentifier{
				GalleryReaderUpdater: db.Gallery,
				StudioReaderWriter:   db.Studio,
				PerformerCreator:     db.Performer,
				TagFinderCreator:     db.Tag,
			}
			gallery := &models.Gallery{
				ID:           1,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			}
			result := &galleryScrapeResult{
				source: GalleryScraperSource{Options: tt.options},
				result: tt.scraped,
			}

			updater, err := identifier.getGalleryUpdater(testCtx, gallery, result)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, updater)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updater)
				if tt.wantIDs == nil {
					assert.Nil(t, updater.Partial.TagIDs)
				} else {
					assert.NotNil(t, updater.Partial.TagIDs)
					assert.Equal(t, tt.wantIDs, updater.Partial.TagIDs.IDs)
				}
			}
		})
	}
}

func Test_GalleryIdentifier_addTagToGallery(t *testing.T) {
	const (
		galleryID = 1
		tagID     = 2
	)
	tagIDStr := strconv.Itoa(tagID)
	tagName := "Skipped"

	tests := []struct {
		name     string
		tagToAdd string
		gallery  *models.Gallery
		setup    func(db *mocks.Database)
		wantErr  bool
	}{
		{
			name:     "invalid tag ID",
			tagToAdd: "invalid",
			gallery:  &models.Gallery{ID: galleryID},
			wantErr:  true,
		},
		{
			name:     "load tag IDs error",
			tagToAdd: tagIDStr,
			gallery:  &models.Gallery{ID: galleryID},
			setup: func(db *mocks.Database) {
				db.Gallery.On("GetTagIDs", mock.Anything, galleryID).Return(nil, errors.New("load tag IDs error"))
			},
			wantErr: true,
		},
		{
			name:     "tag already present",
			tagToAdd: tagIDStr,
			gallery: &models.Gallery{
				ID:     galleryID,
				TagIDs: models.NewRelatedIDs([]int{tagID}),
			},
		},
		{
			name:     "add tag error",
			tagToAdd: tagIDStr,
			gallery: &models.Gallery{
				ID:     galleryID,
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Gallery.On("UpdatePartial", mock.Anything, galleryID, mock.AnythingOfType("models.GalleryPartial")).Return(nil, errors.New("add tag error"))
			},
			wantErr: true,
		},
		{
			name:     "find tag error",
			tagToAdd: tagIDStr,
			gallery: &models.Gallery{
				ID:     galleryID,
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Gallery.On("UpdatePartial", mock.Anything, galleryID, mock.AnythingOfType("models.GalleryPartial")).Return(&models.Gallery{ID: galleryID}, nil)
				db.Tag.On("Find", mock.Anything, tagID).Return(nil, errors.New("find tag error"))
			},
		},
		{
			name:     "find tag nil",
			tagToAdd: tagIDStr,
			gallery: &models.Gallery{
				ID:     galleryID,
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Gallery.On("UpdatePartial", mock.Anything, galleryID, mock.AnythingOfType("models.GalleryPartial")).Return(&models.Gallery{ID: galleryID}, nil)
				db.Tag.On("Find", mock.Anything, tagID).Return(nil, nil)
			},
		},
		{
			name:     "success",
			tagToAdd: tagIDStr,
			gallery: &models.Gallery{
				ID:     galleryID,
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Gallery.On("UpdatePartial", mock.Anything, galleryID, mock.AnythingOfType("models.GalleryPartial")).Return(&models.Gallery{ID: galleryID}, nil)
				db.Tag.On("Find", mock.Anything, tagID).Return(&models.Tag{ID: tagID, Name: tagName}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()
			if tt.setup != nil {
				tt.setup(db)
			}

			identifier := GalleryIdentifier{
				TxnManager:           db,
				GalleryReaderUpdater: db.Gallery,
				TagFinderCreator:     db.Tag,
			}

			err := identifier.addTagToGallery(testCtx, tt.gallery, tt.tagToAdd)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_GalleryIdentifier_modifyGallery(t *testing.T) {
	const galleryID = 1
	scrapedTitle := "scrapedTitle"
	scrapedCode := "scrapedCode"
	invalidID := "invalid"

	tests := []struct {
		name     string
		gallery  *models.Gallery
		scraped  *models.ScrapedGallery
		executor GalleryUpdatePostHookExecutor
		setup    func(db *mocks.Database)
		wantErr  bool
	}{
		{
			name: "load URLs error",
			gallery: &models.Gallery{
				ID:           galleryID,
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Gallery.On("GetURLs", mock.Anything, galleryID).Return(nil, errors.New("load URLs error"))
			},
			wantErr: true,
		},
		{
			name: "load performers error",
			gallery: &models.Gallery{
				ID:     galleryID,
				URLs:   models.NewRelatedStrings([]string{}),
				TagIDs: models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Gallery.On("GetPerformerIDs", mock.Anything, galleryID).Return(nil, errors.New("load performers error"))
			},
			wantErr: true,
		},
		{
			name: "load tags error",
			gallery: &models.Gallery{
				ID:           galleryID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
			},
			setup: func(db *mocks.Database) {
				db.Gallery.On("GetTagIDs", mock.Anything, galleryID).Return(nil, errors.New("load tags error"))
			},
			wantErr: true,
		},
		{
			name: "updater error",
			gallery: &models.Gallery{
				ID:           galleryID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			},
			scraped: &models.ScrapedGallery{
				Studio: &models.ScrapedStudio{StoredID: &invalidID},
			},
			wantErr: true,
		},
		{
			name: "empty update",
			gallery: &models.Gallery{
				ID:           galleryID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			},
			scraped: &models.ScrapedGallery{},
		},
		{
			name: "update error",
			gallery: &models.Gallery{
				ID:           galleryID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			},
			scraped: &models.ScrapedGallery{Title: &scrapedTitle},
			setup: func(db *mocks.Database) {
				db.Gallery.On("UpdatePartial", mock.Anything, galleryID, mock.AnythingOfType("models.GalleryPartial")).Return(nil, errors.New("update error"))
			},
			wantErr: true,
		},
		{
			name: "success with title",
			gallery: &models.Gallery{
				ID:           galleryID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			},
			scraped: &models.ScrapedGallery{Title: &scrapedTitle},
			setup: func(db *mocks.Database) {
				db.Gallery.On("UpdatePartial", mock.Anything, galleryID, mock.AnythingOfType("models.GalleryPartial")).Return(&models.Gallery{ID: galleryID}, nil)
			},
		},
		{
			name: "success without title",
			gallery: &models.Gallery{
				ID:           galleryID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			},
			scraped:  &models.ScrapedGallery{Code: &scrapedCode},
			executor: mockHookExecutor{},
			setup: func(db *mocks.Database) {
				db.Gallery.On("UpdatePartial", mock.Anything, galleryID, mock.AnythingOfType("models.GalleryPartial")).Return(&models.Gallery{ID: galleryID}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDatabase()
			if tt.setup != nil {
				tt.setup(db)
			}

			identifier := GalleryIdentifier{
				TxnManager:           db,
				GalleryReaderUpdater: db.Gallery,
				StudioReaderWriter:   db.Studio,
				PerformerCreator:     db.Performer,
				TagFinderCreator:     db.Tag,
				PostHookExecutor:     tt.executor,
			}
			result := &galleryScrapeResult{
				source: GalleryScraperSource{},
				result: tt.scraped,
			}

			err := identifier.modifyGallery(testCtx, tt.gallery, result)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_GalleryIdentifier_Identify_paths(t *testing.T) {
	const (
		errID = iota
		missingID
		foundID
		errUpdateID
		multiID
		multiTagID
		multiInvalidTagID
	)

	var (
		skipTagID     = 1
		skipTagIDStr  = strconv.Itoa(skipTagID)
		scrapedTitle  = "scrapedTitle"
		scrapedTitle2 = "scrapedTitle2"
		invalidTagID  = "invalid"
		boolTrue      = true
	)

	skipMultiple := &MetadataOptions{SkipMultipleMatches: &boolTrue}
	sources := []GalleryScraperSource{
		{
			Name:    "test-gallery-scraper",
			Options: skipMultiple,
			Scraper: mockGalleryScraper{
				errIDs: []int{errID},
				results: map[int][]*models.ScrapedGallery{
					foundID: {{
						Title: &scrapedTitle,
					}},
					errUpdateID: {{
						Title: &scrapedTitle,
					}},
					multiID: {
						{Title: &scrapedTitle},
						{Title: &scrapedTitle2},
					},
					multiTagID: {
						{Title: &scrapedTitle},
						{Title: &scrapedTitle2},
					},
					multiInvalidTagID: {
						{Title: &scrapedTitle},
						{Title: &scrapedTitle2},
					},
				},
			},
		},
	}

	db := mocks.NewDatabase()
	db.Gallery.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(id int) bool {
		return id == errUpdateID
	}), mock.AnythingOfType("models.GalleryPartial")).Return(nil, errors.New("update error"))
	db.Gallery.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(id int) bool {
		return id != errUpdateID
	}), mock.AnythingOfType("models.GalleryPartial")).Return(&models.Gallery{ID: foundID}, nil)
	db.Tag.On("Find", mock.Anything, skipTagID).Return(&models.Tag{
		ID:   skipTagID,
		Name: skipTagIDStr,
	}, nil)

	tests := []struct {
		name      string
		galleryID int
		options   *MetadataOptions
		wantErr   bool
	}{
		{
			name:      "error scraping",
			galleryID: errID,
		},
		{
			name:      "not found",
			galleryID: missingID,
		},
		{
			name:      "found",
			galleryID: foundID,
		},
		{
			name:      "error modifying",
			galleryID: errUpdateID,
			wantErr:   true,
		},
		{
			name:      "multiple found",
			galleryID: multiID,
		},
		{
			name:      "multiple found with tag",
			galleryID: multiTagID,
			options: &MetadataOptions{
				SkipMultipleMatchTag: &skipTagIDStr,
			},
		},
		{
			name:      "multiple found with invalid tag",
			galleryID: multiInvalidTagID,
			options: &MetadataOptions{
				SkipMultipleMatchTag: &invalidTagID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identifier := GalleryIdentifier{
				TxnManager:           db,
				GalleryReaderUpdater: db.Gallery,
				StudioReaderWriter:   db.Studio,
				PerformerCreator:     db.Performer,
				TagFinderCreator:     db.Tag,
				Sources:              sources,
				PostHookExecutor:     mockHookExecutor{},
			}
			if tt.options != nil {
				identifier.DefaultOptions = tt.options
			}

			gallery := &models.Gallery{
				ID:           tt.galleryID,
				URLs:         models.NewRelatedStrings([]string{}),
				PerformerIDs: models.NewRelatedIDs([]int{}),
				TagIDs:       models.NewRelatedIDs([]int{}),
			}

			err := identifier.Identify(testCtx, gallery)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// mockGalleryScraper implements GalleryScraper for testing.
type mockGalleryScraper struct {
	errIDs  []int
	results map[int][]*models.ScrapedGallery
}

func (m mockGalleryScraper) ScrapeGalleries(ctx context.Context, galleryID int) ([]*models.ScrapedGallery, error) {
	if slices.Contains(m.errIDs, galleryID) {
		return nil, errors.New("scrape gallery error")
	}
	return m.results[galleryID], nil
}
