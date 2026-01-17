package file

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type folderRenameCandidate struct {
	folder *models.Folder
	found  int
	files  int
}

type folderRenameDetector struct {
	// candidates is a map of folder id to the number of files that match
	candidates map[models.FolderID]folderRenameCandidate
	// rejects is a set of folder ids which were found to still exist
	rejects map[models.FolderID]struct{}
}

func (d *folderRenameDetector) isReject(id models.FolderID) bool {
	_, ok := d.rejects[id]
	return ok
}

func (d *folderRenameDetector) getCandidate(id models.FolderID) *folderRenameCandidate {
	c, ok := d.candidates[id]
	if !ok {
		return nil
	}

	return &c
}

func (d *folderRenameDetector) setCandidate(c folderRenameCandidate) {
	d.candidates[c.folder.ID] = c
}

func (d *folderRenameDetector) reject(id models.FolderID) {
	d.rejects[id] = struct{}{}
}

// bestCandidate returns the folder that is the best candidate for a rename.
// This is the folder that has the largest number of its original files that
// are still present in the new location.
func (d *folderRenameDetector) bestCandidate() *models.Folder {
	if len(d.candidates) == 0 {
		return nil
	}

	var best *folderRenameCandidate

	for _, c := range d.candidates {
		// ignore folders that have less than 50% of their original files
		if c.found < c.files/2 {
			continue
		}

		// prefer the folder with the most files if the ratio is the same
		if best == nil || c.found > best.found {
			cc := c
			best = &cc
		}
	}

	if best == nil {
		return nil
	}

	return best.folder
}

func (s *Scanner) detectFolderMove(ctx context.Context, file ScannedFile) (*models.Folder, error) {
	// in order for a folder to be considered moved, the existing folder must be
	// missing, and the majority of the old folder's files must be present, unchanged,
	// in the new folder.

	detector := folderRenameDetector{
		candidates: make(map[models.FolderID]folderRenameCandidate),
		rejects:    make(map[models.FolderID]struct{}),
	}
	// rejects is a set of folder ids which were found to still exist

	r := s.Repository

	if err := SymWalk(file.FS, file.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// don't let errors prevent scanning
			logger.Errorf("error scanning %s: %v", path, err)
			return nil
		}

		// ignore root
		if path == file.Path {
			return nil
		}

		// ignore directories
		if d.IsDir() {
			return fs.SkipDir
		}

		info, err := d.Info()
		if err != nil {
			logger.Errorf("reading info for %q: %v", path, err)
			return nil
		}

		if !s.AcceptEntry(ctx, path, info) {
			return nil
		}

		size, err := GetFileSize(file.FS, path, info)
		if err != nil {
			return fmt.Errorf("getting file size for %q: %w", path, err)
		}

		// check if the file exists in the database based on basename, size and mod time
		existing, err := r.File.FindByFileInfo(ctx, info, size)
		if err != nil {
			return fmt.Errorf("checking for existing file %q: %w", path, err)
		}

		for _, e := range existing {
			// ignore files in zip files
			if e.Base().ZipFileID != nil {
				continue
			}

			parentFolderID := e.Base().ParentFolderID

			if detector.isReject(parentFolderID) {
				// folder was found to still exist, not a candidate
				continue
			}

			c := detector.getCandidate(parentFolderID)

			if c == nil {
				// need to check if the folder exists in the filesystem
				pf, err := r.Folder.Find(ctx, e.Base().ParentFolderID)
				if err != nil {
					return fmt.Errorf("getting parent folder %d: %w", e.Base().ParentFolderID, err)
				}

				if pf == nil {
					// shouldn't happen, but just in case
					continue
				}

				// parent folder must be missing
				_, err = file.FS.Lstat(pf.Path)
				if err == nil {
					// parent folder exists, not a candidate
					detector.reject(parentFolderID)
					continue
				}

				if !errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("checking for parent folder %q: %w", pf.Path, err)
				}

				// parent folder is missing, possible candidate
				// count the total number of files in the existing folder
				count, err := r.File.CountByFolderID(ctx, parentFolderID)
				if err != nil {
					return fmt.Errorf("counting files in folder %d: %w", parentFolderID, err)
				}

				if count == 0 {
					// no files in the folder, not a candidate
					detector.reject(parentFolderID)
					continue
				}

				c = &folderRenameCandidate{
					folder: pf,
					found:  0,
					files:  count,
				}
			}

			// increment the count and set it in the map
			c.found++
			detector.setCandidate(*c)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("walking filesystem for folder rename detection: %w", err)
	}

	return detector.bestCandidate(), nil
}
