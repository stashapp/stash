package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/internal/identify"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
)

// RegisterScraperTools adds scraper-management + bulk-identify tools. They use the
// global manager singleton (the same way the GraphQL resolvers do), so they need no
// Deps wiring.
func RegisterScraperTools(reg *Registry) {
	reg.Register(listScrapersTool())
	reg.Register(reloadScrapersTool())
	reg.Register(identifyScenesTool())
	reg.Register(identifyScenesFastTool())
}

type scraperOut struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Types []string `json:"types"`
}

func scraperSupportedTypes(s *scraper.Scraper) []string {
	var t []string
	if s.Scene != nil {
		t = append(t, "scene")
	}
	if s.Performer != nil {
		t = append(t, "performer")
	}
	if s.Gallery != nil {
		t = append(t, "gallery")
	}
	if s.Image != nil {
		t = append(t, "image")
	}
	if s.Group != nil {
		t = append(t, "movie")
	}
	return t
}

func hasString(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

func listScrapersTool() *Tool {
	return &Tool{
		Name: "list_scrapers",
		Description: "List installed metadata scrapers. Optional `query` filters by name (case-insensitive); " +
			"optional `type` (scene|performer|gallery|image|movie) filters by supported content type. " +
			"Returns id, name, and supported types, plus the total installed count.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"query":{"type":"string"},
				"type":{"type":"string","enum":["scene","performer","gallery","image","movie"]},
				"limit":{"type":"integer","minimum":1,"maximum":200}
			},
			"additionalProperties":false
		}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in struct {
				Query string `json:"query"`
				Type  string `json:"type"`
				Limit int    `json:"limit"`
			}
			if len(raw) > 0 {
				if err := json.Unmarshal(raw, &in); err != nil {
					return "", fmt.Errorf("invalid input: %w", err)
				}
			}
			m := manager.GetInstance()
			if m == nil || m.ScraperCache == nil {
				return "", fmt.Errorf("scraper cache is not available")
			}

			all := m.ScraperCache.ListScrapers(scraper.AllScrapeContentType)
			limit := in.Limit
			if limit <= 0 || limit > 200 {
				limit = 100
			}
			q := strings.ToLower(in.Query)
			typeFilter := strings.ToLower(in.Type)

			out := make([]scraperOut, 0, limit)
			for _, s := range all {
				types := scraperSupportedTypes(s)
				if typeFilter != "" && !hasString(types, typeFilter) {
					continue
				}
				if q != "" && !strings.Contains(strings.ToLower(s.Name), q) {
					continue
				}
				if len(out) < limit {
					out = append(out, scraperOut{ID: s.ID, Name: s.Name, Types: types})
				}
			}
			return resultJSON(map[string]any{
				"total_installed": len(all),
				"returned":        len(out),
				"scrapers":        out,
			})
		},
	}
}

type identifyScenesInput struct {
	SceneIDs            []int    `json:"scene_ids"`
	Sources             []string `json:"sources"`
	SetOrganized        *bool    `json:"set_organized"`
	SkipMultipleMatches *bool    `json:"skip_multiple_matches"`
	FlagMultipleAsTag   string   `json:"flag_multiple_as_tag"`
}

// resolveStashBoxEndpoint matches a user-supplied source (name or endpoint,
// case-insensitive) against the configured stash-box endpoints.
func resolveStashBoxEndpoint(boxes []*models.StashBox, want string) string {
	w := strings.ToLower(strings.TrimSpace(want))
	for _, b := range boxes {
		if strings.ToLower(b.Name) == w || strings.ToLower(b.Endpoint) == w || strings.Contains(strings.ToLower(b.Endpoint), w) {
			return b.Endpoint
		}
	}
	return ""
}

// ensureTagID finds a tag by name (creating it if missing) and returns its id.
func ensureTagID(ctx context.Context, m *manager.Manager, name string) (int, error) {
	repo := m.Repository
	var id int
	err := repo.WithTxn(ctx, func(ctx context.Context) error {
		existing, e := repo.Tag.FindByName(ctx, name, true)
		if e != nil {
			return e
		}
		if existing != nil {
			id = existing.ID
			return nil
		}
		nt := &models.Tag{Name: name}
		if e := repo.Tag.Create(ctx, &models.CreateTagInput{Tag: nt}); e != nil {
			return e
		}
		id = nt.ID
		return nil
	})
	return id, err
}

func identifyScenesTool() *Tool {
	return &Tool{
		Name: "identify_scenes",
		Description: "Bulk-identify scenes against the configured stash-box endpoints (e.g. StashDB, " +
			"ThePornDB), matching by fingerprint and applying metadata. This is the proper mass-tagging " +
			"mechanism — far more capable than the per-scene Tagger view — and runs as a background job. " +
			"Defaults: all scenes; all configured stash-boxes in priority order (first match wins); " +
			"ambiguous multi-match scenes are skipped. Use flag_multiple_as_tag to instead tag skipped " +
			"multi-match scenes (created if missing) for manual review. Check fingerprint_coverage first.",
		Writes: true,
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"scene_ids":{"type":"array","items":{"type":"integer"},"description":"defaults to all scenes"},
				"sources":{"type":"array","items":{"type":"string"},"description":"stash-box names/endpoints in priority order; defaults to all configured"},
				"set_organized":{"type":"boolean"},
				"skip_multiple_matches":{"type":"boolean","description":"default true; ambiguous scenes are skipped"},
				"flag_multiple_as_tag":{"type":"string","description":"tag name to mark skipped multi-match scenes for review"}
			},
			"additionalProperties":false
		}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var in identifyScenesInput
			if len(raw) > 0 {
				if err := json.Unmarshal(raw, &in); err != nil {
					return "", fmt.Errorf("invalid input: %w", err)
				}
			}
			m := manager.GetInstance()
			if m == nil {
				return "", fmt.Errorf("manager is not available")
			}
			boxes := config.GetInstance().GetStashBoxes()
			if len(boxes) == 0 {
				return "", fmt.Errorf("no stash-box endpoints are configured")
			}

			// resolve source endpoints (in priority order)
			var endpoints []string
			if len(in.Sources) == 0 {
				for _, b := range boxes {
					endpoints = append(endpoints, b.Endpoint)
				}
			} else {
				for _, s := range in.Sources {
					ep := resolveStashBoxEndpoint(boxes, s)
					if ep == "" {
						names := make([]string, len(boxes))
						for i, b := range boxes {
							names[i] = b.Name
						}
						return "", fmt.Errorf("unknown stash-box source %q (configured: %s)", s, strings.Join(names, ", "))
					}
					endpoints = append(endpoints, ep)
				}
			}

			opts := &identify.MetadataOptions{}
			if in.SetOrganized != nil {
				opts.SetOrganized = in.SetOrganized
			}
			if in.SkipMultipleMatches != nil {
				opts.SkipMultipleMatches = in.SkipMultipleMatches
			}
			if in.FlagMultipleAsTag != "" {
				tagID, err := ensureTagID(ctx, m, in.FlagMultipleAsTag)
				if err != nil {
					return "", fmt.Errorf("resolving flag tag: %w", err)
				}
				idStr := strconv.Itoa(tagID)
				opts.SkipMultipleMatchTag = &idStr
				skip := true
				opts.SkipMultipleMatches = &skip
			}

			sources := make([]*identify.Source, 0, len(endpoints))
			for i := range endpoints {
				ep := endpoints[i]
				sources = append(sources, &identify.Source{Source: &scraper.Source{StashBoxEndpoint: &ep}})
			}

			identOpts := identify.Options{Sources: sources, Options: opts}
			scope := "all scenes"
			if len(in.SceneIDs) > 0 {
				identOpts.SceneIDs = intsToStrings(in.SceneIDs)
				scope = fmt.Sprintf("%d scene(s)", len(in.SceneIDs))
			}

			// Submit as a background job (independent of this request's lifetime).
			t := manager.CreateIdentifyJob(identOpts)
			jobID := m.JobManager.Add(context.Background(), "Identifying...", t)

			return resultJSON(map[string]any{
				"job_id": jobID, "scope": scope, "sources": endpoints,
				"note": "identification is running in the background; check Settings → Tasks for progress",
			})
		},
	}
}

// ── identify_scenes_fast: batched, parallel-to-job-queue identifier ────────

type identifyFastInput struct {
	Limit         int  `json:"limit"`
	SetOrganized  bool `json:"set_organized"`
	AllowMultiple bool `json:"allow_multiple"`
	AllScenes     bool `json:"all_scenes"`
	DryRun        bool `json:"dry_run"`
}

// summaryRE parses the python tool's final line:
//   summary: matched=180 applied=175 skipped_multiple=4 no_match=1 (dry-run — re-run with --apply)
var summaryRE = regexp.MustCompile(`matched=(\d+)\s+applied=(\d+)\s+skipped_multiple=(\d+)\s+no_match=(\d+)`)

// identifyExternalPath is where the Dockerfile installs the bundled script.
const identifyExternalPath = "/usr/local/bin/identify_external.py"

func identifyScenesFastTool() *Tool {
	return &Tool{
		Name: "identify_scenes_fast",
		Description: "FAST batched scene identification — runs the bundled external identifier in a " +
			"separate process so it does NOT queue behind stash's other tasks (Generate→Phash etc.), " +
			"and uses batched stash-box fingerprint queries (~40× fewer round-trips than native). " +
			"Matches by oshash (no phash required). Default: 200 unorganized scenes, sets organized, " +
			"skips ambiguous multi-matches. Returns summary counts. Use this for SPEED; use " +
			"identify_scenes for native full-fidelity matching when no urgency.",
		Writes: true,
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"limit":{"type":"integer","minimum":1,"maximum":500,"description":"max scenes per run (default 200; capped so the run fits within the chat turn timeout)"},
				"set_organized":{"type":"boolean","description":"mark matched scenes organized (default true)"},
				"allow_multiple":{"type":"boolean","description":"apply first match when a scene has several (default false; skip ambiguous)"},
				"all_scenes":{"type":"boolean","description":"include already-organized scenes too (default false)"},
				"dry_run":{"type":"boolean","description":"preview only, no writes (default false; the tool is invoked deliberately)"}
			},
			"additionalProperties":false
		}`),
		Run: func(ctx context.Context, raw json.RawMessage) (string, error) {
			in := identifyFastInput{Limit: 200, SetOrganized: true}
			if len(raw) > 0 {
				// re-decode to preserve explicit false defaults if caller sent them
				if err := json.Unmarshal(raw, &in); err != nil {
					return "", fmt.Errorf("invalid input: %w", err)
				}
			}
			if in.Limit <= 0 {
				in.Limit = 200
			}
			if in.Limit > 500 {
				in.Limit = 500
			}

			args := []string{
				identifyExternalPath,
				"--stash-url", "http://localhost:9999",
				"--limit", strconv.Itoa(in.Limit),
			}
			if !in.DryRun {
				args = append(args, "--apply")
			}
			if in.SetOrganized {
				args = append(args, "--set-organized")
			}
			if in.AllowMultiple {
				args = append(args, "--allow-multiple")
			}
			if in.AllScenes {
				args = append(args, "--all-scenes")
			}

			// 100s ceiling — well below the 120s LLM client timeout. limit=200 typically
			// finishes in <60s; the ceiling protects against a stuck stash-box endpoint.
			cmdCtx, cancel := context.WithTimeout(ctx, 100*time.Second)
			defer cancel()
			cmd := exec.CommandContext(cmdCtx, "python3", args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			start := time.Now()
			err := cmd.Run()
			elapsed := time.Since(start).Round(time.Millisecond)
			out := stdout.String()

			if err != nil {
				if cmdCtx.Err() == context.DeadlineExceeded {
					return "", fmt.Errorf("identifier timed out after %s — try a smaller limit", elapsed)
				}
				// keep the response actionable: tail of stderr, then exit error
				return "", fmt.Errorf("identifier failed (%s) in %s: %s", err, elapsed, truncate(stderr.String(), 300))
			}

			// parse the final summary line; fall back to truncated stdout if missing
			result := map[string]any{
				"limit":   in.Limit,
				"dry_run": in.DryRun,
				"elapsed": elapsed.String(),
			}
			if m := summaryRE.FindStringSubmatch(out); m != nil {
				result["matched"], _ = strconv.Atoi(m[1])
				result["applied"], _ = strconv.Atoi(m[2])
				result["skipped_multiple"], _ = strconv.Atoi(m[3])
				result["no_match"], _ = strconv.Atoi(m[4])
			} else {
				// unusual output — surface a tail so we can debug from the tool result
				result["raw_tail"] = tailLines(out, 8)
			}
			if in.DryRun {
				result["note"] = "dry-run; nothing was written. Re-run without dry_run to apply."
			}
			return resultJSON(result)
		},
	}
}

func tailLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

func reloadScrapersTool() *Tool {
	return &Tool{
		Name: "reload_scrapers",
		Description: "Reload all metadata scrapers from disk (e.g. after new scraper files were added under " +
			"the config scrapers/ directory). Hot reload — no restart. Returns the number of scrapers loaded.",
		Writes: true,
		Schema: json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`),
		Run: func(ctx context.Context, _ json.RawMessage) (string, error) {
			m := manager.GetInstance()
			if m == nil {
				return "", fmt.Errorf("manager is not available")
			}
			m.RefreshScraperCache()
			loaded := 0
			if m.ScraperCache != nil {
				loaded = len(m.ScraperCache.ListScrapers(scraper.AllScrapeContentType))
			}
			return resultJSON(map[string]any{"reloaded": true, "scrapers_loaded": loaded})
		},
	}
}
