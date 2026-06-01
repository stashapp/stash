// Command stash-worker runs the external GPU media-generation pipeline against
// a stash instance. See ../../README.md and ../../docs/llm/EXTERNAL-WORKERS.md.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Ryokushen/stash/worker/internal"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("[stash-worker] ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("interrupted; finishing current scene and exiting (^C again to force)")
		cancel()
		<-sigCh
		os.Exit(130)
	}()

	if err := run(ctx, cfg); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("done (cancelled)")
			os.Exit(130)
		}
		log.Fatal(err)
	}
}

// ── flags ──────────────────────────────────────────────────────────────────

type runConfig struct {
	stashURL    string
	stashAPIKey string

	mediaRewriter     *internal.PrefixRewriter
	generatedRewriter *internal.PrefixRewriter

	ffmpegPath string

	tasks       []string // ordered task names: "previews", "covers", "sprites", "phash"
	limit       int
	maxFailures int
	perPage     int
	watch       bool
	dryRun      bool

	// verifyPhash > 0 runs the phash VERIFICATION GATE instead of the normal
	// loop: recompute phashes for N files that already have a native stash
	// phash and compare. Used to prove bit-exactness before any bulk phash run.
	verifyPhash int
}

func parseFlags() (*runConfig, error) {
	fs := flag.NewFlagSet("stash-worker", flag.ExitOnError)
	stashURL := fs.String("stash-url", "http://localhost:9999", "Base URL of the stash instance.")
	apiKey := fs.String("stash-api-key", os.Getenv("STASH_API_KEY"), "Optional stash API key (also via STASH_API_KEY env).")
	mediaPrefix := fs.String("media-prefix", "", "STASH=WORKER prefix rewrite for media file paths (e.g. \"/data=\\\\overwatch-stash\\torrents\").")
	generatedPrefix := fs.String("generated-prefix", "", "STASH=WORKER prefix rewrite for the generated/ dir.")
	ffmpegPath := fs.String("ffmpeg", "ffmpeg", "Path to ffmpeg.exe with NVENC + NVDEC support.")
	tasks := fs.String("tasks", "previews", "Comma-separated tasks to run, in order. Available: previews, covers, sprites, phash.")
	verifyPhash := fs.Int("verify-phash", 0, "VERIFICATION GATE: recompute N files that already have a native stash phash and compare (proves bit-exactness). Skips the normal loop.")
	limit := fs.Int("limit", 0, "Cap on items ENCODED/applied per run across ALL tasks (0 = unbounded). Useful for tests.")
	maxFailures := fs.Int("max-failures", 5, "Abort the run after N consecutive failures (catches systemic problems).")
	perPage := fs.Int("per-page", 200, "GraphQL pagination size for scene enumeration.")
	watch := fs.Bool("watch", false, "Keep polling for new work instead of exiting when the queue is empty.")
	dryRun := fs.Bool("dry-run", false, "Print actions without writing.")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	cfg := &runConfig{
		stashURL:    strings.TrimSpace(*stashURL),
		stashAPIKey: strings.TrimSpace(*apiKey),
		ffmpegPath:  strings.TrimSpace(*ffmpegPath),
		limit:       *limit,
		maxFailures: *maxFailures,
		perPage:     *perPage,
		watch:       *watch,
		dryRun:      *dryRun,
		verifyPhash: *verifyPhash,
	}
	if cfg.stashURL == "" {
		return nil, errors.New("--stash-url is required")
	}
	if *mediaPrefix != "" {
		mr, err := internal.ParsePrefixRewriter(*mediaPrefix)
		if err != nil {
			return nil, fmt.Errorf("--media-prefix: %w", err)
		}
		cfg.mediaRewriter = mr
	}

	// verify-phash mode bypasses the task loop entirely — it only needs to read
	// media (for recompute) and the API (for stored hashes), never generated/.
	if cfg.verifyPhash > 0 {
		return cfg, nil
	}

	for _, t := range strings.Split(*tasks, ",") {
		t = strings.TrimSpace(strings.ToLower(t))
		if t == "" {
			continue
		}
		if t != "previews" && t != "covers" && t != "sprites" && t != "phash" && t != "transcode" {
			return nil, fmt.Errorf("--tasks: unknown task %q (available: previews, covers, sprites, phash, transcode)", t)
		}
		cfg.tasks = append(cfg.tasks, t)
	}
	if len(cfg.tasks) == 0 {
		return nil, errors.New("--tasks: at least one task required")
	}

	// Only previews and sprites write into generated/. covers (API upload) and
	// phash (fingerprint mutation) don't, so --generated-prefix is optional for
	// those.
	needsGenerated := false
	for _, t := range cfg.tasks {
		if t == "previews" || t == "sprites" || t == "transcode" {
			needsGenerated = true
		}
	}
	if needsGenerated {
		if *generatedPrefix == "" {
			return nil, errors.New("--generated-prefix is required for previews/sprites (where to write generated files)")
		}
		gr, err := internal.ParsePrefixRewriter(*generatedPrefix)
		if err != nil {
			return nil, fmt.Errorf("--generated-prefix: %w", err)
		}
		cfg.generatedRewriter = gr
	} else if *generatedPrefix != "" {
		// Accept it if supplied (harmless) so mixed task sets still work.
		gr, err := internal.ParsePrefixRewriter(*generatedPrefix)
		if err != nil {
			return nil, fmt.Errorf("--generated-prefix: %w", err)
		}
		cfg.generatedRewriter = gr
	}
	return cfg, nil
}

// ── main loop ──────────────────────────────────────────────────────────────

type outcome int

const (
	outcomeEncoded outcome = iota // counts as work done (cover applied OR preview encoded)
	outcomeSkipped
	outcomeFailed
)

// totals tracks counters across all tasks in a single invocation.
type totals struct {
	processed, encoded, skipped, failed int
}

func run(ctx context.Context, c *runConfig) error {
	stash := internal.NewStashClient(c.stashURL, c.stashAPIKey)
	enc := internal.NewEncoder(c.ffmpegPath, c.dryRun)

	// Verification gate: prove worker phashes match native stash phashes bit-for-bit
	// before any bulk phash run. Bypasses the normal task loop.
	if c.verifyPhash > 0 {
		return runVerifyPhash(ctx, c, stash, enc)
	}

	scfg, err := stash.FetchConfig(ctx)
	if err != nil {
		return fmt.Errorf("fetch stash config: %w", err)
	}
	log.Printf("stash config: algorithm=%s previewSegments=%d previewDur=%.2fs audio=%v",
		scfg.VideoFileNamingAlgorithm, scfg.PreviewSegments, scfg.PreviewSegmentDuration, scfg.PreviewAudio)
	log.Printf("tasks: %s", strings.Join(c.tasks, ", "))

	var gen internal.GeneratedPaths
	if c.generatedRewriter != nil {
		gen = internal.GeneratedPaths{Root: c.generatedRewriter.To()}
		if err := internal.EnsureTmpDir(gen); err != nil {
			return fmt.Errorf("ensure tmp dir %s: %w", gen.TmpDir(), err)
		}
	}
	t := &totals{}
	emptyPasses := 0

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		passEncoded := 0
		for _, taskName := range c.tasks {
			if c.limit > 0 && t.encoded >= c.limit {
				break
			}
			n, err := runTask(ctx, taskName, c, stash, scfg, gen, enc, t)
			if err != nil {
				return err
			}
			passEncoded += n
		}
		log.Printf("pass complete: encoded=%d failed=%d (processed %d total)", t.encoded, t.failed, t.processed)

		if !c.watch {
			break
		}
		if c.limit > 0 && t.encoded >= c.limit {
			log.Printf("limit %d reached; exiting watch loop", c.limit)
			break
		}
		if passEncoded == 0 {
			emptyPasses++
			if emptyPasses >= 2 {
				log.Printf("two consecutive empty passes; exiting")
				break
			}
		} else {
			emptyPasses = 0
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(60 * time.Second):
		}
	}

	log.Printf("done: encoded=%d skipped=%d failed=%d", t.encoded, t.skipped, t.failed)
	return nil
}

// runTask drives one task's enumeration + per-scene dispatch. Returns the
// count of items encoded in this pass (used to break the watch loop's empty
// detection).
func runTask(
	ctx context.Context,
	taskName string,
	c *runConfig,
	stash *internal.StashClient,
	scfg *internal.StashConfig,
	gen internal.GeneratedPaths,
	enc *internal.Encoder,
	t *totals,
) (int, error) {
	// is_missing filter on the stash query for tasks where stash tracks
	// missingness. Previews aren't tracked — we filesystem-check instead.
	// shrinksOnSuccess = true means: as the worker applies items, those items
	// drop out of the server-side filter result set. The pagination strategy
	// MUST always re-fetch page 1 in that case — otherwise pages 2+ silently
	// skip items (the resultset shifted under us). For tasks that use a
	// client-side filter (file existence check), the resultset is stable and
	// normal page++ pagination is correct.
	missingFilter := ""
	shrinksOnSuccess := false
	switch taskName {
	case "covers":
		missingFilter = "cover"
		shrinksOnSuccess = true
	case "phash":
		// is_missing:"phash" left-joins fingerprints and selects files with no
		// phash row; as the worker sets phashes the result set shrinks, exactly
		// like covers — so the same re-fetch-page-1 strategy is required.
		missingFilter = "phash"
		shrinksOnSuccess = true
	}

	page := 1
	passEncoded := 0
	consecutiveFailures := 0
	firstPageLogged := false
	for {
		scenes, total, err := stash.FetchScenesPage(ctx, page, c.perPage, missingFilter)
		if err != nil {
			return passEncoded, fmt.Errorf("[%s] page %d: %w", taskName, page, err)
		}
		if len(scenes) == 0 {
			break
		}
		if !firstPageLogged {
			log.Printf("[%s] scanning %d scenes (per-page %d)", taskName, total, c.perPage)
			firstPageLogged = true
		}
		pageProgressed := 0 // count successful encodes within THIS page
		for _, s := range scenes {
			if c.limit > 0 && t.encoded >= c.limit {
				return passEncoded, nil
			}
			if c.maxFailures > 0 && consecutiveFailures >= c.maxFailures {
				return passEncoded, fmt.Errorf("[%s] aborting: %d consecutive failures", taskName, consecutiveFailures)
			}
			t.processed++
			out := processScene(ctx, taskName, &s, scfg, c, gen, enc, stash)
			switch out {
			case outcomeEncoded:
				t.encoded++
				passEncoded++
				pageProgressed++
				consecutiveFailures = 0
			case outcomeSkipped:
				t.skipped++
			case outcomeFailed:
				t.failed++
				consecutiveFailures++
			}
		}
		// Pagination strategy:
		//   shrinksOnSuccess + page made progress → re-fetch page 1 (filter has shrunk;
		//     advancing the page number would skip items that just shifted into the
		//     lower positions).
		//   shrinksOnSuccess + page made NO progress → advance, otherwise we'd loop
		//     forever on a page full of un-processable scenes (all skipped/failed).
		//   stable resultset → straight page++.
		if shrinksOnSuccess && pageProgressed > 0 {
			page = 1
		} else if len(scenes) < c.perPage {
			break
		} else {
			page++
		}
	}
	return passEncoded, nil
}

// processScene dispatches by task. Per-scene errors are logged and turned into
// outcomes — one bad scene doesn't kill the run.
func processScene(
	ctx context.Context,
	taskName string,
	s *internal.Scene,
	scfg *internal.StashConfig,
	c *runConfig,
	gen internal.GeneratedPaths,
	enc *internal.Encoder,
	stash *internal.StashClient,
) outcome {
	switch taskName {
	case "previews":
		return processPreview(ctx, s, scfg, c, gen, enc)
	case "covers":
		return processCover(ctx, s, c, enc, stash)
	case "sprites":
		return processSprite(ctx, s, scfg, c, gen, enc)
	case "phash":
		return processPhash(ctx, s, c, enc, stash)
	case "transcode":
		return processTranscode(ctx, s, scfg, c, gen, enc)
	default:
		log.Printf("scene %s: unknown task %q; skipping", s.ID, taskName)
		return outcomeSkipped
	}
}

// ── previews ───────────────────────────────────────────────────────────────

func processPreview(
	ctx context.Context,
	s *internal.Scene,
	scfg *internal.StashConfig,
	c *runConfig,
	gen internal.GeneratedPaths,
	enc *internal.Encoder,
) outcome {
	if len(s.Files) == 0 {
		log.Printf("scene %s: no files; skipping", s.ID)
		return outcomeSkipped
	}
	hash, err := s.PrimaryHash(scfg.VideoFileNamingAlgorithm)
	if err != nil {
		log.Printf("scene %s: %v; skipping", s.ID, err)
		return outcomeSkipped
	}
	exists, err := internal.PreviewExists(gen, hash)
	if err != nil {
		log.Printf("scene %s: stat preview: %v", s.ID, err)
		return outcomeFailed
	}
	if exists {
		return outcomeSkipped
	}

	source := s.Files[0].Path
	if c.mediaRewriter != nil {
		source = c.mediaRewriter.Rewrite(source)
	}
	output := gen.VideoPreview(hash)

	log.Printf("scene %s: encoding preview %s", s.ID, filepath.Base(output))
	if err := enc.EncodePreview(ctx, scfg, source, output, gen.TmpDir()); err != nil {
		log.Printf("scene %s: encode failed: %v", s.ID, err)
		return outcomeFailed
	}
	return outcomeEncoded
}

// ── covers ─────────────────────────────────────────────────────────────────

// processCover extracts a single JPEG frame at stash's 20% timestamp and
// uploads it via sceneUpdate. Stash decodes the base64 and writes to its blob
// storage (DB or blobs/ dir, per config).
func processCover(
	ctx context.Context,
	s *internal.Scene,
	c *runConfig,
	enc *internal.Encoder,
	stash *internal.StashClient,
) outcome {
	if len(s.Files) == 0 {
		log.Printf("scene %s: no files; skipping", s.ID)
		return outcomeSkipped
	}
	pf := s.Files[0]
	if pf.Duration <= 0 {
		log.Printf("scene %s: no duration metadata; skipping cover", s.ID)
		return outcomeSkipped
	}
	source := pf.Path
	if c.mediaRewriter != nil {
		source = c.mediaRewriter.Rewrite(source)
	}
	at := internal.CoverScreenshotProportion * pf.Duration

	log.Printf("scene %s: extracting cover at %.1fs (%.0f%% of %.0fs)",
		s.ID, at, internal.CoverScreenshotProportion*100, pf.Duration)
	jpeg, err := enc.ExtractCover(ctx, source, at)
	if err != nil {
		log.Printf("scene %s: cover extract failed: %v", s.ID, err)
		return outcomeFailed
	}
	if c.dryRun {
		log.Printf("scene %s: (dry-run) would upload %d-byte JPEG", s.ID, len(jpeg))
		return outcomeEncoded
	}
	if err := stash.SceneUpdateCover(ctx, s.ID, jpeg); err != nil {
		log.Printf("scene %s: upload cover: %v", s.ID, err)
		return outcomeFailed
	}
	return outcomeEncoded
}

// ── sprites ────────────────────────────────────────────────────────────────

// processSprite generates a sprite + VTT pair for scenes that don't already
// have both. Like previews, stash detects sprite existence lazily by file
// path, so the filter is filesystem-based.
func processSprite(
	ctx context.Context,
	s *internal.Scene,
	scfg *internal.StashConfig,
	c *runConfig,
	gen internal.GeneratedPaths,
	enc *internal.Encoder,
) outcome {
	if len(s.Files) == 0 {
		log.Printf("scene %s: no files; skipping", s.ID)
		return outcomeSkipped
	}
	hash, err := s.PrimaryHash(scfg.VideoFileNamingAlgorithm)
	if err != nil {
		log.Printf("scene %s: %v; skipping", s.ID, err)
		return outcomeSkipped
	}
	exists, err := internal.SpriteExists(gen, hash)
	if err != nil {
		log.Printf("scene %s: stat sprite: %v", s.ID, err)
		return outcomeFailed
	}
	if exists {
		return outcomeSkipped
	}

	source := s.Files[0].Path
	if c.mediaRewriter != nil {
		source = c.mediaRewriter.Rewrite(source)
	}
	spritePath := gen.SpriteImage(hash)
	vttPath := gen.SpriteVTT(hash)

	log.Printf("scene %s: generating sprite + VTT", s.ID)
	if err := enc.GenerateSprite(ctx, source, spritePath, vttPath, gen.TmpDir()); err != nil {
		log.Printf("scene %s: sprite failed: %v", s.ID, err)
		return outcomeFailed
	}
	return outcomeEncoded
}

// ── transcode ──────────────────────────────────────────────────────────────

// processTranscode pre-generates a browser-friendly h264/aac MP4 into
// generated/transcodes/<hash>.mp4 for scenes stash can't stream directly, so the
// NAS never has to live-transcode them. Skips scenes that already have a transcode
// or are already streamable. Like previews, existence is a filesystem check.
func processTranscode(
	ctx context.Context,
	s *internal.Scene,
	scfg *internal.StashConfig,
	c *runConfig,
	gen internal.GeneratedPaths,
	enc *internal.Encoder,
) outcome {
	if len(s.Files) == 0 {
		log.Printf("scene %s: no files; skipping", s.ID)
		return outcomeSkipped
	}
	hash, err := s.PrimaryHash(scfg.VideoFileNamingAlgorithm)
	if err != nil {
		log.Printf("scene %s: %v; skipping", s.ID, err)
		return outcomeSkipped
	}
	exists, err := internal.TranscodeExists(gen, hash)
	if err != nil {
		log.Printf("scene %s: stat transcode: %v", s.ID, err)
		return outcomeFailed
	}
	if exists {
		return outcomeSkipped
	}
	pf := s.Files[0]
	if !internal.NeedsTranscode(pf.VideoCodec, pf.AudioCodec, pf.Format) {
		return outcomeSkipped // already directly browser-streamable
	}

	source := pf.Path
	if c.mediaRewriter != nil {
		source = c.mediaRewriter.Rewrite(source)
	}
	output := gen.Transcode(hash)
	if !c.dryRun {
		if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
			log.Printf("scene %s: mkdir transcodes: %v", s.ID, err)
			return outcomeFailed
		}
	}

	mode := "remux (copy video)"
	if strings.ToLower(pf.VideoCodec) != "h264" {
		mode = "NVENC h264 re-encode"
	}
	log.Printf("scene %s: transcoding [%s] %s/%s/%s -> %s",
		s.ID, mode, pf.VideoCodec, pf.AudioCodec, pf.Format, filepath.Base(output))
	if err := enc.Transcode(ctx, source, output, gen.TmpDir(), pf.VideoCodec, pf.AudioCodec); err != nil {
		log.Printf("scene %s: transcode failed: %v", s.ID, err)
		return outcomeFailed
	}
	return outcomeEncoded
}

// ── phash ──────────────────────────────────────────────────────────────────

// filePhash returns the file's existing phash fingerprint value (hex), if any.
func filePhash(f internal.SceneFile) (string, bool) {
	for _, fp := range f.Fingerprints {
		if strings.EqualFold(fp.Type, "phash") && fp.Value != "" {
			return fp.Value, true
		}
	}
	return "", false
}

// processPhash computes stash's perceptual hash for every file of a scene that
// doesn't already have one, and writes it back via fileSetFingerprints. phash is
// a per-FILE fingerprint, so a multi-file scene may produce several. Uses the
// stash-stored duration (NOT a fresh probe) so the sampled timestamps — and thus
// the hash — match what stash itself would compute.
func processPhash(
	ctx context.Context,
	s *internal.Scene,
	c *runConfig,
	enc *internal.Encoder,
	stash *internal.StashClient,
) outcome {
	if len(s.Files) == 0 {
		log.Printf("scene %s: no files; skipping", s.ID)
		return outcomeSkipped
	}
	did := false
	for i := range s.Files {
		f := s.Files[i]
		if _, has := filePhash(f); has {
			continue // already phashed
		}
		if f.Duration <= 0 {
			log.Printf("scene %s file %s: no duration metadata; skipping phash", s.ID, f.ID)
			continue
		}
		source := f.Path
		if c.mediaRewriter != nil {
			source = c.mediaRewriter.Rewrite(source)
		}
		log.Printf("scene %s file %s: computing phash (%.0fs)", s.ID, f.ID, f.Duration)
		hash, err := enc.GeneratePhash(ctx, source, f.Duration)
		if err != nil {
			log.Printf("scene %s file %s: phash failed: %v", s.ID, f.ID, err)
			return outcomeFailed
		}
		if c.dryRun {
			log.Printf("scene %s file %s: (dry-run) phash=%s", s.ID, f.ID, strconv.FormatUint(hash, 16))
			did = true
			continue
		}
		if err := stash.SetFilePhash(ctx, f.ID, hash); err != nil {
			log.Printf("scene %s file %s: set phash: %v", s.ID, f.ID, err)
			return outcomeFailed
		}
		log.Printf("scene %s file %s: phash=%s set", s.ID, f.ID, strconv.FormatUint(hash, 16))
		did = true
	}
	if did {
		return outcomeEncoded
	}
	return outcomeSkipped
}

// runVerifyPhash is the bit-exactness GATE. It recomputes phashes for files that
// ALREADY have a native stash phash and compares. If any differ, the worker's
// hashes won't match StashDB/TPDB fingerprints and a bulk run is worthless — so
// this returns a non-zero error and you must NOT proceed to --tasks phash.
func runVerifyPhash(ctx context.Context, c *runConfig, stash *internal.StashClient, enc *internal.Encoder) error {
	target := c.verifyPhash
	log.Printf("PHASH VERIFICATION GATE: recomputing up to %d files that already have a native phash", target)

	var checked, matched, mismatched, skipped int
	page := 1
	for checked < target {
		if err := ctx.Err(); err != nil {
			return err
		}
		scenes, total, err := stash.FetchScenesPage(ctx, page, c.perPage, "")
		if err != nil {
			return fmt.Errorf("verify: page %d: %w", page, err)
		}
		if len(scenes) == 0 {
			break
		}
		if page == 1 {
			log.Printf("verify: library has %d scenes", total)
		}
		for _, s := range scenes {
			if checked >= target {
				break
			}
			for _, f := range s.Files {
				if checked >= target {
					break
				}
				native, has := filePhash(f)
				if !has || f.Duration <= 0 {
					continue
				}
				nativeU, err := strconv.ParseUint(native, 16, 64)
				if err != nil {
					log.Printf("verify: scene %s file %s: bad native phash %q; skipping", s.ID, f.ID, native)
					skipped++
					continue
				}
				source := f.Path
				if c.mediaRewriter != nil {
					source = c.mediaRewriter.Rewrite(source)
				}
				got, err := enc.GeneratePhash(ctx, source, f.Duration)
				if err != nil {
					log.Printf("verify: scene %s file %s: recompute failed: %v", s.ID, f.ID, err)
					skipped++
					continue
				}
				checked++
				if got == nativeU {
					matched++
					log.Printf("MATCH    scene %s file %s  %s", s.ID, f.ID, native)
				} else {
					mismatched++
					log.Printf("MISMATCH scene %s file %s  native=%s computed=%s",
						s.ID, f.ID, native, strconv.FormatUint(got, 16))
				}
			}
		}
		if len(scenes) < c.perPage {
			break
		}
		page++
	}

	log.Printf("verification: %d checked, %d matched, %d mismatched, %d skipped",
		checked, matched, mismatched, skipped)
	if checked == 0 {
		return errors.New("no files with a native phash found — generate some phashes in stash first, then re-run the gate")
	}
	if mismatched > 0 {
		return fmt.Errorf("VERIFICATION FAILED: %d/%d phashes did not match native — DO NOT run --tasks phash", mismatched, checked)
	}
	log.Printf("VERIFICATION PASSED: all %d phashes are bit-identical to native — safe to run --tasks phash", checked)
	return nil
}
