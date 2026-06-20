package coverfetch

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

type fakeFinder struct {
	fingerprints models.Fingerprints
	queries      []string
	scenes       []*models.ScrapedScene
	queryScenes  map[string][]*models.ScrapedScene
}

func (f *fakeFinder) FindSceneByFingerprints(_ context.Context, fingerprints models.Fingerprints) ([]*models.ScrapedScene, error) {
	f.fingerprints = fingerprints
	return f.scenes, nil
}

func (f *fakeFinder) QueryScene(_ context.Context, query string) ([]*models.ScrapedScene, error) {
	f.queries = append(f.queries, query)
	return f.queryScenes[query], nil
}

func TestFetchWritesOfficialCoverFromStashBox(t *testing.T) {
	image := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString([]byte("jpg"))
	remoteID := "remote-scene"
	finder := &fakeFinder{
		scenes: []*models.ScrapedScene{{
			Image:        &image,
			RemoteSiteID: &remoteID,
		}},
	}

	var updatedID int
	var updatedCover []byte
	service := &Service{
		finder: finder,
		getFingerprints: func(_ context.Context, sceneID int) (models.Fingerprints, error) {
			return models.Fingerprints{{Type: models.FingerprintTypeOshash, Fingerprint: "abc"}}, nil
		},
		updateCover: func(_ context.Context, sceneID int, cover []byte) error {
			updatedID = sceneID
			updatedCover = cover
			return nil
		},
		processImage: defaultProcessImage,
	}

	result, err := service.Fetch(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}

	if updatedID != 42 {
		t.Fatalf("updatedID = %d, want 42", updatedID)
	}
	if string(updatedCover) != "jpg" {
		t.Fatalf("updatedCover = %q, want jpg", updatedCover)
	}
	if result.RemoteSiteID != remoteID {
		t.Fatalf("RemoteSiteID = %q, want %q", result.RemoteSiteID, remoteID)
	}
	if finder.fingerprints[0].String() != "abc" {
		t.Fatalf("fingerprints = %#v", finder.fingerprints)
	}
}

func TestFetchFallsBackToExactCodeSearchWhenFingerprintMisses(t *testing.T) {
	image := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString([]byte("jpg"))
	remoteID := "remote-scene"
	finder := &fakeFinder{
		queryScenes: map[string][]*models.ScrapedScene{
			"SSPD-082": {{
				Code:         stringPtr("SSPD-082"),
				Image:        &image,
				RemoteSiteID: &remoteID,
				Duration:     intPtr(6940),
			}},
		},
	}

	var updatedCover []byte
	service := &Service{
		finder: finder,
		getFingerprints: func(_ context.Context, sceneID int) (models.Fingerprints, error) {
			return models.Fingerprints{{Type: models.FingerprintTypeOshash, Fingerprint: "77b44c45691f809b"}}, nil
		},
		getSceneInfo: func(_ context.Context, sceneID int) (sceneInfo, error) {
			return sceneInfo{Basename: "SSPD-082-07032012.mp4", Duration: 6940.44}, nil
		},
		updateCover: func(_ context.Context, sceneID int, cover []byte) error {
			updatedCover = cover
			return nil
		},
		processImage: defaultProcessImage,
	}

	result, err := service.Fetch(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}

	if result.RemoteSiteID != remoteID {
		t.Fatalf("RemoteSiteID = %q, want %q", result.RemoteSiteID, remoteID)
	}
	if string(updatedCover) != "jpg" {
		t.Fatalf("updatedCover = %q, want jpg", updatedCover)
	}
	if len(finder.queries) != 1 || finder.queries[0] != "SSPD-082" {
		t.Fatalf("queries = %#v, want [SSPD-082]", finder.queries)
	}
}

func TestFetchDoesNotUseInexactCodeSearchResult(t *testing.T) {
	image := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString([]byte("jpg"))
	finder := &fakeFinder{
		queryScenes: map[string][]*models.ScrapedScene{
			"WANZ-107": {{
				Code:     stringPtr("SSIS-107"),
				Image:    &image,
				Duration: intPtr(7103),
			}},
		},
	}

	service := &Service{
		finder: finder,
		getFingerprints: func(_ context.Context, sceneID int) (models.Fingerprints, error) {
			return models.Fingerprints{{Type: models.FingerprintTypeOshash, Fingerprint: "602c617369e2a8da"}}, nil
		},
		getSceneInfo: func(_ context.Context, sceneID int) (sceneInfo, error) {
			return sceneInfo{Basename: "Wanz-107.mp4", Duration: 7281.85}, nil
		},
		updateCover: func(_ context.Context, sceneID int, cover []byte) error {
			t.Fatal("should not update cover")
			return nil
		},
		processImage: defaultProcessImage,
	}

	_, err := service.Fetch(context.Background(), 42)
	if !errors.Is(err, ErrNoMatch) {
		t.Fatalf("err = %v, want ErrNoMatch", err)
	}
}

func TestFetchRequiresConfiguredFinder(t *testing.T) {
	service := &Service{}

	_, err := service.Fetch(context.Background(), 42)
	if err == nil {
		t.Fatal("expected error")
	}
}

func stringPtr(v string) *string {
	return &v
}

func intPtr(v int) *int {
	return &v
}
