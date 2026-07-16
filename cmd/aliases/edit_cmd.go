package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/aliases/internal/adapters/tty"
	"github.com/aliases/internal/bookmark"
	"github.com/aliases/internal/config"
	"github.com/aliases/internal/domain"
	"github.com/aliases/internal/flags"
	"github.com/aliases/internal/ui"
)

type editOptions struct {
	configPath string
}

// @docs-command:
//
//	name: edit
//	description:
//		The edit command opens a bookmark for editing or opens the entire bookmarks file in the editor.
//	example:
//		```bash
//		# Open all bookmarks in editor
//		bookmark edit
//
//		# Open specific bookmark in form
//		bookmark edit my-alias
//		```
func newEditCmd() *cobra.Command {
	opts := &editOptions{}
	cmd := &cobra.Command{
		Use:   "edit [alias]",
		Short: "Edit a bookmark or open bookmarks file in editor",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()
			cfg := config.Load(cwd, opts.configPath)
			return runEditCommand(cmd, args, opts, cfg)
		},
	}

	flags.Set(cmd, &opts.configPath, "config", "c", "config file path")

	return cmd
}

func runEditCommand(cmd *cobra.Command, args []string, opts *editOptions, cfg domain.Config) error {
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
