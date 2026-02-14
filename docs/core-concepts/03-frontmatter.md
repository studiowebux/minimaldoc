---
title: Frontmatter
description: YAML frontmatter schema for page metadata and configuration
tags:
  - core-concepts
  - frontmatter
---

# Frontmatter

YAML frontmatter at the top of Markdown files defines page metadata.

## Basic Format

```yaml
---
title: Page Title
description: Brief description for SEO
---

# Page Content

Markdown content here...
```

## Full Schema

```yaml
---
# Required
title: Page Title

# Recommended
description: Brief page description for SEO and search
tags:
  - tag1
  - tag2

# Optional metadata
author: Author Name
date: 2025-01-15

# Navigation
menu_title: Short Title      # Override title in navigation
menu_order: 10               # Sort order (lower = earlier)
hidden: false                # Hide from navigation

# Stale warning
stale_warning: true          # Override site setting
stale_threshold_days: 730    # Override threshold for this page

# SEO overrides
seo:
  title: SEO Title Override
  description: SEO description override
  keywords:
    - keyword1
    - keyword2
  image: /images/og-image.png
  canonical: https://example.com/canonical-url
  noindex: false
  nofollow: false
---
```

## Field Reference

### Core Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Page title, used in navigation and `<title>` |
| `description` | string | No | Page description for SEO and search results |
| `tags` | []string | No | Tags for categorization and search |
| `author` | string | No | Page author |
| `date` | string | No | Publication date (ISO 8601) |

### Navigation Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `menu_title` | string | `title` | Shorter title for navigation |
| `menu_order` | int | 0 | Sort order within section |
| `hidden` | bool | false | Hide from navigation |

### Stale Warning Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `stale_warning` | bool | site default | Enable/disable stale warning |
| `stale_threshold_days` | int | site default | Days before considered stale |

### SEO Fields

| Field | Type | Description |
|-------|------|-------------|
| `seo.title` | string | Override `<title>` tag |
| `seo.description` | string | Override meta description |
| `seo.keywords` | []string | Meta keywords |
| `seo.image` | string | Open Graph image URL |
| `seo.canonical` | string | Canonical URL |
| `seo.noindex` | bool | Add noindex meta |
| `seo.nofollow` | bool | Add nofollow meta |

### Version Fields

Used when multi-version documentation is enabled.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `versions` | []string | all | Versions this page appears in |
| `since` | string | - | Version where feature was introduced |
| `deprecated_in` | string | - | Version where feature was deprecated |
| `removed_in` | string | - | Version where feature was removed |
| `version_note` | string | - | Version to show special note for |

## Examples

### Minimal Page

```yaml
---
title: Quick Start
---
```

### Documentation Page

```yaml
---
title: API Authentication
description: How to authenticate with the API
tags:
  - api
  - authentication
  - security
author: API Team
---
```

### Hidden Page

```yaml
---
title: Internal Notes
description: Notes not shown in navigation
hidden: true
---
```

### Custom Navigation Order

```yaml
---
title: Introduction
menu_order: 1
---
```

```yaml
---
title: Getting Started
menu_order: 2
---
```

```yaml
---
title: Advanced Topics
menu_order: 99
---
```

### SEO Optimized

```yaml
---
title: MinimalDoc Documentation
description: Complete guide to MinimalDoc
seo:
  title: MinimalDoc - Static Site Generator for Docs
  description: Build beautiful documentation sites with MinimalDoc
  keywords:
    - documentation
    - static site generator
    - markdown
  image: /images/social-card.png
  canonical: https://minimaldoc.dev/docs
---
```

### Stable Content (No Stale Warning)

```yaml
---
title: License Agreement
description: Software license terms
stale_warning: false
---
```

### Extended Threshold

```yaml
---
title: Architecture Overview
description: System architecture that rarely changes
stale_threshold_days: 730
---
```

### Version-Specific Page

```yaml
---
title: New API Endpoint
description: Available in v2 and later
versions:
  - v2
  - v3
since: "v2.0"
---
```

### Deprecated Feature

```yaml
---
title: Legacy Authentication
description: Use OAuth2 instead
deprecated_in: "v2.0"
removed_in: "v3.0"
---
```

## Date Formats

Supported date formats:

```yaml
date: 2025-01-15
date: 2025-01-15T10:30:00Z
date: 2025-01-15T10:30:00+02:00
```

## Validation

MinimalDoc validates frontmatter and reports errors:

```
Error: docs/page.md: missing required field 'title'
Error: docs/other.md: invalid date format '15/01/2025'
```

Missing optional fields use defaults silently.
