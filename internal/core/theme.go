package core

// ThemeConfig holds theme customization settings
type ThemeConfig struct {
	Name   string      `yaml:"name"`   // Theme name: "default", "dark", or custom
	Colors ThemeColors `yaml:"colors"` // Custom color overrides
	Fonts  ThemeFonts  `yaml:"fonts"`  // Custom font settings
	Hero   ThemeHero   `yaml:"hero"`   // Hero section styling
}

// ThemeColors defines customizable color variables
type ThemeColors struct {
	// Light mode colors
	Light ThemeColorSet `yaml:"light"`
	// Dark mode colors
	Dark ThemeColorSet `yaml:"dark"`
}

// ThemeColorSet defines a complete set of theme colors
type ThemeColorSet struct {
	// Backgrounds
	BgPrimary   string `yaml:"bg_primary"`
	BgSecondary string `yaml:"bg_secondary"`
	BgTertiary  string `yaml:"bg_tertiary"`
	BgCode      string `yaml:"bg_code"`
	BgHover     string `yaml:"bg_hover"`

	// Text
	TextPrimary   string `yaml:"text_primary"`
	TextSecondary string `yaml:"text_secondary"`
	TextTertiary  string `yaml:"text_tertiary"`
	TextMuted     string `yaml:"text_muted"`

	// Borders
	BorderPrimary   string `yaml:"border_primary"`
	BorderSecondary string `yaml:"border_secondary"`

	// Accents
	AccentPrimary string `yaml:"accent_primary"`
	AccentHover   string `yaml:"accent_hover"`

	// Links
	LinkColor string `yaml:"link_color"`
	LinkHover string `yaml:"link_hover"`

	// Status colors
	ColorSuccess string `yaml:"color_success"`
	ColorWarning string `yaml:"color_warning"`
	ColorDanger  string `yaml:"color_danger"`
	ColorInfo    string `yaml:"color_info"`
}

// ThemeFonts defines custom font settings
type ThemeFonts struct {
	Heading   string `yaml:"heading"`    // Font family for headings
	Body      string `yaml:"body"`       // Font family for body text
	Code      string `yaml:"code"`       // Font family for code
	GoogleURL string `yaml:"google_url"` // Google Fonts import URL
}

// ThemeHero defines hero section styling
type ThemeHero struct {
	BackgroundImage   string `yaml:"background_image"`   // Background image URL
	BackgroundOverlay string `yaml:"background_overlay"` // Overlay color (e.g., "rgba(0,0,0,0.6)")
	TextAlign         string `yaml:"text_align"`         // Text alignment: left, center, right
	MinHeight         string `yaml:"min_height"`         // Minimum height (e.g., "80vh")
}

// DefaultThemeConfig returns a ThemeConfig with sensible defaults
func DefaultThemeConfig() ThemeConfig {
	return ThemeConfig{
		Name: "default",
		Colors: ThemeColors{
			Light: ThemeColorSet{
				BgPrimary:       "#fafafa",
				BgSecondary:     "#f5f5f5",
				BgTertiary:      "#eeeeee",
				BgCode:          "#f8f8f8",
				BgHover:         "#e8e8e8",
				TextPrimary:     "#1a1a1a",
				TextSecondary:   "#4a4a4a",
				TextTertiary:    "#6a6a6a",
				TextMuted:       "#9a9a9a",
				BorderPrimary:   "#e0e0e0",
				BorderSecondary: "#d0d0d0",
				AccentPrimary:   "#2a2a2a",
				AccentHover:     "#3a3a3a",
				LinkColor:       "#2563eb",
				LinkHover:       "#1d4ed8",
				ColorSuccess:    "#10b981",
				ColorWarning:    "#f59e0b",
				ColorDanger:     "#ef4444",
				ColorInfo:       "#3b82f6",
			},
			Dark: ThemeColorSet{
				BgPrimary:       "#1a1a1a",
				BgSecondary:     "#2a2a2a",
				BgTertiary:      "#333333",
				BgCode:          "#252525",
				BgHover:         "#3a3a3a",
				TextPrimary:     "#ffffff",
				TextSecondary:   "#e8e8e8",
				TextTertiary:    "#d0d0d0",
				TextMuted:       "#b0b0b0",
				BorderPrimary:   "#3a3a3a",
				BorderSecondary: "#4a4a4a",
				AccentPrimary:   "#e5e5e5",
				AccentHover:     "#f5f5f5",
				LinkColor:       "#7bb3ff",
				LinkHover:       "#a5cfff",
				ColorSuccess:    "#34d399",
				ColorWarning:    "#fbbf24",
				ColorDanger:     "#f87171",
				ColorInfo:       "#60a5fa",
			},
		},
		Fonts: ThemeFonts{
			Heading:   "",
			Body:      "",
			Code:      "",
			GoogleURL: "",
		},
		Hero: ThemeHero{
			BackgroundImage:   "",
			BackgroundOverlay: "",
			TextAlign:         "center",
			MinHeight:         "",
		},
	}
}

// HasCustomColors returns true if any custom colors are defined
func (t *ThemeConfig) HasCustomColors() bool {
	return t.Colors.Light.BgPrimary != "" || t.Colors.Dark.BgPrimary != ""
}

// HasCustomFonts returns true if any custom fonts are defined
func (t *ThemeConfig) HasCustomFonts() bool {
	return t.Fonts.Heading != "" || t.Fonts.Body != "" || t.Fonts.GoogleURL != ""
}

// HasHeroBackground returns true if hero background is configured
func (t *ThemeConfig) HasHeroBackground() bool {
	return t.Hero.BackgroundImage != ""
}
