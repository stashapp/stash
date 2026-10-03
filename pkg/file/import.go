package file

import (
	"context"
	gojson "encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

var ErrZipFileNotExist = errors.New("zip file does not exist")

type Importer struct {
	ReaderWriter models.FileFinderCreator
	FolderStore  models.FolderFinderCreator
	Input        jsonschema.DirEntry

	file   models.File
	folder *models.Folder
}

func (i *Importer) PreImport(ctx context.Context) error {
	var err error

	switch ff := i.Input.(type) {
	case *jsonschema.BaseDirEntry:
		i.folder, err = i.folderJSONToFolder(ctx, ff)
	default:
		i.file, err = i.fileJSONToFile(ctx, i.Input)
	}

	return err
}

func (i *Importer) folderJSONToFolder(ctx context.Context, baseJSON *jsonschema.BaseDirEntry) (*models.Folder, error) {
	ret := models.Folder{
		DirEntry: models.DirEntry{
			ModTime: baseJSON.ModTime.GetTime(),
		},
		Path:      baseJSON.Path,
		CreatedAt: baseJSON.CreatedAt.GetTime(),
		UpdatedAt: baseJSON.CreatedAt.GetTime(),
	}

	if err := i.populateZipFileID(ctx, &ret.DirEntry); err != nil {
		return nil, err
	}

	// set parent folder id during the creation process

	return &ret, nil
}

func (i *Importer) fileJSONToFile(ctx context.Context, fileJSON jsonschema.DirEntry) (models.File, error) {
	switch ff := fileJSON.(type) {
	case *jsonschema.VideoFile:
		baseFile, err := i.baseFileJSONToBaseFile(ctx, ff.BaseFile)
		if err != nil {
			return nil, err
		}
		return &models.VideoFile{
			BaseFile:         baseFile,
			Format:           ff.Format,
			Width:            ff.Width,
			Height:           ff.Height,
			Duration:         ff.Duration,
			VideoCodec:       ff.VideoCodec,
			AudioCodec:       ff.AudioCodec,
			FrameRate:        ff.FrameRate,
			BitRate:          ff.BitRate,
			Interactive:      ff.Interactive,
			InteractiveSpeed: ff.InteractiveSpeed,
			Projection:       projectionPtr(ff.Projection),
			StereoMode:       stereoModePtr(ff.StereoMode),
			VRCorrections:    vrCorrectionsFromJSON(ff.VRCorrections),
		}, nil
	case *jsonschema.ImageFile:
		baseFile, err := i.baseFileJSONToBaseFile(ctx, ff.BaseFile)
		if err != nil {
			return nil, err
		}
		return &models.ImageFile{
			BaseFile: baseFile,
			Format:   ff.Format,
			Width:    ff.Width,
			Height:   ff.Height,
		}, nil
	case *jsonschema.BaseFile:
		return i.baseFileJSONToBaseFile(ctx, ff)
	}

	return nil, errors.New("unknown file type")
}

func unmarshalFingerprintValue(fp gojson.RawMessage) (interface{}, error) {
	// try to unmarshal as string first
	var str string
	if err := gojson.Unmarshal(fp, &str); err == nil {
		return str, nil
	}

	// if that fails, try to unmarshal as number
	var num gojson.Number
	if err := gojson.Unmarshal(fp, &num); err == nil {
		// reject floating point values
		if strings.Contains(num.String(), ".") {
			return nil, fmt.Errorf("floating point fingerprint values are not supported: %s", num.String())
		}

		ret, err := num.Int64()
		if err != nil {
			return nil, fmt.Errorf("error parsing fingerprint number: %v", err)
		}
		return ret, nil
	}

	return nil, fmt.Errorf("unable to unmarshal fingerprint value: %s", string(fp))
}

func (i *Importer) baseFileJSONToBaseFile(ctx context.Context, baseJSON *jsonschema.BaseFile) (*models.BaseFile, error) {
	baseFile := models.BaseFile{
		DirEntry: models.DirEntry{
			ModTime: baseJSON.ModTime.GetTime(),
		},
		Basename:  filepath.Base(baseJSON.Path),
		Size:      baseJSON.Size,
		CreatedAt: baseJSON.CreatedAt.GetTime(),
		UpdatedAt: baseJSON.CreatedAt.GetTime(),
	}

	for _, fp := range baseJSON.Fingerprints {
		fingerprintValue, err := unmarshalFingerprintValue(fp.Fingerprint)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling fingerprint value for type %q: %v", fp.Type, err)
		}

		// Handle phash: convert hex string to int64, or use legacy int64 values
		if fp.Type == models.FingerprintTypePhash {
			switch v := fingerprintValue.(type) {
			case string:
				// New format: hex string
				phash, err := utils.StringToPhash(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing phash hex string %q: %v", v, err)
				} else {
					fingerprintValue = phash
				}
			case int64:
				// Old format: int64 number
				// nothing to do
			default:
				return nil, fmt.Errorf("unexpected type for phash fingerprint: %T", v)
			}
		}
		baseFile.Fingerprints = append(baseFile.Fingerprints, models.Fingerprint{
			Type:        fp.Type,
			Fingerprint: fingerprintValue,
		})
	}

	if err := i.populateZipFileID(ctx, &baseFile.DirEntry); err != nil {
		return nil, err
	}

	return &baseFile, nil
}

func (i *Importer) populateZipFileID(ctx context.Context, f *models.DirEntry) error {
	zipFilePath := i.Input.DirEntry().ZipFile
	if zipFilePath != "" {
		zf, err := i.ReaderWriter.FindByPath(ctx, zipFilePath, true)
		if err != nil {
			return fmt.Errorf("error finding file by path %q: %v", zipFilePath, err)
		}

		if zf == nil {
			return ErrZipFileNotExist
		}

		id := zf.Base().ID
		f.ZipFileID = &id
	}

	return nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	return nil
}

func (i *Importer) Name() string {
	return i.Input.DirEntry().Path
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	path := i.Input.DirEntry().Path
	existing, err := i.ReaderWriter.FindByPath(ctx, path, true)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := int(existing.Base().ID)
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) createFolderHierarchy(ctx context.Context, p string) (*models.Folder, error) {
	parentPath := filepath.Dir(p)

	if parentPath == p {
		// get or create this folder
		return i.getOrCreateFolder(ctx, p, nil)
	}

	parent, err := i.createFolderHierarchy(ctx, parentPath)
	if err != nil {
		return nil, err
	}

	return i.getOrCreateFolder(ctx, p, parent)
}

func (i *Importer) getOrCreateFolder(ctx context.Context, path string, parent *models.Folder) (*models.Folder, error) {
	folder, err := i.FolderStore.FindByPath(ctx, path, true)
	if err != nil {
		return nil, err
	}

	if folder != nil {
		return folder, nil
	}

	now := time.Now()

	folder = &models.Folder{
		Path:      path,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if parent != nil {
		folder.ZipFileID = parent.ZipFileID
		folder.ParentFolderID = &parent.ID
	}

	if err := i.FolderStore.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (i *Importer) Create(ctx context.Context) (*int, error) {
	// create folder hierarchy and set parent folder id
	path := i.Input.DirEntry().Path
	path = filepath.Dir(path)
	folder, err := i.createFolderHierarchy(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("creating folder hierarchy for %q: %w", path, err)
	}

	if i.folder != nil {
		return i.createFolder(ctx, folder)
	}

	return i.createFile(ctx, folder)
}

func (i *Importer) createFile(ctx context.Context, parentFolder *models.Folder) (*int, error) {
	if parentFolder != nil {
		i.file.Base().ParentFolderID = parentFolder.ID
	}

	if err := i.ReaderWriter.Create(ctx, i.file); err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}

	id := int(i.file.Base().ID)
	return &id, nil
}

func (i *Importer) createFolder(ctx context.Context, parentFolder *models.Folder) (*int, error) {
	if parentFolder != nil {
		i.folder.ParentFolderID = &parentFolder.ID
	}

	if err := i.FolderStore.Create(ctx, i.folder); err != nil {
		return nil, fmt.Errorf("error creating folder: %w", err)
	}

	id := int(i.folder.ID)
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	// update not supported
	return nil
}

func projectionPtr(s *string) *models.ProjectionEnum {
	if s == nil {
		return nil
	}
	p := models.ProjectionEnum(*s)
	return &p
}

func stereoModePtr(s *string) *models.StereoModeEnum {
	if s == nil {
		return nil
	}
	m := models.StereoModeEnum(*s)
	return &m
}

func vrCorrectionsFromJSON(jc *jsonschema.VRCorrections) *models.VRCorrections {
	if jc == nil {
		return nil
	}
	return &models.VRCorrections{
		HorizontalOffset: jc.HorizontalOffset,
		VerticalOffset:   jc.VerticalOffset,
		AlphaMode:        jc.AlphaMode,
	}
}
