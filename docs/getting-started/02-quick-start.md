---
title: Quick Start
description: Create your first MinimalDoc site in minutes
tags:
  - getting-started
  - quick-start
---

# Quick Start

## Initialize a New Site

```bash
minimaldoc init my-docs
cd my-docs
```

Generated structure:

```
my-docs/
├── config.yaml
├── index.md
├── getting-started/
│   ├── installation.md
│   └── quick-start.md
└── guides/
    └── deployment.md
```

For a full-featured site with status page, changelog, and OpenAPI:

```bash
minimaldoc init my-docs --full
```

This adds:

```
my-docs/
├── config.yaml               # All config here
├── TOC.md                    # Custom navigation
├── api/openapi.yaml          # OpenAPI spec
├── __status__/               # Status page
│   ├── components.yaml
│   └── incidents/
└── __changelog__/            # Changelog
    └── releases/
```

## Build the Site

```bash
minimaldoc build my-docs
```

Output written to `public/`:

```
public/
├── index.html
├── getting-started/
│   ├── installation.html
│   └── quick-start.html
├── guides/
│   └── deployment.html
├── search-index.json
├── sitemap.xml
└── llms.txt
```

## Preview

Open `public/index.html` in your browser, or serve locally:

```bash
# Python
python -m http.server -d public 8080

# Go
go run github.com/nicholaswilde/go-serve@latest -d public

# Node
npx serve public
```

Visit `http://localhost:8080`.

## Add a Page

Create `my-docs/guides/getting-help.md`:

```markdown
---
title: Getting Help
description: Where to find help and support
---

# Getting Help

Visit our [GitHub Issues](https://github.com/studiowebux/minimaldoc/issues).
```

Rebuild:

```bash
minimaldoc build
```

The page appears in navigation automatically.

## Edit Configuration

The `config.yaml` file is created automatically. Edit it to customize your site:

```yaml
title: My Project Docs
description: Documentation for my project
base_url: https://docs.myproject.com
theme: default

# Enable optional features
enable_llms: true
enable_search: true

# Uncomment to enable OpenAPI docs
# openapi:
#   enabled: true
#   spec_files:
#     - "api/openapi.yaml"

# Uncomment to enable status page
# status:
#   enabled: true
```

Rebuild to apply changes.

## Next Steps

- [Configuration](03-configuration.md) - Full configuration reference
- [Project Structure](../core-concepts/01-project-structure.md) - Understand the file layout
- [Navigation](../core-concepts/04-navigation.md) - Customize navigation
