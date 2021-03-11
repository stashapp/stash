// +build ignore

package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/database"
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
	performers int
	studios    int
	tags       int
}

func main() {
	f, err := ioutil.TempFile(".", "*.sqlite")
	if err != nil {
		panic(fmt.Sprintf("Could not create temporary file: %s", err.Error()))
	}

	f.Close()
	databaseFile := f.Name()
	database.Initialize(databaseFile)

	populateDB(config{
		scenes: 30000,
		images: 150000,
		// galleries: 1500,
		markers:    300,
		performers: 10000,
		studios:    500,
		tags:       1500,
	})
}

func populateDB(c config) {
	makeTags(c.tags)
	makeStudios(c.studios)
	makePerformers(c.performers)
	makeScenes(c.scenes)
}

func withTxn(f func(r models.Repository) error) error {
	t := sqlite.NewTransactionManager()
	return t.WithTxn(context.TODO(), f)
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
	if err := withTxn(func(r models.Repository) error {
		for i := 0; i < n; i++ {
			scene := generateScene(i)

			// TODO - set tags, performers, studio etc

			if _, err := r.Scene().Create(scene); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		panic(err)
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
	rand.Seed(533)

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
