package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	formAlias = iota
	formDesc
	formPath
	formTmux
	formScript
)

var fieldTitles = []string{
	"Alias",
	"Description (optional)",
	"Path",
	"Tmux Window (optional)",
	"Post-jump script (optional)",
}

var fieldDescs = []string{
	"Short name for the bookmark",
	"",
	"Directory path to bookmark",
	"Tmux window name to create/switch to",
	"Script/command to run after jumping",
}

// BookmarkFormModel is a form for creating a new bookmark.
type BookmarkFormModel struct {
	inputs     []textinput.Model
	focused    int
	theme      Theme
	responsive *ResponsiveManager
	completed  bool
	cancelled  bool
}

// NewBookmarkFormModel creates a new bookmark form with optional default values.
func NewBookmarkFormModel(theme Theme, defaultAlias, defaultPath string) BookmarkFormModel {
	inputs := make([]textinput.Model, 5)

	inputs[formAlias] = textinput.New()
	inputs[formAlias].Placeholder = defaultAlias
	inputs[formAlias].Focus()
	inputs[formAlias].Prompt = ""

	inputs[formDesc] = textinput.New()
	inputs[formDesc].Placeholder = "Optional description"
	inputs[formDesc].Prompt = ""

	inputs[formPath] = textinput.New()
	inputs[formPath].Placeholder = defaultPath
	inputs[formPath].Prompt = ""

	inputs[formTmux] = textinput.New()
	inputs[formTmux].Placeholder = "Optional tmux window name"
	inputs[formTmux].Prompt = ""

	inputs[formScript] = textinput.New()
	inputs[formScript].Placeholder = "Optional post-jump script"
	inputs[formScript].Prompt = ""

	return BookmarkFormModel{
		inputs:     inputs,
		focused:    0,
		theme:      theme,
		responsive: NewResponsiveManager(80),
	}
}

// Init initializes the form.
func (m BookmarkFormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the form.
func (m BookmarkFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.responsive.SetWidth(msg.Width)
		// Update inputs width based on content width
		inputWidth := m.responsive.MaxContentWidth() - 6
		if inputWidth < 20 {
			inputWidth = 20
		}
		for i := range m.inputs {
			m.inputs[i].Width = inputWidth
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			if m.focused == len(m.inputs)-1 {
				m.completed = true
				return m, tea.Quit
			}
			m.inputs[m.focused].Blur()
			m.focused++
			m.inputs[m.focused].Focus()
			return m, nil

		case "shift+tab", "up":
			if m.focused > 0 {
				m.inputs[m.focused].Blur()
				m.focused--
				m.inputs[m.focused].Focus()
			}
			return m, nil

		case "tab", "down":
			if m.focused < len(m.inputs)-1 {
				m.inputs[m.focused].Blur()
				m.focused++
				m.inputs[m.focused].Focus()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

// View renders the form.
func (m BookmarkFormModel) View() string {
	if m.completed || m.cancelled {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Headings).
		Render("Create New Bookmark")

	help := lipgloss.NewStyle().
		Foreground(m.theme.Muted).
		Render("↑/↓ or tab/shift+tab: navigate • enter: next/submit • esc: cancel")

	// Calculate sliding window for exactly 3 items
	start := m.focused - 1
	if start < 0 {
		start = 0
	}
	if start+3 > len(m.inputs) {
		start = len(m.inputs) - 3
	}
	end := start + 3

	var items []string
	for i := start; i < end; i++ {
		var itemLines []string

		// Title
		titleStyle := lipgloss.NewStyle().Bold(true)
		if i == m.focused {
			titleStyle = titleStyle.Foreground(m.theme.Secondary)
		} else {
			titleStyle = titleStyle.Foreground(m.theme.Text)
		}
		itemLines = append(itemLines, titleStyle.Render(fieldTitles[i]))

		// Description
		if fieldDescs[i] != "" {
			descStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
			itemLines = append(itemLines, descStyle.Render(fieldDescs[i]))
		}

		// Input wrapped in border
		// We set the width of the input wrapper to the responsive content width
		borderWidth := m.responsive.MaxContentWidth() - 2
		if borderWidth < 22 {
			borderWidth = 22
		}
		inputStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Width(borderWidth).
			BorderStyle(lipgloss.RoundedBorder())

		if i == m.focused {
			inputStyle = inputStyle.BorderForeground(m.theme.Secondary)
		} else {
			inputStyle = inputStyle.BorderForeground(m.theme.Border)
		}

		inputBox := inputStyle.Render(m.inputs[i].View())
		itemLines = append(itemLines, inputBox)

		items = append(items, strings.Join(itemLines, "\n"))
	}

	content := strings.Join([]string{
		title,
		"",
		strings.Join(items, "\n\n"),
		"",
		help,
	}, "\n")

	return m.responsive.AdaptiveFrameStyle(m.theme).Render(content)
}

// Values returns the form values.
func (m BookmarkFormModel) Values() (alias, path, desc, tmuxWindowName, postJumpScript string) {
	alias = m.inputs[formAlias].Value()
	if alias == "" {
		alias = m.inputs[formAlias].Placeholder
	}
	path = m.inputs[formPath].Value()
	if path == "" {
		path = m.inputs[formPath].Placeholder
	}
	desc = m.inputs[formDesc].Value()
	tmuxWindowName = m.inputs[formTmux].Value()
	postJumpScript = m.inputs[formScript].Value()

	return strings.TrimSpace(alias),
		strings.TrimSpace(path),
		strings.TrimSpace(desc),
		strings.TrimSpace(tmuxWindowName),
		strings.TrimSpace(postJumpScript)
}

// IsCompleted returns true if the form was completed successfully.
func (m BookmarkFormModel) IsCompleted() bool {
	return m.completed && !m.cancelled
}

// IsCancelled returns true if the form was cancelled.
func (m BookmarkFormModel) IsCancelled() bool {
	return m.cancelled
}
