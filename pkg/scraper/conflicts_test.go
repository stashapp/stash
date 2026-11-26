package scraper

import (
	"sort"
	"testing"
)

func byURL(urls []string) []*ByURLDefinition {
	var ret []*ByURLDefinition
	for _, u := range urls {
		ret = append(ret, &ByURLDefinition{URL: []string{u}})
	}
	return ret
}

// test scrapers use their ID as their name
func src(id, pattern string) URLPatternSource {
	return URLPatternSource{ScraperID: id, ScraperName: id, Pattern: pattern}
}

func makeTestScraper(id string, sceneURLs []string, performerURLs []string) scraper {
	return scraperFromDefinition(Definition{
		ID:             id,
		Name:           id,
		SceneByURL:     byURL(sceneURLs),
		PerformerByURL: byURL(performerURLs),
	}, nil)
}

func sortConflicts(c []*URLConflict) {
	key := func(x *URLConflict) string {
		return x.Scraper.ScraperID + "\x00" + x.Scraper.Pattern + "\x00" + x.Other.ScraperID
	}
	sort.Slice(c, func(i, j int) bool { return key(c[i]) < key(c[j]) })
}

func assertConflicts(t *testing.T, got []*URLConflict, want []URLConflict) {
	t.Helper()
	sortConflicts(got)

	if len(got) != len(want) {
		t.Fatalf("FindURLConflicts() returned %d conflicts, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if *got[i] != want[i] {
			t.Errorf("FindURLConflicts()[%d] = %+v, want %+v", i, *got[i], want[i])
		}
	}
}

func TestCache_FindURLConflicts_Subsumption(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com"}, nil),
			"b": makeTestScraper("b", []string{"example.com/scene"}, nil),
		},
	}

	// the longer pattern contains the shorter one, so it is reported first
	assertConflicts(t, c.FindURLConflicts(), []URLConflict{
		{
			Type:    ScrapeContentTypeScene,
			Scraper: src("b", "example.com/scene"),
			Other:   src("a", "example.com"),
		},
	})
}

func TestCache_FindURLConflicts_DisjointPatternsNoConflict(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com/scenes/"}, nil),
			"b": makeTestScraper("b", []string{"example.com/movies/"}, nil),
		},
	}

	assertConflicts(t, c.FindURLConflicts(), nil)
}

func TestCache_FindURLConflicts_EqualPatternsReportedOnce(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com"}, nil),
			"b": makeTestScraper("b", []string{"example.com"}, nil),
		},
	}

	assertConflicts(t, c.FindURLConflicts(), []URLConflict{
		{
			Type:    ScrapeContentTypeScene,
			Scraper: src("a", "example.com"),
			Other:   src("b", "example.com"),
		},
	})
}

func TestCache_FindURLConflicts_ThreeScrapersSamePattern(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com"}, nil),
			"b": makeTestScraper("b", []string{"example.com"}, nil),
			"c": makeTestScraper("c", []string{"example.com"}, nil),
		},
	}

	scene := ScrapeContentTypeScene
	assertConflicts(t, c.FindURLConflicts(), []URLConflict{
		{Type: scene, Scraper: src("a", "example.com"), Other: src("b", "example.com")},
		{Type: scene, Scraper: src("a", "example.com"), Other: src("c", "example.com")},
		{Type: scene, Scraper: src("b", "example.com"), Other: src("c", "example.com")},
	})
}

func TestCache_FindURLConflicts_SameScraperNoSelfConflict(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com", "example.com/scene"}, nil),
		},
	}

	assertConflicts(t, c.FindURLConflicts(), nil)
}

func TestCache_FindURLConflicts_ScopedPerContentType(t *testing.T) {
	// ScrapeURL dispatches independently per content type
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com"}, nil),
			"b": makeTestScraper("b", nil, []string{"example.com"}),
		},
	}

	assertConflicts(t, c.FindURLConflicts(), nil)
}

func TestCache_FindURLConflicts_MovieAndGroupShareBucket(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": scraperFromDefinition(Definition{ID: "a", Name: "a", MovieByURL: byURL([]string{"example.com"})}, nil),
			"b": scraperFromDefinition(Definition{ID: "b", Name: "b", GroupByURL: byURL([]string{"example.com"})}, nil),
		},
	}

	assertConflicts(t, c.FindURLConflicts(), []URLConflict{
		{
			Type:    ScrapeContentTypeGroup,
			Scraper: src("a", "example.com"),
			Other:   src("b", "example.com"),
		},
	})
}

func TestIDFromFile(t *testing.T) {
	tests := []struct {
		path   string
		wantID string
		wantOK bool
	}{
		{"MojoHost.yml", "MojoHost", true},
		{"MojoHost/MojoHost.yml", "MojoHost", true},
		{"MojoHost/MojoHost.py", "", false},
		{"MojoHost/requirements.txt", "", false},
	}

	for _, tt := range tests {
		id, ok := IDFromFile(tt.path)
		if id != tt.wantID || ok != tt.wantOK {
			t.Errorf("IDFromFile(%q) = %q, %v; want %q, %v", tt.path, id, ok, tt.wantID, tt.wantOK)
		}
	}
}
