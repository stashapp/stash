package models

import (
	"fmt"
	"strings"
)

type StashBoxes []*StashBox

func (sb StashBoxes) ResolveStashBox(source ScraperSourceInput) (*StashBox, error) {
	if source.StashBoxIndex != nil {
		index := source.StashBoxIndex
		if *index < 0 || *index >= len(sb) {
			return nil, fmt.Errorf("%w: invalid stash_box_index: %d", ErrScraperSource, index)
		}

		return sb[*index], nil
	}

	if source.StashBoxEndpoint != nil {
		var ret *StashBox
		endpoint := *source.StashBoxEndpoint
		for _, b := range sb {
			if strings.EqualFold(endpoint, b.Endpoint) {
				ret = b
			}
		}

		if ret == nil {
			return nil, fmt.Errorf(`%w: stash-box with endpoint "%s"`, ErrNotFound, endpoint)
		}

		return ret, nil
	}

	// neither stash-box inputs were provided, so assume it is a scraper

	return nil, nil
}
