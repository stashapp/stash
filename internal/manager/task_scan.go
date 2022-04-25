package manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/remeh/sizedwaitgroup"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scene/generate"
	"github.com/stashapp/stash/pkg/utils"
)

const scanQueueSize = 200000

type ScanJob struct {
	txnManager    models.TransactionManager
	input         ScanMetadataInput
	subscriptions *subscriptionManager
}

type scanFile struct {
	path            string
	info            os.FileInfo
	caseSensitiveFs bool
}

func (j *ScanJob) Execute(ctx context.Context, progress *job.Progress) {
	input := j.input
	paths := getScanPaths(input.Paths)

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	start := time.Now()
	config := config.GetInstance()
	parallelTasks := config.GetParallelTasksWithAutoDetection()

	logger.Infof("Scan started with %d parallel tasks", parallelTasks)

	fileQueue := make(chan scanFile, scanQueueSize)
	go func() {
		total, newFiles := j.queueFiles(ctx, paths, fileQueue, parallelTasks)

		if !job.IsCancelled(ctx) {
			progress.SetTotal(total)
			logger.Infof("Finished counting files. Total files to scan: %d, %d new files found", total, newFiles)
		}
	}()

	wg := sizedwaitgroup.New(parallelTasks)

	fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
	calculateMD5 := config.IsCalculateMD5()

	var err error

	var galleries []string

	mutexManager := utils.NewMutexManager()

	for f := range fileQueue {
		if job.IsCancelled(ctx) {
			break
		}

		if isGallery(f.path) {
			galleries = append(galleries, f.path)
		}

		if err := instance.Paths.Generated.EnsureTmpDir(); err != nil {
			logger.Warnf("couldn't create temporary directory: %v", err)
		}

		wg.Add()
		task := ScanTask{
			TxnManager:           j.txnManager,
			file:                 file.FSFile(f.path, f.info),
			UseFileMetadata:      input.UseFileMetadata,
			StripFileExtension:   input.StripFileExtension,
			fileNamingAlgorithm:  fileNamingAlgo,
			calculateMD5:         calculateMD5,
			GeneratePreview:      input.ScanGeneratePreviews,
			GenerateImagePreview: input.ScanGenerateImagePreviews,
			GenerateSprite:       input.ScanGenerateSprites,
			GeneratePhash:        input.ScanGeneratePhashes,
			GenerateThumbnails:   input.ScanGenerateThumbnails,
			progress:             progress,
			CaseSensitiveFs:      f.caseSensitiveFs,
			mutexManager:         mutexManager,
		}

		go func() {
			task.Start(ctx)
			wg.Done()
			progress.Increment()
		}()
	}

	wg.Wait()

	if err := instance.Paths.Generated.EmptyTmpDir(); err != nil {
		logger.Warnf("couldn't empty temporary directory: %v", err)
	}

	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Scan finished (%s)", elapsed))

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	if err != nil {
		return
	}

	progress.ExecuteTask("Associating galleries", func() {
		for _, path := range galleries {
			wg.Add()
			task := ScanTask{
				TxnManager:      j.txnManager,
				file:            file.FSFile(path, nil), // hopefully info is not needed
				UseFileMetadata: false,
			}

			go task.associateGallery(ctx, &wg)
			wg.Wait()
		}
		logger.Info("Finished gallery association")
	})

	j.subscriptions.notify()
}

func (j *ScanJob) queueFiles(ctx context.Context, paths []*config.StashConfig, scanQueue chan<- scanFile, parallelTasks int) (total int, newFiles int) {
	defer close(scanQueue)

	var minModTime time.Time
	if j.input.Filter != nil && j.input.Filter.MinModTime != nil {
		minModTime = *j.input.Filter.MinModTime
	}

	wg := sizedwaitgroup.New(parallelTasks)

	for _, sp := range paths {
		csFs, er := fsutil.IsFsPathCaseSensitive(sp.Path)
		if er != nil {
			logger.Warnf("Cannot determine fs case sensitivity: %s", er.Error())
		}

		err := walkFilesToScan(sp, func(path string, info os.FileInfo, err error) error {
			// check stop
			if job.IsCancelled(ctx) {
				return context.Canceled
			}

			// exit early on cutoff
			if info.Mode().IsRegular() && info.ModTime().Before(minModTime) {
				return nil
			}

			wg.Add()

			go func() {
				defer wg.Done()

				// #1756 - skip zero length files and directories
				if info.IsDir() {
					return
				}

				if info.Size() == 0 {
					logger.Infof("Skipping zero-length file: %s", path)
					return
				}

				total++
				if !j.doesPathExist(ctx, path) {
					newFiles++
				}

				scanQueue <- scanFile{
					path:            path,
					info:            info,
					caseSensitiveFs: csFs,
				}
			}()

			return nil
		})

		wg.Wait()

		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf("Error encountered queuing files to scan: %s", err.Error())
			return
		}
	}

	return
}

func (j *ScanJob) doesPathExist(ctx context.Context, path string) bool {
	config := config.GetInstance()
	vidExt := config.GetVideoExtensions()
	imgExt := config.GetImageExtensions()
	gExt := config.GetGalleryExtensions()

	ret := false
	txnErr := j.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		switch {
		case fsutil.MatchExtension(path, gExt):
			g, _ := r.Gallery().FindByPath(path)
			if g != nil {
				ret = true
			}
		case fsutil.MatchExtension(path, vidExt):
			s, _ := r.Scene().FindByPath(path)
			if s != nil {
				ret = true
			}
		case fsutil.MatchExtension(path, imgExt):
			i, _ := r.Image().FindByPath(path)
			if i != nil {
				ret = true
			}
		}

		return nil
	})
	if txnErr != nil {
		logger.Warnf("error checking if file exists in database: %v", txnErr)
	}

	return ret
}

type ScanTask struct {
	TxnManager           models.TransactionManager
	file                 file.SourceFile
	UseFileMetadata      bool
	StripFileExtension   bool
	calculateMD5         bool
	fileNamingAlgorithm  models.HashAlgorithm
	GenerateSprite       bool
	GeneratePhash        bool
	GeneratePreview      bool
	GenerateImagePreview bool
	GenerateThumbnails   bool
	zipGallery           *models.Gallery
	progress             *job.Progress
	CaseSensitiveFs      bool

	mutexManager *utils.MutexManager
}

func (t *ScanTask) Start(ctx context.Context) {
	var s *models.Scene
	path := t.file.Path()
	t.progress.ExecuteTask("Scanning "+path, func() {
		switch {
		case isGallery(path):
			t.scanGallery(ctx)
		case isVideo(path):
			s = t.scanScene(ctx)
		case isImage(path):
			t.scanImage(ctx)
		case isCaptions(path):
			t.associateCaptions(ctx)
		}
	})

	if s == nil {
		return
	}

	// Handle the case of a scene
	iwg := sizedwaitgroup.New(2)

	if t.GenerateSprite {
		iwg.Add()

		go t.progress.ExecuteTask(fmt.Sprintf("Generating sprites for %s", path), func() {
			taskSprite := GenerateSpriteTask{
				Scene:               *s,
				Overwrite:           false,
				fileNamingAlgorithm: t.fileNamingAlgorithm,
			}
			taskSprite.Start(ctx)
			iwg.Done()
		})
	}

	if t.GeneratePhash {
		iwg.Add()

		go t.progress.ExecuteTask(fmt.Sprintf("Generating phash for %s", path), func() {
			taskPhash := GeneratePhashTask{
				Scene:               *s,
				fileNamingAlgorithm: t.fileNamingAlgorithm,
				txnManager:          t.TxnManager,
			}
			taskPhash.Start(ctx)
			iwg.Done()
		})
	}

	if t.GeneratePreview {
		iwg.Add()

		go t.progress.ExecuteTask(fmt.Sprintf("Generating preview for %s", path), func() {
			options := getGeneratePreviewOptions(GeneratePreviewOptionsInput{})
			const overwrite = false

			g := &generate.Generator{
				Encoder:     instance.FFMPEG,
				LockManager: instance.ReadLockManager,
				MarkerPaths: instance.Paths.SceneMarkers,
				ScenePaths:  instance.Paths.Scene,
				Overwrite:   overwrite,
			}

			taskPreview := GeneratePreviewTask{
				Scene:               *s,
				ImagePreview:        t.GenerateImagePreview,
				Options:             options,
				Overwrite:           overwrite,
				fileNamingAlgorithm: t.fileNamingAlgorithm,
				generator:           g,
			}
			taskPreview.Start(ctx)
			iwg.Done()
		})
	}

	iwg.Wait()
}

func walkFilesToScan(s *config.StashConfig, f filepath.WalkFunc) error {
	config := config.GetInstance()
	vidExt := config.GetVideoExtensions()
	imgExt := config.GetImageExtensions()
	gExt := config.GetGalleryExtensions()
	capExt := scene.CaptionExts
	excludeVidRegex := generateRegexps(config.GetExcludes())
	excludeImgRegex := generateRegexps(config.GetImageExcludes())

	// don't scan zip images directly
	if file.IsZipPath(s.Path) {
		logger.Warnf("Cannot rescan zip image %s. Rescan zip gallery instead.", s.Path)
		return nil
	}

	generatedPath := config.GetGeneratedPath()

	return fsutil.SymWalk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warnf("error scanning %s: %s", path, err.Error())
			return nil
		}

		if info.IsDir() {
			// #1102 - ignore files in generated path
			if fsutil.IsPathInDir(generatedPath, path) {
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

		if !s.ExcludeVideo && fsutil.MatchExtension(path, vidExt) && !matchFileRegex(path, excludeVidRegex) {
			return f(path, info, err)
		}

		if !s.ExcludeImage {
			if (fsutil.MatchExtension(path, imgExt) || fsutil.MatchExtension(path, gExt)) && !matchFileRegex(path, excludeImgRegex) {
				return f(path, info, err)
			}
		}

		if fsutil.MatchExtension(path, capExt) {
			return f(path, info, err)
		}

		return nil
	})
}
