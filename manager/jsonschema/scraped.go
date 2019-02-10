package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"
)

type ScrapedItem struct {
	Title           string    `json:"title,omitempty"`
	Description     string    `json:"description,omitempty"`
	Url             string    `json:"url,omitempty"`
	Date            string    `json:"date,omitempty"`
	Rating          string    `json:"rating,omitempty"`
	Tags            string    `json:"tags,omitempty"`
	Models          string    `json:"models,omitempty"`
	Episode         int       `json:"episode,omitempty"`
	GalleryFilename string    `json:"gallery_filename,omitempty"`
	GalleryUrl      string    `json:"gallery_url,omitempty"`
	VideoFilename   string    `json:"video_filename,omitempty"`
	VideoUrl        string    `json:"video_url,omitempty"`
	Studio          string    `json:"studio,omitempty"`
	UpdatedAt       RailsTime `json:"updated_at,omitempty"`
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