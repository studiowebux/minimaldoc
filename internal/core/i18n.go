package core

// I18nConfig holds the configuration for internationalization
type I18nConfig struct {
	Enabled           bool           `yaml:"enabled"`
	DefaultLocale     string         `yaml:"default_locale"`      // Default locale code (e.g., "en")
	Locales           []LocaleInfo   `yaml:"locales"`             // Available locales
	Fallback          string         `yaml:"fallback"`            // Fallback locale for missing translations
	HideDefaultLocale bool           `yaml:"hide_default_locale"` // Use /page.html instead of /en/page.html for default
	ShowUntranslated  bool           `yaml:"show_untranslated"`   // Show fallback content with warning badge
	Selector          LocaleSelector `yaml:"selector"`            // UI selector configuration
}

// LocaleInfo represents a single locale
type LocaleInfo struct {
	Code      string `yaml:"code"`      // Locale code (e.g., "en", "fr", "de")
	Name      string `yaml:"name"`      // Display name (e.g., "English", "Francais")
	Direction string `yaml:"direction"` // Text direction: "ltr" or "rtl"
	Flag      string `yaml:"flag"`      // Optional flag emoji or icon
}

// LocaleSelector configures the locale selector UI
type LocaleSelector struct {
	Position  string `yaml:"position"`   // "header" or "sidebar"
	ShowFlags bool   `yaml:"show_flags"` // Show flag icons/emojis
}

// LocalizedSite represents a site built for a specific locale
type LocalizedSite struct {
	Locale     LocaleInfo // Current locale being built
	Pages      []*Page    // Pages for this locale
	OutputRoot string     // Output directory for this locale
}

// DefaultI18nConfig returns an I18nConfig with sensible defaults
func DefaultI18nConfig() I18nConfig {
	return I18nConfig{
		Enabled:           false,
		DefaultLocale:     "en",
		Locales:           []LocaleInfo{},
		Fallback:          "en",
		HideDefaultLocale: true,
		ShowUntranslated:  true,
		Selector: LocaleSelector{
			Position:  "header",
			ShowFlags: false,
		},
	}
}

// IsRTL returns true if the locale uses right-to-left text direction
func (l *LocaleInfo) IsRTL() bool {
	return l.Direction == "rtl"
}

// GetDirection returns the text direction, defaulting to "ltr"
func (l *LocaleInfo) GetDirection() string {
	if l.Direction == "" {
		return "ltr"
	}
	return l.Direction
}

// GetLocale returns the LocaleInfo for a given locale code
func (c *I18nConfig) GetLocale(code string) *LocaleInfo {
	for i := range c.Locales {
		if c.Locales[i].Code == code {
			return &c.Locales[i]
		}
	}
	return nil
}

// GetDefaultLocale returns the default LocaleInfo
func (c *I18nConfig) GetDefaultLocale() *LocaleInfo {
	if c.DefaultLocale == "" && len(c.Locales) > 0 {
		return &c.Locales[0]
	}
	return c.GetLocale(c.DefaultLocale)
}

// GetFallbackLocale returns the fallback LocaleInfo
func (c *I18nConfig) GetFallbackLocale() *LocaleInfo {
	if c.Fallback == "" {
		return c.GetDefaultLocale()
	}
	return c.GetLocale(c.Fallback)
}

// I18nSourceDir is the directory name for locale-specific translations
const I18nSourceDir = "__translations__"
