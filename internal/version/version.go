package version

// Version is the current version of MinimalDoc.
// Set at build time via ldflags:
//
//	go build -ldflags "-X github.com/studiowebux/minimaldoc/internal/version.Version=1.8.0"
//
// Falls back to "dev" for local development builds without ldflags.
var Version = "dev"
