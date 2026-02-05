---
title: MinimalDoc
description: A modern, minimal static site generator for beautiful documentation
tags:
  - home
  - documentation
  - static-site-generator
---

# MinimalDoc

A fast, minimal static site generator designed for documentation. Write in Markdown, generate beautiful sites.

## Why MinimalDoc?

MinimalDoc focuses on what matters: your content. No complex configuration, no bloated dependencies, no unnecessary features.

| Feature | MinimalDoc |
|---------|------------|
| Build Speed | Sub-second for most sites |
| Configuration | Single YAML file |
| Dependencies | Zero runtime dependencies |
| Output | Pure static HTML/CSS/JS |
| Deployment | Any static host |

## Core Features

**Content**
- Markdown with YAML frontmatter
- GitHub Flavored Markdown (tables, task lists, strikethrough)
- Syntax highlighting for 100+ languages
- Custom admonitions (info, warning, danger, success)

**Navigation**
- Automatic navigation from folder structure
- Custom navigation via TOC.md
- Auto-generated page table of contents
- Keyboard-accessible search (Cmd+K)

**Design**
- Clean, responsive layout
- Dark/light mode with persistence
- CSS-only theming
- Mobile-first approach

**Integrations**
- OpenAPI/Swagger documentation
- Status page with incidents and uptime
- Changelog with RSS feed
- LLM-friendly output (llms.txt)

## Quick Start

```bash
# Install
go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest

# Initialize (basic)
minimaldoc init my-docs

# Or initialize with all features (status page, changelog, OpenAPI)
minimaldoc init my-docs --full

# Build
minimaldoc build my-docs
```

Output in `public/`. Open `public/index.html` to preview.

## Documentation

- [Installation](getting-started/01-installation.md) - Install MinimalDoc
- [Quick Start](getting-started/02-quick-start.md) - Create your first site
- [Configuration](getting-started/03-configuration.md) - Configure your site
- [Features Overview](features/01-overview.md) - Explore all features

## License

MIT License. See [LICENSE](https://github.com/studiowebux/minimaldoc/blob/main/LICENSE).
