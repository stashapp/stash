package client_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		require.Equal(t, `{"query":"user(id:$id){name}","variables":{"id":1}}`, string(b))

		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"name": "bob",
			},
		})
		if err != nil {
			panic(err)
		}
	})

	c := client.New(h)

	var resp struct {
		Name string
	}

	c.MustPost("user(id:$id){name}", &resp, client.Var("id", 1))

	require.Equal(t, "bob", resp.Name)
}

func TestClientMultipartFormData(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Contains(t, string(bodyBytes), `Content-Disposition: form-data; name="operations"`)
		require.Contains(t, string(bodyBytes), `{"query":"mutation ($input: Input!) {}","variables":{"file":{}}`)
		require.Contains(t, string(bodyBytes), `Content-Disposition: form-data; name="map"`)
		require.Contains(t, string(bodyBytes), `{"0":["variables.file"]}`)
		require.Contains(t, string(bodyBytes), `Content-Disposition: form-data; name="0"; filename="example.txt"`)
		require.Contains(t, string(bodyBytes), `Content-Type: text/plain`)
		require.Contains(t, string(bodyBytes), `Hello World`)

		w.Write([]byte(`{}`))
	})

	c := client.New(h)

	var resp struct{}
	c.MustPost("{ id }", &resp,
		func(bd *client.Request) {
			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			bodyWriter.WriteField("operations", `{"query":"mutation ($input: Input!) {}","variables":{"file":{}}`)
			bodyWriter.WriteField("map", `{"0":["variables.file"]}`)

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", `form-data; name="0"; filename="example.txt"`)
			h.Set("Content-Type", "text/plain")
			ff, _ := bodyWriter.CreatePart(h)
			ff.Write([]byte("Hello World"))
			bodyWriter.Close()

			bd.HTTP.Body = ioutil.NopCloser(bodyBuf)
			bd.HTTP.Header.Set("Content-Type", bodyWriter.FormDataContentType())
		},
	)
}

func TestAddHeader(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "ASDF", r.Header.Get("Test-Key"))

		w.Write([]byte(`{}`))
	})

	c := client.New(h)

	var resp struct{}
	c.MustPost("{ id }", &resp,
		client.AddHeader("Test-Key", "ASDF"),
	)
}

func TestAddClientHeader(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "ASDF", r.Header.Get("Test-Key"))

		w.Write([]byte(`{}`))
	})

	c := client.New(h, client.AddHeader("Test-Key", "ASDF"))

	var resp struct{}
	c.MustPost("{ id }", &resp)
}

func TestBasicAuth(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		require.True(t, ok)
		require.Equal(t, "user", user)
		require.Equal(t, "pass", pass)

		w.Write([]byte(`{}`))
	})

	c := client.New(h)

	var resp struct{}
	c.MustPost("{ id }", &resp,
		client.BasicAuth("user", "pass"),
	)
}

func TestAddCookie(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("foo")
		require.NoError(t, err)
		require.Equal(t, "value", c.Value)

		w.Write([]byte(`{}`))
	})

	c := client.New(h)

	var resp struct{}
	c.MustPost("{ id }", &resp,
		client.AddCookie(&http.Cookie{Name: "foo", Value: "value"}),
	)
}

func TestAddExtensions(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		require.Equal(t, `{"query":"user(id:1){name}","extensions":{"persistedQuery":{"sha256Hash":"ceec2897e2da519612279e63f24658c3e91194cbb2974744fa9007a7e1e9f9e7","version":1}}}`, string(b))
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"Name": "Bob",
			},
		})
		if err != nil {
			panic(err)
		}
	})

	c := client.New(h)

	var resp struct {
		Name string
	}
	c.MustPost("user(id:1){name}", &resp,
		client.Extensions(map[string]interface{}{"persistedQuery": map[string]interface{}{"version": 1, "sha256Hash": "ceec2897e2da519612279e63f24658c3e91194cbb2974744fa9007a7e1e9f9e7"}}),
	)
}
