package database

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/internal/cli/cover"
	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func TestOpenRequiresExistingDatabase(t *testing.T) {
	initTestLogger()

	dbPath := t.TempDir() + "/stash-go.sqlite"
	_, err := Open(dbPath, Options{})
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("err = %v, want os.ErrNotExist", err)
	}
	if _, statErr := os.Stat(dbPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("database was created: stat err = %v", statErr)
	}
}

func TestOpenConfiguresBlobStore(t *testing.T) {
	initTestLogger()

	dbPath := createExistingDatabase(t)
	store, err := Open(dbPath, Options{BlobStorage: BlobStorageDatabase})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	service := cover.New(store.Repo, t.TempDir())
	if _, err := service.Load(context.Background(), cover.Request{SceneID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestOpenExistingDatabase(t *testing.T) {
	initTestLogger()

	dbPath := createExistingDatabase(t)
	store, err := Open(dbPath, Options{BlobStorage: BlobStorageDatabase})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if store.Schema == 0 {
		t.Fatal("Schema should be populated")
	}
	if err := store.DB.Ready(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenRequiresPath(t *testing.T) {
	_, err := Open("", Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOpenReadsFilesystemBlobs(t *testing.T) {
	initTestLogger()

	dbPath := t.TempDir() + "/stash-go.sqlite"
	blobsPath := filepath.Join(t.TempDir(), "blobs")
	db := sqlite.NewDatabase()
	db.SetBlobStoreOptions(sqlite.BlobStoreOptions{UseFilesystem: true, Path: blobsPath})
	if err := db.Open(dbPath); err != nil {
		t.Fatal(err)
	}
	scene := models.NewScene()
	repo := db.Repository()
	if err := repo.WithTxn(context.Background(), func(ctx context.Context) error {
		if err := repo.Scene.Create(ctx, &scene, nil); err != nil {
			return err
		}
		return repo.Scene.UpdateCover(ctx, scene.ID, []byte("jpg"))
	}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	store, err := Open(dbPath, Options{BlobStorage: BlobStorageFilesystem, BlobsPath: blobsPath})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var loaded []byte
	if err := store.Repo.WithReadTxn(context.Background(), func(ctx context.Context) error {
		var err error
		loaded, err = store.Repo.Scene.GetCover(ctx, scene.ID)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if string(loaded) != "jpg" {
		t.Fatalf("cover = %q, want jpg", loaded)
	}
}

func createExistingDatabase(t *testing.T) string {
	t.Helper()

	dbPath := t.TempDir() + "/stash-go.sqlite"
	db := sqlite.NewDatabase()
	db.SetBlobStoreOptions(sqlite.BlobStoreOptions{UseDatabase: true})
	if err := db.Open(dbPath); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	return dbPath
}

func initTestLogger() {
	if logger.Logger != nil {
		return
	}
	l := log.NewLogger()
	l.Init("", true, "Error", 0)
	logger.Logger = l
}
