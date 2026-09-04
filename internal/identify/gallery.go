package identify

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type GalleryReaderUpdater interface {
	models.GalleryUpdater
	models.PerformerIDLoader
	models.TagIDLoader
	models.URLLoader
}

type galleryRelationships struct {
	galleryReader            GalleryReaderUpdater
	studioReaderWriter       models.StudioReaderWriter
	performerCreator         PerformerCreator
	tagCreator               models.TagCreator
	gallery                  *models.Gallery
	scraped                  *models.ScrapedGallery
	remoteSite               string
	fieldOptions             map[string]*FieldOptions
	skipSingleNamePerformers bool
}

func (r galleryRelationships) studio(ctx context.Context) (*int, error) {
	existingID := r.gallery.StudioID
	fieldStrategy := r.fieldOptions["studio"]
	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)

	scraped := r.scraped.Studio
	endpoint := r.remoteSite

	if scraped == nil || !shouldSetSingleValueField(fieldStrategy, existingID != nil) {
		return nil, nil
	}

	if scraped.StoredID != nil {
		// existing studio, just set it
		studioID, err := parseStudioID(*scraped.StoredID)
		if err != nil {
			return nil, err
		}

		// only return value if different to current
		if existingID == nil || *existingID != studioID {
			return &studioID, nil
		}
	} else if createMissing {
		return createMissingStudio(ctx, endpoint, r.studioReaderWriter, scraped)
	}

	return nil, nil
}

func (r galleryRelationships) performers(ctx context.Context, allowedGenders []models.GenderEnum) ([]int, error) {
	fieldStrategy := r.fieldOptions["performers"]
	scraped := r.scraped.Performers

	// just check if ignored
	if len(scraped) == 0 || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	endpoint := r.remoteSite

	var performerIDs []int
	originalPerformerIDs := r.gallery.PerformerIDs.List()

	if strategy == FieldStrategyMerge {
		// add to existing
		performerIDs = originalPerformerIDs
	}

	singleNamePerformerSkipped := false

	for _, p := range scraped {
		if performerGenderExcluded(allowedGenders, p.Gender) {
			continue
		}

		performerID, err := getPerformerID(ctx, endpoint, r.performerCreator, p, createMissing, r.skipSingleNamePerformers)
		if err != nil {
			if errors.Is(err, ErrSkipSingleNamePerformer) {
				singleNamePerformerSkipped = true
				continue
			}
			return nil, err
		}

		if performerID != nil {
			performerIDs = sliceutil.AppendUnique(performerIDs, *performerID)
		}
	}

	// don't return if nothing was added
	if sliceutil.SliceSame(originalPerformerIDs, performerIDs) {
		if singleNamePerformerSkipped {
			return nil, ErrSkipSingleNamePerformer
		}
		return nil, nil
	}

	if singleNamePerformerSkipped {
		return performerIDs, ErrSkipSingleNamePerformer
	}
	return performerIDs, nil
}

func (r galleryRelationships) tags(ctx context.Context) ([]int, error) {
	fieldStrategy := r.fieldOptions["tags"]
	scraped := r.scraped.Tags
	target := r.gallery

	// just check if ignored
	if len(scraped) == 0 || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var tagIDs []int
	originalTagIDs := target.TagIDs.List()

	if strategy == FieldStrategyMerge {
		// add to existing
		tagIDs = originalTagIDs
	}

	endpoint := r.remoteSite

	for _, t := range scraped {
		if t.StoredID != nil {
			// existing tag, just add it
			tagID, err := parseTagID(*t.StoredID)
			if err != nil {
				return nil, err
			}

			tagIDs = sliceutil.AppendUnique(tagIDs, tagID)
		} else if createMissing {
			newTag := t.ToTag(endpoint, nil)

			err := r.tagCreator.Create(ctx, &models.CreateTagInput{
				Tag: newTag,
			})
			if err != nil {
				return nil, fmt.Errorf("error creating tag: %w", err)
			}

			tagIDs = append(tagIDs, newTag.ID)
		}
	}

	// don't return if nothing was added
	if sliceutil.SliceSame(originalTagIDs, tagIDs) {
		return nil, nil
	}

	return tagIDs, nil
}

type GalleryScraper interface {
	ScrapeGalleries(ctx context.Context, galleryID int) ([]*models.ScrapedGallery, error)
}

type GalleryScraperSource struct {
	Name       string
	Options    *MetadataOptions
	Scraper    GalleryScraper
	RemoteSite string
}

type GalleryIdentifier struct {
	TxnManager           txn.Manager
	GalleryReaderUpdater GalleryReaderUpdater
	StudioReaderWriter   models.StudioReaderWriter
	PerformerCreator     PerformerCreator
	TagFinderCreator     models.TagFinderCreator

	DefaultOptions *MetadataOptions
	Sources        []GalleryScraperSource

	PostHookExecutor GalleryUpdatePostHookExecutor
}

func (t *GalleryIdentifier) Identify(ctx context.Context, galleryObj *models.Gallery) error {
	result, err := t.scrapeGallery(ctx, galleryObj)
	var multipleMatchErr *MultipleMatchesFoundError
	if err != nil {
		if !errors.As(err, &multipleMatchErr) {
			return err
		}
	}

	if result == nil {
		if multipleMatchErr != nil {
			logger.Debugf("Identify skipped because multiple results returned for %s", galleryObj.DisplayName())

			src := multipleMatchErr.Source
			options := t.getOptions(GalleryScraperSource{
				Name:       src.Name,
				Options:    src.Options,
				RemoteSite: src.RemoteSite,
			})
			if options.SkipMultipleMatchTag != nil && len(*options.SkipMultipleMatchTag) > 0 {
				err := t.addTagToGallery(ctx, galleryObj, *options.SkipMultipleMatchTag)
				if err != nil {
					return err
				}
				return nil
			}
		} else {
			logger.Debugf("Unable to identify %s", galleryObj.DisplayName())
		}
		return nil
	}

	if err := t.modifyGallery(ctx, galleryObj, result); err != nil {
		return fmt.Errorf("error modifying gallery: %w", err)
	}

	return nil
}

type galleryScrapeResult struct {
	result *models.ScrapedGallery
	source GalleryScraperSource
}

func (t *GalleryIdentifier) scrapeGallery(ctx context.Context, galleryObj *models.Gallery) (*galleryScrapeResult, error) {
	for _, source := range t.Sources {
		results, err := source.Scraper.ScrapeGalleries(ctx, galleryObj.ID)
		if err != nil {
			logger.Errorf("error scraping from %v: %v", source.Scraper, err)
			continue
		}

		if len(results) > 0 {
			options := t.getOptions(source)
			if len(results) > 1 && utils.IsTrue(options.SkipMultipleMatches) {
				return nil, &MultipleMatchesFoundError{
					Source: ScraperSource{
						Name:       source.Name,
						Options:    source.Options,
						RemoteSite: source.RemoteSite,
					},
				}
			} else {
				return &galleryScrapeResult{
					result: results[0],
					source: source,
				}, nil
			}
		}
	}

	return nil, nil
}

func (t *GalleryIdentifier) getOptions(source GalleryScraperSource) MetadataOptions {
	var options MetadataOptions
	if t.DefaultOptions != nil {
		options = *t.DefaultOptions
	}
	if source.Options == nil {
		return options
	}

	return mergeSharedMetadataOptions(options, source.Options)
}

func (t *GalleryIdentifier) getGalleryUpdater(ctx context.Context, g *models.Gallery, result *galleryScrapeResult) (*gallery.UpdateSet, error) {
	ret := &gallery.UpdateSet{
		ID: g.ID,
	}

	allOptions := []MetadataOptions{}
	if result.source.Options != nil {
		allOptions = append(allOptions, *result.source.Options)
	}
	if t.DefaultOptions != nil {
		allOptions = append(allOptions, *t.DefaultOptions)
	}

	fieldOptions := getFieldOptions(allOptions)
	options := t.getOptions(result.source)

	scraped := result.result

	rel := galleryRelationships{
		galleryReader:            t.GalleryReaderUpdater,
		studioReaderWriter:       t.StudioReaderWriter,
		performerCreator:         t.PerformerCreator,
		tagCreator:               t.TagFinderCreator,
		gallery:                  g,
		scraped:                  scraped,
		remoteSite:               result.source.RemoteSite,
		fieldOptions:             fieldOptions,
		skipSingleNamePerformers: utils.IsTrue(options.SkipSingleNamePerformers),
	}

	setOrganized := utils.IsTrue(options.SetOrganized)
	ret.Partial = getGalleryPartial(g, scraped, fieldOptions, setOrganized)

	studioID, err := rel.studio(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting studio: %w", err)
	}

	if studioID != nil {
		ret.Partial.StudioID = models.NewOptionalInt(*studioID)
	}

	// nil allowedGenders means include all performers
	allowedGenders := resolveAllowedGenders(options)

	addSkipSingleNamePerformerTag := false
	performerIDs, err := rel.performers(ctx, allowedGenders)
	if err != nil {
		if errors.Is(err, ErrSkipSingleNamePerformer) {
			addSkipSingleNamePerformerTag = true
		} else {
			return nil, err
		}
	}
	if performerIDs != nil {
		ret.Partial.PerformerIDs = &models.UpdateIDs{
			IDs:  performerIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	}

	tagIDs, err := rel.tags(ctx)
	if err != nil {
		return nil, err
	}
	if addSkipSingleNamePerformerTag && options.SkipSingleNamePerformerTag != nil {
		tagID, err := parseTagID(*options.SkipSingleNamePerformerTag)
		if err != nil {
			return nil, err
		}

		tagIDs = sliceutil.AppendUnique(tagIDs, tagID)
	}
	if tagIDs != nil {
		ret.Partial.TagIDs = &models.UpdateIDs{
			IDs:  tagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	}

	return ret, nil
}

func (t *GalleryIdentifier) modifyGallery(ctx context.Context, g *models.Gallery, result *galleryScrapeResult) error {
	var updater *gallery.UpdateSet
	if err := txn.WithTxn(ctx, t.TxnManager, func(ctx context.Context) error {
		if err := g.LoadURLs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}
		if err := g.LoadPerformerIDs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}
		if err := g.LoadTagIDs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}

		var err error
		updater, err = t.getGalleryUpdater(ctx, g, result)
		if err != nil {
			return err
		}

		if updater.IsEmpty() {
			logger.Debugf("Nothing to set for %s", g.DisplayName())
			return nil
		}

		if _, err := updater.Update(ctx, t.GalleryReaderUpdater); err != nil {
			return fmt.Errorf("error updating gallery: %w", err)
		}

		as := ""
		title := updater.Partial.Title
		if title.Ptr() != nil {
			as = fmt.Sprintf(" as %s", title.Value)
		}
		logger.Infof("Successfully identified %s%s using %s", g.DisplayName(), as, result.source.Name)

		return nil
	}); err != nil {
		return err
	}

	// fire post-update hooks
	if !updater.IsEmpty() && t.PostHookExecutor != nil {
		updateInput := updater.UpdateInput()
		fields := utils.NotNilFields(updateInput, "json")
		t.PostHookExecutor.ExecuteGalleryUpdatePostHooks(ctx, updateInput, fields)
	}

	return nil
}

func (t *GalleryIdentifier) addTagToGallery(ctx context.Context, g *models.Gallery, tagToAdd string) error {
	if err := txn.WithTxn(ctx, t.TxnManager, func(ctx context.Context) error {
		tagID, err := parseTagID(tagToAdd)
		if err != nil {
			return err
		}

		if err := g.LoadTagIDs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}
		existing := g.TagIDs.List()

		if slices.Contains(existing, tagID) {
			return nil
		}

		if err := gallery.AddTag(ctx, t.GalleryReaderUpdater, g, tagID); err != nil {
			return err
		}

		ret, err := t.TagFinderCreator.Find(ctx, tagID)
		if err != nil || ret == nil {
			logger.Infof("Added tag id %s to skipped gallery %s", tagToAdd, g.DisplayName())
		} else {
			logger.Infof("Added tag %s to skipped gallery %s", ret.Name, g.DisplayName())
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func getGalleryPartial(gallery *models.Gallery, scraped *models.ScrapedGallery, fieldOptions map[string]*FieldOptions, setOrganized bool) models.GalleryPartial {
	partial := models.GalleryPartial{}

	if scraped.Title != nil && (gallery.Title != *scraped.Title) {
		if shouldSetSingleValueField(fieldOptions["title"], gallery.Title != "") {
			partial.Title = models.NewOptionalString(*scraped.Title)
		}
	}
	if scraped.Code != nil && (gallery.Code != *scraped.Code) {
		if shouldSetSingleValueField(fieldOptions["code"], gallery.Code != "") {
			partial.Code = models.NewOptionalString(*scraped.Code)
		}
	}
	if scraped.Details != nil && (gallery.Details != *scraped.Details) {
		if shouldSetSingleValueField(fieldOptions["details"], gallery.Details != "") {
			partial.Details = models.NewOptionalString(*scraped.Details)
		}
	}
	if scraped.Photographer != nil && (gallery.Photographer != *scraped.Photographer) {
		if shouldSetSingleValueField(fieldOptions["photographer"], gallery.Photographer != "") {
			partial.Photographer = models.NewOptionalString(*scraped.Photographer)
		}
	}
	if scraped.Date != nil && (gallery.Date == nil || gallery.Date.String() != *scraped.Date) {
		if shouldSetSingleValueField(fieldOptions["date"], gallery.Date != nil) {
			d, err := models.ParseDate(*scraped.Date)
			if err == nil {
				partial.Date = models.NewOptionalDate(d)
			}
		}
	}
	if len(scraped.URLs) > 0 && shouldSetSingleValueField(fieldOptions["url"], false) {
		switch getFieldStrategy(fieldOptions["url"]) {
		case FieldStrategyOverwrite:
			if !sliceutil.SliceSame(scraped.URLs, gallery.URLs.List()) {
				partial.URLs = &models.UpdateStrings{
					Values: scraped.URLs,
					Mode:   models.RelationshipUpdateModeSet,
				}
			}
		case FieldStrategyMerge:
			urls := sliceutil.AppendUniques(gallery.URLs.List(), scraped.URLs)
			if len(urls) != len(gallery.URLs.List()) {
				partial.URLs = &models.UpdateStrings{
					Values: urls,
					Mode:   models.RelationshipUpdateModeSet,
				}
			}
		}
	}

	if setOrganized && !gallery.Organized {
		partial.Organized = models.NewOptionalBool(true)
	}

	return partial
}
