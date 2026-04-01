package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
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

	// 2. 加载配置文件（尝试多个可能的位置）
	configFile := findConfigFile()
	if configFile != "" {
		if data, err := os.ReadFile(configFile); err == nil {
			fmt.Fprintf(os.Stderr, "[config] 成功读取配置文件: %s\n", configFile)
			if err := json.Unmarshal(data, cfg); err != nil {
				return fmt.Errorf("解析配置文件失败：%w", err)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "[config] 配置文件不存在，请创建或设置环境变量\n")
		fmt.Fprintf(os.Stderr, "[config] 尝试过的路径:\n")
		for _, p := range getConfigSearchPaths() {
			fmt.Fprintf(os.Stderr, "[config]   - %s\n", p)
		}
	}

	// 3. 环境变量覆盖配置文件（始终检查，作为中间优先级）
	// 优先级：命令行标志 > 环境变量 > 配置文件 > 默认值
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

// getConfigSearchPaths 返回所有可能配置文件的搜索路径
func getConfigSearchPaths() []string {
	var paths []string

	// 方式1: 操作系统用户数据库中的真实主目录（不受环境变量影响）
	if usr, err := user.Current(); err == nil && usr.HomeDir != "" {
		paths = append(paths, filepath.Join(usr.HomeDir, ".config", "oho", "config.json"))
	}

	// 方式2: os.UserHomeDir()（某些情况下会读取 HOME 环境变量）
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		paths = append(paths, filepath.Join(home, ".config", "oho", "config.json"))
	}

	// 方式3: $HOME 环境变量
	if home := os.Getenv("HOME"); home != "" {
		paths = append(paths, filepath.Join(home, ".config", "oho", "config.json"))
	}

	// 方式4: $USERPROFILE (Windows)
	if home := os.Getenv("USERPROFILE"); home != "" {
		paths = append(paths, filepath.Join(home, ".config", "oho", "config.json"))
	}

	// 方式5: 相对于当前工作目录
	paths = append(paths, filepath.Join(".", ".config", "oho", "config.json"))

	// 去重
	seen := make(map[string]bool)
	var unique []string
	for _, p := range paths {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}
	return unique
}

// findConfigFile 查找配置文件，返回第一个存在的路径
func findConfigFile() string {
	for _, p := range getConfigSearchPaths() {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// getConfigPath 获取配置文件路径（用于保存）
func getConfigPath() string {
	// 优先使用 os.UserHomeDir() 对应的路径
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", "oho", "config.json")
	}
	// fallback 到 $HOME
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config", "oho", "config.json")
	}
	// fallback 到 $USERPROFILE
	if home := os.Getenv("USERPROFILE"); home != "" {
		return filepath.Join(home, ".config", "oho", "config.json")
	}
	return filepath.Join(".", ".config", "oho", "config.json")
}
