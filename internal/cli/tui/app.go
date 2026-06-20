package tui

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strconv"
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
	"github.com/stashapp/stash/internal/cli/scanner"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type Browser interface {
	Search(context.Context, browse.Query) (browse.Result, error)
}

type Editor interface {
	SetSceneRating(context.Context, int, int) error
	DeleteScene(context.Context, int) error
	SetPerformerRating(context.Context, int, int) error
}

type CoverLoader interface {
	Load(context.Context, cover.Request) (cover.Cover, error)
}

type Player interface {
	Play(context.Context, browse.SceneItem) error
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
	ctx        context.Context
	browser    Browser
	editor     Editor
	covers     CoverLoader
	player     Player
	scanner    Scanner
	pic        picture.Model
	gridPics   map[int]*gridCover
	forceKitty bool

	mode              ViewMode
	query             browse.Query
	result            browse.Result
	cursor            int
	gridStart         int
	input             string
	commandMode       bool
	completionQuery   string
	completionMatches []string
	completionIndex   int
	showDetails       bool
	performerCursor   int
	detailEditMode    string
	detailEditInput   string
	confirmDelete     bool
	status            string
	scan              *scanState
	width             int
	height            int
}

type gridCoverState int

const (
	gridCoverMissing gridCoverState = iota
	gridCoverLoading
	gridCoverReady
	gridCoverEmpty
	gridCoverFailed
)

type gridCover struct {
	pic   picture.Model
	state gridCoverState
}

type Deps struct {
	Browser    Browser
	Editor     Editor
	Covers     CoverLoader
	Player     Player
	Scanner    Scanner
	ForceKitty bool
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
		ctx:        ctx,
		browser:    deps.Browser,
		editor:     deps.Editor,
		covers:     deps.Covers,
		player:     deps.Player,
		scanner:    deps.Scanner,
		pic:        pic,
		gridPics:   map[int]*gridCover{},
		forceKitty: deps.ForceKitty,
		mode:       mode,
		status:     "Press : for commands",
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
		m.gridStart = normalizeGridStart(m.gridStart, len(m.result.Items), gridColumns(m.width), visibleGridRows(m.height))
		return m, tea.Batch(m.resizeGridPictures(), m.loadVisibleCovers())
	case resultMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.result = msg.result
		m.cursor = 0
		m.gridStart = 0
		if msg.status != "" {
			m.status = msg.status
		} else {
			m.status = fmt.Sprintf("%d results", msg.result.Total)
		}
		return m, m.loadVisibleCovers()
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
		m.gridStart = m.gridStartForCursor(m.cursor)
		m.status = fmt.Sprintf("%d results", m.result.Total)
		return m, m.loadVisibleCovers()
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
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	case coverMsg:
		if msg.sceneID != 0 {
			gridPic := m.gridPic(msg.sceneID)
			if msg.err != nil {
				gridPic.state = gridCoverFailed
				m.status = msg.err.Error()
				return m, nil
			}
			if msg.image == nil {
				gridPic.state = gridCoverEmpty
				return m, nil
			}
			gridPic.state = gridCoverReady
			return m, m.ensureKitty(gridPic.pic.SetImage(msg.image))
		}
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		cmd := m.pic.SetImage(msg.image)
		return m, m.ensureKitty(cmd)
	case playMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.status = "Played: " + msg.item.Title
	}

	return m, m.updatePictures(msg)
}

func (m Model) View() tea.View {
	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Render("stash-cli")
	fmt.Fprintf(&b, "%s\n", title)
	fmt.Fprintf(&b, "%s\n\n", m.status)

	if len(m.result.Items) == 0 {
		b.WriteString("No scenes to display\n")
	} else {
		b.WriteString(m.renderSceneGrid())
		if m.showDetails {
			b.WriteString("\n")
			b.WriteString(m.renderSelectedDetails())
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	if m.commandMode {
		if matches := m.commandCompletionMatches(); len(matches) > 0 {
			b.WriteString(lipgloss.NewStyle().Faint(true).Render("matches: " + strings.Join(matches, " ")))
			b.WriteString("\n")
		}
		b.WriteString(lipgloss.NewStyle().Faint(true).Render(command.Help()))
		b.WriteString("\n:")
		b.WriteString(m.input)
	} else {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("h/j/k/l browse, enter play, space details, : command, :q quit"))
	}

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

const (
	gridTileWidth = 28
	gridTileGap   = 2
	gridCoverRows = 8
	gridTileRows  = gridCoverRows + 5
)

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.commandMode {
		return m.handleCommandKey(msg)
	}
	if m.showDetails {
		return m.handleDetailsKey(msg)
	}
	return m.handleNormalKey(msg)
}

func (m Model) handleCommandKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.commandMode = false
		m.input = ""
		m.clearCompletion()
		return m, nil
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
		m.clearCompletion()
	case "tab":
		m.applyCommandCompletion()
	case "enter":
		return m.executeInput()
	default:
		if text := msg.Key().Text; text != "" {
			m.input += text
			m.clearCompletion()
		}
	}

	return m, nil
}

func (m *Model) clearCompletion() {
	m.completionQuery = ""
	m.completionMatches = nil
	m.completionIndex = 0
}

func (m *Model) applyCommandCompletion() {
	query := m.completionSeed()
	if query == "" {
		return
	}

	if m.completionQuery != query || len(m.completionMatches) == 0 {
		m.completionQuery = query
		m.completionMatches = fuzzyCommandMatches(query)
		m.completionIndex = 0
	} else if len(m.completionMatches) > 0 {
		m.completionIndex = (m.completionIndex + 1) % len(m.completionMatches)
	}

	if len(m.completionMatches) == 0 {
		return
	}
	m.input = m.completionMatches[m.completionIndex] + " "
}

func (m Model) commandCompletionMatches() []string {
	if m.completionQuery != "" && len(m.completionMatches) > 0 {
		return append([]string(nil), m.completionMatches...)
	}
	return fuzzyCommandMatches(m.completionSeed())
}

func (m Model) completionSeed() string {
	if m.completionQuery != "" {
		trimmed := strings.TrimSpace(m.input)
		for _, match := range m.completionMatches {
			if trimmed == match {
				return m.completionQuery
			}
		}
	}

	input := strings.TrimSpace(m.input)
	if input == "" || strings.Contains(input, " ") {
		return ""
	}
	if strings.HasPrefix(input, "/") {
		input = strings.TrimPrefix(input, "/")
	}
	return input
}

func fuzzyCommandMatches(query string) []string {
	query = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(query, "/")))
	if query == "" {
		return nil
	}

	var matches []string
	for _, name := range command.CompletableCommands() {
		if fuzzyMatch(query, name) {
			matches = append(matches, "/"+name)
		}
	}
	return matches
}

func fuzzyMatch(query, candidate string) bool {
	if query == "" {
		return true
	}
	candidate = strings.ToLower(candidate)
	pos := 0
	for _, r := range query {
		found := false
		for pos < len(candidate) {
			if rune(candidate[pos]) == r {
				pos++
				found = true
				break
			}
			pos++
		}
		if !found {
			return false
		}
	}
	return true
}

func (m Model) handleNormalKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.showDetails = false
	case ":":
		m.commandMode = true
		m.input = ""
		m.clearCompletion()
	case "enter":
		return m.executePlay()
	case " ", "space":
		m.showDetails = !m.showDetails
	case "up", "k":
		return m.moveCursor(-gridColumns(m.width))
	case "down", "j":
		return m.moveCursor(gridColumns(m.width))
	case "left", "h":
		return m.moveCursor(-1)
	case "right", "l":
		return m.moveCursor(1)
	}

	return m, nil
}

func (m Model) handleDetailsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.detailEditMode != "" {
		return m.handleDetailInputKey(msg)
	}
	if m.confirmDelete {
		switch msg.String() {
		case "y", "Y":
			return m.executeDeleteSelectedScene()
		case "n", "N", "esc", " ", "space":
			m.confirmDelete = false
			m.status = "Delete cancelled"
			return m, nil
		}
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", " ", "space":
		m.showDetails = false
		m.confirmDelete = false
		m.detailEditMode = ""
		m.detailEditInput = ""
	case "j", "down":
		m.movePerformerCursor(1)
	case "k", "up":
		m.movePerformerCursor(-1)
	case "r":
		m.detailEditMode = "scene-rating"
		m.detailEditInput = ""
		m.status = "Scene rating: enter 0-100"
	case "R":
		item, ok := m.selectedItem()
		if !ok || len(item.Performers) == 0 {
			m.status = "No performer selected"
			return m, nil
		}
		m.detailEditMode = "performer-rating"
		m.detailEditInput = ""
		m.status = "Performer rating: enter 0-100"
	case "d":
		m.confirmDelete = true
		m.status = "Delete scene? y/N"
	}

	return m, nil
}

func (m Model) handleDetailInputKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.detailEditMode = ""
		m.detailEditInput = ""
		m.status = "Edit cancelled"
	case "backspace":
		if len(m.detailEditInput) > 0 {
			m.detailEditInput = m.detailEditInput[:len(m.detailEditInput)-1]
		}
	case "enter":
		rating, err := strconv.Atoi(strings.TrimSpace(m.detailEditInput))
		if err != nil || rating < 0 || rating > 100 {
			m.status = "Rating must be between 0 and 100"
			return m, nil
		}
		if m.detailEditMode == "scene-rating" {
			return m.executeSetSceneRating(rating)
		}
		if m.detailEditMode == "performer-rating" {
			return m.executeSetPerformerRating(rating)
		}
	default:
		if text := msg.Key().Text; text != "" {
			m.detailEditInput += text
		}
	}

	return m, nil
}

func (m *Model) movePerformerCursor(delta int) {
	item, ok := m.selectedItem()
	if !ok || len(item.Performers) == 0 {
		m.performerCursor = 0
		return
	}
	next := m.performerCursor + delta
	if next < 0 {
		next = 0
	}
	if next >= len(item.Performers) {
		next = len(item.Performers) - 1
	}
	m.performerCursor = next
}

func (m Model) moveCursor(delta int) (tea.Model, tea.Cmd) {
	if delta == 0 || len(m.result.Items) == 0 {
		return m, nil
	}

	next := m.cursor + delta
	if next >= 0 && next < len(m.result.Items) {
		m.cursor = next
		m.gridStart = m.gridStartForCursor(next)
		return m, m.loadVisibleCovers()
	}
	if delta > 0 && len(m.result.Items) < m.result.Total {
		return m, m.loadMoreResults()
	}

	return m, nil
}

func gridColumns(width int) int {
	if width <= 0 {
		return 1
	}
	cols := (width + gridTileGap) / (gridTileWidth + gridTileGap)
	if cols < 1 {
		return 1
	}
	return cols
}

func visibleGridStart(cursor, total, cols, rows int) int {
	if total <= 0 || cols <= 0 || rows <= 0 {
		return 0
	}
	visible := cols * rows
	if visible >= total {
		return 0
	}
	cursorRow := cursor / cols
	firstRow := 0
	lastVisibleRow := rows - 1
	if cursorRow > lastVisibleRow {
		firstRow = cursorRow - lastVisibleRow
	}
	start := firstRow * cols
	maxStart := total - visible
	if start > maxStart {
		start = maxStart
	}
	if start < 0 {
		return 0
	}
	return start
}

func normalizeGridStart(start, total, cols, rows int) int {
	if total <= 0 || cols <= 0 || rows <= 0 {
		return 0
	}
	visible := cols * rows
	if visible >= total {
		return 0
	}

	start = (start / cols) * cols
	lastRow := (total - 1) / cols
	maxStart := (lastRow - rows + 1) * cols
	if maxStart < 0 {
		maxStart = 0
	}
	if start > maxStart {
		return maxStart
	}
	if start < 0 {
		return 0
	}
	return start
}

func (m Model) gridStartForCursor(cursor int) int {
	cols := gridColumns(m.width)
	rows := visibleGridRows(m.height)
	start := normalizeGridStart(m.gridStart, len(m.result.Items), cols, rows)
	visible := cols * rows
	if cursor < start {
		rowStart := (cursor / cols) * cols
		return normalizeGridStart(rowStart-(rows-1)*cols, len(m.result.Items), cols, rows)
	}
	if cursor >= start+visible {
		rowStart := (cursor / cols) * cols
		return normalizeGridStart(rowStart, len(m.result.Items), cols, rows)
	}
	return start
}

func visibleGridRows(height int) int {
	rows := (height - 7) / gridTileRows
	if rows < 1 {
		return 1
	}
	return rows
}

func compactText(s string, width int) string {
	runes := []rune(strings.TrimSpace(s))
	if width <= 0 || len(runes) <= width {
		return string(runes)
	}
	if width <= 3 {
		return string(runes[:width])
	}
	left := (width - 3) / 2
	right := width - 3 - left
	return string(runes[:left]) + "..." + string(runes[len(runes)-right:])
}

func (m Model) renderSceneGrid() string {
	cols := gridColumns(m.width)
	rows := visibleGridRows(m.height)
	start := m.gridStartForCursor(m.cursor)
	end := start + cols*rows
	if end > len(m.result.Items) {
		end = len(m.result.Items)
	}

	var b strings.Builder
	for rowStart := start; rowStart < end; rowStart += cols {
		rowEnd := rowStart + cols
		if rowEnd > end {
			rowEnd = end
		}
		var tiles []string
		for i := rowStart; i < rowEnd; i++ {
			tiles = append(tiles, m.renderSceneTile(i, m.result.Items[i]))
		}
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tiles...))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderSceneTile(index int, item browse.SceneItem) string {
	selected := index == m.cursor
	border := lipgloss.NormalBorder()
	style := lipgloss.NewStyle().
		Width(gridTileWidth).
		Height(gridTileRows).
		Border(border).
		Padding(0, 1)
	if selected {
		style = style.BorderForeground(lipgloss.Color("12")).Bold(true)
	} else {
		style = style.BorderForeground(lipgloss.Color("8"))
	}

	coverView := "[cover loading]"
	if gridPic, ok := m.gridPics[item.ID]; ok {
		switch gridPic.state {
		case gridCoverReady:
			coverView = fitBlock(gridPic.pic.View().Content, gridTileWidth-4, gridCoverRows)
		case gridCoverEmpty:
			coverView = "[no cover]"
		case gridCoverFailed:
			coverView = "[cover error]"
		default:
			coverView = "[cover loading]"
		}
	}

	meta := compactText(item.Title, gridTileWidth-4)
	performer := compactText(performerSummary(item.Performers), gridTileWidth-4)
	details := strings.TrimSpace(strings.Join([]string{formatDuration(item.Duration), item.Date}, " "))
	details = compactText(details, gridTileWidth-4)

	return style.Render(strings.Join([]string{
		fitBlock(coverView, gridTileWidth-4, gridCoverRows),
		meta,
		performer,
		details,
	}, "\n")) + strings.Repeat(" ", gridTileGap)
}

func performerSummary(performers []browse.PerformerItem) string {
	if len(performers) == 0 {
		return "No performers"
	}
	if len(performers) == 1 {
		return performers[0].Name
	}
	return fmt.Sprintf("%s <%d omitted>", performers[0].Name, len(performers)-1)
}

func (m Model) renderSelectedDetails() string {
	item, ok := m.selectedItem()
	if !ok {
		return ""
	}

	lines := []string{
		"Title: " + item.Title,
		"Path: " + item.Path,
		"Rating: " + formatRating(item.Rating),
	}
	if item.Duration > 0 {
		lines = append(lines, "Duration: "+formatDuration(item.Duration))
	}
	if item.Date != "" {
		lines = append(lines, "Date: "+item.Date)
	}
	if item.Studio != "" {
		lines = append(lines, "Studio: "+item.Studio)
	}
	lines = append(lines, m.renderPerformerDetails(item.Performers)...)
	if len(item.Tags) > 0 {
		lines = append(lines, "Tags: "+strings.Join(item.Tags, ", "))
	}
	lines = append(lines, "")
	lines = append(lines, "Scene: r rating, d delete")
	lines = append(lines, "Performers:")
	lines = append(lines, m.renderPerformerDetails(item.Performers)...)
	if m.detailEditMode != "" {
		lines = append(lines, "Input: "+m.detailEditInput)
	}
	if m.confirmDelete {
		lines = append(lines, "Delete scene? y/N")
	}

	width := m.width - 2
	if width < 60 {
		width = 60
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(0, 1).
		Width(width).
		Render(strings.Join(lines, "\n"))
}

func (m Model) renderPerformerDetails(performers []browse.PerformerItem) []string {
	if len(performers) == 0 {
		return []string{"  No performers"}
	}
	cursor := m.performerCursor
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(performers) {
		cursor = len(performers) - 1
	}

	lines := make([]string, 0, len(performers))
	for i, performer := range performers {
		prefix := "  "
		if i == cursor {
			prefix = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%s rating:%s", prefix, performer.Name, formatRating(performer.Rating)))
	}
	return lines
}

func formatRating(rating *int) string {
	if rating == nil {
		return "--"
	}
	return strconv.Itoa(*rating)
}

func fitBlock(content string, width, height int) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	if height < 0 {
		height = 0
	}
	var b strings.Builder
	for i := 0; i < height; i++ {
		line := ""
		if i < len(lines) {
			line = lines[i]
			if !strings.Contains(line, "\x1b") {
				line = compactText(line, width)
			}
		}
		b.WriteString(line)
		if pad := width - lipgloss.Width(line); pad > 0 {
			b.WriteString(strings.Repeat(" ", pad))
		}
		if i != height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func (m Model) executeInput() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.input)
	m.input = ""
	m.commandMode = false
	if input == "" {
		return m, nil
	}
	if input == "q" || input == "quit" {
		return m, tea.Quit
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
	case "random":
		return m.executeRandom(cmd.Args)
	case "clear":
		m.query = browse.Query{Page: 1, PerPage: 40}
		return m, m.refresh()
	case "view":
		if len(cmd.Args) != 1 || cmd.Args[0] != string(ViewGrid) {
			m.status = "Usage: /view grid"
			return m, nil
		}
		m.mode = ViewMode(cmd.Args[0])
		m.status = "View mode changed"
		return m, m.loadVisibleCovers()
	case "scan":
		return m.executeScan()
	case "cover":
		m.status = "Unknown command: cover"
	case "open":
		m.status = "Unknown command: open"
	case "play":
		m.status = "Unknown command: play"
	case "edit":
		m.status = "Unknown command: edit"
	case "help":
		m.status = command.Help()
	case "quit", "q":
		return m, tea.Quit
	default:
		m.status = "Unknown command: " + cmd.Name
	}

	return m, nil
}

func (m Model) executeSetSceneRating(rating int) (tea.Model, tea.Cmd) {
	if m.editor == nil {
		m.status = "Editing is unavailable"
		return m, nil
	}
	item, ok := m.selectedItem()
	if !ok {
		m.status = "No scene selected"
		return m, nil
	}
	m.detailEditMode = ""
	m.detailEditInput = ""
	return m, func() tea.Msg {
		if err := m.editor.SetSceneRating(m.ctx, item.ID, rating); err != nil {
			return statusMsg{status: err.Error()}
		}
		result, err := m.browser.Search(m.ctx, m.query)
		return resultMsg{result: result, err: err, status: "Scene rating updated"}
	}
}

func (m Model) executeDeleteSelectedScene() (tea.Model, tea.Cmd) {
	if m.editor == nil {
		m.status = "Editing is unavailable"
		return m, nil
	}
	item, ok := m.selectedItem()
	if !ok {
		m.status = "No scene selected"
		return m, nil
	}
	m.confirmDelete = false
	m.showDetails = false
	return m, func() tea.Msg {
		if err := m.editor.DeleteScene(m.ctx, item.ID); err != nil {
			return statusMsg{status: err.Error()}
		}
		result, err := m.browser.Search(m.ctx, m.query)
		return resultMsg{result: result, err: err, status: "Scene deleted"}
	}
}

func (m Model) executeSetPerformerRating(rating int) (tea.Model, tea.Cmd) {
	if m.editor == nil {
		m.status = "Editing is unavailable"
		return m, nil
	}
	item, ok := m.selectedItem()
	if !ok || len(item.Performers) == 0 {
		m.status = "No performer selected"
		return m, nil
	}
	cursor := m.performerCursor
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(item.Performers) {
		cursor = len(item.Performers) - 1
	}
	performer := item.Performers[cursor]
	m.detailEditMode = ""
	m.detailEditInput = ""
	return m, func() tea.Msg {
		if err := m.editor.SetPerformerRating(m.ctx, performer.ID, rating); err != nil {
			return statusMsg{status: err.Error()}
		}
		result, err := m.browser.Search(m.ctx, m.query)
		return resultMsg{result: result, err: err, status: "Performer rating updated"}
	}
}

func (m Model) executeRandom(args []string) (tea.Model, tea.Cmd) {
	if len(args) > 1 {
		m.status = "Usage: /random <n>"
		return m, nil
	}

	limit := m.visibleGridCapacity()
	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil || n <= 0 {
			m.status = "Usage: /random <n>"
			return m, nil
		}
		limit = n
	}

	m.query = browse.Query{
		Page:      1,
		PerPage:   limit,
		Sort:      "random",
		Direction: models.SortDirectionEnumAsc,
	}
	return m, m.refresh()
}

func (m Model) visibleGridCapacity() int {
	capacity := gridColumns(m.width) * visibleGridRows(m.height)
	if capacity < 1 {
		return 1
	}
	return capacity
}

func (m Model) executePlay() (tea.Model, tea.Cmd) {
	if m.player == nil {
		m.status = "Playback is unavailable: configure ffplay_path"
		return m, nil
	}
	item, ok := m.selectedItem()
	if !ok {
		m.status = "No scene selected"
		return m, nil
	}

	m.status = "Playing: " + item.Title
	return m, func() tea.Msg {
		err := m.player.Play(m.ctx, item)
		return playMsg{item: item, err: err}
	}
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
	sceneID int
	image   image.Image
	err     error
}

type playMsg struct {
	item browse.SceneItem
	err  error
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

func (m *Model) loadVisibleCovers() tea.Cmd {
	if m.covers == nil {
		return nil
	}

	var cmds []tea.Cmd
	for _, item := range m.visibleItems() {
		gridPic := m.gridPic(item.ID)
		if gridPic.state != gridCoverMissing {
			continue
		}
		gridPic.state = gridCoverLoading
		cmds = append(cmds, m.loadCover(item))
	}

	return tea.Batch(cmds...)
}

func (m Model) loadCover(item browse.SceneItem) tea.Cmd {
	return func() tea.Msg {
		loaded, err := m.covers.Load(m.ctx, cover.Request{
			SceneID:  item.ID,
			Path:     item.Path,
			Duration: item.Duration,
		})
		if err != nil {
			return coverMsg{sceneID: item.ID, err: err}
		}
		if len(loaded.Data) == 0 {
			return coverMsg{sceneID: item.ID}
		}
		img, _, err := image.Decode(bytes.NewReader(loaded.Data))
		if err != nil {
			return coverMsg{sceneID: item.ID, err: fmt.Errorf("decode cover: %w", err)}
		}
		bounds := img.Bounds()
		logger.Infof("[stash-cli] loaded cover: scene_id=%d source=%s width=%d height=%d bytes=%d render_mode=%s kitty=%s", item.ID, loaded.Source, bounds.Dx(), bounds.Dy(), len(loaded.Data), pictureModeString(m.gridPicMode(item.ID)), kittyCapabilityString(picture.KittySupported()))
		return coverMsg{sceneID: item.ID, image: img}
	}
}

func (m Model) visibleItems() []browse.SceneItem {
	if len(m.result.Items) == 0 {
		return nil
	}
	cols := gridColumns(m.width)
	rows := visibleGridRows(m.height)
	start := m.gridStartForCursor(m.cursor)
	end := start + cols*rows
	if end > len(m.result.Items) {
		end = len(m.result.Items)
	}
	return m.result.Items[start:end]
}

func (m *Model) gridPic(sceneID int) *gridCover {
	if m.gridPics == nil {
		m.gridPics = map[int]*gridCover{}
	}
	if existing := m.gridPics[sceneID]; existing != nil {
		return existing
	}

	pic := picture.NewWithConfig(picture.Config{
		KittyID: 8000 + sceneID,
		Fit:     picture.FitCover,
		Anchor:  picture.AnchorCenter,
	})
	if m.forceKitty {
		_ = pic.Toggle()
	}
	_ = pic.SetSize(gridTileWidth-4, gridCoverRows)
	gridPic := &gridCover{pic: pic}
	m.gridPics[sceneID] = gridPic
	return gridPic
}

func (m *Model) resizeGridPictures() tea.Cmd {
	var cmds []tea.Cmd
	for _, gridPic := range m.gridPics {
		if cmd := gridPic.pic.SetSize(gridTileWidth-4, gridCoverRows); cmd != nil {
			cmds = append(cmds, m.ensureKitty(cmd))
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) updatePictures(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if m.mode == ViewGrid && picture.KittySupported() == picture.KittyCapabilitySupported && m.pic.Mode() == picture.PictureGlyph {
		cmds = append(cmds, m.pic.Toggle())
	}
	if cmd := m.pic.Update(msg); cmd != nil {
		cmds = append(cmds, m.ensureKitty(cmd))
	}
	for _, gridPic := range m.gridPics {
		if cmd := gridPic.pic.Update(msg); cmd != nil {
			cmds = append(cmds, m.ensureKitty(cmd))
		}
	}
	return tea.Batch(cmds...)
}

func (m Model) gridPicMode(sceneID int) picture.PictureMode {
	if gridPic := m.gridPics[sceneID]; gridPic != nil {
		return gridPic.pic.Mode()
	}
	return m.pic.Mode()
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
