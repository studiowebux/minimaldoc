package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"
)

// Templates provides email template rendering.
type Templates struct {
	siteName string
	baseURL  string
}

// NewTemplates creates a new template renderer.
func NewTemplates(siteName, baseURL string) *Templates {
	return &Templates{
		siteName: siteName,
		baseURL:  baseURL,
	}
}

// VerificationEmail generates a newsletter verification email.
func (t *Templates) VerificationEmail(email, siteID, token string) *Message {
	verifyURL := fmt.Sprintf("%s/api/newsletter/verify?site_id=%s&token=%s", t.baseURL, siteID, token)

	data := struct {
		SiteName  string
		Email     string
		VerifyURL string
		Year      int
	}{
		SiteName:  t.siteName,
		Email:     email,
		VerifyURL: verifyURL,
		Year:      time.Now().Year(),
	}

	html := renderTemplate(verificationHTMLTemplate, data)
	text := renderTemplate(verificationTextTemplate, data)

	return &Message{
		To:       email,
		Subject:  fmt.Sprintf("Confirm your subscription to %s", t.siteName),
		HTMLBody: html,
		TextBody: text,
	}
}

// UnsubscribeEmail generates an unsubscribe confirmation email.
func (t *Templates) UnsubscribeEmail(email string) *Message {
	data := struct {
		SiteName string
		Email    string
		Year     int
	}{
		SiteName: t.siteName,
		Email:    email,
		Year:     time.Now().Year(),
	}

	html := renderTemplate(unsubscribeHTMLTemplate, data)
	text := renderTemplate(unsubscribeTextTemplate, data)

	return &Message{
		To:       email,
		Subject:  fmt.Sprintf("You've been unsubscribed from %s", t.siteName),
		HTMLBody: html,
		TextBody: text,
	}
}

// WelcomeEmail generates a welcome email after verification.
// unsubscribeToken is an HMAC-signed token from auth.SignUnsubscribeToken.
func (t *Templates) WelcomeEmail(email, unsubscribeToken string) *Message {
	unsubscribeURL := fmt.Sprintf("%s/api/newsletter/unsubscribe?token=%s", t.baseURL, unsubscribeToken)

	data := struct {
		SiteName       string
		Email          string
		UnsubscribeURL string
		Year           int
	}{
		SiteName:       t.siteName,
		Email:          email,
		UnsubscribeURL: unsubscribeURL,
		Year:           time.Now().Year(),
	}

	html := renderTemplate(welcomeHTMLTemplate, data)
	text := renderTemplate(welcomeTextTemplate, data)

	return &Message{
		To:       email,
		Subject:  fmt.Sprintf("Welcome to %s!", t.siteName),
		HTMLBody: html,
		TextBody: text,
	}
}

func renderTemplate(tmpl string, data any) string {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse email template: %v\n", err)
		return ""
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to render email template: %v\n", err)
		return ""
	}
	return buf.String()
}

// Email templates

const verificationHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Confirm your subscription</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
  <div style="background: #f8f9fa; padding: 30px; border-radius: 8px;">
    <h1 style="margin: 0 0 20px; font-size: 24px; color: #1a1a1a;">Confirm your subscription</h1>
    <p style="margin: 0 0 20px;">You requested to subscribe to <strong>{{.SiteName}}</strong> with this email address: <strong>{{.Email}}</strong></p>
    <p style="margin: 0 0 30px;">Click the button below to confirm your subscription:</p>
    <a href="{{.VerifyURL}}" style="display: inline-block; background: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: 500;">Confirm Subscription</a>
    <p style="margin: 30px 0 0; font-size: 14px; color: #666;">If you didn't request this, you can safely ignore this email.</p>
  </div>
  <p style="margin: 20px 0 0; font-size: 12px; color: #999; text-align: center;">&copy; {{.Year}} {{.SiteName}}</p>
</body>
</html>`

const verificationTextTemplate = `Confirm your subscription to {{.SiteName}}

You requested to subscribe with this email address: {{.Email}}

Click the link below to confirm your subscription:
{{.VerifyURL}}

If you didn't request this, you can safely ignore this email.

© {{.Year}} {{.SiteName}}`

const unsubscribeHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Unsubscribed</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
  <div style="background: #f8f9fa; padding: 30px; border-radius: 8px;">
    <h1 style="margin: 0 0 20px; font-size: 24px; color: #1a1a1a;">You've been unsubscribed</h1>
    <p style="margin: 0 0 20px;">Your email address <strong>{{.Email}}</strong> has been removed from the <strong>{{.SiteName}}</strong> mailing list.</p>
    <p style="margin: 0; color: #666;">We're sorry to see you go. If this was a mistake, you can always subscribe again.</p>
  </div>
  <p style="margin: 20px 0 0; font-size: 12px; color: #999; text-align: center;">&copy; {{.Year}} {{.SiteName}}</p>
</body>
</html>`

const unsubscribeTextTemplate = `You've been unsubscribed from {{.SiteName}}

Your email address {{.Email}} has been removed from our mailing list.

We're sorry to see you go. If this was a mistake, you can always subscribe again.

© {{.Year}} {{.SiteName}}`

const welcomeHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Welcome!</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
  <div style="background: #f8f9fa; padding: 30px; border-radius: 8px;">
    <h1 style="margin: 0 0 20px; font-size: 24px; color: #1a1a1a;">Welcome to {{.SiteName}}!</h1>
    <p style="margin: 0 0 20px;">Your subscription has been confirmed. You'll now receive updates from us.</p>
    <p style="margin: 0; color: #666;">Thank you for subscribing!</p>
  </div>
  <p style="margin: 20px 0 0; font-size: 12px; color: #999; text-align: center;">
    <a href="{{.UnsubscribeURL}}" style="color: #999;">Unsubscribe</a> &bull; &copy; {{.Year}} {{.SiteName}}
  </p>
</body>
</html>`

const welcomeTextTemplate = `Welcome to {{.SiteName}}!

Your subscription has been confirmed. You'll now receive updates from us.

Thank you for subscribing!

To unsubscribe: {{.UnsubscribeURL}}

© {{.Year}} {{.SiteName}}`
