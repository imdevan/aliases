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
	configPath  string
	tmux        bool
	tmuxName    string
	description string
	yes         bool
	file        string
	execute     string
	source      string
}

func newAddCmd() *cobra.Command {
	opts := &addOptions{}
	cmd := &cobra.Command{
		Use:   "add [alias]",
		Short: "Add a new bookmark",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddInteractive(cmd, args, opts)
		},
	}
	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "config file path")
	cmd.Flags().BoolVarP(&opts.tmux, "tmux", "t", false, "set tmux window name (same as alias)")
	cmd.Flags().StringVarP(&opts.tmuxName, "tmux-name", "T", "", "custom tmux window name")
	cmd.Flags().StringVarP(&opts.description, "description", "d", "", "bookmark description")
	cmd.Flags().BoolVarP(&opts.yes, "yes", "y", false, "skip form, save directly")
	cmd.Flags().StringVarP(&opts.file, "file", "f", "", "file to open in editor after navigation")
	cmd.Flags().StringVarP(&opts.execute, "execute", "x", "", "command to execute after navigation")
	cmd.Flags().StringVarP(&opts.source, "source", "s", "", "path to bookmark (instead of current directory)")
	return cmd
}

func runAddInteractive(cmd *cobra.Command, args []string, opts *addOptions) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	cfg := config.Load(cwd, opts.configPath)

	targetPath := cwd
	if opts.source != "" {
		targetPath = opts.source
	}

	alias := generateAlias(args, targetPath, cfg)
	bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)

	// -y: skip form, save directly
	if opts.yes {
		bm := buildAddBookmark(alias, targetPath, opts)
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
		printSuccess(cfg, action, alias, targetPath)
		return nil
	}

	// Open form prefilled with args/flags values
	tmuxName := opts.tmuxName
	if opts.tmux && tmuxName == "" {
		tmuxName = alias
	}
	prefill := domain.Bookmark{
		Alias:          alias,
		Path:           targetPath,
		Description:    opts.description,
		File:           opts.file,
		TmuxWindowName: tmuxName,
	}

	theme := ui.ThemeFromConfig(cfg)
	m := ui.NewBookmarkFormModelEdit(theme, prefill).WithTitle("Add Bookmark")

	progOpts := tty.GetProgramOptions(tea.WithoutSignalHandler())
	p := tea.NewProgram(m, progOpts...)
	result, err := p.Run()
	if err != nil {
		return err
	}

	fm, ok := result.(ui.BookmarkFormModel)
	if !ok || !fm.IsCompleted() {
		cmd.Println(ui.CanceledMessage(theme, "Add"))
		return nil
	}

	fAlias, fPath, fDesc, fFile, fTmux, fScript := fm.Values()
	bm := domain.Bookmark{
		Alias:          fAlias,
		Path:           fPath,
		Description:    fDesc,
		File:           fFile,
		TmuxWindowName: fTmux,
		PostJumpScript: fScript,
		Execute:        opts.execute,
	}

	exists, err := bmManager.Exists(bm.Alias)
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
	printSuccess(cfg, action, bm.Alias, bm.Path)
	return nil
}

func buildAddBookmark(alias, path string, opts *addOptions) domain.Bookmark {
	bm := domain.Bookmark{
		Alias:       alias,
		Path:        path,
		Description: opts.description,
		File:        opts.file,
		Execute:     opts.execute,
	}
	if opts.tmux {
		bm.TmuxWindowName = alias
	}
	if opts.tmuxName != "" {
		bm.TmuxWindowName = opts.tmuxName
	}
	return bm
}
