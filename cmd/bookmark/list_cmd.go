package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bookmark/internal/bookmark"
	"bookmark/internal/config"
	"bookmark/internal/flags"
)

/*
newListCmd creates the list command for displaying all bookmarks.

The list command shows all bookmarks in a formatted table with:
  - Alias: The bookmark name
  - Path: The directory path
  - Description: Optional bookmark description

The output is formatted with proper alignment for easy reading.

Examples:

	# List all bookmarks
	bookmark list

	# Use with custom config
	bookmark list -c ~/.config/bookmark/custom.toml
*/
func newListCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all bookmarks",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := cmd.Flags().GetString("cwd")
			if err != nil {
				cwd = "."
			}

			cfg := config.Load(cwd, configPath)

			bmManager := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
			bookmarks, err := bmManager.Load()
			if err != nil {
				return err
			}

			if len(bookmarks) == 0 {
				cmd.Println("No bookmarks found")
				return nil
			}

			// Find max alias length for alignment (including description)
			maxAlias := 0
			for _, bm := range bookmarks {
				aliasLen := len(bm.Alias)
				if bm.Description != "" {
					aliasLen += len(" # " + bm.Description)
				}
				if aliasLen > maxAlias {
					maxAlias = aliasLen
				}
			}

			for _, bm := range bookmarks {
				alias := bm.Alias
				if bm.Description != "" {
					alias += " # " + bm.Description
				}
				line := fmt.Sprintf("%-*s  %s", maxAlias, alias, bm.Path)
				cmd.Println(line)
			}

			return nil
		},
	}

	flags.Set(cmd, &configPath, "config", "c", "config file path")

	return cmd
}
