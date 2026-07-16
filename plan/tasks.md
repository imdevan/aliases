# Context

`aliases` is an alias manager for your favorite shell. Companion to `b` (bookmark manager). Lives at `cmd/aliases`. Binary name TBD: `a` | `am` | `aliases`.

# Definitions

- **alias file** — default storage at `~/.aliases/aliases.zsh`; where `am add` writes new aliases
- **index folders** — user-defined glob patterns of dirs to scan for alias files (e.g. `~/dotfiles`)
- **wrapper function** — shell function sourced by user's rc file; re-sources alias file after add/edit/delete and handles `am` output (e.g. copy to clipboard)

# v0.1.0

## Feature 1: Simplify
At it's core aliases is a simpler version of `bookmark`.
- [ ] 1.1 rename bookmark / bm to aliases / al 
- [ ] 1.2 remove directory, file, and tmux rename flags. 
- the bookmark script becomes the alias value
- alias is only name, value, description. 
- aliases add name value description


## Feature 2: Project bootstrap
  - [x] 2.1 Set up `cmd/aliases` entrypoint and wire cobra root command
  - [ ] 2.2 Define domain types: `Alias` (name, value, description, source file) simplify bookmark
  - [ ] 2.3 Config scaffold: `~/.config/aliases/config.toml`, options: `shell`, `alias_file`, `index_folders`, `cache_interval`, `script_icons`

## Feature 3: Alias storage
  - [ ] 3.1 Read/write aliases from default alias file (`~/.aliases/aliases.zsh`) basically same as bookmark, but change the location
  - [ ] 3.2 Parse alias syntax from shell files (`alias name='value' # description`)
  - [ ] 3.3 Write alias back in correct shell syntax on add/edit/delete

## Feature 4: `am add`
  - [ ] 4.1 Same as bookmark but simpler
  - [ ] 4.2 Write to alias file
## Feature 5: `am edit`
  - [ ] 5.1 Same as bookmark but simpler
  - [ ] 5.2 Pre-fill form with current values
  - [ ] 5.3 Write updated alias back to source file

## Feature 6: `am delete`
  - [ ] 6.1 Same as bookmark 
  - [ ] 6.2 Confirm prompt

## Feature 7: `am` (interactive list)
  - [ ] 7.1 Same as bookmark but simpler

## Feature 8: Shell wrapper
  - [ ] 8.1 Generate shell wrapper function (`am()`) for zsh
  - [ ] 8.2 Re-source alias file after mutating commands
  - [ ] 8.3 `am config init` prints or installs wrapper snippet to rc file

## Feature 9: `am config` / `am config init`
  - [ ] 9.1 `am config` — show current config
  - [ ] 9.2 `am config init` — interactive setup: shell, alias file path, index folders

