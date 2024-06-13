package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models/json"
)

type Movie struct {
	Name       string        `json:"name,omitempty"`
	Aliases    string        `json:"aliases,omitempty"`
	Duration   int           `json:"duration,omitempty"`
	Date       string        `json:"date,omitempty"`
	Rating     int           `json:"rating,omitempty"`
	Director   string        `json:"director,omitempty"`
	Synopsis   string        `json:"synopsis,omitempty"`
	FrontImage string        `json:"front_image,omitempty"`
	BackImage  string        `json:"back_image,omitempty"`
	URLs       []string      `json:"urls,omitempty"`
	Studio     string        `json:"studio,omitempty"`
	CreatedAt  json.JSONTime `json:"created_at,omitempty"`
	UpdatedAt  json.JSONTime `json:"updated_at,omitempty"`

	// deprecated - for import only
	URL string `json:"url,omitempty"`
}

func (s Movie) Filename() string {
	return fsutil.SanitiseBasename(s.Name) + ".json"
}

// Backwards Compatible synopsis for the movie
type MovieSynopsisBC struct {
	Synopsis string `json:"sypnopsis,omitempty"`
}

func LoadMovieFile(filePath string) (*Movie, error) {
	var movie Movie
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&movie)
	if err != nil {
		return nil, err
	}
	if movie.Synopsis == "" {
		// keep backwards compatibility with pre #2664 builds
		// attempt to get the synopsis from the alternate (sypnopsis) key

		_, err = file.Seek(0, 0) // seek to start of file
		if err == nil {
			var synopsis MovieSynopsisBC
			err = jsonParser.Decode(&synopsis)
			if err == nil {
				movie.Synopsis = synopsis.Synopsis
				if movie.Synopsis != "" {
					logger.Debug("Movie synopsis retrieved from alternate key")
				}
			}
		}
	}
	return &movie, nil
}

func SaveMovieFile(filePath string, movie *Movie) error {
	if movie == nil {
		return fmt.Errorf("movie must not be nil")
	}
	return marshalToFile(filePath, movie)
}
