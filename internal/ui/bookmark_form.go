package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bm "github.com/aliases/internal/bookmark"
	"github.com/aliases/internal/domain"
)

const (
	formAlias = iota
	formPath
	formDesc
	formTmux
	formFile
	formScript
)

var fieldTitles = []string{
	"Alias",
	"Path",
	"Description (optional)",
	"Tmux Window (optional)",
	"File (optional)",
	"Post-jump script (optional)",
}

var fieldDescs = []string{
	"Short name for the bookmark",
	"Directory path to bookmark",
	"",
	"Tmux window name to create/switch to",
	"File to open after navigation",
	"Script/command to run after jumping",
}

// BookmarkFormModel is a form for creating a new bookmark.
type BookmarkFormModel struct {
	inputs           []textinput.Model
	focused          int
	theme            Theme
	responsive       *ResponsiveManager
	completed        bool
	cancelled        bool
	title            string
	validationErrors map[int]string
}

// WithTitle sets a custom title for the form model.
func (m BookmarkFormModel) WithTitle(title string) BookmarkFormModel {
	m.title = title
	return m
}

// NewBookmarkFormModel creates a new bookmark form with optional default values.
func NewBookmarkFormModel(theme Theme, defaultAlias, defaultPath string) BookmarkFormModel {
	inputs := make([]textinput.Model, 6)

	inputs[formAlias] = textinput.New()
	inputs[formAlias].Placeholder = defaultAlias
	inputs[formAlias].Focus()
	inputs[formAlias].Prompt = ""

	inputs[formPath] = textinput.New()
	inputs[formPath].Placeholder = defaultPath
	inputs[formPath].Prompt = ""

	inputs[formDesc] = textinput.New()
	inputs[formDesc].Placeholder = "Optional description"
	inputs[formDesc].Prompt = ""

	inputs[formTmux] = textinput.New()
	inputs[formTmux].Placeholder = "Optional tmux window name"
	inputs[formTmux].Prompt = ""

	inputs[formFile] = textinput.New()
	inputs[formFile].Placeholder = "Optional file to open"
	inputs[formFile].Prompt = ""

	inputs[formScript] = textinput.New()
	inputs[formScript].Placeholder = "Optional post-jump script"
	inputs[formScript].Prompt = ""

	return BookmarkFormModel{
		inputs:           inputs,
		focused:          0,
		theme:            theme,
		responsive:       NewResponsiveManager(80),
		title:            "Add Bookmark",
		validationErrors: make(map[int]string),
	}
}

// NewBookmarkFormModelEdit creates a new bookmark form prefilled with the values of an existing bookmark.
func NewBookmarkFormModelEdit(theme Theme, bm domain.Bookmark) BookmarkFormModel {
	inputs := make([]textinput.Model, 6)

	inputs[formAlias] = textinput.New()
	inputs[formAlias].Placeholder = bm.Alias
	inputs[formAlias].SetValue(bm.Alias)
	inputs[formAlias].Focus()
	inputs[formAlias].Prompt = ""

	inputs[formPath] = textinput.New()
	inputs[formPath].Placeholder = bm.Path
	inputs[formPath].SetValue(bm.Path)
	inputs[formPath].Prompt = ""

	inputs[formDesc] = textinput.New()
	inputs[formDesc].Placeholder = "Optional description"
	inputs[formDesc].SetValue(bm.Description)
	inputs[formDesc].Prompt = ""

	inputs[formTmux] = textinput.New()
	inputs[formTmux].Placeholder = "Optional tmux window name"
	inputs[formTmux].SetValue(bm.TmuxWindowName)
	inputs[formTmux].Prompt = ""

	inputs[formFile] = textinput.New()
	inputs[formFile].Placeholder = "Optional file to open"
	inputs[formFile].SetValue(bm.File)
	inputs[formFile].Prompt = ""

	inputs[formScript] = textinput.New()
	inputs[formScript].Placeholder = "Optional post-jump script"
	inputs[formScript].SetValue(bm.PostJumpScript)
	inputs[formScript].Prompt = ""

	return BookmarkFormModel{
		inputs:           inputs,
		focused:          0,
		theme:            theme,
		responsive:       NewResponsiveManager(80),
		title:            "Edit Bookmark",
		validationErrors: make(map[int]string),
	}
}

// Init initializes the form.
func (m BookmarkFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *BookmarkFormModel) validateField(index int) bool {
	delete(m.validationErrors, index)

	val := m.inputs[index].Value()
	if val == "" {
		val = m.inputs[index].Placeholder
	}
	val = strings.TrimSpace(val)

	switch index {
	case formAlias:
		if val == "" {
			m.validationErrors[formAlias] = "* Alias cannot be empty"
			return false
		}
		if !bm.IsValidAlias(val) {
			m.validationErrors[formAlias] = "* Alias must contain only alphanumeric characters, hyphens, and underscores"
			return false
		}
		if bm.IsReservedKeyword(val) {
			m.validationErrors[formAlias] = "* Alias cannot be a shell reserved keyword"
			return false
		}
	case formPath:
		if val == "" {
			m.validationErrors[formPath] = "* Path cannot be empty"
			return false
		}
	}
	return true
}

func (m *BookmarkFormModel) validateAll() bool {
	valid := true
	for i := range m.inputs {
		if !m.validateField(i) {
			valid = false
		}
	}
	return valid
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
			// Validate current field first
			if !m.validateField(m.focused) {
				return m, nil
			}
			if m.focused == len(m.inputs)-1 {
				if m.validateAll() {
					m.completed = true
					return m, tea.Quit
				}
				// Focus first invalid field
				for i := range m.inputs {
					if _, exists := m.validationErrors[i]; exists {
						m.inputs[m.focused].Blur()
						m.focused = i
						m.inputs[m.focused].Focus()
						break
					}
				}
				return m, nil
			}
			m.inputs[m.focused].Blur()
			m.focused++
			m.inputs[m.focused].Focus()
			return m, nil

		case "alt+enter":
			if m.validateAll() {
				m.completed = true
				return m, tea.Quit
			}
			// Focus first invalid field
			for i := range m.inputs {
				if _, exists := m.validationErrors[i]; exists {
					m.inputs[m.focused].Blur()
					m.focused = i
					m.inputs[m.focused].Focus()
					break
				}
			}
			return m, nil

		case "shift+tab", "up":
			if m.focused > 0 {
				m.validateField(m.focused)
				m.inputs[m.focused].Blur()
				m.focused--
				m.inputs[m.focused].Focus()
			}
			return m, nil

		case "tab", "down":
			if m.focused < len(m.inputs)-1 {
				m.validateField(m.focused)
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

	titleText := m.title
	if titleText == "" {
		titleText = "Create New Bookmark"
	}
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Headings).
		Render(titleText)

	help := lipgloss.NewStyle().
		Foreground(m.theme.Muted).
		Render("↑/↓ or tab/shift+tab: navigate • enter: next • alt+enter: submit • esc: cancel")

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

		// Validation Error (using m.theme.Error)
		if errMsg, exists := m.validationErrors[i]; exists && errMsg != "" {
			errStyle := lipgloss.NewStyle().
				Foreground(m.theme.Error).
				Bold(true)
			itemLines = append(itemLines, errStyle.Render(errMsg))
		}

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

		if _, exists := m.validationErrors[i]; exists {
			inputStyle = inputStyle.BorderForeground(m.theme.Error)
		} else if i == m.focused {
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
func (m BookmarkFormModel) Values() (alias, path, desc, file, tmuxWindowName, postJumpScript string) {
	alias = m.inputs[formAlias].Value()
	if alias == "" {
		alias = m.inputs[formAlias].Placeholder
	}
	path = m.inputs[formPath].Value()
	if path == "" {
		path = m.inputs[formPath].Placeholder
	}
	desc = m.inputs[formDesc].Value()
	file = m.inputs[formFile].Value()
	tmuxWindowName = m.inputs[formTmux].Value()
	postJumpScript = m.inputs[formScript].Value()

	return strings.TrimSpace(alias),
		strings.TrimSpace(path),
		strings.TrimSpace(desc),
		strings.TrimSpace(file),
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
