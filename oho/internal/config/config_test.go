package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 设置测试环境变量
	os.Setenv("OPENCODE_SERVER_HOST", "127.0.0.1")
	os.Setenv("OPENCODE_SERVER_PORT", "4096")
	os.Setenv("OPENCODE_SERVER_USERNAME", "opencode")
	os.Setenv("OPENCODE_SERVER_PASSWORD", "test")

	// 初始化配置
	Init()

	m.Run()
}
