---
title: Quick Start Guide
description: Get your documentation site running in 5 minutes
tags:
  - getting-started
  - tutorial
---

# Quick Start Guide

Get your documentation site up and running in just a few minutes!

## Step 1: Create Your Project

Create a new directory for your documentation:

```bash
mkdir my-docs
cd my-docs
```

## Step 2: Create Your First Page

Create an `index.md` file:

```markdown
---
title: My Documentation
description: Welcome to my docs
---

# Welcome

This is my documentation homepage!
```

## Step 3: Build Your Site

Run the build command:

```bash
minimaldoc build . --title="My Documentation"
```

This will:
- Parse all Markdown files
- Generate HTML pages
- Create navigation
- Generate search index
- Create sitemap

## Step 4: Preview Your Site

Open `public/index.html` in your browser to preview your site.

:::success
**Success!** Your documentation site is ready!
:::

## Step 5: Add More Pages

Create additional pages by adding more Markdown files:

```bash
mkdir guides
echo "# User Guide" > guides/user-guide.md
```

Rebuild to see your changes:

```bash
minimaldoc build .
```

## Build Options

Customize your build with these options:

| Option | Description | Example |
|--------|-------------|---------|
| `--title` | Site title | `--title="My Docs"` |
| `--description` | Site description | `--description="Documentation"` |
| `--base-url` | Base URL for sitemap | `--base-url="https://example.com"` |
| `--output` | Output directory | `--output="dist"` |
| `--theme` | Theme name | `--theme="default"` |

## Next Steps

- [Learn about configuration](configuration.html)
- [Explore features](../features/overview.html)
- [Understand navigation](../features/navigation.html)

:::warning
**Important**: Always rebuild after making changes to your Markdown files.
:::
