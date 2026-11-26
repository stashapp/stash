package manager

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
)

func createStashIgnoreTestFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}

	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
}

func TestStashIgnoreUsesTopmostLibraryRootWithNestedLibraries(t *testing.T) {
	root := filepath.Join(t.TempDir(), "library")
	nested := filepath.Join(root, "images")
	ignoredFile := filepath.Join(nested, "secret", "image.jpg")

	createStashIgnoreTestFile(t, ignoredFile)
	if err := os.WriteFile(filepath.Join(root, ".stashignore"), []byte("images/secret/\n"), 0644); err != nil {
		t.Fatalf("failed to create .stashignore: %v", err)
	}

	info, err := os.Stat(ignoredFile)
	if err != nil {
		t.Fatalf("failed to stat %s: %v", ignoredFile, err)
	}

	stashPaths := config.StashConfigs{
		{Path: root, ExcludeImage: true},
		{Path: nested},
	}

	scannerFilter := &scanFilter{
		stashPaths:        stashPaths,
		generatedPath:     filepath.Join(root, "generated"),
		stashIgnoreFilter: file.NewStashIgnoreFilter(),
	}

	if scannerFilter.Accept(context.Background(), ignoredFile, info, "") {
		t.Fatalf("expected scan filter to reject file due to parent .stashignore")
	}

	verifyFilter := &verifyFilter{
		scanFilter: scanFilter{
			stashPaths:        stashPaths,
			generatedPath:     filepath.Join(root, "generated"),
			stashIgnoreFilter: file.NewStashIgnoreFilter(),
		},
	}

	if verifyFilter.Accept(context.Background(), ignoredFile, info, "") {
		t.Fatalf("expected clean filter to reject file due to parent .stashignore")
	}
}
