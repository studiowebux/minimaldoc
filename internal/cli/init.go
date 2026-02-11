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
		if err := os.MkdirAll(dir, 0755); err != nil {
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
		filepath.Join(targetDir, "config.yaml"):                       finalConfig,
		filepath.Join(targetDir, "index.md"):                          indexContent,
		filepath.Join(targetDir, "getting-started", "installation.md"): installContent,
		filepath.Join(targetDir, "getting-started", "quick-start.md"): quickstartContent,
		filepath.Join(targetDir, "guides", "deployment.md"):           deploymentContent,
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
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
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
var configContent = `# MinimalDoc Configuration
# Full reference: https://minimaldoc.dev/getting-started/configuration.html

# Site Information
title: My Documentation
description: Documentation for my project
base_url: ""                    # Set for production (e.g., https://docs.example.com)
author: Your Name

# Theme & Appearance
theme: default                  # default | minimal | dark
dark_mode: false                # Start in dark mode

# Features
enable_llms: true               # Generate llms.txt for AI tools
enable_search: true             # Enable client-side search (Cmd+K)
clean_urls: false               # Use /page/ instead of /page.html

# Stale Content Warnings
stale_warning:
  enabled: false
  threshold_days: 365
  show_update_date: true

# Social Links (shown in sidebar)
social_links:
  - name: GitHub
    url: https://github.com/your-org/your-project
    icon: github

# Available icons: github, twitter, linkedin, youtube, discord, mastodon, rss, email, website
`

// OpenAPI config section (appended when --with-openapi is used)
var openAPIConfigSection = `
# OpenAPI Documentation
openapi:
  enabled: true
  spec_files:
    - "api/openapi.yaml"
  default_view: "path"            # path | tag | flat
  enable_testing: true            # Show "Try It" interface
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
