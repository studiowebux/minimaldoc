---
title: Scheduled Maintenance Title
description: Brief description of maintenance work
affected_components:
  - database
  - api
scheduled_start: 2025-02-01T02:00:00Z
scheduled_end: 2025-02-01T06:00:00Z
status: scheduled
---

Brief overview of the maintenance work being performed.

## Expected Impact

- Service A will be unavailable for approximately X minutes
- Service B will operate in read-only mode
- Users may experience brief interruptions

## Preparation

Any actions users should take before the maintenance window.

## Contact

For questions, contact support@example.com.

<!--
FRONTMATTER REFERENCE:

title: (required) Short description of maintenance
description: (optional) Longer description
affected_components: (optional) List of component IDs from components.yaml
scheduled_start: (required) ISO 8601 timestamp for start
scheduled_end: (required) ISO 8601 timestamp for end
status: (required) scheduled | in_progress | completed

FILENAME FORMAT:
YYYY-MM-DD-short-description.md
Example: 2025-02-01-database-upgrade.md

NOTES:
- Maintenance with status: scheduled and future dates appear in "Upcoming Maintenance"
- Update status to in_progress when maintenance begins
- Update status to completed when finished
-->
