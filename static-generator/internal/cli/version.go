package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/studiowebux/minimaldoc/static-generator/internal/version"
)

var checkUpdate bool

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Display the current version of MinimalDoc.

Use --check to check for updates from GitHub releases.`,
	RunE: runVersion,
}

func init() {
	VersionCmd.Flags().BoolVar(&checkUpdate, "check", false, "Check for updates from GitHub")
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("minimaldoc version %s\n", version.Version)

	if checkUpdate {
		fmt.Println()
		fmt.Println("Checking for updates...")

		result, err := version.CheckForUpdate()
		if err != nil {
			fmt.Printf("Failed to check for updates: %v\n", err)
			return nil
		}

		if result.UpdateAvailable {
			fmt.Println()
			fmt.Printf("Update available: %s -> %s\n", result.CurrentVersion, result.LatestVersion)
			fmt.Printf("Release URL: %s\n", result.ReleaseURL)
			fmt.Println()
			fmt.Println("To update:")
			fmt.Println("  go install github.com/studiowebux/minimaldoc/static-generator/cmd/minimaldoc@latest")
			fmt.Println()
			fmt.Println("Or download from GitHub releases.")
		} else {
			fmt.Println()
			fmt.Println("You are running the latest version.")
		}
	}

	return nil
}
