---
title: Features Overview
description: Explore all MinimalDoc features
tags:
  - features
  - overview
---

# Features Overview

MinimalDoc comes packed with features to make your documentation beautiful and functional.

## Core Features

### Markdown Support

Write your documentation in standard Markdown with GitHub Flavored Markdown extensions:

- Tables
- Task lists
- Strikethrough
- Automatic URL linking

[Learn more about Markdown syntax →](markdown.html)

### Syntax Highlighting

Beautiful syntax highlighting for code blocks using Chroma:

```javascript
const greet = (name) => {
  console.log(`Hello, ${name}!`);
};

greet('World');
```

Supports 100+ languages including JavaScript, Python, Go, Rust, and more.

### Custom Admonitions

Create eye-catching callout blocks:

:::info
This is an info block for general information.
:::

:::warning
This is a warning block for important notices.
:::

:::danger
This is a danger block for critical warnings.
:::

:::success
This is a success block for positive messages.
:::

:::question
This is a question block for FAQs.
:::

:::note
This is a note block for side notes.
:::

[Learn more about admonitions →](markdown.html#admonitions)

## Navigation Features

### Automatic Navigation

Navigation is automatically generated from your file structure. Folders become sections, files become links.

### Custom TOC

Define custom navigation with a `TOC.md` file for complete control over structure and order.

### Table of Contents

Automatic TOC generation for each page based on headings. Includes smooth scrolling and active section highlighting.

[Learn more about navigation →](navigation.html)

## Search

### Full-Text Search

Client-side search across all your documentation. No server required!

### Keyboard Shortcuts

- `Cmd+K` (Mac) or `Ctrl+K` (Windows/Linux) - Open search
- `↑` / `↓` - Navigate results
- `Enter` - Go to selected page
- `Esc` - Close search

[Learn more about search →](search.html)

## Design Features

### Responsive Design

Your documentation looks great on all devices:
- Desktop (3-column layout)
- Tablet (2-column layout)
- Mobile (single column with hamburger menu)

### Dark Mode

Built-in dark mode with smooth transitions. Toggle with the theme button in the sidebar.

### Clean Typography

Carefully chosen fonts and spacing for maximum readability.

## SEO & Metadata

### Automatic Sitemap

XML sitemap automatically generated for search engines.

### Meta Tags

Comprehensive meta tags including:
- Open Graph (Facebook)
- Twitter Cards
- SEO-optimized titles and descriptions

### Structured URLs

Clean, hierarchical URLs that match your content structure.

## LLM Integration

### LLM-Friendly Output

Automatically generates `llms.txt` - a single file containing all your documentation optimized for AI assistants and chatbots.

Perfect for:
- ChatGPT
- Claude
- GitHub Copilot
- Other AI tools

## Performance

### Fast Builds

Generate hundreds of pages in seconds. No slow JavaScript bundlers.

### Static Output

Pure HTML, CSS, and JavaScript. Deploy anywhere:
- GitHub Pages
- Netlify
- Vercel
- Any static host

### Minimal JavaScript

Only ~10KB of JavaScript for interactive features. No heavy frameworks.

## Developer Experience

### Zero Configuration

Works out of the box with sensible defaults.

### File Watching

(Coming soon) Automatic rebuilds when files change.

### Live Preview

(Coming soon) Built-in development server with live reload.

## Next Steps

Explore specific features in detail:

- [Markdown Syntax Guide](markdown.html)
- [Navigation System](navigation.html)
- [Search Functionality](search.html)
- [Theming & Customization](theming.html)
