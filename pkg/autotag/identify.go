package autotag

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type SceneScraper interface {
	ScrapeScene(sceneID int) (*models.ScrapedScene, error)
}

type ScraperSource struct {
	Options    *models.IdentifyMetadataOptionsInput
	Scraper    SceneScraper
	RemoteSite string
}

type IdentifySceneTask struct {
	DefaultOptions *models.IdentifyMetadataOptionsInput
	Sources        []ScraperSource
	Scene          *models.Scene

	Repo models.Repository
}

func (t *IdentifySceneTask) Execute(ctx context.Context) error {
	// iterate through the input sources
	for _, source := range t.Sources {
		// scrape using the source
		scraped, err := source.Scraper.ScrapeScene(t.Scene.ID)
		if err != nil {
			return fmt.Errorf("error scraping from %v: %v", source.Scraper, err)
		}

		// if results were found then modify the scene
		if scraped != nil {
			options := modifySceneOptions{
				scene:    t.Scene,
				scraped:  scraped,
				endpoint: source.RemoteSite,
			}

			options.setOptions(source.Options, t.DefaultOptions)

			if err := t.modifyScene(ctx, options); err != nil {
				return fmt.Errorf("error modifying scene: %v", err)
			}

			// scene was matched
			return nil
		}
	}

	logger.Infof("Unable to identify %s", t.Scene.Path)
	return nil
}

func (t *IdentifySceneTask) modifyScene(ctx context.Context, options modifySceneOptions) error {
	target := options.scene

	partial := models.ScenePartial{
		ID: target.ID,
	}

	fieldsSet := t.setSceneFields(&partial, options)

	studioSet, err := t.setSceneStudio(&partial, options)
	if err != nil {
		return fmt.Errorf("error setting studio: %w", err)
	}

	performerIDs, err := t.getScenePerformerIDs(options)
	if err != nil {
		return err
	}

	tagIDs, err := t.getSceneTagIDs(options)
	if err != nil {
		return err
	}

	stashIDs, err := t.getSceneStashIDs(options)
	if err != nil {
		return err
	}

	coverImage, err := t.getSceneCover(ctx, options)
	if err != nil {
		return err
	}

	// don't update anything if nothing was set
	if !fieldsSet && !studioSet && len(performerIDs) == 0 && len(tagIDs) == 0 && len(stashIDs) == 0 && len(coverImage) == 0 {
		logger.Infof("Nothing to set for %s", target.Path)
		return nil
	}

	qb := t.Repo.Scene()
	partial.UpdatedAt = &models.SQLiteTimestamp{
		Timestamp: time.Now(),
	}
	_, err = qb.Update(partial)
	if err != nil {
		return fmt.Errorf("error updating scene: %w", err)
	}

	if len(performerIDs) > 0 {
		if err := qb.UpdatePerformers(t.Scene.ID, performerIDs); err != nil {
			return fmt.Errorf("error updating scene performers: %w", err)
		}
	}

	if len(tagIDs) > 0 {
		if err := qb.UpdateTags(t.Scene.ID, tagIDs); err != nil {
			return fmt.Errorf("error updating scene tags: %w", err)
		}
	}

	if len(stashIDs) > 0 {
		if err := qb.UpdateStashIDs(t.Scene.ID, stashIDs); err != nil {
			return fmt.Errorf("error updating scene stash_ids: %w", err)
		}
	}

	if len(coverImage) > 0 {
		if err := qb.UpdateCover(t.Scene.ID, coverImage); err != nil {
			return fmt.Errorf("error updating scene cover: %w", err)
		}
	}

	as := ""
	if partial.Title != nil {
		as = fmt.Sprintf(" as %s", partial.Title.String)
	}
	logger.Infof("Successfully identified %s%s", target.Path, as)

	return nil
}

// setSceneFields sets scene fields based on field strategies. Returns true if at least one
// field was set.
func (t *IdentifySceneTask) setSceneFields(partial *models.ScenePartial, input modifySceneOptions) bool {
	scraped := input.scraped
	fieldStrategies := input.fieldOptions
	target := input.scene
	set := false

	if scraped.Title != nil && target.Title.String != *scraped.Title {
		if t.shouldSetSingleValueField(fieldStrategies["title"], target.Title.String != "") {
			partial.Title = models.NullStringPtr(*scraped.Title)
			set = true
		}
	}
	if scraped.Date != nil && target.Date.String != *scraped.Date {
		if t.shouldSetSingleValueField(fieldStrategies["date"], target.Date.Valid) {
			partial.Date = &models.SQLiteDate{
				String: *scraped.Date,
			}
			set = true
		}
	}
	if scraped.Details != nil && target.Details.String != *scraped.Details {
		if t.shouldSetSingleValueField(fieldStrategies["details"], target.Details.String != "") {
			partial.Details = models.NullStringPtr(*scraped.Details)
			set = true
		}
	}
	if scraped.URL != nil && target.URL.String != *scraped.URL {
		if t.shouldSetSingleValueField(fieldStrategies["url"], target.URL.String != "") {
			partial.URL = models.NullStringPtr(*scraped.URL)
			set = true
		}
	}

	setOrganized := false
	for _, o := range input.options {
		if o.SetOrganized != nil {
			setOrganized = *o.SetOrganized
			break
		}
	}
	if setOrganized && !target.Organized {
		// just reuse the boolean since we know it's true
		partial.Organized = &setOrganized
		set = true
	}

	return set
}

func (t *IdentifySceneTask) setSceneStudio(partial *models.ScenePartial, input modifySceneOptions) (bool, error) {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["studio"]
	target := input.scene
	r := t.Repo
	set := false

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)

	if scraped.Studio != nil {
		if t.shouldSetSingleValueField(fieldStrategy, target.StudioID.Valid) {
			if scraped.Studio.StoredID != nil {
				// existing studio, just set it
				studioID, err := strconv.ParseInt(*scraped.Studio.StoredID, 10, 64)
				if err != nil {
					return false, fmt.Errorf("error converting studio ID %s: %w", *scraped.Studio.StoredID, err)
				}

				if target.StudioID.Int64 == studioID {
					partial.StudioID = &sql.NullInt64{
						Int64: studioID,
						Valid: true,
					}
					set = true
				}
			} else if createMissing {
				created, err := r.Studio().Create(scrapedToStudioInput(scraped.Studio))
				if err != nil {
					return false, fmt.Errorf("error creating studio: %w", err)
				}

				if input.endpoint != "" && scraped.RemoteSiteID != nil {
					if err := r.Studio().UpdateStashIDs(created.ID, []models.StashID{
						{
							Endpoint: input.endpoint,
							StashID:  *scraped.Studio.RemoteSiteID,
						},
					}); err != nil {
						return false, fmt.Errorf("error setting studio stash id: %w", err)
					}
				}

				partial.StudioID = &sql.NullInt64{
					Int64: int64(created.ID),
					Valid: true,
				}
				set = true
			}
		}
	}

	return set, nil
}

func scrapedToStudioInput(studio *models.ScrapedStudio) models.Studio {
	currentTime := time.Now()
	ret := models.Studio{
		Name:      sql.NullString{String: studio.Name, Valid: true},
		Checksum:  utils.MD5FromString(studio.Name),
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}

	if studio.URL != nil {
		ret.URL = sql.NullString{String: *studio.URL, Valid: true}
	}

	return ret
}

func idListEquals(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aa := range a {
		if !utils.IntInclude(b, aa) {
			return false
		}
	}

	return true
}

func (t *IdentifySceneTask) getScenePerformerIDs(input modifySceneOptions) ([]int, error) {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["performers"]
	target := input.scene
	r := t.Repo

	// just check if ignored
	if !t.shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := models.IdentifyFieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	ignoreMale := false
	for _, o := range input.options {
		if o.IncludeMalePerformers != nil {
			ignoreMale = !*o.IncludeMalePerformers
			break
		}
	}

	var originalPerformerIDs []int
	var performerIDs []int
	if len(scraped.Performers) > 0 {
		var err error
		originalPerformerIDs, err = r.Scene().GetPerformerIDs(target.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting scene performers: %w", err)
		}

		if strategy == models.IdentifyFieldStrategyMerge {
			// add to existing
			performerIDs = originalPerformerIDs
		}

		for _, p := range scraped.Performers {
			if ignoreMale && p.Gender != nil && strings.EqualFold(*p.Gender, models.GenderEnumMale.String()) {
				continue
			}

			if p.StoredID != nil {
				// existing performer, just add it
				performerID, err := strconv.ParseInt(*p.StoredID, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error converting performer ID %s: %w", *p.StoredID, err)
				}

				performerIDs = utils.IntAppendUnique(performerIDs, int(performerID))
			} else if createMissing && p.Name != nil { // name is mandatory
				created, err := r.Performer().Create(scrapedToPerformerInput(p))
				if err != nil {
					return nil, fmt.Errorf("error creating performer: %w", err)
				}

				if input.endpoint != "" && scraped.RemoteSiteID != nil {
					if err := r.Performer().UpdateStashIDs(created.ID, []models.StashID{
						{
							Endpoint: input.endpoint,
							StashID:  *scraped.Studio.RemoteSiteID,
						},
					}); err != nil {
						return nil, fmt.Errorf("error setting performer stash id: %w", err)
					}
				}

				performerIDs = append(performerIDs, created.ID)
			}
		}
	}

	// don't return if nothing was added
	if idListEquals(originalPerformerIDs, performerIDs) {
		return nil, nil
	}

	return performerIDs, nil
}

func scrapedToPerformerInput(performer *models.ScrapedPerformer) models.Performer {
	currentTime := time.Now()
	ret := models.Performer{
		Name:      sql.NullString{String: *performer.Name, Valid: true},
		Checksum:  utils.MD5FromString(*performer.Name),
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		Favorite:  sql.NullBool{Bool: false, Valid: true},
	}
	if performer.Birthdate != nil {
		ret.Birthdate = models.SQLiteDate{String: *performer.Birthdate, Valid: true}
	}
	if performer.DeathDate != nil {
		ret.DeathDate = models.SQLiteDate{String: *performer.DeathDate, Valid: true}
	}
	if performer.Gender != nil {
		ret.Gender = sql.NullString{String: *performer.Gender, Valid: true}
	}
	if performer.Ethnicity != nil {
		ret.Ethnicity = sql.NullString{String: *performer.Ethnicity, Valid: true}
	}
	if performer.Country != nil {
		ret.Country = sql.NullString{String: *performer.Country, Valid: true}
	}
	if performer.EyeColor != nil {
		ret.EyeColor = sql.NullString{String: *performer.EyeColor, Valid: true}
	}
	if performer.HairColor != nil {
		ret.HairColor = sql.NullString{String: *performer.HairColor, Valid: true}
	}
	if performer.Height != nil {
		ret.Height = sql.NullString{String: *performer.Height, Valid: true}
	}
	if performer.Measurements != nil {
		ret.Measurements = sql.NullString{String: *performer.Measurements, Valid: true}
	}
	if performer.FakeTits != nil {
		ret.FakeTits = sql.NullString{String: *performer.FakeTits, Valid: true}
	}
	if performer.CareerLength != nil {
		ret.CareerLength = sql.NullString{String: *performer.CareerLength, Valid: true}
	}
	if performer.Tattoos != nil {
		ret.Tattoos = sql.NullString{String: *performer.Tattoos, Valid: true}
	}
	if performer.Piercings != nil {
		ret.Piercings = sql.NullString{String: *performer.Piercings, Valid: true}
	}
	if performer.Aliases != nil {
		ret.Aliases = sql.NullString{String: *performer.Aliases, Valid: true}
	}
	if performer.Twitter != nil {
		ret.Twitter = sql.NullString{String: *performer.Twitter, Valid: true}
	}
	if performer.Instagram != nil {
		ret.Instagram = sql.NullString{String: *performer.Instagram, Valid: true}
	}

	return ret
}

func (t *IdentifySceneTask) getSceneTagIDs(input modifySceneOptions) ([]int, error) {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["tags"]
	target := input.scene
	r := t.Repo

	// just check if ignored
	if !t.shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := models.IdentifyFieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var originalTagIDs []int
	var tagIDs []int
	if len(scraped.Tags) > 0 {
		var err error
		originalTagIDs, err = r.Scene().GetTagIDs(target.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting scene tag: %w", err)
		}

		if strategy == models.IdentifyFieldStrategyMerge {
			// add to existing
			tagIDs = originalTagIDs
		}

		for _, t := range scraped.Tags {
			if t.StoredID != nil {
				// existing tag, just add it
				tagID, err := strconv.ParseInt(*t.StoredID, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error converting tag ID %s: %w", *t.StoredID, err)
				}

				tagIDs = utils.IntAppendUnique(tagIDs, int(tagID))
			} else if createMissing {
				now := time.Now()
				created, err := r.Tag().Create(models.Tag{
					Name:      t.Name,
					CreatedAt: models.SQLiteTimestamp{Timestamp: now},
					UpdatedAt: models.SQLiteTimestamp{Timestamp: now},
				})
				if err != nil {
					return nil, fmt.Errorf("error creating tag: %w", err)
				}

				tagIDs = append(tagIDs, created.ID)
			}
		}
	}

	// don't return if nothing was added
	if idListEquals(originalTagIDs, tagIDs) {
		return nil, nil
	}

	return tagIDs, nil
}

func stashIDListEquals(a, b []models.StashID) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aa := range a {
		found := false
		for _, bb := range b {
			if aa == bb {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (t *IdentifySceneTask) getSceneStashIDs(input modifySceneOptions) ([]models.StashID, error) {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["stash_ids"]
	target := input.scene
	r := t.Repo

	endpoint := input.endpoint

	// just check if ignored
	if endpoint == "" || !t.shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	strategy := models.IdentifyFieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var originalStashIDs []models.StashID
	var stashIDs []models.StashID
	if scraped.RemoteSiteID != nil {
		stashIDPtrs, err := r.Scene().GetStashIDs(target.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting scene tag: %w", err)
		}

		// convert existing to non-pointer types
		for _, stashID := range stashIDPtrs {
			originalStashIDs = append(stashIDs, *stashID)
		}

		if strategy == models.IdentifyFieldStrategyMerge {
			// add to existing
			stashIDs = originalStashIDs
		}

		remoteSiteID := *scraped.RemoteSiteID
		for _, stashID := range stashIDs {
			if endpoint == stashID.Endpoint {
				// if stashID is the same, then don't set
				if stashID.StashID == remoteSiteID {
					return nil, nil
				}

				// replace the stash id and return
				stashID.StashID = remoteSiteID
				return stashIDs, nil
			}
		}

		// not found, create new entry
		stashIDs = append(stashIDs, models.StashID{
			StashID:  remoteSiteID,
			Endpoint: endpoint,
		})
	}

	if stashIDListEquals(originalStashIDs, stashIDs) {
		return nil, nil
	}

	return stashIDs, nil
}

func (t *IdentifySceneTask) getSceneCover(ctx context.Context, input modifySceneOptions) ([]byte, error) {
	scraped := input.scraped
	target := input.scene
	r := t.Repo

	if scraped.Image == nil {
		return nil, nil
	}

	setCoverImage := false
	for _, o := range input.options {
		if o.SetCoverImage != nil {
			setCoverImage = *o.SetCoverImage
			break
		}
	}

	// just check if ignored
	if !setCoverImage {
		return nil, nil
	}

	// always overwrite if present
	existingCover, err := r.Scene().GetCover(target.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene cover: %w", err)
	}

	data, err := utils.ProcessImageInput(ctx, *scraped.Image)
	if err != nil {
		return nil, fmt.Errorf("error processing image input: %w", err)
	}

	// only return if different
	if !bytes.Equal(existingCover, data) {
		return data, nil
	}

	return nil, nil
}

func (t *IdentifySceneTask) shouldSetSingleValueField(strategy *models.IdentifyFieldOptionsInput, hasExistingValue bool) bool {
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

type modifySceneOptions struct {
	scene        *models.Scene
	scraped      *models.ScrapedScene
	endpoint     string
	fieldOptions map[string]*models.IdentifyFieldOptionsInput
	options      []models.IdentifyMetadataOptionsInput
}

func (o *modifySceneOptions) setOptions(sourceOptions *models.IdentifyMetadataOptionsInput, defaultOptions *models.IdentifyMetadataOptionsInput) {
	options := sourceOptions

	// set up options in order of preference
	if options != nil {
		o.options = []models.IdentifyMetadataOptionsInput{*options}
	}
	if defaultOptions != nil {
		o.options = append(o.options, *defaultOptions)
	}

	o.setFieldStrategies()
}

func (o *modifySceneOptions) setFieldStrategies() {
	// prefer source-specific field strategies, then the defaults
	ret := make(map[string]*models.IdentifyFieldOptionsInput)
	for _, oo := range o.options {
		for _, f := range oo.FieldOptions {
			if _, found := ret[f.Field]; !found {
				ret[f.Field] = f
			}
		}
	}

	o.fieldOptions = ret
}
