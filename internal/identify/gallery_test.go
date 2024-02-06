package identify

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/scraper"
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
		s := args.Get(1).(*models.Studio)
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
			tr.result = &galleryScrapeResult{
				source: GalleryScraperSource{
					RemoteSite: "endpoint",
				},
				result: &scraper.ScrapedGallery{
					Studio: tt.result,
				},
			}

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

	tr := galleryRelationships{
		fieldOptions: make(map[string]*FieldOptions),
	}

	tests := []struct {
		name         string
		gallery      *models.Gallery
		fieldOptions *FieldOptions
		scraped      []*models.ScrapedPerformer
		ignoreMale   bool
		want         []int
		wantErr      bool
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
			false,
			nil,
			false,
		},
		{
			"none",
			emptyGallery,
			defaultOptions,
			[]*models.ScrapedPerformer{},
			false,
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
			false,
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
			false,
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
			true,
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
			false,
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
			true,
			[]int{validStoredIDInt},
			false,
		},
		{
			"error getting tag ID",
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
			false,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.gallery = tt.gallery
			tr.fieldOptions["performers"] = tt.fieldOptions
			tr.result = &galleryScrapeResult{
				result: &scraper.ScrapedGallery{
					Performers: tt.scraped,
				},
			}

			got, err := tr.performers(testCtx, tt.ignoreMale)
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

	db.Tag.On("Create", testCtx, mock.MatchedBy(func(p *models.Tag) bool {
		return p.Name == validName
	})).Run(func(args mock.Arguments) {
		t := args.Get(1).(*models.Tag)
		t.ID = validStoredIDInt
	}).Return(nil)
	db.Tag.On("Create", testCtx, mock.MatchedBy(func(p *models.Tag) bool {
		return p.Name == invalidName
	})).Return(errors.New("error creating tag"))

	tr := galleryRelationships{
		tagCreator:   db.Tag,
		fieldOptions: make(map[string]*FieldOptions),
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
			tr.result = &galleryScrapeResult{
				result: &scraper.ScrapedGallery{
					Tags: tt.scraped,
				},
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
