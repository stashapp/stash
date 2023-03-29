package manager

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models/jsonschema"
)

// HACK: this is all here because of an import loop in jsonschema -> models -> file

var errZipFileNotExist = errors.New("zip file does not exist")

type fileFolderImporter struct {
	ReaderWriter file.Store
	FolderStore  file.FolderStore
	Input        jsonschema.DirEntry

	file   file.File
	folder *file.Folder
}

func (i *fileFolderImporter) PreImport(ctx context.Context) error {
	var err error

	switch ff := i.Input.(type) {
	case *jsonschema.BaseDirEntry:
		i.folder, err = i.folderJSONToFolder(ctx, ff)
	default:
		i.file, err = i.fileJSONToFile(ctx, i.Input)
	}

	return err
}

func (i *fileFolderImporter) folderJSONToFolder(ctx context.Context, baseJSON *jsonschema.BaseDirEntry) (*file.Folder, error) {
	ret := file.Folder{
		DirEntry: file.DirEntry{
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

func (i *fileFolderImporter) fileJSONToFile(ctx context.Context, fileJSON jsonschema.DirEntry) (file.File, error) {
	switch ff := fileJSON.(type) {
	case *jsonschema.VideoFile:
		baseFile, err := i.baseFileJSONToBaseFile(ctx, ff.BaseFile)
		if err != nil {
			return nil, err
		}
		return &file.VideoFile{
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
		}, nil
	case *jsonschema.ImageFile:
		baseFile, err := i.baseFileJSONToBaseFile(ctx, ff.BaseFile)
		if err != nil {
			return nil, err
		}
		return &file.ImageFile{
			BaseFile: baseFile,
			Format:   ff.Format,
			Width:    ff.Width,
			Height:   ff.Height,
		}, nil
	case *jsonschema.BaseFile:
		return i.baseFileJSONToBaseFile(ctx, ff)
	}

	return nil, fmt.Errorf("unknown file type")
}

func (i *fileFolderImporter) baseFileJSONToBaseFile(ctx context.Context, baseJSON *jsonschema.BaseFile) (*file.BaseFile, error) {
	baseFile := file.BaseFile{
		DirEntry: file.DirEntry{
			ModTime: baseJSON.ModTime.GetTime(),
		},
		Basename:  filepath.Base(baseJSON.Path),
		Size:      baseJSON.Size,
		CreatedAt: baseJSON.CreatedAt.GetTime(),
		UpdatedAt: baseJSON.CreatedAt.GetTime(),
	}

	for _, fp := range baseJSON.Fingerprints {
		baseFile.Fingerprints = append(baseFile.Fingerprints, file.Fingerprint{
			Type:        fp.Type,
			Fingerprint: fp.Fingerprint,
		})
	}

	if err := i.populateZipFileID(ctx, &baseFile.DirEntry); err != nil {
		return nil, err
	}

	return &baseFile, nil
}

func (i *fileFolderImporter) populateZipFileID(ctx context.Context, f *file.DirEntry) error {
	zipFilePath := i.Input.DirEntry().ZipFile
	if zipFilePath != "" {
		zf, err := i.ReaderWriter.FindByPath(ctx, zipFilePath)
		if err != nil {
			return fmt.Errorf("error finding file by path %q: %v", zipFilePath, err)
		}

		if zf == nil {
			return errZipFileNotExist
		}

		id := zf.Base().ID
		f.ZipFileID = &id
	}

	return nil
}

func (i *fileFolderImporter) PostImport(ctx context.Context, id int) error {
	return nil
}

func (i *fileFolderImporter) Name() string {
	return i.Input.DirEntry().Path
}

func (i *fileFolderImporter) FindExistingID(ctx context.Context) (*int, error) {
	path := i.Input.DirEntry().Path
	existing, err := i.ReaderWriter.FindByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := int(existing.Base().ID)
		return &id, nil
	}

	return nil, nil
}

func (i *fileFolderImporter) createFolderHierarchy(ctx context.Context, p string) (*file.Folder, error) {
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

func (i *fileFolderImporter) getOrCreateFolder(ctx context.Context, path string, parent *file.Folder) (*file.Folder, error) {
	folder, err := i.FolderStore.FindByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	if folder != nil {
		return folder, nil
	}

	now := time.Now()

	folder = &file.Folder{
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

func (i *fileFolderImporter) Create(ctx context.Context) (*int, error) {
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

func (i *fileFolderImporter) createFile(ctx context.Context, parentFolder *file.Folder) (*int, error) {
	if parentFolder != nil {
		i.file.Base().ParentFolderID = parentFolder.ID
	}

	if err := i.ReaderWriter.Create(ctx, i.file); err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}

	id := int(i.file.Base().ID)
	return &id, nil
}

func (i *fileFolderImporter) createFolder(ctx context.Context, parentFolder *file.Folder) (*int, error) {
	if parentFolder != nil {
		i.folder.ParentFolderID = &parentFolder.ID
	}

	if err := i.FolderStore.Create(ctx, i.folder); err != nil {
		return nil, fmt.Errorf("error creating folder: %w", err)
	}

	id := int(i.folder.ID)
	return &id, nil
}

func (i *fileFolderImporter) Update(ctx context.Context, id int) error {
	// update not supported
	return nil
}
