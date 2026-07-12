---
title: config init
description: Generate a default config file
---

The config init command creates a new configuration file with default values
at the standard XDG config location ($XDG_CONFIG_HOME/bookmark/config.toml).

### Example

```bash
# Generate default config
bookmark config init

# Overwrite existing config
bookmark config init --force

# Generate and open in editor
bookmark config init --editor
```

:::note
The generated config file includes commented examples for all available options.
:::

## Usage

```bash
bookmark config init
```

## Flags




| Flag | Type | Description |
|------|------|-------------|
| -e, --editor | bool | open config in editor after creation |
| -f, --force | bool | overwrite existing config |


## Source

See [config_init.go](https://github.com/imdevan/bookmark/blob/main/cmd/bookmark/config_init.go) for implementation details.
