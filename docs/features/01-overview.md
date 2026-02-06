---
title: Features Overview
description: Complete feature list and capabilities of MinimalDoc
tags:
  - features
  - overview
---

# Features Overview

MinimalDoc provides everything needed for modern documentation sites.

## Feature Matrix

| Category | Feature | Description |
|----------|---------|-------------|
| **Content** | Markdown | Full Markdown with GFM extensions |
| | Frontmatter | YAML metadata for pages |
| | Admonitions | Callout blocks (info, warning, danger, success, note, question) |
| | Syntax Highlighting | 100+ languages via Chroma with copy button |
| | Stale Warnings | Configurable warnings for outdated content |
| **Navigation** | Auto Navigation | From folder structure |
| | Custom TOC | Via TOC.md file |
| | Page TOC | Auto-generated from headings |
| | Scrollspy | Active heading highlight |
| **Search** | Full-Text | Client-side search index |
| | Keyboard | Cmd+K / Ctrl+K shortcut |
| | Fuzzy Matching | Tolerant search |
| **Design** | Responsive | Mobile, tablet, desktop |
| | Dark Mode | Toggle with localStorage persistence |
| | Themes | CSS-only customization (default, yellow) |
| **Pages** | Landing | Marketing homepage with hero, features, steps |
| | Portfolio | Project showcase with tags and filtering |
| | Contact | Contact page with email and info |
| | Status | Service health with incidents, maintenance, uptime |
| | Changelog | Version history with RSS feed |
| **SEO** | Sitemap | Automatic sitemap.xml generation |
| | Meta Tags | Open Graph, Twitter Cards |
| | Canonical URLs | Configurable base URL |
| | LLM Output | llms.txt and llms-full.txt generation |
| **API** | OpenAPI | Interactive docs with testing UI |
| | Code Samples | curl, JavaScript, Go, Python, Swift |
| | Schema Viewer | Request/response visualization |

## Design Philosophy

### Minimal Configuration

Works out of the box. Add `config.yaml` only when needed.

```bash
# This works
minimaldoc build ./docs

# Configuration optional
minimaldoc build ./docs --title "My Docs"
```

### Static Output

Pure HTML, CSS, JavaScript. No server runtime.

```
public/
├── index.html      # Static HTML
├── css/main.css    # Compiled CSS
├── js/search.js    # Vanilla JS
└── ...
```

Deploy anywhere: GitHub Pages, Netlify, S3, Nginx, Apache.

### Fast Builds

Written in Go. Builds hundreds of pages in seconds.

```bash
$ time minimaldoc build ./large-docs
Built 247 pages in 1.2s
```

### CSS-Only Themes

Themes are CSS variables only. No template editing.

```css
:root[data-theme="light"] {
  --bg-primary: #fafafa;
  --text-primary: #1a1a1a;
}
```

## Comparison

| Feature | MinimalDoc | Hugo | MkDocs | Docusaurus |
|---------|------------|------|--------|------------|
| Language | Go | Go | Python | Node.js |
| Config | YAML | TOML/YAML | YAML | JS |
| Themes | CSS-only | Templates | Templates | React |
| Search | Built-in | Plugin | Plugin | Built-in |
| OpenAPI | Built-in | Plugin | Plugin | Plugin |
| Status Page | Built-in | No | No | No |
| Changelog | Built-in | No | No | No |
| Portfolio | Built-in | No | No | No |
| Landing Page | Built-in | Theme | Theme | Built-in |
| LLM Output | Built-in | No | No | No |

## Core Strengths

**Simplicity**
- Single binary, no dependencies
- Minimal configuration
- Sensible defaults

**Speed**
- Fast builds (Go)
- Small output (no framework bloat)
- Client-side search (no server)

**Documentation Focus**
- Built for docs, not blogs
- Navigation-first design
- Technical content features

**Integration Ready**
- OpenAPI/Swagger native
- Status page included
- CI/CD friendly
