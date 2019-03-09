package jsonschema

import (
	"encoding/json"
	"fmt"
	"github.com/stashapp/stash/pkg/models"
	"os"
)

type ScrapedItem struct {
	Title           string          `json:"title,omitempty"`
	Description     string          `json:"description,omitempty"`
	URL             string          `json:"url,omitempty"`
	Date            string          `json:"date,omitempty"`
	Rating          string          `json:"rating,omitempty"`
	Tags            string          `json:"tags,omitempty"`
	Models          string          `json:"models,omitempty"`
	Episode         int             `json:"episode,omitempty"`
	GalleryFilename string          `json:"gallery_filename,omitempty"`
	GalleryURL      string          `json:"gallery_url,omitempty"`
	VideoFilename   string          `json:"video_filename,omitempty"`
	VideoURL        string          `json:"video_url,omitempty"`
	Studio          string          `json:"studio,omitempty"`
	UpdatedAt       models.JSONTime `json:"updated_at,omitempty"`
}

func LoadScrapedFile(filePath string) ([]ScrapedItem, error) {
	var scraped []ScrapedItem
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&scraped)
	if err != nil {
		return nil, err
	}
	return scraped, nil
}

func SaveScrapedFile(filePath string, scrapedItems []ScrapedItem) error {
	if scrapedItems == nil {
		return fmt.Errorf("scraped items must not be nil")
	}
	return marshalToFile(filePath, scrapedItems)
}
