package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

// resolveSavedFilter is a helper that looks up a saved filter by ID, enforces that it matches
// the expected mode (e.g. SCENES), and returns the populated object filter (e.g., SceneFilterType)
// and an updated find filter that merges the saved one with the overrides.
func (r *queryResolver) resolveSavedFilter(ctx context.Context, savedFilterID string, mode models.FilterMode, outObjectFilter interface{}, currentFindFilter *models.FindFilterType) (*models.FindFilterType, error) {
	id, err := strconv.Atoi(savedFilterID)
	if err != nil {
		return nil, fmt.Errorf("invalid saved_filter_id: %w", err)
	}

	var savedFilter *models.SavedFilter
	err = r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		savedFilter, err = r.repository.SavedFilter.Find(ctx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch saved filter: %w", err)
	}
	if savedFilter == nil {
		return nil, fmt.Errorf("saved filter %s not found", savedFilterID)
	}

	if savedFilter.Mode != mode {
		return nil, fmt.Errorf("saved filter is of mode %s, but expected %s", savedFilter.Mode, mode)
	}

	if savedFilter.ObjectFilter != nil {
		b, err := json.Marshal(savedFilter.ObjectFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal object filter: %w", err)
		}
		if err := json.Unmarshal(b, outObjectFilter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal object filter into target struct: %w", err)
		}
	}

	// Merge find filter
	finalFindFilter := savedFilter.FindFilter
	if finalFindFilter == nil {
		finalFindFilter = &models.FindFilterType{}
	}

	if currentFindFilter != nil {
		if currentFindFilter.Q != nil {
			finalFindFilter.Q = currentFindFilter.Q
		}
		if currentFindFilter.Page != nil {
			finalFindFilter.Page = currentFindFilter.Page
		}
		if currentFindFilter.PerPage != nil {
			finalFindFilter.PerPage = currentFindFilter.PerPage
		}
		if currentFindFilter.Sort != nil {
			finalFindFilter.Sort = currentFindFilter.Sort
		}
		if currentFindFilter.Direction != nil {
			finalFindFilter.Direction = currentFindFilter.Direction
		}
	}

	return finalFindFilter, nil
}
