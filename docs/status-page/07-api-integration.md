---
title: API & Integrations
description: JSON API, RSS feeds, and external integrations
tags:
  - status-page
  - api
  - rss
---

# API & Integrations

Integrate your status page with external tools and services.

## JSON API

Machine-readable status at `/status/status.json`:

```json
{
  "title": "Service Status",
  "description": "Current operational status",
  "overall_status": "operational",
  "last_updated": "2025-01-28T15:00:00Z",
  "components": [
    {
      "id": "api",
      "name": "API",
      "description": "Core REST API",
      "status": "operational",
      "group": "Core Services",
      "url": "https://api.example.com",
      "health_endpoint": "/health",
      "health_interval": 30,
      "uptime": {
        "mode": "incidents",
        "period_days": 90,
        "uptime_percent": 99.95,
        "sla": {
          "target": 99.9,
          "current_7d": 100.0,
          "current_30d": 99.97,
          "current_90d": 99.95
        }
      }
    }
  ],
  "active_incidents": [
    {
      "id": "2025-01-28-api-latency",
      "title": "API Latency Issues",
      "status": "monitoring",
      "severity": "minor",
      "affected_components": ["api"],
      "created_at": "2025-01-28T10:00:00Z",
      "url": "/status/incident/2025-01-28-api-latency.html"
    }
  ],
  "scheduled_maintenance": [
    {
      "id": "2025-02-01-upgrade",
      "title": "Database Upgrade",
      "affected_components": ["database"],
      "scheduled_start": "2025-02-01T02:00:00Z",
      "scheduled_end": "2025-02-01T04:00:00Z",
      "status": "scheduled"
    }
  ]
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Status page title |
| `description` | string | Status page description |
| `overall_status` | string | Aggregate status |
| `last_updated` | datetime | Last build time |
| `components` | array | All components |
| `active_incidents` | array | Open incidents |
| `scheduled_maintenance` | array | Upcoming maintenance |

### Use Cases

- Custom dashboards
- Monitoring tools
- Slack/Discord bots
- Mobile apps

## RSS Feed

Subscribe to updates at `/status/feed.xml`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Service Status</title>
    <link>https://example.com/status</link>
    <description>Service status updates</description>
    <item>
      <title>API Latency Issues</title>
      <link>https://example.com/status/incident/2025-01-28-api-latency.html</link>
      <pubDate>Tue, 28 Jan 2025 10:00:00 GMT</pubDate>
      <description>Investigating API latency issues...</description>
    </item>
  </channel>
</rss>
```

### Configuration

```yaml
status:
  rss_enabled: true
```

### Feed Contents

- Active incidents
- Recently resolved incidents (last 10)
- Scheduled maintenance

### Subscribe

Add to any RSS reader:

```
https://your-site.com/status/feed.xml
```

## Slack Integration

### Webhook Notification

```bash
#!/bin/bash
# check-status.sh

STATUS=$(curl -s https://example.com/status/status.json | jq -r '.overall_status')

if [ "$STATUS" != "operational" ]; then
  curl -X POST -H 'Content-type: application/json' \
    --data '{"text":"Status Alert: '$STATUS'"}' \
    $SLACK_WEBHOOK_URL
fi
```

### RSS to Slack

Use Slack's RSS app to subscribe to your feed.

## Discord Integration

### Webhook

```python
import requests
import json

status = requests.get('https://example.com/status/status.json').json()

if status['overall_status'] != 'operational':
    webhook_url = 'https://discord.com/api/webhooks/...'
    message = {
        'content': f"Status Alert: {status['overall_status']}",
        'embeds': [{
            'title': 'Active Incidents',
            'description': '\n'.join([
                f"- {i['title']}" for i in status['active_incidents']
            ])
        }]
    }
    requests.post(webhook_url, json=message)
```

## Monitoring Integration

### Prometheus

Expose metrics endpoint:

```python
# Custom exporter
from prometheus_client import Gauge

status_gauge = Gauge('service_status', 'Service status', ['component'])

def update_metrics():
    status = requests.get('https://example.com/status/status.json').json()
    status_map = {'operational': 1, 'degraded': 0.5, 'outage': 0}
    for component in status['components']:
        status_gauge.labels(component=component['id']).set(
            status_map.get(component['status'], 0)
        )
```

### Datadog

```python
from datadog import statsd

def check_status():
    status = requests.get('https://example.com/status/status.json').json()
    for component in status['components']:
        statsd.gauge(
            f'status.{component["id"]}',
            1 if component['status'] == 'operational' else 0
        )
```

## Widget Embedding

Embed status in other pages:

```html
<div id="status-widget"></div>
<script>
  fetch('https://example.com/status/status.json')
    .then(r => r.json())
    .then(data => {
      const widget = document.getElementById('status-widget');
      const status = data.overall_status;
      const color = status === 'operational' ? 'green' : 'red';
      widget.innerHTML = `
        <span style="color: ${color}">
          Status: ${status}
        </span>
      `;
    });
</script>
```

## Automation

### GitHub Actions

Check status before deployment:

```yaml
- name: Check Status
  run: |
    STATUS=$(curl -s https://example.com/status/status.json | jq -r '.overall_status')
    if [ "$STATUS" != "operational" ]; then
      echo "Services not operational, aborting deployment"
      exit 1
    fi
```

### CI/CD Gate

```bash
#!/bin/bash
# pre-deploy-check.sh

STATUS_URL="https://example.com/status/status.json"
STATUS=$(curl -s $STATUS_URL | jq -r '.overall_status')

case $STATUS in
  operational)
    echo "All systems operational, proceeding..."
    exit 0
    ;;
  degraded)
    echo "Warning: Systems degraded, proceeding with caution..."
    exit 0
    ;;
  *)
    echo "Error: Systems in $STATUS state, aborting deployment"
    exit 1
    ;;
esac
```

## External Status Pages

### Link to External

For components managed by third parties:

```yaml
- id: aws
  name: AWS
  description: Amazon Web Services
  status: operational
  url: https://status.aws.amazon.com
```

No health check, just a link to their status page.

### Aggregate Multiple Pages

```python
def aggregate_status():
    sources = [
        'https://example.com/status/status.json',
        'https://api.example.com/status/status.json',
    ]

    statuses = []
    for url in sources:
        data = requests.get(url).json()
        statuses.append(data['overall_status'])

    # Worst status wins
    if 'major_outage' in statuses:
        return 'major_outage'
    if 'partial_outage' in statuses:
        return 'partial_outage'
    if 'degraded' in statuses:
        return 'degraded'
    return 'operational'
```

## Webhooks (Future)

For real-time notifications, consider:
- GitHub webhooks on incident file changes
- CI/CD triggered notifications
- Custom webhook endpoints

See [Prometheus Integration](../guides/04-prometheus-integration.md) for automated incident creation.
