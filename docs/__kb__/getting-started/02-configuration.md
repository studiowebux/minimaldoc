---
title: "Configuration Basics"
description: "Learn how to configure your MinimalDoc site"
tags: ["configuration", "yaml", "settings"]
---

## Config File

MinimalDoc uses a `config.yaml` file in your docs directory. Here's a minimal example:

```yaml
title: My Documentation
description: Documentation for my project
base_url: https://example.com
```

## Common Options

| Option | Description | Default |
|--------|-------------|---------|
| `title` | Site title | "Documentation" |
| `description` | Site description | "" |
| `theme` | Theme name | "default" |
| `dark_mode` | Enable dark mode by default | false |

## Enabling Features

Enable optional features in your config:

```yaml
enable_search: true
enable_llms: true

changelog:
  enabled: true

status:
  enabled: true
```

## Environment-Specific Config

Use different base URLs for development and production:

```yaml
# Development
base_url: http://localhost:8889

# Production (uncomment when deploying)
# base_url: https://docs.example.com
```
