package manager

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scene/generate"
	"github.com/stashapp/stash/pkg/txn"
)

type scanner interface {
	Scan(ctx context.Context, handlers []file.Handler, options file.ScanOptions, progressReporter file.ProgressReporter)
}

type ScanJob struct {
	scanner       scanner
	input         ScanMetadataInput
	subscriptions *subscriptionManager
}

func (j *ScanJob) Execute(ctx context.Context, progress *job.Progress) {
	input := j.input

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	sp := getScanPaths(input.Paths)
	paths := make([]string, len(sp))
	for i, p := range sp {
		paths[i] = p.Path
	}

	start := time.Now()

	const taskQueueSize = 200000
	taskQueue := job.NewTaskQueue(ctx, progress, taskQueueSize, instance.Config.GetParallelTasksWithAutoDetection())

	var minModTime time.Time
	if j.input.Filter != nil && j.input.Filter.MinModTime != nil {
		minModTime = *j.input.Filter.MinModTime
	}

	j.scanner.Scan(ctx, getScanHandlers(j.input, taskQueue, progress), file.ScanOptions{
		Paths:             paths,
		ScanFilters:       []file.PathFilter{newScanFilter(instance.Config, minModTime)},
		ZipFileExtensions: instance.Config.GetGalleryExtensions(),
		ParallelTasks:     instance.Config.GetParallelTasksWithAutoDetection(),
		HandlerRequiredFilters: []file.Filter{
			newHandlerRequiredFilter(instance.Config),
		},
	}, progress)

	taskQueue.Close()

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Scan finished (%s)", elapsed))

	j.subscriptions.notify()
}

type extensionConfig struct {
	vidExt []string
	imgExt []string
	zipExt []string
}

func newExtensionConfig(c *config.Instance) extensionConfig {
	return extensionConfig{
		vidExt: c.GetVideoExtensions(),
		imgExt: c.GetImageExtensions(),
		zipExt: c.GetGalleryExtensions(),
	}
}

type fileCounter interface {
	CountByFileID(ctx context.Context, fileID models.FileID) (int, error)
}

type galleryFinder interface {
	fileCounter
	FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Gallery, error)
}

type sceneFinder interface {
	fileCounter
	FindByPrimaryFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error)
}

// handlerRequiredFilter returns true if a File's handler needs to be executed despite the file not being updated.
type handlerRequiredFilter struct {
	extensionConfig
	txnManager     txn.Manager
	SceneFinder    sceneFinder
	ImageFinder    fileCounter
	GalleryFinder  galleryFinder
	CaptionUpdater video.CaptionUpdater

	FolderCache *lru.LRU

	videoFileNamingAlgorithm models.HashAlgorithm
}

func newHandlerRequiredFilter(c *config.Instance) *handlerRequiredFilter {
	db := instance.Database
	processes := c.GetParallelTasksWithAutoDetection()

	return &handlerRequiredFilter{
		extensionConfig:          newExtensionConfig(c),
		txnManager:               db,
		SceneFinder:              db.Scene,
		ImageFinder:              db.Image,
		GalleryFinder:            db.Gallery,
		CaptionUpdater:           db.File,
		FolderCache:              lru.New(processes * 2),
		videoFileNamingAlgorithm: c.GetVideoFileNamingAlgorithm(),
	}
}

func (f *handlerRequiredFilter) Accept(ctx context.Context, ff models.File) bool {
	path := ff.Base().Path
	isVideoFile := useAsVideo(path)
	isImageFile := useAsImage(path)
	isZipFile := fsutil.MatchExtension(path, f.zipExt)

	var counter fileCounter

	switch {
	case isVideoFile:
		// return true if there are no scenes associated
		counter = f.SceneFinder
	case isImageFile:
		counter = f.ImageFinder
	case isZipFile:
		counter = f.GalleryFinder
	}

	if counter == nil {
		return false
	}

	n, err := counter.CountByFileID(ctx, ff.Base().ID)
	if err != nil {
		// just ignore
		return false
	}

	// execute handler if there are no related objects
	if n == 0 {
		return true
	}

	// if create galleries from folder is enabled and the file is not in a zip
	// file, then check if there is a folder-based gallery for the file's
	// directory
	if isImageFile && instance.Config.GetCreateGalleriesFromFolders() && ff.Base().ZipFileID == nil {
		// only do this for the first time it encounters the folder
		// the first instance should create the gallery
		_, found := f.FolderCache.Get(ctx, ff.Base().ParentFolderID.String())
		if found {
			// should already be handled
			return false
		}

		g, _ := f.GalleryFinder.FindByFolderID(ctx, ff.Base().ParentFolderID)
		f.FolderCache.Add(ctx, ff.Base().ParentFolderID.String(), true)

		if len(g) == 0 {
			// no folder gallery. Return true so that it creates one.
			return true
		}
	}

	if isVideoFile {
		// TODO - check if the cover exists
		// hash := scene.GetHash(ff, f.videoFileNamingAlgorithm)
		// ssPath := instance.Paths.Scene.GetScreenshotPath(hash)
		// if exists, _ := fsutil.FileExists(ssPath); !exists {
		// 	// if not, check if the file is a primary file for a scene
		// 	scenes, err := f.SceneFinder.FindByPrimaryFileID(ctx, ff.Base().ID)
		// 	if err != nil {
		// 		// just ignore
		// 		return false
		// 	}

		// 	if len(scenes) > 0 {
		// 		// if it is, then it needs to be re-generated
		// 		return true
		// 	}
		// }

		// clean captions - scene handler handles this as well, but
		// unchanged files aren't processed by the scene handler
		videoFile, _ := ff.(*models.VideoFile)
		if videoFile != nil {
			if err := video.CleanCaptions(ctx, videoFile, f.txnManager, f.CaptionUpdater); err != nil {
				logger.Errorf("Error cleaning captions: %v", err)
			}
		}
	}

	return false
}

type scanFilter struct {
	extensionConfig
	stashPaths        config.StashConfigs
	generatedPath     string
	videoExcludeRegex []*regexp.Regexp
	imageExcludeRegex []*regexp.Regexp
	minModTime        time.Time
}

func newScanFilter(c *config.Instance, minModTime time.Time) *scanFilter {
	return &scanFilter{
		extensionConfig:   newExtensionConfig(c),
		stashPaths:        c.GetStashPaths(),
		generatedPath:     c.GetGeneratedPath(),
		videoExcludeRegex: generateRegexps(c.GetExcludes()),
		imageExcludeRegex: generateRegexps(c.GetImageExcludes()),
		minModTime:        minModTime,
	}
}

func (f *scanFilter) Accept(ctx context.Context, path string, info fs.FileInfo) bool {
	if fsutil.IsPathInDir(f.generatedPath, path) {
		logger.Warnf("Skipping %q as it overlaps with the generated folder", path)
		return false
	}

	// exit early on cutoff
	if info.Mode().IsRegular() && info.ModTime().Before(f.minModTime) {
		return false
	}

	isVideoFile := useAsVideo(path)
	isImageFile := useAsImage(path)
	isZipFile := fsutil.MatchExtension(path, f.zipExt)

	// handle caption files
	if fsutil.MatchExtension(path, video.CaptionExts) {
		// we don't include caption files in the file scan, but we do need
		// to handle them
		video.AssociateCaptions(ctx, path, instance.Repository, instance.Database.File, instance.Database.File)

		return false
	}

	if !info.IsDir() && !isVideoFile && !isImageFile && !isZipFile {
		logger.Debugf("Skipping %s as it does not match any known file extensions", path)
		return false
	}

	// #1756 - skip zero length files
	if !info.IsDir() && info.Size() == 0 {
		logger.Infof("Skipping zero-length file: %s", path)
		return false
	}

	s := f.stashPaths.GetStashFromDirPath(path)

	if s == nil {
		logger.Debugf("Skipping %s as it is not in the stash library", path)
		return false
	}

	// shortcut: skip the directory entirely if it matches both exclusion patterns
	// add a trailing separator so that it correctly matches against patterns like path/.*
	pathExcludeTest := path + string(filepath.Separator)
	if (matchFileRegex(pathExcludeTest, f.videoExcludeRegex)) && (s.ExcludeImage || matchFileRegex(pathExcludeTest, f.imageExcludeRegex)) {
		logger.Debugf("Skipping directory %s as it matches video and image exclusion patterns", path)
		return false
	}

	if isVideoFile && (s.ExcludeVideo || matchFileRegex(path, f.videoExcludeRegex)) {
		logger.Debugf("Skipping %s as it matches video exclusion patterns", path)
		return false
	} else if (isImageFile || isZipFile) && (s.ExcludeImage || matchFileRegex(path, f.imageExcludeRegex)) {
		logger.Debugf("Skipping %s as it matches image exclusion patterns", path)
		return false
	}

	return true
}

type scanConfig struct {
	isGenerateThumbnails   bool
	isGenerateClipPreviews bool
}

func (c *scanConfig) GetCreateGalleriesFromFolders() bool {
	return instance.Config.GetCreateGalleriesFromFolders()
}

func getScanHandlers(options ScanMetadataInput, taskQueue *job.TaskQueue, progress *job.Progress) []file.Handler {
	db := instance.Database
	pluginCache := instance.PluginCache

	return []file.Handler{
		&file.FilteredHandler{
			Filter: file.FilterFunc(imageFileFilter),
			Handler: &image.ScanHandler{
				CreatorUpdater: db.Image,
				GalleryFinder:  db.Gallery,
				ScanGenerator: &imageGenerators{
					input:     options,
					taskQueue: taskQueue,
					progress:  progress,
				},
				ScanConfig: &scanConfig{
					isGenerateThumbnails:   options.ScanGenerateThumbnails,
					isGenerateClipPreviews: options.ScanGenerateClipPreviews,
				},
				PluginCache: pluginCache,
				Paths:       instance.Paths,
			},
		},
		&file.FilteredHandler{
			Filter: file.FilterFunc(galleryFileFilter),
			Handler: &gallery.ScanHandler{
				CreatorUpdater:     db.Gallery,
				SceneFinderUpdater: db.Scene,
				ImageFinderUpdater: db.Image,
				PluginCache:        pluginCache,
			},
		},
		&file.FilteredHandler{
			Filter: file.FilterFunc(videoFileFilter),
			Handler: &scene.ScanHandler{
				CreatorUpdater: db.Scene,
				PluginCache:    pluginCache,
				CaptionUpdater: db.File,
				ScanGenerator: &sceneGenerators{
					input:     options,
					taskQueue: taskQueue,
					progress:  progress,
				},
				FileNamingAlgorithm: instance.Config.GetVideoFileNamingAlgorithm(),
				Paths:               instance.Paths,
			},
		},
	}
}

type imageGenerators struct {
	input     ScanMetadataInput
	taskQueue *job.TaskQueue
	progress  *job.Progress
}

func (g *imageGenerators) Generate(ctx context.Context, i *models.Image, f models.File) error {
	const overwrite = false

	progress := g.progress
	t := g.input
	path := f.Base().Path
	config := instance.Config
	sequentialScanning := config.GetSequentialScanning()

	if t.ScanGenerateThumbnails {
		// this should be quick, so always generate sequentially
		if err := g.generateThumbnail(ctx, i, f); err != nil {
			logger.Errorf("Error generating thumbnail for %s: %v", path, err)
		}
	}

	// avoid adding a task if the file isn't a video file
	_, isVideo := f.(*models.VideoFile)
	if isVideo && t.ScanGenerateClipPreviews {
		// this is a bit of a hack: the task requires files to be loaded, but
		// we don't really need to since we already have the file
		ii := *i
		ii.Files = models.NewRelatedFiles([]models.File{f})

		progress.AddTotal(1)
		previewsFn := func(ctx context.Context) {
			taskPreview := GenerateClipPreviewTask{
				Image:     ii,
				Overwrite: overwrite,
			}

			taskPreview.Start(ctx)
			progress.Increment()
		}

		if sequentialScanning {
			previewsFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating preview for %s", path), previewsFn)
		}
	}

	return nil
}

func (g *imageGenerators) generateThumbnail(ctx context.Context, i *models.Image, f models.File) error {
	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(i.Checksum, models.DefaultGthumbWidth)
	exists, _ := fsutil.FileExists(thumbPath)
	if exists {
		return nil
	}

	path := f.Base().Path

	vf, ok := f.(models.VisualFile)
	if !ok {
		return fmt.Errorf("file %s is not a visual file", path)
	}

	if vf.GetHeight() <= models.DefaultGthumbWidth && vf.GetWidth() <= models.DefaultGthumbWidth {
		return nil
	}

	logger.Debugf("Generating thumbnail for %s", path)

	clipPreviewOptions := image.ClipPreviewOptions{
		InputArgs:  instance.Config.GetTranscodeInputArgs(),
		OutputArgs: instance.Config.GetTranscodeOutputArgs(),
		Preset:     instance.Config.GetPreviewPreset().String(),
	}

	encoder := image.NewThumbnailEncoder(instance.FFMPEG, instance.FFProbe, clipPreviewOptions)
	data, err := encoder.GetThumbnail(f, models.DefaultGthumbWidth)

	if err != nil {
		// don't log for animated images
		if !errors.Is(err, image.ErrNotSupportedForThumbnail) {
			return fmt.Errorf("getting thumbnail for image %s: %w", path, err)
		}
		return nil
	}

	err = fsutil.WriteFile(thumbPath, data)
	if err != nil {
		return fmt.Errorf("writing thumbnail for image %s: %w", path, err)
	}

	return nil
}

type sceneGenerators struct {
	input     ScanMetadataInput
	taskQueue *job.TaskQueue
	progress  *job.Progress
}

func (g *sceneGenerators) Generate(ctx context.Context, s *models.Scene, f *models.VideoFile) error {
	const overwrite = false

	progress := g.progress
	t := g.input
	path := f.Path
	config := instance.Config
	fileNamingAlgorithm := config.GetVideoFileNamingAlgorithm()
	sequentialScanning := config.GetSequentialScanning()

	if t.ScanGenerateSprites {
		progress.AddTotal(1)
		spriteFn := func(ctx context.Context) {
			taskSprite := GenerateSpriteTask{
				Scene:               *s,
				Overwrite:           overwrite,
				fileNamingAlgorithm: fileNamingAlgorithm,
			}
			taskSprite.Start(ctx)
			progress.Increment()
		}

		if sequentialScanning {
			spriteFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating sprites for %s", path), spriteFn)
		}
	}

	if t.ScanGeneratePhashes {
		progress.AddTotal(1)
		phashFn := func(ctx context.Context) {
			taskPhash := GeneratePhashTask{
				File:                f,
				fileNamingAlgorithm: fileNamingAlgorithm,
				txnManager:          instance.Database,
				fileUpdater:         instance.Database.File,
				Overwrite:           overwrite,
			}
			taskPhash.Start(ctx)
			progress.Increment()
		}

		if sequentialScanning {
			phashFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating phash for %s", path), phashFn)
		}
	}

	if t.ScanGeneratePreviews {
		progress.AddTotal(1)
		previewsFn := func(ctx context.Context) {
			options := getGeneratePreviewOptions(GeneratePreviewOptionsInput{})

			g := &generate.Generator{
				Encoder:      instance.FFMPEG,
				FFMpegConfig: instance.Config,
				LockManager:  instance.ReadLockManager,
				MarkerPaths:  instance.Paths.SceneMarkers,
				ScenePaths:   instance.Paths.Scene,
				Overwrite:    overwrite,
			}

			taskPreview := GeneratePreviewTask{
				Scene:               *s,
				ImagePreview:        t.ScanGenerateImagePreviews,
				Options:             options,
				Overwrite:           overwrite,
				fileNamingAlgorithm: fileNamingAlgorithm,
				generator:           g,
			}
			taskPreview.Start(ctx)
			progress.Increment()
		}

		if sequentialScanning {
			previewsFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating preview for %s", path), previewsFn)
		}
	}

	if t.ScanGenerateCovers {
		progress.AddTotal(1)
		g.taskQueue.Add(fmt.Sprintf("Generating cover for %s", path), func(ctx context.Context) {
			taskCover := GenerateCoverTask{
				Scene:      *s,
				txnManager: instance.Repository,
				Overwrite:  overwrite,
			}
			taskCover.Start(ctx)
			progress.Increment()
		})
	}

	return nil
}
