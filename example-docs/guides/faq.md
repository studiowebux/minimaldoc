---
title: Frequently Asked Questions
description: Common questions and answers about MinimalDoc
tags:
  - faq
  - help
  - troubleshooting
---

# Frequently Asked Questions

Common questions and answers about MinimalDoc.

## General Questions

### What is MinimalDoc?

MinimalDoc is a static site generator specifically designed for documentation. It converts Markdown files into beautiful, responsive documentation websites with built-in search, navigation, and dark mode.

### Why choose MinimalDoc over other generators?

MinimalDoc is:
- **Simpler** - No complex configuration
- **Faster** - Lightning-fast builds
- **Smaller** - Minimal dependencies
- **Focused** - Built specifically for docs

### Is MinimalDoc free?

Yes! MinimalDoc is open source and completely free to use.

## Installation & Setup

### How do I install MinimalDoc?

```bash
go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
```

See the [Installation Guide](../getting-started/installation.html) for details.

### What are the system requirements?

- Go 1.20 or later
- Any operating system (Windows, macOS, Linux)

### Can I use MinimalDoc without Go installed?

Currently, no. We plan to release standalone binaries in the future.

## Usage Questions

### How do I create a new site?

1. Create a directory for your docs
2. Add an `index.md` file
3. Run `minimaldoc build .`
4. Open `public/index.html`

See the [Quick Start Guide](../getting-started/quick-start.html).

### How do I add new pages?

Simply create new `.md` files in your docs directory and rebuild:

```bash
echo "# New Page" > new-page.md
minimaldoc build .
```

### How do I control page order?

Use numbered prefixes:

```
01-introduction.md
02-getting-started.md
03-advanced.md
```

Or use `TOC.md` for complete control. See [Navigation](../features/navigation.html).

### Can I use subdirectories?

Yes! Subdirectories become sections in your navigation:

```
docs/
├── getting-started/
│   ├── installation.md
│   └── quick-start.md
└── guides/
    └── user-guide.md
```

## Features

### Does MinimalDoc support search?

Yes! Built-in client-side search with `Cmd+K` / `Ctrl+K` shortcut. See [Search](../features/search.html).

### Can I customize the theme?

Currently, customization is limited. Custom themes are planned for a future release. See [Theming](../features/theming.html).

### Does it support syntax highlighting?

Yes! Supports 100+ languages with beautiful syntax highlighting. See [Markdown Guide](../features/markdown.html).

### What Markdown features are supported?

- GitHub Flavored Markdown
- Tables
- Task lists
- Syntax highlighting
- Custom admonitions
- And more!

See the [Markdown Syntax Guide](../features/markdown.html).

### Does it work on mobile?

Yes! Fully responsive with a mobile-friendly hamburger menu.

### Is there a dark mode?

Yes! Built-in dark mode with automatic theme switching.

## Deployment

### Where can I host my site?

Anywhere that hosts static sites:
- GitHub Pages
- Netlify
- Vercel
- Cloudflare Pages
- Your own server

See the [Deployment Guide](deployment.html).

### How do I deploy to GitHub Pages?

```bash
# Build
minimaldoc build docs/ --base-url="https://username.github.io/repo"

# Deploy
cd public
git init
git add -A
git commit -m "Deploy"
git push -f git@github.com:username/repo.git main:gh-pages
```

See the [Deployment Guide](deployment.html) for automated options.

### Do I need a server?

No! MinimalDoc generates static HTML files that work on any web server or static hosting platform.

### How do I use a custom domain?

Configure your DNS to point to your hosting platform. Most platforms provide automatic HTTPS. See [Deployment](deployment.html#custom-domain).

## Troubleshooting

### Build fails with "command not found"

Make sure MinimalDoc is installed and in your PATH:

```bash
which minimaldoc
```

If not found, reinstall or add to PATH.

### Pages show 404 errors

Check:
1. Files exist in `public/` directory
2. `--base-url` matches your deployment URL
3. Server configuration for .html extension
4. File paths are correct

### Search doesn't work

Check these common issues:

1. Rebuild to regenerate search index
2. Check `public/search-index.json` exists
3. Clear browser cache
4. Check browser console for errors

### Navigation is in wrong order

Use numbered prefixes or create a `TOC.md` file to control order. See [Navigation](../features/navigation.html).

### Dark mode doesn't work

Check:
1. JavaScript is enabled
2. Browser supports localStorage
3. No JavaScript errors in console

### Links are broken

Use relative paths and include `.html` extension:

```markdown
[Link](../guides/user-guide.html)
```

Not:
```markdown
[Link](../guides/user-guide)
```

## Performance

### How fast is MinimalDoc?

Very fast! Builds hundreds of pages in seconds. The output is pure HTML with minimal JavaScript (~10KB).

### Does it work offline?

Yes! Once built, the site is static HTML and works completely offline.

### How big are the generated files?

Very small:
- ~3KB per HTML page (compressed)
- ~5KB CSS (total)
- ~10KB JavaScript (total)
- Search index varies by content

## Development

### Can I watch for file changes?

Not yet. File watching and live reload are planned features.

Current workaround:

```bash
# Install watchexec
brew install watchexec  # macOS
# or: apt install watchexec  # Linux

# Watch for changes
watchexec -e md minimaldoc build docs/
```

### Can I extend MinimalDoc?

Plugin system is planned for future releases.

### How do I contribute?

Visit the [GitHub repository](https://github.com/studiowebux/minimaldoc) to:
- Report issues
- Suggest features
- Submit pull requests

## Comparison

### MinimalDoc vs MkDocs?

MinimalDoc:
- Simpler configuration
- Faster builds
- No Python required
- Smaller output

MkDocs:
- More themes
- More plugins
- Larger ecosystem

### MinimalDoc vs Hugo?

MinimalDoc:
- Focused on docs only
- Simpler to use
- Built-in search

Hugo:
- More versatile (blogs, etc.)
- More themes
- More complex

### MinimalDoc vs Docusaurus?

MinimalDoc:
- Much simpler
- No Node.js required
- Faster builds
- Smaller output

Docusaurus:
- React-based
- More features
- Larger ecosystem

## Future Plans

### What features are planned?

- [ ] File watching
- [ ] Live preview server
- [ ] Custom themes
- [ ] Plugin system
- [ ] Standalone binaries
- [ ] More customization options

### How can I request a feature?

Open an issue on [GitHub](https://github.com/studiowebux/minimaldoc/issues) with your suggestion.

## Getting Help

### Where can I get help?

1. Read this FAQ
2. Check the [documentation](../index.html)
3. Search [GitHub issues](https://github.com/studiowebux/minimaldoc/issues)
4. Open a new issue

### How do I report a bug?

Open an issue on [GitHub](https://github.com/studiowebux/minimaldoc/issues) with:
- Description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version)

## Still Have Questions?

Can't find your answer?

- [Read the documentation](../index.html)
- [Check GitHub issues](https://github.com/studiowebux/minimaldoc/issues)
- [Ask a question](https://github.com/studiowebux/minimaldoc/discussions)
