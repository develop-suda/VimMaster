package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/develop-suda/VimMaster/internal/buffer"
	"github.com/develop-suda/VimMaster/internal/stage"
	"github.com/develop-suda/VimMaster/internal/ui"
	"github.com/develop-suda/VimMaster/internal/vim"
)

// Screen represents the current screen state.
type Screen int

const (
	ScreenTitle Screen = iota
	ScreenStageSelect
	ScreenGame
	ScreenClear
)

// Model is the main Bubble Tea model for VimMaster.
type Model struct {
	// App state
	screen       Screen
	stages       []stage.Stage
	currentStage int
	cleared      []bool

	// Game state
	buf       *buffer.Buffer
	mode      vim.Mode
	pendingOp vim.PendingOp
	strokes   int
	showHint  bool

	// UI state
	width  int
	height int
	rating stage.Rating
}

// New creates a new VimMaster application model.
func New() (Model, error) {
	stages, err := stage.LoadAllStages()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load stages: %w", err)
	}

	return Model{
		screen:  ScreenTitle,
		stages:  stages,
		cleared: make([]bool, len(stages)),
		width:   80,
		height:  24,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.screen {
		case ScreenTitle:
			return m.updateTitle(msg)
		case ScreenStageSelect:
			return m.updateStageSelect(msg)
		case ScreenGame:
			return m.updateGame(msg)
		case ScreenClear:
			return m.updateClear(msg)
		}
	}

	return m, nil
}

func (m Model) updateTitle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		m.screen = ScreenStageSelect
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateStageSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.currentStage < len(m.stages)-1 {
			m.currentStage++
		}
	case "k", "up":
		if m.currentStage > 0 {
			m.currentStage--
		}
	case "enter":
		m.startStage()
		m.screen = ScreenGame
	case "q", "esc":
		m.screen = ScreenTitle
	}
	return m, nil
}

func (m *Model) startStage() {
	s := m.stages[m.currentStage]
	m.buf = buffer.NewBuffer(s.InitialText)
	m.mode = vim.NormalMode
	m.pendingOp = vim.OpNone
	m.strokes = 0
	m.showHint = false
}

func (m Model) updateGame(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Special keys in game
	if msg.String() == "ctrl+q" {
		m.screen = ScreenStageSelect
		return m, nil
	}

	// Toggle hint
	if m.mode == vim.NormalMode && msg.String() == "?" {
		m.showHint = !m.showHint
		return m, nil
	}

	// Handle key based on mode
	var result vim.HandleResult
	switch m.mode {
	case vim.NormalMode:
		result = vim.HandleNormalMode(msg, m.buf, &m.pendingOp)
	case vim.InsertMode:
		result = vim.HandleInsertMode(msg, m.buf)
	}

	m.strokes += result.Strokes

	if result.ModeChanged {
		m.mode = result.NewMode
	}

	// Check clear condition
	s := m.stages[m.currentStage]
	if stage.CheckClear(m.buf, s.ExpectedText) {
		m.rating = stage.CalculateRating(m.strokes, s.MinStrokes)
		m.cleared[m.currentStage] = true
		m.screen = ScreenClear
	}

	return m, nil
}

func (m Model) updateClear(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		// Go to next stage if available
		if m.currentStage < len(m.stages)-1 {
			m.currentStage++
			m.startStage()
			m.screen = ScreenGame
		} else {
			m.screen = ScreenStageSelect
		}
	case "r":
		// Retry current stage
		m.startStage()
		m.screen = ScreenGame
	case "q", "esc":
		m.screen = ScreenStageSelect
	}
	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case ScreenTitle:
		return m.viewTitle()
	case ScreenStageSelect:
		return m.viewStageSelect()
	case ScreenGame:
		return m.viewGame()
	case ScreenClear:
		return m.viewClear()
	}
	return ""
}

func (m Model) viewTitle() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Align(lipgloss.Center)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Align(lipgloss.Center)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Align(lipgloss.Center)

	logo := `
 ██╗   ██╗██╗███╗   ███╗
 ██║   ██║██║████╗ ████║
 ██║   ██║██║██╔████╔██║
 ╚██╗ ██╔╝██║██║╚██╔╝██║
  ╚████╔╝ ██║██║ ╚═╝ ██║
   ╚═══╝  ╚═╝╚═╝     ╚═╝
  ███╗   ███╗ █████╗ ███████╗████████╗███████╗██████╗
  ████╗ ████║██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗
  ██╔████╔██║███████║███████╗   ██║   █████╗  ██████╔╝
  ██║╚██╔╝██║██╔══██║╚════██║   ██║   ██╔══╝  ██╔══██╗
  ██║ ╚═╝ ██║██║  ██║███████║   ██║   ███████╗██║  ██║
  ╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝`

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(logo),
		"",
		subtitleStyle.Render("vimtutor よりも楽しく、ゲーム感覚で指に覚えさせる"),
		"",
		instructionStyle.Render("[Enter] スタート  [q] 終了"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) viewStageSelect() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(1, 0)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	clearedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1, 0)

	var content string
	content += titleStyle.Render("📖 ステージ選択") + "\n\n"

	for i, s := range m.stages {
		prefix := "  "
		if i == m.currentStage {
			prefix = "▸ "
		}

		status := ""
		if m.cleared[i] {
			status = " ✅"
		}

		line := fmt.Sprintf("%s[%s] %s%s", prefix, s.Category, s.Name, status)

		if i == m.currentStage {
			content += selectedStyle.Render(line) + "\n"
		} else if m.cleared[i] {
			content += clearedStyle.Render(line) + "\n"
		} else {
			content += normalStyle.Render(line) + "\n"
		}
	}

	content += footerStyle.Render("\n[j/k] 選択  [Enter] 開始  [q] 戻る")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) viewGame() string {
	if m.buf == nil {
		return "Loading..."
	}

	s := m.stages[m.currentStage]
	clearedCount := 0
	for _, c := range m.cleared {
		if c {
			clearedCount++
		}
	}

	header := ui.RenderHeader(s.Name, clearedCount, len(m.stages), m.width)
	editor := ui.RenderEditor(m.buf, s.Description, m.width)
	statusLine := ui.RenderStatusLine(m.mode, m.buf.CursorRow, m.buf.CursorCol, m.pendingOp, m.strokes, m.width)
	hint := ui.RenderHint(s.Hint, m.showHint, m.width)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	footer := footerStyle.Render("  [Ctrl+Q] ステージ選択に戻る  [?] ヒント表示/非表示")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		editor,
		statusLine,
		hint,
		footer,
	)
}

func (m Model) viewClear() string {
	s := m.stages[m.currentStage]

	clearStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42")).
		Align(lipgloss.Center)

	ratingStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		Align(lipgloss.Center)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Align(lipgloss.Center)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Align(lipgloss.Center)

	ratingDisplay := fmt.Sprintf("評価: %s", m.rating.String())
	switch m.rating {
	case stage.RatingS:
		ratingDisplay = "🏆 評価: S (パーフェクト!)"
	case stage.RatingA:
		ratingDisplay = "🥇 評価: A (すばらしい!)"
	case stage.RatingB:
		ratingDisplay = "🥈 評価: B (よくできました!)"
	default:
		ratingDisplay = "🥉 評価: C (クリア!)"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		clearStyle.Render("🎉 ステージクリア! 🎉"),
		"",
		infoStyle.Render(s.Name),
		"",
		ratingStyle.Render(ratingDisplay),
		"",
		infoStyle.Render(fmt.Sprintf("手数: %d (最小: %d)", m.strokes, s.MinStrokes)),
		"",
		instructionStyle.Render("[Enter] 次のステージ  [r] リトライ  [q] ステージ選択"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
