package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/scraper"
)

// RegisterScraperTools adds scraper-management tools. They use the global manager
// singleton (the same way the GraphQL resolvers reach the scraper cache), so they
// need no Deps wiring.
func RegisterScraperTools(reg *Registry) {
	reg.Register(listScrapersTool())
	reg.Register(reloadScrapersTool())
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
