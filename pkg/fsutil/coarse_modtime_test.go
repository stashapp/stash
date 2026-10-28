package fsutil_test

import (
	"testing"

	"github.com/stashapp/stash/pkg/fsutil"
)

func TestCoarseModtimeFilesystem_LocalTempDirNotMarkedCoarse(t *testing.T) {
	dir := t.TempDir()
	coarse, err := fsutil.CoarseModtimeFilesystem(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coarse {
		t.Fatal("expected a local temp directory not to be classified as coarse-modtime FS")
	}
}
