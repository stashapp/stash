package file

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// fakeRenamerRemover records Rename/Remove calls and delegates to real OS operations.
type fakeRenamerRemover struct {
	renames []renameOp
	removes []string
}

type renameOp struct{ from, to string }

func (f *fakeRenamerRemover) Rename(oldpath, newpath string) error {
	f.renames = append(f.renames, renameOp{oldpath, newpath})
	return os.Rename(oldpath, newpath)
}

func (f *fakeRenamerRemover) Remove(name string) error {
	f.removes = append(f.removes, name)
	return os.Remove(name)
}

func (f *fakeRenamerRemover) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (f *fakeRenamerRemover) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func newDeleterWithFake(t *testing.T) (*Deleter, *fakeRenamerRemover) {
	t.Helper()
	fake := &fakeRenamerRemover{}
	d := NewDeleter()
	d.RenamerRemover = fake
	return d, fake
}

func TestDeleter_ZipEntry_Accumulates(t *testing.T) {
	d, _ := newDeleterWithFake(t)

	d.ZipEntry("/path/to/gallery.zip", "a.jpg")
	d.ZipEntry("/path/to/gallery.zip", "b.jpg")

	entries := d.zipEntries["/path/to/gallery.zip"]
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d: %v", len(entries), entries)
	}
}

func TestDeleter_ZipEntry_MultipleZips(t *testing.T) {
	d, _ := newDeleterWithFake(t)

	d.ZipEntry("/zip/one.zip", "img1.jpg")
	d.ZipEntry("/zip/two.zip", "img2.jpg")
	d.ZipEntry("/zip/one.zip", "img3.jpg")

	if len(d.zipEntries["/zip/one.zip"]) != 2 {
		t.Errorf("one.zip: expected 2 entries, got %d", len(d.zipEntries["/zip/one.zip"]))
	}
	if len(d.zipEntries["/zip/two.zip"]) != 1 {
		t.Errorf("two.zip: expected 1 entry, got %d", len(d.zipEntries["/zip/two.zip"]))
	}
}

func TestDeleter_Rollback_ClearsZipEntries(t *testing.T) {
	d, _ := newDeleterWithFake(t)

	d.ZipEntry("/path/to/gallery.zip", "a.jpg")
	d.Rollback()

	if d.zipEntries != nil {
		t.Errorf("expected zipEntries to be nil after Rollback, got %v", d.zipEntries)
	}
}

func TestDeleter_Commit_RewritesZip(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "gallery.zip")
	makeTestZip(t, zipPath, "", map[string]string{
		"keep.jpg":   "keep",
		"remove.jpg": "remove",
	})

	d, fake := newDeleterWithFake(t)
	d.ZipEntry(zipPath, "remove.jpg")
	d.Commit()

	// Rename should have been called once: tmpPath → zipPath
	if len(fake.renames) != 1 {
		t.Fatalf("expected 1 Rename call, got %d", len(fake.renames))
	}
	if fake.renames[0].to != zipPath {
		t.Errorf("Rename target = %q, want %q", fake.renames[0].to, zipPath)
	}

	// The original zip should now contain only keep.jpg
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("open rewritten zip: %v", err)
	}
	defer r.Close()

	names := make(map[string]bool)
	for _, f := range r.File {
		names[f.Name] = true
	}
	if names["remove.jpg"] {
		t.Error("remove.jpg still present after Commit")
	}
	if !names["keep.jpg"] {
		t.Error("keep.jpg missing after Commit")
	}
}

func TestDeleter_Commit_ClearsZipEntries(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "gallery.zip")
	makeTestZip(t, zipPath, "", map[string]string{
		"a.jpg": "a",
	})

	d, _ := newDeleterWithFake(t)
	d.ZipEntry(zipPath, "a.jpg")
	d.Commit()

	if d.zipEntries != nil {
		t.Errorf("expected zipEntries to be nil after Commit, got %v", d.zipEntries)
	}
}

func TestDeleter_Commit_MultipleZips(t *testing.T) {
	dir := t.TempDir()

	zip1 := filepath.Join(dir, "one.zip")
	zip2 := filepath.Join(dir, "two.zip")
	makeTestZip(t, zip1, "", map[string]string{"a.jpg": "a", "b.jpg": "b"})
	makeTestZip(t, zip2, "", map[string]string{"c.jpg": "c", "d.jpg": "d"})

	d, fake := newDeleterWithFake(t)
	d.ZipEntry(zip1, "a.jpg")
	d.ZipEntry(zip2, "c.jpg")
	d.Commit()

	if len(fake.renames) != 2 {
		t.Fatalf("expected 2 Rename calls, got %d", len(fake.renames))
	}

	for _, zp := range []string{zip1, zip2} {
		r, err := zip.OpenReader(zp)
		if err != nil {
			t.Fatalf("open zip %s: %v", zp, err)
		}
		names := make(map[string]bool)
		for _, f := range r.File {
			names[f.Name] = true
		}
		r.Close()

		switch zp {
		case zip1:
			if names["a.jpg"] {
				t.Error("one.zip: a.jpg should have been removed")
			}
			if !names["b.jpg"] {
				t.Error("one.zip: b.jpg should be present")
			}
		case zip2:
			if names["c.jpg"] {
				t.Error("two.zip: c.jpg should have been removed")
			}
			if !names["d.jpg"] {
				t.Error("two.zip: d.jpg should be present")
			}
		}
	}
}
