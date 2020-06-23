// +build integration

package manager

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

const testName = "Foo's Bar"
const testExtension = ".mp4"
const existingStudioName = "ExistingStudio"

const existingStudioSceneName = testName + ".dontChangeStudio" + testExtension

var existingStudioID int

var testSeparators = []string{
	".",
	"-",
	"_",
	" ",
}

var testEndSeparators = []string{
	"{",
	"}",
	"(",
	")",
	",",
}

func generateNamePatterns(name, separator string) []string {
	var ret []string
	ret = append(ret, fmt.Sprintf("%s%saaa"+testExtension, name, separator))
	ret = append(ret, fmt.Sprintf("aaa%s%s"+testExtension, separator, name))
	ret = append(ret, fmt.Sprintf("aaa%s%s%sbbb"+testExtension, separator, name, separator))
	ret = append(ret, fmt.Sprintf("dir/%s%saaa"+testExtension, name, separator))
	ret = append(ret, fmt.Sprintf("dir\\%s%saaa"+testExtension, name, separator))
	ret = append(ret, fmt.Sprintf("%s%saaa/dir/bbb"+testExtension, name, separator))
	ret = append(ret, fmt.Sprintf("%s%saaa\\dir\\bbb"+testExtension, name, separator))
	ret = append(ret, fmt.Sprintf("dir/%s%s/aaa"+testExtension, name, separator))
	ret = append(ret, fmt.Sprintf("dir\\%s%s\\aaa"+testExtension, name, separator))

	return ret
}

func generateFalseNamePattern(name string, separator string) string {
	splitted := strings.Split(name, " ")

	return fmt.Sprintf("%s%saaa%s%s"+testExtension, splitted[0], separator, separator, splitted[1])
}

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

func createPerformer(tx *sqlx.Tx) error {
	// create the performer
	pqb := models.NewPerformerQueryBuilder()

	performer := models.Performer{
		Checksum: testName,
		Name:     sql.NullString{Valid: true, String: testName},
		Favorite: sql.NullBool{Valid: true, Bool: false},
	}

	_, err := pqb.Create(performer, tx)
	if err != nil {
		return err
	}

	return nil
}

func createStudio(tx *sqlx.Tx, name string) (*models.Studio, error) {
	// create the studio
	qb := models.NewStudioQueryBuilder()

	studio := models.Studio{
		Checksum: name,
		Name:     sql.NullString{Valid: true, String: testName},
	}

	return qb.Create(studio, tx)
}

func createTag(tx *sqlx.Tx) error {
	// create the studio
	qb := models.NewTagQueryBuilder()

	tag := models.Tag{
		Name: testName,
	}

	_, err := qb.Create(tag, tx)
	if err != nil {
		return err
	}

	return nil
}

func createScenes(tx *sqlx.Tx) error {
	sqb := models.NewSceneQueryBuilder()

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
		err := createScene(sqb, tx, makeScene(fn, true))
		if err != nil {
			return err
		}
	}
	for _, fn := range falseScenePatterns {
		err := createScene(sqb, tx, makeScene(fn, false))
		if err != nil {
			return err
		}
	}

	// create scene with existing studio io
	studioScene := makeScene(existingStudioSceneName, true)
	studioScene.StudioID = sql.NullInt64{Valid: true, Int64: int64(existingStudioID)}
	err := createScene(sqb, tx, studioScene)
	if err != nil {
		return err
	}

	return nil
}

func makeScene(name string, expectedResult bool) *models.Scene {
	scene := &models.Scene{
		Checksum: utils.MD5FromString(name),
		Path:     name,
	}

	// if expectedResult is true then we expect it to match, set the title accordingly
	if expectedResult {
		scene.Title = sql.NullString{Valid: true, String: name}
	}

	return scene
}

func createScene(sqb models.SceneQueryBuilder, tx *sqlx.Tx, scene *models.Scene) error {
	_, err := sqb.Create(*scene, tx)

	if err != nil {
		return fmt.Errorf("Failed to create scene with name '%s': %s", scene.Path, err.Error())
	}

	return nil
}

func populateDB() error {
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	err := createPerformer(tx)
	if err != nil {
		return err
	}

	_, err = createStudio(tx, testName)
	if err != nil {
		return err
	}

	// create existing studio
	existingStudio, err := createStudio(tx, existingStudioName)
	if err != nil {
		return err
	}

	existingStudioID = existingStudio.ID

	err = createTag(tx)
	if err != nil {
		return err
	}

	err = createScenes(tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func TestParsePerformers(t *testing.T) {
	pqb := models.NewPerformerQueryBuilder()
	performers, err := pqb.All()

	if err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	task := AutoTagPerformerTask{
		performer: performers[0],
	}

	var wg sync.WaitGroup
	wg.Add(1)
	task.Start(&wg)

	// verify that scenes were tagged correctly
	sqb := models.NewSceneQueryBuilder()

	scenes, err := sqb.All()

	for _, scene := range scenes {
		performers, err := pqb.FindBySceneID(scene.ID, nil)

		if err != nil {
			t.Errorf("Error getting scene performers: %s", err.Error())
			return
		}

		// title is only set on scenes where we expect performer to be set
		if scene.Title.String == scene.Path && len(performers) == 0 {
			t.Errorf("Did not set performer '%s' for path '%s'", testName, scene.Path)
		} else if scene.Title.String != scene.Path && len(performers) > 0 {
			t.Errorf("Incorrectly set performer '%s' for path '%s'", testName, scene.Path)
		}
	}
}

func TestParseStudios(t *testing.T) {
	studioQuery := models.NewStudioQueryBuilder()
	studios, err := studioQuery.All()

	if err != nil {
		t.Errorf("Error getting studio: %s", err)
		return
	}

	task := AutoTagStudioTask{
		studio: studios[0],
	}

	var wg sync.WaitGroup
	wg.Add(1)
	task.Start(&wg)

	// verify that scenes were tagged correctly
	sqb := models.NewSceneQueryBuilder()

	scenes, err := sqb.All()

	for _, scene := range scenes {
		// check for existing studio id scene first
		if scene.Path == existingStudioSceneName {
			if scene.StudioID.Int64 != int64(existingStudioID) {
				t.Error("Incorrectly overwrote studio ID for scene with existing studio ID")
			}
		} else {
			// title is only set on scenes where we expect studio to be set
			if scene.Title.String == scene.Path && scene.StudioID.Int64 != int64(studios[0].ID) {
				t.Errorf("Did not set studio '%s' for path '%s'", testName, scene.Path)
			} else if scene.Title.String != scene.Path && scene.StudioID.Int64 == int64(studios[0].ID) {
				t.Errorf("Incorrectly set studio '%s' for path '%s'", testName, scene.Path)
			}
		}
	}
}

func TestParseTags(t *testing.T) {
	tagQuery := models.NewTagQueryBuilder()
	tags, err := tagQuery.All()

	if err != nil {
		t.Errorf("Error getting performer: %s", err)
		return
	}

	task := AutoTagTagTask{
		tag: tags[0],
	}

	var wg sync.WaitGroup
	wg.Add(1)
	task.Start(&wg)

	// verify that scenes were tagged correctly
	sqb := models.NewSceneQueryBuilder()

	scenes, err := sqb.All()

	for _, scene := range scenes {
		tags, err := tagQuery.FindBySceneID(scene.ID, nil)

		if err != nil {
			t.Errorf("Error getting scene tags: %s", err.Error())
			return
		}

		// title is only set on scenes where we expect performer to be set
		if scene.Title.String == scene.Path && len(tags) == 0 {
			t.Errorf("Did not set tag '%s' for path '%s'", testName, scene.Path)
		} else if scene.Title.String != scene.Path && len(tags) > 0 {
			t.Errorf("Incorrectly set tag '%s' for path '%s'", testName, scene.Path)
		}
	}
}
