---
title: edit
description: Edit a bookmark or open bookmarks file in editor
---

The edit command opens a bookmark for editing or opens the entire bookmarks file in the editor.

### Example

```bash
# Open all bookmarks in editor
bookmark edit

# Open specific bookmark in form
bookmark edit my-alias
```

## Usage

```bash
bookmark edit [alias]
```

## Flags




| Flag | Type | Description |
|------|------|-------------|
| -c, --config | string | config file path |


## Source

See [edit_cmd.go](https://github.com/imdevan/bookmark/blob/main/cmd/bookmark/edit_cmd.go) for implementation details.
