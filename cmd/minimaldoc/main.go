package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studiowebux/minimaldoc/internal/cli"
)

var (
	version = "0.0.2"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "minimaldoc",
	Short: "Minimal Doc - A modern static site generator for documentation",
	Long: `Minimal Doc is a fast, minimal static site generator for documentation.

It converts your Markdown files into a beautiful, responsive documentation website
with automatic navigation, table of contents, and dark mode support.

Features:
  • Automatic navigation from folder structure
  • Table of contents with scrollspy
  • Light/Dark theme toggle
  • LLM-friendly markdown export
  • Clean, minimal design
  • SEO optimized
  • GitHub Pages ready`,
	Version: version,
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(cli.BuildCmd)
	rootCmd.AddCommand(cli.InitCmd)

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
