# Aliases

An alias manager for your favorite shell

## Features

- An alias manager that works WITH your shell
- Integrates with your favorite shell!
## Install

```bash
# With homebrew
homebrew install imdevan/aliases/aliases

# With AUR
yay -S aliasesaa

# Manual
git clone https://github.com/imdevan/aliases.git
cd aliases
just build && sudo just install
```

See [install](https://devan.gg/aliases/install/)  docs for more information.


## aliases your favorite folder

```bash
~/Projects/favorite-project
aliases
"aliases fp created!"

# Pass a name
aliases foo
"aliases foo created!"

# Rename tmux window on navigation
aliases -t

#  Rename tmux custom window
aliases -T foo
```


## How `aliases` works

A **aliases** in this case is an alias that is sourced into the shell on load time. 

Aliases live in `~/.aliasess/aliasess.sh` by default. The location can be changed via config options. 

aliases is essentially a light weight wrapper around that file. 

### Escaping and Quoting

When aliasess are saved, the tool automatically escapes and quotes all fields (such as directory paths, descriptions, tmux window names, and associated files or scripts) to ensure they are written and parsed correctly regardless of special characters like single/double quotes, pipes (`|`), or newlines. The escaping behavior is adapted automatically to your configured shell (POSIX, Fish, or Nushell). 

## Using different shell? 

aliases is set up for zsh first but works with any shell. Run `aliases config init` to create a custom config file.

See configuration options for more info.

## Commands


```bash
aliases                 # aliases a file
bm                       # Interactive aliases list
aliases add             # Interactive form to add a new aliases
aliases config          # View or edit configuration
aliases config init     # Generate default config file
aliases completion      # Generate shell completion scripts
```

## Configuration

aliases is designed to be highly customizable. 

```bash
aliases config # open config file location
```

Configuration file location: `~/.config/aliases/config.toml`  (`$XDG_CONFIG_HOME/aliases/config.toml`)

See [configuration](./CONFIGURATION.md) for installation options.

## Installation

See [install](./INSTALL.md) for installation options.

## Customization

This template is designed to be customized for your specific CLI tool needs:
1. Edit `package.toml` with your project details (name, module, description, etc.)

2. Run `just sync` to sync changes across all files
3. Review changes with `git diff`
4. Build and test: `just build && just test`

The `package.toml` file is the single source of truth for project metadata. The sync script will update:
- Go module name in `go.mod` and all import paths
- Binary name in justfile and build scripts
- Config paths in `internal/utils/paths.go`
- Completion examples
- README description
- Version in root.go

## Architecture

- `cmd/`              - CLI entrypoint and commands
- `internal/config`   - Configuration management
- `internal/domain`   - Domain models
- `internal/ui`       - Bubble Tea UI components
- `internal/utils`    - Utility functions
- `internal/adapters` - External service adapters (editor, clipboard)

# Thank you!

This project was made by deconstructing a another cli project of mine [Prompter](http://devan.gg/prompter-cli/). Check it out if you like fiddling with coding agents and want a more vim centric way of managing your prompting!


