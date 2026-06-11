package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"bookmark/internal/bookmark"
	"bookmark/internal/config"
	"bookmark/internal/domain"
	"bookmark/internal/ui"
)

type addOptions struct {
	configPath string
}

/*
newAddCmd creates the add command for interactively adding bookmarks.

The add command provides an interactive form to create a new bookmark with all available options:
  - Alias (auto-generated or custom)
  - Path (current directory or custom)
  - Description
  - Tmux window name
  - Execute command
  - Post-jump script
  - File to open

Examples:

	# Interactive add with form
	bookmark add

	# Add with config override
	bookmark add -c ~/.config/bookmark/custom.toml
*/
func newAddCmd() *cobra.Command {
	opts := &addOptions{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Interactively add a new bookmark",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddInteractive(cmd, opts)
		},
	}
	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "config file path")
	return cmd
}

func runAddInteractive(cmd *cobra.Command, opts *addOptions) error {
	// Load config
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	manager := config.NewManager(cwd)
	var cfg domain.Config
	if opts.configPath != "" {
		cfg, err = manager.LoadWithOverride(opts.configPath)
	} else {
		cfg, err = manager.Load()
	}
	if err != nil {
		cfg = domain.DefaultConfig()
	}

	// Create bookmark manager
	bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)

	// Generate default alias
	defaultAlias := bookmark.GenerateAlias(cwd, cfg.AutoAliasSeparator, cfg.AutoAliasLowercase, cfg.DefaultAliasPartLength)

	// Run interactive form
	m := newAddModel(cwd, defaultAlias, cfg, bmManager)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Check if user cancelled
	if finalModel.(addModel).cancelled {
		cmd.Println("Cancelled")
		return nil
	}

	// Get the bookmark from the model
	bm := finalModel.(addModel).toBookmark()

	// Check if bookmark exists
	exists, err := bmManager.Exists(bm.Alias)
	if err != nil {
		return err
	}

	if exists {
		cmd.Printf("⚠️  Bookmark '%s' already exists and will be updated\n", bm.Alias)
	}

	// Save bookmark
	if err := bmManager.Add(bm); err != nil {
		return err
	}

	if exists {
		cmd.Printf("✓ Updated bookmark '%s' → %s\n", bm.Alias, bm.Path)
	} else {
		cmd.Printf("✓ Created bookmark '%s' → %s\n", bm.Alias, bm.Path)
	}

	return nil
}

type addModel struct {
	inputs    []textinput.Model
	focused   int
	cfg       domain.Config
	bmManager *bookmark.Manager
	cancelled bool
	cwd       string
}

const (
	inputAlias = iota
	inputPath
	inputDescription
	inputTmuxName
	inputExecute
	inputPostJump
	inputFile
)

func newAddModel(cwd, defaultAlias string, cfg domain.Config, bmManager *bookmark.Manager) addModel {
	m := addModel{
		inputs:    make([]textinput.Model, 7),
		cfg:       cfg,
		bmManager: bmManager,
		cwd:       cwd,
	}

	// Alias input
	m.inputs[inputAlias] = textinput.New()
	m.inputs[inputAlias].Placeholder = defaultAlias
	m.inputs[inputAlias].Focus()
	m.inputs[inputAlias].Prompt = "Alias:        "
	m.inputs[inputAlias].CharLimit = 50

	// Path input
	m.inputs[inputPath] = textinput.New()
	m.inputs[inputPath].Placeholder = cwd
	m.inputs[inputPath].Prompt = "Path:         "
	m.inputs[inputPath].CharLimit = 500

	// Description input
	m.inputs[inputDescription] = textinput.New()
	m.inputs[inputDescription].Placeholder = "Optional description"
	m.inputs[inputDescription].Prompt = "Description:  "
	m.inputs[inputDescription].CharLimit = 200

	// Tmux window name input
	m.inputs[inputTmuxName] = textinput.New()
	m.inputs[inputTmuxName].Placeholder = "Optional tmux window name"
	m.inputs[inputTmuxName].Prompt = "Tmux Window:  "
	m.inputs[inputTmuxName].CharLimit = 50

	// Execute command input
	m.inputs[inputExecute] = textinput.New()
	m.inputs[inputExecute].Placeholder = "Optional command to execute"
	m.inputs[inputExecute].Prompt = "Execute:      "
	m.inputs[inputExecute].CharLimit = 500

	// Post-jump script input
	m.inputs[inputPostJump] = textinput.New()
	m.inputs[inputPostJump].Placeholder = "Optional post-jump script"
	m.inputs[inputPostJump].Prompt = "Post-jump:    "
	m.inputs[inputPostJump].CharLimit = 500

	// File input
	m.inputs[inputFile] = textinput.New()
	m.inputs[inputFile].Placeholder = "Optional file to open"
	m.inputs[inputFile].Prompt = "File:         "
	m.inputs[inputFile].CharLimit = 500

	return m
}

func (m addModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			// Move to next input or submit
			if m.focused < len(m.inputs)-1 {
				m.inputs[m.focused].Blur()
				m.focused++
				return m, m.inputs[m.focused].Focus()
			}
			// Submit on last field
			return m, tea.Quit

		case "shift+tab", "up":
			// Move to previous input
			if m.focused > 0 {
				m.inputs[m.focused].Blur()
				m.focused--
				return m, m.inputs[m.focused].Focus()
			}

		case "tab", "down":
			// Move to next input
			if m.focused < len(m.inputs)-1 {
				m.inputs[m.focused].Blur()
				m.focused++
				return m, m.inputs[m.focused].Focus()
			}
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

func (m addModel) View() string {
	theme := ui.ThemeFromConfig(m.cfg)

	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Headings)).
		Bold(true).
		MarginBottom(1)

	b.WriteString(titleStyle.Render("Add New Bookmark"))
	b.WriteString("\n\n")

	for i, input := range m.inputs {
		b.WriteString(input.View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted)).
		MarginTop(2)

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓ or tab/shift+tab: navigate • enter: next/submit • esc: cancel"))

	return b.String()
}

func (m addModel) toBookmark() domain.Bookmark {
	alias := m.inputs[inputAlias].Value()
	if alias == "" {
		alias = m.inputs[inputAlias].Placeholder
	}

	path := m.inputs[inputPath].Value()
	if path == "" {
		path = m.inputs[inputPath].Placeholder
	}

	return domain.Bookmark{
		Alias:          alias,
		Path:           path,
		Description:    m.inputs[inputDescription].Value(),
		TmuxWindowName: m.inputs[inputTmuxName].Value(),
		Execute:        m.inputs[inputExecute].Value(),
		PostJumpScript: m.inputs[inputPostJump].Value(),
		File:           m.inputs[inputFile].Value(),
	}
}
