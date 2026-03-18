//go:build integration
// +build integration

package manager

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/pkg/file"

	// Necessary to register custom migrations.
	_ "github.com/stashapp/stash/pkg/sqlite/migrations"
)

// stashIgnorePathFilter wraps StashIgnoreFilter to implement PathFilter for testing.
// It provides a fixed library root for the filter.
type stashIgnorePathFilter struct {
	filter      *file.StashIgnoreFilter
	libraryRoot string
}

func (f *stashIgnorePathFilter) Accept(ctx context.Context, path string, info fs.FileInfo, zipFilePath string) bool {
	return f.filter.Accept(ctx, path, info, f.libraryRoot, zipFilePath)
}

// createTestFileOnDisk creates a file with some content.
func createTestFileOnDisk(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}
	// Write some content so the file has a non-zero size.
	if err := os.WriteFile(path, []byte("test content for "+name), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
	return path
}

// createStashIgnoreFile creates a .stashignore file with the given content.
func createStashIgnoreFile(t *testing.T, dir, content string) {
	t.Helper()
	path := filepath.Join(dir, ".stashignore")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create .stashignore: %v", err)
	}
}

func TestScannerWithStashIgnore(t *testing.T) {
	// Create temp directory structure.
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "video1.mp4")
	createTestFileOnDisk(t, tmpDir, "video2.mp4")
	createTestFileOnDisk(t, tmpDir, "ignore_me.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/video3.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/skip_this.mp4")
	createTestFileOnDisk(t, tmpDir, "excluded_dir/video4.mp4")
	createTestFileOnDisk(t, tmpDir, "temp/processing.mp4")

	// Create .stashignore file.
	stashignore := `# Ignore specific files
ignore_me.mp4
subdir/skip_this.mp4

# Ignore directories
excluded_dir/
temp/
`
	createStashIgnoreFile(t, tmpDir, stashignore)

	// Create stashignore filter with library root.
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Create scanner.
	scanner := &file.Scanner{
		ScanFilters: []file.PathFilter{stashIgnoreFilter},
	}

	testScenarios := []struct {
		path     string
		accepted bool
	}{
		{filepath.Join(tmpDir, "video1.mp4"), true},
		{filepath.Join(tmpDir, "video2.mp4"), true},
		{filepath.Join(tmpDir, "ignore_me.mp4"), false},
		{filepath.Join(tmpDir, "subdir/video3.mp4"), true},
		{filepath.Join(tmpDir, "subdir/skip_this.mp4"), false},
		{filepath.Join(tmpDir, "excluded_dir/video4.mp4"), false},
		{filepath.Join(tmpDir, "temp/processing.mp4"), false},
	}

	ctx := context.Background()

	for _, scenario := range testScenarios {
		info, err := os.Stat(scenario.path)
		if err != nil {
			t.Fatalf("failed to stat file %s: %v", scenario.path, err)
		}
		accepted := scanner.AcceptEntry(ctx, scenario.path, info, "")

		if accepted != scenario.accepted {
			t.Errorf("unexpected accept result for %s: expected %v, got %v",
				scenario.path, scenario.accepted, accepted)
		}
	}
}

func TestScannerWithNestedStashIgnore(t *testing.T) {
	// Create temp directory structure.
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "root.mp4")
	createTestFileOnDisk(t, tmpDir, "root.tmp")
	createTestFileOnDisk(t, tmpDir, "subdir/sub.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/sub.log")
	createTestFileOnDisk(t, tmpDir, "subdir/sub.tmp")

	// Root .stashignore excludes *.tmp.
	createStashIgnoreFile(t, tmpDir, "*.tmp\n")

	// Subdir .stashignore excludes *.log.
	createStashIgnoreFile(t, filepath.Join(tmpDir, "subdir"), "*.log\n")

	// Create stashignore filter with library root.
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Create scanner.
	scanner := &file.Scanner{
		ScanFilters: []file.PathFilter{stashIgnoreFilter},
	}

	testScenarios := []struct {
		path     string
		accepted bool
	}{
		{filepath.Join(tmpDir, "root.mp4"), true},
		{filepath.Join(tmpDir, "root.tmp"), false},
		{filepath.Join(tmpDir, "subdir/sub.mp4"), true},
		{filepath.Join(tmpDir, "subdir/sub.log"), false},
		{filepath.Join(tmpDir, "subdir/sub.tmp"), false},
	}

	ctx := context.Background()

	for _, scenario := range testScenarios {
		info, err := os.Stat(scenario.path)
		if err != nil {
			t.Fatalf("failed to stat file %s: %v", scenario.path, err)
		}
		accepted := scanner.AcceptEntry(ctx, scenario.path, info, "")

		if accepted != scenario.accepted {
			t.Errorf("unexpected accept result for %s: expected %v, got %v",
				scenario.path, scenario.accepted, accepted)
		}
	}
}

func TestScannerWithoutStashIgnore(t *testing.T) {
	// Create temp directory structure (no .stashignore).
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "video1.mp4")
	createTestFileOnDisk(t, tmpDir, "video2.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/video3.mp4")

	// Create stashignore filter with library root (but no .stashignore file exists).
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Create scanner.
	scanner := &file.Scanner{
		ScanFilters: []file.PathFilter{stashIgnoreFilter},
	}

	testScenarios := []struct {
		path     string
		accepted bool
	}{
		{filepath.Join(tmpDir, "video1.mp4"), true},
		{filepath.Join(tmpDir, "video2.mp4"), true},
		{filepath.Join(tmpDir, "subdir/video3.mp4"), true},
	}

	ctx := context.Background()

	for _, scenario := range testScenarios {
		info, err := os.Stat(scenario.path)
		if err != nil {
			t.Fatalf("failed to stat file %s: %v", scenario.path, err)
		}
		accepted := scanner.AcceptEntry(ctx, scenario.path, info, "")

		if accepted != scenario.accepted {
			t.Errorf("unexpected accept result for %s: expected %v, got %v",
				scenario.path, scenario.accepted, accepted)
		}
	}
}

func TestScannerWithNegationPattern(t *testing.T) {
	// Create temp directory structure.
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "file1.tmp")
	createTestFileOnDisk(t, tmpDir, "file2.tmp")
	createTestFileOnDisk(t, tmpDir, "keep_this.tmp")
	createTestFileOnDisk(t, tmpDir, "video.mp4")

	// Create .stashignore with negation.
	stashignore := `*.tmp
!keep_this.tmp
`
	createStashIgnoreFile(t, tmpDir, stashignore)

	// Create stashignore filter with library root.
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Create scanner.
	scanner := &file.Scanner{
		ScanFilters: []file.PathFilter{stashIgnoreFilter},
	}

	testScenarios := []struct {
		path     string
		accepted bool
	}{
		{filepath.Join(tmpDir, "file1.tmp"), false},
		{filepath.Join(tmpDir, "file2.tmp"), false},
		{filepath.Join(tmpDir, "keep_this.tmp"), true},
		{filepath.Join(tmpDir, "video.mp4"), true},
	}

	ctx := context.Background()

	for _, scenario := range testScenarios {
		info, err := os.Stat(scenario.path)
		if err != nil {
			t.Fatalf("failed to stat file %s: %v", scenario.path, err)
		}
		accepted := scanner.AcceptEntry(ctx, scenario.path, info, "")

		if accepted != scenario.accepted {
			t.Errorf("unexpected accept result for %s: expected %v, got %v",
				scenario.path, scenario.accepted, accepted)
		}
	}
}
