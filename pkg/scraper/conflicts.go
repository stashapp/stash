package scraper

import "strings"

// ScraperURLConflict describes two installed scrapers whose URL patterns for the
// same content type overlap, such that Cache.ScrapeURL cannot deterministically
// choose between them
type ScraperURLConflict struct {
	Type ScrapeContentType `json:"type"`

	ScraperID string `json:"scraperID"`
	Pattern   string `json:"pattern"`

	OtherScraperID string `json:"otherScraperID"`
	OtherPattern   string `json:"otherPattern"`
}

// urlsForType returns the URL patterns declared by spec for the given content
// type. Movie and Group are treated as the same bucket, since they are
// matched identically by Definition.matchesURL
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
	case ScrapeContentTypeGroup, ScrapeContentTypeMovie:
		s = spec.Group
	}

	if s == nil {
		return nil
	}

	return s.Urls
}

// FindURLConflicts finds pairs of installed scrapers whose URL patterns for
// the same content type overlap in a way that makes ScrapeURL's dispatch
// ambiguous: one pattern is equal to, or a substring of, another pattern
// belonging to a different scraper
//
// Each conflicting pattern pair is reported once
func (c Cache) FindURLConflicts() []ScraperURLConflict {
	var ret []ScraperURLConflict

	tys := make([]ScrapeContentType, 0, len(AllScrapeContentType))
	for _, ty := range AllScrapeContentType {
		if ty != ScrapeContentTypeMovie {
			tys = append(tys, ty)
		}
	}

	type patternSource struct {
		scraperID string
		pattern   string
	}

	for _, ty := range tys {
		var patterns []patternSource
		for _, s := range c.scrapers {
			spec := s.spec()
			for _, u := range urlsForType(spec, ty) {
				patterns = append(patterns, patternSource{scraperID: spec.ID, pattern: u})
			}
		}

		for i, a := range patterns {
			for j, b := range patterns {
				if i == j || a.scraperID == b.scraperID {
					continue
				}

				if !strings.Contains(a.pattern, b.pattern) {
					continue
				}

				// deliberately compare 'greater than' to guarantee a total ordering
				if a.pattern == b.pattern && a.scraperID > b.scraperID {
					continue
				}

				ret = append(ret, ScraperURLConflict{
					Type:           ty,
					ScraperID:      a.scraperID,
					Pattern:        a.pattern,
					OtherScraperID: b.scraperID,
					OtherPattern:   b.pattern,
				})
			}
		}
	}

	return ret
}
