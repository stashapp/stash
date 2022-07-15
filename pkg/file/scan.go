package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/txn"
)

const scanQueueSize = 200000

// Repository provides access to storage methods for files and folders.
type Repository struct {
	txn.Manager
	txn.DatabaseProvider
	Store

	FolderStore FolderStore
}

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
	FS                    FS
	Repository            Repository
	FingerprintCalculator FingerprintCalculator

	// FileDecorators are applied to files as they are scanned.
	FileDecorators []Decorator
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
	dbQueue        chan func(ctx context.Context) error
	retryList      []scanFile
	retrying       bool
	folderPathToID sync.Map
	zipPathToID    sync.Map
	count          int

	txnMutex sync.Mutex
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
}

// Scan starts the scanning process.
func (s *Scanner) Scan(ctx context.Context, handlers []Handler, options ScanOptions, progressReporter ProgressReporter) {
	job := &scanJob{
		Scanner:         s,
		handlers:        handlers,
		ProgressReports: progressReporter,
		options:         options,
	}

	job.execute(ctx)
}

type scanFile struct {
	*BaseFile
	fs      FS
	info    fs.FileInfo
	zipFile *scanFile
}

func (s *scanJob) withTxn(ctx context.Context, fn func(ctx context.Context) error) error {
	// get exclusive access to the database
	s.txnMutex.Lock()
	defer s.txnMutex.Unlock()
	return txn.WithTxn(ctx, s.Repository, fn)
}

func (s *scanJob) withDB(ctx context.Context, fn func(ctx context.Context) error) error {
	return txn.WithDatabase(ctx, s.Repository, fn)
}

func (s *scanJob) execute(ctx context.Context) {
	paths := s.options.Paths
	logger.Infof("scanning %d paths", len(paths))
	s.startTime = time.Now()

	s.fileQueue = make(chan scanFile, scanQueueSize)
	s.dbQueue = make(chan func(ctx context.Context) error, scanQueueSize)

	go func() {
		if err := s.queueFiles(ctx, paths); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			logger.Errorf("error queuing files for scan: %v", err)
			return
		}

		logger.Infof("Finished adding files to queue. %d files queued", s.count)
	}()

	done := make(chan struct{}, 1)

	go func() {
		if err := s.processDBOperations(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			logger.Errorf("error processing database operations for scan: %v", err)
		}

		close(done)
	}()

	if err := s.processQueue(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		logger.Errorf("error scanning files: %v", err)
		return
	}

	// wait for database operations to complete
	<-done
}

func (s *scanJob) queueFiles(ctx context.Context, paths []string) error {
	var err error
	s.ProgressReports.ExecuteTask("Walking directory tree", func() {
		for _, p := range paths {
			err = symWalk(s.FS, p, s.queueFileFunc(ctx, s.FS, nil))
			if err != nil {
				return
			}
		}
	})

	close(s.fileQueue)

	if s.ProgressReports != nil {
		s.ProgressReports.AddTotal(s.count)
		s.ProgressReports.Definite()
	}

	return err
}

func (s *scanJob) queueFileFunc(ctx context.Context, f FS, zipFile *scanFile) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("reading info for %q: %w", path, err)
		}

		if !s.acceptEntry(ctx, path, info) {
			if info.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		ff := scanFile{
			BaseFile: &BaseFile{
				DirEntry: DirEntry{
					ModTime: modTime(info),
				},
				Path:     path,
				Basename: filepath.Base(path),
				Size:     info.Size(),
			},
			fs:   f,
			info: info,
			// there is no guarantee that the zip file has been scanned
			// so we can't just plug in the id.
			zipFile: zipFile,
		}

		if info.IsDir() {
			// handle folders immediately
			if err := s.handleFolder(ctx, ff); err != nil {
				logger.Errorf("error processing %q: %v", path, err)
				// skip the directory since we won't be able to process the files anyway
				return fs.SkipDir
			}

			return nil
		}

		// if zip file is present, we handle immediately
		if zipFile != nil {
			s.ProgressReports.ExecuteTask("Scanning "+path, func() {
				if err := s.handleFile(ctx, ff); err != nil {
					logger.Errorf("error processing %q: %v", path, err)
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
	zipFS, err := f.fs.OpenZip(f.Path)
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

	wg.Wait()
	s.retrying = true
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

	wg.Wait()

	close(s.dbQueue)

	return nil
}

func (s *scanJob) incrementProgress() {
	if s.ProgressReports != nil {
		s.ProgressReports.Increment()
	}
}

func (s *scanJob) processDBOperations(ctx context.Context) error {
	for fn := range s.dbQueue {
		if err := ctx.Err(); err != nil {
			return err
		}

		_ = s.withTxn(ctx, fn)
	}

	return nil
}

func (s *scanJob) processQueueItem(ctx context.Context, f scanFile) {
	s.ProgressReports.ExecuteTask("Scanning "+f.Path, func() {
		var err error
		if f.info.IsDir() {
			err = s.handleFolder(ctx, f)
		} else {
			err = s.handleFile(ctx, f)
		}

		if err != nil {
			logger.Errorf("error processing %q: %v", f.Path, err)
		}
	})
}

func (s *scanJob) getFolderID(ctx context.Context, path string) (*FolderID, error) {
	// check the folder cache first
	if f, ok := s.folderPathToID.Load(path); ok {
		v := f.(FolderID)
		return &v, nil
	}

	ret, err := s.Repository.FolderStore.FindByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, nil
	}

	s.folderPathToID.Store(path, ret.ID)
	return &ret.ID, nil
}

func (s *scanJob) getZipFileID(ctx context.Context, zipFile *scanFile) (*ID, error) {
	if zipFile == nil {
		return nil, nil
	}

	if zipFile.ID != 0 {
		return &zipFile.ID, nil
	}

	path := zipFile.Path

	// check the folder cache first
	if f, ok := s.zipPathToID.Load(path); ok {
		v := f.(ID)
		return &v, nil
	}

	ret, err := s.Repository.FindByPath(ctx, path)
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
		defer s.incrementProgress()

		// determine if folder already exists in data store (by path)
		f, err := s.Repository.FolderStore.FindByPath(ctx, path)
		if err != nil {
			return fmt.Errorf("checking for existing folder %q: %w", path, err)
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

func (s *scanJob) onNewFolder(ctx context.Context, file scanFile) (*Folder, error) {
	now := time.Now()

	toCreate := &Folder{
		DirEntry: DirEntry{
			ModTime: file.ModTime,
		},
		Path:      file.Path,
		CreatedAt: now,
		UpdatedAt: now,
	}

	zipFileID, err := s.getZipFileID(ctx, file.zipFile)
	if err != nil {
		return nil, err
	}

	if zipFileID != nil {
		toCreate.ZipFileID = zipFileID
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

	logger.Infof("%s doesn't exist. Creating new folder entry...", file.Path)
	if err := s.Repository.FolderStore.Create(ctx, toCreate); err != nil {
		return nil, fmt.Errorf("creating folder %q: %w", file.Path, err)
	}

	return toCreate, nil
}

func (s *scanJob) onExistingFolder(ctx context.Context, f scanFile, existing *Folder) (*Folder, error) {
	// check if the mod time is changed
	entryModTime := f.ModTime

	if !entryModTime.Equal(existing.ModTime) {
		// update entry in store
		existing.ModTime = entryModTime

		var err error
		if err = s.Repository.FolderStore.Update(ctx, existing); err != nil {
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
	var ff File
	// don't use a transaction to check if new or existing
	if err := s.withDB(ctx, func(ctx context.Context) error {
		// determine if file already exists in data store
		var err error
		ff, err = s.Repository.FindByPath(ctx, f.Path)
		if err != nil {
			return fmt.Errorf("checking for existing file %q: %w", f.Path, err)
		}

		if ff == nil {
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
		if err := s.scanZipFile(ctx, f); err != nil {
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

func (s *scanJob) onNewFile(ctx context.Context, f scanFile) (File, error) {
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
			s.incrementProgress()
			return nil, fmt.Errorf("parent folder for %q doesn't exist", path)
		}

		s.retryList = append(s.retryList, f)
		return nil, nil
	}

	baseFile.ParentFolderID = *parentFolderID

	zipFileID, err := s.getZipFileID(ctx, f.zipFile)
	if err != nil {
		s.incrementProgress()
		return nil, err
	}

	if zipFileID != nil {
		baseFile.ZipFileID = zipFileID
	}

	fp, err := s.calculateFingerprints(f.fs, baseFile, path)
	if err != nil {
		s.incrementProgress()
		return nil, err
	}

	baseFile.SetFingerprints(fp)

	file, err := s.fireDecorators(ctx, f.fs, baseFile)
	if err != nil {
		s.incrementProgress()
		return nil, err
	}

	// determine if the file is renamed from an existing file in the store
	// do this after decoration so that missing fields can be populated
	renamed, err := s.handleRename(ctx, file, fp)
	if err != nil {
		s.incrementProgress()
		return nil, err
	}

	if renamed != nil {
		return renamed, nil
	}

	// if not renamed, queue file for creation
	if err := s.queueDBOperation(ctx, path, func(ctx context.Context) error {
		if err := s.Repository.Create(ctx, file); err != nil {
			return fmt.Errorf("creating file %q: %w", path, err)
		}

		if err := s.fireHandlers(ctx, file); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *scanJob) queueDBOperation(ctx context.Context, path string, fn func(ctx context.Context) error) error {
	// perform immediately if it is a zip file
	if s.isZipFile(path) {
		return s.withTxn(ctx, fn)
	}

	s.dbQueue <- fn

	return nil
}

func (s *scanJob) fireDecorators(ctx context.Context, fs FS, f File) (File, error) {
	for _, h := range s.FileDecorators {
		var err error
		f, err = h.Decorate(ctx, fs, f)
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

func (s *scanJob) fireHandlers(ctx context.Context, f File) error {
	for _, h := range s.handlers {
		if err := h.Handle(ctx, f); err != nil {
			return err
		}
	}

	return nil
}

func (s *scanJob) calculateFingerprints(fs FS, f *BaseFile, path string) ([]Fingerprint, error) {
	logger.Infof("Calculating fingerprints for %s ...", path)

	// calculate primary fingerprint for the file
	fp, err := s.FingerprintCalculator.CalculateFingerprints(f, &fsOpener{
		fs:   fs,
		name: path,
	})
	if err != nil {
		return nil, fmt.Errorf("calculating fingerprint for file %q: %w", path, err)
	}

	return fp, nil
}

func appendFileUnique(v []File, toAdd []File) []File {
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

func (s *scanJob) getFileFS(f *BaseFile) (FS, error) {
	if f.ZipFile == nil {
		return s.FS, nil
	}

	fs, err := s.getFileFS(f.ZipFile.Base())
	if err != nil {
		return nil, err
	}

	zipPath := f.ZipFile.Base().Path
	return fs.OpenZip(zipPath)
}

func (s *scanJob) handleRename(ctx context.Context, f File, fp []Fingerprint) (File, error) {
	var others []File

	for _, tfp := range fp {
		thisOthers, err := s.Repository.FindByFingerprint(ctx, tfp)
		if err != nil {
			return nil, fmt.Errorf("getting files by fingerprint %v: %w", tfp, err)
		}

		others = appendFileUnique(others, thisOthers)
	}

	var missing []File

	for _, other := range others {
		// if file does not exist, then update it to the new path
		// TODO - handle #1426 scenario
		fs, err := s.getFileFS(other.Base())
		if err != nil {
			return nil, fmt.Errorf("getting FS for %q: %w", other.Base().Path, err)
		}

		if _, err := fs.Lstat(other.Base().Path); err != nil {
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
	other := missing[0]
	otherBase := other.Base()

	fBase := f.Base()

	logger.Infof("%s moved to %s. Updating path...", otherBase.Path, fBase.Path)
	fBase.ID = otherBase.ID
	fBase.CreatedAt = otherBase.CreatedAt
	fBase.Fingerprints = otherBase.Fingerprints

	if err := s.queueDBOperation(ctx, fBase.Path, func(ctx context.Context) error {
		if err := s.Repository.Update(ctx, f); err != nil {
			return fmt.Errorf("updating file for rename %q: %w", fBase.Path, err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return f, nil
}

func (s *scanJob) isHandlerRequired(ctx context.Context, f File) bool {
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

// returns a file only if it was updated
func (s *scanJob) onExistingFile(ctx context.Context, f scanFile, existing File) (File, error) {
	base := existing.Base()
	path := base.Path

	fileModTime := f.ModTime
	updated := !fileModTime.Equal(base.ModTime)

	if !updated {
		handlerRequired := false
		if err := s.withDB(ctx, func(ctx context.Context) error {
			// check if the handler needs to be run
			handlerRequired = s.isHandlerRequired(ctx, existing)
			return nil
		}); err != nil {
			return nil, err
		}

		if !handlerRequired {
			s.incrementProgress()
			return nil, nil
		}

		if err := s.queueDBOperation(ctx, path, func(ctx context.Context) error {
			if err := s.fireHandlers(ctx, existing); err != nil {
				return err
			}

			s.incrementProgress()
			return nil
		}); err != nil {
			return nil, err
		}

		return nil, nil
	}

	logger.Infof("%s has been updated: rescanning", path)
	base.ModTime = fileModTime
	base.Size = f.Size
	base.UpdatedAt = time.Now()

	// calculate and update fingerprints for the file
	fp, err := s.calculateFingerprints(f.fs, base, path)
	if err != nil {
		s.incrementProgress()
		return nil, err
	}

	existing.SetFingerprints(fp)

	existing, err = s.fireDecorators(ctx, f.fs, existing)
	if err != nil {
		s.incrementProgress()
		return nil, err
	}

	// queue file for update
	if err := s.queueDBOperation(ctx, path, func(ctx context.Context) error {
		if err := s.Repository.Update(ctx, existing); err != nil {
			return fmt.Errorf("updating file %q: %w", path, err)
		}

		if err := s.fireHandlers(ctx, existing); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return existing, nil
}
