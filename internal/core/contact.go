package core

// ContactPage represents the complete contact page data
type ContactPage struct {
	Config ContactConfig
}

// ContactConfig holds configuration for the contact page
type ContactConfig struct {
	Enabled     bool          `yaml:"enabled"`
	Title       string        `yaml:"title"`
	Description string        `yaml:"description"`
	Path        string        `yaml:"path"`
	Email       string        `yaml:"email"`
	Info        []ContactInfo `yaml:"info"`
}

// DefaultContactConfig returns a ContactConfig with sensible defaults
func DefaultContactConfig() ContactConfig {
	return ContactConfig{
		Enabled:     false,
		Title:       "Contact",
		Description: "Get in touch",
		Path:        "contact",
	}
}

// ContactInfo represents a contact information item
type ContactInfo struct {
	Icon string `yaml:"icon"`
	Text string `yaml:"text"`
}
