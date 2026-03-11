---
title: Analytics
description: Plugin-based analytics system with support for multiple providers
menu_order: 11
---

# Analytics

MinimalDoc includes a plugin-based analytics system that supports multiple providers simultaneously. Each provider is independently configurable, and custom/self-hosted solutions can be integrated using the `custom` provider type.

## Configuration

Enable analytics in your `config.yaml`:

```yaml
analytics:
  enabled: true
  providers:
    - type: plausible
      enabled: true
      config:
        domain: "docs.example.com"
```

## Supported Providers

| Provider | Type Key | Required Config | Script Location |
|----------|----------|-----------------|-----------------|
| MinimalDoc Backend | `minimaldoc` | `endpoint`, `site_id` | body |
| Google Analytics 4 | `ga4` | `measurement_id` | head |
| Plausible | `plausible` | `domain` | head |
| Umami | `umami` | `website_id`, `src` | head |
| Matomo | `matomo` | `url`, `site_id` | head |
| Fathom | `fathom` | `site_id` | head |
| Simple Analytics | `simple` | (none) | head |
| Custom | `custom` | `src` | configurable |

## Provider Examples

### MinimalDoc Backend

Self-hosted, cookie-free analytics with the optional MinimalDoc backend server. See [Backend](/features/backend.html) for full documentation.

```yaml
analytics:
  enabled: true
  providers:
    - type: minimaldoc
      enabled: true
      config:
        endpoint: "https://your-minimaldoc-server.com"
        site_id: "your-site-id-from-bootstrap"
        features: "analytics,feedback,newsletter"  # optional, default: analytics
        debug: "false"  # optional, enables console logging
```

Features available:
- `analytics` - Cookie-free page view tracking
- `feedback` - Page rating widget (requires `data-minimaldoc-feedback` element)
- `newsletter` - Newsletter signup forms (requires `data-minimaldoc-newsletter` form)

### Google Analytics 4

```yaml
analytics:
  enabled: true
  providers:
    - type: ga4
      enabled: true
      config:
        measurement_id: "G-XXXXXXXXXX"
```

### Plausible

Standard Plausible.io:

```yaml
analytics:
  enabled: true
  providers:
    - type: plausible
      enabled: true
      config:
        domain: "docs.example.com"
```

Self-hosted Plausible:

```yaml
analytics:
  enabled: true
  providers:
    - type: plausible
      enabled: true
      config:
        domain: "docs.example.com"
        src: "https://analytics.your-domain.com/js/script.js"
```

### Umami

```yaml
analytics:
  enabled: true
  providers:
    - type: umami
      enabled: true
      config:
        website_id: "your-website-uuid"
        src: "https://analytics.example.com/script.js"
```

### Matomo

```yaml
analytics:
  enabled: true
  providers:
    - type: matomo
      enabled: true
      config:
        url: "https://matomo.example.com"
        site_id: "1"
```

### Fathom

```yaml
analytics:
  enabled: true
  providers:
    - type: fathom
      enabled: true
      config:
        site_id: "ABCDEFGH"
```

Self-hosted:

```yaml
analytics:
  enabled: true
  providers:
    - type: fathom
      enabled: true
      config:
        site_id: "ABCDEFGH"
        src: "https://analytics.your-domain.com/tracker.js"
```

### Simple Analytics

```yaml
analytics:
  enabled: true
  providers:
    - type: simple
      enabled: true
```

### Custom Provider

For any analytics service not directly supported, use the `custom` type:

```yaml
analytics:
  enabled: true
  providers:
    - type: custom
      enabled: true
      config:
        src: "https://my-analytics.example.com/tracker.js"
        defer: "true"
        async: "false"
        location: "head"  # or "body"
        attrs:
          data-site-id: "my-site"
          data-cookieless: "true"
          data-no-pii: "true"
```

The `attrs` map allows passing arbitrary data attributes to the script tag.

## Multiple Providers

Run multiple analytics providers simultaneously:

```yaml
analytics:
  enabled: true
  providers:
    - type: plausible
      enabled: true
      config:
        domain: "docs.example.com"

    - type: ga4
      enabled: false  # Disabled but kept in config
      config:
        measurement_id: "G-XXXXXXXXXX"

    - type: custom
      enabled: true
      config:
        src: "https://internal-analytics.example.com/t.js"
        defer: "true"
        attrs:
          data-site: "docs"
```

Only providers with `enabled: true` will have their scripts injected.

## Script Injection

Analytics scripts are injected in two locations:

1. **Head scripts** (before `</head>`) - Most providers inject here
2. **Body scripts** (before `</body>`) - Custom providers can use `location: body`

Scripts are injected on all page types:
- Documentation pages
- Landing pages
- Status pages
- Changelog
- FAQ
- Contact
- Legal pages
- Knowledge Base
- Portfolio
- OpenAPI documentation

## Privacy Considerations

When choosing an analytics provider, consider:

- **GDPR compliance** - Some providers require cookie consent banners
- **Cookieless tracking** - Plausible, Fathom, and Simple Analytics are cookieless by default
- **Self-hosting** - Plausible, Umami, and Matomo can be self-hosted for full data control
- **Data retention** - Check provider policies on how long data is stored
