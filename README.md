# MinimalDoc

A minimal static site generator for documentation. Fast, clean, and easy to use.

> **Work in Progress**
>
> The **static site generator** (`minimaldoc build`) is stable and
> production-ready. The **backend server** (`minimaldoc-server`) is in early
> development — not ready for production use. APIs may change, features are
> largely untested end-to-end.

## Features

| Category       | Features                                                                                              | Status |
| -------------- | ----------------------------------------------------------------------------------------------------- | ------ |
| **Content**    | Markdown (GFM), Frontmatter, Admonitions, Syntax Highlighting, Link Checking, Copy-to-Markdown       | Stable |
| **Navigation** | Auto-generated from folders, Custom TOC.md, Page TOC with scrollspy                                  | Stable |
| **Search**     | Full-text index, Cmd+K shortcut, Fuzzy matching, Section-level results                               | Stable |
| **Design**     | Dark/Light mode, Responsive, Custom themes via YAML, Custom fonts                                    | Stable |
| **Versioning** | Multi-version docs, Version selector, EOL warnings                                                    | Stable |
| **i18n**       | Multi-locale, Translation files, RTL support                                                          | Stable |
| **Pages**      | Landing, Knowledge Base, Portfolio, FAQ, Legal, Status, Changelog, Roadmap                           | Stable |
| **API Docs**   | OpenAPI/Swagger, Live testing, Code samples, Schema viewer                                            | Stable |
| **SEO**        | Sitemap, Open Graph, Twitter Cards, llms.txt                                                          | Stable |
| **Analytics**  | GA4, Plausible, Umami, Matomo provider integrations                                                   | Stable |
| **Backend**    | Optional server (`minimaldoc-server`): cookie-free analytics, feedback widget, newsletter, blog CMS (RBAC, scheduling, RSS, comments), private docs, image uploads (local + S3), OAuth/JWT auth, admin portal, custom events, PDF export | WIP    |
| **Forum**      | Via backend server: categories, topics, replies, likes, bookmarks, reputation, moderation, bans, notifications, search | WIP    |

## Quick Start

```bash
# Install
git clone https://github.com/studiowebux/minimaldoc
cd minimaldoc && go build -o minimaldoc cmd/minimaldoc/main.go

# Create and build site
minimaldoc init my-docs && cd my-docs && minimaldoc build
```

## Documentation

The `docs/` directory is a **demo site** — it exists to showcase what MinimalDoc can generate and to test the generator itself. It is not the authoritative documentation for this project.

| Topic           | Link                                                       |
| --------------- | ---------------------------------------------------------- |
| Getting Started | [docs/getting-started/](docs/getting-started/)             |
| Configuration   | [docs/core-concepts/](docs/core-concepts/)                 |
| Features        | [docs/features/](docs/features/)                           |
| OpenAPI         | [docs/openapi/](docs/openapi/)                             |
| Status Page     | [docs/status-page/](docs/status-page/)                     |
| Backend Server  | [docs/features/12-backend.md](docs/features/12-backend.md) |
| Forum           | [docs/features/13-forum.md](docs/features/13-forum.md)     |
| Deployment      | [deploy/](deploy/)                                         |

## License

See [LICENSE](LICENSE) file.

## Credits

Created by [Studio Webux](https://studiowebux.com)
