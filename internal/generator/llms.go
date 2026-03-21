package generator

import (
	"bytes"
	"fmt"
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

// Generate creates llms.txt (index), llms-full.txt (content), and companion .md files
func (g *LLMSGenerator) Generate() error {
	if !g.site.Config.EnableLLMS {
		return nil
	}

	fmt.Println("Generating llms.txt...")

	// Generate index file (llms.txt)
	if err := g.generateIndex(); err != nil {
		return fmt.Errorf("failed to generate llms.txt: %w", err)
	}

	// Generate full content file (llms-full.txt)
	if err := g.generateFullVersion(); err != nil {
		return fmt.Errorf("failed to generate llms-full.txt: %w", err)
	}

	// Generate companion .md files for each page
	if err := g.generateMarkdownFiles(); err != nil {
		return fmt.Errorf("failed to generate companion .md files: %w", err)
	}

	fmt.Println("Generated llms.txt, llms-full.txt, and companion .md files")
	return nil
}

// generateIndex creates the main llms.txt index file following the spec
// Format: H1 title, blockquote summary, H2 sections with file lists
func (g *LLMSGenerator) generateIndex() error {
	var buf bytes.Buffer

	// H1: Project name (required)
	buf.WriteString(fmt.Sprintf("# %s\n\n", g.site.Config.Title))

	// Blockquote: Summary (optional but recommended)
	if g.site.Config.Description != "" {
		buf.WriteString(fmt.Sprintf("> %s\n\n", g.site.Config.Description))
	}

	// Sort pages by order
	sortedPages := make([]*core.Page, len(g.site.Pages))
	copy(sortedPages, g.site.Pages)
	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].Order != sortedPages[j].Order {
			return sortedPages[i].Order < sortedPages[j].Order
		}
		return sortedPages[i].Slug < sortedPages[j].Slug
	})

	// Group pages by section (first directory)
	sections := make(map[string][]*core.Page)
	var sectionOrder []string

	for _, page := range sortedPages {
		if page.IsHidden() {
			continue
		}

		section := "Documentation"
		if idx := strings.Index(page.Slug, "/"); idx > 0 {
			section = titleCase(strings.ReplaceAll(page.Slug[:idx], "-", " "))
		}

		if _, exists := sections[section]; !exists {
			sectionOrder = append(sectionOrder, section)
		}
		sections[section] = append(sections[section], page)
	}

	// H2 sections with file lists
	baseURL := strings.TrimSuffix(g.site.Config.BaseURL, "/")

	for _, section := range sectionOrder {
		pages := sections[section]
		buf.WriteString(fmt.Sprintf("## %s\n\n", section))

		for _, page := range pages {
			url := fmt.Sprintf("%s/%s.html", baseURL, page.Slug)
			if page.Metadata.Description != "" {
				buf.WriteString(fmt.Sprintf("- [%s](%s): %s\n", page.Title(), url, page.Metadata.Description))
			} else {
				buf.WriteString(fmt.Sprintf("- [%s](%s)\n", page.Title(), url))
			}
		}
		buf.WriteString("\n")
	}

	// Optional section for full content
	buf.WriteString("## Optional\n\n")
	buf.WriteString(fmt.Sprintf("- [Complete Documentation](%s/llms-full.txt): All documentation in a single file\n", baseURL))

	// Write file
	outputPath := filepath.Join(g.site.OutputRoot, "llms.txt")
	return writeWebFile(outputPath, buf.Bytes())
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

	// Create output path (page.html.md for LLM consumption)
	mdPath := filepath.Join(g.site.OutputRoot, page.Slug+".html.md")
	mdDir := filepath.Dir(mdPath)

	// Create directory
	if err := makeWebDir(mdDir); err != nil {
		return err
	}

	// Write file
	return writeWebFile(mdPath, buf.Bytes())
}

// generateFullVersion creates llms-full.txt with all documentation content
func (g *LLMSGenerator) generateFullVersion() error {
	var buf bytes.Buffer

	// H1: Project name
	buf.WriteString(fmt.Sprintf("# %s\n\n", g.site.Config.Title))

	// Blockquote: Summary
	if g.site.Config.Description != "" {
		buf.WriteString(fmt.Sprintf("> %s\n>\n", g.site.Config.Description))
		buf.WriteString("> This file contains the complete documentation.\n\n")
	}

	// Sort pages by order
	sortedPages := make([]*core.Page, len(g.site.Pages))
	copy(sortedPages, g.site.Pages)
	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].Order != sortedPages[j].Order {
			return sortedPages[i].Order < sortedPages[j].Order
		}
		return sortedPages[i].Slug < sortedPages[j].Slug
	})

	// Concatenate all pages as H2 sections
	for _, page := range sortedPages {
		if page.IsHidden() {
			continue
		}

		// H2: Page title
		buf.WriteString(fmt.Sprintf("## %s\n\n", page.Title()))

		if page.Metadata.Description != "" {
			buf.WriteString(fmt.Sprintf("*%s*\n\n", page.Metadata.Description))
		}

		// Content (shift headings down: # -> ###, ## -> ####, etc.)
		content := shiftHeadings(string(page.RawMD), 2)
		buf.WriteString(content)
		buf.WriteString("\n\n")
	}

	// Write file
	outputPath := filepath.Join(g.site.OutputRoot, "llms-full.txt")
	return writeWebFile(outputPath, buf.Bytes())
}

// shiftHeadings shifts markdown heading levels down by n levels
func shiftHeadings(content string, levels int) string {
	lines := strings.Split(content, "\n")
	inCodeBlock := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
		}
		if inCodeBlock {
			continue
		}
		if strings.HasPrefix(line, "#") {
			// Count existing hashes
			hashes := 0
			for _, c := range line {
				if c == '#' {
					hashes++
				} else {
					break
				}
			}
			if hashes > 0 && hashes <= 6 {
				// Shift heading (cap at h6)
				newLevel := hashes + levels
				if newLevel > 6 {
					newLevel = 6
				}
				lines[i] = strings.Repeat("#", newLevel) + line[hashes:]
			}
		}
	}

	return strings.Join(lines, "\n")
}

// titleCase converts string to title case (first letter of each word uppercase)
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}
