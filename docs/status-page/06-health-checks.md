---
title: Health Checks
description: Static and real-time approaches to service monitoring
tags:
  - status-page
  - health-checks
  - monitoring
---

# Health Checks

Two approaches for monitoring service health in MinimalDoc.

## Approach Comparison

| Aspect | Static (CI/CD) | Real-Time (Client) |
|--------|----------------|-------------------|
| Updates | On pipeline trigger | Every page load |
| Server load | Zero | Every visitor polls endpoints |
| Latency display | No | Yes |
| External dependency | Prometheus, Datadog, etc. | CORS-enabled endpoints |
| Cost | Build/deploy only | None |
| Accuracy | Delayed (pipeline latency) | Live |

## Approach 1: Static (CI/CD Driven)

Status updates only when your CI/CD pipeline runs. External monitoring triggers the update.

### How It Works

```
Prometheus detects issue → Webhook triggers pipeline →
Pipeline edits incident file → Site rebuilds → Static files uploaded
```

### Configuration

No `health_endpoint` needed. Set status manually or via automation:

```yaml
# components.yaml
- id: api
  name: API
  status: operational  # Updated by CI/CD
```

### Automation Example

GitHub Actions workflow triggered by external monitoring:

```yaml
name: Update Status
on:
  repository_dispatch:
    types: [status-change]

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Update component status
        run: |
          yq -i '.[] | select(.id == "${{ github.event.client_payload.component }}") .status = "${{ github.event.client_payload.status }}"' \
            docs/__status__/components.yaml
      - name: Create incident if needed
        if: github.event.client_payload.status != 'operational'
        run: |
          cat > docs/__status__/incidents/$(date +%Y-%m-%d)-${{ github.event.client_payload.component }}.md << EOF
          ---
          title: ${{ github.event.client_payload.title }}
          status: investigating
          severity: ${{ github.event.client_payload.severity }}
          affected_components:
            - ${{ github.event.client_payload.component }}
          created_at: $(date -u +%Y-%m-%dT%H:%M:%SZ)
          ---
          Investigating issue detected by monitoring.
          EOF
      - name: Build and deploy
        run: |
          minimaldoc build
          # deploy to hosting...
```

### Benefits

- No client-side requests to your services
- Works with any monitoring stack (Prometheus, Datadog, Pingdom, etc.)
- Zero runtime cost after build
- Full control over when status changes

### Trade-offs

- Status updates have pipeline latency (seconds to minutes)
- Requires external monitoring infrastructure
- No real-time latency display

---

## Approach 2: Real-Time (Client Polling)

Browser polls endpoints directly for live status and latency.

### How It Works

1. Browser loads status page
2. JavaScript polls each configured endpoint
3. Status and latency displayed
4. Countdown shows time to next check
5. Repeat at configured interval

### Configuration

```yaml
# components.yaml
- id: api
  name: API
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30
```

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Base service URL |
| `health_endpoint` | string | Health check path |
| `health_interval` | int | Seconds between checks |

### Benefits

- Live latency display
- Immediate status updates
- No external monitoring needed
- Simple setup

### Trade-offs

- Every visitor polls your endpoints
- Requires CORS headers on health endpoints
- Can increase load on services
- Browser must be open to see updates

---

## Real-Time Details

### Display

```
┌─────────────────────────────────────────────────┐
│  API                                             │
│  https://api.example.com                        │
│                                                  │
│  Status: ● Healthy         Latency: 24ms        │
│  Last check: 10s ago       Next: 20s            │
└─────────────────────────────────────────────────┘
```

## Configuration

```yaml
# components.yaml
- id: api
  name: API
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30
```

| Field             | Type   | Description            |
| ----------------- | ------ | ---------------------- |
| `url`             | string | Base service URL       |
| `health_endpoint` | string | Health check path      |
| `health_interval` | int    | Seconds between checks |

## How It Works

1. Browser loads status page
2. JavaScript polls each configured endpoint
3. Status and latency displayed
4. Countdown shows time to next check
5. Repeat at configured interval

## Display

```
┌─────────────────────────────────────────────────┐
│  API                                             │
│  https://api.example.com                        │
│                                                  │
│  Status: ● Healthy         Latency: 24ms        │
│  Last check: 10s ago       Next: 20s            │
└─────────────────────────────────────────────────┘
```

### Status Indicators

| Status    | Indicator       | Condition           |
| --------- | --------------- | ------------------- |
| Healthy   | Green dot       | 2xx response        |
| Degraded  | Yellow dot      | Slow response (>1s) |
| Unhealthy | Red dot         | Non-2xx or timeout  |
| Checking  | Pulse animation | Request in progress |

### Latency Display

| Range     | Display |
| --------- | ------- |
| < 100ms   | Green   |
| 100-500ms | Yellow  |
| 500ms-1s  | Orange  |
| > 1s      | Red     |

## Health Endpoint

Your service should expose a health endpoint:

### Simple Response

```json
GET /health

200 OK
{
  "status": "ok"
}
```

### Detailed Response

```json
GET /health

200 OK
{
  "status": "ok",
  "version": "1.2.3",
  "uptime": 86400,
  "checks": {
    "database": "ok",
    "cache": "ok",
    "queue": "ok"
  }
}
```

### Unhealthy Response

```json
GET /health

503 Service Unavailable
{
  "status": "unhealthy",
  "checks": {
    "database": "error",
    "cache": "ok"
  }
}
```

## Status Determination

| Response      | Status             |
| ------------- | ------------------ |
| 200-299       | Healthy            |
| 300-399       | Healthy (redirect) |
| 400-499       | Unhealthy          |
| 500-599       | Unhealthy          |
| Timeout       | Unhealthy          |
| Network error | Unhealthy          |

## Configuration Examples

### Basic

```yaml
- id: api
  name: API
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30
```

### Frequent Checks

```yaml
- id: critical-service
  name: Critical Service
  url: https://critical.example.com
  health_endpoint: /health
  health_interval: 10 # Every 10 seconds
```

### Custom Endpoint

```yaml
- id: legacy-service
  name: Legacy Service
  url: https://legacy.example.com
  health_endpoint: /api/v1/status
  health_interval: 60
```

### Multiple Services

```yaml
- id: api
  name: API
  url: https://api.example.com
  health_endpoint: /health
  health_interval: 30

- id: web
  name: Web App
  url: https://app.example.com
  health_endpoint: /health
  health_interval: 30

- id: cdn
  name: CDN
  url: https://cdn.example.com
  health_endpoint: /health
  health_interval: 60
```

## CORS Requirements

Health endpoints must allow cross-origin requests:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, HEAD, OPTIONS
```

## Choosing an Approach

| Use Case | Recommended |
|----------|-------------|
| High traffic status page | Static (CI/CD) |
| Internal tools, low traffic | Real-Time (Client) |
| External services you don't control | Static (CI/CD) |
| Services with CORS-enabled health endpoints | Either |
| No monitoring infrastructure | Real-Time (Client) |

## Static-Only Components

Components without health configuration use the static approach:

```yaml
- id: external-service
  name: External Service
  status: operational  # Updated by CI/CD only
```

No live polling, no latency display. Status changes via pipeline.

## Fallback Behavior

If health check fails:

1. Component shows configured static status
2. Error indicator displayed
3. Retry at next interval

## JSON Output

Health configuration included in `status.json`:

```json
{
  "components": [
    {
      "id": "api",
      "name": "API",
      "url": "https://api.example.com",
      "health_endpoint": "/health",
      "health_interval": 30
    }
  ]
}
```

External tools can use this for their own monitoring.
