package identify

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

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
	SceneReaderUpdater SceneReaderUpdater
	StudioCreator      StudioCreator
	PerformerCreator   PerformerCreator
	TagCreator         TagCreator

	DefaultOptions              *MetadataOptions
	Sources                     []ScraperSource
	SceneUpdatePostHookExecutor SceneUpdatePostHookExecutor
}

func (t *SceneIdentifier) Identify(ctx context.Context, txnManager txn.Manager, scene *models.Scene) error {
	result, err := t.scrapeScene(ctx, txnManager, scene)
	if err != nil {
		return err
	}

	if result == nil {
		logger.Debugf("Unable to identify %s", scene.Path)
		return nil
	}

	// results were found, modify the scene
	if err := t.modifyScene(ctx, txnManager, scene, result); err != nil {
		return fmt.Errorf("error modifying scene: %v", err)
	}

	return nil
}

type scrapeResult struct {
	result *scraper.ScrapedScene
	source ScraperSource
}

func (t *SceneIdentifier) scrapeScene(ctx context.Context, txnManager txn.Manager, scene *models.Scene) (*scrapeResult, error) {
	// iterate through the input sources
	for _, source := range t.Sources {
		// scrape using the source
		results, err := source.Scraper.ScrapeScenes(ctx, scene.ID)
		if err != nil {
			logger.Errorf("error scraping from %v: %v", source.Scraper, err)
			continue
		}

		if len(results) > 0 {
			options := t.getOptions(&source)
			if len(results) > 1 && *options.SkipMultipleMatches {
				if options.SkipMultipleMatchTag != nil && len(*options.SkipMultipleMatchTag) > 0 {
					// Tag it with the multiple results tag and ignore
					err := t.addTagToScene(ctx, txnManager, scene, options.SkipMultipleMatchTag)
					if err != nil {
						return nil, err
					}
					return nil, nil
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
func (t *SceneIdentifier) getOptions(source *ScraperSource) MetadataOptions {
	options := t.DefaultOptions
	if source.Options.SetCoverImage != nil {
		options.SetCoverImage = source.Options.SetCoverImage
	}
	if source.Options.SetOrganized != nil {
		options.SetOrganized = source.Options.SetOrganized
	}
	if source.Options.IncludeMalePerformers != nil {
		options.IncludeMalePerformers = source.Options.IncludeMalePerformers
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
	return *options
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
	options := t.getOptions(&result.source)

	scraped := result.result

	rel := sceneRelationships{
		sceneReader:              t.SceneReaderUpdater,
		studioCreator:            t.StudioCreator,
		performerCreator:         t.PerformerCreator,
		tagCreator:               t.TagCreator,
		scene:                    s,
		result:                   result,
		fieldOptions:             fieldOptions,
		skipSingleNamePerformers: *options.SkipSingleNamePerformers,
	}

	ret.Partial = getScenePartial(s, scraped, fieldOptions, *options.SetOrganized)

	studioID, err := rel.studio(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting studio: %w", err)
	}

	if studioID != nil {
		ret.Partial.StudioID = models.NewOptionalInt(*studioID)
	}

	performerIDs, err := rel.performers(ctx, !*options.IncludeMalePerformers)
	if err != nil {
		return nil, err
	}
	addSkipSingleNamePerformerTag := false
	if performerIDs != nil {
		// If there is a -1 id, that means a performer was skipped due to SkipSingleNamePerformers
		i := slices.Index(performerIDs, -1)
		if i >= 0 {
			performerIDs = slices.Delete(performerIDs, i, i+1)

			if options.SkipSingleNamePerformerTag != nil && len(*options.SkipSingleNamePerformerTag) > 0 {
				// Tag it with the skipped single name performers tag
				addSkipSingleNamePerformerTag = true
			}
		}
		ret.Partial.PerformerIDs = &models.UpdateIDs{
			IDs:  performerIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	}

	tagIDs, err := rel.tags(ctx)
	if err != nil {
		return nil, err
	}
	if addSkipSingleNamePerformerTag {
		tagID, err := strconv.ParseInt(*options.SkipSingleNamePerformerTag, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting tag ID %s: %w", *options.SkipSingleNamePerformerTag, err)
		}

		tagIDs = intslice.IntAppendUnique(tagIDs, int(tagID))
	}
	if tagIDs != nil {
		ret.Partial.TagIDs = &models.UpdateIDs{
			IDs:  tagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	}

	stashIDs, err := rel.stashIDs(ctx)
	if err != nil {
		return nil, err
	}
	if stashIDs != nil {
		ret.Partial.StashIDs = &models.UpdateStashIDs{
			StashIDs: stashIDs,
			Mode:     models.RelationshipUpdateModeSet,
		}
	}

	if *options.SetCoverImage {
		ret.CoverImage, err = rel.cover(ctx)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (t *SceneIdentifier) modifyScene(ctx context.Context, txnManager txn.Manager, s *models.Scene, result *scrapeResult) error {
	var updater *scene.UpdateSet
	if err := txn.WithTxn(ctx, txnManager, func(ctx context.Context) error {
		// load scene relationships
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

func (t *SceneIdentifier) addTagToScene(ctx context.Context, txnManager txn.Manager, s *models.Scene, tagToAdd *string) error {
	if err := txn.WithTxn(ctx, txnManager, func(ctx context.Context) error {
		if err := s.LoadTagIDs(ctx, t.SceneReaderUpdater); err != nil {
			return err
		}

		ret := &scene.UpdateSet{
			ID: s.ID,
		}
		ret.Partial = models.NewScenePartial()

		// add to the existing tags
		originalTagIDs := s.TagIDs.List()
		var tagIDs []int
		tagIDs = originalTagIDs

		tagID, err2 := strconv.ParseInt(*tagToAdd, 10, 64)
		if err2 != nil {
			return fmt.Errorf("error converting tag ID %s: %w", *tagToAdd, err2)
		}
		tagIDs = intslice.IntAppendUnique(tagIDs, int(tagID))

		// skip if the scene was already tagged
		if sliceutil.SliceSame(originalTagIDs, tagIDs) {
			return nil
		}

		ret.Partial.TagIDs = &models.UpdateIDs{
			IDs:  tagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}

		if _, err := ret.Update(ctx, t.SceneReaderUpdater); err != nil {
			return fmt.Errorf("error updating scene: %w", err)
		}

		logger.Infof("Added tag %s to skipped scene %s", *tagToAdd, s.Path)

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
			d := models.NewDate(*scraped.Date)
			partial.Date = models.NewOptionalDate(d)
		}
	}
	if scraped.Details != nil && (scene.Details != *scraped.Details) {
		if shouldSetSingleValueField(fieldOptions["details"], scene.Details != "") {
			partial.Details = models.NewOptionalString(*scraped.Details)
		}
	}
	if scraped.URL != nil && (scene.URL != *scraped.URL) {
		if shouldSetSingleValueField(fieldOptions["url"], scene.URL != "") {
			partial.URL = models.NewOptionalString(*scraped.URL)
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

func shouldSetSingleValueField(strategy *FieldOptions, hasExistingValue bool) bool {
	// if unset then default to MERGE
	fs := FieldStrategyMerge

	if strategy != nil && strategy.Strategy.IsValid() {
		fs = strategy.Strategy
	}

	if fs == FieldStrategyIgnore {
		return false
	}

	return !hasExistingValue || fs == FieldStrategyOverwrite
}
