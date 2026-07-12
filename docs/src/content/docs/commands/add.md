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
| -c, --config | string | config file path |
| -d, --description | string | bookmark description |
| -x, --execute | string | command to execute after navigation |
| -f, --file | string | file to open in editor after navigation |
| -s, --source | string | path to bookmark (instead of current directory) |
| -t, --tmux | bool | set tmux window name (same as alias) |
| -T, --tmux-name | string | custom tmux window name |
| -y, --yes | bool | skip form, save directly |


## Source

See [add_cmd.go](https://github.com/imdevan/bookmark/blob/main/cmd/bookmark/add_cmd.go) for implementation details.
