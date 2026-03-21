package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/studiowebux/minimaldoc/internal/core"
)

var (
	withTOC       bool
	withStatus    bool
	withChangelog bool
	withOpenAPI   bool
	fullInit      bool
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new documentation site",
	Long: `Initialize creates a new documentation site with example files.

This will create:
- config.yaml with site configuration
- Sample markdown files with frontmatter
- A basic directory structure

Optional features (use flags to enable):
- --with-toc: Custom navigation (TOC.md)
- --with-status: Status page structure
- --with-changelog: Changelog structure
- --with-openapi: OpenAPI specification example
- --full: Include all optional features`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	InitCmd.Flags().BoolVar(&withTOC, "with-toc", false, "Create TOC.md for custom navigation")
	InitCmd.Flags().BoolVar(&withStatus, "with-status", false, "Create status page structure")
	InitCmd.Flags().BoolVar(&withChangelog, "with-changelog", false, "Create changelog structure")
	InitCmd.Flags().BoolVar(&withOpenAPI, "with-openapi", false, "Create OpenAPI specification example")
	InitCmd.Flags().BoolVar(&fullInit, "full", false, "Create all optional features")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Handle --full flag
	if fullInit {
		withTOC = true
		withStatus = true
		withChangelog = true
		withOpenAPI = true
	}

	// Determine target directory
	targetDir := "docs"
	if len(args) > 0 {
		targetDir = args[0]
	}

	fmt.Println("Initializing new documentation site...")

	// Create base directory structure
	dirs := []string{
		targetDir,
		filepath.Join(targetDir, "getting-started"),
		filepath.Join(targetDir, "guides"),
	}

	// Add optional directories
	if withOpenAPI {
		dirs = append(dirs, filepath.Join(targetDir, "api"))
	}
	if withStatus {
		dirs = append(dirs,
			filepath.Join(targetDir, core.StatusSourceDir),
			filepath.Join(targetDir, core.StatusSourceDir, "incidents"),
			filepath.Join(targetDir, core.StatusSourceDir, "maintenance"),
		)
	}
	if withChangelog {
		dirs = append(dirs,
			filepath.Join(targetDir, core.ChangelogSourceDir),
			filepath.Join(targetDir, core.ChangelogSourceDir, "releases"),
		)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil { // #nosec G301 -- project source directories need broad access
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Build config content based on flags
	finalConfig := configContent
	if withOpenAPI {
		finalConfig += openAPIConfigSection
	}
	if withStatus {
		finalConfig += statusConfigSection
	}
	if withChangelog {
		finalConfig += changelogConfigSection
	}

	// Create core files
	files := map[string]string{
		filepath.Join(targetDir, "config.yaml"):                        finalConfig,
		filepath.Join(targetDir, "index.md"):                           indexContent,
		filepath.Join(targetDir, "getting-started", "installation.md"): installContent,
		filepath.Join(targetDir, "getting-started", "quick-start.md"):  quickstartContent,
		filepath.Join(targetDir, "guides", "deployment.md"):            deploymentContent,
	}

	// Add optional files
	if withTOC {
		files[filepath.Join(targetDir, "TOC.md")] = tocContent
	}

	if withOpenAPI {
		files[filepath.Join(targetDir, "api", "openapi.yaml")] = openapiContent
	}

	if withStatus {
		files[filepath.Join(targetDir, core.StatusSourceDir, "components.yaml")] = componentsContent
		files[filepath.Join(targetDir, core.StatusSourceDir, "incidents", incidentFilename())] = incidentContent
	}

	if withChangelog {
		files[filepath.Join(targetDir, core.ChangelogSourceDir, "releases", "0.1.0.md")] = releaseContent
	}

	for path, content := range files {
		if _, err := os.Stat(path); err == nil {
			fmt.Fprintf(os.Stderr, "Warning: %s already exists, skipping\n", path)
			continue
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil { // #nosec G306 -- documentation source files are not sensitive
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
	}

	fmt.Println()
	fmt.Println("Documentation site initialized!")
	fmt.Printf("Location: %s\n", targetDir)
	fmt.Println()
	fmt.Println("Created files:")
	fmt.Println("  - config.yaml (site configuration)")
	fmt.Println("  - index.md (homepage)")
	fmt.Println("  - getting-started/ (starter pages)")
	fmt.Println("  - guides/ (guide pages)")
	if withTOC {
		fmt.Println("  - TOC.md (custom navigation)")
	}
	if withOpenAPI {
		fmt.Println("  - api/openapi.yaml (OpenAPI specification)")
	}
	if withStatus {
		fmt.Printf("  - %s/ (status page)\n", core.StatusSourceDir)
	}
	if withChangelog {
		fmt.Printf("  - %s/ (changelog)\n", core.ChangelogSourceDir)
	}
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit config.yaml to customize your site")
	fmt.Println("  2. Edit the markdown files")
	fmt.Println("  3. Run 'minimaldoc build " + targetDir + "' to generate your site")
	fmt.Println("  4. Open public/index.html in your browser")
	fmt.Println()

	return nil
}

func incidentFilename() string {
	return time.Now().AddDate(0, 0, -7).Format("2006-01-02") + "-example-incident.md"
}

// Config template (base)
// Every supported key is present. Core fields are active; optional sections
// are commented out so the user can see what exists and uncomment as needed.
var configContent = `# ──────────────────────────────────────────────────────────────
# MinimalDoc Configuration
# Docs: https://minimaldoc.dev/getting-started/configuration.html
# ──────────────────────────────────────────────────────────────

# ── Site Information ──────────────────────────────────────────
title: My Documentation
description: Documentation for my project
author: Your Name

# IMPORTANT: Set base_url before deploying. It controls how asset
# and navigation URLs are built. Leave empty for local preview.
#   Root domain  → https://docs.example.com
#   Subdirectory → https://example.com/docs
base_url: ""

# ── Theme & Appearance ───────────────────────────────────────
theme: default                    # default | yellow
dark_mode: false                  # Start in dark mode

# Custom theme colors and fonts (uncomment to override)
# theme_config:
#   colors:
#     light:
#       bg_primary: "#fafafa"
#       text_primary: "#1a1a1a"
#       link_color: "#2563eb"
#       accent_primary: "#2a2a2a"
#     dark:
#       bg_primary: "#1a1a1a"
#       text_primary: "#ffffff"
#       link_color: "#7bb3ff"
#       accent_primary: "#e5e5e5"
#   fonts:
#     heading: "Inter, sans-serif"
#     body: "Inter, sans-serif"
#     code: "Fira Code, monospace"
#     google_url: "https://fonts.googleapis.com/css2?family=Inter&display=swap"
#   hero:
#     background_image: ""
#     background_overlay: "rgba(0,0,0,0.6)"
#     text_align: center            # left | center | right
#     min_height: ""                # e.g. "80vh"

# ── Features ─────────────────────────────────────────────────
enable_llms: true                 # Generate llms.txt for AI tools
enable_search: true               # Client-side search (Cmd+K / Ctrl+K)
clean_urls: false                 # /page/ instead of /page.html
# entrypoint: ""                  # Custom homepage file (default: index.md)
# nav_depth: 0                    # Max navigation tree depth (0 = unlimited)

# ── Stale Content Warnings ───────────────────────────────────
# Show a banner on pages that haven't been updated in a while.
stale_warning:
  enabled: false
  threshold_days: 365
  show_update_date: true
  # message: ""                   # Custom warning text

# ── Social Links ─────────────────────────────────────────────
# Shown in the sidebar. Icons: github, twitter, linkedin, youtube,
# discord, mastodon, rss, email, website
social_links:
  - name: GitHub
    url: https://github.com/your-org/your-project
    icon: github

# ── OpenAPI / Swagger ────────────────────────────────────────
# Interactive API docs from OpenAPI 3.x specs.
# openapi:
#   enabled: true
#   spec_files:                   # Local spec files (relative to docs dir)
#     - "api/openapi.yaml"
#   # spec_urls:                  # Remote specs (fetched on build)
#   #   - "https://api.example.com/openapi.json"
#   default_view: "path"          # path | tag | flat
#   # sync_on_build: false        # Re-fetch remote specs every build
#   # cache_dir: ".openapi-cache" # Cache directory for remote specs
#   enable_testing: true          # Show "Try It" request tester
#   enable_export: true           # cURL / restcli export buttons
#   enable_code_samples: true     # Code samples panel
#   # lazy_load_chunk_size: 0     # Split large specs (0 = no split)

# ── MCP Server Documentation ────────────────────────────────
# Generate docs from MCP server JSON manifests.
# mcp:
#   enabled: true
#   spec_files:
#     - "mcp/server-manifest.json"
#   path: "mcp"                   # Output URL path

# ── Status Page ──────────────────────────────────────────────
# Service status with components, incidents, and maintenance windows.
# Reads from __status__/ directory.
# status:
#   enabled: true
#   title: "Service Status"
#   description: "Current operational status"
#   path: "status"
#   show_history: true
#   history_months: 12
#   rss_enabled: true

# ── Changelog ────────────────────────────────────────────────
# Release notes. Reads from __changelog__/releases/ directory.
# changelog:
#   enabled: true
#   title: "Changelog"
#   description: "All notable changes to this project"
#   path: "changelog"
#   rss_enabled: true
#   # repository: "https://github.com/your-org/your-project"

# ── Roadmap ──────────────────────────────────────────────────
# Public roadmap board or timeline.
# roadmap:
#   enabled: true
#   title: "Roadmap"
#   description: "What we're building and where we're headed."
#   path: "roadmap"
#   layout: "board"               # board | timeline
#   show_versions: true
#   columns:
#     - id: planned
#       label: Planned
#     - id: in-progress
#       label: In Progress
#     - id: shipped
#       label: Shipped
#   items:
#     - title: "Feature X"
#       description: "Short description"
#       status: planned            # Must match a column id
#       version: "1.0"
#       tags: [backend]

# ── Landing Page ─────────────────────────────────────────────
# Full-featured landing page with hero, features, CTA, etc.
# Replaces index.md as the homepage when enabled.
# landing:
#   enabled: true
#   nav:
#     - text: Docs
#       url: /getting-started/installation.html
#     - text: Features
#       url: /features/overview.html
#   hero:
#     title: "Build Docs Fast"
#     subtitle: "A minimal static site generator for documentation."
#     buttons:
#       - text: Get Started
#         url: /getting-started/installation.html
#         primary: true
#       - text: View on GitHub
#         url: https://github.com/your-org/your-project
#   features:
#     title: Features
#     items:
#       - title: Fast
#         description: "Builds in milliseconds"
#         icon: "~"
#   # steps, cta, testimonials, opensource, links sections also available

# ── Portfolio ────────────────────────────────────────────────
# Project showcase page. Reads markdown from __portfolio__/ directory.
# portfolio:
#   enabled: true
#   title: "Portfolio"
#   description: "Projects and experiments"
#   path: "portfolio"

# ── FAQ ──────────────────────────────────────────────────────
# Accordion-style FAQ. Reads markdown from __faq__/ directory or
# define inline categories below.
# faq:
#   enabled: true
#   title: "FAQ"
#   description: "Frequently asked questions"
#   path: "faq"
#   # categories:
#   #   - name: General
#   #     items:
#   #       - question: "What is this?"
#   #         answer: "A documentation tool."

# ── Contact ──────────────────────────────────────────────────
# Simple contact page with email and info items.
# contact:
#   enabled: true
#   title: "Contact"
#   description: "Get in touch"
#   path: "contact"
#   email: "hello@example.com"
#   info:
#     - icon: email
#       text: "hello@example.com"

# ── Knowledge Base ───────────────────────────────────────────
# Categorized help articles. Reads from __knowledgebase__/ directory.
# knowledgebase:
#   enabled: true
#   title: "Knowledge Base"
#   description: "Find answers and solutions"
#   path: "kb"
#   search:
#     enabled: true
#     placeholder: "Search articles..."

# ── Legal Pages ──────────────────────────────────────────────
# Privacy policy, terms of service, etc. Reads from __legal__/ directory.
# legal:
#   enabled: true
#   path: "legal"
#   footer_group: "Legal"         # Footer column header for legal links

# ── PDF Export ───────────────────────────────────────────────
# Add an "Export PDF" button to documentation pages.
# pdf_export:
#   enabled: true
#   page_break_level: 1           # Page break before: 1=h1, 2=h1+h2, 0=none

# ── Claude Assist ────────────────────────────────────────────
# "Ask Claude" button that copies page context for AI assistance.
# claude_assist:
#   enabled: true
#   label: "Ask Claude"
#   prompt: ""                    # Custom prompt prefix

# ── Analytics ────────────────────────────────────────────────
# Support for GA4, Plausible, Umami, Matomo, Fathom, Simple Analytics, and custom providers.
# analytics:
#   enabled: true
#   providers:
#     - type: plausible
#       enabled: true
#       config:
#         domain: "docs.example.com"
#         src: "https://plausible.io/js/script.js"
#     # - type: ga4
#     #   enabled: true
#     #   config:
#     #     measurement_id: "G-XXXXXXXXXX"
#     # - type: umami
#     #   enabled: true
#     #   config:
#     #     website_id: "your-website-id"
#     #     src: "https://analytics.example.com/script.js"

# ── Multi-Version Documentation ──────────────────────────────
# Serve multiple doc versions side-by-side. Each version reads from
# a subdirectory under __versions__/ (e.g., __versions__/v1/).
# versions:
#   enabled: true
#   default: "v2"
#   list:
#     - name: v2
#       label: "2.x (Latest)"
#       path: v2
#     - name: v1
#       label: "1.x"
#       path: v1
#       eol: "2025-12-31"         # Show end-of-life warning
#   selector:
#     position: header            # header | sidebar
#     show_eol_warning: true

# ── Internationalization (i18n) ──────────────────────────────
# Multi-language documentation. Each locale reads from a subdirectory
# under __translations__/ (e.g., __translations__/fr/).
# i18n:
#   enabled: true
#   default_locale: en
#   fallback: en
#   hide_default_locale: true     # /page.html instead of /en/page.html
#   show_untranslated: true       # Show fallback content with warning
#   locales:
#     - code: en
#       name: English
#     - code: fr
#       name: "Français"
#     - code: ar
#       name: "العربية"
#       direction: rtl
#   selector:
#     position: header
#     show_flags: false

# ── Footer ───────────────────────────────────────────────────
# Footer for landing and standalone pages.
# footer:
#   copyright: "© 2026 Your Company. All rights reserved."
#   hideVersion: false            # Hide "minimaldoc vX.Y.Z" in footer
#   links:
#     - title: Product
#       items:
#         - text: Documentation
#           url: /
#         - text: Changelog
#           url: /changelog/
#     - title: Company
#       items:
#         - text: About
#           url: https://example.com/about
#   social:
#     - name: GitHub
#       url: https://github.com/your-org
#       icon: github
#   badges:                       # "Powered by" badges
#     - text: "Built with MinimalDoc"
#       url: https://minimaldoc.dev

# ── Link Checker ─────────────────────────────────────────────
# Validates internal (and optionally external) links at build time.
link_check:
  enabled: true                   # Run link checker during build
  mode: warn                      # error | warn | ignore
  check_external: false
  # external_timeout: 5           # Seconds per external request
  # ignore_patterns: []           # Glob patterns to skip
  # allowed_broken: []            # Known broken links to ignore

# ── Custom Fields ────────────────────────────────────────────
# Arbitrary key-value pairs accessible in templates via .Site.Config.Custom
# custom:
#   company_name: "Acme Corp"
#   support_email: "support@example.com"
`

// OpenAPI config section (appended when --with-openapi is used)
var openAPIConfigSection = `
# OpenAPI Documentation
openapi:
  enabled: true
  spec_files:
    - "api/openapi.yaml"
  default_view: "path"            # path | tag | flat
  enable_testing: true            # Show "Try It" request tester
  enable_export: true             # cURL / restcli export buttons
  enable_code_samples: true       # Show code samples panel
`

// Status config section (appended when --with-status is used)
var statusConfigSection = `
# Status Page
status:
  enabled: true
  title: "Service Status"
  description: "Current operational status"
  path: "status"
  show_history: true
  history_months: 12
  rss_enabled: true
`

// Changelog config section (appended when --with-changelog is used)
var changelogConfigSection = `
# Changelog
changelog:
  enabled: true
  title: "Changelog"
  description: "All notable changes to this project"
  path: "changelog"
  rss_enabled: true
`

// Index page
var indexContent = `---
title: Welcome
description: Welcome to the documentation
---

# Welcome

Welcome to your documentation site built with MinimalDoc.

## Getting Started

- [Installation](getting-started/installation.html) - Install and set up
- [Quick Start](getting-started/quick-start.html) - Build your first site

## Guides

- [Deployment](guides/deployment.html) - Deploy to production

## Features

- Markdown documentation with frontmatter
- Automatic navigation from folder structure
- Dark mode support
- Full-text search
- LLM-friendly output (llms.txt)
`

// Installation page
var installContent = `---
title: Installation
description: How to install MinimalDoc
---

# Installation

## Prerequisites

- Go 1.21 or higher

## Install from Source

` + "```bash" + `
git clone https://github.com/studiowebux/minimaldoc.git
cd minimaldoc
go build -o minimaldoc ./cmd/minimaldoc
` + "```" + `

## Verify Installation

` + "```bash" + `
minimaldoc --version
` + "```" + `

## Next Steps

Continue to [Quick Start](quick-start.html) to create your first site.
`

// Quick start page
var quickstartContent = `---
title: Quick Start
description: Create your first documentation site
---

# Quick Start

Create a documentation site in minutes.

## Initialize

` + "```bash" + `
minimaldoc init docs
` + "```" + `

## Build

` + "```bash" + `
minimaldoc build docs
` + "```" + `

## Preview

Open ` + "`public/index.html`" + ` in your browser.

## Directory Structure

` + "```" + `
docs/
├── config.yaml      # Site configuration
├── index.md         # Homepage
├── getting-started/
│   ├── installation.md
│   └── quick-start.md
└── guides/
    └── deployment.md
` + "```" + `

## Configuration

Edit ` + "`config.yaml`" + ` to customize:

` + "```yaml" + `
title: My Documentation
description: Documentation for my project
theme: default
` + "```" + `
`

// Deployment guide
var deploymentContent = `---
title: Deployment
description: Deploy your documentation site
---

# Deployment

MinimalDoc generates static HTML files that can be hosted anywhere.

## Build for Production

` + "```bash" + `
minimaldoc build docs -o public
` + "```" + `

## GitHub Pages

1. Push the ` + "`public/`" + ` directory to ` + "`gh-pages`" + ` branch
2. Enable GitHub Pages in repository settings

## Netlify

1. Connect your repository
2. Set build command: ` + "`minimaldoc build docs`" + `
3. Set publish directory: ` + "`public`" + `

## Vercel

1. Import your repository
2. Override build command: ` + "`minimaldoc build docs`" + `
3. Set output directory: ` + "`public`" + `
`

// TOC template
var tocContent = `# Table of Contents

Custom navigation structure for your documentation.

- [Home](index.md)
- Getting Started
  - [Installation](getting-started/installation.md)
  - [Quick Start](getting-started/quick-start.md)
- Guides
  - [Deployment](guides/deployment.md)

<!--
TOC.md overrides automatic navigation.
Use indentation (2 spaces) for nested items.
Plain text creates section headers.
-->
`

// OpenAPI template
var openapiContent = `openapi: 3.0.3
info:
  title: Example API
  description: Example API specification
  version: 1.0.0
servers:
  - url: https://api.example.com/v1
    description: Production
  - url: https://staging-api.example.com/v1
    description: Staging
paths:
  /users:
    get:
      summary: List users
      operationId: listUsers
      tags:
        - Users
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            default: 10
      responses:
        '200':
          description: List of users
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'
    post:
      summary: Create user
      operationId: createUser
      tags:
        - Users
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
  /users/{id}:
    get:
      summary: Get user by ID
      operationId: getUser
      tags:
        - Users
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: User details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '404':
          description: User not found
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          example: "usr_123"
        name:
          type: string
          example: "John Doe"
        email:
          type: string
          format: email
          example: "john@example.com"
        created_at:
          type: string
          format: date-time
    CreateUserRequest:
      type: object
      required:
        - name
        - email
      properties:
        name:
          type: string
        email:
          type: string
          format: email
`

// Components template
var componentsContent = `# Status Page Components

- id: api
  name: API
  description: Core REST API
  status: operational
  group: Core Services
  order: 1

- id: web-app
  name: Web Application
  description: Main web interface
  status: operational
  group: Core Services
  order: 2

- id: database
  name: Database
  description: Primary database
  status: operational
  group: Infrastructure
  order: 1

# STATUS VALUES:
# operational, degraded, partial_outage, major_outage, maintenance
`

// Incident template
var incidentContent = `---
title: Example Incident
status: resolved
severity: minor
affected_components:
  - api
created_at: ` + time.Now().AddDate(0, 0, -7).Format("2006-01-02T15:04:05Z07:00") + `
resolved_at: ` + time.Now().AddDate(0, 0, -7).Add(2*time.Hour).Format("2006-01-02T15:04:05Z07:00") + `
---

## Resolved

Issue has been resolved. All services operating normally.

## Investigating

We are investigating reports of degraded API performance.
`

// Release template
var releaseContent = `---
version: 0.1.0
date: ` + time.Now().Format("2006-01-02") + `
---

# 0.1.0

Initial release.

## Added

- Basic documentation structure
- Markdown support
- Navigation generation

## Changed

- N/A

## Fixed

- N/A
`
