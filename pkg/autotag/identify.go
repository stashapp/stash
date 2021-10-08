package autotag

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
	"github.com/stashapp/stash/pkg/utils"
)

type IdentifySceneTask struct {
	Input   models.IdentifyMetadataInput
	SceneID int

	Ctx          context.Context
	TxnManager   models.TransactionManager
	ScraperCache *scraper.Cache
	StashBoxes   models.StashBoxes
}

func (t *IdentifySceneTask) Execute(ctx context.Context, progress *job.Progress) {
	var scene *models.Scene

	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		// find the scene
		var err error
		scene, err = r.Scene().Find(t.SceneID)
		if err != nil {
			return fmt.Errorf("error finding scene with id %d: %w", t.SceneID, err)
		}

		if scene == nil {
			return fmt.Errorf("no scene found with id %d", t.SceneID)
		}

		progress.ExecuteTask("Identifying "+scene.Path, func() {
			// iterate through the input sources
			// TODO - use default sources if not provided
			for _, source := range t.Input.Sources {
				var stashBox *models.StashBox
				stashBox, err = t.getStashBox(source.Source)

				// scrape using the source
				var scraped *models.ScrapedScene
				if stashBox != nil {
					scraped, err = t.scrapeUsingStashBox(stashBox)
				} else {
					scraped, err = t.ScraperCache.ScrapeScene(*source.Source.ScraperID, t.SceneID)
				}

				if err != nil {
					return
				}

				// if results were found then modify the scene
				if scraped != nil {
					options := modifySceneOptions{
						scene:    scene,
						scraped:  scraped,
						stashBox: stashBox,
						repo:     r,
					}

					options.setOptions(*source, t.Input.Options)

					err = t.modifyScene(r, options)
					return
				}
			}

			logger.Infof("Unable to identify %s", scene.Path)
		})

		return err
	}); err != nil {
		if scene == nil {
			logger.Error(err.Error())
		} else {
			logger.Errorf("Error encountered identifying %s: %v", scene.Path, err)
		}
	}
}

func (t *IdentifySceneTask) getStashBox(scraperSrc *models.ScraperSourceInput) (*models.StashBox, error) {
	if scraperSrc.ScraperID != nil {
		return nil, nil
	}

	// must be stash-box
	if scraperSrc.StashBoxIndex == nil && scraperSrc.StashBoxEndpoint == nil {
		return nil, errors.New("stash_box_index or stash_box_endpoint or scraper_id must be set")
	}

	return t.StashBoxes.ResolveStashBox(*scraperSrc)
}

func (t *IdentifySceneTask) scrapeUsingStashBox(box *models.StashBox) (*models.ScrapedScene, error) {
	client := stashbox.NewClient(*box, t.TxnManager)
	results, err := client.FindStashBoxScenesByFingerprintsFlat([]string{strconv.Itoa(t.SceneID)})
	if err != nil {
		return nil, fmt.Errorf("error querying stash-box using scene ID %d: %w", t.SceneID, err)
	}

	if len(results) > 0 {
		return results[0], nil
	}

	return nil, nil
}

func (t *IdentifySceneTask) modifyScene(r models.Repository, options modifySceneOptions) error {
	target := options.scene

	partial := models.ScenePartial{
		ID: target.ID,
	}

	t.setSceneFields(&partial, options)

	if err := t.setSceneStudio(&partial, options); err != nil {
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

	coverImage, err := t.getSceneCover(options)
	if err != nil {
		return err
	}

	qb := r.Scene()
	_, err = qb.Update(partial)
	if err != nil {
		return fmt.Errorf("error updating scene: %w", err)
	}

	if len(performerIDs) > 0 {
		if err := qb.UpdatePerformers(t.SceneID, performerIDs); err != nil {
			return fmt.Errorf("error updating scene performers: %w", err)
		}
	}

	if len(tagIDs) > 0 {
		if err := qb.UpdateTags(t.SceneID, tagIDs); err != nil {
			return fmt.Errorf("error updating scene tags: %w", err)
		}
	}

	if len(stashIDs) > 0 {
		if err := qb.UpdateStashIDs(t.SceneID, stashIDs); err != nil {
			return fmt.Errorf("error updating scene stash_ids: %w", err)
		}
	}

	if len(coverImage) > 0 {
		if err := qb.UpdateCover(t.SceneID, coverImage); err != nil {
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

func (t *IdentifySceneTask) setSceneFields(partial *models.ScenePartial, input modifySceneOptions) {
	scraped := input.scraped
	fieldStrategies := input.fieldOptions
	target := input.scene

	if scraped.Title != nil {
		if t.shouldSetSingleValueField(fieldStrategies["title"], target.Title.String != "") {
			partial.Title = models.NullStringPtr(*scraped.Title)
		}
	}
	if scraped.Date != nil {
		if t.shouldSetSingleValueField(fieldStrategies["date"], target.Date.Valid) {
			partial.Date = &models.SQLiteDate{
				String: *scraped.Date,
			}
		}
	}
	if scraped.Details != nil {
		if t.shouldSetSingleValueField(fieldStrategies["details"], target.Details.String != "") {
			partial.Title = models.NullStringPtr(*scraped.Details)
		}
	}
	if scraped.URL != nil {
		if t.shouldSetSingleValueField(fieldStrategies["url"], target.URL.String != "") {
			partial.URL = models.NullStringPtr(*scraped.URL)
		}
	}

	setOrganized := false
	for _, o := range input.options {
		if o.SetOrganized != nil {
			setOrganized = *o.SetOrganized
			break
		}
	}
	if setOrganized {
		// just reuse the boolean since we know it's true
		partial.Organized = &setOrganized
	}

}

func (t *IdentifySceneTask) setSceneStudio(partial *models.ScenePartial, input modifySceneOptions) error {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["studio"]
	target := input.scene
	r := input.repo

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)

	if scraped.Studio != nil {
		if t.shouldSetSingleValueField(fieldStrategy, target.StudioID.Valid) {
			if scraped.Studio.StoredID != nil {
				// existing studio, just set it
				studioID, err := strconv.ParseInt(*scraped.Studio.StoredID, 10, 64)
				if err != nil {
					return fmt.Errorf("error converting studio ID %s: %w", *scraped.Studio.StoredID, err)
				}

				target.StudioID = sql.NullInt64{
					Int64: studioID,
					Valid: true,
				}
			} else if createMissing {
				created, err := r.Studio().Create(scrapedToStudioInput(scraped.Studio))
				if err != nil {
					return fmt.Errorf("error creating studio: %w", err)
				}

				if input.stashBox != nil && scraped.RemoteSiteID != nil {
					if err := r.Studio().UpdateStashIDs(created.ID, []models.StashID{
						{
							Endpoint: input.stashBox.Endpoint,
							StashID:  *scraped.Studio.RemoteSiteID,
						},
					}); err != nil {
						return fmt.Errorf("error setting studio stash id: %w", err)
					}
				}

				target.StudioID = sql.NullInt64{
					Int64: int64(created.ID),
					Valid: true,
				}
			}
		}
	}

	return nil
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

func (t *IdentifySceneTask) getScenePerformerIDs(input modifySceneOptions) ([]int, error) {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["performers"]
	target := input.scene
	r := input.repo

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

	var performerIDs []int
	if len(scraped.Performers) > 0 {
		var err error

		if strategy == models.IdentifyFieldStrategyMerge {
			// add to existing
			performerIDs, err = r.Scene().GetPerformerIDs(target.ID)
			if err != nil {
				return nil, fmt.Errorf("error getting scene performers: %w", err)
			}
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

				if input.stashBox != nil && scraped.RemoteSiteID != nil {
					if err := r.Performer().UpdateStashIDs(created.ID, []models.StashID{
						{
							Endpoint: input.stashBox.Endpoint,
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
	r := input.repo

	// just check if ignored
	if !t.shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := models.IdentifyFieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var tagIDs []int
	if len(scraped.Tags) > 0 {
		var err error

		if strategy == models.IdentifyFieldStrategyMerge {
			// add to existing
			tagIDs, err = r.Scene().GetTagIDs(target.ID)
			if err != nil {
				return nil, fmt.Errorf("error getting scene tag: %w", err)
			}
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

	return tagIDs, nil
}

func (t *IdentifySceneTask) getSceneStashIDs(input modifySceneOptions) ([]models.StashID, error) {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["stash_ids"]
	target := input.scene
	r := input.repo

	stashBox := input.stashBox

	// just check if ignored
	if stashBox == nil || !t.shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	strategy := models.IdentifyFieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var stashIDs []models.StashID
	if scraped.RemoteSiteID != nil {
		var err error

		if strategy == models.IdentifyFieldStrategyMerge {
			// add to existing
			var stashIDPtrs []*models.StashID
			stashIDPtrs, err = r.Scene().GetStashIDs(target.ID)
			if err != nil {
				return nil, fmt.Errorf("error getting scene tag: %w", err)
			}

			// convert existing to non-pointer types
			for _, stashID := range stashIDPtrs {
				stashIDs = append(stashIDs, *stashID)
			}
		}

		remoteSiteID := *scraped.RemoteSiteID
		for _, stashID := range stashIDs {
			if stashBox.Endpoint == stashID.Endpoint {
				// replace the stash id and return
				stashID.Endpoint = remoteSiteID
				return stashIDs, nil
			}
		}

		// not found, create new entry
		stashIDs = append(stashIDs, models.StashID{
			StashID:  remoteSiteID,
			Endpoint: stashBox.Endpoint,
		})
	}

	return stashIDs, nil
}

func (t *IdentifySceneTask) getSceneCover(input modifySceneOptions) ([]byte, error) {
	scraped := input.scraped
	fieldStrategy := input.fieldOptions["cover_image"]

	// just check if ignored
	if !t.shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	if scraped.Image != nil {
		// always overwrite if present
		return utils.ProcessImageInput(*scraped.Image)
	}

	return nil, nil
}

func (t *IdentifySceneTask) shouldSetSingleValueField(strategy *models.IdentifyFieldOptions, hasExistingValue bool) bool {
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
	stashBox     *models.StashBox
	fieldOptions map[string]*models.IdentifyFieldOptions
	options      []models.IdentifyMetadataOptionsInput
	repo         models.Repository
}

func (o *modifySceneOptions) setOptions(source models.IdentifySourceInput, defaultOptions *models.IdentifyMetadataOptionsInput) {
	options := source.Options

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
	ret := make(map[string]*models.IdentifyFieldOptions)
	for _, oo := range o.options {
		for _, f := range oo.FieldOptions {
			if _, found := ret[f.Field]; !found {
				ret[f.Field] = f
			}
		}
	}

	o.fieldOptions = ret
}
