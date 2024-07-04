package manager

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/group"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
)

type Resetter interface {
	Reset() error
}

type ImportTask struct {
	repository models.Repository
	resetter   Resetter
	json       jsonUtils

	BaseDir             string
	TmpZip              string
	Reset               bool
	DuplicateBehaviour  ImportDuplicateEnum
	MissingRefBehaviour models.ImportMissingRefEnum

	fileNamingAlgorithm models.HashAlgorithm
}

type ImportObjectsInput struct {
	File                graphql.Upload              `json:"file"`
	DuplicateBehaviour  ImportDuplicateEnum         `json:"duplicateBehaviour"`
	MissingRefBehaviour models.ImportMissingRefEnum `json:"missingRefBehaviour"`
}

func CreateImportTask(a models.HashAlgorithm, input ImportObjectsInput) (*ImportTask, error) {
	baseDir, err := instance.Paths.Generated.TempDir("import")
	if err != nil {
		logger.Errorf("error creating temporary directory for import: %v", err)
		return nil, err
	}

	tmpZip := ""
	if input.File.File != nil {
		tmpZip = filepath.Join(baseDir, "import.zip")
		out, err := os.Create(tmpZip)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(out, input.File.File)
		out.Close()
		if err != nil {
			return nil, err
		}
	}

	mgr := GetInstance()
	return &ImportTask{
		repository:          mgr.Repository,
		resetter:            mgr.Database,
		BaseDir:             baseDir,
		TmpZip:              tmpZip,
		Reset:               false,
		DuplicateBehaviour:  input.DuplicateBehaviour,
		MissingRefBehaviour: input.MissingRefBehaviour,
		fileNamingAlgorithm: a,
	}, nil
}

func (t *ImportTask) GetDescription() string {
	return "Importing..."
}

func (t *ImportTask) Start(ctx context.Context) {
	if t.TmpZip != "" {
		defer func() {
			err := fsutil.RemoveDir(t.BaseDir)
			if err != nil {
				logger.Errorf("error removing directory %s: %v", t.BaseDir, err)
			}
		}()

		if err := t.unzipFile(); err != nil {
			logger.Errorf("error unzipping provided file for import: %v", err)
			return
		}
	}

	t.json = jsonUtils{
		json: *paths.GetJSONPaths(t.BaseDir),
	}

	// set default behaviour if not provided
	if !t.DuplicateBehaviour.IsValid() {
		t.DuplicateBehaviour = ImportDuplicateEnumFail
	}
	if !t.MissingRefBehaviour.IsValid() {
		t.MissingRefBehaviour = models.ImportMissingRefEnumFail
	}

	if t.Reset {
		err := t.resetter.Reset()

		if err != nil {
			logger.Errorf("Error resetting database: %v", err)
			return
		}
	}

	t.ImportTags(ctx)
	t.ImportPerformers(ctx)
	t.ImportStudios(ctx)
	t.ImportGroups(ctx)
	t.ImportFiles(ctx)
	t.ImportGalleries(ctx)

	t.ImportScenes(ctx)
	t.ImportImages(ctx)
}

func (t *ImportTask) unzipFile() error {
	defer func() {
		err := os.Remove(t.TmpZip)
		if err != nil {
			logger.Errorf("error removing temporary zip file %s: %v", t.TmpZip, err)
		}
	}()

	// now we can read the zip file
	r, err := zip.OpenReader(t.TmpZip)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fn := filepath.Join(t.BaseDir, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fn, os.ModePerm); err != nil {
				logger.Warnf("couldn't create directory %v while unzipping import file: %v", fn, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fn), os.ModePerm); err != nil {
			return err
		}

		o, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		i, err := f.Open()
		if err != nil {
			o.Close()
			return err
		}

		if _, err := io.Copy(o, i); err != nil {
			o.Close()
			i.Close()
			return err
		}

		o.Close()
		i.Close()
	}

	return nil
}

func (t *ImportTask) ImportPerformers(ctx context.Context) {
	logger.Info("[performers] importing")

	path := t.json.json.Performers
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[performers] failed to read performers directory: %v", err)
		}

		return
	}

	r := t.repository

	for i, fi := range files {
		index := i + 1
		performerJSON, err := jsonschema.LoadPerformerFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Errorf("[performers] failed to read json: %v", err)
			continue
		}

		logger.Progressf("[performers] %d of %d", index, len(files))

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			importer := &performer.Importer{
				ReaderWriter: r.Performer,
				TagWriter:    r.Tag,
				Input:        *performerJSON,
			}

			return performImport(ctx, importer, t.DuplicateBehaviour)
		}); err != nil {
			logger.Errorf("[performers] <%s> import failed: %v", fi.Name(), err)
		}
	}

	logger.Info("[performers] import complete")
}

func (t *ImportTask) ImportStudios(ctx context.Context) {
	pendingParent := make(map[string][]*jsonschema.Studio)

	logger.Info("[studios] importing")

	path := t.json.json.Studios
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[studios] failed to read studios directory: %v", err)
		}

		return
	}

	r := t.repository

	for i, fi := range files {
		index := i + 1
		studioJSON, err := jsonschema.LoadStudioFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Errorf("[studios] failed to read json: %v", err)
			continue
		}

		logger.Progressf("[studios] %d of %d", index, len(files))

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			return t.importStudio(ctx, studioJSON, pendingParent)
		}); err != nil {
			if errors.Is(err, studio.ErrParentStudioNotExist) {
				// add to the pending parent list so that it is created after the parent
				s := pendingParent[studioJSON.ParentStudio]
				s = append(s, studioJSON)
				pendingParent[studioJSON.ParentStudio] = s
				continue
			}

			logger.Errorf("[studios] <%s> failed to create: %v", fi.Name(), err)
			continue
		}
	}

	// create the leftover studios, warning for missing parents
	if len(pendingParent) > 0 {
		logger.Warnf("[studios] importing studios with missing parents")

		for _, s := range pendingParent {
			for _, orphanStudioJSON := range s {
				if err := r.WithTxn(ctx, func(ctx context.Context) error {
					return t.importStudio(ctx, orphanStudioJSON, nil)
				}); err != nil {
					logger.Errorf("[studios] <%s> failed to create: %v", orphanStudioJSON.Name, err)
					continue
				}
			}
		}
	}

	logger.Info("[studios] import complete")
}

func (t *ImportTask) importStudio(ctx context.Context, studioJSON *jsonschema.Studio, pendingParent map[string][]*jsonschema.Studio) error {
	r := t.repository

	importer := &studio.Importer{
		ReaderWriter:        t.repository.Studio,
		TagWriter:           r.Tag,
		Input:               *studioJSON,
		MissingRefBehaviour: t.MissingRefBehaviour,
	}

	// first phase: return error if parent does not exist
	if pendingParent != nil {
		importer.MissingRefBehaviour = models.ImportMissingRefEnumFail
	}

	if err := performImport(ctx, importer, t.DuplicateBehaviour); err != nil {
		return err
	}

	// now create the studios pending this studios creation
	s := pendingParent[studioJSON.Name]
	for _, childStudioJSON := range s {
		// map is nil since we're not checking parent studios at this point
		if err := t.importStudio(ctx, childStudioJSON, nil); err != nil {
			return fmt.Errorf("failed to create child studio <%s>: %v", childStudioJSON.Name, err)
		}
	}

	// delete the entry from the map so that we know its not left over
	delete(pendingParent, studioJSON.Name)

	return nil
}

func (t *ImportTask) ImportGroups(ctx context.Context) {
	logger.Info("[groups] importing")

	path := t.json.json.Groups
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[groups] failed to read movies directory: %v", err)
		}

		return
	}

	r := t.repository

	for i, fi := range files {
		index := i + 1
		groupJSON, err := jsonschema.LoadGroupFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Errorf("[groups] failed to read json: %v", err)
			continue
		}

		logger.Progressf("[groups] %d of %d", index, len(files))

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			groupImporter := &group.Importer{
				ReaderWriter:        r.Group,
				StudioWriter:        r.Studio,
				TagWriter:           r.Tag,
				Input:               *groupJSON,
				MissingRefBehaviour: t.MissingRefBehaviour,
			}

			return performImport(ctx, groupImporter, t.DuplicateBehaviour)
		}); err != nil {
			logger.Errorf("[groups] <%s> import failed: %v", fi.Name(), err)
			continue
		}
	}

	logger.Info("[groups] import complete")
}

func (t *ImportTask) ImportFiles(ctx context.Context) {
	logger.Info("[files] importing")

	path := t.json.json.Files
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[files] failed to read files directory: %v", err)
		}

		return
	}

	r := t.repository

	pendingParent := make(map[string][]jsonschema.DirEntry)

	for i, fi := range files {
		index := i + 1
		fileJSON, err := jsonschema.LoadFileFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Errorf("[files] failed to read json: %v", err)
			continue
		}

		logger.Progressf("[files] %d of %d", index, len(files))

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			return t.importFile(ctx, fileJSON, pendingParent)
		}); err != nil {
			if errors.Is(err, file.ErrZipFileNotExist) {
				// add to the pending parent list so that it is created after the parent
				s := pendingParent[fileJSON.DirEntry().ZipFile]
				s = append(s, fileJSON)
				pendingParent[fileJSON.DirEntry().ZipFile] = s
				continue
			}

			logger.Errorf("[files] <%s> failed to create: %v", fi.Name(), err)
			continue
		}
	}

	// create the leftover studios, warning for missing parents
	if len(pendingParent) > 0 {
		logger.Warnf("[files] importing files with missing zip files")

		for _, s := range pendingParent {
			for _, orphanFileJSON := range s {
				if err := r.WithTxn(ctx, func(ctx context.Context) error {
					return t.importFile(ctx, orphanFileJSON, nil)
				}); err != nil {
					logger.Errorf("[files] <%s> failed to create: %v", orphanFileJSON.DirEntry().Path, err)
					continue
				}
			}
		}
	}

	logger.Info("[files] import complete")
}

func (t *ImportTask) importFile(ctx context.Context, fileJSON jsonschema.DirEntry, pendingParent map[string][]jsonschema.DirEntry) error {
	r := t.repository

	fileImporter := &file.Importer{
		ReaderWriter: r.File,
		FolderStore:  r.Folder,
		Input:        fileJSON,
	}

	// ignore duplicate files - don't overwrite
	if err := performImport(ctx, fileImporter, ImportDuplicateEnumIgnore); err != nil {
		return err
	}

	// now create the files pending this file's creation
	s := pendingParent[fileJSON.DirEntry().Path]
	for _, childFileJSON := range s {
		// map is nil since we're not checking parent studios at this point
		if err := t.importFile(ctx, childFileJSON, nil); err != nil {
			return fmt.Errorf("failed to create child file <%s>: %v", childFileJSON.DirEntry().Path, err)
		}
	}

	// delete the entry from the map so that we know its not left over
	delete(pendingParent, fileJSON.DirEntry().Path)

	return nil
}

func (t *ImportTask) ImportGalleries(ctx context.Context) {
	logger.Info("[galleries] importing")

	path := t.json.json.Galleries
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[galleries] failed to read galleries directory: %v", err)
		}

		return
	}

	r := t.repository

	for i, fi := range files {
		index := i + 1
		galleryJSON, err := jsonschema.LoadGalleryFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Errorf("[galleries] failed to read json: %v", err)
			continue
		}

		logger.Progressf("[galleries] %d of %d", index, len(files))

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			galleryImporter := &gallery.Importer{
				ReaderWriter:        r.Gallery,
				FolderFinder:        r.Folder,
				FileFinder:          r.File,
				PerformerWriter:     r.Performer,
				StudioWriter:        r.Studio,
				TagWriter:           r.Tag,
				Input:               *galleryJSON,
				MissingRefBehaviour: t.MissingRefBehaviour,
			}

			if err := performImport(ctx, galleryImporter, t.DuplicateBehaviour); err != nil {
				return err
			}

			// import the gallery chapters
			for _, m := range galleryJSON.Chapters {
				chapterImporter := &gallery.ChapterImporter{
					GalleryID:           galleryImporter.ID,
					Input:               m,
					MissingRefBehaviour: t.MissingRefBehaviour,
					ReaderWriter:        r.GalleryChapter,
				}

				if err := performImport(ctx, chapterImporter, t.DuplicateBehaviour); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			logger.Errorf("[galleries] <%s> import failed to commit: %v", fi.Name(), err)
			continue
		}
	}

	logger.Info("[galleries] import complete")
}

func (t *ImportTask) ImportTags(ctx context.Context) {
	pendingParent := make(map[string][]*jsonschema.Tag)
	logger.Info("[tags] importing")

	path := t.json.json.Tags
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[tags] failed to read tags directory: %v", err)
		}

		return
	}

	r := t.repository

	for i, fi := range files {
		index := i + 1
		tagJSON, err := jsonschema.LoadTagFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Errorf("[tags] failed to read json: %v", err)
			continue
		}

		logger.Progressf("[tags] %d of %d", index, len(files))

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			return t.importTag(ctx, tagJSON, pendingParent, false)
		}); err != nil {
			var parentError tag.ParentTagNotExistError
			if errors.As(err, &parentError) {
				pendingParent[parentError.MissingParent()] = append(pendingParent[parentError.MissingParent()], tagJSON)
				continue
			}

			logger.Errorf("[tags] <%s> failed to import: %v", fi.Name(), err)
			continue
		}
	}

	for _, s := range pendingParent {
		for _, orphanTagJSON := range s {
			if err := r.WithTxn(ctx, func(ctx context.Context) error {
				return t.importTag(ctx, orphanTagJSON, nil, true)
			}); err != nil {
				logger.Errorf("[tags] <%s> failed to create: %v", orphanTagJSON.Name, err)
				continue
			}
		}
	}

	logger.Info("[tags] import complete")
}

func (t *ImportTask) importTag(ctx context.Context, tagJSON *jsonschema.Tag, pendingParent map[string][]*jsonschema.Tag, fail bool) error {
	importer := &tag.Importer{
		ReaderWriter:        t.repository.Tag,
		Input:               *tagJSON,
		MissingRefBehaviour: t.MissingRefBehaviour,
	}

	// first phase: return error if parent does not exist
	if !fail {
		importer.MissingRefBehaviour = models.ImportMissingRefEnumFail
	}

	if err := performImport(ctx, importer, t.DuplicateBehaviour); err != nil {
		return err
	}

	for _, childTagJSON := range pendingParent[tagJSON.Name] {
		if err := t.importTag(ctx, childTagJSON, pendingParent, fail); err != nil {
			var parentError tag.ParentTagNotExistError
			if errors.As(err, &parentError) {
				pendingParent[parentError.MissingParent()] = append(pendingParent[parentError.MissingParent()], childTagJSON)
				continue
			}

			return fmt.Errorf("failed to create child tag <%s>: %v", childTagJSON.Name, err)
		}
	}

	delete(pendingParent, tagJSON.Name)

	return nil
}

func (t *ImportTask) ImportScenes(ctx context.Context) {
	logger.Info("[scenes] importing")

	path := t.json.json.Scenes
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[scenes] failed to read scenes directory: %v", err)
		}

		return
	}

	r := t.repository

	for i, fi := range files {
		index := i + 1

		logger.Progressf("[scenes] %d of %d", index, len(files))

		sceneJSON, err := jsonschema.LoadSceneFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Infof("[scenes] <%s> json parse failure: %v", fi.Name(), err)
			continue
		}

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			sceneImporter := &scene.Importer{
				ReaderWriter: r.Scene,
				Input:        *sceneJSON,
				FileFinder:   r.File,

				FileNamingAlgorithm: t.fileNamingAlgorithm,
				MissingRefBehaviour: t.MissingRefBehaviour,

				GalleryFinder:   r.Gallery,
				GroupWriter:     r.Group,
				PerformerWriter: r.Performer,
				StudioWriter:    r.Studio,
				TagWriter:       r.Tag,
			}

			if err := performImport(ctx, sceneImporter, t.DuplicateBehaviour); err != nil {
				return err
			}

			// import the scene markers
			for _, m := range sceneJSON.Markers {
				markerImporter := &scene.MarkerImporter{
					SceneID:             sceneImporter.ID,
					Input:               m,
					MissingRefBehaviour: t.MissingRefBehaviour,
					ReaderWriter:        r.SceneMarker,
					TagWriter:           r.Tag,
				}

				if err := performImport(ctx, markerImporter, t.DuplicateBehaviour); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			logger.Errorf("[scenes] <%s> import failed: %v", fi.Name(), err)
		}
	}

	logger.Info("[scenes] import complete")
}

func (t *ImportTask) ImportImages(ctx context.Context) {
	logger.Info("[images] importing")

	path := t.json.json.Images
	files, err := os.ReadDir(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Errorf("[images] failed to read images directory: %v", err)
		}

		return
	}

	r := t.repository

	for i, fi := range files {
		index := i + 1

		logger.Progressf("[images] %d of %d", index, len(files))

		imageJSON, err := jsonschema.LoadImageFile(filepath.Join(path, fi.Name()))
		if err != nil {
			logger.Infof("[images] <%s> json parse failure: %v", fi.Name(), err)
			continue
		}

		if err := r.WithTxn(ctx, func(ctx context.Context) error {
			imageImporter := &image.Importer{
				ReaderWriter: r.Image,
				FileFinder:   r.File,
				Input:        *imageJSON,

				MissingRefBehaviour: t.MissingRefBehaviour,

				GalleryFinder:   r.Gallery,
				PerformerWriter: r.Performer,
				StudioWriter:    r.Studio,
				TagWriter:       r.Tag,
			}

			return performImport(ctx, imageImporter, t.DuplicateBehaviour)
		}); err != nil {
			logger.Errorf("[images] <%s> import failed: %v", fi.Name(), err)
		}
	}

	logger.Info("[images] import complete")
}
