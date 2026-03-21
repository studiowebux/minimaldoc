---
title: Changelog Format
description: Release file format, categories, and semver
tags:
  - changelog
  - format
  - semver
---

# Changelog Format

Release files use Markdown with YAML frontmatter.

## File Structure

```
__changelog__/
└── releases/
    ├── 2.0.0.md
    ├── 1.1.0.md
    ├── 1.0.0.md
    └── 0.1.0.md
```

Filename: `{version}.md`

## Basic Format

```markdown
---
version: 1.1.0
date: 2025-01-28
---

## Added

- New authentication system
- Dark mode support
- Export to PDF feature

## Changed

- Improved search performance
- Updated dependencies

## Fixed

- Fixed login timeout issue
- Corrected date formatting
```

## Frontmatter Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | string | Yes | Semver version |
| `date` | string | Yes | Release date (YYYY-MM-DD) |
| `title` | string | No | Release title |
| `prerelease` | bool | No | Mark as prerelease |

## Categories

Standard categories from Keep a Changelog:

| Category | Description | Badge Color |
|----------|-------------|-------------|
| `Added` | New features | Green |
| `Changed` | Changes in existing functionality | Blue |
| `Deprecated` | Features to be removed | Yellow |
| `Removed` | Removed features | Red |
| `Fixed` | Bug fixes | Purple |
| `Security` | Security fixes | Orange |

### Category Headings

Use H2 headings:

```markdown
## Added
- Feature one
- Feature two

## Fixed
- Bug fix one
```

### Entry Format

Use list items:

```markdown
## Added

- Brief description of feature
- Another feature with more detail that explains
  what changed and why
- Feature with [link to issue](https://github.com/org/repo/issues/123)
```

## Semver

Versions follow [Semantic Versioning](https://semver.org):

```
MAJOR.MINOR.PATCH
```

| Component | When to Increment |
|-----------|-------------------|
| MAJOR | Breaking changes |
| MINOR | New features (backward compatible) |
| PATCH | Bug fixes (backward compatible) |

### Examples

```
1.0.0    # Initial release
1.0.1    # Bug fix
1.1.0    # New feature
2.0.0    # Breaking change
```

### Prereleases

```
2.0.0-alpha.1
2.0.0-beta.1
2.0.0-rc.1
```

Mark in frontmatter:

```yaml
---
version: 2.0.0-beta.1
date: 2025-01-20
prerelease: true
---
```

## Sorting

Releases sorted by semver:

```
2.0.0         # First (newest)
1.1.0
1.0.1
1.0.0
0.2.0
0.1.0         # Last (oldest)
```

Prereleases sort before their release:

```
2.0.0
2.0.0-rc.1
2.0.0-beta.1
2.0.0-alpha.1
1.0.0
```

## Examples

### Initial Release

```markdown
---
version: 1.0.0
date: 2025-01-15
title: Initial Release
---

## Added

- Core documentation features
- Markdown support with GFM
- Automatic navigation
- Dark mode support
- Search functionality
- Syntax highlighting
```

### Feature Release

```markdown
---
version: 1.1.0
date: 2025-01-28
---

## Added

- OpenAPI/Swagger documentation support
- Status page feature
- Social links in sidebar

## Changed

- Improved search performance
- Updated navigation styling

## Fixed

- Fixed mobile menu toggle
- Corrected sitemap generation
```

### Bug Fix Release

```markdown
---
version: 1.0.1
date: 2025-01-20
---

## Fixed

- Fixed syntax highlighting for YAML
- Corrected link handling in TOC
- Fixed dark mode persistence
```

### Breaking Change

```markdown
---
version: 2.0.0
date: 2025-02-15
title: Configuration Overhaul
---

## Changed

- **BREAKING**: Renamed `config.yml` to `config.yaml`
- **BREAKING**: Changed theme structure

## Added

- New configuration options
- Plugin system

## Removed

- **BREAKING**: Removed legacy theme format

## Migration

See [migration guide](/docs/migration-v2.html) for upgrade instructions.
```

### Security Release

```markdown
---
version: 1.0.2
date: 2025-01-22
---

## Security

- Fixed XSS vulnerability in search input
- Updated dependencies with security patches

## Fixed

- Minor UI fixes
```

### Prerelease

```markdown
---
version: 2.0.0-beta.1
date: 2025-02-01
prerelease: true
---

## Added

- Preview of new plugin system
- Experimental features

## Changed

- Updated internal architecture

## Known Issues

- Plugin hot-reload not working
- Some themes may not render correctly
```

## Compare Links

With repository configured:

```yaml
changelog:
  repository: "https://github.com/org/repo"
```

Links generated:

| Version | Compare Link |
|---------|--------------|
| 2.0.0 | `/compare/v1.1.0...v2.0.0` |
| 1.1.0 | `/compare/v1.0.0...v1.1.0` |
| 1.0.0 | `/releases/tag/v1.0.0` |

Assumes tags follow `v{version}` format.

