package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// Cleaner scans through stored file and folder instances and removes those that are no longer present on disk.
type Cleaner struct {
	FS         models.FS
	Repository Repository

	Handlers  []CleanHandler
	TrashPath string
}

type cleanJob struct {
	*Cleaner

	progress *job.Progress
	options  CleanOptions
}

// CleanOptions provides options for scanning files.
type CleanOptions struct {
	Paths []string

	// Do a dry run. Don't delete any files
	DryRun bool

	// PathFilter are used to determine if a file should be included.
	// Excluded files are marked for cleaning.
	PathFilter PathFilter
}

// Clean starts the clean process.
func (s *Cleaner) Clean(ctx context.Context, options CleanOptions, progress *job.Progress) {
	j := &cleanJob{
		Cleaner:  s,
		progress: progress,
		options:  options,
	}

	if err := j.execute(ctx); err != nil {
		logger.Errorf("error cleaning files: %v", err)
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

func (j *cleanJob) execute(ctx context.Context) error {
	progress := j.progress

	toDelete := newDeleteSet()

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

	if err := j.assessFiles(ctx, &toDelete); err != nil {
		return err
	}

	if err := j.assessFolders(ctx, &toDelete); err != nil {
		return err
	}

	if j.options.DryRun && toDelete.len() > 0 {
		// add progress for files that would've been deleted
		progress.AddProcessed(toDelete.len())
		return nil
	}

	progress.ExecuteTask(fmt.Sprintf("Cleaning %d files and folders", toDelete.len()), func() {
		for _, ff := range toDelete.orderedList {
			if job.IsCancelled(ctx) {
				return
			}

			if ff.fileID != 0 {
				j.deleteFile(ctx, ff.fileID, toDelete.fileIDSet[ff.fileID])
			}
			if ff.folderID != 0 {
				j.deleteFolder(ctx, ff.folderID, toDelete.folderIDSet[ff.folderID])
			}

			progress.Increment()
		}
	})

	return nil
}

func (j *cleanJob) assessFiles(ctx context.Context, toDelete *deleteSet) error {
	const batchSize = 1000
	offset := 0
	progress := j.progress

	more := true
	r := j.Repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		for more {
			if job.IsCancelled(ctx) {
				return nil
			}

			files, err := r.File.FindAllInPaths(ctx, j.options.Paths, batchSize, offset)
			if err != nil {
				return fmt.Errorf("error querying for files: %w", err)
			}

			for _, f := range files {
				path := f.Base().Path
				err = nil
				fileID := f.Base().ID

				// short-cut, don't assess if already added
				if toDelete.has(fileID) {
					continue
				}

				progress.ExecuteTask(fmt.Sprintf("Assessing file %s for clean", path), func() {
					if j.shouldClean(ctx, f) {
						err = j.flagFileForDelete(ctx, toDelete, f)
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
	}); err != nil {
		return err
	}

	return nil
}

// flagFolderForDelete adds folders to the toDelete set, with the leaf folders added first
func (j *cleanJob) flagFileForDelete(ctx context.Context, toDelete *deleteSet, f models.File) error {
	r := j.Repository
	// add contained files first
	containedFiles, err := r.File.FindByZipFileID(ctx, f.Base().ID)
	if err != nil {
		return fmt.Errorf("error finding contained files for %q: %w", f.Base().Path, err)
	}

	for _, cf := range containedFiles {
		logger.Infof("Marking contained file %q to clean", cf.Base().Path)
		toDelete.add(cf.Base().ID, cf.Base().Path)
	}

	// add contained folders as well
	containedFolders, err := r.Folder.FindByZipFileID(ctx, f.Base().ID)
	if err != nil {
		return fmt.Errorf("error finding contained folders for %q: %w", f.Base().Path, err)
	}

	for _, cf := range containedFolders {
		logger.Infof("Marking contained folder %q to clean", cf.Path)
		toDelete.addFolder(cf.ID, cf.Path)
	}

	toDelete.add(f.Base().ID, f.Base().Path)

	return nil
}

func (j *cleanJob) assessFolders(ctx context.Context, toDelete *deleteSet) error {
	const batchSize = 1000
	offset := 0
	progress := j.progress

	more := true
	r := j.Repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		for more {
			if job.IsCancelled(ctx) {
				return nil
			}

			folders, err := r.Folder.FindAllInPaths(ctx, j.options.Paths, batchSize, offset)
			if err != nil {
				return fmt.Errorf("error querying for folders: %w", err)
			}

			for _, f := range folders {
				path := f.Path
				folderID := f.ID

				// short-cut, don't assess if already added
				if toDelete.hasFolder(folderID) {
					continue
				}

				err = nil
				progress.ExecuteTask(fmt.Sprintf("Assessing folder %s for clean", path), func() {
					if j.shouldCleanFolder(ctx, f) {
						if err = j.flagFolderForDelete(ctx, toDelete, f); err != nil {
							return
						}
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
	}); err != nil {
		return err
	}

	return nil
}

func (j *cleanJob) flagFolderForDelete(ctx context.Context, toDelete *deleteSet, folder *models.Folder) error {
	// it is possible that child folders may be included while parent folders are not
	// so we need to check child folders separately
	toDelete.addFolder(folder.ID, folder.Path)

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

func (j *cleanJob) shouldClean(ctx context.Context, f models.File) bool {
	path := f.Base().Path

	info, err := f.Base().Info(j.FS)
	if err != nil && !isNotFound(err) {
		logger.Errorf("error getting file info for %q, not cleaning: %v", path, err)
		return false
	}

	if info == nil {
		// info is nil - file not exist
		logger.Infof("File not found. Marking to clean: \"%s\"", path)
		return true
	}

	// run through path filter, if returns false then the file should be cleaned
	filter := j.options.PathFilter

	// don't log anything - assume filter will have logged the reason
	return !filter.Accept(ctx, path, info)
}

func (j *cleanJob) shouldCleanFolder(ctx context.Context, f *models.Folder) bool {
	path := f.Path

	info, err := f.Info(j.FS)

	if err != nil && !isNotFound(err) {
		logger.Errorf("error getting folder info for %q, not cleaning: %v", path, err)
		return false
	}

	if info == nil {
		// info is nil - file not exist
		logger.Infof("Folder not found. Marking to clean: \"%s\"", path)
		return true
	}

	// #3261 - handle symlinks
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		finalPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			// don't bail out if symlink is invalid
			logger.Infof("Invalid symlink. Marking to clean: \"%s\"", path)
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

	// don't log anything - assume filter will have logged the reason
	return !filter.Accept(ctx, path, info)
}

func (j *cleanJob) deleteFile(ctx context.Context, fileID models.FileID, fn string) {
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

func (j *cleanJob) deleteFolder(ctx context.Context, folderID models.FolderID, fn string) {
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

func (j *cleanJob) fireHandlers(ctx context.Context, fileDeleter *Deleter, fileID models.FileID) error {
	for _, h := range j.Handlers {
		if err := h.HandleFile(ctx, fileDeleter, fileID); err != nil {
			return err
		}
	}

	return nil
}

func (j *cleanJob) fireFolderHandlers(ctx context.Context, fileDeleter *Deleter, folderID models.FolderID) error {
	for _, h := range j.Handlers {
		if err := h.HandleFolder(ctx, fileDeleter, folderID); err != nil {
			return err
		}
	}

	return nil
}
