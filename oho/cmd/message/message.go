package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	sessionID    string
	messageID    string
	model        string
	agent        string
	noReply      bool
	systemPrompt string
	tools        []string
	commandArgs  []string
	shellCommand string
	files        []string
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
	_ = Cmd.MarkPersistentFlagRequired("session")

	// list 命令标志
	listCmd.Flags().IntP("limit", "l", 0, "限制消息数量")

	// add 命令标志
	addCmd.Flags().StringVar(&messageID, "message", "", "消息 ID")
	addCmd.Flags().StringVar(&model, "model", "", "模型 ID")
	addCmd.Flags().StringVar(&agent, "agent", "", "代理 ID")
	addCmd.Flags().BoolVar(&noReply, "no-reply", false, "不等待响应")
	addCmd.Flags().StringVar(&systemPrompt, "system", "", "系统提示")
	addCmd.Flags().StringSliceVar(&tools, "tools", nil, "工具列表")
	addCmd.Flags().StringSliceVar(&files, "file", nil, "附件文件路径 (可多次使用)")

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

// convertModel converts a model string to the appropriate format (string or object)
func convertModel(model string) interface{} {
	if model == "" {
		return nil
	}

	// Check if the model string contains a colon, which indicates provider:model format
	if strings.Contains(model, ":") {
		parts := strings.SplitN(model, ":", 2)
		if len(parts) == 2 {
			return types.Model{
				ProviderID: parts[0],
				ModelID:    parts[1],
			}
		}
	}

	// If no colon, treat as simple string model (backward compatibility)
	return model
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
			return fmt.Errorf("解析消息列表失败：%w", err)
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
		if len(args) == 0 && len(files) == 0 {
			// 从 stdin 读取
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				data, err := os.ReadFile("/dev/stdin")
				if err != nil {
					return fmt.Errorf("读取 stdin 失败：%w", err)
				}
				args = []string{string(data)}
			} else {
				return fmt.Errorf("请提供消息内容或文件，例如：oho message add -s <session> \"你好\" 或 oho message add -s <session> --file image.jpg")
			}
		}

		c := client.NewClient()
		ctx := context.Background()

		// 构建 parts 数组
		var parts []types.Part

		// 添加文本部分
		if len(args) > 0 && args[0] != "" {
			text := args[0]
			parts = append(parts, types.Part{
				Type: "text",
				Text: &text,
			})
		}

		// 添加文件部分
		for _, filePath := range files {
			// 检查文件是否存在
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("文件不存在：%s", filePath)
			}

			// 读取文件内容
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("读取文件失败：%s: %w", filePath, err)
			}

			// 检测 MIME 类型
			mimeType := detectMimeType(filePath)

			// 将文件编码为 base64 data URL
			base64Data := base64.StdEncoding.EncodeToString(fileData)
			dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

			// 使用 base64 data URL
			parts = append(parts, types.Part{
				Type: "file",
				URL:  dataURL,
				Mime: mimeType,
			})
		}

		req := types.MessageRequest{
			MessageID: messageID,
			Model:     convertModel(model),
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

		// 服务器返回空响应时处理
		if len(resp) == 0 {
			fmt.Println("消息已发送")
			return nil
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(resp, &result); err != nil {
			return fmt.Errorf("解析响应失败：%w", err)
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
			if part.Text != nil {
				fmt.Printf("%s\n", *part.Text)
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
			return fmt.Errorf("解析响应失败：%w", err)
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
				Text: &args[0],
			},
		}

		req := types.MessageRequest{
			MessageID: messageID,
			Model:     convertModel(model),
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
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				argMap[parts[0]] = parts[1]
			}
		}

		req := types.CommandRequest{
			MessageID: messageID,
			Agent:     agent,
			Model:     convertModel(model),
			Command:   args[0],
			Arguments: argMap,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/command", sessionID), req)
		if err != nil {
			return err
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(resp, &result); err != nil {
			return fmt.Errorf("解析响应失败：%w", err)
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
			if part.Text != nil {
				fmt.Printf("%s\n", *part.Text)
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

		// 优先使用 --command 标志，否则使用第一个位置参数
		cmdStr := shellCommand
		if cmdStr == "" && len(args) > 0 {
			cmdStr = args[0]
		}

		if cmdStr == "" {
			return fmt.Errorf("请提供要执行的 shell 命令")
		}

		req := types.ShellRequest{
			Agent:   agent,
			Model:   convertModel(model),
			Command: cmdStr,
		}

		resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/shell", sessionID), req)
		if err != nil {
			return err
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(resp, &result); err != nil {
			return fmt.Errorf("解析响应失败：%w", err)
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
			if part.Text != nil {
				fmt.Printf("%s\n", *part.Text)
			}
		}

		return nil
	},
}

// detectMimeType 根据文件扩展名检测 MIME 类型
func detectMimeType(filePath string) string {
	ext := strings.ToLower(filePath[strings.LastIndex(filePath, "."):])

	mimeTypes := map[string]string{
		// 图片
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".bmp":  "image/bmp",
		".svg":  "image/svg+xml",

		// 文档
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",

		// 文本
		".txt":  "text/plain",
		".md":   "text/markdown",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".yaml": "application/x-yaml",
		".yml":  "application/x-yaml",

		// 代码
		".py":   "text/x-python",
		".go":   "text/x-go",
		".java": "text/x-java",
		".c":    "text/x-c",
		".cpp":  "text/x-c++",
		".h":    "text/x-c",
		".rs":   "text/x-rust",
		".ts":   "text/x-typescript",
		".tsx":  "text/x-typescript",

		// 其他
		".zip": "application/zip",
		".tar": "application/x-tar",
		".gz":  "application/gzip",
		".mp3": "audio/mpeg",
		".mp4": "video/mp4",
		".wav": "audio/wav",
	}

	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}

	// 默认返回 octet-stream
	return "application/octet-stream"
}
