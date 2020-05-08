package manager

import (
	"path/filepath"
	"strconv"
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
	extLen := len(extensionsToScan)
	glb := "{"
	for i := 0; i < extLen-1; i++ { // append extensions and commas
		glb += extensionsToScan[i] + ","
	}
	if extLen >= 1 { // append last extension without comma
		glb += extensionsToScan[extLen-1]
	}
	glb += "}"
	return glb

}

func isGallery(pathname string) bool {
	for _, ext := range extensionsGallery {
		if filepath.Ext(pathname) == "."+ext {
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
		for i, path := range results {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}
			wg.Add(1)
			task := ScanTask{FilePath: path, UseFileMetadata: useFileMetadata}
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
		task := ImportTask{}
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
		task := ExportTask{}
		go task.Start(&wg)
		wg.Wait()
	}()
}

func (s *singleton) Generate(sprites bool, previews bool, markers bool, transcodes bool, thumbnails bool) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Generate)
	s.Status.indefiniteProgress()

	qb := models.NewSceneQueryBuilder()
	qg := models.NewGalleryQueryBuilder()
	//this.job.total = await ObjectionUtils.getCount(Scene);
	instance.Paths.Generated.EnsureTmpDir()

	go func() {
		defer s.returnToIdleState()

		scenes, err := qb.All()
		var galleries []*models.Gallery
		var gqErr error

		if err != nil {
			logger.Errorf("failed to get scenes for generate")
			return
		}

		delta := utils.Btoi(sprites) + utils.Btoi(previews) + utils.Btoi(markers) + utils.Btoi(transcodes)
		var wg sync.WaitGroup
		s.Status.Progress = 0
		lenScenes := len(scenes)
		total := lenScenes
		if thumbnails {
			galleries, gqErr = qg.All()
			if gqErr != nil {
				logger.Errorf("failed to get galleries for generate")
				return
			}
			total += len(galleries)
		}

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}
		totalsNeeded := s.neededGenerate(scenes, sprites, previews, markers, transcodes)
		if totalsNeeded == nil {
			logger.Infof("Taking too long to count content. Skipping...")
			logger.Infof("Generating content")
		} else {
			logger.Infof("Generating %d sprites %d previews %d markers %d transcodes", totalsNeeded.sprites, totalsNeeded.previews, totalsNeeded.markers, totalsNeeded.transcodes)
		}
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
			if sprites || previews || markers {
				instance.Paths.Generated.EmptyTmpDir()
			}

			if sprites {
				task := GenerateSpriteTask{Scene: *scene}
				go task.Start(&wg)
			}

			if previews {
				task := GeneratePreviewTask{Scene: *scene}
				go task.Start(&wg)
			}

			if markers {
				task := GenerateMarkersTask{Scene: *scene}
				go task.Start(&wg)
			}

			if transcodes {
				task := GenerateTranscodeTask{Scene: *scene}
				go task.Start(&wg)
			}

			wg.Wait()
		}

		if thumbnails {
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
				task := GenerateGthumbsTask{Gallery: *gallery}
				go task.Start(&wg)
				wg.Wait()
			}
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
			Scene:        *scene,
			ScreenshotAt: at,
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

			task := CleanTask{Scene: scene}
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
	sprites    int64
	previews   int64
	markers    int64
	transcodes int64
}

func (s *singleton) neededGenerate(scenes []*models.Scene, sprites, previews, markers, transcodes bool) *totalsGenerate {

	var totals totalsGenerate
	const timeout = 90 * time.Second

	// create a control channel through which to signal the counting loop when the timeout is reached
	chTimeout := make(chan struct{})

	//run the timeout function in a separate thread
	go func() {
		time.Sleep(timeout)
		chTimeout <- struct{}{}
	}()

	logger.Infof("Counting content to generate...")
	for _, scene := range scenes {
		if scene != nil {
			if sprites {
				task := GenerateSpriteTask{Scene: *scene}
				if !task.doesSpriteExist(task.Scene.Checksum) {
					totals.sprites++
				}
			}

			if previews {
				task := GeneratePreviewTask{Scene: *scene}
				if !task.doesPreviewExist(task.Scene.Checksum) {
					totals.previews++
				}
			}

			if markers {
				task := GenerateMarkersTask{Scene: *scene}
				totals.markers += int64(task.isMarkerNeeded())

			}
			if transcodes {
				task := GenerateTranscodeTask{Scene: *scene}
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
