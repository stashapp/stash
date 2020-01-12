package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/stashapp/stash/pkg/models"
)

type Dvd struct {
	Name        string          `json:"name,omitempty"`
	Aliases     string          `json:"aliases,omitempty"`
	Durationdvd string          `json:"durationdvd,omitempty"`
	Year        string          `json:"year,omitempty"`
	Director    string          `json:"director,omitempty"`
	Synopsis    string          `json:"sypnopsis,omitempty"`
	FrontImage  string          `json:"frontimage,omitempty"`
	BackImage   string          `json:"backimage,omitempty"`
	URL         string          `json:"url,omitempty"`
	CreatedAt   models.JSONTime `json:"created_at,omitempty"`
	UpdatedAt   models.JSONTime `json:"updated_at,omitempty"`
}

func LoadDvdFile(filePath string) (*Dvd, error) {
	var dvd Dvd
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&dvd)
	if err != nil {
		return nil, err
	}
	return &dvd, nil
}

func SaveDvdFile(filePath string, dvd *Dvd) error {
	if dvd == nil {
		return fmt.Errorf("dvd must not be nil")
	}
	return marshalToFile(filePath, dvd)
}
