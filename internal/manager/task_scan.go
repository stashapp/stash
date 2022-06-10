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

	j.scanner.Scan(ctx, getScanHandlers(j.input), file.ScanOptions{
		Paths:             paths,
		ScanFilters:       []file.PathFilter{newScanFilter(instance.Config)},
		ZipFileExtensions: instance.Config.GetGalleryExtensions(),
	}, progress)

	// FIXME - handle generate jobs after scanning

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Scan finished (%s)", elapsed))

	j.subscriptions.notify()
}

type scanFilter struct {
	stashPaths        []*config.StashConfig
	generatedPath     string
	vidExt            []string
	imgExt            []string
	zipExt            []string
	videoExcludeRegex []*regexp.Regexp
	imageExcludeRegex []*regexp.Regexp
}

func newScanFilter(c *config.Instance) *scanFilter {
	return &scanFilter{
		stashPaths:        c.GetStashPaths(),
		generatedPath:     c.GetGeneratedPath(),
		vidExt:            c.GetVideoExtensions(),
		imgExt:            c.GetImageExtensions(),
		zipExt:            c.GetGalleryExtensions(),
		videoExcludeRegex: generateRegexps(c.GetExcludes()),
		imageExcludeRegex: generateRegexps(c.GetImageExcludes()),
	}
}

func (f *scanFilter) Accept(ctx context.Context, path string, info fs.FileInfo) bool {
	if fsutil.IsPathInDir(f.generatedPath, path) {
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

func getScanHandlers(options ScanMetadataInput) []file.Handler {
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
