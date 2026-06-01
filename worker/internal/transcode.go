package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Transcode pre-generation. stash live-transcodes any scene whose codec/container
// the browser can't play directly — brutal on the NAS's weak CPU. This writes a
// browser-friendly h264/aac MP4 to generated/transcodes/<hash>.mp4, which stash
// then serves directly (GetStreamPath, pkg/models/paths/paths_scenes.go:28-35) —
// no live transcode ever runs for that scene.

// streamable containers (ffprobe reports comma-lists, e.g. "mov,mp4,m4a,3gp,3g2,mj2").
func isMP4Container(format string) bool {
	f := strings.ToLower(format)
	return strings.Contains(f, "mp4") || strings.Contains(f, "mov") || strings.Contains(f, "m4v")
}

func isBrowserAudio(audioCodec string) bool {
	switch strings.ToLower(strings.TrimSpace(audioCodec)) {
	case "aac", "mp3", "":
		return true
	}
	return false
}

// NeedsTranscode reports whether a scene should be pre-transcoded. Mirrors
// pkg/ffmpeg/browser.go IsStreamable, but deliberately treats HEVC/h265 as
// NON-streamable: stash considers it streamable and serves it direct, yet
// Chrome/Firefox usually can't decode it — a pre-generated h264 transcode (served
// because the file exists) fixes that. Streamable == h264 video + mp4 container +
// aac/mp3 audio; everything else (HEVC, mpeg4, wmv, vp6, ac3/dts/wmav2, mkv/avi/wmv)
// gets a transcode.
func NeedsTranscode(videoCodec, audioCodec, format string) bool {
	v := strings.ToLower(strings.TrimSpace(videoCodec))
	return !(v == "h264" && isMP4Container(format) && isBrowserAudio(audioCodec))
}

// Transcode produces output (an .mp4) from source. If the video is already h264
// (only the audio or container is the problem) it stream-COPIES the video — fast
// and lossless — and only fixes audio/container. Otherwise it re-encodes video
// with NVENC h264 at full resolution. Decode is left to the CPU (no -hwaccel) so
// exotic source codecs (mpeg4/wmv/vp6/10-bit HEVC) all decode reliably; only the
// encode is on the GPU. Atomic via tmpDir + rename.
func (e *Encoder) Transcode(ctx context.Context, source, output, tmpDir, videoCodec, audioCodec string) error {
	stem := strings.TrimSuffix(filepath.Base(output), filepath.Ext(output))
	partial := filepath.Join(tmpDir, stem+".mp4.partial")

	args := []string{"-loglevel", "error", "-y", "-i", source}

	if strings.ToLower(strings.TrimSpace(videoCodec)) == "h264" {
		// Remux: keep the (already-fine) h264 video bit-for-bit.
		args = append(args, "-c:v", "copy")
		if strings.ToLower(strings.TrimSpace(audioCodec)) == "aac" {
			args = append(args, "-c:a", "copy")
		} else {
			args = append(args, "-c:a", "aac", "-b:a", "160k", "-ac", "2")
		}
	} else {
		// Full re-encode: NVENC h264, full resolution (no scale filter).
		args = append(args,
			"-c:v", "h264_nvenc",
			"-preset", "p6",
			"-tune", "hq",
			"-rc", "vbr",
			"-cq", "21",
			"-profile:v", "high",
			"-pix_fmt", "yuv420p",
			"-c:a", "aac", "-b:a", "160k", "-ac", "2",
		)
	}

	args = append(args, "-movflags", "+faststart", "-f", "mp4", partial)

	if err := e.runFFmpeg(ctx, args); err != nil {
		return fmt.Errorf("transcode: %w", err)
	}
	if e.DryRun {
		return nil
	}
	if err := os.Rename(partial, output); err != nil {
		_ = os.Remove(partial)
		return fmt.Errorf("rename %s -> %s: %w", partial, output, err)
	}
	return nil
}

// TranscodeExists reports whether a non-empty transcode already exists for hash.
func TranscodeExists(gen GeneratedPaths, checksum string) (bool, error) {
	st, err := os.Stat(gen.Transcode(checksum))
	if err == nil {
		return st.Size() > 0, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
