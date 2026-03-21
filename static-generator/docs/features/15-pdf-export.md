---
title: PDF Export
description: Export documentation pages as PDF files
tags:
  - features
  - pdf
  - export
---

# PDF Export

Add an "Export PDF" button to each documentation page. Uses the browser's native print-to-PDF.

## Configuration

```yaml
pdf_export:
  enabled: true
  page_break_level: 2
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Show the Export PDF button on each page |
| `page_break_level` | int | `1` | Heading level for page breaks: 0=none, 1=h1 only, 2=h1+h2 |

## How It Works

The button triggers `window.print()` with a print-optimized stylesheet. The generated PDF:

- Strips sidebar, navigation, and interactive elements
- Preserves code blocks with syntax highlighting
- Inserts page breaks before headings (configurable level)
- Includes the page title and URL in headers/footers (browser-dependent)

## Page Break Levels

| Value | Behavior |
|-------|----------|
| `0` | No forced page breaks |
| `1` | Page break before every `h1` |
| `2` | Page break before every `h1` and `h2` |

## Button Placement

The button appears in the page title row alongside other action buttons (Copy MD, Ask Claude). On mobile, buttons wrap below the title.

## Print Stylesheet

MinimalDoc includes a `@media print` stylesheet that:

- Hides the sidebar, header, footer, and search
- Expands content to full width
- Preserves code block backgrounds
- Forces black text for readability
- Removes interactive elements (buttons, toggles)

## Limitations

- Output depends on the browser's print engine (Chrome recommended)
- Complex tables may not paginate cleanly
- Embedded images must be accessible at print time
- Header/footer content varies by browser
