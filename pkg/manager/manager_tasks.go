package manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/remeh/sizedwaitgroup"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func isGallery(pathname string) bool {
	gExt := config.GetInstance().GetGalleryExtensions()
	return matchExtension(pathname, gExt)
}

func isVideo(pathname string) bool {
	vidExt := config.GetInstance().GetVideoExtensions()
	return matchExtension(pathname, vidExt)
}

func isImage(pathname string) bool {
	imgExt := config.GetInstance().GetImageExtensions()
	return matchExtension(pathname, imgExt)
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

func (t *TaskStatus) setProgressPercent(progress float64) {
	if progress != t.Progress {
		t.Progress = progress
		t.updated()
	}
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

func getScanPaths(inputPaths []string) []*models.StashConfig {
	if len(inputPaths) == 0 {
		return config.GetInstance().GetStashPaths()
	}

	var ret []*models.StashConfig
	for _, p := range inputPaths {
		s := getStashFromDirPath(p)
		if s == nil {
			logger.Warnf("%s is not in the configured stash paths", p)
			continue
		}

		// make a copy, changing the path
		ss := *s
		ss.Path = p
		ret = append(ret, &ss)
	}

	return ret
}

func (s *singleton) neededScan(paths []*models.StashConfig) (total *int, newFiles *int) {
	const timeout = 90 * time.Second

	// create a control channel through which to signal the counting loop when the timeout is reached
	chTimeout := time.After(timeout)

	logger.Infof("Counting files to scan...")

	t := 0
	n := 0

	timeoutErr := errors.New("timed out")

	for _, sp := range paths {
		err := walkFilesToScan(sp, func(path string, info os.FileInfo, err error) error {
			t++
			task := ScanTask{FilePath: path, TxnManager: s.TxnManager}
			if !task.doesPathExist() {
				n++
			}

			//check for timeout
			select {
			case <-chTimeout:
				return timeoutErr
			default:
			}

			// check stop
			if s.Status.stopping {
				return timeoutErr
			}

			return nil
		})

		if err == timeoutErr {
			// timeout should return nil counts
			return nil, nil
		}

		if err != nil {
			logger.Errorf("Error encountered counting files to scan: %s", err.Error())
			return nil, nil
		}
	}

	return &t, &n
}

func (s *singleton) Scan(input models.ScanMetadataInput) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Scan)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		paths := getScanPaths(input.Paths)

		total, newFiles := s.neededScan(paths)

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		if total == nil || newFiles == nil {
			logger.Infof("Taking too long to count content. Skipping...")
			logger.Infof("Starting scan")
		} else {
			logger.Infof("Starting scan of %d files. %d New files found", *total, *newFiles)
		}

		start := time.Now()
		config := config.GetInstance()
		parallelTasks := config.GetParallelTasksWithAutoDetection()
		logger.Infof("Scan started with %d parallel tasks", parallelTasks)
		wg := sizedwaitgroup.New(parallelTasks)

		s.Status.Progress = 0
		fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
		calculateMD5 := config.IsCalculateMD5()

		i := 0
		stoppingErr := errors.New("stopping")
		var err error

		var galleries []string

		for _, sp := range paths {
			err = walkFilesToScan(sp, func(path string, info os.FileInfo, err error) error {
				if total != nil {
					s.Status.setProgress(i, *total)
					i++
				}

				if s.Status.stopping {
					return stoppingErr
				}

				if isGallery(path) {
					galleries = append(galleries, path)
				}

				instance.Paths.Generated.EnsureTmpDir()

				wg.Add()
				task := ScanTask{
					TxnManager:           s.TxnManager,
					FilePath:             path,
					UseFileMetadata:      utils.IsTrue(input.UseFileMetadata),
					StripFileExtension:   utils.IsTrue(input.StripFileExtension),
					fileNamingAlgorithm:  fileNamingAlgo,
					calculateMD5:         calculateMD5,
					GeneratePreview:      utils.IsTrue(input.ScanGeneratePreviews),
					GenerateImagePreview: utils.IsTrue(input.ScanGenerateImagePreviews),
					GenerateSprite:       utils.IsTrue(input.ScanGenerateSprites),
					GeneratePhash:        utils.IsTrue(input.ScanGeneratePhashes),
				}
				go task.Start(&wg)

				return nil
			})

			if err == stoppingErr {
				logger.Info("Stopping due to user request")
				break
			}

			if err != nil {
				logger.Errorf("Error encountered scanning files: %s", err.Error())
				break
			}
		}

		wg.Wait()
		instance.Paths.Generated.EmptyTmpDir()
		elapsed := time.Since(start)
		logger.Info(fmt.Sprintf("Scan finished (%s)", elapsed))

		if s.Status.stopping || err != nil {
			return
		}

		for _, path := range galleries {
			wg.Add()
			task := ScanTask{
				TxnManager:      s.TxnManager,
				FilePath:        path,
				UseFileMetadata: false,
			}
			go task.associateGallery(&wg)
			wg.Wait()
		}
		logger.Info("Finished gallery association")
	}()
}

func (s *singleton) Import() error {
	config := config.GetInstance()
	metadataPath := config.GetMetadataPath()
	if metadataPath == "" {
		return errors.New("metadata path must be set in config")
	}

	if s.Status.Status != Idle {
		return nil
	}
	s.Status.SetStatus(Import)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var wg sync.WaitGroup
		wg.Add(1)

		task := ImportTask{
			txnManager:          s.TxnManager,
			BaseDir:             metadataPath,
			Reset:               true,
			DuplicateBehaviour:  models.ImportDuplicateEnumFail,
			MissingRefBehaviour: models.ImportMissingRefEnumFail,
			fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm(),
		}
		go task.Start(&wg)
		wg.Wait()
	}()

	return nil
}

func (s *singleton) Export() error {
	config := config.GetInstance()
	metadataPath := config.GetMetadataPath()
	if metadataPath == "" {
		return errors.New("metadata path must be set in config")
	}

	if s.Status.Status != Idle {
		return nil
	}
	s.Status.SetStatus(Export)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var wg sync.WaitGroup
		wg.Add(1)
		task := ExportTask{
			txnManager:          s.TxnManager,
			full:                true,
			fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm(),
		}
		go task.Start(&wg)
		wg.Wait()
	}()

	return nil
}

func (s *singleton) RunSingleTask(t Task) (*sync.WaitGroup, error) {
	if s.Status.Status != Idle {
		return nil, errors.New("task already running")
	}

	s.Status.SetStatus(t.GetStatus())
	s.Status.indefiniteProgress()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer s.returnToIdleState()

		go t.Start(&wg)
		wg.Wait()
	}()

	return &wg, nil
}

func setGeneratePreviewOptionsInput(optionsInput *models.GeneratePreviewOptionsInput) {
	config := config.GetInstance()
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

	//this.job.total = await ObjectionUtils.getCount(Scene);
	instance.Paths.Generated.EnsureTmpDir()

	sceneIDs, err := utils.StringSliceToIntSlice(input.SceneIDs)
	if err != nil {
		logger.Error(err.Error())
	}
	markerIDs, err := utils.StringSliceToIntSlice(input.MarkerIDs)
	if err != nil {
		logger.Error(err.Error())
	}

	go func() {
		defer s.returnToIdleState()

		var scenes []*models.Scene
		var err error
		var markers []*models.SceneMarker

		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			qb := r.Scene()
			if len(sceneIDs) > 0 {
				scenes, err = qb.FindMany(sceneIDs)
			} else {
				scenes, err = qb.All()
			}

			if err != nil {
				return err
			}

			if len(markerIDs) > 0 {
				markers, err = r.SceneMarker().FindMany(markerIDs)
				if err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			return
		}

		config := config.GetInstance()
		parallelTasks := config.GetParallelTasksWithAutoDetection()

		logger.Infof("Generate started with %d parallel tasks", parallelTasks)
		wg := sizedwaitgroup.New(parallelTasks)

		s.Status.Progress = 0
		lenScenes := len(scenes)
		total := lenScenes + len(markers)

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		totalsNeeded := s.neededGenerate(scenes, input)
		if totalsNeeded == nil {
			logger.Infof("Taking too long to count content. Skipping...")
			logger.Infof("Generating content")
		} else {
			logger.Infof("Generating %d sprites %d previews %d image previews %d markers %d transcodes %d phashes", totalsNeeded.sprites, totalsNeeded.previews, totalsNeeded.imagePreviews, totalsNeeded.markers, totalsNeeded.transcodes, totalsNeeded.phashes)
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

		// Start measuring how long the generate has taken. (consider moving this up)
		start := time.Now()
		instance.Paths.Generated.EnsureTmpDir()

		for i, scene := range scenes {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				wg.Wait()
				instance.Paths.Generated.EmptyTmpDir()
				return
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping generate")
				continue
			}

			if input.Sprites {
				task := GenerateSpriteTask{
					Scene:               *scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				wg.Add()
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
				wg.Add()
				go task.Start(&wg)
			}

			if input.Markers {
				wg.Add()
				task := GenerateMarkersTask{
					TxnManager:          s.TxnManager,
					Scene:               scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				go task.Start(&wg)
			}

			if input.Transcodes {
				wg.Add()
				task := GenerateTranscodeTask{
					Scene:               *scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				go task.Start(&wg)
			}

			if input.Phashes {
				task := GeneratePhashTask{
					Scene:               *scene,
					fileNamingAlgorithm: fileNamingAlgo,
					txnManager:          s.TxnManager,
				}
				wg.Add()
				go task.Start(&wg)
			}
		}

		wg.Wait()

		for i, marker := range markers {
			s.Status.setProgress(lenScenes+i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				wg.Wait()
				instance.Paths.Generated.EmptyTmpDir()
				elapsed := time.Since(start)
				logger.Info(fmt.Sprintf("Generate finished (%s)", elapsed))
				return
			}

			if marker == nil {
				logger.Errorf("nil marker, skipping generate")
				continue
			}

			wg.Add()
			task := GenerateMarkersTask{
				TxnManager:          s.TxnManager,
				Marker:              marker,
				Overwrite:           overwrite,
				fileNamingAlgorithm: fileNamingAlgo,
			}
			go task.Start(&wg)
		}

		wg.Wait()

		instance.Paths.Generated.EmptyTmpDir()
		elapsed := time.Since(start)
		logger.Info(fmt.Sprintf("Generate finished (%s)", elapsed))
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

	instance.Paths.Generated.EnsureTmpDir()

	go func() {
		defer s.returnToIdleState()

		sceneIdInt, err := strconv.Atoi(sceneId)
		if err != nil {
			logger.Errorf("Error parsing scene id %s: %s", sceneId, err.Error())
			return
		}

		var scene *models.Scene
		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			var err error
			scene, err = r.Scene().Find(sceneIdInt)
			return err
		}); err != nil || scene == nil {
			logger.Errorf("failed to get scene for generate: %s", err.Error())
			return
		}

		task := GenerateScreenshotTask{
			txnManager:          s.TxnManager,
			Scene:               *scene,
			ScreenshotAt:        at,
			fileNamingAlgorithm: config.GetInstance().GetVideoFileNamingAlgorithm(),
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go task.Start(&wg)

		wg.Wait()

		logger.Infof("Generate screenshot finished")
	}()
}

func (s *singleton) AutoTag(input models.AutoTagMetadataInput) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(AutoTag)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		performerIds := input.Performers
		studioIds := input.Studios
		tagIds := input.Tags

		// calculate work load
		performerCount := len(performerIds)
		studioCount := len(studioIds)
		tagCount := len(tagIds)

		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			performerQuery := r.Performer()
			studioQuery := r.Studio()
			tagQuery := r.Tag()

			const wildcard = "*"
			var err error
			if performerCount == 1 && performerIds[0] == wildcard {
				performerCount, err = performerQuery.Count()
				if err != nil {
					return fmt.Errorf("Error getting performer count: %s", err.Error())
				}
			}
			if studioCount == 1 && studioIds[0] == wildcard {
				studioCount, err = studioQuery.Count()
				if err != nil {
					return fmt.Errorf("Error getting studio count: %s", err.Error())
				}
			}
			if tagCount == 1 && tagIds[0] == wildcard {
				tagCount, err = tagQuery.Count()
				if err != nil {
					return fmt.Errorf("Error getting tag count: %s", err.Error())
				}
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			return
		}

		total := performerCount + studioCount + tagCount
		s.Status.setProgress(0, total)

		s.autoTagPerformers(input.Paths, performerIds)
		s.autoTagStudios(input.Paths, studioIds)
		s.autoTagTags(input.Paths, tagIds)
	}()
}

func (s *singleton) autoTagPerformers(paths []string, performerIds []string) {
	var wg sync.WaitGroup
	for _, performerId := range performerIds {
		var performers []*models.Performer

		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			performerQuery := r.Performer()

			if performerId == "*" {
				var err error
				performers, err = performerQuery.All()
				if err != nil {
					return fmt.Errorf("Error querying performers: %s", err.Error())
				}
			} else {
				performerIdInt, err := strconv.Atoi(performerId)
				if err != nil {
					return fmt.Errorf("Error parsing performer id %s: %s", performerId, err.Error())
				}

				performer, err := performerQuery.Find(performerIdInt)
				if err != nil {
					return fmt.Errorf("Error finding performer id %s: %s", performerId, err.Error())
				}
				performers = append(performers, performer)
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			continue
		}

		for _, performer := range performers {
			wg.Add(1)
			task := AutoTagPerformerTask{
				AutoTagTask: AutoTagTask{
					txnManager: s.TxnManager,
					paths:      paths,
				},
				performer: performer,
			}
			go task.Start(&wg)
			wg.Wait()

			s.Status.incrementProgress()
		}
	}
}

func (s *singleton) autoTagStudios(paths []string, studioIds []string) {
	var wg sync.WaitGroup
	for _, studioId := range studioIds {
		var studios []*models.Studio

		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			studioQuery := r.Studio()
			if studioId == "*" {
				var err error
				studios, err = studioQuery.All()
				if err != nil {
					return fmt.Errorf("Error querying studios: %s", err.Error())
				}
			} else {
				studioIdInt, err := strconv.Atoi(studioId)
				if err != nil {
					return fmt.Errorf("Error parsing studio id %s: %s", studioId, err.Error())
				}

				studio, err := studioQuery.Find(studioIdInt)
				if err != nil {
					return fmt.Errorf("Error finding studio id %s: %s", studioId, err.Error())
				}
				studios = append(studios, studio)
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			continue
		}

		for _, studio := range studios {
			wg.Add(1)
			task := AutoTagStudioTask{
				AutoTagTask: AutoTagTask{
					txnManager: s.TxnManager,
					paths:      paths,
				},
				studio: studio,
			}
			go task.Start(&wg)
			wg.Wait()

			s.Status.incrementProgress()
		}
	}
}

func (s *singleton) autoTagTags(paths []string, tagIds []string) {
	var wg sync.WaitGroup
	for _, tagId := range tagIds {
		var tags []*models.Tag
		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			tagQuery := r.Tag()
			if tagId == "*" {
				var err error
				tags, err = tagQuery.All()
				if err != nil {
					return fmt.Errorf("Error querying tags: %s", err.Error())
				}
			} else {
				tagIdInt, err := strconv.Atoi(tagId)
				if err != nil {
					return fmt.Errorf("Error parsing tag id %s: %s", tagId, err.Error())
				}

				tag, err := tagQuery.Find(tagIdInt)
				if err != nil {
					return fmt.Errorf("Error finding tag id %s: %s", tagId, err.Error())
				}
				tags = append(tags, tag)
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			continue
		}

		for _, tag := range tags {
			wg.Add(1)
			task := AutoTagTagTask{
				AutoTagTask: AutoTagTask{
					txnManager: s.TxnManager,
					paths:      paths,
				},
				tag: tag,
			}
			go task.Start(&wg)
			wg.Wait()

			s.Status.incrementProgress()
		}
	}
}

func (s *singleton) Clean(input models.CleanMetadataInput) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Clean)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var scenes []*models.Scene
		var images []*models.Image
		var galleries []*models.Gallery

		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			qb := r.Scene()
			iqb := r.Image()
			gqb := r.Gallery()

			logger.Infof("Starting cleaning of tracked files")
			if input.DryRun {
				logger.Infof("Running in Dry Mode")
			}
			var err error

			scenes, err = qb.All()

			if err != nil {
				return errors.New("failed to fetch list of scenes for cleaning")
			}

			images, err = iqb.All()
			if err != nil {
				return errors.New("failed to fetch list of images for cleaning")
			}

			galleries, err = gqb.All()
			if err != nil {
				return errors.New("failed to fetch list of galleries for cleaning")
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			return
		}

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		var wg sync.WaitGroup
		s.Status.Progress = 0
		total := len(scenes) + len(images) + len(galleries)
		fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()
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

			task := CleanTask{
				TxnManager:          s.TxnManager,
				Scene:               scene,
				fileNamingAlgorithm: fileNamingAlgo,
			}
			go task.Start(&wg, input.DryRun)
			wg.Wait()
		}

		for i, img := range images {
			s.Status.setProgress(len(scenes)+i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if img == nil {
				logger.Errorf("nil image, skipping Clean")
				continue
			}

			wg.Add(1)

			task := CleanTask{
				TxnManager: s.TxnManager,
				Image:      img,
			}
			go task.Start(&wg, input.DryRun)
			wg.Wait()
		}

		for i, gallery := range galleries {
			s.Status.setProgress(len(scenes)+len(galleries)+i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if gallery == nil {
				logger.Errorf("nil gallery, skipping Clean")
				continue
			}

			wg.Add(1)

			task := CleanTask{
				TxnManager: s.TxnManager,
				Gallery:    gallery,
			}
			go task.Start(&wg, input.DryRun)
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

	go func() {
		defer s.returnToIdleState()

		fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()
		logger.Infof("Migrating generated files for %s naming hash", fileNamingAlgo.String())

		var scenes []*models.Scene
		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			var err error
			scenes, err = r.Scene().All()
			return err
		}); err != nil {
			logger.Errorf("failed to fetch list of scenes for migration: %s", err.Error())
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

type totalsGenerate struct {
	sprites       int64
	previews      int64
	imagePreviews int64
	markers       int64
	transcodes    int64
	phashes       int64
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

	fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()
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
					TxnManager:          s.TxnManager,
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

			if input.Phashes {
				task := GeneratePhashTask{
					Scene:               *scene,
					fileNamingAlgorithm: fileNamingAlgo,
				}

				if task.shouldGenerate() {
					totals.phashes++
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
