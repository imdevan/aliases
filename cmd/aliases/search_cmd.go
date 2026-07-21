package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/aliases/internal/config"
	"github.com/aliases/internal/domain"
	"github.com/aliases/internal/flags"
	"github.com/aliases/internal/index"
	"github.com/aliases/internal/ui"
)

func newSearchCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search aliases in the index",
		Long:  "Search the SQLite alias index by name, value, or description.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				cwd = "."
			}

			cfg := config.Load(cwd, configPath)

			store, err := index.NewStore()
			if err != nil {
				return fmt.Errorf("failed to open index: %w", err)
			}
			defer store.Close()

			// Ensure index is fresh before searching.
			indexer := index.NewIndexer(store, cfg)
			if indexer.NeedsRefresh() {
				indexer.Refresh()
			}

			if len(args) == 0 {
				return listAllFromIndex(cmd, store, cfg)
			}

			query := args[0]
			return searchIndex(cmd, store, cfg, query)
		},
	}

	flags.Set(cmd, &configPath, "config", "c", "config file path")

	return cmd
}

func searchIndex(cmd *cobra.Command, store *index.Store, cfg domain.Config, query string) error {
	aliases, err := store.Search(query)
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		cmd.Println("No matching aliases found")
		return nil
	}

	return printGroupedAliases(cmd, aliases, cfg)
}

func listAllFromIndex(cmd *cobra.Command, store *index.Store, cfg domain.Config) error {
	aliases, err := store.All()
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		cmd.Println("Index is empty. Run 'aliases search' after adding some aliases.")
		return nil
	}

	return printGroupedAliases(cmd, aliases, cfg)
}

func printGroupedAliases(cmd *cobra.Command, aliases []domain.Alias, cfg domain.Config) error {
	// Group by source file.
	groups := make(map[string][]domain.Alias)
	var order []string
	for _, a := range aliases {
		src := a.SourceFile
		if src == "" {
			src = "(unknown)"
		}
		if _, exists := groups[src]; !exists {
			order = append(order, src)
		}
		groups[src] = append(groups[src], a)
	}

	theme := ui.ThemeFromConfig(cfg)
	mutedStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	first := true
	for _, src := range order {
		if !first {
			cmd.Println()
		}
		first = false

		// Print source header.
		cmd.Println(mutedStyle.Render(fmt.Sprintf("─── %s", src)))

		group := groups[src]

		// Find max name length for alignment.
		maxName := 0
		for _, al := range group {
			nameLen := len(al.Name)
			if al.Description != "" {
				nameLen += len(" • " + al.Description)
			}
			if nameLen > maxName {
				maxName = nameLen
			}
		}

		for _, al := range group {
			name := al.Name
			descExtra := 0
			if al.Description != "" {
				descStr := " • " + al.Description
				descExtra = len(descStr)
				name += mutedStyle.Render(descStr)
			}
			plainLen := len(al.Name) + descExtra
			padding := strings.Repeat(" ", maxName-plainLen+2)
			cmd.Printf("  %s%s%s\n", name, padding, al.Value)
		}
	}

	return nil
}
