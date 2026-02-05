# MinimalDoc Templates

Starter templates for common documentation patterns.

## Usage

Copy templates to your docs directory and customize as needed.

## Templates

| File | Description | Location |
|------|-------------|----------|
| `config.yaml` | Site configuration | `docs/config.yaml` |
| `page.md` | Documentation page | `docs/any-page.md` |
| `changelog.md` | Keep a Changelog format | `docs/changelog.md` |
| `toc.md` | Custom navigation | `docs/TOC.md` |
| `components.yaml` | Status page components | `docs/__status__/components.yaml` |
| `status-config.yaml` | Status page config | `docs/__status__/config.yaml` |
| `incident.md` | Incident report | `docs/__status__/incidents/YYYY-MM-DD-title.md` |
| `maintenance.md` | Scheduled maintenance | `docs/__status__/maintenance/YYYY-MM-DD-title.md` |

## Directory Structure

```
docs/
├── config.yaml           # Site configuration
├── TOC.md                # Custom navigation (optional)
├── index.md              # Homepage
├── changelog.md          # Project changelog
├── getting-started/
│   └── *.md
├── features/
│   └── *.md
├── api/
│   └── openapi.yaml      # OpenAPI specs (if enabled)
└── __status__/           # Status page (if enabled)
    ├── config.yaml
    ├── components.yaml
    ├── incidents/
    │   └── YYYY-MM-DD-*.md
    └── maintenance/
        └── YYYY-MM-DD-*.md
```
