---
title: Page Title
description: Brief description for SEO and page header
tags:
  - tag1
  - tag2
author: Author Name
menu_title: Short Nav Title
menu_order: 10
hidden: false
# full_width: true
# no_header: true
# stale_warning: false
# stale_threshold_days: 730
seo:
  title: SEO Title Override
  description: SEO description override
  keywords:
    - keyword1
    - keyword2
  image: /images/og-image.png
  canonical: https://example.com/canonical-url
  noindex: false
  nofollow: false
---

# Page Title

Introduction paragraph describing the page content.

## Section 1

Content for section 1.

### Subsection 1.1

More detailed content.

## Section 2

Content for section 2.

:::info
Use admonitions for callouts: info, warning, danger, success, note, question
:::

## Code Examples

```language
// Code block with syntax highlighting
```

## Tables

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Data     | Data     | Data     |

<!--
FRONTMATTER REFERENCE:

title: (optional) Page title, defaults to first h1 or filename
description: (optional) Shown in page header and meta description
tags: (optional) Array of tags displayed on page
author: (optional) Page author
menu_title: (optional) Shorter title for navigation
menu_order: (optional) Sort order in navigation (lower = earlier)
hidden: (optional) Set true to exclude from navigation
seo: (optional) SEO overrides
  title: Override page title in <title> tag
  description: Override meta description
  keywords: Meta keywords array
  image: Open Graph image URL
  canonical: Canonical URL
  noindex: Set true for noindex robots tag
  nofollow: Set true for nofollow robots tag

NAVIGATION ORDER:
1. Files with menu_order (ascending)
2. Files with numeric prefixes (01-, 02-)
3. Alphabetically by filename
-->
