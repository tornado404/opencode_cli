package add

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/types"
)

// Flag variables for add command
var (
	addTitle      string
	addParent     string
	addAgent      string
	addModel      string
	addNoReply    bool
	addSystem     string
	addTools      []string
	addFiles      []string
	addDirectory  string
	addJSONOutput bool
)

// Cmd add 命令 - 创建会话并发送消息
var Cmd = &cobra.Command{
	Use:   "add [message]",
	Short: "Create a new session and send a message in one command",
	Long: `Create a new session in the current directory and send a message to it in one step.

This command combines session creation and message sending into a single operation.
By default, it uses the current working directory as the session directory and
generates an automatic title.

Examples:
  oho add "帮我分析这个项目"
  oho add "修复登录 bug" --title "Bug 修复"
  oho add "测试功能" --no-reply --agent default
  oho add "分析日志" --file /var/log/app.log`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAdd,
}

func init() {
	// Session-related flags
	Cmd.Flags().StringVar(&addTitle, "title", "", "Session title (auto-generated if not provided)")
	Cmd.Flags().StringVar(&addParent, "parent", "", "Parent session ID (for creating sub-session)")
	Cmd.Flags().StringVar(&addDirectory, "directory", "", "Working directory for the session (default: current directory)")

	// Message-related flags
	Cmd.Flags().StringVar(&addAgent, "agent", "", "Agent ID for message")
	Cmd.Flags().StringVar(&addModel, "model", "", "Model ID for message")
	Cmd.Flags().BoolVar(&addNoReply, "no-reply", false, "Don't wait for AI response")
	Cmd.Flags().StringVar(&addSystem, "system", "", "System prompt")
	Cmd.Flags().StringSliceVar(&addTools, "tools", nil, "Tools list (can be specified multiple times)")
	Cmd.Flags().StringSliceVar(&addFiles, "file", nil, "File attachments (can be specified multiple times)")

	// Output format
	Cmd.Flags().BoolVarP(&addJSONOutput, "json", "j", false, "Output in JSON format")
}

func runAdd(cmd *cobra.Command, args []string) error {
	c := client.NewClient()
	ctx := context.Background()

	// Step 1: Get current working directory
	sessionDir := addDirectory
	if sessionDir == "" {
		var err error
		sessionDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Step 2: Generate title if not provided
	sessionTitle := addTitle
	if sessionTitle == "" {
		sessionTitle = fmt.Sprintf("New session - %s", time.Now().Format("2006-01-02T15:04:05"))
	}

	// Step 3: Create session
	sessionID, err := createSession(c, ctx, sessionTitle, addParent, sessionDir)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Step 4: Send message
	message := args[0]
	messageID, err := sendMessage(c, ctx, sessionID, message, addAgent, addModel, addNoReply, addSystem, addTools, addFiles)
	if err != nil {
		// Message send failed, but session was created
		if addJSONOutput {
			output := map[string]interface{}{
				"sessionId": sessionID,
				"status":    "partial",
				"error":     fmt.Sprintf("failed to send message: %v", err),
			}
			data, _ := json.MarshalIndent(output, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf("Session created: %s\n", sessionID)
			fmt.Printf("Warning: Message send failed: %v\n", err)
		}
		return nil
	}

	// Step 5: Output result
	if addJSONOutput {
		output := map[string]interface{}{
			"sessionId": sessionID,
			"messageId": messageID,
			"directory": sessionDir,
			"title":     sessionTitle,
			"status":    "success",
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Session created: %s\n", sessionID)
		fmt.Printf("Message sent: %s\n", messageID)
	}

	return nil
}

// createSession creates a new session and returns the session ID
func createSession(c *client.Client, ctx context.Context, title, parentID, directory string) (string, error) {
	req := map[string]interface{}{}
	if title != "" {
		req["title"] = title
	}
	if parentID != "" {
		req["parentID"] = parentID
	}

	// Use directory as query parameter (per OpenCode SDK spec)
	queryParams := map[string]string{"directory": directory}
	resp, err := c.PostWithQuery(ctx, "/session", queryParams, req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	var session types.Session
	if err := json.Unmarshal(resp, &session); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return session.ID, nil
}

// sendMessage sends a message to the session and returns the message ID
func sendMessage(c *client.Client, ctx context.Context, sessionID, message, agent, model string, noReply bool, system string, tools, files []string) (string, error) {
	// Build message parts
	var parts []types.Part

	// Add text part
	text := message
	parts = append(parts, types.Part{
		Type: "text",
		Text: &text,
	})

	// Add file parts
	for _, filePath := range files {
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %s", filePath)
		}

		// Read file content
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
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

	// Build message request
	msgReq := types.MessageRequest{
		Model:   model,
		Agent:   agent,
		NoReply: noReply,
		System:  system,
		Tools:   tools,
		Parts:   parts,
	}

	resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/message", sessionID), msgReq)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	// Handle empty response (no-reply mode)
	if len(resp) == 0 {
		return "pending", nil
	}

	var result types.MessageWithParts
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Info.ID, nil
}

// detectMimeType detects MIME type based on file extension
func detectMimeType(filePath string) string {
	ext := strings.ToLower(filePath[strings.LastIndex(filePath, "."):])

	mimeTypes := map[string]string{
		// Images
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".bmp":  "image/bmp",
		".svg":  "image/svg+xml",

		// Documents
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",

		// Text
		".txt":  "text/plain",
		".md":   "text/markdown",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".yaml": "application/x-yaml",
		".yml":  "application/x-yaml",

		// Code
		".py":   "text/x-python",
		".go":   "text/x-go",
		".java": "text/x-java",
		".c":    "text/x-c",
		".cpp":  "text/x-c++",
		".h":    "text/x-c",
		".rs":   "text/x-rust",
		".ts":   "text/x-typescript",
		".tsx":  "text/x-typescript",

		// Other
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

	// Default to octet-stream
	return "application/octet-stream"
}
