package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type NameMapping struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
}

type PathMapping struct {
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
}

type Mappings struct {
	Tags       []NameMapping `json:"tags"`
	Performers []NameMapping `json:"performers"`
	Studios    []NameMapping `json:"studios"`
	Movies     []NameMapping `json:"movies"`
	Galleries  []PathMapping `json:"galleries"`
	Scenes     []PathMapping `json:"scenes"`
}

func LoadMappingsFile(filePath string) (*Mappings, error) {
	var mappings Mappings
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
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
