//go:build integration
// +build integration

package autotag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// necessary to register custom migrations
	_ "github.com/stashapp/stash/pkg/sqlite/migrations"
)

const testName = "Foo's Bar"
const existingStudioName = "ExistingStudio"

const existingStudioSceneName = testName + ".dontChangeStudio.mp4"
const existingStudioImageName = testName + ".dontChangeStudio.png"
const existingStudioGalleryName = testName + ".dontChangeStudio.zip"

var existingStudioID int

const expectedMatchTitle = "expected match"

var db *sqlite.Database
var r models.Repository

func testTeardown(databaseFile string) {
	err := db.Close()

	if err != nil {
		panic(err)
	}

	err = os.Remove(databaseFile)
	if err != nil {
		panic(err)
	}
}

func runTests(m *testing.M) int {
	// create the database file
	f, err := os.CreateTemp("", "*.sqlite")
	if err != nil {
		panic(fmt.Sprintf("Could not create temporary file: %s", err.Error()))
	}

	f.Close()
	databaseFile := f.Name()
	db = sqlite.NewDatabase()
	if err := db.Open(databaseFile); err != nil {
		panic(fmt.Sprintf("Could not initialize database: %s", err.Error()))
	}

	r = db.Repository()

	// defer close and delete the database
	defer testTeardown(databaseFile)

	err = populateDB()
	if err != nil {
		panic(fmt.Sprintf("Could not populate database: %s", err.Error()))
	} else {
		// run the tests
		return m.Run()
	}
}

func TestMain(m *testing.M) {
	// initialise empty config - needed by some db migrations
	_ = config.InitializeEmpty()

	ret := runTests(m)
	os.Exit(ret)
}

func createPerformer(ctx context.Context, pqb models.PerformerWriter) error {
	// create the performer
	performer := models.Performer{
		Name: testName,
	}

	err := pqb.Create(ctx, &models.CreatePerformerInput{Performer: &performer})
	if err != nil {
		return err
	}

	return nil
}

func createStudio(ctx context.Context, qb models.StudioWriter, name string) (*models.Studio, error) {
	// create the studio
	studio := models.Studio{
		Name: name,
	}

	err := qb.Create(ctx, &studio)
	if err != nil {
		return nil, err
	}

	return &studio, nil
}

func createTag(ctx context.Context, qb models.TagWriter) error {
	// create the studio
	tag := models.Tag{
		Name: testName,
	}

	err := qb.Create(ctx, &tag)
	if err != nil {
		return err
	}

	return nil
}

func createScenes(ctx context.Context, sqb models.SceneReaderWriter, folderStore models.FolderFinderCreator, fileCreator models.FileCreator) error {
	// create the scenes
	scenePatterns, falseScenePatterns := generateTestPaths(testName, sceneExt)

	for _, fn := range scenePatterns {
		f, err := createSceneFile(ctx, fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = true
		if err := createScene(ctx, sqb, makeScene(expectedResult), f); err != nil {
			return err
		}
	}

	for _, fn := range falseScenePatterns {
		f, err := createSceneFile(ctx, fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = false
		if err := createScene(ctx, sqb, makeScene(expectedResult), f); err != nil {
			return err
		}
	}

	// add organized scenes
	for _, fn := range scenePatterns {
		f, err := createSceneFile(ctx, "organized"+fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = false
		s := makeScene(expectedResult)
		s.Organized = true
		if err := createScene(ctx, sqb, s, f); err != nil {
			return err
		}
	}

	// create scene with existing studio io
	f, err := createSceneFile(ctx, existingStudioSceneName, folderStore, fileCreator)
	if err != nil {
		return err
	}

	s := &models.Scene{
		Title:    expectedMatchTitle,
		Code:     existingStudioSceneName,
		StudioID: &existingStudioID,
	}
	if err := createScene(ctx, sqb, s, f); err != nil {
		return err
	}

	return nil
}

func makeScene(expectedResult bool) *models.Scene {
	s := &models.Scene{}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		s.Title = expectedMatchTitle
	}

	return s
}

func createSceneFile(ctx context.Context, name string, folderStore models.FolderFinderCreator, fileCreator models.FileCreator) (*models.VideoFile, error) {
	folderPath := filepath.Dir(name)
	basename := filepath.Base(name)

	folder, err := getOrCreateFolder(ctx, folderStore, folderPath)
	if err != nil {
		return nil, err
	}

	folderID := folder.ID

	f := &models.VideoFile{
		BaseFile: &models.BaseFile{
			Basename:       basename,
			ParentFolderID: folderID,
		},
	}

	if err := fileCreator.Create(ctx, f); err != nil {
		return nil, fmt.Errorf("creating scene file %q: %w", name, err)
	}

	return f, nil
}

func getOrCreateFolder(ctx context.Context, folderStore models.FolderFinderCreator, folderPath string) (*models.Folder, error) {
	f, err := folderStore.FindByPath(ctx, folderPath, true)
	if err != nil {
		return nil, fmt.Errorf("getting folder by path: %w", err)
	}

	if f != nil {
		return f, nil
	}

	var parentID models.FolderID
	dir := filepath.Dir(folderPath)
	if dir != "." {
		parent, err := getOrCreateFolder(ctx, folderStore, dir)
		if err != nil {
			return nil, err
		}

		parentID = parent.ID
	}

	f = &models.Folder{
		Path: folderPath,
	}

	if parentID != 0 {
		f.ParentFolderID = &parentID
	}

	if err := folderStore.Create(ctx, f); err != nil {
		return nil, fmt.Errorf("creating folder: %w", err)
	}

	return f, nil
}

func createScene(ctx context.Context, sqb models.SceneWriter, s *models.Scene, f *models.VideoFile) error {
	err := sqb.Create(ctx, s, []models.FileID{f.ID})

	if err != nil {
		return fmt.Errorf("Failed to create scene with path '%s': %s", f.Path, err.Error())
	}

	return nil
}

func createImages(ctx context.Context, w models.ImageReaderWriter, folderStore models.FolderFinderCreator, fileCreator models.FileCreator) error {
	// create the images
	imagePatterns, falseImagePatterns := generateTestPaths(testName, imageExt)

	for _, fn := range imagePatterns {
		f, err := createImageFile(ctx, fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = true
		if err := createImage(ctx, w, makeImage(expectedResult), f); err != nil {
			return err
		}
	}
	for _, fn := range falseImagePatterns {
		f, err := createImageFile(ctx, fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = false
		if err := createImage(ctx, w, makeImage(expectedResult), f); err != nil {
			return err
		}
	}

	// add organized images
	for _, fn := range imagePatterns {
		f, err := createImageFile(ctx, "organized"+fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = false
		s := makeImage(expectedResult)
		s.Organized = true
		if err := createImage(ctx, w, s, f); err != nil {
			return err
		}
	}

	// create image with existing studio io
	f, err := createImageFile(ctx, existingStudioImageName, folderStore, fileCreator)
	if err != nil {
		return err
	}

	s := &models.Image{
		Title:    existingStudioImageName,
		StudioID: &existingStudioID,
	}
	if err := createImage(ctx, w, s, f); err != nil {
		return err
	}

	return nil
}

func createImageFile(ctx context.Context, name string, folderStore models.FolderFinderCreator, fileCreator models.FileCreator) (*models.ImageFile, error) {
	folderPath := filepath.Dir(name)
	basename := filepath.Base(name)

	folder, err := getOrCreateFolder(ctx, folderStore, folderPath)
	if err != nil {
		return nil, err
	}

	folderID := folder.ID

	f := &models.ImageFile{
		BaseFile: &models.BaseFile{
			Basename:       basename,
			ParentFolderID: folderID,
		},
	}

	if err := fileCreator.Create(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

func makeImage(expectedResult bool) *models.Image {
	o := &models.Image{}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		o.Title = expectedMatchTitle
	}

	return o
}

func createImage(ctx context.Context, w models.ImageWriter, o *models.Image, f *models.ImageFile) error {
	err := w.Create(ctx, o, []models.FileID{f.ID})

	if err != nil {
		return fmt.Errorf("Failed to create image with path '%s': %s", f.Path, err.Error())
	}

	return nil
}

func createGalleries(ctx context.Context, w models.GalleryReaderWriter, folderStore models.FolderFinderCreator, fileCreator models.FileCreator) error {
	// create the galleries
	galleryPatterns, falseGalleryPatterns := generateTestPaths(testName, galleryExt)

	for _, fn := range galleryPatterns {
		f, err := createGalleryFile(ctx, fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = true
		if err := createGallery(ctx, w, makeGallery(expectedResult), f); err != nil {
			return err
		}
	}
	for _, fn := range falseGalleryPatterns {
		f, err := createGalleryFile(ctx, fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = false
		if err := createGallery(ctx, w, makeGallery(expectedResult), f); err != nil {
			return err
		}
	}

	// add organized galleries
	for _, fn := range galleryPatterns {
		f, err := createGalleryFile(ctx, "organized"+fn, folderStore, fileCreator)
		if err != nil {
			return err
		}

		const expectedResult = false
		s := makeGallery(expectedResult)
		s.Organized = true
		if err := createGallery(ctx, w, s, f); err != nil {
			return err
		}
	}

	// create gallery with existing studio io
	f, err := createGalleryFile(ctx, existingStudioGalleryName, folderStore, fileCreator)
	if err != nil {
		return err
	}

	s := &models.Gallery{
		Title:    existingStudioGalleryName,
		StudioID: &existingStudioID,
	}
	if err := createGallery(ctx, w, s, f); err != nil {
		return err
	}

	return nil
}

func createGalleryFile(ctx context.Context, name string, folderStore models.FolderFinderCreator, fileCreator models.FileCreator) (*models.BaseFile, error) {
	folderPath := filepath.Dir(name)
	basename := filepath.Base(name)

	folder, err := getOrCreateFolder(ctx, folderStore, folderPath)
	if err != nil {
		return nil, err
	}

	folderID := folder.ID

	f := &models.BaseFile{
		Basename:       basename,
		ParentFolderID: folderID,
	}

	if err := fileCreator.Create(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

func makeGallery(expectedResult bool) *models.Gallery {
	o := &models.Gallery{}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		o.Title = expectedMatchTitle
	}

	return o
}

func createGallery(ctx context.Context, w models.GalleryWriter, o *models.Gallery, f *models.BaseFile) error {
	err := w.Create(ctx, o, []models.FileID{f.ID})
	if err != nil {
		return fmt.Errorf("Failed to create gallery with path '%s': %s", f.Path, err.Error())
	}

	return nil
}

func withTxn(f func(ctx context.Context) error) error {
	return txn.WithTxn(testCtx, db, f)
}

func withDB(f func(ctx context.Context) error) error {
	return txn.WithDatabase(testCtx, db, f)
}

func populateDB() error {
	if err := withTxn(func(ctx context.Context) error {
		err := createPerformer(ctx, r.Performer)
		if err != nil {
			return err
		}

		_, err = createStudio(ctx, r.Studio, testName)
		if err != nil {
			return err
		}

		// create existing studio
		existingStudio, err := createStudio(ctx, r.Studio, existingStudioName)
		if err != nil {
			return err
		}

		existingStudioID = existingStudio.ID

		err = createTag(ctx, r.Tag)
		if err != nil {
			return err
		}

		err = createScenes(ctx, r.Scene, r.Folder, r.File)
		if err != nil {
			return err
		}

		err = createImages(ctx, r.Image, r.Folder, r.File)
		if err != nil {
			return err
		}

		err = createGalleries(ctx, r.Gallery, r.Folder, r.File)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func TestParsePerformerScenes(t *testing.T) {
	var performers []*models.Performer
	if err := withTxn(func(ctx context.Context) error {
		var err error
		performers, err = r.Performer.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, p := range performers {
		if err := withDB(func(ctx context.Context) error {
			if err := p.LoadAliases(ctx, r.Performer); err != nil {
				return err
			}
			return tagger.PerformerScenes(ctx, p, nil, r.Scene)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that scenes were tagged correctly
	withTxn(func(ctx context.Context) error {
		pqb := r.Performer

		scenes, err := r.Scene.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		for _, scene := range scenes {
			performers, err := pqb.FindBySceneID(ctx, scene.ID)

			if err != nil {
				t.Errorf("Error getting scene performers: %s", err.Error())
			}

			// title is only set on scenes where we expect performer to be set
			if scene.Title == expectedMatchTitle && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, scene.Path)
			} else if scene.Title != expectedMatchTitle && len(performers) > 0 {
				t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, scene.Path)
			}
		}

		return nil
	})
}

func TestParseStudioScenes(t *testing.T) {
	var studios []*models.Studio
	if err := withTxn(func(ctx context.Context) error {
		var err error
		studios, err = r.Studio.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting studio: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, s := range studios {
		if err := withDB(func(ctx context.Context) error {
			aliases, err := r.Studio.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return tagger.StudioScenes(ctx, s, nil, aliases, r.Scene)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that scenes were tagged correctly
	withTxn(func(ctx context.Context) error {
		scenes, err := r.Scene.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		for _, scene := range scenes {
			// check for existing studio id scene first
			if scene.Code == existingStudioSceneName {
				if scene.StudioID == nil || *scene.StudioID != existingStudioID {
					t.Error("Incorrectly overwrote studio ID for scene with existing studio ID")
				}
			} else {
				// title is only set on scenes where we expect studio to be set
				if scene.Title == expectedMatchTitle {
					if scene.StudioID == nil {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, scene.Path)
					} else if scene.StudioID != nil && *scene.StudioID != studios[1].ID {
						t.Errorf("Incorrect studio id %d set for path '%s'", scene.StudioID, scene.Path)
					}

				} else if scene.Title != expectedMatchTitle && scene.StudioID != nil && *scene.StudioID == studios[1].ID {
					t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, scene.Path)
				}
			}
		}

		return nil
	})
}

func TestParseTagScenes(t *testing.T) {
	var tags []*models.Tag
	if err := withTxn(func(ctx context.Context) error {
		var err error
		tags, err = r.Tag.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, s := range tags {
		if err := withDB(func(ctx context.Context) error {
			aliases, err := r.Tag.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return tagger.TagScenes(ctx, s, nil, aliases, r.Scene)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that scenes were tagged correctly
	withTxn(func(ctx context.Context) error {
		scenes, err := r.Scene.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		tqb := r.Tag

		for _, scene := range scenes {
			tags, err := tqb.FindBySceneID(ctx, scene.ID)

			if err != nil {
				t.Errorf("Error getting scene tags: %s", err.Error())
			}

			// title is only set on scenes where we expect tag to be set
			if scene.Title == expectedMatchTitle && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, scene.Path)
			} else if (scene.Title != expectedMatchTitle) && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, scene.Path)
			}
		}

		return nil
	})
}

func TestParsePerformerImages(t *testing.T) {
	var performers []*models.Performer
	if err := withTxn(func(ctx context.Context) error {
		var err error
		performers, err = r.Performer.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, p := range performers {
		if err := withDB(func(ctx context.Context) error {
			if err := p.LoadAliases(ctx, r.Performer); err != nil {
				return err
			}
			return tagger.PerformerImages(ctx, p, nil, r.Image)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that images were tagged correctly
	withTxn(func(ctx context.Context) error {
		pqb := r.Performer

		images, err := r.Image.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		for _, image := range images {
			performers, err := pqb.FindByImageID(ctx, image.ID)

			if err != nil {
				t.Errorf("Error getting image performers: %s", err.Error())
			}

			// title is only set on images where we expect performer to be set
			expectedMatch := image.Title == expectedMatchTitle || image.Title == existingStudioImageName
			if expectedMatch && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, image.Path)
			} else if !expectedMatch && len(performers) > 0 {
				t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, image.Path)
			}
		}

		return nil
	})
}

func TestParseStudioImages(t *testing.T) {
	var studios []*models.Studio
	if err := withTxn(func(ctx context.Context) error {
		var err error
		studios, err = r.Studio.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting studio: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, s := range studios {
		if err := withDB(func(ctx context.Context) error {
			aliases, err := r.Studio.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return tagger.StudioImages(ctx, s, nil, aliases, r.Image)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that images were tagged correctly
	withTxn(func(ctx context.Context) error {
		images, err := r.Image.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		for _, image := range images {
			// check for existing studio id image first
			if image.Title == existingStudioImageName {
				if *image.StudioID != existingStudioID {
					t.Error("Incorrectly overwrote studio ID for image with existing studio ID")
				}
			} else {
				// title is only set on images where we expect studio to be set
				if image.Title == expectedMatchTitle {
					if image.StudioID == nil {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, image.Path)
					} else if *image.StudioID != studios[1].ID {
						t.Errorf("Incorrect studio id %d set for path '%s'", *image.StudioID, image.Path)
					}

				} else if image.Title != expectedMatchTitle && image.StudioID != nil && *image.StudioID == studios[1].ID {
					t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, image.Path)
				}
			}
		}

		return nil
	})
}

func TestParseTagImages(t *testing.T) {
	var tags []*models.Tag
	if err := withTxn(func(ctx context.Context) error {
		var err error
		tags, err = r.Tag.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, s := range tags {
		if err := withDB(func(ctx context.Context) error {
			aliases, err := r.Tag.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return tagger.TagImages(ctx, s, nil, aliases, r.Image)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that images were tagged correctly
	withTxn(func(ctx context.Context) error {
		images, err := r.Image.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		tqb := r.Tag

		for _, image := range images {
			tags, err := tqb.FindByImageID(ctx, image.ID)

			if err != nil {
				t.Errorf("Error getting image tags: %s", err.Error())
			}

			// title is only set on images where we expect performer to be set
			expectedMatch := image.Title == expectedMatchTitle || image.Title == existingStudioImageName
			if expectedMatch && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, image.Path)
			} else if !expectedMatch && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, image.Path)
			}
		}

		return nil
	})
}

func TestParsePerformerGalleries(t *testing.T) {
	var performers []*models.Performer
	if err := withTxn(func(ctx context.Context) error {
		var err error
		performers, err = r.Performer.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, p := range performers {
		if err := withDB(func(ctx context.Context) error {
			if err := p.LoadAliases(ctx, r.Performer); err != nil {
				return err
			}
			return tagger.PerformerGalleries(ctx, p, nil, r.Gallery)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that galleries were tagged correctly
	withTxn(func(ctx context.Context) error {
		pqb := r.Performer

		galleries, err := r.Gallery.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		for _, gallery := range galleries {
			performers, err := pqb.FindByGalleryID(ctx, gallery.ID)

			if err != nil {
				t.Errorf("Error getting gallery performers: %s", err.Error())
			}

			// title is only set on galleries where we expect performer to be set
			expectedMatch := gallery.Title == expectedMatchTitle || gallery.Title == existingStudioGalleryName
			if expectedMatch && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, gallery.Path)
			} else if !expectedMatch && len(performers) > 0 {
				t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, gallery.Path)
			}
		}

		return nil
	})
}

func TestParseStudioGalleries(t *testing.T) {
	var studios []*models.Studio
	if err := withTxn(func(ctx context.Context) error {
		var err error
		studios, err = r.Studio.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting studio: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, s := range studios {
		if err := withDB(func(ctx context.Context) error {
			aliases, err := r.Studio.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return tagger.StudioGalleries(ctx, s, nil, aliases, r.Gallery)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that galleries were tagged correctly
	withTxn(func(ctx context.Context) error {
		galleries, err := r.Gallery.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		for _, gallery := range galleries {
			// check for existing studio id gallery first
			if gallery.Title == existingStudioGalleryName {
				if *gallery.StudioID != existingStudioID {
					t.Error("Incorrectly overwrote studio ID for gallery with existing studio ID")
				}
			} else {
				// title is only set on galleries where we expect studio to be set
				if gallery.Title == expectedMatchTitle {
					if gallery.StudioID == nil {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, gallery.Path)
					} else if *gallery.StudioID != studios[1].ID {
						t.Errorf("Incorrect studio id %d set for path '%s'", *gallery.StudioID, gallery.Path)
					}

				} else if gallery.Title != expectedMatchTitle && (gallery.StudioID != nil && *gallery.StudioID == studios[1].ID) {
					t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, gallery.Path)
				}
			}
		}

		return nil
	})
}

func TestParseTagGalleries(t *testing.T) {
	var tags []*models.Tag
	if err := withTxn(func(ctx context.Context) error {
		var err error
		tags, err = r.Tag.All(ctx)
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	tagger := Tagger{
		TxnManager: db,
	}

	for _, s := range tags {
		if err := withDB(func(ctx context.Context) error {
			aliases, err := r.Tag.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return tagger.TagGalleries(ctx, s, nil, aliases, r.Gallery)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that galleries were tagged correctly
	withTxn(func(ctx context.Context) error {
		galleries, err := r.Gallery.All(ctx)
		if err != nil {
			t.Error(err.Error())
		}

		tqb := r.Tag

		for _, gallery := range galleries {
			tags, err := tqb.FindByGalleryID(ctx, gallery.ID)

			if err != nil {
				t.Errorf("Error getting gallery tags: %s", err.Error())
			}

			// title is only set on galleries where we expect performer to be set
			expectedMatch := gallery.Title == expectedMatchTitle || gallery.Title == existingStudioGalleryName
			if expectedMatch && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, gallery.Path)
			} else if !expectedMatch && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, gallery.Path)
			}
		}

		return nil
	})
}
