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
// 优先级：命令行标志 > 配置文件 > 环境变量 > 默认值
func Init() error {
	// 1. 初始化默认值（最低优先级）
	cfg = &Config{
		Host:     "127.0.0.1",
		Port:     4096,
		Username: "opencode",
		Password: "",
		JSON:     false,
	}

	// 2. 检查配置文件是否存在
	configFile := getConfigPath()
	if data, err := os.ReadFile(configFile); err == nil {
		// 配置文件存在，加载配置
		if err := json.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("解析配置文件失败：%w", err)
		}
	} else {
		// 3. 配置文件不存在时，使用环境变量
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
	}

	return nil
}

// BindFlags 绑定命令行标志到配置（最高优先级）
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
