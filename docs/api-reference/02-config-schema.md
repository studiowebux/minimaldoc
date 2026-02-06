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

# Landing Page Configuration
landing:
  enabled: boolean               # Enable landing page (default: false)
  nav: [{text, url}]             # Navigation links
  hero:
    title: string                # Hero title
    subtitle: string             # Hero subtitle
    buttons: [{text, url, primary}]
  features:
    title: string
    items: [{emoji, title, description}]
  steps:
    title: string
    items: [{title, description, code}]
  links:
    title: string
    items: [{icon, title, description, url}]
  opensource:
    title: string
    description: string
    links: [{text, url}]

# Portfolio Configuration
portfolio:
  enabled: boolean               # Enable portfolio (default: false)
  title: string                  # Portfolio page title
  description: string            # Portfolio description
  path: string                   # Output path (default: "projects")

# Contact Configuration
contact:
  enabled: boolean               # Enable contact page (default: false)
  title: string                  # Contact page title
  description: string            # Contact description
  path: string                   # Output path (default: "contact")
  email: string                  # Contact email
  info: [{icon, text}]           # Info items (mail, location, phone)

# FAQ Configuration
faq:
  enabled: boolean               # Enable FAQ page (default: false)
  title: string                  # FAQ page title (default: "FAQ")
  description: string            # FAQ description
  path: string                   # Output path (default: "faq")
  categories:                    # FAQ categories with items
    - name: string               # Category name
      order: integer             # Sort order
      items:                     # Questions in category
        - question: string       # Question text
          answer: string         # Answer text (plain or markdown)
          order: integer         # Sort order within category
          tags: [string]         # Tags for search

# Legal Configuration
legal:
  enabled: boolean               # Enable legal pages (default: false)
  path: string                   # Output path prefix (default: "legal")
  footer_group: string           # Footer section title (default: "Legal")

# Footer Configuration
footer:
  copyright: string              # Copyright text
  links:
    - title: string              # Section title
      items: [{text, url}]       # Links in section
  social: [{name, url, icon}]    # Social links

# Social Links (sidebar)
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

## Landing Page Settings

### landing.enabled

Enable landing page generation.

```yaml
landing:
  enabled: true
```

Type: `boolean`
Default: `false`

### landing.nav

Navigation links in the landing header.

```yaml
landing:
  nav:
    - text: "Docs"
      url: "/getting-started/installation.html"
    - text: "Features"
      url: "/features/overview.html"
```

Type: `array of object`

### landing.hero

Hero section content.

```yaml
landing:
  hero:
    title: "Documentation Made Simple"
    subtitle: "A minimal static site generator"
    buttons:
      - text: "Get Started"
        url: "/getting-started/installation.html"
        primary: true
      - text: "View on GitHub"
        url: "https://github.com/org/repo"
        primary: false
```

### landing.features

Feature grid section.

```yaml
landing:
  features:
    title: "Features"
    items:
      - emoji: "~"
        title: "Fast"
        description: "Sub-second builds"
      - emoji: "@"
        title: "Dark Mode"
        description: "Built-in theme switching"
```

### landing.steps

Quick start steps section.

```yaml
landing:
  steps:
    title: "Quick Start"
    items:
      - title: "Install"
        description: "Download the binary"
        code: "curl -sSL https://..."
      - title: "Build"
        description: "Generate your site"
        code: "minimaldoc build"
```

### landing.links

Resource links section.

```yaml
landing:
  links:
    title: "Resources"
    items:
      - icon: "github"
        title: "GitHub"
        description: "Source code"
        url: "https://github.com/org/repo"
```

### landing.opensource

Open source section.

```yaml
landing:
  opensource:
    title: "Open Source"
    description: "MIT Licensed. Self-host with full control."
    links:
      - text: "GitHub Repository"
        url: "https://github.com/org/repo"
```

## Portfolio Settings

### portfolio.enabled

Enable portfolio page generation.

```yaml
portfolio:
  enabled: true
```

Type: `boolean`
Default: `false`

### portfolio.title

Portfolio page title.

```yaml
portfolio:
  title: "Projects"
```

Type: `string`
Default: `"Portfolio"`

### portfolio.description

Portfolio page description.

```yaml
portfolio:
  description: "Sites built with MinimalDoc"
```

Type: `string`
Default: `""`

### portfolio.path

Output path for portfolio.

```yaml
portfolio:
  path: "projects"
```

Type: `string`
Default: `"projects"`

## Contact Settings

### contact.enabled

Enable contact page generation.

```yaml
contact:
  enabled: true
```

Type: `boolean`
Default: `false`

### contact.title

Contact page title.

```yaml
contact:
  title: "Contact Us"
```

Type: `string`
Default: `"Contact"`

### contact.description

Contact page description.

```yaml
contact:
  description: "Get in touch with our team"
```

Type: `string`
Default: `""`

### contact.path

Output path for contact page.

```yaml
contact:
  path: "contact"
```

Type: `string`
Default: `"contact"`

### contact.email

Contact email address (used by form).

```yaml
contact:
  email: "hello@example.com"
```

Type: `string`
Default: `""`

### contact.info

Contact information items.

```yaml
contact:
  info:
    - icon: "mail"
      text: "hello@example.com"
    - icon: "location"
      text: "Remote, Worldwide"
    - icon: "phone"
      text: "+1 555 123 4567"
```

Type: `array of object`

### Contact Info Icons

| Icon | Purpose |
|------|---------|
| `mail` | Email address |
| `location` | Physical location |
| `phone` | Phone number |

## FAQ Settings

### faq.enabled

Enable FAQ page generation.

```yaml
faq:
  enabled: true
```

Type: `boolean`
Default: `false`

### faq.title

FAQ page title.

```yaml
faq:
  title: "FAQ"
```

Type: `string`
Default: `"FAQ"`

### faq.description

FAQ page description.

```yaml
faq:
  description: "Frequently asked questions"
```

Type: `string`
Default: `"Frequently asked questions"`

### faq.path

Output path for FAQ page.

```yaml
faq:
  path: "faq"
```

Type: `string`
Default: `"faq"`

### faq.categories

FAQ categories with questions and answers.

```yaml
faq:
  categories:
    - name: "General"
      order: 1
      items:
        - question: "What is this?"
          answer: "A documentation generator."
          tags:
            - general
        - question: "Is it free?"
          answer: "Yes, open source under GPL-3.0."
    - name: "Technical"
      order: 2
      items:
        - question: "What languages are supported?"
          answer: "100+ via Chroma syntax highlighting."
```

Type: `array of object`

Categories can also be defined via markdown files in `__faq__/` directory.

## Legal Settings

### legal.enabled

Enable legal pages generation.

```yaml
legal:
  enabled: true
```

Type: `boolean`
Default: `false`

### legal.path

Output path prefix for legal pages.

```yaml
legal:
  path: "legal"
```

Type: `string`
Default: `"legal"`

### legal.footer_group

Footer section title for auto-generated legal links.

```yaml
legal:
  footer_group: "Legal"
```

Type: `string`
Default: `"Legal"`

Legal pages are defined via markdown files in `__legal__/` directory.

## Footer Settings

### footer.copyright

Copyright text displayed in footer.

```yaml
footer:
  copyright: "2025 Your Company. MIT License."
```

Type: `string`
Default: `""`

### footer.links

Footer link sections.

```yaml
footer:
  links:
    - title: "Documentation"
      items:
        - text: "Getting Started"
          url: "/getting-started/installation.html"
        - text: "Features"
          url: "/features/overview.html"
    - title: "Resources"
      items:
        - text: "Status"
          url: "/status/"
        - text: "GitHub"
          url: "https://github.com/org/repo"
```

Type: `array of object`

### footer.social

Social links in footer.

```yaml
footer:
  social:
    - name: GitHub
      url: https://github.com/org/repo
      icon: github
    - name: Twitter
      url: https://twitter.com/handle
      icon: twitter
```

Type: `array of object`
