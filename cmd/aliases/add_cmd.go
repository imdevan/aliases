package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/aliases/internal/adapters/tty"
	"github.com/aliases/internal/alias"
	"github.com/aliases/internal/config"
	"github.com/aliases/internal/domain"
	"github.com/aliases/internal/flags"
	"github.com/aliases/internal/ui"
)

type addOptions struct {
	configPath string
	yes        bool
}

func newAddCmd() *cobra.Command {
	opts := &addOptions{}
	cmd := &cobra.Command{
		Use:   "add [name] [value] [description]",
		Short: "Add a new alias",
		Args:  cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddInteractive(cmd, args, opts)
		},
	}

	flags.Set(cmd, &opts.configPath, "config", "c", "config file path")
	flags.Set(cmd, &opts.yes, "yes", "y", "skip form, save directly")

	return cmd
}

func runAddInteractive(cmd *cobra.Command, args []string, opts *addOptions) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	cfg := config.Load(cwd, opts.configPath)
	aliasManager := alias.NewManager(cfg.ResolvedAliasFile(), cfg.Shell, cfg.FunctionAlias, cfg.InteractiveAlias, cfg.IndexFolders)

	var name, value, description string
	if len(args) >= 1 {
		name = args[0]
	}
	if len(args) >= 2 {
		value = args[1]
	}
	if len(args) >= 3 {
		description = args[2]
	}

	// -y: skip form, save directly
	if opts.yes {
		if name == "" || value == "" {
			return fmt.Errorf("name and value are required when using --yes flag")
		}
		alItem := domain.Alias{
			Name:        name,
			Value:       value,
			Description: description,
		}
		exists, err := aliasManager.Exists(name)
		if err != nil {
			return err
		}
		if err := aliasManager.Add(alItem); err != nil {
			return err
		}
		action := "created"
		if exists {
			action = "updated"
		}
		printSuccess(cfg, action, name, value)
		return nil
	}

	theme := ui.ThemeFromConfig(cfg)
	var m ui.AliasFormModel
	if name != "" || value != "" || description != "" {
		m = ui.NewAliasFormModelEdit(theme, domain.Alias{
			Name:        name,
			Value:       value,
			Description: description,
		}).WithTitle("Add Alias")
	} else {
		m = ui.NewAliasFormModel(theme, "", "").WithTitle("Add Alias")
	}

	progOpts := tty.GetProgramOptions(tea.WithoutSignalHandler())
	p := tea.NewProgram(m, progOpts...)
	result, err := p.Run()
	if err != nil {
		return err
	}

	fm, ok := result.(ui.AliasFormModel)
	if !ok || !fm.IsCompleted() {
		cmd.Println(ui.CanceledMessage(theme, "Add"))
		return nil
	}

	fName, fValue, fDesc := fm.Values()
	alItem := domain.Alias{
		Name:        fName,
		Value:       fValue,
		Description: fDesc,
	}

	exists, err := aliasManager.Exists(alItem.Name)
	if err != nil {
		return err
	}
	if err := aliasManager.Add(alItem); err != nil {
		return err
	}
	action := "created"
	if exists {
		action = "updated"
	}
	printSuccess(cfg, action, alItem.Name, alItem.Value)
	return nil
}
