package checker

import (
	"fmt"
	"sort"
	"strings"

	"github.com/studiowebux/minimaldoc/static-generator/internal/core"
)

// Reporter formats and outputs link check results
type Reporter struct {
	mode core.LinkCheckMode
}

// NewReporter creates a new reporter
func NewReporter(mode core.LinkCheckMode) *Reporter {
	return &Reporter{mode: mode}
}

// Report outputs the link check results
func (r *Reporter) Report(result *core.LinkCheckResult) {
	if !result.HasErrors() {
		r.printSuccess(result)
		return
	}

	r.printErrors(result)
}

// printSuccess prints a success message
func (r *Reporter) printSuccess(result *core.LinkCheckResult) {
	fmt.Printf("Link check passed: %d links verified", result.TotalLinks)

	details := []string{}
	if result.SkippedLinks > 0 {
		details = append(details, fmt.Sprintf("%d skipped", result.SkippedLinks))
	}
	if result.ExternalLinks > 0 {
		details = append(details, fmt.Sprintf("%d external", result.ExternalLinks))
	}

	if len(details) > 0 {
		fmt.Printf(" (%s)", strings.Join(details, ", "))
	}
	fmt.Println()
}

// printErrors prints broken link errors grouped by file
func (r *Reporter) printErrors(result *core.LinkCheckResult) {
	prefix := "Warning"
	if r.mode == core.LinkCheckError {
		prefix = "Error"
	}

	fmt.Printf("\n%s: Found %d broken links\n\n", prefix, result.BrokenCount())

	// Group by source file
	byFile := make(map[string][]core.BrokenLink)
	for _, broken := range result.BrokenLinks {
		byFile[broken.Link.SourceFile] = append(byFile[broken.Link.SourceFile], broken)
	}

	// Sort files for consistent output
	files := make([]string, 0, len(byFile))
	for file := range byFile {
		files = append(files, file)
	}
	sort.Strings(files)

	// Print errors grouped by file
	for _, file := range files {
		fmt.Printf("%s:\n", file)

		// Sort by line number
		links := byFile[file]
		sort.Slice(links, func(i, j int) bool {
			return links[i].Link.Line < links[j].Link.Line
		})

		for _, broken := range links {
			r.printBrokenLink(broken)
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("Link check %s: %d broken links in %d files\n",
		r.statusWord(), result.BrokenCount(), len(byFile))
}

// printBrokenLink prints a single broken link
func (r *Reporter) printBrokenLink(broken core.BrokenLink) {
	// Format: Line X: URL (reason)
	fmt.Printf("  Line %d: %s (%s)",
		broken.Link.Line,
		broken.Link.URL,
		broken.Reason,
	)

	if broken.Suggestion != "" {
		fmt.Printf(" - %s", broken.Suggestion)
	}

	fmt.Println()
}

// statusWord returns "failed" or "completed with warnings"
func (r *Reporter) statusWord() string {
	if r.mode == core.LinkCheckError {
		return "failed"
	}
	return "completed with warnings"
}

// FormatSummary returns a one-line summary suitable for build output
func (r *Reporter) FormatSummary(result *core.LinkCheckResult) string {
	if !result.HasErrors() {
		return fmt.Sprintf("Link check: %d links OK", result.TotalLinks-result.SkippedLinks)
	}

	return fmt.Sprintf("Link check: %d broken links found", result.BrokenCount())
}
