---
title: Navigation
description: Configure navigation with folder structure or custom TOC.md
tags:
  - core-concepts
  - navigation
---

# Navigation

MinimalDoc supports two navigation modes: automatic folder-based and custom TOC.md.

## Folder-Based Navigation

Default mode. Navigation is generated from your directory structure.

```
docs/
├── index.md                    # Home
├── 01-getting-started/         # Section
│   ├── 01-installation.md      # Page
│   └── 02-quick-start.md       # Page
├── 02-features/                # Section
│   ├── 01-search.md
│   └── 02-theming.md
└── 03-api/                     # Section
    └── reference.md
```

Generates:

```
- Home
- Getting Started
  - Installation
  - Quick Start
- Features
  - Search
  - Theming
- API
  - Reference
```

### Ordering Rules

1. `menu_order` frontmatter (ascending, lowest first)
2. Numeric prefix (01-, 02-, etc.)
3. Alphabetical by filename

Numeric prefixes are stripped from display:
- `01-getting-started/` → "Getting Started"
- `02-installation.md` → "Installation"

### Frontmatter Override

```yaml
---
title: Advanced Configuration
menu_title: Advanced       # Shorter nav title
menu_order: 50             # Custom sort position
---
```

## Custom Navigation with TOC.md

Create `TOC.md` in your docs root for full control.

```markdown
# My Project

- [Home](index.md)
- Getting Started
  - [Installation](getting-started/installation.md)
  - [Quick Start](getting-started/quick-start.md)
- Features
  - [Search](features/search.md)
  - [Theming](features/theming.md)
- [API Reference](api/reference.md)
- Resources
  - [GitHub](https://github.com/org/repo)
  - [Issues](https://github.com/org/repo/issues)
```

### Format Rules

| Syntax | Result |
|--------|--------|
| `- [Title](path.md)` | Page link |
| `- Section Name` | Section header (no link) |
| `- [External](https://...)` | External link |
| 2-space indent | Nesting level |

### Nesting

```markdown
- Level 1
  - Level 2
    - Level 3
      - Level 4
```

No practical nesting limit, but 3-4 levels recommended for usability.

### External Links

```markdown
- Resources
  - [GitHub](https://github.com/studiowebux/minimaldoc)
  - [Discord](https://discord.gg/invite)
```

External links display with an icon and open in new tabs.

### Mixed Navigation

```markdown
- [Home](index.md)
- Documentation
  - [Installation](docs/install.md)
  - [Configuration](docs/config.md)
- [API Docs](api/index.md)
- External
  - [GitHub](https://github.com/example)
```

## Page Table of Contents

Each page generates a right-side TOC from headings.

### Included Headings

- H2 through H6 are included
- H1 is typically the page title, excluded from TOC

### Scrollspy

Active heading is highlighted as you scroll. Smooth scrolling enabled.

### Anchor Links

Click any heading to copy its anchor link:

```
https://example.com/docs/page.html#section-name
```

## Hidden Pages

Hide pages from navigation but keep them accessible:

```yaml
---
title: Internal Notes
hidden: true
---
```

Page is built and accessible via direct URL but not shown in navigation.

## Navigation Behavior

### Active State

Current page and its parent sections are highlighted.

### Collapse

Sections collapse/expand on click. Current section auto-expands.

### Mobile

On mobile, navigation is hidden behind a hamburger menu.

### Keyboard

- `Tab` navigates between links
- `Enter` activates links
- `Escape` closes mobile menu

## Examples

### Documentation Site

```markdown
- [Overview](index.md)
- Getting Started
  - [Installation](getting-started/installation.md)
  - [Quick Start](getting-started/quick-start.md)
  - [Configuration](getting-started/configuration.md)
- Guides
  - [Deployment](guides/deployment.md)
  - [Customization](guides/customization.md)
- API Reference
  - [REST API](api/rest.md)
  - [SDK](api/sdk.md)
- [Changelog](changelog.md)
```

### Product Documentation

```markdown
- [Home](index.md)
- User Guide
  - [Dashboard](user/dashboard.md)
  - [Reports](user/reports.md)
  - [Settings](user/settings.md)
- Admin Guide
  - [Setup](admin/setup.md)
  - [Users](admin/users.md)
  - [Security](admin/security.md)
- Developer
  - [API](dev/api.md)
  - [Webhooks](dev/webhooks.md)
- Support
  - [FAQ](support/faq.md)
  - [Contact](https://support.example.com)
```
