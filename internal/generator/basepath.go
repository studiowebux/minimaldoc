package generator

import "strings"

// featurePath returns configPath if non-empty, otherwise defaultPath.
// Replaces the repeated pattern: if path == "" { path = "default" }
func featurePath(configPath, defaultPath string) string {
	if configPath == "" {
		return defaultPath
	}
	return configPath
}

// trimBaseURL strips the trailing slash from a base URL for concatenation.
// Replaces the repeated pattern: strings.TrimSuffix(g.site.Config.BaseURL, "/")
func trimBaseURL(baseURL string) string {
	return strings.TrimSuffix(baseURL, "/")
}

// GetBasePath extracts the path component from a BaseURL for asset linking.
// Examples:
//   - "https://example.com/docs/" → "/docs"
//   - "https://example.com/" → ""
//   - "" → ""
func GetBasePath(baseURL string) string {
	if baseURL == "" {
		return ""
	}

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
