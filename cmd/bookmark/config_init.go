package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"bookmark/internal/adapters/editor"
	"bookmark/internal/config"
	"bookmark/internal/domain"
	"bookmark/internal/flags"
	"bookmark/internal/utils"
)

type configInitOptions struct {
	force        bool
	openInEditor bool
}

// @docs-command:
//
// name: config init
// description:
//
//	The config init command creates a new configuration file with default values
//	at the standard XDG config location ($XDG_CONFIG_HOME/bookmark/config.toml).
//
// example:
//
//	```bash
//	# Generate default config
//	bookmark config init
//
//	# Overwrite existing config
//	bookmark config init --force
//
//	# Generate and open in editor
//	bookmark config init --editor
//	```
//
// note:
//
//	The generated config file includes commented examples for all available options.
func newConfigInitCmd() *cobra.Command {
	opts := &configInitOptions{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a default config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigInit(cmd, opts)
		},
	}
	flags.Set(cmd, &opts.force, "force", "f", "overwrite existing config")
	flags.Set(cmd, &opts.openInEditor, "editor", "e", "open config in editor after creation")
	return cmd
}

func runConfigInit(cmd *cobra.Command, opts *configInitOptions) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	exists, err := config.NewManager(cwd).Exists()
	if err != nil {
		return err
	}
	if exists && !opts.force {
		return fmt.Errorf("config already exists at %s (use --force to overwrite)", utils.ConfigPathGlobal())
	}
	cfg := domain.DefaultConfig()
	path := utils.ConfigPathGlobal()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	content := renderConfigTemplate(cfg)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	if opts.openInEditor {
		editorAdapter := editor.New(cfg.Editor)
		if err := editorAdapter.Open(path); err != nil {
			return err
		}
	}
	cmd.Printf("Wrote config to %s\n", utils.ConfigPathGlobal())
	return nil
}

func renderConfigTemplate(cfg domain.Config) string {
	var builder strings.Builder
	builder.WriteString("# Bookmark Manager Configuration\n\n")
	builder.WriteString("# General\n")
	builder.WriteString(fmt.Sprintf("# editor = %q\n", cfg.Editor))
	builder.WriteString("\n# CLI behavior\n")
	builder.WriteString(fmt.Sprintf("# interactive_default = %t\n", cfg.InteractiveDefault))
	builder.WriteString("# plain_text = false  # When true: output plain text confirmation messages and errors instead of styled output\n")
	builder.WriteString(fmt.Sprintf("# confirm_delete = %t  # When false: delete bookmarks without confirmation prompt\n", cfg.ConfirmDelete))

	builder.WriteString("\n# Bookmark settings\n")
	builder.WriteString(fmt.Sprintf("# bookmark_location = %q  # Directory for bookmark files\n", cfg.BookmarkLocation))
	builder.WriteString(fmt.Sprintf("# navigation_tool = %q  # Options: cd, z, zoxide, none\n", cfg.NavigationTool))
	builder.WriteString(fmt.Sprintf("# shell = %q  # Options: bash, zsh, fish, nu\n", cfg.Shell))

	builder.WriteString("\n# Auto-alias generation\n")
	builder.WriteString(fmt.Sprintf("# auto_alias_separator = %q  # Character between first letters (empty = \"mcp\", \"-\" = \"m-c-p\")\n", cfg.AutoAliasSeparator))
	builder.WriteString(fmt.Sprintf("# auto_alias_lowercase = %t\n", cfg.AutoAliasLowercase))
	builder.WriteString(fmt.Sprintf("# default_alias_part_length = %d  # Number of characters to take from each part (1 = first letter only)\n", cfg.DefaultAliasPartLength))

	builder.WriteString("\n# Function alias wrapper\n")
	builder.WriteString("# When enabled, adds a shell function that auto-sources bookmarks after running bookmark commands\n")
	builder.WriteString("# Options: \"true\" (use default name \"bookmark\"), \"custom_name\" (use custom function name), \"false\" (disabled)\n")
	builder.WriteString(fmt.Sprintf("# function_alias = %q\n", cfg.FunctionAlias))

	builder.WriteString("\n# Interactive alias wrapper\n")
	builder.WriteString("# When enabled, adds a shell function for interactive bookmark navigation (bookmark -i)\n")
	builder.WriteString("# The function displays the TUI and executes the selected bookmark command\n")
	builder.WriteString("# Options: \"bm\" (default), \"custom_name\" (use custom function name), \"false\" (disabled)\n")
	builder.WriteString(fmt.Sprintf("# interactive_alias = %q\n", cfg.InteractiveAlias))

	builder.WriteString("\n# Display\n")
	builder.WriteString(fmt.Sprintf("# home_icon = %q  # Icon to represent home directory in list view (can be nerd font icon like \"\")\n", cfg.HomeIcon))
	builder.WriteString(fmt.Sprintf("# default_sort_by = %q  # Options: newest, oldest, a-z, z-a\n", cfg.DefaultSortBy))

	builder.WriteString("\n# UI\n")
	builder.WriteString("# list_spacing options: compact (title only), tight (title + description, no margin), space (default, with spacing)\n")
	builder.WriteString(fmt.Sprintf("# list_spacing = %q\n", cfg.ListSpacing))
	builder.WriteString("\n# Colors\n")
	builder.WriteString("# Colors support named, numeric, or hex values (ex: 7, 13, \"#ff8800\").\n")
	builder.WriteString(fmt.Sprintf("# headings = %q\n", cfg.Headings))
	builder.WriteString(fmt.Sprintf("# primary = %q\n", cfg.Primary))
	builder.WriteString(fmt.Sprintf("# secondary = %q\n", cfg.Secondary))
	builder.WriteString(fmt.Sprintf("# text = %q\n", cfg.Text))
	builder.WriteString(fmt.Sprintf("# text_highlight = %q\n", cfg.TextHighlight))
	builder.WriteString(fmt.Sprintf("# description_highlight = %q\n", cfg.DescriptionHighlight))
	builder.WriteString(fmt.Sprintf("# tags = %q\n", cfg.Tags))
	builder.WriteString(fmt.Sprintf("# flags = %q\n", cfg.Flags))
	builder.WriteString(fmt.Sprintf("# muted = %q\n", cfg.Muted))
	builder.WriteString(fmt.Sprintf("# border = %q\n", cfg.Border))
	builder.WriteString(fmt.Sprintf("# error = %q\n", cfg.Error))
	builder.WriteString(fmt.Sprintf("# success = %q\n", cfg.Success))
	return builder.String()
}
