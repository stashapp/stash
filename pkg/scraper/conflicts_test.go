package scraper

import (
	"sort"
	"testing"
)

func makeTestScraper(id string, sceneURLs []string, performerURLs []string) scraper {
	def := Definition{
		ID:   id,
		Name: id,
	}

	for _, u := range sceneURLs {
		def.SceneByURL = append(def.SceneByURL, &ByURLDefinition{URL: []string{u}})
	}
	for _, u := range performerURLs {
		def.PerformerByURL = append(def.PerformerByURL, &ByURLDefinition{URL: []string{u}})
	}

	return scraperFromDefinition(def, nil)
}

func sortConflicts(c []ScraperURLConflict) {
	sort.Slice(c, func(i, j int) bool {
		if c[i].ScraperID != c[j].ScraperID {
			return c[i].ScraperID < c[j].ScraperID
		}
		if c[i].Pattern != c[j].Pattern {
			return c[i].Pattern < c[j].Pattern
		}
		return c[i].OtherScraperID < c[j].OtherScraperID
	})
}

func TestCache_FindURLConflicts_Subsumption(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com"}, nil),
			"b": makeTestScraper("b", []string{"example.com/scene"}, nil),
		},
	}

	got := c.FindURLConflicts()
	sortConflicts(got)

	// the longer, more specific pattern ("example.com/scene") contains the
	// shorter, broader one ("example.com"), so it is reported as the primary
	// side of the pair: consumers should treat the pair as unordered
	want := []ScraperURLConflict{
		{
			Type:           ScrapeContentTypeScene,
			ScraperID:      "b",
			Pattern:        "example.com/scene",
			OtherScraperID: "a",
			OtherPattern:   "example.com",
		},
	}

	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("FindURLConflicts() = %+v, want %+v", got, want)
	}
}

func TestCache_FindURLConflicts_DisjointPatternsNoConflict(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com/scenes/"}, nil),
			"b": makeTestScraper("b", []string{"example.com/movies/"}, nil),
		},
	}

	got := c.FindURLConflicts()
	if len(got) != 0 {
		t.Errorf("FindURLConflicts() = %+v, want no conflicts for disjoint patterns on the same domain", got)
	}
}

func TestCache_FindURLConflicts_EqualPatternsReportedOnce(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com"}, nil),
			"b": makeTestScraper("b", []string{"example.com"}, nil),
		},
	}

	got := c.FindURLConflicts()
	if len(got) != 1 {
		t.Fatalf("FindURLConflicts() returned %d conflicts for identical patterns, want exactly 1", len(got))
	}
	if got[0].ScraperID != "a" || got[0].OtherScraperID != "b" {
		t.Errorf("FindURLConflicts() = %+v, want scraperID a, otherScraperID b (canonical ordering)", got[0])
	}
}

func TestCache_FindURLConflicts_SameScraperNoSelfConflict(t *testing.T) {
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com", "example.com/scene"}, nil),
		},
	}

	got := c.FindURLConflicts()
	if len(got) != 0 {
		t.Errorf("FindURLConflicts() = %+v, want no conflicts within a single scraper's own patterns", got)
	}
}

func TestCache_FindURLConflicts_ScopedPerContentType(t *testing.T) {
	// a scene scraper and a performer scraper on the same domain is not a
	// conflict, since ScrapeURL dispatches independently per content type
	c := Cache{
		scrapers: map[string]scraper{
			"a": makeTestScraper("a", []string{"example.com"}, nil),
			"b": makeTestScraper("b", nil, []string{"example.com"}),
		},
	}

	got := c.FindURLConflicts()
	if len(got) != 0 {
		t.Errorf("FindURLConflicts() = %+v, want no conflicts across different content types", got)
	}
}
