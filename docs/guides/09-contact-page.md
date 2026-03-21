---
title: Contact Page
description: Add a contact page with info cards
tags:
  - guides
  - contact
---

# Contact Page

Add a contact page with info cards displaying email, location, and other details.

## Configuration

```yaml
contact:
  enabled: true
  title: Contact
  description: Get in touch
  path: contact
  email: hello@example.com
  info:
    - icon: mail
      text: hello@example.com
    - icon: location
      text: Montreal, Canada
    - icon: website
      text: https://example.com
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable contact page |
| `title` | string | `"Contact"` | Page title |
| `description` | string | `""` | Page description |
| `path` | string | `"contact"` | Output path |
| `email` | string | `""` | Contact email address |
| `info` | array | `[]` | Info cards with icon and text |

## Info Card Icons

Available icon identifiers:

| Icon | Description |
|------|-------------|
| `mail` | Email |
| `location` | Physical location |
| `phone` | Phone number |
| `website` | Website URL |
| `github` | GitHub profile |
| `twitter` | Twitter/X profile |
| `discord` | Discord server |

## Output

Generated at `/<path>/index.html`. A link appears in the site footer if the contact page is enabled.
