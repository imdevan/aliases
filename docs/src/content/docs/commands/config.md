---
title: config
description: View or edit configuration
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
bookmark config
```

## Source

See [config.go](https://github.com/imdevan/bookmark/blob/main/cmd/bookmark/config.go) for implementation details.
