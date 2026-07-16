---
title: newCompletionCmd
description: Generate shell completion scripts
---

newCompletionCmd creates the completion command for generating shell completion scripts.

The completion command generates shell completion scripts for various shells.
This enables tab-completion for aliases commands and aliases.

Supported shells:
  - bash
  - zsh
  - fish
  - powershell

Examples:

	# Generate bash completion
	aliases completion bash > /etc/bash_completion.d/aliases

	# Generate zsh completion
	aliases completion zsh > ~/.zsh/completion/_aliases

	# Generate fish completion
	aliases completion fish > ~/.config/fish/completions/aliases.fish

	# Generate powershell completion
	aliases completion powershell > aliases.ps1

## Usage

```bash
aliases completion [bash|zsh|fish|powershell]
```

## Source

See [completion.go](https://github.com/imdevan/aliases/blob/main/cmd/aliases/completion.go) for implementation details.
