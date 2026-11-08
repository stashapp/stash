package match

import (
	"context"
	"slices"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/mock"
)

// Path-matching semantic tests that lock in the behavior of
// PathTo{Performers,Studio,Tags} via the generated testify mocks in
// pkg/models/mocks. These are the regression guard when the candidate-
// lookup strategy changes (e.g., replacing the SQL prefilter with an
// in-memory matcher): each case runs against both cache=nil and a
// preloaded cache, asserting identical output.

// --- mock setup helpers ---

// preloadFilter matches the filter PreloadX passes: IgnoreAutoTag=false.
// singleLetterFilter matches the filter the single-letter-cache path
// passes: a regex in Name. Keeping them disjoint means testify will
// route each Query call to the right stub regardless of declaration
// order.
func performerPreloadFilter() interface{} {
	return mock.MatchedBy(func(f *models.PerformerFilterType) bool {
		return f != nil && f.IgnoreAutoTag != nil && !*f.IgnoreAutoTag
	})
}
func performerSingleLetterFilter() interface{} {
	return mock.MatchedBy(func(f *models.PerformerFilterType) bool {
		return f != nil && f.Name != nil
	})
}
func studioPreloadFilter() interface{} {
	return mock.MatchedBy(func(f *models.StudioFilterType) bool {
		return f != nil && f.IgnoreAutoTag != nil && !*f.IgnoreAutoTag
	})
}
func studioSingleLetterFilter() interface{} {
	return mock.MatchedBy(func(f *models.StudioFilterType) bool {
		return f != nil && f.Name != nil
	})
}
func tagPreloadFilter() interface{} {
	return mock.MatchedBy(func(f *models.TagFilterType) bool {
		return f != nil && f.IgnoreAutoTag != nil && !*f.IgnoreAutoTag
	})
}
func tagSingleLetterFilter() interface{} {
	return mock.MatchedBy(func(f *models.TagFilterType) bool {
		return f != nil && f.Name != nil
	})
}

// primePerformerMock sets up a PerformerReaderWriter to serve both the
// no-preload path (QueryForAutoTag returns all non-ignored; single-letter
// Query returns nothing) and the preload path (Query with IgnoreAutoTag
// filter returns all non-ignored). All expectations are .Maybe() because
// which ones fire depends on whether the test passes a cache.
func primePerformerMock(m *mocks.PerformerReaderWriter, performers []*models.Performer) {
	var nonIgnored []*models.Performer
	for _, p := range performers {
		if !p.IgnoreAutoTag {
			nonIgnored = append(nonIgnored, p)
		}
	}
	m.On("QueryForAutoTag", mock.Anything, mock.Anything).Return(nonIgnored, nil).Maybe()
	m.On("Query", mock.Anything, performerPreloadFilter(), mock.Anything).Return(nonIgnored, len(nonIgnored), nil).Maybe()
	m.On("Query", mock.Anything, performerSingleLetterFilter(), mock.Anything).Return(nil, 0, nil).Maybe()
}

func primeStudioMock(m *mocks.StudioReaderWriter, studios []*models.Studio, aliases map[int][]string) {
	var nonIgnored []*models.Studio
	for _, s := range studios {
		if !s.IgnoreAutoTag {
			nonIgnored = append(nonIgnored, s)
		}
	}
	m.On("QueryForAutoTag", mock.Anything, mock.Anything).Return(nonIgnored, nil).Maybe()
	m.On("Query", mock.Anything, studioPreloadFilter(), mock.Anything).Return(nonIgnored, len(nonIgnored), nil).Maybe()
	m.On("Query", mock.Anything, studioSingleLetterFilter(), mock.Anything).Return(nil, 0, nil).Maybe()
	for _, s := range studios {
		m.On("GetAliases", mock.Anything, s.ID).Return(aliases[s.ID], nil).Maybe()
	}
}

func primeTagMock(m *mocks.TagReaderWriter, tags []*models.Tag, aliases map[int][]string) {
	var nonIgnored []*models.Tag
	for _, t := range tags {
		if !t.IgnoreAutoTag {
			nonIgnored = append(nonIgnored, t)
		}
	}
	m.On("QueryForAutoTag", mock.Anything, mock.Anything).Return(nonIgnored, nil).Maybe()
	m.On("Query", mock.Anything, tagPreloadFilter(), mock.Anything).Return(nonIgnored, len(nonIgnored), nil).Maybe()
	m.On("Query", mock.Anything, tagSingleLetterFilter(), mock.Anything).Return(nil, 0, nil).Maybe()
	for _, t := range tags {
		m.On("GetAliases", mock.Anything, t.ID).Return(aliases[t.ID], nil).Maybe()
	}
}

// --- helpers ---

func perfIDs(ps []*models.Performer) []int {
	ids := make([]int, 0, len(ps))
	for _, p := range ps {
		ids = append(ids, p.ID)
	}
	slices.Sort(ids)
	return ids
}

func tagIDs(ts []*models.Tag) []int {
	ids := make([]int, 0, len(ts))
	for _, t := range ts {
		ids = append(ids, t.ID)
	}
	slices.Sort(ids)
	return ids
}

// --- tests ---

func TestPathToPerformers_Semantics(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	alice := &models.Performer{ID: 1, Name: "alice smith"}
	bob := &models.Performer{ID: 2, Name: "bob jones"}
	unicodeP := &models.Performer{ID: 3, Name: "伏字"}
	ignored := &models.Performer{ID: 4, Name: "ignored person", IgnoreAutoTag: true}
	substr := &models.Performer{ID: 5, Name: "ali"} // substring of "alice" - should NOT match "alice smith.jpg"

	performers := []*models.Performer{alice, bob, unicodeP, ignored, substr}
	db := mocks.NewDatabase()
	primePerformerMock(db.Performer, performers)

	tests := []struct {
		name    string
		path    string
		wantIDs []int
	}{
		{"plain name match", "/media/alice smith.jpg", []int{1}},
		{"separator variants", "/media/alice.smith.jpg", []int{1}},
		{"separator variants 2", "/media/alice_smith.jpg", []int{1}},
		{"multiple matches", "/media/alice smith and bob jones.jpg", []int{1, 2}},
		{"case insensitive", "/media/ALICE SMITH.jpg", []int{1}},
		{"unicode", "/media/伏字.jpg", []int{3}},
		{"ignore_auto_tag skipped", "/media/ignored person.jpg", nil},
		{"no substring match", "/media/alicent.jpg", nil},
		{"short name does NOT match inside longer", "/media/alice smith.jpg", []int{1}}, // 'ali' should not match
		{"short name matches exact", "/media/ali.jpg", []int{5}},
		{"no match", "/media/nobody here.jpg", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/no-preload", func(t *testing.T) {
			got, err := PathToPerformers(ctx, tt.path, db.Performer, nil, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotIDs := perfIDs(got); !slices.Equal(gotIDs, tt.wantIDs) {
				t.Errorf("got %v, want %v", gotIDs, tt.wantIDs)
			}
		})
		t.Run(tt.name+"/preloaded", func(t *testing.T) {
			cache := &Cache{}
			if err := cache.PreloadPerformers(ctx, db.Performer); err != nil {
				t.Fatalf("preload: %v", err)
			}
			got, err := PathToPerformers(ctx, tt.path, db.Performer, cache, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotIDs := perfIDs(got); !slices.Equal(gotIDs, tt.wantIDs) {
				t.Errorf("got %v, want %v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func TestPathToStudio_Semantics(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s1 := &models.Studio{ID: 1, Name: "first studio"}
	s2 := &models.Studio{ID: 2, Name: "second"}
	s3 := &models.Studio{ID: 3, Name: "third", IgnoreAutoTag: true}

	studios := []*models.Studio{s1, s2, s3}
	aliases := map[int][]string{2: {"second alias"}}
	db := mocks.NewDatabase()
	primeStudioMock(db.Studio, studios, aliases)

	tests := []struct {
		name   string
		path   string
		wantID int // 0 == no match
	}{
		{"primary name", "/first studio/scene.mp4", 1},
		{"alias matches", "/second alias/scene.mp4", 2},
		{"ignore_auto_tag studio skipped", "/third/scene.mp4", 0},
		{"multiple matches - rightmost wins", "/first studio/second/scene.mp4", 2},
		{"no match", "/unknown/scene.mp4", 0},
	}

	runCase := func(t *testing.T, path string, wantID int, cache *Cache) {
		got, err := PathToStudio(ctx, path, db.Studio, cache, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var gotID int
		if got != nil {
			gotID = got.ID
		}
		if gotID != wantID {
			t.Errorf("got %d, want %d", gotID, wantID)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name+"/no-preload", func(t *testing.T) {
			runCase(t, tt.path, tt.wantID, nil)
		})
		t.Run(tt.name+"/preloaded", func(t *testing.T) {
			cache := &Cache{}
			if err := cache.PreloadStudios(ctx, db.Studio); err != nil {
				t.Fatalf("preload: %v", err)
			}
			runCase(t, tt.path, tt.wantID, cache)
		})
	}
}

func TestPathToTags_Semantics(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t1 := &models.Tag{ID: 1, Name: "anime"}
	t2 := &models.Tag{ID: 2, Name: "docs"}
	t3 := &models.Tag{ID: 3, Name: "skip me", IgnoreAutoTag: true}

	tags := []*models.Tag{t1, t2, t3}
	aliases := map[int][]string{2: {"documentary"}}
	db := mocks.NewDatabase()
	primeTagMock(db.Tag, tags, aliases)

	tests := []struct {
		name    string
		path    string
		wantIDs []int
	}{
		{"name match", "/media/anime/x.mp4", []int{1}},
		{"alias match", "/media/documentary/x.mp4", []int{2}},
		{"multiple matches", "/media/anime-documentary/x.mp4", []int{1, 2}},
		{"ignore_auto_tag skipped", "/media/skip me/x.mp4", nil},
		{"no match", "/media/comedy/x.mp4", nil},
	}

	runCase := func(t *testing.T, path string, wantIDs []int, cache *Cache) {
		got, err := PathToTags(ctx, path, db.Tag, cache, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotIDs := tagIDs(got); !slices.Equal(gotIDs, wantIDs) {
			t.Errorf("got %v, want %v", gotIDs, wantIDs)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name+"/no-preload", func(t *testing.T) {
			runCase(t, tt.path, tt.wantIDs, nil)
		})
		t.Run(tt.name+"/preloaded", func(t *testing.T) {
			cache := &Cache{}
			if err := cache.PreloadTags(ctx, db.Tag); err != nil {
				t.Fatalf("preload: %v", err)
			}
			runCase(t, tt.path, tt.wantIDs, cache)
		})
	}
}

// Performer whose name starts with a single-letter word (e.g., "X Man")
// can't be reached via 2-rune prefix lookup (getPathWords drops 1-char
// words). The preload must put them in the alwaysCheck list so they're
// still regex-tested.
func TestPathToPerformers_SingleLetterFirstWord(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	xman := &models.Performer{ID: 1, Name: "X Man"}
	other := &models.Performer{ID: 2, Name: "alice smith"}

	db := mocks.NewDatabase()
	primePerformerMock(db.Performer, []*models.Performer{xman, other})

	cache := &Cache{}
	if err := cache.PreloadPerformers(ctx, db.Performer); err != nil {
		t.Fatal(err)
	}

	got, err := PathToPerformers(ctx, "/media/X Man.mp4", db.Performer, cache, false)
	if err != nil {
		t.Fatal(err)
	}
	if ids := perfIDs(got); !slices.Equal(ids, []int{1}) {
		t.Errorf("expected [1], got %v", ids)
	}
}

// A studio whose name shares no prefix with its aliases must be reachable
// by alias prefix. "Acme Corp" with alias "Widgets Inc" must match a path
// containing "widgets inc".
func TestPathToStudio_AliasPrefixDistinctFromName(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := &models.Studio{ID: 1, Name: "Acme Corp"}

	db := mocks.NewDatabase()
	primeStudioMock(db.Studio, []*models.Studio{s}, map[int][]string{1: {"Widgets Inc"}})

	cache := &Cache{}
	if err := cache.PreloadStudios(ctx, db.Studio); err != nil {
		t.Fatal(err)
	}

	got, err := PathToStudio(ctx, "/media/Widgets Inc/scene.mp4", db.Studio, cache, false)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != 1 {
		t.Errorf("expected studio 1, got %v", got)
	}
}

// Same for tags.
func TestPathToTags_AliasPrefixDistinctFromName(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db := mocks.NewDatabase()
	primeTagMock(db.Tag, []*models.Tag{{ID: 1, Name: "documentary"}}, map[int][]string{1: {"film"}})

	cache := &Cache{}
	if err := cache.PreloadTags(ctx, db.Tag); err != nil {
		t.Fatal(err)
	}

	got, err := PathToTags(ctx, "/media/film/x.mp4", db.Tag, cache, false)
	if err != nil {
		t.Fatal(err)
	}
	if ids := tagIDs(got); !slices.Equal(ids, []int{1}) {
		t.Errorf("expected [1], got %v", ids)
	}
}

// Two aliases on the same studio with different prefixes should each
// reach the studio. Index bucket must dedupe inside the bucket.
func TestPathToStudio_MultipleAliasesDedup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := &models.Studio{ID: 1, Name: "Primary Name"}

	db := mocks.NewDatabase()
	primeStudioMock(db.Studio, []*models.Studio{s}, map[int][]string{1: {"Primary Nickname", "Primary Alt"}})

	cache := &Cache{}
	if err := cache.PreloadStudios(ctx, db.Studio); err != nil {
		t.Fatal(err)
	}
	// Studio "Primary Name" and both aliases all share prefix "pr".
	// The bucket should contain it exactly once.
	if got := len(cache.studioByPrefix["pr"]); got != 1 {
		t.Errorf("bucket 'pr' should have 1 entry, got %d", got)
	}
}

// Equivalence test: the function must return the same result regardless of
// whether a match.Cache is passed in. This is the invariant that any
// caching-based optimization must preserve.
func TestPathToPerformers_CachedVsUncached(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	perfs := []*models.Performer{
		{ID: 1, Name: "alice smith"},
		{ID: 2, Name: "bob jones"},
		{ID: 3, Name: "charlie"},
		{ID: 4, Name: "david wong"},
	}
	db := mocks.NewDatabase()
	primePerformerMock(db.Performer, perfs)

	paths := []string{
		"/media/alice smith.jpg",
		"/media/bob_jones.jpg",
		"/media/alice smith and charlie.jpg",
		"/media/nobody.jpg",
		"/media/alice smith.jpg", // repeat: cached regex should not change outcome
	}

	var noCache, withCache [][]int
	cache := &Cache{}
	for _, p := range paths {
		uc, err := PathToPerformers(ctx, p, db.Performer, nil, false)
		if err != nil {
			t.Fatal(err)
		}
		wc, err := PathToPerformers(ctx, p, db.Performer, cache, false)
		if err != nil {
			t.Fatal(err)
		}
		noCache = append(noCache, perfIDs(uc))
		withCache = append(withCache, perfIDs(wc))
	}

	for i := range paths {
		if !slices.Equal(noCache[i], withCache[i]) {
			t.Errorf("path %q: no-cache %v vs cached %v", paths[i], noCache[i], withCache[i])
		}
	}
}
