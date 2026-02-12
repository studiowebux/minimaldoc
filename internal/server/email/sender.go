// Package email provides pluggable email sending for minimaldoc-server.
// Supports AWS SES, SMTP, SendGrid, and other providers.
package email

import (
	"context"
	"fmt"

	"github.com/studiowebux/minimaldoc/internal/server/config"
)

// Sender defines the interface for email providers.
type Sender interface {
	// Send sends an email to a single recipient.
	Send(ctx context.Context, msg *Message) error

	// SendBulk sends emails to multiple recipients.
	SendBulk(ctx context.Context, msgs []*Message) error
}

// Message represents an email message.
type Message struct {
	To          string
	Subject     string
	HTMLBody    string
	TextBody    string
	FromAddress string // Optional override
	FromName    string // Optional override
	ReplyTo     string
	Headers     map[string]string
}

// NewSender creates an email sender based on configuration.
// Currently supports: smtp, mock
// Future: ses, sendgrid (add when needed)
func NewSender(cfg config.EmailConfig) (Sender, error) {
	switch cfg.Provider {
	case "smtp":
		return NewSMTPSender(cfg)
	case "mock", "test", "":
		return NewMockSender(), nil
	default:
		return nil, fmt.Errorf("unknown email provider: %s (supported: smtp, mock)", cfg.Provider)
	}
}

// TemplateData contains common template variables.
type TemplateData struct {
	SiteName        string
	SiteURL         string
	RecipientEmail  string
	RecipientName   string
	VerificationURL string
	UnsubscribeURL  string
	Year            int
}
