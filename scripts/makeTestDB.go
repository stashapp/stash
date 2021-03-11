// uild ignore

package main

import (
	"context"
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

// create an example database by generating a number of scenes, markers,
// performers, studios and tags, and associating between them all

type config struct {
	scenes     int
	markers    int
	images     int
	galleries  int
	performers int
	studios    int
	tags       int
}

var txnManager models.TransactionManager
var c config

func main() {
	database.Initialize("generated.sqlite")

	c = config{
		scenes:     30000,
		images:     150000,
		galleries:  1500,
		markers:    300,
		performers: 10000,
		studios:    500,
		tags:       1500,
	}

	populateDB()
}

func populateDB() {
	makeTags(c.tags)
	makeStudios(c.studios)
	makePerformers(c.performers)
	makeScenes(c.scenes)
	makeImages(c.images)
	makeGalleries(c.galleries)
}

func withTxn(f func(r models.Repository) error) error {
	if txnManager == nil {
		txnManager = sqlite.NewTransactionManager()
	}

	return txnManager.WithTxn(context.TODO(), f)
}

func makeTags(n int) {
	if err := withTxn(func(r models.Repository) error {
		for i := 0; i < n; i++ {
			name := utils.MD5FromString("tag/" + strconv.Itoa(i))
			tag := models.Tag{
				Name: name,
			}

			if _, err := r.Tag().Create(tag); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func makeStudios(n int) {
	if err := withTxn(func(r models.Repository) error {
		for i := 0; i < n; i++ {
			name := utils.MD5FromString("studio/" + strconv.Itoa(i))
			studio := models.Studio{
				Name:     sql.NullString{String: name, Valid: true},
				Checksum: utils.MD5FromString(name),
			}

			if _, err := r.Studio().Create(studio); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func makePerformers(n int) {
	if err := withTxn(func(r models.Repository) error {
		for i := 0; i < n; i++ {
			name := utils.MD5FromString("performer/" + strconv.Itoa(i))
			performer := models.Performer{
				Name:     sql.NullString{String: name, Valid: true},
				Checksum: utils.MD5FromString(name),
				Favorite: sql.NullBool{
					Bool:  false,
					Valid: true,
				},
			}

			// TODO - set tags

			if _, err := r.Performer().Create(performer); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func makeScenes(n int) {
	logger.Infof("creating %d scenes...", n)
	rand.Seed(533)
	for i := 0; i < n; i++ {
		if i > 0 && i%100 == 0 {
			logger.Infof("... created %d scenes", i)
		}

		if err := withTxn(func(r models.Repository) error {
			scene := generateScene(i)
			scene.StudioID = getRandomStudioID(r)

			created, err := r.Scene().Create(scene)
			if err != nil {
				return err
			}

			makeSceneRelationships(r, created.ID)
			return nil
		}); err != nil {
			panic(err)
		}
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
	path := utils.MD5FromString("scene/" + strconv.Itoa(i))
	w, h := getResolution()

	return models.Scene{
		Path:     path,
		Checksum: sql.NullString{String: utils.MD5FromString(path), Valid: true},
		OSHash:   sql.NullString{String: utils.MD5FromString(path), Valid: true},
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
	rand.Seed(1293)
	for i := 0; i < n; i++ {
		if i > 0 && i%100 == 0 {
			logger.Infof("... created %d images", i)
		}
		if err := withTxn(func(r models.Repository) error {
			image := generateImage(i)
			image.StudioID = getRandomStudioID(r)

			created, err := r.Image().Create(image)
			if err != nil {
				return err
			}

			makeImageRelationships(r, created.ID)
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func generateImage(i int) models.Image {
	path := utils.MD5FromString("image/" + strconv.Itoa(i))

	w, h := getResolution()

	return models.Image{
		Path:     path,
		Checksum: utils.MD5FromString(path),
		Height:   models.NullInt64(h),
		Width:    models.NullInt64(w),
	}
}

func makeGalleries(n int) {
	logger.Infof("creating %d galleries...", n)
	rand.Seed(92113)
	for i := 0; i < n; i++ {
		if i > 0 && i%100 == 0 {
			logger.Infof("... created %d galleries", i)
		}

		if err := withTxn(func(r models.Repository) error {
			gallery := generateGallery(i)
			gallery.StudioID = getRandomStudioID(r)

			created, err := r.Gallery().Create(gallery)
			if err != nil {
				return err
			}

			makeGalleryRelationships(r, created.ID)
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func generateGallery(i int) models.Gallery {
	path := utils.MD5FromString("gallery/" + strconv.Itoa(i))

	return models.Gallery{
		Path:     sql.NullString{String: path, Valid: true},
		Checksum: utils.MD5FromString(path),
		Date: models.SQLiteDate{
			String: getDate(),
			Valid:  true,
		},
	}
}

func getRandomFilter(n int) *models.FindFilterType {
	sortBy := "random"
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
		Int64: int64(rand.Int63n(int64(c.studios)) + 1),
		Valid: true,
	}
}

func makeSceneRelationships(r models.Repository, id int) {
	// add tags
	tagIDs := getRandomTags(r)
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
	// add tags
	tagIDs := getRandomTags(r)
	if len(tagIDs) > 0 {
		if err := r.Image().UpdateTags(id, tagIDs); err != nil {
			panic(err)
		}
	}

	// add performers
	performerIDs := getRandomPerformers(r)
	if len(tagIDs) > 0 {
		if err := r.Image().UpdatePerformers(id, performerIDs); err != nil {
			panic(err)
		}
	}
}

func makeGalleryRelationships(r models.Repository, id int) {
	// add tags
	tagIDs := getRandomTags(r)
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
	// 		ret = utils.IntAppendUnique(ret, pp.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = utils.IntAppendUnique(ret, rand.Intn(c.performers)+1)
	}

	return ret
}

func getRandomTags(r models.Repository) []int {
	n := rand.Intn(15)

	var ret []int
	// if n > 0 {
	// 	t, _, err := r.Tag().Query(nil, getRandomFilter(n))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, tt := range t {
	// 		ret = utils.IntAppendUnique(ret, tt.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = utils.IntAppendUnique(ret, rand.Intn(c.tags)+1)
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
	// 		ret = utils.IntAppendUnique(ret, tt.ID)
	// 	}
	// }

	for i := 0; i < n; i++ {
		ret = utils.IntAppendUnique(ret, rand.Intn(c.images)+1)
	}

	return ret
}
