---
title: "Dark Mode"
description: "Configure and customize dark mode appearance"
tags: ["theme", "dark-mode", "accessibility"]
---

## Automatic Theme Detection

MinimalDoc respects the user's system preference by default. If the operating system is set to dark mode, the documentation will automatically display in dark mode.

## Manual Toggle

Users can override the system preference using the theme toggle button in the header. The preference is saved to localStorage and persists across sessions.

## Default Theme

Set the default theme in your configuration:

```yaml
# Start in dark mode
dark_mode: true
```

This sets the initial theme for users who haven't made a selection yet.

## CSS Variables

Dark mode uses CSS custom properties. The main variables are:

```css
[data-theme="dark"] {
  --bg-primary: #1a1a2e;
  --bg-secondary: #16213e;
  --text-primary: #eaeaea;
  --text-secondary: #b8b8b8;
  --accent-color: #4fc3f7;
  --border-color: #2d2d44;
}
```

## Custom Themes

Create a custom theme by overriding CSS variables in your own stylesheet. Place custom CSS in your docs directory and it will be copied to the output.
