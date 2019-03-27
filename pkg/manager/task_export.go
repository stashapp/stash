package manager

import (
	"context"
	"fmt"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"math"
	"strconv"
	"sync"
)

type ExportTask struct {
	Mappings *jsonschema.Mappings
	Scraped  []jsonschema.ScrapedItem
}

func (t *ExportTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	// @manager.total = Scene.count + Gallery.count + Performer.count + Studio.count

	t.Mappings = &jsonschema.Mappings{}
	t.Scraped = []jsonschema.ScrapedItem{}

	ctx := context.TODO()

	t.ExportScenes(ctx)
	t.ExportGalleries(ctx)
	t.ExportPerformers(ctx)
	t.ExportStudios(ctx)

	if err := instance.JSON.saveMappings(t.Mappings); err != nil {
		logger.Errorf("[mappings] failed to save json: %s", err.Error())
	}

	t.ExportScrapedItems(ctx)
}

func (t *ExportTask) ExportScenes(ctx context.Context) {
	tx := database.DB.MustBeginTx(ctx, nil)
	defer tx.Commit()
	qb := models.NewSceneQueryBuilder()
	studioQB := models.NewStudioQueryBuilder()
	galleryQB := models.NewGalleryQueryBuilder()
	performerQB := models.NewPerformerQueryBuilder()
	tagQB := models.NewTagQueryBuilder()
	sceneMarkerQB := models.NewSceneMarkerQueryBuilder()
	scenes, err := qb.All()
	if err != nil {
		logger.Errorf("[scenes] failed to fetch all scenes: %s", err.Error())
	}

	logger.Info("[scenes] exporting")

	for i, scene := range scenes {
		index := i + 1
		logger.Progressf("[scenes] %d of %d", index, len(scenes))

		t.Mappings.Scenes = append(t.Mappings.Scenes, jsonschema.PathMapping{Path: scene.Path, Checksum: scene.Checksum})
		newSceneJSON := jsonschema.Scene{
			CreatedAt: models.JSONTime{Time: scene.CreatedAt.Timestamp},
			UpdatedAt: models.JSONTime{Time: scene.UpdatedAt.Timestamp},
		}

		var studioName string
		if scene.StudioID.Valid {
			studio, _ := studioQB.Find(int(scene.StudioID.Int64), tx)
			if studio != nil {
				studioName = studio.Name.String
			}
		}

		var galleryChecksum string
		gallery, _ := galleryQB.FindBySceneID(scene.ID, tx)
		if gallery != nil {
			galleryChecksum = gallery.Checksum
		}

		performers, _ := performerQB.FindBySceneID(scene.ID, tx)
		tags, _ := tagQB.FindBySceneID(scene.ID, tx)
		sceneMarkers, _ := sceneMarkerQB.FindBySceneID(scene.ID, tx)

		if scene.Title.Valid {
			newSceneJSON.Title = scene.Title.String
		}
		if studioName != "" {
			newSceneJSON.Studio = studioName
		}
		if scene.URL.Valid {
			newSceneJSON.URL = scene.URL.String
		}
		if scene.Date.Valid {
			newSceneJSON.Date = utils.GetYMDFromDatabaseDate(scene.Date.String)
		}
		if scene.Rating.Valid {
			newSceneJSON.Rating = int(scene.Rating.Int64)
		}
		if scene.Details.Valid {
			newSceneJSON.Details = scene.Details.String
		}
		if galleryChecksum != "" {
			newSceneJSON.Gallery = galleryChecksum
		}

		newSceneJSON.Performers = t.getPerformerNames(performers)
		newSceneJSON.Tags = t.getTagNames(tags)

		for _, sceneMarker := range sceneMarkers {
			var primaryTagID int
			if sceneMarker.PrimaryTagID.Valid {
				primaryTagID = int(sceneMarker.PrimaryTagID.Int64)
			}
			primaryTag, err := tagQB.Find(primaryTagID, tx)
			if err != nil {
				logger.Errorf("[scenes] <%s> invalid primary tag for scene marker: %s", scene.Checksum, err.Error())
				continue
			}
			sceneMarkerTags, err := tagQB.FindBySceneMarkerID(sceneMarker.ID, tx)
			if err != nil {
				logger.Errorf("[scenes] <%s> invalid tags for scene marker: %s", scene.Checksum, err.Error())
				continue
			}
			if sceneMarker.Title == "" || sceneMarker.Seconds == 0 || primaryTag.Name == "" {
				logger.Errorf("[scenes] invalid scene marker: %v", sceneMarker)
			}

			sceneMarkerJSON := jsonschema.SceneMarker{
				Title:      sceneMarker.Title,
				Seconds:    t.getDecimalString(sceneMarker.Seconds),
				PrimaryTag: primaryTag.Name,
				Tags:       t.getTagNames(sceneMarkerTags),
				CreatedAt:  models.JSONTime{Time: sceneMarker.CreatedAt.Timestamp},
				UpdatedAt:  models.JSONTime{Time: sceneMarker.UpdatedAt.Timestamp},
			}

			newSceneJSON.Markers = append(newSceneJSON.Markers, sceneMarkerJSON)
		}

		newSceneJSON.File = &jsonschema.SceneFile{}
		if scene.Size.Valid {
			newSceneJSON.File.Size = scene.Size.String
		}
		if scene.Duration.Valid {
			newSceneJSON.File.Duration = t.getDecimalString(scene.Duration.Float64)
		}
		if scene.VideoCodec.Valid {
			newSceneJSON.File.VideoCodec = scene.VideoCodec.String
		}
		if scene.AudioCodec.Valid {
			newSceneJSON.File.AudioCodec = scene.AudioCodec.String
		}
		if scene.Width.Valid {
			newSceneJSON.File.Width = int(scene.Width.Int64)
		}
		if scene.Height.Valid {
			newSceneJSON.File.Height = int(scene.Height.Int64)
		}
		if scene.Framerate.Valid {
			newSceneJSON.File.Framerate = t.getDecimalString(scene.Framerate.Float64)
		}
		if scene.Bitrate.Valid {
			newSceneJSON.File.Bitrate = int(scene.Bitrate.Int64)
		}

		sceneJSON, err := instance.JSON.getScene(scene.Checksum)
		if err != nil {
			logger.Debugf("[scenes] error reading scene json: %s", err.Error())
		} else if jsonschema.CompareJSON(*sceneJSON, newSceneJSON) {
			continue
		}

		if err := instance.JSON.saveScene(scene.Checksum, &newSceneJSON); err != nil {
			logger.Errorf("[scenes] <%s> failed to save json: %s", scene.Checksum, err.Error())
		}
	}

	logger.Infof("[scenes] export complete")
}

func (t *ExportTask) ExportGalleries(ctx context.Context) {
	qb := models.NewGalleryQueryBuilder()
	galleries, err := qb.All()
	if err != nil {
		logger.Errorf("[galleries] failed to fetch all galleries: %s", err.Error())
	}

	logger.Info("[galleries] exporting")

	for i, gallery := range galleries {
		index := i + 1
		logger.Progressf("[galleries] %d of %d", index, len(galleries))
		t.Mappings.Galleries = append(t.Mappings.Galleries, jsonschema.PathMapping{Path: gallery.Path, Checksum: gallery.Checksum})
	}

	logger.Infof("[galleries] export complete")
}

func (t *ExportTask) ExportPerformers(ctx context.Context) {
	qb := models.NewPerformerQueryBuilder()
	performers, err := qb.All()
	if err != nil {
		logger.Errorf("[performers] failed to fetch all performers: %s", err.Error())
	}

	logger.Info("[performers] exporting")

	for i, performer := range performers {
		index := i + 1
		logger.Progressf("[performers] %d of %d", index, len(performers))

		t.Mappings.Performers = append(t.Mappings.Performers, jsonschema.NameMapping{Name: performer.Name.String, Checksum: performer.Checksum})

		newPerformerJSON := jsonschema.Performer{
			CreatedAt: models.JSONTime{Time: performer.CreatedAt.Timestamp},
			UpdatedAt: models.JSONTime{Time: performer.UpdatedAt.Timestamp},
		}

		if performer.Name.Valid {
			newPerformerJSON.Name = performer.Name.String
		}
		if performer.URL.Valid {
			newPerformerJSON.URL = performer.URL.String
		}
		if performer.Birthdate.Valid {
			newPerformerJSON.Birthdate = utils.GetYMDFromDatabaseDate(performer.Birthdate.String)
		}
		if performer.Ethnicity.Valid {
			newPerformerJSON.Ethnicity = performer.Ethnicity.String
		}
		if performer.Country.Valid {
			newPerformerJSON.Country = performer.Country.String
		}
		if performer.EyeColor.Valid {
			newPerformerJSON.EyeColor = performer.EyeColor.String
		}
		if performer.Height.Valid {
			newPerformerJSON.Height = performer.Height.String
		}
		if performer.Measurements.Valid {
			newPerformerJSON.Measurements = performer.Measurements.String
		}
		if performer.FakeTits.Valid {
			newPerformerJSON.FakeTits = performer.FakeTits.String
		}
		if performer.CareerLength.Valid {
			newPerformerJSON.CareerLength = performer.CareerLength.String
		}
		if performer.Tattoos.Valid {
			newPerformerJSON.Tattoos = performer.Tattoos.String
		}
		if performer.Piercings.Valid {
			newPerformerJSON.Piercings = performer.Piercings.String
		}
		if performer.Aliases.Valid {
			newPerformerJSON.Aliases = performer.Aliases.String
		}
		if performer.Twitter.Valid {
			newPerformerJSON.Twitter = performer.Twitter.String
		}
		if performer.Instagram.Valid {
			newPerformerJSON.Instagram = performer.Instagram.String
		}
		if performer.Favorite.Valid {
			newPerformerJSON.Favorite = performer.Favorite.Bool
		}

		newPerformerJSON.Image = utils.GetBase64StringFromData(performer.Image)

		performerJSON, err := instance.JSON.getPerformer(performer.Checksum)
		if err != nil {
			logger.Debugf("[performers] error reading performer json: %s", err.Error())
		} else if jsonschema.CompareJSON(*performerJSON, newPerformerJSON) {
			continue
		}

		if err := instance.JSON.savePerformer(performer.Checksum, &newPerformerJSON); err != nil {
			logger.Errorf("[performers] <%s> failed to save json: %s", performer.Checksum, err.Error())
		}
	}

	logger.Infof("[performers] export complete")
}

func (t *ExportTask) ExportStudios(ctx context.Context) {
	qb := models.NewStudioQueryBuilder()
	studios, err := qb.All()
	if err != nil {
		logger.Errorf("[studios] failed to fetch all studios: %s", err.Error())
	}

	logger.Info("[studios] exporting")

	for i, studio := range studios {
		index := i + 1
		logger.Progressf("[studios] %d of %d", index, len(studios))

		t.Mappings.Studios = append(t.Mappings.Studios, jsonschema.NameMapping{Name: studio.Name.String, Checksum: studio.Checksum})

		newStudioJSON := jsonschema.Studio{
			CreatedAt: models.JSONTime{Time: studio.CreatedAt.Timestamp},
			UpdatedAt: models.JSONTime{Time: studio.UpdatedAt.Timestamp},
		}

		if studio.Name.Valid {
			newStudioJSON.Name = studio.Name.String
		}
		if studio.URL.Valid {
			newStudioJSON.URL = studio.URL.String
		}

		newStudioJSON.Image = utils.GetBase64StringFromData(studio.Image)

		studioJSON, err := instance.JSON.getStudio(studio.Checksum)
		if err != nil {
			logger.Debugf("[studios] error reading studio json: %s", err.Error())
		} else if jsonschema.CompareJSON(*studioJSON, newStudioJSON) {
			continue
		}

		if err := instance.JSON.saveStudio(studio.Checksum, &newStudioJSON); err != nil {
			logger.Errorf("[studios] <%s> failed to save json: %s", studio.Checksum, err.Error())
		}
	}

	logger.Infof("[studios] export complete")
}

func (t *ExportTask) ExportScrapedItems(ctx context.Context) {
	tx := database.DB.MustBeginTx(ctx, nil)
	defer tx.Commit()
	qb := models.NewScrapedItemQueryBuilder()
	sqb := models.NewStudioQueryBuilder()
	scrapedItems, err := qb.All()
	if err != nil {
		logger.Errorf("[scraped sites] failed to fetch all items: %s", err.Error())
	}

	logger.Info("[scraped sites] exporting")

	for i, scrapedItem := range scrapedItems {
		index := i + 1
		logger.Progressf("[scraped sites] %d of %d", index, len(scrapedItems))

		var studioName string
		if scrapedItem.StudioID.Valid {
			studio, _ := sqb.Find(int(scrapedItem.StudioID.Int64), tx)
			if studio != nil {
				studioName = studio.Name.String
			}
		}

		newScrapedItemJSON := jsonschema.ScrapedItem{}

		if scrapedItem.Title.Valid {
			newScrapedItemJSON.Title = scrapedItem.Title.String
		}
		if scrapedItem.Description.Valid {
			newScrapedItemJSON.Description = scrapedItem.Description.String
		}
		if scrapedItem.URL.Valid {
			newScrapedItemJSON.URL = scrapedItem.URL.String
		}
		if scrapedItem.Date.Valid {
			newScrapedItemJSON.Date = utils.GetYMDFromDatabaseDate(scrapedItem.Date.String)
		}
		if scrapedItem.Rating.Valid {
			newScrapedItemJSON.Rating = scrapedItem.Rating.String
		}
		if scrapedItem.Tags.Valid {
			newScrapedItemJSON.Tags = scrapedItem.Tags.String
		}
		if scrapedItem.Models.Valid {
			newScrapedItemJSON.Models = scrapedItem.Models.String
		}
		if scrapedItem.Episode.Valid {
			newScrapedItemJSON.Episode = int(scrapedItem.Episode.Int64)
		}
		if scrapedItem.GalleryFilename.Valid {
			newScrapedItemJSON.GalleryFilename = scrapedItem.GalleryFilename.String
		}
		if scrapedItem.GalleryURL.Valid {
			newScrapedItemJSON.GalleryURL = scrapedItem.GalleryURL.String
		}
		if scrapedItem.VideoFilename.Valid {
			newScrapedItemJSON.VideoFilename = scrapedItem.VideoFilename.String
		}
		if scrapedItem.VideoURL.Valid {
			newScrapedItemJSON.VideoURL = scrapedItem.VideoURL.String
		}

		newScrapedItemJSON.Studio = studioName
		updatedAt := models.JSONTime{Time: scrapedItem.UpdatedAt.Timestamp} // TODO keeping ruby format
		newScrapedItemJSON.UpdatedAt = updatedAt

		t.Scraped = append(t.Scraped, newScrapedItemJSON)
	}

	scrapedJSON, err := instance.JSON.getScraped()
	if err != nil {
		logger.Debugf("[scraped sites] error reading json: %s", err.Error())
	}
	if !jsonschema.CompareJSON(scrapedJSON, t.Scraped) {
		if err := instance.JSON.saveScaped(t.Scraped); err != nil {
			logger.Errorf("[scraped sites] failed to save json: %s", err.Error())
		}
	}

	logger.Infof("[scraped sites] export complete")
}

func (t *ExportTask) getPerformerNames(performers []models.Performer) []string {
	if len(performers) == 0 {
		return nil
	}

	var results []string
	for _, performer := range performers {
		if performer.Name.Valid {
			results = append(results, performer.Name.String)
		}
	}

	return results
}

func (t *ExportTask) getTagNames(tags []models.Tag) []string {
	if len(tags) == 0 {
		return nil
	}

	var results []string
	for _, tag := range tags {
		if tag.Name != "" {
			results = append(results, tag.Name)
		}
	}

	return results
}

func (t *ExportTask) getDecimalString(num float64) string {
	if num == 0 {
		return ""
	}

	precision := getPrecision(num)
	if precision == 0 {
		precision = 1
	}
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num)
}

func getPrecision(num float64) int {
	if num == 0 {
		return 0
	}

	e := 1.0
	p := 0
	for (math.Round(num*e) / e) != num {
		e *= 10
		p++
	}
	return p
}
