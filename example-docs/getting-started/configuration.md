---
title: Configuration
description: Configure MinimalDoc for your needs
tags:
  - configuration
  - settings
---

# Configuration

MinimalDoc works with sensible defaults but can be customized through command-line flags.

## Command-Line Options

### Basic Options

```bash
minimaldoc build [docs-directory] [flags]
```

**Available Flags:**

- `--title` - Site title (default: "Documentation")
- `--description` - Site description
- `--base-url` - Base URL for sitemap and links
- `--output, -o` - Output directory (default: "public")
- `--theme, -t` - Theme name (default: "default")

### Advanced Options

- `--llms, -l` - Generate LLM-friendly files (default: true)
- `--clean-urls` - Use clean URLs (/page/ instead of /page.html)

## Example Commands

### Basic Build

```bash
minimaldoc build docs/
```

### Full Configuration

```bash
minimaldoc build docs/ \
  --title="My Project Docs" \
  --description="Complete documentation for My Project" \
  --base-url="https://docs.myproject.com" \
  --output="dist" \
  --llms=true
```

### Production Build

```bash
minimaldoc build . \
  --title="ProductName Docs" \
  --base-url="https://docs.product.com" \
  --clean-urls
```

## Frontmatter Configuration

Configure individual pages using YAML frontmatter:

```yaml
---
title: Page Title
description: Page description
tags:
  - tag1
  - tag2
hidden: false
menu_title: Short Title
menu_order: 5
seo:
  title: SEO Title
  description: SEO Description
  keywords:
    - keyword1
    - keyword2
---
```

### Frontmatter Options

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Page title |
| `description` | string | Page description |
| `tags` | array | Page tags |
| `hidden` | boolean | Hide from navigation |
| `menu_title` | string | Custom nav title |
| `menu_order` | number | Sort order |
| `seo.title` | string | SEO title tag |
| `seo.description` | string | Meta description |
| `seo.keywords` | array | Meta keywords |

## File Organization

MinimalDoc discovers pages automatically from your file structure:

```
docs/
├── index.md              # Homepage
├── getting-started/
│   ├── installation.md
│   └── quick-start.md
├── guides/
│   ├── user-guide.md
│   └── admin-guide.md
└── api/
    ├── overview.md
    └── reference.md
```

:::info
**Navigation Order**: By default, pages are sorted alphabetically. Use numbered prefixes (01-, 02-) or the `menu_order` frontmatter field to control order.
:::

## Custom Navigation

Create a `TOC.md` file to define custom navigation structure:

```markdown
# Table of Contents

- [Home](index.md)
- Getting Started
  - [Installation](getting-started/installation.md)
  - [Quick Start](getting-started/quick-start.md)
- Guides
  - [User Guide](guides/user-guide.md)
```

:::note
When `TOC.md` exists, it takes precedence over automatic navigation generation.
:::

## Next Steps

- [Learn about features](../features/overview.html)
- [Understand navigation](../features/navigation.html)
- [Explore Markdown syntax](../features/markdown.html)
