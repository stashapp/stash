package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/stashapp/stash/pkg/models"
)

type Movie struct {
	Name       string          `json:"name,omitempty"`
	Aliases    string          `json:"aliases,omitempty"`
	Duration   int             `json:"duration,omitempty"`
	Date       string          `json:"date,omitempty"`
	Rating     int             `json:"rating,omitempty"`
	Director   string          `json:"director,omitempty"`
	Synopsis   string          `json:"sypnopsis,omitempty"`
	FrontImage string          `json:"front_image,omitempty"`
	BackImage  string          `json:"back_image,omitempty"`
	URL        string          `json:"url,omitempty"`
	Studio     string          `json:"studio,omitempty"`
	CreatedAt  models.JSONTime `json:"created_at,omitempty"`
	UpdatedAt  models.JSONTime `json:"updated_at,omitempty"`
}

func LoadMovieFile(filePath string) (*Movie, error) {
	var movie Movie
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func SaveMovieFile(filePath string, movie *Movie) error {
	if movie == nil {
		return fmt.Errorf("movie must not be nil")
	}
	return marshalToFile(filePath, movie)
}
