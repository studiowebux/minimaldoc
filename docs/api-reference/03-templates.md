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
| `components.yaml` | Status components | `docs/__status__/components.yaml` |
| `status-config.yaml` | Status config | `docs/__status__/config.yaml` |
| `incident.md` | Incident report | `docs/__status__/incidents/*.md` |
| `maintenance.md` | Maintenance notice | `docs/__status__/maintenance/*.md` |
| `changelog.md` | Release notes | `docs/__changelog__/releases/*.md` |

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

# Social Links
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

## Status Config Template

```yaml
# docs/__status__/config.yaml

enabled: true
title: "System Status"
description: "Current operational status of our services"
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
├── config.yaml
├── TOC.md
├── index.md
├── getting-started/
│   ├── 01-installation.md
│   ├── 02-quick-start.md
│   └── 03-configuration.md
├── features/
│   └── *.md
├── guides/
│   └── *.md
├── api/
│   ├── reference.md
│   └── openapi.yaml
├── __status__/
│   ├── config.yaml
│   ├── components.yaml
│   ├── incidents/
│   │   └── YYYY-MM-DD-*.md
│   └── maintenance/
│       └── YYYY-MM-DD-*.md
└── __changelog__/
    ├── config.yaml
    └── releases/
        └── *.md
```

## Usage

Copy templates to your docs directory:

```bash
# Copy config template
cp templates/config.yaml docs/

# Copy page template
cp templates/page.md docs/my-new-page.md

# Copy status templates
mkdir -p docs/__status__/{incidents,maintenance}
cp templates/components.yaml docs/__status__/
cp templates/status-config.yaml docs/__status__/config.yaml
```
