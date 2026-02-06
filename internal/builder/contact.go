package builder

import (
	"github.com/studiowebux/minimaldoc/internal/core"
)

// ContactBuilder handles building the contact page
type ContactBuilder struct{}

// NewContactBuilder creates a new contact builder
func NewContactBuilder() *ContactBuilder {
	return &ContactBuilder{}
}

// Build creates the contact page from configuration
func (cb *ContactBuilder) Build(config core.ContactConfig) (*core.ContactPage, error) {
	if !config.Enabled {
		return nil, nil
	}

	page := &core.ContactPage{
		Config: config,
	}

	return page, nil
}
