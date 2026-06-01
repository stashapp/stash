package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// PreviewWidth is the fixed width stash uses for preview video (scene_preview_width in
// pkg/scene/generate/preview.go:19). Output preserves aspect ratio.
const PreviewWidth = 640

// Encoder runs the per-segment NVENC + concat-demuxer preview pipeline.
type Encoder struct {
	FFmpegPath string // path to ffmpeg.exe (must support nvenc + cuda hwaccel)
	FFprobe    string // path to ffprobe; defaults to "ffprobe" alongside FFmpegPath
	DryRun     bool   // log commands but don't execute or touch files
}

// NewEncoder builds an encoder. ffprobe is resolved to the sibling of FFmpegPath
// if not supplied.
func NewEncoder(ffmpegPath string, dryRun bool) *Encoder {
	ff := strings.TrimSpace(ffmpegPath)
	if ff == "" {
		ff = "ffmpeg"
	}
	probe := ""
	if dir, file := filepath.Split(ff); file == "ffmpeg" || file == "ffmpeg.exe" {
		ext := filepath.Ext(file)
		probe = filepath.Join(dir, "ffprobe"+ext)
	}
	if probe == "" {
		probe = "ffprobe"
	}
	return &Encoder{FFmpegPath: ff, FFprobe: probe, DryRun: dryRun}
}

// EncodePreview produces a preview MP4 at finalPath for the given source video.
// Workflow: probe → compute segment times → encode each segment (with retry) →
// concat via demuxer → atomic rename. All intermediate files live in tmpDir.
//
// Honors the same parameters stash's CPU pipeline uses
// (pkg/scene/generate/preview.go:182-215), substituting NVENC for libx264.
func (e *Encoder) EncodePreview(
	ctx context.Context,
	cfg *StashConfig,
	sourcePath, finalPath, tmpDir string,
) error {
	probe, err := e.probeVideo(ctx, sourcePath)
	if err != nil {
		return fmt.Errorf("probe %s: %w", sourcePath, err)
	}

	segments := segmentPlan(probe.Duration, cfg)
	if len(segments) == 0 {
		return errors.New("video too short for a preview (no segments computed)")
	}

	// Working directory for this scene's segment files + concat list.
	// Keyed by the final filename to keep concurrent encodes from colliding.
	stem := strings.TrimSuffix(filepath.Base(finalPath), filepath.Ext(finalPath))
	workDir := filepath.Join(tmpDir, "preview-"+stem)
	if !e.DryRun {
		if err := os.MkdirAll(workDir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", workDir, err)
		}
		defer os.RemoveAll(workDir)
	}

	// Encode each segment. On per-segment failure, retry once with slow-seek
	// per task_generate_preview.go:67-72 (the "fallback" pass).
	segFiles := make([]string, 0, len(segments))
	for i, seg := range segments {
		segPath := filepath.Join(workDir, fmt.Sprintf("seg_%03d.mp4", i))
		if err := e.encodeSegment(ctx, sourcePath, segPath, seg, cfg, probe, false); err != nil {
			if err := e.encodeSegment(ctx, sourcePath, segPath, seg, cfg, probe, true); err != nil {
				return fmt.Errorf("encode segment %d (%.2fs): %w", i, seg.Start, err)
			}
		}
		segFiles = append(segFiles, segPath)
	}

	// Concat via demuxer (no re-encode), then atomic rename into place.
	partial := finalPath + ".partial"
	if err := e.concatSegments(ctx, segFiles, partial); err != nil {
		return fmt.Errorf("concat: %w", err)
	}
	if e.DryRun {
		return nil
	}
	if err := os.Rename(partial, finalPath); err != nil {
		_ = os.Remove(partial)
		return fmt.Errorf("rename %s -> %s: %w", partial, finalPath, err)
	}
	return nil
}

// ── segment timing ─────────────────────────────────────────────────────────

type segment struct {
	Start    float64 // seconds
	Duration float64 // seconds (== cfg.PreviewSegmentDuration)
}

// segmentPlan computes the segment start times, mirroring preview.go:57-69
// (getStepSizeAndOffset) honoring previewExcludeStart/End (seconds or "%").
func segmentPlan(duration float64, cfg *StashConfig) []segment {
	if cfg.PreviewSegments <= 0 || cfg.PreviewSegmentDuration <= 0 || duration <= 0 {
		return nil
	}
	excludeStart := parseExclude(cfg.PreviewExcludeStart, duration)
	excludeEnd := parseExclude(cfg.PreviewExcludeEnd, duration)
	usable := duration - excludeStart - excludeEnd - cfg.PreviewSegmentDuration
	if usable <= 0 {
		return nil
	}
	n := cfg.PreviewSegments
	step := usable / float64(n)
	if step <= 0 {
		return nil
	}
	out := make([]segment, 0, n)
	for i := 0; i < n; i++ {
		t := excludeStart + float64(i)*step
		out = append(out, segment{Start: t, Duration: cfg.PreviewSegmentDuration})
	}
	return out
}

// parseExclude accepts "30" (seconds), "5%" (proportion of duration), or empty
// (zero). Matches preview.go:38-48 (getExcludeValue).
func parseExclude(raw string, duration float64) float64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	if strings.HasSuffix(raw, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(raw, "%"), 64)
		if err != nil {
			return 0
		}
		return duration * v / 100.0
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return v
}

// ── ffmpeg invocations ─────────────────────────────────────────────────────

// encodeSegment runs ffmpeg for a single segment. When slowSeek is false this
// uses fast-seek (`-ss` before `-i`) plus `-xerror`; on failure, the caller
// re-invokes with slowSeek=true (matching task_generate_preview.go:67-72).
func (e *Encoder) encodeSegment(
	ctx context.Context,
	source, output string,
	seg segment,
	cfg *StashConfig,
	probe *probeResult,
	slowSeek bool,
) error {
	args := []string{
		"-loglevel", "error",
		"-y",
	}
	if !slowSeek {
		args = append(args, "-xerror")
	}

	// NVDEC + NVENC: keep frames on the GPU through the seek if possible.
	args = append(args, "-hwaccel", "cuda", "-hwaccel_output_format", "cuda")

	startStr := strconv.FormatFloat(seg.Start, 'f', 3, 64)
	durStr := strconv.FormatFloat(seg.Duration, 'f', 3, 64)
	if slowSeek {
		args = append(args, "-i", source, "-ss", startStr, "-t", durStr)
	} else {
		args = append(args, "-ss", startStr, "-i", source, "-t", durStr)
	}

	// Pull frames back to system memory for the scale filter (libswscale CPU
	// scale; doing this on the GPU needs scale_cuda which isn't in every BtbN
	// build — sticking with hwdownload keeps compatibility broad).
	args = append(args,
		"-vf", "hwdownload,format=nv12,scale="+strconv.Itoa(PreviewWidth)+":-2",
	)

	// VFR detection: when probe reports near-zero framerate (broken or VFR
	// source), tell ffmpeg to drop duplicates. Mirrors the pipeline's
	// useVsync2 flag (task_generate_preview.go:62-65).
	if probe.FrameRate > 0 && probe.FrameRate <= 0.01 {
		args = append(args, "-vsync", "2")
	}

	// Video encode params: NVENC analogs of stash's libx264 settings
	// (preview.go:182-190). -cq 21 approximates -crf 21; -preset p6 is the
	// slow-quality analog of x264's "slow" preset.
	args = append(args,
		"-c:v", "h264_nvenc",
		"-preset", "p6",
		"-tune", "hq",
		"-rc", "vbr",
		"-cq", "21",
		"-profile:v", "high",
		"-level", "4.2",
		"-pix_fmt", "yuv420p",
	)

	// Audio: mirror preview.go:215 (128k AAC, downmix to stereo) when enabled.
	if cfg.PreviewAudio {
		args = append(args, "-c:a", "aac", "-b:a", "128k", "-ac", "2")
	} else {
		args = append(args, "-an")
	}

	args = append(args, "-movflags", "+faststart", output)

	return e.runFFmpeg(ctx, args)
}

// concatSegments builds a concat-demuxer list file and stream-copies the
// segments into the final output (no re-encode), which is how stash's own
// transcoder.Splice produces previews.
func (e *Encoder) concatSegments(ctx context.Context, segFiles []string, output string) error {
	listFile := output + ".concat.txt"
	if !e.DryRun {
		var sb strings.Builder
		for _, p := range segFiles {
			// ffmpeg concat demuxer expects forward-slash, single-quoted paths
			// with internal single quotes escaped.
			safe := strings.ReplaceAll(filepath.ToSlash(p), "'", "'\\''")
			fmt.Fprintf(&sb, "file '%s'\n", safe)
		}
		if err := os.WriteFile(listFile, []byte(sb.String()), 0o644); err != nil {
			return fmt.Errorf("write concat list: %w", err)
		}
		defer os.Remove(listFile)
	}

	args := []string{
		"-loglevel", "error", "-y",
		"-f", "concat", "-safe", "0",
		"-i", listFile,
		"-c", "copy",
		"-movflags", "+faststart",
		output,
	}
	return e.runFFmpeg(ctx, args)
}

func (e *Encoder) runFFmpeg(ctx context.Context, args []string) error {
	if e.DryRun {
		fmt.Println(append([]string{e.FFmpegPath}, args...))
		return nil
	}
	cmd := exec.CommandContext(ctx, e.FFmpegPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %s", err, truncate(strings.TrimSpace(stderr.String()), 500))
	}
	return nil
}

// ── ffprobe ───────────────────────────────────────────────────────────────

type probeResult struct {
	Duration  float64 // seconds
	FrameRate float64 // average fps; 0 if unknown
}

type probeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
	Streams []struct {
		CodecType  string `json:"codec_type"`
		AvgFrameRate string `json:"avg_frame_rate"`
	} `json:"streams"`
}

// probeVideo returns the source's duration + average frame rate.
func (e *Encoder) probeVideo(ctx context.Context, source string) (*probeResult, error) {
	if e.DryRun {
		return &probeResult{Duration: 60, FrameRate: 30}, nil
	}
	cmd := exec.CommandContext(ctx, e.FFprobe,
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		source,
	)
	cmd.Env = append(cmd.Env, "LC_ALL=C")
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe: %w", err)
	}
	var po probeOutput
	if err := json.Unmarshal(out, &po); err != nil {
		return nil, fmt.Errorf("decode ffprobe output: %w", err)
	}
	d, _ := strconv.ParseFloat(po.Format.Duration, 64)
	if d <= 0 {
		return nil, fmt.Errorf("invalid duration: %q", po.Format.Duration)
	}
	res := &probeResult{Duration: d}
	for _, s := range po.Streams {
		if s.CodecType == "video" {
			res.FrameRate = parseFrameRate(s.AvgFrameRate)
			break
		}
	}
	return res, nil
}

func parseFrameRate(s string) float64 {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	num, _ := strconv.ParseFloat(parts[0], 64)
	den, _ := strconv.ParseFloat(parts[1], 64)
	if den == 0 {
		return 0
	}
	return num / den
}

// ── exposed for the main loop ──────────────────────────────────────────────

// PreviewExists reports whether stash would consider this scene's preview as
// already generated. The check is purely file-system based, matching what stash
// itself does (task_generate_preview.go:100-105).
func PreviewExists(generatedPaths GeneratedPaths, checksum string) (bool, error) {
	st, err := os.Stat(generatedPaths.VideoPreview(checksum))
	if err == nil {
		return st.Size() > 0, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// EnsureTmpDir creates the worker's scratch dir under stash's generated/tmp/.
// Should be called once at worker startup. We use a short retry because SMB
// permissions glitches can occasionally race here.
func EnsureTmpDir(generatedPaths GeneratedPaths) error {
	var lastErr error
	for i := 0; i < 3; i++ {
		err := os.MkdirAll(generatedPaths.TmpDir(), 0o755)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(200 * time.Millisecond)
	}
	return lastErr
}
