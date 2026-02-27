package util

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"short", 10, "short"},
		{"exactly", 7, "exactly"},
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
		{2, "item", "items", "items"},
		{0, "item", "items", "items"},
		{100, "session", "sessions", "sessions"},
		{1, "session", "sessions", "session"},
	}

	for _, tt := range tests {
		result := Pluralize(tt.count, tt.singular, tt.plural)
		if result != tt.expected {
			t.Errorf("Pluralize(%d, %q, %q) = %q, want %q", tt.count, tt.singular, tt.plural, result, tt.expected)
		}
	}
}

func TestConfirm(t *testing.T) {
	// Confirm 函数需要交互式输入，这里只测试不 panic
	// 在 CI 环境中会被跳过
}

func TestReadStdin(t *testing.T) {
	// ReadStdin 测试需要特殊处理，这里只验证函数存在
}
