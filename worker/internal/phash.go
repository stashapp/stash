package internal

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"math"
	"os/exec"
	"strconv"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
)

// Phash generation — a BIT-FOR-BIT replica of stash's
// pkg/hash/videophash/phash.go. The output uint64 must be identical to what
// stash computes, or the hash won't match StashDB/TPDB fingerprints and the
// whole offload is worthless. Every constant, timestamp, frame size, montage
// layout, and library (goimagehash v1.1.0 + imaging v1.6.2 + x/image/bmp) is
// matched deliberately. DO NOT "optimize" this without re-running the
// verification gate (--verify-phash) against known native hashes.
//
// Critical fidelity rules (mirroring stash exactly):
//   - 25 frames (5x5), each scaled to width 160 (height -2, even, AR-preserved).
//   - timestamps: offset = 0.05*dur; step = 0.9*dur/25; t_i = offset + i*step.
//   - duration is stash's STORED duration (from the GraphQL API), NOT a fresh
//     ffprobe — a few-ms difference shifts every timestamp and changes the hash.
//   - frames are CPU-decoded BMP (lossless). NO -hwaccel: NVDEC can yield
//     subtly different pixels than libavcodec's software decoder, which would
//     change the perceptual hash.
//   - -ss value is formatted via fmt.Sprint to match stash's Args.Seek exactly.
const (
	phashScreenshotSize = 160
	phashColumns        = 5
	phashRows           = 5
	phashChunkCount     = phashColumns * phashRows // 25
)

// GeneratePhash computes stash's perceptual hash for a single video file.
// duration MUST be the stash-stored duration for this file (seconds).
// Returns the raw uint64 hash; the caller formats it as hex for the API.
func (e *Encoder) GeneratePhash(ctx context.Context, source string, duration float64) (uint64, error) {
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration %.4f", duration)
	}

	// Sample 25 frames offset 5% in from each end (matches generateSprite).
	offset := 0.05 * duration
	stepSize := (0.9 * duration) / float64(phashChunkCount)

	images := make([]image.Image, 0, phashChunkCount)
	for i := 0; i < phashChunkCount; i++ {
		t := offset + (float64(i) * stepSize)
		img, err := e.phashFrame(ctx, source, t)
		if err != nil {
			return 0, fmt.Errorf("phash frame %d at %gs: %w", i, t, err)
		}
		images = append(images, img)
	}

	sprite := combinePhashImages(images)
	hash, err := goimagehash.PerceptionHash(sprite)
	if err != nil {
		return 0, fmt.Errorf("computing phash from sprite: %w", err)
	}
	return hash.GetHash(), nil
}

// phashFrame extracts ONE frame at time t as a lossless BMP via stdout, decoding
// it in Go. The ffmpeg argument vector is a byte-for-byte match of stash's
// transcoder.ScreenshotTime(input, t, {Width:160, OutputType: BMP}):
//
//	-v error -y -ss <t> -i <src> -frames:v 1 -vf scale=160:-2 -c:v bmp -f rawvideo -
//
// where <t> is fmt.Sprint(t) (stash's Args.Seek formatting).
func (e *Encoder) phashFrame(ctx context.Context, source string, t float64) (image.Image, error) {
	args := []string{
		"-v", "error",
		"-y",
		"-ss", fmt.Sprint(t),
		"-i", source,
		"-frames:v", "1",
		"-vf", "scale=" + strconv.Itoa(phashScreenshotSize) + ":-2",
		"-c:v", "bmp",
		"-f", "rawvideo",
		"-",
	}
	cmd := exec.CommandContext(ctx, e.FFmpegPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg: %w: %s", err, truncate(stderr.String(), 200))
	}

	img, err := bmp.Decode(bytes.NewReader(stdout.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("decoding bmp (%d bytes): %w", stdout.Len(), err)
	}
	return img, nil
}

// combinePhashImages tiles the 25 frames into one montage, replicating stash's
// combineImages exactly: an imaging.New NRGBA canvas of (w*5)x(h*5) with each
// frame pasted at x=w*(i%columns), y=h*floor(i/rows).
func combinePhashImages(images []image.Image) image.Image {
	width := images[0].Bounds().Size().X
	height := images[0].Bounds().Size().Y
	canvasWidth := width * phashColumns
	canvasHeight := height * phashRows
	montage := imaging.New(canvasWidth, canvasHeight, color.NRGBA{})
	for index := 0; index < len(images); index++ {
		x := width * (index % phashColumns)
		y := height * int(math.Floor(float64(index)/float64(phashRows)))
		montage = imaging.Paste(montage, images[index], image.Pt(x, y))
	}
	return montage
}
