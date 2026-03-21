---
title: Knowledge Base
description: Categorized help articles with search and structured navigation
tags:
  - guides
  - knowledge-base
---

# Knowledge Base

A categorized article system with search, icons, and structured navigation.

## Configuration

```yaml
knowledgebase:
  enabled: true
  title: Knowledge Base
  description: Find answers and guides
  path: kb
  search:
    enabled: true
    placeholder: "Search articles..."
  categories:
    getting-started:
      name: Getting Started
      description: Installation and setup
      icon: rocket
      order: 1
    features:
      name: Features
      description: Search, themes, code blocks
      icon: code
      order: 2
    troubleshooting:
      name: Troubleshooting
      description: Common issues and solutions
      icon: wrench
      order: 3
```

## Article Files

Create markdown files in `__kb__/<category>/` directories:

```
docs/
  __kb__/
    getting-started/
      01-installation.md
      02-configuration.md
    features/
      01-search.md
      02-dark-mode.md
    troubleshooting/
      01-build-errors.md
```

The directory name must match a key in the `categories` map.

### Article Frontmatter

```yaml
---
title: How to Configure Search
description: Enable and customize client-side search
tags:
  - search
  - configuration
order: 1
---
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | yes | Article title |
| `description` | string | no | Short description shown on listing |
| `tags` | []string | no | Article tags for search |
| `order` | int | no | Sort order within category |

## Category Configuration

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Display name |
| `description` | string | Category description shown on landing |
| `icon` | string | Icon identifier (rocket, code, wrench, book, etc.) |
| `order` | int | Sort order on landing page |

## KB Search

When `search.enabled` is true, a search input appears on the KB landing page. It searches across all article titles, descriptions, and tags.

## Output

- Landing page: `/<path>/index.html`
- Category page: `/<path>/<category>/index.html`
- Article page: `/<path>/<category>/<article>.html`
