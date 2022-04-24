//go:build integration
// +build integration

package autotag

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const testName = "Foo's Bar"
const existingStudioName = "ExistingStudio"

const existingStudioSceneName = testName + ".dontChangeStudio.mp4"
const existingStudioImageName = testName + ".dontChangeStudio.mp4"
const existingStudioGalleryName = testName + ".dontChangeStudio.mp4"

var existingStudioID int

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
	db = &sqlite.Database{}
	if err := db.Open(databaseFile); err != nil {
		panic(fmt.Sprintf("Could not initialize database: %s", err.Error()))
	}

	r = db.TxnRepository()

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
	ret := runTests(m)
	os.Exit(ret)
}

func createPerformer(ctx context.Context, pqb models.PerformerWriter) error {
	// create the performer
	performer := models.Performer{
		Checksum: testName,
		Name:     sql.NullString{Valid: true, String: testName},
		Favorite: sql.NullBool{Valid: true, Bool: false},
	}

	_, err := pqb.Create(ctx, performer)
	if err != nil {
		return err
	}

	return nil
}

func createStudio(ctx context.Context, qb models.StudioWriter, name string) (*models.Studio, error) {
	// create the studio
	studio := models.Studio{
		Checksum: name,
		Name:     sql.NullString{Valid: true, String: name},
	}

	return qb.Create(ctx, studio)
}

func createTag(ctx context.Context, qb models.TagWriter) error {
	// create the studio
	tag := models.Tag{
		Name: testName,
	}

	_, err := qb.Create(ctx, tag)
	if err != nil {
		return err
	}

	return nil
}

func createScenes(ctx context.Context, sqb models.SceneReaderWriter) error {
	// create the scenes
	scenePatterns, falseScenePatterns := generateTestPaths(testName, sceneExt)

	for _, fn := range scenePatterns {
		err := createScene(ctx, sqb, makeScene(fn, true))
		if err != nil {
			return err
		}
	}
	for _, fn := range falseScenePatterns {
		err := createScene(ctx, sqb, makeScene(fn, false))
		if err != nil {
			return err
		}
	}

	// add organized scenes
	for _, fn := range scenePatterns {
		s := makeScene("organized"+fn, false)
		s.Organized = true
		err := createScene(ctx, sqb, s)
		if err != nil {
			return err
		}
	}

	// create scene with existing studio io
	studioScene := makeScene(existingStudioSceneName, true)
	studioScene.StudioID = &existingStudioID
	err := createScene(ctx, sqb, studioScene)
	if err != nil {
		return err
	}

	return nil
}

func makeScene(name string, expectedResult bool) *models.Scene {
	checksum := md5.FromString(name)
	scene := &models.Scene{
		Checksum: &checksum,
		Path:     name,
	}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		scene.Title = name
	}

	return scene
}

func createScene(ctx context.Context, sqb models.SceneWriter, scene *models.Scene) error {
	err := sqb.Create(ctx, scene)

	if err != nil {
		return fmt.Errorf("Failed to create scene with name '%s': %s", scene.Path, err.Error())
	}

	return nil
}

func createImages(ctx context.Context, sqb models.ImageReaderWriter) error {
	// create the images
	imagePatterns, falseImagePatterns := generateTestPaths(testName, imageExt)

	for _, fn := range imagePatterns {
		err := createImage(ctx, sqb, makeImage(fn, true))
		if err != nil {
			return err
		}
	}
	for _, fn := range falseImagePatterns {
		err := createImage(ctx, sqb, makeImage(fn, false))
		if err != nil {
			return err
		}
	}

	// add organized images
	for _, fn := range imagePatterns {
		s := makeImage("organized"+fn, false)
		s.Organized = true
		err := createImage(ctx, sqb, s)
		if err != nil {
			return err
		}
	}

	// create image with existing studio io
	studioImage := makeImage(existingStudioImageName, true)
	studioImage.StudioID = &existingStudioID
	err := createImage(ctx, sqb, studioImage)
	if err != nil {
		return err
	}

	return nil
}

func makeImage(name string, expectedResult bool) *models.Image {
	image := &models.Image{
		Checksum: md5.FromString(name),
		Path:     name,
	}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		image.Title = name
	}

	return image
}

func createImage(ctx context.Context, sqb models.ImageWriter, image *models.Image) error {
	if err := sqb.Create(ctx, image); err != nil {
		return fmt.Errorf("Failed to create image with name '%s': %s", image.Path, err.Error())
	}

	return nil
}

func createGalleries(ctx context.Context, sqb models.GalleryReaderWriter) error {
	// create the galleries
	galleryPatterns, falseGalleryPatterns := generateTestPaths(testName, galleryExt)

	for _, fn := range galleryPatterns {
		err := createGallery(ctx, sqb, makeGallery(fn, true))
		if err != nil {
			return err
		}
	}
	for _, fn := range falseGalleryPatterns {
		err := createGallery(ctx, sqb, makeGallery(fn, false))
		if err != nil {
			return err
		}
	}

	// add organized galleries
	for _, fn := range galleryPatterns {
		s := makeGallery("organized"+fn, false)
		s.Organized = true
		err := createGallery(ctx, sqb, s)
		if err != nil {
			return err
		}
	}

	// create gallery with existing studio io
	studioGallery := makeGallery(existingStudioGalleryName, true)
	studioGallery.StudioID = &existingStudioID
	err := createGallery(ctx, sqb, studioGallery)
	if err != nil {
		return err
	}

	return nil
}

func makeGallery(name string, expectedResult bool) *models.Gallery {
	gallery := &models.Gallery{
		Checksum: md5.FromString(name),
		Path:     &name,
	}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		gallery.Title = name
	}

	return gallery
}

func createGallery(ctx context.Context, sqb models.GalleryWriter, gallery *models.Gallery) error {
	err := sqb.Create(ctx, gallery)

	if err != nil {
		return fmt.Errorf("Failed to create gallery with name '%s': %s", *gallery.Path, err.Error())
	}

	return nil
}

func withTxn(f func(ctx context.Context) error) error {
	return txn.WithTxn(context.TODO(), db, f)
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

		err = createScenes(ctx, r.Scene)
		if err != nil {
			return err
		}

		err = createImages(ctx, r.Image)
		if err != nil {
			return err
		}

		err = createGalleries(ctx, r.Gallery)
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

	for _, p := range performers {
		if err := withTxn(func(ctx context.Context) error {
			return PerformerScenes(ctx, p, nil, r.Scene, nil)
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
			if scene.Title == scene.Path && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, scene.Path)
			} else if scene.Title != scene.Path && len(performers) > 0 {
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

	for _, s := range studios {
		if err := withTxn(func(ctx context.Context) error {
			aliases, err := r.Studio.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return StudioScenes(ctx, s, nil, aliases, r.Scene, nil)
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
			if scene.Path == existingStudioSceneName {
				if scene.StudioID == nil || *scene.StudioID != existingStudioID {
					t.Error("Incorrectly overwrote studio ID for scene with existing studio ID")
				}
			} else {
				// title is only set on scenes where we expect studio to be set
				if scene.Title == scene.Path {
					if scene.StudioID == nil {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, scene.Path)
					} else if scene.StudioID != nil && *scene.StudioID != studios[1].ID {
						t.Errorf("Incorrect studio id %d set for path '%s'", scene.StudioID, scene.Path)
					}

				} else if scene.Title != scene.Path && scene.StudioID != nil && *scene.StudioID == studios[1].ID {
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

	for _, s := range tags {
		if err := withTxn(func(ctx context.Context) error {
			aliases, err := r.Tag.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return TagScenes(ctx, s, nil, aliases, r.Scene, nil)
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
			if scene.Title == scene.Path && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, scene.Path)
			} else if (scene.Title != scene.Path) && len(tags) > 0 {
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

	for _, p := range performers {
		if err := withTxn(func(ctx context.Context) error {
			return PerformerImages(ctx, p, nil, r.Image, nil)
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
			if image.Title == image.Path && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, image.Path)
			} else if image.Title != image.Path && len(performers) > 0 {
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

	for _, s := range studios {
		if err := withTxn(func(ctx context.Context) error {
			aliases, err := r.Studio.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return StudioImages(ctx, s, nil, aliases, r.Image, nil)
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
			if image.Path == existingStudioImageName {
				if *image.StudioID != existingStudioID {
					t.Error("Incorrectly overwrote studio ID for image with existing studio ID")
				}
			} else {
				// title is only set on images where we expect studio to be set
				if image.Title == image.Path {
					if image.StudioID == nil {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, image.Path)
					} else if *image.StudioID != studios[1].ID {
						t.Errorf("Incorrect studio id %d set for path '%s'", *image.StudioID, image.Path)
					}

				} else if image.Title != image.Path && image.StudioID != nil && *image.StudioID == studios[1].ID {
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

	for _, s := range tags {
		if err := withTxn(func(ctx context.Context) error {
			aliases, err := r.Tag.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return TagImages(ctx, s, nil, aliases, r.Image, nil)
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
			if image.Title == image.Path && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, image.Path)
			} else if image.Title != image.Path && len(tags) > 0 {
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

	for _, p := range performers {
		if err := withTxn(func(ctx context.Context) error {
			return PerformerGalleries(ctx, p, nil, r.Gallery, nil)
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
			if gallery.Title == *gallery.Path && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, *gallery.Path)
			} else if gallery.Title != *gallery.Path && len(performers) > 0 {
				t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, *gallery.Path)
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

	for _, s := range studios {
		if err := withTxn(func(ctx context.Context) error {
			aliases, err := r.Studio.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return StudioGalleries(ctx, s, nil, aliases, r.Gallery, nil)
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
			if gallery.Path != nil && *gallery.Path == existingStudioGalleryName {
				if *gallery.StudioID != existingStudioID {
					t.Error("Incorrectly overwrote studio ID for gallery with existing studio ID")
				}
			} else {
				// title is only set on galleries where we expect studio to be set
				if gallery.Title == *gallery.Path {
					if gallery.StudioID == nil {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, *gallery.Path)
					} else if *gallery.StudioID != studios[1].ID {
						t.Errorf("Incorrect studio id %d set for path '%s'", *gallery.StudioID, *gallery.Path)
					}

				} else if gallery.Title != *gallery.Path && (gallery.StudioID != nil && *gallery.StudioID == studios[1].ID) {
					t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, *gallery.Path)
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

	for _, s := range tags {
		if err := withTxn(func(ctx context.Context) error {
			aliases, err := r.Tag.GetAliases(ctx, s.ID)
			if err != nil {
				return err
			}

			return TagGalleries(ctx, s, nil, aliases, r.Gallery, nil)
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
			if gallery.Title == *gallery.Path && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, *gallery.Path)
			} else if gallery.Title != *gallery.Path && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, *gallery.Path)
			}
		}

		return nil
	})
}
