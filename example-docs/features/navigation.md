---
title: Navigation System
description: How navigation works in MinimalDoc
tags:
  - navigation
  - toc
---

# Navigation System

MinimalDoc provides flexible navigation options to organize your documentation.

## Automatic Navigation

By default, navigation is generated automatically from your file structure.

### File Structure Example

```
docs/
├── index.md
├── getting-started/
│   ├── installation.md
│   └── quick-start.md
├── guides/
│   ├── user-guide.md
│   └── admin-guide.md
└── api/
    └── reference.md
```

This creates:
- Home
- Getting Started
  - Installation
  - Quick Start
- Guides
  - User Guide
  - Admin Guide
- API
  - Reference

### Ordering

Control the order of pages and sections using numbered prefixes:

```
docs/
├── index.md
├── 01-getting-started/
│   ├── 01-installation.md
│   └── 02-quick-start.md
├── 02-guides/
│   ├── 01-user-guide.md
│   └── 02-admin-guide.md
└── 03-api/
    └── 01-reference.md
```

:::info
Number prefixes (01-, 02-, etc.) are automatically removed from navigation titles.
:::

## Custom Navigation with TOC.md

For complete control over navigation, create a `TOC.md` file in your docs root:

```markdown
# Table of Contents

- [Home](index.md)
- Getting Started
  - [Installation](getting-started/installation.md)
  - [Quick Start](getting-started/quick-start.md)
  - [Configuration](getting-started/configuration.md)
- Features
  - [Overview](features/overview.md)
  - [Markdown Guide](features/markdown.md)
  - [Navigation](features/navigation.md)
  - [Search](features/search.md)
- API Reference
  - [Reference](api/reference.md)
```

### TOC.md Features

- **Custom order** - Define exactly how items appear
- **Custom hierarchy** - Create any nesting structure
- **Section headers** - Add sections without links
- **Flexible paths** - Use relative paths to your markdown files

:::note
When `TOC.md` exists, it takes precedence over automatic navigation. Files not listed in `TOC.md` won't appear in navigation.
:::

## Page Table of Contents

Each page automatically generates a table of contents from its headings.

### How It Works

The right sidebar shows:
- All H2-H6 headings
- Nested hierarchy
- Active section highlighting
- Auto-scroll to follow reading position

### Controlling TOC Display

Use headings strategically:

```markdown
# Main Title (Not in TOC)
## Section 1 (In TOC)
### Subsection 1.1 (In TOC, nested)
#### Details (In TOC, further nested)
```

:::tip
Keep heading hierarchy clean for the best TOC experience. Don't skip levels (H2 → H4).
:::

## Navigation Frontmatter

Control navigation behavior per page:

```yaml
---
title: Page Title
menu_title: Short Nav Title
menu_order: 5
hidden: false
---
```

### Frontmatter Options

| Field | Description | Example |
|-------|-------------|---------|
| `menu_title` | Custom title in navigation | `menu_title: "Quick Start"` |
| `menu_order` | Sort order (overrides prefix) | `menu_order: 1` |
| `hidden` | Hide from navigation | `hidden: true` |

## Previous/Next Links

Pages automatically get previous/next links based on:
1. TOC.md order (if defined)
2. File structure order (automatic)
3. menu_order frontmatter
4. Alphabetical order

:::success
Sequential navigation helps readers progress through your documentation naturally.
:::

## Mobile Navigation

On mobile devices (< 768px width):
- Left sidebar becomes a slide-in menu
- Hamburger button toggles menu
- Backdrop closes menu on tap
- Smooth animations

## Keyboard Navigation

Use keyboard shortcuts for better navigation:

- `Cmd+K` / `Ctrl+K` - Open search
- `↑` / `↓` - Navigate search results
- `Enter` - Go to page
- `Esc` - Close modals

## Best Practices

:::info Navigation Guidelines

1. **Keep it shallow** - Max 3 levels of nesting
2. **Group logically** - Related content together
3. **Clear titles** - Descriptive, scannable names
4. **Consistent structure** - Similar patterns throughout
5. **Use TOC.md** - For complex sites, use custom navigation
:::

### Good Navigation Structure

```
- Getting Started (How to begin)
  - Installation
  - Quick Start
  - Configuration
- Guides (How to use)
  - User Guide
  - Admin Guide
  - Advanced Topics
- Reference (What exists)
  - API Reference
  - CLI Reference
  - Configuration Reference
```

### Poor Navigation Structure

```
- Docs (Too vague)
  - Page 1 (Not descriptive)
  - Stuff (Unclear)
  - More Docs
    - Things
      - Other Things
        - More (Too deep)
```

## Examples

### Minimal Site

```
docs/
├── index.md
├── guide.md
└── api.md
```

### Medium Site

```
docs/
├── index.md
├── getting-started/
│   └── installation.md
├── guides/
│   ├── user-guide.md
│   └── admin-guide.md
└── api/
    └── reference.md
```

### Large Site with TOC.md

```
docs/
├── TOC.md (Custom navigation)
├── index.md
├── getting-started/
│   ├── installation.md
│   ├── quick-start.md
│   └── configuration.md
├── tutorials/
│   ├── beginner.md
│   ├── intermediate.md
│   └── advanced.md
├── guides/
│   ├── user-guide.md
│   └── admin-guide.md
├── api/
│   ├── overview.md
│   └── reference.md
└── faq.md
```

## Next Steps

- [Learn about search](search.html)
- [Explore theming](theming.html)
- [See API reference](../api/reference.md)
