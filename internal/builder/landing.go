package builder

import (
	"github.com/studiowebux/minimaldoc/internal/core"
)

// LandingBuilder handles building the landing page
type LandingBuilder struct{}

// NewLandingBuilder creates a new landing builder
func NewLandingBuilder() *LandingBuilder {
	return &LandingBuilder{}
}

// Build creates the landing page from configuration
func (lb *LandingBuilder) Build(config core.LandingConfig) (*core.LandingPage, error) {
	if !config.Enabled {
		return nil, nil
	}

	page := &core.LandingPage{
		Config: config,
	}

	// Set sections if they have content
	if config.Hero.Title != "" {
		page.Hero = &config.Hero
	}

	if len(config.Features.Items) > 0 {
		page.Features = &config.Features
	}

	if len(config.Steps.Items) > 0 {
		page.Steps = &config.Steps
	}

	if config.CTA.Title != "" {
		page.CTA = &config.CTA
	}

	if len(config.Testimonials.Items) > 0 {
		page.Testimonials = &config.Testimonials
	}

	if config.OpenSource.Title != "" {
		page.OpenSource = &config.OpenSource
	}

	if len(config.Links.Items) > 0 {
		page.Links = &config.Links
	}

	return page, nil
}
