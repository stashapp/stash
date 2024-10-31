// Package identify provides the scene identification functionality for the application.
// The identify functionality uses scene scrapers to identify a given scene and
// set its metadata based on the scraped data.
package identify

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

var (
	ErrSkipSingleNamePerformer = errors.New("a performer was skipped because they only had a single name and no disambiguation")
)

type MultipleMatchesFoundError struct {
	Source ScraperSource
}

func (e *MultipleMatchesFoundError) Error() string {
	return fmt.Sprintf("multiple matches found for %s", e.Source.Name)
}

type SceneScraper interface {
	ScrapeScenes(ctx context.Context, sceneID int) ([]*scraper.ScrapedScene, error)
}

type SceneUpdatePostHookExecutor interface {
	ExecuteSceneUpdatePostHooks(ctx context.Context, input models.SceneUpdateInput, inputFields []string)
}

type ScraperSource struct {
	Name       string
	Options    *MetadataOptions
	Scraper    SceneScraper
	RemoteSite string
}

type SceneIdentifier struct {
	TxnManager         txn.Manager
	SceneReaderUpdater SceneReaderUpdater
	StudioReaderWriter models.StudioReaderWriter
	PerformerCreator   PerformerCreator
	TagFinderCreator   models.TagFinderCreator

	DefaultOptions              *MetadataOptions
	Sources                     []ScraperSource
	SceneUpdatePostHookExecutor SceneUpdatePostHookExecutor
}

func (t *SceneIdentifier) Identify(ctx context.Context, scene *models.Scene) error {
	result, err := t.scrapeScene(ctx, scene)
	var multipleMatchErr *MultipleMatchesFoundError
	if err != nil {
		if !errors.As(err, &multipleMatchErr) {
			return err
		}
	}

	if result == nil {
		if multipleMatchErr != nil {
			logger.Debugf("Identify skipped because multiple results returned for %s", scene.Path)

			// find if the scene should be tagged for multiple results
			options := t.getOptions(multipleMatchErr.Source)
			if options.SkipMultipleMatchTag != nil && len(*options.SkipMultipleMatchTag) > 0 {
				// Tag it with the multiple results tag
				err := t.addTagToScene(ctx, scene, *options.SkipMultipleMatchTag)
				if err != nil {
					return err
				}
				return nil
			}
		} else {
			logger.Debugf("Unable to identify %s", scene.Path)
		}
		return nil
	}

	// results were found, modify the scene
	if err := t.modifyScene(ctx, scene, result); err != nil {
		return fmt.Errorf("error modifying scene: %v", err)
	}

	return nil
}

type scrapeResult struct {
	result *scraper.ScrapedScene
	source ScraperSource
}

func (t *SceneIdentifier) scrapeScene(ctx context.Context, scene *models.Scene) (*scrapeResult, error) {
	// iterate through the input sources
	for _, source := range t.Sources {
		// scrape using the source
		results, err := source.Scraper.ScrapeScenes(ctx, scene.ID)
		if err != nil {
			logger.Errorf("error scraping from %v: %v", source.Scraper, err)
			continue
		}

		if len(results) > 0 {
			options := t.getOptions(source)
			if len(results) > 1 && utils.IsTrue(options.SkipMultipleMatches) {
				return nil, &MultipleMatchesFoundError{
					Source: source,
				}
			} else {
				// if results were found then return
				return &scrapeResult{
					result: results[0],
					source: source,
				}, nil
			}
		}
	}

	return nil, nil
}

// Returns a MetadataOptions object with any default options overwritten by source specific options
func (t *SceneIdentifier) getOptions(source ScraperSource) MetadataOptions {
	var options MetadataOptions
	if t.DefaultOptions != nil {
		options = *t.DefaultOptions
	}
	if source.Options == nil {
		return options
	}

	if source.Options.SetCoverImage != nil {
		options.SetCoverImage = source.Options.SetCoverImage
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

func (t *SceneIdentifier) getSceneUpdater(ctx context.Context, s *models.Scene, result *scrapeResult) (*scene.UpdateSet, error) {
	ret := &scene.UpdateSet{
		ID: s.ID,
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

	rel := sceneRelationships{
		sceneReader:              t.SceneReaderUpdater,
		studioReaderWriter:       t.StudioReaderWriter,
		performerCreator:         t.PerformerCreator,
		tagCreator:               t.TagFinderCreator,
		scene:                    s,
		result:                   result,
		fieldOptions:             fieldOptions,
		skipSingleNamePerformers: utils.IsTrue(options.SkipSingleNamePerformers),
	}

	setOrganized := utils.IsTrue(options.SetOrganized)
	ret.Partial = getScenePartial(s, scraped, fieldOptions, setOrganized)

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

	// SetCoverImage defaults to true if unset
	if options.SetCoverImage == nil || *options.SetCoverImage {
		ret.CoverImage, err = rel.cover(ctx)
		if err != nil {
			return nil, err
		}
	}

	// if anything changed, also update the updated at time on the applicable stash id
	changed := !ret.IsEmpty()

	stashIDs, err := rel.stashIDs(ctx, changed)
	if err != nil {
		return nil, err
	}
	if stashIDs != nil {
		ret.Partial.StashIDs = &models.UpdateStashIDs{
			StashIDs: stashIDs,
			Mode:     models.RelationshipUpdateModeSet,
		}
	}

	return ret, nil
}

func (t *SceneIdentifier) modifyScene(ctx context.Context, s *models.Scene, result *scrapeResult) error {
	var updater *scene.UpdateSet
	if err := txn.WithTxn(ctx, t.TxnManager, func(ctx context.Context) error {
		// load scene relationships
		if err := s.LoadURLs(ctx, t.SceneReaderUpdater); err != nil {
			return err
		}
		if err := s.LoadPerformerIDs(ctx, t.SceneReaderUpdater); err != nil {
			return err
		}
		if err := s.LoadTagIDs(ctx, t.SceneReaderUpdater); err != nil {
			return err
		}
		if err := s.LoadStashIDs(ctx, t.SceneReaderUpdater); err != nil {
			return err
		}

		var err error
		updater, err = t.getSceneUpdater(ctx, s, result)
		if err != nil {
			return err
		}

		// don't update anything if nothing was set
		if updater.IsEmpty() {
			logger.Debugf("Nothing to set for %s", s.Path)
			return nil
		}

		if _, err := updater.Update(ctx, t.SceneReaderUpdater); err != nil {
			return fmt.Errorf("error updating scene: %w", err)
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
		t.SceneUpdatePostHookExecutor.ExecuteSceneUpdatePostHooks(ctx, updateInput, fields)
	}

	return nil
}

func (t *SceneIdentifier) addTagToScene(ctx context.Context, s *models.Scene, tagToAdd string) error {
	if err := txn.WithTxn(ctx, t.TxnManager, func(ctx context.Context) error {
		tagID, err := strconv.Atoi(tagToAdd)
		if err != nil {
			return fmt.Errorf("error converting tag ID %s: %w", tagToAdd, err)
		}

		if err := s.LoadTagIDs(ctx, t.SceneReaderUpdater); err != nil {
			return err
		}
		existing := s.TagIDs.List()

		if slices.Contains(existing, tagID) {
			// skip if the scene was already tagged
			return nil
		}

		if err := scene.AddTag(ctx, t.SceneReaderUpdater, s, tagID); err != nil {
			return err
		}

		ret, err := t.TagFinderCreator.Find(ctx, tagID)
		if err != nil {
			logger.Infof("Added tag id %s to skipped scene %s", tagToAdd, s.Path)
		} else {
			logger.Infof("Added tag %s to skipped scene %s", ret.Name, s.Path)
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func getFieldOptions(options []MetadataOptions) map[string]*FieldOptions {
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

func getScenePartial(scene *models.Scene, scraped *scraper.ScrapedScene, fieldOptions map[string]*FieldOptions, setOrganized bool) models.ScenePartial {
	partial := models.ScenePartial{}

	if scraped.Title != nil && (scene.Title != *scraped.Title) {
		if shouldSetSingleValueField(fieldOptions["title"], scene.Title != "") {
			partial.Title = models.NewOptionalString(*scraped.Title)
		}
	}
	if scraped.Date != nil && (scene.Date == nil || scene.Date.String() != *scraped.Date) {
		if shouldSetSingleValueField(fieldOptions["date"], scene.Date != nil) {
			d, err := models.ParseDate(*scraped.Date)
			if err == nil {
				partial.Date = models.NewOptionalDate(d)
			}
		}
	}
	if scraped.Details != nil && (scene.Details != *scraped.Details) {
		if shouldSetSingleValueField(fieldOptions["details"], scene.Details != "") {
			partial.Details = models.NewOptionalString(*scraped.Details)
		}
	}
	if len(scraped.URLs) > 0 && shouldSetSingleValueField(fieldOptions["url"], false) {
		// if overwrite, then set over the top
		switch getFieldStrategy(fieldOptions["url"]) {
		case FieldStrategyOverwrite:
			// only overwrite if not equal
			if len(sliceutil.Exclude(scraped.URLs, scene.URLs.List())) != 0 {
				partial.URLs = &models.UpdateStrings{
					Values: scraped.URLs,
					Mode:   models.RelationshipUpdateModeSet,
				}
			}
		case FieldStrategyMerge:
			// if merge, add if not already present
			urls := sliceutil.AppendUniques(scene.URLs.List(), scraped.URLs)

			if len(urls) != len(scene.URLs.List()) {
				partial.URLs = &models.UpdateStrings{
					Values: urls,
					Mode:   models.RelationshipUpdateModeSet,
				}
			}
		}
	}
	if scraped.Director != nil && (scene.Director != *scraped.Director) {
		if shouldSetSingleValueField(fieldOptions["director"], scene.Director != "") {
			partial.Director = models.NewOptionalString(*scraped.Director)
		}
	}
	if scraped.Code != nil && (scene.Code != *scraped.Code) {
		if shouldSetSingleValueField(fieldOptions["code"], scene.Code != "") {
			partial.Code = models.NewOptionalString(*scraped.Code)
		}
	}

	if setOrganized && !scene.Organized {
		partial.Organized = models.NewOptionalBool(true)
	}

	return partial
}

func getFieldStrategy(strategy *FieldOptions) FieldStrategy {
	// if unset then default to MERGE
	fs := FieldStrategyMerge

	if strategy != nil && strategy.Strategy.IsValid() {
		fs = strategy.Strategy
	}

	return fs
}

func shouldSetSingleValueField(strategy *FieldOptions, hasExistingValue bool) bool {
	// if unset then default to MERGE
	fs := getFieldStrategy(strategy)

	if fs == FieldStrategyIgnore {
		return false
	}

	return !hasExistingValue || fs == FieldStrategyOverwrite
}
