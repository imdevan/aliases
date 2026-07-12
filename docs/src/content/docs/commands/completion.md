---
title: completion
description: Generate shell completion scripts
---

Examples:
  bookmark completion bash > /etc/bash_completion.d/bookmark
  bookmark completion zsh > ~/.zsh/completion/_bookmark
  bookmark completion fish > ~/.config/fish/completions/bookmark.fish
  bookmark completion powershell > bookmark.ps1

## Usage

```bash
bookmark completion [bash|zsh|fish|powershell]
```

## Source

See [completion.go](https://github.com/imdevan/bookmark/blob/main/cmd/bookmark/completion.go) for implementation details.
