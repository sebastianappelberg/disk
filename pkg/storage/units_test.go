package storage

import (
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{512, "512B"},
		{1024, "1kB"},
		{1536, "1kB"}, // Should round down (integer division)
		{1048576, "1MB"},
		{2097152, "2MB"},
		{1073741824, "1GB"},
		{2147483648, "2GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatSize(tt.input)
			if result != tt.expected {
				t.Errorf("FormatSize(%d) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		input       string
		expected    int64
		expectError bool
	}{
		{"0B", 0, false},
		{"512B", 512, false},
		{"1kB", 1024, false},
		{"2kB", 2048, false},
		{"1MB", 1048576, false},
		{"2MB", 2097152, false},
		{"1GB", 1073741824, false},
		{"3GB", 3221225472, false},
		{"100B", 100, false},
		{"1234kB", 1263616, false},
		{"InvalidSize", 0, true}, // Invalid string
		{"5TB", 0, true},         // Unsupported unit
		{"-1MB", 0, true},        // Negative value
		{"1.5MB", 0, true},       // Decimal values not supported
		{"ABCkB", 0, true},       // Non-numeric values
		{"1B1", 0, true},         // Invalid format
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseSize(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParseSize(%q) expected an error, but got none. Got %d", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("ParseSize(%q) unexpected error: %v", tt.input, err)
				} else if result != tt.expected {
					t.Errorf("ParseSize(%q) = %d; want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestSplitToInt(t *testing.T) {
	tests := []struct {
		input       string
		unit        string
		expected    int64
		expectError bool
	}{
		{"100kB", "kB", 100, false},
		{"5MB", "MB", 5, false},
		{"0GB", "GB", 0, false},
		{"123B", "B", 123, false},
		{"ABCkB", "kB", 0, true},   // Invalid number
		{"100.5kB", "kB", 0, true}, // Decimal value
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := splitToInt(tt.input, tt.unit)
			if tt.expectError {
				if err == nil {
					t.Errorf("splitToInt(%q, %q) expected an error, but got none", tt.input, tt.unit)
				}
			} else {
				if err != nil {
					t.Errorf("splitToInt(%q, %q) unexpected error: %v", tt.input, tt.unit, err)
				} else if result != tt.expected {
					t.Errorf("splitToInt(%q, %q) = %d; want %d", tt.input, tt.unit, result, tt.expected)
				}
			}
		})
	}
}
