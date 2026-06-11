# Bookmark

<img width="407" height="270" alt="screenshot-2026-03-06_14-10-16" src="https://github.com/user-attachments/assets/081eeec6-bbbd-4687-b29a-b659892a5e58" />

A bookmark manager for your favorite shell

## Features

- A bookmark manager that works WITH your shell
- Integrates with your favorite shell!
- Integrates with TMUX and your favorite editor!
## Install

```bash
# With homebrew
homebrew install imdevan/bookmark/bookmark

# With AUR
yay -S bookmark-plus

# Manual
git clone https://github.com/imdevan/bookmark.git
cd bookmark
just build && sudo just install
```

See [install](https://devan.gg/bookmark/install/)  docs for more information.


## Bookmark your favorite folder

```bash
~/Projects/favorite-project
bookmark
"bookmark fb created!"

# Pass a name
bookmark foo
"bookmark foo created!"

# Rename tmux window on navigation
bookmark -t

#  Rename tmux custom window
bookmark -T foo
```

Bookmarks are created by default at `~/.bookmarks/bookmarks.sh`

## Using different shell? 

Bookmark is set up for zsh first but works with any shell. Run `bookmark config init` to create a custom config file.

See configuration options for more info.

## Commands


```bash
bookmark                 # Bookmark a file
bm                       # Interactive bookmark list
bookmark add             # Interactive form to add a new bookmark
bookmark config          # View or edit configuration
bookmark config init     # Generate default config file
bookmark completion      # Generate shell completion scripts
```

## Development

```bash
just sync            # Sync project from package.toml
just build           # Build the binary
just build-run       # Build and run the binary
just dev-build       # Build with debug symbols
just test            # Run tests
just test-verbose    # Run tests with verbose output
just watch           # Watch for changes and rebuild
just cross-platform  # Build for multiple platforms
just install         # Install to /usr/local/bin
just clean           # Remove build artifacts
```

## Configuration

Configuration file location: `$XDG_CONFIG_HOME/go-cli-template/config.toml`

See `example-config.toml` for available configuration options.

## Installation

See `INSTALL.md` for installation options.

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


