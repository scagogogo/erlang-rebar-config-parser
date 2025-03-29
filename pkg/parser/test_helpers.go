package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// compareConfigs compares two RebarConfig structs by comparing their terms
// This is a common helper used across different test files
func compareConfigs(c1, c2 *RebarConfig) bool {
	if len(c1.Terms) != len(c2.Terms) {
		return false
	}
	for i := range c1.Terms {
		if !c1.Terms[i].Compare(c2.Terms[i]) {
			return false
		}
	}
	return true
}

// createTempConfigFile creates a temporary file with the given content and returns its path
func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_rebar.config")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return filePath
}
