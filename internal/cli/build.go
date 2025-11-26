package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	minimaldoc "github.com/studiowebux/minimaldoc"
	"github.com/studiowebux/minimaldoc/internal/builder"
	"github.com/studiowebux/minimaldoc/internal/config"
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
	outputDir      string
	themeName      string
	enableLLMS     bool
	cleanURLs      bool
	siteTitle      string
	siteDesc       string
	baseURL        string
	enableOpenAPI  bool
	openapiSpecDir string
)

func init() {
	BuildCmd.Flags().StringVarP(&outputDir, "output", "o", "public", "Output directory for generated site")
	BuildCmd.Flags().StringVarP(&themeName, "theme", "t", "default", "Theme to use")
	BuildCmd.Flags().BoolVarP(&enableLLMS, "llms", "l", true, "Generate llms.txt")
	BuildCmd.Flags().BoolVar(&cleanURLs, "clean-urls", false, "Use clean URLs (/page/ instead of /page.html)")
	BuildCmd.Flags().StringVar(&siteTitle, "title", "Documentation", "Site title")
	BuildCmd.Flags().StringVar(&siteDesc, "description", "", "Site description")
	BuildCmd.Flags().StringVar(&baseURL, "base-url", "", "Base URL for the site")
	BuildCmd.Flags().BoolVar(&enableOpenAPI, "openapi", false, "Enable OpenAPI/Swagger documentation")
	BuildCmd.Flags().StringVar(&openapiSpecDir, "openapi-dir", "api", "Directory containing OpenAPI spec files (relative to docs root)")
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
	fmt.Println("║     Minimal Doc - Static Site        ║")
	fmt.Println("║          Generator                   ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	// Load config.yaml if it exists
	fileConfig, err := config.LoadConfig(docsDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Track which flags were explicitly set by the user
	flagsSet := make(map[string]bool)
	cmd.Flags().Visit(func(f *pflag.Flag) {
		flagsSet[f.Name] = true
	})

	// Create site configuration
	openapiConfig := core.DefaultOpenAPIConfig()
	openapiConfig.Enabled = enableOpenAPI
	if enableOpenAPI {
		// Auto-discover spec files in the specified directory
		openapiConfig.SpecFiles = []string{openapiSpecDir + "/*.yaml", openapiSpecDir + "/*.yml", openapiSpecDir + "/*.json"}
		openapiConfig.DefaultView = "path"
		openapiConfig.EnableTesting = true
		openapiConfig.EnableExport = true
	}

	siteConfig := core.SiteConfig{
		Title:       siteTitle,
		Description: siteDesc,
		BaseURL:     baseURL,
		Theme:       themeName,
		EnableLLMS:  enableLLMS,
		CleanURLs:   cleanURLs,
		NavDepth:    0,
		OpenAPI:     openapiConfig,
	}

	// Merge with config.yaml if it exists (CLI flags take precedence)
	if fileConfig != nil {
		siteConfig = fileConfig.MergeWithCLI(siteConfig, flagsSet)
		fmt.Println("Loaded configuration from config.yaml")
	}

	// Create site
	site := core.NewSite(docsDir, outputDir, siteConfig)

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

	// Generate OpenAPI documentation (if enabled)
	openapiGen, err := generator.NewOpenAPIGenerator(site, minimaldoc.ThemeFS)
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI generator: %w", err)
	}
	if openapiGen != nil {
		if err := openapiGen.Generate(); err != nil {
			return fmt.Errorf("OpenAPI generation failed: %w", err)
		}
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
