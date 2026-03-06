package file

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/stashapp/stash/pkg/logger"
)

const stashIgnoreFilename = ".stashignore"

// entriesCacheSize is the size of the LRU cache for collected ignore entries.
// This cache stores the computed list of ignore entries per directory, avoiding
// repeated directory tree walks for files in the same directory.
const entriesCacheSize = 500

// StashIgnoreFilter implements PathFilter to exclude files/directories
// based on .stashignore files with gitignore-style patterns.
type StashIgnoreFilter struct {
	// cache stores compiled ignore patterns per directory.
	cache sync.Map // map[string]*ignoreEntry
	// entriesCache stores collected ignore entries per (dir, libraryRoot) pair.
	// This avoids recomputing the entry list for every file in the same directory.
	entriesCache *lru.Cache[string, []*ignoreEntry]
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
	// Create the LRU cache for collected entries.
	// Ignore error as it only fails if size <= 0.
	entriesCache, _ := lru.New[string, []*ignoreEntry](entriesCacheSize)
	return &StashIgnoreFilter{
		entriesCache: entriesCache,
	}
}

// Accept returns true if the path should be included in the scan.
// It checks for .stashignore files in the directory hierarchy and
// applies gitignore-style pattern matching.
// The libraryRoot parameter bounds the search for .stashignore files -
// only directories within the library root are checked.
func (f *StashIgnoreFilter) Accept(ctx context.Context, path string, info fs.FileInfo, libraryRoot string) bool {
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
			entryRelPath += "/"
		}

		if entry.patterns.MatchesPath(entryRelPath) {
			ignored = true
		}
	}

	return !ignored
}

// collectIgnoreEntries gathers all ignore entries from library root to the given directory.
// It walks up the directory tree from dir to libraryRoot and returns entries in order
// from root to most specific. Results are cached to avoid repeated computation for
// files in the same directory.
func (f *StashIgnoreFilter) collectIgnoreEntries(dir string, libraryRoot string) []*ignoreEntry {
	// Clean paths for consistent comparison and cache key generation.
	dir = filepath.Clean(dir)
	libraryRoot = filepath.Clean(libraryRoot)

	// Build cache key from dir and libraryRoot.
	cacheKey := dir + "\x00" + libraryRoot

	// Check the entries cache first.
	if cached, ok := f.entriesCache.Get(cacheKey); ok {
		return cached
	}

	// Try subdirectory shortcut: if parent's entries are cached, extend them.
	if dir != libraryRoot {
		parent := filepath.Dir(dir)
		if isPathInOrEqual(libraryRoot, parent) {
			parentKey := parent + "\x00" + libraryRoot
			if parentEntries, ok := f.entriesCache.Get(parentKey); ok {
				// Parent is cached - just check if current dir has a .stashignore.
				entries := parentEntries
				if entry := f.getOrLoadIgnoreEntry(dir); entry != nil {
					// Copy parent slice and append to avoid mutating cached slice.
					entries = make([]*ignoreEntry, len(parentEntries), len(parentEntries)+1)
					copy(entries, parentEntries)
					entries = append(entries, entry)
				}
				f.entriesCache.Add(cacheKey, entries)
				return entries
			}
		}
	}

	// No cache hit - compute from scratch.
	// Walk up from dir to library root, collecting directories.
	var dirs []string
	current := dir
	for {
		// Check if we're still within the library root.
		if !isPathInOrEqual(libraryRoot, current) {
			break
		}

		dirs = append(dirs, current)

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

	// Reverse to get root-to-leaf order.
	for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	}

	// Check each directory for .stashignore files.
	var entries []*ignoreEntry
	for _, d := range dirs {
		if entry := f.getOrLoadIgnoreEntry(d); entry != nil {
			entries = append(entries, entry)
		}
	}

	// Cache the result.
	f.entriesCache.Add(cacheKey, entries)

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
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Warnf("Failed to load .stashignore from %s: %v", dir, err)
		}
		f.cache.Store(dir, &ignoreEntry{patterns: nil, dir: dir})
		return nil
	}
	if patterns == nil {
		// File exists but has no patterns (empty or only comments).
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
