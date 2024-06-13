//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"

	// necessary to register custom migrations
	_ "github.com/stashapp/stash/pkg/sqlite/migrations"
)

const (
	spacedSceneTitle = "zzz yyy xxx"
)

const (
	folderIdxWithSubFolder = iota
	folderIdxWithParentFolder
	folderIdxWithFiles
	folderIdxInZip

	folderIdxForObjectFiles
	folderIdxWithImageFiles
	folderIdxWithGalleryFiles
	folderIdxWithSceneFiles

	totalFolders
)

const (
	fileIdxZip = iota
	fileIdxInZip

	fileIdxStartVideoFiles
	fileIdxStartImageFiles
	fileIdxStartGalleryFiles

	totalFiles
)

const (
	sceneIdxWithMovie = iota
	sceneIdxWithGallery
	sceneIdxWithPerformer
	sceneIdx1WithPerformer
	sceneIdx2WithPerformer
	sceneIdxWithTwoPerformers
	sceneIdxWithThreePerformers
	sceneIdxWithTag
	sceneIdxWithTwoTags
	sceneIdxWithThreeTags
	sceneIdxWithMarkerAndTag
	sceneIdxWithMarkerTwoTags
	sceneIdxWithStudio
	sceneIdx1WithStudio
	sceneIdx2WithStudio
	sceneIdxWithMarkers
	sceneIdxWithPerformerTag
	sceneIdxWithTwoPerformerTag
	sceneIdxWithPerformerTwoTags
	sceneIdxWithSpacedName
	sceneIdxWithStudioPerformer
	sceneIdxWithGrandChildStudio
	sceneIdxMissingPhash
	sceneIdxWithPerformerParentTag
	// new indexes above
	lastSceneIdx

	totalScenes = lastSceneIdx + 3
)

const dupeScenePhashes = 2

const (
	imageIdxWithGallery = iota
	imageIdx1WithGallery
	imageIdx2WithGallery
	imageIdxWithTwoGalleries
	imageIdxWithPerformer
	imageIdx1WithPerformer
	imageIdx2WithPerformer
	imageIdxWithTwoPerformers
	imageIdxWithThreePerformers
	imageIdxWithTag
	imageIdxWithTwoTags
	imageIdxWithThreeTags
	imageIdxWithStudio
	imageIdx1WithStudio
	imageIdx2WithStudio
	imageIdxWithStudioPerformer
	imageIdxInZip
	imageIdxWithPerformerTag
	imageIdxWithTwoPerformerTag
	imageIdxWithPerformerTwoTags
	imageIdxWithGrandChildStudio
	imageIdxWithPerformerParentTag
	// new indexes above
	totalImages
)

const (
	performerIdxWithScene = iota
	performerIdx1WithScene
	performerIdx2WithScene
	performerIdx3WithScene
	performerIdxWithTwoScenes
	performerIdxWithImage
	performerIdxWithTwoImages
	performerIdx1WithImage
	performerIdx2WithImage
	performerIdx3WithImage
	performerIdxWithTag
	performerIdx2WithTag
	performerIdxWithTwoTags
	performerIdxWithGallery
	performerIdxWithTwoGalleries
	performerIdx1WithGallery
	performerIdx2WithGallery
	performerIdx3WithGallery
	performerIdxWithSceneStudio
	performerIdxWithImageStudio
	performerIdxWithGalleryStudio
	performerIdxWithParentTag
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
	galleryIdxWithChapters
	galleryIdxWithImage
	galleryIdx1WithImage
	galleryIdx2WithImage
	galleryIdxWithTwoImages
	galleryIdxWithPerformer
	galleryIdx1WithPerformer
	galleryIdx2WithPerformer
	galleryIdxWithTwoPerformers
	galleryIdxWithThreePerformers
	galleryIdxWithTag
	galleryIdxWithTwoTags
	galleryIdxWithThreeTags
	galleryIdxWithStudio
	galleryIdx1WithStudio
	galleryIdx2WithStudio
	galleryIdxWithPerformerTag
	galleryIdxWithTwoPerformerTag
	galleryIdxWithPerformerTwoTags
	galleryIdxWithStudioPerformer
	galleryIdxWithGrandChildStudio
	galleryIdxWithoutFile
	galleryIdxWithPerformerParentTag
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
	tagIdx3WithImage
	tagIdxWithPerformer
	tagIdx1WithPerformer
	tagIdx2WithPerformer
	tagIdxWithGallery
	tagIdx1WithGallery
	tagIdx2WithGallery
	tagIdx3WithGallery
	tagIdxWithChildTag
	tagIdxWithParentTag
	tagIdxWithGrandChild
	tagIdxWithParentAndChild
	tagIdxWithGrandParent
	tagIdx2WithMarkers
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
	chapterIdxWithGallery = iota
	totalChapters
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
	folderIDs      []models.FolderID
	fileIDs        []models.FileID
	sceneFileIDs   []models.FileID
	imageFileIDs   []models.FileID
	galleryFileIDs []models.FileID
	chapterIDs     []int

	sceneIDs       []int
	imageIDs       []int
	performerIDs   []int
	movieIDs       []int
	galleryIDs     []int
	tagIDs         []int
	studioIDs      []int
	markerIDs      []int
	savedFilterIDs []int

	folderPaths []string

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
	folderParentFolders = map[int]int{
		folderIdxWithParentFolder: folderIdxWithSubFolder,
		folderIdxWithSceneFiles:   folderIdxForObjectFiles,
		folderIdxWithImageFiles:   folderIdxForObjectFiles,
		folderIdxWithGalleryFiles: folderIdxForObjectFiles,
	}

	fileFolders = map[int]int{
		fileIdxZip:   folderIdxWithFiles,
		fileIdxInZip: folderIdxInZip,
	}

	folderZipFiles = map[int]int{
		folderIdxInZip: fileIdxZip,
	}

	fileZipFiles = map[int]int{
		fileIdxInZip: fileIdxZip,
	}
)

var (
	sceneTags = linkMap{
		sceneIdxWithTag:           {tagIdxWithScene},
		sceneIdxWithTwoTags:       {tagIdx1WithScene, tagIdx2WithScene},
		sceneIdxWithThreeTags:     {tagIdx1WithScene, tagIdx2WithScene, tagIdx3WithScene},
		sceneIdxWithMarkerAndTag:  {tagIdx3WithScene},
		sceneIdxWithMarkerTwoTags: {tagIdx2WithScene, tagIdx3WithScene},
	}

	scenePerformers = linkMap{
		sceneIdxWithPerformer:          {performerIdxWithScene},
		sceneIdxWithTwoPerformers:      {performerIdx1WithScene, performerIdx2WithScene},
		sceneIdxWithThreePerformers:    {performerIdx1WithScene, performerIdx2WithScene, performerIdx3WithScene},
		sceneIdxWithPerformerTag:       {performerIdxWithTag},
		sceneIdxWithTwoPerformerTag:    {performerIdxWithTag, performerIdx2WithTag},
		sceneIdxWithPerformerTwoTags:   {performerIdxWithTwoTags},
		sceneIdx1WithPerformer:         {performerIdxWithTwoScenes},
		sceneIdx2WithPerformer:         {performerIdxWithTwoScenes},
		sceneIdxWithStudioPerformer:    {performerIdxWithSceneStudio},
		sceneIdxWithPerformerParentTag: {performerIdxWithParentTag},
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
		{sceneIdxWithMarkers, tagIdxWithPrimaryMarkers, []int{tagIdx2WithMarkers}},
		{sceneIdxWithMarkers, tagIdxWithPrimaryMarkers, []int{tagIdxWithMarkers, tagIdx2WithMarkers}},
		{sceneIdxWithMarkerAndTag, tagIdxWithPrimaryMarkers, nil},
		{sceneIdxWithMarkerTwoTags, tagIdxWithPrimaryMarkers, nil},
	}
)

type chapterSpec struct {
	galleryIdx int
	title      string
	imageIndex int
}

var (
	// indexed by chapter
	chapterSpecs = []chapterSpec{
		{galleryIdxWithChapters, "Test1", 10},
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
		imageIdxWithTag:       {tagIdxWithImage},
		imageIdxWithTwoTags:   {tagIdx1WithImage, tagIdx2WithImage},
		imageIdxWithThreeTags: {tagIdx1WithImage, tagIdx2WithImage, tagIdx3WithImage},
	}
	imagePerformers = linkMap{
		imageIdxWithPerformer:          {performerIdxWithImage},
		imageIdxWithTwoPerformers:      {performerIdx1WithImage, performerIdx2WithImage},
		imageIdxWithThreePerformers:    {performerIdx1WithImage, performerIdx2WithImage, performerIdx3WithImage},
		imageIdxWithPerformerTag:       {performerIdxWithTag},
		imageIdxWithTwoPerformerTag:    {performerIdxWithTag, performerIdx2WithTag},
		imageIdxWithPerformerTwoTags:   {performerIdxWithTwoTags},
		imageIdx1WithPerformer:         {performerIdxWithTwoImages},
		imageIdx2WithPerformer:         {performerIdxWithTwoImages},
		imageIdxWithStudioPerformer:    {performerIdxWithImageStudio},
		imageIdxWithPerformerParentTag: {performerIdxWithParentTag},
	}
)

var (
	galleryPerformers = linkMap{
		galleryIdxWithPerformer:          {performerIdxWithGallery},
		galleryIdxWithTwoPerformers:      {performerIdx1WithGallery, performerIdx2WithGallery},
		galleryIdxWithThreePerformers:    {performerIdx1WithGallery, performerIdx2WithGallery, performerIdx3WithGallery},
		galleryIdxWithPerformerTag:       {performerIdxWithTag},
		galleryIdxWithTwoPerformerTag:    {performerIdxWithTag, performerIdx2WithTag},
		galleryIdxWithPerformerTwoTags:   {performerIdxWithTwoTags},
		galleryIdx1WithPerformer:         {performerIdxWithTwoGalleries},
		galleryIdx2WithPerformer:         {performerIdxWithTwoGalleries},
		galleryIdxWithStudioPerformer:    {performerIdxWithGalleryStudio},
		galleryIdxWithPerformerParentTag: {performerIdxWithParentTag},
	}

	galleryStudios = map[int]int{
		galleryIdxWithStudio:           studioIdxWithGallery,
		galleryIdx1WithStudio:          studioIdxWithTwoGalleries,
		galleryIdx2WithStudio:          studioIdxWithTwoGalleries,
		galleryIdxWithStudioPerformer:  studioIdxWithGalleryPerformer,
		galleryIdxWithGrandChildStudio: studioIdxWithGrandParent,
	}

	galleryTags = linkMap{
		galleryIdxWithTag:       {tagIdxWithGallery},
		galleryIdxWithTwoTags:   {tagIdx1WithGallery, tagIdx2WithGallery},
		galleryIdxWithThreeTags: {tagIdx1WithGallery, tagIdx2WithGallery, tagIdx3WithGallery},
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
	performerTags = linkMap{
		performerIdxWithTag:       {tagIdxWithPerformer},
		performerIdx2WithTag:      {tagIdx2WithPerformer},
		performerIdxWithTwoTags:   {tagIdx1WithPerformer, tagIdx2WithPerformer},
		performerIdxWithParentTag: {tagIdxWithParentAndChild},
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
	ret := make([]int, len(indexes))
	for i, idx := range indexes {
		ret[i] = ids[idx]
	}

	return ret
}

func indexFromID(ids []int, id int) int {
	for i, v := range ids {
		if v == id {
			return i
		}
	}

	return -1
}

var db *sqlite.Database

func TestMain(m *testing.M) {
	// initialise empty config - needed by some migrations
	_ = config.InitializeEmpty()

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
	db = sqlite.NewDatabase()
	db.SetBlobStoreOptions(sqlite.BlobStoreOptions{
		UseDatabase: true,
		// don't use filesystem
	})

	if err := db.Open(databaseFile); err != nil {
		panic(fmt.Sprintf("Could not initialize database: %s", err.Error()))
	}

	// defer close and delete the database
	defer testTeardown(databaseFile)

	err = populateDB()
	if err != nil {
		panic(fmt.Sprintf("Could not populate database: %s", err.Error()))
	}

	// run the tests
	return m.Run()
}

func populateDB() error {
	if err := withTxn(func(ctx context.Context) error {
		if err := createFolders(ctx); err != nil {
			return fmt.Errorf("creating folders: %w", err)
		}

		if err := createFiles(ctx); err != nil {
			return fmt.Errorf("creating files: %w", err)
		}

		// TODO - link folders to zip files

		if err := createMovies(ctx, db.Movie, moviesNameCase, moviesNameNoCase); err != nil {
			return fmt.Errorf("error creating movies: %s", err.Error())
		}

		if err := createTags(ctx, db.Tag, tagsNameCase, tagsNameNoCase); err != nil {
			return fmt.Errorf("error creating tags: %s", err.Error())
		}

		if err := createPerformers(ctx, performersNameCase, performersNameNoCase); err != nil {
			return fmt.Errorf("error creating performers: %s", err.Error())
		}

		if err := createStudios(ctx, studiosNameCase, studiosNameNoCase); err != nil {
			return fmt.Errorf("error creating studios: %s", err.Error())
		}

		if err := createGalleries(ctx, totalGalleries); err != nil {
			return fmt.Errorf("error creating galleries: %s", err.Error())
		}

		if err := createScenes(ctx, totalScenes); err != nil {
			return fmt.Errorf("error creating scenes: %s", err.Error())
		}

		if err := createImages(ctx, totalImages); err != nil {
			return fmt.Errorf("error creating images: %s", err.Error())
		}

		if err := addTagImage(ctx, db.Tag, tagIdxWithCoverImage); err != nil {
			return fmt.Errorf("error adding tag image: %s", err.Error())
		}

		if err := createSavedFilters(ctx, db.SavedFilter, totalSavedFilters); err != nil {
			return fmt.Errorf("error creating saved filters: %s", err.Error())
		}

		if err := linkMovieStudios(ctx, db.Movie); err != nil {
			return fmt.Errorf("error linking movie studios: %s", err.Error())
		}

		if err := linkStudiosParent(ctx); err != nil {
			return fmt.Errorf("error linking studios parent: %s", err.Error())
		}

		if err := linkTagsParent(ctx, db.Tag); err != nil {
			return fmt.Errorf("error linking tags parent: %s", err.Error())
		}

		for _, ms := range markerSpecs {
			if err := createMarker(ctx, db.SceneMarker, ms); err != nil {
				return fmt.Errorf("error creating scene marker: %s", err.Error())
			}
		}
		for _, cs := range chapterSpecs {
			if err := createChapter(ctx, db.GalleryChapter, cs); err != nil {
				return fmt.Errorf("error creating gallery chapter: %s", err.Error())
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func getFolderPath(index int, parentFolderIdx *int) string {
	path := getPrefixedStringValue("folder", index, pathField)

	if parentFolderIdx != nil {
		return filepath.Join(folderPaths[*parentFolderIdx], path)
	}

	return path
}

func getFolderModTime(index int) time.Time {
	return time.Date(2000, 1, (index%10)+1, 0, 0, 0, 0, time.UTC)
}

func makeFolder(i int) models.Folder {
	var folderID *models.FolderID
	var folderIdx *int
	if pidx, ok := folderParentFolders[i]; ok {
		folderIdx = &pidx
		v := folderIDs[pidx]
		folderID = &v
	}

	return models.Folder{
		ParentFolderID: folderID,
		DirEntry: models.DirEntry{
			// zip files have to be added after creating files
			ModTime: getFolderModTime(i),
		},
		Path: getFolderPath(i, folderIdx),
	}
}

func createFolders(ctx context.Context) error {
	qb := db.Folder

	for i := 0; i < totalFolders; i++ {
		folder := makeFolder(i)

		if err := qb.Create(ctx, &folder); err != nil {
			return fmt.Errorf("Error creating folder [%d] %v+: %s", i, folder, err.Error())
		}

		folderIDs = append(folderIDs, folder.ID)
		folderPaths = append(folderPaths, folder.Path)
	}

	return nil
}

func getFileBaseName(index int) string {
	return getPrefixedStringValue("file", index, "basename")
}

func getFileStringValue(index int, field string) string {
	return getPrefixedStringValue("file", index, field)
}

func getFileModTime(index int) time.Time {
	return getFolderModTime(index)
}

func getFileFingerprints(index int) []models.Fingerprint {
	return []models.Fingerprint{
		{
			Type:        "MD5",
			Fingerprint: getPrefixedStringValue("file", index, "md5"),
		},
		{
			Type:        "OSHASH",
			Fingerprint: getPrefixedStringValue("file", index, "oshash"),
		},
	}
}

func getFileSize(index int) int64 {
	return int64(index) * 10
}

func getFileDuration(index int) float64 {
	duration := (index % 4) + 1
	duration = duration * 100

	return float64(duration) + 0.432
}

func makeFile(i int) models.File {
	folderID := folderIDs[fileFolders[i]]
	if folderID == 0 {
		folderID = folderIDs[folderIdxWithFiles]
	}

	var zipFileID *models.FileID
	if zipFileIndex, found := fileZipFiles[i]; found {
		zipFileID = &fileIDs[zipFileIndex]
	}

	var ret models.File
	baseFile := &models.BaseFile{
		Basename:       getFileBaseName(i),
		ParentFolderID: folderID,
		DirEntry: models.DirEntry{
			// zip files have to be added after creating files
			ModTime:   getFileModTime(i),
			ZipFileID: zipFileID,
		},
		Fingerprints: getFileFingerprints(i),
		Size:         getFileSize(i),
	}

	ret = baseFile

	if i >= fileIdxStartVideoFiles && i < fileIdxStartImageFiles {
		ret = &models.VideoFile{
			BaseFile:   baseFile,
			Format:     getFileStringValue(i, "format"),
			Width:      getWidth(i),
			Height:     getHeight(i),
			Duration:   getFileDuration(i),
			VideoCodec: getFileStringValue(i, "videoCodec"),
			AudioCodec: getFileStringValue(i, "audioCodec"),
			FrameRate:  getFileDuration(i) * 2,
			BitRate:    int64(getFileDuration(i)) * 3,
		}
	} else if i >= fileIdxStartImageFiles && i < fileIdxStartGalleryFiles {
		ret = &models.ImageFile{
			BaseFile: baseFile,
			Format:   getFileStringValue(i, "format"),
			Width:    getWidth(i),
			Height:   getHeight(i),
		}
	}

	return ret
}

func createFiles(ctx context.Context) error {
	qb := db.File

	for i := 0; i < totalFiles; i++ {
		file := makeFile(i)

		if err := qb.Create(ctx, file); err != nil {
			return fmt.Errorf("Error creating file [%d] %v+: %s", i, file, err.Error())
		}

		fileIDs = append(fileIDs, file.Base().ID)
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

func getScenePhash(index int, field string) int64 {
	return int64(index % (totalScenes - dupeScenePhashes) * 1234)
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
	return sql.NullInt64{Int64: int64(rating * 20), Valid: rating > 0}
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

func getSceneDuration(index int) float64 {
	duration := index + 1
	duration = duration * 100

	return float64(duration) + 0.432
}

func getHeight(index int) int {
	heights := []int{200, 240, 300, 480, 700, 720, 800, 1080, 1500, 2160, 3000}
	height := heights[index%len(heights)]
	return height
}

func getWidth(index int) int {
	height := getHeight(index)
	return height * 2
}

func getObjectDate(index int) *models.Date {
	dates := []string{"null", "2000-01-01", "0001-01-01", "2001-02-03"}
	date := dates[index%len(dates)]

	if date == "null" {
		return nil
	}

	ret, _ := models.ParseDate(date)
	return &ret
}

func sceneStashID(i int) models.StashID {
	return models.StashID{
		StashID:  getSceneStringValue(i, "stashid"),
		Endpoint: getSceneStringValue(i, "endpoint"),
	}
}

func getSceneBasename(index int) string {
	return getSceneStringValue(index, pathField)
}

func makeSceneFile(i int) *models.VideoFile {
	fp := []models.Fingerprint{
		{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: getSceneStringValue(i, checksumField),
		},
		{
			Type:        models.FingerprintTypeOshash,
			Fingerprint: getSceneStringValue(i, "oshash"),
		},
	}

	if i != sceneIdxMissingPhash {
		fp = append(fp, models.Fingerprint{
			Type:        models.FingerprintTypePhash,
			Fingerprint: getScenePhash(i, "phash"),
		})
	}

	return &models.VideoFile{
		BaseFile: &models.BaseFile{
			Path:           getFilePath(folderIdxWithSceneFiles, getSceneBasename(i)),
			Basename:       getSceneBasename(i),
			ParentFolderID: folderIDs[folderIdxWithSceneFiles],
			Fingerprints:   fp,
		},
		Duration: getSceneDuration(i),
		Height:   getHeight(i),
		Width:    getWidth(i),
	}
}

func getScenePlayDuration(index int) float64 {
	if index%5 == 0 {
		return 0
	}

	return float64(index%5) * 123.4
}

func getSceneResumeTime(index int) float64 {
	if index%5 == 0 {
		return 0
	}

	return float64(index%5) * 1.2
}

func makeScene(i int) *models.Scene {
	title := getSceneTitle(i)
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

	movies := make([]models.MoviesScenes, len(mids))
	for i, m := range mids {
		movies[i] = models.MoviesScenes{
			MovieID: m,
		}
	}

	rating := getRating(i)

	return &models.Scene{
		Title:   title,
		Details: details,
		URLs: models.NewRelatedStrings([]string{
			getSceneEmptyString(i, urlField),
		}),
		Rating:       getIntPtr(rating),
		Date:         getObjectDate(i),
		StudioID:     studioID,
		GalleryIDs:   models.NewRelatedIDs(gids),
		PerformerIDs: models.NewRelatedIDs(pids),
		TagIDs:       models.NewRelatedIDs(tids),
		Movies:       models.NewRelatedMovies(movies),
		StashIDs: models.NewRelatedStashIDs([]models.StashID{
			sceneStashID(i),
		}),
		PlayDuration: getScenePlayDuration(i),
		ResumeTime:   getSceneResumeTime(i),
	}
}

func createScenes(ctx context.Context, n int) error {
	sqb := db.Scene
	fqb := db.File

	for i := 0; i < n; i++ {
		f := makeSceneFile(i)
		if err := fqb.Create(ctx, f); err != nil {
			return fmt.Errorf("creating scene file: %w", err)
		}
		sceneFileIDs = append(sceneFileIDs, f.ID)

		scene := makeScene(i)

		if err := sqb.Create(ctx, scene, []models.FileID{f.ID}); err != nil {
			return fmt.Errorf("Error creating scene %v+: %s", scene, err.Error())
		}

		sceneIDs = append(sceneIDs, scene.ID)
	}

	return nil
}

func getImageStringValue(index int, field string) string {
	return fmt.Sprintf("image_%04d_%s", index, field)
}

func getImageNullStringPtr(index int, field string) *string {
	return getStringPtrFromNullString(getPrefixedNullStringValue("image", index, field))
}

func getImageEmptyString(index int, field string) string {
	v := getImageNullStringPtr(index, field)
	if v == nil {
		return ""
	}

	return *v
}

func getImageBasename(index int) string {
	return getImageStringValue(index, pathField)
}

func makeImageFile(i int) *models.ImageFile {
	return &models.ImageFile{
		BaseFile: &models.BaseFile{
			Path:           getFilePath(folderIdxWithImageFiles, getImageBasename(i)),
			Basename:       getImageBasename(i),
			ParentFolderID: folderIDs[folderIdxWithImageFiles],
			Fingerprints: []models.Fingerprint{
				{
					Type:        models.FingerprintTypeMD5,
					Fingerprint: getImageStringValue(i, checksumField),
				},
			},
		},
		Height: getHeight(i),
		Width:  getWidth(i),
	}
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
		Title:  title,
		Rating: getIntPtr(getRating(i)),
		Date:   getObjectDate(i),
		URLs: models.NewRelatedStrings([]string{
			getImageEmptyString(i, urlField),
		}),
		OCounter:     getOCounter(i),
		StudioID:     studioID,
		GalleryIDs:   models.NewRelatedIDs(gids),
		PerformerIDs: models.NewRelatedIDs(pids),
		TagIDs:       models.NewRelatedIDs(tids),
	}
}

func createImages(ctx context.Context, n int) error {
	qb := db.Image
	fqb := db.File

	for i := 0; i < n; i++ {
		f := makeImageFile(i)
		if i == imageIdxInZip {
			f.ZipFileID = &fileIDs[fileIdxZip]
		}

		if err := fqb.Create(ctx, f); err != nil {
			return fmt.Errorf("creating image file: %w", err)
		}
		imageFileIDs = append(imageFileIDs, f.ID)

		image := makeImage(i)

		err := qb.Create(ctx, image, []models.FileID{f.ID})

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
	return getStringPtrFromNullString(getPrefixedNullStringValue("gallery", index, field))
}

func getGalleryEmptyString(index int, field string) string {
	v := getGalleryNullStringPtr(index, field)
	if v == nil {
		return ""
	}

	return *v
}

func getGalleryBasename(index int) string {
	return getGalleryStringValue(index, pathField)
}

func makeGalleryFile(i int) *models.BaseFile {
	return &models.BaseFile{
		Path:           getFilePath(folderIdxWithGalleryFiles, getGalleryBasename(i)),
		Basename:       getGalleryBasename(i),
		ParentFolderID: folderIDs[folderIdxWithGalleryFiles],
		Fingerprints: []models.Fingerprint{
			{
				Type:        models.FingerprintTypeMD5,
				Fingerprint: getGalleryStringValue(i, checksumField),
			},
		},
	}
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
		Title: getGalleryStringValue(i, titleField),
		URLs: models.NewRelatedStrings([]string{
			getGalleryEmptyString(i, urlField),
		}),
		Rating:       getIntPtr(getRating(i)),
		Date:         getObjectDate(i),
		StudioID:     studioID,
		PerformerIDs: models.NewRelatedIDs(pids),
		TagIDs:       models.NewRelatedIDs(tids),
	}

	if includeScenes {
		ret.SceneIDs = models.NewRelatedIDs(indexesToIDs(sceneIDs, sceneGalleries.reverseLookup(i)))
	}

	return ret
}

func createGalleries(ctx context.Context, n int) error {
	gqb := db.Gallery
	fqb := db.File

	for i := 0; i < n; i++ {
		var fileIDs []models.FileID
		if i != galleryIdxWithoutFile {
			f := makeGalleryFile(i)
			if err := fqb.Create(ctx, f); err != nil {
				return fmt.Errorf("creating gallery file: %w", err)
			}
			galleryFileIDs = append(galleryFileIDs, f.ID)
			fileIDs = []models.FileID{f.ID}
		} else {
			galleryFileIDs = append(galleryFileIDs, 0)
		}

		// gallery relationship will be created with galleries
		const includeScenes = false
		gallery := makeGallery(i, includeScenes)

		err := gqb.Create(ctx, gallery, fileIDs)

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

func getMovieNullStringValue(index int, field string) string {
	ret := getPrefixedNullStringValue("movie", index, field)

	return ret.String
}

func getMovieEmptyString(index int, field string) string {
	v := getPrefixedNullStringValue("movie", index, field)
	if !v.Valid {
		return ""
	}

	return v.String
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
			Name: name,
			URLs: models.NewRelatedStrings([]string{
				getMovieEmptyString(i, urlField),
			}),
		}

		err := mqb.Create(ctx, &movie)

		if err != nil {
			return fmt.Errorf("Error creating movie [%d] %v+: %s", i, movie, err.Error())
		}

		movieIDs = append(movieIDs, movie.ID)
		movieNames = append(movieNames, movie.Name)
	}

	return nil
}

func getPerformerStringValue(index int, field string) string {
	return getPrefixedStringValue("performer", index, field)
}

func getPerformerNullStringValue(index int, field string) string {
	ret := getPrefixedNullStringValue("performer", index, field)

	return ret.String
}

func getPerformerBoolValue(index int) bool {
	index = index % 2
	return index == 1
}

func getPerformerBirthdate(index int) *models.Date {
	const minAge = 18
	birthdate := time.Now()
	birthdate = birthdate.AddDate(-minAge-index, -1, -1)

	ret := models.Date{
		Time: birthdate,
	}
	return &ret
}

func getPerformerDeathDate(index int) *models.Date {
	if index != 5 {
		return nil
	}

	deathDate := time.Now()
	deathDate = deathDate.AddDate(-index+1, -1, -1)

	ret := models.Date{
		Time: deathDate,
	}
	return &ret
}

func getPerformerCareerLength(index int) *string {
	if index%5 == 0 {
		return nil
	}

	ret := fmt.Sprintf("20%2d", index)
	return &ret
}

func getPerformerPenisLength(index int) *float64 {
	if index%5 == 0 {
		return nil
	}

	ret := float64(index)
	return &ret
}

func getPerformerCircumcised(index int) *models.CircumisedEnum {
	var ret models.CircumisedEnum
	switch {
	case index%3 == 0:
		return nil
	case index%3 == 1:
		ret = models.CircumisedEnumCut
	default:
		ret = models.CircumisedEnumUncut
	}

	return &ret
}

func getIgnoreAutoTag(index int) bool {
	return index%5 == 0
}

func performerStashID(i int) models.StashID {
	return models.StashID{
		StashID:  getPerformerStringValue(i, "stashid"),
		Endpoint: getPerformerStringValue(i, "endpoint"),
	}
}

func performerAliases(i int) []string {
	if i%5 == 0 {
		return []string{}
	}

	return []string{getPerformerStringValue(i, "alias")}
}

// createPerformers creates n performers with plain Name and o performers with camel cased NaMe included
func createPerformers(ctx context.Context, n int, o int) error {
	pqb := db.Performer

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

		tids := indexesToIDs(tagIDs, performerTags[i])

		performer := models.Performer{
			Name:           getPerformerStringValue(index, name),
			Disambiguation: getPerformerStringValue(index, "disambiguation"),
			Aliases:        models.NewRelatedStrings(performerAliases(index)),
			URL:            getPerformerNullStringValue(i, urlField),
			Favorite:       getPerformerBoolValue(i),
			Birthdate:      getPerformerBirthdate(i),
			DeathDate:      getPerformerDeathDate(i),
			Details:        getPerformerStringValue(i, "Details"),
			Ethnicity:      getPerformerStringValue(i, "Ethnicity"),
			PenisLength:    getPerformerPenisLength(i),
			Circumcised:    getPerformerCircumcised(i),
			Rating:         getIntPtr(getRating(i)),
			IgnoreAutoTag:  getIgnoreAutoTag(i),
			TagIDs:         models.NewRelatedIDs(tids),
		}

		careerLength := getPerformerCareerLength(i)
		if careerLength != nil {
			performer.CareerLength = *careerLength
		}

		if (index+1)%5 != 0 {
			performer.StashIDs = models.NewRelatedStashIDs([]models.StashID{
				performerStashID(i),
			})
		}

		err := pqb.Create(ctx, &performer)

		if err != nil {
			return fmt.Errorf("Error creating performer %v+: %s", performer, err.Error())
		}

		performerIDs = append(performerIDs, performer.ID)
		performerNames = append(performerNames, performer.Name)
	}

	return nil
}
func getTagBoolValue(index int) bool {
	index = index % 2
	return index == 1
}
func getTagStringValue(index int, field string) string {
	return "tag_" + strconv.FormatInt(int64(index), 10) + "_" + field
}

func getTagSceneCount(id int) int {
	idx := indexFromID(tagIDs, id)
	return len(sceneTags.reverseLookup(idx))
}

func getTagMarkerCount(id int) int {
	count := 0
	idx := indexFromID(tagIDs, id)
	for _, s := range markerSpecs {
		if s.primaryTagIdx == idx || sliceutil.Contains(s.tagIdxs, idx) {
			count++
		}
	}

	return count
}

func getTagImageCount(id int) int {
	idx := indexFromID(tagIDs, id)
	return len(imageTags.reverseLookup(idx))
}

func getTagGalleryCount(id int) int {
	idx := indexFromID(tagIDs, id)
	return len(galleryTags.reverseLookup(idx))
}

func getTagPerformerCount(id int) int {
	idx := indexFromID(tagIDs, id)
	return len(performerTags.reverseLookup(idx))
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

// createTags creates n tags with plain Name and o tags with camel cased NaMe included
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

		err := tqb.Create(ctx, &tag)

		if err != nil {
			return fmt.Errorf("Error creating tag %v+: %s", tag, err.Error())
		}

		// add alias
		alias := getTagStringValue(i, "Alias")
		if err := tqb.UpdateAliases(ctx, tag.ID, []string{alias}); err != nil {
			return fmt.Errorf("error setting tag alias: %s", err.Error())
		}

		tagIDs = append(tagIDs, tag.ID)
		tagNames = append(tagNames, tag.Name)
	}

	return nil
}

func getStudioStringValue(index int, field string) string {
	return getPrefixedStringValue("studio", index, field)
}

func getStudioNullStringValue(index int, field string) string {
	ret := getPrefixedNullStringValue("studio", index, field)

	return ret.String
}

func createStudio(ctx context.Context, sqb *sqlite.StudioStore, name string, parentID *int) (*models.Studio, error) {
	studio := models.Studio{
		Name: name,
	}

	if parentID != nil {
		studio.ParentID = parentID
	}

	err := createStudioFromModel(ctx, sqb, &studio)
	if err != nil {
		return nil, err
	}

	return &studio, nil
}

func createStudioFromModel(ctx context.Context, sqb *sqlite.StudioStore, studio *models.Studio) error {
	err := sqb.Create(ctx, studio)

	if err != nil {
		return fmt.Errorf("Error creating studio %v+: %s", studio, err.Error())
	}

	return nil
}

func getStudioBoolValue(index int) bool {
	index = index % 2
	return index == 1
}

// createStudios creates n studios with plain Name and o studios with camel cased NaMe included
func createStudios(ctx context.Context, n int, o int) error {
	sqb := db.Studio
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
			Name:          name,
			URL:           getStudioStringValue(index, urlField),
			Favorite:      getStudioBoolValue(index),
			IgnoreAutoTag: getIgnoreAutoTag(i),
		}
		// only add aliases for some scenes
		if i == studioIdxWithMovie || i%5 == 0 {
			alias := getStudioStringValue(i, "Alias")
			studio.Aliases = models.NewRelatedStrings([]string{alias})
		}
		err := createStudioFromModel(ctx, sqb, &studio)

		if err != nil {
			return err
		}

		studioIDs = append(studioIDs, studio.ID)
		studioNames = append(studioNames, studio.Name)
	}

	return nil
}

func createMarker(ctx context.Context, mqb models.SceneMarkerReaderWriter, markerSpec markerSpec) error {
	marker := models.SceneMarker{
		SceneID:      sceneIDs[markerSpec.sceneIdx],
		PrimaryTagID: tagIDs[markerSpec.primaryTagIdx],
	}

	err := mqb.Create(ctx, &marker)

	if err != nil {
		return fmt.Errorf("error creating marker %v+: %w", marker, err)
	}

	markerIDs = append(markerIDs, marker.ID)

	if len(markerSpec.tagIdxs) > 0 {
		newTagIDs := []int{}

		for _, tagIdx := range markerSpec.tagIdxs {
			newTagIDs = append(newTagIDs, tagIDs[tagIdx])
		}

		if err := mqb.UpdateTags(ctx, marker.ID, newTagIDs); err != nil {
			return fmt.Errorf("error creating marker/tag join: %w", err)
		}
	}

	return nil
}

func createChapter(ctx context.Context, mqb models.GalleryChapterReaderWriter, chapterSpec chapterSpec) error {
	chapter := models.GalleryChapter{
		GalleryID:  sceneIDs[chapterSpec.galleryIdx],
		Title:      chapterSpec.title,
		ImageIndex: chapterSpec.imageIndex,
	}

	err := mqb.Create(ctx, &chapter)

	if err != nil {
		return fmt.Errorf("error creating chapter %v+: %w", chapter, err)
	}

	chapterIDs = append(chapterIDs, chapter.ID)

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
		filterQ := ""
		filterPage := i
		filterPerPage := i * 40
		filterSort := "date"
		filterDirection := models.SortDirectionEnumAsc
		findFilter := models.FindFilterType{
			Q:         &filterQ,
			Page:      &filterPage,
			PerPage:   &filterPerPage,
			Sort:      &filterSort,
			Direction: &filterDirection,
		}
		savedFilter := models.SavedFilter{
			Mode:       getSavedFilterMode(i),
			Name:       getSavedFilterName(i),
			FindFilter: &findFilter,
			ObjectFilter: map[string]interface{}{
				"test": "object",
			},
			UIOptions: map[string]interface{}{
				"display_mode": 1,
				"zoom_index":   1,
			},
		}

		err := qb.Create(ctx, &savedFilter)

		if err != nil {
			return fmt.Errorf("Error creating saved filter %v+: %s", savedFilter, err.Error())
		}

		savedFilterIDs = append(savedFilterIDs, savedFilter.ID)
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

func linkMovieStudios(ctx context.Context, mqb models.MovieWriter) error {
	return doLinks(movieStudioLinks, func(movieIndex, studioIndex int) error {
		movie := models.MoviePartial{
			StudioID: models.NewOptionalInt(studioIDs[studioIndex]),
		}
		_, err := mqb.UpdatePartial(ctx, movieIDs[movieIndex], movie)

		return err
	})
}

func linkStudiosParent(ctx context.Context) error {
	qb := db.Studio
	return doLinks(studioParentLinks, func(parentIndex, childIndex int) error {
		input := &models.StudioPartial{
			ID:       studioIDs[childIndex],
			ParentID: models.NewOptionalInt(studioIDs[parentIndex]),
		}
		_, err := qb.UpdatePartial(ctx, *input)

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
	return qb.UpdateImage(ctx, tagIDs[tagIndex], []byte("image"))
}
