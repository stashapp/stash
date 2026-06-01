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

	// Graceful Ctrl-C: first signal cancels the context (lets in-flight ffmpeg
	// finish or be killed by its own cancellation), second forces exit.
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

	limit   int
	perPage int
	watch   bool
	dryRun  bool
}

func parseFlags() (*runConfig, error) {
	fs := flag.NewFlagSet("stash-worker", flag.ExitOnError)
	stashURL := fs.String("stash-url", "http://localhost:9999", "Base URL of the stash instance.")
	apiKey := fs.String("stash-api-key", os.Getenv("STASH_API_KEY"), "Optional stash API key (also via STASH_API_KEY env).")
	mediaPrefix := fs.String("media-prefix", "", "STASH=WORKER prefix rewrite for media file paths (e.g. \"/data=\\\\overwatch-stash\\torrents\").")
	generatedPrefix := fs.String("generated-prefix", "", "STASH=WORKER prefix rewrite for the generated/ dir.")
	ffmpegPath := fs.String("ffmpeg", "ffmpeg", "Path to ffmpeg.exe with NVENC + NVDEC support.")
	limit := fs.Int("limit", 0, "Cap on scenes processed per run (0 = unbounded).")
	perPage := fs.Int("per-page", 200, "GraphQL pagination size for scene enumeration.")
	watch := fs.Bool("watch", false, "Keep polling for new work instead of exiting when the queue is empty.")
	dryRun := fs.Bool("dry-run", false, "Print ffmpeg commands without executing or writing any files.")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	cfg := &runConfig{
		stashURL:    strings.TrimSpace(*stashURL),
		stashAPIKey: strings.TrimSpace(*apiKey),
		ffmpegPath:  strings.TrimSpace(*ffmpegPath),
		limit:       *limit,
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
	return cfg, nil
}

// ── main loop ──────────────────────────────────────────────────────────────

// run drives the preview-generation loop. It reads stash's config, enumerates
// scenes in id-ASC order, and encodes a preview for each scene whose target
// file is missing. With --watch, it idles between empty passes; otherwise it
// exits after the first empty pass.
func run(ctx context.Context, c *runConfig) error {
	stash := internal.NewStashClient(c.stashURL, c.stashAPIKey)
	scfg, err := stash.FetchConfig(ctx)
	if err != nil {
		return fmt.Errorf("fetch stash config: %w", err)
	}
	log.Printf("stash config: algorithm=%s previewSegments=%d previewDur=%.2fs audio=%v",
		scfg.VideoFileNamingAlgorithm, scfg.PreviewSegments, scfg.PreviewSegmentDuration, scfg.PreviewAudio)

	// Worker's view of stash's generated/ dir is the rewriter's right-hand side.
	gen := internal.GeneratedPaths{Root: c.generatedRewriter.To()}
	if err := internal.EnsureTmpDir(gen); err != nil {
		return fmt.Errorf("ensure tmp dir %s: %w", gen.TmpDir(), err)
	}
	enc := internal.NewEncoder(c.ffmpegPath, c.dryRun)

	processed := 0
	encoded := 0
	skipped := 0
	failed := 0
	emptyPasses := 0

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		page := 1
		passEncoded := 0
		for {
			scenes, total, err := stash.FetchScenesPage(ctx, page, c.perPage)
			if err != nil {
				return fmt.Errorf("page %d: %w", page, err)
			}
			if len(scenes) == 0 {
				break
			}
			if page == 1 {
				log.Printf("scanning %d scenes (page %d, per-page %d)", total, page, c.perPage)
			}
			for _, s := range scenes {
				if c.limit > 0 && processed >= c.limit {
					break
				}
				processed++
				outcome := processScene(ctx, &s, scfg, c, gen, enc)
				switch outcome {
				case outcomeEncoded:
					encoded++
					passEncoded++
				case outcomeSkipped:
					skipped++
				case outcomeFailed:
					failed++
				}
			}
			if c.limit > 0 && processed >= c.limit {
				break
			}
			if len(scenes) < c.perPage {
				break // last page
			}
			page++
		}
		log.Printf("pass complete: encoded=%d skipped=%d failed=%d (processed %d total)",
			passEncoded, skipped, failed, processed)

		if !c.watch {
			break
		}
		if c.limit > 0 && processed >= c.limit {
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
		// Sleep before next pass; honor cancellation.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(60 * time.Second):
		}
	}

	log.Printf("done: encoded=%d skipped=%d failed=%d", encoded, skipped, failed)
	return nil
}

type outcome int

const (
	outcomeEncoded outcome = iota
	outcomeSkipped
	outcomeFailed
)

// processScene handles one scene: checks if the preview is already on disk,
// and if not, runs the encode pipeline. All errors are logged and converted
// to outcomes — one bad scene doesn't kill the whole run.
func processScene(
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

	stashSourcePath := s.Files[0].Path
	source := stashSourcePath
	if c.mediaRewriter != nil {
		source = c.mediaRewriter.Rewrite(stashSourcePath)
	}
	output := gen.VideoPreview(hash)

	log.Printf("scene %s: encoding preview %s", s.ID, filepath.Base(output))
	if err := enc.EncodePreview(ctx, scfg, source, output, gen.TmpDir()); err != nil {
		log.Printf("scene %s: encode failed: %v", s.ID, err)
		return outcomeFailed
	}
	return outcomeEncoded
}

