package config

import (
	"path/filepath"
	"testing"
)

func TestStashConfigsGetStashFromDirPathReturnsMostSpecificPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "library")
	nested := filepath.Join(root, "images")

	parent := &StashConfig{
		Path:         root,
		ExcludeImage: true,
	}
	child := &StashConfig{
		Path:         nested,
		ExcludeVideo: true,
	}

	stashes := StashConfigs{parent, child}

	got := stashes.GetStashFromDirPath(filepath.Join(nested, "set"))
	if got != child {
		t.Fatalf("expected nested stash config, got %#v", got)
	}
}

func TestStashConfigsGetStashFromDirPathReturnsMostSpecificPathRegardlessOfOrder(t *testing.T) {
	root := filepath.Join(t.TempDir(), "library")
	nested := filepath.Join(root, "images")

	parent := &StashConfig{
		Path:         root,
		ExcludeImage: true,
	}
	child := &StashConfig{
		Path:         nested,
		ExcludeVideo: true,
	}

	stashes := StashConfigs{child, parent}

	got := stashes.GetStashFromDirPath(filepath.Join(nested, "set"))
	if got != child {
		t.Fatalf("expected nested stash config, got %#v", got)
	}
}

func TestStashConfigsGetStashFromPathReturnsMostSpecificPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "library")
	nested := filepath.Join(root, "images")

	parent := &StashConfig{
		Path:         root,
		ExcludeImage: true,
	}
	child := &StashConfig{
		Path:         nested,
		ExcludeVideo: true,
	}

	stashes := StashConfigs{parent, child}

	got := stashes.GetStashFromPath(filepath.Join(nested, "image.jpg"))
	if got != child {
		t.Fatalf("expected nested stash config, got %#v", got)
	}
}

func TestStashConfigsGetStashFromDirPathReturnsNilOutsideLibraries(t *testing.T) {
	root := filepath.Join(t.TempDir(), "library")
	outside := filepath.Join(t.TempDir(), "outside")

	stashes := StashConfigs{{Path: root}}

	got := stashes.GetStashFromDirPath(outside)
	if got != nil {
		t.Fatalf("expected nil stash config, got %#v", got)
	}
}
