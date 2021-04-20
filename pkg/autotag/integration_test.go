// +build integration

package autotag

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const testName = "Foo's Bar"
const existingStudioName = "ExistingStudio"

const existingStudioSceneName = testName + ".dontChangeStudio.mp4"

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
	f, err := ioutil.TempFile("", "*.sqlite")
	if err != nil {
		panic(fmt.Sprintf("Could not create temporary file: %s", err.Error()))
	}

	f.Close()
	databaseFile := f.Name()
	database.Initialize(databaseFile)

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
	var scenePatterns []string
	var falseScenePatterns []string

	separators := append(testSeparators, testEndSeparators...)

	for _, separator := range separators {
		scenePatterns = append(scenePatterns, generateNamePatterns(testName, separator)...)
		scenePatterns = append(scenePatterns, generateNamePatterns(strings.ToLower(testName), separator)...)
		falseScenePatterns = append(falseScenePatterns, generateFalseNamePattern(testName, separator))
	}

	// add test cases for intra-name separators
	for _, separator := range testSeparators {
		if separator != " " {
			scenePatterns = append(scenePatterns, generateNamePatterns(strings.Replace(testName, " ", separator, -1), separator)...)
		}
	}

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
		Checksum: sql.NullString{String: utils.MD5FromString(name), Valid: true},
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

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func TestParsePerformers(t *testing.T) {
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
			return PerformerScenes(p, nil, r.Scene())
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

func TestParseStudios(t *testing.T) {
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
			return StudioScenes(s, nil, r.Scene())
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

func TestParseTags(t *testing.T) {
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
			return TagScenes(s, nil, r.Scene())
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

			// title is only set on scenes where we expect performer to be set
			if scene.Title.String == scene.Path && len(tags) == 0 {
				t.Errorf("Did not set tag '%s' for path '%s'", testName, scene.Path)
			} else if scene.Title.String != scene.Path && len(tags) > 0 {
				t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, scene.Path)
			}
		}

		return nil
	})
}
