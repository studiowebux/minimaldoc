---
title: Landing Pages
description: Create beautiful landing pages, portfolios, and contact pages
menu_order: 8
---

# Landing Pages

MinimalDoc supports landing pages, portfolio pages, and contact pages. Configure them in `config.yaml` and MinimalDoc renders elegant templates.

## Landing Page

The landing page replaces the default documentation index with a marketing-style homepage.

### Configuration

```yaml
landing:
  enabled: true

  hero:
    title: "Your Product Name"
    subtitle: "A compelling tagline for your product"
    buttons:
      - text: "Get Started"
        url: "/docs/getting-started/"
        primary: true
      - text: "View on GitHub"
        url: "https://github.com/..."
        primary: false
    image: "/assets/hero.png"  # optional

  features:
    title: "Features"
    items:
      - emoji: "~"
        title: "Fast"
        description: "Sub-second builds"
      - icon: "moon"
        title: "Dark Mode"
        description: "Built-in theme switching"

  steps:
    title: "Quick Start"
    items:
      - title: "Install"
        description: "Download the binary"
        code: "go install github.com/..."
      - title: "Initialize"
        description: "Create your docs folder"
        code: "minimaldoc init my-docs"

  links:
    title: "Resources"
    items:
      - icon: "github"
        title: "GitHub"
        description: "View source code"
        url: "https://github.com/..."

  testimonials:
    items:
      - quote: "MinimalDoc changed how we write docs"
        author: "Jane Doe"
        role: "Lead Developer"
        avatar: "/assets/jane.png"

  cta:
    title: "Ready to get started?"
    description: "Join thousands of developers"
    buttons:
      - text: "Download Now"
        url: "/releases/"
        primary: true

  opensource:
    title: "Open Source"
    description: "MIT Licensed. Free forever."
    links:
      - text: "GitHub"
        url: "https://github.com/..."
```

### Sections

| Section | Description |
|---------|-------------|
| `hero` | Main header with title, subtitle, buttons, and optional image |
| `features` | Grid of feature cards with icons/emojis |
| `steps` | Numbered quick start steps with code blocks |
| `links` | Grid of resource link cards |
| `testimonials` | Customer quotes with avatars |
| `cta` | Call-to-action section |
| `opensource` | Open source info and links |

All sections are optional. Only configured sections are rendered.

## Portfolio Page

Display projects parsed from markdown files in `__portfolio__/`.

### Configuration

```yaml
portfolio:
  enabled: true
  title: "Portfolio"
  description: "Projects and experiments"
  path: "portfolio"
```

### Project Files

Create markdown files in `docs/__portfolio__/`:

```markdown
---
title: Project Name
description: Short description
image: /assets/project.png
tags:
  - go
  - cli
  - open-source
links:
  - text: GitHub
    url: https://github.com/...
  - text: Demo
    url: https://demo.example.com
date: 2026-01-15
featured: true
---

Longer description here. Supports full markdown.
```

### Frontmatter Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Project title |
| `description` | string | Short description for cards |
| `image` | string | Project image URL |
| `tags` | array | Tags for filtering |
| `links` | array | Project links (text, url) |
| `date` | string | Date (YYYY-MM-DD) |
| `featured` | boolean | Show in featured section |
| `menu_order` | number | Sort order |

### Features

- Responsive card grid layout
- Client-side tag filtering
- Individual project detail pages
- Hover effects

## Contact Page

Simple contact page with email link and contact info.

### Configuration

```yaml
contact:
  enabled: true
  title: "Contact Us"
  description: "Get in touch"
  path: "contact"
  email: "hello@example.com"
  info:
    - icon: "mail"
      text: "hello@example.com"
    - icon: "location"
      text: "Remote, Worldwide"
    - icon: "phone"
      text: "+1 555 123 4567"
```

### Available Icons

- `mail` - Email icon
- `location` - Map pin icon
- `phone` - Phone icon

## Footer

Configure a multi-column footer for landing pages.

### Configuration

```yaml
footer:
  copyright: "2026 Your Company. All rights reserved."

  links:
    - title: "Product"
      items:
        - text: "Features"
          url: "/features/"
        - text: "Documentation"
          url: "/docs/"
    - title: "Company"
      items:
        - text: "About"
          url: "/about/"
        - text: "Blog"
          url: "/blog/"

  social:
    - icon: "github"
      name: "GitHub"
      url: "https://github.com/..."
    - icon: "twitter"
      name: "Twitter"
      url: "https://twitter.com/..."

  badges:
    - text: "Powered by MinimalDoc"
      url: "https://minimaldoc.dev"
```

### Social Icons

Available icons: `github`, `twitter`, `linkedin`, `youtube`, `discord`, `mail`

## Styling

All landing page elements use CSS variables for theming:

- Light/dark mode automatic switching
- Responsive design (mobile-first)
- Clean, minimal aesthetic
- Consistent with documentation theme

Custom CSS can be added via theme overrides.
