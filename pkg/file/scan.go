package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

const mutexType = "file"

type SourceFile interface {
	Open() (io.ReadCloser, error)
	Path() string
	FileInfo() fs.FileInfo
	ZipFile() *models.File
}

type FileBased interface {
	File() models.File
}

type Statter interface {
	Stat(reader models.FileReader, f models.File) (fs.FileInfo, error)
}

type MetadataGenerator interface {
	GenerateMetadata(dest *models.File, src SourceFile) error
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
	Hasher            Hasher
	Statter           Statter
	MutexManager      *utils.MutexManager
	MetadataGenerator MetadataGenerator

	CalculateMD5    bool
	CalculateOSHash bool

	Done chan struct{}
}

func (o Scanner) ApplyChanges(rw models.FileWriter, scanned *Scanned) error {
	if scanned.Old == nil {
		// create the new file
		created, err := rw.Create(*scanned.New)
		if err != nil {
			return fmt.Errorf("creating file database entry: %w", err)
		}

		scanned.New.ID = created.ID
		return nil
	}

	_, err := rw.UpdateFull(*scanned.New)
	if err != nil {
		return fmt.Errorf("updating file database entry: %w", err)
	}
	return nil
}

func (o *Scanner) Close() {
	close(o.Done)
}

func fullPath(file SourceFile) string {
	zipFile := file.ZipFile()
	if zipFile != nil {
		return filepath.Join(zipFile.Path, file.Path())
	}

	return file.Path()
}

func (o Scanner) Scan(reader models.FileReader, file SourceFile) (h *Scanned, err error) {
	existing, err := reader.FindByPath(fullPath(file))
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return o.scanExisting(existing, file)
	}

	return o.scanNew(reader, file)
}

func (o Scanner) scanExisting(existingFile *models.File, file SourceFile) (h *Scanned, err error) {
	info := file.FileInfo()
	h = &Scanned{}

	h.Old = existingFile

	updatedFile := *existingFile
	h.New = &updatedFile

	// update existing data if needed
	// truncate to seconds, since we don't store beyond that in the database
	updatedFile.FileModTime = info.ModTime().Truncate(time.Second)
	updatedFile.Size = info.Size()

	modTimeChanged := !existingFile.FileModTime.Equal(updatedFile.FileModTime)

	// regenerate hash(es) if missing or file mod time changed
	if _, err = o.generateHashes(&updatedFile, file, modTimeChanged); err != nil {
		return nil, err
	}

	// only regenerate metdata if file has changed
	if o.MetadataGenerator != nil && h.ContentsChanged() {
		if err := o.MetadataGenerator.GenerateMetadata(h.New, file); err != nil {
			return nil, fmt.Errorf("generating metadata for %q: %w", file.Path(), err)
		}
	}

	// notify of changes as needed
	// object exists, no further processing required
	return
}

func (o Scanner) scanNew(reader models.FileReader, file SourceFile) (*Scanned, error) {
	info := file.FileInfo()
	modTime := info.ModTime()
	zipFile := file.ZipFile()
	f := &models.File{
		Path:        fullPath(file),
		Size:        info.Size(),
		FileModTime: modTime,
	}
	if zipFile != nil {
		f.ZipFileID = models.NullInt64(int64(zipFile.ID))
	}

	if _, err := o.generateHashes(f, file, true); err != nil {
		return nil, err
	}

	renamed, err := o.detectRename(reader, *f)
	if err != nil {
		return nil, err
	}

	if renamed != nil {
		return renamed, nil
	}

	ret := &Scanned{
		New: f,
	}

	if o.MetadataGenerator != nil {
		if err := o.MetadataGenerator.GenerateMetadata(ret.New, file); err != nil {
			return nil, fmt.Errorf("generating metadata for %q: %w", file.Path(), err)
		}
	}

	return ret, nil
}

// detectRename performs rename detection - find files that have the same hash
// and ensure all exist in the file system. For the first one that does not
// exist, treat as a rename. Returns the old and new file objects if a rename
// is detected.
func (o Scanner) detectRename(reader models.FileReader, f models.File) (*Scanned, error) {
	ret := &Scanned{}

	var existingFiles []*models.File
	var err error
	if o.CalculateOSHash {
		existingFiles, err = reader.FindByOSHash(f.OSHash)
	} else {
		existingFiles, err = reader.FindByChecksum(f.Checksum)
	}

	if err != nil {
		return nil, err
	}

	for _, ff := range existingFiles {
		_, err := o.Statter.Stat(reader, *ff)
		if errors.Is(err, fs.ErrNotExist) {
			// treat as a rename
			logger.Infof("Detected move: %s -> %s", ff.Path, f.Path)

			ret.Old = ff
			ret.New = &f
			f.ID = ff.ID

			return ret, nil
		} else if err != nil {
			return nil, err
		}
	}

	// treat as new, duplicate file
	return nil, nil
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

	if changed {
		o.claimHashes(f)
	}

	return
}

// claimHashes claims the hashes for the provided file, to ensure that no
// other threads can operate on files with these hashes.
func (o Scanner) claimHashes(f *models.File) {
	if f.OSHash != "" {
		o.MutexManager.Claim(mutexType, f.OSHash, o.Done)
	}
	if f.Checksum != "" {
		o.MutexManager.Claim(mutexType, f.Checksum, o.Done)
	}
}
