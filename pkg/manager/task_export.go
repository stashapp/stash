package manager

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

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

type ExportTask struct {
	full bool

	baseDir string
	json    jsonUtils

	Mappings            *jsonschema.Mappings
	fileNamingAlgorithm models.HashAlgorithm

	scenes     *exportSpec
	performers *exportSpec
	movies     *exportSpec
	tags       *exportSpec
	studios    *exportSpec
	galleries  *exportSpec

	includeDependencies bool

	DownloadHash string
}

type exportSpec struct {
	IDs []int
	all bool
}

func newExportSpec(input *models.ExportObjectTypeInput) *exportSpec {
	if input == nil {
		return &exportSpec{}
	}

	ret := &exportSpec{
		IDs: utils.StringSliceToIntSlice(input.Ids),
	}

	if input.All != nil {
		ret.all = *input.All
	}

	return ret
}

func CreateExportTask(a models.HashAlgorithm, input models.ExportObjectsInput) *ExportTask {
	includeDeps := false
	if input.IncludeDependencies != nil {
		includeDeps = *input.IncludeDependencies
	}

	return &ExportTask{
		fileNamingAlgorithm: a,
		scenes:              newExportSpec(input.Scenes),
		performers:          newExportSpec(input.Performers),
		movies:              newExportSpec(input.Movies),
		tags:                newExportSpec(input.Tags),
		studios:             newExportSpec(input.Studios),
		galleries:           newExportSpec(input.Galleries),
		includeDependencies: includeDeps,
	}
}

func (t *ExportTask) GetStatus() JobStatus {
	return Export
}

func (t *ExportTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	// @manager.total = Scene.count + Gallery.count + Performer.count + Studio.count + Movie.count
	workerCount := runtime.GOMAXPROCS(0) // set worker count to number of cpus available

	t.Mappings = &jsonschema.Mappings{}

	startTime := time.Now()

	if t.full {
		t.baseDir = config.GetMetadataPath()
	} else {
		var err error
		t.baseDir, err = instance.Paths.Generated.TempDir("export")
		if err != nil {
			logger.Errorf("error creating temporary directory for export: %s", err.Error())
			return
		}

		defer func() {
			err := utils.RemoveDir(t.baseDir)
			if err != nil {
				logger.Errorf("error removing directory %s: %s", t.baseDir, err.Error())
			}
		}()
	}

	t.json = jsonUtils{
		json: *paths.GetJSONPaths(t.baseDir),
	}

	paths.EnsureJSONDirs(t.baseDir)

	t.ExportScenes(workerCount)
	t.ExportGalleries()
	t.ExportPerformers(workerCount)
	t.ExportStudios(workerCount)
	t.ExportMovies(workerCount)
	t.ExportTags(workerCount)

	if err := t.json.saveMappings(t.Mappings); err != nil {
		logger.Errorf("[mappings] failed to save json: %s", err.Error())
	}

	if t.full {
		t.ExportScrapedItems()
	} else {
		err := t.generateDownload()
		if err != nil {
			logger.Errorf("error generating download link: %s", err.Error())
			return
		}
	}
	logger.Infof("Export complete in %s.", time.Since(startTime))
}

func (t *ExportTask) generateDownload() error {
	// zip the files and register a download link
	utils.EnsureDir(instance.Paths.Generated.Downloads)
	z, err := ioutil.TempFile(instance.Paths.Generated.Downloads, "export*.zip")
	if err != nil {
		return err
	}
	defer z.Close()

	err = t.zipFiles(z)
	if err != nil {
		return err
	}

	t.DownloadHash = instance.DownloadStore.RegisterFile(z.Name(), "", false)
	logger.Debugf("Generated zip file %s with hash %s", z.Name(), t.DownloadHash)
	return nil
}

func (t *ExportTask) zipFiles(w io.Writer) error {
	z := zip.NewWriter(w)
	defer z.Close()

	u := jsonUtils{
		json: *paths.GetJSONPaths(""),
	}

	// write the mappings file
	err := t.zipFile(t.json.json.MappingsFile, "", z)
	if err != nil {
		return err
	}

	filepath.Walk(t.json.json.Tags, t.zipWalkFunc(u.json.Tags, z))
	filepath.Walk(t.json.json.Galleries, t.zipWalkFunc(u.json.Galleries, z))
	filepath.Walk(t.json.json.Performers, t.zipWalkFunc(u.json.Performers, z))
	filepath.Walk(t.json.json.Studios, t.zipWalkFunc(u.json.Studios, z))
	filepath.Walk(t.json.json.Movies, t.zipWalkFunc(u.json.Movies, z))
	filepath.Walk(t.json.json.Scenes, t.zipWalkFunc(u.json.Scenes, z))

	return nil
}

func (t *ExportTask) zipWalkFunc(outDir string, z *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return t.zipFile(path, outDir, z)
	}
}

func (t *ExportTask) zipFile(fn, outDir string, z *zip.Writer) error {
	bn := filepath.Base(fn)

	f, err := z.Create(filepath.Join(outDir, bn))
	if err != nil {
		return fmt.Errorf("error creating zip entry for %s: %s", fn, err.Error())
	}

	i, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("error opening %s: %s", fn, err.Error())
	}

	defer i.Close()

	if _, err := io.Copy(f, i); err != nil {
		return fmt.Errorf("error writing %s to zip: %s", fn, err.Error())
	}

	return nil
}

func (t *ExportTask) ExportScenes(workers int) {
	var scenesWg sync.WaitGroup

	sceneReader := models.NewSceneReaderWriter(nil)

	var scenes []*models.Scene
	var err error
	all := t.full || (t.scenes != nil && t.scenes.all)
	if all {
		scenes, err = sceneReader.All()
	} else if t.scenes != nil && len(t.scenes.IDs) > 0 {
		scenes, err = sceneReader.FindMany(t.scenes.IDs)
	}

	if err != nil {
		logger.Errorf("[scenes] failed to fetch scenes: %s", err.Error())
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

		sceneGallery, err := galleryReader.FindBySceneID(s.ID)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene gallery: %s", sceneHash, err.Error())
			continue
		}

		if sceneGallery != nil {
			newSceneJSON.Gallery = sceneGallery.Checksum
		}

		performers, err := performerReader.FindBySceneID(s.ID)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene performer names: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Performers = performer.GetNames(performers)

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

		if t.includeDependencies {
			if s.StudioID.Valid {
				t.studios.IDs = utils.IntAppendUnique(t.studios.IDs, int(s.StudioID.Int64))
			}

			if sceneGallery != nil {
				t.galleries.IDs = utils.IntAppendUnique(t.galleries.IDs, sceneGallery.ID)
			}

			tagIDs, err := scene.GetDependentTagIDs(tagReader, joinReader, sceneMarkerReader, s)
			if err != nil {
				logger.Errorf("[scenes] <%s> error getting scene tags: %s", sceneHash, err.Error())
				continue
			}
			t.tags.IDs = utils.IntAppendUniques(t.tags.IDs, tagIDs)

			movieIDs, err := scene.GetDependentMovieIDs(joinReader, s)
			if err != nil {
				logger.Errorf("[scenes] <%s> error getting scene movies: %s", sceneHash, err.Error())
				continue
			}
			t.movies.IDs = utils.IntAppendUniques(t.movies.IDs, movieIDs)

			t.performers.IDs = utils.IntAppendUniques(t.performers.IDs, performer.GetIDs(performers))
		}

		sceneJSON, err := t.json.getScene(sceneHash)
		if err == nil && jsonschema.CompareJSON(*sceneJSON, *newSceneJSON) {
			continue
		}

		if err := t.json.saveScene(sceneHash, newSceneJSON); err != nil {
			logger.Errorf("[scenes] <%s> failed to save json: %s", sceneHash, err.Error())
		}
	}
}

func (t *ExportTask) ExportGalleries() {
	reader := models.NewGalleryReaderWriter(nil)

	var galleries []*models.Gallery
	var err error
	all := t.full || (t.galleries != nil && t.galleries.all)
	if all {
		galleries, err = reader.All()
	} else if t.galleries != nil && len(t.galleries.IDs) > 0 {
		galleries, err = reader.FindMany(t.galleries.IDs)
	}

	if err != nil {
		logger.Errorf("[galleries] failed to fetch galleries: %s", err.Error())
	}

	logger.Info("[galleries] exporting")

	for i, gallery := range galleries {
		index := i + 1
		logger.Progressf("[galleries] %d of %d", index, len(galleries))
		t.Mappings.Galleries = append(t.Mappings.Galleries, jsonschema.PathMapping{Path: gallery.Path, Checksum: gallery.Checksum})
	}

	logger.Infof("[galleries] export complete")
}

func (t *ExportTask) ExportPerformers(workers int) {
	var performersWg sync.WaitGroup

	reader := models.NewPerformerReaderWriter(nil)
	var performers []*models.Performer
	var err error
	all := t.full || (t.performers != nil && t.performers.all)
	if all {
		performers, err = reader.All()
	} else if t.performers != nil && len(t.performers.IDs) > 0 {
		performers, err = reader.FindMany(t.performers.IDs)
	}

	if err != nil {
		logger.Errorf("[performers] failed to fetch performers: %s", err.Error())
	}
	jobCh := make(chan *models.Performer, workers*2) // make a buffered channel to feed workers

	logger.Info("[performers] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Performer workers
		performersWg.Add(1)
		go t.exportPerformer(&performersWg, jobCh)
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

func (t *ExportTask) exportPerformer(wg *sync.WaitGroup, jobChan <-chan *models.Performer) {
	defer wg.Done()

	performerReader := models.NewPerformerReaderWriter(nil)

	for p := range jobChan {
		newPerformerJSON, err := performer.ToJSON(performerReader, p)

		if err != nil {
			logger.Errorf("[performers] <%s> error getting performer JSON: %s", p.Checksum, err.Error())
			continue
		}

		performerJSON, err := t.json.getPerformer(p.Checksum)
		if err != nil {
			logger.Debugf("[performers] error reading performer json: %s", err.Error())
		} else if jsonschema.CompareJSON(*performerJSON, *newPerformerJSON) {
			continue
		}

		if err := t.json.savePerformer(p.Checksum, newPerformerJSON); err != nil {
			logger.Errorf("[performers] <%s> failed to save json: %s", p.Checksum, err.Error())
		}
	}
}

func (t *ExportTask) ExportStudios(workers int) {
	var studiosWg sync.WaitGroup

	reader := models.NewStudioReaderWriter(nil)
	var studios []*models.Studio
	var err error
	all := t.full || (t.studios != nil && t.studios.all)
	if all {
		studios, err = reader.All()
	} else if t.studios != nil && len(t.studios.IDs) > 0 {
		studios, err = reader.FindMany(t.studios.IDs)
	}

	if err != nil {
		logger.Errorf("[studios] failed to fetch studios: %s", err.Error())
	}

	logger.Info("[studios] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Studio, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Studio workers
		studiosWg.Add(1)
		go t.exportStudio(&studiosWg, jobCh)
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

func (t *ExportTask) exportStudio(wg *sync.WaitGroup, jobChan <-chan *models.Studio) {
	defer wg.Done()

	studioReader := models.NewStudioReaderWriter(nil)

	for s := range jobChan {
		newStudioJSON, err := studio.ToJSON(studioReader, s)

		if err != nil {
			logger.Errorf("[studios] <%s> error getting studio JSON: %s", s.Checksum, err.Error())
			continue
		}

		studioJSON, err := t.json.getStudio(s.Checksum)
		if err == nil && jsonschema.CompareJSON(*studioJSON, *newStudioJSON) {
			continue
		}

		if err := t.json.saveStudio(s.Checksum, newStudioJSON); err != nil {
			logger.Errorf("[studios] <%s> failed to save json: %s", s.Checksum, err.Error())
		}
	}
}

func (t *ExportTask) ExportTags(workers int) {
	var tagsWg sync.WaitGroup

	reader := models.NewTagReaderWriter(nil)
	var tags []*models.Tag
	var err error
	all := t.full || (t.tags != nil && t.tags.all)
	if all {
		tags, err = reader.All()
	} else if t.tags != nil && len(t.tags.IDs) > 0 {
		tags, err = reader.FindMany(t.tags.IDs)
	}

	if err != nil {
		logger.Errorf("[tags] failed to fetch tags: %s", err.Error())
	}

	logger.Info("[tags] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Tag, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Tag workers
		tagsWg.Add(1)
		go t.exportTag(&tagsWg, jobCh)
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

func (t *ExportTask) exportTag(wg *sync.WaitGroup, jobChan <-chan *models.Tag) {
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

		tagJSON, err := t.json.getTag(checksum)
		if err == nil && jsonschema.CompareJSON(*tagJSON, *newTagJSON) {
			continue
		}

		if err := t.json.saveTag(checksum, newTagJSON); err != nil {
			logger.Errorf("[tags] <%s> failed to save json: %s", checksum, err.Error())
		}
	}
}

func (t *ExportTask) ExportMovies(workers int) {
	var moviesWg sync.WaitGroup

	reader := models.NewMovieReaderWriter(nil)
	var movies []*models.Movie
	var err error
	all := t.full || (t.movies != nil && t.movies.all)
	if all {
		movies, err = reader.All()
	} else if t.movies != nil && len(t.movies.IDs) > 0 {
		movies, err = reader.FindMany(t.movies.IDs)
	}

	if err != nil {
		logger.Errorf("[movies] failed to fetch movies: %s", err.Error())
	}

	logger.Info("[movies] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Movie, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Studio workers
		moviesWg.Add(1)
		go t.exportMovie(&moviesWg, jobCh)
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
func (t *ExportTask) exportMovie(wg *sync.WaitGroup, jobChan <-chan *models.Movie) {
	defer wg.Done()

	movieReader := models.NewMovieReaderWriter(nil)
	studioReader := models.NewStudioReaderWriter(nil)

	for m := range jobChan {
		newMovieJSON, err := movie.ToJSON(movieReader, studioReader, m)

		if err != nil {
			logger.Errorf("[movies] <%s> error getting tag JSON: %s", m.Checksum, err.Error())
			continue
		}

		movieJSON, err := t.json.getMovie(m.Checksum)
		if err != nil {
			logger.Debugf("[movies] error reading movie json: %s", err.Error())
		} else if jsonschema.CompareJSON(*movieJSON, *newMovieJSON) {
			continue
		}

		if err := t.json.saveMovie(m.Checksum, newMovieJSON); err != nil {
			logger.Errorf("[movies] <%s> failed to save json: %s", m.Checksum, err.Error())
		}
	}
}

func (t *ExportTask) ExportScrapedItems() {
	qb := models.NewScrapedItemQueryBuilder()
	sqb := models.NewStudioQueryBuilder()
	scrapedItems, err := qb.All()
	if err != nil {
		logger.Errorf("[scraped sites] failed to fetch all items: %s", err.Error())
	}

	logger.Info("[scraped sites] exporting")

	scraped := []jsonschema.ScrapedItem{}

	for i, scrapedItem := range scrapedItems {
		index := i + 1
		logger.Progressf("[scraped sites] %d of %d", index, len(scrapedItems))

		var studioName string
		if scrapedItem.StudioID.Valid {
			studio, _ := sqb.Find(int(scrapedItem.StudioID.Int64), nil)
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

		scraped = append(scraped, newScrapedItemJSON)
	}

	scrapedJSON, err := t.json.getScraped()
	if err != nil {
		logger.Debugf("[scraped sites] error reading json: %s", err.Error())
	}
	if !jsonschema.CompareJSON(scrapedJSON, scraped) {
		if err := t.json.saveScaped(scraped); err != nil {
			logger.Errorf("[scraped sites] failed to save json: %s", err.Error())
		}
	}

	logger.Infof("[scraped sites] export complete")
}
