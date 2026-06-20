package cover

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCache(t *testing.T) {
	service := &Service{cacheDir: t.TempDir()}
	path, err := service.WriteCache(42, []byte("jpg"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "jpg" {
		t.Fatalf("cached data = %q", got)
	}
}

func TestLoadDoesNotGenerateCoverWhenCacheIsMissing(t *testing.T) {
	service := &Service{
		cacheDir: t.TempDir(),
	}

	got, err := service.Load(context.Background(), Request{
		SceneID:  42,
		Path:     filepath.Join(t.TempDir(), "scene.mp4"),
		Duration: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got.Source != SourceMissing {
		t.Fatalf("Source = %q, want %q", got.Source, SourceMissing)
	}
}

func TestLoadIgnoresLegacyGeneratedCache(t *testing.T) {
	service := &Service{cacheDir: t.TempDir()}
	if _, err := service.WriteCache(42, []byte("old-generated-jpg")); err != nil {
		t.Fatal(err)
	}

	got, err := service.Load(context.Background(), Request{SceneID: 42})
	if err != nil {
		t.Fatal(err)
	}

	if got.Source != SourceMissing {
		t.Fatalf("Source = %q, want %q", got.Source, SourceMissing)
	}
	if len(got.Data) != 0 {
		t.Fatalf("Data length = %d, want 0", len(got.Data))
	}
}
