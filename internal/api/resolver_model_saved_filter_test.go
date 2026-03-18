package api

import (
	"context"
	"github.com/stashapp/stash/pkg/models"
	"testing"
)

// We verify the `LabelMapping` function handles parsing the interface mapping without panic and extracts correct lists.
func TestSavedFilterLabelMappingEmpty(t *testing.T) {
	// Basic instantiation to just ensure it does not panic and returns empty correctly.
	resolver := &savedFilterResolver{}

	obj := &models.SavedFilter{
		ObjectFilter: nil,
	}

	mapping, err := resolver.LabelMapping(context.Background(), obj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mapping == nil || mapping.Tags != nil {
		t.Errorf("expected empty mapping, got %v", mapping)
	}
}
