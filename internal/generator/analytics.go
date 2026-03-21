package generator

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// AnalyticsFuncMap returns template functions for analytics rendering.
func AnalyticsFuncMap() template.FuncMap {
	return template.FuncMap{
		"analyticsHeadScripts":     analyticsHeadScripts,
		"analyticsBodyScripts":     analyticsBodyScripts,
		"hasMinimalDocFeedback":    hasMinimalDocFeedback,
		"hasMinimalDocNewsletter":  hasMinimalDocNewsletter,
		"minimalDocFeedbackWidget": minimalDocFeedbackWidget,
		"minimalDocNewsletterForm": minimalDocNewsletterForm,
		"feedbackWidgetFor":        feedbackWidgetFor,
		"newsletterFormFor":        newsletterFormFor,
	}
}

// hasMinimalDocFeedback checks if feedback feature is enabled
func hasMinimalDocFeedback(config core.AnalyticsConfig) bool {
	return config.HasMinimalDocFeature("feedback")
}

// hasMinimalDocNewsletter checks if newsletter feature is enabled
func hasMinimalDocNewsletter(config core.AnalyticsConfig) bool {
	return config.HasMinimalDocFeature("newsletter")
}

// minimalDocFeedbackWidget renders the feedback widget placeholder (for docs pages - backward compat)
func minimalDocFeedbackWidget(config core.AnalyticsConfig) template.HTML {
	return feedbackWidgetFor(config, "docs")
}

// minimalDocNewsletterForm renders the newsletter form placeholder (for docs pages - backward compat)
func minimalDocNewsletterForm(config core.AnalyticsConfig) template.HTML {
	return newsletterFormFor(config, "docs")
}

// feedbackWidgetFor renders the feedback widget for a specific page type
func feedbackWidgetFor(config core.AnalyticsConfig, pageType string) template.HTML {
	if !config.ShouldShowFeedback(pageType) {
		return ""
	}
	return template.HTML(`<div class="minimaldoc-feedback-wrapper" data-minimaldoc-feedback></div>`)
}

// newsletterFormFor renders the newsletter form for a specific page type
func newsletterFormFor(config core.AnalyticsConfig, pageType string) template.HTML {
	if !config.ShouldShowNewsletter(pageType) {
		return ""
	}
	return template.HTML(`<div class="minimaldoc-newsletter-wrapper">
<form data-minimaldoc-newsletter class="minimaldoc-newsletter-form">
<label for="newsletter-email">Subscribe to updates</label>
<div class="newsletter-input-group">
<input type="email" id="newsletter-email" placeholder="Enter your email" required>
<button type="submit">Subscribe</button>
</div>
</form>
</div>`)
}

// analyticsHeadScripts generates script tags for analytics providers that go in <head>
func analyticsHeadScripts(config core.AnalyticsConfig) template.HTML {
	if !config.Enabled {
		return ""
	}

	var scripts strings.Builder
	for _, provider := range config.GetEnabledProviders() {
		script := renderProviderScript(provider, "head")
		if script != "" {
			scripts.WriteString(script)
			scripts.WriteString("\n")
		}
	}
	return template.HTML(scripts.String())
}

// analyticsBodyScripts generates script tags for analytics providers that go before </body>
func analyticsBodyScripts(config core.AnalyticsConfig) template.HTML {
	if !config.Enabled {
		return ""
	}

	var scripts strings.Builder
	for _, provider := range config.GetEnabledProviders() {
		script := renderProviderScript(provider, "body")
		if script != "" {
			scripts.WriteString(script)
			scripts.WriteString("\n")
		}
	}
	return template.HTML(scripts.String())
}

// renderProviderScript renders the script tag(s) for a specific provider
func renderProviderScript(provider core.AnalyticsProvider, location string) string {
	switch provider.Type {
	case "ga4":
		if location != "head" {
			return ""
		}
		return renderGA4(provider)
	case "plausible":
		if location != "head" {
			return ""
		}
		return renderPlausible(provider)
	case "umami":
		if location != "head" {
			return ""
		}
		return renderUmami(provider)
	case "matomo":
		if location != "head" {
			return ""
		}
		return renderMatomo(provider)
	case "fathom":
		if location != "head" {
			return ""
		}
		return renderFathom(provider)
	case "simple":
		if location != "head" {
			return ""
		}
		return renderSimpleAnalytics(provider)
	case "minimaldoc":
		if location != "body" {
			return ""
		}
		return renderMinimalDoc(provider)
	case "custom":
		return renderCustom(provider, location)
	default:
		return ""
	}
}

// renderGA4 renders Google Analytics 4 script tags
func renderGA4(provider core.AnalyticsProvider) string {
	measurementID := provider.GetConfigString("measurement_id")
	if measurementID == "" {
		return ""
	}
	return fmt.Sprintf(`<script async src="https://www.googletagmanager.com/gtag/js?id=%s"></script>
<script>window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','%s');</script>`,
		escapeHTML(measurementID), escapeJS(measurementID))
}

// renderPlausible renders Plausible Analytics script tag
func renderPlausible(provider core.AnalyticsProvider) string {
	domain := provider.GetConfigString("domain")
	if domain == "" {
		return ""
	}
	src := provider.GetConfigString("src")
	if src == "" {
		src = "https://plausible.io/js/script.js"
	}
	return fmt.Sprintf(`<script defer data-domain="%s" src="%s"></script>`,
		escapeHTML(domain), escapeHTML(src))
}

// renderUmami renders Umami Analytics script tag
func renderUmami(provider core.AnalyticsProvider) string {
	websiteID := provider.GetConfigString("website_id")
	src := provider.GetConfigString("src")
	if websiteID == "" || src == "" {
		return ""
	}
	return fmt.Sprintf(`<script defer data-website-id="%s" src="%s"></script>`,
		escapeHTML(websiteID), escapeHTML(src))
}

// renderMatomo renders Matomo Analytics script tag
func renderMatomo(provider core.AnalyticsProvider) string {
	url := provider.GetConfigString("url")
	siteID := provider.GetConfigString("site_id")
	if url == "" || siteID == "" {
		return ""
	}
	// Ensure URL ends without trailing slash for consistency
	url = strings.TrimSuffix(url, "/")
	return fmt.Sprintf(`<script>var _paq=window._paq=window._paq||[];_paq.push(['trackPageView']);_paq.push(['enableLinkTracking']);(function(){var u="%s/";_paq.push(['setTrackerUrl',u+'matomo.php']);_paq.push(['setSiteId','%s']);var d=document,g=d.createElement('script'),s=d.getElementsByTagName('script')[0];g.async=true;g.src=u+'matomo.js';s.parentNode.insertBefore(g,s);})();</script>`,
		escapeJS(url), escapeJS(siteID))
}

// renderFathom renders Fathom Analytics script tag
func renderFathom(provider core.AnalyticsProvider) string {
	siteID := provider.GetConfigString("site_id")
	if siteID == "" {
		return ""
	}
	src := provider.GetConfigString("src")
	if src == "" {
		src = "https://cdn.usefathom.com/script.js"
	}
	return fmt.Sprintf(`<script src="%s" data-site="%s" defer></script>`,
		escapeHTML(src), escapeHTML(siteID))
}

// renderSimpleAnalytics renders Simple Analytics script tag
func renderSimpleAnalytics(provider core.AnalyticsProvider) string {
	src := provider.GetConfigString("src")
	if src == "" {
		src = "https://scripts.simpleanalyticscdn.com/latest.js"
	}
	return fmt.Sprintf(`<script async defer src="%s"></script>
<noscript><img src="https://queue.simpleanalyticscdn.com/noscript.gif" alt="" referrerpolicy="no-referrer-when-downgrade"/></noscript>`,
		escapeHTML(src))
}

// renderMinimalDoc renders the MinimalDoc backend tracking script
func renderMinimalDoc(provider core.AnalyticsProvider) string {
	endpoint := provider.GetConfigString("endpoint")
	siteID := provider.GetConfigString("site_id")
	if endpoint == "" || siteID == "" {
		return ""
	}

	// Features to enable (default: analytics only)
	features := provider.GetConfigString("features")
	if features == "" {
		features = "analytics"
	}

	// Debug mode
	debug := ""
	if provider.GetConfigString("debug") == "true" {
		debug = ` data-debug`
	}

	return fmt.Sprintf(`<script src="%s/minimaldoc.js" data-endpoint="%s" data-site-id="%s" data-features="%s"%s defer></script>`,
		escapeHTML(endpoint), escapeHTML(endpoint), escapeHTML(siteID), escapeHTML(features), debug)
}

// renderCustom renders a custom analytics script tag with arbitrary attributes
func renderCustom(provider core.AnalyticsProvider, location string) string {
	// Check if this provider should render in the specified location
	configLocation := provider.GetConfigString("location")
	if configLocation == "" {
		configLocation = "head" // Default to head
	}
	if configLocation != location {
		return ""
	}

	src := provider.GetConfigString("src")
	if src == "" {
		return ""
	}

	var attrs strings.Builder
	attrMap := provider.GetConfigMap("attrs")
	for k, v := range attrMap {
		attrs.WriteString(fmt.Sprintf(` %s="%s"`, escapeHTML(k), escapeHTML(v)))
	}

	// Check for defer/async attributes
	if provider.GetConfigString("defer") == "true" {
		attrs.WriteString(" defer")
	}
	if provider.GetConfigString("async") == "true" {
		attrs.WriteString(" async")
	}

	return fmt.Sprintf(`<script%s src="%s"></script>`, attrs.String(), escapeHTML(src))
}

// escapeHTML escapes special HTML characters for use in attributes.
func escapeHTML(s string) string {
	return html.EscapeString(s)
}

// escapeJS escapes a string for safe embedding inside a JavaScript string literal
// within a <script> tag (backslash-escaping, not HTML entity escaping).
func escapeJS(s string) string {
	return template.JSEscapeString(s)
}
