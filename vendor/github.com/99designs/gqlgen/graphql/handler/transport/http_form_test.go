package transport_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestFileUpload(t *testing.T) {
	es := &graphql.ExecutableSchemaMock{
		ExecFunc: func(ctx context.Context) graphql.ResponseHandler {
			return graphql.OneShot(graphql.ErrorResponse(ctx, "not implemented"))
		},
		SchemaFunc: func() *ast.Schema {
			return gqlparser.MustLoadSchema(&ast.Source{Input: `
				type Mutation {
					singleUpload(file: Upload!): String!
					singleUploadWithPayload(req: UploadFile!): String!
					multipleUpload(files: [Upload!]!): String!
					multipleUploadWithPayload(req: [UploadFile!]!): String!
				}
				scalar Upload
				scalar UploadFile
			`})
		},
	}

	h := handler.New(es)
	multipartForm := transport.MultipartForm{}
	h.AddTransport(&multipartForm)

	t.Run("valid single file upload", func(t *testing.T) {
		es.ExecFunc = func(ctx context.Context) graphql.ResponseHandler {
			op := graphql.GetOperationContext(ctx).Operation
			require.Equal(t, len(op.VariableDefinitions), 1)
			require.Equal(t, op.VariableDefinitions[0].Variable, "file")
			return graphql.OneShot(&graphql.Response{Data: []byte(`{"singleUpload":"test"}`)})
		}

		operations := `{ "query": "mutation ($file: Upload!) { singleUpload(file: $file) }", "variables": { "file": null } }`
		mapData := `{ "0": ["variables.file"] }`
		files := []file{
			{
				mapKey:      "0",
				name:        "a.txt",
				content:     "test1",
				contentType: "text/plain",
			},
		}
		req := createUploadRequest(t, operations, mapData, files)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"data":{"singleUpload":"test"}}`, resp.Body.String())
	})

	t.Run("valid single file upload with payload", func(t *testing.T) {
		es.ExecFunc = func(ctx context.Context) graphql.ResponseHandler {
			op := graphql.GetOperationContext(ctx).Operation
			require.Equal(t, len(op.VariableDefinitions), 1)
			require.Equal(t, op.VariableDefinitions[0].Variable, "req")
			return graphql.OneShot(&graphql.Response{Data: []byte(`{"singleUploadWithPayload":"test"}`)})
		}

		operations := `{ "query": "mutation ($req: UploadFile!) { singleUploadWithPayload(req: $req) }", "variables": { "req": {"file": null, "id": 1 } } }`
		mapData := `{ "0": ["variables.req.file"] }`
		files := []file{
			{
				mapKey:      "0",
				name:        "a.txt",
				content:     "test1",
				contentType: "text/plain",
			},
		}
		req := createUploadRequest(t, operations, mapData, files)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"data":{"singleUploadWithPayload":"test"}}`, resp.Body.String())
	})

	t.Run("valid file list upload", func(t *testing.T) {
		es.ExecFunc = func(ctx context.Context) graphql.ResponseHandler {
			op := graphql.GetOperationContext(ctx).Operation
			require.Equal(t, len(op.VariableDefinitions), 1)
			require.Equal(t, op.VariableDefinitions[0].Variable, "files")
			return graphql.OneShot(&graphql.Response{Data: []byte(`{"multipleUpload":[{"id":1},{"id":2}]}`)})
		}

		operations := `{ "query": "mutation($files: [Upload!]!) { multipleUpload(files: $files) }", "variables": { "files": [null, null] } }`
		mapData := `{ "0": ["variables.files.0"], "1": ["variables.files.1"] }`
		files := []file{
			{
				mapKey:      "0",
				name:        "a.txt",
				content:     "test1",
				contentType: "text/plain",
			},
			{
				mapKey:      "1",
				name:        "b.txt",
				content:     "test2",
				contentType: "text/plain",
			},
		}
		req := createUploadRequest(t, operations, mapData, files)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"data":{"multipleUpload":[{"id":1},{"id":2}]}}`, resp.Body.String())
	})

	t.Run("valid file list upload with payload", func(t *testing.T) {
		es.ExecFunc = func(ctx context.Context) graphql.ResponseHandler {
			op := graphql.GetOperationContext(ctx).Operation
			require.Equal(t, len(op.VariableDefinitions), 1)
			require.Equal(t, op.VariableDefinitions[0].Variable, "req")
			return graphql.OneShot(&graphql.Response{Data: []byte(`{"multipleUploadWithPayload":[{"id":1},{"id":2}]}`)})
		}

		operations := `{ "query": "mutation($req: [UploadFile!]!) { multipleUploadWithPayload(req: $req) }", "variables": { "req": [ { "id": 1, "file": null }, { "id": 2, "file": null } ] } }`
		mapData := `{ "0": ["variables.req.0.file"], "1": ["variables.req.1.file"] }`
		files := []file{
			{
				mapKey:      "0",
				name:        "a.txt",
				content:     "test1",
				contentType: "text/plain",
			},
			{
				mapKey:      "1",
				name:        "b.txt",
				content:     "test2",
				contentType: "text/plain",
			},
		}
		req := createUploadRequest(t, operations, mapData, files)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, `{"data":{"multipleUploadWithPayload":[{"id":1},{"id":2}]}}`, resp.Body.String())
	})

	t.Run("valid file list upload with payload and file reuse", func(t *testing.T) {
		test := func(uploadMaxMemory int64) {
			es.ExecFunc = func(ctx context.Context) graphql.ResponseHandler {
				op := graphql.GetOperationContext(ctx).Operation
				require.Equal(t, len(op.VariableDefinitions), 1)
				require.Equal(t, op.VariableDefinitions[0].Variable, "req")
				return graphql.OneShot(&graphql.Response{Data: []byte(`{"multipleUploadWithPayload":[{"id":1},{"id":2}]}`)})
			}
			multipartForm.MaxMemory = uploadMaxMemory

			operations := `{ "query": "mutation($req: [UploadFile!]!) { multipleUploadWithPayload(req: $req) }", "variables": { "req": [ { "id": 1, "file": null }, { "id": 2, "file": null } ] } }`
			mapData := `{ "0": ["variables.req.0.file", "variables.req.1.file"] }`
			files := []file{
				{
					mapKey:      "0",
					name:        "a.txt",
					content:     "test1",
					contentType: "text/plain",
				},
			}
			req := createUploadRequest(t, operations, mapData, files)

			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
			require.Equal(t, `{"data":{"multipleUploadWithPayload":[{"id":1},{"id":2}]}}`, resp.Body.String())
		}

		t.Run("payload smaller than UploadMaxMemory, stored in memory", func(t *testing.T) {
			test(5000)
		})

		t.Run("payload bigger than UploadMaxMemory, persisted to disk", func(t *testing.T) {
			test(2)
		})
	})

	validOperations := `{ "query": "mutation ($file: Upload!) { singleUpload(file: $file) }", "variables": { "file": null } }`
	validMap := `{ "0": ["variables.file"] }`
	validFiles := []file{
		{
			mapKey:      "0",
			name:        "a.txt",
			content:     "test1",
			contentType: "text/plain",
		},
	}

	t.Run("failed to parse multipart", func(t *testing.T) {
		req := &http.Request{
			Method: "POST",
			Header: http.Header{"Content-Type": {`multipart/form-data; boundary="foo123"`}},
			Body:   ioutil.NopCloser(new(bytes.Buffer)),
		}
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"failed to parse multipart form"}],"data":null}`, resp.Body.String())
	})

	t.Run("fail parse operation", func(t *testing.T) {
		operations := `invalid operation`
		req := createUploadRequest(t, operations, validMap, validFiles)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"operations form field could not be decoded"}],"data":null}`, resp.Body.String())
	})

	t.Run("fail parse map", func(t *testing.T) {
		mapData := `invalid map`
		req := createUploadRequest(t, validOperations, mapData, validFiles)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"map form field could not be decoded"}],"data":null}`, resp.Body.String())
	})

	t.Run("fail missing file", func(t *testing.T) {
		var files []file
		req := createUploadRequest(t, validOperations, validMap, files)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"failed to get key 0 from form"}],"data":null}`, resp.Body.String())
	})

	t.Run("fail map entry with invalid operations paths prefix", func(t *testing.T) {
		mapData := `{ "0": ["var.file"] }`
		req := createUploadRequest(t, validOperations, mapData, validFiles)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"invalid operations paths for key 0"}],"data":null}`, resp.Body.String())
	})

	t.Run("fail parse request big body", func(t *testing.T) {
		multipartForm.MaxUploadSize = 2
		req := createUploadRequest(t, validOperations, validMap, validFiles)

		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
		require.Equal(t, `{"errors":[{"message":"failed to parse multipart form, request body too large"}],"data":null}`, resp.Body.String())
	})
}

type file struct {
	mapKey      string
	name        string
	content     string
	contentType string
}

func createUploadRequest(t *testing.T, operations, mapData string, files []file) *http.Request {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := bodyWriter.WriteField("operations", operations)
	require.NoError(t, err)

	err = bodyWriter.WriteField("map", mapData)
	require.NoError(t, err)

	for i := range files {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, files[i].mapKey, files[i].name))
		h.Set("Content-Type", files[i].contentType)
		ff, err := bodyWriter.CreatePart(h)
		require.NoError(t, err)
		_, err = ff.Write([]byte(files[i].content))
		require.NoError(t, err)
	}
	err = bodyWriter.Close()
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/graphql", bodyBuf)
	require.NoError(t, err)

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	return req
}
