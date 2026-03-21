---
title: Components
description: Define and configure service components
tags:
  - status-page
  - components
---

# Components

Components represent services displayed on your status page.

## Basic Configuration

Create `__status__/components.yaml`:

```yaml
- id: api
  name: API
  status: operational

- id: web
  name: Web Application
  status: operational

- id: database
  name: Database
  status: operational
```

## Full Configuration

```yaml
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
```

## Field Reference

| Field             | Type   | Required | Description                             |
| ----------------- | ------ | -------- | --------------------------------------- |
| `id`              | string | Yes      | Unique identifier                       |
| `name`            | string | Yes      | Display name                            |
| `description`     | string | No       | Component description                   |
| `status`          | string | No       | Current status (default: `operational`) |
| `group`           | string | No       | Group for organization                  |
| `order`           | int    | No       | Sort order within group                 |
| `url`             | string | No       | Service URL                             |
| `health_endpoint` | string | No       | Health check path                       |
| `health_interval` | int    | No       | Check interval (seconds)                |
| `uptime`          | object | No       | Uptime tracking config                  |

## Status Values

| Status           | Display         |
| ---------------- | --------------- |
| `operational`    | Green checkmark |
| `degraded`       | Yellow warning  |
| `partial_outage` | Orange alert    |
| `major_outage`   | Red X           |
| `maintenance`    | Blue wrench     |

## Grouping

Organize components into groups:

```yaml
- id: api
  name: API
  group: Core Services
  order: 1

- id: web
  name: Web App
  group: Core Services
  order: 2

- id: cdn
  name: CDN
  group: Infrastructure
  order: 1

- id: storage
  name: Storage
  group: Infrastructure
  order: 2
```

Renders as:

```
Core Services
├── API            ✓
└── Web App        ✓

Infrastructure
├── CDN            ✓
└── Storage        ✓
```

## Ordering

Components sorted by:

1. Group name (alphabetical)
2. `order` field (ascending)
3. Component name (alphabetical)

```yaml
# Explicit ordering
- id: api
  name: API
  group: Core
  order: 1 # First

- id: web
  name: Web
  group: Core
  order: 2 # Second

- id: admin
  name: Admin
  group: Core
  order: 3 # Third
```

## Health Checks

Configure browser-based health polling:

```yaml
- id: api
  name: API
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30
```

| Field             | Description             |
| ----------------- | ----------------------- |
| `url`             | Base URL of service     |
| `health_endpoint` | Path to health endpoint |
| `health_interval` | Seconds between checks  |

The status page polls endpoints and shows:

- Current status indicator
- Response latency
- Last check time

See [Health Checks](06-health-checks.md) for details.

## Uptime Configuration

Track uptime per component:

```yaml
- id: api
  name: API
  uptime:
    mode: incidents # or "api"
    sla_target: 99.9
    period_days: 90
```

| Field         | Description                       |
| ------------- | --------------------------------- |
| `mode`        | Data source: `incidents` or `api` |
| `sla_target`  | Target uptime percentage          |
| `period_days` | Tracking period                   |

See [Uptime Tracking](05-uptime-tracking.md) for details.

## Examples

### Minimal

```yaml
- id: website
  name: Website
```

### Basic Service

```yaml
- id: api
  name: REST API
  description: Public API endpoints
  status: operational
  group: Backend
```

### With Health Check

```yaml
- id: api
  name: API
  status: operational
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30
```

### With Uptime Tracking

```yaml
- id: api
  name: API
  status: operational
  uptime:
    mode: incidents
    sla_target: 99.9
    period_days: 90
```

### Full Configuration

```yaml
- id: api
  name: REST API
  description: Core API for all client applications
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
```

## Complete Example

```yaml
# __status__/components.yaml

# Core Services
- id: api
  name: API
  description: REST API
  status: operational
  group: Core Services
  order: 1
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30
  uptime:
    mode: incidents
    sla_target: 99.9

- id: web
  name: Web Application
  description: Main web interface
  status: operational
  group: Core Services
  order: 2
  url: https://app.example.com
  health_endpoint: /health
  health_interval: 30

- id: mobile-api
  name: Mobile API
  description: API for mobile apps
  status: operational
  group: Core Services
  order: 3

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

- id: cdn
  name: CDN
  description: Content delivery network
  status: operational
  group: Infrastructure
  order: 2

- id: storage
  name: Object Storage
  description: S3-compatible storage
  status: operational
  group: Infrastructure
  order: 3
```
