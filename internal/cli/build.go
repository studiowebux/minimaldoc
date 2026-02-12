package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/studiowebux/minimaldoc/internal/assets"
	"github.com/studiowebux/minimaldoc/internal/builder"
	"github.com/studiowebux/minimaldoc/internal/checker"
	"github.com/studiowebux/minimaldoc/internal/config"
	"github.com/studiowebux/minimaldoc/internal/core"
	"github.com/studiowebux/minimaldoc/internal/generator"
	"github.com/studiowebux/minimaldoc/internal/version"
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
	outputDir        string
	themeName        string
	enableLLMS       bool
	cleanURLs        bool
	siteTitle        string
	siteDesc         string
	baseURL          string
	enableOpenAPI    bool
	openapiSpecDir   string
	enableStatus     bool
	statusTitle      string
	statusPath       string
	enableChangelog  bool
	changelogTitle   string
	changelogPath    string
	enableStaleWarn  bool
	staleThreshold   int
	linkCheckMode    string
	checkExternal    bool
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
	BuildCmd.Flags().BoolVar(&enableStatus, "status", false, "Enable status page")
	BuildCmd.Flags().StringVar(&statusTitle, "status-title", "Service Status", "Status page title")
	BuildCmd.Flags().StringVar(&statusPath, "status-path", "status", "Status page output path")
	BuildCmd.Flags().BoolVar(&enableChangelog, "changelog", false, "Enable changelog")
	BuildCmd.Flags().StringVar(&changelogTitle, "changelog-title", "Changelog", "Changelog page title")
	BuildCmd.Flags().StringVar(&changelogPath, "changelog-path", "changelog", "Changelog output path")
	BuildCmd.Flags().BoolVar(&enableStaleWarn, "stale-warning", false, "Enable stale content warnings")
	BuildCmd.Flags().IntVar(&staleThreshold, "stale-threshold", 365, "Days before content is considered stale")
	BuildCmd.Flags().StringVar(&linkCheckMode, "link-check", "warn", "Link check mode: error, warn, ignore")
	BuildCmd.Flags().BoolVar(&checkExternal, "check-external", false, "Check external URLs (slower)")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Determine docs directory or single file
	docsDir := "."
	singleFile := ""
	if len(args) > 0 {
		docsDir = args[0]
	}

	// Check if argument is a single markdown file
	info, err := os.Stat(docsDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", docsDir)
	}
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	// If it's a file (not directory) and ends with .md, use single-file mode
	if !info.IsDir() && strings.HasSuffix(strings.ToLower(docsDir), ".md") {
		singleFile = filepath.Base(docsDir)
		docsDir = filepath.Dir(docsDir)
		if docsDir == "" {
			docsDir = "."
		}
		fmt.Printf("Single-file mode: %s\n", singleFile)
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

	// Create status configuration
	statusConfig := core.DefaultStatusConfig()
	statusConfig.Enabled = enableStatus
	if enableStatus {
		statusConfig.Title = statusTitle
		statusConfig.Path = statusPath
	}

	// Create changelog configuration
	changelogConfig := core.DefaultChangelogConfig()
	changelogConfig.Enabled = enableChangelog
	if enableChangelog {
		changelogConfig.Title = changelogTitle
		changelogConfig.Path = changelogPath
	}

	// Create stale warning configuration
	staleWarningConfig := core.DefaultStaleWarningConfig()
	staleWarningConfig.Enabled = enableStaleWarn
	if enableStaleWarn {
		staleWarningConfig.ThresholdDays = staleThreshold
	}

	// Create link check configuration
	linkCheckConfig := core.DefaultLinkCheckConfig()
	linkCheckConfig.Mode = core.LinkCheckMode(linkCheckMode)
	linkCheckConfig.CheckExternal = checkExternal
	if linkCheckMode == "ignore" {
		linkCheckConfig.Enabled = false
	}

	siteConfig := core.SiteConfig{
		Title:        siteTitle,
		Description:  siteDesc,
		BaseURL:      baseURL,
		Theme:        themeName,
		EnableLLMS:   enableLLMS,
		CleanURLs:    cleanURLs,
		NavDepth:     0,
		OpenAPI:      openapiConfig,
		Status:       statusConfig,
		Changelog:    changelogConfig,
		StaleWarning: staleWarningConfig,
		LinkCheck:    linkCheckConfig,
	}

	// Merge with config.yaml if it exists (CLI flags take precedence)
	if fileConfig != nil {
		siteConfig = fileConfig.MergeWithCLI(siteConfig, flagsSet)
		fmt.Println("Loaded configuration from config.yaml")
	}

	// In single-file mode, set the entrypoint to the file
	if singleFile != "" {
		siteConfig.Entrypoint = singleFile
		siteConfig.SingleFileMode = true
	}

	// Create site
	site := core.NewSite(docsDir, outputDir, siteConfig)

	// Build site (parse and process all markdown)
	siteBuilder := builder.NewBuilder(site)
	if err := siteBuilder.Build(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Generate HTML
	htmlGen, err := generator.NewHTMLGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create HTML generator: %w", err)
	}

	if err := htmlGen.Generate(); err != nil {
		return fmt.Errorf("HTML generation failed: %w", err)
	}

	// Generate OpenAPI documentation (if enabled)
	openapiGen, err := generator.NewOpenAPIGenerator(site, assets.ThemeFS, version.Version)
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

	// Generate custom theme CSS (if configured)
	themeGen := generator.NewThemeGenerator(site)
	if err := themeGen.Generate(); err != nil {
		return fmt.Errorf("theme CSS generation failed: %w", err)
	}

	// Generate status page (if enabled)
	statusGen, err := generator.NewStatusGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create status generator: %w", err)
	}
	if statusGen != nil {
		if err := statusGen.Generate(); err != nil {
			return fmt.Errorf("status page generation failed: %w", err)
		}
	}

	// Generate changelog (if enabled)
	changelogGen, err := generator.NewChangelogGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create changelog generator: %w", err)
	}
	if changelogGen != nil {
		if err := changelogGen.Generate(); err != nil {
			return fmt.Errorf("changelog generation failed: %w", err)
		}
	}

	// Generate landing page (if enabled)
	landingGen, err := generator.NewLandingGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create landing generator: %w", err)
	}
	if landingGen != nil {
		if err := landingGen.Generate(); err != nil {
			return fmt.Errorf("landing page generation failed: %w", err)
		}
	}

	// Generate portfolio (if enabled)
	portfolioGen, err := generator.NewPortfolioGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create portfolio generator: %w", err)
	}
	if portfolioGen != nil {
		if err := portfolioGen.Generate(); err != nil {
			return fmt.Errorf("portfolio generation failed: %w", err)
		}
	}

	// Generate contact page (if enabled)
	contactGen, err := generator.NewContactGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create contact generator: %w", err)
	}
	if contactGen != nil {
		if err := contactGen.Generate(); err != nil {
			return fmt.Errorf("contact page generation failed: %w", err)
		}
	}

	// Generate FAQ page (if enabled)
	faqGen, err := generator.NewFaqGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create FAQ generator: %w", err)
	}
	if faqGen != nil {
		if err := faqGen.Generate(); err != nil {
			return fmt.Errorf("FAQ page generation failed: %w", err)
		}
	}

	// Generate legal pages (if enabled)
	legalGen, err := generator.NewLegalGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create legal generator: %w", err)
	}
	if legalGen != nil {
		if err := legalGen.Generate(); err != nil {
			return fmt.Errorf("legal pages generation failed: %w", err)
		}
	}

	// Generate knowledge base (if enabled)
	kbGen, err := generator.NewKBGenerator(site, assets.ThemeFS, version.Version)
	if err != nil {
		return fmt.Errorf("failed to create KB generator: %w", err)
	}
	if kbGen != nil {
		if err := kbGen.Generate(); err != nil {
			return fmt.Errorf("knowledge base generation failed: %w", err)
		}
	}

	// Generate versioned documentation (if enabled)
	if site.Config.Versions.Enabled {
		versionGen, err := generator.NewVersionGenerator(site, assets.ThemeFS, version.Version)
		if err != nil {
			return fmt.Errorf("failed to create version generator: %w", err)
		}
		if err := versionGen.Generate(); err != nil {
			return fmt.Errorf("versioned documentation generation failed: %w", err)
		}
	}

	// Run link checker
	linkChecker := checker.NewLinkChecker(site, site.Config.LinkCheck)
	if err := linkChecker.Check(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✓ Build complete!")
	fmt.Printf("✓ Output: %s\n", outputDir)
	fmt.Println()

	return nil
}
