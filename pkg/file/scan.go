package file

import (
	"os"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type FileBased interface {
	File() models.File
}

type Hasher interface {
	OSHash(path string) (string, error)
	MD5(path string) (string, error)
}

type Scanned struct {
	Old *models.File
	New *models.File
}

func (s Scanned) FileUpdated() bool {
	if s.Old == nil || s.New == nil {
		return false
	}

	return s.Old.Equal(*s.New)
}

func (s Scanned) ContentsChanged() bool {
	if s.Old == nil || s.New == nil {
		return false
	}

	if s.Old.Checksum != nil && s.New.Checksum != nil && *s.Old.Checksum != *s.New.Checksum {
		return true
	}

	if s.Old.OSHash != nil && s.New.OSHash != nil && *s.Old.OSHash != *s.New.OSHash {
		return true
	}

	return false
}

type Scanner struct {
	Hasher Hasher

	CalculateMD5    bool
	CalculateOSHash bool
}

func (o Scanner) ScanExisting(existing FileBased, path string, info os.FileInfo) (h *Scanned, err error) {
	h = &Scanned{}

	existingFile := existing.File()
	h.Old = &existingFile

	updatedFile := existingFile
	h.New = &updatedFile

	//  update existing data if needed
	//  if item file mod time not set, then set it
	if existingFile.FileModTime == nil {
		t := info.ModTime()
		updatedFile.FileModTime = &t
	}

	modTimeChanged := existingFile.FileModTime == nil || existingFile.FileModTime.Equal(info.ModTime())

	//  regenerate hash(es)
	if _, err = o.generateHashes(&updatedFile, modTimeChanged); err != nil {
		return nil, err
	}

	// notify of changes as needed
	// object exists, no further processing required
	return
}

func (o Scanner) ScanNew(path string, info os.FileInfo) (h *Scanned, err error) {
	sizeStr := strconv.FormatInt(info.Size(), 10)
	modTime := info.ModTime()
	f := models.File{
		Path:        path,
		Size:        &sizeStr,
		FileModTime: &modTime,
	}

	_, err = o.generateHashes(&f, false)
	if err != nil {
		return nil, err
	}

	return
}

// generateHashes regenerates and sets the hashes in the provided File.
// It will not recalculate unless specified.
func (o Scanner) generateHashes(file *models.File, regenerate bool) (changed bool, err error) {
	existing := *file

	if o.CalculateOSHash && (regenerate || file.OSHash == nil) {
		logger.Infof("Calculating oshash for %s ...", file.Path)
		// regenerate hash
		var oshash string
		oshash, err = o.Hasher.OSHash(file.Path)
		if err != nil {
			return
		}

		file.OSHash = &oshash
	}

	// always generate if MD5 is nil
	// only regenerate MD5 if:
	// - OSHash was not calculated, or
	// - existing OSHash is different to generated one
	// or if it was different to the previous version
	if o.CalculateMD5 && (file.Checksum == nil || (regenerate && (!o.CalculateOSHash || existing.OSHash == nil || *existing.OSHash != *file.OSHash))) {
		logger.Infof("Calculating checksum for %s...", file.Path)

		// regenerate checksum
		var checksum string
		checksum, err = o.Hasher.MD5(file.Path)
		if err != nil {
			return
		}

		file.Checksum = &checksum
	}

	changed = (o.CalculateOSHash && (existing.OSHash == nil || *file.OSHash != *existing.OSHash)) || (o.CalculateMD5 && (existing.Checksum == nil || *file.Checksum != *existing.Checksum))

	return
}
