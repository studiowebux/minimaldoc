---
title: Templates
description: Starter templates for documentation content
tags:
  - api-reference
  - templates
---

# Templates

Starter templates for common documentation patterns.

## Available Templates

| Template | Description | Location |
|----------|-------------|----------|
| `config.yaml` | Site configuration | `docs/config.yaml` |
| `page.md` | Documentation page | `docs/*.md` |
| `toc.md` | Custom navigation | `docs/TOC.md` |
| `portfolio-item.md` | Portfolio project | `docs/__portfolio__/*.md` |
| `faq-item.md` | FAQ question/answer | `docs/__faq__/category/*.md` |
| `kb-article.md` | Knowledge base article | `docs/__kb__/category/*.md` |
| `legal-page.md` | Legal page | `docs/__legal__/*.md` |
| `openapi-page.md` | OpenAPI endpoint page | `docs/api/*.md` |
| `version-page.md` | Versioned documentation | `docs/__versions__/v*/*.md` |
| `changelog-release.md` | Individual release | `docs/__changelog__/releases/*.md` |
| `components.yaml` | Status components | `docs/__status__/components.yaml` |
| `incident.md` | Incident report | `docs/__status__/incidents/*.md` |
| `maintenance.md` | Maintenance notice | `docs/__status__/maintenance/*.md` |

## Page Template

```markdown
---
title: Page Title
description: Brief description for SEO and search
tags:
  - tag1
  - tag2
author: Author Name
---

# Page Title

Introduction paragraph explaining the page content.

## Section One

Content for section one.

### Subsection

More detailed content.

## Section Two

Content for section two.

## Next Steps

- [Related Page](related.md)
- [Another Page](another.md)
```

## TOC.md Template

```markdown
# My Documentation

- [GitHub](https://github.com/org/repo)

- [Home](index.md)
- Getting Started
  - [Installation](getting-started/installation.md)
  - [Quick Start](getting-started/quick-start.md)
  - [Configuration](getting-started/configuration.md)
- Features
  - [Overview](features/overview.md)
  - [Feature A](features/feature-a.md)
  - [Feature B](features/feature-b.md)
- Guides
  - [Deployment](guides/deployment.md)
  - [Customization](guides/customization.md)
- API Reference
  - [REST API](api/rest.md)
  - [SDK](api/sdk.md)
- Resources
  - [GitHub](https://github.com/org/repo)
  - [Support](https://support.example.com)
```

## config.yaml Template

```yaml
# Site Information
title: My Documentation
description: Documentation for my project
base_url: https://docs.example.com
author: Documentation Team

# Theme
theme: default
dark_mode: false

# Features
enable_llms: true
enable_search: true
clean_urls: false

# OpenAPI (optional)
openapi:
  enabled: false
  spec_files:
    - "api/openapi.yaml"
  spec_urls: []
  default_view: "path"
  enable_testing: true
  enable_code_samples: true

# Status Page (optional)
status:
  enabled: false
  title: "Service Status"
  path: "status"
  show_history: true
  rss_enabled: true

# Changelog (optional)
changelog:
  enabled: false
  title: "Changelog"
  path: "changelog"
  repository: "https://github.com/org/repo"

# Stale Warning
stale_warning:
  enabled: true
  threshold_days: 365
  show_update_date: true

# Landing Page (optional)
landing:
  enabled: false
  nav:
    - text: "Docs"
      url: "/getting-started/installation.html"
  hero:
    title: "Documentation Made Simple"
    subtitle: "A minimal static site generator"
    buttons:
      - text: "Get Started"
        url: "/getting-started/installation.html"
        primary: true
  features:
    title: "Features"
    items:
      - emoji: "~"
        title: "Fast"
        description: "Sub-second builds"

# Portfolio (optional)
portfolio:
  enabled: false
  title: "Projects"
  description: "Showcase of projects"
  path: "projects"

# Contact (optional)
contact:
  enabled: false
  title: "Contact"
  description: "Get in touch"
  path: "contact"
  email: "hello@example.com"
  info:
    - icon: "mail"
      text: "hello@example.com"
    - icon: "location"
      text: "Remote"

# Footer
footer:
  copyright: "2025 Your Company"
  links:
    - title: "Documentation"
      items:
        - text: "Getting Started"
          url: "/getting-started/installation.html"
  social:
    - name: GitHub
      url: https://github.com/org/repo
      icon: github

# Social Links (sidebar)
social_links:
  - name: GitHub
    url: https://github.com/org/repo
    icon: github
```

## Status Components Template

```yaml
# docs/__status__/components.yaml

# Core Services
- id: api
  name: API
  description: Core REST API
  status: operational
  group: Core Services
  order: 1
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30
  uptime:
    mode: incidents
    sla_target: 99.9
    period_days: 90

- id: web
  name: Web Application
  description: Main web interface
  status: operational
  group: Core Services
  order: 2
  url: https://app.example.com
  health_endpoint: /health
  health_interval: 30

# Infrastructure
- id: database
  name: Database
  description: PostgreSQL cluster
  status: operational
  group: Infrastructure
  order: 1
  uptime:
    mode: incidents
    sla_target: 99.99

- id: cache
  name: Cache
  description: Redis cache cluster
  status: operational
  group: Infrastructure
  order: 2
```

## Incident Template

```markdown
---
title: Brief Description of Issue
status: investigating
severity: minor
affected_components:
  - api
  - web
created_at: 2025-01-28T10:00:00Z
---

## Update - 10:00 UTC

Initial report of the issue. We are investigating.

Describe what is known so far.
```

### Incident Status Values

| Status | When to Use |
|--------|-------------|
| `investigating` | Initial investigation |
| `identified` | Root cause found |
| `monitoring` | Fix applied, monitoring |
| `resolved` | Fully resolved |

### Incident Severity Values

| Severity | Description |
|----------|-------------|
| `minor` | Low impact |
| `major` | Significant impact |
| `critical` | Service down |

## Maintenance Template

```markdown
---
title: Maintenance Title
description: Brief description of maintenance
affected_components:
  - database
  - api
scheduled_start: 2025-02-01T02:00:00Z
scheduled_end: 2025-02-01T04:00:00Z
status: scheduled
---

Description of the maintenance and its purpose.

## Expected Impact

- Service A will be unavailable
- Service B will have reduced performance
- Expected duration: 2 hours

## Preparation

Instructions for users, if any.
```

### Maintenance Status Values

| Status | When to Use |
|--------|-------------|
| `scheduled` | Planned, not started |
| `in_progress` | Currently underway |
| `completed` | Finished |

## Portfolio Project Template

```markdown
---
title: Project Name
description: Brief project description
tags:
  - web
  - documentation
image: /images/project-screenshot.png
links:
  - text: Live Demo
    url: https://example.com
  - text: Source Code
    url: https://github.com/user/project
date: 2025-01-15
featured: true
order: 1
---

Project description with details about the implementation,
technologies used, and outcomes achieved.
```

### Portfolio Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Project name |
| `description` | string | No | Brief description |
| `image` | string | No | Thumbnail image URL or path |
| `tags` | array | No | Filterable tags |
| `links` | array | No | Array of `{text, url}` objects |
| `date` | string | No | Project date (ISO 8601) |
| `featured` | boolean | No | Show in featured section |
| `order` | integer | No | Display order |

## FAQ Item Template

```markdown
---
question: How do I configure the site?
category: Getting Started
order: 1
tags:
  - configuration
---

Edit `config.yaml` in your docs directory.
```

### FAQ Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `question` | string | Yes | The FAQ question text |
| `category` | string | Yes | Category name for grouping |
| `order` | integer | No | Sort order within category |
| `tags` | array | No | Tags for search |

## Knowledge Base Article Template

```markdown
---
title: Article Title
description: Brief summary of what this article covers
tags:
  - troubleshooting
  - setup
order: 1
---

## Problem

Describe the problem or topic.

## Solution

Step-by-step solution or explanation.
```

### Knowledge Base Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Article title |
| `description` | string | No | Summary shown in article cards |
| `tags` | array | No | Tags for search and related articles |
| `order` | integer | No | Sort order within category |

Category is determined by the directory name under `__kb__/`.

## Legal Page Template

```markdown
---
title: Privacy Policy
description: How we handle your data
order: 1
---

**Effective date:** YYYY-MM-DD

## Information We Collect

Describe what data is collected.
```

### Legal Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Page title |
| `description` | string | No | Brief description |
| `order` | integer | No | Sort order in footer links |

## OpenAPI Endpoint Page Template

```markdown
---
title: Create User
description: API endpoint for creating a new user
openapi_spec: api/openapi.yaml
openapi_path: /users
openapi_method: POST
---

Additional context or usage notes for this endpoint.
```

### OpenAPI Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `openapi_spec` | string | Yes | Path to OpenAPI spec file or URL |
| `openapi_path` | string | Yes | API endpoint path (e.g., `/users`) |
| `openapi_method` | string | Yes | HTTP method (`GET`, `POST`, etc.) |

Without `openapi_path` and `openapi_method`, the page embeds the full spec.

## Versioned Page Template

```markdown
---
title: Feature Name
description: Feature available since v2.0
versions:
  - v2
  - v3
since: "v2.0"
# deprecated_in: "v3.0"
# removed_in: "v4.0"
# version_note: Important version-specific information
---

Feature documentation content.
```

### Version Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `versions` | array | No | Versions this page appears in (empty = all) |
| `since` | string | No | Version where feature was introduced |
| `deprecated_in` | string | No | Version where feature was deprecated |
| `removed_in` | string | No | Version where feature was removed |
| `version_note` | string | No | Note displayed for version context |

## Changelog Release Template

```markdown
---
version: 1.0.0
date: 2025-01-28
title: Optional Release Title
---

## Added

- New feature A
- New feature B

## Changed

- Updated feature C
- Improved performance of D

## Fixed

- Fixed bug in E
- Corrected behavior of F

## Security

- Patched vulnerability in G
```

### Changelog Categories

| Category | Description | Color |
|----------|-------------|-------|
| `Added` | New features | Green |
| `Changed` | Changes to existing features | Blue |
| `Deprecated` | Features to be removed | Yellow |
| `Removed` | Removed features | Red |
| `Fixed` | Bug fixes | Purple |
| `Security` | Security fixes | Orange |

## Directory Structure

Complete documentation structure:

```
docs/
├── config.yaml               # All configuration
├── TOC.md                    # Custom navigation (optional)
├── index.md                  # Homepage
├── getting-started/
│   └── *.md
├── features/
│   └── *.md
├── api/
│   ├── openapi.yaml          # OpenAPI specs (if enabled)
│   └── endpoint.md           # Per-endpoint pages (optional)
├── __portfolio__/            # Portfolio (if enabled)
│   └── project-name.md
├── __faq__/                  # FAQ (if enabled)
│   └── category-name/
│       └── question.md
├── __kb__/                   # Knowledge base (if enabled)
│   └── category-name/
│       └── article.md
├── __legal__/                # Legal pages (if enabled)
│   └── privacy-policy.md
├── __changelog__/            # Per-release changelog (if enabled)
│   └── releases/
│       └── vX.Y.Z.md
├── __versions__/             # Versioned docs (if enabled)
│   ├── v1/
│   │   └── *.md
│   └── v2/
│       └── *.md
└── __status__/               # Status page (if enabled)
    ├── components.yaml
    ├── incidents/
    │   └── YYYY-MM-DD-*.md
    └── maintenance/
        └── YYYY-MM-DD-*.md
```
