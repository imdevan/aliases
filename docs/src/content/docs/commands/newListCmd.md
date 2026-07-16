---
title: newListCmd
description: List all bookmarks
---

newListCmd creates the list command for displaying all bookmarks.

The list command shows all bookmarks in a formatted table with:
  - Alias: The bookmark name
  - Path: The directory path
  - Description: Optional bookmark description

The output is formatted with proper alignment for easy reading.

Examples:

	# List all bookmarks
	bookmark list

	# Use with custom config
	bookmark list -c ~/.config/bookmark/custom.toml

## Usage

```bash
aliases list
```

## Source

See [list_cmd.go](https://github.com/imdevan/aliases/blob/main/cmd/aliases/list_cmd.go) for implementation details.
