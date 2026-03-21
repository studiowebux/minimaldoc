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
| `changelog-release.md` | Individual release | `docs/__changelog__/releases/vX.Y.Z.md` |
| `toc.md` | Custom navigation | `docs/TOC.md` |
| `portfolio-item.md` | Portfolio project | `docs/__portfolio__/project-name.md` |
| `faq-item.md` | FAQ question/answer | `docs/__faq__/category/question.md` |
| `kb-article.md` | Knowledge base article | `docs/__kb__/category/article.md` |
| `legal-page.md` | Legal page | `docs/__legal__/page-name.md` |
| `openapi-page.md` | OpenAPI endpoint page | `docs/api/endpoint.md` |
| `version-page.md` | Versioned documentation | `docs/__versions__/v2/page.md` |
| `components.yaml` | Status page components | `docs/__status__/components.yaml` |
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
│   ├── openapi.yaml      # OpenAPI specs (if enabled)
│   └── endpoint.md       # Per-endpoint pages (optional)
├── __portfolio__/        # Portfolio (if enabled)
│   └── project-name.md
├── __faq__/              # FAQ (if enabled)
│   └── category-name/
│       └── question.md
├── __kb__/               # Knowledge base (if enabled)
│   └── category-name/
│       └── article.md
├── __legal__/            # Legal pages (if enabled)
│   └── privacy-policy.md
├── __changelog__/        # Per-release changelog (if enabled)
│   └── releases/
│       └── vX.Y.Z.md
├── __versions__/         # Versioned docs (if enabled)
│   ├── v1/
│   │   └── *.md
│   └── v2/
│       └── *.md
└── __status__/           # Status page (if enabled)
    ├── config.yaml
    ├── components.yaml
    ├── incidents/
    │   └── YYYY-MM-DD-*.md
    └── maintenance/
        └── YYYY-MM-DD-*.md
```
