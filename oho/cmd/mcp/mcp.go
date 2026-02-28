package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd MCP 命令
var Cmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP 服务器管理",
	Long:  "管理 MCP 服务器",
}

var (
	mcpConfig string

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "列出 MCP 服务器状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.NewClient()
			ctx := context.Background()

			resp, err := c.Get(ctx, "/mcp")
			if err != nil {
				return err
			}

			var status map[string]types.MCPStatus
			if err := json.Unmarshal(resp, &status); err != nil {
				return err
			}

			if config.Get().JSON {
				data, _ := json.MarshalIndent(status, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			if len(status) == 0 {
				fmt.Println("没有 MCP 服务器")
				return nil
			}

			fmt.Println("MCP 服务器状态:")
			for name, s := range status {
				icon := "❌"
				if s.Status == "running" {
					icon = "✅"
				}
				fmt.Printf("%s %s (状态：%s)\n", icon, name, s.Status)
				if s.Error != "" {
					fmt.Printf("   错误：%s\n", s.Error)
				}
			}

			return nil
		},
	}

	addCmd = &cobra.Command{
		Use:   "add <name>",
		Short: "添加 MCP 服务器",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if mcpConfig == "" {
				return fmt.Errorf("请提供 --config 参数 (JSON 格式)")
			}

			c := client.NewClient()
			ctx := context.Background()

			// 解析配置
			var configData map[string]interface{}
			if err := json.Unmarshal([]byte(mcpConfig), &configData); err != nil {
				return fmt.Errorf("解析配置失败：%w", err)
			}

			req := types.MCPConfig{
				Name:   args[0],
				Config: configData,
			}

			resp, err := c.Post(ctx, "/mcp", req)
			if err != nil {
				return err
			}

			var status map[string]types.MCPStatus
			if err := json.Unmarshal(resp, &status); err != nil {
				return err
			}

			fmt.Printf("MCP 服务器 %s 已添加\n", args[0])
			return nil
		},
	}
)

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&mcpConfig, "config", "", "MCP 服务器配置 (JSON 格式)")
	_ = addCmd.MarkFlagRequired("config")
}
