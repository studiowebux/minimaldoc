---
title: Custom Themes
description: Create custom themes with CSS variables
tags:
  - guides
  - theming
  - css
---

# Custom Themes

Create themes by defining CSS variables. No template editing required.

## Theme Structure

```
themes/
├── common/                 # Shared (don't modify)
│   ├── templates/
│   └── static/
│       ├── css/
│       └── js/
├── default/                # Default theme
│   └── static/css/main.css
└── my-theme/               # Your custom theme
    └── static/css/main.css
```

## Creating a Theme

### 1. Create Directory

```bash
mkdir -p themes/my-theme/static/css
```

### 2. Create main.css

```css
/* themes/my-theme/static/css/main.css */

/* Light mode */
:root[data-theme="light"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f8fafc;
  --bg-tertiary: #f1f5f9;
  --text-primary: #0f172a;
  --text-secondary: #475569;
  --text-tertiary: #94a3b8;
  --text-muted: #cbd5e1;
  --border-primary: #e2e8f0;
  --border-secondary: #f1f5f9;
  --accent-primary: #3b82f6;
  --accent-hover: #2563eb;
  --link-color: #2563eb;
  --link-hover: #1d4ed8;
  --bg-code: #f8fafc;
  --bg-hover: #f1f5f9;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.1);
}

/* Dark mode */
:root[data-theme="dark"] {
  --bg-primary: #0f172a;
  --bg-secondary: #1e293b;
  --bg-tertiary: #334155;
  --text-primary: #f8fafc;
  --text-secondary: #cbd5e1;
  --text-tertiary: #94a3b8;
  --text-muted: #64748b;
  --border-primary: #334155;
  --border-secondary: #1e293b;
  --accent-primary: #60a5fa;
  --accent-hover: #3b82f6;
  --link-color: #60a5fa;
  --link-hover: #93c5fd;
  --bg-code: #1e293b;
  --bg-hover: #334155;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.3);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.5);
}
```

### 3. Use Your Theme

```bash
minimaldoc build --theme my-theme
```

Or in config:

```yaml
theme: my-theme
```

## CSS Variables Reference

### Backgrounds

| Variable | Purpose |
|----------|---------|
| `--bg-primary` | Page background |
| `--bg-secondary` | Sidebar, cards |
| `--bg-tertiary` | Nested elements, highlights |
| `--bg-code` | Code block background |
| `--bg-hover` | Hover states |

### Text

| Variable | Purpose |
|----------|---------|
| `--text-primary` | Main content |
| `--text-secondary` | Secondary text, descriptions |
| `--text-tertiary` | Muted text |
| `--text-muted` | Very subtle text |

### Borders

| Variable | Purpose |
|----------|---------|
| `--border-primary` | Main borders |
| `--border-secondary` | Subtle borders |

### Accents

| Variable | Purpose |
|----------|---------|
| `--accent-primary` | Buttons, active states |
| `--accent-hover` | Hover on accent elements |
| `--link-color` | Link text |
| `--link-hover` | Link hover |

### Shadows

| Variable | Purpose |
|----------|---------|
| `--shadow-sm` | Subtle depth |
| `--shadow-md` | Cards, dropdowns |

## Custom Fonts

### Google Fonts

```css
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

:root {
  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
}

body {
  font-family: var(--font-sans);
}
```

### Monospace Font

```css
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono&display=swap');

code, pre {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}
```

### Self-Hosted Fonts

```css
@font-face {
  font-family: 'CustomFont';
  src: url('/fonts/custom.woff2') format('woff2');
  font-weight: normal;
  font-style: normal;
  font-display: swap;
}
```

## Theme Examples

### Warm Theme

```css
:root[data-theme="light"] {
  --bg-primary: #fffbeb;
  --bg-secondary: #fef3c7;
  --text-primary: #78350f;
  --accent-primary: #f59e0b;
  --link-color: #d97706;
}

:root[data-theme="dark"] {
  --bg-primary: #1c1917;
  --bg-secondary: #292524;
  --text-primary: #fef3c7;
  --accent-primary: #fbbf24;
  --link-color: #fcd34d;
}
```

### Cool Theme

```css
:root[data-theme="light"] {
  --bg-primary: #f0f9ff;
  --bg-secondary: #e0f2fe;
  --text-primary: #0c4a6e;
  --accent-primary: #0284c7;
  --link-color: #0369a1;
}

:root[data-theme="dark"] {
  --bg-primary: #0c1929;
  --bg-secondary: #1e3a5f;
  --text-primary: #e0f2fe;
  --accent-primary: #38bdf8;
  --link-color: #7dd3fc;
}
```

### High Contrast

```css
:root[data-theme="light"] {
  --bg-primary: #ffffff;
  --text-primary: #000000;
  --accent-primary: #0000ff;
  --border-primary: #000000;
}

:root[data-theme="dark"] {
  --bg-primary: #000000;
  --text-primary: #ffffff;
  --accent-primary: #00ffff;
  --border-primary: #ffffff;
}
```

## Admonition Colors

Customize admonition appearance:

```css
/* Info */
.admonition.info {
  --admonition-bg: #dbeafe;
  --admonition-border: #3b82f6;
  --admonition-text: #1e40af;
}

[data-theme="dark"] .admonition.info {
  --admonition-bg: rgba(59, 130, 246, 0.15);
  --admonition-text: #93c5fd;
}

/* Warning */
.admonition.warning {
  --admonition-bg: #fef3c7;
  --admonition-border: #f59e0b;
  --admonition-text: #92400e;
}

/* Danger */
.admonition.danger {
  --admonition-bg: #fee2e2;
  --admonition-border: #ef4444;
  --admonition-text: #991b1b;
}
```

## Code Block Colors

Override syntax highlighting:

```css
/* Code block background */
pre {
  background: var(--bg-code) !important;
}

/* Line numbers */
.chroma .ln {
  color: var(--text-muted);
}
```

## Testing Themes

### Local Development

```bash
# Build with theme
minimaldoc build ./docs --theme my-theme --output dist

# Serve locally
python -m http.server -d dist 8080
```

### Test Both Modes

Toggle dark mode in browser to verify both color schemes.

### Checklist

- [ ] Light mode readable
- [ ] Dark mode readable
- [ ] Links visible
- [ ] Code blocks legible
- [ ] Admonitions distinct
- [ ] Hover states visible
- [ ] Borders appropriate
- [ ] Search input styled

## Distributing Themes

### Share Theme

Package as a directory:

```
my-theme/
└── static/css/main.css
```

### Install Theme

Copy to project's `themes/` directory:

```bash
cp -r my-theme ./themes/
```

### Community Themes

Share on GitHub with tag `minimaldoc-theme`.
