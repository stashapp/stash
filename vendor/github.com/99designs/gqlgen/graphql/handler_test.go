package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddUploadToOperations(t *testing.T) {
	key := "0"

	t.Run("fail missing all variables", func(t *testing.T) {
		file, _ := os.Open("path/to/file")
		params := &RawParams{}

		upload := Upload{
			File:        file,
			Filename:    "a.txt",
			Size:        int64(5),
			ContentType: "text/plain",
		}
		path := "variables.req.0.file"
		err := params.AddUpload(upload, key, path)
		require.NotNil(t, err)
		require.Equal(t, "input: path is missing \"variables.\" prefix, key: 0, path: variables.req.0.file", err.Error())
	})

	t.Run("valid variable", func(t *testing.T) {
		file, _ := os.Open("path/to/file")
		request := &RawParams{
			Variables: map[string]interface{}{
				"file": nil,
			},
		}

		upload := Upload{
			File:        file,
			Filename:    "a.txt",
			Size:        int64(5),
			ContentType: "text/plain",
		}

		expected := &RawParams{
			Variables: map[string]interface{}{
				"file": upload,
			},
		}

		path := "variables.file"
		err := request.AddUpload(upload, key, path)
		require.Nil(t, err)

		require.Equal(t, request, expected)
	})

	t.Run("valid nested variable", func(t *testing.T) {
		file, _ := os.Open("path/to/file")
		request := &RawParams{
			Variables: map[string]interface{}{
				"req": []interface{}{
					map[string]interface{}{
						"file": nil,
					},
				},
			},
		}

		upload := Upload{
			File:        file,
			Filename:    "a.txt",
			Size:        int64(5),
			ContentType: "text/plain",
		}

		expected := &RawParams{
			Variables: map[string]interface{}{
				"req": []interface{}{
					map[string]interface{}{
						"file": upload,
					},
				},
			},
		}

		path := "variables.req.0.file"
		err := request.AddUpload(upload, key, path)
		require.Nil(t, err)

		require.Equal(t, request, expected)
	})
}
