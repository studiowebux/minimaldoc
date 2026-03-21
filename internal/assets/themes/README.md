# Theme Structure

Themes in minimaldoc are CSS-only. All HTML structure, JavaScript functionality, and templates are shared. Themes only define colors and visual styling through CSS.

## Directory Structure

```
themes/
├── common/               # Shared assets (structure, behavior, templates)
│   ├── static/
│   │   ├── css/         # Base styles (structure, layout, components)
│   │   │   ├── base.css
│   │   │   └── openapi.css
│   │   └── js/          # All JavaScript functionality
│   └── templates/       # All HTML structure
│       ├── layout.html
│       ├── openapi-page.html
│       └── partials/
│           ├── nav.html
│           ├── openapi-nav.html
│           └── toc.html
├── default/             # Default theme (colors only)
│   └── static/
│       └── css/
│           └── main.css # 65 lines (CSS variables only)
└── yellow/              # Yellow theme (colors only)
    └── static/
        └── css/
            └── main.css # 67 lines (CSS variables + font import)
```

## Common Assets

### Templates (themes/common/templates/)
All HTML structure is shared:
- `layout.html` - Main page layout
- `openapi-page.html` - OpenAPI documentation page layout
- `partials/nav.html` - Navigation structure
- `partials/toc.html` - Table of contents structure
- `partials/openapi-nav.html` - OpenAPI navigation structure

### JavaScript (themes/common/static/js/)
All behavior is shared:
- `nav-collapse.js` - Navigation collapse/expand with localStorage persistence
- `anchor-links.js` - Anchor link handling
- `external-links.js` - External link handling
- `scrollspy.js` - Scrollspy functionality
- `search.js` - Search modal and functionality
- `theme-toggle.js` - Theme switcher (light/dark mode)
- `mobile-menu.js` - Mobile menu with auto-detection for OpenAPI features

## How Theme Loading Works

1. **Templates**: All templates loaded from `themes/common/`
2. **JavaScript**: Common JS loaded from `themes/common/`, then theme-specific JS (if any) from `themes/{theme}/`
3. **CSS**: Only theme-specific CSS loaded from `themes/{theme}/static/css/`

## Creating a New Theme

Themes are CSS-only. No HTML or JavaScript needed.

1. Create a new directory:
   ```
   themes/your-theme/
   └── static/
       └── css/
           └── main.css
   ```

2. Define your colors in `main.css`:
   ```css
   /* Light mode colors */
   :root[data-theme="light"] {
       --bg-primary: #fff;
       --text-primary: #000;
       --accent-primary: #your-color;
       --link-color: #your-link-color;
       /* ... other variables */
   }

   /* Dark mode colors */
   :root[data-theme="dark"] {
       --bg-primary: #000;
       --text-primary: #fff;
       --accent-primary: #your-color;
       --link-color: #your-link-color;
       /* ... other variables */
   }

   /* Optional: Custom font */
   body {
       font-family: "Your Font", -apple-system, ...;
   }
   ```

3. Build with your theme:
   ```
   ./minimaldoc build --theme your-theme
   ```

## Available CSS Variables

Themes should define these CSS variables for both light and dark modes:

**Backgrounds:**
- `--bg-primary` - Main background
- `--bg-secondary` - Secondary background (sidebars)
- `--bg-tertiary` - Tertiary background (hover states)
- `--bg-code` - Code block background
- `--bg-hover` - Hover state background

**Text:**
- `--text-primary` - Main text color
- `--text-secondary` - Secondary text
- `--text-tertiary` - Tertiary text
- `--text-muted` - Muted/disabled text

**Borders:**
- `--border-primary` - Main borders
- `--border-secondary` - Secondary borders

**Accents:**
- `--accent-primary` - Primary accent (active states)
- `--accent-hover` - Accent hover state

**Links:**
- `--link-color` - Link color
- `--link-hover` - Link hover color

**Shadows:**
- `--shadow-sm` - Small shadow
- `--shadow-md` - Medium shadow

## Theme Examples

### Default Theme
Neutral grayscale design with blue links

### Yellow Theme
Bright yellow accents with Montserrat font

## Best Practices

- Only modify CSS variables, never add new selectors
- Provide both light and dark mode colors
- Test color contrast for accessibility
- Keep accent colors consistent
- Consider font pairing carefully
