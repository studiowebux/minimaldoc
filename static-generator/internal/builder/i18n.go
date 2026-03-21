package builder

import (
	"os"
	"path/filepath"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
	"github.com/studiowebux/minimaldoc/static-generator/internal/parser"
)

// I18nBuilder handles building locale-specific page sets
type I18nBuilder struct {
	frontmatterParser *parser.FrontmatterParser
	markdownParser    *parser.MarkdownParser
	tocParser         *parser.TOCParser
}

// NewI18nBuilder creates a new i18n builder
func NewI18nBuilder() *I18nBuilder {
	return &I18nBuilder{
		frontmatterParser: parser.NewFrontmatterParser(),
		markdownParser:    parser.NewMarkdownParser(),
		tocParser:         parser.NewTOCParser(),
	}
}

// Build builds locale-specific page sets from shared pages and translations
func (ib *I18nBuilder) Build(site *core.Site) error {
	if !site.Config.I18n.Enabled || len(site.Config.I18n.Locales) == 0 {
		return nil
	}

	// Initialize the localized pages map
	site.LocalizedPages = make(map[string][]*core.Page)

	// For each configured locale, build its page set
	for _, locale := range site.Config.I18n.Locales {
		pages, err := ib.buildLocalePages(site, locale)
		if err != nil {
			return err
		}
		site.LocalizedPages[locale.Code] = pages
	}

	return nil
}

// buildLocalePages builds the page set for a specific locale
func (ib *I18nBuilder) buildLocalePages(site *core.Site, locale core.LocaleInfo) ([]*core.Page, error) {
	var pages []*core.Page

	// 1. Start with shared pages (default locale content as fallback)
	for _, page := range site.Pages {
		// Clone the page for this locale
		localePage := ib.clonePage(page, locale, site.Config.I18n)
		pages = append(pages, localePage)
	}

	// 2. Look for locale-specific translations in __translations__/{locale}/
	translationDir := filepath.Join(site.DocsRoot, core.I18nSourceDir, locale.Code)
	if _, err := os.Stat(translationDir); err == nil {
		translations, err := ib.discoverTranslations(translationDir, site.DocsRoot, locale, site.Config.I18n)
		if err != nil {
			return nil, err
		}

		// Merge translations with existing pages (translation wins)
		pages = ib.mergeTranslations(pages, translations)
	}

	return pages, nil
}

// clonePage creates a copy of a page for a specific locale
func (ib *I18nBuilder) clonePage(page *core.Page, locale core.LocaleInfo, config core.I18nConfig) *core.Page {
	clone := &core.Page{
		SourcePath:      page.SourcePath,
		RelPath:         page.RelPath,
		Slug:            page.Slug,
		OutputPath:      page.OutputPath,
		ModTime:         page.ModTime,
		Metadata:        cloneMetadata(page.Metadata),
		RawMD:           page.RawMD,
		HTML:            page.HTML,
		Order:           page.Order,
		TOC:             page.TOC,
		IsStale:         page.IsStale,
		StaleAge:        page.StaleAge,
		DaysSinceUpdate: page.DaysSinceUpdate,
	}

	// Update output path for localized output
	// For non-default locales or if HideDefaultLocale is false, prepend locale path
	if locale.Code != config.DefaultLocale || !config.HideDefaultLocale {
		clone.OutputPath = filepath.Join(locale.Code, page.OutputPath)
	}

	return clone
}

// discoverTranslations finds locale-specific translation pages
func (ib *I18nBuilder) discoverTranslations(translationDir, _ string, locale core.LocaleInfo, config core.I18nConfig) ([]*core.Page, error) {
	var pages []*core.Page

	err := filepath.WalkDir(translationDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Create page with relative path from translation directory
		relPath, _ := filepath.Rel(translationDir, path)
		page := core.NewPage(path, translationDir)
		page.RelPath = relPath

		// Parse frontmatter
		meta, content, err := ib.frontmatterParser.ParseFile(path)
		if err != nil {
			return err
		}
		page.Metadata = meta
		page.RawMD = content

		// Override order if specified
		if meta.MenuOrder >= 0 {
			page.Order = meta.MenuOrder
		} else {
			page.Order = extractOrder(filepath.Base(path))
		}

		// Parse markdown with link transformation
		html, err := ib.markdownParser.ParseWithContext(content, relPath, "")
		if err != nil {
			return err
		}
		page.HTML = html

		// Generate TOC
		toc, err := ib.tocParser.Parse(content)
		if err != nil {
			return err
		}
		page.TOC = toc

		// Set output path for localized output
		if locale.Code != config.DefaultLocale || !config.HideDefaultLocale {
			page.OutputPath = filepath.Join(locale.Code, page.Slug+".html")
		} else {
			page.OutputPath = page.Slug + ".html"
		}

		pages = append(pages, page)
		return nil
	})

	return pages, err
}

// mergeTranslations merges locale-specific translations with shared pages
func (ib *I18nBuilder) mergeTranslations(shared, translations []*core.Page) []*core.Page {
	// Create a map of shared pages by slug for quick lookup
	pageMap := make(map[string]*core.Page)
	for _, page := range shared {
		pageMap[page.Slug] = page
	}

	// Override or add pages from translations
	for _, translation := range translations {
		pageMap[translation.Slug] = translation
	}

	// Convert map back to slice
	result := make([]*core.Page, 0, len(pageMap))
	for _, page := range pageMap {
		result = append(result, page)
	}

	return result
}
