---
title: Project Structure
description: Understanding MinimalDoc directory layout and file organization
tags:
  - core-concepts
  - structure
---

# Project Structure

## Directory Layout

```
docs/
├── config.yaml              # Site configuration
├── TOC.md                   # Custom navigation (optional)
├── index.md                 # Homepage
├── getting-started/         # Section directory
│   ├── 01-installation.md   # Ordered page
│   └── 02-quick-start.md
├── guides/
│   └── deployment.md
├── api/                     # OpenAPI specs (if enabled)
│   └── openapi.yaml
├── __status__/              # Status page content
│   ├── config.yaml
│   ├── components.yaml
│   ├── incidents/
│   └── maintenance/
└── __changelog__/           # Changelog content
    ├── config.yaml
    └── releases/
```

## Reserved Directories

| Directory | Purpose |
|-----------|---------|
| `__status__` | Status page content (incidents, components) |
| `__changelog__` | Changelog releases |
| `api/` | OpenAPI specification files (configurable) |

These directories are processed specially and don't appear in regular navigation.

## File Types

| Extension | Processing |
|-----------|------------|
| `.md` | Markdown pages, converted to HTML |
| `.yaml`, `.json` | Configuration and data files |
| Other | Copied to output as static assets |

## Build Output

```
public/
├── index.html               # Generated HTML pages
├── getting-started/
│   ├── installation.html
│   └── quick-start.html
├── guides/
│   └── deployment.html
├── css/                     # Theme styles
│   └── main.css
├── js/                      # JavaScript assets
│   ├── search.js
│   └── theme-toggle.js
├── search-index.json        # Search index
├── sitemap.xml              # SEO sitemap
├── llms.txt                 # LLM-friendly content
├── status/                  # Status page (if enabled)
│   ├── index.html
│   ├── history.html
│   ├── status.json
│   └── feed.xml
└── changelog/               # Changelog (if enabled)
    ├── index.html
    └── feed.xml
```

## Naming Conventions

### Numeric Prefixes

Use numeric prefixes for ordering:

```
01-introduction.md    # Displays as "Introduction"
02-installation.md    # Displays as "Installation"
10-advanced.md        # Displays as "Advanced"
```

Prefixes are stripped from display titles and URLs:
- `01-getting-started/` → `/getting-started/`
- `02-quick-start.md` → `/quick-start.html`

### Section Directories

Directories become navigation sections:

```
getting-started/      # Section: "Getting Started"
  01-intro.md         # Page in section
  02-install.md
api-reference/        # Section: "API Reference"
  endpoints.md
```

## Static Assets

Place static files anywhere in your docs directory:

```
docs/
├── images/
│   └── logo.png
├── downloads/
│   └── example.zip
└── pages/
    └── about.md
```

Reference in markdown:

```markdown
![Logo](/images/logo.png)
[Download](/downloads/example.zip)
```

Assets are copied to output preserving directory structure.

## Hidden Files

Files and directories starting with `.` or `_` are ignored:

```
docs/
├── _drafts/          # Ignored
│   └── wip.md
├── .private/         # Ignored
└── published/        # Processed
    └── page.md
```

Exception: `__status__` and `__changelog__` use double underscores intentionally.
