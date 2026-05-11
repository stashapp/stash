package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/internal/manager/config"
)

func TestCleanSceneUploadFilename(t *testing.T) {
	t.Parallel()

	mgr := &Manager{}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "plain filename",
			input: "scene.mp4",
			want:  "scene.mp4",
		},
		{
			name:  "browser path",
			input: `C:\fakepath\scene.mp4`,
			want:  "scene.mp4",
		},
		{
			name:  "unix path",
			input: "/tmp/scene.mp4",
			want:  "scene.mp4",
		},
		{
			name:    "empty filename",
			input:   "",
			wantErr: true,
		},
		{
			name:    "current directory",
			input:   ".",
			wantErr: true,
		},
		{
			name:    "null byte",
			input:   "scene\x00.mp4",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mgr.cleanSceneUploadFilename(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("cleanSceneUploadFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("cleanSceneUploadFilename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUniqueUploadPath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "scene.mp4"), []byte("one"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "scene (1).mp4"), []byte("two"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := uniqueUploadPath(dir, "scene.mp4")
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(dir, "scene (2).mp4")
	if got != want {
		t.Fatalf("uniqueUploadPath() = %q, want %q", got, want)
	}
}

func TestUniqueUploadPathSkipsDirectories(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "scene.mp4"), 0755); err != nil {
		t.Fatal(err)
	}

	got, err := uniqueUploadPath(dir, "scene.mp4")
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(dir, "scene (1).mp4")
	if got != want {
		t.Fatalf("uniqueUploadPath() = %q, want %q", got, want)
	}
}

func TestSceneUploadDestination(t *testing.T) {
	libraryDir := t.TempDir()
	cfg := config.InitializeEmpty()
	cfg.SetInterface(config.Stash, []*config.StashConfig{
		{Path: libraryDir},
	})

	mgr := &Manager{
		Config: cfg,
	}

	got, err := mgr.sceneUploadDestination(nil)
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(libraryDir, "uploads")
	if got != want {
		t.Fatalf("sceneUploadDestination() = %q, want %q", got, want)
	}

	customDestination := filepath.Join(libraryDir, "custom")
	got, err = mgr.sceneUploadDestination(&customDestination)
	if err != nil {
		t.Fatal(err)
	}
	if got != customDestination {
		t.Fatalf("sceneUploadDestination() = %q, want %q", got, customDestination)
	}

	outsideDestination := filepath.Join(t.TempDir(), "outside")
	if _, err = mgr.sceneUploadDestination(&outsideDestination); err == nil {
		t.Fatal("sceneUploadDestination() expected error for destination outside library")
	}
}

func TestSceneUploadDestinationRejectsVideoExcludedPath(t *testing.T) {
	libraryDir := t.TempDir()
	cfg := config.InitializeEmpty()
	cfg.SetInterface(config.Stash, []*config.StashConfig{
		{Path: libraryDir, ExcludeVideo: true},
	})

	mgr := &Manager{
		Config: cfg,
	}

	if _, err := mgr.sceneUploadDestination(nil); err == nil {
		t.Fatal("sceneUploadDestination() expected error without video-enabled library paths")
	}

	customDestination := filepath.Join(libraryDir, "custom")
	if _, err := mgr.sceneUploadDestination(&customDestination); err == nil {
		t.Fatal("sceneUploadDestination() expected error for video-excluded library path")
	}
}

func TestWriteSceneUpload(t *testing.T) {
	cfg := config.InitializeEmpty()
	cfg.SetInterface(config.VideoExtensions, []string{"mp4"})
	mgr := &Manager{
		Config: cfg,
	}

	dir := t.TempDir()
	got, err := mgr.writeSceneUpload(dir, graphql.Upload{
		Filename: `C:\fakepath\scene.mp4`,
		File:     strings.NewReader("video"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got != filepath.Join(dir, "scene.mp4") {
		t.Fatalf("writeSceneUpload() = %q, want scene.mp4 destination", got)
	}

	contents, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "video" {
		t.Fatalf("uploaded file contents = %q, want %q", string(contents), "video")
	}
}

func TestWriteSceneUploadRejectsUnsupportedExtension(t *testing.T) {
	cfg := config.InitializeEmpty()
	cfg.SetInterface(config.VideoExtensions, []string{"mp4"})
	mgr := &Manager{
		Config: cfg,
	}

	_, err := mgr.writeSceneUpload(t.TempDir(), graphql.Upload{
		Filename: "scene.txt",
		File:     strings.NewReader("text"),
	})
	if err == nil {
		t.Fatal("writeSceneUpload() expected error for unsupported extension")
	}
}
