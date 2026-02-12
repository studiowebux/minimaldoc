---
title: Features Overview
description: Complete feature list and capabilities of MinimalDoc
tags:
  - features
  - overview
---

# Features Overview

MinimalDoc provides everything needed for modern documentation sites.
Pure static output designed for CI/CD integration—automate updates to documentation, status pages, and changelogs on every commit.

## Feature Matrix

| Category       | Feature             | Description                                                     |
| -------------- | ------------------- | --------------------------------------------------------------- |
| **Content**    | Markdown            | Full Markdown with GFM extensions                               |
|                | Frontmatter         | YAML metadata for pages                                         |
|                | Admonitions         | Callout blocks (info, warning, danger, success, note, question) |
|                | Syntax Highlighting | 100+ languages via Chroma with copy button                      |
|                | Stale Warnings      | Configurable warnings for outdated content                      |
|                | Link Checking       | Validate internal/external links during build                   |
| **Navigation** | Auto Navigation     | From folder structure                                           |
|                | Custom TOC          | Via TOC.md file                                                 |
|                | Page TOC            | Auto-generated from headings                                    |
|                | Scrollspy           | Active heading highlight                                        |
| **Search**     | Full-Text           | Client-side search index                                        |
|                | Keyboard            | Cmd+K / Ctrl+K shortcut                                         |
|                | Fuzzy Matching      | Tolerant search                                                 |
| **Design**     | Responsive          | Mobile, tablet, desktop                                         |
|                | Dark Mode           | Toggle with localStorage persistence                            |
|                | Themes              | CSS-only customization (default, yellow)                        |
| **Pages**      | Landing             | Marketing homepage with hero, features, steps                   |
|                | Portfolio           | Project showcase with tags and filtering                        |
|                | Contact             | Contact page with email and info                                |
|                | FAQ                 | Collapsible Q&A with categories and search integration          |
|                | Legal               | Privacy policy, terms of service with auto footer links         |
|                | Knowledge Base      | Self-service support hub with categories and scoped search      |
|                | Status              | Service health with incidents, maintenance, uptime              |
|                | Changelog           | Version history with RSS feed                                   |
| **SEO**        | Sitemap             | Automatic sitemap.xml generation                                |
|                | Meta Tags           | Open Graph, Twitter Cards                                       |
|                | Canonical URLs      | Configurable base URL                                           |
|                | LLM Output          | llms.txt and llms-full.txt generation                           |
| **API**        | OpenAPI             | Interactive docs with testing UI                                |
|                | Code Samples        | curl, JavaScript, Go, Python, Swift                             |
|                | Schema Viewer       | Request/response visualization                                  |

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

No web server required. Content updates only when you push changes. Your CI/CD pipeline rebuilds the site, uploads static files, and you're done. No runtime costs, no database, no server maintenance. Data is generated and consumed only when needed.

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
