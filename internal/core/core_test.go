package core

import (
	"testing"
)

// ---------------------------------------------------------------------------
// DefaultSiteConfig
// ---------------------------------------------------------------------------

func TestDefaultSiteConfig_BasicFields(t *testing.T) {
	cfg := DefaultSiteConfig()

	if cfg.Title != "Documentation" {
		t.Errorf("Title = %q, want %q", cfg.Title, "Documentation")
	}
	if cfg.Description != "Documentation site powered by Minimal Doc" {
		t.Errorf("Description = %q, want %q", cfg.Description, "Documentation site powered by Minimal Doc")
	}
	if cfg.Theme != "default" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "default")
	}
	if cfg.DarkMode {
		t.Error("DarkMode should be false by default")
	}
	if !cfg.EnableLLMS {
		t.Error("EnableLLMS should be true by default")
	}
	if cfg.EnableSearch {
		t.Error("EnableSearch should be false by default")
	}
	if cfg.NavDepth != 0 {
		t.Errorf("NavDepth = %d, want 0", cfg.NavDepth)
	}
	if cfg.CleanURLs {
		t.Error("CleanURLs should be false by default")
	}
	if cfg.Custom == nil {
		t.Error("Custom map should be initialized, got nil")
	}
}

func TestDefaultSiteConfig_SubConfigsNotNil(t *testing.T) {
	cfg := DefaultSiteConfig()

	// Verify sub-configs are initialized (not zero-valued)
	if cfg.ThemeConfig.Name != "default" {
		t.Errorf("ThemeConfig.Name = %q, want %q", cfg.ThemeConfig.Name, "default")
	}
	if cfg.Changelog.Title != "Changelog" {
		t.Errorf("Changelog.Title = %q, want %q", cfg.Changelog.Title, "Changelog")
	}
	if cfg.LinkCheck.ExternalTimeout != 5 {
		t.Errorf("LinkCheck.ExternalTimeout = %d, want 5", cfg.LinkCheck.ExternalTimeout)
	}
}

// ---------------------------------------------------------------------------
// DefaultStaleWarningConfig
// ---------------------------------------------------------------------------

func TestDefaultStaleWarningConfig(t *testing.T) {
	cfg := DefaultStaleWarningConfig()

	if cfg.Enabled {
		t.Error("Enabled should be false")
	}
	if cfg.ThresholdDays != 365 {
		t.Errorf("ThresholdDays = %d, want 365", cfg.ThresholdDays)
	}
	if cfg.Message != "" {
		t.Errorf("Message = %q, want empty string", cfg.Message)
	}
	if !cfg.ShowUpdateDate {
		t.Error("ShowUpdateDate should be true")
	}
}

// ---------------------------------------------------------------------------
// DefaultPDFExportConfig
// ---------------------------------------------------------------------------

func TestDefaultPDFExportConfig(t *testing.T) {
	cfg := DefaultPDFExportConfig()

	if cfg.Enabled {
		t.Error("Enabled should be false")
	}
	if cfg.PageBreakLevel != 1 {
		t.Errorf("PageBreakLevel = %d, want 1", cfg.PageBreakLevel)
	}
}

// ---------------------------------------------------------------------------
// DefaultClaudeAssistConfig
// ---------------------------------------------------------------------------

func TestDefaultClaudeAssistConfig(t *testing.T) {
	cfg := DefaultClaudeAssistConfig()

	if cfg.Enabled {
		t.Error("Enabled should be false")
	}
	if cfg.Prompt != "" {
		t.Errorf("Prompt = %q, want empty string", cfg.Prompt)
	}
	if cfg.Label != "Ask Claude" {
		t.Errorf("Label = %q, want %q", cfg.Label, "Ask Claude")
	}
}

// ---------------------------------------------------------------------------
// DefaultChangelogConfig
// ---------------------------------------------------------------------------

func TestDefaultChangelogConfig(t *testing.T) {
	cfg := DefaultChangelogConfig()

	if cfg.Enabled {
		t.Error("Enabled should be false")
	}
	if cfg.Title != "Changelog" {
		t.Errorf("Title = %q, want %q", cfg.Title, "Changelog")
	}
	if cfg.Path != "changelog" {
		t.Errorf("Path = %q, want %q", cfg.Path, "changelog")
	}
	if !cfg.RSSEnabled {
		t.Error("RSSEnabled should be true")
	}
}

// ---------------------------------------------------------------------------
// DefaultLinkCheckConfig
// ---------------------------------------------------------------------------

func TestDefaultLinkCheckConfig(t *testing.T) {
	cfg := DefaultLinkCheckConfig()

	if !cfg.Enabled {
		t.Error("Enabled should be true")
	}
	if cfg.Mode != LinkCheckWarn {
		t.Errorf("Mode = %q, want %q", cfg.Mode, LinkCheckWarn)
	}
	if cfg.CheckExternal {
		t.Error("CheckExternal should be false")
	}
	if cfg.ExternalTimeout != 5 {
		t.Errorf("ExternalTimeout = %d, want 5", cfg.ExternalTimeout)
	}
}

// ---------------------------------------------------------------------------
// DefaultConfigRegistry
// ---------------------------------------------------------------------------

func TestDefaultConfigRegistry_AllEntriesCallable(t *testing.T) {
	for name, fn := range DefaultConfigRegistry {
		result := fn()
		if result == nil {
			t.Errorf("DefaultConfigRegistry[%q] returned nil", name)
		}
	}
}

func TestDefaultConfigs_DelegatesToDefaultSiteConfig(t *testing.T) {
	a := DefaultConfigs()
	b := DefaultSiteConfig()

	if a.Title != b.Title || a.Theme != b.Theme || a.Description != b.Description {
		t.Error("DefaultConfigs() should return the same result as DefaultSiteConfig()")
	}
}

// ---------------------------------------------------------------------------
// DefaultMetadata
// ---------------------------------------------------------------------------

func TestDefaultMetadata(t *testing.T) {
	m := DefaultMetadata()

	if m.Title != "Untitled" {
		t.Errorf("Title = %q, want %q", m.Title, "Untitled")
	}
	if m.MenuOrder != -1 {
		t.Errorf("MenuOrder = %d, want -1", m.MenuOrder)
	}
	if m.Hidden {
		t.Error("Hidden should be false")
	}
	if m.Tags == nil {
		t.Error("Tags should be initialized, got nil")
	}
	if m.Custom == nil {
		t.Error("Custom should be initialized, got nil")
	}
}

// ---------------------------------------------------------------------------
// Page.IsHidden
// ---------------------------------------------------------------------------

func TestPage_IsHidden(t *testing.T) {
	p := &Page{Metadata: Metadata{Hidden: false}}
	if p.IsHidden() {
		t.Error("IsHidden() = true, want false when Hidden=false")
	}

	p.Metadata.Hidden = true
	if !p.IsHidden() {
		t.Error("IsHidden() = false, want true when Hidden=true")
	}
}

// ---------------------------------------------------------------------------
// Page.Title — priority: MenuTitle > Metadata.Title > filename
// ---------------------------------------------------------------------------

func TestPage_Title_MenuTitle(t *testing.T) {
	p := &Page{
		SourcePath: "/docs/01-intro.md",
		Metadata: Metadata{
			Title:     "Introduction",
			MenuTitle: "Intro",
		},
	}
	if got := p.Title(); got != "Intro" {
		t.Errorf("Title() = %q, want %q (MenuTitle takes precedence)", got, "Intro")
	}
}

func TestPage_Title_MetadataTitle(t *testing.T) {
	p := &Page{
		SourcePath: "/docs/01-intro.md",
		Metadata:   Metadata{Title: "Introduction"},
	}
	if got := p.Title(); got != "Introduction" {
		t.Errorf("Title() = %q, want %q", got, "Introduction")
	}
}

func TestPage_Title_FallbackToFilename(t *testing.T) {
	p := &Page{
		SourcePath: "/docs/01-getting-started.md",
		Metadata:   Metadata{},
	}
	got := p.Title()
	if got != "Getting Started" {
		t.Errorf("Title() = %q, want %q", got, "Getting Started")
	}
}

// ---------------------------------------------------------------------------
// GenerateSlugFromPath
// ---------------------------------------------------------------------------

func TestGenerateSlugFromPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"01-getting-started.md", "getting-started"},
		{"02_api/index.md", "api/index"},
		{"README.md", "readme"},
		{"My Page.md", "my-page"},
		{"100-advanced.md", "advanced"},
		{"simple.md", "simple"},
	}

	for _, tt := range tests {
		got := GenerateSlugFromPath(tt.input)
		if got != tt.want {
			t.Errorf("GenerateSlugFromPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// ChangeType.Color
// ---------------------------------------------------------------------------

func TestChangeType_Color(t *testing.T) {
	tests := []struct {
		ct   ChangeType
		want string
	}{
		{ChangeAdded, "green"},
		{ChangeChanged, "blue"},
		{ChangeDeprecated, "yellow"},
		{ChangeRemoved, "red"},
		{ChangeFixed, "purple"},
		{ChangeSecurity, "orange"},
		{ChangeType("Unknown"), "gray"},
	}

	for _, tt := range tests {
		if got := tt.ct.Color(); got != tt.want {
			t.Errorf("ChangeType(%q).Color() = %q, want %q", tt.ct, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// ChangeType.IsValid
// ---------------------------------------------------------------------------

func TestChangeType_IsValid(t *testing.T) {
	for _, ct := range AllChangeTypes() {
		if !ct.IsValid() {
			t.Errorf("IsValid() = false for %q, want true", ct)
		}
	}
	if ChangeType("bogus").IsValid() {
		t.Error("IsValid() = true for invalid type, want false")
	}
}

// ---------------------------------------------------------------------------
// ParseChangeType
// ---------------------------------------------------------------------------

func TestParseChangeType(t *testing.T) {
	tests := []struct {
		input string
		want  ChangeType
		ok    bool
	}{
		{"added", ChangeAdded, true},
		{"CHANGED", ChangeChanged, true},
		{"  Fixed  ", ChangeFixed, true},
		{"Security", ChangeSecurity, true},
		{"deprecated", ChangeDeprecated, true},
		{"removed", ChangeRemoved, true},
		{"invalid", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got, ok := ParseChangeType(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseChangeType(%q) = (%q, %v), want (%q, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

// ---------------------------------------------------------------------------
// ChangeType.Label
// ---------------------------------------------------------------------------

func TestChangeType_Label(t *testing.T) {
	if got := ChangeAdded.Label(); got != "Added" {
		t.Errorf("Label() = %q, want %q", got, "Added")
	}
}

// ---------------------------------------------------------------------------
// AllChangeTypes
// ---------------------------------------------------------------------------

func TestAllChangeTypes_Count(t *testing.T) {
	types := AllChangeTypes()
	if len(types) != 6 {
		t.Errorf("AllChangeTypes() returned %d types, want 6", len(types))
	}
}

// ---------------------------------------------------------------------------
// ThemeConfig.HasCustomColors
// ---------------------------------------------------------------------------

func TestThemeConfig_HasCustomColors(t *testing.T) {
	// Default config has colors set, so HasCustomColors should be true
	cfg := DefaultThemeConfig()
	if !cfg.HasCustomColors() {
		t.Error("HasCustomColors() = false for default config with colors set, want true")
	}

	// Empty config should return false
	empty := ThemeConfig{}
	if empty.HasCustomColors() {
		t.Error("HasCustomColors() = true for zero-value config, want false")
	}

	// Only dark BgPrimary set
	darkOnly := ThemeConfig{Colors: ThemeColors{Dark: ThemeColorSet{BgPrimary: "#000"}}}
	if !darkOnly.HasCustomColors() {
		t.Error("HasCustomColors() = false when Dark.BgPrimary is set, want true")
	}
}

// ---------------------------------------------------------------------------
// ThemeConfig.HasCustomFonts
// ---------------------------------------------------------------------------

func TestThemeConfig_HasCustomFonts(t *testing.T) {
	// Zero-value should return false
	empty := ThemeConfig{}
	if empty.HasCustomFonts() {
		t.Error("HasCustomFonts() = true for zero-value config, want false")
	}

	// Heading font set
	withHeading := ThemeConfig{Fonts: ThemeFonts{Heading: "Inter"}}
	if !withHeading.HasCustomFonts() {
		t.Error("HasCustomFonts() = false when Heading is set, want true")
	}

	// Body font set
	withBody := ThemeConfig{Fonts: ThemeFonts{Body: "Roboto"}}
	if !withBody.HasCustomFonts() {
		t.Error("HasCustomFonts() = false when Body is set, want true")
	}

	// GoogleURL set
	withURL := ThemeConfig{Fonts: ThemeFonts{GoogleURL: "https://fonts.googleapis.com/css2?family=Inter"}}
	if !withURL.HasCustomFonts() {
		t.Error("HasCustomFonts() = false when GoogleURL is set, want true")
	}
}

// ---------------------------------------------------------------------------
// ThemeConfig.HasHeroBackground
// ---------------------------------------------------------------------------

func TestThemeConfig_HasHeroBackground(t *testing.T) {
	empty := ThemeConfig{}
	if empty.HasHeroBackground() {
		t.Error("HasHeroBackground() = true for zero-value config, want false")
	}

	withBg := ThemeConfig{Hero: ThemeHero{BackgroundImage: "/img/hero.jpg"}}
	if !withBg.HasHeroBackground() {
		t.Error("HasHeroBackground() = false when BackgroundImage is set, want true")
	}
}

// ---------------------------------------------------------------------------
// DefaultThemeConfig
// ---------------------------------------------------------------------------

func TestDefaultThemeConfig(t *testing.T) {
	cfg := DefaultThemeConfig()

	if cfg.Name != "default" {
		t.Errorf("Name = %q, want %q", cfg.Name, "default")
	}
	if cfg.Hero.TextAlign != "center" {
		t.Errorf("Hero.TextAlign = %q, want %q", cfg.Hero.TextAlign, "center")
	}
	// Light mode accent color should be set
	if cfg.Colors.Light.AccentPrimary == "" {
		t.Error("Light.AccentPrimary should not be empty")
	}
	// Dark mode accent color should be set
	if cfg.Colors.Dark.AccentPrimary == "" {
		t.Error("Dark.AccentPrimary should not be empty")
	}
}

// ---------------------------------------------------------------------------
// NewSite
// ---------------------------------------------------------------------------

func TestNewSite(t *testing.T) {
	cfg := DefaultSiteConfig()
	site := NewSite("/docs", "/public", cfg)

	if site.DocsRoot != "/docs" {
		t.Errorf("DocsRoot = %q, want %q", site.DocsRoot, "/docs")
	}
	if site.OutputRoot != "/public" {
		t.Errorf("OutputRoot = %q, want %q", site.OutputRoot, "/public")
	}
	if site.Pages == nil {
		t.Error("Pages should be initialized")
	}
	if site.Navigation == nil {
		t.Error("Navigation should be initialized")
	}
	if site.VersionedPages == nil {
		t.Error("VersionedPages should be initialized")
	}
	if site.LocalizedPages == nil {
		t.Error("LocalizedPages should be initialized")
	}
}

// ---------------------------------------------------------------------------
// LinkCheckResult methods
// ---------------------------------------------------------------------------

func TestLinkCheckResult_HasErrors(t *testing.T) {
	r := &LinkCheckResult{BrokenLinks: []BrokenLink{}}
	if r.HasErrors() {
		t.Error("HasErrors() = true with empty BrokenLinks, want false")
	}

	r.BrokenLinks = append(r.BrokenLinks, BrokenLink{Reason: "not found"})
	if !r.HasErrors() {
		t.Error("HasErrors() = false with broken links, want true")
	}
}

func TestLinkCheckResult_BrokenCount(t *testing.T) {
	r := &LinkCheckResult{BrokenLinks: []BrokenLink{{}, {}, {}}}
	if got := r.BrokenCount(); got != 3 {
		t.Errorf("BrokenCount() = %d, want 3", got)
	}
}

// ---------------------------------------------------------------------------
// LinkType.String
// ---------------------------------------------------------------------------

func TestLinkType_String(t *testing.T) {
	tests := []struct {
		lt   LinkType
		want string
	}{
		{LinkTypeInternalPage, "internal page"},
		{LinkTypeInternalAnchor, "anchor"},
		{LinkTypeInternalAsset, "asset"},
		{LinkTypeExternal, "external"},
		{LinkTypeEmail, "email"},
		{LinkTypeOther, "other"},
	}

	for _, tt := range tests {
		if got := tt.lt.String(); got != tt.want {
			t.Errorf("LinkType(%d).String() = %q, want %q", tt.lt, got, tt.want)
		}
	}
}
