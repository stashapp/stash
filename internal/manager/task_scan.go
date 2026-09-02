package manager

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/remeh/sizedwaitgroup"
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
	"github.com/stashapp/stash/pkg/utils"
)

type ScanJob struct {
	scanner       *file.Scanner
	input         ScanMetadataInput
	subscriptions *subscriptionManager

	fileQueue chan file.ScannedFile
	dirQueue  chan file.ScannedFile
	count     int

	unmatchedCaptionFiles utils.MutexField[[]string]

	// unchangedDirs records directories whose set of direct children has not
	// changed since the last scan. Files directly inside these directories are
	// skipped to avoid unnecessary DB lookups. Subdirectories are still walked
	// recursively so that deeper changes (new files, renamed dirs) are found.
	unchangedDirs sync.Map

	// skipUnchangedFolders is set per-path before walking; false for network
	// filesystems where ModTime may not be reliable.
	skipUnchangedFolders bool

	// dirCheckAttempts and dirCheckHits track the rolling hit rate of
	// CheckFolder calls. shouldCheckFolder uses these to gate calls once the
	// hit rate drops below checkFolderThreshold.
	dirCheckAttempts     atomic.Int64
	dirCheckHits         atomic.Int64
	checkFolderThreshold float64
}

// CheckFolder hit-rate gating — assumptions and context:
//
//   - After warmup, cumulative hits/attempts proxies whether more CheckFolder
//     calls pay off; walk order can skew this (single global ratio).
//
//   - A hit is CheckFolder returning unchanged: DB's folder ModTime equals this
//     directory's ModTime from disk (see Scanner.CheckFolder). Skipping enqueue for
//     direct files then assumes that invariant means nothing relevant changed among
//     those children. This is safe for local filesystems, but not for network or
//     coarse-modtime filesystems.
//
//   - When hits are rare (cold scan / folders absent or changed in DB), extra
//     CheckFolder calls are mostly overhead versus processing dirs/files without this
//     shortcut.
//
//   - Warmup assumes the first dirCheckWarmup CheckFolder calls are a representative
//     sample of folder outcomes for this walk. It is treated as forecasting the whole
//     job's hit rate. After warmup, a low cumulative rate stops further CheckFolder
//     (e.g. first scan where almost every folder is new).
//
// Basic idea: we always CheckFolder until `dirCheckWarmup` folders have been checked,
// then we continue to do so as long as the hit rate remains >= checkFolderHitRateThreshold
const (
	// dirCheckWarmup is the number of CheckFolder calls made before the
	// hit-rate signal is trusted. Keeps the first few directories enabled
	// unconditionally so the estimate is not dominated by noise.
	dirCheckWarmup = 50

	// checkFolderHitRateThreshold is the minimum hit rate below which
	// CheckFolder is skipped. Derived from benchmarks: break-even is ~0.50
	// (C_check_windows ≈ 139ms, C_save_per_dir ≈ 280ms); 0.30 is conservative,
	// favouring the warm-scan optimisation in mixed-state libraries.
	checkFolderHitRateThreshold = 0.30
)

// shouldCheckFolder returns true when calling CheckFolder is expected to pay
// for itself. It is always true during the warm-up window. After that it
// returns true only while the cumulative hit rate is at or above the threshold.
// Returns false unconditionally when the modtime optimisation is disabled
// (network FS) or when a forced rescan is requested (Rescan flag).
func (j *ScanJob) shouldCheckFolder() bool {
	if !j.skipUnchangedFolders || j.scanner.Rescan {
		return false
	}
	total := j.dirCheckAttempts.Load()
	if total < dirCheckWarmup {
		return true
	}
	return float64(j.dirCheckHits.Load())/float64(total) >= j.checkFolderThreshold
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

	nTasks := cfg.GetParallelTasksWithAutoDetection()

	const taskQueueSize = 200000
	taskQueue := job.NewTaskQueue(ctx, progress, taskQueueSize, nTasks)

	var minModTime time.Time
	if j.input.Filter != nil && j.input.Filter.MinModTime != nil {
		minModTime = *j.input.Filter.MinModTime
	}

	// HACK - these should really be set in the scanner initialization
	j.scanner.FileHandlers = getScanHandlers(j.input, taskQueue, progress)
	j.scanner.ScanFilters = []file.PathFilter{newScanFilter(c, repo, minModTime)}
	j.scanner.HandlerRequiredFilters = []file.Filter{newHandlerRequiredFilter(cfg, repo)}

	logger.Infof("Starting scan of %d paths with %d parallel tasks", len(paths), nTasks)

	j.runJob(ctx, paths, nTasks, progress)

	taskQueue.Close()

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}

	elapsed := time.Since(start)
	logger.Infof("Scan finished (%s)", elapsed)

	j.subscriptions.notify()
	return nil
}

func (j *ScanJob) runJob(ctx context.Context, paths []string, nTasks int, progress *job.Progress) {
	j.fileQueue = make(chan file.ScannedFile, scanQueueSize)
	j.dirQueue = make(chan file.ScannedFile, scanQueueSize)

	var walkWg sync.WaitGroup
	walkWg.Add(1)
	go func() {
		defer func() {
			walkWg.Done()

			// handle panics in goroutine
			if p := recover(); p != nil {
				logger.Errorf("panic while queuing files for scan: %v", p)
				logger.Errorf(string(debug.Stack()))
			}
		}()

		if err := j.queueFiles(ctx, paths, progress); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			logger.Errorf("error queuing files for scan: %v", err)
			return
		}

		logger.Infof("Finished adding files to queue. %d files queued", j.count)
	}()

	// dir writer: single goroutine drains dirQueue sequentially, preserving DFS
	// order so parent directories are always committed before their children.
	var dirWg sync.WaitGroup
	dirWg.Add(1)
	go func() {
		defer func() {
			dirWg.Done()

			if p := recover(); p != nil {
				logger.Errorf("panic while processing dir queue: %v", p)
				logger.Errorf(string(debug.Stack()))
			}
		}()

		for d := range j.dirQueue {
			j.processQueueItem(ctx, d, progress)
		}
	}()

	defer walkWg.Wait()
	defer dirWg.Wait()

	j.processQueue(ctx, nTasks, progress)
}

const scanQueueSize = 200000

func (j *ScanJob) queueFiles(ctx context.Context, paths []string, progress *job.Progress) error {
	fs := &file.OsFS{}

	defer func() {
		close(j.fileQueue)
		close(j.dirQueue)

		progress.AddTotal(j.count)
		progress.Definite()
	}()

	var err error
	progress.ExecuteTask("Walking directory tree", func() {
		j.dirCheckAttempts.Store(0)
		j.dirCheckHits.Store(0)

		for _, p := range paths {
			skipUnchanged := true
			isNet, netErr := fsutil.IsNetworkFS(p)
			if netErr != nil {
				logger.Warnf("Could not detect FS type for %s, disabling folder skip optimisation: %v", p, netErr)
				skipUnchanged = false
			} else if isNet {
				logger.Infof("Network filesystem detected for %s; disabling folder ModTime skip optimisation", p)
				skipUnchanged = false
			}
			if skipUnchanged {
				coarse, coarseErr := fsutil.CoarseModtimeFilesystem(p)
				if coarseErr != nil {
					logger.Warnf("Could not detect coarse-modtime filesystem for %s, disabling folder skip optimisation: %v", p, coarseErr)
					skipUnchanged = false
				} else if coarse {
					logger.Infof("FAT/exFAT or coarse-modtime filesystem detected for %s; disabling folder ModTime skip optimisation", p)
					skipUnchanged = false
				}
			}
			j.skipUnchangedFolders = skipUnchanged
			j.checkFolderThreshold = checkFolderHitRateThreshold

			err = file.SymWalk(fs, p, j.queueFileFunc(ctx, fs, nil, progress))
			if err != nil {
				return
			}
		}
	})

	return err
}

func (j *ScanJob) queueFileFunc(ctx context.Context, f models.FS, zipFile *file.ScannedFile, progress *job.Progress) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// don't let errors prevent scanning
			logger.Errorf("error scanning %s: %v", path, err)
			return nil
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			logger.Errorf("reading info for %q: %v", path, err)
			return nil
		}

		zipFilePath := ""
		if zipFile != nil {
			zipFilePath = zipFile.Path
		}

		if !j.scanner.AcceptEntry(ctx, path, info, zipFilePath) {
			if info.IsDir() {
				logger.Debugf("Skipping directory %s", path)
				return fs.SkipDir
			}

			// we don't include caption files in the file scan, but we do need
			// to handle them
			if fsutil.MatchExtension(path, video.CaptionExts) {
				fileRepo := j.scanner.Repository.File
				matched := video.AssociateCaptions(ctx, path, j.scanner.Repository.TxnManager, fileRepo, fileRepo)

				if !matched {
					logger.Debugf("No matching video file found for caption file %s", path)
					j.unmatchedCaptionFiles.SetFunc(func(files []string) []string {
						return append(files, path)
					})
				}

				return nil
			}

			logger.Debugf("Skipping file %s", path)
			return nil
		}

		size, err := file.GetFileSize(f, path, info)
		if err != nil {
			return err
		}

		// #4425 - store the NFC-normalized path so search matches user input.
		// Filtering and filesystem access above use the original path; zip
		// entries are matched byte-exact so are left unnormalized.
		storedPath := path
		if zipFile == nil {
			storedPath = fsutil.NormalizePath(path)
		}

		ff := file.ScannedFile{
			BaseFile: &models.BaseFile{
				DirEntry: models.DirEntry{
					ModTime: file.ModTime(info),
				},
				Path:     storedPath,
				Basename: filepath.Base(storedPath),
				Size:     size,
			},
			FS:   f,
			Info: info,
		}

		if zipFile != nil {
			ff.ZipFileID = &zipFile.ID
			ff.ZipFile = zipFile
		}

		if info.IsDir() {
			if j.shouldCheckFolder() {
				j.dirCheckAttempts.Add(1)
				_, unchanged, err := j.scanner.CheckFolder(ctx, ff)
				if err != nil {
					if !errors.Is(err, context.Canceled) {
						logger.Errorf("error checking folder %q: %v", path, err)
					}
					return fs.SkipDir
				}
				if unchanged {
					j.dirCheckHits.Add(1)
					j.unchangedDirs.Store(ff.Path, struct{}{})
					return nil
				}
			}
			j.dirQueue <- ff
			return nil
		}

		// Skip files whose containing directory is known-unchanged. Their set
		// of direct children has not changed, so no file was added here. We
		// still walked into the parent to allow deeper subdirectories to be
		// checked. Use Rescan to force re-examining file content.
		if j.skipUnchangedFolders && !j.scanner.Rescan {
			if _, isUnchanged := j.unchangedDirs.Load(filepath.Dir(path)); isUnchanged {
				logger.Tracef("Skipping file %s in unchanged directory", path)
				return nil
			}
		}

		// if zip file is present, we handle immediately
		if zipFile != nil {
			progress.ExecuteTask("Scanning "+path, func() {
				// don't increment progress in zip files
				if err := j.handleFile(ctx, ff, nil); err != nil {
					if !errors.Is(err, context.Canceled) {
						logger.Errorf("error processing %q: %v", path, err)
					}
					// don't return an error, just skip the file
				}
			})

			return nil
		}

		logger.Tracef("Queueing file %s for scanning", path)
		j.fileQueue <- ff

		j.count++

		return nil
	}
}

func (j *ScanJob) processQueue(ctx context.Context, parallelTasks int, progress *job.Progress) {
	if parallelTasks < 1 {
		parallelTasks = 1
	}

	wg := sizedwaitgroup.New(parallelTasks)

	func() {
		defer func() {
			wg.Wait()

			// handle panics in goroutine
			if p := recover(); p != nil {
				logger.Errorf("panic while scanning files: %v", p)
				logger.Errorf(string(debug.Stack()))
			}
		}()

		for f := range j.fileQueue {
			logger.Tracef("Processing queued file %s", f.Path)
			if ctx.Err() != nil {
				// Keep receiving until queueFiles closes the channel; otherwise
				// the walker can block on send (full buffer) and never finish.
				continue
			}

			wg.Add()
			ff := f
			go func() {
				defer wg.Done()
				j.processQueueItem(ctx, ff, progress)
			}()
		}
	}()
}

func (j *ScanJob) processQueueItem(ctx context.Context, f file.ScannedFile, progress *job.Progress) {
	progress.ExecuteTask("Scanning "+f.Path, func() {
		var err error
		if f.Info.IsDir() {
			err = j.handleFolder(ctx, f, progress)
		} else {
			err = j.handleFile(ctx, f, progress)
		}

		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf("error processing %q: %v", f.Path, err)
		}
	})
}

func (j *ScanJob) handleFolder(ctx context.Context, f file.ScannedFile, progress *job.Progress) error {
	if progress != nil {
		defer progress.Increment()
	}

	_, changed, err := j.scanner.ScanFolder(ctx, f)
	if err != nil {
		return err
	}

	// Record unchanged directories so their direct file children can be
	// skipped. We never return fs.SkipDir here: subdirectories must still be
	// walked because their own modtime may have changed (a new file added to
	// /a/b only updates /a/b/'s modtime, not /a/'s).
	if !changed && j.skipUnchangedFolders && !j.scanner.Rescan {
		j.unchangedDirs.Store(f.Path, struct{}{})
	}

	return nil
}

func (j *ScanJob) handleFile(ctx context.Context, f file.ScannedFile, progress *job.Progress) error {
	if progress != nil {
		defer progress.Increment()
	}

	r, err := j.scanner.ScanFile(ctx, f)
	if err != nil {
		return err
	}

	// if this is a new video file, match it with any unmatched caption files
	if r.New && len(j.unmatchedCaptionFiles.Get()) > 0 {
		videoFile, _ := r.File.(*models.VideoFile)

		if videoFile != nil {
			// try to match any unmatched caption files to this video file
			for _, captionPath := range j.unmatchedCaptionFiles.Get() {
				if video.MatchesCaption(videoFile.Path, captionPath) {
					video.AssociateCaptions(ctx, captionPath, j.scanner.Repository.TxnManager, j.scanner.Repository.File, j.scanner.Repository.File)

					// remove from the unmatched list
					j.unmatchedCaptionFiles.SetFunc(func(files []string) []string {
						newFiles := make([]string, 0, len(files)-1)
						for _, f := range files {
							if f != captionPath {
								newFiles = append(newFiles, f)
							}
						}
						return newFiles
					})
				}
			}
		}
	}

	// clean captions - scene handler handles this as well, but
	// unchanged files aren't processed by the scene handler
	if r.IsUnchanged() {
		videoFile, _ := r.File.(*models.VideoFile)

		if videoFile != nil {
			txnMgr := j.scanner.Repository.TxnManager
			fileRepo := j.scanner.Repository.File
			if err := txn.WithDatabase(ctx, txnMgr, func(ctx context.Context) error {
				return video.CleanCaptions(ctx, videoFile, txnMgr, fileRepo)
			}); err != nil {
				logger.Errorf("Error cleaning captions: %v", err)
			}
		}
	}

	// handle rename should have already handled the contents of the zip file
	// so shouldn't need to scan it again.
	// Only scan zip contents if the file is new, the fingerprint changed,
	// if a force rescan was requested, or if the handler was required because
	// a related object (e.g. a deleted gallery) is missing and needs to be
	// recreated from the contents.

	if j.scanner.IsZipFile(f.Info.Name()) && (r.New || r.FingerprintChanged || j.scanner.Rescan || r.HandlerRequired) {
		ff := r.File
		f.BaseFile = ff.Base()

		// scan zip files with a different context that is not cancellable
		// cancelling while scanning zip file contents results in the scan
		// contents being partially completed
		zipCtx := context.WithoutCancel(ctx)

		if err := j.scanZipFile(zipCtx, f, progress); err != nil {
			logger.Errorf("Error scanning zip file %q: %v", f.Path, err)
		}
	} else if r.Updated && j.scanner.IsZipFile(f.Info.Name()) {
		logger.Debugf("Skipping zip file scan for %q: fingerprint unchanged", f.Path)
	}

	return nil
}

func (j *ScanJob) scanZipFile(ctx context.Context, f file.ScannedFile, progress *job.Progress) error {
	zipFS, err := f.FS.OpenZip(f.Path, f.Size)
	if err != nil {
		if errors.Is(err, file.ErrNotReaderAt) {
			// can't walk the zip file
			// just return
			logger.Debugf("Skipping zip file %q as it cannot be opened for walking", f.Path)
			return nil
		}

		return err
	}

	defer zipFS.Close()

	return file.SymWalk(zipFS, f.Path, j.queueFileFunc(ctx, zipFS, &f, progress))
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
	txnManager    txn.Manager
	SceneFinder   sceneFinder
	ImageFinder   fileCounter
	GalleryFinder galleryFinder

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

	return false
}

type scanFilter struct {
	extensionConfig
	txnManager txn.Manager

	stashPaths        config.StashConfigs
	generatedPath     string
	videoExcludeRegex []*regexp.Regexp
	imageExcludeRegex []*regexp.Regexp
	minModTime        time.Time
	stashIgnoreFilter *file.StashIgnoreFilter
}

// shouldSkipDir returns true if dir matches the video (and image) exclude patterns
// and no configured stash root is at or beneath dir.
func shouldSkipDir(dir string, s *config.StashConfig, allPaths config.StashConfigs, videoRE, imageRE []*regexp.Regexp) bool {
	test := dir + string(filepath.Separator)

	if !matchFileRegex(test, videoRE) {
		return false
	}

	if !s.ExcludeImage && !matchFileRegex(test, imageRE) {
		return false
	}

	cleanDir := filepath.Clean(dir)
	for _, sc := range allPaths {
		if fsutil.IsPathInDir(cleanDir, filepath.Clean(sc.Path)) {
			return false
		}
	}

	return true
}

func newScanFilter(c *config.Config, repo models.Repository, minModTime time.Time) *scanFilter {
	return &scanFilter{
		extensionConfig:   newExtensionConfig(c),
		txnManager:        repo.TxnManager,
		stashPaths:        c.GetStashPaths(),
		generatedPath:     c.GetGeneratedPath(),
		videoExcludeRegex: generateRegexps(c.GetExcludes()),
		imageExcludeRegex: generateRegexps(c.GetImageExcludes()),
		minModTime:        minModTime,
		stashIgnoreFilter: file.NewStashIgnoreFilter(),
	}
}

func (f *scanFilter) Accept(ctx context.Context, path string, info fs.FileInfo, zipFilePath string) bool {
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

	// Check .stashignore files, bounded to the library root.
	if !f.stashIgnoreFilter.Accept(ctx, path, info, f.stashPaths.GetStashRootFromDirPath(path), zipFilePath) {
		logger.Debugf("Skipping %s due to .stashignore", path)
		return false
	}

	isVideoFile := useAsVideo(path)
	isImageFile := useAsImage(path)
	isZipFile := fsutil.MatchExtension(path, f.zipExt)

	if !info.IsDir() && !isVideoFile && !isImageFile && !isZipFile {
		logger.Debugf("Skipping %s as it does not match any known file extensions", path)
		return false
	}

	// #1756 - skip zero length files
	if !info.IsDir() && info.Size() == 0 {
		logger.Infof("Skipping zero-length file: %s", path)
		return false
	}

	if info.IsDir() && shouldSkipDir(path, s, f.stashPaths, f.videoExcludeRegex, f.imageExcludeRegex) {
		logger.Debugf("Skipping directory %s as it matches exclude patterns", path)
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
				CreatorUpdater:     r.Image,
				GalleryFinder:      r.Gallery,
				SceneFinderUpdater: r.Scene,
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
				CreatorUpdater:       r.Scene,
				GalleryFinderUpdater: r.Gallery,
				CaptionUpdater:       r.File,
				PluginCache:          pluginCache,
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

	if t.ScanGenerateImagePhashes {
		progress.AddTotal(1)
		phashFn := func(ctx context.Context) {
			mgr := GetInstance()
			// Only generate phash for image files, not video files
			if imageFile, ok := f.(*models.ImageFile); ok {
				taskPhash := GenerateImagePhashTask{
					repository: mgr.Repository,
					File:       imageFile,
					Overwrite:  overwrite,
				}
				taskPhash.Start(ctx)
			}
			progress.Increment()
		}

		if g.sequentialScanning {
			phashFn(ctx)
		} else {
			g.taskQueue.Add(fmt.Sprintf("Generating phash for %s", path), phashFn)
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
