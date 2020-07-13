package manager

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type ExportTask struct {
	Mappings *jsonschema.Mappings
	Scraped  []jsonschema.ScrapedItem
}

func (t *ExportTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	// @manager.total = Scene.count + Gallery.count + Performer.count + Studio.count + Movie.count
	workerCount := runtime.GOMAXPROCS(0) // set worker count to number of cpus available

	t.Mappings = &jsonschema.Mappings{}
	t.Scraped = []jsonschema.ScrapedItem{}

	ctx := context.TODO()
	startTime := time.Now()

	paths.EnsureJSONDirs()

	t.ExportScenes(ctx, workerCount)
	t.ExportGalleries(ctx)
	t.ExportPerformers(ctx, workerCount)
	t.ExportStudios(ctx, workerCount)
	t.ExportMovies(ctx, workerCount)
	t.ExportTags(ctx, workerCount)

	if err := instance.JSON.saveMappings(t.Mappings); err != nil {
		logger.Errorf("[mappings] failed to save json: %s", err.Error())
	}

	t.ExportScrapedItems(ctx)
	logger.Infof("Export complete in %s.", time.Since(startTime))
}

func (t *ExportTask) ExportScenes(ctx context.Context, workers int) {
	var scenesWg sync.WaitGroup

	qb := models.NewSceneQueryBuilder()

	scenes, err := qb.All()
	if err != nil {
		logger.Errorf("[scenes] failed to fetch all scenes: %s", err.Error())
	}

	jobCh := make(chan *models.Scene, workers*2) // make a buffered channel to feed workers

	logger.Info("[scenes] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Scene workers
		scenesWg.Add(1)
		go exportScene(&scenesWg, jobCh, t, nil) // no db data is changed so tx is set to nil
	}

	for i, scene := range scenes {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[scenes] %d of %d", index, len(scenes))
		}
		t.Mappings.Scenes = append(t.Mappings.Scenes, jsonschema.PathMapping{Path: scene.Path, Checksum: scene.Checksum})
		jobCh <- scene // feed workers
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	scenesWg.Wait()

	logger.Infof("[scenes] export complete in %s. %d workers used.", time.Since(startTime), workers)
}
func exportScene(wg *sync.WaitGroup, jobChan <-chan *models.Scene, t *ExportTask, tx *sqlx.Tx) {
	defer wg.Done()
	sceneQB := models.NewSceneQueryBuilder()
	studioQB := models.NewStudioQueryBuilder()
	movieQB := models.NewMovieQueryBuilder()
	galleryQB := models.NewGalleryQueryBuilder()
	performerQB := models.NewPerformerQueryBuilder()
	tagQB := models.NewTagQueryBuilder()
	sceneMarkerQB := models.NewSceneMarkerQueryBuilder()
	joinQB := models.NewJoinsQueryBuilder()

	for scene := range jobChan {
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

		performers, _ := performerQB.FindNameBySceneID(scene.ID, tx)
		sceneMovies, _ := joinQB.GetSceneMovies(scene.ID, tx)
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

		newSceneJSON.OCounter = scene.OCounter

		if scene.Details.Valid {
			newSceneJSON.Details = scene.Details.String
		}
		if galleryChecksum != "" {
			newSceneJSON.Gallery = galleryChecksum
		}

		newSceneJSON.Performers = t.getPerformerNames(performers)
		newSceneJSON.Tags = t.getTagNames(tags)

		for _, sceneMarker := range sceneMarkers {
			primaryTag, err := tagQB.Find(sceneMarker.PrimaryTagID, tx)
			if err != nil {
				logger.Errorf("[scenes] <%s> invalid primary tag for scene marker: %s", scene.Checksum, err.Error())
				continue
			}
			sceneMarkerTags, err := tagQB.FindBySceneMarkerID(sceneMarker.ID, tx)
			if err != nil {
				logger.Errorf("[scenes] <%s> invalid tags for scene marker: %s", scene.Checksum, err.Error())
				continue
			}
			if sceneMarker.Seconds == 0 || primaryTag.Name == "" {
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

		for _, sceneMovie := range sceneMovies {
			movie, _ := movieQB.Find(sceneMovie.MovieID, tx)

			if movie.Name.Valid {
				sceneMovieJSON := jsonschema.SceneMovie{
					MovieName:  movie.Name.String,
					SceneIndex: int(sceneMovie.SceneIndex.Int64),
				}
				newSceneJSON.Movies = append(newSceneJSON.Movies, sceneMovieJSON)
			}
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
		if scene.Format.Valid {
			newSceneJSON.File.Format = scene.Format.String
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

		cover, err := sceneQB.GetSceneCover(scene.ID, tx)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene cover: %s", scene.Checksum, err.Error())
			continue
		}

		if len(cover) > 0 {
			newSceneJSON.Cover = utils.GetBase64StringFromData(cover)
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

func (t *ExportTask) ExportPerformers(ctx context.Context, workers int) {
	var performersWg sync.WaitGroup

	qb := models.NewPerformerQueryBuilder()
	performers, err := qb.All()
	if err != nil {
		logger.Errorf("[performers] failed to fetch all performers: %s", err.Error())
	}
	jobCh := make(chan *models.Performer, workers*2) // make a buffered channel to feed workers

	logger.Info("[performers] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Performer workers
		performersWg.Add(1)
		go exportPerformer(&performersWg, jobCh)
	}

	for i, performer := range performers {
		index := i + 1
		logger.Progressf("[performers] %d of %d", index, len(performers))

		t.Mappings.Performers = append(t.Mappings.Performers, jsonschema.NameMapping{Name: performer.Name.String, Checksum: performer.Checksum})
		jobCh <- performer // feed workers
	}

	close(jobCh) // close channel so workers will know that no more jobs are available
	performersWg.Wait()

	logger.Infof("[performers] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func exportPerformer(wg *sync.WaitGroup, jobChan <-chan *models.Performer) {
	defer wg.Done()

	performerQB := models.NewPerformerQueryBuilder()

	for performer := range jobChan {
		newPerformerJSON := jsonschema.Performer{
			CreatedAt: models.JSONTime{Time: performer.CreatedAt.Timestamp},
			UpdatedAt: models.JSONTime{Time: performer.UpdatedAt.Timestamp},
		}

		if performer.Name.Valid {
			newPerformerJSON.Name = performer.Name.String
		}
		if performer.Gender.Valid {
			newPerformerJSON.Gender = performer.Gender.String
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

		image, err := performerQB.GetPerformerImage(performer.ID, nil)
		if err != nil {
			logger.Errorf("[performers] <%s> error getting performers image: %s", performer.Checksum, err.Error())
			continue
		}

		if len(image) > 0 {
			newPerformerJSON.Image = utils.GetBase64StringFromData(image)
		}

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
}

func (t *ExportTask) ExportStudios(ctx context.Context, workers int) {
	var studiosWg sync.WaitGroup

	qb := models.NewStudioQueryBuilder()
	studios, err := qb.All()
	if err != nil {
		logger.Errorf("[studios] failed to fetch all studios: %s", err.Error())
	}

	logger.Info("[studios] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Studio, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Studio workers
		studiosWg.Add(1)
		go exportStudio(&studiosWg, jobCh)
	}

	for i, studio := range studios {
		index := i + 1
		logger.Progressf("[studios] %d of %d", index, len(studios))

		t.Mappings.Studios = append(t.Mappings.Studios, jsonschema.NameMapping{Name: studio.Name.String, Checksum: studio.Checksum})
		jobCh <- studio // feed workers
	}

	close(jobCh)
	studiosWg.Wait()

	logger.Infof("[studios] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func exportStudio(wg *sync.WaitGroup, jobChan <-chan *models.Studio) {
	defer wg.Done()

	studioQB := models.NewStudioQueryBuilder()

	for studio := range jobChan {

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
		if studio.ParentID.Valid {
			parent, _ := studioQB.Find(int(studio.ParentID.Int64), nil)
			if parent != nil {
				newStudioJSON.ParentStudio = parent.Name.String
			}
		}

		image, err := studioQB.GetStudioImage(studio.ID, nil)
		if err != nil {
			logger.Errorf("[studios] <%s> error getting studio image: %s", studio.Checksum, err.Error())
			continue
		}

		if len(image) > 0 {
			newStudioJSON.Image = utils.GetBase64StringFromData(image)
		}

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
}

func (t *ExportTask) ExportTags(ctx context.Context, workers int) {
	var tagsWg sync.WaitGroup

	qb := models.NewTagQueryBuilder()
	tags, err := qb.All()
	if err != nil {
		logger.Errorf("[tags] failed to fetch all tags: %s", err.Error())
	}

	logger.Info("[tags] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Tag, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Tag workers
		tagsWg.Add(1)
		go exportTag(&tagsWg, jobCh)
	}

	for i, tag := range tags {
		index := i + 1
		logger.Progressf("[tags] %d of %d", index, len(tags))

		// generate checksum on the fly by name, since we don't store it
		checksum := utils.MD5FromString(tag.Name)

		t.Mappings.Tags = append(t.Mappings.Tags, jsonschema.NameMapping{Name: tag.Name, Checksum: checksum})
		jobCh <- tag // feed workers
	}

	close(jobCh)
	tagsWg.Wait()

	logger.Infof("[tags] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func exportTag(wg *sync.WaitGroup, jobChan <-chan *models.Tag) {
	defer wg.Done()

	tagQB := models.NewTagQueryBuilder()

	for tag := range jobChan {

		newTagJSON := jsonschema.Tag{
			Name:      tag.Name,
			CreatedAt: models.JSONTime{Time: tag.CreatedAt.Timestamp},
			UpdatedAt: models.JSONTime{Time: tag.UpdatedAt.Timestamp},
		}

		image, err := tagQB.GetTagImage(tag.ID, nil)
		if err != nil {
			logger.Errorf("[tags] <%s> error getting tag image: %s", tag.Name, err.Error())
			continue
		}

		if len(image) > 0 {
			newTagJSON.Image = utils.GetBase64StringFromData(image)
		}

		// generate checksum on the fly by name, since we don't store it
		checksum := utils.MD5FromString(tag.Name)

		tagJSON, err := instance.JSON.getTag(checksum)
		if err != nil {
			logger.Debugf("[tags] error reading tag json: %s", err.Error())
		} else if jsonschema.CompareJSON(*tagJSON, newTagJSON) {
			continue
		}

		if err := instance.JSON.saveTag(checksum, &newTagJSON); err != nil {
			logger.Errorf("[tags] <%s> failed to save json: %s", checksum, err.Error())
		}
	}
}

func (t *ExportTask) ExportMovies(ctx context.Context, workers int) {
	var moviesWg sync.WaitGroup

	qb := models.NewMovieQueryBuilder()
	movies, err := qb.All()
	if err != nil {
		logger.Errorf("[movies] failed to fetch all movies: %s", err.Error())
	}

	logger.Info("[movies] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Movie, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Studio workers
		moviesWg.Add(1)
		go exportMovie(&moviesWg, jobCh)
	}

	for i, movie := range movies {
		index := i + 1
		logger.Progressf("[movies] %d of %d", index, len(movies))

		t.Mappings.Movies = append(t.Mappings.Movies, jsonschema.NameMapping{Name: movie.Name.String, Checksum: movie.Checksum})
		jobCh <- movie // feed workers
	}

	close(jobCh)
	moviesWg.Wait()

	logger.Infof("[movies] export complete in %s. %d workers used.", time.Since(startTime), workers)

}
func exportMovie(wg *sync.WaitGroup, jobChan <-chan *models.Movie) {
	defer wg.Done()

	movieQB := models.NewMovieQueryBuilder()
	studioQB := models.NewStudioQueryBuilder()

	for movie := range jobChan {
		newMovieJSON := jsonschema.Movie{
			CreatedAt: models.JSONTime{Time: movie.CreatedAt.Timestamp},
			UpdatedAt: models.JSONTime{Time: movie.UpdatedAt.Timestamp},
		}

		if movie.Name.Valid {
			newMovieJSON.Name = movie.Name.String
		}
		if movie.Aliases.Valid {
			newMovieJSON.Aliases = movie.Aliases.String
		}
		if movie.Date.Valid {
			newMovieJSON.Date = utils.GetYMDFromDatabaseDate(movie.Date.String)
		}
		if movie.Rating.Valid {
			newMovieJSON.Rating = int(movie.Rating.Int64)
		}
		if movie.Duration.Valid {
			newMovieJSON.Duration = int(movie.Duration.Int64)
		}

		if movie.Director.Valid {
			newMovieJSON.Director = movie.Director.String
		}

		if movie.Synopsis.Valid {
			newMovieJSON.Synopsis = movie.Synopsis.String
		}

		if movie.URL.Valid {
			newMovieJSON.URL = movie.URL.String
		}

		if movie.StudioID.Valid {
			studio, _ := studioQB.Find(int(movie.StudioID.Int64), nil)
			if studio != nil {
				newMovieJSON.Studio = studio.Name.String
			}
		}

		frontImage, err := movieQB.GetFrontImage(movie.ID, nil)
		if err != nil {
			logger.Errorf("[movies] <%s> error getting movie front image: %s", movie.Checksum, err.Error())
			continue
		}

		if len(frontImage) > 0 {
			newMovieJSON.FrontImage = utils.GetBase64StringFromData(frontImage)
		}

		backImage, err := movieQB.GetBackImage(movie.ID, nil)
		if err != nil {
			logger.Errorf("[movies] <%s> error getting movie back image: %s", movie.Checksum, err.Error())
			continue
		}

		if len(backImage) > 0 {
			newMovieJSON.BackImage = utils.GetBase64StringFromData(backImage)
		}

		movieJSON, err := instance.JSON.getMovie(movie.Checksum)
		if err != nil {
			logger.Debugf("[movies] error reading movie json: %s", err.Error())
		} else if jsonschema.CompareJSON(*movieJSON, newMovieJSON) {
			continue
		}

		if err := instance.JSON.saveMovie(movie.Checksum, &newMovieJSON); err != nil {
			logger.Errorf("[movies] <%s> failed to save json: %s", movie.Checksum, err.Error())
		}
	}
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

func (t *ExportTask) getPerformerNames(performers []*models.Performer) []string {
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

func (t *ExportTask) getTagNames(tags []*models.Tag) []string {
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
