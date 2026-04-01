package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitWithMissingConfig(t *testing.T) {
	// 清除所有可能影响测试的环境变量
	os.Unsetenv("OPENCODE_SERVER_HOST")
	os.Unsetenv("OPENCODE_SERVER_PORT")
	os.Unsetenv("OPENCODE_SERVER_USERNAME")
	os.Unsetenv("OPENCODE_SERVER_PASSWORD")

	// 确保没有配置文件在默认路径
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "oho", "config.json")
	os.RemoveAll(filepath.Dir(configPath))

	// 初始化配置
	err := Init()
	if err != nil {
		t.Fatalf("Init() returned error: %v", err)
	}

	// 验证默认值被使用
	cfg := Get()
	if cfg.Host != "127.0.0.1" {
		t.Errorf("Expected Host '127.0.0.1', got '%s'", cfg.Host)
	}
	if cfg.Port != 4096 {
		t.Errorf("Expected Port 4096, got %d", cfg.Port)
	}
	if cfg.Username != "opencode" {
		t.Errorf("Expected Username 'opencode', got '%s'", cfg.Username)
	}
	if cfg.Password != "" {
		t.Errorf("Expected empty Password, got '%s'", cfg.Password)
	}
}

func TestInitWithEnvOverrides(t *testing.T) {
	// 设置环境变量
	os.Setenv("OPENCODE_SERVER_HOST", "192.168.1.1")
	os.Setenv("OPENCODE_SERVER_PORT", "8080")
	os.Setenv("OPENCODE_SERVER_USERNAME", "testuser")
	os.Setenv("OPENCODE_SERVER_PASSWORD", "testpass")
	defer func() {
		os.Unsetenv("OPENCODE_SERVER_HOST")
		os.Unsetenv("OPENCODE_SERVER_PORT")
		os.Unsetenv("OPENCODE_SERVER_USERNAME")
		os.Unsetenv("OPENCODE_SERVER_PASSWORD")
	}()

	err := Init()
	if err != nil {
		t.Fatalf("Init() returned error: %v", err)
	}

	cfg := Get()
	if cfg.Host != "192.168.1.1" {
		t.Errorf("Expected Host '192.168.1.1', got '%s'", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("Expected Port 8080, got %d", cfg.Port)
	}
	if cfg.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", cfg.Username)
	}
	if cfg.Password != "testpass" {
		t.Errorf("Expected Password 'testpass', got '%s'", cfg.Password)
	}
}

func TestGetBaseURL(t *testing.T) {
	os.Unsetenv("OPENCODE_SERVER_HOST")
	os.Unsetenv("OPENCODE_SERVER_PORT")
	os.Unsetenv("OPENCODE_SERVER_USERNAME")
	os.Unsetenv("OPENCODE_SERVER_PASSWORD")

	Init()

	baseURL := GetBaseURL()
	expected := "http://127.0.0.1:4096"
	if baseURL != expected {
		t.Errorf("Expected BaseURL '%s', got '%s'", expected, baseURL)
	}
}

func TestGetBaseURLWithCustomHostPort(t *testing.T) {
	os.Setenv("OPENCODE_SERVER_HOST", "10.0.0.1")
	os.Setenv("OPENCODE_SERVER_PORT", "3000")
	defer func() {
		os.Unsetenv("OPENCODE_SERVER_HOST")
		os.Unsetenv("OPENCODE_SERVER_PORT")
	}()

	Init()

	baseURL := GetBaseURL()
	expected := "http://10.0.0.1:3000"
	if baseURL != expected {
		t.Errorf("Expected BaseURL '%s', got '%s'", expected, baseURL)
	}
}

func TestGetConfigPathFallback(t *testing.T) {
	// 测试 getConfigPath 返回有效路径
	path := getConfigPath()
	if path == "" {
		t.Error("getConfigPath() returned empty string")
	}

	// 验证路径包含 oho 和 config.json
	if !filepath.IsAbs(path) && path != filepath.Join(".", ".config", "oho", "config.json") {
		// 如果不是绝对路径，应该是相对路径 fallback
		expectedRelative := filepath.Join(".", ".config", "oho", "config.json")
		if path != expectedRelative {
			t.Errorf("Unexpected relative path: got %s, expected %s", path, expectedRelative)
		}
	}

	// 如果是绝对路径，应该包含正确的目录结构
	if filepath.IsAbs(path) {
		if !contains(path, ".config") || !contains(path, "oho") || !contains(path, "config.json") {
			t.Errorf("Path does not contain expected components: %s", path)
		}
	}
}

// contains 是一个简单的字符串包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
