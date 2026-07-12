---
title: Installation
description: Installation instructions for bookmark
---


## Package Managers

### Homebrew (macOS/Linux)

```bash
# Add the tap
brew tap imdevan/bookmark-plus

# Install
brew install bookmark-plus
```

### Arch Linux (AUR)

```bash
# Using yay
yay -S bookmark-plus

# Using paru
paru -S bookmark-plus
```

## GitHub Releases

Download the latest release for your platform from the [releases page](https://github.com/imdevan/bookmark/releases).

### Linux/macOS

```bash
# Download the binary (replace VERSION and PLATFORM with actual values)
# Example: v0.1.0 and linux-amd64
curl -LO https://github.com/imdevan/bookmark/releases/download/VERSION/bookmark-PLATFORM

# Make it executable
chmod +x bookmark-PLATFORM

# Move to your PATH
sudo mv bookmark-PLATFORM /usr/local/bin/bookmark
```

### Windows

Download the `.exe` file from the releases page and add it to your PATH.

## Manual Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/imdevan/bookmark.git
cd bookmark

# Build using just
just build

# Install to /usr/local/bin
sudo just install
```

Or build manually:

```bash
go build -o bookmark ./cmd/bookmark
sudo mv bookmark /usr/local/bin/bookmark
```

## Verify Installation

```bash
bookmark --version
```

## Post-Installation Setup

Initialize your configuration:

```bash
# Create default config
bookmark config init

# Edit config (optional)
bookmark config
```

The bookmark functions will be automatically sourced in your shell. If not, add this to your shell's RC file:

```bash
# For bash/zsh (~/.bashrc or ~/.zshrc)
source ~/.bookmarks/bookmarks.sh

# For fish (~/.config/fish/config.fish)
source ~/.bookmarks/bookmarks.fish

# For nushell (~/.config/nushell/config.nu)
source ~/.bookmarks/bookmarks.nu
```

See [configuration docs](/configuration) for more details.

