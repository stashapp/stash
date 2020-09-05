package manager

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/movie"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/utils"
)

type ImportTask struct {
	json jsonUtils

	BaseDir             string
	Reset               bool
	DuplicateBehaviour  models.ImportDuplicateEnum
	MissingRefBehaviour models.ImportMissingRefEnum

	mappings            *jsonschema.Mappings
	scraped             []jsonschema.ScrapedItem
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *ImportTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.json = jsonUtils{
		json: *paths.GetJSONPaths(t.BaseDir),
	}

	// set default behaviour if not provided
	if !t.DuplicateBehaviour.IsValid() {
		t.DuplicateBehaviour = models.ImportDuplicateEnumFail
	}
	if !t.MissingRefBehaviour.IsValid() {
		t.MissingRefBehaviour = models.ImportMissingRefEnumFail
	}

	t.mappings, _ = t.json.getMappings()
	if t.mappings == nil {
		logger.Error("missing mappings json")
		return
	}
	scraped, _ := t.json.getScraped()
	if scraped == nil {
		logger.Warn("missing scraped json")
	}
	t.scraped = scraped

	if t.Reset {
		err := database.Reset(config.GetDatabasePath())

		if err != nil {
			logger.Errorf("Error resetting database: %s", err.Error())
			return
		}
	}

	ctx := context.TODO()

	t.ImportTags(ctx)
	t.ImportPerformers(ctx)
	t.ImportStudios(ctx)
	t.ImportMovies(ctx)
	t.ImportGalleries(ctx)

	t.ImportScrapedItems(ctx)
	t.ImportScenes(ctx)
}

func (t *ImportTask) ImportPerformers(ctx context.Context) {
	logger.Info("[performers] importing")

	for i, mappingJSON := range t.mappings.Performers {
		index := i + 1
		performerJSON, err := t.json.getPerformer(mappingJSON.Checksum)
		if err != nil {
			logger.Errorf("[performers] failed to read json: %s", err.Error())
			continue
		}

		logger.Progressf("[performers] %d of %d", index, len(t.mappings.Performers))

		tx := database.DB.MustBeginTx(ctx, nil)
		readerWriter := models.NewPerformerReaderWriter(tx)
		importer := &performer.Importer{
			ReaderWriter: readerWriter,
			Input:        *performerJSON,
		}

		if err := performImport(importer, t.DuplicateBehaviour); err != nil {
			tx.Rollback()
			logger.Errorf("[performers] <%s> failed to import: %s", mappingJSON.Checksum, err.Error())
			continue
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Errorf("[performers] <%s> import failed to commit: %s", mappingJSON.Checksum, err.Error())
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Errorf("[performers] <%s> import failed to commit: %s", mappingJSON.Checksum, err.Error())
		}
	}

	logger.Info("[performers] import complete")
}

func (t *ImportTask) ImportStudios(ctx context.Context) {
	pendingParent := make(map[string][]*jsonschema.Studio)

	logger.Info("[studios] importing")

	for i, mappingJSON := range t.mappings.Studios {
		index := i + 1
		studioJSON, err := t.json.getStudio(mappingJSON.Checksum)
		if err != nil {
			logger.Errorf("[studios] failed to read json: %s", err.Error())
			continue
		}

		logger.Progressf("[studios] %d of %d", index, len(t.mappings.Studios))

		tx := database.DB.MustBeginTx(ctx, nil)

		// fail on missing parent studio to begin with
		if err := t.ImportStudio(studioJSON, pendingParent, tx); err != nil {
			tx.Rollback()

			if err == studio.ErrParentStudioNotExist {
				// add to the pending parent list so that it is created after the parent
				s := pendingParent[studioJSON.ParentStudio]
				s = append(s, studioJSON)
				pendingParent[studioJSON.ParentStudio] = s
				continue
			} else {
				logger.Errorf("[studios] <%s> failed to create: %s", mappingJSON.Checksum, err.Error())
				continue
			}
		}

		if err := tx.Commit(); err != nil {
			logger.Errorf("[studios] import failed to commit: %s", err.Error())
			continue
		}
	}

	// create the leftover studios, warning for missing parents
	if len(pendingParent) > 0 {
		logger.Warnf("[studios] importing studios with missing parents")

		for _, s := range pendingParent {
			for _, orphanStudioJSON := range s {
				tx := database.DB.MustBeginTx(ctx, nil)

				if err := t.ImportStudio(orphanStudioJSON, nil, tx); err != nil {
					tx.Rollback()
					logger.Errorf("[studios] <%s> failed to create: %s", orphanStudioJSON.Name, err.Error())
					continue
				}

				if err := tx.Commit(); err != nil {
					logger.Errorf("[studios] import failed to commit: %s", err.Error())
					continue
				}
			}
		}
	}

	logger.Info("[studios] import complete")
}

func (t *ImportTask) ImportStudio(studioJSON *jsonschema.Studio, pendingParent map[string][]*jsonschema.Studio, tx *sqlx.Tx) error {
	readerWriter := models.NewStudioReaderWriter(tx)
	importer := &studio.Importer{
		ReaderWriter:        readerWriter,
		Input:               *studioJSON,
		MissingRefBehaviour: t.MissingRefBehaviour,
	}

	// first phase: return error if parent does not exist
	if pendingParent != nil {
		importer.MissingRefBehaviour = models.ImportMissingRefEnumFail
	}

	if err := performImport(importer, t.DuplicateBehaviour); err != nil {
		return err
	}

	// now create the studios pending this studios creation
	s := pendingParent[studioJSON.Name]
	for _, childStudioJSON := range s {
		// map is nil since we're not checking parent studios at this point
		if err := t.ImportStudio(childStudioJSON, nil, tx); err != nil {
			return fmt.Errorf("failed to create child studio <%s>: %s", childStudioJSON.Name, err.Error())
		}
	}

	// delete the entry from the map so that we know its not left over
	delete(pendingParent, studioJSON.Name)

	return nil
}

func (t *ImportTask) ImportMovies(ctx context.Context) {
	logger.Info("[movies] importing")

	for i, mappingJSON := range t.mappings.Movies {
		index := i + 1
		movieJSON, err := t.json.getMovie(mappingJSON.Checksum)
		if err != nil {
			logger.Errorf("[movies] failed to read json: %s", err.Error())
			continue
		}

		logger.Progressf("[movies] %d of %d", index, len(t.mappings.Movies))

		tx := database.DB.MustBeginTx(ctx, nil)
		readerWriter := models.NewMovieReaderWriter(tx)
		studioReaderWriter := models.NewStudioReaderWriter(tx)

		movieImporter := &movie.Importer{
			ReaderWriter:        readerWriter,
			StudioWriter:        studioReaderWriter,
			Input:               *movieJSON,
			MissingRefBehaviour: t.MissingRefBehaviour,
		}

		if err := performImport(movieImporter, t.DuplicateBehaviour); err != nil {
			tx.Rollback()
			logger.Errorf("[movies] <%s> failed to import: %s", mappingJSON.Checksum, err.Error())
			continue
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Errorf("[movies] <%s> import failed to commit: %s", mappingJSON.Checksum, err.Error())
			continue
		}
	}

	logger.Info("[movies] import complete")
}

func (t *ImportTask) ImportGalleries(ctx context.Context) {
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewGalleryQueryBuilder()

	for i, mappingJSON := range t.mappings.Galleries {
		index := i + 1
		if mappingJSON.Checksum == "" || mappingJSON.Path == "" {
			return
		}

		logger.Progressf("[galleries] %d of %d", index, len(t.mappings.Galleries))

		// Populate a new gallery from the input
		currentTime := time.Now()
		newGallery := models.Gallery{
			Checksum:  mappingJSON.Checksum,
			Path:      mappingJSON.Path,
			CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		}

		_, err := qb.Create(newGallery, tx)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("[galleries] <%s> failed to create: %s", mappingJSON.Checksum, err.Error())
			return
		}
	}

	logger.Info("[galleries] importing")
	if err := tx.Commit(); err != nil {
		logger.Errorf("[galleries] import failed to commit: %s", err.Error())
	}
	logger.Info("[galleries] import complete")
}

func (t *ImportTask) ImportTags(ctx context.Context) {
	logger.Info("[tags] importing")

	for i, mappingJSON := range t.mappings.Tags {
		index := i + 1
		tagJSON, err := t.json.getTag(mappingJSON.Checksum)
		if err != nil {
			logger.Errorf("[tags] failed to read json: %s", err.Error())
			continue
		}

		logger.Progressf("[tags] %d of %d", index, len(t.mappings.Tags))

		tx := database.DB.MustBeginTx(ctx, nil)
		readerWriter := models.NewTagReaderWriter(tx)

		tagImporter := &tag.Importer{
			ReaderWriter: readerWriter,
			Input:        *tagJSON,
		}

		if err := performImport(tagImporter, t.DuplicateBehaviour); err != nil {
			tx.Rollback()
			logger.Errorf("[tags] <%s> failed to import: %s", mappingJSON.Checksum, err.Error())
			continue
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Errorf("[tags] <%s> import failed to commit: %s", mappingJSON.Checksum, err.Error())
		}
	}

	logger.Info("[tags] import complete")
}

func (t *ImportTask) ImportScrapedItems(ctx context.Context) {
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewScrapedItemQueryBuilder()
	sqb := models.NewStudioQueryBuilder()
	currentTime := time.Now()

	for i, mappingJSON := range t.scraped {
		index := i + 1
		logger.Progressf("[scraped sites] %d of %d", index, len(t.mappings.Scenes))

		newScrapedItem := models.ScrapedItem{
			Title:           sql.NullString{String: mappingJSON.Title, Valid: true},
			Description:     sql.NullString{String: mappingJSON.Description, Valid: true},
			URL:             sql.NullString{String: mappingJSON.URL, Valid: true},
			Date:            models.SQLiteDate{String: mappingJSON.Date, Valid: true},
			Rating:          sql.NullString{String: mappingJSON.Rating, Valid: true},
			Tags:            sql.NullString{String: mappingJSON.Tags, Valid: true},
			Models:          sql.NullString{String: mappingJSON.Models, Valid: true},
			Episode:         sql.NullInt64{Int64: int64(mappingJSON.Episode), Valid: true},
			GalleryFilename: sql.NullString{String: mappingJSON.GalleryFilename, Valid: true},
			GalleryURL:      sql.NullString{String: mappingJSON.GalleryURL, Valid: true},
			VideoFilename:   sql.NullString{String: mappingJSON.VideoFilename, Valid: true},
			VideoURL:        sql.NullString{String: mappingJSON.VideoURL, Valid: true},
			CreatedAt:       models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:       models.SQLiteTimestamp{Timestamp: t.getTimeFromJSONTime(mappingJSON.UpdatedAt)},
		}

		studio, err := sqb.FindByName(mappingJSON.Studio, tx, false)
		if err != nil {
			logger.Errorf("[scraped sites] failed to fetch studio: %s", err.Error())
		}
		if studio != nil {
			newScrapedItem.StudioID = sql.NullInt64{Int64: int64(studio.ID), Valid: true}
		}

		_, err = qb.Create(newScrapedItem, tx)
		if err != nil {
			logger.Errorf("[scraped sites] <%s> failed to create: %s", newScrapedItem.Title.String, err.Error())
		}
	}

	logger.Info("[scraped sites] importing")
	if err := tx.Commit(); err != nil {
		logger.Errorf("[scraped sites] import failed to commit: %s", err.Error())
	}
	logger.Info("[scraped sites] import complete")
}

func (t *ImportTask) ImportScenes(ctx context.Context) {
	tx := database.DB.MustBeginTx(ctx, nil)
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	for i, mappingJSON := range t.mappings.Scenes {
		index := i + 1
		if mappingJSON.Checksum == "" || mappingJSON.Path == "" {
			_ = tx.Rollback()
			logger.Warn("[scenes] scene mapping without checksum or path: ", mappingJSON)
			return
		}

		logger.Progressf("[scenes] %d of %d", index, len(t.mappings.Scenes))

		sceneJSON, err := t.json.getScene(mappingJSON.Checksum)
		if err != nil {
			logger.Infof("[scenes] <%s> json parse failure: %s", mappingJSON.Checksum, err.Error())
			continue
		}

		sceneHash := mappingJSON.Checksum

		newScene := models.Scene{
			Checksum: sql.NullString{String: sceneJSON.Checksum, Valid: sceneJSON.Checksum != ""},
			OSHash:   sql.NullString{String: sceneJSON.OSHash, Valid: sceneJSON.OSHash != ""},
			Path:     mappingJSON.Path,
		}

		// Process the base 64 encoded cover image string
		var coverImageData []byte
		if sceneJSON.Cover != "" {
			_, coverImageData, err = utils.ProcessBase64Image(sceneJSON.Cover)
			if err != nil {
				logger.Warnf("[scenes] <%s> invalid cover image: %s", sceneHash, err.Error())
			}
			if len(coverImageData) > 0 {
				if err = SetSceneScreenshot(sceneHash, coverImageData); err != nil {
					logger.Warnf("[scenes] <%s> failed to create cover image: %s", sceneHash, err.Error())
				}

				// write the cover image data after creating the scene
			}
		}

		// Populate scene fields
		if sceneJSON != nil {
			if sceneJSON.Title != "" {
				newScene.Title = sql.NullString{String: sceneJSON.Title, Valid: true}
			}
			if sceneJSON.Details != "" {
				newScene.Details = sql.NullString{String: sceneJSON.Details, Valid: true}
			}
			if sceneJSON.URL != "" {
				newScene.URL = sql.NullString{String: sceneJSON.URL, Valid: true}
			}
			if sceneJSON.Date != "" {
				newScene.Date = models.SQLiteDate{String: sceneJSON.Date, Valid: true}
			}
			if sceneJSON.Rating != 0 {
				newScene.Rating = sql.NullInt64{Int64: int64(sceneJSON.Rating), Valid: true}
			}

			newScene.OCounter = sceneJSON.OCounter
			newScene.CreatedAt = models.SQLiteTimestamp{Timestamp: t.getTimeFromJSONTime(sceneJSON.CreatedAt)}
			newScene.UpdatedAt = models.SQLiteTimestamp{Timestamp: t.getTimeFromJSONTime(sceneJSON.UpdatedAt)}

			if sceneJSON.File != nil {
				if sceneJSON.File.Size != "" {
					newScene.Size = sql.NullString{String: sceneJSON.File.Size, Valid: true}
				}
				if sceneJSON.File.Duration != "" {
					duration, _ := strconv.ParseFloat(sceneJSON.File.Duration, 64)
					newScene.Duration = sql.NullFloat64{Float64: duration, Valid: true}
				}
				if sceneJSON.File.VideoCodec != "" {
					newScene.VideoCodec = sql.NullString{String: sceneJSON.File.VideoCodec, Valid: true}
				}
				if sceneJSON.File.AudioCodec != "" {
					newScene.AudioCodec = sql.NullString{String: sceneJSON.File.AudioCodec, Valid: true}
				}
				if sceneJSON.File.Format != "" {
					newScene.Format = sql.NullString{String: sceneJSON.File.Format, Valid: true}
				}
				if sceneJSON.File.Width != 0 {
					newScene.Width = sql.NullInt64{Int64: int64(sceneJSON.File.Width), Valid: true}
				}
				if sceneJSON.File.Height != 0 {
					newScene.Height = sql.NullInt64{Int64: int64(sceneJSON.File.Height), Valid: true}
				}
				if sceneJSON.File.Framerate != "" {
					framerate, _ := strconv.ParseFloat(sceneJSON.File.Framerate, 64)
					newScene.Framerate = sql.NullFloat64{Float64: framerate, Valid: true}
				}
				if sceneJSON.File.Bitrate != 0 {
					newScene.Bitrate = sql.NullInt64{Int64: int64(sceneJSON.File.Bitrate), Valid: true}
				}
			} else {
				// TODO: Get FFMPEG data?
			}
		}

		// Populate the studio ID
		if sceneJSON.Studio != "" {
			sqb := models.NewStudioQueryBuilder()
			studio, err := sqb.FindByName(sceneJSON.Studio, tx, false)
			if err != nil {
				logger.Warnf("[scenes] error getting studio <%s>: %s", sceneJSON.Studio, err.Error())
			} else if studio == nil {
				logger.Warnf("[scenes] studio <%s> does not exist", sceneJSON.Studio)
			} else {
				newScene.StudioID = sql.NullInt64{Int64: int64(studio.ID), Valid: true}
			}
		}

		// Create the scene in the DB
		scene, err := qb.Create(newScene, tx)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("[scenes] <%s> failed to create: %s", sceneHash, err.Error())
			return
		}
		if scene.ID == 0 {
			_ = tx.Rollback()
			logger.Errorf("[scenes] <%s> invalid id after scene creation", sceneHash)
			return
		}

		// Add the scene cover if set
		if len(coverImageData) > 0 {
			if err := qb.UpdateSceneCover(scene.ID, coverImageData, tx); err != nil {
				_ = tx.Rollback()
				logger.Errorf("[scenes] <%s> error setting scene cover: %s", sceneHash, err.Error())
				return
			}
		}

		// Relate the scene to the gallery
		if sceneJSON.Gallery != "" {
			gqb := models.NewGalleryQueryBuilder()
			gallery, err := gqb.FindByChecksum(sceneJSON.Gallery, tx)
			if err != nil {
				logger.Warnf("[scenes] gallery <%s> does not exist: %s", sceneJSON.Gallery, err.Error())
			} else {
				gallery.SceneID = sql.NullInt64{Int64: int64(scene.ID), Valid: true}
				_, err := gqb.Update(*gallery, tx)
				if err != nil {
					logger.Errorf("[scenes] <%s> failed to update gallery: %s", sceneHash, err.Error())
				}
			}
		}

		// Relate the scene to the performers
		if len(sceneJSON.Performers) > 0 {
			performers, err := t.getPerformers(sceneJSON.Performers, tx)
			if err != nil {
				logger.Warnf("[scenes] <%s> failed to fetch performers: %s", sceneHash, err.Error())
			} else {
				var performerJoins []models.PerformersScenes
				for _, performer := range performers {
					join := models.PerformersScenes{
						PerformerID: performer.ID,
						SceneID:     scene.ID,
					}
					performerJoins = append(performerJoins, join)
				}
				if err := jqb.CreatePerformersScenes(performerJoins, tx); err != nil {
					logger.Errorf("[scenes] <%s> failed to associate performers: %s", sceneHash, err.Error())
				}
			}
		}

		// Relate the scene to the movies
		if len(sceneJSON.Movies) > 0 {
			moviesScenes, err := t.getMoviesScenes(sceneJSON.Movies, scene.ID, tx)
			if err != nil {
				logger.Warnf("[scenes] <%s> failed to fetch movies: %s", sceneHash, err.Error())
			} else {
				if err := jqb.CreateMoviesScenes(moviesScenes, tx); err != nil {
					logger.Errorf("[scenes] <%s> failed to associate movies: %s", sceneHash, err.Error())
				}
			}
		}

		// Relate the scene to the tags
		if len(sceneJSON.Tags) > 0 {
			tags, err := t.getTags(sceneHash, sceneJSON.Tags, tx)
			if err != nil {
				logger.Warnf("[scenes] <%s> failed to fetch tags: %s", sceneHash, err.Error())
			} else {
				var tagJoins []models.ScenesTags
				for _, tag := range tags {
					join := models.ScenesTags{
						SceneID: scene.ID,
						TagID:   tag.ID,
					}
					tagJoins = append(tagJoins, join)
				}
				if err := jqb.CreateScenesTags(tagJoins, tx); err != nil {
					logger.Errorf("[scenes] <%s> failed to associate tags: %s", sceneHash, err.Error())
				}
			}
		}

		// Relate the scene to the scene markers
		if len(sceneJSON.Markers) > 0 {
			smqb := models.NewSceneMarkerQueryBuilder()
			tqb := models.NewTagQueryBuilder()
			for _, marker := range sceneJSON.Markers {
				seconds, _ := strconv.ParseFloat(marker.Seconds, 64)
				newSceneMarker := models.SceneMarker{
					Title:     marker.Title,
					Seconds:   seconds,
					SceneID:   sql.NullInt64{Int64: int64(scene.ID), Valid: true},
					CreatedAt: models.SQLiteTimestamp{Timestamp: t.getTimeFromJSONTime(marker.CreatedAt)},
					UpdatedAt: models.SQLiteTimestamp{Timestamp: t.getTimeFromJSONTime(marker.UpdatedAt)},
				}

				primaryTag, err := tqb.FindByName(marker.PrimaryTag, tx, false)
				if err != nil {
					logger.Errorf("[scenes] <%s> failed to find primary tag for marker: %s", sceneHash, err.Error())
				} else {
					newSceneMarker.PrimaryTagID = primaryTag.ID
				}

				// Create the scene marker in the DB
				sceneMarker, err := smqb.Create(newSceneMarker, tx)
				if err != nil {
					logger.Warnf("[scenes] <%s> failed to create scene marker: %s", sceneHash, err.Error())
					continue
				}
				if sceneMarker.ID == 0 {
					logger.Warnf("[scenes] <%s> invalid scene marker id after scene marker creation", sceneHash)
					continue
				}

				// Get the scene marker tags and create the joins
				tags, err := t.getTags(sceneHash, marker.Tags, tx)
				if err != nil {
					logger.Warnf("[scenes] <%s> failed to fetch scene marker tags: %s", sceneHash, err.Error())
				} else {
					var tagJoins []models.SceneMarkersTags
					for _, tag := range tags {
						join := models.SceneMarkersTags{
							SceneMarkerID: sceneMarker.ID,
							TagID:         tag.ID,
						}
						tagJoins = append(tagJoins, join)
					}
					if err := jqb.CreateSceneMarkersTags(tagJoins, tx); err != nil {
						logger.Errorf("[scenes] <%s> failed to associate scene marker tags: %s", sceneHash, err.Error())
					}
				}
			}
		}
	}

	logger.Info("[scenes] importing")
	if err := tx.Commit(); err != nil {
		logger.Errorf("[scenes] import failed to commit: %s", err.Error())
	}
	logger.Info("[scenes] import complete")
}

func (t *ImportTask) getPerformers(names []string, tx *sqlx.Tx) ([]*models.Performer, error) {
	pqb := models.NewPerformerQueryBuilder()
	performers, err := pqb.FindByNames(names, tx, false)
	if err != nil {
		return nil, err
	}

	var pluckedNames []string
	for _, performer := range performers {
		if !performer.Name.Valid {
			continue
		}
		pluckedNames = append(pluckedNames, performer.Name.String)
	}

	missingPerformers := utils.StrFilter(names, func(name string) bool {
		return !utils.StrInclude(pluckedNames, name)
	})

	for _, missingPerformer := range missingPerformers {
		logger.Warnf("[scenes] performer %s does not exist", missingPerformer)
	}

	return performers, nil
}

func (t *ImportTask) getMoviesScenes(input []jsonschema.SceneMovie, sceneID int, tx *sqlx.Tx) ([]models.MoviesScenes, error) {
	mqb := models.NewMovieQueryBuilder()

	var movies []models.MoviesScenes
	for _, inputMovie := range input {
		movie, err := mqb.FindByName(inputMovie.MovieName, tx, false)
		if err != nil {
			return nil, err
		}

		if movie == nil {
			logger.Warnf("[scenes] movie %s does not exist", inputMovie.MovieName)
		} else {
			toAdd := models.MoviesScenes{
				MovieID: movie.ID,
				SceneID: sceneID,
			}

			if inputMovie.SceneIndex != 0 {
				toAdd.SceneIndex = sql.NullInt64{
					Int64: int64(inputMovie.SceneIndex),
					Valid: true,
				}
			}

			movies = append(movies, toAdd)
		}
	}

	return movies, nil
}

func (t *ImportTask) getTags(sceneChecksum string, names []string, tx *sqlx.Tx) ([]*models.Tag, error) {
	tqb := models.NewTagQueryBuilder()
	tags, err := tqb.FindByNames(names, tx, false)
	if err != nil {
		return nil, err
	}

	var pluckedNames []string
	for _, tag := range tags {
		if tag.Name == "" {
			continue
		}
		pluckedNames = append(pluckedNames, tag.Name)
	}

	missingTags := utils.StrFilter(names, func(name string) bool {
		return !utils.StrInclude(pluckedNames, name)
	})

	for _, missingTag := range missingTags {
		logger.Warnf("[scenes] <%s> tag %s does not exist", sceneChecksum, missingTag)
	}

	return tags, nil
}

// https://www.reddit.com/r/golang/comments/5ia523/idiomatic_way_to_remove_duplicates_in_a_slice/db6qa2e
func (t *ImportTask) getUnique(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

var currentLocation = time.Now().Location()

func (t *ImportTask) getTimeFromJSONTime(jsonTime models.JSONTime) time.Time {
	if currentLocation != nil {
		if jsonTime.IsZero() {
			return time.Now().In(currentLocation)
		} else {
			return jsonTime.Time.In(currentLocation)
		}
	} else {
		if jsonTime.IsZero() {
			return time.Now()
		} else {
			return jsonTime.Time
		}
	}
}
