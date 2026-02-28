package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/cmd/agent"
	"github.com/anomalyco/oho/cmd/auth"
	"github.com/anomalyco/oho/cmd/command"
	"github.com/anomalyco/oho/cmd/configcmd"
	"github.com/anomalyco/oho/cmd/file"
	"github.com/anomalyco/oho/cmd/find"
	"github.com/anomalyco/oho/cmd/formatter"
	"github.com/anomalyco/oho/cmd/global"
	"github.com/anomalyco/oho/cmd/lsp"
	"github.com/anomalyco/oho/cmd/mcp"
	"github.com/anomalyco/oho/cmd/mcpserver"
	"github.com/anomalyco/oho/cmd/message"
	"github.com/anomalyco/oho/cmd/project"
	"github.com/anomalyco/oho/cmd/provider"
	"github.com/anomalyco/oho/cmd/session"
	"github.com/anomalyco/oho/cmd/tool"
	"github.com/anomalyco/oho/cmd/tui"
	"github.com/anomalyco/oho/internal/config"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var versionStr string
var commitStr string
var dateStr string

// rootCmd 是 oho CLI 的根命令
var rootCmd = &cobra.Command{
	Use:   "oho",
	Short: "OpenCode CLI - HTTP 客户端工具",
	Long: `oho 是 OpenCode Server 的命令行客户端工具。
	
它提供了对 OpenCode Server API 的完整访问，允许你通过命令行管理会话、消息、配置等。

示例:
  oho session create              # 创建新会话
  oho message add -s session123   # 添加消息到会话
  oho session list                # 列出所有会话
  oho config get                  # 获取配置
  oho provider list               # 列出所有提供商`,
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo 设置版本信息
func SetVersionInfo(ver, commit, date string) {
	versionStr = ver
	commitStr = commit
	dateStr = date

	if versionStr != "dev" {
		rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", versionStr, commitStr, dateStr)
	} else {
		rootCmd.Version = versionStr
	}
}

func main() {
	SetVersionInfo(Version, Commit, Date)

	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "警告：配置初始化失败：%v\n", err)
	}

	// 全局标志
	rootCmd.PersistentFlags().StringP("host", "", "127.0.0.1", "服务器主机地址")
	rootCmd.PersistentFlags().IntP("port", "p", 4096, "服务器端口")
	rootCmd.PersistentFlags().StringP("password", "", "", "服务器密码 (覆盖环境变量)")
	rootCmd.PersistentFlags().BoolP("json", "j", false, "以 JSON 格式输出")

	// 绑定配置
	config.BindFlags(rootCmd.PersistentFlags())

	// 添加子命令
	rootCmd.AddCommand(
		global.Cmd,
		project.Cmd,
		session.Cmd,
		message.Cmd,
		configcmd.Cmd,
		provider.Cmd,
		file.Cmd,
		find.Cmd,
		tool.Cmd,
		agent.Cmd,
		command.Cmd,
		lsp.Cmd,
		formatter.Cmd,
		mcp.Cmd,
		mcpserver.Cmd,
		tui.Cmd,
		auth.Cmd,
	)

	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
