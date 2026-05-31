package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
