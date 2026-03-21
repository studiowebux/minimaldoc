package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GitHubRelease represents the GitHub API response for a release
type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
	HTMLURL     string `json:"html_url"`
	Body        string `json:"body"`
}

// CheckResult contains the result of a version check
type CheckResult struct {
	CurrentVersion  string
	LatestVersion   string
	UpdateAvailable bool
	ReleaseURL      string
	PublishedAt     string
	ReleaseNotes    string
}

// CheckForUpdate checks if a newer version is available on GitHub
func CheckForUpdate() (*CheckResult, error) {
	const repoURL = "https://api.github.com/repos/studiowebux/minimaldoc/releases/latest"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", repoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "minimaldoc/"+Version)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(Version, "v")

	result := &CheckResult{
		CurrentVersion:  currentVersion,
		LatestVersion:   latestVersion,
		UpdateAvailable: compareVersions(currentVersion, latestVersion) < 0,
		ReleaseURL:      release.HTMLURL,
		PublishedAt:     release.PublishedAt,
		ReleaseNotes:    truncateNotes(release.Body, 500),
	}

	return result, nil
}

// compareVersions compares two semantic versions
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func compareVersions(a, b string) int {
	aParts := parseVersion(a)
	bParts := parseVersion(b)

	for i := range 3 {
		if aParts[i] < bParts[i] {
			return -1
		}
		if aParts[i] > bParts[i] {
			return 1
		}
	}
	return 0
}

// parseVersion parses a semantic version string into [major, minor, patch]
func parseVersion(v string) [3]int {
	parts := strings.Split(v, ".")
	result := [3]int{0, 0, 0}

	for i := 0; i < len(parts) && i < 3; i++ {
		// Handle pre-release suffixes (e.g., "1.0.0-beta")
		numStr := strings.Split(parts[i], "-")[0]
		var num int
		_, _ = fmt.Sscanf(numStr, "%d", &num) // non-numeric string defaults num to 0
		result[i] = num
	}

	return result
}

// truncateNotes truncates release notes to a maximum rune length. UTF-8 safe.
func truncateNotes(notes string, maxLen int) string {
	runes := []rune(notes)
	if len(runes) <= maxLen {
		return notes
	}
	return string(runes[:maxLen]) + "..."
}
