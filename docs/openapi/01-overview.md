---
title: OpenAPI Overview
description: Interactive API documentation from OpenAPI/Swagger specifications
tags:
  - openapi
  - swagger
  - api
---

# OpenAPI Overview

MinimalDoc generates interactive API documentation from OpenAPI 3.x specifications.

## Features

| Feature | Description |
|---------|-------------|
| Multiple Views | Organize by path, tag, or flat list |
| Interactive Testing | Try API calls in browser |
| Authentication | Bearer, API Key, OAuth 2.0 |
| Code Samples | curl, JavaScript, Go, Python, Swift |
| Schema Browser | Explore data models |
| $ref Resolution | Automatic component reference resolution |
| Export | Generate cURL commands |
| Dark Mode | Full theme support |

## Quick Start

1. Place OpenAPI spec in `api/` directory:

```
docs/
└── api/
    └── openapi.yaml
```

2. Enable in config:

```yaml
openapi:
  enabled: true
  spec_files:
    - "api/openapi.yaml"
```

3. Build:

```bash
minimaldoc build --openapi
```

## Supported Formats

| Format | Extension |
|--------|-----------|
| YAML | `.yaml`, `.yml` |
| JSON | `.json` |

## OpenAPI Versions

| Version | Support |
|---------|---------|
| OpenAPI 3.0.x | Full |
| OpenAPI 3.1.x | Full |
| Swagger 2.0 | Partial (convert to 3.x) |

## Generated Output

```
public/
├── api/
│   ├── index.html           # API documentation
│   └── openapi.json         # Processed spec
└── ...
```

## Example Spec

```yaml
openapi: 3.0.3
info:
  title: My API
  version: 1.0.0
  description: API for my service

servers:
  - url: https://api.example.com/v1
    description: Production

paths:
  /users:
    get:
      summary: List users
      operationId: listUsers
      tags:
        - Users
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        email:
          type: string
          format: email
```

## Navigation Integration

When OpenAPI is enabled, an "API Reference" link appears in the sidebar:

```
- Documentation
  - Getting Started
  - Features
- API Reference        # OpenAPI docs
- Status               # Status page (if enabled)
```

## Multiple Specs

Document multiple APIs:

```yaml
openapi:
  enabled: true
  spec_files:
    - "api/v1/openapi.yaml"
    - "api/v2/openapi.yaml"
    - "api/admin/openapi.yaml"
```

Each generates a separate documentation page.

## Remote Specs

Fetch specs from URLs:

```yaml
openapi:
  enabled: true
  spec_urls:
    - "https://api.example.com/openapi.json"
  sync_on_build: true
  cache_dir: ".openapi-cache"
```

Cached locally to avoid repeated fetches.

## Next Steps

- [Configuration](02-configuration.md) - Full configuration options
- [Views](03-views.md) - Organize endpoints
- [Testing](04-testing.md) - Try API calls
- [Authentication](05-authentication.md) - Configure auth
