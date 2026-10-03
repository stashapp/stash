package manager

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// TODO: ideally this would be in the task package, but deferring for now

type verifier interface {
	Verify(ctx context.Context, options file.VerifyOptions, progress *job.Progress)
}

type verifyJob struct {
	verifier     verifier
	repository   models.Repository
	options      file.VerifyOptions
	sceneService SceneService
	imageService ImageService
	scanSubs     *subscriptionManager
}

func (j *verifyJob) Execute(ctx context.Context, progress *job.Progress) error {
	logger.Infof("Starting verification of tracked files")
	start := time.Now()
	if j.options.DryRun {
		logger.Infof("Running in Dry Mode")
	}

	options := j.options
	options.PathFilter = newVerifyFilter(instance.Config)

	j.verifier.Verify(ctx, options, progress)

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}

	// only clean empty galleries if purging
	if !j.options.DryRun && j.options.PurgeMissing {
		// HACK - use purge job to clean empty galleries
		pj := &purgeMissingJob{
			// only need to provide repository
			repository: j.repository,
		}
		pj.cleanEmptyGalleries(ctx)
	}

	j.scanSubs.notify()
	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Finished verifying (%s)", elapsed))
	return nil
}

type verifyFilter struct {
	scanFilter
}

func newVerifyFilter(c *config.Config) *verifyFilter {
	return &verifyFilter{
		scanFilter: scanFilter{
			extensionConfig:   newExtensionConfig(c),
			stashPaths:        c.GetStashPaths(),
			generatedPath:     c.GetGeneratedPath(),
			videoExcludeRegex: generateRegexps(c.GetExcludes()),
			imageExcludeRegex: generateRegexps(c.GetImageExcludes()),
			stashIgnoreFilter: file.NewStashIgnoreFilter(),
		},
	}
}

func (f *verifyFilter) Accept(ctx context.Context, path string, info fs.FileInfo, zipFilePath string) bool {
	//  #1102 - clean anything in generated path
	generatedPath := f.generatedPath

	var stash *config.StashConfig
	fileOrFolder := "File"

	if info.IsDir() {
		fileOrFolder = "Folder"
		stash = f.stashPaths.GetStashFromDirPath(path)
	} else {
		stash = f.stashPaths.GetStashFromPath(path)
	}

	if stash == nil {
		logger.Infof("%s not in any stash library directories. Marking as missing: %q", fileOrFolder, path)
		return false
	}

	if fsutil.IsPathInDir(generatedPath, path) {
		logger.Infof("%s is in generated path. Marking as missing: %q", fileOrFolder, path)
		return false
	}

	// Check .stashignore files, bounded to the library root.
	if !f.stashIgnoreFilter.Accept(ctx, path, info, f.stashPaths.GetStashRootFromDirPath(path), zipFilePath) {
		logger.Infof("%s is excluded due to .stashignore. Marking as missing: %q", fileOrFolder, path)
		return false
	}

	if info.IsDir() {
		return !f.shouldCleanFolder(path, stash)
	}

	return !f.shouldCleanFile(path, info, stash)
}

func (f *verifyFilter) shouldCleanFolder(path string, s *config.StashConfig) bool {
	// only delete folders where it is excluded from everything
	pathExcludeTest := path + string(filepath.Separator)
	if (s.ExcludeVideo || matchFileRegex(pathExcludeTest, f.videoExcludeRegex)) && (s.ExcludeImage || matchFileRegex(pathExcludeTest, f.imageExcludeRegex)) {
		logger.Infof("Folder is excluded from both video and image. Marking as missing: \"%s\"", path)
		return true
	}

	return false
}

func (f *verifyFilter) shouldCleanFile(path string, info fs.FileInfo, stash *config.StashConfig) bool {
	switch {
	case info.IsDir() || fsutil.MatchExtension(path, f.zipExt):
		return f.shouldCleanGallery(path, stash)
	case useAsVideo(path):
		return f.shouldCleanVideoFile(path, stash)
	case useAsImage(path):
		return f.shouldCleanImage(path, stash)
	default:
		logger.Infof("File extension does not match any media extensions. Marking as missing: \"%s\"", path)
		return true
	}
}

func (f *verifyFilter) shouldCleanVideoFile(path string, stash *config.StashConfig) bool {
	if stash.ExcludeVideo {
		logger.Infof("File in stash library that excludes video. Marking as missing: \"%s\"", path)
		return true
	}

	if matchFileRegex(path, f.videoExcludeRegex) {
		logger.Infof("File matched regex. Marking as missing: \"%s\"", path)
		return true
	}

	return false
}

func (f *verifyFilter) shouldCleanGallery(path string, stash *config.StashConfig) bool {
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking as missing: \"%s\"", path)
		return true
	}

	if matchFileRegex(path, f.imageExcludeRegex) {
		logger.Infof("File matched regex. Marking as missing: \"%s\"", path)
		return true
	}

	return false
}

func (f *verifyFilter) shouldCleanImage(path string, stash *config.StashConfig) bool {
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking as missing: \"%s\"", path)
		return true
	}

	if matchFileRegex(path, f.imageExcludeRegex) {
		logger.Infof("File matched regex. Marking as missing: \"%s\"", path)
		return true
	}

	return false
}
