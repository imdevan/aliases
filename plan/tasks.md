# Context

`aliases` is an alias manager for your favorite shell. Companion to `b` (bookmark manager). Lives at `cmd/aliases`. Binary name TBD: `a` | `am` | `aliases`.

# Definitions

- **alias file** — default storage at `~/.aliases/aliases.zsh`; where `am add` writes new aliases
- **index folders** — user-defined glob patterns of dirs to scan for alias files (e.g. `~/dotfiles`)
- **wrapper function** — shell function sourced by user's rc file; re-sources alias file after add/edit/delete and handles `am` output (e.g. copy to clipboard)

# v0.1.0

## Feature 1: Simplify
At it's core aliases is a simpler version of `bookmark`.
- [x] 1.1 rename bookmark / bm to aliases / al 
- [x] 1.2 remove directory, file, and tmux rename flags. 
- the bookmark script becomes the alias value
- alias is only name, value, description. 
- aliases add name value description


## Feature 2: Project bootstrap
  - [x] 2.1 Set up `cmd/aliases` entrypoint and wire cobra root command
  - [x] 2.2 Define domain types: `Alias` (name, value, description, source file) simplify bookmark
  - [x] 2.3 Config scaffold: `~/.config/aliases/config.toml`, options: `shell`, `alias_file`, `index_folders`, `cache_interval`, `script_icons`

## Feature 3: Alias storage
  - [x] 3.1 Read/write aliases from default alias file (`~/.aliases/aliases.zsh`) basically same as bookmark, but change the location
  - [x] 3.2 Parse alias syntax from shell files (`alias name='value' # description`)
  - [x] 3.3 Write alias back in correct shell syntax on add/edit/delete

## Feature 4: `am add`
  - [x] 4.1 Same as bookmark but simpler
  - [x] 4.2 Write to alias file
## Feature 5: `am edit`
  - [x] 5.1 Same as bookmark but simpler
  - [x] 5.2 Pre-fill form with current values
  - [x] 5.3 Write updated alias back to source file

## Feature 6: `am delete`
  - [x] 6.1 Same as bookmark 
  - [x] 6.2 Confirm prompt

## Feature 7: `am` (interactive list)
  - [x] 7.1 Same as bookmark but simpler

## Feature 8: Shell wrapper
  - [x] 8.1 Generate shell wrapper function (`am()`) for zsh
  - [x] 8.2 Re-source alias file after mutating commands
  - [x] 8.3 `am config init` prints or installs wrapper snippet to rc file

## Feature 9: `am config` / `am config init`
  - [x] 9.1 `am config` — show current config
  - [x] 9.2 `am config init` — interactive setup: shell, alias file path, index folders

