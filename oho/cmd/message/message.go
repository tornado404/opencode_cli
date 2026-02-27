package message

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
	"github.com/anomalyco/oho/internal/types"
)

// Cmd 消息命令
var Cmd = &cobra.Command{
	Use:   "message",
	Short: "消息管理命令",
	Long:  "管理 OpenCode 会话消息，包括发送、列表、命令执行等",
}

var (
	sessionID      string
	messageID      string
	model          string
	agent          string
	noReply        bool
	systemPrompt   string
	tools          []string
	asyncMode      bool
	commandName    string
	commandArgs    []string
	shellCommand   string
)

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(promptAsyncCmd)
	Cmd.AddCommand(commandCmd)
	Cmd.AddCommand(shellCmd)

	// 全局标志
	Cmd.PersistentFlags().StringVarP(&sessionID, "session", "s", "", "会话 ID")
	Cmd.MarkPersistentFlagRequired("session")
}

// listCmd 列出消息
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出会话中的消息",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		queryParams := map[string]string{}
		if limit := cmd.Flag("limit"); limit != nil && limit.Value.String() != "" {
			queryParams["limit"] = limit.Value.String()
		}

		resp, err := c.GetWithQuery(ctx, fmt.Sprintf("/session/%s/message", sessionID), queryParams)
		if err != nil {
			return err
		}

		var messages []types.MessageWithParts
		if err := json.Unmarshal(resp, &messages); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(messages, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		for _, msg := range messages {
			fmt.Printf("\n[%s] %s\n", msg.Info.Role, msg.Info.ID)
			if msg.Info.Content != "" {
				fmt.Printf("%s\n", msg.Info.Content)
			}
			for _, part := range msg.Parts {
				fmt.Printf("  └─ 部分类型：%s\n", part.Type)
			}
			fmt.Println("---")
		}

		return nil
	},
}

// addCmd 添加消息
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "发送消息并等待响应",
	Long:  "发送消息到会话并等待 AI 响应",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// 从 stdin 读取
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				data, err := os.ReadFile("/dev/stdin")
				if err != nil {
					return fmt.Errorf("读取 stdin 失败：%w", err)
				}
				args = []string{string(data)}
			} else {
				return fmt.Errorf("请提供消息内容，例如：oho message add -s <session> \"你好\"")
			}
		}

		c := client.NewClient()
		ctx := context.Background()

		parts := []types.Part{
			{
				Type: "text",
				Text: args[0],
			},
		}

		req := types.MessageRequest{
			MessageID: messageID,
			Model:     model,
			Agent:     agent,
			NoReply:   noReply,
			System:    systemPrompt,
			Tools:     tools,
			Parts:     parts,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/message", sessionID), req)
		if err != nil {
			return err
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(resp, &result); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("消息已发送:\n")
		fmt.Printf("  ID: %s\n", result.Info.ID)
		fmt.Printf("  角色：%s\n", result.Info.Role)
		
		for _, part := range result.Parts {
			fmt.Printf("\n[%s]\n", part.Type)
			if text, ok := part.Text.(string); ok {
				fmt.Printf("%s\n", text)
			} else {
				data, _ := json.MarshalIndent(part.Text, "", "  ")
				fmt.Printf("%s\n", string(data))
			}
		}

		return nil
	},
}

// getCmd 获取消息详情
var getCmd = &cobra.Command{
	Use:   "get <messageID>",
	Short: "获取消息详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		resp, err := c.Get(ctx, fmt.Sprintf("/session/%s/message/%s", sessionID, args[0]))
		if err != nil {
			return err
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(resp, &result); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("消息详情:\n")
		fmt.Printf("  ID: %s\n", result.Info.ID)
		fmt.Printf("  会话：%s\n", result.Info.SessionID)
		fmt.Printf("  角色：%s\n", result.Info.Role)
		fmt.Printf("  时间：%d\n", result.Info.CreatedAt)
		
		if result.Info.Content != "" {
			fmt.Printf("\n内容:\n%s\n", result.Info.Content)
		}

		fmt.Printf("\n部分 (%d 个):\n", len(result.Parts))
		for i, part := range result.Parts {
			fmt.Printf("  %d. 类型：%s\n", i+1, part.Type)
		}

		return nil
	},
}

// promptAsyncCmd 异步发送消息
var promptAsyncCmd = &cobra.Command{
	Use:   "prompt-async",
	Short: "异步发送消息（不等待响应）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供消息内容")
		}

		c := client.NewClient()
		ctx := context.Background()

		parts := []types.Part{
			{
				Type: "text",
				Text: args[0],
			},
		}

		req := types.MessageRequest{
			MessageID: messageID,
			Model:     model,
			Agent:     agent,
			NoReply:   false,
			System:    systemPrompt,
			Tools:     tools,
			Parts:     parts,
		}

		_, err := c.Post(ctx, fmt.Sprintf("/session/%s/prompt_async", sessionID), req)
		if err != nil {
			return err
		}

		fmt.Println("消息已异步发送")
		return nil
	},
}

// commandCmd 执行斜杠命令
var commandCmd = &cobra.Command{
	Use:   "command <command>",
	Short: "执行斜杠命令",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient()
		ctx := context.Background()

		// 解析命令参数
		argMap := make(map[string]string)
		for _, arg := range commandArgs {
			if idx := indexOf(arg, "="); idx > 0 {
				key := arg[:idx]
				value := arg[idx+1:]
				argMap[key] = value
			}
		}

		req := types.CommandRequest{
			MessageID: messageID,
			Agent:     agent,
			Model:     model,
			Command:   args[0],
			Arguments: argMap,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/command", sessionID), req)
		if err != nil {
			return err
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(resp, &result); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("命令已执行:\n")
		fmt.Printf("  消息 ID: %s\n", result.Info.ID)
		fmt.Printf("  角色：%s\n", result.Info.Role)
		
		for _, part := range result.Parts {
			fmt.Printf("\n[%s]\n", part.Type)
			if text, ok := part.Text.(string); ok {
				fmt.Printf("%s\n", text)
			}
		}

		return nil
	},
}

// shellCmd 运行 shell 命令
var shellCmd = &cobra.Command{
	Use:   "shell <command>",
	Short: "运行 shell 命令",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if agent == "" {
			return fmt.Errorf("请提供 --agent 参数")
		}

		c := client.NewClient()
		ctx := context.Background()

		cmdStr := shellCommand
		if cmdStr == "" && len(args) > 0 {
			cmdStr = args[0]
		}

		req := types.ShellRequest{
			Agent:   agent,
			Model:   model,
			Command: cmdStr,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/shell", sessionID), req)
		if err != nil {
			return err
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(resp, &result); err != nil {
			return err
		}

		if config.Get().JSON {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Shell 命令已执行:\n")
		fmt.Printf("  消息 ID: %s\n", result.Info.ID)
		
		for _, part := range result.Parts {
			fmt.Printf("\n[%s]\n", part.Type)
			if text, ok := part.Text.(string); ok {
				fmt.Printf("%s\n", text)
			}
		}

		return nil
	},
}

func init() {
	// list 命令标志
	listCmd.Flags().IntP("limit", "l", 0, "限制消息数量")

	// add 命令标志
	addCmd.Flags().StringVar(&messageID, "message", "", "消息 ID")
	addCmd.Flags().StringVar(&model, "model", "", "模型 ID")
	addCmd.Flags().StringVar(&agent, "agent", "", "代理 ID")
	addCmd.Flags().BoolVar(&noReply, "no-reply", false, "不等待响应")
	addCmd.Flags().StringVar(&systemPrompt, "system", "", "系统提示")
	addCmd.Flags().StringSliceVar(&tools, "tools", nil, "工具列表")

	// prompt-async 命令标志
	promptAsyncCmd.Flags().StringVar(&messageID, "message", "", "消息 ID")
	promptAsyncCmd.Flags().StringVar(&model, "model", "", "模型 ID")
	promptAsyncCmd.Flags().StringVar(&agent, "agent", "", "代理 ID")
	promptAsyncCmd.Flags().StringVar(&systemPrompt, "system", "", "系统提示")
	promptAsyncCmd.Flags().StringSliceVar(&tools, "tools", nil, "工具列表")

	// command 命令标志
	commandCmd.Flags().StringVar(&messageID, "message", "", "消息 ID")
	commandCmd.Flags().StringVar(&agent, "agent", "", "代理 ID")
	commandCmd.Flags().StringVar(&model, "model", "", "模型 ID")
	commandCmd.Flags().StringArrayVar(&commandArgs, "args", nil, "命令参数 (key=value)")

	// shell 命令标志
	shellCmd.Flags().StringVar(&agent, "agent", "", "代理 ID (必需)")
	shellCmd.Flags().StringVar(&model, "model", "", "模型 ID")
	shellCmd.Flags().StringVar(&shellCommand, "command", "", "Shell 命令")
}

func indexOf(s string, substr string) int {
	for i := 0; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
