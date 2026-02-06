package core

// LandingPage represents the complete landing page data
type LandingPage struct {
	Config       LandingConfig
	Hero         *HeroSection
	Features     *FeaturesSection
	Steps        *StepsSection
	CTA          *CTASection
	Testimonials *TestimonialsSection
	OpenSource   *OpenSourceSection
	Links        *LinksSection
}

// LandingConfig holds configuration for the landing page
type LandingConfig struct {
	Enabled      bool                `yaml:"enabled"`
	Nav          []LandingNavLink    `yaml:"nav"`
	Hero         HeroSection         `yaml:"hero"`
	Features     FeaturesSection     `yaml:"features"`
	Steps        StepsSection        `yaml:"steps"`
	CTA          CTASection          `yaml:"cta"`
	Testimonials TestimonialsSection `yaml:"testimonials"`
	OpenSource   OpenSourceSection   `yaml:"opensource"`
	Links        LinksSection        `yaml:"links"`
}

// LandingNavLink represents a navigation link in the landing header
type LandingNavLink struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}

// DefaultLandingConfig returns a LandingConfig with sensible defaults
func DefaultLandingConfig() LandingConfig {
	return LandingConfig{
		Enabled: false,
	}
}

// HeroSection represents the hero/header section of landing page
type HeroSection struct {
	Title    string       `yaml:"title"`
	Subtitle string       `yaml:"subtitle"`
	Buttons  []HeroButton `yaml:"buttons"`
	Image    string       `yaml:"image"`
}

// HeroButton represents a CTA button in the hero section
type HeroButton struct {
	Text    string `yaml:"text"`
	URL     string `yaml:"url"`
	Primary bool   `yaml:"primary"`
}

// FeaturesSection represents the features grid section
type FeaturesSection struct {
	Title string        `yaml:"title"`
	Items []FeatureItem `yaml:"items"`
}

// FeatureItem represents a single feature in the grid
type FeatureItem struct {
	Icon        string `yaml:"icon"`
	Emoji       string `yaml:"emoji"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// StepsSection represents the how-to/quick start section
type StepsSection struct {
	Title string     `yaml:"title"`
	Items []StepItem `yaml:"items"`
}

// StepItem represents a single step
type StepItem struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Code        string `yaml:"code"`
}

// CTASection represents a call-to-action section
type CTASection struct {
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Buttons     []HeroButton `yaml:"buttons"`
}

// TestimonialsSection represents the testimonials section
type TestimonialsSection struct {
	Title string        `yaml:"title"`
	Items []Testimonial `yaml:"items"`
}

// Testimonial represents a single testimonial
type Testimonial struct {
	Quote  string `yaml:"quote"`
	Author string `yaml:"author"`
	Role   string `yaml:"role"`
	Avatar string `yaml:"avatar"`
}

// OpenSourceSection represents the open source/credits section
type OpenSourceSection struct {
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Links       []SimpleLink `yaml:"links"`
}

// SimpleLink represents a simple text/url link
type SimpleLink struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}

// LinksSection represents a grid of resource links
type LinksSection struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	Items       []LinkItem `yaml:"items"`
}

// LinkItem represents a single link card
type LinkItem struct {
	Icon        string `yaml:"icon"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}
