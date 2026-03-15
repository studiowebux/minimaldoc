---
title: Portfolio
description: Project showcase with filterable tags and individual project pages
tags:
  - guides
  - portfolio
---

# Portfolio

Display a project showcase with filterable tags, images, and individual project pages.

## Configuration

```yaml
portfolio:
  enabled: true
  title: Projects
  description: Sites built with our tool
  path: projects
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable portfolio page |
| `title` | string | `"Projects"` | Page title |
| `description` | string | `""` | Page description |
| `path` | string | `"projects"` | Output path |

## Project Files

Create markdown files in `__portfolio__/` directory:

```
docs/
  __portfolio__/
    my-project.md
    another-project.md
```

Each file uses frontmatter to define the project:

```yaml
---
title: My Project
description: A brief description of the project
image: /images/projects/my-project.png
url: https://example.com
tags:
  - documentation
  - open-source
order: 1
---

Detailed project description in markdown.
```

### Frontmatter Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | yes | Project name |
| `description` | string | no | Short description shown on the listing |
| `image` | string | no | Preview image URL |
| `url` | string | no | External project link |
| `tags` | []string | no | Filterable tags |
| `order` | int | no | Sort order (lower = first) |

## Tag Filtering

Tags are rendered as clickable filters on the portfolio listing page. Clicking a tag shows only projects with that tag.

## Output

- Listing page: `/<path>/index.html`
- Individual project: `/<path>/<slug>/index.html`
