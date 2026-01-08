package file

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/stashapp/stash/pkg/logger"
)

const stashIgnoreFilename = ".stashignore"

// StashIgnoreFilter implements PathFilter to exclude files/directories
// based on .stashignore files with gitignore-style patterns.
type StashIgnoreFilter struct {
	// cache stores compiled ignore patterns per directory.
	cache sync.Map // map[string]*ignoreEntry
}

// ignoreEntry holds the compiled ignore patterns for a directory.
type ignoreEntry struct {
	// patterns is the compiled gitignore matcher for this directory.
	patterns *ignore.GitIgnore
	// dir is the directory this entry applies to.
	dir string
}

// NewStashIgnoreFilter creates a new StashIgnoreFilter.
func NewStashIgnoreFilter() *StashIgnoreFilter {
	return &StashIgnoreFilter{}
}

// Accept returns true if the path should be included in the scan.
// It checks for .stashignore files in the directory hierarchy and
// applies gitignore-style pattern matching.
// The libraryRoot parameter bounds the search for .stashignore files -
// only directories within the library root are checked.
func (f *StashIgnoreFilter) Accept(ctx context.Context, path string, info fs.FileInfo, libraryRoot string) bool {
	// Always accept .stashignore files themselves so they can be read.
	if filepath.Base(path) == stashIgnoreFilename {
		return true
	}

	// If no library root provided, accept the file (safety fallback).
	if libraryRoot == "" {
		return true
	}

	// Get the directory containing this path.
	dir := filepath.Dir(path)

	// Collect all applicable ignore entries from library root to this directory.
	entries := f.collectIgnoreEntries(dir, libraryRoot)

	// If no .stashignore files found, accept the file.
	if len(entries) == 0 {
		return true
	}

	// Check each ignore entry in order (from root to most specific).
	// Later entries can override earlier ones with negation patterns.
	ignored := false
	for _, entry := range entries {
		// Get path relative to the ignore file's directory.
		entryRelPath, err := filepath.Rel(entry.dir, path)
		if err != nil {
			continue
		}
		entryRelPath = filepath.ToSlash(entryRelPath)
		if info.IsDir() {
			entryRelPath = entryRelPath + "/"
		}

		if entry.patterns.MatchesPath(entryRelPath) {
			ignored = true
		}
	}

	return !ignored
}

// collectIgnoreEntries gathers all ignore entries from library root to the given directory.
// It walks up the directory tree from dir to libraryRoot and returns entries in order
// from root to most specific.
func (f *StashIgnoreFilter) collectIgnoreEntries(dir string, libraryRoot string) []*ignoreEntry {
	// Collect directories from library root down to current dir.
	var dirs []string

	// Clean paths for consistent comparison.
	dir = filepath.Clean(dir)
	libraryRoot = filepath.Clean(libraryRoot)

	// Walk up from dir to library root, collecting directories.
	current := dir
	for {
		// Check if we're still within the library root.
		if !isPathInOrEqual(libraryRoot, current) {
			break
		}

		dirs = append([]string{current}, dirs...) // Prepend to maintain root-to-leaf order.

		// Stop if we've reached the library root.
		if current == libraryRoot {
			break
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root without finding library root.
			break
		}
		current = parent
	}

	// Check each directory for .stashignore files.
	var entries []*ignoreEntry
	for _, d := range dirs {
		if entry := f.getOrLoadIgnoreEntry(d); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries
}

// isPathInOrEqual checks if path is equal to or inside root.
func isPathInOrEqual(root, path string) bool {
	if path == root {
		return true
	}
	// Check if path starts with root + separator.
	return strings.HasPrefix(path, root+string(filepath.Separator))
}

// getOrLoadIgnoreEntry returns the cached ignore entry for a directory, or loads it.
func (f *StashIgnoreFilter) getOrLoadIgnoreEntry(dir string) *ignoreEntry {
	// Check cache first.
	if cached, ok := f.cache.Load(dir); ok {
		entry := cached.(*ignoreEntry)
		if entry.patterns == nil {
			return nil // Cached negative result.
		}
		return entry
	}

	// Try to load .stashignore from this directory.
	stashIgnorePath := filepath.Join(dir, stashIgnoreFilename)
	patterns, err := f.loadIgnoreFile(stashIgnorePath)
	if err != nil || patterns == nil {
		// Cache negative result (file doesn't exist or has no patterns).
		f.cache.Store(dir, &ignoreEntry{patterns: nil, dir: dir})
		return nil
	}

	logger.Debugf("Loaded .stashignore from %s", dir)

	entry := &ignoreEntry{
		patterns: patterns,
		dir:      dir,
	}
	f.cache.Store(dir, entry)
	return entry
}

// loadIgnoreFile loads and compiles a .stashignore file.
func (f *StashIgnoreFilter) loadIgnoreFile(path string) (*ignore.GitIgnore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var patterns []string

	for _, line := range lines {
		// Trim trailing whitespace (but preserve leading for patterns).
		line = strings.TrimRight(line, " \t\r")

		// Skip empty lines.
		if line == "" {
			continue
		}

		// Skip comments (but not escaped #).
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "\\#") {
			continue
		}

		patterns = append(patterns, line)
	}

	if len(patterns) == 0 {
		// File exists but has no patterns (e.g., only comments).
		return nil, nil
	}

	return ignore.CompileIgnoreLines(patterns...), nil
}
