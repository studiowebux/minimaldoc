---
title: Roadmap Pages
description: Add a visual roadmap to your documentation site
tags:
  - guides
  - roadmap
---

# Roadmap Pages

Display a public roadmap with board or timeline layout, version badges, and tag filtering.

## Configuration

```yaml
roadmap:
  enabled: true
  title: Roadmap
  description: What we're building next
  layout: board
  show_versions: true
  path: roadmap
  columns:
    - id: planned
      label: Planned
    - id: in_progress
      label: In Progress
    - id: shipped
      label: Shipped
  items:
    - title: Dark mode
      description: System-aware theme switching with manual override
      status: shipped
      version: v1.0
      tags: [ui, theme]
      shipped_date: "2025-01-15"
      changelog_url: "/changelog/#v1.0"
    - title: Plugin system
      description: Extensible plugin architecture for custom page types
      status: planned
      version: v2.0
      tags: [architecture]
```

## Layouts

### Board

Kanban-style columns. Each column maps to a `columns[].id`.

```yaml
layout: board
columns:
  - id: planned
    label: Planned
  - id: in_progress
    label: In Progress
  - id: shipped
    label: Shipped
```

Items appear in the column matching their `status` field.

### Timeline

Chronological view grouped by version. Enable with:

```yaml
layout: timeline
show_versions: true
```

## Item Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | yes | Item title |
| `description` | string | no | Short description |
| `status` | string | yes | Must match a `columns[].id` |
| `version` | string | no | Version badge (e.g., `v1.0`) |
| `tags` | []string | no | Filterable tags |
| `shipped_date` | string | no | Ship date (YYYY-MM-DD), shown for shipped items |
| `changelog_url` | string | no | Link to changelog entry |

## Tag Filtering

Tags are rendered as clickable filters above the board. Clicking a tag shows only items with that tag. Tags are extracted automatically from all items.

## Linking to Changelog

Use `changelog_url` to link shipped items to their release notes:

```yaml
- title: Status pages
  status: shipped
  version: v1.1
  changelog_url: "/changelog/#v1.1"
```

## Custom Columns

Columns are fully configurable. Common patterns:

**Simple three-column:**

```yaml
columns:
  - id: planned
    label: Planned
  - id: in_progress
    label: In Progress
  - id: shipped
    label: Shipped
```

**With backlog:**

```yaml
columns:
  - id: backlog
    label: Backlog
  - id: next
    label: Up Next
  - id: in_progress
    label: Building
  - id: shipped
    label: Shipped
```

## Output

The roadmap page is generated at `/<path>/index.html` (default: `/roadmap/index.html`). A link appears automatically in the site navigation.
