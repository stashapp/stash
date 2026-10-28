package manager

import (
	"path/filepath"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
)

// Synthetic absolute paths for shouldSkipDir tests (no relation to the host layout).
var (
	testLibraryRoot    = filepath.FromSlash("/stash-test/library/media")
	testExcludedPrefix = filepath.FromSlash("/stash-test/library/apps")
	testSiblingPrefix  = filepath.FromSlash("/stash-test/library/appsdata")
	testWideRoot       = filepath.FromSlash("/stash-test")
	testNestedLibrary  = filepath.FromSlash("/stash-test/library/apps/nested")
)

func makeStashConfigs(paths ...string) config.StashConfigs {
	cfgs := make(config.StashConfigs, len(paths))
	for i, p := range paths {
		cfgs[i] = &config.StashConfig{Path: p}
	}
	return cfgs
}

func TestShouldSkipDir_MatchesBothExcludes_NoDescendantRoot(t *testing.T) {
	dir := testExcludedPrefix
	s := &config.StashConfig{Path: testLibraryRoot}
	allPaths := makeStashConfigs(testLibraryRoot)
	videoRE := generateRegexps([]string{testExcludedPrefix})
	imageRE := generateRegexps([]string{testExcludedPrefix})

	if !shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=true for excluded dir with no descendant root")
	}
}

func TestShouldSkipDir_DirIsLibraryRoot(t *testing.T) {
	dir := testLibraryRoot
	s := &config.StashConfig{Path: testLibraryRoot}
	allPaths := makeStashConfigs(testLibraryRoot)
	videoRE := generateRegexps([]string{testLibraryRoot})
	imageRE := generateRegexps([]string{testLibraryRoot})

	if shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=false when dir is a library root")
	}
}

func TestShouldSkipDir_LibraryRootIsDescendant(t *testing.T) {
	dir := testExcludedPrefix
	s := &config.StashConfig{Path: testWideRoot}
	allPaths := makeStashConfigs(testNestedLibrary)
	videoRE := generateRegexps([]string{testExcludedPrefix})
	imageRE := generateRegexps([]string{testExcludedPrefix})

	if shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=false when a library root is beneath the excluded dir")
	}
}

func TestShouldSkipDir_SiblingMatchesPrefix(t *testing.T) {
	// Exclude pattern for .../apps also matches .../appsdata — expected regex prefix behavior.
	dir := testSiblingPrefix
	s := &config.StashConfig{Path: testLibraryRoot}
	allPaths := makeStashConfigs(testLibraryRoot)
	videoRE := generateRegexps([]string{testExcludedPrefix})
	imageRE := generateRegexps([]string{testExcludedPrefix})

	if !shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=true: prefix match on sibling is expected regex behavior")
	}
}

func TestShouldSkipDir_ExcludeImageFlag(t *testing.T) {
	dir := testExcludedPrefix
	s := &config.StashConfig{ExcludeImage: true}
	allPaths := makeStashConfigs(testLibraryRoot)
	videoRE := generateRegexps([]string{testExcludedPrefix})
	imageRE := generateRegexps(nil) // no image exclude patterns

	if !shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=true when ExcludeImage=true and video regex matches")
	}
}

func TestShouldSkipDir_OnlyVideoExcludeMatches(t *testing.T) {
	dir := testExcludedPrefix
	s := &config.StashConfig{ExcludeImage: false}
	allPaths := makeStashConfigs(testLibraryRoot)
	videoRE := generateRegexps([]string{testExcludedPrefix})
	imageRE := generateRegexps(nil) // no image exclude patterns

	if shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=false when only video exclude matches and ExcludeImage=false")
	}
}

func TestShouldSkipDir_EmptyExcludes(t *testing.T) {
	dir := testExcludedPrefix
	s := &config.StashConfig{}
	allPaths := makeStashConfigs(testLibraryRoot)
	videoRE := generateRegexps(nil)
	imageRE := generateRegexps(nil)

	if shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=false when no exclude patterns")
	}
}

func TestShouldSkipDir_TrailingSeparatorOnRoot(t *testing.T) {
	// allPaths entry has a trailing separator; filepath.Clean normalises it.
	dir := testExcludedPrefix
	s := &config.StashConfig{}
	allPaths := makeStashConfigs(testExcludedPrefix + string(filepath.Separator)) // trailing sep
	videoRE := generateRegexps([]string{testExcludedPrefix})
	imageRE := generateRegexps([]string{testExcludedPrefix})

	if shouldSkipDir(dir, s, allPaths, videoRE, imageRE) {
		t.Error("expected shouldSkipDir=false: trailing separator on root path should be cleaned")
	}
}
