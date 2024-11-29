package jsonschema

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func loadFile[T any](filePath string) (*T, error) {
	var ret T
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func saveFile[T any](filePath string, obj *T) error {
	if obj == nil {
		return fmt.Errorf("object must not be nil")
	}
	return marshalToFile(filePath, obj)
}
