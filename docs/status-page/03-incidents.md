---
title: Incidents
description: Document and track service incidents
tags:
  - status-page
  - incidents
---

# Incidents

Incidents document service issues and their resolution.

## File Structure

Place incidents in `__status__/incidents/`:

```
__status__/
└── incidents/
    ├── 2025-01-28-api-outage.md
    ├── 2025-01-15-database-latency.md
    └── 2025-01-10-cdn-issue.md
```

Filename format: `YYYY-MM-DD-short-description.md`

## Basic Format

```markdown
---
title: API Service Outage
status: resolved
severity: major
affected_components:
  - api
created_at: 2025-01-28T10:00:00Z
resolved_at: 2025-01-28T12:30:00Z
---

## Update - 12:30 UTC

Issue resolved. The root cause was identified as a misconfigured load balancer.

## Update - 11:15 UTC

Identified the issue. Deploying fix.

## Update - 10:00 UTC

Investigating reports of API unavailability.
```

## Frontmatter Fields

| Field                 | Type     | Required | Description           |
| --------------------- | -------- | -------- | --------------------- |
| `title`               | string   | Yes      | Incident title        |
| `status`              | string   | Yes      | Current status        |
| `severity`            | string   | Yes      | Severity level        |
| `affected_components` | []string | No       | Component IDs         |
| `created_at`          | datetime | Yes      | Start time (ISO 8601) |
| `updated_at`          | datetime | No       | Last update time      |
| `resolved_at`         | datetime | No       | Resolution time       |

## Status Values

| Status          | Description             | Display |
| --------------- | ----------------------- | ------- |
| `investigating` | Initial investigation   | Yellow  |
| `identified`    | Root cause found        | Orange  |
| `monitoring`    | Fix applied, monitoring | Blue    |
| `resolved`      | Fully resolved          | Green   |

## Severity Levels

| Severity   | Description        | Impact              |
| ---------- | ------------------ | ------------------- |
| `minor`    | Low impact         | Some users affected |
| `major`    | Significant impact | Most users affected |
| `critical` | Service down       | All users affected  |

## Update Timeline

Write updates as H2 headings with timestamps:

```markdown
## Update - 14:30 UTC

Service fully restored. Post-mortem scheduled.

## Update - 13:45 UTC

Fix deployed to production. Monitoring.

## Update - 12:00 UTC

Identified root cause as database connection pool exhaustion.

## Update - 11:30 UTC

Investigating elevated error rates on API endpoints.
```

Updates display newest first on the status page.

### Timestamp Formats

```markdown
## Update - 14:30 UTC

## Update - 2:30 PM UTC

## Update - January 28, 14:30 UTC

## Update - 2025-01-28 14:30 UTC
```

## Affected Components

Link incidents to components:

```yaml
affected_components:
  - api
  - web
  - mobile-api
```

Component IDs must match those in `components.yaml`.

Effects:

- Components show incident indicator
- Uptime calculations include downtime
- Filtering by component

## Examples

### Investigating

```markdown
---
title: Elevated API Latency
status: investigating
severity: minor
affected_components:
  - api
created_at: 2025-01-28T15:00:00Z
---

## Update - 15:00 UTC

We are investigating reports of slow API responses. Some users may experience delays.
```

### Identified

```markdown
---
title: Database Connection Issues
status: identified
severity: major
affected_components:
  - database
  - api
created_at: 2025-01-28T10:00:00Z
---

## Update - 10:45 UTC

Root cause identified: connection pool exhaustion due to long-running queries. Working on fix.

## Update - 10:00 UTC

Investigating database connectivity issues affecting API responses.
```

### Monitoring

```markdown
---
title: CDN Outage
status: monitoring
severity: critical
affected_components:
  - cdn
created_at: 2025-01-28T08:00:00Z
---

## Update - 09:30 UTC

CDN service restored. Monitoring for stability.

## Update - 08:30 UTC

Identified issue with CDN provider. Working with vendor.

## Update - 08:00 UTC

CDN serving errors for static assets.
```

### Resolved

```markdown
---
title: Authentication Service Outage
status: resolved
severity: major
affected_components:
  - api
  - web
created_at: 2025-01-28T06:00:00Z
resolved_at: 2025-01-28T07:15:00Z
---

## Update - 07:15 UTC

Service fully restored. Root cause was an expired TLS certificate on the auth service.

## Update - 06:45 UTC

Certificate renewed. Restarting services.

## Update - 06:15 UTC

Identified expired certificate causing auth failures.

## Update - 06:00 UTC

Users unable to log in. Investigating.
```

## Post-Mortem

Include post-mortem in resolved incidents:

```markdown
---
title: Major API Outage
status: resolved
severity: critical
affected_components:
  - api
created_at: 2025-01-28T10:00:00Z
resolved_at: 2025-01-28T14:00:00Z
---

## Post-Mortem

**Duration:** 4 hours

**Root Cause:** A deployment introduced a memory leak that caused API servers to crash under load.

**Impact:** API was unavailable for all users during the incident.

**Timeline:**

- 10:00 - Deployment completed
- 10:15 - First alerts triggered
- 10:30 - Rollback initiated
- 14:00 - Service restored

**Action Items:**

- Add memory usage monitoring
- Improve deployment canary process
- Update runbook for similar incidents

## Update - 14:00 UTC

Service fully restored after rollback.

## Update - 10:30 UTC

Rolling back recent deployment.

## Update - 10:00 UTC

Investigating API unavailability.
```

## Display

Incidents appear on:

1. **Dashboard** - Active and recent incidents
2. **History Page** - All incidents, grouped by month
3. **Individual Pages** - Full incident details
4. **RSS Feed** - Updates for subscribers
