package identify

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type SceneScraper interface {
	ScrapeScene(sceneID int) (*models.ScrapedScene, error)
}

type ScraperSource struct {
	Name       string
	Options    *models.IdentifyMetadataOptionsInput
	Scraper    SceneScraper
	RemoteSite string
}

type SceneIdentifier struct {
	DefaultOptions *models.IdentifyMetadataOptionsInput
	Sources        []ScraperSource
}

func (t *SceneIdentifier) Identify(ctx context.Context, repo models.Repository, scene *models.Scene) error {
	result, err := t.scrapeScene(scene)
	if err != nil {
		return err
	}

	if result == nil {
		logger.Infof("Unable to identify %s", scene.Path)
		return nil
	}

	// results were found, modify the scene
	if err := t.modifyScene(ctx, repo, scene, result); err != nil {
		return fmt.Errorf("error modifying scene: %v", err)
	}

	return nil
}

type scrapeResult struct {
	result *models.ScrapedScene
	source ScraperSource
}

func (t *SceneIdentifier) scrapeScene(scene *models.Scene) (*scrapeResult, error) {
	// iterate through the input sources
	for _, source := range t.Sources {
		// scrape using the source
		scraped, err := source.Scraper.ScrapeScene(scene.ID)
		if err != nil {
			return nil, fmt.Errorf("error scraping from %v: %v", source.Scraper, err)
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

func (t *SceneIdentifier) getSceneUpdater(ctx context.Context, s *models.Scene, result *scrapeResult, repo models.Repository) (*scene.Updater, error) {
	ret := &scene.Updater{
		ID: s.ID,
	}

	options := []models.IdentifyMetadataOptionsInput{}
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
		repo:         repo,
		scene:        s,
		result:       result,
		fieldOptions: fieldOptions,
	}

	ret.Partial = getScenePartial(s, scraped, fieldOptions, setOrganized)

	studioID, err := rel.studio()
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

	ret.PerformerIDs, err = rel.performers(ignoreMale)
	if err != nil {
		return nil, err
	}

	ret.TagIDs, err = rel.tags()
	if err != nil {
		return nil, err
	}

	ret.StashIDs, err = rel.stashIDs()
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

func (t *SceneIdentifier) modifyScene(ctx context.Context, repo models.Repository, scene *models.Scene, result *scrapeResult) error {
	updater, err := t.getSceneUpdater(ctx, scene, result, repo)
	if err != nil {
		return err
	}

	// don't update anything if nothing was set
	if updater.IsEmpty() {
		logger.Infof("Nothing to set for %s", scene.Path)
		return nil
	}

	_, err = updater.Update(repo.Scene())
	if err != nil {
		return fmt.Errorf("error updating scene: %w", err)
	}

	as := ""
	title := updater.Partial.Title
	if title != nil {
		as = fmt.Sprintf(" as %s", title.String)
	}
	logger.Infof("Successfully identified %s%s using %s", scene.Path, as, result.source.Name)

	return nil
}

func getFieldOptions(options []models.IdentifyMetadataOptionsInput) map[string]*models.IdentifyFieldOptionsInput {
	// prefer source-specific field strategies, then the defaults
	ret := make(map[string]*models.IdentifyFieldOptionsInput)
	for _, oo := range options {
		for _, f := range oo.FieldOptions {
			if _, found := ret[f.Field]; !found {
				ret[f.Field] = f
			}
		}
	}

	return ret
}

func getScenePartial(scene *models.Scene, scraped *models.ScrapedScene, fieldOptions map[string]*models.IdentifyFieldOptionsInput, setOrganized bool) models.ScenePartial {
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

func shouldSetSingleValueField(strategy *models.IdentifyFieldOptionsInput, hasExistingValue bool) bool {
	// if unset then default to MERGE
	fs := models.IdentifyFieldStrategyMerge

	if strategy != nil && strategy.Strategy.IsValid() {
		fs = strategy.Strategy
	}

	if fs == models.IdentifyFieldStrategyIgnore {
		return false
	}

	return !hasExistingValue || fs == models.IdentifyFieldStrategyOverwrite
}
