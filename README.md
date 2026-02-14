# MinimalDoc

A minimal static site generator for documentation. Fast, clean, and easy to use.

## Features

| Category | Features |
|----------|----------|
| **Content** | Markdown (GFM), Frontmatter, Admonitions, Syntax Highlighting, Link Checking |
| **Navigation** | Auto-generated from folders, Custom TOC.md, Page TOC with scrollspy |
| **Search** | Full-text index, Cmd+K shortcut, Section-level results |
| **Design** | Dark/Light mode, Responsive, Custom themes via YAML |
| **Versioning** | Multi-version docs, Version selector, EOL warnings |
| **i18n** | Multi-locale, Translation files, RTL support |
| **Pages** | Landing, Knowledge Base, Portfolio, FAQ, Legal, Status, Changelog |
| **API Docs** | OpenAPI/Swagger, Live testing, Code samples, Schema viewer |
| **SEO** | Sitemap, Open Graph, Twitter Cards, llms.txt |
| **Analytics** | GA4, Plausible, Umami, Matomo, MinimalDoc Backend |
| **Backend** | Optional server: analytics, feedback, newsletter, blog CMS |

## Quick Start

```bash
# Install
git clone https://github.com/studiowebux/minimaldoc
cd minimaldoc && go build -o minimaldoc cmd/minimaldoc/main.go

# Create site
minimaldoc init my-docs
cd my-docs

# Build
minimaldoc build

# Output in public/
```

## Usage

```bash
minimaldoc build [docs-dir] [flags]

# Common flags
--output, -o     Output directory (default: public)
--theme, -t      Theme name
--openapi        Enable OpenAPI docs
--status         Enable status page
--clean-urls     Use /page/ instead of /page.html
--link-check     Mode: error, warn, ignore
```

## Configuration

Create `config.yaml` in your docs root:

```yaml
title: My Documentation
description: Project docs
base_url: https://docs.example.com

openapi:
  enabled: true
  spec_files: ["api/openapi.yaml"]

analytics:
  enabled: true
  providers:
    - type: plausible
      config:
        domain: "docs.example.com"
```

## Backend Server (Optional)

Self-hosted server for dynamic features:

```bash
# Build and run
make build-server
AUTH_JWT_SECRET=your-secret BOOTSTRAP_TOKEN=my-token ./minimaldoc-server

# Bootstrap
curl -X POST localhost:8080/api/bootstrap \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "securepass",
    "site_name": "Docs",
    "bootstrap_token": "my-token"
  }'
```

Features: Cookie-free analytics, feedback widget, newsletter, blog CMS with RBAC.

## Documentation

Full documentation available at [docs/](docs/):

| Topic | Link |
|-------|------|
| Getting Started | [docs/getting-started/](docs/getting-started/) |
| Features Index | [docs/features/00-index.md](docs/features/00-index.md) |
| Search | [docs/features/02-search.md](docs/features/02-search.md) |
| Theming | [docs/features/03-theming.md](docs/features/03-theming.md) |
| OpenAPI | [docs/openapi/](docs/openapi/) |
| Multi-Version | [docs/features/09-multi-version.md](docs/features/09-multi-version.md) |
| i18n | [docs/features/10-i18n.md](docs/features/10-i18n.md) |
| Analytics | [docs/features/11-analytics.md](docs/features/11-analytics.md) |
| Backend Server | [docs/features/12-backend.md](docs/features/12-backend.md) |
| Status Page | [docs/status-page/](docs/status-page/) |
| Guides | [docs/guides/](docs/guides/) |

## License

See [LICENSE](LICENSE) file.

## Credits

Created by [Studio Webux](https://studiowebux.com)
