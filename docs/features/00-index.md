---
title: Feature Index
description: Complete list of all MinimalDoc features
menu_order: 0
---

# Feature Index

Complete list of all MinimalDoc features with one-liner descriptions.

## Content

| Feature | Description |
|---------|-------------|
| Markdown | GitHub Flavored Markdown with tables, task lists, strikethrough |
| Frontmatter | YAML metadata for title, description, tags, SEO, ordering |
| Admonitions | Callout blocks: info, warning, danger, success, note, question |
| Syntax Highlighting | 100+ languages via Chroma with one-click copy |
| Stale Warnings | Configurable warnings for content older than N days |
| Link Checking | Validate internal/external links during build |
| Anchor Links | Auto-generated heading IDs for deep linking |
| External Links | Auto-open in new tab with visual indicator |

## Navigation

| Feature | Description |
|---------|-------------|
| Auto Navigation | Generated from folder structure with numeric prefixes |
| Custom TOC | Override navigation order via TOC.md file |
| Page TOC | Right sidebar table of contents from headings |
| Scrollspy | Active heading highlight while scrolling |
| Collapsible Sections | Expand/collapse navigation groups with state persistence |
| Breadcrumbs | Path navigation in Knowledge Base and versioned docs |

## Search

| Feature | Description |
|---------|-------------|
| Full-Text Search | Client-side inverted index with weighted scoring |
| Keyboard Shortcut | Cmd+K / Ctrl+K to open search modal |
| Prefix Matching | Find results while typing partial words |
| Section Results | Jump directly to matching headings within pages |
| Per-Version Index | Separate search indexes for each documentation version |
| KB/FAQ Indexing | Knowledge Base and FAQ items included in search |

## Design

| Feature | Description |
|---------|-------------|
| Responsive | Mobile, tablet, desktop layouts with hamburger menu |
| Dark Mode | Toggle with localStorage persistence and system preference |
| Built-in Themes | Default (neutral) and Yellow themes included |
| Custom Colors | Configure all color variables in YAML |
| Custom Fonts | Google Fonts integration via config |
| Hero Backgrounds | Background images with overlay for landing pages |
| Skip Links | Accessibility: skip to main content link |
| Focus Indicators | Visible focus states for keyboard navigation |

## Pages

| Feature | Description |
|---------|-------------|
| Landing Page | Marketing homepage with hero, features, steps, CTA sections |
| Markdown Landing | Define landing sections in `__landing__/*.md` files |
| Knowledge Base | Self-service support hub with categories and article counts |
| Portfolio | Project showcase with tags, images, and filtering |
| Contact | Contact information page with email and details |
| FAQ | Collapsible Q&A with categories and deep linking |
| Legal | Privacy policy, terms of service with auto footer links |
| Status Page | Service health dashboard with incidents and maintenance |
| Changelog | Version history with categories and RSS feed |

## Versioning

| Feature | Description |
|---------|-------------|
| Multi-Version | Maintain docs for multiple software versions |
| Version Selector | Dropdown to switch between versions |
| Version Overrides | Override specific pages in `__versions__/{version}/` |
| EOL Warnings | Warning banners for end-of-life versions |
| Since/Deprecated | Badges showing when features were added or deprecated |
| Per-Version Search | Separate search indexes per version |

## Internationalization

| Feature | Description |
|---------|-------------|
| Multi-Locale | Support multiple languages with locale selector |
| Translation Files | Translations in `__translations__/{locale}/` directory |
| Fallback Content | Show default locale content when translation missing |
| RTL Support | Right-to-left text direction for Arabic, Hebrew, etc. |
| Locale Persistence | Remember language preference in localStorage |
| URL Prefixes | SEO-friendly `/fr/`, `/de/` URL structure |

## API Documentation

| Feature | Description |
|---------|-------------|
| OpenAPI Support | Interactive docs from OpenAPI/Swagger specs |
| Multiple Views | Organize endpoints by path, tag, or flat list |
| Live Testing | Interactive API tester with auth support |
| Code Samples | Auto-generated curl, JavaScript, Go, Python, Swift |
| Schema Viewer | Collapsible explorer for complex types |
| Response Tabs | Tabbed interface for different status codes |
| Export | Generate cURL commands and restcli configs |

## Status Page

| Feature | Description |
|---------|-------------|
| Component Status | Track services with operational/degraded/outage states |
| Incident Timeline | Markdown incidents with update history |
| Maintenance Windows | Display scheduled maintenance with dates |
| Uptime Calendar | 90-day visual uptime grid (GitHub-style) |
| SLA Display | Track targets with 7d/30d/90d breakdowns |
| Health Checks | Browser-based endpoint polling with latency |
| RSS Feed | Subscribe to incident updates |
| JSON API | Machine-readable status at `/status/status.json` |

## SEO

| Feature | Description |
|---------|-------------|
| Sitemap | Automatic sitemap.xml with all pages |
| Meta Tags | Title, description, author, keywords |
| Open Graph | Facebook/LinkedIn preview cards |
| Twitter Cards | Twitter preview cards |
| Canonical URLs | Configurable canonical link tags |
| Robots Meta | noindex/nofollow per page |
| LLM Output | llms.txt and llms-full.txt for AI assistants |

## Analytics

| Feature | Description |
|---------|-------------|
| Plugin System | Multiple providers can run simultaneously |
| MinimalDoc Backend | Self-hosted cookie-free analytics with feedback and newsletter |
| Google Analytics 4 | Native support with measurement ID |
| Plausible | Privacy-friendly analytics (hosted or self-hosted) |
| Umami | Open-source analytics support |
| Matomo | Self-hosted analytics support |
| Fathom | Privacy-focused analytics |
| Simple Analytics | Minimal tracking solution |
| Custom Provider | Arbitrary script injection with custom attributes |

## Backend Server (Optional)

| Feature | Description |
|---------|-------------|
| Self-Hosted | Single binary + SQLite/PostgreSQL, runs anywhere |
| Cookie-Free Analytics | Privacy-first tracking without cookies or PII |
| Feedback Widget | Page rating with optional comments |
| Newsletter | Email subscription with double opt-in verification |
| Email System | SMTP support with verification and welcome emails |
| Admin Portal | Web dashboard for analytics, feedback, subscribers |
| API | REST endpoints for all features |
| Docker Ready | Dockerfile and docker-compose included |

## Build

| Feature | Description |
|---------|-------------|
| Single Binary | No dependencies, just one executable |
| Fast Builds | Go-powered, handles large sites quickly |
| Single File Mode | Build one markdown file as standalone page |
| Config File | Optional config.yaml for persistent settings |
| CLI Overrides | All config options available as flags |
| Clean URLs | Optional `/page/` instead of `/page.html` |
| Custom Entrypoint | Use README.md or any file as homepage |
| Version Check | `minimaldoc version --check` for updates |

## Output

| Feature | Description |
|---------|-------------|
| Static HTML | Pure HTML/CSS/JS, no server runtime |
| Asset Copying | Static files copied from docs directory |
| Icons SVG | Social and UI icons as SVG sprite |
| Version Badge | "Powered by MinimalDoc" footer link |
| Source Maps | CSS source maps for debugging |
