package manager

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

func TestPerformerStoredIDUsesScrapedStoredID(t *testing.T) {
	storedID := "42"
	scraped := &models.ScrapedPerformer{StoredID: &storedID}
	existing := &models.Performer{ID: 7}

	got, err := performerStoredID(scraped, existing)
	if err != nil {
		t.Fatalf("performerStoredID returned error: %v", err)
	}
	if got != 42 {
		t.Fatalf("expected scraped stored id 42, got %d", got)
	}
}

func TestPerformerStoredIDFallsBackToExistingPerformer(t *testing.T) {
	scraped := &models.ScrapedPerformer{}
	existing := &models.Performer{ID: 7}

	got, err := performerStoredID(scraped, existing)
	if err != nil {
		t.Fatalf("performerStoredID returned error: %v", err)
	}
	if got != 7 {
		t.Fatalf("expected existing performer id 7, got %d", got)
	}
}

func TestPreserveExistingPerformerNameAddsRemoteNameAsAlias(t *testing.T) {
	remoteName := "Azusa Maki"
	remoteAliases := "眞木あずさ, 真木あずさ"
	scraped := &models.ScrapedPerformer{
		Name:    &remoteName,
		Aliases: &remoteAliases,
	}
	existing := &models.Performer{
		ID:   7,
		Name: "眞木あずさ",
	}
	partial := scraped.ToPartial("https://stashdb.org/graphql", map[string]bool{}, map[string]bool{}, nil)

	preserveExistingPerformerName(&partial, scraped, existing)

	if partial.Name.Set {
		t.Fatalf("expected performer name to be preserved, got update to %q", partial.Name.Value)
	}
	if partial.Aliases == nil {
		t.Fatal("expected aliases update")
	}
	if partial.Aliases.Mode != models.RelationshipUpdateModeSet {
		t.Fatalf("expected aliases set mode, got %s", partial.Aliases.Mode)
	}

	wantAliases := map[string]bool{
		"眞木あずさ":      true,
		"真木あずさ":      true,
		"Azusa Maki": true,
	}
	for _, alias := range partial.Aliases.Values {
		delete(wantAliases, alias)
	}
	if len(wantAliases) != 0 {
		t.Fatalf("missing aliases after preserving performer name: %v", wantAliases)
	}
}
