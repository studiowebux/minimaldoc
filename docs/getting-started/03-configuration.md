---
title: Configuration
description: Complete config.yaml reference for MinimalDoc
tags:
  - getting-started
  - configuration
---

# Configuration

MinimalDoc is configured via `config.yaml` in your docs root. CLI flags override config values.

## Precedence

1. CLI flags (highest)
2. config.yaml
3. Defaults (lowest)

## Minimal Config

The smallest useful config:

```yaml
title: My Docs
```

Everything else has sensible defaults. Add sections as needed.

## Full Config Reference

### Basic Settings

```yaml
title: My Documentation
description: Project documentation site
base_url: https://docs.example.com
author: Your Name
theme: default
dark_mode: false
enable_llms: true
enable_search: true
clean_urls: false
entrypoint: index.md
nav_depth: 0
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | `"Documentation"` | Site title shown in header and meta tags |
| `description` | string | `""` | Site description for SEO meta tags |
| `base_url` | string | `""` | Base URL for absolute links, sitemap, and RSS feeds |
| `author` | string | `""` | Default author for meta tags |
| `theme` | string | `"default"` | Theme name |
| `dark_mode` | bool | `false` | Start in dark mode by default |
| `enable_llms` | bool | `true` | Generate `llms.txt`, `llms-full.txt`, and per-page `.html.md` files |
| `enable_search` | bool | `false` | Enable client-side search with Cmd+K shortcut |
| `clean_urls` | bool | `false` | Use `/page/` instead of `/page.html` |
| `entrypoint` | string | `""` | Custom homepage file (default: `index.md`) |
| `nav_depth` | int | `0` | Max depth for navigation tree (0 = unlimited) |

### Theme Customization

Override colors, fonts, and hero styling:

```yaml
theme_config:
  colors:
    light:
      bg_primary: "#ffffff"
      bg_secondary: "#f8f9fa"
      text_primary: "#1a1a2e"
      link_color: "#2563eb"
      accent_primary: "#2563eb"
    dark:
      bg_primary: "#0f0f23"
      bg_secondary: "#1a1a2e"
      text_primary: "#e2e8f0"
      link_color: "#60a5fa"
      accent_primary: "#60a5fa"
  fonts:
    heading: "Inter"
    body: "Inter"
    code: "JetBrains Mono"
    google_url: "https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap"
  hero:
    background_image: "/images/hero-bg.jpg"
    background_overlay: "rgba(0,0,0,0.6)"
    text_align: center
    min_height: "80vh"
```

**Color fields** (available in both `light` and `dark`):

| Field | Description |
|-------|-------------|
| `bg_primary` | Main background |
| `bg_secondary` | Sidebar, cards |
| `bg_tertiary` | Code blocks, subtle areas |
| `bg_code` | Inline code background |
| `bg_hover` | Hover state background |
| `text_primary` | Main text |
| `text_secondary` | Subtitles, descriptions |
| `text_tertiary` | Muted text |
| `text_muted` | Least prominent text |
| `border_primary` | Main borders |
| `border_secondary` | Subtle borders |
| `accent_primary` | Buttons, active states |
| `accent_hover` | Accent hover |
| `link_color` | Link text |
| `link_hover` | Link hover |
| `color_success` | Success indicators |
| `color_warning` | Warning indicators |
| `color_danger` | Error indicators |
| `color_info` | Info indicators |

### OpenAPI

Generate interactive API documentation from OpenAPI/Swagger specs:

```yaml
openapi:
  enabled: true
  spec_files:
    - petstore.yaml
    - internal-api.json
  spec_urls:
    - https://api.example.com/openapi.json
  default_view: path
  sync_on_build: false
  cache_dir: .openapi-cache
  enable_testing: true
  enable_export: true
  enable_code_samples: true
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable OpenAPI docs generation |
| `spec_files` | []string | `[]` | Local spec files relative to docs root |
| `spec_urls` | []string | `[]` | Remote spec URLs to fetch |
| `default_view` | string | `"path"` | Organization: `path`, `tag`, or `flat` |
| `sync_on_build` | bool | `false` | Fetch remote specs on every build |
| `cache_dir` | string | `".openapi-cache"` | Cache directory for remote specs |
| `enable_testing` | bool | `true` | Show "Try It" interface for endpoints |
| `enable_export` | bool | `true` | Enable cURL export |
| `enable_code_samples` | bool | `true` | Show code samples (curl, fetch, Go, Python) |

Spec files go in an `api/` subdirectory by default (configurable via `--openapi-dir` CLI flag).

### MCP Server Documentation

Generate reference pages from `.mcp.json` manifests (like Swagger for MCP):

```yaml
mcp:
  enabled: true
  spec_files:
    - "*.mcp.json"
  path: mcp
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable MCP docs generation |
| `spec_files` | []string | `[]` | JSON manifest file paths/globs relative to docs root |
| `path` | string | `"mcp"` | Output URL path |

Place `.mcp.json` files in an `mcp/` subdirectory by default (configurable via `--mcp-dir` CLI flag).

### Status Page

Track service health with incidents, maintenance windows, and uptime:

```yaml
status:
  enabled: true
  title: Service Status
  description: Current system status
  path: status
  show_history: true
  history_months: 12
  rss_enabled: true
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable status page |
| `title` | string | `"Service Status"` | Page title |
| `description` | string | `""` | Page description |
| `path` | string | `"status"` | Output path |
| `show_history` | bool | `true` | Show incident history |
| `history_months` | int | `12` | Months of history to display |
| `rss_enabled` | bool | `true` | Generate RSS feed |

Status data is defined as markdown files in `__status__/` directory. See the status page guide for file format.

### Changelog

Version history with RSS feed:

```yaml
changelog:
  enabled: true
  title: Changelog
  description: Release history
  path: changelog
  rss_enabled: true
  repository: https://github.com/org/repo
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable changelog |
| `title` | string | `"Changelog"` | Page title |
| `description` | string | `""` | Page description |
| `path` | string | `"changelog"` | Output path |
| `rss_enabled` | bool | `true` | Generate RSS feed |
| `repository` | string | `""` | GitHub repository URL for version compare links |

Releases are defined as markdown files in `__changelog__/releases/` directory.

### Landing Page

Marketing-style homepage with hero, features, steps, and CTA sections:

```yaml
landing:
  enabled: true
  nav:
    - text: Docs
      url: /getting-started/installation.html
    - text: Features
      url: /features/overview.html
  hero:
    title: Documentation Made Simple
    subtitle: A minimal static site generator
    buttons:
      - text: Get Started
        url: /getting-started/installation.html
        primary: true
      - text: GitHub
        url: https://github.com/org/repo
        primary: false
  features:
    title: Features
    items:
      - emoji: "~"
        title: Markdown-based
        description: Write in plain Markdown with YAML frontmatter.
      - emoji: "@"
        title: Lightning Fast
        description: Sub-second builds powered by Go.
  steps:
    title: Quick Start
    items:
      - title: Install
        description: Download the binary
        code: "curl -sSL https://example.com/install.sh | bash"
      - title: Build
        description: Generate your site
        code: "minimaldoc build"
  cta:
    title: Ready to start?
    description: Get up and running in minutes.
    buttons:
      - text: Download
        url: https://github.com/org/repo/releases
        primary: true
  opensource:
    title: Open Source
    description: GPL-3.0 Licensed.
    links:
      - text: GitHub
        url: https://github.com/org/repo
  links:
    title: Resources
    items:
      - icon: github
        title: GitHub
        description: Source code
        url: https://github.com/org/repo
```

Landing sections can also be defined as markdown files in a `__landing__/` directory. Markdown files override YAML config values.

Additional landing section types: `testimonials`, `image_text`, `text_blocks`, `links_grid`.

### Portfolio

Project showcase with filterable tags:

```yaml
portfolio:
  enabled: true
  title: Projects
  description: Sites built with our tool
  path: projects
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable portfolio |
| `title` | string | `"Projects"` | Page title |
| `description` | string | `""` | Page description |
| `path` | string | `"projects"` | Output path |

Projects are defined as markdown files in `__portfolio__/` directory with frontmatter for tags, images, and URLs.

### Knowledge Base

Categorized articles with search:

```yaml
knowledgebase:
  enabled: true
  title: Knowledge Base
  description: Find answers and guides
  path: kb
  search:
    enabled: true
    placeholder: "Search articles..."
  categories:
    getting-started:
      name: Getting Started
      description: Installation and setup
      icon: rocket
      order: 1
    troubleshooting:
      name: Troubleshooting
      description: Common issues and solutions
      icon: wrench
      order: 2
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable knowledge base |
| `title` | string | `"Knowledge Base"` | Page title |
| `description` | string | `""` | Page description |
| `path` | string | `"kb"` | Output path |
| `search.enabled` | bool | `true` | Enable search within KB |
| `search.placeholder` | string | `""` | Search input placeholder text |
| `categories` | map | `{}` | Category definitions keyed by directory name |

Articles are markdown files in `__kb__/<category>/` directories.

### Roadmap

Visual roadmap with board or timeline layout:

```yaml
roadmap:
  enabled: true
  title: Roadmap
  description: What we're building next
  layout: board
  show_versions: true
  path: roadmap
  columns:
    - id: planned
      label: Planned
    - id: in_progress
      label: In Progress
    - id: shipped
      label: Shipped
  items:
    - title: Dark mode
      description: System-aware theme switching
      status: shipped
      version: v1.0
      tags: [ui, theme]
      shipped_date: "2025-01-15"
      changelog_url: "/changelog/#v1.0"
    - title: Plugin system
      description: Extensible plugin architecture
      status: planned
      version: v2.0
      tags: [architecture]
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable roadmap page |
| `title` | string | `"Roadmap"` | Page title |
| `description` | string | `""` | Page description |
| `layout` | string | `"board"` | Layout: `board` or `timeline` |
| `show_versions` | bool | `false` | Show version badges |
| `path` | string | `"roadmap"` | Output path |
| `columns` | array | `[]` | Board columns with `id` and `label` |
| `items[].status` | string | | Must match a column `id` |
| `items[].tags` | []string | | Filterable tags |
| `items[].shipped_date` | string | | Date shipped (YYYY-MM-DD) |
| `items[].changelog_url` | string | | Link to changelog entry |

### Waitlist

Pre-launch landing page with email signup:

```yaml
waitlist:
  enabled: true
  title: Coming Soon
  tagline: Sign up for early access
  newsletter_endpoint: https://api.example.com/subscribe
  site_id: my-site
  success_message: "Thanks! We'll be in touch."
  privacy_url: /legal/privacy.html
  social_links:
    - name: github
      url: https://github.com/org/repo
    - name: discord
      url: https://discord.gg/invite
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable waitlist page |
| `title` | string | `""` | Page title |
| `tagline` | string | `""` | Subtitle text |
| `newsletter_endpoint` | string | `""` | POST endpoint for email signup |
| `site_id` | string | `""` | Site identifier sent with signup |
| `success_message` | string | `""` | Message shown after signup |
| `privacy_url` | string | `""` | Link to privacy policy |
| `social_links` | array | `[]` | Social links with `name` and `url` |

### FAQ

Frequently asked questions page:

```yaml
faq:
  enabled: true
  title: FAQ
  description: Frequently asked questions
  path: faq
  categories:
    - name: General
      items:
        - question: What is this?
          answer: A documentation generator.
        - question: Is it free?
          answer: Yes, GPL-3.0 licensed.
    - name: Technical
      items:
        - question: What languages are supported?
          answer: 100+ via Chroma syntax highlighting.
```

FAQ items can also be defined as markdown files in `__faq__/` directory with frontmatter containing `question`, `category`, and `tags` fields.

### Contact

Contact page with info cards:

```yaml
contact:
  enabled: true
  title: Contact
  description: Get in touch
  path: contact
  email: hello@example.com
  info:
    - icon: mail
      text: hello@example.com
    - icon: location
      text: Montreal, Canada
```

### Legal

Legal pages (privacy policy, terms of service):

```yaml
legal:
  enabled: true
  path: legal
  footer_group: Legal
```

Legal pages are markdown files in `__legal__/` directory. Each file becomes a page under the configured path and is automatically added to the footer group.

### Stale Warning

Show warnings on outdated content:

```yaml
stale_warning:
  enabled: true
  threshold_days: 365
  message: "This page may be outdated."
  show_update_date: true
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable stale warnings |
| `threshold_days` | int | `365` | Days before content is flagged stale |
| `message` | string | `""` | Custom warning message (auto-generated if empty) |
| `show_update_date` | bool | `true` | Show last modified date |

### PDF Export

Add an "Export PDF" button to each page:

```yaml
pdf_export:
  enabled: true
  page_break_level: 2
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable PDF export button |
| `page_break_level` | int | `1` | Heading level for page breaks: 0=none, 1=h1 only, 2=h1+h2 |

### Claude Assist

Add an "Ask Claude" button that opens Claude with the page context:

```yaml
claude_assist:
  enabled: true
  label: Ask Claude
  prompt: "Answer based on this documentation:"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable Claude assist button |
| `label` | string | `"Ask Claude"` | Button label text |
| `prompt` | string | `""` | Custom prompt prefix (optional) |

### Versioned Documentation

Serve multiple versions of your docs:

```yaml
versions:
  enabled: true
  default: v1.2
  list:
    - name: v1.2
      label: "1.2.x (Latest)"
      path: ""
    - name: v1.1
      label: "1.1.x"
      path: v1.1
    - name: v1.0
      label: "1.0.x (EOL)"
      path: v1.0
      eol: "2025-06-01"
  selector:
    position: header
    show_eol_warning: true
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable multi-version docs |
| `default` | string | `""` | Default version slug |
| `list[].name` | string | | Version identifier |
| `list[].label` | string | | Display label in version selector |
| `list[].path` | string | | URL path prefix (empty = root) |
| `list[].eol` | string | | End-of-life date (YYYY-MM-DD) |
| `selector.position` | string | `"header"` | `header` or `sidebar` |
| `selector.show_eol_warning` | bool | `false` | Show warning badge for EOL versions |

Versioned content goes in `__versions__/<version-name>/` directories.

### Internationalization (i18n)

Serve docs in multiple languages:

```yaml
i18n:
  enabled: true
  default_locale: en
  hide_default_locale: true
  show_untranslated: true
  fallback: en
  locales:
    - code: en
      name: English
      direction: ltr
    - code: fr
      name: Francais
      direction: ltr
    - code: ar
      name: Arabic
      direction: rtl
  selector:
    position: header
    show_flags: false
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable i18n |
| `default_locale` | string | `"en"` | Default locale code |
| `hide_default_locale` | bool | `false` | Omit locale prefix for default (use `/page.html` not `/en/page.html`) |
| `show_untranslated` | bool | `false` | Show fallback content with untranslated badge |
| `fallback` | string | `""` | Fallback locale for missing translations |
| `locales[].code` | string | | Locale code (e.g., `en`, `fr`, `de`) |
| `locales[].name` | string | | Display name |
| `locales[].direction` | string | | Text direction: `ltr` or `rtl` |
| `selector.position` | string | `"header"` | `header` or `sidebar` |
| `selector.show_flags` | bool | `false` | Show flag icons |

Translated content goes in `__i18n__/<locale-code>/` directories mirroring the main docs structure.

### Analytics

Add analytics tracking:

```yaml
analytics:
  enabled: true
  providers:
    - type: plausible
      enabled: true
      config:
        domain: docs.example.com
    - type: ga4
      enabled: true
      config:
        measurement_id: G-XXXXXXXXXX
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable analytics |
| `providers[].type` | string | | Provider: `ga4`, `plausible`, `umami`, `matomo`, `fathom`, `simple`, `custom` |
| `providers[].enabled` | bool | | Enable this provider |
| `providers[].config` | map | | Provider-specific configuration |

### Link Checking

Validate internal and external links during build:

```yaml
link_check:
  enabled: true
  mode: warn
  check_external: false
  external_timeout: 10
  ignore_patterns:
    - "*/archive/*"
  allowed_broken:
    - "/legacy/old-page.html"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable link checking |
| `mode` | string | `"warn"` | `error` (fail build) or `warn` (report only) |
| `check_external` | bool | `false` | Check external URLs |
| `external_timeout` | int | `10` | Timeout in seconds for external checks |
| `ignore_patterns` | []string | `[]` | Glob patterns to skip |
| `allowed_broken` | []string | `[]` | Known broken links to ignore |

### Footer

Configure the site-wide footer:

```yaml
footer:
  copyright: "2026 Company. GPL-3.0 License."
  hideVersion: false
  links:
    - title: Documentation
      items:
        - text: Getting Started
          url: /getting-started/installation.html
        - text: Configuration
          url: /getting-started/configuration.html
    - title: Resources
      items:
        - text: GitHub
          url: https://github.com/org/repo
  social:
    - name: GitHub
      url: https://github.com/org/repo
      icon: github
    - name: Discord
      url: https://discord.gg/invite
      icon: discord
  badges:
    - text: "Powered by MinimalDoc"
      url: https://minimaldoc.com
```

### Social Links (Sidebar)

Social links shown in the sidebar:

```yaml
social_links:
  - name: GitHub
    url: https://github.com/org/repo
    icon: github
  - name: Discord
    url: https://discord.gg/invite
    icon: discord
```

Available icons: `github`, `twitter`, `linkedin`, `youtube`, `discord`, `mastodon`, `rss`, `email`, `website`.

### Custom Fields

Pass arbitrary data to templates:

```yaml
custom:
  ga_id: G-XXXXXXXXXX
  show_banner: true
```

Access in templates via `{{.Site.Config.Custom.ga_id}}`.

## CLI Flags

CLI flags override config.yaml values:

```bash
minimaldoc build [docs-dir] [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `public` | Output directory |
| `--theme` | `-t` | `default` | Theme name |
| `--title` | | `Documentation` | Site title |
| `--description` | | | Site description |
| `--base-url` | | | Base URL |
| `--llms` | `-l` | `true` | Generate llms.txt |
| `--clean-urls` | | `false` | Use clean URLs |
| `--openapi` | | `false` | Enable OpenAPI |
| `--openapi-dir` | | `api` | OpenAPI specs directory |
| `--mcp` | | `false` | Enable MCP docs |
| `--mcp-dir` | | `mcp` | MCP manifests directory |
| `--status` | | `false` | Enable status page |
| `--status-title` | | `Service Status` | Status page title |
| `--status-path` | | `status` | Status page path |
| `--changelog` | | `false` | Enable changelog |
| `--changelog-title` | | `Changelog` | Changelog page title |
| `--stale-warning` | | `false` | Enable stale warnings |
| `--stale-threshold` | | `365` | Stale threshold in days |
| `--link-check` | | `warn` | Link check mode: `error`, `warn`, `ignore` |

## Common Patterns

### Docs-only (minimal)

```yaml
title: API Reference
base_url: https://docs.example.com
enable_search: true
```

### Full documentation site

```yaml
title: My Project
base_url: https://docs.example.com
enable_search: true
enable_llms: true
clean_urls: true

openapi:
  enabled: true
  spec_files: [openapi.yaml]

changelog:
  enabled: true
  repository: https://github.com/org/repo

status:
  enabled: true

stale_warning:
  enabled: true
  threshold_days: 180

pdf_export:
  enabled: true

claude_assist:
  enabled: true
```

### Marketing site with docs

```yaml
title: My Product
base_url: https://example.com

landing:
  enabled: true
  hero:
    title: Build Faster
    subtitle: Ship documentation in minutes
    buttons:
      - text: Get Started
        url: /getting-started/installation.html
        primary: true

portfolio:
  enabled: true

faq:
  enabled: true

contact:
  enabled: true

roadmap:
  enabled: true
  layout: board
```
