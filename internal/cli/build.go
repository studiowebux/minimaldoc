package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	minimaldoc "github.com/studiowebux/minimaldoc"
	"github.com/studiowebux/minimaldoc/internal/builder"
	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/studiowebux/minimaldoc/internal/generator"
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build [docs-directory]",
	Short: "Build the static documentation site",
	Long: `Build generates a static HTML site from your markdown documentation.

By default, it looks for markdown files in the current directory and outputs
to the 'public' directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBuild,
}

var (
	outputDir   string
	themeName   string
	enableLLMS  bool
	cleanURLs   bool
	siteTitle   string
	siteDesc    string
	baseURL     string
)

func init() {
	BuildCmd.Flags().StringVarP(&outputDir, "output", "o", "public", "Output directory for generated site")
	BuildCmd.Flags().StringVarP(&themeName, "theme", "t", "default", "Theme to use")
	BuildCmd.Flags().BoolVarP(&enableLLMS, "llms", "l", true, "Generate llms.txt")
	BuildCmd.Flags().BoolVar(&cleanURLs, "clean-urls", false, "Use clean URLs (/page/ instead of /page.html)")
	BuildCmd.Flags().StringVar(&siteTitle, "title", "Documentation", "Site title")
	BuildCmd.Flags().StringVar(&siteDesc, "description", "", "Site description")
	BuildCmd.Flags().StringVar(&baseURL, "base-url", "", "Base URL for the site")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Determine docs directory
	docsDir := "."
	if len(args) > 0 {
		docsDir = args[0]
	}

	// Verify docs directory exists
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		return fmt.Errorf("docs directory does not exist: %s", docsDir)
	}

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║     Minimal Doc - Static Site       ║")
	fmt.Println("║          Generator                   ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	// Create site configuration
	config := core.SiteConfig{
		Title:       siteTitle,
		Description: siteDesc,
		BaseURL:     baseURL,
		Theme:       themeName,
		EnableLLMS:  enableLLMS,
		CleanURLs:   cleanURLs,
		NavDepth:    0,
	}

	// Create site
	site := core.NewSite(docsDir, outputDir, config)

	// Build site (parse and process all markdown)
	siteBuilder := builder.NewBuilder(site)
	if err := siteBuilder.Build(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Generate HTML
	htmlGen, err := generator.NewHTMLGenerator(site, minimaldoc.ThemeFS)
	if err != nil {
		return fmt.Errorf("failed to create HTML generator: %w", err)
	}

	if err := htmlGen.Generate(); err != nil {
		return fmt.Errorf("HTML generation failed: %w", err)
	}

	// Generate llms.txt
	if enableLLMS {
		llmsGen := generator.NewLLMSGenerator(site)
		if err := llmsGen.Generate(); err != nil {
			return fmt.Errorf("LLMS generation failed: %w", err)
		}
	}

	// Generate sitemap.xml
	sitemapGen := generator.NewSitemapGenerator(site)
	if err := sitemapGen.Generate(); err != nil {
		return fmt.Errorf("sitemap generation failed: %w", err)
	}

	// Generate search index
	searchGen := generator.NewSearchGenerator(site)
	if err := searchGen.Generate(); err != nil {
		return fmt.Errorf("search index generation failed: %w", err)
	}

	fmt.Println()
	fmt.Println("✓ Build complete!")
	fmt.Printf("✓ Output: %s\n", outputDir)
	fmt.Println()

	return nil
}
