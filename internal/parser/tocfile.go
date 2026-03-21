package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TOCEntry represents an entry in the TOC.md file
type TOCEntry struct {
	Title      string
	FilePath   string
	Level      int
	IsExternal bool
	Children   []*TOCEntry
}

// TOCFileParser parses a TOC.md file to extract navigation structure
type TOCFileParser struct {
	docsDir string
}

// NewTOCFileParser creates a new TOC file parser
func NewTOCFileParser(docsDir string) *TOCFileParser {
	return &TOCFileParser{
		docsDir: docsDir,
	}
}

// Parse parses a TOC.md file and returns the navigation structure
func (p *TOCFileParser) Parse(tocFilePath string) ([]*TOCEntry, error) {
	file, err := os.Open(tocFilePath) // #nosec G304 -- path from trusted user configuration
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []*TOCEntry
	var stack []*TOCEntry
	scanner := bufio.NewScanner(file)

	// Regex to match markdown list items with optional links
	// Matches:  - [Title](path.md)
	//          - Title
	//        - [Title](path.md)
	listItemRegex := regexp.MustCompile(`^(\s*)[-*]\s+(?:\[([^\]]+)\]\(([^\)]+)\)|(.+))$`)

	// Detect indent unit from the first indented line (supports 2-space, 4-space, tabs)
	indentUnit := 0

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Skip empty lines and headings
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		matches := listItemRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		indent := len(matches[1])
		if indent > 0 && indentUnit == 0 {
			indentUnit = indent
		}
		level := 0
		if indentUnit > 0 {
			level = indent / indentUnit
		}

		var title, filePath string
		var isExternal bool
		if matches[2] != "" {
			// Link format: [Title](path)
			title = matches[2]
			filePath = matches[3]
			// Check if path is an external URL
			isExternal = strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://")
			// Normalize relative paths: strip ./ prefix
			if !isExternal {
				filePath = strings.TrimPrefix(filePath, "./")
			}
		} else {
			// Plain text format: Title (section header)
			title = strings.TrimSpace(matches[4])
			filePath = ""
			isExternal = false
		}

		entry := &TOCEntry{
			Title:      title,
			FilePath:   filePath,
			Level:      level,
			IsExternal: isExternal,
			Children:   []*TOCEntry{},
		}

		// Build the hierarchy
		for len(stack) > 0 && stack[len(stack)-1].Level >= level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// Top-level entry
			entries = append(entries, entry)
		} else {
			// Child entry
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, entry)
		}

		stack = append(stack, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// Flatten returns a flat list of all file paths in order
func (p *TOCFileParser) Flatten(entries []*TOCEntry) []string {
	var result []string
	var flatten func([]*TOCEntry)

	flatten = func(entries []*TOCEntry) {
		for _, entry := range entries {
			if entry.FilePath != "" && !entry.IsExternal {
				// Resolve relative path and validate against traversal
				fullPath := filepath.Join(p.docsDir, entry.FilePath)
				absPath, err := filepath.Abs(fullPath)
				if err != nil {
					continue
				}
				absRoot, err := filepath.Abs(p.docsDir)
				if err != nil {
					continue
				}
				if !strings.HasPrefix(absPath, absRoot+string(filepath.Separator)) && absPath != absRoot {
					fmt.Fprintf(os.Stderr, "Warning: TOC path %q escapes docs directory, skipping\n", entry.FilePath)
					continue
				}
				result = append(result, fullPath)
			}
			if len(entry.Children) > 0 {
				flatten(entry.Children)
			}
		}
	}

	flatten(entries)
	return result
}
