package manager

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"

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
		Paths:                  paths,
		ScanFilters:            []file.PathFilter{newScanFilter(instance.Config, minModTime)},
		ZipFileExtensions:      instance.Config.GetGalleryExtensions(),
		ParallelTasks:          instance.Config.GetParallelTasksWithAutoDetection(),
		HandlerRequiredFilters: []file.Filter{newHandlerRequiredFilter(instance.Config)},
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
	CountByFileID(ctx context.Context, fileID file.ID) (int, error)
}

// handlerRequiredFilter returns true if a File's handler needs to be executed despite the file not being updated.
type handlerRequiredFilter struct {
	extensionConfig
	SceneFinder   fileCounter
	ImageFinder   fileCounter
	GalleryFinder fileCounter
}

func newHandlerRequiredFilter(c *config.Instance) *handlerRequiredFilter {
	db := instance.Database

	return &handlerRequiredFilter{
		extensionConfig: newExtensionConfig(c),
		SceneFinder:     db.Scene,
		ImageFinder:     db.Image,
		GalleryFinder:   db.Gallery,
	}
}

func (f *handlerRequiredFilter) Accept(ctx context.Context, ff file.File) bool {
	path := ff.Base().Path
	isVideoFile := fsutil.MatchExtension(path, f.vidExt)
	isImageFile := fsutil.MatchExtension(path, f.imgExt)
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
	return n == 0
}

type scanFilter struct {
	extensionConfig
	stashPaths        []*config.StashConfig
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
		return false
	}

	// exit early on cutoff
	if info.Mode().IsRegular() && info.ModTime().Before(f.minModTime) {
		return false
	}

	isVideoFile := fsutil.MatchExtension(path, f.vidExt)
	isImageFile := fsutil.MatchExtension(path, f.imgExt)
	isZipFile := fsutil.MatchExtension(path, f.zipExt)

	// handle caption files
	if fsutil.MatchExtension(path, video.CaptionExts) {
		// we don't include caption files in the file scan, but we do need
		// to handle them
		video.AssociateCaptions(ctx, path, instance.Repository, instance.Database.File, instance.Database.File)

		return false
	}

	if !info.IsDir() && !isVideoFile && !isImageFile && !isZipFile {
		return false
	}

	// #1756 - skip zero length files
	if !info.IsDir() && info.Size() == 0 {
		logger.Infof("Skipping zero-length file: %s", path)
		return false
	}

	s := getStashFromDirPath(f.stashPaths, path)

	if s == nil {
		return false
	}

	// shortcut: skip the directory entirely if it matches both exclusion patterns
	// add a trailing separator so that it correctly matches against patterns like path/.*
	pathExcludeTest := path + string(filepath.Separator)
	if (s.ExcludeVideo || matchFileRegex(pathExcludeTest, f.videoExcludeRegex)) && (s.ExcludeImage || matchFileRegex(pathExcludeTest, f.imageExcludeRegex)) {
		return false
	}

	if isVideoFile && (s.ExcludeVideo || matchFileRegex(path, f.videoExcludeRegex)) {
		return false
	} else if (isImageFile || isZipFile) && s.ExcludeImage || matchFileRegex(path, f.imageExcludeRegex) {
		return false
	}

	return true
}

type scanConfig struct {
	isGenerateThumbnails bool
}

func (c *scanConfig) GetCreateGalleriesFromFolders() bool {
	return instance.Config.GetCreateGalleriesFromFolders()
}

func (c *scanConfig) IsGenerateThumbnails() bool {
	return c.isGenerateThumbnails
}

func getScanHandlers(options ScanMetadataInput, taskQueue *job.TaskQueue, progress *job.Progress) []file.Handler {
	db := instance.Database
	pluginCache := instance.PluginCache

	return []file.Handler{
		&file.FilteredHandler{
			Filter: file.FilterFunc(imageFileFilter),
			Handler: &image.ScanHandler{
				CreatorUpdater:     db.Image,
				GalleryFinder:      db.Gallery,
				ThumbnailGenerator: &imageThumbnailGenerator{},
				ScanConfig: &scanConfig{
					isGenerateThumbnails: options.ScanGenerateThumbnails,
				},
				PluginCache: pluginCache,
			},
		},
		&file.FilteredHandler{
			Filter: file.FilterFunc(galleryFileFilter),
			Handler: &gallery.ScanHandler{
				CreatorUpdater:     db.Gallery,
				SceneFinderUpdater: db.Scene,
				PluginCache:        pluginCache,
			},
		},
		&file.FilteredHandler{
			Filter: file.FilterFunc(videoFileFilter),
			Handler: &scene.ScanHandler{
				CreatorUpdater: db.Scene,
				PluginCache:    pluginCache,
				CoverGenerator: &coverGenerator{},
				ScanGenerator: &sceneGenerators{
					input:     options,
					taskQueue: taskQueue,
					progress:  progress,
				},
			},
		},
	}
}

type imageThumbnailGenerator struct{}

func (g *imageThumbnailGenerator) GenerateThumbnail(ctx context.Context, i *models.Image, f *file.ImageFile) error {
	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(i.Checksum(), models.DefaultGthumbWidth)
	exists, _ := fsutil.FileExists(thumbPath)
	if exists {
		return nil
	}

	if f.Height <= models.DefaultGthumbWidth && f.Width <= models.DefaultGthumbWidth {
		return nil
	}

	logger.Debugf("Generating thumbnail for %s", f.Path)

	encoder := image.NewThumbnailEncoder(instance.FFMPEG)
	data, err := encoder.GetThumbnail(f, models.DefaultGthumbWidth)

	if err != nil {
		// don't log for animated images
		if !errors.Is(err, image.ErrNotSupportedForThumbnail) {
			return fmt.Errorf("getting thumbnail for image %s: %w", f.Path, err)
		}
		return nil
	}

	err = fsutil.WriteFile(thumbPath, data)
	if err != nil {
		return fmt.Errorf("writing thumbnail for image %s: %w", f.Path, err)
	}

	return nil
}

type sceneGenerators struct {
	input     ScanMetadataInput
	taskQueue *job.TaskQueue
	progress  *job.Progress
}

func (g *sceneGenerators) Generate(ctx context.Context, s *models.Scene, f *file.VideoFile) error {
	const overwrite = false

	progress := g.progress
	t := g.input
	path := f.Path
	config := instance.Config
	fileNamingAlgorithm := config.GetVideoFileNamingAlgorithm()

	if t.ScanGenerateSprites {
		progress.AddTotal(1)
		g.taskQueue.Add(fmt.Sprintf("Generating sprites for %s", path), func(ctx context.Context) {
			taskSprite := GenerateSpriteTask{
				Scene:               *s,
				Overwrite:           overwrite,
				fileNamingAlgorithm: fileNamingAlgorithm,
			}
			taskSprite.Start(ctx)
			progress.Increment()
		})
	}

	if t.ScanGeneratePhashes {
		progress.AddTotal(1)
		g.taskQueue.Add(fmt.Sprintf("Generating phash for %s", path), func(ctx context.Context) {
			taskPhash := GeneratePhashTask{
				File:                f,
				fileNamingAlgorithm: fileNamingAlgorithm,
				txnManager:          instance.Database,
				fileUpdater:         instance.Database.File,
				Overwrite:           overwrite,
			}
			taskPhash.Start(ctx)
			progress.Increment()
		})
	}

	if t.ScanGeneratePreviews {
		progress.AddTotal(1)
		g.taskQueue.Add(fmt.Sprintf("Generating preview for %s", path), func(ctx context.Context) {
			options := getGeneratePreviewOptions(GeneratePreviewOptionsInput{})

			g := &generate.Generator{
				Encoder:     instance.FFMPEG,
				LockManager: instance.ReadLockManager,
				MarkerPaths: instance.Paths.SceneMarkers,
				ScenePaths:  instance.Paths.Scene,
				Overwrite:   overwrite,
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
		})
	}

	return nil
}
