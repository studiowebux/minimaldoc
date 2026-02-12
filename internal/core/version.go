package core

// VersionConfig holds the configuration for multi-version documentation
type VersionConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Default  string        `yaml:"default"`   // Default version slug (e.g., "v2")
	List     []VersionInfo `yaml:"list"`      // Available versions
	Selector VersionSelector `yaml:"selector"` // UI selector configuration
}

// VersionInfo represents a single documentation version
type VersionInfo struct {
	Name  string `yaml:"name"`  // Version identifier (e.g., "v2")
	Label string `yaml:"label"` // Display label (e.g., "2.x (Latest)")
	Path  string `yaml:"path"`  // URL path prefix (e.g., "v2")
	EOL   string `yaml:"eol"`   // End of life date (YYYY-MM-DD), empty if still supported
}

// VersionSelector configures the version selector UI
type VersionSelector struct {
	Position       string `yaml:"position"`         // "header" or "sidebar"
	ShowEOLWarning bool   `yaml:"show_eol_warning"` // Show warning badge for EOL versions
}

// VersionedSite represents a site built for a specific version
type VersionedSite struct {
	Version    VersionInfo // Current version being built
	Pages      []*Page     // Pages for this version
	OutputRoot string      // Output directory for this version
}

// DefaultVersionConfig returns a VersionConfig with sensible defaults
func DefaultVersionConfig() VersionConfig {
	return VersionConfig{
		Enabled: false,
		Default: "",
		List:    []VersionInfo{},
		Selector: VersionSelector{
			Position:       "header",
			ShowEOLWarning: true,
		},
	}
}

// IsEOL returns true if the version has reached end of life
func (v *VersionInfo) IsEOL() bool {
	return v.EOL != ""
}

// GetVersion returns the VersionInfo for a given version name
func (c *VersionConfig) GetVersion(name string) *VersionInfo {
	for i := range c.List {
		if c.List[i].Name == name {
			return &c.List[i]
		}
	}
	return nil
}

// GetDefaultVersion returns the default VersionInfo
func (c *VersionConfig) GetDefaultVersion() *VersionInfo {
	if c.Default == "" && len(c.List) > 0 {
		return &c.List[0]
	}
	return c.GetVersion(c.Default)
}

// VersionSourceDir is the directory name for version-specific content
const VersionSourceDir = "_versions"
