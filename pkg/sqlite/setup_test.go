// +build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

const totalScenes = 12
const totalImages = 6 // TODO - add one for zip file
const performersNameCase = 9
const performersNameNoCase = 2
const moviesNameCase = 2
const moviesNameNoCase = 1
const totalGalleries = 8
const tagsNameNoCase = 2
const tagsNameCase = 12
const studiosNameCase = 6
const studiosNameNoCase = 1

var sceneIDs []int
var imageIDs []int
var performerIDs []int
var movieIDs []int
var galleryIDs []int
var tagIDs []int
var studioIDs []int
var markerIDs []int

var tagNames []string
var studioNames []string
var movieNames []string
var performerNames []string

const sceneIdxWithMovie = 0
const sceneIdxWithGallery = 1
const sceneIdxWithPerformer = 2
const sceneIdxWithTwoPerformers = 3
const sceneIdxWithTag = 4
const sceneIdxWithTwoTags = 5
const sceneIdxWithStudio = 6
const sceneIdxWithMarker = 7

const imageIdxWithGallery = 0
const imageIdxWithPerformer = 1
const imageIdxWithTwoPerformers = 2
const imageIdxWithTag = 3
const imageIdxWithTwoTags = 4
const imageIdxWithStudio = 5
const imageIdxInZip = 6

const performerIdxWithScene = 0
const performerIdx1WithScene = 1
const performerIdx2WithScene = 2
const performerIdxWithImage = 3
const performerIdx1WithImage = 4
const performerIdx2WithImage = 5
const performerIdxWithGallery = 6
const performerIdx1WithGallery = 7
const performerIdx2WithGallery = 8

// performers with dup names start from the end
const performerIdx1WithDupName = 9
const performerIdxWithDupName = 10

const movieIdxWithScene = 0
const movieIdxWithStudio = 1

// movies with dup names start from the end
const movieIdxWithDupName = 2

const galleryIdxWithScene = 0
const galleryIdxWithImage = 1
const galleryIdxWithPerformer = 2
const galleryIdxWithTwoPerformers = 3
const galleryIdxWithTag = 4
const galleryIdxWithTwoTags = 5
const galleryIdxWithStudio = 6

const tagIdxWithScene = 0
const tagIdx1WithScene = 1
const tagIdx2WithScene = 2
const tagIdxWithPrimaryMarker = 3
const tagIdxWithMarker = 4
const tagIdxWithCoverImage = 5
const tagIdxWithImage = 6
const tagIdx1WithImage = 7
const tagIdx2WithImage = 8
const tagIdxWithGallery = 9
const tagIdx1WithGallery = 10
const tagIdx2WithGallery = 11

// tags with dup names start from the end
const tagIdx1WithDupName = 12
const tagIdxWithDupName = 13

const studioIdxWithScene = 0
const studioIdxWithMovie = 1
const studioIdxWithChildStudio = 2
const studioIdxWithParentStudio = 3
const studioIdxWithImage = 4
const studioIdxWithGallery = 5

// studios with dup names start from the end
const studioIdxWithDupName = 6

const markerIdxWithScene = 0

const pathField = "Path"
const checksumField = "Checksum"
const titleField = "Title"
const zipPath = "zipPath.zip"

func TestMain(m *testing.M) {
	ret := runTests(m)
	os.Exit(ret)
}

func withTxn(f func(r models.Repository) error) error {
	t := sqlite.NewTransactionManager()
	return t.WithTxn(context.TODO(), f)
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

func populateDB() error {
	if err := withTxn(func(r models.Repository) error {
		if err := createScenes(r.Scene(), totalScenes); err != nil {
			return fmt.Errorf("error creating scenes: %s", err.Error())
		}

		if err := createImages(r.Image(), totalImages); err != nil {
			return fmt.Errorf("error creating images: %s", err.Error())
		}

		if err := createGalleries(r.Gallery(), totalGalleries); err != nil {
			return fmt.Errorf("error creating galleries: %s", err.Error())
		}

		if err := createMovies(r.Movie(), moviesNameCase, moviesNameNoCase); err != nil {
			return fmt.Errorf("error creating movies: %s", err.Error())
		}

		if err := createPerformers(r.Performer(), performersNameCase, performersNameNoCase); err != nil {
			return fmt.Errorf("error creating performers: %s", err.Error())
		}

		if err := createTags(r.Tag(), tagsNameCase, tagsNameNoCase); err != nil {
			return fmt.Errorf("error creating tags: %s", err.Error())
		}

		if err := addTagImage(r.Tag(), tagIdxWithCoverImage); err != nil {
			return fmt.Errorf("error adding tag image: %s", err.Error())
		}

		if err := createStudios(r.Studio(), studiosNameCase, studiosNameNoCase); err != nil {
			return fmt.Errorf("error creating studios: %s", err.Error())
		}

		if err := linkSceneGallery(r.Scene(), sceneIdxWithGallery, galleryIdxWithScene); err != nil {
			return fmt.Errorf("error linking scene to gallery: %s", err.Error())
		}

		if err := linkSceneMovie(r.Scene(), sceneIdxWithMovie, movieIdxWithScene); err != nil {
			return fmt.Errorf("error scene to movie: %s", err.Error())
		}

		if err := linkScenePerformers(r.Scene()); err != nil {
			return fmt.Errorf("error linking scene performers: %s", err.Error())
		}

		if err := linkSceneTags(r.Scene()); err != nil {
			return fmt.Errorf("error linking scene tags: %s", err.Error())
		}

		if err := linkSceneStudio(r.Scene(), sceneIdxWithStudio, studioIdxWithScene); err != nil {
			return fmt.Errorf("error linking scene studio: %s", err.Error())
		}

		if err := linkImageGallery(r.Gallery(), imageIdxWithGallery, galleryIdxWithImage); err != nil {
			return fmt.Errorf("error linking gallery images: %s", err.Error())
		}

		if err := linkImagePerformers(r.Image()); err != nil {
			return fmt.Errorf("error linking image performers: %s", err.Error())
		}

		if err := linkImageTags(r.Image()); err != nil {
			return fmt.Errorf("error linking image tags: %s", err.Error())
		}

		if err := linkImageStudio(r.Image(), imageIdxWithStudio, studioIdxWithImage); err != nil {
			return fmt.Errorf("error linking image studio: %s", err.Error())
		}

		if err := linkMovieStudio(r.Movie(), movieIdxWithStudio, studioIdxWithMovie); err != nil {
			return fmt.Errorf("error linking movie studio: %s", err.Error())
		}

		if err := linkStudioParent(r.Studio(), studioIdxWithChildStudio, studioIdxWithParentStudio); err != nil {
			return fmt.Errorf("error linking studio parent: %s", err.Error())
		}

		if err := linkGalleryPerformers(r.Gallery()); err != nil {
			return fmt.Errorf("error linking gallery performers: %s", err.Error())
		}

		if err := linkGalleryTags(r.Gallery()); err != nil {
			return fmt.Errorf("error linking gallery tags: %s", err.Error())
		}

		if err := linkGalleryStudio(r.Gallery(), galleryIdxWithStudio, studioIdxWithGallery); err != nil {
			return fmt.Errorf("error linking gallery studio: %s", err.Error())
		}

		if err := createMarker(r.SceneMarker(), sceneIdxWithMarker, tagIdxWithPrimaryMarker, []int{tagIdxWithMarker}); err != nil {
			return fmt.Errorf("error creating scene marker: %s", err.Error())
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func getSceneStringValue(index int, field string) string {
	return fmt.Sprintf("scene_%04d_%s", index, field)
}

func getRating(index int) sql.NullInt64 {
	rating := index % 6
	return sql.NullInt64{Int64: int64(rating), Valid: rating > 0}
}

func getOCounter(index int) int {
	return index % 3
}

func getSceneDuration(index int) sql.NullFloat64 {
	duration := index % 4
	duration = duration * 100

	return sql.NullFloat64{
		Float64: float64(duration) + 0.432,
		Valid:   duration != 0,
	}
}

func getHeight(index int) sql.NullInt64 {
	heights := []int64{0, 200, 240, 300, 480, 700, 720, 800, 1080, 1500, 2160, 3000}
	height := heights[index%len(heights)]
	return sql.NullInt64{
		Int64: height,
		Valid: height != 0,
	}
}

func getSceneDate(index int) models.SQLiteDate {
	dates := []string{"null", "", "0001-01-01", "2001-02-03"}
	date := dates[index%len(dates)]
	return models.SQLiteDate{
		String: date,
		Valid:  date != "null",
	}
}

func createScenes(sqb models.SceneReaderWriter, n int) error {
	for i := 0; i < n; i++ {
		scene := models.Scene{
			Path:     getSceneStringValue(i, pathField),
			Title:    sql.NullString{String: getSceneStringValue(i, titleField), Valid: true},
			Checksum: sql.NullString{String: getSceneStringValue(i, checksumField), Valid: true},
			Details:  sql.NullString{String: getSceneStringValue(i, "Details"), Valid: true},
			Rating:   getRating(i),
			OCounter: getOCounter(i),
			Duration: getSceneDuration(i),
			Height:   getHeight(i),
			Date:     getSceneDate(i),
		}

		created, err := sqb.Create(scene)

		if err != nil {
			return fmt.Errorf("Error creating scene %v+: %s", scene, err.Error())
		}

		sceneIDs = append(sceneIDs, created.ID)
	}

	return nil
}

func getImageStringValue(index int, field string) string {
	return fmt.Sprintf("image_%04d_%s", index, field)
}

func getImagePath(index int) string {
	// TODO - currently not working
	// if index == imageIdxInZip {
	// 	return image.ZipFilename(zipPath, "image_0001_Path")
	// }

	return getImageStringValue(index, pathField)
}

func createImages(qb models.ImageReaderWriter, n int) error {
	for i := 0; i < n; i++ {
		image := models.Image{
			Path:     getImagePath(i),
			Title:    sql.NullString{String: getImageStringValue(i, titleField), Valid: true},
			Checksum: getImageStringValue(i, checksumField),
			Rating:   getRating(i),
			OCounter: getOCounter(i),
			Height:   getHeight(i),
		}

		created, err := qb.Create(image)

		if err != nil {
			return fmt.Errorf("Error creating image %v+: %s", image, err.Error())
		}

		imageIDs = append(imageIDs, created.ID)
	}

	return nil
}

func getGalleryStringValue(index int, field string) string {
	return "gallery_" + strconv.FormatInt(int64(index), 10) + "_" + field
}

func createGalleries(gqb models.GalleryReaderWriter, n int) error {
	for i := 0; i < n; i++ {
		gallery := models.Gallery{
			Path:     models.NullString(getGalleryStringValue(i, pathField)),
			Checksum: getGalleryStringValue(i, checksumField),
		}

		created, err := gqb.Create(gallery)

		if err != nil {
			return fmt.Errorf("Error creating gallery %v+: %s", gallery, err.Error())
		}

		galleryIDs = append(galleryIDs, created.ID)
	}

	return nil
}

func getMovieStringValue(index int, field string) string {
	return "movie_" + strconv.FormatInt(int64(index), 10) + "_" + field
}

//createMoviees creates n movies with plain Name and o movies with camel cased NaMe included
func createMovies(mqb models.MovieReaderWriter, n int, o int) error {
	const namePlain = "Name"
	const nameNoCase = "NaMe"

	for i := 0; i < n+o; i++ {
		index := i
		name := namePlain

		if i >= n { // i<n tags get normal names
			name = nameNoCase       // i>=n movies get dup names if case is not checked
			index = n + o - (i + 1) // for the name to be the same the number (index) must be the same also
		} // so count backwards to 0 as needed
		// movies [ i ] and [ n + o - i - 1  ] should have similar names with only the Name!=NaMe part different

		name = getMovieStringValue(index, name)
		movie := models.Movie{
			Name:     sql.NullString{String: name, Valid: true},
			Checksum: utils.MD5FromString(name),
		}

		created, err := mqb.Create(movie)

		if err != nil {
			return fmt.Errorf("Error creating movie [%d] %v+: %s", i, movie, err.Error())
		}

		movieIDs = append(movieIDs, created.ID)
		movieNames = append(movieNames, created.Name.String)
	}

	return nil
}

func getPerformerStringValue(index int, field string) string {
	return "performer_" + strconv.FormatInt(int64(index), 10) + "_" + field
}

func getPerformerBoolValue(index int) bool {
	index = index % 2
	return index == 1
}

func getPerformerBirthdate(index int) string {
	const minAge = 18
	birthdate := time.Now()
	birthdate = birthdate.AddDate(-minAge-index, -1, -1)
	return birthdate.Format("2006-01-02")
}

//createPerformers creates n performers with plain Name and o performers with camel cased NaMe included
func createPerformers(pqb models.PerformerReaderWriter, n int, o int) error {
	const namePlain = "Name"
	const nameNoCase = "NaMe"

	name := namePlain

	for i := 0; i < n+o; i++ {
		index := i

		if i >= n { // i<n tags get normal names
			name = nameNoCase       // i>=n performers get dup names if case is not checked
			index = n + o - (i + 1) // for the name to be the same the number (index) must be the same also
		} // so count backwards to 0 as needed
		// performers [ i ] and [ n + o - i - 1  ] should have similar names with only the Name!=NaMe part different

		performer := models.Performer{
			Name:     sql.NullString{String: getPerformerStringValue(index, name), Valid: true},
			Checksum: getPerformerStringValue(i, checksumField),
			Favorite: sql.NullBool{Bool: getPerformerBoolValue(i), Valid: true},
			Birthdate: models.SQLiteDate{
				String: getPerformerBirthdate(i),
				Valid:  true,
			},
		}

		created, err := pqb.Create(performer)

		if err != nil {
			return fmt.Errorf("Error creating performer %v+: %s", performer, err.Error())
		}

		performerIDs = append(performerIDs, created.ID)
		performerNames = append(performerNames, created.Name.String)
	}

	return nil
}

func getTagStringValue(index int, field string) string {
	return "tag_" + strconv.FormatInt(int64(index), 10) + "_" + field
}

func getTagSceneCount(id int) int {
	if id == tagIDs[tagIdx1WithScene] || id == tagIDs[tagIdx2WithScene] || id == tagIDs[tagIdxWithScene] {
		return 1
	}

	return 0
}

func getTagMarkerCount(id int) int {
	if id == tagIDs[tagIdxWithMarker] || id == tagIDs[tagIdxWithPrimaryMarker] {
		return 1
	}

	return 0
}

//createTags creates n tags with plain Name and o tags with camel cased NaMe included
func createTags(tqb models.TagReaderWriter, n int, o int) error {
	const namePlain = "Name"
	const nameNoCase = "NaMe"

	name := namePlain

	for i := 0; i < n+o; i++ {
		index := i

		if i >= n { // i<n tags get normal names
			name = nameNoCase       // i>=n tags get dup names if case is not checked
			index = n + o - (i + 1) // for the name to be the same the number (index) must be the same also
		} // so count backwards to 0 as needed
		// tags [ i ] and [ n + o - i - 1  ] should have similar names with only the Name!=NaMe part different

		tag := models.Tag{
			Name: getTagStringValue(index, name),
		}

		created, err := tqb.Create(tag)

		if err != nil {
			return fmt.Errorf("Error creating tag %v+: %s", tag, err.Error())
		}

		tagIDs = append(tagIDs, created.ID)
		tagNames = append(tagNames, created.Name)
	}

	return nil
}

func getStudioStringValue(index int, field string) string {
	return "studio_" + strconv.FormatInt(int64(index), 10) + "_" + field
}

func createStudio(sqb models.StudioReaderWriter, name string, parentID *int64) (*models.Studio, error) {
	studio := models.Studio{
		Name:     sql.NullString{String: name, Valid: true},
		Checksum: utils.MD5FromString(name),
	}

	if parentID != nil {
		studio.ParentID = sql.NullInt64{Int64: *parentID, Valid: true}
	}

	created, err := sqb.Create(studio)

	if err != nil {
		return nil, fmt.Errorf("Error creating studio %v+: %s", studio, err.Error())
	}

	return created, nil
}

//createStudios creates n studios with plain Name and o studios with camel cased NaMe included
func createStudios(sqb models.StudioReaderWriter, n int, o int) error {
	const namePlain = "Name"
	const nameNoCase = "NaMe"

	for i := 0; i < n+o; i++ {
		index := i
		name := namePlain

		if i >= n { // i<n studios get normal names
			name = nameNoCase       // i>=n studios get dup names if case is not checked
			index = n + o - (i + 1) // for the name to be the same the number (index) must be the same also
		} // so count backwards to 0 as needed
		// studios [ i ] and [ n + o - i - 1  ] should have similar names with only the Name!=NaMe part different

		name = getStudioStringValue(index, name)
		created, err := createStudio(sqb, name, nil)

		if err != nil {
			return err
		}

		studioIDs = append(studioIDs, created.ID)
		studioNames = append(studioNames, created.Name.String)
	}

	return nil
}

func createMarker(mqb models.SceneMarkerReaderWriter, sceneIdx, primaryTagIdx int, tagIdxs []int) error {
	marker := models.SceneMarker{
		SceneID:      sql.NullInt64{Int64: int64(sceneIDs[sceneIdx]), Valid: true},
		PrimaryTagID: tagIDs[primaryTagIdx],
	}

	created, err := mqb.Create(marker)

	if err != nil {
		return fmt.Errorf("Error creating marker %v+: %s", marker, err.Error())
	}

	markerIDs = append(markerIDs, created.ID)

	newTagIDs := []int{}

	for _, tagIdx := range tagIdxs {
		newTagIDs = append(newTagIDs, tagIDs[tagIdx])
	}

	if err := mqb.UpdateTags(created.ID, newTagIDs); err != nil {
		return fmt.Errorf("Error creating marker/tag join: %s", err.Error())
	}

	return nil
}

func linkSceneMovie(qb models.SceneReaderWriter, sceneIndex, movieIndex int) error {
	sceneID := sceneIDs[sceneIndex]
	movies, err := qb.GetMovies(sceneID)
	if err != nil {
		return err
	}

	movies = append(movies, models.MoviesScenes{
		MovieID: movieIDs[movieIndex],
		SceneID: sceneID,
	})
	return qb.UpdateMovies(sceneID, movies)
}

func linkScenePerformers(qb models.SceneReaderWriter) error {
	if err := linkScenePerformer(qb, sceneIdxWithPerformer, performerIdxWithScene); err != nil {
		return err
	}
	if err := linkScenePerformer(qb, sceneIdxWithTwoPerformers, performerIdx1WithScene); err != nil {
		return err
	}
	if err := linkScenePerformer(qb, sceneIdxWithTwoPerformers, performerIdx2WithScene); err != nil {
		return err
	}

	return nil
}

func linkScenePerformer(qb models.SceneReaderWriter, sceneIndex, performerIndex int) error {
	_, err := scene.AddPerformer(qb, sceneIDs[sceneIndex], performerIDs[performerIndex])
	return err
}

func linkSceneGallery(qb models.SceneReaderWriter, sceneIndex, galleryIndex int) error {
	_, err := scene.AddGallery(qb, sceneIDs[sceneIndex], galleryIDs[galleryIndex])
	return err
}

func linkSceneTags(qb models.SceneReaderWriter) error {
	if err := linkSceneTag(qb, sceneIdxWithTag, tagIdxWithScene); err != nil {
		return err
	}
	if err := linkSceneTag(qb, sceneIdxWithTwoTags, tagIdx1WithScene); err != nil {
		return err
	}
	if err := linkSceneTag(qb, sceneIdxWithTwoTags, tagIdx2WithScene); err != nil {
		return err
	}

	return nil
}

func linkSceneTag(qb models.SceneReaderWriter, sceneIndex, tagIndex int) error {
	_, err := scene.AddTag(qb, sceneIDs[sceneIndex], tagIDs[tagIndex])
	return err
}

func linkSceneStudio(sqb models.SceneWriter, sceneIndex, studioIndex int) error {
	scene := models.ScenePartial{
		ID:       sceneIDs[sceneIndex],
		StudioID: &sql.NullInt64{Int64: int64(studioIDs[studioIndex]), Valid: true},
	}
	_, err := sqb.Update(scene)

	return err
}

func linkImageGallery(gqb models.GalleryReaderWriter, imageIndex, galleryIndex int) error {
	return gallery.AddImage(gqb, galleryIDs[galleryIndex], imageIDs[imageIndex])
}

func linkImageTags(iqb models.ImageReaderWriter) error {
	if err := linkImageTag(iqb, imageIdxWithTag, tagIdxWithImage); err != nil {
		return err
	}
	if err := linkImageTag(iqb, imageIdxWithTwoTags, tagIdx1WithImage); err != nil {
		return err
	}
	if err := linkImageTag(iqb, imageIdxWithTwoTags, tagIdx2WithImage); err != nil {
		return err
	}

	return nil
}

func linkImageTag(iqb models.ImageReaderWriter, imageIndex, tagIndex int) error {
	imageID := imageIDs[imageIndex]
	tags, err := iqb.GetTagIDs(imageID)
	if err != nil {
		return err
	}

	tags = append(tags, tagIDs[tagIndex])

	return iqb.UpdateTags(imageID, tags)
}

func linkImageStudio(qb models.ImageWriter, imageIndex, studioIndex int) error {
	image := models.ImagePartial{
		ID:       imageIDs[imageIndex],
		StudioID: &sql.NullInt64{Int64: int64(studioIDs[studioIndex]), Valid: true},
	}
	_, err := qb.Update(image)

	return err
}

func linkImagePerformers(qb models.ImageReaderWriter) error {
	if err := linkImagePerformer(qb, imageIdxWithPerformer, performerIdxWithImage); err != nil {
		return err
	}
	if err := linkImagePerformer(qb, imageIdxWithTwoPerformers, performerIdx1WithImage); err != nil {
		return err
	}
	if err := linkImagePerformer(qb, imageIdxWithTwoPerformers, performerIdx2WithImage); err != nil {
		return err
	}

	return nil
}

func linkImagePerformer(iqb models.ImageReaderWriter, imageIndex, performerIndex int) error {
	imageID := imageIDs[imageIndex]
	performers, err := iqb.GetPerformerIDs(imageID)
	if err != nil {
		return err
	}

	performers = append(performers, performerIDs[performerIndex])

	return iqb.UpdatePerformers(imageID, performers)
}

func linkMovieStudio(mqb models.MovieWriter, movieIndex, studioIndex int) error {
	movie := models.MoviePartial{
		ID:       movieIDs[movieIndex],
		StudioID: &sql.NullInt64{Int64: int64(studioIDs[studioIndex]), Valid: true},
	}
	_, err := mqb.Update(movie)

	return err
}

func linkStudioParent(qb models.StudioWriter, parentIndex, childIndex int) error {
	studio := models.StudioPartial{
		ID:       studioIDs[childIndex],
		ParentID: &sql.NullInt64{Int64: int64(studioIDs[parentIndex]), Valid: true},
	}
	_, err := qb.Update(studio)

	return err
}

func addTagImage(qb models.TagWriter, tagIndex int) error {
	return qb.UpdateImage(tagIDs[tagIndex], models.DefaultTagImage)
}

func linkGalleryTags(iqb models.GalleryReaderWriter) error {
	if err := linkGalleryTag(iqb, galleryIdxWithTag, tagIdxWithGallery); err != nil {
		return err
	}
	if err := linkGalleryTag(iqb, galleryIdxWithTwoTags, tagIdx1WithGallery); err != nil {
		return err
	}
	if err := linkGalleryTag(iqb, galleryIdxWithTwoTags, tagIdx2WithGallery); err != nil {
		return err
	}

	return nil
}

func linkGalleryTag(iqb models.GalleryReaderWriter, galleryIndex, tagIndex int) error {
	galleryID := galleryIDs[galleryIndex]
	tags, err := iqb.GetTagIDs(galleryID)
	if err != nil {
		return err
	}

	tags = append(tags, tagIDs[tagIndex])

	return iqb.UpdateTags(galleryID, tags)
}

func linkGalleryStudio(qb models.GalleryWriter, galleryIndex, studioIndex int) error {
	gallery := models.GalleryPartial{
		ID:       galleryIDs[galleryIndex],
		StudioID: &sql.NullInt64{Int64: int64(studioIDs[studioIndex]), Valid: true},
	}
	_, err := qb.UpdatePartial(gallery)

	return err
}

func linkGalleryPerformers(qb models.GalleryReaderWriter) error {
	if err := linkGalleryPerformer(qb, galleryIdxWithPerformer, performerIdxWithGallery); err != nil {
		return err
	}
	if err := linkGalleryPerformer(qb, galleryIdxWithTwoPerformers, performerIdx1WithGallery); err != nil {
		return err
	}
	if err := linkGalleryPerformer(qb, galleryIdxWithTwoPerformers, performerIdx2WithGallery); err != nil {
		return err
	}

	return nil
}

func linkGalleryPerformer(iqb models.GalleryReaderWriter, galleryIndex, performerIndex int) error {
	galleryID := galleryIDs[galleryIndex]
	performers, err := iqb.GetPerformerIDs(galleryID)
	if err != nil {
		return err
	}

	performers = append(performers, performerIDs[performerIndex])

	return iqb.UpdatePerformers(galleryID, performers)
}
