---
title: Features Overview
description: Complete feature list and capabilities of MinimalDoc
tags:
  - features
  - overview
since: "v1.0"
---

# Features Overview

MinimalDoc provides everything needed for modern documentation sites. Pure
static output designed for CI/CD integration—automate updates to documentation,
status pages, and changelogs on every commit.

## Feature Matrix

| Category       | Feature             | Description                                                     | Since |
| -------------- | ------------------- | --------------------------------------------------------------- | ----- |
| **Content**    | Markdown            | Full Markdown with GFM extensions                               | v1.0  |
|                | Frontmatter         | YAML metadata for pages                                         | v1.0  |
|                | Admonitions         | Callout blocks (info, warning, danger, success, note, question) | v1.0  |
|                | Syntax Highlighting | 100+ languages via Chroma with copy button                      | v1.0  |
|                | Stale Warnings      | Configurable warnings for outdated content                      | v1.0  |
|                | Link Checking       | Validate internal/external links during build                   | v1.2  |
| **Navigation** | Auto Navigation     | From folder structure                                           | v1.0  |
|                | Custom TOC          | Via TOC.md file                                                 | v1.0  |
|                | Page TOC            | Auto-generated from headings                                    | v1.0  |
|                | Scrollspy           | Active heading highlight                                        | v1.0  |
| **Search**     | Full-Text           | Client-side search index                                        | v1.0  |
|                | Keyboard            | Cmd+K / Ctrl+K shortcut                                         | v1.0  |
|                | Fuzzy Matching      | Tolerant search                                                 | v1.0  |
| **Design**     | Responsive          | Mobile, tablet, desktop                                         | v1.0  |
|                | Dark Mode           | Toggle with localStorage persistence                            | v1.0  |
|                | Themes              | CSS-only customization (default, yellow)                        | v1.0  |
|                | Custom Themes       | Configure colors, fonts, backgrounds in YAML                    | v1.2  |
| **Pages**      | Landing             | Marketing homepage with hero, features, steps                   | v1.0  |
|                | MD Landing Sections | Define landing sections in `__landing__/*.md` files             | v1.1  |
|                | Portfolio           | Project showcase with tags and filtering                        | v1.0  |
|                | Contact             | Contact page with email and info                                | v1.0  |
|                | FAQ                 | Collapsible Q&A with categories and search integration          | v1.0  |
|                | Legal               | Privacy policy, terms of service with auto footer links         | v1.0  |
|                | Knowledge Base      | Self-service support hub with categories and scoped search      | v1.1  |
|                | Status              | Service health with incidents, maintenance, uptime              | v1.0  |
|                | Changelog           | Version history with RSS feed                                   | v1.0  |
|                | Roadmap             | Version-based timeline with tag filtering                       | v1.4  |
| **Versioning** | Multi-Version       | Documentation for multiple software versions                    | v1.2  |
|                | Version Selector    | Dropdown to switch between versions                             | v1.2  |
|                | EOL Warnings        | Warning banners for end-of-life versions                        | v1.2  |
| **SEO**        | Sitemap             | Automatic sitemap.xml generation                                | v1.0  |
|                | Meta Tags           | Open Graph, Twitter Cards                                       | v1.0  |
|                | Canonical URLs      | Configurable base URL                                           | v1.0  |
|                | LLM Output          | llms.txt and llms-full.txt generation                           | v1.0  |
| **API**        | OpenAPI             | Interactive docs with testing UI (OAS 3.0 and 3.1)              | v1.0  |
|                | Code Samples        | curl, JavaScript, Go, Python, Swift                             | v1.0  |
|                | Schema Viewer       | Request/response visualization                                  | v1.0  |
|                | MCP Server Docs     | Static reference pages from JSON manifests (no server needed)   | v1.5  |

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

Deploy anywhere.

### Build on Change

No web server required. Content updates only when you push changes. Your CI/CD
pipeline rebuilds the site, uploads static files, and you're done. No runtime
costs, no database, no server maintenance. Data is generated and consumed only
when needed.

```
git push → CI builds → Static files uploaded → Done
```

### CSS-Only Themes

Themes are CSS variables only. No template editing.

```css
:root[data-theme="light"] {
  --bg-primary: #fafafa;
  --text-primary: #1a1a1a;
}
```

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
