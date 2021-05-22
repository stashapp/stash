package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type PathNameMapping struct {
	Path     string `json:"path,omitempty"`
	Name     string `json:"name,omitempty"`
	Checksum string `json:"checksum"`
}

type Mappings struct {
	Tags       []PathNameMapping `json:"tags"`
	Performers []PathNameMapping `json:"performers"`
	Studios    []PathNameMapping `json:"studios"`
	Movies     []PathNameMapping `json:"movies"`
	Galleries  []PathNameMapping `json:"galleries"`
	Scenes     []PathNameMapping `json:"scenes"`
	Images     []PathNameMapping `json:"images"`
}

func LoadMappingsFile(filePath string) (*Mappings, error) {
	var mappings Mappings
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&mappings)
	if err != nil {
		return nil, err
	}
	return &mappings, nil
}

func SaveMappingsFile(filePath string, mappings *Mappings) error {
	if mappings == nil {
		return fmt.Errorf("mappings must not be nil")
	}
	return marshalToFile(filePath, mappings)
}
