package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studiowebux/minimaldoc/internal/cli"
	"github.com/studiowebux/minimaldoc/internal/version"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "minimaldoc",
	Short: "Minimal Doc - A modern static site generator for documentation",
	Long: `Minimal Doc is a fast, minimal static site generator for documentation.

It converts your Markdown files into a beautiful, responsive documentation site
with automatic navigation, search, dark mode, and more.

Full documentation: https://minimaldoc.studiowebux.com`,
	Version: version.Version,
}

func init() {
	// Prevent cobra from printing errors (we handle them in main)
	rootCmd.SilenceErrors = true

	// Add subcommands
	rootCmd.AddCommand(cli.BuildCmd)
	rootCmd.AddCommand(cli.InitCmd)
	rootCmd.AddCommand(cli.VersionCmd)

	// Set version template
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
