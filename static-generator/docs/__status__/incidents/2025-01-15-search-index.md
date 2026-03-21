---
title: Search Index Generation Issue
status: resolved
severity: minor
affected_components:
  - search
  - build-system
created_at: 2025-01-15T14:30:00Z
resolved_at: 2025-01-15T15:45:00Z
---

## Update - 15:45 UTC

Issue resolved. Search functionality fully restored. Root cause was a character encoding issue in certain page titles that caused the index generator to fail silently.

## Update - 15:15 UTC

Identified the issue. A recent documentation update included special characters that weren't being handled correctly by the search indexer. Deploying fix.

## Update - 14:30 UTC

We're investigating reports that some pages are not appearing in search results. Search is functional but may return incomplete results.
