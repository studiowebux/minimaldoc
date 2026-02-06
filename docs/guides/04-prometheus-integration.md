---
title: Prometheus Integration
description: Automate status page updates with Prometheus AlertManager
tags:
  - guides
  - prometheus
  - alertmanager
  - automation
---

# Prometheus Integration

Automate status page updates using Prometheus AlertManager and GitHub Actions.

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   Prometheus    │────▶│  AlertManager    │────▶│ GitHub Webhook  │
│   (Monitoring)  │     │  (Alert Router)  │     │  (Receiver)     │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                                                          │
                                                          ▼
                                                 ┌─────────────────┐
                                                 │ GitHub Actions  │
                                                 │ (Create/Update) │
                                                 └─────────────────┘
                                                          │
                                                          ▼
                                                 ┌─────────────────┐
                                                 │  Status Page    │
                                                 │ (Rebuilt Site)  │
                                                 └─────────────────┘
```

## Components

| Component      | Purpose                                |
| -------------- | -------------------------------------- |
| Prometheus     | Monitor services, evaluate alert rules |
| AlertManager   | Route alerts, send webhooks            |
| GitHub Actions | Process alerts, update incident files  |
| MinimalDoc     | Rebuild status page                    |

## Prometheus Alert Rules

Create alert rules for your services:

```yaml
# prometheus/rules/alerts.yml
groups:
  - name: service-alerts
    rules:
      # Service Down
      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.job }} is down"
          description: "{{ $labels.job }} has been down for more than 1 minute"
          component_id: "{{ $labels.component }}"

      # High Latency
      - alert: HighLatency
        expr: http_request_duration_seconds{quantile="0.99"} > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency on {{ $labels.job }}"
          description: "99th percentile latency is {{ $value }}s"
          component_id: "{{ $labels.component }}"

      # High Error Rate
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m])) by (job)
          /
          sum(rate(http_requests_total[5m])) by (job)
          > 0.05
        for: 5m
        labels:
          severity: major
        annotations:
          summary: "High error rate on {{ $labels.job }}"
          description: "Error rate is {{ $value | humanizePercentage }}"
          component_id: "{{ $labels.component }}"

      # Database Connection Pool
      - alert: DatabaseConnectionPoolExhausted
        expr: pg_stat_activity_count >= pg_settings_max_connections * 0.9
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Database connection pool near exhaustion"
          description: "{{ $value }} of {{ $labels.max_connections }} connections in use"
          component_id: "database"
```

## AlertManager Configuration

Configure AlertManager to send webhooks to GitHub:

```yaml
# alertmanager/config.yml
global:
  resolve_timeout: 5m

route:
  receiver: "github-status"
  group_by: ["alertname", "component_id"]
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h

  routes:
    # Critical alerts - immediate
    - match:
        severity: critical
      receiver: "github-status"
      group_wait: 10s

    # Warning alerts - batched
    - match:
        severity: warning
      receiver: "github-status"
      group_wait: 2m

receivers:
  - name: "github-status"
    webhook_configs:
      - url: "https://api.github.com/repos/YOUR_ORG/YOUR_REPO/dispatches"
        http_config:
          authorization:
            type: Bearer
            credentials_file: /etc/alertmanager/secrets/github-token
        send_resolved: true
        max_alerts: 10
```

### GitHub Token

Create a fine-grained Personal Access Token with:

- Repository access: Your docs repository
- Permissions: Contents (read/write), Actions (write)

Store securely:

```bash
# Create secret file
echo "ghp_your_token_here" > /etc/alertmanager/secrets/github-token
chmod 600 /etc/alertmanager/secrets/github-token
```

## GitHub Actions Workflow

Create workflow to process alerts:

```yaml
# .github/workflows/status-update.yml
name: Update Status Page

on:
  repository_dispatch:
    types: [alertmanager]

jobs:
  process-alert:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Process Alert Payload
        id: process
        run: |
          # Extract alert data
          ALERTS='${{ toJson(github.event.client_payload.alerts) }}'

          # Get first alert details
          STATUS=$(echo "$ALERTS" | jq -r '.[0].status')
          ALERTNAME=$(echo "$ALERTS" | jq -r '.[0].labels.alertname')
          SEVERITY=$(echo "$ALERTS" | jq -r '.[0].labels.severity')
          COMPONENT=$(echo "$ALERTS" | jq -r '.[0].annotations.component_id // "unknown"')
          SUMMARY=$(echo "$ALERTS" | jq -r '.[0].annotations.summary')
          DESCRIPTION=$(echo "$ALERTS" | jq -r '.[0].annotations.description')

          # Generate incident ID
          DATE=$(date +%Y-%m-%d)
          INCIDENT_ID="${DATE}-${ALERTNAME,,}"

          # Set outputs
          echo "status=$STATUS" >> $GITHUB_OUTPUT
          echo "alertname=$ALERTNAME" >> $GITHUB_OUTPUT
          echo "severity=$SEVERITY" >> $GITHUB_OUTPUT
          echo "component=$COMPONENT" >> $GITHUB_OUTPUT
          echo "summary=$SUMMARY" >> $GITHUB_OUTPUT
          echo "description=$DESCRIPTION" >> $GITHUB_OUTPUT
          echo "incident_id=$INCIDENT_ID" >> $GITHUB_OUTPUT
          echo "date=$DATE" >> $GITHUB_OUTPUT

      - name: Create or Update Incident
        run: |
          INCIDENT_FILE="docs/__status__/incidents/${{ steps.process.outputs.incident_id }}.md"
          TIMESTAMP=$(date -u +"%H:%M UTC")

          if [ "${{ steps.process.outputs.status }}" == "firing" ]; then
            # Create new incident or add update
            if [ ! -f "$INCIDENT_FILE" ]; then
              # New incident
              cat > "$INCIDENT_FILE" << EOF
          ---
          title: ${{ steps.process.outputs.summary }}
          status: investigating
          severity: ${{ steps.process.outputs.severity }}
          affected_components:
            - ${{ steps.process.outputs.component }}
          created_at: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
          ---

          ## Update - $TIMESTAMP

          ${{ steps.process.outputs.description }}

          Alert triggered by automated monitoring.
          EOF
            else
              # Add update to existing incident
              CURRENT=$(cat "$INCIDENT_FILE")
              FRONTMATTER=$(echo "$CURRENT" | sed -n '/^---$/,/^---$/p')
              BODY=$(echo "$CURRENT" | sed '1,/^---$/d' | sed '1,/^---$/d')

              # Update status to identified if still investigating
              FRONTMATTER=$(echo "$FRONTMATTER" | sed 's/status: investigating/status: identified/')

              cat > "$INCIDENT_FILE" << EOF
          $FRONTMATTER

          ## Update - $TIMESTAMP

          ${{ steps.process.outputs.description }}

          $BODY
          EOF
            fi
          else
            # Alert resolved
            if [ -f "$INCIDENT_FILE" ]; then
              CURRENT=$(cat "$INCIDENT_FILE")

              # Update frontmatter
              CURRENT=$(echo "$CURRENT" | sed 's/status: investigating/status: resolved/')
              CURRENT=$(echo "$CURRENT" | sed 's/status: identified/status: resolved/')
              CURRENT=$(echo "$CURRENT" | sed 's/status: monitoring/status: resolved/')

              # Add resolved_at if not present
              if ! grep -q "resolved_at:" "$INCIDENT_FILE"; then
                CURRENT=$(echo "$CURRENT" | sed "/^created_at:/a resolved_at: $(date -u +"%Y-%m-%dT%H:%M:%SZ")")
              fi

              # Add resolution update
              BODY=$(echo "$CURRENT" | sed '1,/^---$/d' | sed '1,/^---$/d')
              FRONTMATTER=$(echo "$CURRENT" | sed -n '/^---$/,/^---$/p')

              cat > "$INCIDENT_FILE" << EOF
          $FRONTMATTER

          ## Update - $TIMESTAMP

          Issue resolved. Services operating normally.

          $BODY
          EOF
            fi
          fi

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.24"

      - name: Install MinimalDoc
        run: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest

      - name: Rebuild Status Page
        run: |
          minimaldoc build ./docs \
            --base-url "${{ vars.DOCS_BASE_URL }}" \
            --output dist \
            --status

      - name: Commit Changes
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add docs/__status__/incidents/
          git diff --staged --quiet || git commit -m "Update incident: ${{ steps.process.outputs.incident_id }}"
          git push

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: dist

      - name: Deploy to GitHub Pages
        uses: actions/deploy-pages@v4
```

## Component ID Mapping

Map Prometheus job labels to MinimalDoc component IDs:

```yaml
# prometheus/prometheus.yml
scrape_configs:
  - job_name: "api"
    static_configs:
      - targets: ["api.example.com:9090"]
        labels:
          component: "api" # Maps to components.yaml id

  - job_name: "web"
    static_configs:
      - targets: ["web.example.com:9090"]
        labels:
          component: "web"

  - job_name: "database"
    static_configs:
      - targets: ["db.example.com:9187"]
        labels:
          component: "database"
```

## Severity Mapping

Map Prometheus severities to MinimalDoc:

| Prometheus | MinimalDoc | Component Status |
| ---------- | ---------- | ---------------- |
| `critical` | `critical` | `major_outage`   |
| `major`    | `major`    | `partial_outage` |
| `warning`  | `minor`    | `degraded`       |

## Testing the Integration

### Test Webhook Locally

```bash
# Simulate AlertManager webhook
curl -X POST \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/YOUR_ORG/YOUR_REPO/dispatches \
  -d '{
    "event_type": "alertmanager",
    "client_payload": {
      "alerts": [{
        "status": "firing",
        "labels": {
          "alertname": "TestAlert",
          "severity": "warning"
        },
        "annotations": {
          "summary": "Test Alert",
          "description": "This is a test alert",
          "component_id": "api"
        }
      }]
    }
  }'
```

### Verify Workflow

1. Check Actions tab for workflow run
2. Verify incident file created
3. Check status page updated

