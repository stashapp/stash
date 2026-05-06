package api

import (
	"fmt"

	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

// TODO - apply handleIDs to other resolvers that accept ID lists

// handleIDList validates and converts a list of string IDs to integers
func handleIDList(idList []string, field string) ([]int, error) {
	if err := validateIDList(idList); err != nil {
		return nil, fmt.Errorf("validating %s: %w", field, err)
	}

	ids, err := stringslice.StringSliceToIntSlice(idList)
	if err != nil {
		return nil, fmt.Errorf("converting %s: %w", field, err)
	}

	return ids, nil
}

// validateIDList returns an error if there are any duplicate ids in the list
func validateIDList(ids []string) error {
	seen := make(map[string]struct{})
	for _, id := range ids {
		if _, exists := seen[id]; exists {
			return fmt.Errorf("duplicate id found: %s", id)
		}
		seen[id] = struct{}{}
	}
	return nil
}
