package core

// SectionBackground holds background styling for any section
type SectionBackground struct {
	Image       string `yaml:"image"`        // Background image URL
	Overlay     string `yaml:"overlay"`      // Overlay color (e.g., "rgba(0,0,0,0.6)")
	Color       string `yaml:"color"`        // Background color
	Position    string `yaml:"position"`     // Background position (default: "center")
	Size        string `yaml:"size"`         // Background size (default: "cover")
	Attachment  string `yaml:"attachment"`   // Background attachment (scroll/fixed)
}

// HasBackground returns true if any background is configured
func (s *SectionBackground) HasBackground() bool {
	return s.Image != "" || s.Color != ""
}

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
	ImageText    []*ImageTextSection
	TextBlocks   []*TextSection
	LinksGrid    []*LinksGridSection
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
	ImageText    []ImageTextSection  `yaml:"image_text"`
	TextBlocks   []TextSection       `yaml:"text_blocks"`
	LinksGrid    []LinksGridSection  `yaml:"links_grid"`
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
	Title      string            `yaml:"title"`
	Subtitle   string            `yaml:"subtitle"`
	Buttons    []HeroButton      `yaml:"buttons"`
	Image      string            `yaml:"image"`
	Background SectionBackground `yaml:"background"`
}

// HeroButton represents a CTA button in the hero section
type HeroButton struct {
	Text    string `yaml:"text"`
	URL     string `yaml:"url"`
	Primary bool   `yaml:"primary"`
}

// FeaturesSection represents the features grid section
type FeaturesSection struct {
	Title      string            `yaml:"title"`
	Items      []FeatureItem     `yaml:"items"`
	Background SectionBackground `yaml:"background"`
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
	Title      string            `yaml:"title"`
	Items      []StepItem        `yaml:"items"`
	Background SectionBackground `yaml:"background"`
}

// StepItem represents a single step
type StepItem struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Code        string `yaml:"code"`
}

// CTASection represents a call-to-action section
type CTASection struct {
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	Buttons     []HeroButton      `yaml:"buttons"`
	Background  SectionBackground `yaml:"background"`
}

// TestimonialsSection represents the testimonials section
type TestimonialsSection struct {
	Title      string            `yaml:"title"`
	Items      []Testimonial     `yaml:"items"`
	Background SectionBackground `yaml:"background"`
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
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	Links       []SimpleLink      `yaml:"links"`
	Background  SectionBackground `yaml:"background"`
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

// ImageTextSection represents an image + content side-by-side section
type ImageTextSection struct {
	ID            string            `yaml:"id"`             // Section ID for ordering
	Title         string            `yaml:"title"`
	Description   string            `yaml:"description"`
	Content       string            `yaml:"content"`        // Markdown content
	Image         string            `yaml:"image"`          // Image URL
	ImageAlt      string            `yaml:"image_alt"`      // Image alt text
	ImagePosition string            `yaml:"image_position"` // left or right (default: right)
	Items         []ImageTextItem   `yaml:"items"`          // Optional bullet points
	Buttons       []HeroButton      `yaml:"buttons"`        // Optional CTA buttons
	Background    SectionBackground `yaml:"background"`
	Order         int               `yaml:"order"`          // Display order
}

// ImageTextItem represents a bullet point in an image-text section
type ImageTextItem struct {
	Icon        string `yaml:"icon"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// TextSection represents a simple text block section
type TextSection struct {
	ID          string            `yaml:"id"`          // Section ID for ordering
	Title       string            `yaml:"title"`
	Subtitle    string            `yaml:"subtitle"`
	Content     string            `yaml:"content"`     // Markdown content
	Alignment   string            `yaml:"alignment"`   // left, center, right (default: center)
	MaxWidth    string            `yaml:"max_width"`   // Max width (e.g., "800px")
	Buttons     []HeroButton      `yaml:"buttons"`     // Optional CTA buttons
	Background  SectionBackground `yaml:"background"`
	Order       int               `yaml:"order"`       // Display order
}

// LinksGridSection represents a grid of external link cards
type LinksGridSection struct {
	ID          string            `yaml:"id"`          // Section ID for ordering
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	Columns     int               `yaml:"columns"`     // Number of columns (default: 4)
	Items       []LinksGridItem   `yaml:"items"`
	Background  SectionBackground `yaml:"background"`
	Order       int               `yaml:"order"`       // Display order
}

// LinksGridItem represents a single card in the links grid
type LinksGridItem struct {
	Icon        string `yaml:"icon"`        // Icon identifier or emoji
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
	External    bool   `yaml:"external"`    // Opens in new tab (default: true for http/https)
}
