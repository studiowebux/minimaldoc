---
title: Uptime Tracking
description: Track service availability with calendar grid and SLA display
tags:
  - status-page
  - uptime
  - sla
---

# Uptime Tracking

Track service availability with visual calendar and SLA metrics.

## Overview

Uptime tracking provides:
- 90-day calendar grid (GitHub-style)
- SLA percentage display
- Period breakdowns (7d, 30d, 90d)
- Per-component tracking

## Data Sources

Two modes for uptime data:

| Mode | Source | Use Case |
|------|--------|----------|
| `incidents` | Derived from incident files | Static sites, no backend |
| `api` | External API endpoint | Real-time monitoring data |

## Incidents Mode

Calculate uptime from incident history at build time.

### Configuration

```yaml
# components.yaml
- id: api
  name: API
  uptime:
    mode: incidents
    sla_target: 99.9
    period_days: 90
```

### How It Works

1. Scan incidents where `affected_components` includes this component
2. For each day, calculate downtime from incident timestamps
3. Determine daily status based on downtime
4. Calculate uptime percentage

### Calculation

```
Daily Uptime = (1440 - downtime_minutes) / 1440 * 100

Period Uptime = Average of daily uptimes
```

Where 1440 = minutes in a day.

### Downtime Mapping

| Incident Duration | Daily Status |
|-------------------|--------------|
| 0 minutes | Operational |
| 1-30 minutes | Degraded |
| 31-240 minutes | Partial Outage |
| 241+ minutes | Major Outage |
| Maintenance window | Maintenance |

## API Mode

Fetch uptime data from external monitoring service.

### Configuration

```yaml
# components.yaml
- id: api
  name: API
  url: https://api.example.com
  uptime:
    mode: api
    endpoint: /uptime
    sla_target: 99.9
```

### API Response Format

Your endpoint must return:

```json
{
  "component_id": "api",
  "period_days": 90,
  "uptime_percent": 99.95,
  "history": [
    {
      "date": "2025-01-28",
      "status": "operational",
      "downtime_minutes": 0
    },
    {
      "date": "2025-01-27",
      "status": "partial_outage",
      "downtime_minutes": 45
    }
  ],
  "sla": {
    "target": 99.9,
    "current_7d": 100.0,
    "current_30d": 99.95,
    "current_90d": 99.92
  }
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `component_id` | string | Component identifier |
| `period_days` | int | History period |
| `uptime_percent` | float | Overall uptime |
| `history` | array | Daily status entries |
| `sla` | object | SLA metrics |

### History Entry

| Field | Type | Description |
|-------|------|-------------|
| `date` | string | Date (YYYY-MM-DD) |
| `status` | string | Daily status |
| `downtime_minutes` | int | Minutes of downtime |

### SLA Object

| Field | Type | Description |
|-------|------|-------------|
| `target` | float | Target SLA percentage |
| `current_7d` | float | Last 7 days uptime |
| `current_30d` | float | Last 30 days uptime |
| `current_90d` | float | Last 90 days uptime |

## Calendar Grid

Visual display of daily status:

```
     Jan 2025
Mo Tu We Th Fr Sa Su
       1  2  3  4  5
 6  7  8  9 10 11 12
13 14 15 16 17 18 19
20 21 22 23 24 25 26
27 28 29 30 31

Legend:
■ Operational
■ Degraded
■ Partial Outage
■ Major Outage
■ Maintenance
□ No Data
```

### Colors

| Status | Light Mode | Dark Mode |
|--------|------------|-----------|
| Operational | Green | Dark Green |
| Degraded | Yellow | Dark Yellow |
| Partial Outage | Orange | Dark Orange |
| Major Outage | Red | Dark Red |
| Maintenance | Blue | Dark Blue |

### Hover Tooltip

Hover over any day for details:

```
January 28, 2025
Status: Partial Outage
Downtime: 45 minutes
Incidents: 1
```

## SLA Display

Shows uptime metrics:

```
┌─────────────────────────────────┐
│  API Uptime                     │
│                                 │
│       99.95%                    │
│    Current Uptime               │
│                                 │
│  7 days   30 days   90 days    │
│  100.0%   99.97%    99.95%     │
│                                 │
│  Target: 99.9% ✓ Met            │
└─────────────────────────────────┘
```

### Target Comparison

| Condition | Display |
|-----------|---------|
| Current >= Target | Green "Met" |
| Current < Target | Red "Below Target" |

## Configuration Options

```yaml
uptime:
  mode: incidents          # or "api"
  endpoint: /uptime        # API mode only
  sla_target: 99.9         # Target percentage
  period_days: 90          # Tracking period
```

| Option | Default | Description |
|--------|---------|-------------|
| `mode` | `incidents` | Data source mode |
| `endpoint` | `/uptime` | API endpoint path |
| `sla_target` | `99.9` | Target uptime % |
| `period_days` | `90` | Days of history |

## Examples

### Basic Incidents Mode

```yaml
- id: api
  name: API
  uptime:
    mode: incidents
    sla_target: 99.9
```

### API Mode with Custom Endpoint

```yaml
- id: api
  name: API
  url: https://api.example.com
  uptime:
    mode: api
    endpoint: /metrics/uptime
    sla_target: 99.99
    period_days: 30
```

### High-Availability Service

```yaml
- id: database
  name: Database
  uptime:
    mode: incidents
    sla_target: 99.99    # Four nines
    period_days: 90
```

## Build Output

Uptime data included in `status.json`:

```json
{
  "components": [
    {
      "id": "api",
      "name": "API",
      "uptime": {
        "mode": "incidents",
        "period_days": 90,
        "uptime_percent": 99.95,
        "history": [...],
        "sla": {
          "target": 99.9,
          "current_7d": 100.0,
          "current_30d": 99.97,
          "current_90d": 99.95
        }
      }
    }
  ]
}
```

## Best Practices

### Realistic Targets

| Service Type | Typical SLA |
|--------------|-------------|
| Non-critical | 99.0% |
| Standard | 99.9% |
| High-availability | 99.99% |
| Mission-critical | 99.999% |

### Downtime Budget

| SLA | Monthly Downtime |
|-----|------------------|
| 99.0% | 7.3 hours |
| 99.9% | 43.8 minutes |
| 99.99% | 4.4 minutes |
| 99.999% | 26 seconds |

### Incident Documentation

For accurate incident-mode tracking:
- Always set `created_at` and `resolved_at`
- List all `affected_components`
- Document all incidents, even minor ones
