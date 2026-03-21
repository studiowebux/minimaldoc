---
title: FAQ Pages
description: Frequently asked questions with categories and collapsible answers
tags:
  - guides
  - faq
---

# FAQ Pages

Add a FAQ page with categorized questions and collapsible answers.

## Configuration

Define questions inline in config.yaml:

```yaml
faq:
  enabled: true
  title: FAQ
  description: Frequently asked questions
  path: faq
  categories:
    - name: General
      items:
        - question: What is MinimalDoc?
          answer: A minimal static site generator for documentation.
        - question: Is it free?
          answer: Yes, GPL-3.0 licensed.
    - name: Technical
      items:
        - question: What languages are supported for syntax highlighting?
          answer: 100+ languages via Chroma.
```

## Markdown-Based FAQ

Alternatively, define FAQ items as markdown files in `__faq__/` directory:

```
docs/
  __faq__/
    01-what-is-minimaldoc.md
    02-installation-issues.md
```

Each file uses frontmatter:

```yaml
---
question: What is MinimalDoc?
category: General
tags:
  - general
  - getting-started
order: 1
---

MinimalDoc is a minimal static site generator for documentation.
It converts Markdown files into a fast, clean HTML site.

Markdown content here supports full formatting:
- Lists
- Code blocks
- Links
```

Markdown files allow richer answers than inline YAML.

## Category Configuration

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Category heading |
| `order` | int | Sort order |
| `items[].question` | string | Question text |
| `items[].answer` | string | Answer text (plain text in YAML) |

## Output

The FAQ page is generated at `/<path>/index.html`. Questions are grouped by category with collapsible answers.
