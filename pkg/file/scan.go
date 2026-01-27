package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

const (
	scanQueueSize = 200000
	// maximum number of times to retry in the event of a locked database
	// use -1 to retry forever
	maxRetries = -1
)

// Scanner scans files into the database.
//
// The scan process works using two goroutines. The first walks through the provided paths
// in the filesystem. It runs each directory entry through the provided ScanFilters. If none
// of the filter Accept methods return true, then the file/directory is ignored.
// Any folders found are handled immediately. Files inside zip files are also handled immediately.
// All other files encountered are sent to the second goroutine queue.
//
// Folders are handled by checking if the folder exists in the database, by its full path.
// If a folder entry already exists, then its mod time is updated (if applicable).
// If the folder does not exist in the database, then a new folder entry its created.
//
// Files are handled by first querying for the file by its path. If the file entry exists in the
// database, then the mod time is compared to the value in the database. If the mod time is different
// then file is marked as updated - it recalculates any fingerprints and fires decorators, then
// the file entry is updated and any applicable handlers are fired.
//
// If the file entry does not exist in the database, then fingerprints are calculated for the file.
// It then determines if the file is a rename of an existing file by querying for file entries with
// the same fingerprint. If any are found, it checks each to see if any are missing in the file
// system. If one is, then the file is treated as renamed and its path is updated. If none are missing,
// or many are, then the file is treated as a new file.
//
// If the file is not a renamed file, then the decorators are fired and the file is created, then
// the applicable handlers are fired.
type Scanner struct {
	FS                    models.FS
	Repository            Repository
	FingerprintCalculator FingerprintCalculator

	// FileDecorators are applied to files as they are scanned.
	FileDecorators []Decorator
}

// FingerprintCalculator calculates a fingerprint for the provided file.
type FingerprintCalculator interface {
	CalculateFingerprints(f *models.BaseFile, o Opener, useExisting bool) ([]models.Fingerprint, error)
}

// Decorator wraps the Decorate method to add additional functionality while scanning files.
type Decorator interface {
	Decorate(ctx context.Context, fs models.FS, f models.File) (models.File, error)
	IsMissingMetadata(ctx context.Context, fs models.FS, f models.File) bool
}

type FilteredDecorator struct {
	Decorator
	Filter
}

// Decorate runs the decorator if the filter accepts the file.
func (d *FilteredDecorator) Decorate(ctx context.Context, fs models.FS, f models.File) (models.File, error) {
	if d.Accept(ctx, f) {
		return d.Decorator.Decorate(ctx, fs, f)
	}
	return f, nil
}

func (d *FilteredDecorator) IsMissingMetadata(ctx context.Context, fs models.FS, f models.File) bool {
	if d.Accept(ctx, f) {
		return d.Decorator.IsMissingMetadata(ctx, fs, f)
	}

	return false
}

// ProgressReporter is used to report progress of the scan.
type ProgressReporter interface {
	AddTotal(total int)
	Increment()
	Definite()
	ExecuteTask(description string, fn func())
}

type scanJob struct {
	*Scanner

	// handlers are called after a file has been scanned.
	handlers []Handler

	ProgressReports ProgressReporter
	options         ScanOptions

	startTime      time.Time
	fileQueue      chan scanFile
	retryList      []scanFile
	retrying       bool
	folderPathToID sync.Map
	zipPathToID    sync.Map
	count          int

	txnRetryer txn.Retryer
}

// ScanOptions provides options for scanning files.
type ScanOptions struct {
	Paths []string

	// ZipFileExtensions is a list of file extensions that are considered zip files.
	// Extension does not include the . character.
	ZipFileExtensions []string

	// ScanFilters are used to determine if a file should be scanned.
	ScanFilters []PathFilter

	// HandlerRequiredFilters are used to determine if an unchanged file needs to be handled
	HandlerRequiredFilters []Filter

	ParallelTasks int

	// When true files in path will be rescanned even if they haven't changed
	Rescan bool
}

// Scan starts the scanning process.
func (s *Scanner) Scan(ctx context.Context, handlers []Handler, options ScanOptions, progressReporter ProgressReporter) {
	job := &scanJob{
		Scanner:         s,
		handlers:        handlers,
		ProgressReports: progressReporter,
		options:         options,
		txnRetryer: txn.Retryer{
			Manager: s.Repository.TxnManager,
			Retries: maxRetries,
		},
	}

	job.execute(ctx)
}

type scanFile struct {
	*models.BaseFile
	fs   models.FS
	info fs.FileInfo
}

func (s *scanJob) withTxn(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.txnRetryer.WithTxn(ctx, fn)
}

func (s *scanJob) withDB(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.Repository.WithDB(ctx, fn)
}

func (s *scanJob) execute(ctx context.Context) {
	paths := s.options.Paths
	logger.Infof("scanning %d paths", len(paths))
	s.startTime = time.Now()

	s.fileQueue = make(chan scanFile, scanQueueSize)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()

			// handle panics in goroutine
			if p := recover(); p != nil {
				logger.Errorf("panic while queuing files for scan: %v", p)
				logger.Errorf(string(debug.Stack()))
			}
		}()

		if err := s.queueFiles(ctx, paths); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			logger.Errorf("error queuing files for scan: %v", err)
			return
		}

		logger.Infof("Finished adding files to queue. %d files queued", s.count)
	}()

	defer wg.Wait()

	if err := s.processQueue(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		logger.Errorf("error scanning files: %v", err)
		return
	}
}

func (s *scanJob) queueFiles(ctx context.Context, paths []string) error {
	defer func() {
		close(s.fileQueue)

		if s.ProgressReports != nil {
			s.ProgressReports.AddTotal(s.count)
			s.ProgressReports.Definite()
		}
	}()

	var err error
	s.ProgressReports.ExecuteTask("Walking directory tree", func() {
		for _, p := range paths {
			err = symWalk(s.FS, p, s.queueFileFunc(ctx, s.FS, nil))
			if err != nil {
				return
			}
		}
	})

	return err
}

func (s *scanJob) queueFileFunc(ctx context.Context, f models.FS, zipFile *scanFile) fs.WalkDirFunc {
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

		if !s.acceptEntry(ctx, path, info) {
			if info.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		size, err := getFileSize(f, path, info)
		if err != nil {
			return err
		}

		ff := scanFile{
			BaseFile: &models.BaseFile{
				DirEntry: models.DirEntry{
					ModTime: modTime(info),
				},
				Path:     path,
				Basename: filepath.Base(path),
				Size:     size,
			},
			fs:   f,
			info: info,
		}

		if zipFile != nil {
			zipFileID, err := s.getZipFileID(ctx, zipFile)
			if err != nil {
				return err
			}
			ff.ZipFileID = zipFileID
			ff.ZipFile = zipFile
		}

		if info.IsDir() {
			// handle folders immediately
			if err := s.handleFolder(ctx, ff); err != nil {
				if !errors.Is(err, context.Canceled) {
					logger.Errorf("error processing %q: %v", path, err)
				}

				// skip the directory since we won't be able to process the files anyway
				return fs.SkipDir
			}

			return nil
		}

		// if zip file is present, we handle immediately
		if zipFile != nil {
			s.ProgressReports.ExecuteTask("Scanning "+path, func() {
				if err := s.handleFile(ctx, ff); err != nil {
					if !errors.Is(err, context.Canceled) {
						logger.Errorf("error processing %q: %v", path, err)
					}
					// don't return an error, just skip the file
				}
			})

			return nil
		}

		s.fileQueue <- ff

		s.count++

		return nil
	}
}

func getFileSize(f models.FS, path string, info fs.FileInfo) (int64, error) {
	// #2196/#3042 - replace size with target size if file is a symlink
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		targetInfo, err := f.Stat(path)
		if err != nil {
			return 0, fmt.Errorf("reading info for symlink %q: %w", path, err)
		}
		return targetInfo.Size(), nil
	}

	return info.Size(), nil
}

func (s *scanJob) acceptEntry(ctx context.Context, path string, info fs.FileInfo) bool {
	// always accept if there's no filters
	accept := len(s.options.ScanFilters) == 0
	for _, filter := range s.options.ScanFilters {
		// accept if any filter accepts the file
		if filter.Accept(ctx, path, info) {
			accept = true
			break
		}
	}

	return accept
}

func (s *scanJob) scanZipFile(ctx context.Context, f scanFile) error {
	zipFS, err := f.fs.OpenZip(f.Path, f.Size)
	if err != nil {
		if errors.Is(err, errNotReaderAt) {
			// can't walk the zip file
			// just return
			return nil
		}

		return err
	}

	defer zipFS.Close()

	return symWalk(zipFS, f.Path, s.queueFileFunc(ctx, zipFS, &f))
}

func (s *scanJob) processQueue(ctx context.Context) error {
	parallelTasks := s.options.ParallelTasks
	if parallelTasks < 1 {
		parallelTasks = 1
	}

	wg := sizedwaitgroup.New(parallelTasks)

	if err := func() error {
		defer wg.Wait()

		for f := range s.fileQueue {
			if err := ctx.Err(); err != nil {
				return err
			}

			wg.Add()
			ff := f
			go func() {
				defer wg.Done()
				s.processQueueItem(ctx, ff)
			}()
		}

		return nil
	}(); err != nil {
		return err
	}

	s.retrying = true

	if err := func() error {
		defer wg.Wait()

		for _, f := range s.retryList {
			if err := ctx.Err(); err != nil {
				return err
			}

			wg.Add()
			ff := f
			go func() {
				defer wg.Done()
				s.processQueueItem(ctx, ff)
			}()
		}

		return nil
	}(); err != nil {
		return err
	}

	return nil
}

func (s *scanJob) incrementProgress(f scanFile) {
	// don't increment for files inside zip files since these aren't
	// counted during the initial walking
	if s.ProgressReports != nil && f.ZipFile == nil {
		s.ProgressReports.Increment()
	}
}

func (s *scanJob) processQueueItem(ctx context.Context, f scanFile) {
	s.ProgressReports.ExecuteTask("Scanning "+f.Path, func() {
		var err error
		if f.info.IsDir() {
			err = s.handleFolder(ctx, f)
		} else {
			err = s.handleFile(ctx, f)
		}

		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf("error processing %q: %v", f.Path, err)
		}
	})
}

func (s *scanJob) getFolderID(ctx context.Context, path string) (*models.FolderID, error) {
	// check the folder cache first
	if f, ok := s.folderPathToID.Load(path); ok {
		v := f.(models.FolderID)
		return &v, nil
	}

	// assume case sensitive when searching for the folder
	const caseSensitive = true

	ret, err := s.Repository.Folder.FindByPath(ctx, path, caseSensitive)
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, nil
	}

	s.folderPathToID.Store(path, ret.ID)
	return &ret.ID, nil
}

func (s *scanJob) getZipFileID(ctx context.Context, zipFile *scanFile) (*models.FileID, error) {
	if zipFile == nil {
		return nil, nil
	}

	if zipFile.ID != 0 {
		return &zipFile.ID, nil
	}

	path := zipFile.Path

	// check the folder cache first
	if f, ok := s.zipPathToID.Load(path); ok {
		v := f.(models.FileID)
		return &v, nil
	}

	// assume case sensitive when searching for the zip file
	const caseSensitive = true

	ret, err := s.Repository.File.FindByPath(ctx, path, caseSensitive)
	if err != nil {
		return nil, fmt.Errorf("getting zip file ID for %q: %w", path, err)
	}

	if ret == nil {
		return nil, fmt.Errorf("zip file %q doesn't exist in database", zipFile.Path)
	}

	s.zipPathToID.Store(path, ret.Base().ID)
	return &ret.Base().ID, nil
}

func (s *scanJob) handleFolder(ctx context.Context, file scanFile) error {
	path := file.Path

	return s.withTxn(ctx, func(ctx context.Context) error {
		defer s.incrementProgress(file)

		// determine if folder already exists in data store (by path)
		// assume case sensitive by default
		f, err := s.Repository.Folder.FindByPath(ctx, path, true)
		if err != nil {
			return fmt.Errorf("checking for existing folder %q: %w", path, err)
		}

		// #1426 / #6326 - if folder is in a case-insensitive filesystem, then try
		// case insensitive searching
		// assume case sensitive if in zip
		if f == nil && file.ZipFileID == nil {
			caseSensitive, _ := file.fs.IsPathCaseSensitive(file.Path)

			if !caseSensitive {
				f, err = s.Repository.Folder.FindByPath(ctx, path, false)
				if err != nil {
					return fmt.Errorf("checking for existing folder %q: %w", path, err)
				}
			}
		}

		// if folder not exists, create it
		if f == nil {
			f, err = s.onNewFolder(ctx, file)
		} else {
			f, err = s.onExistingFolder(ctx, file, f)
		}

		if err != nil {
			return err
		}

		if f != nil {
			s.folderPathToID.Store(f.Path, f.ID)
		}

		return nil
	})
}

func (s *scanJob) onNewFolder(ctx context.Context, file scanFile) (*models.Folder, error) {
	renamed, err := s.handleFolderRename(ctx, file)
	if err != nil {
		return nil, err
	}

	if renamed != nil {
		return renamed, nil
	}

	now := time.Now()

	toCreate := &models.Folder{
		DirEntry:  file.DirEntry,
		Path:      file.Path,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dir := filepath.Dir(file.Path)
	if dir != "." {
		parentFolderID, err := s.getFolderID(ctx, dir)
		if err != nil {
			return nil, fmt.Errorf("getting parent folder %q: %w", dir, err)
		}

		// if parent folder doesn't exist, assume it's a top-level folder
		// this may not be true if we're using multiple goroutines
		if parentFolderID != nil {
			toCreate.ParentFolderID = parentFolderID
		}
	}

	txn.AddPostCommitHook(ctx, func(ctx context.Context) {
		// log at the end so that if anything fails above due to a locked database
		// error and the transaction must be retried, then we shouldn't get multiple
		// logs of the same thing.
		logger.Infof("%s doesn't exist. Creating new folder entry...", file.Path)
	})

	if err := s.Repository.Folder.Create(ctx, toCreate); err != nil {
		return nil, fmt.Errorf("creating folder %q: %w", file.Path, err)
	}

	return toCreate, nil
}

func (s *scanJob) handleFolderRename(ctx context.Context, file scanFile) (*models.Folder, error) {
	// ignore folders in zip files
	if file.ZipFileID != nil {
		return nil, nil
	}

	// check if the folder was moved from elsewhere
	renamedFrom, err := s.detectFolderMove(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("detecting folder move: %w", err)
	}

	if renamedFrom == nil {
		return nil, nil
	}

	// if the folder was moved, update the existing folder
	logger.Infof("%s moved to %s. Updating path...", renamedFrom.Path, file.Path)
	renamedFrom.Path = file.Path

	// update the parent folder ID
	// find the parent folder
	parentFolderID, err := s.getFolderID(ctx, filepath.Dir(file.Path))
	if err != nil {
		return nil, fmt.Errorf("getting parent folder for %q: %w", file.Path, err)
	}

	renamedFrom.ParentFolderID = parentFolderID

	if err := s.Repository.Folder.Update(ctx, renamedFrom); err != nil {
		return nil, fmt.Errorf("updating folder for rename %q: %w", renamedFrom.Path, err)
	}

	// #4146 - correct sub-folders to have the correct path
	if err := correctSubFolderHierarchy(ctx, s.Repository.Folder, renamedFrom); err != nil {
		return nil, fmt.Errorf("correcting sub folder hierarchy for %q: %w", renamedFrom.Path, err)
	}

	return renamedFrom, nil
}

func (s *scanJob) onExistingFolder(ctx context.Context, f scanFile, existing *models.Folder) (*models.Folder, error) {
	update := false

	// update if mod time is changed
	entryModTime := f.ModTime
	if !entryModTime.Equal(existing.ModTime) {
		existing.Path = f.Path
		existing.ModTime = entryModTime
		update = true
	}

	// #6326 - update if path has changed - should only happen if case is
	// changed and filesystem is case insensitive
	if existing.Path != f.Path {
		existing.Path = f.Path
		update = true
	}

	// update if zip file ID has changed
	fZfID := f.ZipFileID
	existingZfID := existing.ZipFileID
	if fZfID != existingZfID {
		if fZfID == nil {
			existing.ZipFileID = nil
			update = true
		} else if existingZfID == nil || *fZfID != *existingZfID {
			existing.ZipFileID = fZfID
			update = true
		}
	}

	if update {
		var err error
		if err = s.Repository.Folder.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("updating folder %q: %w", f.Path, err)
		}
	}

	return existing, nil
}

func modTime(info fs.FileInfo) time.Time {
	// truncate to seconds, since we don't store beyond that in the database
	return info.ModTime().Truncate(time.Second)
}

func (s *scanJob) handleFile(ctx context.Context, f scanFile) error {
	defer s.incrementProgress(f)

	var ff models.File

	// don't use a transaction to check if new or existing
	if err := s.withDB(ctx, func(ctx context.Context) error {
		// determine if file already exists in data store
		// assume case sensitive when searching for the file to begin with
		var err error
		ff, err = s.Repository.File.FindByPath(ctx, f.Path, true)
		if err != nil {
			return fmt.Errorf("checking for existing file %q: %w", f.Path, err)
		}

		// #1426 / #6326 - if file is in a case-insensitive filesystem, then try
		// case insensitive search
		// assume case sensitive if in zip
		if ff == nil && f.ZipFileID != nil {
			caseSensitive, _ := f.fs.IsPathCaseSensitive(f.Path)

			if !caseSensitive {
				ff, err = s.Repository.File.FindByPath(ctx, f.Path, false)
				if err != nil {
					return fmt.Errorf("checking for existing file %q: %w", f.Path, err)
				}
			}
		}

		if ff == nil {
			// returns a file only if it is actually new
			ff, err = s.onNewFile(ctx, f)
			return err
		}

		ff, err = s.onExistingFile(ctx, f, ff)
		return err
	}); err != nil {
		return err
	}

	if ff != nil && s.isZipFile(f.info.Name()) {
		f.BaseFile = ff.Base()

		// scan zip files with a different context that is not cancellable
		// cancelling while scanning zip file contents results in the scan
		// contents being partially completed
		zipCtx := context.WithoutCancel(ctx)

		if err := s.scanZipFile(zipCtx, f); err != nil {
			logger.Errorf("Error scanning zip file %q: %v", f.Path, err)
		}
	}

	return nil
}

func (s *scanJob) isZipFile(path string) bool {
	fExt := filepath.Ext(path)
	for _, ext := range s.options.ZipFileExtensions {
		if strings.EqualFold(fExt, "."+ext) {
			return true
		}
	}

	return false
}

func (s *scanJob) onNewFile(ctx context.Context, f scanFile) (models.File, error) {
	now := time.Now()

	baseFile := f.BaseFile
	path := baseFile.Path

	baseFile.CreatedAt = now
	baseFile.UpdatedAt = now

	// find the parent folder
	parentFolderID, err := s.getFolderID(ctx, filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("getting parent folder for %q: %w", path, err)
	}

	if parentFolderID == nil {
		// if parent folder doesn't exist, assume it's not yet created
		// add this file to the queue to be created later
		if s.retrying {
			// if we're retrying and the folder still doesn't exist, then it's a problem
			return nil, fmt.Errorf("parent folder for %q doesn't exist", path)
		}

		s.retryList = append(s.retryList, f)
		return nil, nil
	}

	baseFile.ParentFolderID = *parentFolderID

	const useExisting = false
	fp, err := s.calculateFingerprints(f.fs, baseFile, path, useExisting)
	if err != nil {
		return nil, err
	}

	baseFile.SetFingerprints(fp)

	file, err := s.fireDecorators(ctx, f.fs, baseFile)
	if err != nil {
		return nil, err
	}

	// determine if the file is renamed from an existing file in the store
	// do this after decoration so that missing fields can be populated
	renamed, err := s.handleRename(ctx, file, fp)
	if err != nil {
		return nil, err
	}

	if renamed != nil {
		// handle rename should have already handled the contents of the zip file
		// so shouldn't need to scan it again
		// return nil so it doesn't
		return nil, nil
	}

	// if not renamed, queue file for creation
	if err := s.withTxn(ctx, func(ctx context.Context) error {
		if err := s.Repository.File.Create(ctx, file); err != nil {
			return fmt.Errorf("creating file %q: %w", path, err)
		}

		if err := s.fireHandlers(ctx, file, nil); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *scanJob) fireDecorators(ctx context.Context, fs models.FS, f models.File) (models.File, error) {
	for _, h := range s.FileDecorators {
		var err error
		f, err = h.Decorate(ctx, fs, f)
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

func (s *scanJob) fireHandlers(ctx context.Context, f models.File, oldFile models.File) error {
	for _, h := range s.handlers {
		if err := h.Handle(ctx, f, oldFile); err != nil {
			return err
		}
	}

	return nil
}

func (s *scanJob) calculateFingerprints(fs models.FS, f *models.BaseFile, path string, useExisting bool) (models.Fingerprints, error) {
	// only log if we're (re)calculating fingerprints
	if !useExisting {
		logger.Infof("Calculating fingerprints for %s ...", path)
	}

	// calculate primary fingerprint for the file
	fp, err := s.FingerprintCalculator.CalculateFingerprints(f, &fsOpener{
		fs:   fs,
		name: path,
	}, useExisting)
	if err != nil {
		return nil, fmt.Errorf("calculating fingerprint for file %q: %w", path, err)
	}

	return fp, nil
}

func appendFileUnique(v []models.File, toAdd []models.File) []models.File {
	for _, f := range toAdd {
		found := false
		id := f.Base().ID
		for _, vv := range v {
			if vv.Base().ID == id {
				found = true
				break
			}
		}

		if !found {
			v = append(v, f)
		}
	}

	return v
}

func (s *scanJob) getFileFS(f *models.BaseFile) (models.FS, error) {
	if f.ZipFile == nil {
		return s.FS, nil
	}

	fs, err := s.getFileFS(f.ZipFile.Base())
	if err != nil {
		return nil, err
	}

	zipPath := f.ZipFile.Base().Path
	zipSize := f.ZipFile.Base().Size
	return fs.OpenZip(zipPath, zipSize)
}

func (s *scanJob) handleRename(ctx context.Context, f models.File, fp []models.Fingerprint) (models.File, error) {
	var others []models.File

	for _, tfp := range fp {
		thisOthers, err := s.Repository.File.FindByFingerprint(ctx, tfp)
		if err != nil {
			return nil, fmt.Errorf("getting files by fingerprint %v: %w", tfp, err)
		}

		others = appendFileUnique(others, thisOthers)
	}

	var missing []models.File

	fZipID := f.Base().ZipFileID
	for _, other := range others {
		// if file is from a zip file, then only rename if both files are from the same zip file
		otherZipID := other.Base().ZipFileID
		if otherZipID != nil && (fZipID == nil || *otherZipID != *fZipID) {
			continue
		}

		// if file does not exist, then update it to the new path
		fs, err := s.getFileFS(other.Base())
		if err != nil {
			missing = append(missing, other)
			continue
		}

		info, err := fs.Lstat(other.Base().Path)
		switch {
		case err != nil:
			missing = append(missing, other)
		case strings.EqualFold(f.Base().Path, other.Base().Path):
			// #1426 - if file exists but is a case-insensitive match for the
			// original filename, and the filesystem is case-insensitive
			// then treat it as a move
			// #6326 - this should now be handled earlier, and this shouldn't be necessary
			if caseSensitive, _ := fs.IsPathCaseSensitive(other.Base().Path); !caseSensitive {
				// treat as a move
				missing = append(missing, other)
			}
		case !s.acceptEntry(ctx, other.Base().Path, info):
			// #4393 - if the file is no longer in the configured library paths, treat it as a move
			logger.Debugf("File %q no longer in library paths. Treating as a move.", other.Base().Path)
			missing = append(missing, other)
		}
	}

	n := len(missing)
	if n == 0 {
		// no missing files, not a rename
		return nil, nil
	}

	// assume does not exist, update existing file
	// it's possible that there may be multiple missing files.
	// just use the first one to rename.
	// #4775 - using the new file instance means that any changes made to the existing
	// file will be lost. Update the existing file instead.
	other := missing[0]
	updated := other.Clone()
	updatedBase := updated.Base()

	fBaseCopy := *(f.Base())

	oldPath := updatedBase.Path
	newPath := fBaseCopy.Path

	logger.Infof("%s moved to %s. Updating path...", oldPath, newPath)
	fBaseCopy.ID = updatedBase.ID
	fBaseCopy.CreatedAt = updatedBase.CreatedAt
	fBaseCopy.Fingerprints = updatedBase.Fingerprints
	*updatedBase = fBaseCopy

	if err := s.withTxn(ctx, func(ctx context.Context) error {
		if err := s.Repository.File.Update(ctx, updated); err != nil {
			return fmt.Errorf("updating file for rename %q: %w", newPath, err)
		}

		if s.isZipFile(updatedBase.Basename) {
			if err := transferZipHierarchy(ctx, s.Repository.Folder, s.Repository.File, updatedBase.ID, oldPath, newPath); err != nil {
				return fmt.Errorf("moving zip hierarchy for renamed zip file %q: %w", newPath, err)
			}
		}

		if err := s.fireHandlers(ctx, updated, other); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *scanJob) isHandlerRequired(ctx context.Context, f models.File) bool {
	accept := len(s.options.HandlerRequiredFilters) == 0
	for _, filter := range s.options.HandlerRequiredFilters {
		// accept if any filter accepts the file
		if filter.Accept(ctx, f) {
			accept = true
			break
		}
	}

	return accept
}

// isMissingMetadata returns true if the provided file is missing metadata.
// Missing metadata should only occur after the 32 schema migration.
// Looks for special values. For numbers, this will be -1. For strings, this
// will be 'unset'.
// Missing metadata includes the following:
// - file size
// - image format, width or height
// - video codec, audio codec, format, width, height, framerate or bitrate
func (s *scanJob) isMissingMetadata(ctx context.Context, f scanFile, existing models.File) bool {
	for _, h := range s.FileDecorators {
		if h.IsMissingMetadata(ctx, f.fs, existing) {
			return true
		}
	}

	return false
}

func (s *scanJob) setMissingMetadata(ctx context.Context, f scanFile, existing models.File) (models.File, error) {
	path := existing.Base().Path
	logger.Infof("Updating metadata for %s", path)

	existing.Base().Size = f.Size

	var err error
	existing, err = s.fireDecorators(ctx, f.fs, existing)
	if err != nil {
		return nil, err
	}

	// queue file for update
	if err := s.withTxn(ctx, func(ctx context.Context) error {
		if err := s.Repository.File.Update(ctx, existing); err != nil {
			return fmt.Errorf("updating file %q: %w", path, err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *scanJob) setMissingFingerprints(ctx context.Context, f scanFile, existing models.File) (models.File, error) {
	const useExisting = true
	fp, err := s.calculateFingerprints(f.fs, existing.Base(), f.Path, useExisting)
	if err != nil {
		return nil, err
	}

	if fp.ContentsChanged(existing.Base().Fingerprints) {
		existing.SetFingerprints(fp)

		if err := s.withTxn(ctx, func(ctx context.Context) error {
			if err := s.Repository.File.Update(ctx, existing); err != nil {
				return fmt.Errorf("updating file %q: %w", f.Path, err)
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return existing, nil
}

// returns a file only if it was updated
func (s *scanJob) onExistingFile(ctx context.Context, f scanFile, existing models.File) (models.File, error) {
	base := existing.Base()
	path := base.Path

	fileModTime := f.ModTime
	// #6326 - also force a rescan if the basename changed
	updated := !fileModTime.Equal(base.ModTime) || base.Basename != f.Basename
	forceRescan := s.options.Rescan

	if !updated && !forceRescan {
		return s.onUnchangedFile(ctx, f, existing)
	}

	oldBase := *base

	if !updated && forceRescan {
		logger.Infof("rescanning %s", path)
	} else {
		logger.Infof("%s has been updated: rescanning", path)
	}

	// #6326 - update basename in case it changed
	base.Basename = f.Basename
	base.ModTime = fileModTime
	base.Size = f.Size
	base.UpdatedAt = time.Now()

	// calculate and update fingerprints for the file
	const useExisting = false
	fp, err := s.calculateFingerprints(f.fs, base, path, useExisting)
	if err != nil {
		return nil, err
	}

	s.removeOutdatedFingerprints(existing, fp)
	existing.SetFingerprints(fp)

	existing, err = s.fireDecorators(ctx, f.fs, existing)
	if err != nil {
		return nil, err
	}

	// queue file for update
	if err := s.withTxn(ctx, func(ctx context.Context) error {
		if err := s.Repository.File.Update(ctx, existing); err != nil {
			return fmt.Errorf("updating file %q: %w", path, err)
		}

		if err := s.fireHandlers(ctx, existing, &oldBase); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *scanJob) removeOutdatedFingerprints(existing models.File, fp models.Fingerprints) {
	// HACK - if no MD5 fingerprint was returned, and the oshash is changed
	// then remove the MD5 fingerprint
	oshash := fp.For(models.FingerprintTypeOshash)
	if oshash == nil {
		return
	}

	existingOshash := existing.Base().Fingerprints.For(models.FingerprintTypeOshash)
	if existingOshash == nil || *existingOshash == *oshash {
		// missing oshash or same oshash - nothing to do
		return
	}

	md5 := fp.For(models.FingerprintTypeMD5)

	if md5 != nil {
		// nothing to do
		return
	}

	// oshash has changed, MD5 is missing - remove MD5 from the existing fingerprints
	logger.Infof("Removing outdated checksum from %s", existing.Base().Path)
	b := existing.Base()
	b.Fingerprints = b.Fingerprints.Remove(models.FingerprintTypeMD5)
}

// returns a file only if it was updated
func (s *scanJob) onUnchangedFile(ctx context.Context, f scanFile, existing models.File) (models.File, error) {
	var err error

	isMissingMetdata := s.isMissingMetadata(ctx, f, existing)
	// set missing information
	if isMissingMetdata {
		existing, err = s.setMissingMetadata(ctx, f, existing)
		if err != nil {
			return nil, err
		}
	}

	// calculate missing fingerprints
	existing, err = s.setMissingFingerprints(ctx, f, existing)
	if err != nil {
		return nil, err
	}

	handlerRequired := false
	if err := s.withDB(ctx, func(ctx context.Context) error {
		// check if the handler needs to be run
		handlerRequired = s.isHandlerRequired(ctx, existing)
		return nil
	}); err != nil {
		return nil, err
	}

	if !handlerRequired {
		// if this file is a zip file, then we need to rescan the contents
		// as well. We do this by returning the file, instead of nil.
		if isMissingMetdata {
			return existing, nil
		}

		return nil, nil
	}

	if err := s.withTxn(ctx, func(ctx context.Context) error {
		if err := s.fireHandlers(ctx, existing, nil); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// if this file is a zip file, then we need to rescan the contents
	// as well. We do this by returning the file, instead of nil.
	return existing, nil
}
