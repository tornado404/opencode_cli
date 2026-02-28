package util

import (
	"os"
	"testing"

	"github.com/anomalyco/oho/internal/config"
)

func TestMain(m *testing.M) {
	// Initialize config for tests
	os.Setenv("OPENCODE_SERVER_HOST", "127.0.0.1")
	os.Setenv("OPENCODE_SERVER_PORT", "4096")
	os.Setenv("OPENCODE_SERVER_USERNAME", "opencode")
	os.Setenv("OPENCODE_SERVER_PASSWORD", "test")
	config.Init()

	m.Run()
}

func TestOutputJSON(t *testing.T) {
	// Save original config
	origJSON := config.Get().JSON
	defer func() { config.Get().JSON = origJSON }()

	// Enable JSON mode
	config.Get().JSON = true

	// Test JSON output
	data := map[string]string{"key": "value"}
	err := OutputJSON(data)
	if err != nil {
		t.Errorf("OutputJSON failed: %v", err)
	}
}

func TestOutputText(t *testing.T) {
	// Save original config
	origJSON := config.Get().JSON
	defer func() { config.Get().JSON = origJSON }()

	// Disable JSON mode
	config.Get().JSON = false

	// Test text output - should not panic
	OutputText("test %s", "value")
}

func TestOutputLine(t *testing.T) {
	// Save original config
	origJSON := config.Get().JSON
	defer func() { config.Get().JSON = origJSON }()

	// Disable JSON mode
	config.Get().JSON = false

	// Test line output - should not panic
	OutputLine("test line")
}

func TestOutputTable(t *testing.T) {
	// Save original config
	origJSON := config.Get().JSON
	defer func() { config.Get().JSON = origJSON }()

	// Disable JSON mode
	config.Get().JSON = false

	headers := []string{"Name", "Age"}
	rows := [][]string{
		{"Alice", "30"},
		{"Bob", "25"},
	}

	// Should not panic
	OutputTable(headers, rows)
}

func TestOutputTableJSON(t *testing.T) {
	// Save original config
	origJSON := config.Get().JSON
	defer func() { config.Get().JSON = origJSON }()

	// Enable JSON mode
	config.Get().JSON = true

	headers := []string{"Name", "Age"}
	rows := [][]string{
		{"Alice", "30"},
		{"Bob", "25"},
	}

	// Should not panic
	OutputTable(headers, rows)
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello", 3, "..."},
		{"hi", 5, "hi"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		result := Truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		count    int
		singular string
		plural   string
		expected string
	}{
		{1, "item", "items", "item"},
		{0, "item", "items", "items"},
		{2, "item", "items", "items"},
		{100, "session", "sessions", "sessions"},
	}

	for _, tt := range tests {
		result := Pluralize(tt.count, tt.singular, tt.plural)
		if result != tt.expected {
			t.Errorf("Pluralize(%d, %q, %q) = %q, want %q", tt.count, tt.singular, tt.plural, result, tt.expected)
		}
	}
}
