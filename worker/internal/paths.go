package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

// PrefixRewriter translates stash-side paths to worker-side paths.
//
// Example: stash sees a media file at "/data/completed/foo.mp4" because it's
// mounted at /data in the stash container. The worker accesses the same file
// over SMB at "\\overwatch-stash\torrents\completed\foo.mp4". The rewriter
// strips "/data" and prepends "\\overwatch-stash\torrents".
type PrefixRewriter struct {
	from string // stash-side prefix (e.g. "/data")
	to   string // worker-side prefix (e.g. \\overwatch-stash\torrents)
}

// ParsePrefixRewriter parses a CLI flag value "stash=worker" (or "stash:worker")
// into a PrefixRewriter. Returns an error if the format is wrong.
func ParsePrefixRewriter(s string) (*PrefixRewriter, error) {
	// Prefer "=" as the separator; ":" is ambiguous with Windows drive letters.
	idx := strings.Index(s, "=")
	if idx < 1 || idx == len(s)-1 {
		return nil, fmt.Errorf("expected STASH_PREFIX=WORKER_PREFIX, got %q", s)
	}
	from := strings.TrimSpace(s[:idx])
	to := strings.TrimSpace(s[idx+1:])
	return &PrefixRewriter{
		from: strings.TrimRight(from, "/\\"),
		to:   strings.TrimRight(to, "/\\"),
	}, nil
}

// To returns the worker-side prefix (the right-hand side of the CLI value). Useful
// when the worker needs to know its own view of a stash-managed root directory
// (e.g. the generated/ dir) without translating an actual stash-side path.
func (r *PrefixRewriter) To() string {
	if r == nil {
		return ""
	}
	return r.to
}

// Rewrite returns the worker-side equivalent of a stash-side path, or the input
// unchanged if it doesn't match the prefix.
func (r *PrefixRewriter) Rewrite(stashPath string) string {
	if r == nil || r.from == "" {
		return stashPath
	}
	// Match the from-prefix at the start, with a directory separator after.
	if stashPath == r.from {
		return r.to
	}
	if strings.HasPrefix(stashPath, r.from+"/") || strings.HasPrefix(stashPath, r.from+"\\") {
		rest := stashPath[len(r.from):]
		// Convert separators to whatever the worker side wants (heuristic: if
		// the to-prefix uses backslashes, use those; else forward slashes).
		if strings.Contains(r.to, "\\") {
			rest = strings.ReplaceAll(rest, "/", "\\")
		}
		return r.to + rest
	}
	return stashPath
}

// GeneratedPaths derives the on-disk locations of generated artifacts for a
// scene, mirroring pkg/models/paths/paths_scenes.go in the stash source.
type GeneratedPaths struct {
	// Root is the WORKER-side path to stash's generated/ dir
	// (e.g. \\overwatch-stash\generated).
	Root string
}

// VideoPreview returns "<Root>/screenshots/<checksum>.mp4".
// See paths_scenes.go:38 (GetVideoPreviewPath).
func (g GeneratedPaths) VideoPreview(checksum string) string {
	return filepath.Join(g.Root, "screenshots", checksum+".mp4")
}

// WebPPreview returns "<Root>/screenshots/<checksum>.webp".
// See paths_scenes.go:42 (GetWebpPreviewPath).
func (g GeneratedPaths) WebPPreview(checksum string) string {
	return filepath.Join(g.Root, "screenshots", checksum+".webp")
}

// SpriteImage returns "<Root>/vtt/<checksum>_sprite.jpg".
// See paths_scenes.go:46 (GetSpriteImageFilePath).
func (g GeneratedPaths) SpriteImage(checksum string) string {
	return filepath.Join(g.Root, "vtt", checksum+"_sprite.jpg")
}

// SpriteVTT returns "<Root>/vtt/<checksum>_thumbs.vtt".
// See paths_scenes.go:50 (GetSpriteVttFilePath).
func (g GeneratedPaths) SpriteVTT(checksum string) string {
	return filepath.Join(g.Root, "vtt", checksum+"_thumbs.vtt")
}

// TmpDir returns "<Root>/tmp", stash's own scratch dir (paths_generated.go:34).
// Worker writes ".partial" files here, then atomically renames to the
// destination. Using a different directory than the destination avoids stash
// briefly seeing a zero-byte file via path-existence checks.
func (g GeneratedPaths) TmpDir() string {
	return filepath.Join(g.Root, "tmp")
}
