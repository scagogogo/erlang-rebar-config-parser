package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCompareConfigs tests the compareConfigs helper function
func TestCompareConfigs(t *testing.T) {
	t.Run("Equal configs", func(t *testing.T) {
		config1, _ := Parse(`{erl_opts, [debug_info]}.`)
		config2, _ := Parse(`{erl_opts, [debug_info]}.`)

		if !compareConfigs(config1, config2) {
			t.Error("Expected configs to be equal")
		}
	})

	t.Run("Different number of terms", func(t *testing.T) {
		config1, _ := Parse(`{erl_opts, [debug_info]}.`)
		config2, _ := Parse(`{erl_opts, [debug_info]}. {deps, []}.`)

		if compareConfigs(config1, config2) {
			t.Error("Expected configs to be different (different number of terms)")
		}
	})

	t.Run("Different term values", func(t *testing.T) {
		config1, _ := Parse(`{erl_opts, [debug_info]}.`)
		config2, _ := Parse(`{erl_opts, [warnings_as_errors]}.`)

		if compareConfigs(config1, config2) {
			t.Error("Expected configs to be different (different term values)")
		}
	})

	t.Run("Empty configs", func(t *testing.T) {
		config1, _ := Parse(``)
		config2, _ := Parse(``)

		if !compareConfigs(config1, config2) {
			t.Error("Expected empty configs to be equal")
		}
	})
}

// TestCreateTempConfigFile tests the createTempConfigFile helper function
func TestCreateTempConfigFile(t *testing.T) {
	t.Run("Create valid temp file", func(t *testing.T) {
		content := `{erl_opts, [debug_info]}.`
		filePath := createTempConfigFile(t, content)

		// Check that file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Temp file was not created: %s", filePath)
		}

		// Check file content
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}

		if string(fileContent) != content {
			t.Errorf("File content mismatch. Expected: %s, Got: %s", content, string(fileContent))
		}

		// Check file is in temp directory
		if !filepath.IsAbs(filePath) {
			t.Error("Expected absolute path for temp file")
		}
	})

	t.Run("Create empty temp file", func(t *testing.T) {
		content := ""
		filePath := createTempConfigFile(t, content)

		// Check that file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Temp file was not created: %s", filePath)
		}

		// Check file content
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}

		if string(fileContent) != content {
			t.Errorf("File content mismatch. Expected empty, Got: %s", string(fileContent))
		}
	})

	t.Run("Create temp file with complex content", func(t *testing.T) {
		content := `{erl_opts, [debug_info, warnings_as_errors]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "3.1.0"}
]}.
{profiles, [
    {test, [{deps, [{meck, "0.9.0"}]}]}
]}.`
		filePath := createTempConfigFile(t, content)

		// Check that file exists and can be parsed
		config, err := ParseFile(filePath)
		if err != nil {
			t.Fatalf("Failed to parse created temp file: %v", err)
		}

		if len(config.Terms) != 3 {
			t.Errorf("Expected 3 terms in parsed config, got %d", len(config.Terms))
		}
	})

	t.Run("Create temp file with very large content", func(t *testing.T) {
		// Create a very large content to test edge cases
		var builder strings.Builder
		for i := 0; i < 1000; i++ {
			builder.WriteString(fmt.Sprintf("{dep%d, \"version%d\"}.\n", i, i))
		}
		content := builder.String()

		filePath := createTempConfigFile(t, content)

		// Check that file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Large temp file was not created: %s", filePath)
		}

		// Check file can be read
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read large temp file: %v", err)
		}

		if len(fileContent) == 0 {
			t.Error("Large temp file is empty")
		}
	})
}
