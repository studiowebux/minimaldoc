package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/studiowebux/minimaldoc/internal/server/config"
)

// SMTPSender implements email sending via SMTP.
type SMTPSender struct {
	host        string
	port        int
	user        string
	pass        string
	fromAddress string
	fromName    string
}

// NewSMTPSender creates a new SMTP email sender.
func NewSMTPSender(cfg config.EmailConfig) (*SMTPSender, error) {
	return &SMTPSender{
		host:        cfg.SMTPHost,
		port:        cfg.SMTPPort,
		user:        cfg.SMTPUser,
		pass:        cfg.SMTPPass,
		fromAddress: cfg.FromAddress,
		fromName:    cfg.FromName,
	}, nil
}

// Send sends an email to a single recipient.
func (s *SMTPSender) Send(ctx context.Context, msg *Message) error {
	from := s.formatFrom(msg)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Build email content
	headers := fmt.Sprintf("From: %s\r\n", from)
	headers += fmt.Sprintf("To: %s\r\n", msg.To)
	headers += fmt.Sprintf("Subject: %s\r\n", msg.Subject)
	headers += "MIME-Version: 1.0\r\n"

	if msg.ReplyTo != "" {
		headers += fmt.Sprintf("Reply-To: %s\r\n", msg.ReplyTo)
	}

	var body string
	if msg.HTMLBody != "" {
		headers += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
		body = msg.HTMLBody
	} else {
		headers += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
		body = msg.TextBody
	}

	content := headers + "\r\n" + body

	// Create auth
	var auth smtp.Auth
	if s.user != "" && s.pass != "" {
		auth = smtp.PlainAuth("", s.user, s.pass, s.host)
	}

	// Send email
	err := smtp.SendMail(addr, auth, s.fromAddress, []string{msg.To}, []byte(content))
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	return nil
}

// SendBulk sends emails to multiple recipients.
func (s *SMTPSender) SendBulk(ctx context.Context, msgs []*Message) error {
	for _, msg := range msgs {
		if err := s.Send(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// formatFrom formats the From address with optional name.
func (s *SMTPSender) formatFrom(msg *Message) string {
	fromAddr := s.fromAddress
	fromName := s.fromName

	if msg.FromAddress != "" {
		fromAddr = msg.FromAddress
	}
	if msg.FromName != "" {
		fromName = msg.FromName
	}

	if fromName != "" {
		return fmt.Sprintf("%s <%s>", fromName, fromAddr)
	}
	return fromAddr
}
