---
title: Search
description: Client-side full-text search with keyboard shortcuts
tags:
  - features
  - search
---

# Search

MinimalDoc includes client-side full-text search with fuzzy matching.

## Usage

### Keyboard Shortcut

Press `Cmd+K` (macOS) or `Ctrl+K` (Windows/Linux) to open search.

### Click

Click the search input in the sidebar.

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate results |
| `Enter` | Go to selected result |
| `Escape` | Close search |

## Search Index

A `search-index.json` file is generated during build:

```json
[
  {
    "title": "Installation",
    "description": "How to install MinimalDoc",
    "url": "/getting-started/installation.html",
    "tags": ["getting-started", "installation"],
    "content": "First 500 characters of content..."
  }
]
```

### Indexed Fields

| Field | Weight | Source |
|-------|--------|--------|
| Title | Highest | Frontmatter `title` |
| Description | High | Frontmatter `description` |
| Tags | Medium | Frontmatter `tags` |
| Content | Lower | Page body (first 500 chars) |

## Configuration

Search is enabled by default:

```yaml
# config.yaml
enable_search: true
```

Disable search:

```yaml
enable_search: false
```

## Fuzzy Matching

Search tolerates typos and partial matches:

| Query | Matches |
|-------|---------|
| `instal` | "Installation" |
| `confg` | "Configuration" |
| `api ref` | "API Reference" |

## Index Size

Typical index sizes:

| Pages | Index Size |
|-------|------------|
| 10 | ~5 KB |
| 50 | ~25 KB |
| 100 | ~50 KB |
| 500 | ~250 KB |

## Optimization Tips

### Good Titles

```yaml
# Good - Specific, searchable
title: Configure OAuth 2.0 Authentication

# Avoid - Generic, hard to find
title: Configuration
```

### Descriptive Tags

```yaml
tags:
  - authentication
  - oauth
  - security
  - api
```

### Front-Load Content

Put important keywords in the first 500 characters, as that's what gets indexed.

### Unique Descriptions

```yaml
# Good - Unique, descriptive
description: Set up OAuth 2.0 authentication with Google, GitHub, or custom providers

# Avoid - Generic
description: Authentication documentation
```

## Excluding Pages

Hidden pages are excluded from search:

```yaml
---
title: Internal Notes
hidden: true
---
```

## Search UI

The search interface includes:

- Search input with placeholder
- Real-time results as you type
- Result highlighting
- Keyboard navigation indicators
- Result count
- "No results" state

## Technical Details

### Algorithm

1. Tokenize query into words
2. For each word, find fuzzy matches in index
3. Score matches by field weight and match quality
4. Sort by combined score
5. Return top N results

### Performance

- Index loaded once, cached in memory
- Search executes in <10ms for most sites
- No server round-trips
