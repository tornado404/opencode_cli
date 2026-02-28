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


