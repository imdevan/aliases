package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"bookmark/internal/adapters/tty"
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
		cmd.Println(ui.ExitMessage(theme, "Cancelled", true))
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
