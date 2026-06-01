package internal

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg" // ensure JPEG decoder is registered for image.Decode
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Sprite parameters matching pkg/scene/generate/sprite.go.
const (
	SpriteScreenshotWidth = 160  // per-cell pixel width; sprite_screenshot_width default
	SpriteIntervalSec     = 30.0 // sprite_interval default; one frame per ~30s
	SpriteMinFrames       = 10
	SpriteMaxFrames       = 500
	// SpriteJPEGQuality — for the FINAL tiled output (Go's image/jpeg quality, 1-100).
	// Higher than ffmpeg's -q:v 5 (which is libavcodec's 1-31 scale, very high quality);
	// 82 is a reasonable visual match for stash's libvips output.
	SpriteJPEGQuality = 82

	// spriteFrameConcurrency caps how many per-frame ffmpeg processes run at
	// once for a single scene. Higher = faster per scene but more NAS read
	// pressure. 4 is a balance — see EXTERNAL-WORKERS.md.
	spriteFrameConcurrency = 4
)

// SpritePlan captures the layout decisions for one scene's sprite + VTT.
type SpritePlan struct {
	FrameCount  int     // N frames sampled across the duration
	IntervalSec float64 // duration / N
	GridSize    int     // ceil(sqrt(N))
}

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
	return SpritePlan{
		FrameCount:  n,
		IntervalSec: duration / float64(n),
		GridSize:    int(math.Ceil(math.Sqrt(float64(n)))),
	}
}

// GenerateSprite produces the scrubber sprite + WebVTT pair for one scene.
//
// The fast path: for each of N timestamps, run an ffmpeg "fast seek" (-ss
// BEFORE -i, plus NVDEC for hardware-accelerated keyframe seek and decode)
// returning ONE scaled JPEG via stdout. Then tile them all in Go using
// image/draw. This avoids the full-pass decode that the older fps-filter
// pipeline forced.
//
// Atomic via the worker's tmp dir, like the preview pipeline.
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

	// Extract frames concurrently (bounded).
	frames, err := e.extractSpriteFrames(ctx, source, plan, probe.Duration)
	if err != nil {
		return fmt.Errorf("extract frames: %w", err)
	}
	if e.DryRun {
		// Still emit the VTT plan so a dry-run shows the layout.
		_ = e.writeVTT(plan, 0, 0, vttPath, partialVTT, stem)
		return nil
	}

	// Tile in Go: decode each JPEG, draw into a single canvas.
	cellW, cellH, err := tileFramesToJPEG(frames, plan.GridSize, SpriteJPEGQuality, partialSprite)
	if err != nil {
		return fmt.Errorf("tile sprite: %w", err)
	}

	// VTT references the sprite by sibling filename; cells use the dims we just
	// computed (matches pkg/scene/generate/sprite.go's gridSize math).
	if err := e.writeVTT(plan, cellW, cellH, vttPath, partialVTT, stem); err != nil {
		return fmt.Errorf("write VTT: %w", err)
	}

	if err := os.Rename(partialSprite, spritePath); err != nil {
		_ = os.Remove(partialSprite)
		_ = os.Remove(partialVTT)
		return fmt.Errorf("rename sprite to %s: %w", spritePath, err)
	}
	if err := os.Rename(partialVTT, vttPath); err != nil {
		_ = os.Remove(partialVTT)
		return fmt.Errorf("rename VTT to %s: %w", vttPath, err)
	}
	return nil
}

// extractSpriteFrames runs N concurrent ffmpeg invocations (bounded by
// spriteFrameConcurrency), each seeking to its timestamp and emitting one
// scaled JPEG via stdout. Returns the JPEG bytes in frame-order.
func (e *Encoder) extractSpriteFrames(
	ctx context.Context, source string, plan SpritePlan, duration float64,
) ([][]byte, error) {
	frames := make([][]byte, plan.FrameCount)
	errs := make([]error, plan.FrameCount)

	sem := make(chan struct{}, spriteFrameConcurrency)
	var wg sync.WaitGroup
	for i := 0; i < plan.FrameCount; i++ {
		i := i
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			t := float64(i) * plan.IntervalSec
			// Clamp the last timestamp slightly below duration so we don't seek
			// past EOF (some files won't decode the very last frame).
			if t >= duration {
				t = duration * 0.999
			}
			b, err := e.extractOneSpriteFrame(ctx, source, t)
			if err != nil {
				errs[i] = err
				return
			}
			frames[i] = b
		}()
	}
	wg.Wait()

	// Fail only if EVERY frame failed — for partial failures we keep going
	// with blank cells. (Rare, but better than aborting a 100-frame sprite
	// because of one weird timestamp.)
	allFailed := true
	for _, err := range errs {
		if err == nil {
			allFailed = false
			break
		}
	}
	if allFailed {
		return nil, fmt.Errorf("all %d frame extractions failed; first error: %v", plan.FrameCount, errs[0])
	}
	return frames, nil
}

// extractOneSpriteFrame runs ffmpeg to produce ONE JPEG via stdout. Uses fast
// seek (`-ss` before `-i`) and NVDEC where available for keyframe decode speed.
func (e *Encoder) extractOneSpriteFrame(ctx context.Context, source string, t float64) ([]byte, error) {
	args := []string{
		"-loglevel", "error", "-y",
		"-hwaccel", "cuda", "-hwaccel_output_format", "cuda",
		"-ss", strconv.FormatFloat(t, 'f', 3, 64),
		"-i", source,
		"-frames:v", "1",
		// Drop back to CPU memory for scaling — scale_cuda isn't in every
		// BtbN build; staying with libswscale keeps the pipeline portable.
		"-vf", "hwdownload,format=nv12,scale=" + strconv.Itoa(SpriteScreenshotWidth) + ":-2",
		// libavcodec quality scale 2-31 (lower = better). 4 keeps thumbnails
		// crisp enough for scrubber use without ballooning JPEG sizes.
		"-q:v", "4",
		"-f", "mjpeg",
		"pipe:1",
	}
	cmd := exec.CommandContext(ctx, e.FFmpegPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg seek at %.2fs: %w: %s", t, err, truncate(stderr.String(), 200))
	}
	b := stdout.Bytes()
	if len(b) < 4 || b[0] != 0xFF || b[1] != 0xD8 {
		return nil, fmt.Errorf("seek at %.2fs produced no valid JPEG (%d bytes)", t, len(b))
	}
	// Copy out of the bytes.Buffer-backed slice so the buffer can be GC'd.
	out := make([]byte, len(b))
	copy(out, b)
	return out, nil
}

// tileFramesToJPEG decodes each frame's JPEG, lays them out in a gridSize×gridSize
// canvas, and writes the result as a single JPEG. Returns the per-cell pixel
// dimensions so the caller can build a matching VTT. Empty cells (when N <
// gridSize²) are left transparent / black — stash's VTT only references the
// N valid cues so this never shows in the UI.
func tileFramesToJPEG(frames [][]byte, gridSize, jpegQuality int, outPath string) (cellW, cellH int, err error) {
	if len(frames) == 0 || gridSize == 0 {
		return 0, 0, fmt.Errorf("no frames or zero grid size")
	}

	// First successful frame defines the per-cell dimensions; all frames are
	// the same size from ffmpeg's scale=160:-2 invocation.
	var sampleCellW, sampleCellH int
	for _, f := range frames {
		if len(f) == 0 {
			continue
		}
		img, err := jpeg.Decode(bytes.NewReader(f))
		if err != nil {
			continue
		}
		b := img.Bounds()
		sampleCellW = b.Dx()
		sampleCellH = b.Dy()
		break
	}
	if sampleCellW == 0 || sampleCellH == 0 {
		return 0, 0, fmt.Errorf("no decodable frames")
	}

	canvas := image.NewRGBA(image.Rect(0, 0, sampleCellW*gridSize, sampleCellH*gridSize))

	for i, f := range frames {
		if len(f) == 0 {
			continue // blank cell; the VTT will skip it
		}
		img, err := jpeg.Decode(bytes.NewReader(f))
		if err != nil {
			continue
		}
		col := i % gridSize
		row := i / gridSize
		dest := image.Rect(col*sampleCellW, row*sampleCellH, (col+1)*sampleCellW, (row+1)*sampleCellH)
		draw.Draw(canvas, dest, img, image.Point{}, draw.Src)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return 0, 0, err
	}
	defer out.Close()
	if err := jpeg.Encode(out, canvas, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return 0, 0, err
	}
	return sampleCellW, sampleCellH, nil
}

// writeVTT emits a WebVTT cue file. Each cue maps a time range to a region of
// the sprite image. See pkg/scene/generate/sprite.go.
func (e *Encoder) writeVTT(plan SpritePlan, cellW, cellH int, finalPath, partialPath, stem string) error {
	var sb strings.Builder
	sb.WriteString("WEBVTT\n\n")
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

// SpriteExists reports whether the sprite + VTT pair both exist + are non-empty.
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
