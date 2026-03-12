# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.1] - 2026-03-12

### Fixed

- `code-copy.js` was incorrectly excluded from build output when OpenAPI was
  disabled, causing a 404 on code block copy buttons across all deployed sites
- Documentation article pages were capped at 768px width; now matches the
  category page width at 1280px

## [1.4.0] - 2026-03-11

### Added

- Forum system: categories, topics, posts, tags, likes, reputation, moderation,
  bans
- Public authentication: login/register pages, email verification
- OAuth on public router for newsletter subscribe
- Blog comment authentication: auto-fill name/email for logged-in users
- Custom analytics events with persistent storage
- Waitlist landing page type with newsletter CTA
- Copy-to-markdown button for documentation pages
- llms.txt footer link for AI discoverability
- Roadmap page generator (board and timeline layouts, tag filtering)
- CSS design system with unified token architecture
- Structured logging with `slog` (JSON in production, text in development)
- Health check endpoints: `/healthz` (liveness), `/readyz` (readiness with DB
  ping)
- Makefile targets: `lint`, `check`, `security`, `coverage`, `ci`
- CI pipeline with GitHub Actions (lint, test, security)

### Changed

- Split `queries.go` (1,471 lines) into domain-specific files
- Split `handlers.go` (1,337 lines) into domain-specific files
- Replaced all `log.*` calls in server code with `slog` structured logging
- Added `SERVER_ENV` environment variable (default: `development`)
- Eliminated CSS duplication, tokenized all hardcoded values

### Fixed

- 12 security vulnerabilities (CSP headers, input validation, error discarding)
- Rate limiting for login, API, and submission endpoints
- Code deduplication across handler files

## [1.3.0] - 2025-12-15

### Added

- Backend server (`minimaldoc-server`) with dual-port architecture (public
  :8080, admin :8090)
- Cookie-free analytics (page views, duration, bounce rate, session tracking)
- Feedback collection widget
- Newsletter subscription with email verification
- Authentication: JWT, bcrypt, OAuth 2.0/OIDC (Cognito, Auth0, Google, GitHub)
- Admin dashboard with HTMX fragments
- Blog CMS: markdown editor, RBAC, comments, RSS, scheduling, related posts
- Private docs: path-based access rules, role requirements
- Image uploads: local filesystem and S3 providers
- Background scheduler for scheduled post publishing
- PDF export for documentation pages
- Analytics plugin system (GA4, Plausible, Umami, Matomo, MinimalDoc)

## [1.2.0] - 2025-10-01

### Added

- Multi-version documentation with version selector
- Internationalization (i18n) with directory-based structure
- RTL language support

## [1.1.0] - 2025-08-01

### Added

- Knowledge base page type
- Landing page builder (YAML + Markdown)
- Footer version display

## [1.0.0] - 2025-06-01

### Added

- Static site generator for documentation
- FAQ, legal pages, knowledge base page types
- Broken link checker
- Dark/light theme support
- Search (client-side JSON index)
- Responsive design
- CLI tool (`minimaldoc build`)

[Unreleased]: https://github.com/studiowebux/minimaldoc/compare/1.4.0...HEAD
[1.4.0]: https://github.com/studiowebux/minimaldoc/compare/1.3.0...1.4.0
[1.3.0]: https://github.com/studiowebux/minimaldoc/compare/1.2.0...1.3.0
[1.2.0]: https://github.com/studiowebux/minimaldoc/compare/1.1.0...1.2.0
[1.1.0]: https://github.com/studiowebux/minimaldoc/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/studiowebux/minimaldoc/releases/tag/1.0.0
