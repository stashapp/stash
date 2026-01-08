package file

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ignore "github.com/sabhiram/go-gitignore"
)

const stashIgnoreFilename = ".stashignore"

// StashIgnoreFilter implements PathFilter to exclude files/directories
// based on .stashignore files with gitignore-style patterns.
type StashIgnoreFilter struct {
	// root is the root directory being scanned.
	root string

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

// NewStashIgnoreFilter creates a new StashIgnoreFilter for the given root directory.
func NewStashIgnoreFilter(root string) *StashIgnoreFilter {
	return &StashIgnoreFilter{
		root: root,
	}
}

// Accept returns true if the path should be included in the scan.
// It checks for .stashignore files in the directory hierarchy and
// applies gitignore-style pattern matching.
func (f *StashIgnoreFilter) Accept(ctx context.Context, path string, info fs.FileInfo) bool {
	// Always accept .stashignore files themselves so they can be read.
	if filepath.Base(path) == stashIgnoreFilename {
		return true
	}

	// Get the directory containing this path.
	var dir string
	if info.IsDir() {
		dir = filepath.Dir(path)
	} else {
		dir = filepath.Dir(path)
	}

	// Collect all applicable ignore entries from root to this directory.
	entries := f.collectIgnoreEntries(dir)

	// Check if any pattern matches (and isn't negated).
	relPath, err := filepath.Rel(f.root, path)
	if err != nil {
		// If we can't get relative path, accept the file.
		return true
	}

	// Normalise to forward slashes for consistent matching.
	relPath = filepath.ToSlash(relPath)

	// For directories, also check with trailing slash.
	if info.IsDir() {
		relPath = relPath + "/"
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
		// Check negation by testing without the directory suffix.
		// The library handles negation internally.
	}

	return !ignored
}

// collectIgnoreEntries gathers all ignore entries from root to the given directory.
func (f *StashIgnoreFilter) collectIgnoreEntries(dir string) []*ignoreEntry {
	var entries []*ignoreEntry

	// Walk from root to current directory.
	current := f.root
	relDir, err := filepath.Rel(f.root, dir)
	if err != nil {
		return entries
	}

	// Check root directory first.
	if entry := f.getOrLoadIgnoreEntry(current); entry != nil {
		entries = append(entries, entry)
	}

	// Then check each subdirectory.
	if relDir != "." {
		parts := strings.Split(filepath.ToSlash(relDir), "/")
		for _, part := range parts {
			current = filepath.Join(current, part)
			if entry := f.getOrLoadIgnoreEntry(current); entry != nil {
				entries = append(entries, entry)
			}
		}
	}

	return entries
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
	patterns, err := f.loadIgnoreFile(stashIgnorePath, dir)
	if err != nil {
		// Cache negative result.
		f.cache.Store(dir, &ignoreEntry{patterns: nil, dir: dir})
		return nil
	}

	entry := &ignoreEntry{
		patterns: patterns,
		dir:      dir,
	}
	f.cache.Store(dir, entry)
	return entry
}

// loadIgnoreFile loads and compiles a .stashignore file.
func (f *StashIgnoreFilter) loadIgnoreFile(path string, dir string) (*ignore.GitIgnore, error) {
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
		return nil, os.ErrNotExist
	}

	return ignore.CompileIgnoreLines(patterns...), nil
}
