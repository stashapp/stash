package manager

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	"github.com/stashapp/stash/pkg/savedfilter"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
)

type Resetter interface {
	Reset() error
}

type ImportTask struct {
	repository        models.Repository
	resetter          Resetter
	json              jsonUtils
	BaseDir           string
	TmpZip            string
	Reset             bool
	DuplicateBehaviour ImportDuplicateEnum
	MissingRefBehaviour models.ImportMissingRefEnum
	fileNamingAlgorithm models.HashAlgorithm
}

type ImportObjectsInput struct {
	File               graphql.Upload `json:"file"`
	DuplicateBehaviour ImportDuplicateEnum `json:"duplicateBehaviour"`
	MissingRefBehaviour models.ImportMissingRefEnum `json:"missingRefBehaviour"`
}

func CreateImportTask(a models.HashAlgorithm, input ImportObjectsInput) (*ImportTask, error) {
	baseDir, err := instance.Paths.Generated.TempDir("import")
	if err != nil {
		logger.Errorf("error creating temporary directory import: %v", err)
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
		repository:         mgr.Repository,
		resetter:           mgr.Database,
		BaseDir:            baseDir,
		TmpZip:             tmpZip,
		Reset:              false,
		DuplicateBehaviour: input.DuplicateBehaviour,
		MissingRefBehaviour: input.MissingRefBehaviour,
		fileNamingAlgorithm: nil,
	}, nil
}

func (t *ImportTask) GetDescription() string {
	return "Importing..."
}

func (t *ImportTask) Start(ctx context.Context) {
	defer func() {
		err := fsutil.RemoveDir(t.BaseDir)
		if err != nil {
			logger.Errorf("error removing directory %s: %v", t.BaseDir, err)
		}
	}()

	err := t.unzipFile()
	if err != nil {
		logger.Errorf("error unzipping provided file import: %v", err)
		return
	}

	t.json = jsonUtils{
		json: paths.GetJSONPaths(t.BaseDir),
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

	t.ImportSavedFilters(ctx)
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

	// now read zip file
	r, err := zip.OpenReader(t.TmpZip)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fn := filepath.Join(t.BaseDir, f.Name)

		// Prevent Zip-Slip attacks by ensuring the resolved path stays within BaseDir
		rel, err := filepath.Rel(t.BaseDir, fn)
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("zip-slip attempt detected: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fn, os.ModePerm); err != nil {
				logger.Warnf("couldn't create directory while unzipping import file: %v", fn, err)
			}
			continue
		}

		if err := t.unzipFileEntry(fn); err != nil {
			return err
		}
	}

	return nil
}

func (t *ImportTask) unzipFileEntry(fn string) error {
	if err := os.MkdirAll(filepath.Dir(fn), os.ModePerm); err != nil {
		return err
	}

	o, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer o.Close()

	i, err := f.Open()
	if err != nil {
		return err
	}
	defer i.Close()

	_, err = io.Copy(o, i)
	return err
}
