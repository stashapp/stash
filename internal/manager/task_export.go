package manager

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/group"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/savedfilter"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
)

type ExportTask struct {
	repository models.Repository
	full       bool

	baseDir string
	json    jsonUtils

	fileNamingAlgorithm models.HashAlgorithm

	scenes     *exportSpec
	images     *exportSpec
	performers *exportSpec
	groups     *exportSpec
	tags       *exportSpec
	studios    *exportSpec
	galleries  *exportSpec

	includeDependencies bool

	DownloadHash string
}

type ExportObjectTypeInput struct {
	Ids []string `json:"ids"`
	All *bool    `json:"all"`
}

type ExportObjectsInput struct {
	Scenes              *ExportObjectTypeInput `json:"scenes"`
	Images              *ExportObjectTypeInput `json:"images"`
	Studios             *ExportObjectTypeInput `json:"studios"`
	Performers          *ExportObjectTypeInput `json:"performers"`
	Tags                *ExportObjectTypeInput `json:"tags"`
	Groups              *ExportObjectTypeInput `json:"groups"`
	Movies              *ExportObjectTypeInput `json:"movies"` // deprecated
	Galleries           *ExportObjectTypeInput `json:"galleries"`
	IncludeDependencies *bool                  `json:"includeDependencies"`
}

type exportSpec struct {
	IDs []int
	all bool
}

func newExportSpec(input *ExportObjectTypeInput) *exportSpec {
	if input == nil {
		return &exportSpec{}
	}

	ids, _ := stringslice.StringSliceToIntSlice(input.Ids)

	ret := &exportSpec{
		IDs: ids,
	}

	if input.All != nil {
		ret.all = *input.All
	}

	return ret
}

func CreateExportTask(a models.HashAlgorithm, input ExportObjectsInput) *ExportTask {
	includeDeps := false
	if input.IncludeDependencies != nil {
		includeDeps = *input.IncludeDependencies
	}

	// handle deprecated Movies field
	groupSpec := input.Groups
	if groupSpec == nil && input.Movies != nil {
		groupSpec = input.Movies
	}

	return &ExportTask{
		repository:          GetInstance().Repository,
		fileNamingAlgorithm: a,
		scenes:              newExportSpec(input.Scenes),
		images:              newExportSpec(input.Images),
		performers:          newExportSpec(input.Performers),
		groups:              newExportSpec(groupSpec),
		tags:                newExportSpec(input.Tags),
		studios:             newExportSpec(input.Studios),
		galleries:           newExportSpec(input.Galleries),
		includeDependencies: includeDeps,
	}
}

func (t *ExportTask) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// @manager.total = Scene.count + Gallery.count + Performer.count + Studio.count + Group.count
	workerCount := runtime.GOMAXPROCS(0) // set worker count to number of cpus available

	startTime := time.Now()

	if t.full {
		t.baseDir = config.GetInstance().GetMetadataPath()
	} else {
		var err error
		t.baseDir, err = instance.Paths.Generated.TempDir("export")
		if err != nil {
			logger.Errorf("error creating temporary directory for export: %v", err)
			return
		}

		defer func() {
			err := fsutil.RemoveDir(t.baseDir)
			if err != nil {
				logger.Errorf("error removing directory %s: %v", t.baseDir, err)
			}
		}()
	}

	if t.baseDir == "" {
		logger.Errorf("baseDir must not be empty")
		return
	}

	t.json = jsonUtils{
		json: *paths.GetJSONPaths(t.baseDir),
	}

	paths.EmptyJSONDirs(t.baseDir)
	paths.EnsureJSONDirs(t.baseDir)

	txnErr := t.repository.WithTxn(ctx, func(ctx context.Context) error {
		// include group scenes and gallery images
		if !t.full {
			// only include group scenes if includeDependencies is also set
			if !t.scenes.all && t.includeDependencies {
				t.populateGroupScenes(ctx)
			}

			// always export gallery images
			if !t.images.all {
				t.populateGalleryImages(ctx)
			}
		}

		t.ExportScenes(ctx, workerCount)
		t.ExportImages(ctx, workerCount)
		t.ExportGalleries(ctx, workerCount)
		t.ExportGroups(ctx, workerCount)
		t.ExportPerformers(ctx, workerCount)
		t.ExportStudios(ctx, workerCount)
		t.ExportTags(ctx, workerCount)
		t.ExportSavedFilters(ctx, workerCount)

		return nil
	})
	if txnErr != nil {
		logger.Warnf("error while running export transaction: %v", txnErr)
	}

	if !t.full {
		err := t.generateDownload()
		if err != nil {
			logger.Errorf("error generating download link: %v", err)
			return
		}
	}
	logger.Infof("Export complete in %s.", time.Since(startTime))
}

func (t *ExportTask) generateDownload() error {
	// zip the files and register a download link
	if err := fsutil.EnsureDir(instance.Paths.Generated.Downloads); err != nil {
		return err
	}
	z, err := os.CreateTemp(instance.Paths.Generated.Downloads, "export*.zip")
	if err != nil {
		return err
	}
	defer z.Close()

	err = t.zipFiles(z)
	if err != nil {
		return err
	}

	t.DownloadHash, err = instance.DownloadStore.RegisterFile(z.Name(), "", false)
	if err != nil {
		return fmt.Errorf("error registering file for download: %w", err)
	}
	logger.Debugf("Generated zip file %s with hash %s", z.Name(), t.DownloadHash)
	return nil
}

func (t *ExportTask) zipFiles(w io.Writer) error {
	z := zip.NewWriter(w)
	defer z.Close()

	u := jsonUtils{
		json: *paths.GetJSONPaths(""),
	}

	walkWarn(t.json.json.Tags, t.zipWalkFunc(u.json.Tags, z))
	walkWarn(t.json.json.Galleries, t.zipWalkFunc(u.json.Galleries, z))
	walkWarn(t.json.json.Performers, t.zipWalkFunc(u.json.Performers, z))
	walkWarn(t.json.json.Studios, t.zipWalkFunc(u.json.Studios, z))
	walkWarn(t.json.json.Groups, t.zipWalkFunc(u.json.Groups, z))
	walkWarn(t.json.json.Scenes, t.zipWalkFunc(u.json.Scenes, z))
	walkWarn(t.json.json.Images, t.zipWalkFunc(u.json.Images, z))

	return nil
}

// like filepath.Walk but issue a warning on error
func walkWarn(root string, fn filepath.WalkFunc) {
	if err := filepath.Walk(root, fn); err != nil {
		logger.Warnf("error walking structure %v: %v", root, err)
	}
}

func (t *ExportTask) zipWalkFunc(outDir string, z *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return t.zipFile(path, outDir, z)
	}
}

func (t *ExportTask) zipFile(fn, outDir string, z *zip.Writer) error {
	bn := filepath.Base(fn)

	p := filepath.Join(outDir, bn)
	p = filepath.ToSlash(p)

	f, err := z.Create(p)
	if err != nil {
		return fmt.Errorf("error creating zip entry for %s: %v", fn, err)
	}

	i, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", fn, err)
	}

	defer i.Close()

	if _, err := io.Copy(f, i); err != nil {
		return fmt.Errorf("error writing %s to zip: %v", fn, err)
	}

	return nil
}

func (t *ExportTask) populateGroupScenes(ctx context.Context) {
	r := t.repository
	reader := r.Group
	sceneReader := r.Scene

	var groups []*models.Group
	var err error
	all := t.full || (t.groups != nil && t.groups.all)
	if all {
		groups, err = reader.All(ctx)
	} else if t.groups != nil && len(t.groups.IDs) > 0 {
		groups, err = reader.FindMany(ctx, t.groups.IDs)
	}

	if err != nil {
		logger.Errorf("[groups] failed to fetch groups: %v", err)
	}

	for _, m := range groups {
		scenes, err := sceneReader.FindByGroupID(ctx, m.ID)
		if err != nil {
			logger.Errorf("[groups] <%s> failed to fetch scenes for group: %v", m.Name, err)
			continue
		}

		for _, s := range scenes {
			t.scenes.IDs = sliceutil.AppendUnique(t.scenes.IDs, s.ID)
		}
	}
}

func (t *ExportTask) populateGalleryImages(ctx context.Context) {
	r := t.repository
	reader := r.Gallery
	imageReader := r.Image

	var galleries []*models.Gallery
	var err error
	all := t.full || (t.galleries != nil && t.galleries.all)
	if all {
		galleries, err = reader.All(ctx)
	} else if t.galleries != nil && len(t.galleries.IDs) > 0 {
		galleries, err = reader.FindMany(ctx, t.galleries.IDs)
	}

	if err != nil {
		logger.Errorf("[galleries] failed to fetch galleries: %v", err)
	}

	for _, g := range galleries {
		if err := g.LoadFiles(ctx, reader); err != nil {
			logger.Errorf("[galleries] <%s> failed to fetch files for gallery: %v", g.DisplayName(), err)
			continue
		}

		images, err := imageReader.FindByGalleryID(ctx, g.ID)
		if err != nil {
			logger.Errorf("[galleries] <%s> failed to fetch images for gallery: %v", g.DisplayName(), err)
			continue
		}

		for _, i := range images {
			t.images.IDs = sliceutil.AppendUnique(t.images.IDs, i.ID)
		}
	}
}

func (t *ExportTask) ExportScenes(ctx context.Context, workers int) {
	var scenesWg sync.WaitGroup

	sceneReader := t.repository.Scene

	var scenes []*models.Scene
	var err error
	all := t.full || (t.scenes != nil && t.scenes.all)
	if all {
		scenes, err = sceneReader.All(ctx)
	} else if t.scenes != nil && len(t.scenes.IDs) > 0 {
		scenes, err = sceneReader.FindMany(ctx, t.scenes.IDs)
	}

	if err != nil {
		logger.Errorf("[scenes] failed to fetch scenes: %v", err)
	}

	jobCh := make(chan *models.Scene, workers*2) // make a buffered channel to feed workers

	logger.Info("[scenes] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Scene workers
		scenesWg.Add(1)
		go t.exportScene(ctx, &scenesWg, jobCh)
	}

	for i, scene := range scenes {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[scenes] %d of %d", index, len(scenes))
		}
		jobCh <- scene // feed workers
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	scenesWg.Wait()

	logger.Infof("[scenes] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportFile(f models.File) {
	newFileJSON := fileToJSON(f)

	fn := newFileJSON.Filename()

	if err := t.json.saveFile(fn, newFileJSON); err != nil {
		logger.Errorf("[files] <%s> failed to save json: %v", fn, err)
	}
}

func fileToJSON(f models.File) jsonschema.DirEntry {
	bf := f.Base()

	base := jsonschema.BaseFile{
		BaseDirEntry: jsonschema.BaseDirEntry{
			Type:      jsonschema.DirEntryTypeFile,
			ModTime:   json.JSONTime{Time: bf.ModTime},
			Path:      bf.Path,
			CreatedAt: json.JSONTime{Time: bf.CreatedAt},
			UpdatedAt: json.JSONTime{Time: bf.UpdatedAt},
		},
		Size: bf.Size,
	}

	if bf.ZipFile != nil {
		base.ZipFile = bf.ZipFile.Base().Path
	}

	for _, fp := range bf.Fingerprints {
		base.Fingerprints = append(base.Fingerprints, jsonschema.Fingerprint{
			Type:        fp.Type,
			Fingerprint: fp.Fingerprint,
		})
	}

	switch ff := f.(type) {
	case *models.VideoFile:
		base.Type = jsonschema.DirEntryTypeVideo
		return jsonschema.VideoFile{
			BaseFile:         &base,
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
		}
	case *models.ImageFile:
		base.Type = jsonschema.DirEntryTypeImage
		return jsonschema.ImageFile{
			BaseFile: &base,
			Format:   ff.Format,
			Width:    ff.Width,
			Height:   ff.Height,
		}
	}

	return &base
}

func (t *ExportTask) exportFolder(f models.Folder) {
	newFileJSON := folderToJSON(f)

	fn := newFileJSON.Filename()

	if err := t.json.saveFile(fn, newFileJSON); err != nil {
		logger.Errorf("[files] <%s> failed to save json: %v", fn, err)
	}
}

func folderToJSON(f models.Folder) jsonschema.DirEntry {
	base := jsonschema.BaseDirEntry{
		Type:      jsonschema.DirEntryTypeFolder,
		ModTime:   json.JSONTime{Time: f.ModTime},
		Path:      f.Path,
		CreatedAt: json.JSONTime{Time: f.CreatedAt},
		UpdatedAt: json.JSONTime{Time: f.UpdatedAt},
	}

	if f.ZipFile != nil {
		base.ZipFile = f.ZipFile.Base().Path
	}

	return &base
}

func (t *ExportTask) exportScene(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.Scene) {
	defer wg.Done()

	r := t.repository
	sceneReader := r.Scene
	studioReader := r.Studio
	groupReader := r.Group
	galleryReader := r.Gallery
	performerReader := r.Performer
	tagReader := r.Tag
	sceneMarkerReader := r.SceneMarker

	for s := range jobChan {
		sceneHash := s.GetHash(t.fileNamingAlgorithm)

		if err := s.LoadRelationships(ctx, sceneReader); err != nil {
			logger.Errorf("[scenes] <%s> error loading scene relationships: %v", sceneHash, err)
		}

		newSceneJSON, err := scene.ToBasicJSON(ctx, sceneReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene JSON: %v", sceneHash, err)
			continue
		}

		// export files
		for _, f := range s.Files.List() {
			t.exportFile(f)
		}

		newSceneJSON.Studio, err = scene.GetStudioName(ctx, studioReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene studio name: %v", sceneHash, err)
			continue
		}

		galleries, err := galleryReader.FindBySceneID(ctx, s.ID)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene gallery checksums: %v", sceneHash, err)
			continue
		}

		for _, g := range galleries {
			if err := g.LoadFiles(ctx, galleryReader); err != nil {
				logger.Errorf("[scenes] <%s> error getting scene gallery files: %v", sceneHash, err)
				continue
			}
		}

		newSceneJSON.Galleries = gallery.GetRefs(galleries)

		newSceneJSON.ResumeTime = s.ResumeTime
		newSceneJSON.PlayDuration = s.PlayDuration

		performers, err := performerReader.FindBySceneID(ctx, s.ID)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene performer names: %v", sceneHash, err)
			continue
		}

		newSceneJSON.Performers = performer.GetNames(performers)

		newSceneJSON.Tags, err = scene.GetTagNames(ctx, tagReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene tag names: %v", sceneHash, err)
			continue
		}

		newSceneJSON.Markers, err = scene.GetSceneMarkersJSON(ctx, sceneMarkerReader, tagReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene markers JSON: %v", sceneHash, err)
			continue
		}

		newSceneJSON.Groups, err = scene.GetSceneGroupsJSON(ctx, groupReader, s)
		if err != nil {
			logger.Errorf("[scenes] <%s> error getting scene groups JSON: %v", sceneHash, err)
			continue
		}

		if t.includeDependencies {
			if s.StudioID != nil {
				t.studios.IDs = sliceutil.AppendUnique(t.studios.IDs, *s.StudioID)
			}

			t.galleries.IDs = sliceutil.AppendUniques(t.galleries.IDs, gallery.GetIDs(galleries))

			tagIDs, err := scene.GetDependentTagIDs(ctx, tagReader, sceneMarkerReader, s)
			if err != nil {
				logger.Errorf("[scenes] <%s> error getting scene tags: %v", sceneHash, err)
				continue
			}
			t.tags.IDs = sliceutil.AppendUniques(t.tags.IDs, tagIDs)

			groupIDs, err := scene.GetDependentGroupIDs(ctx, s)
			if err != nil {
				logger.Errorf("[scenes] <%s> error getting scene groups: %v", sceneHash, err)
				continue
			}
			t.groups.IDs = sliceutil.AppendUniques(t.groups.IDs, groupIDs)

			t.performers.IDs = sliceutil.AppendUniques(t.performers.IDs, performer.GetIDs(performers))
		}

		basename := filepath.Base(s.Path)
		hash := s.OSHash

		fn := newSceneJSON.Filename(s.ID, basename, hash)

		if err := t.json.saveScene(fn, newSceneJSON); err != nil {
			logger.Errorf("[scenes] <%s> failed to save json: %v", sceneHash, err)
		}
	}
}

func (t *ExportTask) ExportImages(ctx context.Context, workers int) {
	var imagesWg sync.WaitGroup

	r := t.repository
	imageReader := r.Image

	var images []*models.Image
	var err error
	all := t.full || (t.images != nil && t.images.all)
	if all {
		images, err = imageReader.All(ctx)
	} else if t.images != nil && len(t.images.IDs) > 0 {
		images, err = imageReader.FindMany(ctx, t.images.IDs)
	}

	if err != nil {
		logger.Errorf("[images] failed to fetch images: %v", err)
	}

	jobCh := make(chan *models.Image, workers*2) // make a buffered channel to feed workers

	logger.Info("[images] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Image workers
		imagesWg.Add(1)
		go t.exportImage(ctx, &imagesWg, jobCh)
	}

	for i, image := range images {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[images] %d of %d", index, len(images))
		}
		jobCh <- image // feed workers
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	imagesWg.Wait()

	logger.Infof("[images] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportImage(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.Image) {
	defer wg.Done()

	r := t.repository
	studioReader := r.Studio
	galleryReader := r.Gallery
	performerReader := r.Performer
	tagReader := r.Tag

	for s := range jobChan {
		imageHash := s.Checksum

		if err := s.LoadFiles(ctx, r.Image); err != nil {
			logger.Errorf("[images] <%s> error getting image files: %v", imageHash, err)
			continue
		}

		if err := s.LoadURLs(ctx, r.Image); err != nil {
			logger.Errorf("[images] <%s> error getting image urls: %v", imageHash, err)
			continue
		}

		newImageJSON := image.ToBasicJSON(s)

		// export files
		for _, f := range s.Files.List() {
			t.exportFile(f)
		}

		var err error
		newImageJSON.Studio, err = image.GetStudioName(ctx, studioReader, s)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image studio name: %v", imageHash, err)
			continue
		}

		imageGalleries, err := galleryReader.FindByImageID(ctx, s.ID)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image galleries: %v", imageHash, err)
			continue
		}

		for _, g := range imageGalleries {
			if err := g.LoadFiles(ctx, galleryReader); err != nil {
				logger.Errorf("[images] <%s> error getting image gallery files: %v", imageHash, err)
				continue
			}
		}

		newImageJSON.Galleries = gallery.GetRefs(imageGalleries)

		performers, err := performerReader.FindByImageID(ctx, s.ID)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image performer names: %v", imageHash, err)
			continue
		}

		newImageJSON.Performers = performer.GetNames(performers)

		tags, err := tagReader.FindByImageID(ctx, s.ID)
		if err != nil {
			logger.Errorf("[images] <%s> error getting image tag names: %v", imageHash, err)
			continue
		}

		newImageJSON.Tags = tag.GetNames(tags)

		if t.includeDependencies {
			if s.StudioID != nil {
				t.studios.IDs = sliceutil.AppendUnique(t.studios.IDs, *s.StudioID)
			}

			t.galleries.IDs = sliceutil.AppendUniques(t.galleries.IDs, gallery.GetIDs(imageGalleries))
			t.tags.IDs = sliceutil.AppendUniques(t.tags.IDs, tag.GetIDs(tags))
			t.performers.IDs = sliceutil.AppendUniques(t.performers.IDs, performer.GetIDs(performers))
		}

		fn := newImageJSON.Filename(filepath.Base(s.Path), s.Checksum)

		if err := t.json.saveImage(fn, newImageJSON); err != nil {
			logger.Errorf("[images] <%s> failed to save json: %v", imageHash, err)
		}
	}
}

func (t *ExportTask) ExportGalleries(ctx context.Context, workers int) {
	var galleriesWg sync.WaitGroup

	reader := t.repository.Gallery

	var galleries []*models.Gallery
	var err error
	all := t.full || (t.galleries != nil && t.galleries.all)
	if all {
		galleries, err = reader.All(ctx)
	} else if t.galleries != nil && len(t.galleries.IDs) > 0 {
		galleries, err = reader.FindMany(ctx, t.galleries.IDs)
	}

	if err != nil {
		logger.Errorf("[galleries] failed to fetch galleries: %v", err)
	}

	jobCh := make(chan *models.Gallery, workers*2) // make a buffered channel to feed workers

	logger.Info("[galleries] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Scene workers
		galleriesWg.Add(1)
		go t.exportGallery(ctx, &galleriesWg, jobCh)
	}

	for i, gallery := range galleries {
		index := i + 1

		if (i % 100) == 0 { // make progress easier to read
			logger.Progressf("[galleries] %d of %d", index, len(galleries))
		}

		jobCh <- gallery
	}

	close(jobCh) // close channel so that workers will know no more jobs are available
	galleriesWg.Wait()

	logger.Infof("[galleries] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportGallery(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.Gallery) {
	defer wg.Done()

	r := t.repository
	studioReader := r.Studio
	performerReader := r.Performer
	tagReader := r.Tag
	galleryChapterReader := r.GalleryChapter

	for g := range jobChan {
		if err := g.LoadFiles(ctx, r.Gallery); err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery files: %v", g.DisplayName(), err)
			continue
		}

		if err := g.LoadURLs(ctx, r.Gallery); err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery urls: %v", g.DisplayName(), err)
			continue
		}

		newGalleryJSON, err := gallery.ToBasicJSON(g)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery JSON: %v", g.DisplayName(), err)
			continue
		}

		// export files
		for _, f := range g.Files.List() {
			t.exportFile(f)
		}

		// export folder if necessary
		if g.FolderID != nil {
			folder, err := r.Folder.Find(ctx, *g.FolderID)
			if err != nil {
				logger.Errorf("[galleries] <%s> error getting gallery folder: %v", g.DisplayName(), err)
				continue
			}

			if folder == nil {
				logger.Errorf("[galleries] <%s> unable to find gallery folder", g.DisplayName())
				continue
			}

			t.exportFolder(*folder)
		}

		newGalleryJSON.Studio, err = gallery.GetStudioName(ctx, studioReader, g)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery studio name: %v", g.DisplayName(), err)
			continue
		}

		performers, err := performerReader.FindByGalleryID(ctx, g.ID)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery performer names: %v", g.DisplayName(), err)
			continue
		}

		newGalleryJSON.Performers = performer.GetNames(performers)

		tags, err := tagReader.FindByGalleryID(ctx, g.ID)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery tag names: %v", g.DisplayName(), err)
			continue
		}

		newGalleryJSON.Chapters, err = gallery.GetGalleryChaptersJSON(ctx, galleryChapterReader, g)
		if err != nil {
			logger.Errorf("[galleries] <%s> error getting gallery chapters JSON: %v", g.DisplayName(), err)
			continue
		}

		newGalleryJSON.Tags = tag.GetNames(tags)

		if t.includeDependencies {
			if g.StudioID != nil {
				t.studios.IDs = sliceutil.AppendUnique(t.studios.IDs, *g.StudioID)
			}

			t.tags.IDs = sliceutil.AppendUniques(t.tags.IDs, tag.GetIDs(tags))
			t.performers.IDs = sliceutil.AppendUniques(t.performers.IDs, performer.GetIDs(performers))
		}

		basename := ""
		// use id in case multiple galleries with the same basename
		hash := strconv.Itoa(g.ID)

		switch {
		case g.Path != "":
			basename = filepath.Base(g.Path)
		default:
			basename = g.Title
		}

		fn := newGalleryJSON.Filename(basename, hash)

		if err := t.json.saveGallery(fn, newGalleryJSON); err != nil {
			logger.Errorf("[galleries] <%s> failed to save json: %v", g.DisplayName(), err)
		}
	}
}

func (t *ExportTask) ExportPerformers(ctx context.Context, workers int) {
	var performersWg sync.WaitGroup

	reader := t.repository.Performer
	var performers []*models.Performer
	var err error
	all := t.full || (t.performers != nil && t.performers.all)
	if all {
		performers, err = reader.All(ctx)
	} else if t.performers != nil && len(t.performers.IDs) > 0 {
		performers, err = reader.FindMany(ctx, t.performers.IDs)
	}

	if err != nil {
		logger.Errorf("[performers] failed to fetch performers: %v", err)
	}
	jobCh := make(chan *models.Performer, workers*2) // make a buffered channel to feed workers

	logger.Info("[performers] exporting")
	startTime := time.Now()

	for w := 0; w < workers; w++ { // create export Performer workers
		performersWg.Add(1)
		go t.exportPerformer(ctx, &performersWg, jobCh)
	}

	for i, performer := range performers {
		index := i + 1
		logger.Progressf("[performers] %d of %d", index, len(performers))

		jobCh <- performer // feed workers
	}

	close(jobCh) // close channel so workers will know that no more jobs are available
	performersWg.Wait()

	logger.Infof("[performers] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportPerformer(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.Performer) {
	defer wg.Done()

	r := t.repository
	performerReader := r.Performer

	for p := range jobChan {
		newPerformerJSON, err := performer.ToJSON(ctx, performerReader, p)

		if err != nil {
			logger.Errorf("[performers] <%s> error getting performer JSON: %v", p.Name, err)
			continue
		}

		tags, err := r.Tag.FindByPerformerID(ctx, p.ID)
		if err != nil {
			logger.Errorf("[performers] <%s> error getting performer tags: %v", p.Name, err)
			continue
		}

		newPerformerJSON.Tags = tag.GetNames(tags)

		if t.includeDependencies {
			t.tags.IDs = sliceutil.AppendUniques(t.tags.IDs, tag.GetIDs(tags))
		}

		fn := newPerformerJSON.Filename()

		if err := t.json.savePerformer(fn, newPerformerJSON); err != nil {
			logger.Errorf("[performers] <%s> failed to save json: %v", p.Name, err)
		}
	}
}

func (t *ExportTask) ExportStudios(ctx context.Context, workers int) {
	var studiosWg sync.WaitGroup

	reader := t.repository.Studio
	var studios []*models.Studio
	var err error
	all := t.full || (t.studios != nil && t.studios.all)
	if all {
		studios, err = reader.All(ctx)
	} else if t.studios != nil && len(t.studios.IDs) > 0 {
		studios, err = reader.FindMany(ctx, t.studios.IDs)
	}

	if err != nil {
		logger.Errorf("[studios] failed to fetch studios: %v", err)
	}

	logger.Info("[studios] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Studio, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Studio workers
		studiosWg.Add(1)
		go t.exportStudio(ctx, &studiosWg, jobCh)
	}

	for i, studio := range studios {
		index := i + 1
		logger.Progressf("[studios] %d of %d", index, len(studios))

		jobCh <- studio // feed workers
	}

	close(jobCh)
	studiosWg.Wait()

	logger.Infof("[studios] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportStudio(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.Studio) {
	defer wg.Done()

	r := t.repository
	studioReader := t.repository.Studio

	for s := range jobChan {
		newStudioJSON, err := studio.ToJSON(ctx, studioReader, s)

		if err != nil {
			logger.Errorf("[studios] <%s> error getting studio JSON: %v", s.Name, err)
			continue
		}

		tags, err := r.Tag.FindByStudioID(ctx, s.ID)
		if err != nil {
			logger.Errorf("[studios] <%s> error getting studio tags: %s", s.Name, err.Error())
			continue
		}

		newStudioJSON.Tags = tag.GetNames(tags)

		if t.includeDependencies {
			t.tags.IDs = sliceutil.AppendUniques(t.tags.IDs, tag.GetIDs(tags))
		}

		fn := newStudioJSON.Filename()

		if err := t.json.saveStudio(fn, newStudioJSON); err != nil {
			logger.Errorf("[studios] <%s> failed to save json: %v", s.Name, err)
		}
	}
}

func (t *ExportTask) ExportTags(ctx context.Context, workers int) {
	var tagsWg sync.WaitGroup

	reader := t.repository.Tag
	var tags []*models.Tag
	var err error
	all := t.full || (t.tags != nil && t.tags.all)
	if all {
		tags, err = reader.All(ctx)
	} else if t.tags != nil && len(t.tags.IDs) > 0 {
		tags, err = reader.FindMany(ctx, t.tags.IDs)
	}

	if err != nil {
		logger.Errorf("[tags] failed to fetch tags: %v", err)
	}

	logger.Info("[tags] exporting")
	startTime := time.Now()

	tagIdx := 0
	if t.tags != nil {
		tagIdx = len(t.tags.IDs)
	}

	for {
		jobCh := make(chan *models.Tag, workers*2) // make a buffered channel to feed workers

		for w := 0; w < workers; w++ { // create export Tag workers
			tagsWg.Add(1)
			go t.exportTag(ctx, &tagsWg, jobCh)
		}

		for i, tag := range tags {
			index := i + 1 + tagIdx
			logger.Progressf("[tags] %d of %d", index, len(tags)+tagIdx)

			jobCh <- tag // feed workers
		}

		close(jobCh)
		tagsWg.Wait()

		// if more tags were added, we need to export those too
		if t.tags == nil || len(t.tags.IDs) == tagIdx {
			break
		}

		newTags, err := reader.FindMany(ctx, t.tags.IDs[tagIdx:])
		if err != nil {
			logger.Errorf("[tags] failed to fetch tags: %v", err)
		}

		tags = newTags
		tagIdx = len(t.tags.IDs)
	}

	logger.Infof("[tags] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportTag(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.Tag) {
	defer wg.Done()

	tagReader := t.repository.Tag

	for thisTag := range jobChan {
		newTagJSON, err := tag.ToJSON(ctx, tagReader, thisTag)

		if err != nil {
			logger.Errorf("[tags] <%s> error getting tag JSON: %v", thisTag.Name, err)
			continue
		}

		if t.includeDependencies {
			tagIDs, err := tag.GetDependentTagIDs(ctx, tagReader, thisTag)
			if err != nil {
				logger.Errorf("[tags] <%s> error getting dependent tags: %v", thisTag.Name, err)
				continue
			}
			t.tags.IDs = sliceutil.AppendUniques(t.tags.IDs, tagIDs)
		}

		fn := newTagJSON.Filename()

		if err := t.json.saveTag(fn, newTagJSON); err != nil {
			logger.Errorf("[tags] <%s> failed to save json: %v", fn, err)
		}
	}
}

func (t *ExportTask) ExportGroups(ctx context.Context, workers int) {
	var groupsWg sync.WaitGroup

	reader := t.repository.Group
	var groups []*models.Group
	var err error
	all := t.full || (t.groups != nil && t.groups.all)
	if all {
		groups, err = reader.All(ctx)
	} else if t.groups != nil && len(t.groups.IDs) > 0 {
		groups, err = reader.FindMany(ctx, t.groups.IDs)
	}

	if err != nil {
		logger.Errorf("[groups] failed to fetch groups: %v", err)
	}

	logger.Info("[groups] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.Group, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Studio workers
		groupsWg.Add(1)
		go t.exportGroup(ctx, &groupsWg, jobCh)
	}

	for i, group := range groups {
		index := i + 1
		logger.Progressf("[groups] %d of %d", index, len(groups))

		jobCh <- group // feed workers
	}

	close(jobCh)
	groupsWg.Wait()

	logger.Infof("[groups] export complete in %s. %d workers used.", time.Since(startTime), workers)

}
func (t *ExportTask) exportGroup(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.Group) {
	defer wg.Done()

	r := t.repository
	groupReader := r.Group
	studioReader := r.Studio
	tagReader := r.Tag

	for m := range jobChan {
		if err := m.LoadURLs(ctx, r.Group); err != nil {
			logger.Errorf("[groups] <%s> error getting group urls: %v", m.Name, err)
			continue
		}
		if err := m.LoadSubGroupIDs(ctx, r.Group); err != nil {
			logger.Errorf("[groups] <%s> error getting group sub-groups: %v", m.Name, err)
			continue
		}

		newGroupJSON, err := group.ToJSON(ctx, groupReader, studioReader, m)

		if err != nil {
			logger.Errorf("[groups] <%s> error getting tag JSON: %v", m.Name, err)
			continue
		}

		tags, err := tagReader.FindByGroupID(ctx, m.ID)
		if err != nil {
			logger.Errorf("[groups] <%s> error getting image tag names: %v", m.Name, err)
			continue
		}

		newGroupJSON.Tags = tag.GetNames(tags)

		subGroups := m.SubGroups.List()
		if err := func() error {
			for _, sg := range subGroups {
				subGroup, err := groupReader.Find(ctx, sg.GroupID)
				if err != nil {
					return fmt.Errorf("error getting sub group: %v", err)
				}

				newGroupJSON.SubGroups = append(newGroupJSON.SubGroups, jsonschema.SubGroupDescription{
					// TODO - this won't be unique
					Group:       subGroup.Name,
					Description: sg.Description,
				})
			}
			return nil
		}(); err != nil {
			logger.Errorf("[groups] <%s> %v", m.Name, err)
		}

		if t.includeDependencies {
			if m.StudioID != nil {
				t.studios.IDs = sliceutil.AppendUnique(t.studios.IDs, *m.StudioID)
			}
		}

		fn := newGroupJSON.Filename()

		if err := t.json.saveGroup(fn, newGroupJSON); err != nil {
			logger.Errorf("[groups] <%s> failed to save json: %v", m.Name, err)
		}
	}
}

func (t *ExportTask) ExportSavedFilters(ctx context.Context, workers int) {
	// don't export saved filters unless we're doing a full export
	if !t.full {
		return
	}

	var wg sync.WaitGroup

	reader := t.repository.SavedFilter
	var filters []*models.SavedFilter
	var err error
	filters, err = reader.All(ctx)

	if err != nil {
		logger.Errorf("[saved filters] failed to fetch saved filters: %v", err)
	}

	logger.Info("[saved filters] exporting")
	startTime := time.Now()

	jobCh := make(chan *models.SavedFilter, workers*2) // make a buffered channel to feed workers

	for w := 0; w < workers; w++ { // create export Saved Filter workers
		wg.Add(1)
		go t.exportSavedFilter(ctx, &wg, jobCh)
	}

	for i, savedFilter := range filters {
		index := i + 1
		logger.Progressf("[saved filters] %d of %d", index, len(filters))

		jobCh <- savedFilter // feed workers
	}

	close(jobCh)
	wg.Wait()

	logger.Infof("[saved filters] export complete in %s. %d workers used.", time.Since(startTime), workers)
}

func (t *ExportTask) exportSavedFilter(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan *models.SavedFilter) {
	defer wg.Done()

	for thisFilter := range jobChan {
		newJSON, err := savedfilter.ToJSON(ctx, thisFilter)

		if err != nil {
			logger.Errorf("[saved filter] <%s> error getting saved filter JSON: %v", thisFilter.Name, err)
			continue
		}

		fn := newJSON.Filename()

		if err := t.json.saveSavedFilter(fn, newJSON); err != nil {
			logger.Errorf("[saved filter] <%s> failed to save json: %v", fn, err)
		}
	}
}
