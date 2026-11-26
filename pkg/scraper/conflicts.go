package scraper

import (
	"path/filepath"
	"strings"
)

// URLPatternSource is a URL pattern declared by a particular scraper
type URLPatternSource struct {
	ScraperID   string
	ScraperName string
	Pattern     string
}

// URLConflict describes two installed scrapers whose URL patterns for the
// same content type overlap, such that Cache.ScrapeURL cannot deterministically
// choose between them. The pair is unordered.
type URLConflict struct {
	Type    ScrapeContentType
	Scraper URLPatternSource
	Other   URLPatternSource
}

// IDFromFile returns the ID of the scraper defined by the given file path,
// mirroring how loadConfigFromYAMLFile assigns IDs. Returns false if the file
// is not a scraper definition.
func IDFromFile(path string) (string, bool) {
	if filepath.Ext(path) != ".yml" {
		return "", false
	}
	return strings.TrimSuffix(filepath.Base(path), ".yml"), true
}

// urlsForType returns the URL patterns declared by spec for the given content
// type. Movie is not handled since it shares Group's patterns.
func urlsForType(spec Scraper, ty ScrapeContentType) []string {
	var s *ScraperSpec

	switch ty {
	case ScrapeContentTypeScene:
		s = spec.Scene
	case ScrapeContentTypePerformer:
		s = spec.Performer
	case ScrapeContentTypeGallery:
		s = spec.Gallery
	case ScrapeContentTypeImage:
		s = spec.Image
	case ScrapeContentTypeGroup:
		s = spec.Group
	}

	if s == nil {
		return nil
	}

	return s.Urls
}

// FindURLConflicts finds pairs of installed scrapers whose URL patterns for
// the same content type overlap: one pattern is equal to, or a substring of,
// another pattern belonging to a different scraper. Any URL matching the
// longer pattern then also matches the shorter one.
//
// Patterns that overlap only on some URLs (neither contains the other) are not
// detected.
//
// Each conflicting pattern pair is reported once. A pattern conflicting with
// several scrapers produces one conflict per scraper.
func (c Cache) FindURLConflicts() []*URLConflict {
	var ret []*URLConflict

	for _, ty := range AllScrapeContentType {
		// Movie is an alias of Group
		if ty == ScrapeContentTypeMovie {
			continue
		}

		var patterns []URLPatternSource
		for _, s := range c.scrapers {
			spec := s.spec()
			for _, u := range urlsForType(spec, ty) {
				patterns = append(patterns, URLPatternSource{ScraperID: spec.ID, ScraperName: spec.Name, Pattern: u})
			}
		}

		for _, a := range patterns {
			for _, b := range patterns {
				if a.ScraperID == b.ScraperID || !strings.Contains(a.Pattern, b.Pattern) {
					continue
				}

				// identical patterns match both ways; keep only one ordering
				if a.Pattern == b.Pattern && a.ScraperID > b.ScraperID {
					continue
				}

				ret = append(ret, &URLConflict{
					Type:    ty,
					Scraper: a,
					Other:   b,
				})
			}
		}
	}

	return ret
}
