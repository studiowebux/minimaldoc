---
title: Theming
description: Customize appearance with CSS variables and custom themes
tags:
  - features
  - theming
  - css
---

# Theming

MinimalDoc themes are CSS-only. Customize colors without editing templates.

## Built-in Themes

| Theme | Description |
|-------|-------------|
| `default` | Clean, neutral colors |
| `yellow` | Warm accent color |

Set theme in config:

```yaml
theme: default
```

Or via CLI:

```bash
minimaldoc build --theme yellow
```

## Dark Mode

Dark mode is built into every theme.

### Default Behavior

Respects system preference (`prefers-color-scheme`).

### Force Dark Mode

```yaml
dark_mode: true
```

### Toggle

Users can toggle via the button in the header. Preference persists in localStorage.

## CSS Variables

Themes define CSS variables for light and dark modes:

```css
/* Light mode */
:root[data-theme="light"] {
  --bg-primary: #fafafa;
  --bg-secondary: #f0f0f0;
  --bg-tertiary: #e8e8e8;
  --text-primary: #1a1a1a;
  --text-secondary: #666666;
  --text-tertiary: #999999;
  --text-muted: #bbbbbb;
  --border-primary: #e0e0e0;
  --border-secondary: #f0f0f0;
  --accent-primary: #007bff;
  --accent-hover: #0056b3;
  --link-color: #0066cc;
  --link-hover: #0052a3;
  --bg-code: #f5f5f5;
  --bg-hover: #f0f0f0;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.1);
}

/* Dark mode */
:root[data-theme="dark"] {
  --bg-primary: #1a1a1a;
  --bg-secondary: #2d2d2d;
  --bg-tertiary: #3a3a3a;
  --text-primary: #f5f5f5;
  --text-secondary: #b0b0b0;
  --text-tertiary: #888888;
  --text-muted: #666666;
  --border-primary: #444444;
  --border-secondary: #333333;
  --accent-primary: #4a9eff;
  --accent-hover: #357abd;
  --link-color: #66b3ff;
  --link-hover: #99ccff;
  --bg-code: #0d1117;
  --bg-hover: #2d2d2d;
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.3);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.5);
}
```

## Variable Reference

### Backgrounds

| Variable | Usage |
|----------|-------|
| `--bg-primary` | Main page background |
| `--bg-secondary` | Sidebar, cards |
| `--bg-tertiary` | Nested elements |
| `--bg-code` | Code blocks |
| `--bg-hover` | Hover states |

### Text

| Variable | Usage |
|----------|-------|
| `--text-primary` | Main content |
| `--text-secondary` | Secondary text |
| `--text-tertiary` | Muted text |
| `--text-muted` | Very light text |

### Borders

| Variable | Usage |
|----------|-------|
| `--border-primary` | Main borders |
| `--border-secondary` | Subtle borders |

### Accents

| Variable | Usage |
|----------|-------|
| `--accent-primary` | Buttons, highlights |
| `--accent-hover` | Hover states |
| `--link-color` | Link text |
| `--link-hover` | Link hover |

### Shadows

| Variable | Usage |
|----------|-------|
| `--shadow-sm` | Subtle elevation |
| `--shadow-md` | Cards, dropdowns |

## Config-Based Theming

Customize colors, fonts, and hero styling directly in `config.yaml` without creating theme files.

### Theme Configuration

```yaml
theme_config:
  name: "custom"
  colors:
    light:
      bg_primary: "#ffffff"
      bg_secondary: "#f8f9fa"
      text_primary: "#212529"
      link_color: "#0066cc"
    dark:
      bg_primary: "#0f172a"
      bg_secondary: "#1e293b"
      text_primary: "#f8fafc"
      link_color: "#60a5fa"
  fonts:
    heading: "Inter"
    body: "Inter"
    code: "JetBrains Mono"
    google_url: "https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&family=JetBrains+Mono&display=swap"
  hero:
    background_image: "/images/hero-bg.jpg"
    background_overlay: "rgba(0,0,0,0.6)"
    text_align: "center"
    min_height: "80vh"
```

### Color Options

| Field | Description |
|-------|-------------|
| `bg_primary` | Main page background |
| `bg_secondary` | Sidebar, cards |
| `bg_tertiary` | Nested elements |
| `bg_code` | Code block background |
| `bg_hover` | Hover states |
| `text_primary` | Main text |
| `text_secondary` | Secondary text |
| `text_tertiary` | Muted text |
| `text_muted` | Very light text |
| `border_primary` | Main borders |
| `border_secondary` | Subtle borders |
| `accent_primary` | Buttons, highlights |
| `accent_hover` | Accent hover state |
| `link_color` | Link text |
| `link_hover` | Link hover |

### Font Options

| Field | Description |
|-------|-------------|
| `heading` | Font family for headings |
| `body` | Font family for body text |
| `code` | Font family for code |
| `google_url` | Google Fonts import URL |

### Hero Options

| Field | Description |
|-------|-------------|
| `background_image` | Hero background image URL |
| `background_overlay` | Overlay color (e.g., `rgba(0,0,0,0.6)`) |
| `text_align` | Text alignment: `left`, `center`, `right` |
| `min_height` | Minimum hero height (e.g., `80vh`) |

### Partial Configuration

Only specify values you want to override. Unspecified values use theme defaults:

```yaml
theme_config:
  colors:
    light:
      link_color: "#ff6600"
    dark:
      link_color: "#ff9933"
```

## Creating Custom Themes (File-Based)

1. Create theme directory:

```
themes/
└── my-theme/
    └── static/
        └── css/
            └── main.css
```

2. Define variables in `main.css`:

```css
/* my-theme/static/css/main.css */

:root[data-theme="light"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f8f9fa;
  --text-primary: #212529;
  --accent-primary: #6366f1;
  /* ... other variables */
}

:root[data-theme="dark"] {
  --bg-primary: #0f172a;
  --bg-secondary: #1e293b;
  --text-primary: #f1f5f9;
  --accent-primary: #818cf8;
  /* ... other variables */
}
```

3. Use your theme:

```bash
minimaldoc build --theme my-theme
```

## Custom Fonts

Add font import and override:

```css
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap');

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
}

code, pre {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}
```

## Theme Structure

```
themes/
├── common/                 # Shared (don't modify)
│   ├── templates/          # HTML templates
│   └── static/
│       ├── css/            # Base styles
│       └── js/             # JavaScript
├── default/                # Default theme
│   └── static/css/main.css # Variables only (~65 lines)
└── yellow/                 # Yellow theme
    └── static/css/main.css # Variables + font
```

Themes only contain CSS variables. All HTML and JS is shared.

## Admonition Colors

Admonition colors can be customized:

```css
.admonition.info {
  --admonition-bg: #e3f2fd;
  --admonition-border: #2196f3;
  --admonition-text: #1565c0;
}

.admonition.warning {
  --admonition-bg: #fff3e0;
  --admonition-border: #ff9800;
  --admonition-text: #e65100;
}

/* Add dark mode variants */
[data-theme="dark"] .admonition.info {
  --admonition-bg: rgba(33, 150, 243, 0.1);
  --admonition-border: #2196f3;
  --admonition-text: #90caf9;
}
```
