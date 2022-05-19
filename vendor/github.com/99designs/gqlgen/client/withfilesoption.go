package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

type fileFormDataMap struct {
	mapKey string
	file   *os.File
}

func findFiles(parentMapKey string, variables map[string]interface{}) []*fileFormDataMap {
	files := []*fileFormDataMap{}
	for key, value := range variables {
		if v, ok := value.(map[string]interface{}); ok {
			files = append(files, findFiles(parentMapKey+"."+key, v)...)
		} else if v, ok := value.([]map[string]interface{}); ok {
			for i, arr := range v {
				files = append(files, findFiles(fmt.Sprintf(`%s.%s.%d`, parentMapKey, key, i), arr)...)
			}
		} else if v, ok := value.([]*os.File); ok {
			for i, file := range v {
				files = append(files, &fileFormDataMap{
					mapKey: fmt.Sprintf(`%s.%s.%d`, parentMapKey, key, i),
					file:   file,
				})
			}
		} else if v, ok := value.(*os.File); ok {
			files = append(files, &fileFormDataMap{
				mapKey: parentMapKey + "." + key,
				file:   v,
			})
		}
	}

	return files
}

// WithFiles encodes the outgoing request body as multipart form data for file variables
func WithFiles() Option {
	return func(bd *Request) {
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)

		//-b7955bd2e1d17b67ac157b9e9ddb6238888caefc6f3541920a1debad284d
		// Content-Disposition: form-data; name="operations"
		//
		// {"query":"mutation ($input: Input!) {}","variables":{"input":{"file":{}}}
		requestBody, _ := json.Marshal(bd)
		bodyWriter.WriteField("operations", string(requestBody))

		// --b7955bd2e1d17b67ac157b9e9ddb6238888caefc6f3541920a1debad284d
		// Content-Disposition: form-data; name="map"
		//
		// `{ "0":["variables.input.file"] }`
		// or
		// `{ "0":["variables.input.files.0"], "1":["variables.input.files.1"] }`
		// or
		// `{ "0": ["variables.input.0.file"], "1": ["variables.input.1.file"] }`
		// or
		// `{ "0": ["variables.req.0.file", "variables.req.1.file"] }`
		mapData := ""
		filesData := findFiles("variables", bd.Variables)
		filesGroup := [][]*fileFormDataMap{}
		for _, fd := range filesData {
			foundDuplicate := false
			for j, fg := range filesGroup {
				f1, _ := fd.file.Stat()
				f2, _ := fg[0].file.Stat()
				if os.SameFile(f1, f2) {
					foundDuplicate = true
					filesGroup[j] = append(filesGroup[j], fd)
				}
			}

			if !foundDuplicate {
				filesGroup = append(filesGroup, []*fileFormDataMap{fd})
			}
		}
		if len(filesGroup) > 0 {
			mapDataFiles := []string{}

			for i, fileData := range filesGroup {
				mapDataFiles = append(
					mapDataFiles,
					fmt.Sprintf(`"%d":[%s]`, i, strings.Join(collect(fileData, wrapMapKeyInQuotes), ",")),
				)
			}

			mapData = `{` + strings.Join(mapDataFiles, ",") + `}`
		}
		bodyWriter.WriteField("map", mapData)

		// --b7955bd2e1d17b67ac157b9e9ddb6238888caefc6f3541920a1debad284d
		// Content-Disposition: form-data; name="0"; filename="tempFile"
		// Content-Type: text/plain; charset=utf-8
		// or
		// Content-Type: application/octet-stream
		//
		for i, fileData := range filesGroup {
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%d"; filename="%s"`, i, fileData[0].file.Name()))
			b, _ := ioutil.ReadFile(fileData[0].file.Name())
			h.Set("Content-Type", http.DetectContentType(b))
			ff, _ := bodyWriter.CreatePart(h)
			ff.Write(b)
		}
		bodyWriter.Close()

		bd.HTTP.Body = ioutil.NopCloser(bodyBuf)
		bd.HTTP.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	}
}

func collect(strArr []*fileFormDataMap, f func(s *fileFormDataMap) string) []string {
	result := make([]string, len(strArr))
	for i, str := range strArr {
		result[i] = f(str)
	}
	return result
}

func wrapMapKeyInQuotes(s *fileFormDataMap) string {
	return fmt.Sprintf("\"%s\"", s.mapKey)
}
