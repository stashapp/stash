package manager

import (
	"context"
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
	"github.com/stashapp/stash/pkg/models/paths"
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

func (j *ScanJob) Execute(ctx context.Context, progress *job.Progress) error {
	cfg := config.GetInstance()
	input := j.input

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}

	sp := getScanPaths(input.Paths)
	paths := make([]string, len(sp))
	for i, p := range sp {
		paths[i] = p.Path
	}

	mgr := GetInstance()
	c := mgr.Config
	repo := mgr.Repository

	start := time.Now()

	const taskQueueSize = 200000
	taskQueue := job.NewTaskQueue(ctx, progress, taskQueueSize, cfg.GetParallelTasksWithAutoDetection())

	var minModTime time.Time
	if j.input.Filter != nil && j.input.Filter.MinModTime != nil {
		minModTime = *j.input.Filter.MinModTime
	}

	j.scanner.Scan(ctx, getScanHandlers(j.input, taskQueue, progress), file.ScanOptions{
		Paths:                  paths,
		ScanFilters:            []file.PathFilter{newScanFilter(c, repo, minModTime)},
		ZipFileExtensions:      cfg.GetGalleryExtensions(),
		ParallelTasks:          cfg.GetParallelTasksWithAutoDetection(),
		HandlerRequiredFilters: []file.Filter{newHandlerRequiredFilter(cfg, repo)},
		Rescan:                 j.input.Rescan,
	}, progress)

	taskQueue.Close()

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}

	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Scan finished (%s)", elapsed))

	j.subscriptions.notify()
	return nil
}

type extensionConfig struct {
	vidExt []string
	imgExt []string
	zipExt []string
}

func newExtensionConfig(c *config.Config) extensionConfig {
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

	FolderCache *lru.LRU[bool]

	videoFileNamingAlgorithm models.HashAlgorithm
}

func newHandlerRequiredFilter(c *config.Config, repo models.Repository) *handlerRequiredFilter {
	processes := c.GetParallelTasksWithAutoDetection()

	return &handlerRequiredFilter{
		extensionConfig:          newExtensionConfig(c),
		txnManager:               repo.TxnManager,
		SceneFinder:              repo.Scene,
		ImageFinder:              repo.Image,
		GalleryFinder:            repo.Gallery,
		CaptionUpdater:           repo.File,
		FolderCache:              lru.New[bool](processes * 2),
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
	// #4611 - also check for .forcegallery
	if isImageFile && ff.Base().ZipFileID == nil {
		// only do this for the first time it encounters the folder
		// the first instance should create the gallery
		_, found := f.FolderCache.Get(ctx, ff.Base().ParentFolderID.String())
		if found {
			// should already be handled
			return false
		}

		f.FolderCache.Add(ctx, ff.Base().ParentFolderID.String(), true)

		createGallery := instance.Config.GetCreateGalleriesFromFolders()
		if !createGallery {
			// check for presence of .forcegallery
			forceGalleryPath := filepath.Join(filepath.Dir(path), ".forcegallery")
			if exists, _ := fsutil.FileExists(forceGalleryPath); exists {
				createGallery = true
			}
		}

		if !createGallery {
			return false
		}

		g, _ := f.GalleryFinder.FindByFolderID(ctx, ff.Base().ParentFolderID)

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
	txnManager     txn.Manager
	FileFinder     models.FileFinder
	CaptionUpdater video.CaptionUpdater

	stashPaths        config.StashConfigs
	generatedPath     string
	videoExcludeRegex []*regexp.Regexp
	imageExcludeRegex []*regexp.Regexp
	minModTime        time.Time
}

func newScanFilter(c *config.Config, repo models.Repository, minModTime time.Time) *scanFilter {
	return &scanFilter{
		extensionConfig:   newExtensionConfig(c),
		txnManager:        repo.TxnManager,
		FileFinder:        repo.File,
		CaptionUpdater:    repo.File,
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

	s := f.stashPaths.GetStashFromDirPath(path)
	if s == nil {
		logger.Debugf("Skipping %s as it is not in the stash library", path)
		return false
	}

	isVideoFile := useAsVideo(path)
	isImageFile := useAsImage(path)
	isZipFile := fsutil.MatchExtension(path, f.zipExt)

	// handle caption files
	if fsutil.MatchExtension(path, video.CaptionExts) {
		// we don't include caption files in the file scan, but we do need
		// to handle them
		video.AssociateCaptions(ctx, path, f.txnManager, f.FileFinder, f.CaptionUpdater)

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

	createGalleriesFromFolders bool
}

func (c *scanConfig) GetCreateGalleriesFromFolders() bool {
	return c.createGalleriesFromFolders
}

func videoFileFilter(ctx context.Context, f models.File) bool {
	return useAsVideo(f.Base().Path)
}

func imageFileFilter(ctx context.Context, f models.File) bool {
	return useAsImage(f.Base().Path)
}

func galleryFileFilter(ctx context.Context, f models.File) bool {
	return isZip(f.Base().Basename)
}

func getScanHandlers(options ScanMetadataInput, taskQueue *job.TaskQueue, progress *job.Progress) []file.Handler {
	mgr := GetInstance()
	c := mgr.Config
	r := mgr.Repository
	pluginCache := mgr.PluginCache

	return []file.Handler{
		&file.FilteredHandler{
			Filter: file.FilterFunc(imageFileFilter),
			Handler: &image.ScanHandler{
				CreatorUpdater: r.Image,
				GalleryFinder:  r.Gallery,
				ScanGenerator: &imageGenerators{
					input:              options,
					taskQueue:          taskQueue,
					progress:           progress,
					paths:              mgr.Paths,
					sequentialScanning: c.GetSequentialScanning(),
				},
				ScanConfig: &scanConfig{
					isGenerateThumbnails:       options.ScanGenerateThumbnails,
					isGenerateClipPreviews:     options.ScanGenerateClipPreviews,
					createGalleriesFromFolders: c.GetCreateGalleriesFromFolders(),
				},
				PluginCache: pluginCache,
				Paths:       instance.Paths,
			},
		},
		&file.FilteredHandler{
			Filter: file.FilterFunc(galleryFileFilter),
			Handler: &gallery.ScanHandler{
				CreatorUpdater:     r.Gallery,
				SceneFinderUpdater: r.Scene,
				ImageFinderUpdater: r.Image,
				PluginCache:        pluginCache,
			},
		},
		&file.FilteredHandler{
			Filter: file.FilterFunc(videoFileFilter),
			Handler: &scene.ScanHandler{
				CreatorUpdater: r.Scene,
				CaptionUpdater: r.File,
				PluginCache:    pluginCache,
				ScanGenerator: &sceneGenerators{
					input:               options,
					taskQueue:           taskQueue,
					progress:            progress,
					paths:               mgr.Paths,
					fileNamingAlgorithm: c.GetVideoFileNamingAlgorithm(),
					sequentialScanning:  c.GetSequentialScanning(),
				},
				FileNamingAlgorithm: c.GetVideoFileNamingAlgorithm(),
				Paths:               mgr.Paths,
			},
		},
	}
}

type imageGenerators struct {
	input     ScanMetadataInput
	taskQueue *job.TaskQueue
	progress  *job.Progress

	paths              *paths.Paths
	sequentialScanning bool
}

func (g *imageGenerators) Generate(ctx context.Context, i *models.Image, f models.File) error {
	const overwrite = false

	progress := g.progress
	t := g.input
	path := f.Base().Path

	// this is a bit of a hack: the task requires files to be loaded, but
	// we don't really need to since we already have the file
	ii := *i
	ii.Files = models.NewRelatedFiles([]models.File{f})

	if t.ScanGenerateThumbnails {
		// this should be quick, so always generate sequentially
		taskThumbnail := GenerateImageThumbnailTask{
			Image:     ii,
			Overwrite: overwrite,
		}

		taskThumbnail.Start(ctx)
	}

	// avoid adding a task if the file isn't a video file
	_, isVideo := f.(*models.VideoFile)
	if isVideo && t.ScanGenerateClipPreviews {
		progress.AddTotal(1)
		previewsFn := func(ctx context.Context) {
			taskPreview := GenerateClipPreviewTask{
				Image:     ii,
				Overwrite: overwrite,
			}

			taskPreview.Start(ctx)
			progress.Increment()
		}

		if g.sequentialScanning {
			previewsFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating preview for %s", path), previewsFn)
		}
	}

	return nil
}

type sceneGenerators struct {
	input     ScanMetadataInput
	taskQueue *job.TaskQueue
	progress  *job.Progress

	paths               *paths.Paths
	fileNamingAlgorithm models.HashAlgorithm
	sequentialScanning  bool
}

func (g *sceneGenerators) Generate(ctx context.Context, s *models.Scene, f *models.VideoFile) error {
	const overwrite = false

	progress := g.progress
	t := g.input
	path := f.Path

	mgr := GetInstance()

	if t.ScanGenerateSprites {
		progress.AddTotal(1)
		spriteFn := func(ctx context.Context) {
			taskSprite := GenerateSpriteTask{
				Scene:               *s,
				Overwrite:           overwrite,
				fileNamingAlgorithm: g.fileNamingAlgorithm,
			}
			taskSprite.Start(ctx)
			progress.Increment()
		}

		if g.sequentialScanning {
			spriteFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating sprites for %s", path), spriteFn)
		}
	}

	if t.ScanGeneratePhashes {
		progress.AddTotal(1)
		phashFn := func(ctx context.Context) {
			taskPhash := GeneratePhashTask{
				repository:          mgr.Repository,
				File:                f,
				Overwrite:           overwrite,
				fileNamingAlgorithm: g.fileNamingAlgorithm,
			}
			taskPhash.Start(ctx)
			progress.Increment()
		}

		if g.sequentialScanning {
			phashFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating phash for %s", path), phashFn)
		}
	}

	if t.ScanGeneratePreviews {
		progress.AddTotal(1)
		previewsFn := func(ctx context.Context) {
			options := getGeneratePreviewOptions(GeneratePreviewOptionsInput{})

			generator := &generate.Generator{
				Encoder:      mgr.FFMpeg,
				FFMpegConfig: mgr.Config,
				LockManager:  mgr.ReadLockManager,
				MarkerPaths:  g.paths.SceneMarkers,
				ScenePaths:   g.paths.Scene,
				Overwrite:    overwrite,
			}

			taskPreview := GeneratePreviewTask{
				Scene:               *s,
				ImagePreview:        t.ScanGenerateImagePreviews,
				Options:             options,
				Overwrite:           overwrite,
				fileNamingAlgorithm: g.fileNamingAlgorithm,
				generator:           generator,
			}
			taskPreview.Start(ctx)
			progress.Increment()
		}

		if g.sequentialScanning {
			previewsFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating preview for %s", path), previewsFn)
		}
	}

	if t.ScanGenerateCovers {
		progress.AddTotal(1)
		g.taskQueue.Add(fmt.Sprintf("Generating cover for %s", path), func(ctx context.Context) {
			taskCover := GenerateCoverTask{
				repository: mgr.Repository,
				Scene:      *s,
				Overwrite:  overwrite,
			}
			taskCover.Start(ctx)
			progress.Increment()
		})
	}

	return nil
}
