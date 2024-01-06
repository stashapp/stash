package identify

import (
	"context"
	"errors"
	"fmt"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
	"strconv"
)

type GalleryMultipleMatchesFoundError struct {
	Source GalleryScraperSource
}

type GalleryScraper interface {
	ScrapeGalleries(ctx context.Context, galleryID int) ([]*scraper.ScrapedGallery, error)
}

func (e *GalleryMultipleMatchesFoundError) Error() string {
	return fmt.Sprintf("multiple matches found for %s", e.Source.Name)
}

type GalleryUpdatePostHookExecutor interface {
	ExecuteGalleryUpdatePostHooks(ctx context.Context, input models.GalleryUpdateInput, inputFields []string)
}

type GalleryScraperSource struct {
	Name       string
	Options    *GalleryMetadataOptions
	Scraper    GalleryScraper
	RemoteSite string
}

type GalleryIdentifier struct {
	TxnManager           txn.Manager
	GalleryReaderUpdater GalleryReaderUpdater
	StudioReaderWriter   models.StudioReaderWriter
	PerformerCreator     PerformerCreator
	TagFinderCreator     models.TagFinderCreator

	DefaultOptions                *GalleryMetadataOptions
	Sources                       []GalleryScraperSource
	GalleryUpdatePostHookExecutor GalleryUpdatePostHookExecutor
}

func (t *GalleryIdentifier) Identify(ctx context.Context, gallery *models.Gallery) error {
	result, err := t.scrapeGallery(ctx, gallery)
	var multipleMatchErr *GalleryMultipleMatchesFoundError
	if err != nil {
		if !errors.As(err, &multipleMatchErr) {
			return err
		}
	}

	if result == nil {
		if multipleMatchErr != nil {
			logger.Debugf("Identify skipped because multiple results returned for %s", gallery.Path)

			// find if the gallery should be tagged for multiple results
			options := t.getOptions(multipleMatchErr.Source)
			if options.SkipMultipleMatchTag != nil && len(*options.SkipMultipleMatchTag) > 0 {
				// Tag it with the multiple results tag
				err := t.addTagToGallery(ctx, gallery, *options.SkipMultipleMatchTag)
				if err != nil {
					return err
				}
				return nil
			}
		} else {
			logger.Debugf("Unable to identify %s", gallery.Path)
		}
		return nil
	}

	// results were found, modify the gallery
	if err := t.modifyGallery(ctx, gallery, result); err != nil {
		return fmt.Errorf("error modifying gallery: %v", err)
	}

	return nil
}

type galleryScrapeResult struct {
	result *scraper.ScrapedGallery
	source GalleryScraperSource
}

func (t *GalleryIdentifier) scrapeGallery(ctx context.Context, gallery *models.Gallery) (*galleryScrapeResult, error) {
	// iterate through the input sources
	for _, source := range t.Sources {
		// scrape using the source
		results, err := source.Scraper.ScrapeGalleries(ctx, gallery.ID)
		if err != nil {
			logger.Errorf("error scraping from %v: %v", source.Scraper, err)
			continue
		}

		if len(results) > 0 {
			options := t.getOptions(source)
			if len(results) > 1 && utils.IsTrue(options.SkipMultipleMatches) {
				return nil, &GalleryMultipleMatchesFoundError{
					Source: source,
				}
			} else {
				// if results were found then return
				return &galleryScrapeResult{
					result: results[0],
					source: source,
				}, nil
			}
		}
	}

	return nil, nil
}

// Returns a GalleryMetadataOptions object with any default options overwritten by source specific options
func (t *GalleryIdentifier) getOptions(source GalleryScraperSource) GalleryMetadataOptions {
	var options GalleryMetadataOptions
	if t.DefaultOptions != nil {
		options = *t.DefaultOptions
	}
	if source.Options == nil {
		return options
	}

	if source.Options.SetOrganized != nil {
		options.SetOrganized = source.Options.SetOrganized
	}
	if source.Options.IncludeMalePerformers != nil {
		options.IncludeMalePerformers = source.Options.IncludeMalePerformers
	}
	if source.Options.SkipMultipleMatches != nil {
		options.SkipMultipleMatches = source.Options.SkipMultipleMatches
	}
	if source.Options.SkipMultipleMatchTag != nil && len(*source.Options.SkipMultipleMatchTag) > 0 {
		options.SkipMultipleMatchTag = source.Options.SkipMultipleMatchTag
	}
	if source.Options.SkipSingleNamePerformers != nil {
		options.SkipSingleNamePerformers = source.Options.SkipSingleNamePerformers
	}
	if source.Options.SkipSingleNamePerformerTag != nil && len(*source.Options.SkipSingleNamePerformerTag) > 0 {
		options.SkipSingleNamePerformerTag = source.Options.SkipSingleNamePerformerTag
	}

	return options
}

func (t *GalleryIdentifier) getGalleryUpdater(ctx context.Context, s *models.Gallery, result *galleryScrapeResult) (*gallery.UpdateSet, error) {
	ret := &gallery.UpdateSet{
		ID: s.ID,
	}

	allOptions := []GalleryMetadataOptions{}
	if result.source.Options != nil {
		allOptions = append(allOptions, *result.source.Options)
	}
	if t.DefaultOptions != nil {
		allOptions = append(allOptions, *t.DefaultOptions)
	}

	fieldOptions := getFieldOptionsGallery(allOptions)
	options := t.getOptions(result.source)

	scraped := result.result

	rel := galleryRelationships{
		studioReaderWriter:       t.StudioReaderWriter,
		performerCreator:         t.PerformerCreator,
		tagCreator:               t.TagFinderCreator,
		gallery:                  s,
		result:                   result,
		fieldOptions:             fieldOptions,
		skipSingleNamePerformers: utils.IsTrue(options.SkipSingleNamePerformers),
	}

	setOrganized := utils.IsTrue(options.SetOrganized)
	ret.Partial = getGalleryPartial(s, scraped, fieldOptions, setOrganized)

	studioID, err := rel.studio(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting studio: %w", err)
	}

	if studioID != nil {
		ret.Partial.StudioID = models.NewOptionalInt(*studioID)
	}

	includeMalePerformers := true
	if options.IncludeMalePerformers != nil {
		includeMalePerformers = *options.IncludeMalePerformers
	}

	addSkipSingleNamePerformerTag := false
	performerIDs, err := rel.performers(ctx, !includeMalePerformers)
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
		tagID, err := strconv.ParseInt(*options.SkipSingleNamePerformerTag, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting tag ID %s: %w", *options.SkipSingleNamePerformerTag, err)
		}

		tagIDs = sliceutil.AppendUnique(tagIDs, int(tagID))
	}
	if tagIDs != nil {
		ret.Partial.TagIDs = &models.UpdateIDs{
			IDs:  tagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	}

	return ret, nil
}

func (t *GalleryIdentifier) modifyGallery(ctx context.Context, s *models.Gallery, result *galleryScrapeResult) error {
	var updater *gallery.UpdateSet
	if err := txn.WithTxn(ctx, t.TxnManager, func(ctx context.Context) error {
		// load gallery relationships
		if err := s.LoadURLs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}
		if err := s.LoadPerformerIDs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}
		if err := s.LoadTagIDs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}

		var err error
		updater, err = t.getGalleryUpdater(ctx, s, result)
		if err != nil {
			return err
		}

		// don't update anything if nothing was set
		if updater.IsEmpty() {
			logger.Debugf("Nothing to set for %s", s.Path)
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
		logger.Infof("Successfully identified %s%s using %s", s.Path, as, result.source.Name)

		return nil
	}); err != nil {
		return err
	}

	// fire post-update hooks
	if !updater.IsEmpty() {
		updateInput := updater.UpdateInput()
		fields := utils.NotNilFields(updateInput, "json")
		t.GalleryUpdatePostHookExecutor.ExecuteGalleryUpdatePostHooks(ctx, updateInput, fields)
	}

	return nil
}

func (t *GalleryIdentifier) addTagToGallery(ctx context.Context, s *models.Gallery, tagToAdd string) error {
	if err := txn.WithTxn(ctx, t.TxnManager, func(ctx context.Context) error {
		tagID, err := strconv.Atoi(tagToAdd)
		if err != nil {
			return fmt.Errorf("error converting tag ID %s: %w", tagToAdd, err)
		}

		if err := s.LoadTagIDs(ctx, t.GalleryReaderUpdater); err != nil {
			return err
		}
		existing := s.TagIDs.List()

		if sliceutil.Contains(existing, tagID) {
			// skip if the gallery was already tagged
			return nil
		}

		if err := gallery.AddTag(ctx, t.GalleryReaderUpdater, s, tagID); err != nil {
			return err
		}

		ret, err := t.TagFinderCreator.Find(ctx, tagID)
		if err != nil {
			logger.Infof("Added tag id %s to skipped gallery %s", tagToAdd, s.Path)
		} else {
			logger.Infof("Added tag %s to skipped gallery %s", ret.Name, s.Path)
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func getGalleryPartial(gallery *models.Gallery, scraped *scraper.ScrapedGallery, fieldOptions map[string]*FieldOptions, setOrganized bool) models.GalleryPartial {
	partial := models.GalleryPartial{}

	if scraped.Title != nil && (gallery.Title != *scraped.Title) {
		if shouldSetSingleValueField(fieldOptions["title"], gallery.Title != "") {
			partial.Title = models.NewOptionalString(*scraped.Title)
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
	if scraped.Details != nil && (gallery.Details != *scraped.Details) {
		if shouldSetSingleValueField(fieldOptions["details"], gallery.Details != "") {
			partial.Details = models.NewOptionalString(*scraped.Details)
		}
	}
	if len(scraped.URLs) > 0 && shouldSetSingleValueField(fieldOptions["url"], false) {
		// if overwrite, then set over the top
		switch getFieldStrategy(fieldOptions["url"]) {
		case FieldStrategyOverwrite:
			// only overwrite if not equal
			if len(sliceutil.Exclude(gallery.URLs.List(), scraped.URLs)) != 0 {
				partial.URLs = &models.UpdateStrings{
					Values: scraped.URLs,
					Mode:   models.RelationshipUpdateModeSet,
				}
			}
		case FieldStrategyMerge:
			// if merge, add if not already present
			urls := sliceutil.AppendUniques(gallery.URLs.List(), scraped.URLs)

			if len(urls) != len(gallery.URLs.List()) {
				partial.URLs = &models.UpdateStrings{
					Values: urls,
					Mode:   models.RelationshipUpdateModeSet,
				}
			}
		}
	}
	if scraped.Photographer != nil && (gallery.Photographer != *scraped.Photographer) {
		if shouldSetSingleValueField(fieldOptions["photographer"], gallery.Photographer != "") {
			partial.Photographer = models.NewOptionalString(*scraped.Photographer)
		}
	}
	if scraped.Code != nil && (gallery.Code != *scraped.Code) {
		if shouldSetSingleValueField(fieldOptions["code"], gallery.Code != "") {
			partial.Code = models.NewOptionalString(*scraped.Code)
		}
	}

	if setOrganized && !gallery.Organized {
		partial.Organized = models.NewOptionalBool(true)
	}

	return partial
}

func getFieldOptionsGallery(options []GalleryMetadataOptions) map[string]*FieldOptions {
	// prefer source-specific field strategies, then the defaults
	ret := make(map[string]*FieldOptions)
	for _, oo := range options {
		for _, f := range oo.FieldOptions {
			if _, found := ret[f.Field]; !found {
				ret[f.Field] = f
			}
		}
	}

	return ret
}
