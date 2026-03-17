# Context

This project is a Go CLI application for managing directory bookmarks with shell integration. The tool allows users to quickly save and navigate to frequently-used directories, with support for custom aliases, tmux integration, and post-navigation scripts.

# Definitions

- **Bookmark**: A saved directory path with an associated alias for quick navigation
- **Alias**: A short name used to reference a bookmark (e.g., `bm-proj` for `/home/user/projects`)
- **Post-jump script**: Custom commands executed after navigating to a bookmarked directory
- **CRUD**: Create, Read, Update, Delete operations for managing bookmarks
- **Navigation tool**: The underlying command used for directory changes (cd, z, zoxide, etc.)

# Data Structures

## Bookmark File Format (TOML)

```toml
# ~/.bookmarks/bookmarks.toml

[[bookmarks]]
alias = "proj"
path = "/home/user/projects/myapp"
description = "Main application project"
created_at = "2026-02-24T10:30:00Z"
updated_at = "2026-02-24T10:30:00Z"
tmux_window_name = "myapp"
post_jump_script = "source .env && echo 'Welcome to myapp'"

[[bookmarks]]
alias = "docs"
path = "/home/user/projects/myapp/docs"
description = "Documentation folder"
created_at = "2026-02-24T11:00:00Z"
updated_at = "2026-02-24T11:00:00Z"
```

# Command Examples

## 1. Add Bookmark (Current Directory)

```bash
# Auto-generated alias
$ cd /home/user/projects/my-cool-project
$ bookmark
✓ Bookmark created: mcp → /home/user/my/cool/project

# Auto-generated alias with same tmux name
$ cd /home/user/projects/my-cool-project
$ bookmark -t
✓ Bookmark created: mcp → /home/user/my/cool/project


# Custom alias
$ cd /home/user/projects/webapp
$ bookmark web
✓ Bookmark created: web → /home/user/projects/webapp

# Custom alias with same tmux name
$ cd /home/user/projects/webapp
$ bookmark web
✓ Bookmark created: web → /home/user/projects/webapp

# Custom alias with custom tmux name
$ cd /home/user/projects/webapp
$ bookmark web -t wweb
✓ Bookmark created: web → /home/user/projects/webapp

# With description
$ bookmark web --description "Main web application"
✓ Bookmark created: web → /home/user/projects/webapp

# Overwrite confirmation
$ bookmark web
? Bookmark 'web -> ~/projects/webapp' already exists. Overwrite? (y/N) y
use ui/confirmation
✓ Bookmark updated: web → /home/user/projects/webapp
```

## 2. Interactive Browser

```bash
$ bookmark -i
# or
$ bookmark --interactive

# Output (Bubble Tea UI):
use ui/list
```

## 3. List Bookmarks

```bash
$ bookmark list
proj     /home/user/projects/myapp          Main application project
docs     /home/user/projects/myapp/docs     Documentation folder
web      /home/user/projects/webapp
scripts  /home/user/scripts
```

## 4. Navigate to Bookmark

bookmarks are saved as aliases that are managed by the users shell

```bash
$ proj
# Outputs shell command to stdout:
cd /home/user/projects/myapp

# With tmux integration
$ bookmark go proj -t myapp
tmux rename-window 'myapp' && cd /home/user/projects/myapp

# With post-jump script
$ bookmark go proj
cd /home/user/projects/myapp && source .env && echo 'Welcome to myapp'
```

## 5. Delete Bookmark

```bash
$ bookmark delete web
? Delete bookmark 'web'? (y/N) y
use ui/confirmation
✓ Bookmark deleted: web

$ bookmark delete web --force
✓ Bookmark deleted: web
```

## 6. Edit Bookmark

```bash
$ bookmark edit proj
# Opens bookmark_file in editor at bookmark location

$ bookmark edit proj --description "Updated description"
✓ Bookmark updated: proj

$ bookmark edit proj --path /new/path
✓ Bookmark updated: proj → /new/path
```

## 7. Shell Integration

```bash
# Generate shell function
$ bookmark shell-init bash

echo "source $bookmark_file" >> ."$bookarm_shell"rc
```

# v0.1.0

- [x] 1. Core Bookmark Management
---
  - [x] 1.1 Root command to bookmark current folder
    - notes: `bookmark [alias]` command saves current directory with auto-generated or custom alias
    - example: `bookmark` → creates alias "mcp" for `/my/cool/project`
    - example: `bookmark web` → creates alias "web" for current directory
    - output: `✓ Bookmark created: web → /home/user/projects/webapp`

  - [x] 1.2 Auto-generate alias from directory name
    - notes: Default naming convention uses first letters of each "word" in current dir
    - example: `/my/cool/project` → `mcp`
    - example: `/home/user/dev` → `hud`
    - example: `/projects/web-app` → `pwa`

  - [x] 1.3 Optional custom alias via argument
    - notes: Allow user to pass custom bookmark string: `bookmark my-alias`
    - example: `bookmark web` saves current dir as "web"
    - validation: alias must be alphanumeric + hyphens/underscores only

  - [x] 1.4 Confirmation before overwriting existing bookmark
    - notes: Prompt user if alias already exists before replacing
    - example: `? Bookmark 'web' already exists. Overwrite? (y/N)`
    - behavior: default to No, require explicit confirmation

  - [x] 1.5 -y to accept overwrite from flag
    - bypass confirmations

  - [x] 1.6 -f file flag
    - [x] - f and file to add open file in editor to bookmark, executed after navigation
      - file bookmarks will navigate to folder and open file in configured editor
    - example:
      ```bash
      cd ~/Work/foo-bar

      bookmark -t -f plan.md

      bookmark 
      fb -> ~/Work/foo-bar + tmux + open plan.md
      created!

      created alias: 
      alias fb="$nav_command ~/Work/foo-bar && tmux rename-window fb" && $editor plan.md"
      ```

  - [x] 1.7 -e flag
    - no args: open bookmarks_file in configured editor
    - existing bookmark: open in editor at bookmark line
    - non-existing bookmark: create and open in editor
    - example:
      ```bash
      cd ~/Work/foo-bar

      bookmark -e

      bookmark fb -> ~/Work/foo-bar
      => open bookmark_file at fb location
      ```

  - [x] 1.9 -x/execute flag
    - added to alias after tmux rename but before open 
    - example:
      ```bash
      cd ~/Work/foo-bar

      bookmark -x 'cowsay "toast"'

      bookmark:

      fb -> ~/Work/foo-bar
      cowsay "toast"

      created!
      ```
      


  - [x] 1.10 -s/source flag (or b/bookmark if more appropriate?)
    - location to bookmark 
    - example:
      ```bash
      cd ~/Work/foo-bar

      bookmark -s ~/Documents/bar

      bookmark b -> ~/Documents/bar created
        ```

- [x] 1.11 -T/Tmux flag 
  - pass alternate rename variable
    - example:
      ```bash
      cd ~/Work/foo-bar

      bookmark -T bar

      bookmark b -> ~/Documents/bar created
      tmux rename -> bar
        ```



  - [x] 1.12 Alias creation
    - aliases should always:
      - navigate to folder
      - then rename tmux if present
      - then execute any additional script if present
      - then open associated file
      - comments at end of line
    - this structure should be relied on to use to format the alias visually in the ui


  - [x] 1.13 root tests
    - Reasonably test all features from task 1.01 through 1.7 (only)

- [x] 2. Interactive Bookmark Browser
---
  - [x] 2.1 `-i` flag to view filterable list of bookmarks
    - notes: Display all bookmarks with search/filter capability using Bubble Tea inline UI
    - example input: `bookmark -i` or `bookmark --interactive`
    - example output:
      ```
      ┌─ Bookmarks ────────────────────────────────────┐
      │ Search: proj_                                  │
      ├────────────────────────────────────────────────┤
      │ > proj     ~/projects/myapp                    │
      │   docs     ~/projects/myapp/docs               │
      └────────────────────────────────────────────────┘
      ```
    - implementation: Use `internal/ui/list.go` pattern with Bubble Tea

    - add configs:
    - [x] 2.1.1 home_icon = "~" (default) used in list view for home can be nerd font icon
    - [x] 2.1.2 default_sort_by = "newest" (default) | "latest" | "A to Z" | "Z to A" use as standin for $HOME directory

    - show associated functions
    - [x] 2.1.3 show tmux rename  name
    \uebc8
    - [x] 2.1.5 Define row under description and location metadata
      - metadata contains: all associated "function" calls
      - e.g.: tmux,  script, editor open, shell # show shell last here
     name   name (/e691) some shell function
      - show editor with configured editor
      - example
      ```
       │   bok                                                │
 │   ~/Projects/playground/bookmark                     │
 │    bok •  plan/plan.md •  cowsay 'hi'             │

      ```

  - [ ] 2.2 CRUD operations in interactive mode
    - notes: Support Create, Read, Update, Delete actions from the list view
    - keybindings:
      - `Enter`: Navigate to selected bookmark
      - `e`: Edit selected bookmark
      - `d`: Delete selected bookmark (with confirmation)
      - ~`n`: Create new bookmark~
      - `q`: or `Esc`: Quit
    - example: Press `d` → `? Delete bookmark 'proj'? (y/N)`

    - [x] 2.2.1 enter: Enter`: should execute defined alias (and all associated operations)
    - [x] 2.2.2 e: open selected book mark in editor
            - resource bookmarks on save (if in nvim)
    - [x] 2.2.3 d/D: delete selected bookmark (d: confirmation | D: no conf)
    - [x] 2.2.4 ~n: Create new bookmark~ -> moved to v2
    - [x] 2.2.5 q`: or `Esc`: Quit


  - [x] 2.3 Navigate to selected bookmark
    - notes: Select bookmark from list to navigate to that directory
    - behavior: Output shell command to stdout for evaluation
    - example output: `cd /home/user/projects/myapp`

  - [x] 2.4 Add interact_alias prop config
    - interactive_alias = "bm" default
    - if set adds bm function (currently ./bm.sh)
    - add to the top of the bookmarks alias file. adapt for configured shell
    - under function_alias


- [x] 3. Advanced Features
---
  - [x] 3.1 `-t` flag for tmux window naming
    - notes: Optional flag to define tmux window name when navigating to bookmark
    - example: `bookmark go proj -t myapp`
    - output: `tmux rename-window 'myapp' && cd /home/user/projects/myapp`
    - storage: Save `tmux_window_name` field in bookmark TOML

  - [x] 3.2 Post-jump script execution
    - notes: Define and execute custom scripts after navigation
    - example config in bookmark:
      ```toml
      post_jump_script = "source .env && echo 'Welcome!'"
      ```
    - output: `cd /path && source .env && echo 'Welcome!'`
    - validation: Escape shell special characters for safety

  - [x] 3.3 Bookmark descriptions via comments
    - notes: Support adding descriptions/comments to bookmarks for documentation
    - example: `bookmark web --description "Main web application"`
    - storage: Save `description` field in bookmark TOML
    - display: Show in list view and interactive browser

- [x] 4. Configuration System
---
  - [x] 4.1 Navigation tool selection
    - notes: Config option to choose navigation method: none, cd, z, zoxide, etc.
    - config field: `navigation_tool = "cd"`
    - valid values: `"cd"`, `"z"`, `"zoxide"`, `"none"`
    - behavior: Changes output command format (e.g., `z /path` vs `cd /path`)

  - [x] 4.2 Bookmark storage location
    - notes: Configurable bookmark file location
    - config field: `bookmarks_location = "~/.bookmarks/"`
    - default: `~/.bookmarks/`
    bookmark file will be based on shell: e.g.  `~/.bookmarks/bookmarks.sh`

  - [x] 4.3 Shell type configuration
    - config.editor
    - notes: Define which shell the user uses (bash, zsh, fish, etc.)
    - config field: `shell = "zsh"`
    - valid values: `"bash"`, `"zsh"`, `"fish"`, "nu"
    - usage: Affects shell-init command output format

  - [x] 4.4 Function alias
      add to config example and init
      function_alias = true (default) | string alias_name | false (do not create function alias will not source bookmarks)
      if true add to top of bookmarks.sh, or relative shell configuration
      add comments indicating what is happening and why
      ```bash
      bookmark() {
        command bookmark "$@" && source ~/.bookmarks/bookmarks.sh
      }

#  v0.2.0


- [x] 3 list enter: should execute defined alias (and all associated operations)

4 bookmark list should be more verbose and show associated metadata

# v0.3.0

- [ ] 1 show description next to title in list and interactive as # ...
  - create description with flag for adding and editing if not present
  - show show as text muted color in interactive  and list mode
- [ ] 2 interactive add
- [ ] 3 interactive edit
- [ ] 4 improve bookmark created
- [ ] 5 improve first call
  - bookmark config - should check if inited and ask to do so if not
    - config init should also ask if user would like check / add source to users shellrc file
# v?

- [ ] 1 Support multiple shell types 
  - example: 
  ```bash
  shell = "zsh", "nu"

  # creates:
  ~/.bookmarks/bookmarks.sh
  ~/.bookmarks/bookmarks.nu
  ```

- [ ] 2 bookmark sync command 
  - syncs bookmarks file based on:
    - **these should probably just be vars at the top of the bookmarks file**
    - config.shell
    - prioritize first shell in list as source of truth
    - confirm with user before updating out of sync alternate shells
    - config.home_char
    - config.navigation_tool
    - config.editor

3 Interactive create new bookmark

