package scraper

import (
	"context"
	"regexp"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

func postProcessTags(ctx context.Context, tqb models.TagQueryer, scrapedTags []*models.ScrapedTag) (ret []*models.ScrapedTag, err error) {
	ret = make([]*models.ScrapedTag, 0, len(scrapedTags))

	for _, t := range scrapedTags {
		// Pass empty string for endpoint since this is used by general scrapers, not just stash-box
		err := match.ScrapedTag(ctx, tqb, t, "")
		if err != nil {
			return nil, err
		}
		ret = append(ret, t)
	}

	return ret, err
}

// FilterTags removes tags matching excluded tag patterns from the list of scraped tags
// It returns the filtered list of tags and a list of the excluded tags
func FilterTags(excludeRegexps []*regexp.Regexp, tags []*models.ScrapedTag) (newTags []*models.ScrapedTag, ignoredTags []string) {
	if len(excludeRegexps) == 0 {
		return tags, nil
	}

	newTags = make([]*models.ScrapedTag, 0, len(tags))

	for _, t := range tags {
		ignore := false
		for _, reg := range excludeRegexps {
			if reg.MatchString(strings.ToLower(t.Name)) {
				ignore = true
				ignoredTags = sliceutil.AppendUnique(ignoredTags, t.Name)
				break
			}
		}

		if !ignore {
			newTags = append(newTags, t)
		}
	}

	return newTags, ignoredTags
}

// CompileExclusionRegexps compiles a list of tag exclusion patterns into a list of regular expressions
func CompileExclusionRegexps(patterns []string) []*regexp.Regexp {
	excludePatterns := patterns
	var excludeRegexps []*regexp.Regexp

	for _, excludePattern := range excludePatterns {
		reg, err := regexp.Compile(strings.ToLower(excludePattern))
		if err != nil {
			logger.Errorf("Invalid tag exclusion pattern: %v", err)
		} else {
			excludeRegexps = append(excludeRegexps, reg)
		}
	}

	return excludeRegexps
}

// LogIgnoredTags logs the list of ignored tags
func LogIgnoredTags(ignoredTags []string) {
	if len(ignoredTags) > 0 {
		logger.Debugf("Tags ignored for matching exclusion patterns: %s", strings.Join(ignoredTags, ", "))
	}
}
