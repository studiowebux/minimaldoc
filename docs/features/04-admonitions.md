---
title: Admonitions
description: Callout blocks for important information, warnings, and tips
tags:
  - features
  - admonitions
  - markdown
---

# Admonitions

Admonitions are callout blocks for highlighting important content.

## Syntax

```markdown
:::type
Content here
:::
```

Or with custom title:

```markdown
:::type Custom Title
Content here
:::
```

## Available Types

### Info

```markdown
:::info
This is an informational note.
:::
```

:::info
This is an informational note.
:::

### Warning

```markdown
:::warning
This requires attention.
:::
```

:::warning
This requires attention.
:::

### Danger

```markdown
:::danger
This action is irreversible.
:::
```

:::danger
This action is irreversible.
:::

### Success

```markdown
:::success
Operation completed successfully.
:::
```

:::success
Operation completed successfully.
:::

### Note

```markdown
:::note
A side note for additional context.
:::
```

:::note
A side note for additional context.
:::

### Question

```markdown
:::question
Frequently asked question format.
:::
```

:::question
Frequently asked question format.
:::

## Custom Titles

Override the default title:

```markdown
:::info Prerequisites
Make sure you have Go 1.24+ installed.
:::
```

:::info Prerequisites
Make sure you have Go 1.24+ installed.
:::

```markdown
:::warning Breaking Change
This version changes the API format.
:::
```

:::warning Breaking Change
This version changes the API format.
:::

## Markdown Inside

Full Markdown is supported inside admonitions:

```markdown
:::info
**Bold**, *italic*, and `code` work here.

- List item 1
- List item 2

```bash
echo "Code blocks too"
```
:::
```

:::info
**Bold**, *italic*, and `code` work here.

- List item 1
- List item 2

```bash
echo "Code blocks too"
```
:::

## Color Reference

| Type | Light BG | Dark BG | Use For |
|------|----------|---------|---------|
| `info` | Blue | Dark blue | General information |
| `warning` | Yellow | Dark yellow | Cautions, deprecations |
| `danger` | Red | Dark red | Critical warnings |
| `success` | Green | Dark green | Confirmations |
| `note` | Gray | Dark gray | Side notes |
| `question` | Purple | Dark purple | FAQs |

## Usage Guidelines

### When to Use

**Info** - Background information, tips, explanations
```markdown
:::info
MinimalDoc generates static HTML files.
:::
```

**Warning** - Potential issues, deprecations, important notices
```markdown
:::warning
This feature is deprecated in v2.0.
:::
```

**Danger** - Breaking changes, data loss, security issues
```markdown
:::danger
This will delete all your data.
:::
```

**Success** - Completed steps, confirmations
```markdown
:::success
Your site is now deployed.
:::
```

**Note** - Tangential information, asides
```markdown
:::note
This is optional for basic usage.
:::
```

**Question** - FAQs, common questions
```markdown
:::question How do I deploy?
See the [deployment guide](../guides/deployment.md).
:::
```

### Avoid Overuse

Too many admonitions reduce their impact:

```markdown
<!-- Avoid this -->
:::info
Some info
:::

:::warning
A warning
:::

:::note
A note
:::

Regular paragraph.

:::info
More info
:::
```

Use sparingly for truly important callouts.

## Styling

Admonitions use CSS variables for colors:

```css
.admonition {
  padding: 1rem;
  border-left: 4px solid var(--admonition-border);
  background: var(--admonition-bg);
  border-radius: 4px;
  margin: 1rem 0;
}

.admonition-title {
  font-weight: 600;
  margin-bottom: 0.5rem;
}
```

See [Theming](03-theming.md) for customization.
