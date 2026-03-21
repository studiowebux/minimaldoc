---
title: Multi-Version Documentation
description: Maintain documentation for multiple software versions with version-specific content and overrides
menu_order: 9
since: "v1.2"
versions:
  - v1.2
---

MinimalDoc supports maintaining documentation for multiple versions of your software. This feature allows users to switch between documentation versions while sharing common content across versions.

## Configuration

Enable multi-version documentation in your `config.yaml`:

```yaml
versions:
  enabled: true
  default: "v2"
  list:
    - name: "v2"
      label: "2.x (Latest)"
      path: "v2"
    - name: "v1"
      label: "1.x (LTS)"
      path: "v1"
      eol: "2025-12-31"
  selector:
    position: "header"
    show_eol_warning: true
```

### Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `enabled` | boolean | Enable multi-version documentation |
| `default` | string | Default version to display (used for root URLs) |
| `list` | array | List of version definitions |
| `selector.position` | string | Where to show version selector: `header` or `sidebar` |
| `selector.show_eol_warning` | boolean | Show warning banner for EOL versions |

### Version Definition

Each version in the `list` array supports:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Version identifier (e.g., `v2`) |
| `label` | string | Display label (e.g., `2.x (Latest)`) |
| `path` | string | URL path prefix for this version |
| `eol` | string | End of life date in `YYYY-MM-DD` format (optional) |

## Directory Structure

### Shared Content

By default, all pages in your `docs/` directory are shared across all versions. Use frontmatter to control which versions a page appears in.

### Version-Specific Overrides

Create version-specific content in the `__versions__/` directory:

```
docs/
  getting-started/
    install.md              # Shared across all versions
    config.md               # Shared across all versions
  __versions__/
    v1/
      getting-started/
        config.md           # Overrides config.md for v1
      migration/
        upgrade-guide.md    # Only appears in v1
    v2/
      api/
        breaking-changes.md # Only appears in v2
```

Pages in `__versions__/{version}/` override or extend the shared content for that specific version.

## Frontmatter Fields

Control version behavior with frontmatter:

```yaml
---
title: Configuration Guide
versions:
  - v1
  - v2
since: "v1.5"
deprecated_in: "v3"
removed_in: "v4"
version_note: "v2"
---
```

### Available Fields

| Field | Type | Description |
|-------|------|-------------|
| `versions` | array | Versions this page appears in (empty = all versions) |
| `since` | string | Version where feature was introduced |
| `deprecated_in` | string | Version where feature was deprecated |
| `removed_in` | string | Version where feature was removed (page hidden from this version onward) |
| `version_note` | string | Version to show special note for |

## URL Structure

Multi-version documentation uses URL prefixes:

```
/docs/getting-started/    # Default version (v2)
/v1/docs/getting-started/ # Explicit v1
/v2/docs/getting-started/ # Explicit v2
```

The default version is served from the root URL. Other versions are prefixed with their path.

## Search Integration

Each version gets its own search index:

- `/search-index.json` - Default version
- `/v1/search-index.json` - v1 specific
- `/v2/search-index.json` - v2 specific

Search results are scoped to the currently viewed version.

## Version Selector

A dropdown selector appears in the sidebar allowing users to switch between versions:

- Shows the current version label
- Lists all available versions
- Marks EOL versions with a badge
- Persists selection in localStorage

When switching versions, users are redirected to the same page in the new version (if it exists) or to the version's root page.

## EOL Warnings

For versions that have reached end of life (`eol` date set), MinimalDoc can display:

1. **Badge**: "EOL" badge next to the version in the selector
2. **Warning Banner**: Full-width warning at the top of pages in EOL versions

Enable/disable with `selector.show_eol_warning`.

## Build Output

When versioning is enabled, the build generates:

```
public/
  docs/                   # Default version
  v1/
    docs/
    search-index.json
  v2/
    docs/
    search-index.json
  versions.json           # Version metadata for JS
```

The `versions.json` file contains metadata for client-side version switching.
