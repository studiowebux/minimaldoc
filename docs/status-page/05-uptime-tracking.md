---
title: Uptime Tracking
description: Track service availability with uptime grid and SLA display
tags:
  - status-page
  - uptime
  - sla
---

# Uptime Tracking

Track service availability with uptime grid and SLA metrics.

## Overview

Uptime tracking provides:

- 90-day uptime grid (horizontal colored squares)
- SLA percentage display
- Period breakdowns (7d, 30d, 90d)
- Per-component tracking

## Data Sources

Two modes for uptime data:

| Mode        | Source                      | Use Case                  |
| ----------- | --------------------------- | ------------------------- |
| `incidents` | Derived from incident files | Static sites, no backend  |
| `api`       | External API endpoint       | Real-time monitoring data |

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

| Incident Duration  | Daily Status   |
| ------------------ | -------------- |
| 0 minutes          | Operational    |
| 1-30 minutes       | Degraded       |
| 31-240 minutes     | Partial Outage |
| 241+ minutes       | Major Outage   |
| Maintenance window | Maintenance    |

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

| Field            | Type   | Description          |
| ---------------- | ------ | -------------------- |
| `component_id`   | string | Component identifier |
| `period_days`    | int    | History period       |
| `uptime_percent` | float  | Overall uptime       |
| `history`        | array  | Daily status entries |
| `sla`            | object | SLA metrics          |

### History Entry

| Field              | Type   | Description         |
| ------------------ | ------ | ------------------- |
| `date`             | string | Date (YYYY-MM-DD)   |
| `status`           | string | Daily status        |
| `downtime_minutes` | int    | Minutes of downtime |

### SLA Object

| Field         | Type  | Description           |
| ------------- | ----- | --------------------- |
| `target`      | float | Target SLA percentage |
| `current_7d`  | float | Last 7 days uptime    |
| `current_30d` | float | Last 30 days uptime   |
| `current_90d` | float | Last 90 days uptime   |

## Uptime Grid

A horizontal row of colored squares representing the last 90 days. Each square is 10x10px and shows the status for that day. Hover over any square to see the date, status, and downtime minutes.

### Colors

| Status         | Light Mode | Dark Mode   |
| -------------- | ---------- | ----------- |
| Operational    | Green      | Dark Green  |
| Degraded       | Yellow     | Dark Yellow |
| Partial Outage | Orange     | Dark Orange |
| Major Outage   | Red        | Dark Red    |
| Maintenance    | Blue       | Dark Blue   |

## SLA Display

Shows uptime percentage with period breakdowns (7d, 30d, 90d) and target comparison.

### Target Comparison

| Condition         | Display            |
| ----------------- | ------------------ |
| Current >= Target | Green "Met"        |
| Current < Target  | Red "Below Target" |

## Configuration Options

```yaml
uptime:
  mode: incidents # or "api"
  endpoint: /uptime # API mode only
  sla_target: 99.9 # Target percentage
  period_days: 90 # Tracking period
```

| Option        | Default     | Description       |
| ------------- | ----------- | ----------------- |
| `mode`        | `incidents` | Data source mode  |
| `endpoint`    | `/uptime`   | API endpoint path |
| `sla_target`  | `99.9`      | Target uptime %   |
| `period_days` | `90`        | Days of history   |

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
    sla_target: 99.99 # Four nines
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
