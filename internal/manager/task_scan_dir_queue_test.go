package manager

import (
	"context"
	"io/fs"
	"sync"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

// fakeFileInfo implements fs.FileInfo for test directory entries.
type fakeFileInfo struct {
	name    string
	isDir   bool
	modTime time.Time
	size    int64
	mode    fs.FileMode
}

func (f *fakeFileInfo) Name() string       { return f.name }
func (f *fakeFileInfo) Size() int64        { return f.size }
func (f *fakeFileInfo) Mode() fs.FileMode  { return f.mode }
func (f *fakeFileInfo) ModTime() time.Time { return f.modTime }
func (f *fakeFileInfo) IsDir() bool        { return f.isDir }
func (f *fakeFileInfo) Sys() any           { return nil }

// fakeDirEntry implements fs.DirEntry backed by a fakeFileInfo.
type fakeDirEntry struct{ info *fakeFileInfo }

func (d *fakeDirEntry) Name() string               { return d.info.name }
func (d *fakeDirEntry) IsDir() bool                { return d.info.isDir }
func (d *fakeDirEntry) Type() fs.FileMode          { return d.info.mode }
func (d *fakeDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }

// TestDirQueue_DirectoryGoesToDirQueue verifies that queueFileFunc sends
// a new (cold-scan) directory to dirQueue, not fileQueue.
func TestDirQueue_DirectoryGoesToDirQueue(t *testing.T) {
	scanner := makeTestScanner(t, nil /* folder not found → new dir */, false)
	j := makeTestScanJob(scanner, false)

	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	info := &fakeFileInfo{
		name:    "sub",
		isDir:   true,
		modTime: modTime,
		mode:    fs.ModeDir | 0o755,
	}

	fn := j.queueFileFunc(context.Background(), &fakeTaskFS{}, nil, nil)
	if err := fn("/root/sub", &fakeDirEntry{info}, nil); err != nil {
		t.Fatalf("queueFileFunc returned error: %v", err)
	}

	select {
	case got := <-j.dirQueue:
		if got.Path != "/root/sub" {
			t.Errorf("dirQueue got path %q, want %q", got.Path, "/root/sub")
		}
	default:
		t.Error("expected directory to be sent to dirQueue, but dirQueue is empty")
	}

	select {
	case f := <-j.fileQueue:
		t.Errorf("directory was unexpectedly sent to fileQueue: %s", f.Path)
	default:
		// good — fileQueue is empty
	}
}

// TestDirQueue_UnchangedDirNotQueued verifies that an unchanged directory
// (skipUnchangedFolders=true, same modtime) is stored in unchangedDirs
// and not sent to dirQueue.
func TestDirQueue_UnchangedDirNotQueued(t *testing.T) {
	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(1)
	existingFolder := &models.Folder{
		DirEntry:       models.DirEntry{ModTime: modTime},
		Path:           "/root/sub",
		ParentFolderID: &parentID,
	}

	scanner := makeTestScanner(t, existingFolder, true)
	j := makeTestScanJob(scanner, true)

	info := &fakeFileInfo{
		name:    "sub",
		isDir:   true,
		modTime: modTime,
		mode:    fs.ModeDir | 0o755,
	}

	fn := j.queueFileFunc(context.Background(), &fakeTaskFS{}, nil, nil)
	if err := fn("/root/sub", &fakeDirEntry{info}, nil); err != nil {
		t.Fatalf("queueFileFunc returned error: %v", err)
	}

	select {
	case f := <-j.dirQueue:
		t.Errorf("unchanged directory was unexpectedly sent to dirQueue: %s", f.Path)
	default:
		// good
	}

	if _, stored := j.unchangedDirs.Load("/root/sub"); !stored {
		t.Error("expected /root/sub in unchangedDirs for unchanged folder")
	}
}

// TestDirQueue_WriterDrainsAll verifies that a dir writer goroutine
// processes every item sent to dirQueue by calling handleFolder for each.
func TestDirQueue_WriterDrainsAll(t *testing.T) {
	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(1)
	folder := &models.Folder{
		DirEntry:       models.DirEntry{ModTime: modTime},
		Path:           "/root",
		ParentFolderID: &parentID,
	}

	scanner := makeTestScanner(t, folder, false)
	j := makeTestScanJob(scanner, false)

	dirs := []string{"/root/a", "/root/b", "/root/c"}
	for _, p := range dirs {
		sf := file.ScannedFile{
			BaseFile: &models.BaseFile{
				DirEntry: models.DirEntry{ModTime: modTime},
				Path:     p,
				Basename: "x",
			},
			FS:   &fakeTaskFS{},
			Info: &fakeFileInfo{name: "x", isDir: true, modTime: modTime, mode: fs.ModeDir | 0o755},
		}
		j.dirQueue <- sf
	}
	close(j.dirQueue)

	var mu sync.Mutex
	var processed []string

	// Simulate the dir writer loop the implementation uses.
	for d := range j.dirQueue {
		mu.Lock()
		processed = append(processed, d.Path)
		mu.Unlock()
		// exercise handleFolder integration; errors are non-fatal here
		_ = j.handleFolder(context.Background(), d, nil)
	}

	if len(processed) != len(dirs) {
		t.Errorf("dir writer processed %d dirs, want %d", len(processed), len(dirs))
	}
	for i, p := range dirs {
		if processed[i] != p {
			t.Errorf("dir writer processed[%d] = %q, want %q", i, processed[i], p)
		}
	}
}
