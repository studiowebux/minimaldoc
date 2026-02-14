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

# Create and build site
minimaldoc init my-docs && cd my-docs && minimaldoc build
```

## Documentation

| Topic | Link |
|-------|------|
| Getting Started | [docs/getting-started/](docs/getting-started/) |
| Configuration | [docs/core-concepts/](docs/core-concepts/) |
| Features | [docs/features/](docs/features/) |
| OpenAPI | [docs/openapi/](docs/openapi/) |
| Status Page | [docs/status-page/](docs/status-page/) |
| Backend Server | [docs/features/12-backend.md](docs/features/12-backend.md) |
| Deployment | [deploy/](deploy/) |

## License

See [LICENSE](LICENSE) file.

## Credits

Created by [Studio Webux](https://studiowebux.com)
