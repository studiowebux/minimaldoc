package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/studiowebux/minimaldoc/internal/core"
)

func TestExtractOrder(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"01-intro", 1},
		{"02-guide", 2},
		{"10_advanced", 10},
		{"99-last", 99},
		{"guide", 999},
		{"README", 999},
		{"", 999},
		{"0-zero", 0},
		{"123-big", 123},
	}

	for _, tt := range tests {
		got := extractOrder(tt.input)
		if got != tt.want {
			t.Errorf("extractOrder(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestCleanPathSegment(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"01-getting-started", "Getting Started"},
		{"02_API", "API"},
		{"guide", "Guide"},
		{"my-cool-feature", "My Cool Feature"},
		{"03-hello_world", "Hello World"},
		{"README", "README"},
		{"", ""},
	}

	for _, tt := range tests {
		got := cleanPathSegment(tt.input)
		if got != tt.want {
			t.Errorf("cleanPathSegment(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNavigationBuildTree(t *testing.T) {
	nb := NewNavigationBuilder()

	pages := []*core.Page{
		{RelPath: "index.md", Slug: "index", Metadata: core.Metadata{Title: "Home"}, Order: 0},
		{RelPath: filepath.Join("getting-started", "install.md"), Slug: "getting-started/install", Metadata: core.Metadata{Title: "Installation"}, Order: 1},
		{RelPath: filepath.Join("getting-started", "config.md"), Slug: "getting-started/config", Metadata: core.Metadata{Title: "Configuration"}, Order: 2},
		{RelPath: "changelog.md", Slug: "changelog", Metadata: core.Metadata{Title: "Changelog"}, Order: 99},
	}

	items := nb.buildTree(pages, 0)

	if len(items) == 0 {
		t.Fatal("buildTree returned empty items")
	}

	// Should have root pages + a "Getting Started" section
	foundHome := false
	foundSection := false
	for _, item := range items {
		if item.Title == "Home" {
			foundHome = true
		}
		if item.Title == "Getting Started" {
			foundSection = true
			if len(item.Children) != 2 {
				t.Errorf("Getting Started section has %d children, want 2", len(item.Children))
			}
		}
	}

	if !foundHome {
		t.Error("Home page not found in navigation")
	}
	if !foundSection {
		t.Error("Getting Started section not found in navigation")
	}
}

func TestNavigationBuildTreeMaxDepth(t *testing.T) {
	nb := NewNavigationBuilder()

	pages := []*core.Page{
		{RelPath: "index.md", Slug: "index", Metadata: core.Metadata{Title: "Home"}, Order: 0},
		{RelPath: filepath.Join("level1", "page.md"), Slug: "level1/page", Metadata: core.Metadata{Title: "Level 1 Page"}, Order: 1},
		{RelPath: filepath.Join("level1", "level2", "deep.md"), Slug: "level1/level2/deep", Metadata: core.Metadata{Title: "Deep Page"}, Order: 1},
	}

	// maxDepth=1 should only show root level
	items := nb.buildTree(pages, 1)

	for _, item := range items {
		if len(item.Children) > 0 {
			t.Errorf("maxDepth=1 should not have children, but %q has %d", item.Title, len(item.Children))
		}
	}
}

func TestNavigationBuildSorting(t *testing.T) {
	nb := NewNavigationBuilder()

	pages := []*core.Page{
		{RelPath: "zebra.md", Slug: "zebra", Metadata: core.Metadata{Title: "Zebra"}, Order: 999},
		{RelPath: "alpha.md", Slug: "alpha", Metadata: core.Metadata{Title: "Alpha"}, Order: 999},
		{RelPath: "first.md", Slug: "first", Metadata: core.Metadata{Title: "First"}, Order: 1},
	}

	items := nb.buildTree(pages, 0)

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	// First should come first (order=1), then Alpha before Zebra (alphabetical at order=999)
	if items[0].Title != "First" {
		t.Errorf("items[0] = %q, want First", items[0].Title)
	}
	if items[1].Title != "Alpha" {
		t.Errorf("items[1] = %q, want Alpha", items[1].Title)
	}
	if items[2].Title != "Zebra" {
		t.Errorf("items[2] = %q, want Zebra", items[2].Title)
	}
}

func TestDiscoverPages(t *testing.T) {
	// Create a temporary docs directory with sample files
	tmpDir := t.TempDir()

	// Create docs structure
	docsDir := filepath.Join(tmpDir, "docs")
	os.MkdirAll(filepath.Join(docsDir, "getting-started"), 0o755)
	os.MkdirAll(filepath.Join(docsDir, "__status__", "incidents"), 0o755)
	os.MkdirAll(filepath.Join(docsDir, "__portfolio__"), 0o755)

	// Regular pages
	os.WriteFile(filepath.Join(docsDir, "index.md"), []byte("---\ntitle: Home\n---\n# Home\n"), 0o644)
	os.WriteFile(filepath.Join(docsDir, "getting-started", "install.md"), []byte("---\ntitle: Install\n---\n# Install\n"), 0o644)

	// Special directory files (should be skipped by discoverPages)
	os.WriteFile(filepath.Join(docsDir, "__status__", "incidents", "2025-01-01-test.md"), []byte("---\ntitle: Test\n---\n"), 0o644)
	os.WriteFile(filepath.Join(docsDir, "__portfolio__", "project.md"), []byte("---\ntitle: Project\n---\n"), 0o644)

	// Non-markdown file (should be skipped)
	os.WriteFile(filepath.Join(docsDir, "config.yaml"), []byte("title: Test\n"), 0o644)

	outputDir := filepath.Join(tmpDir, "public")

	site := core.NewSite(docsDir, outputDir, core.DefaultSiteConfig())
	b := NewBuilder(site)

	err := b.discoverPages()
	if err != nil {
		t.Fatalf("discoverPages() error: %v", err)
	}

	// Should find exactly 2 pages (index.md + install.md)
	if len(site.Pages) != 2 {
		t.Errorf("discoverPages found %d pages, want 2", len(site.Pages))
		for _, p := range site.Pages {
			t.Logf("  page: %s", p.RelPath)
		}
	}

	// Verify slugs
	slugs := make(map[string]bool)
	for _, p := range site.Pages {
		slugs[p.Slug] = true
	}

	if !slugs["index"] {
		t.Error("missing index page")
	}
	if !slugs["getting-started/install"] {
		t.Error("missing getting-started/install page")
	}
}

func TestGetOutputPath(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(docsDir, 0o755)

	// Test standard mode (no clean URLs)
	cfg := core.DefaultSiteConfig()
	site := core.NewSite(docsDir, outputDir, cfg)
	b := NewBuilder(site)

	page := &core.Page{Slug: "getting-started/install"}
	got := b.getOutputPath(page)
	want := filepath.Join(outputDir, "getting-started", "install.html")
	if got != want {
		t.Errorf("getOutputPath (standard) = %q, want %q", got, want)
	}

	// Test clean URLs mode
	cfg.CleanURLs = true
	site = core.NewSite(docsDir, outputDir, cfg)
	b = NewBuilder(site)

	got = b.getOutputPath(page)
	want = filepath.Join(outputDir, "getting-started", "install", "index.html")
	if got != want {
		t.Errorf("getOutputPath (clean URLs) = %q, want %q", got, want)
	}

	// Test index page
	indexPage := &core.Page{Slug: "index"}
	got = b.getOutputPath(indexPage)
	want = filepath.Join(outputDir, "index.html")
	if got != want {
		t.Errorf("getOutputPath (index) = %q, want %q", got, want)
	}
}

func TestComputePrevNext(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(docsDir, 0o755)

	site := core.NewSite(docsDir, outputDir, core.DefaultSiteConfig())

	p1 := &core.Page{Slug: "page1", Metadata: core.Metadata{Title: "Page 1"}, Order: 1}
	p2 := &core.Page{Slug: "page2", Metadata: core.Metadata{Title: "Page 2"}, Order: 2}
	p3 := &core.Page{Slug: "page3", Metadata: core.Metadata{Title: "Page 3"}, Order: 3}

	site.Pages = []*core.Page{p1, p2, p3}
	site.Navigation = &core.Navigation{
		Items: []*core.NavItem{
			{Page: p1, Children: []*core.NavItem{}},
			{Page: p2, Children: []*core.NavItem{}},
			{Page: p3, Children: []*core.NavItem{}},
		},
	}

	b := NewBuilder(site)
	b.computePrevNext()

	if p1.Prev != nil {
		t.Error("first page should have nil Prev")
	}
	if p1.Next != p2 {
		t.Error("first page Next should be page2")
	}
	if p2.Prev != p1 {
		t.Error("page2 Prev should be page1")
	}
	if p2.Next != p3 {
		t.Error("page2 Next should be page3")
	}
	if p3.Next != nil {
		t.Error("last page should have nil Next")
	}
}
