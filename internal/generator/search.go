package generator

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// truncate returns s truncated to at most n bytes. Safe when len(s) < n.
func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// SearchPage represents a page in the search index
type SearchPage struct {
	Title string `json:"t"`
	Desc  string `json:"d,omitempty"`
	URL   string `json:"u"`
}

// SearchSection represents a section within a page
type SearchSection struct {
	PageID int    `json:"p"`
	Title  string `json:"t"`
	Anchor string `json:"a,omitempty"`
}

// PostingList is a compact array: [pageID, score, pageID, score, ...]
// Section info removed to save space - we use page-level results
type PostingList []int

// SearchIndex is the inverted index structure
type SearchIndex struct {
	Pages    []SearchPage           `json:"pages"`
	Sections []SearchSection        `json:"sections"`
	Index    map[string]PostingList `json:"idx"`
}

// SearchManifest is the lightweight manifest for sharded search
type SearchManifest struct {
	Pages    []SearchPage    `json:"pages"`
	Sections []SearchSection `json:"sections"`
	Shards   []string        `json:"shards"` // List of shard prefixes (e.g., ["ab", "cd", ...])
}

// SearchShard contains posting lists for terms with a specific prefix
type SearchShard struct {
	Prefix string                 `json:"prefix"`
	Index  map[string]PostingList `json:"idx"`
}

// SearchGenerator generates search index JSON
type SearchGenerator struct {
	site *core.Site
}

// NewSearchGenerator creates a new search generator
func NewSearchGenerator(site *core.Site) *SearchGenerator {
	return &SearchGenerator{site: site}
}

// Common English stop words to exclude from index
var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
	"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
	"or": true, "that": true, "the": true, "to": true, "was": true, "were": true,
	"will": true, "with": true, "this": true, "but": true, "they": true,
	"have": true, "had": true, "what": true, "when": true, "where": true,
	"who": true, "which": true, "why": true, "how": true, "all": true,
	"each": true, "every": true, "both": true, "few": true, "more": true,
	"most": true, "other": true, "some": true, "such": true, "no": true,
	"not": true, "only": true, "own": true, "same": true, "so": true,
	"than": true, "too": true, "very": true, "can": true, "just": true,
	"should": true, "now": true, "also": true, "into": true, "your": true,
	"you": true, "our": true, "we": true, "if": true, "then": true,
}

// Generate creates inverted index
func (g *SearchGenerator) Generate() error {
	fmt.Println("Generating search index...")

	// If versioning is enabled, generate per-version indexes
	if g.site.Config.Versions.Enabled && len(g.site.VersionedPages) > 0 {
		return g.generateVersionedIndexes()
	}

	return g.generateMainIndex()
}

// generateVersionedIndexes creates search indexes for each version
func (g *SearchGenerator) generateVersionedIndexes() error {
	defaultVersion := g.site.Config.Versions.Default
	if defaultVersion == "" && len(g.site.Config.Versions.List) > 0 {
		defaultVersion = g.site.Config.Versions.List[0].Name
	}

	for _, versionInfo := range g.site.Config.Versions.List {
		pages, ok := g.site.VersionedPages[versionInfo.Name]
		if !ok {
			continue
		}

		isDefault := versionInfo.Name == defaultVersion
		if err := g.generateVersionIndex(versionInfo, pages, isDefault); err != nil {
			return err
		}
	}

	return nil
}

// generateVersionIndex creates a search index for a specific version
func (g *SearchGenerator) generateVersionIndex(versionInfo core.VersionInfo, pages []*core.Page, isDefault bool) error {
	basePath := g.getBasePath()
	versionPrefix := ""
	if !isDefault {
		versionPrefix = "/" + versionInfo.Path
	}

	index := SearchIndex{
		Pages:    []SearchPage{},
		Sections: []SearchSection{},
		Index:    make(map[string]PostingList),
	}

	pageID := 0

	for _, page := range pages {
		if page.IsHidden() {
			continue
		}

		index.Pages = append(index.Pages, SearchPage{
			Title: page.Title(),
			Desc:  page.Metadata.Description,
			URL:   basePath + versionPrefix + "/" + page.Slug + ".html",
		})

		g.indexText(index.Index, page.Title(), pageID, 3)

		if page.Metadata.Description != "" {
			g.indexText(index.Index, page.Metadata.Description, pageID, 2)
		}

		for _, tag := range page.Metadata.Tags {
			g.indexText(index.Index, tag, pageID, 2)
		}

		sections := g.extractSections(page.RawMD)
		for _, section := range sections {
			index.Sections = append(index.Sections, SearchSection{
				PageID: pageID,
				Title:  section.title,
				Anchor: section.anchor,
			})

			g.indexText(index.Index, section.title, pageID, 3)
			g.indexText(index.Index, section.content, pageID, 1)
		}

		pageID++
	}

	// Index FAQ items (if enabled)
	if g.site.FaqPage != nil && g.site.Config.Faq.Enabled {
		faqPath := g.site.Config.Faq.Path
		if faqPath == "" {
			faqPath = "faq"
		}

		for _, cat := range g.site.FaqPage.Categories {
			for _, item := range cat.Items {
				index.Pages = append(index.Pages, SearchPage{
					Title: item.Question,
					Desc:  cat.Name,
					URL:   basePath + versionPrefix + "/" + faqPath + "/#" + item.Slug,
				})

				g.indexText(index.Index, item.Question, pageID, 3)
				g.indexText(index.Index, cat.Name, pageID, 2)

				for _, tag := range item.Tags {
					g.indexText(index.Index, tag, pageID, 2)
				}

				answerText := item.Answer
				if item.AnswerHTML != "" {
					answerText = extractPlainText([]byte(item.AnswerHTML))
				}
				g.indexText(index.Index, truncate(answerText, 1000), pageID, 1)

				pageID++
			}
		}
	}

	// Index KB articles (if enabled)
	if g.site.KBPage != nil && g.site.Config.KnowledgeBase.Enabled {
		kbPath := g.site.Config.KnowledgeBase.Path
		if kbPath == "" {
			kbPath = "kb"
		}

		for _, cat := range g.site.KBPage.Categories {
			for _, article := range cat.Articles {
				index.Pages = append(index.Pages, SearchPage{
					Title: article.Title,
					Desc:  cat.Name + " | " + article.Description,
					URL:   basePath + versionPrefix + "/" + kbPath + "/" + cat.Slug + "/" + article.Slug + ".html",
				})

				g.indexText(index.Index, article.Title, pageID, 3)
				g.indexText(index.Index, cat.Name, pageID, 2)

				if article.Description != "" {
					g.indexText(index.Index, article.Description, pageID, 2)
				}

				for _, tag := range article.Tags {
					g.indexText(index.Index, tag, pageID, 2)
				}

				g.indexText(index.Index, truncate(extractPlainText([]byte(article.HTML)), 1000), pageID, 1)

				pageID++
			}
		}
	}

	for word := range index.Index {
		list := index.Index[word]
		sortPostingList(list)
		index.Index[word] = list
	}

	indexData, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal search index: %w", err)
	}

	var indexPath string
	if isDefault {
		indexPath = filepath.Join(g.site.OutputRoot, "search-index.json")
	} else {
		indexPath = filepath.Join(g.site.OutputRoot, versionInfo.Path, "search-index.json")
		if err := makeWebDir(filepath.Dir(indexPath)); err != nil {
			return fmt.Errorf("failed to create version directory: %w", err)
		}
	}

	if err := writeWebFile(indexPath, indexData); err != nil {
		return fmt.Errorf("failed to write search index: %w", err)
	}

	// Also write sharded index
	var shardDir string
	if isDefault {
		shardDir = g.site.OutputRoot
	} else {
		shardDir = filepath.Join(g.site.OutputRoot, versionInfo.Path)
	}

	numShards, err := g.writeShardedIndex(index, shardDir)
	if err != nil {
		return fmt.Errorf("failed to write sharded index: %w", err)
	}

	fmt.Printf("  Generated search index for %s: %d pages, %d terms, %d shards\n",
		versionInfo.Name, len(index.Pages), len(index.Index), numShards)
	return nil
}

// generateMainIndex creates the main search index (non-versioned)
func (g *SearchGenerator) generateMainIndex() error {

	basePath := g.getBasePath()

	index := SearchIndex{
		Pages:    []SearchPage{},
		Sections: []SearchSection{},
		Index:    make(map[string]PostingList),
	}

	pageID := 0

	// Process all visible pages
	for _, page := range g.site.Pages {
		if page.IsHidden() {
			continue
		}

		// Add page metadata
		index.Pages = append(index.Pages, SearchPage{
			Title: page.Title(),
			Desc:  page.Metadata.Description,
			URL:   basePath + "/" + page.Slug + ".html",
		})

		// Index page title (high weight)
		g.indexText(index.Index, page.Title(), pageID, 3)

		// Index page description (medium weight)
		if page.Metadata.Description != "" {
			g.indexText(index.Index, page.Metadata.Description, pageID, 2)
		}

		// Index tags (medium weight)
		for _, tag := range page.Metadata.Tags {
			g.indexText(index.Index, tag, pageID, 2)
		}

		// Extract and index sections
		sections := g.extractSections(page.RawMD)
		for _, section := range sections {
			index.Sections = append(index.Sections, SearchSection{
				PageID: pageID,
				Title:  section.title,
				Anchor: section.anchor,
			})

			// Index section title (high weight)
			g.indexText(index.Index, section.title, pageID, 3)

			// Index section content (low weight)
			g.indexText(index.Index, section.content, pageID, 1)
		}

		pageID++
	}

	// Index FAQ items (if enabled)
	if g.site.FaqPage != nil && g.site.Config.Faq.Enabled {
		faqPath := g.site.Config.Faq.Path
		if faqPath == "" {
			faqPath = "faq"
		}

		for _, cat := range g.site.FaqPage.Categories {
			for _, item := range cat.Items {
				// Add FAQ item as a page
				index.Pages = append(index.Pages, SearchPage{
					Title: item.Question,
					Desc:  cat.Name,
					URL:   basePath + "/" + faqPath + "/#" + item.Slug,
				})

				// Index question (high weight)
				g.indexText(index.Index, item.Question, pageID, 3)

				// Index category (medium weight)
				g.indexText(index.Index, cat.Name, pageID, 2)

				// Index tags (medium weight)
				for _, tag := range item.Tags {
					g.indexText(index.Index, tag, pageID, 2)
				}

				// Index answer content (low weight)
				answerText := item.Answer
				if item.AnswerHTML != "" {
					answerText = extractPlainText([]byte(item.AnswerHTML))
				}
				g.indexText(index.Index, truncate(answerText, 1000), pageID, 1)

				pageID++
			}
		}
	}

	// Index KB articles (if enabled)
	if g.site.KBPage != nil && g.site.Config.KnowledgeBase.Enabled {
		kbPath := g.site.Config.KnowledgeBase.Path
		if kbPath == "" {
			kbPath = "kb"
		}

		for _, cat := range g.site.KBPage.Categories {
			for _, article := range cat.Articles {
				// Add KB article as a page
				index.Pages = append(index.Pages, SearchPage{
					Title: article.Title,
					Desc:  cat.Name + " | " + article.Description,
					URL:   basePath + "/" + kbPath + "/" + cat.Slug + "/" + article.Slug + ".html",
				})

				// Index title (high weight)
				g.indexText(index.Index, article.Title, pageID, 3)

				// Index category (medium weight)
				g.indexText(index.Index, cat.Name, pageID, 2)

				// Index description (medium weight)
				if article.Description != "" {
					g.indexText(index.Index, article.Description, pageID, 2)
				}

				// Index tags (medium weight)
				for _, tag := range article.Tags {
					g.indexText(index.Index, tag, pageID, 2)
				}

				// Index content (low weight)
				g.indexText(index.Index, truncate(extractPlainText([]byte(article.HTML)), 1000), pageID, 1)

				pageID++
			}
		}
	}

	// Sort posting lists by score (descending) for better results
	for word := range index.Index {
		list := index.Index[word]
		// Sort pairs by score (second element of each pair)
		sortPostingList(list)
		index.Index[word] = list
	}

	// Write index
	indexData, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal search index: %w", err)
	}

	indexPath := filepath.Join(g.site.OutputRoot, "search-index.json")
	if err := writeWebFile(indexPath, indexData); err != nil {
		return fmt.Errorf("failed to write search index: %w", err)
	}

	// Also write sharded index
	numShards, err := g.writeShardedIndex(index, g.site.OutputRoot)
	if err != nil {
		return fmt.Errorf("failed to write sharded index: %w", err)
	}

	fmt.Printf("Generated search index: %d pages, %d sections, %d terms, %d shards\n",
		len(index.Pages), len(index.Sections), len(index.Index), numShards)
	return nil
}

// sortPostingList sorts posting list pairs [pageID, score, ...] by score descending.
func sortPostingList(list PostingList) {
	n := len(list) / 2
	if n <= 1 {
		return
	}
	// Build an index slice, sort by score, then reconstruct list in sorted order.
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return list[indices[i]*2+1] > list[indices[j]*2+1]
	})
	tmp := make(PostingList, len(list))
	for i, idx := range indices {
		tmp[i*2] = list[idx*2]
		tmp[i*2+1] = list[idx*2+1]
	}
	copy(list, tmp)
}

// indexText tokenizes text and adds to inverted index
// Format: [pageID, score, pageID, score, ...]
func (g *SearchGenerator) indexText(idx map[string]PostingList, text string, pageID, weight int) {
	words := tokenize(text)
	seen := make(map[string]bool) // dedupe within same text block

	for _, word := range words {
		if seen[word] {
			continue
		}
		seen[word] = true

		// Skip stop words
		if stopWords[word] {
			continue
		}

		// Skip very short words
		if len(word) < 2 {
			continue
		}

		list := idx[word]

		// Check if this page already exists for this word
		found := false
		for i := 0; i < len(list); i += 2 {
			if list[i] == pageID {
				list[i+1] += weight // accumulate score
				found = true
				break
			}
		}

		if !found {
			idx[word] = append(list, pageID, weight)
		}
	}
}

// tokenize splits text into lowercase words
func tokenize(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else if current.Len() > 0 {
			words = append(words, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// sectionData holds extracted section data
type sectionData struct {
	title   string
	anchor  string
	content string
}

// extractSections extracts sections from markdown
func (g *SearchGenerator) extractSections(md []byte) []sectionData {
	lines := strings.Split(string(md), "\n")
	var sections []sectionData
	var currentHeading string
	var currentAnchor string
	var currentContent []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			// Save previous section
			if currentHeading != "" && len(currentContent) > 0 {
				sections = append(sections, sectionData{
					title:   currentHeading,
					anchor:  currentAnchor,
					content: truncate(extractPlainText([]byte(strings.Join(currentContent, "\n"))), 1000),
				})
			}

			// Start new section
			level := 0
			for i := 0; i < len(trimmed) && trimmed[i] == '#'; i++ {
				level++
			}
			currentHeading = strings.TrimSpace(trimmed[level:])
			currentAnchor = generateAnchor(currentHeading)
			currentContent = []string{}
		} else {
			currentContent = append(currentContent, line)
		}
	}

	// Save last section
	if currentHeading != "" && len(currentContent) > 0 {
		sections = append(sections, sectionData{
			title:   currentHeading,
			anchor:  currentAnchor,
			content: truncate(extractPlainText([]byte(strings.Join(currentContent, "\n"))), 1000),
		})
	}

	return sections
}

// extractPlainText removes markdown syntax and extracts plain text
func extractPlainText(md []byte) string {
	text := string(md)

	// Remove code block fences but keep content
	text = removeCodeBlockFences(text)

	// Remove admonition markers
	text = removeAdmonitions(text)

	// Remove inline code backticks
	text = strings.ReplaceAll(text, "`", "")

	// Remove heading markers
	text = strings.ReplaceAll(text, "#", "")

	// Remove emphasis markers
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "_", "")

	// Remove links (keep just the text)
	for strings.Contains(text, "[") && strings.Contains(text, "]") {
		start := strings.Index(text, "[")
		end := strings.Index(text[start:], "]")
		if end == -1 {
			break
		}
		end += start

		linkText := text[start+1 : end]

		if end+1 < len(text) && text[end+1] == '(' {
			urlEnd := strings.Index(text[end+1:], ")")
			if urlEnd != -1 {
				text = text[:start] + linkText + text[end+urlEnd+2:]
			} else {
				text = text[:start] + linkText + text[end+1:]
			}
		} else {
			text = text[:start] + linkText + text[end+1:]
		}
	}

	// Clean up whitespace
	words := strings.Fields(text)
	return strings.Join(words, " ")
}

// removeCodeBlockFences removes fence markers but keeps code content
func removeCodeBlockFences(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// removeAdmonitions removes admonition markers
func removeAdmonitions(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ":::") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// generateAnchor creates a URL-friendly anchor from heading text
// Matches Goldmark's auto-heading-id behavior
func generateAnchor(heading string) string {
	// Convert to lowercase
	anchor := strings.ToLower(heading)

	// Replace spaces with hyphens
	anchor = strings.ReplaceAll(anchor, " ", "-")

	// Remove special characters, keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range anchor {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// getBasePath extracts the path component from BaseURL for asset linking
// Examples:
//   - "https://example.com/docs/" → "/docs"
//   - "https://example.com/" → ""
//   - "" → ""
func (g *SearchGenerator) getBasePath() string {
	baseURL := g.site.Config.BaseURL
	if baseURL == "" {
		return ""
	}

	// Parse the URL to extract the path
	// Remove protocol and domain, keep only the path
	if strings.HasPrefix(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	} else if strings.HasPrefix(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	}

	// Find the first / after the domain
	parts := strings.SplitN(baseURL, "/", 2)
	if len(parts) < 2 {
		return ""
	}

	// Get the path part and ensure it starts with / and doesn't end with /
	path := "/" + parts[1]
	path = strings.TrimSuffix(path, "/")

	// If path is just "/", return empty string
	if path == "/" {
		return ""
	}

	return path
}

// writeShardedIndex writes the search index as sharded files
// Returns the number of shards created
func (g *SearchGenerator) writeShardedIndex(index SearchIndex, outputDir string) (int, error) {
	// Group terms by prefix (first 2 chars, or full term if < 2 chars)
	shardMap := make(map[string]map[string]PostingList)

	for term, postings := range index.Index {
		prefix := getShardPrefix(term)
		if shardMap[prefix] == nil {
			shardMap[prefix] = make(map[string]PostingList)
		}
		shardMap[prefix][term] = postings
	}

	// Collect shard prefixes
	shardPrefixes := make([]string, 0, len(shardMap))
	for prefix := range shardMap {
		shardPrefixes = append(shardPrefixes, prefix)
	}

	// Write manifest (pages + sections + shard list)
	manifest := SearchManifest{
		Pages:    index.Pages,
		Sections: index.Sections,
		Shards:   shardPrefixes,
	}

	manifestData, err := json.Marshal(manifest)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal search manifest: %w", err)
	}

	manifestPath := filepath.Join(outputDir, "search-manifest.json")
	if err := writeWebFile(manifestPath, manifestData); err != nil {
		return 0, fmt.Errorf("failed to write search manifest: %w", err)
	}

	// Create shards directory
	shardsDir := filepath.Join(outputDir, "search-shards")
	if err := makeWebDir(shardsDir); err != nil {
		return 0, fmt.Errorf("failed to create shards directory: %w", err)
	}

	// Write each shard
	for prefix, terms := range shardMap {
		shard := SearchShard{
			Prefix: prefix,
			Index:  terms,
		}

		shardData, err := json.Marshal(shard)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal shard %s: %w", prefix, err)
		}

		shardPath := filepath.Join(shardsDir, prefix+".json")
		if err := writeWebFile(shardPath, shardData); err != nil {
			return 0, fmt.Errorf("failed to write shard %s: %w", prefix, err)
		}
	}

	return len(shardMap), nil
}

// getShardPrefix returns the shard prefix for a term
// Uses first 2 characters, or the full term if shorter
func getShardPrefix(term string) string {
	if len(term) < 2 {
		return term
	}
	return term[:2]
}
