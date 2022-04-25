package file

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/txn"
)

const scanQueueSize = 200000

type Repository struct {
	txn.Manager
	Store
	MissedMarker

	FolderStore FolderStore
}

type Scanner struct {
	FS FS

	Repository          Repository
	FingerprintHandlers FingerprintHandler
	ScanFilters         []Filter
	Handlers            []Handler
	ProgressReports     ScanProgressReporter

	startTime      time.Time
	fileQueue      chan scanFile
	retryList      []scanFile
	retrying       bool
	folderPathToID sync.Map
	zipPathToID    sync.Map
	count          int

	txnMutex sync.Mutex
}

type ScanProgressReporter interface {
	SetTotal(total int)
	Increment()
}

type scanFile struct {
	*BasicFile
	fs      FS
	info    fs.FileInfo
	zipFile *scanFile
}

func (s *Scanner) withTxn(ctx context.Context, fn func(ctx context.Context) error) error {
	// get exclusive access to the database
	s.txnMutex.Lock()
	defer s.txnMutex.Unlock()
	return txn.WithTxn(ctx, s.Repository, fn)
}

func (s *Scanner) Scan(ctx context.Context, paths []string) {
	logger.Infof("scanning %d paths", len(paths))
	s.startTime = time.Now()

	s.fileQueue = make(chan scanFile, scanQueueSize)

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

	if err := s.processQueue(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		logger.Errorf("error scanning files: %v", err)
		return
	}

	// now mark any files not seen as missing
	s.markMissingFiles(ctx)

	logger.Info("Finished scanning files")
}

func (s *Scanner) queueFiles(ctx context.Context, paths []string) error {
	var err error
	for _, p := range paths {
		err = symWalk(s.FS, p, s.queueFileFunc(ctx, s.FS, nil))
	}
	close(s.fileQueue)

	if err != nil {
		return err
	}

	if s.ProgressReports != nil {
		s.ProgressReports.SetTotal(s.count)
	}

	return nil
}

func (s *Scanner) queueFileFunc(ctx context.Context, f FS, zipFile *scanFile) fs.WalkDirFunc {
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

		ff := scanFile{
			BasicFile: &BasicFile{
				DirEntry: DirEntry{
					Path:        path,
					ModTime:     modTime(info),
					LastScanned: time.Now(),
				},
				Basename: filepath.Base(path),
				Size:     info.Size(),
			},
			fs:   f,
			info: info,
			// there is no guarantee that the zip file has been scanned
			// so we can't just plug in the id.
			zipFile: zipFile,
		}

		// always accept if there's no filters
		accept := len(s.ScanFilters) == 0
		for _, filter := range s.ScanFilters {
			// accept if any filter accepts the file
			if filter.Accept(ff.BasicFile) {
				accept = true
				break
			}
		}

		if !accept {
			return fs.SkipDir
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
			if err := s.handleFile(ctx, ff); err != nil {
				logger.Errorf("error processing %q: %v", path, err)
				// don't return an error, just skip the file
			}

			return nil
		}

		s.fileQueue <- ff

		s.count++

		return nil
	}
}

func (s *Scanner) scanZipFile(ctx context.Context, f scanFile) error {
	reader, err := f.fs.Open(f.Path)
	if err != nil {
		return err
	}
	defer reader.Close()

	asReaderAt, _ := reader.(io.ReaderAt)
	if asReaderAt == nil {
		// can't walk the zip file
		// just return
		return nil
	}

	zipReader, err := zip.NewReader(asReaderAt, f.Size)
	if err != nil {
		return err
	}

	zipFS := &ZipFS{
		Reader:  zipReader,
		zipInfo: f.info,
		zipPath: f.Path,
	}

	return symWalk(zipFS, f.Path, s.queueFileFunc(ctx, zipFS, &f))
}

func (s *Scanner) processQueue(ctx context.Context) error {
	for f := range s.fileQueue {
		if err := ctx.Err(); err != nil {
			return err
		}

		s.processQueueItem(ctx, f)
	}

	s.retrying = true
	for _, f := range s.retryList {
		if err := ctx.Err(); err != nil {
			return err
		}

		s.processQueueItem(ctx, f)
	}

	return nil
}

func (s *Scanner) processQueueItem(ctx context.Context, f scanFile) {
	var err error
	if f.info.IsDir() {
		err = s.handleFolder(ctx, f)
	} else {
		err = s.handleFile(ctx, f)
	}

	if err != nil {
		logger.Errorf("error processing %q: %v", f.Path, err)
	}
}

func (s *Scanner) getFolderID(ctx context.Context, path string) (*FolderID, error) {
	// check the folder cache first
	if f, ok := s.folderPathToID.Load(path); ok {
		v := f.(FolderID)
		return &v, nil
	}

	ret, err := s.Repository.FolderStore.GetByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, nil
	}

	s.folderPathToID.Store(path, ret.ID)
	return &ret.ID, nil
}

func (s *Scanner) getZipFileID(ctx context.Context, zipFile *scanFile) (*FileID, error) {
	if zipFile == nil {
		return nil, nil
	}

	path := zipFile.Path

	// check the folder cache first
	if f, ok := s.zipPathToID.Load(path); ok {
		v := f.(FileID)
		return &v, nil
	}

	ret, err := s.Repository.GetByPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("getting zip file ID for %q: %w", path, err)
	}

	if ret == nil {
		return nil, fmt.Errorf("zip file %q doesn't exist in database", zipFile.Path)
	}

	s.zipPathToID.Store(path, ret.Basic().ID)
	return &ret.Basic().ID, nil
}

func (s *Scanner) handleFolder(ctx context.Context, file scanFile) error {
	path := file.Path

	return s.withTxn(ctx, func(ctx context.Context) error {
		// determine if folder already exists in data store (by path)
		f, err := s.Repository.FolderStore.GetByPath(ctx, path)
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

func (s *Scanner) onNewFolder(ctx context.Context, file scanFile) (*Folder, error) {
	now := time.Now()

	toCreate := Folder{
		DirEntry: DirEntry{
			Path:    file.Path,
			ModTime: file.ModTime,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	zipFileID, err := s.getZipFileID(ctx, file.zipFile)
	if err != nil {
		return nil, err
	}

	if zipFileID != nil {
		file.ZipFileID = zipFileID
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
	f, err := s.Repository.FolderStore.Create(ctx, toCreate)
	if err != nil {
		return nil, fmt.Errorf("creating folder %q: %w", file.Path, err)
	}

	return f, nil
}

func (s *Scanner) onExistingFolder(ctx context.Context, f scanFile, existing *Folder) (*Folder, error) {
	// check if the mod time is changed
	entryModTime := f.ModTime

	if !entryModTime.Equal(existing.ModTime) {
		// update entry in store
		existing.ModTime = entryModTime
	}

	existing.scanned()

	var err error
	_, err = s.Repository.FolderStore.Update(ctx, *existing)
	if err != nil {
		return nil, fmt.Errorf("updating folder %q: %w", f.Path, err)
	}

	return existing, nil
}

func modTime(info fs.FileInfo) time.Time {
	// truncate to seconds, since we don't store beyond that in the database
	return info.ModTime().Truncate(time.Second)
}

func (s *Scanner) handleFile(ctx context.Context, f scanFile) error {
	// TODO - ensure file should be included

	var ff File
	if err := s.withTxn(ctx, func(ctx context.Context) error {
		// determine if file already exists in data store
		var err error
		ff, err = s.Repository.GetByPath(ctx, f.Path)
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

	if ff != nil && s.isZipFile(f.info) {
		if err := s.scanZipFile(ctx, f); err != nil {
			logger.Errorf("Error scanning zip file %q: %v", f.Path, err)
		}
	}

	return nil
}

func (s *Scanner) isZipFile(entry fs.FileInfo) bool {
	// TODO - this should be configurable
	return strings.HasSuffix(entry.Name(), ".zip")
}

func (s *Scanner) onNewFile(ctx context.Context, f scanFile) (*BasicFile, error) {
	now := time.Now()

	file := f.BasicFile
	path := file.Path

	file.CreatedAt = now
	file.UpdatedAt = now

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

	file.ParentFolderID = parentFolderID

	zipFileID, err := s.getZipFileID(ctx, f.zipFile)
	if err != nil {
		return nil, err
	}

	if zipFileID != nil {
		file.ZipFileID = zipFileID
	}

	fp, err := s.calculateFingerprint(f.fs, file, path)
	if err != nil {
		return nil, err
	}

	file.SetFingerprint(*fp)

	// determine if the file is renamed from an existing file in the store
	renamed, err := s.handleRename(ctx, file, fp)
	if err != nil {
		return nil, err
	}

	if renamed {
		return file, nil
	}

	// if not renamed, add file to store
	logger.Infof("%s doesn't exist. Creating new file entry...", path)
	if err := s.Repository.Create(ctx, file); err != nil {
		return nil, fmt.Errorf("creating file %q: %w", path, err)
	}

	if err := s.fireHandlers(ctx, f.fs, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *Scanner) fireHandlers(ctx context.Context, fs FS, f File) error {
	for _, h := range s.Handlers {
		if err := h.Handle(ctx, fs, f); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scanner) calculateFingerprint(fs FS, f *BasicFile, path string) (*Fingerprint, error) {
	logger.Infof("Calculating fingerprint for %s ...", path)

	// calculate primary fingerprint for the file
	r, err := fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer r.Close()

	fp, err := s.FingerprintHandlers.CalculateFingerprint(f, r)
	if err != nil {
		return nil, fmt.Errorf("calculating fingerprint for file %q: %w", path, err)
	}

	return fp, nil
}

func (s *Scanner) handleRename(ctx context.Context, f *BasicFile, fp *Fingerprint) (bool, error) {
	others, err := s.Repository.GetByFingerprint(ctx, *fp)
	if err != nil {
		return false, fmt.Errorf("getting files by fingerprint %q: %w", fp, err)
	}

	var missing []File

	for _, other := range others {
		// if file does not exist, then update it to the new path
		// TODO - this could be in a zip file
		if _, err := s.FS.Lstat(other.Basic().Path); err != nil {
			missing = append(missing, other)
		}
	}

	n := len(missing)
	switch {
	case n == 1:
		// assume does not exist, update existing file
		other := missing[0].Basic()

		logger.Infof("%s moved to %s. Updating path...", other.Path, f.Path)
		f.ID = other.ID
		f.CreatedAt = other.CreatedAt
		f.Fingerprints = other.Fingerprints
		if err := s.Repository.Update(ctx, f); err != nil {
			return false, fmt.Errorf("updating file for rename %q: %w", f.Path, err)
		}

		return true, nil
	case n > 1:
		// multiple candidates
		// TODO - mark all as missing and just create a new file
		return false, nil
	}

	return false, nil
}

func (s *Scanner) onExistingFile(ctx context.Context, f scanFile, existing File) (File, error) {
	base := existing.Basic()
	path := base.Path

	base.scanned()

	fileModTime := f.ModTime
	updated := !fileModTime.Equal(base.ModTime)

	if updated {
		logger.Infof("%s has been updated: rescanning", path)
		base.ModTime = fileModTime
		base.Size = f.Size
		base.UpdatedAt = time.Now()

		// calculate and update primary fingerprint for the file
		fp, err := s.calculateFingerprint(f.fs, base, path)
		if err != nil {
			return nil, err
		}

		existing.SetFingerprint(*fp)
	}

	if err := s.Repository.Update(ctx, base); err != nil {
		return nil, fmt.Errorf("updating file %q: %w", path, err)
	}

	if !updated {
		return existing, nil
	}

	if err := s.fireHandlers(ctx, f.fs, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *Scanner) markMissingFiles(ctx context.Context) {
	if err := s.withTxn(ctx, func(ctx context.Context) error {
		return s.Repository.MarkMissing(ctx, s.startTime)
	}); err != nil {
		logger.Errorf("Error marking missing files: %v", err)
	}
}
