//go:build ignore
// +build ignore

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/sqlite"
	"gopkg.in/yaml.v2"
)

const batchSize = 50000

// create an example database by generating a number of scenes, markers,
// performers, studios and tags, and associating between them all

type config struct {
	Database   string       `yaml:"database"`
	Scenes     int          `yaml:"scenes"`
	Markers    int          `yaml:"markers"`
	Images     int          `yaml:"images"`
	Galleries  int          `yaml:"galleries"`
	Performers int          `yaml:"performers"`
	Studios    int          `yaml:"studios"`
	Tags       int          `yaml:"tags"`
	Naming     namingConfig `yaml:"naming"`
}

var txnManager models.TransactionManager
var c *config

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	c, err = loadConfig()
	if err != nil {
		log.Fatalf("couldn't load configuration: %v", err)
	}

	initNaming(*c)

	if err = database.Initialize(c.Database); err != nil {
		log.Fatalf("couldn't initialize database: %v", err)
	}
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
	makeMarkers(c.Markers)
}

func withTxn(f func(r models.Repository) error) error {
	if txnManager == nil {
		txnManager = sqlite.NewTransactionManager()
	}

	return txnManager.WithTxn(context.TODO(), f)
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

func makeTags(n int) {
	for i := 0; i < n; i++ {
		if err := retry(100, func() error {
			return withTxn(func(r models.Repository) error {
				name := names[c.Naming.Tags].generateName(1)
				tag := models.Tag{
					Name: name,
				}

				created, err := r.Tag().Create(tag)
				if err != nil {
					return err
				}

				if rand.Intn(100) > 5 {
					t, _, err := r.Tag().Query(nil, getRandomFilter(1))
					if err != nil {
						return err
					}

					if len(t) > 0 && t[0].ID != created.ID {
						if err := r.Tag().UpdateParentTags(created.ID, []int{t[0].ID}); err != nil {
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
	for i := 0; i < n; i++ {
		if err := retry(100, func() error {
			return withTxn(func(r models.Repository) error {
				name := names[c.Naming.Tags].generateName(rand.Intn(5) + 1)
				studio := models.Studio{
					Name:     sql.NullString{String: name, Valid: true},
					Checksum: md5.FromString(name),
				}

				if rand.Intn(100) > 5 {
					ss, _, err := r.Studio().Query(nil, getRandomFilter(1))
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

				_, err := r.Studio().Create(studio)
				return err
			})
		}); err != nil {
			panic(err)
		}
	}
}

func makePerformers(n int) {
	for i := 0; i < n; i++ {
		if err := retry(100, func() error {
			return withTxn(func(r models.Repository) error {
				name := generatePerformerName()
				performer := models.Performer{
					Name:     sql.NullString{String: name, Valid: true},
					Checksum: md5.FromString(name),
					Favorite: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				}

				// TODO - set tags

				_, err := r.Performer().Create(performer)
				if err != nil {
					err = fmt.Errorf("error creating performer with name: %s: %s", performer.Name.String, err.Error())
				}
				return err
			})
		}); err != nil {
			panic(err)
		}
	}
}

func makeScenes(n int) {
	logger.Infof("creating %d scenes...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize

		if err := withTxn(func(r models.Repository) error {
			for ; i < batch && i < n; i++ {
				scene := generateScene(i)
				scene.StudioID = getRandomStudioID(r)

				created, err := r.Scene().Create(scene)
				if err != nil {
					return err
				}

				makeSceneRelationships(r, created.ID)
			}

			return nil
		}); err != nil {
			panic(err)
		}

		logger.Infof("... created %d scenes", i)
	}
}

func getResolution() (int64, int64) {
	res := models.AllResolutionEnum[rand.Intn(len(models.AllResolutionEnum))]
	h := int64(res.GetMaxResolution())
	var w int64
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

func getDate() string {
	s := rand.Int63n(time.Now().Unix())

	d := time.Unix(s, 0)
	return d.Format("2006-01-02")
}

func generateScene(i int) models.Scene {
	path := md5.FromString("scene/" + strconv.Itoa(i))
	w, h := getResolution()

	return models.Scene{
		Path:     path,
		Title:    sql.NullString{String: names[c.Naming.Scenes].generateName(rand.Intn(7) + 1), Valid: true},
		Checksum: sql.NullString{String: md5.FromString(path), Valid: true},
		OSHash:   sql.NullString{String: md5.FromString(path), Valid: true},
		Duration: sql.NullFloat64{
			Float64: rand.Float64() * 14400,
			Valid:   true,
		},
		Height: models.NullInt64(h),
		Width:  models.NullInt64(w),
		Date: models.SQLiteDate{
			String: getDate(),
			Valid:  true,
		},
	}
}

func makeImages(n int) {
	logger.Infof("creating %d images...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize
		if err := withTxn(func(r models.Repository) error {
			for ; i < batch && i < n; i++ {
				image := generateImage(i)
				image.StudioID = getRandomStudioID(r)

				created, err := r.Image().Create(image)
				if err != nil {
					return err
				}

				makeImageRelationships(r, created.ID)
			}

			logger.Infof("... created %d images", i)

			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func generateImage(i int) models.Image {
	path := md5.FromString("image/" + strconv.Itoa(i))

	w, h := getResolution()

	return models.Image{
		Title:    sql.NullString{String: names[c.Naming.Images].generateName(rand.Intn(7) + 1), Valid: true},
		Path:     path,
		Checksum: md5.FromString(path),
		Height:   models.NullInt64(h),
		Width:    models.NullInt64(w),
	}
}

func makeGalleries(n int) {
	logger.Infof("creating %d galleries...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize

		if err := withTxn(func(r models.Repository) error {
			for ; i < batch && i < n; i++ {
				gallery := generateGallery(i)
				gallery.StudioID = getRandomStudioID(r)

				created, err := r.Gallery().Create(gallery)
				if err != nil {
					return err
				}

				makeGalleryRelationships(r, created.ID)
			}

			return nil
		}); err != nil {
			panic(err)
		}

		logger.Infof("... created %d galleries", i)
	}
}

func generateGallery(i int) models.Gallery {
	path := md5.FromString("gallery/" + strconv.Itoa(i))

	return models.Gallery{
		Title:    sql.NullString{String: names[c.Naming.Galleries].generateName(rand.Intn(7) + 1), Valid: true},
		Path:     sql.NullString{String: path, Valid: true},
		Checksum: md5.FromString(path),
		Date: models.SQLiteDate{
			String: getDate(),
			Valid:  true,
		},
	}
}

func makeMarkers(n int) {
	logger.Infof("creating %d markers...", n)
	for i := 0; i < n; {
		// do in batches of 1000
		batch := i + batchSize
		if err := withTxn(func(r models.Repository) error {
			for ; i < batch && i < n; i++ {
				marker := generateMarker(i)
				marker.SceneID = models.NullInt64(int64(getRandomScene()))
				marker.PrimaryTagID = getRandomTags(r, 1, 1)[0]

				created, err := r.SceneMarker().Create(marker)
				if err != nil {
					return err
				}

				tags := getRandomTags(r, 0, 5)
				// remove primary tag
				tags = intslice.IntExclude(tags, []int{marker.PrimaryTagID})
				if err := r.SceneMarker().UpdateTags(created.ID, tags); err != nil {
					return err
				}
			}

			logger.Infof("... created %d markers", i)

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

func getRandomStudioID(r models.Repository) sql.NullInt64 {
	if rand.Intn(10) == 0 {
		return sql.NullInt64{}
	}

	// s, _, err := r.Studio().Query(nil, getRandomFilter(1))
	// if err != nil {
	// 	panic(err)
	// }

	return sql.NullInt64{
		Int64: int64(rand.Int63n(int64(c.Studios)) + 1),
		Valid: true,
	}
}

func makeSceneRelationships(r models.Repository, id int) {
	// add tags
	tagIDs := getRandomTags(r, 0, 15)
	if len(tagIDs) > 0 {
		if err := r.Scene().UpdateTags(id, tagIDs); err != nil {
			panic(err)
		}
	}

	// add performers
	performerIDs := getRandomPerformers(r)
	if len(tagIDs) > 0 {
		if err := r.Scene().UpdatePerformers(id, performerIDs); err != nil {
			panic(err)
		}
	}
}

func makeImageRelationships(r models.Repository, id int) {
	// there are typically many more images. For performance reasons
	// only a small proportion should have tags/performers

	// add tags
	if rand.Intn(100) == 0 {
		tagIDs := getRandomTags(r, 1, 15)
		if len(tagIDs) > 0 {
			if err := r.Image().UpdateTags(id, tagIDs); err != nil {
				panic(err)
			}
		}
	}

	// add performers
	if rand.Intn(100) <= 1 {
		performerIDs := getRandomPerformers(r)
		if len(performerIDs) > 0 {
			if err := r.Image().UpdatePerformers(id, performerIDs); err != nil {
				panic(err)
			}
		}
	}
}

func makeGalleryRelationships(r models.Repository, id int) {
	// add tags
	tagIDs := getRandomTags(r, 0, 15)
	if len(tagIDs) > 0 {
		if err := r.Gallery().UpdateTags(id, tagIDs); err != nil {
			panic(err)
		}
	}

	// add performers
	performerIDs := getRandomPerformers(r)
	if len(tagIDs) > 0 {
		if err := r.Gallery().UpdatePerformers(id, performerIDs); err != nil {
			panic(err)
		}
	}

	// add images
	imageIDs := getRandomImages(r)
	if len(tagIDs) > 0 {
		if err := r.Gallery().UpdateImages(id, imageIDs); err != nil {
			panic(err)
		}
	}
}

func getRandomPerformers(r models.Repository) []int {
	n := rand.Intn(5)

	var ret []int
	// if n > 0 {
	// 	p, _, err := r.Performer().Query(nil, getRandomFilter(n))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, pp := range p {
	// 		ret = intslice.IntAppendUnique(ret, pp.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = intslice.IntAppendUnique(ret, rand.Intn(c.Performers)+1)
	}

	return ret
}

func getRandomScene() int {
	return rand.Intn(c.Scenes) + 1
}

func getRandomTags(r models.Repository, min, max int) []int {
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
	// 		ret = intslice.IntAppendUnique(ret, tt.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = intslice.IntAppendUnique(ret, rand.Intn(c.Tags)+1)
	}

	return ret
}

func getRandomImages(r models.Repository) []int {
	n := rand.Intn(500)

	var ret []int
	// if n > 0 {
	// 	t, _, err := r.Image().Query(nil, getRandomFilter(n))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, tt := range t {
	// 		ret = intslice.IntAppendUnique(ret, tt.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = intslice.IntAppendUnique(ret, rand.Intn(c.Images)+1)
	}

	return ret
}
