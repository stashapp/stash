package manager

import (
	"archive/zip"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/movie"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/utils"
)

type ImportTask struct {
	json jsonUtils

	BaseDir             string
	ZipFile             io.Reader
	Reset               bool
	DuplicateBehaviour  models.ImportDuplicateEnum
	MissingRefBehaviour models.ImportMissingRefEnum

	mappings            *jsonschema.Mappings
	scraped             []jsonschema.ScrapedItem
	fileNamingAlgorithm models.HashAlgorithm
}

func CreateImportTask(a models.HashAlgorithm, input models.ImportObjectsInput) *ImportTask {
	return &ImportTask{
		ZipFile:             input.File.File,
		Reset:               false,
		DuplicateBehaviour:  input.DuplicateBehaviour,
		MissingRefBehaviour: input.MissingRefBehaviour,
		fileNamingAlgorithm: a,
	}
}

func (t *ImportTask) GetStatus() JobStatus {
	return Import
}

func (t *ImportTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.ZipFile != nil {
		// unzip the file and defer remove the temp directory
		var err error
		t.BaseDir, err = instance.Paths.Generated.TempDir("import")
		if err != nil {
			logger.Errorf("error creating temporary directory for import: %s", err.Error())
			return
		}

		defer func() {
			err := utils.RemoveDir(t.BaseDir)
			if err != nil {
				logger.Errorf("error removing directory %s: %s", t.BaseDir, err.Error())
			}
		}()

		if err := t.unzipFile(); err != nil {
			logger.Errorf("error unzipping provided file for import: %s", err.Error())
			return
		}
	}

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

func (t *ImportTask) unzipFile() error {
	// copy the zip file to the temporary directory
	tmpZip := filepath.Join(t.BaseDir, "import.zip")
	out, err := os.Create(tmpZip)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, t.ZipFile); err != nil {
		out.Close()
		return err
	}

	out.Close()

	defer func() {
		err := os.Remove(tmpZip)
		if err != nil {
			logger.Errorf("error removing temporary zip file %s: %s", tmpZip, err.Error())
		}
	}()

	// now we can read the zip file
	r, err := zip.OpenReader(tmpZip)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fn := filepath.Join(t.BaseDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fn, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fn), os.ModePerm); err != nil {
			return err
		}

		o, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer o.Close()

		i, err := f.Open()
		if err != nil {
			return err
		}
		defer i.Close()

		if _, err := io.Copy(o, i); err != nil {
			return err
		}
	}

	return nil
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
			}

			logger.Errorf("[studios] <%s> failed to create: %s", mappingJSON.Checksum, err.Error())
			continue
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
	logger.Info("[galleries] importing")

	for i, mappingJSON := range t.mappings.Galleries {
		index := i + 1

		logger.Progressf("[galleries] %d of %d", index, len(t.mappings.Galleries))

		tx := database.DB.MustBeginTx(ctx, nil)
		readerWriter := models.NewGalleryReaderWriter(tx)

		galleryImporter := &gallery.Importer{
			ReaderWriter: readerWriter,
			Input:        mappingJSON,
		}

		if err := performImport(galleryImporter, t.DuplicateBehaviour); err != nil {
			tx.Rollback()
			logger.Errorf("[galleries] <%s> failed to import: %s", mappingJSON.Checksum, err.Error())
			continue
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Errorf("[galleries] <%s> import failed to commit: %s", mappingJSON.Checksum, err.Error())
			continue
		}
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
	logger.Info("[scenes] importing")

	for i, mappingJSON := range t.mappings.Scenes {
		index := i + 1

		logger.Progressf("[scenes] %d of %d", index, len(t.mappings.Scenes))

		sceneJSON, err := t.json.getScene(mappingJSON.Checksum)
		if err != nil {
			logger.Infof("[scenes] <%s> json parse failure: %s", mappingJSON.Checksum, err.Error())
			continue
		}

		sceneHash := mappingJSON.Checksum

		tx := database.DB.MustBeginTx(ctx, nil)
		readerWriter := models.NewSceneReaderWriter(tx)
		tagWriter := models.NewTagReaderWriter(tx)
		galleryWriter := models.NewGalleryReaderWriter(tx)
		joinWriter := models.NewJoinReaderWriter(tx)
		movieWriter := models.NewMovieReaderWriter(tx)
		performerWriter := models.NewPerformerReaderWriter(tx)
		studioWriter := models.NewStudioReaderWriter(tx)
		markerWriter := models.NewSceneMarkerReaderWriter(tx)

		sceneImporter := &scene.Importer{
			ReaderWriter: readerWriter,
			Input:        *sceneJSON,
			Path:         mappingJSON.Path,

			FileNamingAlgorithm: t.fileNamingAlgorithm,
			MissingRefBehaviour: t.MissingRefBehaviour,

			GalleryWriter:   galleryWriter,
			JoinWriter:      joinWriter,
			MovieWriter:     movieWriter,
			PerformerWriter: performerWriter,
			StudioWriter:    studioWriter,
			TagWriter:       tagWriter,
		}

		if err := performImport(sceneImporter, t.DuplicateBehaviour); err != nil {
			tx.Rollback()
			logger.Errorf("[scenes] <%s> failed to import: %s", sceneHash, err.Error())
			continue
		}

		// import the scene markers
		failedMarkers := false
		for _, m := range sceneJSON.Markers {
			markerImporter := &scene.MarkerImporter{
				SceneID:             sceneImporter.ID,
				Input:               m,
				MissingRefBehaviour: t.MissingRefBehaviour,
				ReaderWriter:        markerWriter,
				JoinWriter:          joinWriter,
				TagWriter:           tagWriter,
			}

			if err := performImport(markerImporter, t.DuplicateBehaviour); err != nil {
				failedMarkers = true
				logger.Errorf("[scenes] <%s> failed to import markers: %s", sceneHash, err.Error())
				break
			}
		}

		if failedMarkers {
			tx.Rollback()
			continue
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Errorf("[scenes] <%s> import failed to commit: %s", sceneHash, err.Error())
		}
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
