package tui

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/NimbleMarkets/ntcharts/v2/picture"
	_ "golang.org/x/image/webp"

	"github.com/stashapp/stash/internal/cli/browse"
	"github.com/stashapp/stash/internal/cli/command"
	"github.com/stashapp/stash/internal/cli/cover"
	"github.com/stashapp/stash/internal/cli/coverfetch"
	"github.com/stashapp/stash/internal/cli/covergen"
	"github.com/stashapp/stash/internal/cli/edit"
	"github.com/stashapp/stash/internal/cli/scanner"
	"github.com/stashapp/stash/pkg/logger"
)

type Browser interface {
	Search(context.Context, browse.Query) (browse.Result, error)
}

type Editor interface {
	Apply(context.Context, int, edit.Update) error
}

type CoverLoader interface {
	Load(context.Context, cover.Request) (cover.Cover, error)
}

type CoverFetcher interface {
	Fetch(context.Context, int) (coverfetch.Result, error)
}

type CoverGenerator interface {
	Generate(context.Context, covergen.Request) (covergen.Result, error)
}

type Scanner interface {
	ScanWithProgress(context.Context, func(scanner.Progress)) scanner.Result
}

type ViewMode string

const (
	ViewGrid ViewMode = "grid"
	ViewList ViewMode = "list"
)

var scanProgressInterval = time.Second

type Model struct {
	ctx       context.Context
	browser   Browser
	editor    Editor
	covers    CoverLoader
	fetcher   CoverFetcher
	generator CoverGenerator
	scanner   Scanner
	pic       picture.Model

	mode       ViewMode
	query      browse.Query
	result     browse.Result
	cursor     int
	input      string
	status     string
	scan       *scanState
	coverFetch *coverFetchState
	width      int
	height     int
}

type Deps struct {
	Browser        Browser
	Editor         Editor
	Covers         CoverLoader
	CoverFetcher   CoverFetcher
	CoverGenerator CoverGenerator
	Scanner        Scanner
	ForceKitty     bool
}

func New(ctx context.Context, browser Browser, mode ViewMode, editors ...Editor) Model {
	var editor Editor
	if len(editors) > 0 {
		editor = editors[0]
	}

	return NewWithDeps(ctx, Deps{Browser: browser, Editor: editor}, mode)
}

func NewWithDeps(ctx context.Context, deps Deps, mode ViewMode) Model {
	if mode == "" {
		mode = ViewGrid
	}
	if deps.ForceKitty {
		picture.ForceKittyCapability(picture.KittyCapabilitySupported)
	}

	pic := picture.NewWithConfig(picture.Config{
		KittyID: 7001,
		Fit:     picture.FitContain,
	})
	if mode == ViewGrid && deps.ForceKitty {
		_ = pic.Toggle()
	}
	logger.Infof("[stash-cli] tui graphics initialized: view_mode=%s force_kitty=%t render_mode=%s kitty=%s", mode, deps.ForceKitty, pictureModeString(pic.Mode()), kittyCapabilityString(picture.KittySupported()))

	return Model{
		ctx:       ctx,
		browser:   deps.Browser,
		editor:    deps.Editor,
		covers:    deps.Covers,
		fetcher:   deps.CoverFetcher,
		generator: deps.CoverGenerator,
		scanner:   deps.Scanner,
		pic:       pic,
		mode:      mode,
		status:    "Type /help for commands",
		query: browse.Query{
			Page:    1,
			PerPage: 40,
		},
	}
}

func Run(ctx context.Context, deps Deps, mode ViewMode) error {
	_, err := tea.NewProgram(NewWithDeps(ctx, deps, mode)).Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.pic.Init(), m.refresh())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if cmd := m.resizePicture(); cmd != nil {
			return m, cmd
		}
	case resultMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.result = msg.result
		m.cursor = 0
		if msg.status != "" {
			m.status = msg.status
		} else {
			m.status = fmt.Sprintf("%d results", msg.result.Total)
		}
		return m, m.loadSelectedCover()
	case appendResultMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.result.Items = append(m.result.Items, msg.result.Items...)
		if msg.result.Total > 0 {
			m.result.Total = msg.result.Total
		}
		if msg.cursor >= 0 && msg.cursor < len(m.result.Items) {
			m.cursor = msg.cursor
		}
		m.status = fmt.Sprintf("%d results", m.result.Total)
		return m, m.loadSelectedCover()
	case statusMsg:
		m.status = msg.status
	case scanProgressMsg:
		if msg.scan != nil && msg.scan != m.scan {
			return m, nil
		}
		m.status = formatScanProgressStatus(msg.progress)
		if msg.scan != nil && !msg.scan.done() {
			return m, pollScanProgress(msg.scan)
		}
	case scanMsg:
		if msg.scan != nil && msg.scan != m.scan {
			return m, nil
		}
		m.scan = nil
		m.status = formatScanStatus(msg.result)
		return m, m.refresh()
	case coverFetchProgressMsg:
		if msg.fetch != nil && msg.fetch != m.coverFetch {
			return m, nil
		}
		m.status = formatCoverFetchProgressStatus(msg.progress)
		if msg.fetch != nil && !msg.fetch.done() {
			return m, pollCoverFetchProgress(msg.fetch)
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
				return m, m.loadSelectedCover()
			}
		case "down":
			if m.cursor < len(m.result.Items)-1 {
				m.cursor++
				return m, m.loadSelectedCover()
			}
			if len(m.result.Items) < m.result.Total {
				return m, m.loadMoreResults()
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "enter":
			return m.executeInput()
		default:
			if text := msg.Key().Text; text != "" {
				m.input += text
			}
		}
	case coverMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		cmd := m.pic.SetImage(msg.image)
		return m, m.ensureKitty(cmd)
	case coverFetchedMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		if msg.result.RemoteSiteID != "" {
			m.status = "Official cover fetched: " + msg.result.RemoteSiteID
		} else {
			m.status = "Official cover fetched"
		}
		return m, m.loadSelectedCover()
	case coverFetchAllMsg:
		if msg.fetch != nil && msg.fetch != m.coverFetch {
			return m, nil
		}
		m.coverFetch = nil
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.status = formatCoverFetchResultStatus(msg.result)
		return m, m.refreshWithStatus(m.status)
	}

	return m, m.ensureKitty(m.pic.Update(msg))
}

func (m Model) View() tea.View {
	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Render("stash-cli")
	fmt.Fprintf(&b, "%s  mode:%s\n", title, m.mode)
	fmt.Fprintf(&b, "%s\n\n", m.status)

	limit := m.height - 6
	if limit <= 0 || limit > len(m.result.Items) {
		limit = len(m.result.Items)
	}
	start := visibleStart(m.cursor, len(m.result.Items), limit)
	var list strings.Builder
	for i := start; i < start+limit && i < len(m.result.Items); i++ {
		item := m.result.Items[i]
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		fmt.Fprintf(&list, "%s%s  %s  %s\n", prefix, item.Title, formatDuration(item.Duration), item.Date)
	}

	if len(m.result.Items) == 0 {
		list.WriteString("No scenes to display\n")
	}

	if m.mode == ViewGrid && m.covers != nil {
		previewWidth := m.width / 3
		if previewWidth < 20 {
			previewWidth = 20
		}
		b.WriteString(renderGridBody(m.pic.View().Content, previewWidth, list.String()))
	} else {
		b.WriteString(list.String())
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Faint(true).Render(command.Help()))
	b.WriteString("\n:")
	b.WriteString(m.input)

	view := tea.NewView(b.String())
	view.AltScreen = true
	return view
}

func visibleStart(cursor, total, limit int) int {
	if total <= 0 || limit <= 0 || total <= limit || cursor < limit {
		return 0
	}
	start := cursor - limit + 1
	maxStart := total - limit
	if start > maxStart {
		return maxStart
	}
	return start
}

func renderGridBody(preview string, previewWidth int, list string) string {
	previewLines := strings.Split(preview, "\n")
	listLines := strings.Split(strings.TrimRight(list, "\n"), "\n")
	lines := max(len(previewLines), len(listLines))
	if lines == 0 {
		return ""
	}

	var b strings.Builder
	blankPreview := strings.Repeat(" ", previewWidth)
	for i := 0; i < lines; i++ {
		if i < len(previewLines) && previewLines[i] != "" {
			b.WriteString(previewLines[i])
		} else {
			b.WriteString(blankPreview)
		}
		b.WriteString("  ")
		if i < len(listLines) {
			b.WriteString(listLines[i])
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (m Model) executeInput() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.input)
	m.input = ""
	if input == "" {
		return m, nil
	}

	cmd, err := command.Parse(input)
	if err != nil {
		m.query = browse.ParseQuery(input)
		return m, m.refresh()
	}

	switch cmd.Name {
	case "search":
		m.query = browse.ParseQuery(strings.Join(cmd.Args, " "))
		return m, m.refresh()
	case "clear":
		m.query = browse.Query{Page: 1, PerPage: 40}
		return m, m.refresh()
	case "view":
		if len(cmd.Args) != 1 || (cmd.Args[0] != string(ViewGrid) && cmd.Args[0] != string(ViewList)) {
			m.status = "Usage: /view grid or /view list"
			return m, nil
		}
		m.mode = ViewMode(cmd.Args[0])
		m.status = "View mode changed"
	case "scan":
		return m.executeScan()
	case "cover":
		return m.executeCover(cmd.Args)
	case "open":
		m.status = "/open is reserved; terminal playback is out of scope for the first version"
	case "edit":
		return m.executeEdit(cmd.Args)
	case "help":
		m.status = command.Help()
	case "quit", "q":
		return m, tea.Quit
	default:
		m.status = "Unknown command: " + cmd.Name
	}

	return m, nil
}

func (m Model) executeCover(args []string) (tea.Model, tea.Cmd) {
	if len(args) != 1 || (args[0] != "fetch" && args[0] != "fetch-all") {
		m.status = "Usage: /cover fetch or /cover fetch-all"
		return m, nil
	}
	if m.fetcher == nil {
		m.status = "/cover is unavailable: configure [stash_box].endpoint"
		return m, nil
	}
	if args[0] == "fetch-all" {
		return m.executeCoverFetchAll()
	}

	item, ok := m.selectedItem()
	if !ok {
		m.status = "No scene selected"
		return m, nil
	}

	return m, func() tea.Msg {
		result, err := m.fetcher.Fetch(m.ctx, item.ID)
		return coverFetchedMsg{result: result, err: err}
	}
}

func (m Model) executeCoverFetchAll() (tea.Model, tea.Cmd) {
	if len(m.result.Items) == 0 {
		m.status = "No scenes to fetch"
		return m, nil
	}
	if m.coverFetch != nil {
		m.status = formatCoverFetchProgressStatus(m.coverFetch.progress())
		return m, nil
	}

	items := append([]browse.SceneItem(nil), m.result.Items...)
	total := m.result.Total
	if total < len(items) {
		total = len(items)
	}
	fetch := newCoverFetchState(total)
	query := m.query
	m.coverFetch = fetch
	m.status = formatCoverFetchProgressStatus(fetch.progress())
	logger.Infof("[stash-cli] starting official cover fetch for %d scenes", total)
	return m, tea.Batch(func() tea.Msg {
		var err error
		items, err = m.coverFetchItems(m.ctx, query, items, total)
		if err != nil {
			logger.Infof("[stash-cli] official cover fetch failed before fetching scenes: %v", err)
			return coverFetchAllMsg{fetch: fetch, err: err}
		}
		for _, item := range items {
			title := item.Title
			logger.Infof("[stash-cli] fetching official cover: scene_id=%d title=%q", item.ID, title)
			result, err := m.fetcher.Fetch(m.ctx, item.ID)
			if err == nil {
				logger.Infof("[stash-cli] official cover fetched: scene_id=%d title=%q remote_site_id=%q bytes=%d", item.ID, title, result.RemoteSiteID, result.Bytes)
				fetch.incrementOK(title)
				continue
			}
			if generated, genErr := m.generateCoverFallback(item); genErr == nil {
				logger.Infof("[stash-cli] ffmpeg cover generated: scene_id=%d title=%q bytes=%d official_error=%v", item.ID, title, generated.Bytes, err)
				fetch.incrementGenerated(title)
				continue
			} else if genErr != nil {
				logger.Infof("[stash-cli] ffmpeg cover generation failed: scene_id=%d title=%q error=%v official_error=%v", item.ID, title, genErr, err)
			}
			if isCoverFetchSkip(err) {
				logger.Infof("[stash-cli] official cover skipped: scene_id=%d title=%q error=%v", item.ID, title, err)
				fetch.incrementSkipped(title)
				continue
			}
			logger.Infof("[stash-cli] official cover failed: scene_id=%d title=%q error=%v", item.ID, title, err)
			fetch.incrementFailed(title)
		}
		result := fetch.setDone()
		logger.Infof("[stash-cli] official cover fetch complete: total=%d ok=%d generated=%d skipped=%d failed=%d", result.Total, result.OK, result.Generated, result.Skipped, result.Failed)
		return coverFetchAllMsg{fetch: fetch, result: result}
	}, pollCoverFetchProgress(fetch))
}

func (m Model) generateCoverFallback(item browse.SceneItem) (covergen.Result, error) {
	if m.generator == nil {
		return covergen.Result{}, fmt.Errorf("ffmpeg cover generator is not configured")
	}
	return m.generator.Generate(m.ctx, covergen.Request{
		SceneID:  item.ID,
		Path:     item.Path,
		Duration: item.Duration,
	})
}

func (m Model) coverFetchItems(ctx context.Context, query browse.Query, current []browse.SceneItem, total int) ([]browse.SceneItem, error) {
	if total <= len(current) || m.browser == nil {
		return current, nil
	}

	query.Page = 1
	query.PerPage = total
	result, err := m.browser.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return current, nil
	}
	return append([]browse.SceneItem(nil), result.Items...), nil
}

func isCoverFetchSkip(err error) bool {
	return err != nil && (errors.Is(err, coverfetch.ErrNoFingerprints) || errors.Is(err, coverfetch.ErrNoMatch) || errors.Is(err, coverfetch.ErrNoImage))
}

func (m Model) executeScan() (tea.Model, tea.Cmd) {
	if m.scanner == nil {
		m.status = "/scan is unavailable: check media_dirs and ffprobe_path"
		return m, nil
	}
	if m.scan != nil {
		m.status = formatScanProgressStatus(m.scan.progress())
		return m, nil
	}

	scan := &scanState{}
	m.scan = scan
	m.status = formatScanProgressStatus(scanner.Progress{})
	return m, tea.Batch(func() tea.Msg {
		result := m.scanner.ScanWithProgress(m.ctx, scan.setProgress)
		scan.setDone(result)
		return scanMsg{scan: scan, result: result}
	}, pollScanProgress(scan))
}

type scanState struct {
	mu       sync.Mutex
	current  scanner.Progress
	finished bool
	result   scanner.Result
}

func (s *scanState) setProgress(progress scanner.Progress) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = progress
}

func (s *scanState) setDone(result scanner.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.result = result
	s.current = scanner.Progress{
		Directories:  result.Directories,
		FilesSeen:    result.FilesSeen,
		FilesScanned: result.FilesScanned,
		LastFile:     result.LastFile,
		Errors:       len(result.Errors),
	}
	s.finished = true
}

func (s *scanState) progress() scanner.Progress {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

func (s *scanState) done() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.finished
}

type coverFetchProgress struct {
	Total     int
	Attempted int
	OK        int
	Generated int
	Skipped   int
	Failed    int
	LastScene string
}

type coverFetchResult struct {
	Total     int
	OK        int
	Generated int
	Skipped   int
	Failed    int
}

type coverFetchState struct {
	mu       sync.Mutex
	current  coverFetchProgress
	finished bool
	result   coverFetchResult
}

func newCoverFetchState(total int) *coverFetchState {
	return &coverFetchState{current: coverFetchProgress{Total: total}}
}

func (s *coverFetchState) incrementOK(scene string) {
	s.increment(scene, func(progress *coverFetchProgress) {
		progress.OK++
	})
}

func (s *coverFetchState) incrementGenerated(scene string) {
	s.increment(scene, func(progress *coverFetchProgress) {
		progress.Generated++
	})
}

func (s *coverFetchState) incrementSkipped(scene string) {
	s.increment(scene, func(progress *coverFetchProgress) {
		progress.Skipped++
	})
}

func (s *coverFetchState) incrementFailed(scene string) {
	s.increment(scene, func(progress *coverFetchProgress) {
		progress.Failed++
	})
}

func (s *coverFetchState) increment(scene string, fn func(*coverFetchProgress)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current.Attempted++
	s.current.LastScene = scene
	fn(&s.current)
}

func (s *coverFetchState) setDone() coverFetchResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.result = coverFetchResult{
		Total:     s.current.Total,
		OK:        s.current.OK,
		Generated: s.current.Generated,
		Skipped:   s.current.Skipped,
		Failed:    s.current.Failed,
	}
	s.finished = true
	return s.result
}

func (s *coverFetchState) progress() coverFetchProgress {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

func (s *coverFetchState) done() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.finished
}

func (m Model) executeEdit(args []string) (tea.Model, tea.Cmd) {
	if m.editor == nil {
		m.status = "/edit is unavailable"
		return m, nil
	}
	if m.cursor < 0 || m.cursor >= len(m.result.Items) {
		m.status = "No scene selected"
		return m, nil
	}

	update, err := edit.ParseArgs(args)
	if err != nil {
		m.status = err.Error()
		return m, nil
	}

	sceneID := m.result.Items[m.cursor].ID
	query := m.query
	return m, func() tea.Msg {
		err := m.editor.Apply(m.ctx, sceneID, update)
		if err != nil {
			return statusMsg{status: err.Error()}
		}
		result, err := m.browser.Search(m.ctx, query)
		return resultMsg{result: result, err: err, status: "Scene updated"}
	}
}

type resultMsg struct {
	result browse.Result
	err    error
	status string
}

type appendResultMsg struct {
	result browse.Result
	cursor int
	err    error
}

type coverMsg struct {
	image image.Image
	err   error
}

type coverFetchedMsg struct {
	result coverfetch.Result
	err    error
}

type coverFetchAllMsg struct {
	fetch  *coverFetchState
	result coverFetchResult
	err    error
}

type coverFetchProgressMsg struct {
	fetch    *coverFetchState
	progress coverFetchProgress
}

type statusMsg struct {
	status string
}

type scanMsg struct {
	scan   *scanState
	result scanner.Result
}

type scanProgressMsg struct {
	scan     *scanState
	progress scanner.Progress
}

func pollScanProgress(scan *scanState) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(scanProgressInterval)
		return scanProgressMsg{scan: scan, progress: scan.progress()}
	}
}

func pollCoverFetchProgress(fetch *coverFetchState) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(scanProgressInterval)
		return coverFetchProgressMsg{fetch: fetch, progress: fetch.progress()}
	}
}

func formatScanProgressStatus(progress scanner.Progress) string {
	status := fmt.Sprintf("Scanning... %d files scanned, %d videos found, %d directories", progress.FilesScanned, progress.FilesSeen, progress.Directories)
	if progress.LastFile != "" {
		status += ", last: " + compactBasename(filepath.Base(progress.LastFile), 10)
	}
	return status
}

func compactBasename(name string, maxLen int) string {
	if maxLen <= 0 || len(name) <= maxLen {
		return name
	}

	const ellipsis = "..."
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	if len(stem) <= len(ellipsis)+2 {
		return name
	}

	prefixLen := min(4, len(stem)/2)
	tailLen := min(4, len(stem)-prefixLen)
	return stem[:prefixLen] + ellipsis + stem[len(stem)-tailLen:] + ext
}

func formatScanStatus(result scanner.Result) string {
	return fmt.Sprintf("Scan complete: %d files, %d directories, %d errors", result.FilesScanned, result.Directories, len(result.Errors))
}

func formatCoverFetchProgressStatus(progress coverFetchProgress) string {
	status := fmt.Sprintf("Fetching official covers... %d/%d, ok:%d generated:%d skipped:%d failed:%d", progress.Attempted, progress.Total, progress.OK, progress.Generated, progress.Skipped, progress.Failed)
	if progress.LastScene != "" {
		status += ", last: " + compactBasename(progress.LastScene, 12)
	}
	return status
}

func formatCoverFetchResultStatus(result coverFetchResult) string {
	return fmt.Sprintf("Official covers fetched: %d ok, %d generated, %d skipped, %d failed", result.OK, result.Generated, result.Skipped, result.Failed)
}

func (m Model) refresh() tea.Cmd {
	return m.refreshWithStatus("")
}

func (m Model) refreshWithStatus(status string) tea.Cmd {
	query := m.query
	return func() tea.Msg {
		result, err := m.browser.Search(m.ctx, query)
		return resultMsg{result: result, err: err, status: status}
	}
}

func (m Model) loadMoreResults() tea.Cmd {
	if m.browser == nil {
		return nil
	}
	query := m.query
	if query.PerPage == 0 {
		query.PerPage = 40
	}
	query.Page = len(m.result.Items)/query.PerPage + 1
	nextCursor := len(m.result.Items)
	return func() tea.Msg {
		result, err := m.browser.Search(m.ctx, query)
		return appendResultMsg{result: result, cursor: nextCursor, err: err}
	}
}

func (m Model) selectedItem() (browse.SceneItem, bool) {
	if m.cursor < 0 || m.cursor >= len(m.result.Items) {
		return browse.SceneItem{}, false
	}
	return m.result.Items[m.cursor], true
}

func (m Model) loadSelectedCover() tea.Cmd {
	if m.mode != ViewGrid || m.covers == nil {
		return nil
	}
	item, ok := m.selectedItem()
	if !ok {
		return func() tea.Msg { return coverMsg{} }
	}

	return func() tea.Msg {
		loaded, err := m.covers.Load(m.ctx, cover.Request{
			SceneID:  item.ID,
			Path:     item.Path,
			Duration: item.Duration,
		})
		if err != nil {
			return coverMsg{err: err}
		}
		if len(loaded.Data) == 0 {
			return coverMsg{}
		}
		img, _, err := image.Decode(bytes.NewReader(loaded.Data))
		if err != nil {
			return coverMsg{err: fmt.Errorf("decode cover: %w", err)}
		}
		bounds := img.Bounds()
		logger.Infof("[stash-cli] loaded cover: scene_id=%d source=%s width=%d height=%d bytes=%d render_mode=%s kitty=%s", item.ID, loaded.Source, bounds.Dx(), bounds.Dy(), len(loaded.Data), pictureModeString(m.pic.Mode()), kittyCapabilityString(picture.KittySupported()))
		return coverMsg{image: img}
	}
}

func (m *Model) resizePicture() tea.Cmd {
	if m.width == 0 || m.height == 0 {
		return nil
	}

	cols := m.width / 3
	if cols < 20 {
		cols = 20
	}
	rows := m.height - 8
	if rows < 4 {
		rows = 4
	}

	return m.pic.SetSize(cols, rows)
}

func (m *Model) ensureKitty(cmd tea.Cmd) tea.Cmd {
	cmds := []tea.Cmd{cmd}
	if m.mode == ViewGrid && picture.KittySupported() == picture.KittyCapabilitySupported && m.pic.Mode() == picture.PictureGlyph {
		cmds = append(cmds, m.pic.Toggle())
	}
	return tea.Batch(cmds...)
}

func pictureModeString(mode picture.PictureMode) string {
	if mode == picture.PictureKitty {
		return "kitty"
	}
	return "glyph"
}

func kittyCapabilityString(capability picture.KittyCapability) string {
	switch capability {
	case picture.KittyCapabilitySupported:
		return "supported"
	case picture.KittyCapabilityUnsupported:
		return "unsupported"
	default:
		return "unknown"
	}
}

func formatDuration(seconds float64) string {
	if seconds <= 0 {
		return "--:--"
	}
	total := int(seconds)
	h := total / 3600
	min := (total % 3600) / 60
	sec := total % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, min, sec)
	}
	return fmt.Sprintf("%02d:%02d", min, sec)
}
