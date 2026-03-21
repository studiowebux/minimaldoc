# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.8.0] - 2026-03-21

### Added

- 7 new starter templates: portfolio-item, faq-item, kb-article, legal-page,
  openapi-page, version-page, changelog-release
- `featurePath()` and `trimBaseURL()` helpers — replace 32 duplicate patterns
  across 12 generator files
- `nav_depth` is now configurable via config.yaml (was previously undocumented)
- `stale_warning.message` field for custom stale content warning text
- Tests for builder (navigation tree, page discovery, output paths, prev/next)
  and cli (output directory safety checks) — 8/10 packages now have tests
- Version injection via ldflags — Makefile and publish workflow inject version
  from git tag automatically, no manual bump needed
- Auto-deploy docs to GitHub Pages on push to main (path-filtered)
- Detailed documentation for version frontmatter fields with directory structure

### Changed

- **Repository restructured**: all server code, admin UI, Docker configs, and
  deployment files removed. This project now focuses exclusively on the static
  site generator. Backend extracted to `studiowebux/minimaldoc-backend` (private).
- Go module bumped to 1.26.0; CI workflows updated to Go 1.26
- CI actions bumped: upload-artifact v4→v7 (Node 24), upload-pages-artifact v3→v4
- `internal/version/version.go` changed from `const` to `var` with "dev" default
- `templates/config.yaml` expanded from 13 to 34 configuration sections
- Documentation URL updated to minimaldoc.com

### Fixed

- 55+ security fixes from full codebase audit (XSS, SSRF, auth bypass, CSRF,
  input validation, rate limiting, race conditions)
- SSRF protection in OpenAPI spec URL fetching (private IP filtering)
- XSS: analytics and OpenAPI HTML generation now use proper escaping
- Race condition in link checker heading cache (added sync.RWMutex)
- Open redirect prevention in search.js navigation
- Config validation: path traversal checks, URL scheme validation, email format
- TOC parser: path traversal protection, auto-detect indent unit
- search.js manifest polling leak: added 10s timeout to setInterval
- gofmt formatting issues across multiple files

### Removed

- **Backend server** (`minimaldoc-server`) — extracted to a separate private
  repository after a full security audit revealed 52 issues
- Dead features: waitlist, feedback widget, newsletter form, self-hosted
  analytics (all required backend server)
- Server-dependent code: Docker configs, nginx configs, deployment files,
  admin UI, all `internal/server/` packages

## [1.4.2] - 2026-03-12

### Fixed

- Address all excluded gosec rules: G104 (unhandled errors), G107 (URL
  validation before HTTP fetch), G301/G306 (file/directory permissions),
  G304 (file inclusion via variable)
- `internal/generator/fs.go`: centralise web output file/dir helpers with
  explicit nosec annotations for G301/G306
- gosec CI exclusion list reduced from 8 rules to 3 (G101, G115, G203)

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

[Unreleased]: https://github.com/studiowebux/minimaldoc/compare/v1.8.0...HEAD
[1.8.0]: https://github.com/studiowebux/minimaldoc/compare/v1.7.0...v1.8.0
[1.4.0]: https://github.com/studiowebux/minimaldoc/compare/1.3.0...1.4.0
[1.3.0]: https://github.com/studiowebux/minimaldoc/compare/1.2.0...1.3.0
[1.2.0]: https://github.com/studiowebux/minimaldoc/compare/1.1.0...1.2.0
[1.1.0]: https://github.com/studiowebux/minimaldoc/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/studiowebux/minimaldoc/releases/tag/1.0.0
