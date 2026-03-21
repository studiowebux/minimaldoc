---
title: Installation
description: Install MinimalDoc on your system
tags:
  - getting-started
  - installation
---

# Installation

## Prerequisites

- Go 1.24 or later
- Git (optional, for cloning)

## Install from Source

```bash
git clone https://github.com/studiowebux/minimaldoc.git
cd minimaldoc
go build -o minimaldoc ./cmd/minimaldoc
```

Move the binary to your PATH:

```bash
# macOS/Linux
sudo mv minimaldoc /usr/local/bin/

# Or add to user bin
mv minimaldoc ~/.local/bin/
```

## Go Install

```bash
go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
```

Ensure `$GOPATH/bin` is in your PATH.

## Binary Releases

Download pre-built binaries from the [releases page](https://github.com/studiowebux/minimaldoc/releases).

Available for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Verify Installation

```bash
minimaldoc --version
```

Expected output:

```
minimaldoc version 1.0.0
```

## Shell Completion

Generate shell completions:

```bash
# Bash
minimaldoc completion bash > /etc/bash_completion.d/minimaldoc

# Zsh
minimaldoc completion zsh > "${fpath[1]}/_minimaldoc"

# Fish
minimaldoc completion fish > ~/.config/fish/completions/minimaldoc.fish
```

## Next Steps

- [Quick Start](02-quick-start.md) - Create your first documentation site
- [Configuration](03-configuration.md) - Configure your site
