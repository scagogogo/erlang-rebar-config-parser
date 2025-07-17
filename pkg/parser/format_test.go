package parser

import (
	"strings"
	"testing"
)

// MockTerm is a mock implementation of Term interface for testing
type MockTerm struct {
	value string
}

func (m MockTerm) String() string {
	return m.value
}

func (m MockTerm) Compare(other Term) bool {
	if otherMock, ok := other.(MockTerm); ok {
		return m.value == otherMock.value
	}
	return false
}

// TestFormat tests the formatting of Erlang terms
func TestFormat(t *testing.T) {
	// Test simple formatting
	t.Run("Simple formatting", func(t *testing.T) {
		input := `{erl_opts, [debug_info]}. {deps, [{lager, "1.0"}]}.`
		expected := `{erl_opts, [debug_info]}.

{deps, [{lager, "1.0"}]}.
`
		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse input: %v", err)
		}

		formatted := config.Format(4)
		if formatted != expected {
			t.Errorf("Format mismatch:\nExpected:\n%s\nGot:\n%s", expected, formatted)
		}

		// Self-consistency check: Parse the formatted output
		reParsedConfig, err := Parse(formatted)
		if err != nil {
			t.Fatalf("Failed to re-parse formatted output: %v", err)
		}

		// Compare original parsed config with re-parsed config
		if !compareConfigs(config, reParsedConfig) {
			t.Errorf("Re-parsed config structure differs from original")
		}
	})

	// Test nested structures formatting
	t.Run("Complex nested formatting", func(t *testing.T) {
		input := `{deps, [{cowboy, {git, "url", {tag, "2.9"}}}]}.`
		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse input: %v", err)
		}

		formatted := config.Format(2)
		// Parse the formatted output to verify it's valid
		_, err = Parse(formatted)
		if err != nil {
			t.Fatalf("Failed to re-parse formatted output: %v", err)
		}

		// Check that the formatted output contains expected substrings
		if !strings.Contains(formatted, "deps") || !strings.Contains(formatted, "cowboy") ||
			!strings.Contains(formatted, "git") || !strings.Contains(formatted, "tag") {
			t.Errorf("Formatted output missing expected content: %s", formatted)
		}
	})

	// Test mixed list formatting
	t.Run("Mixed list formatting", func(t *testing.T) {
		input := `{opts, [debug, {flag, true}, "string", 123]}.`
		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse input: %v", err)
		}

		formatted := config.Format(2)
		// Parse the formatted output to verify it's valid
		_, err = Parse(formatted)
		if err != nil {
			t.Fatalf("Failed to re-parse formatted output: %v", err)
		}

		// Check that the formatted output contains expected substrings
		if !strings.Contains(formatted, "debug") || !strings.Contains(formatted, "flag") ||
			!strings.Contains(formatted, "string") || !strings.Contains(formatted, "123") {
			t.Errorf("Formatted output missing expected content: %s", formatted)
		}
	})

	// Test quoted atoms and escapes
	t.Run("Quoted atoms and escapes", func(t *testing.T) {
		input := `{'an-atom', "a string with \"escapes\""}.`
		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse input: %v", err)
		}

		formatted := config.Format(4)
		// Parse the formatted output to verify it's valid
		_, err = Parse(formatted)
		if err != nil {
			t.Fatalf("Failed to re-parse formatted output: %v", err)
		}

		if !strings.Contains(formatted, "'an-atom'") || !strings.Contains(formatted, "escapes") {
			t.Errorf("Formatted output missing expected content: %s", formatted)
		}
	})

	// Test empty config
	t.Run("Empty config", func(t *testing.T) {
		config := &RebarConfig{}
		formatted := config.Format(4)
		if formatted != "" {
			t.Errorf("Expected empty formatted output, got: %s", formatted)
		}
	})

	// Test large indentation
	t.Run("Large indentation", func(t *testing.T) {
		input := `{opts, [a, b]}.`
		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse input: %v", err)
		}

		formatted := config.Format(8)
		// Try to format a more complex input that should result in indented lists
		if strings.Contains(formatted, "[a, b]") {
			// Simple formats may not show indentation if formatter keeps short lists on one line
			t.Log("Simple list was formatted on a single line, which is acceptable")
		} else if !strings.Contains(formatted, "        a") {
			t.Errorf("Formatted output doesn't have correct indentation: %s", formatted)
		}
	})

	// Test formatting with different indentation sizes
	t.Run("Different indentation sizes", func(t *testing.T) {
		input := `{deps, [{cowboy, {git, "url", {branch, "master"}}}, {jsx, "3.0"}]}.`
		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse input: %v", err)
		}

		// Test with 2-space indentation
		formatted2 := config.Format(2)
		_, err = Parse(formatted2)
		if err != nil {
			t.Fatalf("Failed to re-parse 2-space formatted output: %v", err)
		}

		// Test with 4-space indentation
		formatted4 := config.Format(4)
		_, err = Parse(formatted4)
		if err != nil {
			t.Fatalf("Failed to re-parse 4-space formatted output: %v", err)
		}

		// Verify that 4-space indentation has more spaces than 2-space
		if len(formatted4) <= len(formatted2) {
			t.Errorf("Expected 4-space indentation to produce longer output than 2-space")
		}
	})

	// Test complex real-world rebar.config formatting
	t.Run("Complex rebar.config", func(t *testing.T) {
		input := `{erl_opts,[debug_info,{parse_transform,lager_transform}]}.{deps,[{cowboy,"2.9.0"},{jsx,"3.0.0"}]}.{profiles,[{dev,[{deps,[{meck,"0.9.0"}]}]},{test,[{deps,[{proper,"1.3.0"}]}]}]}.`

		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse complex input: %v", err)
		}

		// Format with 4-space indentation
		formatted := config.Format(4)

		// Verify it can be parsed back
		_, err = Parse(formatted)
		if err != nil {
			t.Fatalf("Failed to re-parse formatted complex output: %v", err)
		}

		// Check for key structures in the formatted output
		expectedStructures := []string{
			"erl_opts",
			"debug_info",
			"parse_transform",
			"deps",
			"cowboy",
			"jsx",
			"profiles",
			"dev",
			"test",
			"proper",
		}

		for _, str := range expectedStructures {
			if !strings.Contains(formatted, str) {
				t.Errorf("Formatted output missing expected content: %s", str)
			}
		}

		// Check for correct spacing between terms
		lines := strings.Split(formatted, "\n")
		emptyLines := 0
		for _, line := range lines {
			if line == "" {
				emptyLines++
			}
		}

		// Should have empty lines between top-level terms
		if emptyLines < 2 {
			t.Errorf("Expected at least 2 empty lines between top-level terms, got %d", emptyLines)
		}
	})

	// Test formatting of deeply nested structures
	t.Run("Deeply nested structures", func(t *testing.T) {
		input := `{relx,[{release,{my_app,"0.1.0"},[my_app,sasl]},{dev_mode,true},{include_erts,false},{extended_start_script,true},{vm_args,"config/vm.args"},{sys_config,"config/sys.config"}]}.`

		config, err := Parse(input)
		if err != nil {
			t.Fatalf("Failed to parse deeply nested input: %v", err)
		}

		formatted := config.Format(2)

		// Verify it can be parsed back
		_, err = Parse(formatted)
		if err != nil {
			t.Fatalf("Failed to re-parse formatted nested output: %v", err)
		}

		// Check that there is proper nesting in the output
		if !strings.Contains(formatted, "  ") {
			t.Errorf("Expected indentation in nested output")
		}

		// Check for key nested elements
		if !strings.Contains(formatted, "release") ||
			!strings.Contains(formatted, "my_app") ||
			!strings.Contains(formatted, "0.1.0") {
			t.Errorf("Formatted output missing expected nested content")
		}
	})
}

// TestFormatEquivalence tests that formatting preserves semantic structure
func TestFormatEquivalence(t *testing.T) {
	testCases := []string{
		// Simple case
		`{erl_opts, [debug_info]}.`,

		// Complex case with nested structures
		`{deps, [{cowboy, {git, "https://github.com/ninenines/cowboy.git", {tag, "2.9.0"}}}, {jsx, "3.0.0"}]}.`,

		// Case with multiple terms and profiles
		`{erl_opts, [debug_info]}. {deps, [{cowboy, "2.9.0"}]}. {profiles, [{test, [{deps, [{meck, "0.9.0"}]}]}]}.`,

		// Case with quoted atoms and strings
		`{'complex-name', "value with \"quotes\""}. {empty_tuple, {}}.`,
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			// Parse original
			original, err := Parse(tc)
			if err != nil {
				t.Fatalf("Failed to parse test case: %v", err)
			}

			// Format with different indentation sizes
			for _, indent := range []int{2, 4, 8} {
				formatted := original.Format(indent)

				// Parse the formatted output
				parsed, err := Parse(formatted)
				if err != nil {
					t.Fatalf("Failed to parse formatted output with indent %d: %v", indent, err)
				}

				// Check that the parsed structures are equivalent
				if !compareConfigs(original, parsed) {
					t.Errorf("Semantic structure changed after formatting with indent %d:\nOriginal: %v\nFormatted: %s\nParsed: %v",
						indent, original, formatted, parsed)
				}
			}
		})
	}
}

// TestFormatUnknownTerm tests the default case in formatTerm function
func TestFormatUnknownTerm(t *testing.T) {
	// Create a mock term that doesn't match any known types
	mockTerm := MockTerm{value: "mock"}

	// Test formatTerm directly with the mock term
	result := formatTerm(mockTerm, 0, 2)
	expected := "UNKNOWN_TERM"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestFormatTupleDefaultHandling tests the default tuple formatting path
func TestFormatTupleDefaultHandling(t *testing.T) {
	// Create a tuple that doesn't match the special handling conditions
	// (not starting with an atom, or having complex nested structures)
	tuple := Tuple{Elements: []Term{
		Integer{Value: 123},
		String{Value: "test"},
		List{Elements: []Term{
			Atom{Value: "a"},
			Atom{Value: "b"},
			Atom{Value: "c"},
			Atom{Value: "d"}, // This makes it not simple
		}},
	}}

	result := formatTerm(tuple, 0, 2)

	// Should use the default tuple formatting (multi-line)
	if !strings.Contains(result, "{\n") {
		t.Error("Expected default tuple formatting to use multi-line format")
	}
	if !strings.Contains(result, "123") {
		t.Error("Expected tuple to contain the integer value")
	}
	if !strings.Contains(result, "\"test\"") {
		t.Error("Expected tuple to contain the string value")
	}
}

// TestFormatTupleWithMultipleElements tests tuple formatting with more than 2 elements
func TestFormatTupleWithMultipleElements(t *testing.T) {
	// Create a tuple starting with an atom and having more than 2 elements
	// This should trigger the special handling path with multiple elements
	tuple := Tuple{Elements: []Term{
		Atom{Value: "key"},
		String{Value: "value1"},
		String{Value: "value2"},
		String{Value: "value3"},
	}}

	result := formatTerm(tuple, 0, 2)

	// Should use the special atom-based formatting
	if !strings.Contains(result, "{key, ") {
		t.Error("Expected tuple to start with '{key, '")
	}
	if !strings.Contains(result, "\"value1\"") {
		t.Error("Expected tuple to contain value1")
	}
	if !strings.Contains(result, "\"value2\"") {
		t.Error("Expected tuple to contain value2")
	}
	if !strings.Contains(result, "\"value3\"") {
		t.Error("Expected tuple to contain value3")
	}
}

// TestIsSimpleTermEdgeCases tests edge cases for isSimpleTerm function
func TestIsSimpleTermEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		term     Term
		expected bool
	}{
		{
			name:     "Unknown term type",
			term:     MockTerm{value: "unknown"},
			expected: false,
		},
		{
			name: "List with 4 elements (not simple)",
			term: List{Elements: []Term{
				Atom{Value: "a"},
				Atom{Value: "b"},
				Atom{Value: "c"},
				Atom{Value: "d"},
			}},
			expected: false,
		},
		{
			name: "Tuple with 3 elements (not simple)",
			term: Tuple{Elements: []Term{
				Atom{Value: "a"},
				Atom{Value: "b"},
				Atom{Value: "c"},
			}},
			expected: false,
		},
		{
			name: "List with complex nested element",
			term: List{Elements: []Term{
				Atom{Value: "simple"},
				List{Elements: []Term{
					Atom{Value: "a"},
					Atom{Value: "b"},
					Atom{Value: "c"},
					Atom{Value: "d"}, // This makes the nested list not simple
				}},
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSimpleTerm(tt.term)
			if result != tt.expected {
				t.Errorf("isSimpleTerm(%v) = %v, expected %v", tt.term, result, tt.expected)
			}
		})
	}
}
