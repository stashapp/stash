package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type SourceFile interface {
	Open() (io.ReadCloser, error)
	Path() string
	FileInfo() fs.FileInfo
}

type FileBased interface {
	File() models.File
}

type Hasher interface {
	OSHash(src io.ReadSeeker, size int64) (string, error)
	MD5(src io.Reader) (string, error)
}

type Scanned struct {
	Old *models.File
	New *models.File
}

// FileUpdated returns true if both old and new files are present and not equal.
func (s Scanned) FileUpdated() bool {
	if s.Old == nil || s.New == nil {
		return false
	}

	return !s.Old.Equal(*s.New)
}

// ContentsChanged returns true if both old and new files are present and the file content is different.
func (s Scanned) ContentsChanged() bool {
	if s.Old == nil || s.New == nil {
		return false
	}

	if s.Old.Checksum != s.New.Checksum {
		return true
	}

	if s.Old.OSHash != s.New.OSHash {
		return true
	}

	return false
}

type Scanner struct {
	Hasher Hasher

	CalculateMD5    bool
	CalculateOSHash bool
}

func (o Scanner) ScanExisting(existing FileBased, file SourceFile) (h *Scanned, err error) {
	info := file.FileInfo()
	h = &Scanned{}

	existingFile := existing.File()
	h.Old = &existingFile

	updatedFile := existingFile
	h.New = &updatedFile

	// update existing data if needed
	// truncate to seconds, since we don't store beyond that in the database
	updatedFile.FileModTime = info.ModTime().Truncate(time.Second)
	updatedFile.Size = strconv.FormatInt(info.Size(), 10)

	modTimeChanged := !existingFile.FileModTime.Equal(updatedFile.FileModTime)

	// regenerate hash(es) if missing or file mod time changed
	if _, err = o.generateHashes(&updatedFile, file, modTimeChanged); err != nil {
		return nil, err
	}

	// notify of changes as needed
	// object exists, no further processing required
	return
}

func (o Scanner) ScanNew(file SourceFile) (*models.File, error) {
	info := file.FileInfo()
	sizeStr := strconv.FormatInt(info.Size(), 10)
	modTime := info.ModTime()
	f := models.File{
		Path:        file.Path(),
		Size:        sizeStr,
		FileModTime: modTime,
	}

	if _, err := o.generateHashes(&f, file, true); err != nil {
		return nil, err
	}

	return &f, nil
}

// generateHashes regenerates and sets the hashes in the provided File.
// It will not recalculate unless specified.
func (o Scanner) generateHashes(f *models.File, file SourceFile, regenerate bool) (changed bool, err error) {
	existing := *f

	var src io.ReadCloser
	if o.CalculateOSHash && (regenerate || f.OSHash == "") {
		logger.Infof("Calculating oshash for %s ...", f.Path)

		size := file.FileInfo().Size()

		// #2196 for symlinks
		// get the size of the actual file, not the symlink
		if file.FileInfo().Mode()&os.ModeSymlink == os.ModeSymlink {
			fi, err := os.Stat(f.Path)
			if err != nil {
				return false, err
			}
			logger.Debugf("File <%s> is symlink. Size changed from <%d> to <%d>", f.Path, size, fi.Size())
			size = fi.Size()
		}

		src, err = file.Open()
		if err != nil {
			return false, err
		}
		defer src.Close()

		seekSrc, valid := src.(io.ReadSeeker)
		if !valid {
			return false, fmt.Errorf("invalid source file type: %s", file.Path())
		}

		// regenerate hash
		var oshash string
		oshash, err = o.Hasher.OSHash(seekSrc, size)
		if err != nil {
			return false, fmt.Errorf("error generating oshash for %s: %w", file.Path(), err)
		}

		f.OSHash = oshash

		// reset reader to start of file
		_, err = seekSrc.Seek(0, io.SeekStart)
		if err != nil {
			return false, fmt.Errorf("error seeking to start of file in %s: %w", file.Path(), err)
		}
	}

	// always generate if MD5 is nil
	// only regenerate MD5 if:
	// - OSHash was not calculated, or
	// - existing OSHash is different to generated one
	// or if it was different to the previous version
	if o.CalculateMD5 && (f.Checksum == "" || (regenerate && (!o.CalculateOSHash || existing.OSHash != f.OSHash))) {
		logger.Infof("Calculating checksum for %s...", f.Path)

		if src == nil {
			src, err = file.Open()
			if err != nil {
				return false, err
			}
			defer src.Close()
		}

		// regenerate checksum
		var checksum string
		checksum, err = o.Hasher.MD5(src)
		if err != nil {
			return
		}

		f.Checksum = checksum
	}

	changed = (o.CalculateOSHash && (f.OSHash != existing.OSHash)) || (o.CalculateMD5 && (f.Checksum != existing.Checksum))

	return
}
