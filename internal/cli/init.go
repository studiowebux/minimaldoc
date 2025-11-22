package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new documentation site",
	Long: `Initialize creates a new documentation site with example files.

This will create:
- Sample markdown files with frontmatter
- A basic directory structure
- Configuration file (optional)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	fmt.Println("Initializing new documentation site...")

	// Create directory structure
	dirs := []string{
		targetDir,
		filepath.Join(targetDir, "01-getting-started"),
		filepath.Join(targetDir, "02-guides"),
		filepath.Join(targetDir, "03-api"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create example files
	files := map[string]string{
		filepath.Join(targetDir, "index.md"):                                indexContent,
		filepath.Join(targetDir, "01-getting-started", "01-introduction.md"): introContent,
		filepath.Join(targetDir, "01-getting-started", "02-installation.md"): installContent,
		filepath.Join(targetDir, "02-guides", "01-quickstart.md"):            quickstartContent,
		filepath.Join(targetDir, "03-api", "01-reference.md"):                apiContent,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
	}

	fmt.Println()
	fmt.Println("✓ Documentation site initialized!")
	fmt.Printf("✓ Location: %s\n", targetDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit the markdown files in the created directories")
	fmt.Println("  2. Run 'minimaldoc build' to generate your site")
	fmt.Println("  3. Open public/index.html in your browser")
	fmt.Println()

	return nil
}

// Example content for generated files

var indexContent = "---\n" +
	"title: Welcome\n" +
	"description: Welcome to the documentation\n" +
	"tags: [home, welcome]\n" +
	"menu_order: 0\n" +
	"---\n\n" +
	"# Welcome to Your Documentation\n\n" +
	"This is the home page of your documentation site built with **Minimal Doc**.\n\n" +
	"## Features\n\n" +
	"- 📝 Write documentation in Markdown\n" +
	"- 🎨 Beautiful, minimal theme\n" +
	"- 🌓 Light and dark mode\n" +
	"- 📱 Responsive design\n" +
	"- 🔍 Easy navigation\n" +
	"- 🤖 LLM-friendly output\n\n" +
	"## Getting Started\n\n" +
	"Check out the [Introduction](01-getting-started/01-introduction.html) to learn more.\n"

var introContent = "---\n" +
	"title: Introduction\n" +
	"description: Introduction to Minimal Doc\n" +
	"tags: [intro, getting-started]\n" +
	"---\n\n" +
	"# Introduction\n\n" +
	"Welcome to **Minimal Doc**, a simple and elegant static site generator for documentation.\n\n" +
	"## What is Minimal Doc?\n\n" +
	"Minimal Doc converts your Markdown files into a beautiful documentation website.\n\n" +
	"## Next Steps\n\n" +
	"Continue to [Installation](02-installation.html) to get started.\n"

var installContent = "---\n" +
	"title: Installation\n" +
	"description: How to install Minimal Doc\n" +
	"tags: [installation, setup]\n" +
	"---\n\n" +
	"# Installation\n\n" +
	"Installing Minimal Doc is quick and easy.\n\n" +
	"## Build from Source\n\n" +
	"Requires Go 1.21 or higher:\n\n" +
	"```bash\n" +
	"go build -o minimaldoc cmd/minimaldoc/main.go\n" +
	"```\n\n" +
	"## Next Steps\n\n" +
	"Check out the [Quick Start Guide](../02-guides/01-quickstart.html).\n"

var quickstartContent = "---\n" +
	"title: Quick Start Guide\n" +
	"description: Build your first documentation site in minutes\n" +
	"tags: [quickstart, guide, tutorial]\n" +
	"---\n\n" +
	"# Quick Start Guide\n\n" +
	"Let's build your first documentation site in 5 minutes!\n\n" +
	"## Step 1: Initialize\n\n" +
	"```bash\n" +
	"minimaldoc init my-docs\n" +
	"cd my-docs\n" +
	"```\n\n" +
	"## Step 2: Build\n\n" +
	"```bash\n" +
	"minimaldoc build\n" +
	"```\n\n" +
	"## Step 3: Preview\n\n" +
	"Open `public/index.html` in your browser!\n"

var apiContent = "---\n" +
	"title: API Reference\n" +
	"description: Complete API documentation\n" +
	"tags: [api, reference]\n" +
	"---\n\n" +
	"# API Reference\n\n" +
	"## Build Command\n\n" +
	"```bash\n" +
	"minimaldoc build [flags]\n" +
	"```\n\n" +
	"### Flags\n\n" +
	"- `--output, -o` - Output directory (default: public)\n" +
	"- `--theme, -t` - Theme name (default: default)\n" +
	"- `--llms, -l` - Generate llms.txt (default: true)\n"
