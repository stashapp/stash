package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/stashapp/stash/pkg/models"
)

type Movie struct {
	Name           string          `json:"name,omitempty"`
	Aliases        string          `json:"aliases,omitempty"`
	Duration	   string          `json:"duration,omitempty"`
	Date    	   string          `json:"date,omitempty"`
	Rating   	   string          `json:"rating,omitempty"`
	Director       string          `json:"director,omitempty"`
	Synopsis       string          `json:"sypnopsis,omitempty"`
	Front_Image    string          `json:"front_image,omitempty"`
	Back_Image     string          `json:"back_image,omitempty"`
	URL            string          `json:"url,omitempty"`
	CreatedAt      models.JSONTime `json:"created_at,omitempty"`
	UpdatedAt      models.JSONTime `json:"updated_at,omitempty"`
}

func LoadMovieFile(filePath string) (*Movie, error) {
	var movie Movie
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
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
