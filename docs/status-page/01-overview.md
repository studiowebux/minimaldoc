---
title: Status Page Overview
description: Service health dashboard with incidents and uptime tracking
tags:
  - status-page
  - overview
---

# Status Page Overview

MinimalDoc includes a standalone status page for service health monitoring.

## Features

| Feature | Description |
|---------|-------------|
| Component Status | Track multiple services |
| Incidents | Document and track issues |
| Maintenance | Schedule downtime windows |
| Uptime Calendar | 90-day visual grid |
| SLA Display | Track uptime percentages |
| Health Checks | Live endpoint polling |
| RSS Feed | Subscribe to updates |
| JSON API | Machine-readable status |

## Quick Start

1. Create status directory:

```
docs/
└── __status__/
    ├── config.yaml
    ├── components.yaml
    ├── incidents/
    └── maintenance/
```

2. Enable in config:

```yaml
status:
  enabled: true
  title: "Service Status"
  path: "status"
```

3. Build:

```bash
minimaldoc build --status
```

## Directory Structure

```
docs/
└── __status__/
    ├── config.yaml           # Status page settings
    ├── components.yaml       # Service definitions
    ├── incidents/
    │   ├── 2025-01-28-api-outage.md
    │   └── 2025-01-15-slow-response.md
    └── maintenance/
        └── 2025-02-01-upgrade.md
```

The `__status__` directory name avoids conflicts with documentation content.

## Generated Output

```
public/
└── status/
    ├── index.html           # Main dashboard
    ├── history.html         # Incident history
    ├── incident/
    │   └── *.html           # Individual incidents
    ├── status.json          # JSON API
    └── feed.xml             # RSS feed
```

## Status Dashboard

The main dashboard shows:

```
┌─────────────────────────────────────────────────┐
│  System Status                                   │
│  All systems operational                    ✓    │
├─────────────────────────────────────────────────┤
│  Core Services                                   │
│  ├── API                              ✓ 24ms    │
│  ├── Web App                          ✓ 18ms    │
│  └── Database                         ✓ 12ms    │
│                                                  │
│  Infrastructure                                  │
│  ├── CDN                              ✓         │
│  └── Storage                          ✓         │
├─────────────────────────────────────────────────┤
│  Recent Incidents                                │
│  └── No recent incidents                         │
├─────────────────────────────────────────────────┤
│  Scheduled Maintenance                           │
│  └── Database Upgrade - Feb 1, 02:00 UTC        │
└─────────────────────────────────────────────────┘
```

## Sidebar Link

When enabled, a status link appears in the documentation sidebar:

```
- Documentation
  - Getting Started
  - Features
- System Status ●       # Green dot when operational
```

## Status Values

| Status | Color | Description |
|--------|-------|-------------|
| `operational` | Green | Fully functional |
| `degraded` | Yellow | Reduced performance |
| `partial_outage` | Orange | Some features unavailable |
| `major_outage` | Red | Service unavailable |
| `maintenance` | Blue | Planned maintenance |

## Overall Status

Calculated from component statuses:

| Components | Overall Status |
|------------|----------------|
| All operational | Operational |
| Any degraded | Degraded Performance |
| Any partial outage | Partial Outage |
| Any major outage | Major Outage |
| Any in maintenance | Under Maintenance |

## Configuration Options

```yaml
status:
  enabled: true
  title: "Service Status"
  description: "Current status of our services"
  path: "status"
  show_history: true
  history_months: 12
  rss_enabled: true
```

| Option | Default | Description |
|--------|---------|-------------|
| `enabled` | `false` | Enable status page |
| `title` | `"Service Status"` | Page title |
| `description` | `""` | Page description |
| `path` | `"status"` | Output path |
| `show_history` | `true` | Show incident history |
| `history_months` | `12` | Months of history |
| `rss_enabled` | `true` | Generate RSS feed |

## Next Steps

- [Components](02-components.md) - Define your services
- [Incidents](03-incidents.md) - Document issues
- [Maintenance](04-maintenance.md) - Schedule downtime
- [Uptime Tracking](05-uptime-tracking.md) - Track availability
