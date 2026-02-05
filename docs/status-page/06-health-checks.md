---
title: Health Checks
description: Live endpoint polling with status indicators
tags:
  - status-page
  - health-checks
  - monitoring
---

# Health Checks

Browser-based health polling for real-time status display.

## Overview

Health checks provide:
- Live endpoint polling
- Response latency display
- Visual status indicators
- Countdown to next check

## Configuration

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

| Status | Indicator | Condition |
|--------|-----------|-----------|
| Healthy | Green dot | 2xx response |
| Degraded | Yellow dot | Slow response (>1s) |
| Unhealthy | Red dot | Non-2xx or timeout |
| Checking | Pulse animation | Request in progress |

### Latency Display

| Range | Display |
|-------|---------|
| < 100ms | Green |
| 100-500ms | Yellow |
| 500ms-1s | Orange |
| > 1s | Red |

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

| Response | Status |
|----------|--------|
| 200-299 | Healthy |
| 300-399 | Healthy (redirect) |
| 400-499 | Unhealthy |
| 500-599 | Unhealthy |
| Timeout | Unhealthy |
| Network error | Unhealthy |

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
  health_interval: 10    # Every 10 seconds
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

### Server Configuration

**Nginx:**

```nginx
location /health {
    add_header Access-Control-Allow-Origin *;
    return 200 '{"status":"ok"}';
}
```

**Express:**

```javascript
app.get('/health', (req, res) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.json({ status: 'ok' });
});
```

**Go:**

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

## Without Health Checks

Components without health configuration show static status:

```yaml
- id: external-service
  name: External Service
  status: operational    # Manual status only
```

No live polling, no latency display.

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

## Best Practices

### Lightweight Endpoints

```go
// Good - Fast, minimal work
func health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
}

// Avoid - Slow, resource-intensive
func health(w http.ResponseWriter, r *http.Request) {
    db.Query("SELECT 1")  // Adds latency
    cache.Ping()
    json.NewEncoder(w).Encode(fullStatusReport)
}
```

### Appropriate Intervals

| Service Type | Interval |
|--------------|----------|
| Critical | 10-15 seconds |
| Standard | 30 seconds |
| Non-critical | 60 seconds |

### Handle Timeouts

Set reasonable timeout (5-10 seconds):

```javascript
// Status page JS uses 10s timeout
const response = await fetch(url, {
  signal: AbortSignal.timeout(10000)
});
```

### Dedicated Endpoint

Use a dedicated `/health` endpoint, not your main page:

```yaml
# Good
health_endpoint: /health

# Avoid
health_endpoint: /
```
