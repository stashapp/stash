package stashbox

import (
	"testing"

	"github.com/stashapp/stash/pkg/stashbox/graphql"
)

func TestFindMatchingPerformerFragmentMatchesAlias(t *testing.T) {
	name := "眞木あずさ"
	performers := []*graphql.PerformerFragment{
		{
			ID:      "first",
			Name:    "Someone Else",
			Aliases: []string{"別人"},
		},
		{
			ID:      "second",
			Name:    "Azusa Maki",
			Aliases: []string{"眞木あずさ", "真木あずさ"},
		},
	}

	got := findMatchingPerformerFragment(performers, name)

	if got == nil {
		t.Fatal("expected performer alias match")
	}
	if got.ID != "second" {
		t.Fatalf("expected second performer, got %q", got.ID)
	}
}

func TestFindMatchingPerformerFragmentPrefersPrimaryName(t *testing.T) {
	name := "Azusa Maki"
	performers := []*graphql.PerformerFragment{
		{
			ID:      "alias",
			Name:    "Someone Else",
			Aliases: []string{"Azusa Maki"},
		},
		{
			ID:   "primary",
			Name: "Azusa Maki",
		},
	}

	got := findMatchingPerformerFragment(performers, name)

	if got == nil {
		t.Fatal("expected performer primary name match")
	}
	if got.ID != "primary" {
		t.Fatalf("expected primary performer, got %q", got.ID)
	}
}

func TestFindMatchingPerformerFragmentUsesSingleCandidate(t *testing.T) {
	name := "鈴音りおな"
	performers := []*graphql.PerformerFragment{
		{
			ID:      "single",
			Name:    "Suzune Riona",
			Aliases: []string{"鈴木りおな"},
		},
	}

	got := findMatchingPerformerFragment(performers, name)

	if got == nil {
		t.Fatal("expected single performer candidate match")
	}
	if got.ID != "single" {
		t.Fatalf("expected single performer, got %q", got.ID)
	}
}
