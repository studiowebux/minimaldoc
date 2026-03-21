---
title: Scheduled Maintenance
description: Plan and communicate maintenance windows
tags:
  - status-page
  - maintenance
---

# Scheduled Maintenance

Maintenance windows communicate planned downtime to users.

## File Structure

Place maintenance in `__status__/maintenance/`:

```
__status__/
└── maintenance/
    ├── 2025-02-01-database-upgrade.md
    ├── 2025-02-15-infrastructure.md
    └── 2025-03-01-migration.md
```

Filename format: `YYYY-MM-DD-short-description.md`

## Basic Format

```markdown
---
title: Database Upgrade
description: Upgrading PostgreSQL from 15 to 16
affected_components:
  - database
  - api
scheduled_start: 2025-02-01T02:00:00Z
scheduled_end: 2025-02-01T04:00:00Z
status: scheduled
---

We will be upgrading our database cluster to PostgreSQL 16.

## Expected Impact

- Database will be in read-only mode for approximately 30 minutes
- API responses may be slower during migration
- Brief service interruption expected (< 5 minutes)

## Preparation

No action required from users. All data will be preserved.
```

## Frontmatter Fields

| Field                 | Type     | Required | Description           |
| --------------------- | -------- | -------- | --------------------- |
| `title`               | string   | Yes      | Maintenance title     |
| `description`         | string   | No       | Brief description     |
| `affected_components` | []string | No       | Component IDs         |
| `scheduled_start`     | datetime | Yes      | Start time (ISO 8601) |
| `scheduled_end`       | datetime | Yes      | End time (ISO 8601)   |
| `status`              | string   | Yes      | Current status        |

## Status Values

| Status        | Description           | Display        |
| ------------- | --------------------- | -------------- |
| `scheduled`   | Planned, not started  | Blue calendar  |
| `in_progress` | Currently underway    | Yellow spinner |
| `completed`   | Finished successfully | Green check    |

## Content Structure

### Impact Section

Describe what users will experience:

```markdown
## Expected Impact

- Service will be unavailable for approximately 2 hours
- Users will be logged out
- Scheduled jobs will be paused
```

### Preparation Section

Tell users what to do:

```markdown
## Preparation

- Save any unsaved work before the maintenance window
- Download reports you need before 02:00 UTC
- Plan for service unavailability
```

### Updates Section

Add updates as maintenance progresses:

```markdown
## Updates

**02:15 UTC** - Maintenance started. Database backup in progress.

**02:45 UTC** - Backup complete. Starting migration.

**03:30 UTC** - Migration complete. Running verification.

**03:45 UTC** - Verification passed. Bringing services back online.

**04:00 UTC** - Maintenance complete. All services operational.
```

## Examples

### Scheduled

```markdown
---
title: Infrastructure Update
description: Updating server infrastructure
affected_components:
  - api
  - web
scheduled_start: 2025-02-15T03:00:00Z
scheduled_end: 2025-02-15T05:00:00Z
status: scheduled
---

We will be updating our server infrastructure to improve performance and reliability.

## Expected Impact

- Brief interruptions possible (< 1 minute each)
- No data loss expected
- Service will remain mostly available

## Preparation

No action required.
```

### In Progress

```markdown
---
title: Database Migration
description: Migrating to new database cluster
affected_components:
  - database
scheduled_start: 2025-02-01T02:00:00Z
scheduled_end: 2025-02-01T06:00:00Z
status: in_progress
---

Database migration is currently in progress.

## Expected Impact

- Read-only mode until migration completes
- Write operations will fail temporarily

## Updates

**02:00 UTC** - Maintenance started.

**02:30 UTC** - Data export complete. Starting import to new cluster.

**04:00 UTC** - Import 75% complete. On schedule.
```

### Completed

```markdown
---
title: Security Patch
description: Applying critical security patches
affected_components:
  - api
  - web
scheduled_start: 2025-01-20T04:00:00Z
scheduled_end: 2025-01-20T04:30:00Z
status: completed
---

Security patches have been successfully applied.

## Summary

- All patches applied successfully
- No service interruption occurred
- Systems verified and operational

## Updates

**04:00 UTC** - Started applying patches.

**04:20 UTC** - Patches applied. Running verification.

**04:25 UTC** - Verification complete. All systems operational.
```

## Display

### Dashboard

Upcoming maintenance shown in dedicated section:

```
┌─────────────────────────────────────────────────┐
│  Scheduled Maintenance                           │
├─────────────────────────────────────────────────┤
│  Database Upgrade                                │
│  Feb 1, 02:00 - 04:00 UTC                       │
│  Affects: Database, API                          │
├─────────────────────────────────────────────────┤
│  Infrastructure Update                           │
│  Feb 15, 03:00 - 05:00 UTC                      │
│  Affects: API, Web                               │
└─────────────────────────────────────────────────┘
```

### Component Status

During maintenance, affected components show:

```
API          Under Maintenance
Database     Under Maintenance
Web          Operational
```

### RSS Feed

Maintenance announcements included in feed.

## Timezone

All times in ISO 8601 with timezone:

```yaml
# UTC
scheduled_start: 2025-02-01T02:00:00Z

# With offset
scheduled_start: 2025-02-01T02:00:00+00:00

# Local time (converted to UTC)
scheduled_start: 2025-01-31T21:00:00-05:00
```

Display shows user's local timezone.
