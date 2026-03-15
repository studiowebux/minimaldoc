---
title: Legal Pages
description: Privacy policy, terms of service, and other legal pages
tags:
  - guides
  - legal
---

# Legal Pages

Add privacy policy, terms of service, and other legal documents to your site.

## Configuration

```yaml
legal:
  enabled: true
  path: legal
  footer_group: Legal
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable legal pages |
| `path` | string | `"legal"` | Output path prefix |
| `footer_group` | string | `"Legal"` | Footer section title for legal page links |

## Legal Page Files

Create markdown files in `__legal__/` directory:

```
docs/
  __legal__/
    privacy-policy.md
    terms-of-service.md
    cookie-policy.md
```

Each file uses standard frontmatter:

```yaml
---
title: Privacy Policy
description: How we handle your data
order: 1
---

## Data Collection

We collect minimal data necessary for the service...
```

## Footer Integration

Legal pages are automatically added to a footer link group with the title specified by `footer_group`. Links appear in the order defined by the `order` frontmatter field.

## Output

Each file generates a page at `/<path>/<slug>.html`:

- `__legal__/privacy-policy.md` -> `/legal/privacy-policy.html`
- `__legal__/terms-of-service.md` -> `/legal/terms-of-service.html`
