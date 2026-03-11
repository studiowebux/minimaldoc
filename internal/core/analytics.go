package core

import "strings"

// AnalyticsConfig holds configuration for the analytics system
type AnalyticsConfig struct {
	Enabled   bool                `yaml:"enabled"`
	Providers []AnalyticsProvider `yaml:"providers"`
}

// AnalyticsProvider represents a single analytics provider configuration
type AnalyticsProvider struct {
	Type    string         `yaml:"type"`    // ga4, plausible, umami, matomo, fathom, simple, custom
	Enabled bool           `yaml:"enabled"`
	Config  map[string]any `yaml:"config"` // Provider-specific configuration
}

// DefaultAnalyticsConfig returns an AnalyticsConfig with sensible defaults
func DefaultAnalyticsConfig() AnalyticsConfig {
	return AnalyticsConfig{
		Enabled:   false,
		Providers: []AnalyticsProvider{},
	}
}

// GetEnabledProviders returns only the enabled providers
func (c AnalyticsConfig) GetEnabledProviders() []AnalyticsProvider {
	if !c.Enabled {
		return nil
	}
	var enabled []AnalyticsProvider
	for _, p := range c.Providers {
		if p.Enabled {
			enabled = append(enabled, p)
		}
	}
	return enabled
}

// GetConfigString retrieves a string value from provider config
func (p AnalyticsProvider) GetConfigString(key string) string {
	if p.Config == nil {
		return ""
	}
	if v, ok := p.Config[key].(string); ok {
		return v
	}
	return ""
}

// GetConfigMap retrieves a map value from provider config (for attrs)
func (p AnalyticsProvider) GetConfigMap(key string) map[string]string {
	if p.Config == nil {
		return nil
	}
	result := make(map[string]string)
	if v, ok := p.Config[key].(map[string]any); ok {
		for k, val := range v {
			if s, ok := val.(string); ok {
				result[k] = s
			}
		}
	}
	return result
}

// GetConfigSlice retrieves a slice value from provider config
func (p AnalyticsProvider) GetConfigSlice(key string) []string {
	if p.Config == nil {
		return nil
	}
	if v, ok := p.Config[key].([]any); ok {
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// GetMinimalDocProvider returns the minimaldoc provider if enabled, nil otherwise
func (c AnalyticsConfig) GetMinimalDocProvider() *AnalyticsProvider {
	if !c.Enabled {
		return nil
	}
	for _, p := range c.Providers {
		if p.Enabled && p.Type == "minimaldoc" {
			return &p
		}
	}
	return nil
}

// HasMinimalDocFeature checks if minimaldoc provider has a specific feature enabled
func (c AnalyticsConfig) HasMinimalDocFeature(feature string) bool {
	p := c.GetMinimalDocProvider()
	if p == nil {
		return false
	}
	features := p.GetConfigString("features")
	if features == "" {
		return feature == "analytics" // default feature
	}
	for _, f := range splitFeatures(features) {
		if f == feature {
			return true
		}
	}
	return false
}

// ShouldShowFeedback checks if feedback widget should show on a page type
func (c AnalyticsConfig) ShouldShowFeedback(pageType string) bool {
	if !c.HasMinimalDocFeature("feedback") {
		return false
	}
	p := c.GetMinimalDocProvider()
	if p == nil {
		return false
	}
	pages := p.GetConfigSlice("feedback_pages")
	if len(pages) == 0 {
		// Default: show on docs only (backward compatible)
		return pageType == "docs"
	}
	for _, pt := range pages {
		if pt == pageType {
			return true
		}
	}
	return false
}

// ShouldShowNewsletter checks if newsletter widget should show on a page type
func (c AnalyticsConfig) ShouldShowNewsletter(pageType string) bool {
	if !c.HasMinimalDocFeature("newsletter") {
		return false
	}
	p := c.GetMinimalDocProvider()
	if p == nil {
		return false
	}
	pages := p.GetConfigSlice("newsletter_pages")
	if len(pages) == 0 {
		// Default: show on docs only (backward compatible)
		return pageType == "docs"
	}
	for _, pt := range pages {
		if pt == pageType {
			return true
		}
	}
	return false
}

// splitFeatures splits comma-separated feature string
func splitFeatures(s string) []string {
	var result []string
	for _, f := range strings.Split(s, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			result = append(result, f)
		}
	}
	return result
}
