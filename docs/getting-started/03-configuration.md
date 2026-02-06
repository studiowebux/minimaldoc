---
title: Configuration
description: Configure MinimalDoc with config.yaml and CLI flags
tags:
  - getting-started
  - configuration
---

# Configuration

MinimalDoc is configured via `config.yaml` in your docs directory, with CLI flags for overrides.

## config.yaml

Place in your docs root:

```yaml
# Basic
title: My Documentation
description: Project documentation
base_url: https://docs.example.com
author: Your Name

# Theme
theme: default
dark_mode: false

# Features
enable_llms: true
enable_search: true
clean_urls: false

# OpenAPI
openapi:
  enabled: false
  spec_files:
    - "api/openapi.yaml"
  spec_urls: []
  default_view: "path"
  sync_on_build: false
  cache_dir: ".openapi-cache"
  enable_testing: true
  enable_export: true
  enable_code_samples: true

# Status Page
status:
  enabled: false
  title: "Service Status"
  path: "status"
  show_history: true
  rss_enabled: true

# Changelog
changelog:
  enabled: false
  title: "Changelog"
  path: "changelog"
  repository: "https://github.com/org/repo"

# Stale Warning
stale_warning:
  enabled: true
  threshold_days: 365
  show_update_date: true

# Portfolio
portfolio:
  enabled: false
  title: "Projects"
  description: "Project showcase"
  path: "projects"

# Contact
contact:
  enabled: false
  title: "Contact"
  path: "contact"
  email: "hello@example.com"
  info:
    - icon: "mail"
      text: "hello@example.com"
    - icon: "location"
      text: "Location"

# FAQ
faq:
  enabled: false
  title: "FAQ"
  description: "Frequently asked questions"
  path: "faq"
  categories:
    - name: "General"
      items:
        - question: "What is this?"
          answer: "A documentation generator."

# Legal
legal:
  enabled: false
  path: "legal"
  footer_group: "Legal"

# Landing Page
landing:
  enabled: false
  nav:
    - text: "Docs"
      url: "/getting-started/installation.html"
  hero:
    title: "Hero Title"
    subtitle: "Hero subtitle"
    buttons:
      - text: "Get Started"
        url: "/getting-started/"
        primary: true
  features:
    title: "Features"
    items:
      - emoji: "~"
        title: "Feature"
        description: "Description"
  steps:
    title: "Quick Start"
    items:
      - title: "Step"
        description: "Description"
        code: "command"

# Footer
footer:
  copyright: "2026 Company"
  links:
    - title: "Documentation"
      items:
        - text: "Getting Started"
          url: "/getting-started/"
  social:
    - name: GitHub
      url: https://github.com/org/repo
      icon: github

# Social Links (sidebar)
social_links:
  - name: GitHub
    url: https://github.com/org/repo
    icon: github
```

## Configuration Reference

### Basic Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `title` | string | `"Documentation"` | Site title |
| `description` | string | `""` | Site description |
| `base_url` | string | `""` | Base URL for absolute links |
| `author` | string | `""` | Default author |
| `theme` | string | `"default"` | Theme name |
| `dark_mode` | bool | `false` | Default to dark mode |
| `enable_llms` | bool | `true` | Generate llms.txt |
| `enable_search` | bool | `true` | Enable search |
| `clean_urls` | bool | `false` | Use `/page/` instead of `/page.html` |
| `entrypoint` | string | `""` | Custom homepage file (default: index.md) |

### OpenAPI Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `openapi.enabled` | bool | `false` | Enable OpenAPI docs |
| `openapi.spec_files` | []string | `[]` | Local spec files |
| `openapi.spec_urls` | []string | `[]` | Remote spec URLs |
| `openapi.default_view` | string | `"path"` | View mode: path, tag, flat |
| `openapi.sync_on_build` | bool | `false` | Fetch remote specs on build |
| `openapi.cache_dir` | string | `".openapi-cache"` | Cache directory |
| `openapi.enable_testing` | bool | `true` | Show Try It interface |
| `openapi.enable_export` | bool | `true` | Enable cURL export |
| `openapi.enable_code_samples` | bool | `true` | Show code samples |

### Status Page Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `status.enabled` | bool | `false` | Enable status page |
| `status.title` | string | `"Service Status"` | Page title |
| `status.description` | string | `""` | Page description |
| `status.path` | string | `"status"` | Output path |
| `status.show_history` | bool | `true` | Show incident history |
| `status.history_months` | int | `12` | Months of history |
| `status.rss_enabled` | bool | `true` | Generate RSS feed |

### Changelog Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `changelog.enabled` | bool | `false` | Enable changelog |
| `changelog.title` | string | `"Changelog"` | Page title |
| `changelog.path` | string | `"changelog"` | Output path |
| `changelog.rss_enabled` | bool | `true` | Generate RSS feed |
| `changelog.repository` | string | `""` | Repository URL for compare links |

### Stale Warning Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `stale_warning.enabled` | bool | `true` | Enable stale warnings |
| `stale_warning.threshold_days` | int | `365` | Days before considered stale |
| `stale_warning.show_update_date` | bool | `true` | Show last update date |

### Portfolio Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `portfolio.enabled` | bool | `false` | Enable portfolio |
| `portfolio.title` | string | `"Projects"` | Page title |
| `portfolio.description` | string | `""` | Page description |
| `portfolio.path` | string | `"projects"` | Output path |

### Contact Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `contact.enabled` | bool | `false` | Enable contact page |
| `contact.title` | string | `"Contact"` | Page title |
| `contact.path` | string | `"contact"` | Output path |
| `contact.email` | string | `""` | Contact email |
| `contact.info` | array | `[]` | Info items with icon and text |

### FAQ Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `faq.enabled` | bool | `false` | Enable FAQ page |
| `faq.title` | string | `"FAQ"` | Page title |
| `faq.description` | string | `""` | Page description |
| `faq.path` | string | `"faq"` | Output path |
| `faq.categories` | array | `[]` | FAQ categories with items |

FAQ items can also be defined as markdown files in `__faq__/` directory.

### Legal Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `legal.enabled` | bool | `false` | Enable legal pages |
| `legal.path` | string | `"legal"` | Output path prefix |
| `legal.footer_group` | string | `"Legal"` | Footer section title |

Legal pages are defined as markdown files in `__legal__/` directory.

### Landing Page Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `landing.enabled` | bool | `false` | Enable landing page |
| `landing.nav` | array | `[]` | Navigation links |
| `landing.hero` | object | | Hero section config |
| `landing.features` | object | | Features section config |
| `landing.steps` | object | | Quick start steps config |
| `landing.links` | object | | Resource links section |
| `landing.opensource` | object | | Open source section |

### Footer Settings

| Option | Type | Description |
|--------|------|-------------|
| `footer.copyright` | string | Copyright text |
| `footer.links` | array | Link groups with title and items |
| `footer.social` | array | Social links with name, url, icon |

### Social Links

| Option | Type | Description |
|--------|------|-------------|
| `social_links[].name` | string | Display name |
| `social_links[].url` | string | Link URL |
| `social_links[].icon` | string | Icon: github, twitter, linkedin, youtube, discord, mastodon, rss, email, website |

## CLI Flags

CLI flags override config.yaml:

```bash
minimaldoc build [docs-dir] [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output directory (default: `public`) |
| `--theme` | `-t` | Theme name |
| `--title` | | Site title |
| `--description` | | Site description |
| `--base-url` | | Base URL |
| `--llms` | `-l` | Generate llms.txt |
| `--clean-urls` | | Use clean URLs |
| `--openapi` | | Enable OpenAPI |
| `--openapi-dir` | | OpenAPI specs directory |
| `--status` | | Enable status page |
| `--status-title` | | Status page title |
| `--status-path` | | Status page path |
| `--changelog` | | Enable changelog |

## Precedence

1. CLI flags (highest)
2. config.yaml
3. Defaults (lowest)

## Examples

Build with overrides:

```bash
minimaldoc build ./docs \
  --title "My Project" \
  --base-url "https://docs.example.com" \
  --output dist \
  --openapi \
  --status
```

Production build:

```bash
minimaldoc build ./docs \
  --base-url "https://docs.myproject.com" \
  --clean-urls \
  --output dist
```
