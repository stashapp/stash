package identify

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type SceneScraper interface {
	ScrapeScene(ctx context.Context, sceneID int) (*scraper.ScrapedScene, error)
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
	ScreenshotSetter            scene.ScreenshotSetter
	SceneUpdatePostHookExecutor SceneUpdatePostHookExecutor
}

func (t *SceneIdentifier) Identify(ctx context.Context, txnManager txn.Manager, scene *models.Scene) error {
	result, err := t.scrapeScene(ctx, scene)
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

func (t *SceneIdentifier) scrapeScene(ctx context.Context, scene *models.Scene) (*scrapeResult, error) {
	// iterate through the input sources
	for _, source := range t.Sources {
		// scrape using the source
		scraped, err := source.Scraper.ScrapeScene(ctx, scene.ID)
		if err != nil {
			logger.Errorf("error scraping from %v: %v", source.Scraper, err)
			continue
		}

		// if results were found then return
		if scraped != nil {
			return &scrapeResult{
				result: scraped,
				source: source,
			}, nil
		}
	}

	return nil, nil
}

func (t *SceneIdentifier) getSceneUpdater(ctx context.Context, s *models.Scene, result *scrapeResult) (*scene.UpdateSet, error) {
	ret := &scene.UpdateSet{
		ID: s.ID,
	}

	options := []MetadataOptions{}
	if result.source.Options != nil {
		options = append(options, *result.source.Options)
	}
	if t.DefaultOptions != nil {
		options = append(options, *t.DefaultOptions)
	}

	fieldOptions := getFieldOptions(options)

	setOrganized := false
	for _, o := range options {
		if o.SetOrganized != nil {
			setOrganized = *o.SetOrganized
			break
		}
	}

	scraped := result.result

	rel := sceneRelationships{
		sceneReader:      t.SceneReaderUpdater,
		studioCreator:    t.StudioCreator,
		performerCreator: t.PerformerCreator,
		tagCreator:       t.TagCreator,
		scene:            s,
		result:           result,
		fieldOptions:     fieldOptions,
	}

	ret.Partial = getScenePartial(s, scraped, fieldOptions, setOrganized)

	studioID, err := rel.studio(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting studio: %w", err)
	}

	if studioID != nil {
		ret.Partial.StudioID = &sql.NullInt64{
			Int64: *studioID,
			Valid: true,
		}
	}

	ignoreMale := false
	for _, o := range options {
		if o.IncludeMalePerformers != nil {
			ignoreMale = !*o.IncludeMalePerformers
			break
		}
	}

	ret.PerformerIDs, err = rel.performers(ctx, ignoreMale)
	if err != nil {
		return nil, err
	}

	ret.TagIDs, err = rel.tags(ctx)
	if err != nil {
		return nil, err
	}

	ret.StashIDs, err = rel.stashIDs(ctx)
	if err != nil {
		return nil, err
	}

	setCoverImage := false
	for _, o := range options {
		if o.SetCoverImage != nil {
			setCoverImage = *o.SetCoverImage
			break
		}
	}

	if setCoverImage {
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

		_, err = updater.Update(ctx, t.SceneReaderUpdater, t.ScreenshotSetter)
		if err != nil {
			return fmt.Errorf("error updating scene: %w", err)
		}

		as := ""
		title := updater.Partial.Title
		if title != nil {
			as = fmt.Sprintf(" as %s", title.String)
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
	partial := models.ScenePartial{
		ID: scene.ID,
	}

	if scraped.Title != nil && scene.Title.String != *scraped.Title {
		if shouldSetSingleValueField(fieldOptions["title"], scene.Title.String != "") {
			partial.Title = models.NullStringPtr(*scraped.Title)
		}
	}
	if scraped.Date != nil && scene.Date.String != *scraped.Date {
		if shouldSetSingleValueField(fieldOptions["date"], scene.Date.Valid) {
			partial.Date = &models.SQLiteDate{
				String: *scraped.Date,
				Valid:  true,
			}
		}
	}
	if scraped.Details != nil && scene.Details.String != *scraped.Details {
		if shouldSetSingleValueField(fieldOptions["details"], scene.Details.String != "") {
			partial.Details = models.NullStringPtr(*scraped.Details)
		}
	}
	if scraped.URL != nil && scene.URL.String != *scraped.URL {
		if shouldSetSingleValueField(fieldOptions["url"], scene.URL.String != "") {
			partial.URL = models.NullStringPtr(*scraped.URL)
		}
	}

	if setOrganized && !scene.Organized {
		// just reuse the boolean since we know it's true
		partial.Organized = &setOrganized
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
