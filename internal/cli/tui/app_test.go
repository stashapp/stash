package tui

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/ntcharts/v2/picture"

	"github.com/stashapp/stash/internal/cli/browse"
	"github.com/stashapp/stash/internal/cli/coverfetch"
	"github.com/stashapp/stash/internal/cli/covergen"
	"github.com/stashapp/stash/internal/cli/scanner"
)

type fakeBrowser struct {
	query   browse.Query
	queries []browse.Query
	result  browse.Result
}

func (f *fakeBrowser) Search(_ context.Context, q browse.Query) (browse.Result, error) {
	f.query = q
	f.queries = append(f.queries, q)
	if len(f.result.Items) > 0 || f.result.Total > 0 {
		result := f.result
		if result.Total == 0 {
			result.Total = len(result.Items)
		}
		return result, nil
	}
	return browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 1, Title: "Scene", Duration: 65}},
	}, nil
}

type fakeScanner struct {
	called bool
	result scanner.Result
}

func (f *fakeScanner) Scan(_ context.Context) scanner.Result {
	f.called = true
	return f.result
}

func (f *fakeScanner) ScanWithProgress(_ context.Context, progress func(scanner.Progress)) scanner.Result {
	f.called = true
	if progress != nil {
		progress(scanner.Progress{Directories: 2, FilesSeen: 4, FilesScanned: 3, LastFile: "/tmp/example.mp4"})
	}
	return f.result
}

type fakeCoverFetcher struct {
	sceneID  int
	sceneIDs []int
	errs     map[int]error
}

func (f *fakeCoverFetcher) Fetch(_ context.Context, sceneID int) (coverfetch.Result, error) {
	f.sceneID = sceneID
	f.sceneIDs = append(f.sceneIDs, sceneID)
	if f.errs != nil && f.errs[sceneID] != nil {
		return coverfetch.Result{}, f.errs[sceneID]
	}
	return coverfetch.Result{SceneID: sceneID, Bytes: 3, RemoteSiteID: "remote-scene"}, nil
}

type fakeCoverGenerator struct {
	sceneIDs []int
	errs     map[int]error
}

func (f *fakeCoverGenerator) Generate(_ context.Context, req covergen.Request) (covergen.Result, error) {
	f.sceneIDs = append(f.sceneIDs, req.SceneID)
	if f.errs != nil && f.errs[req.SceneID] != nil {
		return covergen.Result{}, f.errs[req.SceneID]
	}
	return covergen.Result{SceneID: req.SceneID, Bytes: 4}, nil
}

type fakePlayer struct {
	items []browse.SceneItem
	err   error
}

func (f *fakePlayer) Play(_ context.Context, item browse.SceneItem) error {
	f.items = append(f.items, item)
	return f.err
}

func TestSearchCommandUpdatesQuery(t *testing.T) {
	browser := &fakeBrowser{}
	model := New(context.Background(), browser, ViewList)
	model.input = "/search tag:demo alice"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected refresh command")
	}
	msg := cmd()
	next, _ := updated.Update(msg)
	m := next.(Model)

	if m.result.Total != 1 {
		t.Fatalf("Total = %d", m.result.Total)
	}
	if browser.query.Tag != "demo" || browser.query.Text != "alice" {
		t.Fatalf("query = %#v", browser.query)
	}
}

func TestQuitCommand(t *testing.T) {
	model := New(context.Background(), &fakeBrowser{}, ViewList)
	model.input = "/quit"

	_, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected quit command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("cmd returned %T, want tea.QuitMsg", cmd())
	}
}

func TestScanCommandRunsScannerAndRefreshesResults(t *testing.T) {
	browser := &fakeBrowser{}
	scanner := &fakeScanner{
		result: scanner.Result{
			Directories:  2,
			FilesScanned: 3,
			Errors:       []error{context.Canceled},
		},
	}
	model := NewWithDeps(context.Background(), Deps{Browser: browser, Scanner: scanner}, ViewList)
	model.input = "/scan"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected scan command")
	}
	msg := runScanCommand(t, cmd)
	if !scanner.called {
		t.Fatal("scanner was not called")
	}
	next, refreshCmd := updated.Update(msg)
	m := next.(Model)

	if m.status != "Scan complete: 3 files, 2 directories, 1 errors" {
		t.Fatalf("status = %q", m.status)
	}
	if refreshCmd == nil {
		t.Fatal("expected refresh command")
	}
	refreshMsg := refreshCmd()
	next, _ = m.Update(refreshMsg)
	m = next.(Model)
	if m.result.Total != 1 {
		t.Fatalf("Total = %d", m.result.Total)
	}
}

func TestCoverFetchCommandFetchesSelectedSceneCover(t *testing.T) {
	fetcher := &fakeCoverFetcher{}
	model := NewWithDeps(context.Background(), Deps{
		Browser:      &fakeBrowser{},
		CoverFetcher: fetcher,
	}, ViewList)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 42, Title: "Scene"}},
	}
	model.input = "/cover fetch"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected cover fetch command")
	}
	msg := cmd()
	next, _ := updated.Update(msg)
	m := next.(Model)

	if fetcher.sceneID != 42 {
		t.Fatalf("sceneID = %d, want 42", fetcher.sceneID)
	}
	if m.status != "Official cover fetched: remote-scene" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestPlayCommandPlaysSelectedScene(t *testing.T) {
	player := &fakePlayer{}
	model := NewWithDeps(context.Background(), Deps{
		Browser: &fakeBrowser{},
		Player:  player,
	}, ViewList)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 42, Title: "Scene", Path: "/tmp/scene.mp4"}},
	}
	model.input = "/play"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected play command")
	}
	msg := cmd()
	next, _ := updated.Update(msg)
	m := next.(Model)

	if len(player.items) != 1 || player.items[0].ID != 42 || player.items[0].Path != "/tmp/scene.mp4" {
		t.Fatalf("played items = %#v, want selected scene", player.items)
	}
	if m.status != "Played: Scene" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestPlayCommandReportsPlayerErrors(t *testing.T) {
	player := &fakePlayer{err: errors.New("ffplay failed")}
	model := NewWithDeps(context.Background(), Deps{
		Browser: &fakeBrowser{},
		Player:  player,
	}, ViewList)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 42, Title: "Scene", Path: "/tmp/scene.mp4"}},
	}
	model.input = "/play"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected play command")
	}
	msg := cmd()
	next, _ := updated.Update(msg)
	m := next.(Model)

	if m.status != "ffplay failed" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestCoverFetchAllCommandFetchesCurrentPageCovers(t *testing.T) {
	fetcher := &fakeCoverFetcher{}
	model := NewWithDeps(context.Background(), Deps{
		Browser:      &fakeBrowser{},
		CoverFetcher: fetcher,
	}, ViewList)
	model.result = browse.Result{
		Total: 2,
		Items: []browse.SceneItem{
			{ID: 42, Title: "Scene 1"},
			{ID: 43, Title: "Scene 2"},
		},
	}
	model.input = "/cover fetch-all"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected cover fetch-all command")
	}
	msg := runCoverFetchAllCommand(t, cmd)
	next, _ := updated.Update(msg)
	m := next.(Model)

	if len(fetcher.sceneIDs) != 2 || fetcher.sceneIDs[0] != 42 || fetcher.sceneIDs[1] != 43 {
		t.Fatalf("sceneIDs = %#v, want [42 43]", fetcher.sceneIDs)
	}
	if m.status != "Official covers fetched: 2 ok, 0 generated, 0 skipped, 0 failed" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestCoverFetchAllCommandFetchesAllMatchingCovers(t *testing.T) {
	fetcher := &fakeCoverFetcher{}
	browser := &fakeBrowser{
		result: browse.Result{
			Total: 45,
			Items: makeSceneItems(45),
		},
	}
	model := NewWithDeps(context.Background(), Deps{
		Browser:      browser,
		CoverFetcher: fetcher,
	}, ViewList)
	model.query = browse.Query{Text: "demo", Page: 1, PerPage: 40}
	model.result = browse.Result{
		Total: 45,
		Items: makeSceneItems(40),
	}
	model.input = "/cover fetch-all"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected cover fetch-all command")
	}
	msg := runCoverFetchAllCommand(t, cmd)
	next, _ := updated.Update(msg)
	m := next.(Model)

	if len(fetcher.sceneIDs) != 45 {
		t.Fatalf("fetched %d scene IDs, want 45", len(fetcher.sceneIDs))
	}
	if len(browser.queries) == 0 || browser.queries[0].PerPage != 45 {
		t.Fatalf("queries = %#v, want first fetch with PerPage 45", browser.queries)
	}
	if m.status != "Official covers fetched: 45 ok, 0 generated, 0 skipped, 0 failed" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestCoverFetchAllCommandContinuesAfterSkipsAndFailures(t *testing.T) {
	fetcher := &fakeCoverFetcher{
		errs: map[int]error{
			43: coverfetch.ErrNoMatch,
			44: errors.New("network unavailable"),
		},
	}
	model := NewWithDeps(context.Background(), Deps{
		Browser:      &fakeBrowser{},
		CoverFetcher: fetcher,
	}, ViewList)
	model.result = browse.Result{
		Total: 3,
		Items: []browse.SceneItem{
			{ID: 42, Title: "Scene 1"},
			{ID: 43, Title: "Scene 2"},
			{ID: 44, Title: "Scene 3"},
		},
	}
	model.input = "/cover fetch-all"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected cover fetch-all command")
	}
	msg := runCoverFetchAllCommand(t, cmd)
	next, _ := updated.Update(msg)
	m := next.(Model)

	if len(fetcher.sceneIDs) != 3 {
		t.Fatalf("sceneIDs = %#v, want 3 attempts", fetcher.sceneIDs)
	}
	if m.status != "Official covers fetched: 1 ok, 0 generated, 1 skipped, 1 failed" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestCoverFetchAllGeneratesCoverWhenOfficialFetchFails(t *testing.T) {
	fetcher := &fakeCoverFetcher{
		errs: map[int]error{
			42: coverfetch.ErrNoMatch,
		},
	}
	generator := &fakeCoverGenerator{}
	model := NewWithDeps(context.Background(), Deps{
		Browser:        &fakeBrowser{},
		CoverFetcher:   fetcher,
		CoverGenerator: generator,
	}, ViewList)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 42, Title: "Scene", Path: "/tmp/scene.mp4", Duration: 100}},
	}
	model.input = "/cover fetch-all"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected cover fetch-all command")
	}
	msg := runCoverFetchAllCommand(t, cmd)
	next, _ := updated.Update(msg)
	m := next.(Model)

	if len(generator.sceneIDs) != 1 || generator.sceneIDs[0] != 42 {
		t.Fatalf("generated scene IDs = %#v, want [42]", generator.sceneIDs)
	}
	if m.status != "Official covers fetched: 0 ok, 1 generated, 0 skipped, 0 failed" {
		t.Fatalf("status = %q", m.status)
	}
}

func makeSceneItems(n int) []browse.SceneItem {
	items := make([]browse.SceneItem, n)
	for i := range items {
		items[i] = browse.SceneItem{ID: 100 + i, Title: "Scene"}
	}
	return items
}

func TestCoverFetchAllProgressStatusShowsCounts(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewList)

	next, _ := model.Update(coverFetchProgressMsg{progress: coverFetchProgress{
		Total:     3,
		Attempted: 2,
		OK:        1,
		Skipped:   1,
		LastScene: "Very Long Scene Title",
	}})
	m := next.(Model)

	if m.status != "Fetching official covers... 2/3, ok:1 generated:0 skipped:1 failed:0, last: Very...itle" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestCoverFetchAllCompletionRefreshesResults(t *testing.T) {
	browser := &fakeBrowser{}
	model := NewWithDeps(context.Background(), Deps{Browser: browser}, ViewList)

	next, refreshCmd := model.Update(coverFetchAllMsg{
		result: coverFetchResult{OK: 2, Skipped: 1, Failed: 0},
	})
	m := next.(Model)

	if m.status != "Official covers fetched: 2 ok, 0 generated, 1 skipped, 0 failed" {
		t.Fatalf("status = %q", m.status)
	}
	if refreshCmd == nil {
		t.Fatal("expected refresh command")
	}
	refreshMsg := refreshCmd()
	next, _ = m.Update(refreshMsg)
	m = next.(Model)
	if m.status != "Official covers fetched: 2 ok, 0 generated, 1 skipped, 0 failed" {
		t.Fatalf("status after refresh = %q", m.status)
	}
	if browser.query.PerPage != 40 {
		t.Fatalf("query = %#v", browser.query)
	}
}

func TestUpdateAutoSwitchesToKittyWhenSupported(t *testing.T) {
	picture.ForceKittyCapability(picture.KittyCapabilitySupported)
	t.Cleanup(func() { picture.ForceKittyCapability(picture.KittyCapabilitySupported) })

	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	if model.pic.Mode() != picture.PictureGlyph {
		t.Fatalf("initial mode = %v, want glyph", model.pic.Mode())
	}

	next, _ := model.Update(struct{}{})
	m := next.(Model)

	if m.pic.Mode() != picture.PictureKitty {
		t.Fatalf("mode = %v, want kitty", m.pic.Mode())
	}
}

func TestRenderGridBodyPreservesPreviewContent(t *testing.T) {
	preview := "abc\n123"
	list := "one\ntwo\nthree"

	got := renderGridBody(preview, 5, list)
	want := "abc  one\n123  two\n       three\n"
	if got != want {
		t.Fatalf("renderGridBody = %q, want %q", got, want)
	}
}

func TestVisibleStartScrollsCursorIntoView(t *testing.T) {
	tests := []struct {
		name   string
		cursor int
		total  int
		limit  int
		want   int
	}{
		{name: "top", cursor: 0, total: 20, limit: 5, want: 0},
		{name: "last visible row", cursor: 4, total: 20, limit: 5, want: 0},
		{name: "scrolls down", cursor: 5, total: 20, limit: 5, want: 1},
		{name: "near bottom", cursor: 19, total: 20, limit: 5, want: 15},
		{name: "short list", cursor: 2, total: 3, limit: 5, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := visibleStart(tt.cursor, tt.total, tt.limit)
			if got != tt.want {
				t.Fatalf("visibleStart = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestViewScrollsListWithCursor(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewList)
	model.height = 11
	model.cursor = 5
	model.result = browse.Result{
		Total: 8,
		Items: []browse.SceneItem{
			{ID: 1, Title: "Scene 1"},
			{ID: 2, Title: "Scene 2"},
			{ID: 3, Title: "Scene 3"},
			{ID: 4, Title: "Scene 4"},
			{ID: 5, Title: "Scene 5"},
			{ID: 6, Title: "Scene 6"},
			{ID: 7, Title: "Scene 7"},
			{ID: 8, Title: "Scene 8"},
		},
	}

	view := model.View().Content
	if strings.Contains(view, "Scene 1") {
		t.Fatalf("view should have scrolled past Scene 1: %q", view)
	}
	if !strings.Contains(view, "> Scene 6") {
		t.Fatalf("view does not show selected row: %q", view)
	}
}

func TestDownAtLoadedEndFetchesNextPage(t *testing.T) {
	browser := &fakeBrowser{
		result: browse.Result{
			Total: 45,
			Items: makeSceneItems(5),
		},
	}
	model := NewWithDeps(context.Background(), Deps{Browser: browser}, ViewList)
	model.query = browse.Query{Page: 1, PerPage: 40}
	model.cursor = 39
	model.result = browse.Result{
		Total: 45,
		Items: makeSceneItems(40),
	}

	next, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if cmd == nil {
		t.Fatal("expected load-more command")
	}
	msg := cmd()
	next, _ = next.Update(msg)
	m := next.(Model)

	if len(m.result.Items) != 45 {
		t.Fatalf("loaded items = %d, want 45", len(m.result.Items))
	}
	if m.cursor != 40 {
		t.Fatalf("cursor = %d, want 40", m.cursor)
	}
	if len(browser.queries) == 0 || browser.queries[0].Page != 2 {
		t.Fatalf("queries = %#v, want page 2", browser.queries)
	}
}

func TestNewWithDepsForceKittyStartsInKittyMode(t *testing.T) {
	picture.ForceKittyCapability(picture.KittyCapabilityUnsupported)
	t.Cleanup(func() { picture.ForceKittyCapability(picture.KittyCapabilitySupported) })

	model := NewWithDeps(context.Background(), Deps{
		Browser:    &fakeBrowser{},
		ForceKitty: true,
	}, ViewGrid)

	if model.pic.Mode() != picture.PictureKitty {
		t.Fatalf("mode = %v, want kitty", model.pic.Mode())
	}
	if got := picture.KittySupported(); got != picture.KittyCapabilitySupported {
		t.Fatalf("KittySupported = %v, want supported", got)
	}
}

func runScanCommand(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()

	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		return msg
	}

	for _, batched := range batch {
		msg = batched()
		if _, ok := msg.(scanMsg); ok {
			return msg
		}
	}

	t.Fatal("batch did not include scanMsg")
	return nil
}

func runCoverFetchAllCommand(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()

	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		return msg
	}

	for _, batched := range batch {
		msg = batched()
		if _, ok := msg.(coverFetchAllMsg); ok {
			return msg
		}
	}

	t.Fatal("batch did not include coverFetchAllMsg")
	return nil
}

func TestScanProgressStatusShowsScannedCount(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewList)

	next, _ := model.Update(scanProgressMsg{progress: scanner.Progress{
		Directories:  2,
		FilesSeen:    4,
		FilesScanned: 3,
		LastFile:     "/media/demo/example-scene.mp4",
	}})
	m := next.(Model)

	if m.status != "Scanning... 3 files scanned, 4 videos found, 2 directories, last: exam...cene.mp4" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestScanProgressPollIntervalIsOneSecond(t *testing.T) {
	if scanProgressInterval.String() != "1s" {
		t.Fatalf("scanProgressInterval = %s, want 1s", scanProgressInterval)
	}
}

func TestCompactBasenameKeepsShortNames(t *testing.T) {
	if got := compactBasename("demo.mp4", 10); got != "demo.mp4" {
		t.Fatalf("compactBasename = %q", got)
	}
}

func TestCompactBasenameOmitsMiddleForLongNames(t *testing.T) {
	if got := compactBasename("example-scene.mp4", 10); got != "exam...cene.mp4" {
		t.Fatalf("compactBasename = %q", got)
	}
}
