---
title: Config Schema
description: Complete config.yaml schema reference
tags:
  - api-reference
  - configuration
  - schema
---

# Config Schema

Complete `config.yaml` schema reference.

## Full Schema

```yaml
# Basic Settings
title: string                    # Site title
description: string              # Site description
base_url: string                 # Base URL for absolute links
author: string                   # Default author

# Theme
theme: string                    # Theme name (default: "default")
dark_mode: boolean               # Default to dark mode (default: false)

# Features
enable_llms: boolean             # Generate llms.txt (default: true)
enable_search: boolean           # Enable search (default: true)
clean_urls: boolean              # Use /page/ URLs (default: false)

# OpenAPI Configuration
openapi:
  enabled: boolean               # Enable OpenAPI docs (default: false)
  spec_files: [string]           # Local spec file paths
  spec_urls: [string]            # Remote spec URLs
  default_view: string           # View mode: path|tag|flat (default: "path")
  sync_on_build: boolean         # Fetch remote specs (default: false)
  cache_dir: string              # Cache directory (default: ".openapi-cache")
  enable_testing: boolean        # Show Try It interface (default: true)
  enable_export: boolean         # Enable cURL export (default: true)
  enable_code_samples: boolean   # Show code samples (default: true)
  lazy_load_chunk_size: integer  # Chunk size in bytes (default: 51200)

# Status Page Configuration
status:
  enabled: boolean               # Enable status page (default: false)
  title: string                  # Page title (default: "Service Status")
  description: string            # Page description
  path: string                   # Output path (default: "status")
  show_history: boolean          # Show incident history (default: true)
  history_months: integer        # Months of history (default: 12)
  rss_enabled: boolean           # Generate RSS feed (default: true)

# Changelog Configuration
changelog:
  enabled: boolean               # Enable changelog (default: false)
  title: string                  # Page title (default: "Changelog")
  path: string                   # Output path (default: "changelog")
  rss_enabled: boolean           # Generate RSS feed (default: true)
  repository: string             # Repository URL for compare links

# Stale Warning Configuration
stale_warning:
  enabled: boolean               # Enable stale warnings (default: true)
  threshold_days: integer        # Days before stale (default: 365)
  show_update_date: boolean      # Show last update (default: true)

# Social Links
social_links:
  - name: string                 # Display name
    url: string                  # Link URL
    icon: string                 # Icon name
```

## Basic Settings

### title

Site title displayed in header and browser tab.

```yaml
title: My Documentation
```

Type: `string`
Default: `"Documentation"`

### description

Site description for SEO and metadata.

```yaml
description: Comprehensive documentation for my project
```

Type: `string`
Default: `""`

### base_url

Base URL for absolute links, sitemap, and canonical URLs.

```yaml
base_url: https://docs.example.com
```

Type: `string`
Default: `""`

### author

Default author for pages without explicit author.

```yaml
author: Documentation Team
```

Type: `string`
Default: `""`

## Theme Settings

### theme

Theme directory name.

```yaml
theme: default
```

Type: `string`
Default: `"default"`
Options: `default`, `yellow`, or custom theme name

### dark_mode

Start in dark mode by default.

```yaml
dark_mode: true
```

Type: `boolean`
Default: `false`

## Feature Toggles

### enable_llms

Generate `llms.txt` file for AI assistants.

```yaml
enable_llms: true
```

Type: `boolean`
Default: `true`

### enable_search

Enable client-side search functionality.

```yaml
enable_search: true
```

Type: `boolean`
Default: `true`

### clean_urls

Use `/page/` instead of `/page.html` URLs.

```yaml
clean_urls: true
```

Type: `boolean`
Default: `false`

Requires server configuration for fallback.

## OpenAPI Settings

### openapi.enabled

Enable OpenAPI documentation generation.

```yaml
openapi:
  enabled: true
```

Type: `boolean`
Default: `false`

### openapi.spec_files

Local OpenAPI specification files.

```yaml
openapi:
  spec_files:
    - "api/openapi.yaml"
    - "api/v2/api.json"
```

Type: `array of string`
Default: `[]`

### openapi.spec_urls

Remote OpenAPI specification URLs.

```yaml
openapi:
  spec_urls:
    - "https://api.example.com/openapi.json"
```

Type: `array of string`
Default: `[]`

### openapi.default_view

Default view mode for endpoints.

```yaml
openapi:
  default_view: tag
```

Type: `string`
Default: `"path"`
Options: `path`, `tag`, `flat`

### openapi.sync_on_build

Fetch remote specs on every build.

```yaml
openapi:
  sync_on_build: true
```

Type: `boolean`
Default: `false`

### openapi.cache_dir

Directory for caching remote specs.

```yaml
openapi:
  cache_dir: ".openapi-cache"
```

Type: `string`
Default: `".openapi-cache"`

### openapi.enable_testing

Show interactive API testing interface.

```yaml
openapi:
  enable_testing: true
```

Type: `boolean`
Default: `true`

### openapi.enable_export

Enable cURL and restcli export.

```yaml
openapi:
  enable_export: true
```

Type: `boolean`
Default: `true`

### openapi.enable_code_samples

Show auto-generated code samples.

```yaml
openapi:
  enable_code_samples: true
```

Type: `boolean`
Default: `true`

### openapi.lazy_load_chunk_size

Chunk size for lazy loading large specs.

```yaml
openapi:
  lazy_load_chunk_size: 102400
```

Type: `integer`
Default: `51200` (50KB)

## Status Page Settings

### status.enabled

Enable status page generation.

```yaml
status:
  enabled: true
```

Type: `boolean`
Default: `false`

### status.title

Status page title.

```yaml
status:
  title: "System Status"
```

Type: `string`
Default: `"Service Status"`

### status.description

Status page description.

```yaml
status:
  description: "Current operational status of our services"
```

Type: `string`
Default: `""`

### status.path

Output path for status page.

```yaml
status:
  path: "status"
```

Type: `string`
Default: `"status"`

### status.show_history

Show incident history page.

```yaml
status:
  show_history: true
```

Type: `boolean`
Default: `true`

### status.history_months

Months of incident history to display.

```yaml
status:
  history_months: 12
```

Type: `integer`
Default: `12`

### status.rss_enabled

Generate RSS feed for status updates.

```yaml
status:
  rss_enabled: true
```

Type: `boolean`
Default: `true`

## Changelog Settings

### changelog.enabled

Enable changelog generation.

```yaml
changelog:
  enabled: true
```

Type: `boolean`
Default: `false`

### changelog.title

Changelog page title.

```yaml
changelog:
  title: "Release Notes"
```

Type: `string`
Default: `"Changelog"`

### changelog.path

Output path for changelog.

```yaml
changelog:
  path: "releases"
```

Type: `string`
Default: `"changelog"`

### changelog.rss_enabled

Generate RSS feed for releases.

```yaml
changelog:
  rss_enabled: true
```

Type: `boolean`
Default: `true`

### changelog.repository

Repository URL for compare links.

```yaml
changelog:
  repository: "https://github.com/org/repo"
```

Type: `string`
Default: `""`

## Stale Warning Settings

### stale_warning.enabled

Enable stale content warnings.

```yaml
stale_warning:
  enabled: true
```

Type: `boolean`
Default: `true`

### stale_warning.threshold_days

Days before content is considered stale.

```yaml
stale_warning:
  threshold_days: 365
```

Type: `integer`
Default: `365`

### stale_warning.show_update_date

Show last update date in warning.

```yaml
stale_warning:
  show_update_date: true
```

Type: `boolean`
Default: `true`

## Social Links

### social_links

Array of social media links.

```yaml
social_links:
  - name: GitHub
    url: https://github.com/org/repo
    icon: github
  - name: Discord
    url: https://discord.gg/invite
    icon: discord
```

Type: `array of object`
Default: `[]`

### Social Link Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Display name |
| `url` | string | Yes | Link URL |
| `icon` | string | Yes | Icon identifier |

### Available Icons

| Icon | Name |
|------|------|
| `github` | GitHub |
| `twitter` | Twitter/X |
| `linkedin` | LinkedIn |
| `youtube` | YouTube |
| `discord` | Discord |
| `mastodon` | Mastodon |
| `rss` | RSS Feed |
| `email` | Email |
| `website` | Generic Website |
