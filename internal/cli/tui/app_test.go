package tui

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/ntcharts/v2/picture"

	"github.com/stashapp/stash/internal/cli/browse"
	"github.com/stashapp/stash/internal/cli/command"
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

type fakePlayer struct {
	items []browse.SceneItem
	err   error
}

func (f *fakePlayer) Play(_ context.Context, item browse.SceneItem) error {
	f.items = append(f.items, item)
	return f.err
}

type fakeEditor struct {
	sceneRatings     map[int]int
	deletedScenes    []int
	performerRatings map[int]int
	err              error
}

func (f *fakeEditor) SetSceneRating(_ context.Context, sceneID int, rating int) error {
	if f.sceneRatings == nil {
		f.sceneRatings = map[int]int{}
	}
	f.sceneRatings[sceneID] = rating
	return f.err
}

func (f *fakeEditor) DeleteScene(_ context.Context, sceneID int) error {
	f.deletedScenes = append(f.deletedScenes, sceneID)
	return f.err
}

func (f *fakeEditor) SetPerformerRating(_ context.Context, performerID int, rating int) error {
	if f.performerRatings == nil {
		f.performerRatings = map[int]int{}
	}
	f.performerRatings[performerID] = rating
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

func TestRandomCommandUsesExplicitLimit(t *testing.T) {
	browser := &fakeBrowser{}
	model := New(context.Background(), browser, ViewGrid)
	model.input = "/random 7"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected random refresh command")
	}
	msg := cmd()
	next, _ := updated.Update(msg)
	m := next.(Model)

	if m.result.Total != 1 {
		t.Fatalf("Total = %d", m.result.Total)
	}
	if browser.query.Sort != "random" || browser.query.PerPage != 7 || browser.query.Page != 1 {
		t.Fatalf("query = %#v, want random page 1 per_page 7", browser.query)
	}
}

func TestRandomCommandUsesVisibleGridCapacityByDefault(t *testing.T) {
	browser := &fakeBrowser{}
	model := New(context.Background(), browser, ViewGrid)
	model.width = 90
	model.height = 35
	model.input = "/random"

	updated, cmd := model.executeInput()
	if cmd == nil {
		t.Fatal("expected random refresh command")
	}
	msg := cmd()
	next, _ := updated.Update(msg)
	m := next.(Model)

	if m.result.Total != 1 {
		t.Fatalf("Total = %d", m.result.Total)
	}
	if browser.query.Sort != "random" || browser.query.PerPage != 6 || browser.query.Page != 1 {
		t.Fatalf("query = %#v, want random page 1 per_page 6", browser.query)
	}
}

func TestRandomCommandRejectsInvalidLimit(t *testing.T) {
	model := New(context.Background(), &fakeBrowser{}, ViewGrid)
	model.input = "/random nope"

	updated, cmd := model.executeInput()
	m := updated.(Model)

	if cmd != nil {
		t.Fatal("expected no command")
	}
	if m.status != "Usage: /random <n>" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestHelpOmitsRemovedCommands(t *testing.T) {
	help := command.Help()
	for _, removed := range []string{"/cover", "/play", "/open", "/edit"} {
		if strings.Contains(help, removed) {
			t.Fatalf("help contains removed command %q: %q", removed, help)
		}
	}
	if !strings.Contains(help, "/random") {
		t.Fatalf("help does not contain /random: %q", help)
	}
}

func TestDetailsPanelShowsPerformerRatings(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{
			ID:     42,
			Title:  "Scene",
			Rating: ptrInt(60),
			Performers: []browse.PerformerItem{
				{ID: 7, Name: "Alice", Rating: ptrInt(80)},
				{ID: 8, Name: "Bob"},
			},
		}},
	}

	next, _ := model.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m := next.(Model)
	view := m.View().Content

	for _, want := range []string{"Rating: 60", "> Alice rating:80", "  Bob rating:--"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view does not contain %q: %q", want, view)
		}
	}
}

func TestDetailsPanelSetsSceneRating(t *testing.T) {
	editor := &fakeEditor{}
	browser := &fakeBrowser{result: browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 42, Title: "Scene"}},
	}}
	model := NewWithDeps(context.Background(), Deps{Browser: browser, Editor: editor}, ViewGrid)
	model.result = browser.result

	next, _ := model.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	next, _ = next.Update(tea.KeyPressMsg{Text: "r"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "8"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "5"})
	updated, cmd := next.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected scene rating command")
	}
	msg := cmd()
	next, _ = updated.Update(msg)
	m := next.(Model)

	if got := editor.sceneRatings[42]; got != 85 {
		t.Fatalf("scene rating = %d, want 85", got)
	}
	if m.status != "Scene rating updated" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestDetailsPanelDeletesSceneAfterConfirmation(t *testing.T) {
	editor := &fakeEditor{}
	browser := &fakeBrowser{result: browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 43, Title: "Next"}},
	}}
	model := NewWithDeps(context.Background(), Deps{Browser: browser, Editor: editor}, ViewGrid)
	model.result = browse.Result{
		Total: 2,
		Items: []browse.SceneItem{
			{ID: 42, Title: "Delete Me"},
			{ID: 43, Title: "Next"},
		},
	}

	next, _ := model.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	next, _ = next.Update(tea.KeyPressMsg{Text: "d"})
	updated, cmd := next.Update(tea.KeyPressMsg{Text: "y"})
	if cmd == nil {
		t.Fatal("expected delete command")
	}
	msg := cmd()
	next, _ = updated.Update(msg)
	m := next.(Model)

	if len(editor.deletedScenes) != 1 || editor.deletedScenes[0] != 42 {
		t.Fatalf("deleted scenes = %#v, want [42]", editor.deletedScenes)
	}
	if m.status != "Scene deleted" {
		t.Fatalf("status = %q", m.status)
	}
	if m.showDetails {
		t.Fatal("details should close after delete")
	}
}

func TestDetailsPanelSetsPerformerRating(t *testing.T) {
	editor := &fakeEditor{}
	browser := &fakeBrowser{result: browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{
			ID:    42,
			Title: "Scene",
			Performers: []browse.PerformerItem{
				{ID: 7, Name: "Alice"},
				{ID: 8, Name: "Bob"},
			},
		}},
	}}
	model := NewWithDeps(context.Background(), Deps{Browser: browser, Editor: editor}, ViewGrid)
	model.result = browser.result

	next, _ := model.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	next, _ = next.Update(tea.KeyPressMsg{Text: "j"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "R"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "9"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "0"})
	updated, cmd := next.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected performer rating command")
	}
	msg := cmd()
	next, _ = updated.Update(msg)
	m := next.(Model)

	if got := editor.performerRatings[8]; got != 90 {
		t.Fatalf("performer rating = %d, want 90", got)
	}
	if m.status != "Performer rating updated" {
		t.Fatalf("status = %q", m.status)
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

func TestNormalEnterPlaysSelectedScene(t *testing.T) {
	player := &fakePlayer{}
	model := NewWithDeps(context.Background(), Deps{
		Browser: &fakeBrowser{},
		Player:  player,
	}, ViewGrid)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 42, Title: "Scene", Path: "/tmp/scene.mp4"}},
	}

	updated, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected play command")
	}
	msg := cmd()
	next, _ := updated.Update(msg)
	m := next.(Model)

	if len(player.items) != 1 || player.items[0].ID != 42 {
		t.Fatalf("played items = %#v, want selected scene", player.items)
	}
	if m.status != "Played: Scene" {
		t.Fatalf("status = %q", m.status)
	}
}

func TestNormalEnterReportsPlayerErrors(t *testing.T) {
	player := &fakePlayer{err: errors.New("ffplay failed")}
	model := NewWithDeps(context.Background(), Deps{
		Browser: &fakeBrowser{},
		Player:  player,
	}, ViewGrid)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{ID: 42, Title: "Scene", Path: "/tmp/scene.mp4"}},
	}

	updated, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
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

func TestColonQQuits(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)

	next, _ := model.Update(tea.KeyPressMsg{Text: ":"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "q"})
	_, cmd := next.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected quit command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatalf("cmd returned %T, want tea.QuitMsg", cmd())
	}
}

func TestCommandTabCompletesFuzzyMatch(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)

	next, _ := model.Update(tea.KeyPressMsg{Text: ":"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "r"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "d"})
	next, _ = next.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m := next.(Model)

	if m.input != "/random " {
		t.Fatalf("input = %q, want /random ", m.input)
	}
}

func TestCommandTabCyclesMatches(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)

	next, _ := model.Update(tea.KeyPressMsg{Text: ":"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "s"})
	next, _ = next.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m := next.(Model)
	if m.input != "/search " {
		t.Fatalf("first completion = %q, want /search ", m.input)
	}

	next, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = next.(Model)
	if m.input != "/scan " {
		t.Fatalf("second completion = %q, want /scan ", m.input)
	}
}

func TestCommandCompletionMatchesRenderInFooter(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)

	next, _ := model.Update(tea.KeyPressMsg{Text: ":"})
	next, _ = next.Update(tea.KeyPressMsg{Text: "s"})
	m := next.(Model)
	view := m.View().Content

	if !strings.Contains(view, "matches: /search /scan") {
		t.Fatalf("view does not show completion matches: %q", view)
	}
}

func makeSceneItems(n int) []browse.SceneItem {
	items := make([]browse.SceneItem, n)
	for i := range items {
		items[i] = browse.SceneItem{ID: 100 + i, Title: "Scene"}
	}
	return items
}

func ptrInt(v int) *int {
	return &v
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

func TestGridColumnsUsesTerminalWidth(t *testing.T) {
	tests := []struct {
		width int
		want  int
	}{
		{width: 0, want: 1},
		{width: 20, want: 1},
		{width: 58, want: 2},
		{width: 90, want: 3},
	}

	for _, tt := range tests {
		if got := gridColumns(tt.width); got != tt.want {
			t.Fatalf("gridColumns(%d) = %d, want %d", tt.width, got, tt.want)
		}
	}
}

func TestVisibleGridStartKeepsCursorRowInView(t *testing.T) {
	tests := []struct {
		cursor int
		total  int
		cols   int
		rows   int
		want   int
	}{
		{cursor: 0, total: 20, cols: 3, rows: 2, want: 0},
		{cursor: 5, total: 20, cols: 3, rows: 2, want: 0},
		{cursor: 6, total: 20, cols: 3, rows: 2, want: 3},
		{cursor: 19, total: 20, cols: 3, rows: 2, want: 14},
	}

	for _, tt := range tests {
		if got := visibleGridStart(tt.cursor, tt.total, tt.cols, tt.rows); got != tt.want {
			t.Fatalf("visibleGridStart(%d, %d, %d, %d) = %d, want %d", tt.cursor, tt.total, tt.cols, tt.rows, got, tt.want)
		}
	}
}

func TestCompactTextOmitsMiddle(t *testing.T) {
	got := compactText("abcdefghijklmnopqrstuvwxyz.mp4", 12)
	if got != "abcd...z.mp4" {
		t.Fatalf("compactText = %q", got)
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

func TestVimKeysMoveGridCursor(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	model.width = 90
	model.cursor = 4
	model.result = browse.Result{Total: 9, Items: makeSceneItems(9)}

	next, _ := model.Update(tea.KeyPressMsg{Text: "h"})
	m := next.(Model)
	if m.cursor != 3 {
		t.Fatalf("after h cursor = %d, want 3", m.cursor)
	}

	next, _ = m.Update(tea.KeyPressMsg{Text: "l"})
	m = next.(Model)
	if m.cursor != 4 {
		t.Fatalf("after l cursor = %d, want 4", m.cursor)
	}

	next, _ = m.Update(tea.KeyPressMsg{Text: "j"})
	m = next.(Model)
	if m.cursor != 7 {
		t.Fatalf("after j cursor = %d, want 7", m.cursor)
	}

	next, _ = m.Update(tea.KeyPressMsg{Text: "k"})
	m = next.(Model)
	if m.cursor != 4 {
		t.Fatalf("after k cursor = %d, want 4", m.cursor)
	}
}

func TestSpaceTogglesSelectedSceneDetails(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{
			ID:         42,
			Title:      "Scene",
			Path:       "/tmp/scene.mp4",
			Duration:   65,
			Date:       "2026-06-20",
			Studio:     "Studio",
			Performers: []browse.PerformerItem{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}},
			Tags:       []string{"demo"},
		}},
	}

	next, _ := model.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m := next.(Model)
	if !m.showDetails {
		t.Fatal("expected details to be shown")
	}
	view := m.View().Content
	for _, want := range []string{"Path: /tmp/scene.mp4", "> Alice rating:--", "  Bob rating:--", "Tags: demo", "Studio: Studio"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view does not contain %q: %q", want, view)
		}
	}

	next, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m = next.(Model)
	if m.showDetails {
		t.Fatal("expected details to be hidden")
	}
}

func TestGridTileDisplaysPerformerSummary(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	model.width = 90
	model.result = browse.Result{
		Total: 1,
		Items: []browse.SceneItem{{
			ID:         42,
			Title:      "Scene",
			Performers: []browse.PerformerItem{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}, {ID: 3, Name: "Cara"}},
		}},
	}

	view := model.View().Content
	if !strings.Contains(view, "Alice <2 omitted>") {
		t.Fatalf("view does not show performer summary: %q", view)
	}
}

func TestViewScrollsGridWithCursor(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	model.width = 30
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
	if !strings.Contains(view, "Scene 6") {
		t.Fatalf("view does not show selected tile: %q", view)
	}
}

func TestGridMovesWithinCurrentWindowBeforePaging(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	model.width = 90
	model.height = 35
	model.cursor = 0
	model.result = browse.Result{
		Total: 12,
		Items: []browse.SceneItem{
			{ID: 1, Title: "Scene 1"},
			{ID: 2, Title: "Scene 2"},
			{ID: 3, Title: "Scene 3"},
			{ID: 4, Title: "Scene 4"},
			{ID: 5, Title: "Scene 5"},
			{ID: 6, Title: "Scene 6"},
			{ID: 7, Title: "Scene 7"},
			{ID: 8, Title: "Scene 8"},
			{ID: 9, Title: "Scene 9"},
			{ID: 10, Title: "Scene 10"},
			{ID: 11, Title: "Scene 11"},
			{ID: 12, Title: "Scene 12"},
		},
	}

	next, _ := model.Update(tea.KeyPressMsg{Text: "j"})
	m := next.(Model)
	view := m.View().Content

	if m.cursor != 3 {
		t.Fatalf("cursor = %d, want 3", m.cursor)
	}
	if !strings.Contains(view, "Scene 1") || !strings.Contains(view, "Scene 6") {
		t.Fatalf("view should keep the first two rows visible: %q", view)
	}
	if strings.Contains(view, "Scene 7") {
		t.Fatalf("view should not scroll in a new row yet: %q", view)
	}
}

func TestGridPagesWhenMovingPastWindowEdge(t *testing.T) {
	model := NewWithDeps(context.Background(), Deps{Browser: &fakeBrowser{}}, ViewGrid)
	model.width = 90
	model.height = 35
	model.cursor = 3
	model.result = browse.Result{
		Total: 12,
		Items: []browse.SceneItem{
			{ID: 1, Title: "Scene 1"},
			{ID: 2, Title: "Scene 2"},
			{ID: 3, Title: "Scene 3"},
			{ID: 4, Title: "Scene 4"},
			{ID: 5, Title: "Scene 5"},
			{ID: 6, Title: "Scene 6"},
			{ID: 7, Title: "Scene 7"},
			{ID: 8, Title: "Scene 8"},
			{ID: 9, Title: "Scene 9"},
			{ID: 10, Title: "Scene 10"},
			{ID: 11, Title: "Scene 11"},
			{ID: 12, Title: "Scene 12"},
		},
	}

	next, _ := model.Update(tea.KeyPressMsg{Text: "j"})
	m := next.(Model)
	view := m.View().Content

	if m.cursor != 6 {
		t.Fatalf("cursor = %d, want 6", m.cursor)
	}
	if strings.Contains(view, "Scene 4") || strings.Contains(view, "Scene 5") || strings.Contains(view, "Scene 6") {
		t.Fatalf("view should switch to the next window instead of appending one row: %q", view)
	}
	if !strings.Contains(view, "Scene 7") || !strings.Contains(view, "Scene 12") {
		t.Fatalf("view should show the next window: %q", view)
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
