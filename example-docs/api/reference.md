---
title: API Reference
description: Command-line interface and API documentation
tags:
  - api
  - cli
  - reference
---

# API Reference

Complete reference for MinimalDoc's command-line interface.

## Command-Line Interface

### build

Build your documentation site from Markdown files.

```bash
minimaldoc build [docs-directory] [flags]
```

#### Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `docs-directory` | Source directory containing .md files | `.` (current dir) |

#### Flags

##### Basic Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--title` | | string | "Documentation" | Site title |
| `--description` | | string | "" | Site description |
| `--base-url` | | string | "" | Base URL for sitemap |
| `--output` | `-o` | string | "public" | Output directory |
| `--theme` | `-t` | string | "default" | Theme name |

##### Advanced Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--llms` | `-l` | bool | true | Generate llms.txt |
| `--clean-urls` | | bool | false | Use clean URLs |

#### Examples

**Basic build:**

```bash
minimaldoc build
```

**Build with custom title:**

```bash
minimaldoc build --title="My Documentation"
```

**Full configuration:**

```bash
minimaldoc build docs/ \
  --title="Product Docs" \
  --description="Complete product documentation" \
  --base-url="https://docs.example.com" \
  --output="dist" \
  --llms=true
```

**Production build:**

```bash
minimaldoc build . \
  --title="MyApp Docs" \
  --base-url="https://docs.myapp.com" \
  --clean-urls
```

#### Output

The build command generates:

```
public/
├── index.html          # Homepage
├── *.html              # Page files
├── css/
│   └── main.css        # Styles
├── js/
│   ├── theme-toggle.js
│   ├── scrollspy.js
│   ├── mobile-menu.js
│   ├── anchor-links.js
│   ├── external-links.js
│   └── search.js
├── search-index.json   # Search index
├── sitemap.xml         # SEO sitemap
└── llms.txt            # LLM-friendly docs
```

### init

Initialize a new documentation project (coming soon).

```bash
minimaldoc init [directory]
```

This will create:
- Sample `index.md`
- Basic directory structure
- Example TOC.md
- .gitignore file

### version

Display MinimalDoc version.

```bash
minimaldoc --version
```

### help

Show help information.

```bash
minimaldoc --help
minimaldoc build --help
```

## Frontmatter API

Configure individual pages using YAML frontmatter.

### Basic Fields

```yaml
---
title: Page Title
description: Page description for SEO and previews
tags:
  - tag1
  - tag2
---
```

### Navigation Fields

```yaml
---
menu_title: Short Nav Title
menu_order: 5
hidden: false
---
```

### SEO Fields

```yaml
---
seo:
  title: Custom SEO Title
  description: Custom meta description
  keywords:
    - keyword1
    - keyword2
  canonical: https://example.com/canonical-url
  noindex: false
  nofollow: false
  image: https://example.com/og-image.jpg
author: Author Name
---
```

### Complete Example

```yaml
---
title: Complete Example
description: This shows all available frontmatter options
tags:
  - example
  - reference
menu_title: Example
menu_order: 10
hidden: false
author: John Doe
seo:
  title: Complete Frontmatter Example | MinimalDoc
  description: Learn all frontmatter options in MinimalDoc
  keywords:
    - frontmatter
    - yaml
    - configuration
  canonical: https://docs.example.com/api/reference.html
  image: https://docs.example.com/images/og-image.jpg
  noindex: false
  nofollow: false
---
```

## TOC.md Format

Define custom navigation structure.

### Basic Format

```markdown
# Table of Contents

- [Page Title](path/to/file.md)
- Section Header
  - [Nested Page](path/to/nested.md)
```

### Features

- **Links**: `[Title](path.md)` creates navigation links
- **Sections**: Text without links creates section headers
- **Nesting**: Indent with 2 spaces per level
- **Paths**: Relative to docs root

### Example

```markdown
# Navigation Structure

- [Home](index.md)
- Getting Started
  - [Installation](getting-started/installation.md)
  - [Quick Start](getting-started/quick-start.md)
  - [Configuration](getting-started/configuration.md)
- Features
  - [Overview](features/overview.md)
  - [Markdown](features/markdown.md)
  - Advanced
    - [Theming](features/theming.md)
    - [Plugins](features/plugins.md)
- API
  - [Reference](api/reference.md)
```

## Generated Files

### search-index.json

JSON file containing the search index:

```json
{
  "entries": [
    {
      "title": "Page Title",
      "description": "Page description",
      "url": "/page.html",
      "content": "First 500 chars of content...",
      "tags": ["tag1", "tag2"]
    }
  ]
}
```

### sitemap.xml

XML sitemap for search engines:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/page.html</loc>
    <lastmod>2025-01-15T10:30:00Z</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.8</priority>
  </url>
</urlset>
```

### llms.txt

LLM-friendly documentation:

```
# Site Title - Complete Documentation

---

This file contains all documentation in a single file.

---

<!-- PAGE: slug -->

# Page Title

Content here...

---

<!-- PAGE: next-slug -->

# Next Page

More content...
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Build failed |
| 2 | Invalid arguments |

## Environment Variables

Currently, MinimalDoc doesn't use environment variables. All configuration is via flags.

## Configuration Files

MinimalDoc doesn't require configuration files. Everything is configured via:
- Command-line flags
- Frontmatter (per-page)
- TOC.md (navigation)

## Best Practices

:::success API Usage Tips

1. **Use --base-url** - Always set for production
2. **Consistent flags** - Use same flags for each build
3. **Script builds** - Create build scripts for repeatability
4. **Version output** - Keep output dir in .gitignore
5. **Test locally** - Build and test before deploying
:::

## Examples

### Development Build

```bash
minimaldoc build docs/ \
  --title="Dev Docs" \
  --output="dev-build"
```

### Production Build

```bash
minimaldoc build docs/ \
  --title="Product Documentation" \
  --description="Official product documentation" \
  --base-url="https://docs.product.com" \
  --output="public" \
  --clean-urls
```

### CI/CD Build

```bash
#!/bin/bash
set -e

echo "Building documentation..."

minimaldoc build docs/ \
  --title="${SITE_TITLE}" \
  --base-url="${SITE_URL}" \
  --output="public"

echo "Build complete!"
```

## Error Messages

### Common Errors

**"docs directory does not exist"**
- The specified directory wasn't found
- Check the path and try again

**"failed to parse frontmatter"**
- YAML frontmatter has syntax errors
- Check for proper indentation and syntax

**"failed to generate HTML"**
- Template rendering failed
- Check for invalid markdown or HTML

## Next Steps

- [Read the Quick Start](../getting-started/quick-start.html)
- [Learn about deployment](../guides/deployment.html)
- [Check the FAQ](../guides/faq.html)
