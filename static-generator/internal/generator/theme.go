package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
)

// ThemeGenerator generates custom theme CSS
type ThemeGenerator struct {
	site *core.Site
}

// NewThemeGenerator creates a new theme generator
func NewThemeGenerator(site *core.Site) *ThemeGenerator {
	return &ThemeGenerator{site: site}
}

// Generate creates custom-theme.css if custom theme settings are configured
func (g *ThemeGenerator) Generate() error {
	cfg := g.site.Config.ThemeConfig

	// Skip if no custom theme configuration
	if !cfg.HasCustomColors() && !cfg.HasCustomFonts() && !cfg.HasHeroBackground() {
		return nil
	}

	fmt.Println("Generating custom theme CSS...")

	var css strings.Builder

	// Generate Google Fonts import if configured
	if cfg.Fonts.GoogleURL != "" {
		css.WriteString(fmt.Sprintf("@import url(\"%s\");\n\n", cfg.Fonts.GoogleURL))
	}

	// Generate light mode variables
	css.WriteString(":root[data-theme=\"light\"] {\n")
	g.writeColorVariables(&css, cfg.Colors.Light)
	css.WriteString("}\n\n")

	// Generate dark mode variables
	css.WriteString(":root[data-theme=\"dark\"] {\n")
	g.writeColorVariables(&css, cfg.Colors.Dark)
	css.WriteString("}\n\n")

	// Generate font variables if configured
	if cfg.HasCustomFonts() {
		css.WriteString(":root {\n")
		if cfg.Fonts.Heading != "" {
			css.WriteString(fmt.Sprintf("    --font-heading: %s;\n", cfg.Fonts.Heading))
		}
		if cfg.Fonts.Body != "" {
			css.WriteString(fmt.Sprintf("    --font-body: %s;\n", cfg.Fonts.Body))
		}
		if cfg.Fonts.Code != "" {
			css.WriteString(fmt.Sprintf("    --font-code: %s;\n", cfg.Fonts.Code))
		}
		css.WriteString("}\n\n")

		// Apply font families to elements
		if cfg.Fonts.Body != "" {
			css.WriteString("body {\n")
			css.WriteString("    font-family: var(--font-body), -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial, sans-serif;\n")
			css.WriteString("}\n\n")
		}

		if cfg.Fonts.Heading != "" {
			css.WriteString(".markdown-body h1,\n")
			css.WriteString(".markdown-body h2,\n")
			css.WriteString(".markdown-body h3,\n")
			css.WriteString(".markdown-body h4,\n")
			css.WriteString(".markdown-body h5,\n")
			css.WriteString(".markdown-body h6,\n")
			css.WriteString(".hero-title,\n")
			css.WriteString(".section-title,\n")
			css.WriteString(".feature-title,\n")
			css.WriteString(".cta-title {\n")
			css.WriteString("    font-family: var(--font-heading), -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial, sans-serif;\n")
			css.WriteString("}\n\n")
		}

		if cfg.Fonts.Code != "" {
			css.WriteString(".markdown-body code,\n")
			css.WriteString(".markdown-body pre,\n")
			css.WriteString(".chroma {\n")
			css.WriteString("    font-family: var(--font-code), \"Consolas\", \"Monaco\", \"Courier New\", monospace;\n")
			css.WriteString("}\n\n")
		}
	}

	// Generate hero background styles if configured
	if cfg.HasHeroBackground() {
		css.WriteString(".hero-section {\n")
		if cfg.Hero.BackgroundImage != "" {
			css.WriteString(fmt.Sprintf("    background-image: url(\"%s\");\n", cfg.Hero.BackgroundImage))
			css.WriteString("    background-size: cover;\n")
			css.WriteString("    background-position: center;\n")
			css.WriteString("    background-repeat: no-repeat;\n")
			css.WriteString("    position: relative;\n")
		}
		if cfg.Hero.MinHeight != "" {
			css.WriteString(fmt.Sprintf("    min-height: %s;\n", cfg.Hero.MinHeight))
		}
		if cfg.Hero.TextAlign != "" && cfg.Hero.TextAlign != "center" {
			css.WriteString(fmt.Sprintf("    text-align: %s;\n", cfg.Hero.TextAlign))
		}
		css.WriteString("}\n\n")

		// Add overlay if configured
		if cfg.Hero.BackgroundOverlay != "" {
			css.WriteString(".hero-section::before {\n")
			css.WriteString("    content: \"\";\n")
			css.WriteString("    position: absolute;\n")
			css.WriteString("    top: 0;\n")
			css.WriteString("    left: 0;\n")
			css.WriteString("    right: 0;\n")
			css.WriteString("    bottom: 0;\n")
			css.WriteString(fmt.Sprintf("    background: %s;\n", cfg.Hero.BackgroundOverlay))
			css.WriteString("    z-index: 0;\n")
			css.WriteString("}\n\n")

			css.WriteString(".hero-content {\n")
			css.WriteString("    position: relative;\n")
			css.WriteString("    z-index: 1;\n")
			css.WriteString("}\n\n")
		}
	}

	// Write to file
	cssPath := filepath.Join(g.site.OutputRoot, "css", "custom-theme.css")

	// Ensure css directory exists
	cssDir := filepath.Dir(cssPath)
	if err := makeWebDir(cssDir); err != nil {
		return fmt.Errorf("failed to create css directory: %w", err)
	}

	if err := writeWebFile(cssPath, []byte(css.String())); err != nil {
		return fmt.Errorf("failed to write custom theme CSS: %w", err)
	}

	fmt.Println("Generated css/custom-theme.css")
	return nil
}

// writeColorVariables writes CSS color variables for a color set
func (g *ThemeGenerator) writeColorVariables(css *strings.Builder, colors core.ThemeColorSet) {
	// Only write non-empty values
	if colors.BgPrimary != "" {
		css.WriteString(fmt.Sprintf("    --bg-primary: %s;\n", colors.BgPrimary))
	}
	if colors.BgSecondary != "" {
		css.WriteString(fmt.Sprintf("    --bg-secondary: %s;\n", colors.BgSecondary))
	}
	if colors.BgTertiary != "" {
		css.WriteString(fmt.Sprintf("    --bg-tertiary: %s;\n", colors.BgTertiary))
	}
	if colors.BgCode != "" {
		css.WriteString(fmt.Sprintf("    --bg-code: %s;\n", colors.BgCode))
	}
	if colors.BgHover != "" {
		css.WriteString(fmt.Sprintf("    --bg-hover: %s;\n", colors.BgHover))
	}
	if colors.TextPrimary != "" {
		css.WriteString(fmt.Sprintf("    --text-primary: %s;\n", colors.TextPrimary))
	}
	if colors.TextSecondary != "" {
		css.WriteString(fmt.Sprintf("    --text-secondary: %s;\n", colors.TextSecondary))
	}
	if colors.TextTertiary != "" {
		css.WriteString(fmt.Sprintf("    --text-tertiary: %s;\n", colors.TextTertiary))
	}
	if colors.TextMuted != "" {
		css.WriteString(fmt.Sprintf("    --text-muted: %s;\n", colors.TextMuted))
	}
	if colors.BorderPrimary != "" {
		css.WriteString(fmt.Sprintf("    --border-primary: %s;\n", colors.BorderPrimary))
	}
	if colors.BorderSecondary != "" {
		css.WriteString(fmt.Sprintf("    --border-secondary: %s;\n", colors.BorderSecondary))
	}
	if colors.AccentPrimary != "" {
		css.WriteString(fmt.Sprintf("    --accent-primary: %s;\n", colors.AccentPrimary))
	}
	if colors.AccentHover != "" {
		css.WriteString(fmt.Sprintf("    --accent-hover: %s;\n", colors.AccentHover))
	}
	if colors.LinkColor != "" {
		css.WriteString(fmt.Sprintf("    --link-color: %s;\n", colors.LinkColor))
	}
	if colors.LinkHover != "" {
		css.WriteString(fmt.Sprintf("    --link-hover: %s;\n", colors.LinkHover))
	}
	if colors.ColorSuccess != "" {
		css.WriteString(fmt.Sprintf("    --color-success: %s;\n", colors.ColorSuccess))
	}
	if colors.ColorWarning != "" {
		css.WriteString(fmt.Sprintf("    --color-warning: %s;\n", colors.ColorWarning))
	}
	if colors.ColorDanger != "" {
		css.WriteString(fmt.Sprintf("    --color-danger: %s;\n", colors.ColorDanger))
	}
	if colors.ColorInfo != "" {
		css.WriteString(fmt.Sprintf("    --color-info: %s;\n", colors.ColorInfo))
	}
}

// HasCustomTheme returns true if custom theme is configured
func (g *ThemeGenerator) HasCustomTheme() bool {
	cfg := g.site.Config.ThemeConfig
	return cfg.HasCustomColors() || cfg.HasCustomFonts() || cfg.HasHeroBackground()
}
