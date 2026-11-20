package stashbox

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/stashbox/graphql"
)

func TestSceneFragmentToScrapedSceneRemoteDeleted(t *testing.T) {
	for _, deleted := range []bool{false, true} {
		c := NewClient(models.StashBox{Endpoint: "http://stash-box.invalid/graphql"})

		got, err := c.sceneFragmentToScrapedScene(context.Background(), &graphql.SceneFragment{ID: "scene-id", Deleted: deleted})
		if err != nil {
			t.Fatalf("sceneFragmentToScrapedScene() error = %v", err)
		}

		if got.RemoteDeleted == nil || *got.RemoteDeleted != deleted {
			t.Errorf("deleted = %v: RemoteDeleted = %v, want %v", deleted, got.RemoteDeleted, deleted)
		}
	}
}
