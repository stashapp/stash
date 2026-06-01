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

	tasks       []string // ordered task names: "previews", "covers"
	limit       int
	maxFailures int
	perPage     int
	watch       bool
	dryRun      bool
}

func parseFlags() (*runConfig, error) {
	fs := flag.NewFlagSet("stash-worker", flag.ExitOnError)
	stashURL := fs.String("stash-url", "http://localhost:9999", "Base URL of the stash instance.")
	apiKey := fs.String("stash-api-key", os.Getenv("STASH_API_KEY"), "Optional stash API key (also via STASH_API_KEY env).")
	mediaPrefix := fs.String("media-prefix", "", "STASH=WORKER prefix rewrite for media file paths (e.g. \"/data=\\\\overwatch-stash\\torrents\").")
	generatedPrefix := fs.String("generated-prefix", "", "STASH=WORKER prefix rewrite for the generated/ dir.")
	ffmpegPath := fs.String("ffmpeg", "ffmpeg", "Path to ffmpeg.exe with NVENC + NVDEC support.")
	tasks := fs.String("tasks", "previews", "Comma-separated tasks to run, in order. Available: previews, covers, sprites.")
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
	}
	if cfg.stashURL == "" {
		return nil, errors.New("--stash-url is required")
	}
	if *generatedPrefix == "" {
		return nil, errors.New("--generated-prefix is required (where to write generated files)")
	}
	gr, err := internal.ParsePrefixRewriter(*generatedPrefix)
	if err != nil {
		return nil, fmt.Errorf("--generated-prefix: %w", err)
	}
	cfg.generatedRewriter = gr
	if *mediaPrefix != "" {
		mr, err := internal.ParsePrefixRewriter(*mediaPrefix)
		if err != nil {
			return nil, fmt.Errorf("--media-prefix: %w", err)
		}
		cfg.mediaRewriter = mr
	}
	for _, t := range strings.Split(*tasks, ",") {
		t = strings.TrimSpace(strings.ToLower(t))
		if t == "" {
			continue
		}
		if t != "previews" && t != "covers" && t != "sprites" {
			return nil, fmt.Errorf("--tasks: unknown task %q (available: previews, covers, sprites)", t)
		}
		cfg.tasks = append(cfg.tasks, t)
	}
	if len(cfg.tasks) == 0 {
		return nil, errors.New("--tasks: at least one task required")
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
	scfg, err := stash.FetchConfig(ctx)
	if err != nil {
		return fmt.Errorf("fetch stash config: %w", err)
	}
	log.Printf("stash config: algorithm=%s previewSegments=%d previewDur=%.2fs audio=%v",
		scfg.VideoFileNamingAlgorithm, scfg.PreviewSegments, scfg.PreviewSegmentDuration, scfg.PreviewAudio)
	log.Printf("tasks: %s", strings.Join(c.tasks, ", "))

	gen := internal.GeneratedPaths{Root: c.generatedRewriter.To()}
	if err := internal.EnsureTmpDir(gen); err != nil {
		return fmt.Errorf("ensure tmp dir %s: %w", gen.TmpDir(), err)
	}
	enc := internal.NewEncoder(c.ffmpegPath, c.dryRun)
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
	missingFilter := ""
	if taskName == "covers" {
		missingFilter = "cover"
	}

	page := 1
	passEncoded := 0
	consecutiveFailures := 0
	for {
		scenes, total, err := stash.FetchScenesPage(ctx, page, c.perPage, missingFilter)
		if err != nil {
			return passEncoded, fmt.Errorf("[%s] page %d: %w", taskName, page, err)
		}
		if len(scenes) == 0 {
			break
		}
		if page == 1 {
			log.Printf("[%s] scanning %d scenes (per-page %d)", taskName, total, c.perPage)
		}
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
				consecutiveFailures = 0
			case outcomeSkipped:
				t.skipped++
			case outcomeFailed:
				t.failed++
				consecutiveFailures++
			}
		}
		if len(scenes) < c.perPage {
			break
		}
		page++
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
