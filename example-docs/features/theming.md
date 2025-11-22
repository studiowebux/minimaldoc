---
title: Theming & Customization
description: Customize the look and feel of your documentation
tags:
  - theming
  - customization
  - design
---

# Theming & Customization

Customize MinimalDoc to match your brand and style preferences.

## Default Theme

MinimalDoc includes a beautiful default theme with:

- Clean, minimal design
- Light and dark mode
- Responsive layout
- Beautiful typography
- Syntax highlighting

## Dark Mode

### Automatic Dark Mode

The theme switcher in the sidebar allows users to toggle between light and dark modes. The preference is saved in localStorage.

### Color Scheme

**Light Mode:**
- Background: Soft white (#fafafa)
- Text: Soft black (#1a1a1a)
- Links: Blue (#2563eb)
- Borders: Light gray (#e0e0e0)

**Dark Mode:**
- Background: Soft black (#1a1a1a)
- Text: Soft white (#f5f5f5)
- Links: Light blue (#60a5fa)
- Borders: Dark gray (#3a3a3a)

## Typography

### Fonts

The default theme uses system fonts for optimal performance:

```css
font-family: -apple-system, BlinkMacSystemFont,
  "Segoe UI", Roboto, "Helvetica Neue",
  Arial, sans-serif;
```

### Font Sizes

- Base: 16px
- H1: 2.5rem (40px)
- H2: 1.75rem (28px)
- H3: 1.5rem (24px)
- H4: 1.25rem (20px)
- H5: 1.125rem (18px)
- H6: 1rem (16px)

## Layout

### Three-Column Layout (Desktop)

```
┌─────────┬─────────────┬──────────┐
│         │             │          │
│  Nav    │   Content   │   TOC    │
│ 280px   │    Fluid    │  240px   │
│         │             │          │
└─────────┴─────────────┴──────────┘
```

### Two-Column Layout (Tablet)

```
┌─────────┬─────────────┐
│         │             │
│  Nav    │   Content   │
│ 260px   │    Fluid    │
│         │             │
└─────────┴─────────────┘
```

### Single Column (Mobile)

```
┌─────────────────┐
│                 │
│    Content      │
│     (Full)      │
│                 │
└─────────────────┘
[≡] Hamburger Menu
```

## Code Highlighting

### Theme

Syntax highlighting uses the Monokai color scheme:

```javascript
// Beautiful syntax highlighting
const colors = {
  keyword: '#66d9ef',
  string: '#e6db74',
  number: '#ae81ff',
  comment: '#75715e',
  function: '#a6e22e'
};
```

### Supported Languages

100+ languages including:
- JavaScript/TypeScript
- Python
- Go
- Rust
- Java
- C/C++
- PHP
- Ruby
- And many more!

## Custom Styling

### CSS Variables

(Future feature) Override theme variables:

```css
:root[data-theme="light"] {
  --link-color: #your-color;
  --bg-primary: #your-bg;
}
```

### Custom CSS

(Future feature) Add custom CSS file:

```bash
minimaldoc build . --custom-css="custom.css"
```

## Customization Roadmap

Planned customization features:

- [ ] Custom theme support
- [ ] CSS variable overrides
- [ ] Custom fonts
- [ ] Logo customization
- [ ] Color scheme editor
- [ ] Multiple theme presets

## Best Practices

:::success Theming Tips

1. **Maintain contrast** - Ensure text is readable
2. **Test both modes** - Check light and dark themes
3. **Use system fonts** - Better performance
4. **Keep it simple** - Don't over-customize
5. **Test on mobile** - Verify responsive design
:::

## Accessibility

The default theme follows accessibility best practices:

- Sufficient color contrast (WCAG AA)
- Focus indicators
- Semantic HTML
- Screen reader friendly
- Keyboard navigable

## Next Steps

- [See the API reference](../api/reference.html)
- [Learn about deployment](../guides/deployment.html)
- [Read the FAQ](../guides/faq.html)
