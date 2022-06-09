package manager

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
)

type scanner interface {
	Scan(ctx context.Context, options file.ScanOptions, progressReporter file.ProgressReporter)
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

	j.scanner.Scan(ctx, file.ScanOptions{
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

	s := getStashFromDirPath(path)

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
	} else if !isVideoFile && s.ExcludeImage || matchFileRegex(path, f.imageExcludeRegex) {
		return false
	}

	return true
}
