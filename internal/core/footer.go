package core

// FooterConfig holds configuration for the landing page footer
type FooterConfig struct {
	Copyright string            `yaml:"copyright"`
	Links     []FooterLinkGroup `yaml:"links"`
	Social    []SocialLink      `yaml:"social"`
	Badges    []FooterBadge     `yaml:"badges"`
}

// DefaultFooterConfig returns a FooterConfig with sensible defaults
func DefaultFooterConfig() FooterConfig {
	return FooterConfig{}
}

// FooterLinkGroup represents a group of footer links
type FooterLinkGroup struct {
	Title string       `yaml:"title"`
	Items []FooterLink `yaml:"items"`
}

// FooterLink represents a single footer link
type FooterLink struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}

// FooterBadge represents a badge or powered-by link
type FooterBadge struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}
