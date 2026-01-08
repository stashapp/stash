package file

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// Helper to create an empty file.
func createTestFile(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
}

// Helper to create a file with content.
func createTestFileWithContent(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
}

// Helper to create a directory.
func createTestDir(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", path, err)
	}
}

// walkAndFilter walks the directory tree and returns paths accepted by the filter.
// Returns paths relative to root for easier assertion.
func walkAndFilter(t *testing.T, root string, filter *StashIgnoreFilter) []string {
	t.Helper()
	var accepted []string
	ctx := context.Background()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself.
		if path == root {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if filter.Accept(ctx, path, info) {
			relPath, _ := filepath.Rel(root, path)
			accepted = append(accepted, relPath)
		} else if info.IsDir() {
			// If directory is rejected, skip it.
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		t.Fatalf("walk failed: %v", err)
	}

	sort.Strings(accepted)
	return accepted
}

// assertPathsEqual checks that the accepted paths match expected.
func assertPathsEqual(t *testing.T, expected, actual []string) {
	t.Helper()
	sort.Strings(expected)

	if len(expected) != len(actual) {
		t.Errorf("path count mismatch:\nexpected %d: %v\nactual %d: %v", len(expected), expected, len(actual), actual)
		return
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("path mismatch at index %d:\nexpected: %s\nactual: %s", i, expected[i], actual[i])
		}
	}
}

func TestStashIgnore_ExactFilename(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "video2.mp4")
	createTestFile(t, tmpDir, "ignore_me.mp4")

	// Create .stashignore that excludes exact filename.
	createTestFileWithContent(t, tmpDir, ".stashignore", "ignore_me.mp4\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"video1.mp4",
		"video2.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_WildcardPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "video2.mp4")
	createTestFile(t, tmpDir, "temp1.tmp")
	createTestFile(t, tmpDir, "temp2.tmp")
	createTestFile(t, tmpDir, "notes.log")

	// Create .stashignore that excludes by extension.
	createTestFileWithContent(t, tmpDir, ".stashignore", "*.tmp\n*.log\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"video1.mp4",
		"video2.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_DirectoryExclusion(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestDir(t, tmpDir, "excluded_dir")
	createTestFile(t, tmpDir, "excluded_dir/video2.mp4")
	createTestFile(t, tmpDir, "excluded_dir/video3.mp4")
	createTestDir(t, tmpDir, "included_dir")
	createTestFile(t, tmpDir, "included_dir/video4.mp4")

	// Create .stashignore that excludes a directory.
	createTestFileWithContent(t, tmpDir, ".stashignore", "excluded_dir/\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"included_dir",
		"included_dir/video4.mp4",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_NegationPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "file1.tmp")
	createTestFile(t, tmpDir, "file2.tmp")
	createTestFile(t, tmpDir, "keep_this.tmp")

	// Create .stashignore that excludes *.tmp but keeps one.
	createTestFileWithContent(t, tmpDir, ".stashignore", "*.tmp\n!keep_this.tmp\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"keep_this.tmp",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_CommentsAndEmptyLines(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "ignore_me.mp4")

	// Create .stashignore with comments and empty lines.
	stashignore := `# This is a comment
ignore_me.mp4

# Another comment

`
	createTestFileWithContent(t, tmpDir, ".stashignore", stashignore)

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_NestedStashIgnoreFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "root_video.mp4")
	createTestFile(t, tmpDir, "root_ignore.tmp")
	createTestDir(t, tmpDir, "subdir")
	createTestFile(t, tmpDir, "subdir/sub_video.mp4")
	createTestFile(t, tmpDir, "subdir/sub_ignore.log")
	createTestFile(t, tmpDir, "subdir/also_tmp.tmp")

	// Root .stashignore excludes *.tmp.
	createTestFileWithContent(t, tmpDir, ".stashignore", "*.tmp\n")

	// Subdir .stashignore excludes *.log.
	createTestFileWithContent(t, tmpDir, "subdir/.stashignore", "*.log\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	// *.tmp from root should apply everywhere.
	// *.log from subdir should only apply in subdir.
	expected := []string{
		".stashignore",
		"root_video.mp4",
		"subdir",
		"subdir/.stashignore",
		"subdir/sub_video.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_PathPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestDir(t, tmpDir, "subdir")
	createTestFile(t, tmpDir, "subdir/video2.mp4")
	createTestFile(t, tmpDir, "subdir/skip_this.mp4")

	// Create .stashignore that excludes a specific path.
	createTestFileWithContent(t, tmpDir, ".stashignore", "subdir/skip_this.mp4\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"subdir",
		"subdir/video2.mp4",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_DoubleStarPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestDir(t, tmpDir, "a")
	createTestFile(t, tmpDir, "a/video2.mp4")
	createTestDir(t, tmpDir, "a/temp")
	createTestFile(t, tmpDir, "a/temp/video3.mp4")
	createTestDir(t, tmpDir, "a/b")
	createTestDir(t, tmpDir, "a/b/temp")
	createTestFile(t, tmpDir, "a/b/temp/video4.mp4")

	// Create .stashignore that excludes temp directories at any level.
	createTestFileWithContent(t, tmpDir, ".stashignore", "**/temp/\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"a",
		"a/b",
		"a/video2.mp4",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_LeadingSlashPattern(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "ignore.mp4")
	createTestDir(t, tmpDir, "subdir")
	createTestFile(t, tmpDir, "subdir/ignore.mp4")

	// Create .stashignore that excludes only at root level.
	createTestFileWithContent(t, tmpDir, ".stashignore", "/ignore.mp4\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	// Only root ignore.mp4 should be excluded.
	expected := []string{
		".stashignore",
		"subdir",
		"subdir/ignore.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_NoStashIgnoreFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files without any .stashignore.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "video2.mp4")
	createTestDir(t, tmpDir, "subdir")
	createTestFile(t, tmpDir, "subdir/video3.mp4")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	// All files should be accepted.
	expected := []string{
		"subdir",
		"subdir/video3.mp4",
		"video1.mp4",
		"video2.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_HiddenDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files including hidden directory.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestDir(t, tmpDir, ".hidden")
	createTestFile(t, tmpDir, ".hidden/video2.mp4")

	// Create .stashignore that excludes hidden directories.
	createTestFileWithContent(t, tmpDir, ".stashignore", ".*\n!.stashignore\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_MultiplePatternsSameLine(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "file.tmp")
	createTestFile(t, tmpDir, "file.log")
	createTestFile(t, tmpDir, "file.bak")

	// Each pattern should be on its own line.
	createTestFileWithContent(t, tmpDir, ".stashignore", "*.tmp\n*.log\n*.bak\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_TrailingSpaces(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "ignore_me.mp4")

	// Pattern with trailing spaces (should be trimmed).
	createTestFileWithContent(t, tmpDir, ".stashignore", "ignore_me.mp4   \n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_EscapedHash(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "#filename.mp4")

	// Escaped hash should match literal # character.
	createTestFileWithContent(t, tmpDir, ".stashignore", "\\#filename.mp4\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"video1.mp4",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_CaseSensitiveMatching(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files - use distinct names that work on all filesystems.
	createTestFile(t, tmpDir, "video_lower.mp4")
	createTestFile(t, tmpDir, "VIDEO_UPPER.mp4")
	createTestFile(t, tmpDir, "other.avi")

	// Pattern should match exactly (case-sensitive).
	createTestFileWithContent(t, tmpDir, ".stashignore", "video_lower.mp4\n")

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	// Only exact match is excluded.
	expected := []string{
		".stashignore",
		"VIDEO_UPPER.mp4",
		"other.avi",
	}

	assertPathsEqual(t, expected, accepted)
}

func TestStashIgnore_ComplexScenario(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a complex directory structure.
	createTestFile(t, tmpDir, "video1.mp4")
	createTestFile(t, tmpDir, "video2.avi")
	createTestFile(t, tmpDir, "thumbnail.jpg")
	createTestFile(t, tmpDir, "metadata.nfo")
	createTestDir(t, tmpDir, "movies")
	createTestFile(t, tmpDir, "movies/movie1.mp4")
	createTestFile(t, tmpDir, "movies/movie1.nfo")
	createTestDir(t, tmpDir, "movies/.thumbnails")
	createTestFile(t, tmpDir, "movies/.thumbnails/thumb1.jpg")
	createTestDir(t, tmpDir, "temp")
	createTestFile(t, tmpDir, "temp/processing.mp4")
	createTestDir(t, tmpDir, "backup")
	createTestFile(t, tmpDir, "backup/video1.mp4.bak")

	// Complex .stashignore.
	stashignore := `# Ignore metadata files
*.nfo

# Ignore hidden directories
.*
!.stashignore

# Ignore temp and backup directories
temp/
backup/

# But keep thumbnails in specific location
!movies/.thumbnails/
`
	createTestFileWithContent(t, tmpDir, ".stashignore", stashignore)

	filter := NewStashIgnoreFilter(tmpDir)
	accepted := walkAndFilter(t, tmpDir, filter)

	expected := []string{
		".stashignore",
		"movies",
		"movies/.thumbnails",
		"movies/.thumbnails/thumb1.jpg",
		"movies/movie1.mp4",
		"thumbnail.jpg",
		"video1.mp4",
		"video2.avi",
	}

	assertPathsEqual(t, expected, accepted)
}
