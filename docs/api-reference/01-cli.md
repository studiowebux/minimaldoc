---
title: CLI Commands
description: Command-line interface reference
tags:
  - api-reference
  - cli
---

# CLI Commands

MinimalDoc command-line interface reference.

## Commands

| Command | Description |
|---------|-------------|
| `build` | Build documentation site |
| `init` | Initialize new documentation |
| `version` | Show version |
| `help` | Show help |

## build

Generate static HTML from Markdown files.

```bash
minimaldoc build [docs-directory] [flags]
```

### Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `docs-directory` | `.` | Path to documentation source |

### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `public` | Output directory |
| `--theme` | `-t` | string | `default` | Theme name |
| `--title` | | string | `Documentation` | Site title |
| `--description` | | string | | Site description |
| `--base-url` | | string | | Base URL for absolute links |
| `--llms` | `-l` | bool | `true` | Generate llms.txt |
| `--clean-urls` | | bool | `false` | Use clean URLs |
| `--openapi` | | bool | `false` | Enable OpenAPI docs |
| `--openapi-dir` | | string | `api` | OpenAPI specs directory |
| `--status` | | bool | `false` | Enable status page |
| `--status-title` | | string | `Service Status` | Status page title |
| `--status-path` | | string | `status` | Status page output path |
| `--changelog` | | bool | `false` | Enable changelog |
| `--changelog-path` | | string | `changelog` | Changelog output path |
| `--stale-warning` | | bool | `true` | Enable stale warnings |
| `--stale-threshold` | | int | `365` | Days before stale |

### Examples

**Basic build:**

```bash
minimaldoc build
```

**Specify source directory:**

```bash
minimaldoc build ./docs
```

**Custom output:**

```bash
minimaldoc build ./docs --output dist
```

**Production build:**

```bash
minimaldoc build ./docs \
  --base-url "https://docs.example.com" \
  --output dist
```

**Full featured:**

```bash
minimaldoc build ./docs \
  --title "My Project" \
  --description "Project documentation" \
  --base-url "https://docs.example.com" \
  --output dist \
  --theme default \
  --openapi \
  --status \
  --changelog \
  --clean-urls
```

**Disable features:**

```bash
minimaldoc build ./docs \
  --llms=false \
  --stale-warning=false
```

## init

Initialize a new documentation site with example files and optional features.

```bash
minimaldoc init [directory] [flags]
```

### Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `directory` | `docs` | Target directory |

### Flags

| Flag | Description |
|------|-------------|
| `--with-toc` | Create TOC.md for custom navigation |
| `--with-status` | Create status page structure |
| `--with-changelog` | Create changelog structure |
| `--with-openapi` | Create OpenAPI specification example |
| `--full` | Include all optional features |

### Generated Structure

**Basic (default):**

```
docs/
├── config.yaml
├── index.md
├── getting-started/
│   ├── installation.md
│   └── quick-start.md
└── guides/
    └── deployment.md
```

**Full (`--full` flag):**

```
docs/
├── config.yaml
├── TOC.md
├── index.md
├── getting-started/
│   ├── installation.md
│   └── quick-start.md
├── guides/
│   └── deployment.md
├── api/
│   └── openapi.yaml
├── __status__/
│   ├── config.yaml
│   ├── components.yaml
│   └── incidents/
│       └── YYYY-MM-DD-example-incident.md
└── __changelog__/
    ├── config.yaml
    └── releases/
        └── 0.1.0.md
```

### Examples

**Basic site:**

```bash
minimaldoc init
```

**Custom directory:**

```bash
minimaldoc init my-docs
```

**With all features:**

```bash
minimaldoc init my-docs --full
```

**With specific features:**

```bash
minimaldoc init my-docs --with-status --with-openapi
```

**API documentation project:**

```bash
minimaldoc init api-docs --with-openapi --with-changelog
```

## version

Show MinimalDoc version.

```bash
minimaldoc version
```

Output:

```
minimaldoc version 1.0.0
```

### Flags

| Flag | Description |
|------|-------------|
| `--short` | Version number only |

## help

Show help for any command.

```bash
minimaldoc help [command]
```

### Examples

```bash
minimaldoc help
minimaldoc help build
minimaldoc build --help
```

## Global Flags

Available for all commands:

| Flag | Description |
|------|-------------|
| `--help` | Show help |
| `--version` | Show version |

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | File not found |
| 4 | Build error |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `MINIMALDOC_CONFIG` | Path to config file |
| `MINIMALDOC_THEME` | Theme name |
| `NO_COLOR` | Disable colored output |

## Configuration Precedence

1. CLI flags (highest priority)
2. Environment variables
3. config.yaml
4. Defaults (lowest priority)

## Shell Completion

Generate completion scripts:

```bash
# Bash
minimaldoc completion bash > /etc/bash_completion.d/minimaldoc

# Zsh
minimaldoc completion zsh > "${fpath[1]}/_minimaldoc"

# Fish
minimaldoc completion fish > ~/.config/fish/completions/minimaldoc.fish

# PowerShell
minimaldoc completion powershell > minimaldoc.ps1
```

## Common Patterns

### Development

```bash
# Quick build for local preview
minimaldoc build ./docs

# Serve locally
python -m http.server -d public 8080
```

### CI/CD

```bash
# Production build with all features
minimaldoc build ./docs \
  --base-url "$DOCS_URL" \
  --output dist \
  --openapi \
  --status \
  --changelog
```

### Docker

```bash
# Build in container
docker run -v $(pwd):/app -w /app golang:1.24-alpine sh -c \
  "go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest && \
   minimaldoc build ./docs --output dist"
```
