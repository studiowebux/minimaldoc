---
title: Changelog Overview
description: Track releases with versioned changelog
tags:
  - changelog
  - releases
---

# Changelog Overview

MinimalDoc includes a changelog system for documenting releases.

## Features

| Feature            | Description                    |
| ------------------ | ------------------------------ |
| Versioned Releases | Semver-sorted release notes    |
| Categories         | Added, Changed, Fixed, etc.    |
| Compare Links      | Link to diffs between versions |
| RSS Feed           | Subscribe to releases          |
| JSON API           | Machine-readable changelog     |

## Quick Start

1. Create changelog directory:

```
docs/
├── config.yaml           # Changelog settings here
└── __changelog__/
    └── releases/
        ├── 1.0.0.md
        └── 0.1.0.md
```

2. Enable in config:

```yaml
changelog:
  enabled: true
  title: "Changelog"
  path: "changelog"
  repository: "https://github.com/org/repo"
```

3. Build:

```bash
minimaldoc build --changelog
```

## Directory Structure

```
docs/
├── config.yaml               # Changelog settings
└── __changelog__/
    └── releases/
        ├── 2.0.0.md          # Newest
        ├── 1.1.0.md
        ├── 1.0.0.md
        └── 0.1.0.md          # Oldest
```

Releases sorted by semver, newest first.

## Generated Output

```
public/
└── changelog/
    ├── index.html           # All releases
    ├── 1.0.0.html           # Individual release
    ├── feed.xml             # RSS feed
    └── changelog.json       # JSON API
```

## Configuration

### Site Config

```yaml
changelog:
  enabled: true
  title: "Changelog"
  path: "changelog"
  rss_enabled: true
  repository: "https://github.com/org/repo"
```

| Option        | Default       | Description                      |
| ------------- | ------------- | -------------------------------- |
| `enabled`     | `false`       | Enable changelog                 |
| `title`       | `"Changelog"` | Page title                       |
| `path`        | `"changelog"` | Output path                      |
| `rss_enabled` | `true`        | Generate RSS feed                |
| `repository`  | `""`          | Repository URL for compare links |

## Compare Links

With repository configured, each release shows diff link:

```
v1.1.0                              [Compare]
```

Links to:

```
https://github.com/org/repo/compare/v1.0.0...v1.1.0
```

## RSS Feed

Subscribe at `/changelog/feed.xml`:

```xml
<rss version="2.0">
  <channel>
    <title>Changelog</title>
    <item>
      <title>v1.1.0</title>
      <pubDate>Tue, 28 Jan 2025 00:00:00 GMT</pubDate>
      <description>Added: Feature A, Feature B...</description>
    </item>
  </channel>
</rss>
```

## JSON API

Machine-readable at `/changelog/changelog.json`:

```json
{
  "releases": [
    {
      "version": "1.1.0",
      "date": "2025-01-28",
      "title": "Feature Update",
      "url": "/changelog/1.1.0.html",
      "categories": {
        "added": ["Feature A", "Feature B"],
        "fixed": ["Bug fix #123"]
      }
    }
  ]
}
```

## Keep a Changelog Format

Based on [keepachangelog.com](https://keepachangelog.com):

- Changelogs are for humans, not machines
- Every version should have an entry
- Group changes by type
- Versions should be linkable
- Latest version comes first
- Release date for each version

## Next Steps

- [Format](02-format.md) - Release file format and categories
