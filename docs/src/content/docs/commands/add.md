---
title: add
description: Add a new bookmark
---

Add a new bookmark

## Usage

```bash
bookmark add [alias]
```

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `-c, --config` | string | config file path |
| `-t, --tmux` | bool | set tmux window name (same as alias) |
| `-T, --tmux-name` | string | custom tmux window name |
| `-d, --description` | string | bookmark description |
| `-y, --yes` | bool | skip form, save directly |
| `-f, --file` | string | file to open in editor after navigation |
| `-x, --execute` | string | command to execute after navigation |
| `-s, --source` | string | path to bookmark (instead of current directory) |

## Source

See [add_cmd.go](https://github.com/imdevan/bookmark//blob/main/cmd/bookmark/add_cmd.go) for implementation details.
