package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// LLMSGenerator generates llms.txt and LLM-friendly markdown files
type LLMSGenerator struct {
	site *core.Site
}

// NewLLMSGenerator creates a new LLMS generator
func NewLLMSGenerator(site *core.Site) *LLMSGenerator {
	return &LLMSGenerator{site: site}
}

// Generate creates llms.txt for LLM consumption
func (g *LLMSGenerator) Generate() error {
	if !g.site.Config.EnableLLMS {
		return nil
	}

	fmt.Println("Generating llms.txt...")

	// Generate full concatenated version only
	if err := g.generateFullVersion(); err != nil {
		return fmt.Errorf("failed to generate llms.txt: %w", err)
	}

	fmt.Println("Generated llms.txt")
	return nil
}

// generateIndex creates the main llms.txt index file
func (g *LLMSGenerator) generateIndex() error {
	var buf bytes.Buffer

	// Header
	buf.WriteString(fmt.Sprintf("# %s\n\n", g.site.Config.Title))

	if g.site.Config.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", g.site.Config.Description))
	}

	// Metadata
	buf.WriteString("## About\n\n")
	buf.WriteString("This is an LLM-friendly index of the documentation.\n\n")
	if g.site.Config.BaseURL != "" {
		buf.WriteString(fmt.Sprintf("Website: %s\n\n", g.site.Config.BaseURL))
	}

	// Navigation structure
	buf.WriteString("## Documentation Structure\n\n")

	// Sort pages by order
	sortedPages := make([]*core.Page, len(g.site.Pages))
	copy(sortedPages, g.site.Pages)
	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].Order != sortedPages[j].Order {
			return sortedPages[i].Order < sortedPages[j].Order
		}
		return sortedPages[i].Slug < sortedPages[j].Slug
	})

	// List all pages with links to their .md versions
	for _, page := range sortedPages {
		if page.IsHidden() {
			continue
		}

		// Create relative path to .md file
		mdPath := fmt.Sprintf("%s.md", page.Slug)

		buf.WriteString(fmt.Sprintf("- [%s](%s)\n", page.Title(), mdPath))

		if page.Metadata.Description != "" {
			buf.WriteString(fmt.Sprintf("  %s\n", page.Metadata.Description))
		}
	}

	buf.WriteString("\n")
	buf.WriteString("## Full Documentation\n\n")
	buf.WriteString("For a single-file version with all content, see [llms-full.txt](llms-full.txt)\n")

	// Write file
	outputPath := filepath.Join(g.site.OutputRoot, "llms.txt")
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

// generateMarkdownFiles creates individual .md files for each page
func (g *LLMSGenerator) generateMarkdownFiles() error {
	for _, page := range g.site.Pages {
		if err := g.generatePageMarkdown(page); err != nil {
			return fmt.Errorf("failed to generate markdown for %s: %w", page.Slug, err)
		}
	}
	return nil
}

// generatePageMarkdown creates a clean markdown file for a single page
func (g *LLMSGenerator) generatePageMarkdown(page *core.Page) error {
	var buf bytes.Buffer

	// Front matter
	buf.WriteString(fmt.Sprintf("# %s\n\n", page.Title()))

	if page.Metadata.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", page.Metadata.Description))
	}

	// Metadata
	if len(page.Metadata.Tags) > 0 {
		buf.WriteString(fmt.Sprintf("**Tags:** %s\n\n", strings.Join(page.Metadata.Tags, ", ")))
	}

	buf.WriteString("---\n\n")

	// Content (raw markdown)
	buf.Write(page.RawMD)

	// Create output path
	mdPath := filepath.Join(g.site.OutputRoot, page.Slug+".md")
	mdDir := filepath.Dir(mdPath)

	// Create directory
	if err := os.MkdirAll(mdDir, 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(mdPath, buf.Bytes(), 0644)
}

// generateFullVersion creates a single file with all documentation
func (g *LLMSGenerator) generateFullVersion() error {
	var buf bytes.Buffer

	// Header
	buf.WriteString(fmt.Sprintf("# %s - Complete Documentation\n\n", g.site.Config.Title))

	if g.site.Config.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", g.site.Config.Description))
	}

	buf.WriteString("---\n\n")
	buf.WriteString("This file contains all documentation in a single file for LLM consumption.\n\n")
	buf.WriteString("---\n\n")

	// Sort pages by order
	sortedPages := make([]*core.Page, len(g.site.Pages))
	copy(sortedPages, g.site.Pages)
	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].Order != sortedPages[j].Order {
			return sortedPages[i].Order < sortedPages[j].Order
		}
		return sortedPages[i].Slug < sortedPages[j].Slug
	})

	// Concatenate all pages
	for i, page := range sortedPages {
		if page.IsHidden() {
			continue
		}

		// Add page marker
		buf.WriteString(fmt.Sprintf("\n\n<!-- PAGE: %s -->\n\n", page.Slug))

		// Add title
		buf.WriteString(fmt.Sprintf("# %s\n\n", page.Title()))

		if page.Metadata.Description != "" {
			buf.WriteString(fmt.Sprintf("%s\n\n", page.Metadata.Description))
		}

		// Add content
		buf.Write(page.RawMD)

		// Add separator between pages (except for last page)
		if i < len(sortedPages)-1 {
			buf.WriteString("\n\n---\n\n")
		}
	}

	// Write file (using llms.txt as the filename)
	outputPath := filepath.Join(g.site.OutputRoot, "llms.txt")
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}
