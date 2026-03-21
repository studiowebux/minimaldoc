---
title: OpenAPI Configuration
description: Configure OpenAPI documentation generation
tags:
  - openapi
  - configuration
---

# OpenAPI Configuration

## Basic Configuration

```yaml
# config.yaml
openapi:
  enabled: true
  spec_files:
    - "api/openapi.yaml"
```

## Full Configuration

```yaml
openapi:
  enabled: true

  # Local spec files (relative to docs root)
  spec_files:
    - "api/openapi.yaml"
    - "api/v2/api.json"

  # Remote specs
  spec_urls:
    - "https://api.example.com/openapi.json"

  # View mode: path, tag, flat
  default_view: "path"

  # Remote spec options
  sync_on_build: false
  cache_dir: ".openapi-cache"

  # Features
  enable_testing: true
  enable_export: true
  enable_code_samples: true

  # Performance
  lazy_load_chunk_size: 51200
```

## Configuration Reference

### Core Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `false` | Enable OpenAPI documentation |
| `spec_files` | []string | `[]` | Local spec file paths |
| `spec_urls` | []string | `[]` | Remote spec URLs |
| `default_view` | string | `"path"` | Default view mode |

### Remote Spec Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `sync_on_build` | bool | `false` | Fetch remote specs on build |
| `cache_dir` | string | `".openapi-cache"` | Cache directory for remote specs |

### Feature Toggles

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enable_testing` | bool | `true` | Show "Try It" interface |
| `enable_export` | bool | `true` | Enable cURL/restcli export |
| `enable_code_samples` | bool | `true` | Show code sample panel |

### Performance

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `lazy_load_chunk_size` | int | `51200` | Chunk size for lazy loading (bytes) |

## CLI Flags

Override config via CLI:

```bash
minimaldoc build \
  --openapi \
  --openapi-dir api
```

| Flag | Description |
|------|-------------|
| `--openapi` | Enable OpenAPI (shorthand) |
| `--openapi-dir` | Directory containing specs |

## Spec File Location

### Default

Place specs in `api/` directory:

```
docs/
└── api/
    ├── openapi.yaml
    └── v2.yaml
```

### Custom Directory

```yaml
openapi:
  spec_files:
    - "specs/main.yaml"
    - "specs/admin.yaml"
```

Or via CLI:

```bash
minimaldoc build --openapi --openapi-dir specs
```

## Multiple Specs

Each spec generates separate documentation:

```yaml
openapi:
  enabled: true
  spec_files:
    - "api/public.yaml"    # /api/public/
    - "api/admin.yaml"     # /api/admin/
    - "api/internal.yaml"  # /api/internal/
```

## Remote Specs

Fetch specs from URLs during build:

```yaml
openapi:
  enabled: true
  spec_urls:
    - "https://api.service1.com/openapi.json"
    - "https://api.service2.com/openapi.yaml"
  sync_on_build: true
  cache_dir: ".openapi-cache"
```

### Caching

- Specs cached locally to avoid repeated fetches
- Cache directory can be gitignored
- Clear cache to force refresh

```bash
rm -rf .openapi-cache
minimaldoc build --openapi
```

## View Modes

Set default view:

```yaml
openapi:
  default_view: "path"  # path | tag | flat
```

| Mode | Description |
|------|-------------|
| `path` | Group by URL path |
| `tag` | Group by OpenAPI tags |
| `flat` | Flat list, no grouping |

Users can switch views in the UI.

## Feature Configuration

### Disable Testing

For read-only documentation:

```yaml
openapi:
  enable_testing: false
```

### Disable Export

Hide cURL/restcli export:

```yaml
openapi:
  enable_export: false
```

### Disable Code Samples

Hide code sample panel:

```yaml
openapi:
  enable_code_samples: false
```

## Examples

### Minimal Setup

```yaml
openapi:
  enabled: true
  spec_files:
    - "api/openapi.yaml"
```

### Read-Only Documentation

```yaml
openapi:
  enabled: true
  spec_files:
    - "api/openapi.yaml"
  enable_testing: false
  enable_export: false
```

### Multiple Remote APIs

```yaml
openapi:
  enabled: true
  spec_urls:
    - "https://api.prod.example.com/openapi.json"
    - "https://api.staging.example.com/openapi.json"
  sync_on_build: true
  cache_dir: ".openapi-cache"
  default_view: "tag"
```

### High-Performance

```yaml
openapi:
  enabled: true
  spec_files:
    - "api/large-api.yaml"
  lazy_load_chunk_size: 102400  # 100KB chunks
```
