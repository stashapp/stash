//go:build integration
// +build integration

package autotag

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const testName = "Foo's Bar"
const existingStudioName = "ExistingStudio"

const existingStudioSceneName = testName + ".dontChangeStudio.mp4"
const existingStudioImageName = testName + ".dontChangeStudio.mp4"
const existingStudioGalleryName = testName + ".dontChangeStudio.mp4"

var existingStudioID int

func testTeardown(databaseFile string) {
	err := database.DB.Close()

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
	if err := database.Initialize(databaseFile); err != nil {
		panic(fmt.Sprintf("Could not initialize database: %s", err.Error()))
	}

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

func createPerformer(pqb models.PerformerWriter) error {
	// create the performer
	performer := models.Performer{
		Checksum: testName,
		Name:     sql.NullString{Valid: true, String: testName},
		Favorite: sql.NullBool{Valid: true, Bool: false},
	}

	_, err := pqb.Create(performer)
	if err != nil {
		return err
	}

	return nil
}

func createStudio(qb models.StudioWriter, name string) (*models.Studio, error) {
	// create the studio
	studio := models.Studio{
		Checksum: name,
		Name:     sql.NullString{Valid: true, String: name},
	}

	return qb.Create(studio)
}

func createTag(qb models.TagWriter) error {
	// create the studio
	tag := models.Tag{
		Name: testName,
	}

	_, err := qb.Create(tag)
	if err != nil {
		return err
	}

	return nil
}

func createScenes(sqb models.SceneReaderWriter) error {
	// create the scenes
	scenePatterns, falseScenePatterns := generateTestPaths(testName, sceneExt)

	for _, fn := range scenePatterns {
		err := createScene(sqb, makeScene(fn, true))
		if err != nil {
			return err
		}
	}
	for _, fn := range falseScenePatterns {
		err := createScene(sqb, makeScene(fn, false))
		if err != nil {
			return err
		}
	}

	// add organized scenes
	for _, fn := range scenePatterns {
		s := makeScene("organized"+fn, false)
		s.Organized = true
		err := createScene(sqb, s)
		if err != nil {
			return err
		}
	}

	// create scene with existing studio io
	studioScene := makeScene(existingStudioSceneName, true)
	studioScene.StudioID = sql.NullInt64{Valid: true, Int64: int64(existingStudioID)}
	err := createScene(sqb, studioScene)
	if err != nil {
		return err
	}

	return nil
}

func makeScene(name string, expectedResult bool) *models.Scene {
	scene := &models.Scene{
		Checksum: sql.NullString{String: md5.FromString(name), Valid: true},
		Path:     name,
	}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		scene.Title = sql.NullString{Valid: true, String: name}
	}

	return scene
}

func createScene(sqb models.SceneWriter, scene *models.Scene) error {
	_, err := sqb.Create(*scene)

	if err != nil {
		return fmt.Errorf("Failed to create scene with name '%s': %s", scene.Path, err.Error())
	}

	return nil
}

func createImages(sqb models.ImageReaderWriter) error {
	// create the images
	imagePatterns, falseImagePatterns := generateTestPaths(testName, imageExt)

	for _, fn := range imagePatterns {
		err := createImage(sqb, makeImage(fn, true))
		if err != nil {
			return err
		}
	}
	for _, fn := range falseImagePatterns {
		err := createImage(sqb, makeImage(fn, false))
		if err != nil {
			return err
		}
	}

	// add organized images
	for _, fn := range imagePatterns {
		s := makeImage("organized"+fn, false)
		s.Organized = true
		err := createImage(sqb, s)
		if err != nil {
			return err
		}
	}

	// create image with existing studio io
	studioImage := makeImage(existingStudioImageName, true)
	studioImage.StudioID = sql.NullInt64{Valid: true, Int64: int64(existingStudioID)}
	err := createImage(sqb, studioImage)
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
		image.Title = sql.NullString{Valid: true, String: name}
	}

	return image
}

func createImage(sqb models.ImageWriter, image *models.Image) error {
	_, err := sqb.Create(*image)

	if err != nil {
		return fmt.Errorf("Failed to create image with name '%s': %s", image.Path, err.Error())
	}

	return nil
}

func createGalleries(sqb models.GalleryReaderWriter) error {
	// create the galleries
	galleryPatterns, falseGalleryPatterns := generateTestPaths(testName, galleryExt)

	for _, fn := range galleryPatterns {
		err := createGallery(sqb, makeGallery(fn, true))
		if err != nil {
			return err
		}
	}
	for _, fn := range falseGalleryPatterns {
		err := createGallery(sqb, makeGallery(fn, false))
		if err != nil {
			return err
		}
	}

	// add organized galleries
	for _, fn := range galleryPatterns {
		s := makeGallery("organized"+fn, false)
		s.Organized = true
		err := createGallery(sqb, s)
		if err != nil {
			return err
		}
	}

	// create gallery with existing studio io
	studioGallery := makeGallery(existingStudioGalleryName, true)
	studioGallery.StudioID = sql.NullInt64{Valid: true, Int64: int64(existingStudioID)}
	err := createGallery(sqb, studioGallery)
	if err != nil {
		return err
	}

	return nil
}

func makeGallery(name string, expectedResult bool) *models.Gallery {
	gallery := &models.Gallery{
		Checksum: md5.FromString(name),
		Path:     models.NullString(name),
	}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		gallery.Title = sql.NullString{Valid: true, String: name}
	}

	return gallery
}

func createGallery(sqb models.GalleryWriter, gallery *models.Gallery) error {
	_, err := sqb.Create(*gallery)

	if err != nil {
		return fmt.Errorf("Failed to create gallery with name '%s': %s", gallery.Path.String, err.Error())
	}

	return nil
}

func withTxn(f func(r models.Repository) error) error {
	t := sqlite.NewTransactionManager()
	return t.WithTxn(context.TODO(), f)
}

func populateDB() error {
	if err := withTxn(func(r models.Repository) error {
		err := createPerformer(r.Performer())
		if err != nil {
			return err
		}

		_, err = createStudio(r.Studio(), testName)
		if err != nil {
			return err
		}

		// create existing studio
		existingStudio, err := createStudio(r.Studio(), existingStudioName)
		if err != nil {
			return err
		}

		existingStudioID = existingStudio.ID

		err = createTag(r.Tag())
		if err != nil {
			return err
		}

		err = createScenes(r.Scene())
		if err != nil {
			return err
		}

		err = createImages(r.Image())
		if err != nil {
			return err
		}

		err = createGalleries(r.Gallery())
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
	if err := withTxn(func(r models.Repository) error {
		var err error
		performers, err = r.Performer().All()
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	for _, p := range performers {
		if err := withTxn(func(r models.Repository) error {
			return PerformerScenes(p, nil, r.Scene(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that scenes were tagged correctly
	withTxn(func(r models.Repository) error {
		pqb := r.Performer()

		scenes, err := r.Scene().All()
		if err != nil {
			t.Error(err.Error())
		}

		for _, scene := range scenes {
			performers, err := pqb.FindBySceneID(scene.ID)

			if err != nil {
				t.Errorf("Error getting scene performers: %s", err.Error())
			}

			// title is only set on scenes where we expect performer to be set
			if scene.Title.String == scene.Path && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, scene.Path)
			} else if scene.Title.String != scene.Path && len(performers) > 0 {
				t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, scene.Path)
			}
		}

		return nil
	})
}

func TestParseStudioScenes(t *testing.T) {
	var studios []*models.Studio
	if err := withTxn(func(r models.Repository) error {
		var err error
		studios, err = r.Studio().All()
		return err
	}); err != nil {
		t.Errorf("Error getting studio: %s", err)
		return
	}

	for _, s := range studios {
		if err := withTxn(func(r models.Repository) error {
			aliases, err := r.Studio().GetAliases(s.ID)
			if err != nil {
				return err
			}

			return StudioScenes(s, nil, aliases, r.Scene(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that scenes were tagged correctly
	withTxn(func(r models.Repository) error {
		scenes, err := r.Scene().All()
		if err != nil {
			t.Error(err.Error())
		}

		for _, scene := range scenes {
			// check for existing studio id scene first
			if scene.Path == existingStudioSceneName {
				if scene.StudioID.Int64 != int64(existingStudioID) {
					t.Error("Incorrectly overwrote studio ID for scene with existing studio ID")
				}
			} else {
				// title is only set on scenes where we expect studio to be set
				if scene.Title.String == scene.Path {
					if !scene.StudioID.Valid {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, scene.Path)
					} else if scene.StudioID.Int64 != int64(studios[1].ID) {
						t.Errorf("Incorrect studio id %d set for path '%s'", scene.StudioID.Int64, scene.Path)
					}

				} else if scene.Title.String != scene.Path && scene.StudioID.Int64 == int64(studios[1].ID) {
					t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, scene.Path)
				}
			}
		}

		return nil
	})
}

func TestParseTagScenes(t *testing.T) {
	var tags []*models.Tag
	if err := withTxn(func(r models.Repository) error {
		var err error
		tags, err = r.Tag().All()
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	for _, s := range tags {
		if err := withTxn(func(r models.Repository) error {
			aliases, err := r.Tag().GetAliases(s.ID)
			if err != nil {
				return err
			}

			return TagScenes(s, nil, aliases, r.Scene(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that scenes were tagged correctly
	withTxn(func(r models.Repository) error {
		scenes, err := r.Scene().All()
		if err != nil {
			t.Error(err.Error())
		}

		tqb := r.Tag()

		for _, scene := range scenes {
			tags, err := tqb.FindBySceneID(scene.ID)

			if err != nil {
				t.Errorf("Error getting scene tags: %s", err.Error())
			}

			// title is only set on scenes where we expect tag to be set
			if scene.Title.String == scene.Path && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, scene.Path)
			} else if scene.Title.String != scene.Path && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, scene.Path)
			}
		}

		return nil
	})
}

func TestParsePerformerImages(t *testing.T) {
	var performers []*models.Performer
	if err := withTxn(func(r models.Repository) error {
		var err error
		performers, err = r.Performer().All()
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	for _, p := range performers {
		if err := withTxn(func(r models.Repository) error {
			return PerformerImages(p, nil, r.Image(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that images were tagged correctly
	withTxn(func(r models.Repository) error {
		pqb := r.Performer()

		images, err := r.Image().All()
		if err != nil {
			t.Error(err.Error())
		}

		for _, image := range images {
			performers, err := pqb.FindByImageID(image.ID)

			if err != nil {
				t.Errorf("Error getting image performers: %s", err.Error())
			}

			// title is only set on images where we expect performer to be set
			if image.Title.String == image.Path && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, image.Path)
			} else if image.Title.String != image.Path && len(performers) > 0 {
				t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, image.Path)
			}
		}

		return nil
	})
}

func TestParseStudioImages(t *testing.T) {
	var studios []*models.Studio
	if err := withTxn(func(r models.Repository) error {
		var err error
		studios, err = r.Studio().All()
		return err
	}); err != nil {
		t.Errorf("Error getting studio: %s", err)
		return
	}

	for _, s := range studios {
		if err := withTxn(func(r models.Repository) error {
			aliases, err := r.Studio().GetAliases(s.ID)
			if err != nil {
				return err
			}

			return StudioImages(s, nil, aliases, r.Image(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that images were tagged correctly
	withTxn(func(r models.Repository) error {
		images, err := r.Image().All()
		if err != nil {
			t.Error(err.Error())
		}

		for _, image := range images {
			// check for existing studio id image first
			if image.Path == existingStudioImageName {
				if image.StudioID.Int64 != int64(existingStudioID) {
					t.Error("Incorrectly overwrote studio ID for image with existing studio ID")
				}
			} else {
				// title is only set on images where we expect studio to be set
				if image.Title.String == image.Path {
					if !image.StudioID.Valid {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, image.Path)
					} else if image.StudioID.Int64 != int64(studios[1].ID) {
						t.Errorf("Incorrect studio id %d set for path '%s'", image.StudioID.Int64, image.Path)
					}

				} else if image.Title.String != image.Path && image.StudioID.Int64 == int64(studios[1].ID) {
					t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, image.Path)
				}
			}
		}

		return nil
	})
}

func TestParseTagImages(t *testing.T) {
	var tags []*models.Tag
	if err := withTxn(func(r models.Repository) error {
		var err error
		tags, err = r.Tag().All()
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	for _, s := range tags {
		if err := withTxn(func(r models.Repository) error {
			aliases, err := r.Tag().GetAliases(s.ID)
			if err != nil {
				return err
			}

			return TagImages(s, nil, aliases, r.Image(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that images were tagged correctly
	withTxn(func(r models.Repository) error {
		images, err := r.Image().All()
		if err != nil {
			t.Error(err.Error())
		}

		tqb := r.Tag()

		for _, image := range images {
			tags, err := tqb.FindByImageID(image.ID)

			if err != nil {
				t.Errorf("Error getting image tags: %s", err.Error())
			}

			// title is only set on images where we expect performer to be set
			if image.Title.String == image.Path && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, image.Path)
			} else if image.Title.String != image.Path && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, image.Path)
			}
		}

		return nil
	})
}

func TestParsePerformerGalleries(t *testing.T) {
	var performers []*models.Performer
	if err := withTxn(func(r models.Repository) error {
		var err error
		performers, err = r.Performer().All()
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	for _, p := range performers {
		if err := withTxn(func(r models.Repository) error {
			return PerformerGalleries(p, nil, r.Gallery(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that galleries were tagged correctly
	withTxn(func(r models.Repository) error {
		pqb := r.Performer()

		galleries, err := r.Gallery().All()
		if err != nil {
			t.Error(err.Error())
		}

		for _, gallery := range galleries {
			performers, err := pqb.FindByGalleryID(gallery.ID)

			if err != nil {
				t.Errorf("Error getting gallery performers: %s", err.Error())
			}

			// title is only set on galleries where we expect performer to be set
			if gallery.Title.String == gallery.Path.String && len(performers) == 0 {
				t.Errorf("Did not set performer '%s' for path '%s'", testName, gallery.Path.String)
			} else if gallery.Title.String != gallery.Path.String && len(performers) > 0 {
				t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, gallery.Path.String)
			}
		}

		return nil
	})
}

func TestParseStudioGalleries(t *testing.T) {
	var studios []*models.Studio
	if err := withTxn(func(r models.Repository) error {
		var err error
		studios, err = r.Studio().All()
		return err
	}); err != nil {
		t.Errorf("Error getting studio: %s", err)
		return
	}

	for _, s := range studios {
		if err := withTxn(func(r models.Repository) error {
			aliases, err := r.Studio().GetAliases(s.ID)
			if err != nil {
				return err
			}

			return StudioGalleries(s, nil, aliases, r.Gallery(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that galleries were tagged correctly
	withTxn(func(r models.Repository) error {
		galleries, err := r.Gallery().All()
		if err != nil {
			t.Error(err.Error())
		}

		for _, gallery := range galleries {
			// check for existing studio id gallery first
			if gallery.Path.String == existingStudioGalleryName {
				if gallery.StudioID.Int64 != int64(existingStudioID) {
					t.Error("Incorrectly overwrote studio ID for gallery with existing studio ID")
				}
			} else {
				// title is only set on galleries where we expect studio to be set
				if gallery.Title.String == gallery.Path.String {
					if !gallery.StudioID.Valid {
						t.Errorf("Did not set studio '%s' for path '%s'", testName, gallery.Path.String)
					} else if gallery.StudioID.Int64 != int64(studios[1].ID) {
						t.Errorf("Incorrect studio id %d set for path '%s'", gallery.StudioID.Int64, gallery.Path.String)
					}

				} else if gallery.Title.String != gallery.Path.String && gallery.StudioID.Int64 == int64(studios[1].ID) {
					t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, gallery.Path.String)
				}
			}
		}

		return nil
	})
}

func TestParseTagGalleries(t *testing.T) {
	var tags []*models.Tag
	if err := withTxn(func(r models.Repository) error {
		var err error
		tags, err = r.Tag().All()
		return err
	}); err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	for _, s := range tags {
		if err := withTxn(func(r models.Repository) error {
			aliases, err := r.Tag().GetAliases(s.ID)
			if err != nil {
				return err
			}

			return TagGalleries(s, nil, aliases, r.Gallery(), nil)
		}); err != nil {
			t.Errorf("Error auto-tagging performers: %s", err)
		}
	}

	// verify that galleries were tagged correctly
	withTxn(func(r models.Repository) error {
		galleries, err := r.Gallery().All()
		if err != nil {
			t.Error(err.Error())
		}

		tqb := r.Tag()

		for _, gallery := range galleries {
			tags, err := tqb.FindByGalleryID(gallery.ID)

			if err != nil {
				t.Errorf("Error getting gallery tags: %s", err.Error())
			}

			// title is only set on galleries where we expect performer to be set
			if gallery.Title.String == gallery.Path.String && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, gallery.Path.String)
			} else if gallery.Title.String != gallery.Path.String && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, gallery.Path.String)
			}
		}

		return nil
	})
}
