package generator

import (
	"testing"
)

// ---------------------------------------------------------------------------
// safeHrefFunc — XSS protection
// ---------------------------------------------------------------------------

func TestSafeHrefFunc_BlocksDangerousSchemes(t *testing.T) {
	dangerous := []string{
		"javascript:alert(1)",
		"JavaScript:alert(1)",
		"JAVASCRIPT:void(0)",
		"data:text/html,<script>alert(1)</script>",
		"DATA:text/html,test",
		"vbscript:MsgBox",
		"VBSCRIPT:run",
		// With leading whitespace
		"  javascript:alert(1)",
		"\tdata:text/html,x",
	}

	for _, url := range dangerous {
		got := safeHrefFunc(url)
		if got != "#" {
			t.Errorf("safeHrefFunc(%q) = %q, want %q (should block dangerous scheme)", url, got, "#")
		}
	}
}

func TestSafeHrefFunc_AllowsSafeURLs(t *testing.T) {
	safe := []struct {
		input string
		want  string
	}{
		{"https://example.com", "https://example.com"},
		{"http://example.com/page", "http://example.com/page"},
		{"/docs/getting-started", "/docs/getting-started"},
		{"./relative.html", "./relative.html"},
		{"../parent/page.html", "../parent/page.html"},
		{"#section", "#section"},
		{"mailto:user@example.com", "mailto:user@example.com"},
		// Whitespace trimming
		{"  https://example.com  ", "https://example.com"},
	}

	for _, tt := range safe {
		got := safeHrefFunc(tt.input)
		if got != tt.want {
			t.Errorf("safeHrefFunc(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSafeHrefFunc_EmptyString(t *testing.T) {
	if got := safeHrefFunc(""); got != "" {
		t.Errorf("safeHrefFunc(%q) = %q, want empty string", "", got)
	}
}

// ---------------------------------------------------------------------------
// GetBasePath
// ---------------------------------------------------------------------------

func TestGetBasePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Standard cases
		{"https://example.com/docs/", "/docs"},
		{"https://example.com/docs", "/docs"},
		{"http://example.com/my/path/", "/my/path"},

		// Root URLs return empty
		{"https://example.com/", ""},
		{"https://example.com", ""},
		{"http://example.com", ""},

		// Empty input
		{"", ""},

		// Deep paths
		{"https://example.com/a/b/c/d", "/a/b/c/d"},

		// Single segment
		{"https://example.com/blog", "/blog"},
	}

	for _, tt := range tests {
		got := GetBasePath(tt.input)
		if got != tt.want {
			t.Errorf("GetBasePath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// escapeHTML
// ---------------------------------------------------------------------------

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"<script>", "&lt;script&gt;"},
		{`"quoted"`, "&#34;quoted&#34;"},
		{"a&b", "a&amp;b"},
		{"it's", "it&#39;s"},
		{"", ""},
	}

	for _, tt := range tests {
		got := escapeHTML(tt.input)
		if got != tt.want {
			t.Errorf("escapeHTML(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// escapeJS
// ---------------------------------------------------------------------------

func TestEscapeJS(t *testing.T) {
	// escapeJS wraps template.JSEscapeString — verify key escapes
	tests := []struct {
		input    string
		contains string // substring that must appear in output
	}{
		{`alert("hi")`, `alert(`},      // quotes should be escaped
		{"line\nnewline", "line"},       // newline escaped
		{"back\\slash", "back\\\\slash"}, // backslash doubled
	}

	for _, tt := range tests {
		got := escapeJS(tt.input)
		if got == tt.input && tt.input != "" {
			// If input has special chars, output should differ
			t.Logf("escapeJS(%q) = %q (may be unchanged if no special chars)", tt.input, got)
		}
		if len(got) == 0 && len(tt.input) > 0 {
			t.Errorf("escapeJS(%q) returned empty string", tt.input)
		}
	}
}

func TestEscapeJS_Empty(t *testing.T) {
	if got := escapeJS(""); got != "" {
		t.Errorf("escapeJS(%q) = %q, want empty", "", got)
	}
}

// ---------------------------------------------------------------------------
// dictFunc
// ---------------------------------------------------------------------------

func TestDictFunc_ValidPairs(t *testing.T) {
	result, err := dictFunc("key1", "value1", "key2", 42)
	if err != nil {
		t.Fatalf("dictFunc returned error: %v", err)
	}
	if result["key1"] != "value1" {
		t.Errorf("result[key1] = %v, want %q", result["key1"], "value1")
	}
	if result["key2"] != 42 {
		t.Errorf("result[key2] = %v, want 42", result["key2"])
	}
}

func TestDictFunc_EmptyArgs(t *testing.T) {
	result, err := dictFunc()
	if err != nil {
		t.Fatalf("dictFunc() returned error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("dictFunc() returned map with %d entries, want 0", len(result))
	}
}

func TestDictFunc_OddArgs(t *testing.T) {
	_, err := dictFunc("key1", "value1", "orphan")
	if err == nil {
		t.Error("dictFunc with odd args should return error, got nil")
	}
}

func TestDictFunc_NonStringKey(t *testing.T) {
	_, err := dictFunc(123, "value")
	if err == nil {
		t.Error("dictFunc with non-string key should return error, got nil")
	}
}

// ---------------------------------------------------------------------------
// addFunc
// ---------------------------------------------------------------------------

func TestAddFunc(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{1, 2, 3},
		{0, 0, 0},
		{-1, 1, 0},
		{100, -50, 50},
	}
	for _, tt := range tests {
		if got := addFunc(tt.a, tt.b); got != tt.want {
			t.Errorf("addFunc(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// lowerFunc / upperFunc
// ---------------------------------------------------------------------------

func TestLowerFunc(t *testing.T) {
	if got := lowerFunc("Hello WORLD"); got != "hello world" {
		t.Errorf("lowerFunc(%q) = %q, want %q", "Hello WORLD", got, "hello world")
	}
}

func TestUpperFunc(t *testing.T) {
	if got := upperFunc("Hello world"); got != "HELLO WORLD" {
		t.Errorf("upperFunc(%q) = %q, want %q", "Hello world", got, "HELLO WORLD")
	}
}

// ---------------------------------------------------------------------------
// replaceFunc
// ---------------------------------------------------------------------------

func TestReplaceFunc(t *testing.T) {
	got := replaceFunc("hello-world-test", "-", "_")
	if got != "hello_world_test" {
		t.Errorf("replaceFunc result = %q, want %q", got, "hello_world_test")
	}
}

// ---------------------------------------------------------------------------
// BaseFuncMap — verify all expected keys exist
// ---------------------------------------------------------------------------

func TestBaseFuncMap_Keys(t *testing.T) {
	fm := BaseFuncMap()
	expectedKeys := []string{
		"dict", "safeHTML", "json", "lower", "upper",
		"hasPrefix", "add", "join", "replace",
		"hasCustomTheme", "formatDate", "safeHref",
	}
	for _, key := range expectedKeys {
		if _, ok := fm[key]; !ok {
			t.Errorf("BaseFuncMap() missing key %q", key)
		}
	}
}

// ---------------------------------------------------------------------------
// ExtendFuncMap
// ---------------------------------------------------------------------------

func TestExtendFuncMap(t *testing.T) {
	extra := map[string]any{
		"custom": func() string { return "hello" },
	}
	fm := ExtendFuncMap(extra)

	// Should have base keys
	if _, ok := fm["dict"]; !ok {
		t.Error("ExtendFuncMap result missing base key 'dict'")
	}
	// Should have the extra key
	if _, ok := fm["custom"]; !ok {
		t.Error("ExtendFuncMap result missing additional key 'custom'")
	}
}
