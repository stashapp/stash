//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"
)

const (
	spacedSceneTitle = "zzz yyy xxx"
)

const (
	sceneIdxWithMovie = iota
	sceneIdxWithGallery
	sceneIdxWithPerformer
	sceneIdx1WithPerformer
	sceneIdx2WithPerformer
	sceneIdxWithTwoPerformers
	sceneIdxWithTag
	sceneIdxWithTwoTags
	sceneIdxWithMarkerAndTag
	sceneIdxWithStudio
	sceneIdx1WithStudio
	sceneIdx2WithStudio
	sceneIdxWithMarkers
	sceneIdxWithPerformerTag
	sceneIdxWithPerformerTwoTags
	sceneIdxWithSpacedName
	sceneIdxWithStudioPerformer
	sceneIdxWithGrandChildStudio
	// new indexes above
	lastSceneIdx

	totalScenes = lastSceneIdx + 3
)

const (
	imageIdxWithGallery = iota
	imageIdx1WithGallery
	imageIdx2WithGallery
	imageIdxWithTwoGalleries
	imageIdxWithPerformer
	imageIdx1WithPerformer
	imageIdx2WithPerformer
	imageIdxWithTwoPerformers
	imageIdxWithTag
	imageIdxWithTwoTags
	imageIdxWithStudio
	imageIdx1WithStudio
	imageIdx2WithStudio
	imageIdxWithStudioPerformer
	imageIdxInZip // TODO - not implemented
	imageIdxWithPerformerTag
	imageIdxWithPerformerTwoTags
	imageIdxWithGrandChildStudio
	// new indexes above
	totalImages
)

const (
	performerIdxWithScene = iota
	performerIdx1WithScene
	performerIdx2WithScene
	performerIdxWithTwoScenes
	performerIdxWithImage
	performerIdxWithTwoImages
	performerIdx1WithImage
	performerIdx2WithImage
	performerIdxWithTag
	performerIdxWithTwoTags
	performerIdxWithGallery
	performerIdxWithTwoGalleries
	performerIdx1WithGallery
	performerIdx2WithGallery
	performerIdxWithSceneStudio
	performerIdxWithImageStudio
	performerIdxWithGalleryStudio
	// new indexes above
	// performers with dup names start from the end
	performerIdx1WithDupName
	performerIdxWithDupName

	performersNameCase   = performerIdx1WithDupName
	performersNameNoCase = 2

	totalPerformers = performersNameCase + performersNameNoCase
)

const (
	movieIdxWithScene = iota
	movieIdxWithStudio
	// movies with dup names start from the end
	// create 10 more basic movies (can remove this if we add more indexes)
	movieIdxWithDupName = movieIdxWithStudio + 10

	moviesNameCase   = movieIdxWithDupName
	moviesNameNoCase = 1
)

const (
	galleryIdxWithScene = iota
	galleryIdxWithImage
	galleryIdx1WithImage
	galleryIdx2WithImage
	galleryIdxWithTwoImages
	galleryIdxWithPerformer
	galleryIdx1WithPerformer
	galleryIdx2WithPerformer
	galleryIdxWithTwoPerformers
	galleryIdxWithTag
	galleryIdxWithTwoTags
	galleryIdxWithStudio
	galleryIdx1WithStudio
	galleryIdx2WithStudio
	galleryIdxWithPerformerTag
	galleryIdxWithPerformerTwoTags
	galleryIdxWithStudioPerformer
	galleryIdxWithGrandChildStudio
	// new indexes above
	lastGalleryIdx

	totalGalleries = lastGalleryIdx + 1
)

const (
	tagIdxWithScene = iota
	tagIdx1WithScene
	tagIdx2WithScene
	tagIdx3WithScene
	tagIdxWithPrimaryMarkers
	tagIdxWithMarkers
	tagIdxWithCoverImage
	tagIdxWithImage
	tagIdx1WithImage
	tagIdx2WithImage
	tagIdxWithPerformer
	tagIdx1WithPerformer
	tagIdx2WithPerformer
	tagIdxWithGallery
	tagIdx1WithGallery
	tagIdx2WithGallery
	tagIdxWithChildTag
	tagIdxWithParentTag
	tagIdxWithGrandChild
	tagIdxWithParentAndChild
	tagIdxWithGrandParent
	// new indexes above
	// tags with dup names start from the end
	tagIdx1WithDupName
	tagIdxWithDupName

	tagsNameNoCase = 2
	tagsNameCase   = tagIdx1WithDupName

	totalTags = tagsNameCase + tagsNameNoCase
)

const (
	studioIdxWithScene = iota
	studioIdxWithTwoScenes
	studioIdxWithMovie
	studioIdxWithChildStudio
	studioIdxWithParentStudio
	studioIdxWithImage
	studioIdxWithTwoImages
	studioIdxWithGallery
	studioIdxWithTwoGalleries
	studioIdxWithScenePerformer
	studioIdxWithImagePerformer
	studioIdxWithGalleryPerformer
	studioIdxWithGrandChild
	studioIdxWithParentAndChild
	studioIdxWithGrandParent
	// new indexes above
	// studios with dup names start from the end
	studioIdxWithDupName

	studiosNameCase   = studioIdxWithDupName
	studiosNameNoCase = 1

	totalStudios = studiosNameCase + studiosNameNoCase
)

const (
	markerIdxWithScene = iota
	markerIdxWithTag
	markerIdxWithSceneTag
	totalMarkers
)

const (
	savedFilterIdxDefaultScene = iota
	savedFilterIdxDefaultImage
	savedFilterIdxScene
	savedFilterIdxImage

	// new indexes above
	totalSavedFilters
)

const (
	pathField            = "Path"
	checksumField        = "Checksum"
	titleField           = "Title"
	urlField             = "URL"
	zipPath              = "zipPath.zip"
	firstSavedFilterName = "firstSavedFilterName"
)

var (
	sceneIDs       []int
	imageIDs       []int
	performerIDs   []int
	movieIDs       []int
	galleryIDs     []int
	tagIDs         []int
	studioIDs      []int
	markerIDs      []int
	savedFilterIDs []int

	tagNames       []string
	studioNames    []string
	movieNames     []string
	performerNames []string
)

type idAssociation struct {
	first  int
	second int
}

type linkMap map[int][]int

func (m linkMap) reverseLookup(idx int) []int {
	var result []int

	for k, v := range m {
		for _, vv := range v {
			if vv == idx {
				result = append(result, k)
			}
		}
	}

	return result
}

var (
	sceneTags = linkMap{
		sceneIdxWithTag:          {tagIdxWithScene},
		sceneIdxWithTwoTags:      {tagIdx1WithScene, tagIdx2WithScene},
		sceneIdxWithMarkerAndTag: {tagIdx3WithScene},
	}

	scenePerformers = linkMap{
		sceneIdxWithPerformer:        {performerIdxWithScene},
		sceneIdxWithTwoPerformers:    {performerIdx1WithScene, performerIdx2WithScene},
		sceneIdxWithPerformerTag:     {performerIdxWithTag},
		sceneIdxWithPerformerTwoTags: {performerIdxWithTwoTags},
		sceneIdx1WithPerformer:       {performerIdxWithTwoScenes},
		sceneIdx2WithPerformer:       {performerIdxWithTwoScenes},
		sceneIdxWithStudioPerformer:  {performerIdxWithSceneStudio},
	}

	sceneGalleries = linkMap{
		sceneIdxWithGallery: {galleryIdxWithScene},
	}

	sceneMovies = linkMap{
		sceneIdxWithMovie: {movieIdxWithScene},
	}

	sceneStudios = map[int]int{
		sceneIdxWithStudio:           studioIdxWithScene,
		sceneIdx1WithStudio:          studioIdxWithTwoScenes,
		sceneIdx2WithStudio:          studioIdxWithTwoScenes,
		sceneIdxWithStudioPerformer:  studioIdxWithScenePerformer,
		sceneIdxWithGrandChildStudio: studioIdxWithGrandParent,
	}
)

type markerSpec struct {
	sceneIdx      int
	primaryTagIdx int
	tagIdxs       []int
}

var (
	// indexed by marker
	markerSpecs = []markerSpec{
		{sceneIdxWithMarkers, tagIdxWithPrimaryMarkers, nil},
		{sceneIdxWithMarkers, tagIdxWithPrimaryMarkers, []int{tagIdxWithMarkers}},
		{sceneIdxWithMarkerAndTag, tagIdxWithPrimaryMarkers, nil},
	}
)

var (
	imageGalleries = linkMap{
		imageIdxWithGallery:      {galleryIdxWithImage},
		imageIdx1WithGallery:     {galleryIdxWithTwoImages},
		imageIdx2WithGallery:     {galleryIdxWithTwoImages},
		imageIdxWithTwoGalleries: {galleryIdx1WithImage, galleryIdx2WithImage},
	}
	imageStudios = map[int]int{
		imageIdxWithStudio:           studioIdxWithImage,
		imageIdx1WithStudio:          studioIdxWithTwoImages,
		imageIdx2WithStudio:          studioIdxWithTwoImages,
		imageIdxWithStudioPerformer:  studioIdxWithImagePerformer,
		imageIdxWithGrandChildStudio: studioIdxWithGrandParent,
	}
	imageTags = linkMap{
		imageIdxWithTag:     {tagIdxWithImage},
		imageIdxWithTwoTags: {tagIdx1WithImage, tagIdx2WithImage},
	}
	imagePerformers = linkMap{
		imageIdxWithPerformer:        {performerIdxWithImage},
		imageIdxWithTwoPerformers:    {performerIdx1WithImage, performerIdx2WithImage},
		imageIdxWithPerformerTag:     {performerIdxWithTag},
		imageIdxWithPerformerTwoTags: {performerIdxWithTwoTags},
		imageIdx1WithPerformer:       {performerIdxWithTwoImages},
		imageIdx2WithPerformer:       {performerIdxWithTwoImages},
		imageIdxWithStudioPerformer:  {performerIdxWithImageStudio},
	}
)

var (
	galleryPerformers = linkMap{
		galleryIdxWithPerformer:        {performerIdxWithGallery},
		galleryIdxWithTwoPerformers:    {performerIdx1WithGallery, performerIdx2WithGallery},
		galleryIdxWithPerformerTag:     {performerIdxWithTag},
		galleryIdxWithPerformerTwoTags: {performerIdxWithTwoTags},
		galleryIdx1WithPerformer:       {performerIdxWithTwoGalleries},
		galleryIdx2WithPerformer:       {performerIdxWithTwoGalleries},
		galleryIdxWithStudioPerformer:  {performerIdxWithGalleryStudio},
	}

	galleryStudios = map[int]int{
		galleryIdxWithStudio:           studioIdxWithGallery,
		galleryIdx1WithStudio:          studioIdxWithTwoGalleries,
		galleryIdx2WithStudio:          studioIdxWithTwoGalleries,
		galleryIdxWithStudioPerformer:  studioIdxWithGalleryPerformer,
		galleryIdxWithGrandChildStudio: studioIdxWithGrandParent,
	}

	galleryTags = linkMap{
		galleryIdxWithTag:     {tagIdxWithGallery},
		galleryIdxWithTwoTags: {tagIdx1WithGallery, tagIdx2WithGallery},
	}
)

var (
	movieStudioLinks = [][2]int{
		{movieIdxWithStudio, studioIdxWithMovie},
	}
)

var (
	studioParentLinks = [][2]int{
		{studioIdxWithChildStudio, studioIdxWithParentStudio},
		{studioIdxWithGrandChild, studioIdxWithParentAndChild},
		{studioIdxWithParentAndChild, studioIdxWithGrandParent},
	}
)

var (
	performerTagLinks = [][2]int{
		{performerIdxWithTag, tagIdxWithPerformer},
		{performerIdxWithTwoTags, tagIdx1WithPerformer},
		{performerIdxWithTwoTags, tagIdx2WithPerformer},
	}
)

var (
	tagParentLinks = [][2]int{
		{tagIdxWithChildTag, tagIdxWithParentTag},
		{tagIdxWithGrandChild, tagIdxWithParentAndChild},
		{tagIdxWithParentAndChild, tagIdxWithGrandParent},
	}
)

func indexesToIDs(ids []int, indexes []int) []int {
	if len(indexes) == 0 {
		return nil
	}

	ret := make([]int, len(indexes))
	for i, idx := range indexes {
		ret[i] = ids[idx]
	}

	return ret
}

var db *sqlite.Database

func TestMain(m *testing.M) {
	ret := runTests(m)
	os.Exit(ret)
}

func withTxn(f func(ctx context.Context) error) error {
	return txn.WithTxn(context.Background(), db, f)
}

func withRollbackTxn(f func(ctx context.Context) error) error {
	var ret error
	withTxn(func(ctx context.Context) error {
		ret = f(ctx)
		return errors.New("fake error for rollback")
	})

	return ret
}

func runWithRollbackTxn(t *testing.T, name string, f func(t *testing.T, ctx context.Context)) {
	withRollbackTxn(func(ctx context.Context) error {
		t.Run(name, func(t *testing.T) {
			f(t, ctx)
		})
		return nil
	})
}

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
	if err := withTxn(func(ctx context.Context) error {
		if err := createMovies(ctx, sqlite.MovieReaderWriter, moviesNameCase, moviesNameNoCase); err != nil {
			return fmt.Errorf("error creating movies: %s", err.Error())
		}

		if err := createPerformers(ctx, sqlite.PerformerReaderWriter, performersNameCase, performersNameNoCase); err != nil {
			return fmt.Errorf("error creating performers: %s", err.Error())
		}

		if err := createTags(ctx, sqlite.TagReaderWriter, tagsNameCase, tagsNameNoCase); err != nil {
			return fmt.Errorf("error creating tags: %s", err.Error())
		}

		if err := createStudios(ctx, sqlite.StudioReaderWriter, studiosNameCase, studiosNameNoCase); err != nil {
			return fmt.Errorf("error creating studios: %s", err.Error())
		}

		if err := createGalleries(ctx, sqlite.GalleryReaderWriter, totalGalleries); err != nil {
			return fmt.Errorf("error creating galleries: %s", err.Error())
		}

		if err := createScenes(ctx, sqlite.SceneReaderWriter, totalScenes); err != nil {
			return fmt.Errorf("error creating scenes: %s", err.Error())
		}

		if err := createImages(ctx, sqlite.ImageReaderWriter, totalImages); err != nil {
			return fmt.Errorf("error creating images: %s", err.Error())
		}

		if err := addTagImage(ctx, sqlite.TagReaderWriter, tagIdxWithCoverImage); err != nil {
			return fmt.Errorf("error adding tag image: %s", err.Error())
		}

		if err := createSavedFilters(ctx, sqlite.SavedFilterReaderWriter, totalSavedFilters); err != nil {
			return fmt.Errorf("error creating saved filters: %s", err.Error())
		}

		if err := linkPerformerTags(ctx, sqlite.PerformerReaderWriter); err != nil {
			return fmt.Errorf("error linking performer tags: %s", err.Error())
		}

		if err := linkMovieStudios(ctx, sqlite.MovieReaderWriter); err != nil {
			return fmt.Errorf("error linking movie studios: %s", err.Error())
		}

		if err := linkStudiosParent(ctx, sqlite.StudioReaderWriter); err != nil {
			return fmt.Errorf("error linking studios parent: %s", err.Error())
		}

		if err := linkTagsParent(ctx, sqlite.TagReaderWriter); err != nil {
			return fmt.Errorf("error linking tags parent: %s", err.Error())
		}

		for _, ms := range markerSpecs {
			if err := createMarker(ctx, sqlite.SceneMarkerReaderWriter, ms); err != nil {
				return fmt.Errorf("error creating scene marker: %s", err.Error())
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func getPrefixedStringValue(prefix string, index int, field string) string {
	return fmt.Sprintf("%s_%04d_%s", prefix, index, field)
}

func getPrefixedNullStringValue(prefix string, index int, field string) sql.NullString {
	if index > 0 && index%5 == 0 {
		return sql.NullString{}
	}
	if index > 0 && index%6 == 0 {
		return sql.NullString{
			String: "",
			Valid:  true,
		}
	}
	return sql.NullString{
		String: getPrefixedStringValue(prefix, index, field),
		Valid:  true,
	}
}

func getSceneStringValue(index int, field string) string {
	return getPrefixedStringValue("scene", index, field)
}

func getSceneStringPtr(index int, field string) *string {
	v := getPrefixedStringValue("scene", index, field)
	return &v
}

func getSceneNullStringPtr(index int, field string) *string {
	return getStringPtrFromNullString(getPrefixedNullStringValue("scene", index, field))
}

func getSceneEmptyString(index int, field string) string {
	v := getSceneNullStringPtr(index, field)
	if v == nil {
		return ""
	}

	return *v
}

func getSceneTitle(index int) string {
	switch index {
	case sceneIdxWithSpacedName:
		return spacedSceneTitle
	default:
		return getSceneStringValue(index, titleField)
	}
}

func getRating(index int) sql.NullInt64 {
	rating := index % 6
	return sql.NullInt64{Int64: int64(rating), Valid: rating > 0}
}

func getIntPtr(r sql.NullInt64) *int {
	if !r.Valid {
		return nil
	}

	v := int(r.Int64)
	return &v
}

func getStringPtrFromNullString(r sql.NullString) *string {
	if !r.Valid || r.String == "" {
		return nil
	}

	v := r.String
	return &v
}

func getStringPtr(r string) *string {
	if r == "" {
		return nil
	}

	return &r
}

func getEmptyStringFromPtr(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

func getOCounter(index int) int {
	return index % 3
}

func getSceneDuration(index int) *float64 {
	duration := index % 4
	duration = duration * 100

	if duration == 0 {
		return nil
	}
	v := float64(duration) + 0.432
	return &v
}

func getHeight(index int) sql.NullInt64 {
	heights := []int64{0, 200, 240, 300, 480, 700, 720, 800, 1080, 1500, 2160, 3000}
	height := heights[index%len(heights)]
	return sql.NullInt64{
		Int64: height,
		Valid: height != 0,
	}
}

func getWidth(index int) sql.NullInt64 {
	height := getHeight(index)
	return sql.NullInt64{
		Int64: height.Int64 * 2,
		Valid: height.Valid,
	}
}

func getObjectDate(index int) models.SQLiteDate {
	dates := []string{"null", "", "0001-01-01", "2001-02-03"}
	date := dates[index%len(dates)]
	return models.SQLiteDate{
		String: date,
		Valid:  date != "null",
	}
}

func getObjectDateObject(index int) *models.Date {
	d := getObjectDate(index)
	if !d.Valid {
		return nil
	}

	ret := models.NewDate(d.String)
	return &ret
}

func sceneStashID(i int) models.StashID {
	return models.StashID{
		StashID:  getSceneStringValue(i, "stashid"),
		Endpoint: getSceneStringValue(i, "endpoint"),
	}
}

func makeScene(i int) *models.Scene {
	title := getSceneTitle(i)
	checksum := getSceneStringValue(i, checksumField)
	oshash := getSceneStringValue(i, "oshash")
	details := getSceneStringValue(i, "Details")

	var studioID *int
	if _, ok := sceneStudios[i]; ok {
		v := studioIDs[sceneStudios[i]]
		studioID = &v
	}

	gids := indexesToIDs(galleryIDs, sceneGalleries[i])
	pids := indexesToIDs(performerIDs, scenePerformers[i])
	tids := indexesToIDs(tagIDs, sceneTags[i])

	mids := indexesToIDs(movieIDs, sceneMovies[i])

	var movies []models.MoviesScenes
	if len(mids) > 0 {
		movies = make([]models.MoviesScenes, len(mids))
		for i, m := range mids {
			movies[i] = models.MoviesScenes{
				MovieID: m,
			}
		}
	}

	return &models.Scene{
		Path:         getSceneStringValue(i, pathField),
		Title:        title,
		Checksum:     &checksum,
		OSHash:       &oshash,
		Details:      details,
		URL:          getSceneEmptyString(i, urlField),
		Rating:       getIntPtr(getRating(i)),
		OCounter:     getOCounter(i),
		Duration:     getSceneDuration(i),
		Height:       getIntPtr(getHeight(i)),
		Width:        getIntPtr(getWidth(i)),
		Date:         getObjectDateObject(i),
		StudioID:     studioID,
		GalleryIDs:   gids,
		PerformerIDs: pids,
		TagIDs:       tids,
		Movies:       movies,
		StashIDs: []models.StashID{
			sceneStashID(i),
		},
	}
}

func createScenes(ctx context.Context, sqb models.SceneReaderWriter, n int) error {
	for i := 0; i < n; i++ {
		scene := makeScene(i)

		if err := sqb.Create(ctx, scene); err != nil {
			return fmt.Errorf("Error creating scene %v+: %s", scene, err.Error())
		}

		sceneIDs = append(sceneIDs, scene.ID)
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

func makeImage(i int) *models.Image {
	title := getImageStringValue(i, titleField)
	var studioID *int
	if _, ok := imageStudios[i]; ok {
		v := studioIDs[imageStudios[i]]
		studioID = &v
	}

	gids := indexesToIDs(galleryIDs, imageGalleries[i])
	pids := indexesToIDs(performerIDs, imagePerformers[i])
	tids := indexesToIDs(tagIDs, imageTags[i])

	return &models.Image{
		Path:         getImagePath(i),
		Title:        title,
		Checksum:     getImageStringValue(i, checksumField),
		Rating:       getIntPtr(getRating(i)),
		OCounter:     getOCounter(i),
		Height:       getIntPtr(getHeight(i)),
		Width:        getIntPtr(getWidth(i)),
		StudioID:     studioID,
		GalleryIDs:   gids,
		PerformerIDs: pids,
		TagIDs:       tids,
	}
}

func createImages(ctx context.Context, qb models.ImageReaderWriter, n int) error {
	for i := 0; i < n; i++ {
		image := makeImage(i)

		err := qb.Create(ctx, image)

		if err != nil {
			return fmt.Errorf("Error creating image %v+: %s", image, err.Error())
		}

		imageIDs = append(imageIDs, image.ID)
	}

	return nil
}

func getGalleryStringValue(index int, field string) string {
	return getPrefixedStringValue("gallery", index, field)
}

func getGalleryNullStringValue(index int, field string) sql.NullString {
	return getPrefixedNullStringValue("gallery", index, field)
}

func getGalleryNullStringPtr(index int, field string) *string {
	return getStringPtr(getPrefixedStringValue("gallery", index, field))
}

func makeGallery(i int, includeScenes bool) *models.Gallery {
	var studioID *int
	if _, ok := galleryStudios[i]; ok {
		v := studioIDs[galleryStudios[i]]
		studioID = &v
	}

	pids := indexesToIDs(performerIDs, galleryPerformers[i])
	tids := indexesToIDs(tagIDs, galleryTags[i])

	ret := &models.Gallery{
		Path:         getStringPtr(getGalleryStringValue(i, pathField)),
		Title:        getGalleryStringValue(i, titleField),
		URL:          getGalleryNullStringValue(i, urlField).String,
		Checksum:     getGalleryStringValue(i, checksumField),
		Rating:       getIntPtr(getRating(i)),
		Date:         getObjectDateObject(i),
		StudioID:     studioID,
		PerformerIDs: pids,
		TagIDs:       tids,
	}

	if includeScenes {
		ret.SceneIDs = indexesToIDs(sceneIDs, sceneGalleries.reverseLookup(i))
	}

	return ret
}

func createGalleries(ctx context.Context, gqb models.GalleryReaderWriter, n int) error {
	for i := 0; i < n; i++ {
		// scene relationship will be created with scenes
		const includeScenes = false
		gallery := makeGallery(i, includeScenes)

		err := gqb.Create(ctx, gallery)

		if err != nil {
			return fmt.Errorf("Error creating gallery %v+: %s", gallery, err.Error())
		}

		galleryIDs = append(galleryIDs, gallery.ID)
	}

	return nil
}

func getMovieStringValue(index int, field string) string {
	return getPrefixedStringValue("movie", index, field)
}

func getMovieNullStringValue(index int, field string) sql.NullString {
	return getPrefixedNullStringValue("movie", index, field)
}

// createMoviees creates n movies with plain Name and o movies with camel cased NaMe included
func createMovies(ctx context.Context, mqb models.MovieReaderWriter, n int, o int) error {
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
			URL:      getMovieNullStringValue(index, urlField),
			Checksum: md5.FromString(name),
		}

		created, err := mqb.Create(ctx, movie)

		if err != nil {
			return fmt.Errorf("Error creating movie [%d] %v+: %s", i, movie, err.Error())
		}

		movieIDs = append(movieIDs, created.ID)
		movieNames = append(movieNames, created.Name.String)
	}

	return nil
}

func getPerformerStringValue(index int, field string) string {
	return getPrefixedStringValue("performer", index, field)
}

func getPerformerNullStringValue(index int, field string) sql.NullString {
	return getPrefixedNullStringValue("performer", index, field)
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

func getPerformerDeathDate(index int) models.SQLiteDate {
	if index != 5 {
		return models.SQLiteDate{}
	}

	deathDate := time.Now()
	deathDate = deathDate.AddDate(-index+1, -1, -1)
	return models.SQLiteDate{
		String: deathDate.Format("2006-01-02"),
		Valid:  true,
	}
}

func getPerformerCareerLength(index int) *string {
	if index%5 == 0 {
		return nil
	}

	ret := fmt.Sprintf("20%2d", index)
	return &ret
}

func getIgnoreAutoTag(index int) bool {
	return index%5 == 0
}

// createPerformers creates n performers with plain Name and o performers with camel cased NaMe included
func createPerformers(ctx context.Context, pqb models.PerformerReaderWriter, n int, o int) error {
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
			URL:      getPerformerNullStringValue(i, urlField),
			Favorite: sql.NullBool{Bool: getPerformerBoolValue(i), Valid: true},
			Birthdate: models.SQLiteDate{
				String: getPerformerBirthdate(i),
				Valid:  true,
			},
			DeathDate:     getPerformerDeathDate(i),
			Details:       sql.NullString{String: getPerformerStringValue(i, "Details"), Valid: true},
			Ethnicity:     sql.NullString{String: getPerformerStringValue(i, "Ethnicity"), Valid: true},
			Rating:        getRating(i),
			IgnoreAutoTag: getIgnoreAutoTag(i),
		}

		careerLength := getPerformerCareerLength(i)
		if careerLength != nil {
			performer.CareerLength = models.NullString(*careerLength)
		}

		created, err := pqb.Create(ctx, performer)

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
	if id == tagIDs[tagIdx1WithScene] || id == tagIDs[tagIdx2WithScene] || id == tagIDs[tagIdxWithScene] || id == tagIDs[tagIdx3WithScene] {
		return 1
	}

	return 0
}

func getTagMarkerCount(id int) int {
	if id == tagIDs[tagIdxWithPrimaryMarkers] {
		return 3
	}

	if id == tagIDs[tagIdxWithMarkers] {
		return 1
	}

	return 0
}

func getTagImageCount(id int) int {
	if id == tagIDs[tagIdx1WithImage] || id == tagIDs[tagIdx2WithImage] || id == tagIDs[tagIdxWithImage] {
		return 1
	}

	return 0
}

func getTagGalleryCount(id int) int {
	if id == tagIDs[tagIdx1WithGallery] || id == tagIDs[tagIdx2WithGallery] || id == tagIDs[tagIdxWithGallery] {
		return 1
	}

	return 0
}

func getTagPerformerCount(id int) int {
	if id == tagIDs[tagIdx1WithPerformer] || id == tagIDs[tagIdx2WithPerformer] || id == tagIDs[tagIdxWithPerformer] {
		return 1
	}

	return 0
}

func getTagParentCount(id int) int {
	if id == tagIDs[tagIdxWithParentTag] || id == tagIDs[tagIdxWithGrandParent] || id == tagIDs[tagIdxWithParentAndChild] {
		return 1
	}

	return 0
}

func getTagChildCount(id int) int {
	if id == tagIDs[tagIdxWithChildTag] || id == tagIDs[tagIdxWithGrandChild] || id == tagIDs[tagIdxWithParentAndChild] {
		return 1
	}

	return 0
}

//createTags creates n tags with plain Name and o tags with camel cased NaMe included
func createTags(ctx context.Context, tqb models.TagReaderWriter, n int, o int) error {
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
			Name:          getTagStringValue(index, name),
			IgnoreAutoTag: getIgnoreAutoTag(i),
		}

		created, err := tqb.Create(ctx, tag)

		if err != nil {
			return fmt.Errorf("Error creating tag %v+: %s", tag, err.Error())
		}

		// add alias
		alias := getTagStringValue(i, "Alias")
		if err := tqb.UpdateAliases(ctx, created.ID, []string{alias}); err != nil {
			return fmt.Errorf("error setting tag alias: %s", err.Error())
		}

		tagIDs = append(tagIDs, created.ID)
		tagNames = append(tagNames, created.Name)
	}

	return nil
}

func getStudioStringValue(index int, field string) string {
	return getPrefixedStringValue("studio", index, field)
}

func getStudioNullStringValue(index int, field string) sql.NullString {
	return getPrefixedNullStringValue("studio", index, field)
}

func createStudio(ctx context.Context, sqb models.StudioReaderWriter, name string, parentID *int64) (*models.Studio, error) {
	studio := models.Studio{
		Name:     sql.NullString{String: name, Valid: true},
		Checksum: md5.FromString(name),
	}

	if parentID != nil {
		studio.ParentID = sql.NullInt64{Int64: *parentID, Valid: true}
	}

	return createStudioFromModel(ctx, sqb, studio)
}

func createStudioFromModel(ctx context.Context, sqb models.StudioReaderWriter, studio models.Studio) (*models.Studio, error) {
	created, err := sqb.Create(ctx, studio)

	if err != nil {
		return nil, fmt.Errorf("Error creating studio %v+: %s", studio, err.Error())
	}

	return created, nil
}

// createStudios creates n studios with plain Name and o studios with camel cased NaMe included
func createStudios(ctx context.Context, sqb models.StudioReaderWriter, n int, o int) error {
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
		studio := models.Studio{
			Name:          sql.NullString{String: name, Valid: true},
			Checksum:      md5.FromString(name),
			URL:           getStudioNullStringValue(index, urlField),
			IgnoreAutoTag: getIgnoreAutoTag(i),
		}
		created, err := createStudioFromModel(ctx, sqb, studio)

		if err != nil {
			return err
		}

		// add alias
		alias := getStudioStringValue(i, "Alias")
		if err := sqb.UpdateAliases(ctx, created.ID, []string{alias}); err != nil {
			return fmt.Errorf("error setting studio alias: %s", err.Error())
		}

		studioIDs = append(studioIDs, created.ID)
		studioNames = append(studioNames, created.Name.String)
	}

	return nil
}

func createMarker(ctx context.Context, mqb models.SceneMarkerReaderWriter, markerSpec markerSpec) error {
	marker := models.SceneMarker{
		SceneID:      sql.NullInt64{Int64: int64(sceneIDs[markerSpec.sceneIdx]), Valid: true},
		PrimaryTagID: tagIDs[markerSpec.primaryTagIdx],
	}

	created, err := mqb.Create(ctx, marker)

	if err != nil {
		return fmt.Errorf("error creating marker %v+: %w", marker, err)
	}

	markerIDs = append(markerIDs, created.ID)

	if len(markerSpec.tagIdxs) > 0 {
		newTagIDs := []int{}

		for _, tagIdx := range markerSpec.tagIdxs {
			newTagIDs = append(newTagIDs, tagIDs[tagIdx])
		}

		if err := mqb.UpdateTags(ctx, created.ID, newTagIDs); err != nil {
			return fmt.Errorf("error creating marker/tag join: %w", err)
		}
	}

	return nil
}

func getSavedFilterMode(index int) models.FilterMode {
	switch index {
	case savedFilterIdxScene, savedFilterIdxDefaultScene:
		return models.FilterModeScenes
	case savedFilterIdxImage, savedFilterIdxDefaultImage:
		return models.FilterModeImages
	default:
		return models.FilterModeScenes
	}
}

func getSavedFilterName(index int) string {
	if index <= savedFilterIdxDefaultImage {
		// empty string for default filters
		return ""
	}

	if index <= savedFilterIdxImage {
		// use the same name for the first two - should be possible
		return firstSavedFilterName
	}

	return getPrefixedStringValue("savedFilter", index, "Name")
}

func createSavedFilters(ctx context.Context, qb models.SavedFilterReaderWriter, n int) error {
	for i := 0; i < n; i++ {
		savedFilter := models.SavedFilter{
			Mode:   getSavedFilterMode(i),
			Name:   getSavedFilterName(i),
			Filter: getPrefixedStringValue("savedFilter", i, "Filter"),
		}

		created, err := qb.Create(ctx, savedFilter)

		if err != nil {
			return fmt.Errorf("Error creating saved filter %v+: %s", savedFilter, err.Error())
		}

		savedFilterIDs = append(savedFilterIDs, created.ID)
	}

	return nil
}

func doLinks(links [][2]int, fn func(idx1, idx2 int) error) error {
	for _, l := range links {
		if err := fn(l[0], l[1]); err != nil {
			return err
		}
	}

	return nil
}

func linkPerformerTags(ctx context.Context, qb models.PerformerReaderWriter) error {
	return doLinks(performerTagLinks, func(performerIndex, tagIndex int) error {
		performerID := performerIDs[performerIndex]
		tagID := tagIDs[tagIndex]
		tagIDs, err := qb.GetTagIDs(ctx, performerID)
		if err != nil {
			return err
		}

		tagIDs = intslice.IntAppendUnique(tagIDs, tagID)

		return qb.UpdateTags(ctx, performerID, tagIDs)
	})
}

func linkMovieStudios(ctx context.Context, mqb models.MovieWriter) error {
	return doLinks(movieStudioLinks, func(movieIndex, studioIndex int) error {
		movie := models.MoviePartial{
			ID:       movieIDs[movieIndex],
			StudioID: &sql.NullInt64{Int64: int64(studioIDs[studioIndex]), Valid: true},
		}
		_, err := mqb.Update(ctx, movie)

		return err
	})
}

func linkStudiosParent(ctx context.Context, qb models.StudioWriter) error {
	return doLinks(studioParentLinks, func(parentIndex, childIndex int) error {
		studio := models.StudioPartial{
			ID:       studioIDs[childIndex],
			ParentID: &sql.NullInt64{Int64: int64(studioIDs[parentIndex]), Valid: true},
		}
		_, err := qb.Update(ctx, studio)

		return err
	})
}

func linkTagsParent(ctx context.Context, qb models.TagReaderWriter) error {
	return doLinks(tagParentLinks, func(parentIndex, childIndex int) error {
		tagID := tagIDs[childIndex]
		parentTags, err := qb.FindByChildTagID(ctx, tagID)
		if err != nil {
			return err
		}

		var parentIDs []int
		for _, parentTag := range parentTags {
			parentIDs = append(parentIDs, parentTag.ID)
		}

		parentIDs = append(parentIDs, tagIDs[parentIndex])

		return qb.UpdateParentTags(ctx, tagID, parentIDs)
	})
}

func addTagImage(ctx context.Context, qb models.TagWriter, tagIndex int) error {
	return qb.UpdateImage(ctx, tagIDs[tagIndex], models.DefaultTagImage)
}
