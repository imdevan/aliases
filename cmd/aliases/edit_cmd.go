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

type editOptions struct {
	configPath string
}

func newEditCmd() *cobra.Command {
	opts := &editOptions{}
	cmd := &cobra.Command{
		Use:   "edit [name]",
		Short: "Edit an alias or open aliases file in editor",
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
	aliasManager := alias.NewManager(cfg.ResolvedAliasFile(), cfg.Shell, cfg.FunctionAlias, cfg.InteractiveAlias, cfg.IndexFolders)

	// If no alias name provided, just open the aliases file in editor
	if len(args) == 0 {
		return openEditor(cfg.Editor, cfg.ResolvedAliasFile(), 0)
	}

	name := args[0]

	// Check if alias exists
	exists, err := aliasManager.Exists(name)
	if err != nil {
		return err
	}

	var al domain.Alias
	if exists {
		al, err = aliasManager.Get(name)
		if err != nil {
			return err
		}
	} else {
		al = domain.Alias{
			Name:  name,
			Value: "",
		}
	}

	theme := ui.ThemeFromConfig(cfg)
	m := ui.NewAliasFormModelEdit(theme, al)
	if !exists {
		m = m.WithTitle(fmt.Sprintf("'%s' Not Found, Add Alias", name))
	}

	progOpts := tty.GetProgramOptions(tea.WithoutSignalHandler())
	p := tea.NewProgram(m, progOpts...)
	result, err := p.Run()
	if err != nil {
		return err
	}

	fm, ok := result.(ui.AliasFormModel)
	if !ok || !fm.IsCompleted() {
		fmt.Println(ui.CanceledMessage(theme, "Edit"))
		return nil
	}

	newName, newValue, newDesc := fm.Values()

	// If the name changed and we are editing an existing one, delete the old one
	if exists && newName != name {
		if err := aliasManager.Delete(name); err != nil {
			return err
		}
	}

	newAl := domain.Alias{
		Name:        newName,
		Value:       newValue,
		Description: newDesc,
		SourceFile:  al.SourceFile,
	}

	if err := aliasManager.Add(newAl); err != nil {
		return err
	}

	action := "created"
	if exists {
		action = "updated"
	}
	printSuccess(cfg, action, newName, newValue)
	return nil
}
