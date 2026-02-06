---
title: Deployment
description: Deploy MinimalDoc sites to any static hosting
tags:
  - guides
  - deployment
---

# Deployment

MinimalDoc generates static HTML that deploys anywhere.

## Build for Production

```bash
minimaldoc build ./docs \
  --base-url "https://docs.example.com" \
  --output dist
```

Output in `dist/` is ready to deploy.

## Base URL

Set `base_url` for correct absolute URLs:

```yaml
# config.yaml
base_url: https://docs.example.com
```

Or via CLI:

```bash
minimaldoc build --base-url "https://docs.example.com"
```

Affects:
- Sitemap URLs
- Canonical links
- Open Graph URLs
- RSS feed links

### Subdirectory Deployment

For `/docs/` subdirectory:

```yaml
base_url: https://example.com/docs
```

## Clean URLs

Enable for `/page/` instead of `/page.html`:

```yaml
clean_urls: true
```

## Hosting

Can be deployed anywhere that supports Static files (HTML, JS, CSS, JSON)
