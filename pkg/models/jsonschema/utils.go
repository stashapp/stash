package jsonschema

import (
	"bytes"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func CompareJSON(a interface{}, b interface{}) bool {
	aBuf, _ := encode(a)
	bBuf, _ := encode(b)
	return bytes.Equal(aBuf, bBuf)
}

func marshalToFile(filePath string, j interface{}) error {
	data, err := encode(j)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func encode(j interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(j); err != nil {
		return nil, err
	}
	// Strip the newline at the end of the file
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}
