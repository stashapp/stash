//go:build tools
// +build tools

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"
)

const batchSize = 50000

// create an example database by generating a number of scenes, markers,
// performers, studios, galleries, chapters and tags, and associating between them all

type config struct {
	Database   string       `yaml:"database"`
	Scenes     int          `yaml:"scenes"`
	Markers    int          `yaml:"markers"`
	Images     int          `yaml:"images"`
	Galleries  int          `yaml:"galleries"`
	Chapters   int          `yaml:"chapters"`
	Performers int          `yaml:"performers"`
	Studios    int          `yaml:"studios"`
	Tags       int          `yaml:"tags"`
	Naming     namingConfig `yaml:"naming"`
}

var (
	repo     models.Repository
	c        *config
	db       *sqlite.Database
	folderID file.FolderID
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	c, err = loadConfig()
	if err != nil {
		log.Fatalf("couldn't load configuration: %v", err)
	}

	initNaming(*c)

	db = sqlite.NewDatabase()
	repo = db.TxnRepository()

	logf("Initializing database...")
	if err = db.Open(c.Database); err != nil {
		log.Fatalf("couldn't initialize database: %v", err)
	}
	logf("Populating database...")
	populateDB()
}

func loadConfig() (*config, error) {
	ret := &config{}

	file, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	parser := yaml.NewDecoder(file)
	parser.SetStrict(true)
	err = parser.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func populateDB() {
	makeTags(c.Tags)
	makeStudios(c.Studios)
	makePerformers(c.Performers)
	makeScenes(c.Scenes)
	makeImages(c.Images)
	makeGalleries(c.Galleries)
	makeChapters(c.Chapters)
	makeMarkers(c.Markers)
}

func withTxn(f func(ctx context.Context) error) error {
	return txn.WithTxn(context.Background(), db, f)
}

func retry(attempts int, fn func() error) error {
	var err error
	for tries := 0; tries < attempts; tries++ {
		err = fn()
		if err == nil {
			return nil
		}
	}

	return err
}

func getOrCreateFolder(ctx context.Context, p string) (*file.Folder, error) {
	ret, err := repo.Folder.FindByPath(ctx, p)
	if err != nil {
		return nil, err
	}

	if ret != nil {
		return ret, nil
	}

	var parentID *file.FolderID

	if p != "." {
		parent := path.Dir(p)
		parentFolder, err := getOrCreateFolder(ctx, parent)
		if err != nil {
			return nil, err
		}

		parentID = &parentFolder.ID
	}

	f := file.Folder{
		Path:           p,
		ParentFolderID: parentID,
	}

	if err := repo.Folder.Create(ctx, &f); err != nil {
		return nil, err
	}

	ret = &f
	return ret, nil
}

func makeTags(n int) {
	logf("creating %d tags...", n)
	for i := 0; i < n; i++ {
		if err := retry(100, func() error {
			return withTxn(func(ctx context.Context) error {
				name := names[c.Naming.Tags].generateName(1)
				tag := models.Tag{
					Name: name,
				}

				created, err := repo.Tag.Create(ctx, tag)
				if err != nil {
					return err
				}

				if rand.Intn(100) > 5 {
					t, _, err := repo.Tag.Query(ctx, nil, getRandomFilter(1))
					if err != nil {
						return err
					}

					if len(t) > 0 && t[0].ID != created.ID {
						if err := repo.Tag.UpdateParentTags(ctx, created.ID, []int{t[0].ID}); err != nil {
							return err
						}
					}
				}

				return nil
			})
		}); err != nil {
			panic(err)
		}
	}
}

func makeStudios(n int) {
	logf("creating %d studios...", n)
	for i := 0; i < n; i++ {
		if err := retry(100, func() error {
			return withTxn(func(ctx context.Context) error {
				name := names[c.Naming.Tags].generateName(rand.Intn(5) + 1)
				studio := models.Studio{
					Name:     sql.NullString{String: name, Valid: true},
					Checksum: md5.FromString(name),
				}

				if rand.Intn(100) > 5 {
					ss, _, err := repo.Studio.Query(ctx, nil, getRandomFilter(1))
					if err != nil {
						return err
					}

					if len(ss) > 0 {
						studio.ParentID = sql.NullInt64{
							Int64: int64(ss[0].ID),
							Valid: true,
						}
					}
				}

				_, err := repo.Studio.Create(ctx, studio)
				return err
			})
		}); err != nil {
			panic(err)
		}
	}
}

func makePerformers(n int) {
	logf("creating %d performers...", n)
	for i := 0; i < n; i++ {
		if err := retry(100, func() error {
			return withTxn(func(ctx context.Context) error {
				name := generatePerformerName()
				performer := &models.Performer{
					Name:     name,
					Checksum: md5.FromString(name),
				}

				// TODO - set tags

				err := repo.Performer.Create(ctx, performer)
				if err != nil {
					err = fmt.Errorf("error creating performer with name: %s: %s", performer.Name, err.Error())
				}
				return err
			})
		}); err != nil {
			panic(err)
		}
	}
}

func generateBaseFile(parentFolderID file.FolderID, path string) *file.BaseFile {
	return &file.BaseFile{
		Basename:       path,
		ParentFolderID: parentFolderID,
		Fingerprints: []file.Fingerprint{
			file.Fingerprint{
				Type:        "md5",
				Fingerprint: md5.FromString(path),
			},
			file.Fingerprint{
				Type:        "oshash",
				Fingerprint: md5.FromString(path),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func generateVideoFile(parentFolderID file.FolderID, path string) file.File {
	w, h := getResolution()

	return &file.VideoFile{
		BaseFile: generateBaseFile(parentFolderID, path),
		Duration: rand.Float64() * 14400,
		Height:   h,
		Width:    w,
	}
}

func makeVideoFile(ctx context.Context, path string) (file.File, error) {
	folderPath := fsutil.GetIntraDir(path, 2, 2)
	parentFolder, err := getOrCreateFolder(ctx, folderPath)
	if err != nil {
		return nil, err
	}

	f := generateVideoFile(parentFolder.ID, path)

	if err := repo.File.Create(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

func logf(f string, args ...interface{}) {
	log.Printf(f+"\n", args...)
}

func makeScenes(n int) {
	logf("creating %d scenes...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize

		if err := withTxn(func(ctx context.Context) error {
			for ; i < batch && i < n; i++ {
				scene := generateScene(i)
				scene.StudioID = getRandomStudioID(ctx)
				makeSceneRelationships(ctx, &scene)

				path := md5.FromString("scene/" + strconv.Itoa(i))
				f, err := makeVideoFile(ctx, path)
				if err != nil {
					return err
				}

				if err := repo.Scene.Create(ctx, &scene, []file.ID{f.Base().ID}); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			panic(err)
		}

		logf("... created %d scenes", i)
	}
}

func getResolution() (int, int) {
	res := models.AllResolutionEnum[rand.Intn(len(models.AllResolutionEnum))]
	h := res.GetMaxResolution()
	var w int
	if h == 240 || h == 480 || rand.Intn(10) == 9 {
		w = h * 4 / 3
	} else {
		w = h * 16 / 9
	}

	if rand.Intn(10) == 9 {
		return h, w
	}

	return w, h
}

func getBool() {
	return rand.Intn(2) == 0
}

func getDate() time.Time {
	s := rand.Int63n(time.Now().Unix())

	return time.Unix(s, 0)
}

func generateScene(i int) models.Scene {
	return models.Scene{
		Title: names[c.Naming.Scenes].generateName(rand.Intn(7) + 1),
		Date: &models.Date{
			Time: getDate(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func generateImageFile(parentFolderID file.FolderID, path string) file.File {
	w, h := getResolution()

	return &file.ImageFile{
		BaseFile: generateBaseFile(parentFolderID, path),
		Height:   h,
		Width:    w,
		Clip:     getBool(),
	}
}

func makeImageFile(ctx context.Context, path string) (file.File, error) {
	folderPath := fsutil.GetIntraDir(path, 2, 2)
	parentFolder, err := getOrCreateFolder(ctx, folderPath)
	if err != nil {
		return nil, err
	}

	f := generateImageFile(parentFolder.ID, path)

	if err := repo.File.Create(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

func makeImages(n int) {
	logf("creating %d images...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize
		if err := withTxn(func(ctx context.Context) error {
			for ; i < batch && i < n; i++ {
				image := generateImage(i)
				image.StudioID = getRandomStudioID(ctx)
				makeImageRelationships(ctx, &image)

				path := md5.FromString("image/" + strconv.Itoa(i))
				f, err := makeImageFile(ctx, path)
				if err != nil {
					return err
				}

				if err := repo.Image.Create(ctx, &models.ImageCreateInput{
					Image:   &image,
					FileIDs: []file.ID{f.Base().ID},
				}); err != nil {
					return err
				}
			}

			logf("... created %d images", i)

			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func generateImage(i int) models.Image {
	return models.Image{
		Title:     names[c.Naming.Images].generateName(rand.Intn(7) + 1),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func makeGalleries(n int) {
	logf("creating %d galleries...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize

		if err := withTxn(func(ctx context.Context) error {
			for ; i < batch && i < n; i++ {
				gallery := generateGallery(i)
				gallery.StudioID = getRandomStudioID(ctx)
				gallery.TagIDs = models.NewRelatedIDs(getRandomTags(ctx, 0, 15))
				gallery.PerformerIDs = models.NewRelatedIDs(getRandomPerformers(ctx))

				path := md5.FromString("gallery/" + strconv.Itoa(i))
				f, err := makeZipFile(ctx, path)
				if err != nil {
					return err
				}

				if err := repo.Gallery.Create(ctx, &gallery, []file.ID{f.Base().ID}); err != nil {
					return err
				}

				makeGalleryRelationships(ctx, &gallery)
			}

			return nil
		}); err != nil {
			panic(err)
		}

		logf("... created %d galleries", i)
	}
}

func generateZipFile(parentFolderID file.FolderID, path string) file.File {
	return generateBaseFile(parentFolderID, path)
}

func makeZipFile(ctx context.Context, path string) (file.File, error) {
	folderPath := fsutil.GetIntraDir(path, 2, 2)
	parentFolder, err := getOrCreateFolder(ctx, folderPath)
	if err != nil {
		return nil, err
	}

	f := generateZipFile(parentFolder.ID, path)

	if err := repo.File.Create(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

func generateGallery(i int) models.Gallery {
	return models.Gallery{
		Title: names[c.Naming.Galleries].generateName(rand.Intn(7) + 1),
		Date: &models.Date{
			Time: getDate(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func makeChapters(n int) {
	logf("creating %d chapters...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize
		if err := withTxn(func(ctx context.Context) error {
			for ; i < batch && i < n; i++ {
				chapter := generateChapter(i)
				chapter.GalleryID = models.NullInt64(int64(getRandomGallery()))

				created, err := repo.GalleryChapter.Create(ctx, chapter)
				if err != nil {
					return err
				}
			}

			logf("... created %d chapters", i)

			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func generateChapter(i int) models.GalleryChapter {
	return models.GalleryChapter{
		Title:      names[c.Naming.Galleries].generateName(rand.Intn(7) + 1),
		ImageIndex: rand.Intn(200),
	}
}

func makeMarkers(n int) {
	logf("creating %d markers...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize
		if err := withTxn(func(ctx context.Context) error {
			for ; i < batch && i < n; i++ {
				marker := generateMarker(i)
				marker.SceneID = models.NullInt64(int64(getRandomScene()))
				marker.PrimaryTagID = getRandomTags(ctx, 1, 1)[0]

				created, err := repo.SceneMarker.Create(ctx, marker)
				if err != nil {
					return err
				}

				tags := getRandomTags(ctx, 0, 5)
				// remove primary tag
				tags = sliceutil.Exclude(tags, []int{marker.PrimaryTagID})
				if err := repo.SceneMarker.UpdateTags(ctx, created.ID, tags); err != nil {
					return err
				}
			}

			logf("... created %d markers", i)

			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func generateMarker(i int) models.SceneMarker {
	return models.SceneMarker{
		Title: names[c.Naming.Scenes].generateName(rand.Intn(7) + 1),
	}
}

func getRandomFilter(n int) *models.FindFilterType {
	seed := math.Floor(rand.Float64() * math.Pow10(8))
	sortBy := fmt.Sprintf("random_%.f", seed)
	return &models.FindFilterType{
		Sort:    &sortBy,
		PerPage: &n,
	}
}

func getRandomStudioID(ctx context.Context) *int {
	if rand.Intn(10) == 0 {
		return nil
	}

	// s, _, err := r.Studio().Query(nil, getRandomFilter(1))
	// if err != nil {
	// 	panic(err)
	// }

	v := rand.Intn(c.Studios) + 1
	return &v
}

func makeSceneRelationships(ctx context.Context, s *models.Scene) {
	// add tags
	s.TagIDs = models.NewRelatedIDs(getRandomTags(ctx, 0, 15))

	// add performers
	s.PerformerIDs = models.NewRelatedIDs(getRandomPerformers(ctx))
}

func makeImageRelationships(ctx context.Context, i *models.Image) {
	// there are typically many more images. For performance reasons
	// only a small proportion should have tags/performers

	// add tags
	if rand.Intn(100) == 0 {
		i.TagIDs = models.NewRelatedIDs(getRandomTags(ctx, 1, 15))
	}

	// add performers
	if rand.Intn(100) <= 1 {
		i.PerformerIDs = models.NewRelatedIDs(getRandomPerformers(ctx))
	}
}

func makeGalleryRelationships(ctx context.Context, g *models.Gallery) {
	// add images
	imageIDs := getRandomImages(ctx)
	if len(imageIDs) > 0 {
		if err := repo.Gallery.UpdateImages(ctx, g.ID, imageIDs); err != nil {
			panic(err)
		}
	}
}

func getRandomPerformers(ctx context.Context) []int {
	n := rand.Intn(5)

	var ret []int
	// if n > 0 {
	// 	p, _, err := r.Performer().Query(nil, getRandomFilter(n))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, pp := range p {
	// 		ret = sliceutil.AppendUnique(ret, pp.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = sliceutil.AppendUnique(ret, rand.Intn(c.Performers)+1)
	}

	return ret
}

func getRandomScene() int {
	return rand.Intn(c.Scenes) + 1
}

func getRandomGallery() int {
	return rand.Intn(c.Galleries) + 1
}

func getRandomTags(ctx context.Context, min, max int) []int {
	var n int
	if min == max {
		n = min
	} else {
		n = rand.Intn(max-min) + min
	}

	var ret []int
	// if n > 0 {
	// 	t, _, err := r.Tag().Query(nil, getRandomFilter(n))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, tt := range t {
	// 		ret = sliceutil.AppendUnique(ret, tt.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = sliceutil.AppendUnique(ret, rand.Intn(c.Tags)+1)
	}

	return ret
}

func getRandomImages(ctx context.Context) []int {
	n := rand.Intn(500)

	var ret []int
	// if n > 0 {
	// 	t, _, err := r.Image().Query(nil, getRandomFilter(n))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, tt := range t {
	// 		ret = sliceutil.AppendUnique(ret, tt.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = sliceutil.AppendUnique(ret, rand.Intn(c.Images)+1)
	}

	return ret
}
