package parser

import (
	"testing"
)

// TestProcessEscapes tests the processing of escape sequences in strings
func TestProcessEscapes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"hello", "hello", "No escapes"},
		{"\\\"quoted\\\"", "\"quoted\"", "Double quotes"},
		{"line1\\nline2", "line1\nline2", "Newline"},
		{"tab\\tchar", "tab\tchar", "Tab"},
		{"return\\rchar", "return\rchar", "Carriage return"},
		{"back\\\\slash", "back\\slash", "Backslash"},
		{"\\\"\\n\\r\\t\\\\", "\"\n\r\t\\", "Multiple escapes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processEscapes(tt.input)
			if result != tt.expected {
				t.Errorf("processEscapes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCharClassifiers tests the character classification functions
func TestCharClassifiers(t *testing.T) {
	// Test isDigit
	t.Run("isDigit", func(t *testing.T) {
		digits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
		nonDigits := []byte{'a', 'A', '_', '-', '+', ' '}

		for _, ch := range digits {
			if !isDigit(ch) {
				t.Errorf("isDigit(%c) should be true", ch)
			}
		}

		for _, ch := range nonDigits {
			if isDigit(ch) {
				t.Errorf("isDigit(%c) should be false", ch)
			}
		}
	})

	// Test isAtomStart
	t.Run("isAtomStart", func(t *testing.T) {
		validStarts := []byte{'a', 'b', 'z', '_'}
		invalidStarts := []byte{'A', 'Z', '0', '9', '-', '+', '@'}

		for _, ch := range validStarts {
			if !isAtomStart(ch) {
				t.Errorf("isAtomStart(%c) should be true", ch)
			}
		}

		for _, ch := range invalidStarts {
			if isAtomStart(ch) {
				t.Errorf("isAtomStart(%c) should be false", ch)
			}
		}
	})

	// Test isAtomChar
	t.Run("isAtomChar", func(t *testing.T) {
		validChars := []byte{'a', 'z', 'A', 'Z', '0', '9', '_', '@'}
		invalidChars := []byte{'-', '+', ' ', ',', '.'}

		for _, ch := range validChars {
			if !isAtomChar(ch) {
				t.Errorf("isAtomChar(%c) should be true", ch)
			}
		}

		for _, ch := range invalidChars {
			if isAtomChar(ch) {
				t.Errorf("isAtomChar(%c) should be false", ch)
			}
		}
	})
}
