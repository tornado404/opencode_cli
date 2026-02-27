package config

import (
	"os"
	"testing"
)

func TestGetBaseURL(t *testing.T) {
	// 保存原始值
	origHost := os.Getenv("OPENCODE_SERVER_HOST")
	origPort := os.Getenv("OPENCODE_SERVER_PORT")
	defer func() {
		if origHost != "" {
			os.Setenv("OPENCODE_SERVER_HOST", origHost)
		}
		if origPort != "" {
			os.Setenv("OPENCODE_SERVER_PORT", origPort)
		}
	}()

	tests := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{
			name:     "default values",
			host:     "127.0.0.1",
			port:     4096,
			expected: "http://127.0.0.1:4096",
		},
		{
			name:     "custom host",
			host:     "192.168.1.1",
			port:     4096,
			expected: "http://192.168.1.1:4096",
		},
		{
			name:     "custom port",
			host:     "127.0.0.1",
			port:     8080,
			expected: "http://127.0.0.1:8080",
		},
		{
			name:     "custom host and port",
			host:     "localhost",
			port:     3000,
			expected: "http://localhost:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.host != "127.0.0.1" {
				os.Setenv("OPENCODE_SERVER_HOST", tt.host)
			}
			if tt.port != 4096 {
				os.Setenv("OPENCODE_SERVER_PORT", "")
			}

			// 重新初始化配置
			cfg = &Config{
				Host:     tt.host,
				Port:     tt.port,
				Username: "opencode",
				Password: "",
			}

			result := GetBaseURL()
			if result != tt.expected {
				t.Errorf("GetBaseURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// 测试默认配置
	cfg = &Config{
		Host:     "127.0.0.1",
		Port:     4096,
		Username: "opencode",
		Password: "",
		JSON:     false,
	}

	result := Get()
	if result == nil {
		t.Error("Get() returned nil")
	}
	if result.Host != "127.0.0.1" {
		t.Errorf("Get().Host = %v, want 127.0.0.1", result.Host)
	}
	if result.Port != 4096 {
		t.Errorf("Get().Port = %v, want 4096", result.Port)
	}
}
