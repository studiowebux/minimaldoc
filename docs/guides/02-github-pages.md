---
title: GitHub Pages
description: Deploy documentation to GitHub Pages
tags:
  - guides
  - github-pages
  - deployment
---

# GitHub Pages

Deploy MinimalDoc to GitHub Pages with automated builds.

## Repository Setup

### Option 1: Docs in Main Repo

```
my-project/
├── src/
├── docs/           # Documentation source
├── README.md
└── .github/
    └── workflows/
        └── docs.yml
```

### Option 2: Dedicated Docs Repo

```
my-project-docs/
├── docs/           # Documentation source
└── .github/
    └── workflows/
        └── docs.yml
```

## GitHub Actions Workflow

Create `.github/workflows/docs.yml`:

```yaml
name: Deploy Documentation

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/workflows/docs.yml'
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Install MinimalDoc
        run: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest

      - name: Build Documentation
        run: |
          minimaldoc build ./docs \
            --base-url "https://${{ github.repository_owner }}.github.io/${{ github.event.repository.name }}" \
            --output dist

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: dist

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Enable GitHub Pages

1. Go to repository Settings
2. Navigate to Pages
3. Source: GitHub Actions

## Custom Domain

### Configure Domain

1. Add `CNAME` file to docs output:

```yaml
- name: Build Documentation
  run: |
    minimaldoc build ./docs --output dist
    echo "docs.example.com" > dist/CNAME
```

2. Or add CNAME in docs folder (gets copied):

```
docs/
├── CNAME              # Contains: docs.example.com
├── config.yaml
└── index.md
```

### DNS Configuration

Add DNS records:

**Apex domain (example.com):**

```
A     @     185.199.108.153
A     @     185.199.109.153
A     @     185.199.110.153
A     @     185.199.111.153
```

**Subdomain (docs.example.com):**

```
CNAME docs  username.github.io
```

### Update Base URL

```yaml
# config.yaml
base_url: https://docs.example.com
```

## Clean URLs

For clean URLs on GitHub Pages:

```yaml
clean_urls: true
```

Add 404 handling in workflow:

```yaml
- name: Build Documentation
  run: |
    minimaldoc build ./docs --output dist --clean-urls
    cp dist/index.html dist/404.html
```

## Organization Sites

For `username.github.io` repository:

```yaml
- name: Build Documentation
  run: |
    minimaldoc build ./docs \
      --base-url "https://username.github.io" \
      --output dist
```

## Project Sites

For `username.github.io/project-name`:

```yaml
- name: Build Documentation
  run: |
    minimaldoc build ./docs \
      --base-url "https://username.github.io/project-name" \
      --output dist
```

## Multiple Documentation Sites

### Monorepo

```yaml
jobs:
  build-docs:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        project: [api, sdk, cli]
    steps:
      - uses: actions/checkout@v4
      - name: Build ${{ matrix.project }}
        run: |
          minimaldoc build ./packages/${{ matrix.project }}/docs \
            --output dist/${{ matrix.project }}
```

## Caching

Speed up builds with caching:

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.24'
    cache: true

- name: Cache MinimalDoc
  uses: actions/cache@v4
  with:
    path: ~/go/bin/minimaldoc
    key: minimaldoc-${{ runner.os }}
```

## Branch Protection

Protect main branch:

1. Settings > Branches > Add rule
2. Require status checks
3. Add "build" as required check

## Manual Deployment

Deploy on-demand:

```yaml
on:
  workflow_dispatch:
    inputs:
      environment:
        description: 'Deployment environment'
        required: true
        default: 'production'
        type: choice
        options:
          - production
          - staging
```

## Preview Deployments

Deploy PRs to preview:

```yaml
name: Preview Documentation

on:
  pull_request:
    paths:
      - 'docs/**'

jobs:
  preview:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build Preview
        run: |
          minimaldoc build ./docs \
            --base-url "https://preview-${{ github.event.number }}.example.com" \
            --output dist

      - name: Deploy Preview
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./dist
          destination_dir: pr-${{ github.event.number }}
```

## Troubleshooting

### 404 Errors

- Check `base_url` matches deployment URL
- Verify paths in TOC.md
- Check file extensions (.html vs clean URLs)

### Build Failures

- Ensure Go version compatibility
- Check docs directory path
- Verify config.yaml syntax

### Deployment Pending

- Check Actions tab for errors
- Verify Pages is enabled
- Check permissions in workflow
