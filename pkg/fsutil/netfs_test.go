package fsutil_test

import (
	"testing"

	"github.com/stashapp/stash/pkg/fsutil"
)

func TestIsNetworkFS_LocalPath(t *testing.T) {
	dir := t.TempDir()
	isNet, err := fsutil.IsNetworkFS(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isNet {
		t.Error("expected a local temp directory not to be classified as a network FS")
	}
}
