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
    image: "/assets/hero.png" # optional

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

| Section        | Description                                                   |
| -------------- | ------------------------------------------------------------- |
| `hero`         | Main header with title, subtitle, buttons, and optional image |
| `features`     | Grid of feature cards with icons/emojis                       |
| `steps`        | Numbered quick start steps with code blocks                   |
| `links`        | Grid of resource link cards                                   |
| `testimonials` | Customer quotes with avatars                                  |
| `cta`          | Call-to-action section                                        |
| `opensource`   | Open source info and links                                    |

All sections are optional. Only configured sections are rendered.

## New Section Types

### Image-Text Section

Side-by-side image and content layout.

```yaml
landing:
  image_text:
    - id: "why-choose"
      title: "Why Choose Us"
      description: "We deliver results"
      image: "/images/feature.png"
      image_alt: "Feature illustration"
      image_position: "right"  # left or right
      order: 1
      items:
        - icon: "~"
          title: "Fast Delivery"
          description: "Quick turnaround times"
        - icon: "@"
          title: "Quality First"
          description: "No compromises"
      buttons:
        - text: "Learn More"
          url: "/about/"
          primary: true
      background:
        color: "#f8f9fa"
```

| Field | Description |
|-------|-------------|
| `id` | Section identifier |
| `title` | Section heading |
| `description` | Section description |
| `image` | Image URL |
| `image_alt` | Image alt text |
| `image_position` | `left` or `right` (default: `right`) |
| `items` | Optional bullet points with icon, title, description |
| `buttons` | Optional CTA buttons |
| `background` | Background styling |
| `order` | Display order |

### Text Section

Simple centered text block.

```yaml
landing:
  text_blocks:
    - id: "mission"
      title: "Our Mission"
      subtitle: "What drives us"
      content: "<p>We believe in making documentation accessible to everyone.</p>"
      alignment: "center"  # left, center, right
      max_width: "800px"
      order: 2
      buttons:
        - text: "Join Us"
          url: "/careers/"
          primary: true
      background:
        image: "/images/bg.jpg"
        overlay: "rgba(0,0,0,0.7)"
```

| Field | Description |
|-------|-------------|
| `id` | Section identifier |
| `title` | Section heading |
| `subtitle` | Optional subtitle |
| `content` | HTML content |
| `alignment` | Text alignment: `left`, `center`, `right` |
| `max_width` | Maximum content width |
| `buttons` | Optional CTA buttons |
| `background` | Background styling |
| `order` | Display order |

### Links Grid Section

Grid of external link cards.

```yaml
landing:
  links_grid:
    - id: "platforms"
      title: "Find Us On"
      description: "Connect with us across platforms"
      columns: 4  # 2, 3, or 4
      order: 3
      items:
        - icon: "github"
          title: "GitHub"
          description: "Source code and issues"
          url: "https://github.com/..."
          external: true
        - icon: "twitter"
          title: "Twitter"
          description: "Updates and news"
          url: "https://twitter.com/..."
      background:
        color: "#1a1a2e"
```

| Field | Description |
|-------|-------------|
| `id` | Section identifier |
| `title` | Section heading |
| `description` | Section description |
| `columns` | Grid columns: 2, 3, or 4 (default: 4) |
| `items` | Array of link cards |
| `background` | Background styling |
| `order` | Display order |

**Link card fields:**

| Field | Description |
|-------|-------------|
| `icon` | Icon or emoji |
| `title` | Card title |
| `description` | Card description |
| `url` | Link URL |
| `external` | Opens in new tab (auto-detected for http/https) |

## Section Backgrounds

All sections support background styling:

```yaml
background:
  image: "/images/bg.jpg"      # Background image URL
  overlay: "rgba(0,0,0,0.6)"   # Overlay color
  color: "#1a1a2e"             # Background color
  position: "center"           # Background position
  size: "cover"                # Background size
  attachment: "fixed"          # scroll or fixed
```

| Field | Description |
|-------|-------------|
| `image` | Background image URL |
| `overlay` | Semi-transparent overlay (use with images) |
| `color` | Solid background color |
| `position` | CSS background-position |
| `size` | CSS background-size |
| `attachment` | CSS background-attachment |

Example with background image:

```yaml
landing:
  cta:
    title: "Ready to Start?"
    background:
      image: "/images/cta-bg.jpg"
      overlay: "rgba(0,0,0,0.7)"
```

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

| Field         | Type    | Description                 |
| ------------- | ------- | --------------------------- |
| `title`       | string  | Project title               |
| `description` | string  | Short description for cards |
| `image`       | string  | Project image URL           |
| `tags`        | array   | Tags for filtering          |
| `links`       | array   | Project links (text, url)   |
| `date`        | string  | Date (YYYY-MM-DD)           |
| `featured`    | boolean | Show in featured section    |
| `menu_order`  | number  | Sort order                  |

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
