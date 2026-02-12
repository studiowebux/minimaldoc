# MinimalDoc Bootstrap

You are a documentation site architect. Your task is to take any repository and create a complete MinimalDoc website from scratch.

---

## Phase 1: Repository Analysis

### Step 1.1: Read Core Documentation

Read these files if they exist:

- `README.md` - Project overview, installation, basic usage
- `CONTRIBUTING.md` - Development setup, contribution guidelines
- `CHANGELOG.md` - Version history, features added over time
- `LICENSE` - License type for footer attribution

### Step 1.2: Identify Project Type

Determine what kind of project this is:

| Type | Indicators | Documentation Focus |
|------|------------|---------------------|
| Library/SDK | `package.json`, `go.mod`, `Cargo.toml`, exports | API reference, integration guides |
| CLI Tool | `main.go`, `bin/`, argument parsing | Commands, flags, usage examples |
| Web Service | `Dockerfile`, routes, endpoints | API docs, deployment, configuration |
| Framework | Plugin system, lifecycle hooks | Getting started, concepts, extensions |
| Monorepo | `packages/`, `apps/`, workspaces | Per-package docs, architecture overview |

### Step 1.3: Extract Feature List

Scan the codebase for:

- Exported functions/classes/types
- Configuration options
- CLI commands and flags
- API endpoints
- Plugin hooks
- Environment variables

Create a feature inventory for documentation planning.

---

## Phase 2: Site Structure

### Step 2.1: Create Directory Structure

```bash
mkdir -p docs
```

### Step 2.2: Generate config.yaml

Create `docs/config.yaml` with appropriate settings:

```yaml
title: "[Project Name]"
description: "[From README first paragraph]"
base_url: ""
author: "[From package.json/go.mod or git config]"
theme: default
dark_mode: false
enable_llms: true
enable_search: true
clean_urls: false

# Enable features based on project type
landing:
  enabled: true
  nav:
    - text: "Docs"
      url: "/getting-started/installation.html"
    - text: "GitHub"
      url: "[repo-url]"
  hero:
    title: "[Project tagline or name]"
    subtitle: "[One-sentence value proposition]"
    buttons:
      - text: "Get Started"
        url: "/getting-started/installation.html"
        primary: true
      - text: "View on GitHub"
        url: "[repo-url]"
  features:
    title: "Features"
    items: []  # Populated from feature inventory

footer:
  copyright: "[Year] [Author]. [License]."
  links:
    - title: "Documentation"
      items:
        - text: "Getting Started"
          url: "/getting-started/installation.html"
    - title: "Resources"
      items:
        - text: "GitHub"
          url: "[repo-url]"
```

### Step 2.3: Create Navigation Structure

Plan documentation sections based on project type:

**Library/SDK:**
```
docs/
  getting-started/
    01-installation.md
    02-quickstart.md
  api-reference/
    01-overview.md
    [one file per module/package]
  guides/
    [use-case guides]
```

**CLI Tool:**
```
docs/
  getting-started/
    01-installation.md
    02-basic-usage.md
  commands/
    [one file per command]
  configuration/
    01-config-file.md
    02-environment-variables.md
```

**Web Service:**
```
docs/
  getting-started/
    01-installation.md
    02-deployment.md
  api/
    01-authentication.md
    02-endpoints.md
  configuration/
    01-environment.md
```

---

## Phase 3: Content Generation

### Step 3.1: Create Index Page

Create `docs/index.md`:

```markdown
---
title: "[Project Name]"
description: "[Project description]"
tags:
  - home
  - overview
---

# [Project Name]

[Expanded description from README]

## Quick Start

[Minimal example to get started]

## Key Features

[Feature highlights with links to detailed docs]

## Documentation

- [Installation](getting-started/installation.html)
- [Quick Start](getting-started/quickstart.html)
- [API Reference](api-reference/overview.html)
```

### Step 3.2: Generate Section Pages

For each planned section, create pages following MinimalDoc format:

**Frontmatter Template:**
```yaml
---
title: "[Clear, descriptive title]"
description: "[1-2 sentence summary]"
tags:
  - [relevant]
  - [keywords]
---
```

**Content Guidelines:**

1. Start with H1 matching the title
2. Brief introduction paragraph
3. Use H2 for major sections
4. Include code examples with proper language tags
5. Use admonitions for warnings/tips
6. Link to related pages

### Step 3.3: API Reference Generation

For libraries/SDKs, document each exported item:

```markdown
---
title: "[Function/Class Name]"
description: "[Brief description]"
tags:
  - api
  - [category]
---

# [Function/Class Name]

[Description of what it does]

## Signature

```[language]
[function signature or class definition]
```

## Parameters

| Name | Type | Description |
|------|------|-------------|
| `param1` | `string` | Description |

## Returns

[Return value description]

## Example

```[language]
[usage example]
```

## Notes

:::tip
[Usage tips]
:::
```

### Step 3.4: Command Documentation

For CLI tools, document each command:

```markdown
---
title: "[command name]"
description: "[What the command does]"
tags:
  - cli
  - command
---

# [command name]

[Description]

## Usage

```bash
[tool] [command] [options]
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `--flag` | Description | `value` |

## Examples

### Basic Usage

```bash
[example command]
```

### Advanced Usage

```bash
[complex example]
```
```

---

## Phase 4: Verification

### Step 4.1: Build Test

```bash
minimaldoc build docs -o public
```

Check for:
- Build completes without errors
- All pages render correctly
- Navigation works
- Search index generates

### Step 4.2: Link Validation

```bash
minimaldoc build docs -o public --link-check=error
```

Fix any broken internal links.

### Step 4.3: Content Review

Verify:
- [ ] All frontmatter has title and description
- [ ] Code examples are complete and runnable
- [ ] Internal links use correct paths
- [ ] Admonitions are used for important notes
- [ ] Landing page features match actual capabilities

---

## Quality Checklist

Before considering the site complete:

- [ ] `config.yaml` has all required fields
- [ ] Landing page hero reflects project purpose
- [ ] Getting started section exists and is complete
- [ ] All major features are documented
- [ ] Code examples work
- [ ] Navigation is logical
- [ ] Search returns relevant results
- [ ] Build completes without warnings
- [ ] Links are valid

---

## Output

When complete, provide:

1. Summary of created structure
2. List of pages generated
3. Any manual follow-up needed (missing information, incomplete sections)
4. Build command to verify the site

---

## Usage

To use this workflow:

```
/bootstrap
```

Then provide the repository path or current directory context.
