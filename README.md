# MinimalDoc

A minimal static site generator for documentation. Fast, clean, and easy to use.

## Features

| Category | Features |
|----------|----------|
| **Content** | Markdown (GFM), Frontmatter, Admonitions, Syntax Highlighting (100+ langs), Stale Warnings, Link Checking |
| **Navigation** | Auto-generated from folders, Custom TOC.md, Page TOC with scrollspy, Collapsible sections |
| **Search** | Full-text inverted index, Cmd+K shortcut, Prefix matching, Section-level results |
| **Design** | Dark mode, Responsive, Built-in themes, Custom colors/fonts via YAML |
| **Versioning** | Multi-version docs, Version selector, EOL warnings, Per-version search |
| **i18n** | Multi-locale, Translation files, Fallback content, RTL support |
| **Pages** | Landing, Knowledge Base, Portfolio, Contact, FAQ, Legal, Status, Changelog |
| **API Docs** | OpenAPI/Swagger, Live testing, Code samples (curl/JS/Go/Python/Swift), Schema viewer |
| **SEO** | Sitemap, Open Graph, Twitter Cards, LLM output (llms.txt) |
| **Analytics** | GA4, Plausible, Umami, Matomo, Fathom, Simple Analytics, MinimalDoc Backend, Custom |
| **Backend** | Optional self-hosted server for analytics, feedback, newsletter (cookie-free) |
| **Build** | Single binary, Fast (Go), Config file or CLI flags, Clean URLs option |

See [Feature Index](docs/features/00-index.md) for the complete list with descriptions

## Installation

### From Source

Requires Go 1.24 or higher:

```bash
git clone https://github.com/studiowebux/minimaldoc
cd minimaldoc
go build -o minimaldoc cmd/minimaldoc/main.go
cp minimaldoc /usr/local/bin/
```

### Binary Releases

Download pre-built binaries from the [Releases page](https://github.com/studiowebux/minimaldoc/releases).

## Quick Start

### 1. Initialize a New Site

```bash
minimaldoc init my-docs
cd my-docs
```

This creates a sample documentation structure:

```
my-docs/
├── index.md
├── 01-getting-started/
│   ├── 01-introduction.md
│   └── 02-installation.md
├── 02-guides/
│   └── 01-quickstart.md
└── 03-api/
    └── 01-reference.md
```

### 2. Build Your Site

```bash
minimaldoc build
```

Output will be in the `public/` directory.

### 3. Preview

Open `public/index.html` in your browser.

## Usage

### Commands

#### `build`

Generate static HTML site from Markdown files.

```bash
minimaldoc build [docs-directory] [flags]
```

**Flags:**

- `--output, -o` - Output directory (default: `public`)
- `--theme, -t` - Theme name (default: `default`)
- `--llms, -l` - Generate llms.txt (default: `true`)
- `--openapi` - Enable OpenAPI/Swagger documentation
- `--openapi-dir` - Directory containing OpenAPI spec files (default: `api`)
- `--status` - Enable status page generation
- `--status-title` - Status page title
- `--status-path` - Status page output path (default: `status`)
- `--clean-urls` - Use clean URLs `/page/` instead of `/page.html`
- `--title` - Site title (default: `Documentation`)
- `--description` - Site description
- `--base-url` - Base URL for the site
- `--link-check` - Link check mode: `error`, `warn`, `ignore` (default: `warn`)
- `--check-external` - Validate external URLs (slower)

**Example:**

```bash
minimaldoc build ./docs --output dist --title "My Project" --description "Project documentation"
```

#### Single File Build

Build a single markdown file as a standalone page:

```bash
minimaldoc build README.md
```

This generates `index.html` from the specified file, useful for quick one-pagers.

#### `init`

Initialize a new documentation site with example files.

```bash
minimaldoc init [directory]
```

#### `version`

Display version information and check for updates:

```bash
minimaldoc version           # Show current version
minimaldoc version --check   # Check GitHub for newer releases
```

## Configuration

MinimalDoc can be configured using either CLI flags or a `config.yaml` file in your docs directory.

### Using config.yaml

Create a `config.yaml` file in your docs root:

```yaml
title: My Project Documentation
description: Comprehensive documentation for my project
base_url: https://example.com/docs
theme: default
enable_llms: true
clean_urls: false
entrypoint: README.md  # Optional: use a custom file as homepage instead of index.md

openapi:
  enabled: true
  spec_files:
    - "api/openapi.yaml"
    - "api/v2/api.yaml"
  spec_urls:
    - "https://api.example.com/openapi.json"
  default_view: "path"  # Options: path, tag, flat
  sync_on_build: false
  cache_dir: ".openapi-cache"
  enable_testing: true
  enable_export: true
  enable_code_samples: true
  lazy_load_chunk_size: 51200

status:
  enabled: true
  title: "Service Status"
  path: "status"
  show_history: true
  rss_enabled: true

social_links:
  - name: GitHub
    url: https://github.com/org/project
    icon: github

footer:
  copyright: "2026 Company"
  hideVersion: false  # Set true to hide version badge
```

**Benefits:**
- No need to pass flags every build
- Commit configuration with your docs
- Easier team collaboration

**Note:** CLI flags always override config.yaml values.

## Documentation Structure

### Navigation

MinimalDoc supports two navigation modes:

#### 1. Folder-Based Navigation (Default)

Automatically generates navigation from your folder structure:

```
docs/
├── index.md                    # Home page
├── 01-getting-started/         # Section 1
│   ├── 01-intro.md            # Page 1.1
│   └── 02-install.md          # Page 1.2
└── 02-guides/                  # Section 2
    └── 01-quickstart.md       # Page 2.1
```

**Ordering:**

- Use numbered prefixes (e.g., `01-`, `02-`) to control order
- Numbers are stripped from display titles
- Unnumbered files appear after numbered ones
- Use `menu_order` in frontmatter to override (0 = first, higher = later)

#### 2. Custom Navigation with TOC.md

For full control, create a `TOC.md` file in your docs root:

```markdown
- [Welcome](index.md)
- Getting Started
  - [Introduction](01-getting-started/01-introduction.md)
  - [Installation](01-getting-started/02-installation.md)
- [Guides](02-guides/01-quickstart.md)
- API
  - [Reference](03-api/01-reference.md)
```

**Format:**

- Use markdown lists (`-` or `*`)
- `[Title](path/to/file.md)` - Links to pages
- `Section Name` - Creates navigation sections (no link)
- Indent with 2 spaces for nested items
- Order in TOC.md determines display order

### Table of Contents

TOC is automatically generated from headings (h1-h6) with:

- Nested structure
- Active scrollspy highlighting
- Smooth scrolling

### Frontmatter

Add YAML frontmatter to customize pages:

```yaml
---
title: Page Title
description: Page description for SEO
tags: [tag1, tag2]
author: Author Name
date: 2024-01-01

# Navigation overrides
menu_title: Custom Menu Name
menu_order: 5
hidden: false

# SEO
seo:
  title: SEO Title
  description: SEO Description
  keywords: [keyword1, keyword2]
  image: /images/og-image.png
  canonical: https://example.com/page
  noindex: false
  nofollow: false
---
# Page Content

Your markdown content here...
```

### External Links

External links automatically open in new tabs with `target="_blank"`.

### Shareable Anchor Links

All headings get anchor IDs for deep linking:

```markdown
## Getting Started
```

Links to: `#getting-started`

## OpenAPI/Swagger Documentation

MinimalDoc provides interactive API documentation from OpenAPI/Swagger specifications.

### Quick Start

1. Place your OpenAPI spec files (`.yaml` or `.json`) in the `api/` directory:

```
docs/
├── index.md
├── api/
│   ├── openapi.yaml
│   └── another-api.json
└── ...
```

2. Build with OpenAPI support:

```bash
minimaldoc build --openapi
```

### Features

- **Multiple Views** - Organize endpoints by path, tag, or flat list
- **Live Testing** - Interactive API tester with authentication support
- **Schema Viewer** - Collapsible schema explorer for complex types
- **Export** - Generate cURL commands and restcli configurations
- **Markdown Support** - Descriptions render with full markdown formatting
- **$ref Resolution** - Automatic component reference resolution
- **Single Endpoint View** - Spotify-style focused documentation
- **Code Samples** - Auto-generated code examples in curl, JavaScript, Go, Swift and Python
- **Response Tabs** - Tabbed interface for viewing different response status codes
- **Collapsible Navigation** - Navigation groups collapsed by default for large APIs

### Authentication

Supports multiple authentication methods:
- Bearer Token
- API Key
- OAuth 2.0 (Authorization Code flow)

### Customization

```bash
# Use a different directory for OpenAPI specs
minimaldoc build --openapi --openapi-dir specs
```

### Schemas Browser

The Schemas tab provides a dedicated view for browsing all reusable data models defined in your OpenAPI specification:

- **Alphabetical listing** of all schemas from `components/schemas`
- **Type indicators** showing object, array, string, etc.
- **Property viewer** with name, type, required markers, and constraints
- **Nested schema expansion** for complex object hierarchies
- **JSON example generation** from schema definitions
- **Constraint display** showing min/max, patterns, enums, defaults

Click any schema in the navigation to view its full definition with properties and generated examples.

## Status Page

MinimalDoc includes a built-in status page for displaying service health, incidents, and maintenance schedules.

### Quick Start

1. Enable in your `config.yaml`:

```yaml
status:
  enabled: true
  title: "Service Status"
  description: "Current operational status"
  path: "status"
  show_history: true
  rss_enabled: true
```

2. Create the status content directory:

```
docs/
└── __status__/
    ├── components.yaml
    ├── incidents/
    │   └── 2025-01-28-api-outage.md
    └── maintenance/
        └── 2025-02-01-upgrade.md
```

3. Build your site:

```bash
minimaldoc build docs
```

### Features

- **Component Status** - Track multiple services with operational/degraded/outage states
- **Incident Timeline** - Markdown-based incidents with update history
- **Scheduled Maintenance** - Display upcoming maintenance windows
- **Uptime Calendar** - 90-day visual uptime grid (GitHub-style)
- **SLA Display** - Track SLA targets with 7d/30d/90d breakdowns
- **Live Health Checks** - Browser-based endpoint polling with latency display
- **RSS Feed** - Subscribe to incident updates
- **JSON API** - Machine-readable status data at `/status/status.json`

### Uptime Tracking

Two modes for uptime data:

```yaml
# Mode 1: Derive from incidents (static)
uptime:
  mode: incidents
  sla_target: 99.9
  period_days: 90

# Mode 2: Fetch from API (dynamic)
uptime:
  mode: api
  endpoint: /uptime
  sla_target: 99.9
```

## Social Links

Add social media links to the sidebar:

```yaml
social_links:
  - name: GitHub
    url: https://github.com/your-org/project
    icon: github
  - name: Discord
    url: https://discord.gg/invite
    icon: discord
```

**Supported icons:** `github`, `twitter`, `linkedin`, `youtube`, `discord`, `mastodon`, `rss`, `email`, `website`

## Landing Pages

Landing pages support two configuration methods:

### YAML Configuration

Define sections in `config.yaml`:

```yaml
landing:
  enabled: true
  hero:
    title: "Documentation Made Simple"
    subtitle: "Build beautiful docs from Markdown"
    buttons:
      - text: "Get Started"
        url: "/getting-started/"
        primary: true
  features:
    title: "Features"
    items:
      - emoji: "~"
        title: "Markdown-based"
        description: "Write in Markdown"
```

### Markdown Files

Create markdown files in `__landing__/` directory:

```
docs/
  __landing__/
    01-hero.md
    02-features.md
    03-steps.md
```

Each file uses frontmatter:

```yaml
---
title: Documentation Made Simple
description: Build beautiful docs from Markdown
section: hero
buttons:
  - text: "Get Started"
    url: "/getting-started/"
    primary: true
---
```

Supported sections: `hero`, `features`, `steps`, `cta`, `testimonials`, `opensource`, `links`.

Markdown files override YAML config (YAML serves as defaults).

## Knowledge Base

Create a self-service support hub with categories and dedicated search.

### Directory Structure

```
docs/
  __kb__/
    getting-started/
      01-quick-start.md
      02-configuration.md
    troubleshooting/
      01-common-issues.md
```

Subdirectories become categories. Markdown files are articles.

### Configuration

```yaml
knowledgebase:
  enabled: true
  title: "Knowledge Base"
  description: "Find answers and solutions"
  path: "kb"
  search:
    enabled: true
    placeholder: "Search articles..."
  categories:
    getting-started:
      name: "Getting Started"
      description: "Installation and setup"
      icon: "rocket"
      order: 1
    troubleshooting:
      name: "Troubleshooting"
      icon: "wrench"
      order: 2
```

### Article Frontmatter

```yaml
---
title: "Quick Start Guide"
description: "Get up and running in minutes"
tags: ["beginner", "setup"]
order: 1
---
```

### Features

- Category landing page with article counts
- Scoped search (`kb-search.json`)
- Breadcrumb navigation
- Related articles by tags
- Previous/Next navigation within categories

## Link Checking

Validate links during build to catch broken references before deployment.

### Modes

```bash
# Warn on broken links (default)
minimaldoc build docs

# Fail build on broken links (CI/CD)
minimaldoc build docs --link-check=error

# Skip link checking
minimaldoc build docs --link-check=ignore

# Also check external URLs
minimaldoc build docs --check-external
```

### Configuration

```yaml
link_check:
  enabled: true
  mode: "warn"           # error, warn, ignore
  check_external: false  # Validate external URLs
  external_timeout: 5    # Seconds
  ignore_patterns:
    - "/api/*"           # Skip generated API docs
  allowed_broken:
    - "example.md"       # Documentation examples
```

### What Gets Checked

- Internal page links (`.md`, `.html`)
- Anchor links (`#section`)
- Asset links (images, files)
- External URLs (optional)

### Output

```
Checking links...

Warning: Found 2 broken links

docs/guides/deploy.md:
  Line 42: /api/reference.html (file not found)
  Line 78: #instalation (anchor not found, did you mean #installation?)

Link check completed with warnings: 2 broken links in 1 files
```

## Multi-Version Documentation

Maintain documentation for multiple versions of your software with version-specific content and overrides.

### Configuration

```yaml
versions:
  enabled: true
  default: "v2"
  list:
    - name: "v2"
      label: "2.x (Latest)"
      path: ""              # Default version at root
    - name: "v1"
      label: "1.x (LTS)"
      path: "v1"
    - name: "v0"
      label: "0.x (EOL)"
      path: "v0"
      eol: "2025-01-01"     # End of life date
  selector:
    position: "header"
    show_eol_warning: true
```

### Directory Structure

Shared content lives in the main `docs/` directory. Version-specific overrides go in `__versions__/`:

```
docs/
  getting-started/
    install.md              # Shared across all versions
  features/
    new-feature.md          # Can use frontmatter to limit versions
  __versions__/
    v1/
      getting-started/
        install.md          # Overrides install.md for v1
      migration/
        upgrade.md          # Only appears in v1
```

### Frontmatter

Control version visibility with frontmatter:

```yaml
---
title: New API Feature
versions:
  - v2                      # Only show in v2
since: "v2.0"               # Badge: "Since v2.0"
deprecated_in: "v3.0"       # Badge: "Deprecated in v3.0"
---
```

### URL Structure

```
/docs/getting-started/      # Default version (v2)
/v1/docs/getting-started/   # v1 version
/v0/docs/getting-started/   # v0 version (with EOL warning)
```

### Features

- Version selector dropdown in sidebar
- EOL badges and warning banners
- Per-version search indexes
- Shared content with version-specific overrides
- `versions.json` metadata file for client-side switching

## Custom Themes

Configure colors, fonts, and backgrounds directly in `config.yaml` without creating theme files.

### Colors

```yaml
theme_config:
  colors:
    light:
      bg_primary: "#ffffff"
      bg_secondary: "#f8fafc"
      text_primary: "#1a1a1a"
      accent_primary: "#2563eb"
      link_color: "#0066cc"
    dark:
      bg_primary: "#0f172a"
      bg_secondary: "#1e293b"
      text_primary: "#f8fafc"
      accent_primary: "#3b82f6"
      link_color: "#60a5fa"
```

### Fonts

```yaml
theme_config:
  fonts:
    heading: "Inter"
    body: "Inter"
    code: "JetBrains Mono"
    google_url: "https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&family=JetBrains+Mono&display=swap"
```

### Hero Background

```yaml
theme_config:
  hero:
    background_image: "/images/hero-bg.jpg"
    background_overlay: "rgba(0,0,0,0.6)"
    text_align: "center"
    min_height: "80vh"
```

### Available Color Variables

| Variable | Description |
|----------|-------------|
| `bg_primary` | Main background |
| `bg_secondary` | Sidebar, cards |
| `bg_tertiary` | Nested elements |
| `bg_code` | Code blocks |
| `bg_hover` | Hover states |
| `text_primary` | Main text |
| `text_secondary` | Secondary text |
| `text_muted` | Muted text |
| `border_primary` | Main borders |
| `accent_primary` | Buttons, links |
| `accent_hover` | Accent hover |
| `link_color` | Link text |
| `link_hover` | Link hover |

### Code Samples

When `enable_code_samples: true` is set, each endpoint displays auto-generated code examples in a right sidebar:

- **curl** - Command-line HTTP requests
- **JavaScript** - Using the Fetch API
- **Swift** - Using URLSession
- **Go** - Using net/http package
- **Python** - Using the requests library

Code samples include:
- Full URL with path parameters
- Authentication headers based on security schemes (Bearer, API Key, OAuth2)
- Request body examples generated from schema

## Default Theme

The default theme uses soft, eye-friendly colors:

**Light Mode:** Background `#fafafa`, Text `#1a1a1a`

**Dark Mode:** Background `#1a1a1a`, Text `#ffffff`, Secondary `#2a2a2a`

For config-based customization, see [Custom Themes](#custom-themes) above.

## Analytics

MinimalDoc supports multiple analytics providers that can run simultaneously.

### Quick Start

```yaml
analytics:
  enabled: true
  providers:
    - type: plausible
      enabled: true
      config:
        domain: "docs.example.com"
```

### Supported Providers

| Provider | Type | Required Config |
|----------|------|-----------------|
| MinimalDoc Backend | `minimaldoc` | `endpoint`, `site_id` |
| Google Analytics 4 | `ga4` | `measurement_id` |
| Plausible | `plausible` | `domain` |
| Umami | `umami` | `website_id`, `src` |
| Matomo | `matomo` | `url`, `site_id` |
| Fathom | `fathom` | `site_id` |
| Simple Analytics | `simple` | (none) |
| Custom | `custom` | `src` |

### Custom Provider

For any analytics service, use the `custom` type with arbitrary attributes:

```yaml
- type: custom
  enabled: true
  config:
    src: "https://analytics.example.com/tracker.js"
    defer: "true"
    attrs:
      data-site-id: "my-site"
      data-cookieless: "true"
```

## LLM-Friendly Output

Minimal Doc generates `llms.txt` - a single file containing all documentation content with navigation structure.

Perfect for:

- AI assistants
- Documentation search
- Content indexing

## Backend Server (Optional)

MinimalDoc includes an optional self-hosted backend for dynamic features. The CLI and static site generation work without it.

### Features

| Feature | Description |
|---------|-------------|
| Analytics | Cookie-free, privacy-first page tracking |
| Feedback | Page rating widget with comments |
| Newsletter | Email subscription with double opt-in verification |
| Email | SMTP support with verification and welcome emails |
| Admin Portal | Web dashboard for all features |

### Quick Start

```bash
# Build the server
make build-server

# Run with SQLite (mock email for testing)
DB_DRIVER=sqlite DATABASE_URL=./data.db \
AUTH_JWT_SECRET=your-32-char-secret \
EMAIL_PROVIDER=mock \
./minimaldoc-server

# Run with SMTP email
DB_DRIVER=sqlite DATABASE_URL=./data.db \
AUTH_JWT_SECRET=your-32-char-secret \
EMAIL_PROVIDER=smtp SMTP_HOST=smtp.example.com \
SMTP_USER=user SMTP_PASS=pass \
EMAIL_BASE_URL=https://api.example.com \
./minimaldoc-server

# Bootstrap (first run)
curl -X POST http://localhost:8080/api/bootstrap \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"pass","site_name":"My Docs"}'
```

### Enable in Static Site

```yaml
# config.yaml
analytics:
  enabled: true
  providers:
    - type: minimaldoc
      enabled: true
      config:
        endpoint: "http://localhost:8080"
        site_id: "your-site-id"
        features: "analytics,feedback,newsletter"
```

See [Backend Documentation](docs/features/12-backend.md) for full setup guide.

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

see [LICENSE](LICENSE) file for details.

## Credits

Created by [Studio Webux](https://studiowebux.com)
