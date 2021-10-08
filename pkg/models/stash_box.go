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
			return nil, fmt.Errorf("invalid stash_box_index: %d", index)
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
			return nil, fmt.Errorf(`stash-box with endpoint "%s" not found`, endpoint)
		}

		return ret, nil
	}

	return nil, nil
}
