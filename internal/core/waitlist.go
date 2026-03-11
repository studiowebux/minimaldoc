package core

// WaitlistSocialLink represents a social link on the waitlist page
type WaitlistSocialLink struct {
	Name string `yaml:"name"` // Platform name (github, twitter, discord, etc.)
	URL  string `yaml:"url"`  // Link URL
}

// WaitlistConfig holds configuration for the waitlist landing page
type WaitlistConfig struct {
	Enabled            bool                 `yaml:"enabled"`
	Title              string               `yaml:"title"`
	Tagline            string               `yaml:"tagline"`
	NewsletterEndpoint string               `yaml:"newsletter_endpoint"`
	SiteID             string               `yaml:"site_id"`
	SuccessMessage     string               `yaml:"success_message"`
	SocialLinks        []WaitlistSocialLink `yaml:"social_links"`
	PrivacyURL         string               `yaml:"privacy_url"`
}

// WaitlistPage represents the waitlist page data passed to the template
type WaitlistPage struct {
	Config WaitlistConfig
}

// DefaultWaitlistConfig returns a WaitlistConfig with sensible defaults
func DefaultWaitlistConfig() WaitlistConfig {
	return WaitlistConfig{
		Enabled:        false,
		Title:          "Coming Soon",
		Tagline:        "Sign up to be notified when we launch.",
		SuccessMessage: "You're on the list.",
	}
}
