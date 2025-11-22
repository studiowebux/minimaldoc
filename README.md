# Minimal Doc

A modern, minimal static site generator for documentation written in Go.

## Features

- **Markdown-based** - Write documentation in Markdown with YAML frontmatter
- **Beautiful Theme** - Clean, minimal design with soft white/black color scheme
- **Dark Mode** - Automatic light/dark theme with localStorage persistence
- **Responsive** - Mobile-friendly, responsive design
- **Auto Navigation** - Automatic navigation tree from folder structure
- **Table of Contents** - Auto-generated TOC with active scrollspy
- **Search** - Client-side search with Cmd+K/Ctrl+K and fuzzy matching
- **Admonitions** - Callout blocks for info, warning, danger, success, note, question
- **Syntax Highlighting** - Code blocks with syntax highlighting
- **GFM Support** - Tables, strikethrough, task lists, and more
- **Sitemap** - Automatic sitemap.xml generation for SEO
- **LLM-Friendly** - Generates `llms.txt` for AI tools
- **Fast Builds** - Lightning-fast site generation
- **GitHub Pages Ready** - Deploy anywhere static sites are supported

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
- `--clean-urls` - Use clean URLs `/page/` instead of `/page.html`
- `--title` - Site title (default: `Documentation`)
- `--description` - Site description
- `--base-url` - Base URL for the site

**Example:**

```bash
minimaldoc build ./docs --output dist --title "My Project" --description "Project documentation"
```

#### `init`

Initialize a new documentation site with example files.

```bash
minimaldoc init [directory]
```

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

## Theme

### Colors

The default theme uses soft, eye-friendly colors:

**Light Mode:**

- Background: `#fafafa` (soft white)
- Text: `#1a1a1a` (soft black)

**Dark Mode:**

- Background: `#1a1a1a` (soft black)
- Text: `#f5f5f5` (soft white)

### Customization

Theme customization via CSS variables:

```css
:root[data-theme="light"] {
  --bg-primary: #fafafa;
  --text-primary: #1a1a1a;
  /* ... */
}
```

## LLM-Friendly Output

Minimal Doc generates `llms.txt` - a single file containing all documentation content with navigation structure.

Perfect for:

- AI assistants
- Documentation search
- Content indexing

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
