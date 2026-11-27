package stashbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/stashbox/graphql"
)

func TestSceneFragmentToScrapedSceneImage(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write([]byte{0xff, 0xd8, 0xff, 0xd9})
	}))
	defer srv.Close()

	fragment := &graphql.SceneFragment{
		ID:     "scene-id",
		Images: []*graphql.ImageFragment{{ID: "image-id", URL: srv.URL + "/cover.jpg"}},
	}

	tests := []struct {
		name      string
		options   []ClientOption
		wantImage bool
	}{
		{"default downloads the cover", nil, true},
		{"SkipImages(true) skips the download", []ClientOption{SkipImages(true)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hits.Store(0)
			c := NewClient(models.StashBox{Endpoint: srv.URL + "/graphql"}, tt.options...)

			got, err := c.sceneFragmentToScrapedScene(context.Background(), fragment)
			if err != nil {
				t.Fatalf("sceneFragmentToScrapedScene() error = %v", err)
			}

			if gotImage := got.Image != nil; gotImage != tt.wantImage {
				t.Errorf("image set = %v, want %v", gotImage, tt.wantImage)
			}
			wantHits := int32(0)
			if tt.wantImage {
				wantHits = 1
			}
			if n := hits.Load(); n != wantHits {
				t.Errorf("image requests = %d, want %d", n, wantHits)
			}
		})
	}
}
