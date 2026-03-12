package generator

import "os"

// writeWebFile writes data to a file with web-readable permissions (0644).
// Web output files must be world-readable for HTTP servers to serve them.
//
//nolint:gosec // G306: web output files must be world-readable
func writeWebFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644) // #nosec G306
}

// makeWebDir creates a directory with web-accessible permissions (0755).
// Web output directories must be world-executable for HTTP servers to traverse them.
//
//nolint:gosec // G301: web output directories must be world-accessible
func makeWebDir(path string) error {
	return os.MkdirAll(path, 0755) // #nosec G301
}
