package mcpserver

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anomalyco/oho/internal/client"
	"github.com/anomalyco/oho/internal/config"
)

// Cmd MCP Server 命令
var Cmd = &cobra.Command{
	Use:   "mcpserver",
	Short: "启动 MCP 服务器",
	Long:  "以 MCP 协议启动服务器，允许外部 MCP 客户端调用 OpenCode API",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMCPServer()
	},
}

// JSON-RPC 2.0 types
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// MCP types
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type CallToolResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError"`
}

func runMCPServer() error {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "警告：配置初始化失败：%v\n", err)
	}

	// 创建 scanner 读取 stdin
	scanner := bufio.NewScanner(os.Stdin)

	// 设置 scanner 分割函数 - 按行读取 JSON-RPC 消息
	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		// 找到换行符作为消息分隔符
		if i := strings.Index(string(data), "\n"); i >= 0 {
			return i + 1, data[:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
	scanner.Split(splitFunc)

	// 服务器状态
	var initialized bool
	capabilities := map[string]interface{}{
		"tools": struct{}{},
	}

	// 处理每条消息
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(nil, -32700, "Parse error")
			continue
		}

		// 处理请求
		switch req.Method {
		case "initialize":
			var params InitializeParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				sendError(req.ID, -32600, "Invalid params")
				continue
			}

			result := InitializeResult{
				ProtocolVersion: "2024-11-05",
				Capabilities:    capabilities,
				ServerInfo: ServerInfo{
					Name:    "oho",
					Version: "1.0.0",
				},
			}
			sendResult(req.ID, result)
			initialized = true

		case "notifications/initialized":
			// 客户端初始化完成通知，不需要响应

		case "tools/list":
			if !initialized {
				sendError(req.ID, -32000, "Server not initialized")
				continue
			}
			result := ToolsListResult{
				Tools: getToolsList(),
			}
			sendResult(req.ID, result)

		case "tools/call":
			if !initialized {
				sendError(req.ID, -32000, "Server not initialized")
				continue
			}
			var params CallToolParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				sendError(req.ID, -32600, "Invalid params")
				continue
			}
			result := handleToolCall(params.Name, params.Arguments)
			sendResult(req.ID, result)

		case "ping":
			sendResult(req.ID, map[string]string{"status": "pong"})

		default:
			sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "读取错误: %v\n", err)
	}

	return nil
}

func sendResult(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
	}
	if result != nil {
		resp.Result, _ = json.Marshal(result)
	}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}

func sendError(id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}

func getToolsList() []Tool {
	return []Tool{
		{
			Name:        "session_list",
			Description: "列出所有 OpenCode 会话",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
		},
		{
			Name:        "session_create",
			Description: "创建新的 OpenCode 会话",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"title": {"type": "string"}, "path": {"type": "string"}}, "required": []}`),
		},
		{
			Name:        "session_get",
			Description: "获取指定会话的详细信息",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"sessionId": {"type": "string"}}, "required": ["sessionId"]}`),
		},
		{
			Name:        "session_delete",
			Description: "删除指定会话",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"sessionId": {"type": "string"}}, "required": ["sessionId"]}`),
		},
		{
			Name:        "session_status",
			Description: "获取所有会话的状态",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
		},
		{
			Name:        "message_list",
			Description: "列出指定会话的所有消息",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"sessionId": {"type": "string"}}, "required": ["sessionId"]}`),
		},
		{
			Name:        "message_add",
			Description: "向指定会话发送消息",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"sessionId": {"type": "string"}, "content": {"type": "string"}}, "required": ["sessionId", "content"]}`),
		},
		{
			Name:        "config_get",
			Description: "获取 OpenCode 配置",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
		},
		{
			Name:        "project_list",
			Description: "列出所有项目",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
		},
		{
			Name:        "project_current",
			Description: "获取当前项目",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
		},
		{
			Name:        "provider_list",
			Description: "列出所有可用的 AI 提供商",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
		},
		{
			Name:        "file_list",
			Description: "列出指定目录的文件",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"path": {"type": "string"}}, "required": []}`),
		},
		{
			Name:        "file_content",
			Description: "读取指定文件的内容",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"path": {"type": "string"}}, "required": ["path"]}`),
		},
		{
			Name:        "find_text",
			Description: "在项目中搜索文本",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"pattern": {"type": "string"}}, "required": ["pattern"]}`),
		},
		{
			Name:        "find_file",
			Description: "根据文件名搜索文件",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {"query": {"type": "string"}}, "required": ["query"]}`),
		},
		{
			Name:        "global_health",
			Description: "检查 OpenCode Server 健康状态",
			InputSchema: json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
		},
	}
}

func handleToolCall(name string, args map[string]interface{}) CallToolResult {
	ctx := context.Background()

	switch name {
	case "session_list":
		return handleSessionList(ctx)
	case "session_create":
		return handleSessionCreate(ctx, args)
	case "session_get":
		return handleSessionGet(ctx, args)
	case "session_delete":
		return handleSessionDelete(ctx, args)
	case "session_status":
		return handleSessionStatus(ctx)
	case "message_list":
		return handleMessageList(ctx, args)
	case "message_add":
		return handleMessageAdd(ctx, args)
	case "config_get":
		return handleConfigGet(ctx)
	case "project_list":
		return handleProjectList(ctx)
	case "project_current":
		return handleProjectCurrent(ctx)
	case "provider_list":
		return handleProviderList(ctx)
	case "file_list":
		return handleFileList(ctx, args)
	case "file_content":
		return handleFileContent(ctx, args)
	case "find_text":
		return handleFindText(ctx, args)
	case "find_file":
		return handleFindFile(ctx, args)
	case "global_health":
		return handleGlobalHealth(ctx)
	default:
		return CallToolResult{
			Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", name)}},
			IsError: true,
		}
	}
}

// Tool handlers

func handleSessionList(ctx context.Context) CallToolResult {
	c := client.NewClient()
	resp, err := c.Get(ctx, "/session")
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleSessionCreate(ctx context.Context, args map[string]interface{}) CallToolResult {
	c := client.NewClient()
	req := map[string]interface{}{}

	if title, ok := args["title"].(string); ok && title != "" {
		req["title"] = title
	}
	if path, ok := args["path"].(string); ok && path != "" {
		req["path"] = path
	}

	resp, err := c.Post(ctx, "/session", req)
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleSessionGet(ctx context.Context, args map[string]interface{}) CallToolResult {
	sessionID, ok := args["sessionId"].(string)
	if !ok || sessionID == "" {
		return errorResult("sessionId is required")
	}

	c := client.NewClient()
	resp, err := c.Get(ctx, fmt.Sprintf("/session/%s", sessionID))
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleSessionDelete(ctx context.Context, args map[string]interface{}) CallToolResult {
	sessionID, ok := args["sessionId"].(string)
	if !ok || sessionID == "" {
		return errorResult("sessionId is required")
	}

	c := client.NewClient()
	resp, err := c.Delete(ctx, fmt.Sprintf("/session/%s", sessionID))
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(fmt.Sprintf("Session %s deleted: %s", sessionID, string(resp)))
}

func handleSessionStatus(ctx context.Context) CallToolResult {
	c := client.NewClient()
	resp, err := c.Get(ctx, "/session/status")
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleMessageList(ctx context.Context, args map[string]interface{}) CallToolResult {
	sessionID, ok := args["sessionId"].(string)
	if !ok || sessionID == "" {
		return errorResult("sessionId is required")
	}

	c := client.NewClient()
	resp, err := c.Get(ctx, fmt.Sprintf("/session/%s/message", sessionID))
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleMessageAdd(ctx context.Context, args map[string]interface{}) CallToolResult {
	sessionID, ok := args["sessionId"].(string)
	if !ok || sessionID == "" {
		return errorResult("sessionId is required")
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		return errorResult("content is required")
	}

	c := client.NewClient()
	req := map[string]interface{}{
		"parts": []map[string]interface{}{
			{"type": "text", "text": content},
		},
	}

	resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/message", sessionID), req)
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleConfigGet(ctx context.Context) CallToolResult {
	c := client.NewClient()
	resp, err := c.Get(ctx, "/config")
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleProjectList(ctx context.Context) CallToolResult {
	c := client.NewClient()
	resp, err := c.Get(ctx, "/project")
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleProjectCurrent(ctx context.Context) CallToolResult {
	c := client.NewClient()
	resp, err := c.Get(ctx, "/project/current")
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleProviderList(ctx context.Context) CallToolResult {
	c := client.NewClient()
	resp, err := c.Get(ctx, "/provider")
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleFileList(ctx context.Context, args map[string]interface{}) CallToolResult {
	path := ""
	if p, ok := args["path"].(string); ok {
		path = p
	}

	c := client.NewClient()
	var resp []byte
	var err error

	if path != "" {
		resp, err = c.Get(ctx, fmt.Sprintf("/file?path=%s", path))
	} else {
		resp, err = c.Get(ctx, "/file")
	}
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleFileContent(ctx context.Context, args map[string]interface{}) CallToolResult {
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return errorResult("path is required")
	}

	c := client.NewClient()
	resp, err := c.Get(ctx, fmt.Sprintf("/file/%s", strings.TrimPrefix(path, "/")))
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleFindText(ctx context.Context, args map[string]interface{}) CallToolResult {
	pattern, ok := args["pattern"].(string)
	if !ok || pattern == "" {
		return errorResult("pattern is required")
	}

	c := client.NewClient()
	resp, err := c.Get(ctx, fmt.Sprintf("/find/text?q=%s", pattern))
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleFindFile(ctx context.Context, args map[string]interface{}) CallToolResult {
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return errorResult("query is required")
	}

	c := client.NewClient()
	resp, err := c.Get(ctx, fmt.Sprintf("/find/file?q=%s", query))
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func handleGlobalHealth(ctx context.Context) CallToolResult {
	c := client.NewClient()
	resp, err := c.Get(ctx, "/global/health")
	if err != nil {
		return errorResult(err.Error())
	}
	return successResult(string(resp))
}

func successResult(content string) CallToolResult {
	return CallToolResult{
		Content: []ToolContent{{Type: "text", Text: content}},
		IsError: false,
	}
}

func errorResult(errMsg string) CallToolResult {
	return CallToolResult{
		Content: []ToolContent{{Type: "text", Text: errMsg}},
		IsError: true,
	}
}
