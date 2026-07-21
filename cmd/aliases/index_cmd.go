package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aliases/internal/config"
	"github.com/aliases/internal/flags"
	"github.com/aliases/internal/index"
)

func newIndexCmd() *cobra.Command {
	var configPath string
	var bg bool

	cmd := &cobra.Command{
		Use:   "index",
		Short: "Refresh the SQLite alias index",
		Long:  "Scan alias files and update the SQLite index for search and fast querying.",
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

			indexer := index.NewIndexer(store, cfg)

			if bg {
				if indexer.NeedsRefresh() {
					indexer.Refresh()
				}
				return nil
			}

			res := indexer.Refresh()
			if len(res.Errors) > 0 {
				for _, e := range res.Errors {
					cmd.Printf("warning: %v\n", e)
				}
			}
			cmd.Printf("Indexed %d aliases from %d files (skipped %d unchanged) in %v\n",
				res.AliasesStored, res.FilesScanned, res.FilesSkipped, res.Duration)
			return nil
		},
	}

	flags.Set(cmd, &configPath, "config", "c", "config file path")
	cmd.Flags().BoolVar(&bg, "bg", false, "run indexer silently in background mode")

	return cmd
}
