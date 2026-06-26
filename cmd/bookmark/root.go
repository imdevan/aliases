package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"bookmark/internal/adapters/editor"
	"bookmark/internal/adapters/icon"
	"bookmark/internal/adapters/tty"
	"bookmark/internal/bookmark"
	"bookmark/internal/config"
	"bookmark/internal/domain"
	pkg "bookmark/internal/package"
	"bookmark/internal/ui"
)

// Metadata loaded from package.toml at build time
var (
	version = pkg.Version()
	name    = pkg.Name()
	short   = pkg.Short()
)

type rootOptions struct {
	configPath  string
	showVersion bool
	interactive bool
	add         bool
	tmux        bool
	tmuxName    string
	description string
	yes         bool
	file        string
	edit        bool
	execute     string
	source      string
}

var rootCmd = newRootCmd()

// Execute is the CLI entrypoint.
func Execute() error {
	return rootCmd.Execute()
}

/*
newRootCmd creates the root command for the bookmark CLI.

The root command serves multiple purposes:
  - Without arguments: Opens interactive bookmark browser (if configured)
  - With alias argument: Navigates to the bookmarked directory
  - With --interactive/-i: Forces interactive mode
  - With --edit/-e: Opens bookmarks file in editor
  - With --version/-v: Prints version information

When adding a bookmark, you can specify:
  - --description/-d: Add a description to the bookmark
  - --tmux/-t: Set tmux window name to match alias
  - --tmux-name/-T: Set custom tmux window name
  - --file/-f: Specify a file to open after navigation
  - --execute/-x: Run a command after navigation
  - --source/-s: Bookmark a different path than current directory
  - --yes/-y: Skip confirmation prompts

Examples:

	# Add bookmark for current directory
	bookmark myproject

	# Add bookmark with description
	bookmark myproject -d "My awesome project"

	# Navigate to bookmark (in interactive mode)
	bookmark

	# Navigate to specific bookmark
	bookmark myproject

	# Edit bookmarks file
	bookmark -e

	# List all bookmarks
	bookmark list
*/
func newRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:   name + " [alias]",
		Short: short,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.showVersion {
				ver := resolvedVersion()
				cmd.Printf("%s\n", ver)
				return nil
			}

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

			// Interactive mode
			if opts.interactive || (len(args) == 0 && cfg.InteractiveDefault && !opts.tmux && opts.description == "" && !opts.edit && !opts.add) {
				return runInteractive(cmd, opts, cfg)
			}

			// Interactive add form
			if opts.add {
				return runAddForm(cmd, opts, cfg, cwd)
			}

			// Edit mode
			if opts.edit {
				return runEdit(cmd, args, opts, cfg)
			}

			// Add bookmark mode
			return runAddBookmark(cmd, args, opts, cfg, cwd)
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.configPath, "config", "c", "", "config file path")
	cmd.Flags().BoolVarP(&opts.showVersion, "version", "v", false, "print version information")
	cmd.Flags().BoolVarP(&opts.interactive, "interactive", "i", false, "interactive bookmark browser")
	cmd.Flags().BoolVarP(&opts.add, "add", "a", false, "interactive add bookmark form")
	cmd.Flags().BoolVarP(&opts.tmux, "tmux", "t", false, "set tmux window name (same as alias)")
	cmd.Flags().StringVarP(&opts.tmuxName, "tmux-name", "T", "", "custom tmux window name")
	cmd.Flags().StringVarP(&opts.description, "description", "d", "", "bookmark description")
	cmd.Flags().BoolVarP(&opts.yes, "yes", "y", false, "skip confirmation prompts")
	cmd.Flags().StringVarP(&opts.file, "file", "f", "", "file to open in editor after navigation")
	cmd.Flags().BoolVarP(&opts.edit, "edit", "e", false, "open bookmarks file in editor")
	cmd.Flags().StringVarP(&opts.execute, "execute", "x", "", "command to execute after navigation")
	cmd.Flags().StringVarP(&opts.source, "source", "s", "", "path to bookmark (instead of current directory)")

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newCompletionCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newAddCmd())

	return cmd
}

func resolvedVersion() string {
	ver := version
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ver
	}
	if ver == "dev" && strings.TrimSpace(info.Main.Version) != "" && info.Main.Version != "(devel)" {
		ver = info.Main.Version
	}
	return ver
}

func runAddBookmark(cmd *cobra.Command, args []string, opts *rootOptions, cfg domain.Config, cwd string) error {
	bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)

	// Use source path if provided, otherwise use current directory
	targetPath := cwd
	if opts.source != "" {
		targetPath = opts.source
	}

	// Generate or use provided alias
	alias := generateAlias(args, targetPath, cfg)

	// Check if bookmark exists and handle confirmation
	exists, err := bmManager.Exists(alias)
	if err != nil {
		return err
	}

	if exists && !opts.yes && !confirmOverwrite(cmd, bmManager, alias, cfg) {
		theme := ui.ThemeFromConfig(cfg)
		cmd.Println(ui.CanceledMessage(theme, "Overwrite"))
		return nil
	}

	// Create and save bookmark
	bm := buildBookmark(alias, targetPath, opts)
	if err := bmManager.Add(bm); err != nil {
		return err
	}

	action := "created"
	if exists {
		action = "updated"
	}
	printSuccess(cfg, action, alias, targetPath)
	return nil
}

func generateAlias(args []string, cwd string, cfg domain.Config) string {
	if len(args) > 0 {
		return args[0]
	}
	return bookmark.GenerateAlias(cwd, cfg.AutoAliasSeparator, cfg.AutoAliasLowercase, cfg.DefaultAliasPartLength)
}

func confirmOverwrite(cmd *cobra.Command, bmManager *bookmark.Manager, alias string, cfg domain.Config) bool {
	existing, _ := bmManager.Get(alias)
	theme := ui.ThemeFromConfig(cfg)
	confirmModel := ui.NewConfirmationModel(
		"Overwrite Bookmark",
		fmt.Sprintf("Bookmark '%s → %s' already exists. Overwrite?", alias, existing.Path),
		theme,
	)

	p := tea.NewProgram(confirmModel, tea.WithoutSignalHandler())
	result, err := p.Run()
	if err != nil {
		return false
	}

	if confirmResult, ok := result.(ui.ConfirmationModel); ok {
		return confirmResult.ChoiceValue()
	}
	return false
}

func buildBookmark(alias, cwd string, opts *rootOptions) domain.Bookmark {
	bm := domain.Bookmark{
		Alias:       alias,
		Path:        cwd,
		Description: opts.description,
		File:        opts.file,
		Execute:     opts.execute,
	}

	// Handle tmux settings
	if opts.tmux {
		bm.TmuxWindowName = alias
	}
	if opts.tmuxName != "" {
		bm.TmuxWindowName = opts.tmuxName
	}

	return bm
}

func printSuccess(cfg domain.Config, action, alias, path string) {
	theme := ui.ThemeFromConfig(cfg)
	inline := action == "deleted"
	var body string
	if path != "" {
		home, _ := os.UserHomeDir()
		displayPath := path
		switch {
		case cfg.HomeIcon == "" || home == "":
		case path == home:
			displayPath = cfg.HomeIcon
		case strings.HasPrefix(path, home+"/"):
			displayPath = cfg.HomeIcon + strings.TrimPrefix(path, home)
		}
		body = fmt.Sprintf("%s → %s", alias, displayPath)
	} else {
		body = alias
	}
	fmt.Println(ui.SuccessMessage(theme, action, body, inline))
}

func runEdit(cmd *cobra.Command, args []string, opts *rootOptions, cfg domain.Config) error {
	bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)

	// If no alias provided, just open the bookmarks file in editor
	if len(args) == 0 {
		return openEditor(cfg.Editor, cfg.BookmarkFile(), 0)
	}

	alias := args[0]

	// Check if bookmark exists
	exists, err := bmManager.Exists(alias)
	if err != nil {
		return err
	}

	var bm domain.Bookmark
	if exists {
		bm, err = bmManager.Get(alias)
		if err != nil {
			return err
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		bm = domain.Bookmark{
			Alias: alias,
			Path:  cwd,
		}
	}

	theme := ui.ThemeFromConfig(cfg)
	m := ui.NewBookmarkFormModelEdit(theme, bm)
	if !exists {
		m = m.WithTitle(fmt.Sprintf("'%s' Not Found, Add Bookmark", alias))
	}

	progOpts := tty.GetProgramOptions(tea.WithoutSignalHandler())
	p := tea.NewProgram(m, progOpts...)
	result, err := p.Run()
	if err != nil {
		return err
	}

	fm, ok := result.(ui.BookmarkFormModel)
	if !ok || !fm.IsCompleted() {
		fmt.Println(ui.CanceledMessage(theme, "Edit"))
		return nil
	}

	newAlias, newPath, newDesc, newFile, tmuxWindowName, postJumpScript := fm.Values()

	// If the alias changed and we are editing an existing one, delete the old one
	if exists && newAlias != alias {
		if err := bmManager.Delete(alias); err != nil {
			return err
		}
	}

	newBm := domain.Bookmark{
		Alias:          newAlias,
		Path:           newPath,
		Description:    newDesc,
		File:           newFile,
		TmuxWindowName: tmuxWindowName,
		PostJumpScript: postJumpScript,
	}
	if exists {
		newBm.CreatedAt = bm.CreatedAt
	}

	if err := bmManager.Add(newBm); err != nil {
		return err
	}

	action := "created"
	if exists {
		action = "updated"
	}
	printSuccess(cfg, action, newAlias, newPath)
	return nil
}

func openEditor(editorName, filePath string, line int) error {
	if editorName == "" {
		return fmt.Errorf("no editor configured")
	}

	// Use the editor adapter
	editorAdapter := editor.New(editorName)
	if line > 0 {
		return editorAdapter.OpenAtLine(filePath, line)
	}
	return editorAdapter.Open(filePath)
}

func runAddForm(cmd *cobra.Command, opts *rootOptions, cfg domain.Config, cwd string) error {
	bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	defaultAlias := bookmark.GenerateAlias(cwd, cfg.AutoAliasSeparator, cfg.AutoAliasLowercase, cfg.DefaultAliasPartLength)

	theme := ui.ThemeFromConfig(cfg)
	m := ui.NewBookmarkFormModel(theme, defaultAlias, cwd)

	progOpts := tty.GetProgramOptions(tea.WithoutSignalHandler())
	p := tea.NewProgram(m, progOpts...)
	result, err := p.Run()
	if err != nil {
		return err
	}

	fm, ok := result.(ui.BookmarkFormModel)
	if !ok || !fm.IsCompleted() {
		fmt.Println(ui.CanceledMessage(theme, "Add"))
		return nil
	}

	alias, path, desc, file, tmuxWindowName, postJumpScript := fm.Values()
	bm := domain.Bookmark{
		Alias:          alias,
		Path:           path,
		Description:    desc,
		File:           file,
		TmuxWindowName: tmuxWindowName,
		PostJumpScript: postJumpScript,
	}

	exists, err := bmManager.Exists(alias)
	if err != nil {
		return err
	}
	if err := bmManager.Add(bm); err != nil {
		return err
	}
	action := "created"
	if exists {
		action = "updated"
	}
	printSuccess(cfg, action, alias, path)
	return nil
}

func runInteractive(cmd *cobra.Command, opts *rootOptions, cfg domain.Config) error {
	bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	bookmarks, err := bmManager.Load()
	if err != nil {
		return err
	}

	if len(bookmarks) == 0 {
		cmd.Println("No bookmarks found. Add one with: bookmark [alias]")
		return nil
	}

	return runBookmarkListing(bookmarks, cfg, bmManager)
}

func sortBookmarks(bookmarks []domain.Bookmark, sortBy string) {
	switch sortBy {
	case "newest":
		// Sort by CreatedAt descending (newest first)
		for i := 0; i < len(bookmarks)-1; i++ {
			for j := i + 1; j < len(bookmarks); j++ {
				if bookmarks[i].CreatedAt.Before(bookmarks[j].CreatedAt) {
					bookmarks[i], bookmarks[j] = bookmarks[j], bookmarks[i]
				}
			}
		}
	case "oldest", "latest":
		// Sort by CreatedAt ascending (oldest first)
		for i := 0; i < len(bookmarks)-1; i++ {
			for j := i + 1; j < len(bookmarks); j++ {
				if bookmarks[i].CreatedAt.After(bookmarks[j].CreatedAt) {
					bookmarks[i], bookmarks[j] = bookmarks[j], bookmarks[i]
				}
			}
		}
	case "a-z", "A to Z":
		// Sort by Alias ascending (A-Z)
		for i := 0; i < len(bookmarks)-1; i++ {
			for j := i + 1; j < len(bookmarks); j++ {
				if strings.ToLower(bookmarks[i].Alias) > strings.ToLower(bookmarks[j].Alias) {
					bookmarks[i], bookmarks[j] = bookmarks[j], bookmarks[i]
				}
			}
		}
	case "z-a", "Z to A":
		// Sort by Alias descending (Z-A)
		for i := 0; i < len(bookmarks)-1; i++ {
			for j := i + 1; j < len(bookmarks); j++ {
				if strings.ToLower(bookmarks[i].Alias) < strings.ToLower(bookmarks[j].Alias) {
					bookmarks[i], bookmarks[j] = bookmarks[j], bookmarks[i]
				}
			}
		}
	}
}

func runBookmarkListing(bookmarks []domain.Bookmark, cfg domain.Config, bmManager *bookmark.Manager) error {
	// Sort bookmarks based on config
	sortBookmarks(bookmarks, cfg.DefaultSortBy)

	items := make([]list.Item, 0, len(bookmarks))
	for _, bm := range bookmarks {
		items = append(items, bookmarkItem{Bookmark: bm, Config: cfg})
	}

	theme := ui.ThemeFromConfig(cfg)
	delegate := ui.NewListDelegate(theme, ui.ListDelegateOptions{
		Spacing:        cfg.ListSpacing,
		ShowMetadata:   true,
		MetadataIndent: 1, // Align with path start
	})

	listModel := ui.NewListModel(items, delegate, 80, 20, theme)
	listModel.Title = fmt.Sprintf("%s Bookmarks (%d)", icon.Bookmarks.String(), len(items))
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(true)

	model := bookmarkListModel{
		list:       listModel,
		theme:      theme,
		responsive: ui.NewResponsiveManager(80),
		manager:    bmManager,
		config:     cfg,
	}

	model.list.AdditionalShortHelpKeys = model.getShortHelpKeys
	model.list.AdditionalFullHelpKeys = model.allHelpKeys

	// Get program options with TTY redirection when needed
	opts := tty.GetProgramOptions(tea.WithoutSignalHandler())

	p := tea.NewProgram(model, opts...)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run interactive list: %w", err)
	}

	return nil
}

type bookmarkItem struct {
	Bookmark domain.Bookmark
	Config   domain.Config
}

func (b bookmarkItem) Title() string {
	title := b.Bookmark.Alias

	// Add description next to title with bullet separator in muted color
	if b.Bookmark.Description != "" {
		mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(b.Config.Muted))
		title += " " + mutedStyle.Render("• "+b.Bookmark.Description)
	}

	return title
}

func (b bookmarkItem) Description() string {
	desc := b.Bookmark.Path

	// Replace home directory with home icon
	if b.Config.HomeIcon != "" {
		home, err := os.UserHomeDir()
		if err == nil && strings.HasPrefix(desc, home) {
			desc = b.Config.HomeIcon + strings.TrimPrefix(desc, home)
		}
	}

	return desc
}

// Metadata implements ui.ItemWithMetadata interface.
func (b bookmarkItem) Metadata() string {
	var parts []string

	// Tmux window name with icon
	if b.Bookmark.TmuxWindowName != "" {
		parts = append(parts, icon.Tmux.String()+" "+b.Bookmark.TmuxWindowName)
	}

	// File to open with icon
	if b.Bookmark.File != "" {
		editorIcon := icon.GetEditorIcon(b.Config.Editor)
		if editorIcon != "" {
			parts = append(parts, editorIcon.String()+" "+b.Bookmark.File)
		} else {
			parts = append(parts, icon.File.String()+" "+b.Bookmark.File)
		}
	}

	// Execute command with icon
	if b.Bookmark.Execute != "" {
		parts = append(parts, icon.Script.String()+" "+b.Bookmark.Execute)
	}
	if b.Bookmark.PostJumpScript != "" {
		parts = append(parts, icon.Script.String()+" "+b.Bookmark.PostJumpScript)
	}

	return strings.Join(parts, " • ")
}

func (b bookmarkItem) FilterValue() string {
	return b.Bookmark.Alias + " " + b.Bookmark.Path + " " + b.Bookmark.Description
}

type bookmarkListModel struct {
	list            list.Model
	theme           ui.Theme
	responsive      *ui.ResponsiveManager
	manager         *bookmark.Manager
	config          domain.Config
	message         string
	confirmMode     bool
	confirmModel    *ui.ConfirmationModel
	addMode         bool
	addModel        *ui.BookmarkFormModel
	editMode        bool
	editModel       *ui.BookmarkFormModel
	editingAlias    string
	pendingAction   string
	pendingItem     bookmarkItem
	pendingBookmark *domain.Bookmark
	screenW         int
	screenH         int
}

func (m *bookmarkListModel) updateTitle() {
	visibleCount := len(m.list.VisibleItems())
	totalCount := len(m.list.Items())

	if m.list.FilterState() == list.Filtering || m.list.FilterState() == list.FilterApplied {
		if visibleCount != totalCount {
			m.list.Title = fmt.Sprintf("%s Bookmarks (%d/%d)", icon.Bookmarks.String(), visibleCount, totalCount)
		} else {
			m.list.Title = fmt.Sprintf("%s Bookmarks (%d)", icon.Bookmarks.String(), totalCount)
		}
	} else {
		m.list.Title = fmt.Sprintf("%s Bookmarks (%d)", icon.Bookmarks.String(), totalCount)
	}
}

func (m bookmarkListModel) allHelpKeys() []key.Binding {
	// Don't show alphabetic keys when filtering to avoid interference
	if m.list.FilterState() == list.Filtering {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "copy cd command"),
			),
		}
	}

	return []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "copy cd command"),
		),
		key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add bookmark"),
		),
		key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit bookmark"),
		),
		key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete bookmark"),
		),
		key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "force delete"),
		),
	}
}

func (m bookmarkListModel) getShortHelpKeys() []key.Binding {
	allKeys := m.allHelpKeys()

	var splitAt int
	switch m.responsive.Breakpoint() {
	case ui.BreakpointXL:
		splitAt = 2
	case ui.BreakpointLG:
		splitAt = 1
	default:
		splitAt = 1
	}

	return allKeys[:splitAt]
}

func (m bookmarkListModel) getFullHelpKeys() []key.Binding {
	allKeys := m.allHelpKeys()

	var splitAt int
	switch m.responsive.Breakpoint() {
	case ui.BreakpointXS:
		splitAt = 1
	case ui.BreakpointSM:
		splitAt = 1
	case ui.BreakpointMD:
		splitAt = 2
	default:
		return []key.Binding{}
	}

	return allKeys[splitAt:]
}

func (m bookmarkListModel) Init() tea.Cmd {
	return nil
}

func (m bookmarkListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editMode && m.editModel != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.screenW = msg.Width
			m.screenH = msg.Height
			m.responsive.SetWidth(msg.Width)
			width, height := m.responsive.GetListDimensions(msg.Width, msg.Height)
			m.list.SetSize(width, height)
			mw, mh := modalDimensions(msg.Width, msg.Height)
			updatedForm, _ := m.editModel.Update(tea.WindowSizeMsg{Width: mw, Height: mh})
			if fm, ok := updatedForm.(ui.BookmarkFormModel); ok {
				m.editModel = &fm
			}
			return m, nil
		}

		var cmd tea.Cmd
		var updatedModel tea.Model
		updatedModel, cmd = m.editModel.Update(msg)
		if updatedForm, ok := updatedModel.(ui.BookmarkFormModel); ok {
			m.editModel = &updatedForm
		}

		if m.editModel.IsCompleted() {
			m.editMode = false
			alias, path, desc, file, tmuxWindowName, postJumpScript := m.editModel.Values()

			// Load old bookmark to preserve CreatedAt
			var oldBm domain.Bookmark
			var loadErr error
			oldBm, loadErr = m.manager.Get(m.editingAlias)

			// If the alias changed, delete the old one
			if m.editingAlias != alias {
				if err := m.manager.Delete(m.editingAlias); err != nil {
					m.message = fmt.Sprintf("✗ Failed to delete old bookmark: %s", err)
					m.editModel = nil
					return m, nil
				}
			}

			bm := domain.Bookmark{
				Alias:          alias,
				Path:           path,
				Description:    desc,
				File:           file,
				TmuxWindowName: tmuxWindowName,
				PostJumpScript: postJumpScript,
			}
			if loadErr == nil {
				bm.CreatedAt = oldBm.CreatedAt
			}

			if err := m.manager.Add(bm); err != nil {
				m.message = fmt.Sprintf("✗ Failed to update: %s", err)
			} else {
				m.reloadList()
			}
			m.editModel = nil
			m.updateTitle()
			return m, nil
		} else if m.editModel.IsCancelled() {
			m.editMode = false
			m.editModel = nil
			return m, nil
		}

		return m, cmd
	}

	if m.addMode && m.addModel != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.screenW = msg.Width
			m.screenH = msg.Height
			m.responsive.SetWidth(msg.Width)
			width, height := m.responsive.GetListDimensions(msg.Width, msg.Height)
			m.list.SetSize(width, height)
			mw, mh := modalDimensions(msg.Width, msg.Height)
			updatedForm, _ := m.addModel.Update(tea.WindowSizeMsg{Width: mw, Height: mh})
			if fm, ok := updatedForm.(ui.BookmarkFormModel); ok {
				m.addModel = &fm
			}
			return m, nil
		}

		var cmd tea.Cmd
		var updatedModel tea.Model
		updatedModel, cmd = m.addModel.Update(msg)
		if updatedForm, ok := updatedModel.(ui.BookmarkFormModel); ok {
			m.addModel = &updatedForm
		}

		if m.addModel.IsCompleted() {
			m.addMode = false
			alias, path, desc, file, tmuxWindowName, postJumpScript := m.addModel.Values()
			bm := domain.Bookmark{
				Alias:          alias,
				Path:           path,
				Description:    desc,
				File:           file,
				TmuxWindowName: tmuxWindowName,
				PostJumpScript: postJumpScript,
			}
			m.addModel = nil

			// Check for existing bookmark — show overwrite confirmation if needed
			exists, _ := m.manager.Exists(alias)
			if exists {
				m.pendingAction = "Overwrite"
				m.pendingBookmark = &bm
				existing, _ := m.manager.Get(alias)
				confirmModel := ui.NewConfirmationModel(
					"Overwrite Bookmark",
					fmt.Sprintf("'%s → %s' exists. Overwrite?", alias, existing.Path),
					m.theme,
				)
				m.confirmModel = &confirmModel
				m.confirmMode = true
				return m, confirmModel.Init()
			}

			if err := m.manager.Add(bm); err != nil {
				m.message = fmt.Sprintf("✗ Failed to add: %s", err)
			} else {
				m.reloadList()
			}
			m.updateTitle()
			return m, nil
		} else if m.addModel.IsCancelled() {
			m.addMode = false
			m.addModel = nil
			return m, nil
		}

		return m, cmd
	}

	if m.confirmMode && m.confirmModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			updated, cmd := m.confirmModel.Update(msg)
			if updatedConfirm, ok := updated.(ui.ConfirmationModel); ok {
				m.confirmModel = &updatedConfirm
				if cmd != nil {
					if _, isQuit := cmd().(tea.QuitMsg); isQuit {
						confirmed := m.confirmModel.ChoiceValue()
						m.confirmMode = false
						if confirmed {
							return m.executeAction()
						} else {
							if m.pendingAction != "Delete" {
								m.message = fmt.Sprintf("%s cancelled", m.pendingAction)
							}
							m.pendingAction = ""
							return m, nil
						}
					}
				}
			}
			return m, cmd
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenW = msg.Width
		m.screenH = msg.Height
		m.responsive.SetWidth(msg.Width)
		width, height := m.responsive.GetListDimensions(msg.Width, msg.Height)
		m.list.SetSize(width, height)
		m.list.AdditionalShortHelpKeys = m.getShortHelpKeys
		return m, nil

	case tea.KeyMsg:
		// When filtering, only allow filter-related keys and enter
		// Block all alphabetic action keys to prevent interference
		if m.list.FilterState() == list.Filtering {
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				return m, tea.Quit
			case "enter":
				if item, ok := m.list.SelectedItem().(bookmarkItem); ok {
					// Execute the full navigation command directly
					fmt.Println(m.manager.BuildNavigationCommand(item.Bookmark))
					return m, tea.Quit
				}
			case "e", "n", "d", "D", "a":
				// Block these keys during filtering - let them pass to filter input
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
			// Pass all other keys to list for filtering
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		// Normal mode - handle action keys
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(bookmarkItem); ok {
				// Execute the full navigation command directly
				fmt.Println(m.manager.BuildNavigationCommand(item.Bookmark))
				return m, tea.Quit
			}
		case "a":
			cwd, err := os.Getwd()
			if err != nil {
				m.message = "✗ Failed to get current directory"
				return m, nil
			}
			defaultAlias := bookmark.GenerateAlias(cwd, m.config.AutoAliasSeparator, m.config.AutoAliasLowercase, m.config.DefaultAliasPartLength)
			formModel := ui.NewBookmarkFormModel(m.theme, defaultAlias, cwd)
			m.addModel = &formModel
			m.addMode = true

			if m.screenW > 0 {
				mw, mh := modalDimensions(m.screenW, m.screenH)
				updatedForm, _ := m.addModel.Update(tea.WindowSizeMsg{Width: mw, Height: mh})
				if fm, ok := updatedForm.(ui.BookmarkFormModel); ok {
					m.addModel = &fm
				}
			}

			return m, formModel.Init()
		case "e":
			if item, ok := m.list.SelectedItem().(bookmarkItem); ok {
				formModel := ui.NewBookmarkFormModelEdit(m.theme, item.Bookmark)
				m.editModel = &formModel
				m.editMode = true
				m.editingAlias = item.Bookmark.Alias

				if m.screenW > 0 {
					mw, mh := modalDimensions(m.screenW, m.screenH)
					updatedForm, _ := m.editModel.Update(tea.WindowSizeMsg{Width: mw, Height: mh})
					if fm, ok := updatedForm.(ui.BookmarkFormModel); ok {
						m.editModel = &fm
					}
				}

				return m, formModel.Init()
			}
		case "d":
			if item, ok := m.list.SelectedItem().(bookmarkItem); ok {
				m.pendingAction = "Delete"
				m.pendingItem = item
				confirmModel := ui.NewConfirmationModel(
					"Delete Bookmark",
					fmt.Sprintf("Delete bookmark '%s'?", item.Bookmark.Alias),
					m.theme,
				).WithTitleColor(m.theme.Error)
				m.confirmModel = &confirmModel
				m.confirmMode = true
				return m, confirmModel.Init()
			}
		case "D":
			if item, ok := m.list.SelectedItem().(bookmarkItem); ok {
				// Force delete without confirmation
				if err := m.manager.Delete(item.Bookmark.Alias); err != nil {
					m.message = fmt.Sprintf("✗ Failed to delete: %s", err)
				} else {
					// Remove from list
					items := m.list.Items()
					filtered := make([]list.Item, 0, len(items))
					for _, listItem := range items {
						if bm, ok := listItem.(bookmarkItem); ok {
							if bm.Bookmark.Alias != item.Bookmark.Alias {
								filtered = append(filtered, bm)
							}
						}
					}
					m.list.SetItems(filtered)
					m.updateTitle()
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.updateTitle()
	return m, cmd
}

func (m bookmarkListModel) View() string {
	listView := m.responsive.FullWidthFrameStyle(m.theme).Render(m.list.View())

	if m.message != "" {
		listView = listView + "\n\n" + m.message
	}

	if m.addMode && m.addModel != nil {
		return ui.Center(listView, m.addModel.View(), m.screenW, m.screenH)
	}

	if m.editMode && m.editModel != nil {
		return ui.Center(listView, m.editModel.View(), m.screenW, m.screenH)
	}

	if m.confirmMode && m.confirmModel != nil {
		return ui.Center(listView, m.confirmModel.View(), m.screenW, m.screenH)
	}

	return listView
}

func (m *bookmarkListModel) reloadList() {
	bookmarks, err := m.manager.Load()
	if err == nil {
		sortBookmarks(bookmarks, m.config.DefaultSortBy)
		items := make([]list.Item, 0, len(bookmarks))
		for _, b := range bookmarks {
			items = append(items, bookmarkItem{Bookmark: b, Config: m.config})
		}
		m.list.SetItems(items)
	}
}

func modalDimensions(screenW, screenH int) (width, height int) {
	const minW, minH = 80, 20
	width = screenW * 75 / 100
	if width < minW {
		width = minW
	}
	height = screenH * 80 / 100
	if height < minH {
		height = minH
	}
	return
}

func (m bookmarkListModel) executeAction() (tea.Model, tea.Cmd) {
	switch m.pendingAction {
	case "Delete":
		if err := m.manager.Delete(m.pendingItem.Bookmark.Alias); err != nil {
			m.message = fmt.Sprintf("✗ Failed to delete: %s", err)
		} else {
			// Remove from list
			items := m.list.Items()
			filtered := make([]list.Item, 0, len(items))
			for _, item := range items {
				if bm, ok := item.(bookmarkItem); ok {
					if bm.Bookmark.Alias != m.pendingItem.Bookmark.Alias {
						filtered = append(filtered, item)
					}
				}
			}
			m.list.SetItems(filtered)
			m.updateTitle()
		}
	case "Overwrite":
		if m.pendingBookmark != nil {
			if err := m.manager.Add(*m.pendingBookmark); err != nil {
				m.message = fmt.Sprintf("✗ Failed to overwrite: %s", err)
			} else {
				m.message = fmt.Sprintf("✓ Overwritten: %s", m.pendingBookmark.Alias)
				m.reloadList()
				m.updateTitle()
			}
			m.pendingBookmark = nil
		}
	}
	m.pendingAction = ""
	return m, nil
}
