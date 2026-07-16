---
title: delete
description: Delete a bookmark
---

The delete command removes a bookmark by its alias.
By default, it will prompt for confirmation before deleting.

### Example

```bash
~/foo
$ bookmark	# create alias "f" that points to ~/foo

~/foo
$ bookmark delete f	# delete alias "f"
```

:::note
Delete confirmation can be skipped by setting `confirm_delete=false` in the [config](https://devan.gg/bookmark/configuration/)
:::

## Usage

```bash
aliases delete <alias>
```

## Source

See [delete_cmd.go](https://github.com/imdevan/aliases/blob/main/cmd/aliases/delete_cmd.go) for implementation details.
