package version

import (
	"testing"
)

// ---------------------------------------------------------------------------
// parseVersion
// ---------------------------------------------------------------------------

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input string
		want  [3]int
	}{
		{"1.2.3", [3]int{1, 2, 3}},
		{"0.0.0", [3]int{0, 0, 0}},
		{"10.20.30", [3]int{10, 20, 30}},

		// Pre-release suffix stripped
		{"1.0.0-beta", [3]int{1, 0, 0}},
		{"2.1.3-rc.1", [3]int{2, 1, 3}},
		{"1.0.0-alpha.2", [3]int{1, 0, 0}},

		// Missing components default to 0
		{"1", [3]int{1, 0, 0}},
		{"1.2", [3]int{1, 2, 0}},

		// Empty string
		{"", [3]int{0, 0, 0}},

		// Non-numeric
		{"abc", [3]int{0, 0, 0}},
		{"1.abc.3", [3]int{1, 0, 3}},
	}

	for _, tt := range tests {
		got := parseVersion(tt.input)
		if got != tt.want {
			t.Errorf("parseVersion(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// compareVersions
// ---------------------------------------------------------------------------

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want int
	}{
		// Equal
		{"1.0.0", "1.0.0", 0},
		{"0.0.0", "0.0.0", 0},
		{"10.20.30", "10.20.30", 0},

		// a < b
		{"1.0.0", "2.0.0", -1},
		{"1.0.0", "1.1.0", -1},
		{"1.0.0", "1.0.1", -1},
		{"0.9.9", "1.0.0", -1},
		{"1.2.3", "1.2.4", -1},

		// a > b
		{"2.0.0", "1.0.0", 1},
		{"1.1.0", "1.0.0", 1},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "0.9.9", 1},
		{"1.2.4", "1.2.3", 1},

		// With pre-release (stripped, so only numeric matters)
		{"1.0.0-beta", "1.0.0", 0},
		{"1.0.0", "1.0.1-beta", -1},
	}

	for _, tt := range tests {
		got := compareVersions(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// truncateNotes
// ---------------------------------------------------------------------------

func TestTruncateNotes(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		// Short enough — no truncation
		{"hello", 10, "hello"},
		{"", 10, ""},

		// Exact length — no truncation
		{"12345", 5, "12345"},

		// Truncated with ellipsis
		{"1234567890", 5, "12345..."},
		{"hello world", 5, "hello..."},
	}

	for _, tt := range tests {
		got := truncateNotes(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncateNotes(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestTruncateNotes_UTF8(t *testing.T) {
	// 5 runes, each multi-byte
	input := "aebec" // plain ASCII for predictable test
	got := truncateNotes(input, 3)
	if got != "aeb..." {
		t.Errorf("truncateNotes(%q, 3) = %q, want %q", input, got, "aeb...")
	}

	// Unicode: truncate at rune boundary, not byte boundary
	unicode := "\u00e9\u00e9\u00e9\u00e9\u00e9" // 5x e-acute
	got2 := truncateNotes(unicode, 3)
	want2 := "\u00e9\u00e9\u00e9..."
	if got2 != want2 {
		t.Errorf("truncateNotes(unicode, 3) = %q, want %q", got2, want2)
	}
}

// ---------------------------------------------------------------------------
// Version constant
// ---------------------------------------------------------------------------

func TestVersionConstant_NotEmpty(t *testing.T) {
	if Version == "" {
		t.Error("Version constant should not be empty")
	}
}

func TestVersionConstant_Parseable(t *testing.T) {
	parts := parseVersion(Version)
	// At least major should be > 0 for a released version
	if parts[0] == 0 && parts[1] == 0 && parts[2] == 0 {
		t.Errorf("Version %q parses to [0,0,0] — unexpected for a released version", Version)
	}
}
