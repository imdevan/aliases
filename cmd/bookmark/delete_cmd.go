package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"bookmark/internal/bookmark"
	"bookmark/internal/config"
	"bookmark/internal/flags"
	"bookmark/internal/ui"
)

/*
newDeleteCmd creates the delete command for removing bookmarks.

The delete command removes a bookmark by its alias.
By default, it will prompt for confirmation before deleting.

Flags:
  - --force/-f: Skip confirmation prompt and delete immediately

Examples:

	# Delete with confirmation
	bookmark delete myproject

	# Force delete without confirmation
	bookmark delete myproject --force
*/
func newDeleteCmd() *cobra.Command {
	var configPath string
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <alias>",
		Short: "Delete a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]

			cwd, _ := os.Getwd()
			cfg := config.Load(cwd, configPath)

			bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)

			// Check if bookmark exists
			bm, err := bmManager.Get(alias)
			if err == bookmark.ErrBookmarkNotFound {
				return fmt.Errorf("bookmark '%s' not found", alias)
			}
			if err != nil {
				return err
			}

			// Confirm deletion unless --force or confirm_delete is false
			if !force && cfg.ConfirmDelete {
				theme := ui.ThemeFromConfig(cfg)
				confirmModel := ui.NewConfirmationModel(
					"Delete Bookmark",
					fmt.Sprintf("Delete bookmark '%s → %s'?", bm.Alias, bm.Path),
					theme,
				).WithTitleColor(theme.Error)

				p := tea.NewProgram(confirmModel, tea.WithoutSignalHandler())
				result, err := p.Run()
				if err != nil {
					return err
				}

				if confirmResult, ok := result.(ui.ConfirmationModel); ok {
					if !confirmResult.ChoiceValue() {
						cmd.Println(ui.CanceledMessage(theme, "Delete"))
						return nil
					}
				}
			}

			if err := bmManager.Delete(alias); err != nil {
				return err
			}

			printSuccess(cfg, "deleted", alias, "")
			return nil
		},
	}

	flags.Set(cmd, &configPath, "config", "c", "config file path", "config")
	flags.Set(cmd, &force, "force", "f", "skip confirmation", "")

	return cmd
}
