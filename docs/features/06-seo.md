---
title: SEO
description: Search engine optimization features including sitemap, meta tags, and Open Graph
tags:
  - features
  - seo
---

# SEO

MinimalDoc generates SEO-friendly output automatically.

## Sitemap

A `sitemap.xml` is generated during build:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/docs/</loc>
    <lastmod>2025-01-15</lastmod>
  </url>
  <url>
    <loc>https://example.com/docs/getting-started/</loc>
    <lastmod>2025-01-14</lastmod>
  </url>
</urlset>
```

Submit to search engines:
- Google Search Console
- Bing Webmaster Tools

## Meta Tags

### Title

Generated from frontmatter:

```yaml
---
title: Installation Guide
---
```

Produces:

```html
<title>Installation Guide - My Docs</title>
```

Override with SEO title:

```yaml
---
title: Installation Guide
seo:
  title: How to Install MyProject - Complete Guide
---
```

### Description

```yaml
---
description: Step-by-step installation instructions
---
```

Produces:

```html
<meta name="description" content="Step-by-step installation instructions">
```

### Keywords

```yaml
---
seo:
  keywords:
    - installation
    - setup
    - getting started
---
```

Produces:

```html
<meta name="keywords" content="installation, setup, getting started">
```

## Open Graph

For social media sharing:

```yaml
---
title: MinimalDoc
description: Static site generator for documentation
seo:
  image: /images/social-card.png
---
```

Produces:

```html
<meta property="og:title" content="MinimalDoc">
<meta property="og:description" content="Static site generator for documentation">
<meta property="og:image" content="https://example.com/images/social-card.png">
<meta property="og:type" content="article">
<meta property="og:url" content="https://example.com/docs/">
```

### Twitter Cards

Also generates Twitter Card meta:

```html
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="MinimalDoc">
<meta name="twitter:description" content="Static site generator for documentation">
<meta name="twitter:image" content="https://example.com/images/social-card.png">
```

## Canonical URLs

Specify canonical URL:

```yaml
---
seo:
  canonical: https://example.com/docs/installation
---
```

Produces:

```html
<link rel="canonical" href="https://example.com/docs/installation">
```

## Robots Control

### Noindex

Hide page from search engines:

```yaml
---
seo:
  noindex: true
---
```

Produces:

```html
<meta name="robots" content="noindex">
```

### Nofollow

Prevent following links:

```yaml
---
seo:
  nofollow: true
---
```

Both:

```yaml
---
seo:
  noindex: true
  nofollow: true
---
```

Produces:

```html
<meta name="robots" content="noindex, nofollow">
```

## Base URL

Configure in `config.yaml`:

```yaml
base_url: https://docs.example.com
```

Used for:
- Sitemap URLs
- Canonical URLs
- Open Graph URLs
- Absolute links

## Clean URLs

Enable for SEO-friendly URLs:

```yaml
clean_urls: true
```

| Setting | URL |
|---------|-----|
| `false` | `/docs/installation.html` |
| `true` | `/docs/installation/` |

## Best Practices

### Page Titles

```yaml
# Good - Specific, unique
title: Configure OAuth 2.0 Authentication

# Avoid - Generic
title: Configuration
```

### Descriptions

```yaml
# Good - Compelling, keyword-rich
description: Learn how to set up OAuth 2.0 authentication with Google, GitHub, and custom providers in MinimalDoc.

# Avoid - Too short or generic
description: Auth docs
```

### Images

- Use descriptive alt text
- Optimize file sizes
- Use WebP format when possible
- Provide social card images (1200x630px)

### URLs

- Use lowercase
- Use hyphens, not underscores
- Keep short and descriptive
- Avoid special characters

### Content

- Use heading hierarchy (H1 → H2 → H3)
- Include target keywords naturally
- Write descriptive link text
- Keep content up to date
