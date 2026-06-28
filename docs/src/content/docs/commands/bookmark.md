---
title: bookmark
description: A bookmark manager for your favorite shell
---

A bookmark manager for your favorite shell

## Usage

```bash
bookmark [alias]
bookmark [command]
```

## Description

newRootCmd creates the root command for the bookmark CLI.

The root command serves multiple purposes:
  - Without arguments: Opens interactive bookmark browser (if configured)
  - With alias argument: Navigates to the bookmarked directory
  - With --interactive/-i: Forces interactive mode
  - With --edit/-e: Opens bookmarks file in editor
  - With --version/-v: Prints version information

When adding a bookmark, you can specify:
  - --description/-d: Add a description to the bookmark
  - --tmux/-t: Set tmux window name to match alias
  - --tmux-name/-T: Set custom tmux window name
  - --file/-f: Specify a file to open after navigation
  - --execute/-x: Run a command after navigation
  - --source/-s: Bookmark a different path than current directory
  - --yes/-y: Skip confirmation prompts

Examples:

	# Add bookmark for current directory
	bookmark myproject

	# Add bookmark with description
	bookmark myproject -d "My awesome project"

	# Navigate to bookmark (in interactive mode)
	bookmark

	# Navigate to specific bookmark
	bookmark myproject

	# Edit bookmarks file
	bookmark -e

	# List all bookmarks
	bookmark list

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `-d, --description` | string | bookmark description |
| `-f, --file` | string | file to open in editor after navigation |
| `-s, --source` | string | path to bookmark (instead of current directory) |
| `-t, --tmux` | bool | set tmux window name (same as alias) |
| `-T, --tmux-name` | string | custom tmux window name |
| `-x, --execute` | string | command to execute after navigation |
| `-c, --config` | string | config file path |
| `-i, --interactive` | bool | interactive bookmark browser |
| `-a, --add` | bool | interactive add bookmark form |
| `-e, --edit` | bool | open bookmarks file in editor |
| `-y, --yes` | bool | skip confirmation, and interactive prompts |
| `-v, --version` | bool | print version information |

## Available Commands

- [`add`](/commands/add) - Add a new bookmark
- [`completion`](/commands/completion) - Generate shell completion scripts
- [`config`](/commands/config) - View or edit configuration
- [`config init`](/commands/config-init) - Generate a default config file
- [`delete`](/commands/delete) - Delete a bookmark
- [`list`](/commands/list) - List all bookmarks

## Source

See [root.go](https://github.com/imdevan/bookmark//blob/main/cmd/bookmark/root.go) for implementation details.
