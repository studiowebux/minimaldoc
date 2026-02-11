---
title: "Search Functionality"
description: "How to use and configure the built-in search"
tags: ["search", "navigation", "keyboard"]
---

## Overview

MinimalDoc includes a client-side search that indexes all your documentation pages, FAQ items, and Knowledge Base articles.

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `/` or `Ctrl+K` | Open search |
| `Escape` | Close search |
| `Arrow keys` | Navigate results |
| `Enter` | Go to selected result |

## Enabling Search

Search is enabled by default. To explicitly configure it:

```yaml
enable_search: true
```

## How It Works

During build, MinimalDoc generates a `search-index.json` file containing:

- Page titles and descriptions
- Section headings
- Content snippets
- Tags and metadata

The search runs entirely in the browser with no server required.

## Search Weighting

Results are ranked by relevance:

- **Title matches**: Highest priority
- **Description matches**: Medium priority
- **Tag matches**: Medium priority
- **Content matches**: Lower priority

## Tips for Better Search

1. Use descriptive page titles
2. Add `description` in frontmatter
3. Use relevant `tags` for categorization
4. Structure content with clear headings
