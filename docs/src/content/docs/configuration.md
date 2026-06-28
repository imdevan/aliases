---
title: Configuration
description: Configuration options for bookmark
---


Configuration file location: `$XDG_CONFIG_HOME/bookmark/config.toml`

Most likely `~/.config/bookmark/config.toml`

## Configuration Options

The following options can be set in your configuration file:

### General Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `editor` | string | `nvim` | Editor to use for editing bookmarks and config files |
| `interactive_default` | bool | `false` | Start in interactive mode by default when no arguments are provided |
| `confirm_delete` | bool | `true` | Prompt for confirmation before deleting a bookmark. When `false`, delete immediately |
| `plain_text` | bool | `false` | Output plain text confirmation messages and errors instead of styled output |

### Bookmark Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `bookmark_location` | string | `~/.bookmarks` | Directory where bookmark files are stored |
| `navigation_tool` | string | `cd` | Tool to use for navigation. Options: `cd`, `z`, `zoxide`, `none` |
| `shell` | string | auto-detected | Shell type. Options: `bash`, `zsh`, `fish`, `nu` |

### Auto-Alias Generation

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `auto_alias_separator` | string | `""` | Character between first letters in auto-generated aliases. Empty = `mcp`, `-` = `m-c-p` |
| `auto_alias_lowercase` | bool | `true` | Convert auto-generated aliases to lowercase |

### Display Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `home_icon` | string | `~` | Icon to represent home directory in list view (can be nerd font icon) |
| `default_sort_by` | string | `newest` | Default sort order. Options: `newest`, `oldest`, `a-z`, `z-a` |
| `list_spacing` | string | `space` | List item spacing. Options: `compact` (title only), `tight` (title + description, no margin), `space` (default, with spacing) |

### Function Aliases

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `function_alias` | string | `true` | Enable function wrapper that auto-sources bookmarks after running bookmark commands. Options: `true` (use default name "bookmark"), `custom_name` (use custom function name), `false` (disabled) |
| `interactive_alias` | string | `bm` | Alias for interactive bookmark navigation. The function displays the TUI and executes the selected bookmark command. Options: `bm` (default), `custom_name` (use custom function name), `false` (disabled) |

### Colors

Colors support named, numeric, or hex values (e.g., `7`, `13`, `"#ff8800"`).

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `headings` | string | `15` | Color for headings |
| `primary` | string | `02` | Primary color |
| `secondary` | string | `06` | Secondary color |
| `text` | string | `07` | Text color |
| `text_highlight` | string | `06` | Highlighted text color |
| `description_highlight` | string | `05` | Highlighted description color |
| `tags` | string | `13` | Tags color |
| `flags` | string | `12` | Flags color |
| `muted` | string | `08` | Muted text color |
| `accent` | string | `13` | Accent color |
| `border` | string | `08` | Border color |
| `error` | string | `01` | Error and validation color |

## Example Configuration

```toml
# General
editor = "nvim"
interactive_default = false
confirm_delete = true
plain_text = false

# Bookmark settings
bookmark_location = "~/.bookmarks"
navigation_tool = "cd"
shell = "bash"

# Auto-alias generation
auto_alias_separator = ""
auto_alias_lowercase = true

# Display
home_icon = "~"
default_sort_by = "newest"

# Function aliases
function_alias = "true"
interactive_alias = "bm"

# UI
list_spacing = "space"

# Colors (numeric or hex values)
headings = "15"
primary = "02"
secondary = "06"
text = "07"
text_highlight = "06"
description_highlight = "05"
tags = "13"
flags = "12"
muted = "08"
border = "08"
error = "01"
```

## Initializing Configuration

To create a new configuration file with default values:

```bash
bookmark config init
```

To overwrite an existing configuration:

```bash
bookmark config init --force
```

To create and immediately open in your editor:

```bash
bookmark config init --editor
```

## Editing Configuration

To edit your configuration file:

```bash
bookmark config
```

This will open the config file in your configured editor.

## See Also

- [Config API Documentation](/api/config) - Detailed API documentation for the config package
- [Domain API Documentation](/api/domain) - Domain models including Config struct
- [config init command](/commands/config-init) - Generate default config file
