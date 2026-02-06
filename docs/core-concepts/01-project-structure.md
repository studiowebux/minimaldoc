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
├── config.yaml              # All site configuration
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
│   ├── components.yaml
│   ├── incidents/
│   └── maintenance/
├── __changelog__/           # Changelog content
│   └── releases/
├── __portfolio__/           # Portfolio projects
│   └── project-name.md
├── __faq__/                 # FAQ markdown files (optional)
│   ├── general/
│   │   └── what-is-this.md
│   └── technical/
│       └── supported-languages.md
└── __legal__/               # Legal pages
    ├── privacy.md
    └── terms.md
```

## Reserved Directories

| Directory | Purpose |
|-----------|---------|
| `__status__` | Status page (components, incidents, maintenance) |
| `__changelog__` | Changelog releases |
| `__portfolio__` | Portfolio project pages |
| `__faq__` | FAQ questions as markdown files |
| `__legal__` | Legal pages (privacy, terms, etc.) |
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
├── index.html               # Landing or docs homepage
├── getting-started/
│   ├── installation.html
│   └── quick-start.html
├── css/
│   └── main.css
├── js/
│   ├── search.js
│   └── theme-toggle.js
├── search-index.json        # Search index
├── sitemap.xml              # SEO sitemap
├── llms.txt                 # LLM-friendly content
├── api/                     # OpenAPI docs (if enabled)
│   └── index.html
├── status/                  # Status page (if enabled)
│   ├── index.html
│   ├── status.json
│   └── feed.xml
├── changelog/               # Changelog (if enabled)
│   ├── index.html
│   └── feed.xml
├── projects/                # Portfolio (if enabled)
│   └── index.html
├── contact/                 # Contact page (if enabled)
│   └── index.html
├── faq/                     # FAQ page (if enabled)
│   └── index.html
└── legal/                   # Legal pages (if enabled)
    ├── privacy/
    │   └── index.html
    └── terms/
        └── index.html
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

Exception: `__status__`, `__changelog__`, `__portfolio__`, `__faq__`, and `__legal__` use double underscores intentionally.
