package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

// CoverScreenshotProportion is the fraction of total duration where stash
// extracts the scene cover (pkg/scene/generate/screenshot.go:17). Mirroring
// this gives us covers visually consistent with stash's own output.
const CoverScreenshotProportion = 0.2

// CoverJPEGQuality matches stash's screenshotQuality const (lower = better;
// ffmpeg -q:v range is 2..31, where 2 is near-lossless JPEG).
const CoverJPEGQuality = 2

// ExtractCover seeks to `atSeconds` and returns one JPEG frame as raw bytes.
// CPU decode is fine for single-frame work — NVDEC's CUDA context overhead
// would outweigh the per-frame savings.
func (e *Encoder) ExtractCover(ctx context.Context, source string, atSeconds float64) ([]byte, error) {
	if atSeconds < 0 {
		atSeconds = 0
	}
	args := []string{
		"-loglevel", "error",
		"-ss", strconv.FormatFloat(atSeconds, 'f', 3, 64),
		"-i", source,
		"-frames:v", "1",
		"-q:v", strconv.Itoa(CoverJPEGQuality),
		"-f", "mjpeg",
		"pipe:1",
	}
	if e.DryRun {
		fmt.Println(append([]string{e.FFmpegPath}, args...))
		// Return a 1-byte stub so callers can still build the upload path; the
		// caller's dry-run check skips the actual stash mutation.
		return []byte{0xFF}, nil
	}
	cmd := exec.CommandContext(ctx, e.FFmpegPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg cover: %w: %s", err, truncate(stderr.String(), 300))
	}
	if stdout.Len() == 0 {
		return nil, errors.New("ffmpeg produced no output")
	}
	// Defensive sanity: JPEG starts with the SOI marker 0xFFD8.
	b := stdout.Bytes()
	if len(b) < 2 || b[0] != 0xFF || b[1] != 0xD8 {
		return nil, fmt.Errorf("unexpected cover bytes (no JPEG SOI): %x...", b[:min(len(b), 4)])
	}
	return b, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

