package generator

import (
	"fmt"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// BuildFooter creates a footer config with all auto-generated content.
// It starts with the user's config and adds:
// - Legal links group (if legal pages exist)
// - Version badge (unless HideVersion is true)
func BuildFooter(site *core.Site, version string) core.FooterConfig {
	footer := site.Config.Footer

	// Auto-generate legal links group
	if site.Config.Legal.Enabled && len(site.LegalPages) > 0 {
		legalPath := featurePath(site.Config.Legal.Path, "legal")

		groupTitle := site.Config.Legal.FooterGroup
		if groupTitle == "" {
			groupTitle = "Legal"
		}

		var legalLinks []core.FooterLink
		for _, page := range site.LegalPages {
			legalLinks = append(legalLinks, core.FooterLink{
				Text: page.Title,
				URL:  "/" + legalPath + "/" + page.Slug + "/",
			})
		}

		footer.Links = append(footer.Links, core.FooterLinkGroup{
			Title: groupTitle,
			Items: legalLinks,
		})
	}

	// Auto-append version badge
	if !footer.HideVersion {
		footer.Badges = append(footer.Badges, core.FooterBadge{
			Text: fmt.Sprintf("minimaldoc v%s", version),
			URL:  "https://minimaldoc.com",
		})
	}

	return footer
}
