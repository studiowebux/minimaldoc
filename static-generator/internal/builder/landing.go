package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
	"github.com/studiowebux/minimaldoc/static-generator/internal/parser"
)

// LandingBuilder handles building the landing page
type LandingBuilder struct {
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
}

// NewLandingBuilder creates a new landing builder
func NewLandingBuilder() *LandingBuilder {
	return &LandingBuilder{
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
	}
}

// landingFile represents a parsed landing markdown file
type landingFile struct {
	Section  string
	Order    int
	Meta     core.Metadata
	Content  string
	HTML     string
	FilePath string
}

// Build creates the landing page from configuration and markdown files
func (lb *LandingBuilder) Build(docsRoot string, config core.LandingConfig, basePath string) (*core.LandingPage, error) {
	if !config.Enabled {
		return nil, nil
	}

	page := &core.LandingPage{
		Config: config,
	}

	// Parse markdown files from landing directory
	landingDir := filepath.Join(docsRoot, core.LandingSourceDir)
	files, err := lb.parseLandingFiles(landingDir, basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse landing files: %w", err)
	}

	// Apply YAML config as base, then override with markdown files
	lb.applyConfig(page, config)
	lb.applyMarkdownFiles(page, files)

	return page, nil
}

// parseLandingFiles parses all markdown files in the landing directory
func (lb *LandingBuilder) parseLandingFiles(landingDir string, basePath string) ([]landingFile, error) {
	var files []landingFile

	if _, err := os.Stat(landingDir); os.IsNotExist(err) {
		return files, nil
	}

	entries, err := os.ReadDir(landingDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(landingDir, entry.Name())
		file, err := lb.parseLandingFile(filePath, basePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse landing file %s: %v\n", filePath, err)
			continue
		}

		files = append(files, *file)
	}

	// Sort by order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Order < files[j].Order
	})

	return files, nil
}

// parseLandingFile parses a single landing markdown file
func (lb *LandingBuilder) parseLandingFile(filePath, basePath string) (*landingFile, error) {
	meta, content, err := lb.frontmatterParser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("frontmatter parse error: %w", err)
	}

	html, err := lb.markdownParser.ParseWithContext(content, "", basePath)
	if err != nil {
		return nil, fmt.Errorf("markdown parse error: %w", err)
	}

	// Extract section from frontmatter or filename
	section := lb.extractSection(filePath, meta)
	order := extractOrder(filepath.Base(filePath))

	return &landingFile{
		Section:  section,
		Order:    order,
		Meta:     meta,
		Content:  string(content),
		HTML:     string(html),
		FilePath: filePath,
	}, nil
}

// extractSection determines the section type from frontmatter or filename
func (lb *LandingBuilder) extractSection(filePath string, meta core.Metadata) string {
	// Check frontmatter for explicit section
	if meta.Custom != nil {
		if section, ok := meta.Custom["section"].(string); ok && section != "" {
			return strings.ToLower(section)
		}
	}

	// Extract from filename: 01-hero.md -> hero
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Remove numeric prefix
	re := regexp.MustCompile(`^\d+-`)
	name = re.ReplaceAllString(name, "")

	return strings.ToLower(name)
}

// applyConfig applies YAML configuration to the landing page
func (lb *LandingBuilder) applyConfig(page *core.LandingPage, config core.LandingConfig) {
	if config.Hero.Title != "" {
		page.Hero = &config.Hero
	}

	if len(config.Features.Items) > 0 {
		page.Features = &config.Features
	}

	if len(config.Steps.Items) > 0 {
		page.Steps = &config.Steps
	}

	if config.CTA.Title != "" {
		page.CTA = &config.CTA
	}

	if len(config.Testimonials.Items) > 0 {
		page.Testimonials = &config.Testimonials
	}

	if config.OpenSource.Title != "" {
		page.OpenSource = &config.OpenSource
	}

	if len(config.Links.Items) > 0 {
		page.Links = &config.Links
	}

	// Apply new section types
	if len(config.ImageText) > 0 {
		for i := range config.ImageText {
			page.ImageText = append(page.ImageText, &config.ImageText[i])
		}
	}

	if len(config.TextBlocks) > 0 {
		for i := range config.TextBlocks {
			page.TextBlocks = append(page.TextBlocks, &config.TextBlocks[i])
		}
	}

	if len(config.LinksGrid) > 0 {
		for i := range config.LinksGrid {
			page.LinksGrid = append(page.LinksGrid, &config.LinksGrid[i])
		}
	}
}

// applyMarkdownFiles overrides sections with markdown file content
func (lb *LandingBuilder) applyMarkdownFiles(page *core.LandingPage, files []landingFile) {
	for _, file := range files {
		switch file.Section {
		case "hero":
			lb.applyHero(page, file)
		case "features":
			lb.applyFeatures(page, file)
		case "steps":
			lb.applySteps(page, file)
		case "cta":
			lb.applyCTA(page, file)
		case "testimonials":
			lb.applyTestimonials(page, file)
		case "opensource":
			lb.applyOpenSource(page, file)
		case "links":
			lb.applyLinks(page, file)
		}
	}
}

// applyHero applies hero section from markdown
func (lb *LandingBuilder) applyHero(page *core.LandingPage, file landingFile) {
	if page.Hero == nil {
		page.Hero = &core.HeroSection{}
	}

	if file.Meta.Title != "" {
		page.Hero.Title = file.Meta.Title
	}
	if file.Meta.Description != "" {
		page.Hero.Subtitle = file.Meta.Description
	}
	if file.Meta.Custom != nil {
		if image, ok := file.Meta.Custom["image"].(string); ok {
			page.Hero.Image = image
		}
		if buttons, ok := file.Meta.Custom["buttons"].([]interface{}); ok {
			page.Hero.Buttons = lb.parseButtons(buttons)
		}
	}
}

// applyFeatures applies features section from markdown
func (lb *LandingBuilder) applyFeatures(page *core.LandingPage, file landingFile) {
	if page.Features == nil {
		page.Features = &core.FeaturesSection{}
	}

	if file.Meta.Title != "" {
		page.Features.Title = file.Meta.Title
	}
	if file.Meta.Custom != nil {
		if items, ok := file.Meta.Custom["items"].([]interface{}); ok {
			page.Features.Items = lb.parseFeatureItems(items)
		}
	}
}

// applySteps applies steps section from markdown
func (lb *LandingBuilder) applySteps(page *core.LandingPage, file landingFile) {
	if page.Steps == nil {
		page.Steps = &core.StepsSection{}
	}

	if file.Meta.Title != "" {
		page.Steps.Title = file.Meta.Title
	}
	if file.Meta.Custom != nil {
		if items, ok := file.Meta.Custom["items"].([]interface{}); ok {
			page.Steps.Items = lb.parseStepItems(items)
		}
	}
}

// applyCTA applies CTA section from markdown
func (lb *LandingBuilder) applyCTA(page *core.LandingPage, file landingFile) {
	if page.CTA == nil {
		page.CTA = &core.CTASection{}
	}

	if file.Meta.Title != "" {
		page.CTA.Title = file.Meta.Title
	}
	if file.Meta.Description != "" {
		page.CTA.Description = file.Meta.Description
	}
	if file.Meta.Custom != nil {
		if buttons, ok := file.Meta.Custom["buttons"].([]interface{}); ok {
			page.CTA.Buttons = lb.parseButtons(buttons)
		}
	}
}

// applyTestimonials applies testimonials section from markdown
func (lb *LandingBuilder) applyTestimonials(page *core.LandingPage, file landingFile) {
	if page.Testimonials == nil {
		page.Testimonials = &core.TestimonialsSection{}
	}

	if file.Meta.Title != "" {
		page.Testimonials.Title = file.Meta.Title
	}
	if file.Meta.Custom != nil {
		if items, ok := file.Meta.Custom["items"].([]interface{}); ok {
			page.Testimonials.Items = lb.parseTestimonialItems(items)
		}
	}
}

// applyOpenSource applies opensource section from markdown
func (lb *LandingBuilder) applyOpenSource(page *core.LandingPage, file landingFile) {
	if page.OpenSource == nil {
		page.OpenSource = &core.OpenSourceSection{}
	}

	if file.Meta.Title != "" {
		page.OpenSource.Title = file.Meta.Title
	}
	if file.Meta.Description != "" {
		page.OpenSource.Description = file.Meta.Description
	}
	// Use markdown body as description if provided
	if file.HTML != "" && file.Meta.Description == "" {
		page.OpenSource.Description = file.HTML
	}
	if file.Meta.Custom != nil {
		if links, ok := file.Meta.Custom["links"].([]interface{}); ok {
			page.OpenSource.Links = lb.parseSimpleLinks(links)
		}
	}
}

// applyLinks applies links section from markdown
func (lb *LandingBuilder) applyLinks(page *core.LandingPage, file landingFile) {
	if page.Links == nil {
		page.Links = &core.LinksSection{}
	}

	if file.Meta.Title != "" {
		page.Links.Title = file.Meta.Title
	}
	if file.Meta.Description != "" {
		page.Links.Description = file.Meta.Description
	}
	if file.Meta.Custom != nil {
		if items, ok := file.Meta.Custom["items"].([]interface{}); ok {
			page.Links.Items = lb.parseLinkItems(items)
		}
	}
}

// parseButtons parses button configurations from frontmatter
func (lb *LandingBuilder) parseButtons(items []interface{}) []core.HeroButton {
	var buttons []core.HeroButton
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			btn := core.HeroButton{}
			if text, ok := m["text"].(string); ok {
				btn.Text = text
			}
			if url, ok := m["url"].(string); ok {
				btn.URL = url
			}
			if primary, ok := m["primary"].(bool); ok {
				btn.Primary = primary
			}
			buttons = append(buttons, btn)
		}
	}
	return buttons
}

// parseFeatureItems parses feature items from frontmatter
func (lb *LandingBuilder) parseFeatureItems(items []interface{}) []core.FeatureItem {
	var features []core.FeatureItem
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			f := core.FeatureItem{}
			if icon, ok := m["icon"].(string); ok {
				f.Icon = icon
			}
			if emoji, ok := m["emoji"].(string); ok {
				f.Emoji = emoji
			}
			if title, ok := m["title"].(string); ok {
				f.Title = title
			}
			if desc, ok := m["description"].(string); ok {
				f.Description = desc
			}
			features = append(features, f)
		}
	}
	return features
}

// parseStepItems parses step items from frontmatter
func (lb *LandingBuilder) parseStepItems(items []interface{}) []core.StepItem {
	var steps []core.StepItem
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			s := core.StepItem{}
			if title, ok := m["title"].(string); ok {
				s.Title = title
			}
			if desc, ok := m["description"].(string); ok {
				s.Description = desc
			}
			if code, ok := m["code"].(string); ok {
				s.Code = code
			}
			steps = append(steps, s)
		}
	}
	return steps
}

// parseTestimonialItems parses testimonial items from frontmatter
func (lb *LandingBuilder) parseTestimonialItems(items []interface{}) []core.Testimonial {
	var testimonials []core.Testimonial
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			t := core.Testimonial{}
			if quote, ok := m["quote"].(string); ok {
				t.Quote = quote
			}
			if author, ok := m["author"].(string); ok {
				t.Author = author
			}
			if role, ok := m["role"].(string); ok {
				t.Role = role
			}
			if avatar, ok := m["avatar"].(string); ok {
				t.Avatar = avatar
			}
			testimonials = append(testimonials, t)
		}
	}
	return testimonials
}

// parseSimpleLinks parses simple link items from frontmatter
func (lb *LandingBuilder) parseSimpleLinks(items []interface{}) []core.SimpleLink {
	var links []core.SimpleLink
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			l := core.SimpleLink{}
			if text, ok := m["text"].(string); ok {
				l.Text = text
			}
			if url, ok := m["url"].(string); ok {
				l.URL = url
			}
			links = append(links, l)
		}
	}
	return links
}

// parseLinkItems parses link card items from frontmatter
func (lb *LandingBuilder) parseLinkItems(items []interface{}) []core.LinkItem {
	var links []core.LinkItem
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			l := core.LinkItem{}
			if icon, ok := m["icon"].(string); ok {
				l.Icon = icon
			}
			if title, ok := m["title"].(string); ok {
				l.Title = title
			}
			if desc, ok := m["description"].(string); ok {
				l.Description = desc
			}
			if url, ok := m["url"].(string); ok {
				l.URL = url
			}
			links = append(links, l)
		}
	}
	return links
}
