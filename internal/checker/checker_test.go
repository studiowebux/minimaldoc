package checker

import (
	"net"
	"testing"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// ---------------------------------------------------------------------------
// isPrivateIP — SSRF protection
// ---------------------------------------------------------------------------

func TestIsPrivateIP_Loopback(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"127.0.0.1", true},
		{"127.0.0.2", true},
		{"127.255.255.255", true},
	}
	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if got := isPrivateIP(ip); got != tt.want {
			t.Errorf("isPrivateIP(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIsPrivateIP_RFC1918(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		// 10.0.0.0/8
		{"10.0.0.1", true},
		{"10.255.255.255", true},
		// 172.16.0.0/12
		{"172.16.0.1", true},
		{"172.31.255.255", true},
		// 192.168.0.0/16
		{"192.168.1.1", true},
		{"192.168.0.0", true},
	}
	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if got := isPrivateIP(ip); got != tt.want {
			t.Errorf("isPrivateIP(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIsPrivateIP_LinkLocal(t *testing.T) {
	ip := net.ParseIP("169.254.1.1")
	if !isPrivateIP(ip) {
		t.Error("isPrivateIP(169.254.1.1) = false, want true (link-local)")
	}
}

func TestIsPrivateIP_IPv6Loopback(t *testing.T) {
	ip := net.ParseIP("::1")
	if !isPrivateIP(ip) {
		t.Error("isPrivateIP(::1) = false, want true")
	}
}

func TestIsPrivateIP_IPv6ULA(t *testing.T) {
	ip := net.ParseIP("fc00::1")
	if !isPrivateIP(ip) {
		t.Error("isPrivateIP(fc00::1) = false, want true (ULA)")
	}
}

func TestIsPrivateIP_IPv6LinkLocal(t *testing.T) {
	ip := net.ParseIP("fe80::1")
	if !isPrivateIP(ip) {
		t.Error("isPrivateIP(fe80::1) = false, want true (link-local)")
	}
}

func TestIsPrivateIP_PublicAddresses(t *testing.T) {
	tests := []string{
		"8.8.8.8",
		"1.1.1.1",
		"93.184.216.34",
		"203.0.113.1",
		"2607:f8b0:4004:800::200e", // Google public IPv6
	}
	for _, ip := range tests {
		parsed := net.ParseIP(ip)
		if isPrivateIP(parsed) {
			t.Errorf("isPrivateIP(%s) = true, want false (public IP)", ip)
		}
	}
}

func TestIsPrivateIP_BoundaryAddresses(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		// Just outside 172.16.0.0/12
		{"172.15.255.255", false},
		{"172.32.0.0", false},
		// Just outside 10.0.0.0/8
		{"11.0.0.0", false},
		{"9.255.255.255", false},
		// Just outside 192.168.0.0/16
		{"192.167.255.255", false},
		{"192.169.0.0", false},
	}
	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if got := isPrivateIP(ip); got != tt.want {
			t.Errorf("isPrivateIP(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// classifyLink
// ---------------------------------------------------------------------------

func TestClassifyLink(t *testing.T) {
	tests := []struct {
		url  string
		want core.LinkType
	}{
		// External
		{"https://example.com", core.LinkTypeExternal},
		{"http://example.com/page", core.LinkTypeExternal},
		{"//cdn.example.com/script.js", core.LinkTypeExternal},

		// Email
		{"mailto:user@example.com", core.LinkTypeEmail},

		// Anchor
		{"#section", core.LinkTypeInternalAnchor},
		{"#top", core.LinkTypeInternalAnchor},

		// Assets
		{"/images/logo.png", core.LinkTypeInternalAsset},
		{"assets/file.pdf", core.LinkTypeInternalAsset},
		{"photo.jpg", core.LinkTypeInternalAsset},
		{"style.css", core.LinkTypeInternalAsset},
		{"app.js", core.LinkTypeInternalAsset},
		{"video.mp4", core.LinkTypeInternalAsset},

		// Internal pages
		{"getting-started", core.LinkTypeInternalPage},
		{"/docs/intro.md", core.LinkTypeInternalPage},
		{"./page.md", core.LinkTypeInternalPage},
		{"/about", core.LinkTypeInternalPage},

		// Other protocols
		{"tel:+1234567890", core.LinkTypeOther},
		{"javascript:alert(1)", core.LinkTypeOther},

		// Empty / whitespace
		{"", core.LinkTypeOther},
		{"   ", core.LinkTypeOther},
	}

	for _, tt := range tests {
		got := classifyLink(tt.url)
		if got != tt.want {
			t.Errorf("classifyLink(%q) = %v (%s), want %v (%s)", tt.url, got, got.String(), tt.want, tt.want.String())
		}
	}
}

// ---------------------------------------------------------------------------
// markdownLinkRegex
// ---------------------------------------------------------------------------

func TestMarkdownLinkRegex(t *testing.T) {
	tests := []struct {
		input    string
		wantText string
		wantURL  string
	}{
		{"[Click here](https://example.com)", "Click here", "https://example.com"},
		{"[API](./api.md)", "API", "./api.md"},
		{`[Title](/page "tooltip")`, "Title", "/page"},
		{"[](empty-text.md)", "", "empty-text.md"},
	}

	for _, tt := range tests {
		matches := markdownLinkRegex.FindStringSubmatch(tt.input)
		if matches == nil {
			t.Errorf("markdownLinkRegex did not match %q", tt.input)
			continue
		}
		if matches[1] != tt.wantText {
			t.Errorf("markdownLinkRegex text = %q, want %q for input %q", matches[1], tt.wantText, tt.input)
		}
		if matches[2] != tt.wantURL {
			t.Errorf("markdownLinkRegex url = %q, want %q for input %q", matches[2], tt.wantURL, tt.input)
		}
	}
}

// ---------------------------------------------------------------------------
// htmlLinkRegex
// ---------------------------------------------------------------------------

func TestHTMLLinkRegex(t *testing.T) {
	tests := []struct {
		input   string
		wantURL string
	}{
		{`<a href="https://example.com">Link</a>`, "https://example.com"},
		{`<a class="btn" href="/page">Page</a>`, "/page"},
		{`<a href='./local.md'>Local</a>`, "./local.md"},
	}

	for _, tt := range tests {
		matches := htmlLinkRegex.FindStringSubmatch(tt.input)
		if matches == nil {
			t.Errorf("htmlLinkRegex did not match %q", tt.input)
			continue
		}
		if matches[1] != tt.wantURL {
			t.Errorf("htmlLinkRegex url = %q, want %q for input %q", matches[1], tt.wantURL, tt.input)
		}
	}
}

// ---------------------------------------------------------------------------
// imageRegex
// ---------------------------------------------------------------------------

func TestImageRegex(t *testing.T) {
	tests := []struct {
		input   string
		wantAlt string
		wantURL string
	}{
		{"![Logo](logo.png)", "Logo", "logo.png"},
		{"![](empty.svg)", "", "empty.svg"},
		{`![Screenshot](/img/shot.jpg "title")`, "Screenshot", "/img/shot.jpg"},
	}

	for _, tt := range tests {
		matches := imageRegex.FindStringSubmatch(tt.input)
		if matches == nil {
			t.Errorf("imageRegex did not match %q", tt.input)
			continue
		}
		if matches[1] != tt.wantAlt {
			t.Errorf("imageRegex alt = %q, want %q for input %q", matches[1], tt.wantAlt, tt.input)
		}
		if matches[2] != tt.wantURL {
			t.Errorf("imageRegex url = %q, want %q for input %q", matches[2], tt.wantURL, tt.input)
		}
	}
}

// ---------------------------------------------------------------------------
// matchesGlob
// ---------------------------------------------------------------------------

func TestMatchesGlob(t *testing.T) {
	tests := []struct {
		s       string
		pattern string
		want    bool
	}{
		{"https://example.com/page", "https://example.com/*", true},
		{"https://other.com/page", "https://example.com/*", false},
		{"/docs/api", "/docs/*", true},
		{"/docs/api", "/docs/api", true},
		{"/docs/api", "/other/*", false},
	}

	for _, tt := range tests {
		got := matchesGlob(tt.s, tt.pattern)
		if got != tt.want {
			t.Errorf("matchesGlob(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// NewLinkCollector
// ---------------------------------------------------------------------------

func TestNewLinkCollector(t *testing.T) {
	c := NewLinkCollector("/docs")
	if c.docsRoot != "/docs" {
		t.Errorf("docsRoot = %q, want %q", c.docsRoot, "/docs")
	}
	if len(c.Links()) != 0 {
		t.Errorf("Links() should be empty initially, got %d", len(c.Links()))
	}
}
