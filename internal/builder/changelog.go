package builder

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/studiowebux/minimaldoc/internal/parser"
)

// ChangelogBuilder handles building the changelog
type ChangelogBuilder struct {
	parser *parser.ChangelogParser
}

// NewChangelogBuilder creates a new changelog builder
func NewChangelogBuilder() *ChangelogBuilder {
	return &ChangelogBuilder{
		parser: parser.NewChangelogParser(),
	}
}

// Build parses and builds the changelog data
func (cb *ChangelogBuilder) Build(docsRoot string, config core.ChangelogConfig) (*core.ChangelogPage, error) {
	changelogDir := filepath.Join(docsRoot, parser.ChangelogSourceDir)

	// Parse all changelog content
	changelogPage, err := cb.parser.ParseChangelogDir(changelogDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse changelog directory: %w", err)
	}

	// Merge config from site config (file config overrides defaults, site config takes precedence)
	changelogPage.Config = cb.mergeConfig(changelogPage.Config, config)

	// Sort releases by semantic version (newest first)
	cb.sortReleases(changelogPage)

	// Generate compare URLs
	cb.generateCompareURLs(changelogPage)

	// Set output paths for releases
	for i := range changelogPage.Releases {
		changelogPage.Releases[i].OutputPath = cb.getReleaseOutputPath(
			changelogPage.Config.Path,
			changelogPage.Releases[i].Slug,
		)
	}

	// Update last updated time
	changelogPage.LastUpdated = time.Now()

	return changelogPage, nil
}

// mergeConfig merges two changelog configs, with override taking precedence for non-zero values
func (cb *ChangelogBuilder) mergeConfig(base, override core.ChangelogConfig) core.ChangelogConfig {
	result := base

	// Override enabled state
	if override.Enabled {
		result.Enabled = override.Enabled
	}

	// Override title if set
	if override.Title != "" && override.Title != "Changelog" {
		result.Title = override.Title
	}

	// Override description if set
	if override.Description != "" && override.Description != "All notable changes to this project" {
		result.Description = override.Description
	}

	// Override path if set
	if override.Path != "" && override.Path != "changelog" {
		result.Path = override.Path
	}

	// Override repository if set
	if override.Repository != "" {
		result.Repository = override.Repository
	}

	// RSS enabled can be explicitly disabled
	if !override.RSSEnabled && base.RSSEnabled {
		// Only disable if explicitly set in override
		result.RSSEnabled = base.RSSEnabled
	} else {
		result.RSSEnabled = override.RSSEnabled
	}

	return result
}

// sortReleases sorts releases by semantic version (newest first)
func (cb *ChangelogBuilder) sortReleases(page *core.ChangelogPage) {
	sort.Slice(page.Releases, func(i, j int) bool {
		// Compare major version
		if page.Releases[i].Major != page.Releases[j].Major {
			return page.Releases[i].Major > page.Releases[j].Major
		}
		// Compare minor version
		if page.Releases[i].Minor != page.Releases[j].Minor {
			return page.Releases[i].Minor > page.Releases[j].Minor
		}
		// Compare patch version
		if page.Releases[i].Patch != page.Releases[j].Patch {
			return page.Releases[i].Patch > page.Releases[j].Patch
		}
		// If versions are equal, sort by date (newer first)
		return page.Releases[i].Date.After(page.Releases[j].Date)
	})
}

// generateCompareURLs creates GitHub compare links between consecutive versions
func (cb *ChangelogBuilder) generateCompareURLs(page *core.ChangelogPage) {
	if page.Config.Repository == "" {
		return
	}

	// Clean up repository URL
	repoURL := strings.TrimSuffix(page.Config.Repository, "/")
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// Generate compare URLs (sorted newest first, so compare current with previous)
	for i := 0; i < len(page.Releases); i++ {
		if i < len(page.Releases)-1 {
			// Compare with previous (older) version
			prevVersion := page.Releases[i+1].Version
			currVersion := page.Releases[i].Version

			// Add 'v' prefix if not present for tag names
			prevTag := prevVersion
			if !strings.HasPrefix(prevTag, "v") {
				prevTag = "v" + prevTag
			}
			currTag := currVersion
			if !strings.HasPrefix(currTag, "v") {
				currTag = "v" + currTag
			}

			page.Releases[i].CompareURL = fmt.Sprintf("%s/compare/%s...%s", repoURL, prevTag, currTag)
		}
	}
}

// getReleaseOutputPath generates the output path for a release
func (cb *ChangelogBuilder) getReleaseOutputPath(changelogPath, slug string) string {
	return filepath.Join(changelogPath, slug+".html")
}
