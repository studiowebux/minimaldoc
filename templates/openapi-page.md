---
title: Create User
description: API endpoint for creating a new user
openapi_spec: api/openapi.yaml
openapi_path: /users
openapi_method: POST
---

Additional context or usage notes for this endpoint.

<!--
FRONTMATTER REFERENCE:

title: (optional) Page title, defaults to operation summary
description: (optional) Page description
openapi_spec: (required) Path to OpenAPI spec file or URL
openapi_path: (required) API endpoint path (e.g., /users, /pets/{id})
openapi_method: (required) HTTP method (GET, POST, PUT, PATCH, DELETE)

SETUP:
Enable OpenAPI in config.yaml:
  openapi:
    enabled: true
    spec_files:
      - "api/openapi.yaml"
    enable_testing: true
    enable_export: true
    enable_code_samples: true

NOTES:
- The spec file path is relative to the docs directory
- Remote URLs are supported (fetched at build time if sync_on_build is true)
- Without openapi_path/openapi_method, the page embeds the full spec
- With both fields, only the specific operation is rendered
-->
