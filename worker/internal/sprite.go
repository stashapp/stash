package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Sprite parameters matching pkg/scene/generate/sprite.go.
const (
	// SpriteScreenshotWidth is the per-cell pixel width before grid tiling
	// (sprite_screenshot_width default, internal/manager/config/config.go).
	SpriteScreenshotWidth = 160
	// SpriteIntervalSec is the target spacing between frames (sprite_interval
	// default). Mirroring config.go.
	SpriteIntervalSec = 30.0
	// SpriteMinFrames and SpriteMaxFrames clamp the auto-computed frame count.
	SpriteMinFrames = 10
	SpriteMaxFrames = 500
	// SpriteJPEGQuality matches stash's libvips JPEG quality for sprites.
	SpriteJPEGQuality = 5
)

// SpritePlan captures the layout decisions for one scene's sprite + VTT pair.
// gridSize, cellW, cellH are filled in after probing the rendered sprite (we
// don't pre-compute cell height because it depends on source aspect ratio,
// which ffmpeg's scale filter resolves at decode time).
type SpritePlan struct {
	FrameCount int     // N frames sampled across the duration
	IntervalSec float64 // duration / N
	GridSize   int     // ceil(sqrt(N))
}

// planSprite computes the sample count + grid size, matching sprite.go's
// GetSpriteGridSize and the count clamping (config defaults applied here).
func planSprite(duration float64) SpritePlan {
	if duration <= 0 {
		return SpritePlan{}
	}
	n := int(math.Ceil(duration / SpriteIntervalSec))
	if n < SpriteMinFrames {
		n = SpriteMinFrames
	}
	if n > SpriteMaxFrames {
		n = SpriteMaxFrames
	}
	grid := int(math.Ceil(math.Sqrt(float64(n))))
	return SpritePlan{
		FrameCount:  n,
		IntervalSec: duration / float64(n),
		GridSize:    grid,
	}
}

// GenerateSprite extracts N frames evenly distributed across the source video
// and tiles them into a single JPEG. Outputs both the sprite image and a
// WebVTT cue file. Atomic via tmpDir + rename, just like the preview pipeline.
// Per-cell dimensions are derived AFTER ffmpeg renders, by probing the result
// (cell height depends on source aspect ratio — see EXTERNAL-WORKERS.md §5).
func (e *Encoder) GenerateSprite(
	ctx context.Context,
	source, spritePath, vttPath, tmpDir string,
) error {
	probe, err := e.probeVideo(ctx, source)
	if err != nil {
		return fmt.Errorf("probe %s: %w", source, err)
	}
	plan := planSprite(probe.Duration)
	if plan.FrameCount == 0 {
		return fmt.Errorf("invalid duration %.2f", probe.Duration)
	}

	stem := strings.TrimSuffix(filepath.Base(spritePath), filepath.Ext(spritePath))
	partialSprite := filepath.Join(tmpDir, stem+".jpg.partial")
	partialVTT := filepath.Join(tmpDir, stem+".vtt.partial")
	if !e.DryRun {
		if err := os.MkdirAll(tmpDir, 0o755); err != nil {
			return fmt.Errorf("mkdir tmp: %w", err)
		}
	}

	// One ffmpeg invocation does the whole thing:
	//   fps=N/duration → sample N frames uniformly (no exclusion windows; matches
	//                    stash's sprite gen which uses full duration)
	//   scale=160:-2  → width=160, height auto, even number for codec safety
	//   tile=GxG:padding=0 → arrange into the grid
	// -frames:v 1 emits just the single tiled output.
	fps := float64(plan.FrameCount) / probe.Duration
	vf := fmt.Sprintf("fps=%.6f,scale=%d:-2,tile=%dx%d:padding=0",
		fps, SpriteScreenshotWidth, plan.GridSize, plan.GridSize)
	args := []string{
		"-loglevel", "error", "-y",
		"-i", source,
		"-vf", vf,
		"-frames:v", "1",
		"-q:v", strconv.Itoa(SpriteJPEGQuality),
		"-f", "mjpeg",
		partialSprite,
	}
	if err := e.runFFmpeg(ctx, args); err != nil {
		return fmt.Errorf("render sprite: %w", err)
	}

	// Probe the rendered sprite for actual dimensions, then derive per-cell
	// size. gridSize × cell == sprite, and ffmpeg's tile filter pads or
	// truncates to fit when frame count < grid². The VTT only emits N cues.
	if e.DryRun {
		// Skip post-render steps in dry-run; we still emit the VTT plan so
		// the user can sanity-check.
		_ = e.writeVTT(plan, 0, 0, vttPath, partialVTT, stem)
		return nil
	}
	spriteW, spriteH, err := e.probeImage(ctx, partialSprite)
	if err != nil {
		return fmt.Errorf("probe sprite: %w", err)
	}
	cellW := spriteW / plan.GridSize
	cellH := spriteH / plan.GridSize
	if cellW == 0 || cellH == 0 {
		return fmt.Errorf("computed zero cell size: spriteW=%d spriteH=%d grid=%d",
			spriteW, spriteH, plan.GridSize)
	}
	if err := e.writeVTT(plan, cellW, cellH, vttPath, partialVTT, stem); err != nil {
		return fmt.Errorf("write VTT: %w", err)
	}
	if err := os.Rename(partialSprite, spritePath); err != nil {
		_ = os.Remove(partialSprite)
		_ = os.Remove(partialVTT)
		return fmt.Errorf("rename sprite to %s: %w", spritePath, err)
	}
	if err := os.Rename(partialVTT, vttPath); err != nil {
		_ = os.Remove(vttPath)
		_ = os.Remove(partialVTT)
		return fmt.Errorf("rename VTT to %s: %w", vttPath, err)
	}
	return nil
}

// writeVTT emits a WebVTT cue file. Each cue maps a time range to a region of
// the sprite image. Stash's player overlays the region while the user scrubs.
// File path "<stem>_sprite.jpg" in the cue is a reference to a SIBLING file in
// the vtt/ directory; stash resolves it relative to the VTT file's location.
func (e *Encoder) writeVTT(plan SpritePlan, cellW, cellH int, finalPath, partialPath, stem string) error {
	var sb strings.Builder
	sb.WriteString("WEBVTT\n\n")
	// strip "_thumbs" suffix from stem if present, then re-form the sprite filename
	spriteName := stem
	if !strings.HasSuffix(spriteName, "_sprite") {
		spriteName += "_sprite"
	}
	spriteName += ".jpg"

	for i := 0; i < plan.FrameCount; i++ {
		t0 := float64(i) * plan.IntervalSec
		t1 := t0 + plan.IntervalSec
		col := i % plan.GridSize
		row := i / plan.GridSize
		fmt.Fprintf(&sb,
			"%s --> %s\n%s#xywh=%d,%d,%d,%d\n\n",
			formatVTT(t0), formatVTT(t1),
			spriteName, col*cellW, row*cellH, cellW, cellH,
		)
	}
	if e.DryRun {
		return nil
	}
	return os.WriteFile(partialPath, []byte(sb.String()), 0o644)
}

// formatVTT renders a duration as "HH:MM:SS.mmm" (WebVTT cue timestamp).
func formatVTT(s float64) string {
	h := int(s) / 3600
	m := (int(s) / 60) % 60
	sec := int(s) % 60
	ms := int((s - math.Floor(s)) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, sec, ms)
}

// probeImage returns the rendered image's dimensions via ffprobe.
func (e *Encoder) probeImage(ctx context.Context, path string) (int, int, error) {
	cmd := exec.CommandContext(ctx, e.FFprobe,
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		"-select_streams", "v:0",
		path,
	)
	cmd.Env = append(cmd.Env, "LC_ALL=C")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, 0, fmt.Errorf("ffprobe: %w: %s", err, truncate(stderr.String(), 200))
	}
	var data struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &data); err != nil {
		return 0, 0, fmt.Errorf("decode ffprobe: %w", err)
	}
	if len(data.Streams) == 0 {
		return 0, 0, fmt.Errorf("no video stream in %s", path)
	}
	return data.Streams[0].Width, data.Streams[0].Height, nil
}

// SpriteExists reports whether the sprite + VTT pair both exist for this
// checksum. Stash's player uses the VTT to discover the sprite, so we need
// both before stash considers the scene complete.
func SpriteExists(gen GeneratedPaths, checksum string) (bool, error) {
	for _, p := range []string{gen.SpriteImage(checksum), gen.SpriteVTT(checksum)} {
		st, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		if st.Size() == 0 {
			return false, nil
		}
	}
	return true, nil
}
