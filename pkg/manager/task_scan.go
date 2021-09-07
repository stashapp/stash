package manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/remeh/sizedwaitgroup"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type ScanJob struct {
	txnManager    models.TransactionManager
	input         models.ScanMetadataInput
	subscriptions *subscriptionManager
}

func (j *ScanJob) Execute(ctx context.Context, progress *job.Progress) {
	input := j.input
	paths := getScanPaths(input.Paths)

	var total *int
	var newFiles *int
	progress.ExecuteTask("Counting files to scan...", func() {
		total, newFiles = j.neededScan(ctx, paths)
	})

	if job.IsCancelled(ctx) {
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

	if total != nil {
		progress.SetTotal(*total)
	}

	fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
	calculateMD5 := config.IsCalculateMD5()

	stoppingErr := errors.New("stopping")
	var err error

	var galleries []string

	for _, sp := range paths {
		csFs, er := utils.IsFsPathCaseSensitive(sp.Path)
		if er != nil {
			logger.Warnf("Cannot determine fs case sensitivity: %s", er.Error())
		}

		err = walkFilesToScan(sp, func(path string, info os.FileInfo, err error) error {
			if job.IsCancelled(ctx) {
				return stoppingErr
			}

			if isGallery(path) {
				galleries = append(galleries, path)
			}

			instance.Paths.Generated.EnsureTmpDir()

			wg.Add()
			task := ScanTask{
				TxnManager:           j.txnManager,
				FilePath:             path,
				UseFileMetadata:      utils.IsTrue(input.UseFileMetadata),
				StripFileExtension:   utils.IsTrue(input.StripFileExtension),
				fileNamingAlgorithm:  fileNamingAlgo,
				calculateMD5:         calculateMD5,
				GeneratePreview:      utils.IsTrue(input.ScanGeneratePreviews),
				GenerateImagePreview: utils.IsTrue(input.ScanGenerateImagePreviews),
				GenerateSprite:       utils.IsTrue(input.ScanGenerateSprites),
				GeneratePhash:        utils.IsTrue(input.ScanGeneratePhashes),
				progress:             progress,
				CaseSensitiveFs:      csFs,
				ctx:                  ctx,
			}

			go func() {
				task.Start(&wg)
				progress.Increment()
			}()

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

	if job.IsCancelled(ctx) || err != nil {
		return
	}

	progress.ExecuteTask("Associating galleries", func() {
		for _, path := range galleries {
			wg.Add()
			task := ScanTask{
				TxnManager:      j.txnManager,
				FilePath:        path,
				UseFileMetadata: false,
			}

			go task.associateGallery(&wg)
			wg.Wait()
		}
		logger.Info("Finished gallery association")
	})

	j.subscriptions.notify()
}

func (j *ScanJob) neededScan(ctx context.Context, paths []*models.StashConfig) (total *int, newFiles *int) {
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
			task := ScanTask{FilePath: path, TxnManager: j.txnManager}
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
			if job.IsCancelled(ctx) {
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

type ScanTask struct {
	ctx                  context.Context
	TxnManager           models.TransactionManager
	FilePath             string
	UseFileMetadata      bool
	StripFileExtension   bool
	calculateMD5         bool
	fileNamingAlgorithm  models.HashAlgorithm
	GenerateSprite       bool
	GeneratePhash        bool
	GeneratePreview      bool
	GenerateImagePreview bool
	zipGallery           *models.Gallery
	progress             *job.Progress
	CaseSensitiveFs      bool
}

func (t *ScanTask) Start(wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()

	var s *models.Scene

	t.progress.ExecuteTask("Scanning "+t.FilePath, func() {
		if isGallery(t.FilePath) {
			t.scanGallery()
		} else if isVideo(t.FilePath) {
			s = t.scanScene()
		} else if isImage(t.FilePath) {
			t.scanImage()
		}
	})

	if s != nil {
		iwg := sizedwaitgroup.New(2)

		if t.GenerateSprite {
			iwg.Add()

			go t.progress.ExecuteTask(fmt.Sprintf("Generating sprites for %s", t.FilePath), func() {
				taskSprite := GenerateSpriteTask{
					Scene:               *s,
					Overwrite:           false,
					fileNamingAlgorithm: t.fileNamingAlgorithm,
				}
				taskSprite.Start(&iwg)
			})
		}

		if t.GeneratePhash {
			iwg.Add()

			go t.progress.ExecuteTask(fmt.Sprintf("Generating phash for %s", t.FilePath), func() {
				taskPhash := GeneratePhashTask{
					Scene:               *s,
					fileNamingAlgorithm: t.fileNamingAlgorithm,
					txnManager:          t.TxnManager,
				}
				taskPhash.Start(&iwg)
			})
		}

		if t.GeneratePreview {
			iwg.Add()

			go t.progress.ExecuteTask(fmt.Sprintf("Generating preview for %s", t.FilePath), func() {
				config := config.GetInstance()
				var previewSegmentDuration = config.GetPreviewSegmentDuration()
				var previewSegments = config.GetPreviewSegments()
				var previewExcludeStart = config.GetPreviewExcludeStart()
				var previewExcludeEnd = config.GetPreviewExcludeEnd()
				var previewPresent = config.GetPreviewPreset()

				// NOTE: the reuse of this model like this is painful.
				previewOptions := models.GeneratePreviewOptionsInput{
					PreviewSegments:        &previewSegments,
					PreviewSegmentDuration: &previewSegmentDuration,
					PreviewExcludeStart:    &previewExcludeStart,
					PreviewExcludeEnd:      &previewExcludeEnd,
					PreviewPreset:          &previewPresent,
				}

				taskPreview := GeneratePreviewTask{
					Scene:               *s,
					ImagePreview:        t.GenerateImagePreview,
					Options:             previewOptions,
					Overwrite:           false,
					fileNamingAlgorithm: t.fileNamingAlgorithm,
				}
				taskPreview.Start(wg)
			})
		}

		iwg.Wait()
	}
}

func (t *ScanTask) getFileModTime() (time.Time, error) {
	fi, err := os.Stat(t.FilePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("error performing stat on %s: %s", t.FilePath, err.Error())
	}

	ret := fi.ModTime()
	// truncate to seconds, since we don't store beyond that in the database
	ret = ret.Truncate(time.Second)

	return ret, nil
}

func (t *ScanTask) isFileModified(fileModTime time.Time, modTime models.NullSQLiteTimestamp) bool {
	return !modTime.Timestamp.Equal(fileModTime)
}

func (t *ScanTask) calculateChecksum() (string, error) {
	logger.Infof("Calculating checksum for %s...", t.FilePath)
	checksum, err := utils.MD5FromFilePath(t.FilePath)
	if err != nil {
		return "", err
	}
	logger.Debugf("Checksum calculated: %s", checksum)
	return checksum, nil
}

func (t *ScanTask) calculateImageChecksum() (string, error) {
	logger.Infof("Calculating checksum for %s...", image.PathDisplayName(t.FilePath))
	// uses image.CalculateMD5 to read files in zips
	checksum, err := image.CalculateMD5(t.FilePath)
	if err != nil {
		return "", err
	}
	logger.Debugf("Checksum calculated: %s", checksum)
	return checksum, nil
}

func (t *ScanTask) doesPathExist() bool {
	config := config.GetInstance()
	vidExt := config.GetVideoExtensions()
	imgExt := config.GetImageExtensions()
	gExt := config.GetGalleryExtensions()

	ret := false
	t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		if matchExtension(t.FilePath, gExt) {
			gallery, _ := r.Gallery().FindByPath(t.FilePath)
			if gallery != nil {
				ret = true
			}
		} else if matchExtension(t.FilePath, vidExt) {
			s, _ := r.Scene().FindByPath(t.FilePath)
			if s != nil {
				ret = true
			}
		} else if matchExtension(t.FilePath, imgExt) {
			i, _ := r.Image().FindByPath(t.FilePath)
			if i != nil {
				ret = true
			}
		}

		return nil
	})

	return ret
}

func walkFilesToScan(s *models.StashConfig, f filepath.WalkFunc) error {
	config := config.GetInstance()
	vidExt := config.GetVideoExtensions()
	imgExt := config.GetImageExtensions()
	gExt := config.GetGalleryExtensions()
	excludeVidRegex := generateRegexps(config.GetExcludes())
	excludeImgRegex := generateRegexps(config.GetImageExcludes())

	// don't scan zip images directly
	if image.IsZipPath(s.Path) {
		logger.Warnf("Cannot rescan zip image %s. Rescan zip gallery instead.", s.Path)
		return nil
	}

	generatedPath := config.GetGeneratedPath()

	return utils.SymWalk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warnf("error scanning %s: %s", path, err.Error())
			return nil
		}

		if info.IsDir() {
			// #1102 - ignore files in generated path
			if utils.IsPathInDir(generatedPath, path) {
				return filepath.SkipDir
			}

			// shortcut: skip the directory entirely if it matches both exclusion patterns
			// add a trailing separator so that it correctly matches against patterns like path/.*
			pathExcludeTest := path + string(filepath.Separator)
			if (s.ExcludeVideo || matchFileRegex(pathExcludeTest, excludeVidRegex)) && (s.ExcludeImage || matchFileRegex(pathExcludeTest, excludeImgRegex)) {
				return filepath.SkipDir
			}

			return nil
		}

		if !s.ExcludeVideo && matchExtension(path, vidExt) && !matchFileRegex(path, excludeVidRegex) {
			return f(path, info, err)
		}

		if !s.ExcludeImage {
			if (matchExtension(path, imgExt) || matchExtension(path, gExt)) && !matchFileRegex(path, excludeImgRegex) {
				return f(path, info, err)
			}
		}

		return nil
	})
}
