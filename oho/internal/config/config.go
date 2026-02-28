package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

// Config 存储 CLI 配置
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	JSON     bool   `json:"json"`
}

var cfg *Config

// Init 初始化配置
func Init() error {
	cfg = &Config{
		Host:     "127.0.0.1",
		Port:     4096,
		Username: "opencode",
		Password: os.Getenv("OPENCODE_SERVER_PASSWORD"),
		JSON:     false,
	}

	// 从配置文件加载
	configFile := getConfigPath()
	if data, err := os.ReadFile(configFile); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("解析配置文件失败：%w", err)
		}
	}

	// 环境变量覆盖（优先级最高）
	if envHost := os.Getenv("OPENCODE_SERVER_HOST"); envHost != "" {
		cfg.Host = envHost
	}
	if envPort := os.Getenv("OPENCODE_SERVER_PORT"); envPort != "" {
		_, _ = fmt.Sscanf(envPort, "%d", &cfg.Port)
	}
	if envUsername := os.Getenv("OPENCODE_SERVER_USERNAME"); envUsername != "" {
		cfg.Username = envUsername
	}
	if envPassword := os.Getenv("OPENCODE_SERVER_PASSWORD"); envPassword != "" {
		cfg.Password = envPassword
	}

	return nil
}

// BindFlags 绑定命令行标志到配置
func BindFlags(flags *pflag.FlagSet) {
	if host, _ := flags.GetString("host"); host != "" {
		cfg.Host = host
	}
	if port, _ := flags.GetInt("port"); port != 4096 {
		cfg.Port = port
	}
	if password, _ := flags.GetString("password"); password != "" {
		cfg.Password = password
	}
	if jsonOut, _ := flags.GetBool("json"); jsonOut {
		cfg.JSON = jsonOut
	}
}

// Get 获取配置
func Get() *Config {
	return cfg
}

// GetBaseURL 获取服务器基础 URL
func GetBaseURL() string {
	return fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
}

// Save 保存配置到文件
func Save() error {
	configFile := getConfigPath()
	dir := filepath.Dir(configFile)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0600)
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "oho", "config.json")
}
