---
title: delete
description: Delete a bookmark
---

The root command serves multiple purposes:
		- Without arguments: Opens interactive bookmark browser (if configured)
		- With alias argument: Navigates to the bookmarked directory

### Example

```bash
~/foo
$ bookmark			# create alias "f" that points to ~/foo

~/foo
$ bookmark bar	# create alias "bar" that points to ~/foo
```

:::note
On first call `~/.bookmark/bookmarks.sh` and `~/.config/bookmark/config.toml` will be created.
:::

## Usage

```bash
bookmark delete <alias>
```

## Flags

### bookmark

Options related to adding a bookmark.
#### Example
```bash
~/foo
$ bookmark foo -t -x "just start-dev" -f "./example.md" -d "an example bookmark"
```

Creates a shell alias `foo` that:
- navigates to `~/foo`
- renames the current tmux window to `foo`
- run script `just start-dev`
- then opens `~/foo/example.md` in the shells default editor
- with a comment description that can be seen when looking at the bookmark list or in the generated .sh file.

| Flag | Type | Description |
|------|------|-------------|

### config

Use a different config other than the standard `~/.config/bookmark/config.toml`"
#### Example
```bash
~/foo
$ bookmark -c ~/foo/local-bookmark-config.toml
```
Creates a shell alias `foo` that uses `~/foo/local-bookmark-config.toml` for config options

| Flag | Type | Description |
|------|------|-------------|
| -c, --config | string | config file path |

### interactive


:::note
`-i` flag only prints the bookmark location. Use `bm` alias for interactive navigation.
:::

| Flag | Type | Description |
|------|------|-------------|

### meta



| Flag | Type | Description |
|------|------|-------------|

### Flags



| Flag | Type | Description |
|------|------|-------------|
| -f, --force | bool | skip confirmation |


## Source

See [root.go](https://github.com/imdevan/bookmark/blob/main/cmd/bookmark/root.go) for implementation details.
