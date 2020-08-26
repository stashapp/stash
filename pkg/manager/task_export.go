package manager

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
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

type ExportTask struct {
	Mappings            *jsonschema.Mappings
	Scraped             []jsonschema.ScrapedItem
	fileNamingAlgorithm models.HashAlgorithm
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
		go exportScene(&scenesWg, jobCh, t)
	}

	for i, scene := range scenes {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[scenes] %d of %d", index, len(scenes))
		}
		t.Mappings.Scenes = append(t.Mappings.Scenes, jsonschema.PathMapping{Path: scene.Path, Checksum: scene.GetHash(t.fileNamingAlgorithm)})
		jobCh <- scene // feed workers
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	scenesWg.Wait()

	logger.Infof("[scenes] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func exportScene(wg *sync.WaitGroup, jobChan <-chan *models.Scene, t *ExportTask) {
	defer wg.Done()
	sceneReader := models.NewSceneReaderWriter(nil)
	studioReader := models.NewStudioReaderWriter(nil)
	movieReader := models.NewMovieReaderWriter(nil)
	galleryReader := models.NewGalleryReaderWriter(nil)
	performerReader := models.NewPerformerReaderWriter(nil)
	tagReader := models.NewTagReaderWriter(nil)
	sceneMarkerReader := models.NewSceneMarkerReaderWriter(nil)
	joinReader := models.NewJoinReaderWriter(nil)

	for s := range jobChan {
		sceneHash := s.GetHash(t.fileNamingAlgorithm)

		newSceneJSON, err := scene.ToBasicJSON(sceneReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene JSON: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Studio, err = scene.GetStudioName(studioReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene studio name: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Gallery, err = scene.GetGalleryChecksum(galleryReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene gallery checksum: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Performers, err = scene.GetPerformerNames(performerReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene performer names: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Tags, err = scene.GetTagNames(tagReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene tag names: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Markers, err = scene.GetSceneMarkersJSON(sceneMarkerReader, tagReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene markers JSON: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Movies, err = scene.GetSceneMoviesJSON(movieReader, joinReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene movies JSON: %s", sceneHash, err.Error())
			continue
		}

		sceneJSON, err := instance.JSON.getScene(sceneHash)
		if err == nil && jsonschema.CompareJSON(*sceneJSON, *newSceneJSON) {
			continue
		}

		if err := instance.JSON.saveScene(sceneHash, newSceneJSON); err != nil {
			logger.Errorf("[scenes] <%s> failed to save json: %s", sceneHash, err.Error())
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

	performerReader := models.NewPerformerReaderWriter(nil)

	for p := range jobChan {
		newPerformerJSON, err := performer.ToJSON(performerReader, p)

		if err != nil {
			logger.Errorf("[performers] <%s> error getting performer JSON: %s", p.Checksum, err.Error())
			continue
		}

		performerJSON, err := instance.JSON.getPerformer(p.Checksum)
		if err != nil {
			logger.Debugf("[performers] error reading performer json: %s", err.Error())
		} else if jsonschema.CompareJSON(*performerJSON, *newPerformerJSON) {
			continue
		}

		if err := instance.JSON.savePerformer(p.Checksum, newPerformerJSON); err != nil {
			logger.Errorf("[performers] <%s> failed to save json: %s", p.Checksum, err.Error())
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

	studioReader := models.NewStudioReaderWriter(nil)

	for s := range jobChan {
		newStudioJSON, err := studio.ToJSON(studioReader, s)

		if err != nil {
			logger.Errorf("[studios] <%s> error getting studio JSON: %s", s.Checksum, err.Error())
			continue
		}

		studioJSON, err := instance.JSON.getStudio(s.Checksum)
		if err == nil && jsonschema.CompareJSON(*studioJSON, *newStudioJSON) {
			continue
		}

		if err := instance.JSON.saveStudio(s.Checksum, newStudioJSON); err != nil {
			logger.Errorf("[studios] <%s> failed to save json: %s", s.Checksum, err.Error())
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

	tagReader := models.NewTagReaderWriter(nil)

	for thisTag := range jobChan {
		newTagJSON, err := tag.ToJSON(tagReader, thisTag)

		if err != nil {
			logger.Errorf("[tags] <%s> error getting tag JSON: %s", thisTag.Name, err.Error())
			continue
		}

		// generate checksum on the fly by name, since we don't store it
		checksum := utils.MD5FromString(thisTag.Name)

		tagJSON, err := instance.JSON.getTag(checksum)
		if err == nil && jsonschema.CompareJSON(*tagJSON, *newTagJSON) {
			continue
		}

		if err := instance.JSON.saveTag(checksum, newTagJSON); err != nil {
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

	movieReader := models.NewMovieReaderWriter(nil)
	studioReader := models.NewStudioReaderWriter(nil)

	for m := range jobChan {
		newMovieJSON, err := movie.ToJSON(movieReader, studioReader, m)

		if err != nil {
			logger.Errorf("[movies] <%s> error getting tag JSON: %s", m.Name, err.Error())
			continue
		}

		movieJSON, err := instance.JSON.getMovie(m.Checksum)
		if err != nil {
			logger.Debugf("[movies] error reading movie json: %s", err.Error())
		} else if jsonschema.CompareJSON(*movieJSON, *newMovieJSON) {
			continue
		}

		if err := instance.JSON.saveMovie(m.Checksum, newMovieJSON); err != nil {
			logger.Errorf("[movies] <%s> failed to save json: %s", m.Checksum, err.Error())
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
