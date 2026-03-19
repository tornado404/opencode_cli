package session

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/types"
)

// submitCmd 提交任务命令
var submitCmd = &cobra.Command{
	Use:   "submit [message]",
	Short: "Submit a task by creating a session and sending a message in one step",
	Long:  "Create a new session in current directory, optionally initialize it with AGENTS.md, and send a message in one command.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Step 1: Validate flags
		if initProject {
			if providerID == "" || modelID == "" {
				return fmt.Errorf("when using --init-project, --provider and --model are required")
			}
		}

		c := client.NewClient()
		ctx := context.Background()

		// Step 2: Create session
		// 根据 OpenCode SDK: directory 是 query 参数，不是 body 参数
		// body 只支持 parentID 和 title
		req := map[string]interface{}{}
		if title != "" {
			req["title"] = title
		}

		// 获取当前工作目录（如果用户未指定）
		sessionDir := directory
		if sessionDir == "" {
			var err error
			sessionDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}

		// 使用 PostWithQuery 发送 directory 作为 query 参数
		queryParams := map[string]string{"directory": sessionDir}
		resp, err := c.PostWithQuery(ctx, "/session", queryParams, req)
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		var session types.Session
		if err := json.Unmarshal(resp, &session); err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		fmt.Printf("Session created: %s\n", session.ID)

		// Step 3: Initialize session (if requested)
		if initProject {
			initReq := map[string]interface{}{
				"providerID": providerID,
				"modelID":    modelID,
			}

			_, err := c.Post(ctx, fmt.Sprintf("/session/%s/init", session.ID), initReq)
			if err != nil {
				return fmt.Errorf("failed to initialize session: %w", err)
			}

			fmt.Println("Session initialized successfully")
		}

		// Step 4: Prepare message parts
		var parts []types.Part

		// Add text part from args[0]
		text := args[0]
		parts = append(parts, types.Part{
			Type: "text",
			Text: &text,
		})

		// Add file parts for each file
		for _, filePath := range files {
			// Check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", filePath)
			}

			// Read file content
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file: %s: %w", filePath, err)
			}

			// Detect MIME type
			mimeType := detectMimeType(filePath)

			// Encode to base64 data URL
			base64Data := base64.StdEncoding.EncodeToString(fileData)
			dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

			parts = append(parts, types.Part{
				Type: "file",
				URL:  dataURL,
				Mime: mimeType,
			})
		}

		// Step 5: Send message
		msgReq := types.MessageRequest{
			MessageID: messageID,
			Model:     messageModel,
			Agent:     messageAgent,
			NoReply:   noReply,
			System:    systemPrompt,
			Tools:     tools,
			Parts:     parts,
		}

		msgResp, err := c.Post(ctx, fmt.Sprintf("/session/%s/message", session.ID), msgReq)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		// Handle empty response
		if len(msgResp) == 0 {
			fmt.Println("Message sent successfully")
			return nil
		}

		var result types.MessageWithParts
		if err := json.Unmarshal(msgResp, &result); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		fmt.Printf("Message sent successfully: %s\n", result.Info.ID)

		// Step 6: Return nil on success
		return nil
	},
}

// Flag variables for submit command (additional ones not in session.go)
var (
	initProject  bool
	directory    string
	messageAgent string
	messageModel string
	noReply      bool
	systemPrompt string
	tools        []string
	files        []string
)

func init() {
	// Session creation flags
	submitCmd.Flags().BoolVar(&initProject, "init-project", false, "Initialize project with AGENTS.md")
	submitCmd.Flags().StringVar(&providerID, "provider", "", "Provider ID for initialization")
	submitCmd.Flags().StringVar(&modelID, "model", "", "Model ID for initialization")
	submitCmd.Flags().StringVar(&title, "title", "", "Session title")
	submitCmd.Flags().StringVar(&directory, "directory", "", "Working directory for the session")

	// Message flags
	submitCmd.Flags().StringVar(&messageAgent, "agent", "", "Agent ID for message")
	submitCmd.Flags().StringVar(&messageModel, "message-model", "", "Model ID for message")
	submitCmd.Flags().BoolVar(&noReply, "no-reply", false, "Don't wait for response")
	submitCmd.Flags().StringVar(&systemPrompt, "system", "", "System prompt")
	submitCmd.Flags().StringSliceVar(&tools, "tools", nil, "Tools list (can be specified multiple times)")
	submitCmd.Flags().StringSliceVar(&files, "file", nil, "File attachments (can be specified multiple times)")
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
