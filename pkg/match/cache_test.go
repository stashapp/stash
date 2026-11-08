package match

import (
	"context"
	"slices"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
)

func TestFirstTwoRunesLower(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"alice smith", "al"},
		{"ALICE", "al"},
		{"Àbc", "àb"},
		{"伏字 name", "伏字"},
		{"ab", "ab"},
		{"a", ""},       // single rune -> no prefix
		{"", ""},        // empty -> no prefix
		{"X Man", "x "}, // space is preserved in 2-rune prefix
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			if got := firstTwoRunesLower(tt.in); got != tt.want {
				t.Errorf("firstTwoRunesLower(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCacheNameRegexpCaches(t *testing.T) {
	t.Parallel()

	c := &Cache{}
	r1 := c.nameRegexp("alice smith", true)
	r2 := c.nameRegexp("alice smith", true)
	if r1 != r2 {
		t.Error("expected cached regexp to be reused across calls")
	}

	// Different useUnicode flag -> different cached regexp.
	r3 := c.nameRegexp("alice smith", false)
	if r3 == r1 {
		t.Error("expected ASCII and unicode variants to be distinct cached entries")
	}

	// Nil cache must still return a valid regexp, just uncached.
	var nilCache *Cache
	if got := nilCache.nameRegexp("alice smith", true); got == nil {
		t.Error("nil cache should still return a regexp")
	}
}

func TestPreloadPerformersBuildsIndex(t *testing.T) {
	t.Parallel()

	alice := &models.Performer{ID: 1, Name: "Alice Smith"}
	bob := &models.Performer{ID: 2, Name: "bob jones"}
	xman := &models.Performer{ID: 3, Name: "X Man"}
	ignored := &models.Performer{ID: 4, Name: "ignored", IgnoreAutoTag: true}

	performers := []*models.Performer{alice, bob, xman, ignored}
	db := mocks.NewDatabase()
	primePerformerMock(db.Performer, performers)

	c := &Cache{}
	if err := c.PreloadPerformers(context.Background(), db.Performer); err != nil {
		t.Fatalf("PreloadPerformers: %v", err)
	}

	// allPerformers excludes IgnoreAutoTag=true.
	if got := len(c.allPerformers); got != 3 {
		t.Errorf("allPerformers len = %d, want 3 (ignored must be excluded)", got)
	}

	// Prefix "al" -> alice, "bo" -> bob, "x " -> xman.
	assertBucket := func(prefix string, wantIDs []int) {
		t.Helper()
		var gotIDs []int
		for _, p := range c.performerByPrefix[prefix] {
			gotIDs = append(gotIDs, p.ID)
		}
		slices.Sort(gotIDs)
		if !slices.Equal(gotIDs, wantIDs) {
			t.Errorf("bucket %q = %v, want %v", prefix, gotIDs, wantIDs)
		}
	}
	assertBucket("al", []int{1})
	assertBucket("bo", []int{2})
	assertBucket("x ", []int{3})

	// Single-letter-first-word performer must also be in alwaysCheck.
	var alwaysIDs []int
	for _, p := range c.performerAlwaysCheck {
		alwaysIDs = append(alwaysIDs, p.ID)
	}
	if !slices.Equal(alwaysIDs, []int{3}) {
		t.Errorf("alwaysCheck IDs = %v, want [3]", alwaysIDs)
	}

	// Idempotent: second call is a no-op.
	if err := c.PreloadPerformers(context.Background(), db.Performer); err != nil {
		t.Fatalf("second PreloadPerformers: %v", err)
	}
	if got := len(c.allPerformers); got != 3 {
		t.Errorf("after idempotent call allPerformers len = %d, want 3", got)
	}
}

func TestPreloadStudiosIndexesAliasPrefixes(t *testing.T) {
	t.Parallel()

	// Name "Acme" shares no prefix with alias "Widgets" — both must be
	// reachable by their own 2-rune prefix.
	s := &models.Studio{ID: 1, Name: "Acme Corp"}
	ignored := &models.Studio{ID: 2, Name: "ignored", IgnoreAutoTag: true}

	db := mocks.NewDatabase()
	primeStudioMock(db.Studio, []*models.Studio{s, ignored}, map[int][]string{1: {"Widgets Inc"}})

	c := &Cache{}
	if err := c.PreloadStudios(context.Background(), db.Studio); err != nil {
		t.Fatalf("PreloadStudios: %v", err)
	}

	if got := len(c.allStudios); got != 1 {
		t.Errorf("allStudios len = %d, want 1 (ignored must be excluded)", got)
	}

	// "ac" bucket has the studio (via name), "wi" bucket has it (via alias).
	if len(c.studioByPrefix["ac"]) != 1 || c.studioByPrefix["ac"][0].Studio.ID != 1 {
		t.Errorf("bucket 'ac' should hold studio 1, got %+v", c.studioByPrefix["ac"])
	}
	if len(c.studioByPrefix["wi"]) != 1 || c.studioByPrefix["wi"][0].Studio.ID != 1 {
		t.Errorf("bucket 'wi' should hold studio 1, got %+v", c.studioByPrefix["wi"])
	}
}

func TestPreloadStudiosDedupsSharedPrefix(t *testing.T) {
	t.Parallel()

	// Name and two aliases all share prefix "pr"; the bucket must contain
	// the studio exactly once.
	s := &models.Studio{ID: 1, Name: "Primary"}
	db := mocks.NewDatabase()
	primeStudioMock(db.Studio, []*models.Studio{s}, map[int][]string{1: {"Primary Nick", "Primary Alt"}})

	c := &Cache{}
	if err := c.PreloadStudios(context.Background(), db.Studio); err != nil {
		t.Fatal(err)
	}

	if got := len(c.studioByPrefix["pr"]); got != 1 {
		t.Errorf("bucket 'pr' should have 1 entry, got %d", got)
	}
}

func TestPreloadTagsIndexesAliasPrefixes(t *testing.T) {
	t.Parallel()

	db := mocks.NewDatabase()
	primeTagMock(db.Tag, []*models.Tag{{ID: 1, Name: "documentary"}}, map[int][]string{1: {"film"}})

	c := &Cache{}
	if err := c.PreloadTags(context.Background(), db.Tag); err != nil {
		t.Fatal(err)
	}

	if len(c.tagByPrefix["do"]) != 1 || c.tagByPrefix["do"][0].Tag.ID != 1 {
		t.Errorf("bucket 'do' should hold tag 1")
	}
	if len(c.tagByPrefix["fi"]) != 1 || c.tagByPrefix["fi"][0].Tag.ID != 1 {
		t.Errorf("bucket 'fi' should hold tag 1 (via alias)")
	}
}

func TestCandidateLookupDedupesAcrossPathWords(t *testing.T) {
	t.Parallel()

	// A performer with name "alabama" falls in bucket "al". If a path has
	// two words that both map to bucket "al" (e.g., from separate tokens),
	// the candidate must appear exactly once.
	p := &models.Performer{ID: 1, Name: "alabama"}
	db := mocks.NewDatabase()
	primePerformerMock(db.Performer, []*models.Performer{p})

	c := &Cache{}
	if err := c.PreloadPerformers(context.Background(), db.Performer); err != nil {
		t.Fatal(err)
	}

	got := c.performerCandidates([]string{"al", "AL", "al"}) // same bucket three times
	if len(got) != 1 {
		t.Errorf("expected 1 candidate after dedup, got %d: %v", len(got), got)
	}
}
