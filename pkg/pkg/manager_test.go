package pkg

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

type testSourcePathGetter struct{}

func (testSourcePathGetter) GetAllSourcePaths() []string        { return []string{"src"} }
func (testSourcePathGetter) GetSourcePath(srcURL string) string { return "src" }

type testPackage struct {
	id       string
	requires []string
	date     string
}

// testSource serves an index of packages and their zips, counting zip downloads.
type testSource struct {
	*httptest.Server
	downloads map[string]int
}

func newTestSource(t *testing.T, pkgs []testPackage) *testSource {
	t.Helper()

	zips := make(map[string][]byte)
	var index strings.Builder
	for _, p := range pkgs {
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		w, err := zw.Create(p.id + ".yml")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(w, "name: %s\n", p.id)
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		zips[p.id] = buf.Bytes()

		date := p.date
		if date == "" {
			date = "2026-01-01 00:00:00"
		}
		fmt.Fprintf(&index, "- id: %s\n  name: %s\n  date: %s\n  path: %s.zip\n  sha256: %x\n  requires: [%s]\n",
			p.id, p.id, date, p.id, sha256.Sum256(zips[p.id]), strings.Join(p.requires, ", "))
	}

	s := &testSource{downloads: make(map[string]int)}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/index.yml" {
			_, _ = w.Write([]byte(index.String()))
			return
		}
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), ".zip")
		s.downloads[id]++
		_, _ = w.Write(zips[id])
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *testSource) indexURL() string { return s.URL + "/index.yml" }

func (s *testSource) spec(id string) models.PackageSpecInput {
	return models.PackageSpecInput{ID: id, SourceURL: s.indexURL()}
}

func newTestManager(t *testing.T, s *testSource) *Manager {
	return &Manager{
		Local:             &Store{BaseDir: t.TempDir(), ManifestFile: ManifestFile},
		PackagePathGetter: testSourcePathGetter{},
		Client:            s.Client(),
	}
}

// preinstall writes a manifest for id as if it were installed with the given date.
func preinstall(t *testing.T, m *Manager, s *testSource, id string, date string) {
	t.Helper()
	store := m.getStore(s.indexURL())
	if err := os.MkdirAll(store.packageDir(id), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	manifest := fmt.Sprintf("id: %s\nname: %s\ndate: %s\nsource_repository: %s\n", id, id, date, s.indexURL())
	if err := os.WriteFile(store.manifestPath(id), []byte(manifest), os.ModePerm); err != nil {
		t.Fatal(err)
	}
}

func assertInstalled(t *testing.T, m *Manager, s *testSource, want map[string]bool) {
	t.Helper()
	store := m.getStore(s.indexURL())
	for id, installed := range want {
		_, err := os.Stat(filepath.Join(store.packageDir(id), id+".yml"))
		if got := err == nil; got != installed {
			t.Errorf("package %s installed = %v, want %v", id, got, installed)
		}
	}
}

func TestInstaller_Install(t *testing.T) {
	ctx := context.Background()

	t.Run("transitive requirements", func(t *testing.T) {
		s := newTestSource(t, []testPackage{
			{id: "a", requires: []string{"b"}},
			{id: "b", requires: []string{"c"}},
			{id: "c"},
		})
		m := newTestManager(t, s)

		if err := m.NewInstaller().Install(ctx, s.spec("a")); err != nil {
			t.Fatal(err)
		}
		assertInstalled(t, m, s, map[string]bool{"a": true, "b": true, "c": true})
	})

	t.Run("current requirement kept, its missing requirement installed", func(t *testing.T) {
		s := newTestSource(t, []testPackage{
			{id: "a", requires: []string{"b"}},
			{id: "b", requires: []string{"c"}},
			{id: "c"},
		})
		m := newTestManager(t, s)
		preinstall(t, m, s, "b", "2026-01-01 00:00:00")

		if err := m.NewInstaller().Install(ctx, s.spec("a")); err != nil {
			t.Fatal(err)
		}
		if s.downloads["b"] != 0 {
			t.Errorf("current requirement b was downloaded %d times", s.downloads["b"])
		}
		assertInstalled(t, m, s, map[string]bool{"a": true, "c": true})
	})

	t.Run("outdated requirement updated", func(t *testing.T) {
		s := newTestSource(t, []testPackage{
			{id: "a", requires: []string{"b"}},
			{id: "b", date: "2026-02-01 00:00:00"},
		})
		m := newTestManager(t, s)
		preinstall(t, m, s, "b", "2026-01-01 00:00:00")

		if err := m.NewInstaller().Install(ctx, s.spec("a")); err != nil {
			t.Fatal(err)
		}
		assertInstalled(t, m, s, map[string]bool{"a": true, "b": true})
	})

	t.Run("missing requirement fails parent but installs siblings", func(t *testing.T) {
		s := newTestSource(t, []testPackage{
			{id: "a", requires: []string{"gone", "b"}},
			{id: "b"},
		})
		m := newTestManager(t, s)

		if err := m.NewInstaller().Install(ctx, s.spec("a")); err == nil {
			t.Fatal("expected error for missing requirement")
		}
		assertInstalled(t, m, s, map[string]bool{"a": false, "b": true})
	})

	t.Run("cycle", func(t *testing.T) {
		s := newTestSource(t, []testPackage{
			{id: "a", requires: []string{"b"}},
			{id: "b", requires: []string{"a"}},
		})
		m := newTestManager(t, s)

		if err := m.NewInstaller().Install(ctx, s.spec("a")); err != nil {
			t.Fatal(err)
		}
		assertInstalled(t, m, s, map[string]bool{"a": true, "b": true})
	})

	t.Run("shared installer installs a requirement once per job", func(t *testing.T) {
		s := newTestSource(t, []testPackage{
			{id: "a", requires: []string{"c"}},
			{id: "b", requires: []string{"c"}},
			{id: "c"},
		})
		m := newTestManager(t, s)
		installer := m.NewInstaller()

		for _, id := range []string{"a", "b", "c"} {
			if err := installer.Install(ctx, s.spec(id)); err != nil {
				t.Fatal(err)
			}
		}
		if s.downloads["c"] != 1 {
			t.Errorf("c downloaded %d times, want 1", s.downloads["c"])
		}
	})

	t.Run("explicit reinstall of a current package", func(t *testing.T) {
		s := newTestSource(t, []testPackage{
			{id: "a", requires: []string{"b"}},
			{id: "b"},
		})
		m := newTestManager(t, s)
		preinstall(t, m, s, "b", "2026-01-01 00:00:00")
		installer := m.NewInstaller()

		for _, id := range []string{"a", "b"} {
			if err := installer.Install(ctx, s.spec(id)); err != nil {
				t.Fatal(err)
			}
		}
		if s.downloads["b"] != 1 {
			t.Errorf("b downloaded %d times, want 1", s.downloads["b"])
		}
	})
}
