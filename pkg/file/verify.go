package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// Verifier scans through stored file and folder instances and marks missing those that are no longer present on disk.
type Verifier struct {
	FS         models.FS
	Repository Repository

	PurgeHandlers []PurgeHandler
	TrashPath     string
}

type missingFileFn func(ctx context.Context, f models.File) (processed bool, err error)
type missingFolderFn func(ctx context.Context, f *models.Folder) (processed bool, err error)

type verifyJob struct {
	*Verifier

	progress *job.Progress
	options  VerifyOptions

	toDelete deleteSet

	missingFileHandler   missingFileFn
	missingFolderHandler missingFolderFn
}

// VerifyOptions provides options for verifying files.
type VerifyOptions struct {
	Paths []string

	// IgnoreZipFileContents will skip checking the contents of zip files when determining whether to verify a file.
	// This can significantly speed up the verify process, but will potentially miss removed files within zip files.
	// Where users do not modify zip files contents directly, this should be safe to use.
	IgnoreZipFileContents bool

	// DryRun indicates if this is a dry run. A dry run will not make any changes and will only log what would have been done.
	DryRun bool

	// PurgeMissing indicates if missing files should be purged. If true, missing files will be deleted from the database and any associated objects will be deleted as well.
	// No effect if DryRun is true.
	PurgeMissing bool

	// PathFilter are used to determine if a file should be included.
	// Excluded files are marked for cleaning.
	PathFilter PathFilter
}

// Verify starts the verify process.
func (s *Verifier) Verify(ctx context.Context, options VerifyOptions, progress *job.Progress) {
	j := &verifyJob{
		Verifier: s,
		progress: progress,
		options:  options,
	}

	if err := j.execute(ctx); err != nil {
		logger.Errorf("error verifying files: %v", err)
		return
	}
}

type fileOrFolder struct {
	fileID   models.FileID
	folderID models.FolderID
}

type deleteSet struct {
	orderedList []fileOrFolder
	fileIDSet   map[models.FileID]string

	folderIDSet map[models.FolderID]string
}

func newDeleteSet() deleteSet {
	return deleteSet{
		fileIDSet:   make(map[models.FileID]string),
		folderIDSet: make(map[models.FolderID]string),
	}
}

func (s *deleteSet) add(id models.FileID, path string) {
	if _, ok := s.fileIDSet[id]; !ok {
		s.orderedList = append(s.orderedList, fileOrFolder{fileID: id})
		s.fileIDSet[id] = path
	}
}

func (s *deleteSet) has(id models.FileID) bool {
	_, ok := s.fileIDSet[id]
	return ok
}

func (s *deleteSet) addFolder(id models.FolderID, path string) {
	if _, ok := s.folderIDSet[id]; !ok {
		s.orderedList = append(s.orderedList, fileOrFolder{folderID: id})
		s.folderIDSet[id] = path
	}
}

func (s *deleteSet) hasFolder(id models.FolderID) bool {
	_, ok := s.folderIDSet[id]
	return ok
}

func (s *deleteSet) len() int {
	return len(s.orderedList)
}

func (j *verifyJob) init() {
	j.toDelete = newDeleteSet()

	switch {
	case j.options.DryRun:
		j.missingFileHandler = j.noopMissingFile
		j.missingFolderHandler = j.noopMissingFolder
	case j.options.PurgeMissing:
		j.missingFileHandler = j.markToDeleteMissingFile
		j.missingFolderHandler = j.markToDeleteMissingFolder
	default:
		j.missingFileHandler = j.markMissingFile
		j.missingFolderHandler = j.markMissingFolder
	}
}

func (j *verifyJob) noopMissingFile(ctx context.Context, f models.File) (processed bool, err error) {
	return true, nil
}

func (j *verifyJob) markMissingFile(ctx context.Context, f models.File) (processed bool, err error) {
	missingSince := time.Now()
	if err := j.Repository.File.SetMissing(ctx, f.Base().ID, &missingSince); err != nil {
		return false, err
	}
	return true, nil
}

func (j *verifyJob) markToDeleteMissingFile(ctx context.Context, f models.File) (processed bool, err error) {
	j.toDelete.add(f.Base().ID, f.Base().Path)
	// needs to be deleted, mark as not processed
	return false, nil
}

func (j *verifyJob) noopMissingFolder(ctx context.Context, f *models.Folder) (processed bool, err error) {
	return true, nil
}

func (j *verifyJob) markMissingFolder(ctx context.Context, f *models.Folder) (processed bool, err error) {
	missingSince := time.Now()
	if err := j.Repository.Folder.SetMissing(ctx, f.ID, &missingSince); err != nil {
		return false, err
	}
	return true, nil
}

func (j *verifyJob) markToDeleteMissingFolder(ctx context.Context, f *models.Folder) (processed bool, err error) {
	j.toDelete.addFolder(f.ID, f.Path)
	// needs to be deleted, mark as not processed
	return false, nil
}

func (j *verifyJob) execute(ctx context.Context) error {
	j.init()
	progress := j.progress

	var (
		fileCount   int
		folderCount int
	)

	r := j.Repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		fileCount, err = r.File.CountAllInPaths(ctx, j.options.Paths)
		if err != nil {
			return err
		}

		folderCount, err = r.Folder.CountAllInPaths(ctx, j.options.Paths)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	progress.AddTotal(fileCount + folderCount)
	progress.Definite()

	if err := j.assessFolders(ctx); err != nil {
		return err
	}

	if err := j.assessFiles(ctx); err != nil {
		return err
	}

	if j.options.DryRun || !j.options.PurgeMissing {
		// nothing further to do
		return nil
	}

	progress.ExecuteTask(fmt.Sprintf("Purging %d files and folders", j.toDelete.len()), func() {
		for _, ff := range j.toDelete.orderedList {
			if job.IsCancelled(ctx) {
				return
			}

			if ff.fileID != 0 {
				j.deleteFile(ctx, ff.fileID, j.toDelete.fileIDSet[ff.fileID])
			}
			if ff.folderID != 0 {
				j.deleteFolder(ctx, ff.folderID, j.toDelete.folderIDSet[ff.folderID])
			}

			progress.Increment()
		}
	})

	return nil
}

func (j *verifyJob) assessFiles(ctx context.Context) error {
	const batchSize = 1000
	offset := 0
	progress := j.progress

	more := true
	r := j.Repository

	includeZipContents := !j.options.IgnoreZipFileContents

	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		var files []models.File
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			files, err = r.File.FindAllInPaths(ctx, j.options.Paths, includeZipContents, batchSize, offset)
			if err != nil {
				return fmt.Errorf("error querying for files: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		for _, f := range files {
			path := f.Base().Path
			fileID := f.Base().ID

			// short-cut, don't assess if already missing
			if f.Base().MissingSince != nil {
				// increment progress, no further processing
				progress.Increment()
				continue
			}

			// skip if already added to delete set
			// don't increment progress here, as it will be incremented when the file is processed
			if j.toDelete.has(fileID) {
				continue
			}

			var err error
			progress.ExecuteTask(fmt.Sprintf("Verifying file %s", path), func() {
				if j.shouldClean(ctx, f) {
					err = j.handleMissingFile(ctx, f)
				} else {
					// increment progress, no further processing
					progress.Increment()
				}
			})
			if err != nil {
				return err
			}
		}

		if len(files) != batchSize {
			more = false
		} else {
			offset += batchSize
		}
	}

	return nil
}

// flagFolderForDelete adds folders to the toDelete set, with the leaf folders added first
func (j *verifyJob) handleMissingFile(ctx context.Context, f models.File) error {
	r := j.Repository

	// do all this in a transaction so that all contained files are marked in a single transaction
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		// add contained files first
		containedFiles, err := r.File.FindByZipFileID(ctx, f.Base().ID)
		if err != nil {
			return fmt.Errorf("error finding contained files for %q: %w", f.Base().Path, err)
		}

		for _, cf := range containedFiles {
			logger.Infof("Marking contained file %q to clean", cf.Base().Path)
			processed, err := j.missingFileHandler(ctx, cf)
			if err != nil {
				return err
			}

			if processed {
				j.progress.Increment()
			}
		}

		// add contained folders as well
		containedFolders, err := r.Folder.FindByZipFileID(ctx, f.Base().ID)
		if err != nil {
			return fmt.Errorf("error finding contained folders for %q: %w", f.Base().Path, err)
		}

		for _, cf := range containedFolders {
			logger.Infof("Marking contained folder %q to clean", cf.Path)
			processed, err := j.missingFolderHandler(ctx, cf)
			if err != nil {
				return err
			}

			if processed {
				j.progress.Increment()
			}
		}

		processed, err := j.missingFileHandler(ctx, f)
		if err != nil {
			return err
		}

		if processed {
			j.progress.Increment()
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (j *verifyJob) assessFolders(ctx context.Context) error {
	const batchSize = 1000
	offset := 0
	progress := j.progress

	includeZipContents := !j.options.IgnoreZipFileContents

	more := true
	r := j.Repository

	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		var folders []*models.Folder
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			var err error
			folders, err = r.Folder.FindAllInPaths(ctx, j.options.Paths, includeZipContents, batchSize, offset)
			if err != nil {
				return fmt.Errorf("error querying for folders: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		for _, f := range folders {
			path := f.Path
			folderID := f.ID

			// don't assess if already missing
			if f.MissingSince != nil {
				// increment progress, no further processing
				progress.Increment()
				continue
			}

			// skip if already added to delete set
			// don't increment progress here, as it will be incremented when the folder is processed
			if j.toDelete.hasFolder(folderID) {
				continue
			}

			var err error
			progress.ExecuteTask(fmt.Sprintf("Verifying folder %s", path), func() {
				if j.shouldCleanFolder(ctx, f) {
					err = j.handleMissingFolder(ctx, f)
				} else {
					// increment progress, no further processing
					progress.Increment()
				}
			})
			if err != nil {
				return err
			}
		}

		if len(folders) != batchSize {
			more = false
		} else {
			offset += batchSize
		}
	}

	return nil
}

func (j *verifyJob) handleMissingFolder(ctx context.Context, folder *models.Folder) error {
	r := j.Repository

	// do all this in a transaction so that all contained files are marked in a single transaction
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		// add contained files first
		containedFiles, err := r.File.FindByFolderID(ctx, folder.ID)
		if err != nil {
			return fmt.Errorf("error finding contained files for %q: %w", folder.Path, err)
		}

		for _, cf := range containedFiles {
			logger.Infof("Marking contained file %q to clean", cf.Base().Path)
			processed, err := j.missingFileHandler(ctx, cf)
			if err != nil {
				return err
			}

			if processed {
				j.progress.Increment()
			}
		}

		// it is possible that child folders may be included while parent folders are not
		// so we need to check child folders separately

		// only use the processed return value for the top-level folder
		processed, err := j.missingFolderHandler(ctx, folder)
		if err != nil {
			return err
		}

		if processed {
			j.progress.Increment()
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func isNotFound(err error) bool {
	// ErrInvalid can occur in zip files where the zip file path changed
	// and the underlying folder did not
	// #3877 - fs.PathError can occur if the network share no longer exists
	var pathErr *fs.PathError
	return err != nil &&
		(errors.Is(err, fs.ErrNotExist) ||
			errors.Is(err, fs.ErrInvalid) ||
			errors.As(err, &pathErr))
}

func (j *verifyJob) shouldClean(ctx context.Context, f models.File) bool {
	path := f.Base().Path

	info, err := f.Base().Info(j.FS)
	if err != nil && !isNotFound(err) {
		logger.Errorf("error getting file info for %q, not cleaning: %v", path, err)
		return false
	}

	if info == nil {
		// info is nil - file not exist
		logger.Infof("File not found. Marking as missing: \"%s\"", path)
		return true
	}

	// run through path filter, if returns false then the file should be cleaned
	filter := j.options.PathFilter

	// need to get the zip file path if present
	zipFilePath := ""
	if f.Base().ZipFile != nil {
		zipFilePath = f.Base().ZipFile.Base().Path
	}

	// don't log anything - assume filter will have logged the reason
	return !filter.Accept(ctx, path, info, zipFilePath)
}

func (j *verifyJob) shouldCleanFolder(ctx context.Context, f *models.Folder) bool {
	path := f.Path

	info, err := f.Info(j.FS)

	if err != nil && !isNotFound(err) {
		logger.Errorf("error getting folder info for %q, not cleaning: %v", path, err)
		return false
	}

	if info == nil {
		// info is nil - file not exist
		logger.Infof("Folder not found. Marking as missing: \"%s\"", path)
		return true
	}

	// #3261 - handle symlinks
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		finalPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			// don't bail out if symlink is invalid
			logger.Infof("Invalid symlink. Marking as missing: \"%s\"", path)
			return true
		}

		info, err = j.FS.Lstat(finalPath)
		if err != nil && !isNotFound(err) {
			logger.Errorf("error getting file info for %q (-> %s), not cleaning: %v", path, finalPath, err)
			return false
		}
	}

	// run through path filter, if returns false then the file should be cleaned
	filter := j.options.PathFilter

	// need to get the zip file path if present
	zipFilePath := ""
	if f.ZipFile != nil {
		zipFilePath = f.ZipFile.Base().Path
	}

	// don't log anything - assume filter will have logged the reason
	return !filter.Accept(ctx, path, info, zipFilePath)
}

func (j *verifyJob) deleteFile(ctx context.Context, fileID models.FileID, fn string) {
	// delete associated objects
	fileDeleter := NewDeleterWithTrash(j.TrashPath)
	r := j.Repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		fileDeleter.RegisterHooks(ctx)

		if err := j.fireHandlers(ctx, fileDeleter, fileID); err != nil {
			return err
		}

		return r.File.Destroy(ctx, fileID)
	}); err != nil {
		logger.Errorf("Error deleting file %q from database: %s", fn, err.Error())
		return
	}
}

func (j *verifyJob) deleteFolder(ctx context.Context, folderID models.FolderID, fn string) {
	// delete associated objects
	fileDeleter := NewDeleterWithTrash(j.TrashPath)
	r := j.Repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		fileDeleter.RegisterHooks(ctx)

		if err := j.fireFolderHandlers(ctx, fileDeleter, folderID); err != nil {
			return err
		}

		return r.Folder.Destroy(ctx, folderID)
	}); err != nil {
		logger.Errorf("Error deleting folder %q from database: %s", fn, err.Error())
		return
	}
}

func (j *verifyJob) fireHandlers(ctx context.Context, fileDeleter *Deleter, fileID models.FileID) error {
	for _, h := range j.PurgeHandlers {
		if err := h.HandleFile(ctx, fileDeleter, fileID); err != nil {
			return err
		}
	}

	return nil
}

func (j *verifyJob) fireFolderHandlers(ctx context.Context, fileDeleter *Deleter, folderID models.FolderID) error {
	for _, h := range j.PurgeHandlers {
		if err := h.HandleFolder(ctx, fileDeleter, folderID); err != nil {
			return err
		}
	}

	return nil
}
