package manager

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/movie"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/utils"
)

type ExportTask struct {
	txnManager models.TransactionManager
	full       bool

	baseDir string
	json    jsonUtils

	Mappings            *jsonschema.Mappings
	fileNamingAlgorithm models.HashAlgorithm

	scenes     *exportSpec
	images     *exportSpec
	performers *exportSpec
	movies     *exportSpec
	tags       *exportSpec
	studios    *exportSpec
	galleries  *exportSpec

	includeDependencies bool

	DownloadHash string
}

type ExportObjectTypeInput struct {
	Ids []string `json:"ids"`
	All *bool    `json:"all"`
}

type ExportObjectsInput struct {
	Scenes              *ExportObjectTypeInput `json:"scenes"`
	Images              *ExportObjectTypeInput `json:"images"`
	Studios             *ExportObjectTypeInput `json:"studios"`
	Performers          *ExportObjectTypeInput `json:"performers"`
	Tags                *ExportObjectTypeInput `json:"tags"`
	Movies              *ExportObjectTypeInput `json:"movies"`
	Galleries           *ExportObjectTypeInput `json:"galleries"`
	IncludeDependencies *bool                  `json:"includeDependencies"`
}

type exportSpec struct {
	IDs []int
	all bool
}

func newExportSpec(input *ExportObjectTypeInput) *exportSpec {
	if input == nil {
		return &exportSpec{}
	}

	ids, _ := stringslice.StringSliceToIntSlice(input.Ids)

	ret := &exportSpec{
		IDs: ids,
	}

	if input.All != nil {
		ret.all = *input.All
	}

	return ret
}

func CreateExportTask(a models.HashAlgorithm, input ExportObjectsInput) *ExportTask {
	includeDeps := false
	if input.IncludeDependencies != nil {
		includeDeps = *input.IncludeDependencies
	}

	return &ExportTask{
		txnManager:          GetInstance().TxnManager,
		fileNamingAlgorithm: a,
		scenes:              newExportSpec(input.Scenes),
		images:              newExportSpec(input.Images),
		performers:          newExportSpec(input.Performers),
		movies:              newExportSpec(input.Movies),
		tags:                newExportSpec(input.Tags),
		studios:             newExportSpec(input.Studios),
		galleries:           newExportSpec(input.Galleries),
		includeDependencies: includeDeps,
	}
}

func (t *ExportTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// @manager.total = Scene.count + Gallery.count + Performer.count + Studio.count + Movie.count
	workerCount := runtime.GOMAXPROCS(0) // set worker count to number of cpus available

	t.Mappings = &jsonschema.Mappings{}

	startTime := time.Now()

	if t.full {
		t.baseDir = config.GetInstance().GetMetadataPath()
	} else {
		var err error
		t.baseDir, err = instance.Paths.Generated.TempDir("export")
		if err != nil {
			logger.Errorf("error creating temporary directory for export: %s", err.Error())
			return
		}

		defer func() {
			err := fsutil.RemoveDir(t.baseDir)
			if err != nil {
				logger.Errorf("error removing directory %s: %s", t.baseDir, err.Error())
			}
		}()
	}

	t.json = jsonUtils{
		json: *paths.GetJSONPaths(t.baseDir),
	}

	paths.EnsureJSONDirs(t.baseDir)

	txnErr := t.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		// include movie scenes and gallery images
		if !t.full {
			// only include movie scenes if includeDependencies is also set
			if !t.scenes.all && t.includeDependencies {
				t.populateMovieScenes(r)
			}

			// always export gallery images
			if !t.images.all {
				t.populateGalleryImages(r)
			}
		}

		t.ExportScenes(workerCount, r)
		t.ExportImages(workerCount, r)
		t.ExportGalleries(workerCount, r)
		t.ExportMovies(workerCount, r)
		t.ExportPerformers(workerCount, r)
		t.ExportStudios(workerCount, r)
		t.ExportTags(workerCount, r)

		if t.full {
			t.ExportScrapedItems(r)
		}

		return nil
	})
	if txnErr != nil {
		logger.Warnf("error while running export transaction: %v", txnErr)
	}

	if err := t.json.saveMappings(t.Mappings); err != nil {
		logger.Errorf("[mappings] failed to save json: %s", err.Error())
	}

	if !t.full {
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
	if err := fsutil.EnsureDir(instance.Paths.Generated.Downloads); err != nil {
		return err
	}
	z, err := os.CreateTemp(instance.Paths.Generated.Downloads, "export*.zip")
	if err != nil {
		return err
	}
	defer z.Close()

	err = t.zipFiles(z)
	if err != nil {
		return err
	}

	t.DownloadHash, err = instance.DownloadStore.RegisterFile(z.Name(), "", false)
	if err != nil {
		return fmt.Errorf("error registering file for download: %w", err)
	}
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

	walkWarn(t.json.json.Tags, t.zipWalkFunc(u.json.Tags, z))
	walkWarn(t.json.json.Galleries, t.zipWalkFunc(u.json.Galleries, z))
	walkWarn(t.json.json.Performers, t.zipWalkFunc(u.json.Performers, z))
	walkWarn(t.json.json.Studios, t.zipWalkFunc(u.json.Studios, z))
	walkWarn(t.json.json.Movies, t.zipWalkFunc(u.json.Movies, z))
	walkWarn(t.json.json.Scenes, t.zipWalkFunc(u.json.Scenes, z))
	walkWarn(t.json.json.Images, t.zipWalkFunc(u.json.Images, z))

	return nil
}

// like filepath.Walk but issue a warning on error
func walkWarn(root string, fn filepath.WalkFunc) {
	if err := filepath.Walk(root, fn); err != nil {
		logger.Warnf("error walking structure %v: %v", root, err)
	}
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

func (t *ExportTask) populateMovieScenes(repo models.ReaderRepository) {
	reader := repo.Movie()
	sceneReader := repo.Scene()

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

	for _, m := range movies {
		scenes, err := sceneReader.FindByMovieID(m.ID)
		if err != nil {
			logger.Errorf("[movies] <%s> failed to fetch scenes for movie: %s", m.Checksum, err.Error())
			continue
		}

		for _, s := range scenes {
			t.scenes.IDs = intslice.IntAppendUnique(t.scenes.IDs, s.ID)
		}
	}
}

func (t *ExportTask) populateGalleryImages(repo models.ReaderRepository) {
	reader := repo.Gallery()
	imageReader := repo.Image()

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

	for _, g := range galleries {
		images, err := imageReader.FindByGalleryID(g.ID)
		if err != nil {
			logger.Errorf("[galleries] <%s> failed to fetch images for gallery: %s", g.Checksum, err.Error())
			continue
		}

		for _, i := range images {
			t.images.IDs = intslice.IntAppendUnique(t.images.IDs, i.ID)
		}
	}
}

func (t *ExportTask) ExportScenes(workers int, repo models.ReaderRepository) {
	var scenesWg sync.WaitGroup

	sceneReader := repo.Scene()

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
		go exportScene(&scenesWg, jobCh, repo, t)
	}

	for i, scene := range scenes {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[scenes] %d of %d", index, len(scenes))
		}
		t.Mappings.Scenes = append(t.Mappings.Scenes, jsonschema.PathNameMapping{Path: scene.Path, Checksum: scene.GetHash(t.fileNamingAlgorithm)})
		jobCh <- scene // feed workers
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	scenesWg.Wait()

	logger.Infof("[scenes] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func exportScene(wg *sync.WaitGroup, jobChan <-chan *models.Scene, repo models.ReaderRepository, t *ExportTask) {
	defer wg.Done()
	sceneReader := repo.Scene()
	studioReader := repo.Studio()
	movieReader := repo.Movie()
	galleryReader := repo.Gallery()
	performerReader := repo.Performer()
	tagReader := repo.Tag()
	sceneMarkerReader := repo.SceneMarker()

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

		galleries, err := galleryReader.FindBySceneID(s.ID)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene gallery checksums: %s", sceneHash, err.Error())
			continue
		}

		newSceneJSON.Galleries = gallery.GetChecksums(galleries)

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

		newSceneJSON.Movies, err = scene.GetSceneMoviesJSON(movieReader, sceneReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene movies JSON: %s", sceneHash, err.Error())
			continue
		}

		if t.includeDependencies {
			if s.StudioID.Valid {
				t.studios.IDs = intslice.IntAppendUnique(t.studios.IDs, int(s.StudioID.Int64))
			}

			t.galleries.IDs = intslice.IntAppendUniques(t.galleries.IDs, gallery.GetIDs(galleries))

			tagIDs, err := scene.GetDependentTagIDs(tagReader, sceneMarkerReader, s)
			if err != nil {
				logger.Errorf("[scenes] <%s> error getting scene tags: %s", sceneHash, err.Error())
				continue
			}
			t.tags.IDs = intslice.IntAppendUniques(t.tags.IDs, tagIDs)

			movieIDs, err := scene.GetDependentMovieIDs(sceneReader, s)
			if err != nil {
				logger.Errorf("[scenes] <%s> error getting scene movies: %s", sceneHash, err.Error())
				continue
			}
			t.movies.IDs = intslice.IntAppendUniques(t.movies.IDs, movieIDs)

			t.performers.IDs = intslice.IntAppendUniques(t.performers.IDs, performer.GetIDs(performers))
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

func (t *ExportTask) ExportImages(workers int, repo models.ReaderRepository) {
	var imagesWg sync.WaitGroup

	imageReader := repo.Image()

	var images []*models.Image
	var err error
	all := t.full || (t.images != nil && t.images.all)
	if all {
		images, err = imageReader.All()
	} else if t.images != nil && len(t.images.IDs) > 0 {
		images, err = imageReader.FindMany(t.images.IDs)
	}

	if err != nil {
		logger.Errorf("[images] failed to fetch images: %s", err.Error())
	}

	jobCh := make(chan *models.Image, workers*2) // make a buffered channel to feed workers

	logger.Info("[images] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Image workers
		imagesWg.Add(1)
		go exportImage(&imagesWg, jobCh, repo, t)
	}

	for i, image := range images {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[images] %d of %d", index, len(images))
		}
		t.Mappings.Images = append(t.Mappings.Images, jsonschema.PathNameMapping{Path: image.Path, Checksum: image.Checksum})
		jobCh <- image // feed workers
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	imagesWg.Wait()

	logger.Infof("[images] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func exportImage(wg *sync.WaitGroup, jobChan <-chan *models.Image, repo models.ReaderRepository, t *ExportTask) {
	defer wg.Done()
	studioReader := repo.Studio()
	galleryReader := repo.Gallery()
	performerReader := repo.Performer()
	tagReader := repo.Tag()

	for s := range jobChan {
		imageHash := s.Checksum

		newImageJSON := image.ToBasicJSON(s)

		var err error
		newImageJSON.Studio, err = image.GetStudioName(studioReader, s)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image studio name: %s", imageHash, err.Error())
			continue
		}

		imageGalleries, err := galleryReader.FindByImageID(s.ID)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image galleries: %s", imageHash, err.Error())
			continue
		}

		newImageJSON.Galleries = t.getGalleryChecksums(imageGalleries)

		performers, err := performerReader.FindByImageID(s.ID)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image performer names: %s", imageHash, err.Error())
			continue
		}

		newImageJSON.Performers = performer.GetNames(performers)

		tags, err := tagReader.FindByImageID(s.ID)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image tag names: %s", imageHash, err.Error())
			continue
		}

		newImageJSON.Tags = tag.GetNames(tags)

		if t.includeDependencies {
			if s.StudioID.Valid {
				t.studios.IDs = intslice.IntAppendUnique(t.studios.IDs, int(s.StudioID.Int64))
			}

			t.galleries.IDs = intslice.IntAppendUniques(t.galleries.IDs, gallery.GetIDs(imageGalleries))
			t.tags.IDs = intslice.IntAppendUniques(t.tags.IDs, tag.GetIDs(tags))
			t.performers.IDs = intslice.IntAppendUniques(t.performers.IDs, performer.GetIDs(performers))
		}

		imageJSON, err := t.json.getImage(imageHash)
		if err == nil && jsonschema.CompareJSON(*imageJSON, *newImageJSON) {
			continue
		}

		if err := t.json.saveImage(imageHash, newImageJSON); err != nil {
			logger.Errorf("[images] <%s> failed to save json: %s", imageHash, err.Error())
		}
	}
}

func (t *ExportTask) getGalleryChecksums(galleries []*models.Gallery) (ret []string) {
	for _, g := range galleries {
		ret = append(ret, g.Checksum)
	}
	return
}

func (t *ExportTask) ExportGalleries(workers int, repo models.ReaderRepository) {
	var galleriesWg sync.WaitGroup

	reader := repo.Gallery()

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

	jobCh := make(chan *models.Gallery, workers*2) // make a buffered channel to feed workers

	logger.Info("[galleries] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Scene workers
		galleriesWg.Add(1)
		go exportGallery(&galleriesWg, jobCh, repo, t)
	}

	for i, gallery := range galleries {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[galleries] %d of %d", index, len(galleries))
		}

		t.Mappings.Galleries = append(t.Mappings.Galleries, jsonschema.PathNameMapping{
			Path:     gallery.Path.String,
			Name:     gallery.Title.String,
			Checksum: gallery.Checksum,
		})
		jobCh <- gallery
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	galleriesWg.Wait()

	logger.Infof("[galleries] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func exportGallery(wg *sync.WaitGroup, jobChan <-chan *models.Gallery, repo models.ReaderRepository, t *ExportTask) {
	defer wg.Done()
	studioReader := repo.Studio()
	performerReader := repo.Performer()
	tagReader := repo.Tag()

	for g := range jobChan {
		galleryHash := g.Checksum

		newGalleryJSON, err := gallery.ToBasicJSON(g)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery JSON: %s", galleryHash, err.Error())
			continue
		}

		newGalleryJSON.Studio, err = gallery.GetStudioName(studioReader, g)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery studio name: %s", galleryHash, err.Error())
			continue
		}

		performers, err := performerReader.FindByGalleryID(g.ID)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery performer names: %s", galleryHash, err.Error())
			continue
		}

		newGalleryJSON.Performers = performer.GetNames(performers)

		tags, err := tagReader.FindByGalleryID(g.ID)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery tag names: %s", galleryHash, err.Error())
			continue
		}

		newGalleryJSON.Tags = tag.GetNames(tags)

		if t.includeDependencies {
			if g.StudioID.Valid {
				t.studios.IDs = intslice.IntAppendUnique(t.studios.IDs, int(g.StudioID.Int64))
			}

			t.tags.IDs = intslice.IntAppendUniques(t.tags.IDs, tag.GetIDs(tags))
			t.performers.IDs = intslice.IntAppendUniques(t.performers.IDs, performer.GetIDs(performers))
		}

		galleryJSON, err := t.json.getGallery(galleryHash)
		if err == nil && jsonschema.CompareJSON(*galleryJSON, *newGalleryJSON) {
			continue
		}

		if err := t.json.saveGallery(galleryHash, newGalleryJSON); err != nil {
			logger.Errorf("[galleries] <%s> failed to save json: %s", galleryHash, err.Error())
		}
	}
}

func (t *ExportTask) ExportPerformers(workers int, repo models.ReaderRepository) {
	var performersWg sync.WaitGroup

	reader := repo.Performer()
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
		go t.exportPerformer(&performersWg, jobCh, repo)
	}

	for i, performer := range performers {
		index := i + 1
		logger.Progressf("[performers] %d of %d", index, len(performers))

		t.Mappings.Performers = append(t.Mappings.Performers, jsonschema.PathNameMapping{Name: performer.Name.String, Checksum: performer.Checksum})
		jobCh <- performer // feed workers
	}

	close(jobCh) // close channel so workers will know that no more jobs are available
	performersWg.Wait()

	logger.Infof("[performers] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportPerformer(wg *sync.WaitGroup, jobChan <-chan *models.Performer, repo models.ReaderRepository) {
	defer wg.Done()

	performerReader := repo.Performer()

	for p := range jobChan {
		newPerformerJSON, err := performer.ToJSON(performerReader, p)

		if err != nil {
			logger.Errorf("[performers] <%s> error getting performer JSON: %s", p.Checksum, err.Error())
			continue
		}

		tags, err := repo.Tag().FindByPerformerID(p.ID)
		if err != nil {
			logger.Errorf("[performers] <%s> error getting performer tags: %s", p.Checksum, err.Error())
			continue
		}

		newPerformerJSON.Tags = tag.GetNames(tags)

		if t.includeDependencies {
			t.tags.IDs = intslice.IntAppendUniques(t.tags.IDs, tag.GetIDs(tags))
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

func (t *ExportTask) ExportStudios(workers int, repo models.ReaderRepository) {
	var studiosWg sync.WaitGroup

	reader := repo.Studio()
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
		go t.exportStudio(&studiosWg, jobCh, repo)
	}

	for i, studio := range studios {
		index := i + 1
		logger.Progressf("[studios] %d of %d", index, len(studios))

		t.Mappings.Studios = append(t.Mappings.Studios, jsonschema.PathNameMapping{Name: studio.Name.String, Checksum: studio.Checksum})
		jobCh <- studio // feed workers
	}

	close(jobCh)
	studiosWg.Wait()

	logger.Infof("[studios] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportStudio(wg *sync.WaitGroup, jobChan <-chan *models.Studio, repo models.ReaderRepository) {
	defer wg.Done()

	studioReader := repo.Studio()

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

func (t *ExportTask) ExportTags(workers int, repo models.ReaderRepository) {
	var tagsWg sync.WaitGroup

	reader := repo.Tag()
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
		go t.exportTag(&tagsWg, jobCh, repo)
	}

	for i, tag := range tags {
		index := i + 1
		logger.Progressf("[tags] %d of %d", index, len(tags))

		// generate checksum on the fly by name, since we don't store it
		checksum := md5.FromString(tag.Name)

		t.Mappings.Tags = append(t.Mappings.Tags, jsonschema.PathNameMapping{Name: tag.Name, Checksum: checksum})
		jobCh <- tag // feed workers
	}

	close(jobCh)
	tagsWg.Wait()

	logger.Infof("[tags] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportTag(wg *sync.WaitGroup, jobChan <-chan *models.Tag, repo models.ReaderRepository) {
	defer wg.Done()

	tagReader := repo.Tag()

	for thisTag := range jobChan {
		newTagJSON, err := tag.ToJSON(tagReader, thisTag)

		if err != nil {
			logger.Errorf("[tags] <%s> error getting tag JSON: %s", thisTag.Name, err.Error())
			continue
		}

		// generate checksum on the fly by name, since we don't store it
		checksum := md5.FromString(thisTag.Name)

		tagJSON, err := t.json.getTag(checksum)
		if err == nil && jsonschema.CompareJSON(*tagJSON, *newTagJSON) {
			continue
		}

		if err := t.json.saveTag(checksum, newTagJSON); err != nil {
			logger.Errorf("[tags] <%s> failed to save json: %s", checksum, err.Error())
		}
	}
}

func (t *ExportTask) ExportMovies(workers int, repo models.ReaderRepository) {
	var moviesWg sync.WaitGroup

	reader := repo.Movie()
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
		go t.exportMovie(&moviesWg, jobCh, repo)
	}

	for i, movie := range movies {
		index := i + 1
		logger.Progressf("[movies] %d of %d", index, len(movies))

		t.Mappings.Movies = append(t.Mappings.Movies, jsonschema.PathNameMapping{Name: movie.Name.String, Checksum: movie.Checksum})
		jobCh <- movie // feed workers
	}

	close(jobCh)
	moviesWg.Wait()

	logger.Infof("[movies] export complete in %s. %d workers used.", time.Since(startTime), workers)

}
func (t *ExportTask) exportMovie(wg *sync.WaitGroup, jobChan <-chan *models.Movie, repo models.ReaderRepository) {
	defer wg.Done()

	movieReader := repo.Movie()
	studioReader := repo.Studio()

	for m := range jobChan {
		newMovieJSON, err := movie.ToJSON(movieReader, studioReader, m)

		if err != nil {
			logger.Errorf("[movies] <%s> error getting tag JSON: %s", m.Checksum, err.Error())
			continue
		}

		if t.includeDependencies {
			if m.StudioID.Valid {
				t.studios.IDs = intslice.IntAppendUnique(t.studios.IDs, int(m.StudioID.Int64))
			}
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

func (t *ExportTask) ExportScrapedItems(repo models.ReaderRepository) {
	qb := repo.ScrapedItem()
	sqb := repo.Studio()
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
			studio, _ := sqb.Find(int(scrapedItem.StudioID.Int64))
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
		updatedAt := json.JSONTime{Time: scrapedItem.UpdatedAt.Timestamp} // TODO keeping ruby format
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
