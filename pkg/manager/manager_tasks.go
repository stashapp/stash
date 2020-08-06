package manager

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

var extensionsToScan = []string{"zip", "m4v", "mp4", "mov", "wmv", "avi", "mpg", "mpeg", "rmvb", "rm", "flv", "asf", "mkv", "webm"}
var extensionsGallery = []string{"zip"}

func constructGlob() string { // create a sequence for glob doublestar from our extensions
	var extList []string
	for _, ext := range extensionsToScan {
		extList = append(extList, strings.ToLower(ext))
		extList = append(extList, strings.ToUpper(ext))
	}
	return "{" + strings.Join(extList, ",") + "}"
}

func isGallery(pathname string) bool {
	for _, ext := range extensionsGallery {
		if strings.ToLower(filepath.Ext(pathname)) == "."+strings.ToLower(ext) {
			return true
		}
	}
	return false
}

type TaskStatus struct {
	Status     JobStatus
	Progress   float64
	LastUpdate time.Time
	stopping   bool
	upTo       int
	total      int
}

func (t *TaskStatus) Stop() bool {
	t.stopping = true
	t.updated()
	return true
}

func (t *TaskStatus) SetStatus(s JobStatus) {
	t.Status = s
	t.updated()
}

func (t *TaskStatus) setProgress(upTo int, total int) {
	if total == 0 {
		t.Progress = 1
	}
	t.upTo = upTo
	t.total = total
	t.Progress = float64(upTo) / float64(total)
	t.updated()
}

func (t *TaskStatus) incrementProgress() {
	t.setProgress(t.upTo+1, t.total)
}

func (t *TaskStatus) indefiniteProgress() {
	t.Progress = -1
	t.updated()
}

func (t *TaskStatus) updated() {
	t.LastUpdate = time.Now()
}

func (s *singleton) Scan(useFileMetadata bool) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Scan)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var results []string
		for _, path := range config.GetStashPaths() {
			globPath := filepath.Join(path, "**/*."+constructGlob())
			globResults, _ := doublestar.Glob(globPath)
			results = append(results, globResults...)
		}

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		results, _ = excludeFiles(results, config.GetExcludes())
		total := len(results)
		logger.Infof("Starting scan of %d files. %d New files found", total, s.neededScan(results))

		var wg sync.WaitGroup
		s.Status.Progress = 0
		fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
		calculateMD5 := config.IsCalculateMD5()
		for i, path := range results {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}
			wg.Add(1)
			task := ScanTask{FilePath: path, UseFileMetadata: useFileMetadata, fileNamingAlgorithm: fileNamingAlgo, calculateMD5: calculateMD5}
			go task.Start(&wg)
			wg.Wait()
		}

		logger.Info("Finished scan")
		for _, path := range results {
			if isGallery(path) {
				wg.Add(1)
				task := ScanTask{FilePath: path, UseFileMetadata: false}
				go task.associateGallery(&wg)
				wg.Wait()
			}
		}
		logger.Info("Finished gallery association")
	}()
}

func (s *singleton) Import() {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Import)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var wg sync.WaitGroup
		wg.Add(1)
		task := ImportTask{fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm()}
		go task.Start(&wg)
		wg.Wait()
	}()
}

func (s *singleton) Export() {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Export)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var wg sync.WaitGroup
		wg.Add(1)
		task := ExportTask{fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm()}
		go task.Start(&wg)
		wg.Wait()
	}()
}

func setGeneratePreviewOptionsInput(optionsInput *models.GeneratePreviewOptionsInput) {
	if optionsInput.PreviewSegments == nil {
		val := config.GetPreviewSegments()
		optionsInput.PreviewSegments = &val
	}

	if optionsInput.PreviewSegmentDuration == nil {
		val := config.GetPreviewSegmentDuration()
		optionsInput.PreviewSegmentDuration = &val
	}

	if optionsInput.PreviewExcludeStart == nil {
		val := config.GetPreviewExcludeStart()
		optionsInput.PreviewExcludeStart = &val
	}

	if optionsInput.PreviewExcludeEnd == nil {
		val := config.GetPreviewExcludeEnd()
		optionsInput.PreviewExcludeEnd = &val
	}

	if optionsInput.PreviewPreset == nil {
		val := config.GetPreviewPreset()
		optionsInput.PreviewPreset = &val
	}
}

func (s *singleton) Generate(input models.GenerateMetadataInput) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Generate)
	s.Status.indefiniteProgress()

	qb := models.NewSceneQueryBuilder()
	qg := models.NewGalleryQueryBuilder()
	mqb := models.NewSceneMarkerQueryBuilder()

	//this.job.total = await ObjectionUtils.getCount(Scene);
	instance.Paths.Generated.EnsureTmpDir()

	galleryIDs := utils.StringSliceToIntSlice(input.GalleryIDs)
	sceneIDs := utils.StringSliceToIntSlice(input.SceneIDs)
	markerIDs := utils.StringSliceToIntSlice(input.MarkerIDs)

	go func() {
		defer s.returnToIdleState()

		var scenes []*models.Scene
		var err error

		if len(sceneIDs) > 0 {
			scenes, err = qb.FindMany(sceneIDs)
		} else {
			scenes, err = qb.All()
		}

		if err != nil {
			logger.Errorf("failed to get scenes for generate")
			return
		}

		delta := utils.Btoi(input.Sprites) + utils.Btoi(input.Previews) + utils.Btoi(input.Markers) + utils.Btoi(input.Transcodes)
		var wg sync.WaitGroup

		s.Status.Progress = 0
		lenScenes := len(scenes)
		total := lenScenes

		var galleries []*models.Gallery
		if input.Thumbnails {
			if len(galleryIDs) > 0 {
				galleries, err = qg.FindMany(galleryIDs)
			} else {
				galleries, err = qg.All()
			}

			if err != nil {
				logger.Errorf("failed to get galleries for generate")
				return
			}
			total += len(galleries)
		}

		var markers []*models.SceneMarker
		if len(markerIDs) > 0 {
			markers, err = mqb.FindMany(markerIDs)

			total += len(markers)
		}

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		totalsNeeded := s.neededGenerate(scenes, input)
		if totalsNeeded == nil {
			logger.Infof("Taking too long to count content. Skipping...")
			logger.Infof("Generating content")
		} else {
			logger.Infof("Generating %d sprites %d previews %d image previews %d markers %d transcodes", totalsNeeded.sprites, totalsNeeded.previews, totalsNeeded.imagePreviews, totalsNeeded.markers, totalsNeeded.transcodes)
		}

		fileNamingAlgo := config.GetVideoFileNamingAlgorithm()

		overwrite := false
		if input.Overwrite != nil {
			overwrite = *input.Overwrite
		}

		generatePreviewOptions := input.PreviewOptions
		if generatePreviewOptions == nil {
			generatePreviewOptions = &models.GeneratePreviewOptionsInput{}
		}
		setGeneratePreviewOptionsInput(generatePreviewOptions)

		for i, scene := range scenes {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping generate")
				continue
			}

			wg.Add(delta)

			// Clear the tmp directory for each scene
			if input.Sprites || input.Previews || input.Markers {
				instance.Paths.Generated.EmptyTmpDir()
			}

			if input.Sprites {
				task := GenerateSpriteTask{Scene: *scene, Overwrite: overwrite, fileNamingAlgorithm: fileNamingAlgo}
				go task.Start(&wg)
			}

			if input.Previews {
				task := GeneratePreviewTask{
					Scene:               *scene,
					ImagePreview:        input.ImagePreviews,
					Options:             *generatePreviewOptions,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				go task.Start(&wg)
			}

			if input.Markers {
				task := GenerateMarkersTask{Scene: scene, Overwrite: overwrite, fileNamingAlgorithm: fileNamingAlgo}
				go task.Start(&wg)
			}

			if input.Transcodes {
				task := GenerateTranscodeTask{Scene: *scene, Overwrite: overwrite, fileNamingAlgorithm: fileNamingAlgo}
				go task.Start(&wg)
			}

			wg.Wait()
		}

		if input.Thumbnails {
			logger.Infof("Generating thumbnails for the galleries")
			for i, gallery := range galleries {
				s.Status.setProgress(lenScenes+i, total)
				if s.Status.stopping {
					logger.Info("Stopping due to user request")
					return
				}

				if gallery == nil {
					logger.Errorf("nil gallery, skipping generate")
					continue
				}

				wg.Add(1)
				task := GenerateGthumbsTask{Gallery: *gallery, Overwrite: overwrite}
				go task.Start(&wg)
				wg.Wait()
			}
		}

		for i, marker := range markers {
			s.Status.setProgress(lenScenes+len(galleries)+i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if marker == nil {
				logger.Errorf("nil marker, skipping generate")
				continue
			}

			wg.Add(1)
			task := GenerateMarkersTask{Marker: marker, Overwrite: overwrite, fileNamingAlgorithm: fileNamingAlgo}
			go task.Start(&wg)
			wg.Wait()
		}

		logger.Infof("Generate finished")
	}()
}

func (s *singleton) GenerateDefaultScreenshot(sceneId string) {
	s.generateScreenshot(sceneId, nil)
}

func (s *singleton) GenerateScreenshot(sceneId string, at float64) {
	s.generateScreenshot(sceneId, &at)
}

// generate default screenshot if at is nil
func (s *singleton) generateScreenshot(sceneId string, at *float64) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Generate)
	s.Status.indefiniteProgress()

	qb := models.NewSceneQueryBuilder()
	instance.Paths.Generated.EnsureTmpDir()

	go func() {
		defer s.returnToIdleState()

		sceneIdInt, err := strconv.Atoi(sceneId)
		if err != nil {
			logger.Errorf("Error parsing scene id %s: %s", sceneId, err.Error())
			return
		}

		scene, err := qb.Find(sceneIdInt)
		if err != nil || scene == nil {
			logger.Errorf("failed to get scene for generate")
			return
		}

		task := GenerateScreenshotTask{
			Scene:               *scene,
			ScreenshotAt:        at,
			fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm(),
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go task.Start(&wg)

		wg.Wait()

		logger.Infof("Generate finished")
	}()
}

func (s *singleton) AutoTag(performerIds []string, studioIds []string, tagIds []string) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(AutoTag)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		// calculate work load
		performerCount := len(performerIds)
		studioCount := len(studioIds)
		tagCount := len(tagIds)

		performerQuery := models.NewPerformerQueryBuilder()
		studioQuery := models.NewTagQueryBuilder()
		tagQuery := models.NewTagQueryBuilder()

		const wildcard = "*"
		var err error
		if performerCount == 1 && performerIds[0] == wildcard {
			performerCount, err = performerQuery.Count()
			if err != nil {
				logger.Errorf("Error getting performer count: %s", err.Error())
			}
		}
		if studioCount == 1 && studioIds[0] == wildcard {
			studioCount, err = studioQuery.Count()
			if err != nil {
				logger.Errorf("Error getting studio count: %s", err.Error())
			}
		}
		if tagCount == 1 && tagIds[0] == wildcard {
			tagCount, err = tagQuery.Count()
			if err != nil {
				logger.Errorf("Error getting tag count: %s", err.Error())
			}
		}

		total := performerCount + studioCount + tagCount
		s.Status.setProgress(0, total)

		s.autoTagPerformers(performerIds)
		s.autoTagStudios(studioIds)
		s.autoTagTags(tagIds)
	}()
}

func (s *singleton) autoTagPerformers(performerIds []string) {
	performerQuery := models.NewPerformerQueryBuilder()

	var wg sync.WaitGroup
	for _, performerId := range performerIds {
		var performers []*models.Performer
		if performerId == "*" {
			var err error
			performers, err = performerQuery.All()
			if err != nil {
				logger.Errorf("Error querying performers: %s", err.Error())
				continue
			}
		} else {
			performerIdInt, err := strconv.Atoi(performerId)
			if err != nil {
				logger.Errorf("Error parsing performer id %s: %s", performerId, err.Error())
				continue
			}

			performer, err := performerQuery.Find(performerIdInt)
			if err != nil {
				logger.Errorf("Error finding performer id %s: %s", performerId, err.Error())
				continue
			}
			performers = append(performers, performer)
		}

		for _, performer := range performers {
			wg.Add(1)
			task := AutoTagPerformerTask{performer: performer}
			go task.Start(&wg)
			wg.Wait()

			s.Status.incrementProgress()
		}
	}
}

func (s *singleton) autoTagStudios(studioIds []string) {
	studioQuery := models.NewStudioQueryBuilder()

	var wg sync.WaitGroup
	for _, studioId := range studioIds {
		var studios []*models.Studio
		if studioId == "*" {
			var err error
			studios, err = studioQuery.All()
			if err != nil {
				logger.Errorf("Error querying studios: %s", err.Error())
				continue
			}
		} else {
			studioIdInt, err := strconv.Atoi(studioId)
			if err != nil {
				logger.Errorf("Error parsing studio id %s: %s", studioId, err.Error())
				continue
			}

			studio, err := studioQuery.Find(studioIdInt, nil)
			if err != nil {
				logger.Errorf("Error finding studio id %s: %s", studioId, err.Error())
				continue
			}
			studios = append(studios, studio)
		}

		for _, studio := range studios {
			wg.Add(1)
			task := AutoTagStudioTask{studio: studio}
			go task.Start(&wg)
			wg.Wait()

			s.Status.incrementProgress()
		}
	}
}

func (s *singleton) autoTagTags(tagIds []string) {
	tagQuery := models.NewTagQueryBuilder()

	var wg sync.WaitGroup
	for _, tagId := range tagIds {
		var tags []*models.Tag
		if tagId == "*" {
			var err error
			tags, err = tagQuery.All()
			if err != nil {
				logger.Errorf("Error querying tags: %s", err.Error())
				continue
			}
		} else {
			tagIdInt, err := strconv.Atoi(tagId)
			if err != nil {
				logger.Errorf("Error parsing tag id %s: %s", tagId, err.Error())
				continue
			}

			tag, err := tagQuery.Find(tagIdInt, nil)
			if err != nil {
				logger.Errorf("Error finding tag id %s: %s", tagId, err.Error())
				continue
			}
			tags = append(tags, tag)
		}

		for _, tag := range tags {
			wg.Add(1)
			task := AutoTagTagTask{tag: tag}
			go task.Start(&wg)
			wg.Wait()

			s.Status.incrementProgress()
		}
	}
}

func (s *singleton) Clean() {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Clean)
	s.Status.indefiniteProgress()

	qb := models.NewSceneQueryBuilder()
	gqb := models.NewGalleryQueryBuilder()
	go func() {
		defer s.returnToIdleState()

		logger.Infof("Starting cleaning of tracked files")
		scenes, err := qb.All()
		if err != nil {
			logger.Errorf("failed to fetch list of scenes for cleaning")
			return
		}

		galleries, err := gqb.All()
		if err != nil {
			logger.Errorf("failed to fetch list of galleries for cleaning")
			return
		}

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		var wg sync.WaitGroup
		s.Status.Progress = 0
		total := len(scenes) + len(galleries)
		fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
		for i, scene := range scenes {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping Clean")
				continue
			}

			wg.Add(1)

			task := CleanTask{Scene: scene, fileNamingAlgorithm: fileNamingAlgo}
			go task.Start(&wg)
			wg.Wait()
		}

		for i, gallery := range galleries {
			s.Status.setProgress(len(scenes)+i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if gallery == nil {
				logger.Errorf("nil gallery, skipping Clean")
				continue
			}

			wg.Add(1)

			task := CleanTask{Gallery: gallery}
			go task.Start(&wg)
			wg.Wait()
		}

		logger.Info("Finished Cleaning")
	}()
}

func (s *singleton) MigrateHash() {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Migrate)
	s.Status.indefiniteProgress()

	qb := models.NewSceneQueryBuilder()

	go func() {
		defer s.returnToIdleState()

		fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
		logger.Infof("Migrating generated files for %s naming hash", fileNamingAlgo.String())

		scenes, err := qb.All()
		if err != nil {
			logger.Errorf("failed to fetch list of scenes for migration")
			return
		}

		var wg sync.WaitGroup
		s.Status.Progress = 0
		total := len(scenes)

		for i, scene := range scenes {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping migrate")
				continue
			}

			wg.Add(1)

			task := MigrateHashTask{Scene: scene, fileNamingAlgorithm: fileNamingAlgo}
			go task.Start(&wg)
			wg.Wait()
		}

		logger.Info("Finished migrating")
	}()
}

func (s *singleton) returnToIdleState() {
	if r := recover(); r != nil {
		logger.Info("recovered from ", r)
	}

	if s.Status.Status == Generate {
		instance.Paths.Generated.RemoveTmpDir()
	}
	s.Status.SetStatus(Idle)
	s.Status.indefiniteProgress()
	s.Status.stopping = false
}

func (s *singleton) neededScan(paths []string) int64 {
	var neededScans int64

	for _, path := range paths {
		task := ScanTask{FilePath: path}
		if !task.doesPathExist() {
			neededScans++
		}
	}
	return neededScans
}

type totalsGenerate struct {
	sprites       int64
	previews      int64
	imagePreviews int64
	markers       int64
	transcodes    int64
}

func (s *singleton) neededGenerate(scenes []*models.Scene, input models.GenerateMetadataInput) *totalsGenerate {

	var totals totalsGenerate
	const timeout = 90 * time.Second

	// create a control channel through which to signal the counting loop when the timeout is reached
	chTimeout := make(chan struct{})

	//run the timeout function in a separate thread
	go func() {
		time.Sleep(timeout)
		chTimeout <- struct{}{}
	}()

	fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
	overwrite := false
	if input.Overwrite != nil {
		overwrite = *input.Overwrite
	}

	logger.Infof("Counting content to generate...")
	for _, scene := range scenes {
		if scene != nil {
			if input.Sprites {
				task := GenerateSpriteTask{
					Scene:               *scene,
					fileNamingAlgorithm: fileNamingAlgo,
				}

				if overwrite || task.required() {
					totals.sprites++
				}
			}

			if input.Previews {
				task := GeneratePreviewTask{
					Scene:               *scene,
					ImagePreview:        input.ImagePreviews,
					fileNamingAlgorithm: fileNamingAlgo,
				}

				sceneHash := scene.GetHash(task.fileNamingAlgorithm)
				if overwrite || !task.doesVideoPreviewExist(sceneHash) {
					totals.previews++
				}

				if input.ImagePreviews && (overwrite || !task.doesImagePreviewExist(sceneHash)) {
					totals.imagePreviews++
				}
			}

			if input.Markers {
				task := GenerateMarkersTask{
					Scene:               scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				totals.markers += int64(task.isMarkerNeeded())
			}

			if input.Transcodes {
				task := GenerateTranscodeTask{
					Scene:               *scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				if task.isTranscodeNeeded() {
					totals.transcodes++
				}
			}
		}
		//check for timeout
		select {
		case <-chTimeout:
			return nil
		default:
		}

	}
	return &totals
}
