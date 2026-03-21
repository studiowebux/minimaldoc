// Package core provides consolidated default configuration functions.
// All feature-specific defaults are defined in their respective files
// (e.g., DefaultStatusConfig in status.go) and re-exported here for convenience.
package core

// DefaultConfigs returns a fully initialized SiteConfig with all defaults applied.
// Delegates to DefaultSiteConfig to ensure a single source of truth.
func DefaultConfigs() SiteConfig {
	return DefaultSiteConfig()
}

// DefaultConfigRegistry lists all default config functions for documentation
// and tooling purposes. Each entry maps a config name to its default function.
var DefaultConfigRegistry = map[string]func() any{
	"site":          func() any { return DefaultSiteConfig() },
	"openapi":       func() any { return DefaultOpenAPIConfig() },
	"status":        func() any { return DefaultStatusConfig() },
	"uptime":        func() any { return DefaultUptimeConfig() },
	"changelog":     func() any { return DefaultChangelogConfig() },
	"stale_warning": func() any { return DefaultStaleWarningConfig() },
	"landing":       func() any { return DefaultLandingConfig() },
	"theme":         func() any { return DefaultThemeConfig() },
	"portfolio":     func() any { return DefaultPortfolioConfig() },
	"contact":       func() any { return DefaultContactConfig() },
	"faq":           func() any { return DefaultFaqConfig() },
	"legal":         func() any { return DefaultLegalConfig() },
	"knowledgebase": func() any { return DefaultKBConfig() },
	"waitlist":      func() any { return DefaultWaitlistConfig() },
	"roadmap":       func() any { return DefaultRoadmapConfig() },
	"footer":        func() any { return DefaultFooterConfig() },
	"link_check":    func() any { return DefaultLinkCheckConfig() },
	"versions":      func() any { return DefaultVersionConfig() },
	"i18n":          func() any { return DefaultI18nConfig() },
	"pdf_export":    func() any { return DefaultPDFExportConfig() },
	"claude_assist": func() any { return DefaultClaudeAssistConfig() },
	"analytics":     func() any { return DefaultAnalyticsConfig() },
}
