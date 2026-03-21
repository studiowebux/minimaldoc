package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckOutputDir_NonExistent(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "does-not-exist")

	err := checkOutputDir(dir, false)
	if err != nil {
		t.Errorf("non-existent dir should be safe, got: %v", err)
	}
}

func TestCheckOutputDir_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	err := checkOutputDir(dir, false)
	if err != nil {
		t.Errorf("empty dir should be safe, got: %v", err)
	}
}

func TestCheckOutputDir_WithMarker(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, buildMarkerFile), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html></html>"), 0o644)

	err := checkOutputDir(dir, false)
	if err != nil {
		t.Errorf("dir with marker should be safe, got: %v", err)
	}
}

func TestCheckOutputDir_WithFilesNoMarker(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "important.txt"), []byte("data"), 0o644)

	err := checkOutputDir(dir, false)
	if err == nil {
		t.Error("dir with files but no marker should fail without --force")
	}
}

func TestCheckOutputDir_WithFilesNoMarkerForce(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "important.txt"), []byte("data"), 0o644)

	err := checkOutputDir(dir, true)
	if err != nil {
		t.Errorf("dir with files and --force should be safe, got: %v", err)
	}
}

func TestCheckOutputDir_NotADirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "afile")
	os.WriteFile(file, []byte("data"), 0o644)

	err := checkOutputDir(file, false)
	if err == nil {
		t.Error("path that is a file should fail")
	}
}
