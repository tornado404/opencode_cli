package testutil

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/anomalyco/oho/internal/types"
)

// NewMockServer 创建模拟 HTTP 服务器
func NewMockServer(handlers map[string]http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	for path, handler := range handlers {
		mux.Handle(path, handler)
	}
	return httptest.NewServer(mux)
}

// MockResponse 创建 JSON 响应
func MockResponse(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

// MockSessionsResponse 模拟会话列表响应
func MockSessionsResponse() []byte {
	sessions := []types.Session{
		{ID: "session1", Title: "Test Session 1", Model: "gpt-4"},
		{ID: "session2", Title: "Test Session 2", Model: "gpt-3.5-turbo"},
	}
	return MockResponse(sessions)
}

// MockSessionResponse 模拟单个会话响应
func MockSessionResponse() []byte {
	session := types.Session{
		ID:        "session1",
		Title:     "Test Session",
		Model:     "gpt-4",
		Agent:     "default",
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	}
	return MockResponse(session)
}

// MockSessionStatusResponse 模拟会话状态响应
func MockSessionStatusResponse() []byte {
	status := map[string]types.SessionStatus{
		"session1": {Status: "idle", IsReady: true, IsWorking: false},
		"session2": {Status: "working", IsReady: true, IsWorking: true, MessageID: "msg1"},
	}
	return MockResponse(status)
}

// MockMessagesResponse 模拟消息列表响应
func MockMessagesResponse() []byte {
	messages := []types.MessageWithParts{
		{Info: types.Message{ID: "msg1", Role: "user", Content: "Hello"}},
		{Info: types.Message{ID: "msg2", Role: "assistant", Content: "Hi there"}},
	}
	return MockResponse(messages)
}

// MockMessageResponse 模拟单个消息响应
func MockMessageResponse() []byte {
	msg := types.MessageWithParts{
		Info: types.Message{
			ID:        "msg1",
			SessionID: "session1",
			Role:      "user",
			Content:   "Hello world",
			CreatedAt: 1234567890,
		},
		Parts: []types.Part{
			{Type: "text", Data: "Hello world"},
		},
	}
	return MockResponse(msg)
}

// MockHealthResponse 模拟健康检查响应
func MockHealthResponse() []byte {
	health := types.HealthResponse{Healthy: true, Version: "1.0.0"}
	return MockResponse(health)
}

// MockConfigResponse 模拟配置响应
func MockConfigResponse() []byte {
	cfg := types.Config{
		Providers:    map[string]interface{}{},
		DefaultModel: "gpt-4",
		Theme:        "dark",
		Language:     "en",
		AutoApprove:  []string{},
		MaxTokens:    4096,
		Temperature:  0.7,
	}
	return MockResponse(cfg)
}

// MockProvidersResponse 模拟提供商列表响应
func MockProvidersResponse() []byte {
	providers := []types.Provider{
		{ID: "openai", Name: "OpenAI", BaseURL: "https://api.openai.com", Models: []string{"gpt-4", "gpt-3.5-turbo"}, AuthType: "api_key"},
		{ID: "anthropic", Name: "Anthropic", BaseURL: "https://api.anthropic.com", Models: []string{"claude-3"}, AuthType: "api_key"},
	}
	return MockResponse(providers)
}

// MockProjectsResponse 模拟项目列表响应
func MockProjectsResponse() []byte {
	projects := []types.Project{
		{ID: "proj1", Name: "Project 1", Path: "/home/user/project1", Vcs: "git"},
		{ID: "proj2", Name: "Project 2", Path: "/home/user/project2", Vcs: "none"},
	}
	return MockResponse(projects)
}

// MockPathResponse 模拟路径响应
func MockPathResponse() []byte {
	path := types.Path{
		Current: "/home/user/current",
		Home:    "/home/user",
		IsGit:   true,
	}
	return MockResponse(path)
}

// MockVCSResponse 模拟 VCS 信息响应
func MockVCSResponse() []byte {
	vcs := types.VcsInfo{
		Type:    "git",
		Branch:  "main",
		Commit:  "abc123",
		Remote:  "origin",
		IsDirty: false,
	}
	return MockResponse(vcs)
}

// MockFileListResponse 模拟文件列表响应
func MockFileListResponse() []byte {
	files := []types.FileNode{
		{Name: "src", Path: "/src", Type: "directory"},
		{Name: "main.go", Path: "/main.go", Type: "file"},
	}
	return MockResponse(files)
}

// MockFileContentResponse 模拟文件内容响应
func MockFileContentResponse() []byte {
	content := types.FileContent{
		Path:     "/main.go",
		Content:  "package main\n\nfunc main() {}",
		Encoding: "utf-8",
	}
	return MockResponse(content)
}

// MockFileStatusResponse 模拟文件状态响应
func MockFileStatusResponse() []byte {
	files := []types.File{
		{Path: "main.go", Status: "modified"},
		{Path: "go.mod", Status: "unchanged"},
	}
	return MockResponse(files)
}

// MockAgentsResponse 模拟代理列表响应
func MockAgentsResponse() []byte {
	agents := []types.Agent{
		{ID: "default", Name: "Default Agent", Description: "Default programming agent", Tools: []string{"Read", "Edit", "Bash"}},
		{ID: "review", Name: "Review Agent", Description: "Code review agent", Tools: []string{"Read", "Grep"}},
	}
	return MockResponse(agents)
}

// MockCommandsResponse 模拟命令列表响应
func MockCommandsResponse() []byte {
	commands := []types.Command{
		{Name: "test", Description: "Run tests", Usage: "test [package]"},
		{Name: "build", Description: "Build the project", Usage: "build [target]"},
	}
	return MockResponse(commands)
}

// MockToolIDsResponse 模拟工具 ID 列表响应
func MockToolIDsResponse() []byte {
	ids := types.ToolIDs{
		IDs: []string{"Read", "Edit", "Write", "Bash", "Grep", "Glob"},
	}
	return MockResponse(ids)
}

// MockToolListResponse 模拟工具列表响应
func MockToolListResponse() []byte {
	tools := types.ToolList{
		Tools: []types.Tool{
			{Name: "Read", Description: "Read file contents", Schema: map[string]interface{}{}},
			{Name: "Edit", Description: "Edit file contents", Schema: map[string]interface{}{}},
		},
	}
	return MockResponse(tools)
}

// MockLSPStatusResponse 模拟 LSP 状态响应
func MockLSPStatusResponse() []byte {
	status := []types.LSPStatus{
		{Name: "gopls", Status: "running", Port: 1234},
		{Name: "tsserver", Status: "stopped", Port: 0},
	}
	return MockResponse(status)
}

// MockFormatterStatusResponse 模拟格式化器状态响应
func MockFormatterStatusResponse() []byte {
	status := []types.FormatterStatus{
		{Name: "gofmt", Status: "available"},
		{Name: "prettier", Status: "available"},
	}
	return MockResponse(status)
}

// MockMCPStatusResponse 模拟 MCP 状态响应
func MockMCPStatusResponse() []byte {
	status := []types.MCPStatus{
		{Name: "filesystem", Status: "running"},
		{Name: "github", Status: "error", Error: "auth required"},
	}
	return MockResponse(status)
}

// MockTodoResponse 模拟待办事项响应
func MockTodoResponse() []byte {
	todos := []types.Todo{
		{ID: "todo1", Content: "Implement feature X", Status: "pending", MessageID: "msg1"},
		{ID: "todo2", Content: "Fix bug Y", Status: "completed", MessageID: "msg2"},
	}
	return MockResponse(todos)
}

// MockDiffResponse 模拟差异响应
func MockDiffResponse() []byte {
	diffs := []types.FileDiff{
		{Path: "main.go", Before: "func main() {}", After: "func main() {\n}", Status: "modified"},
	}
	return MockResponse(diffs)
}

// MockSymbolsResponse 模拟符号搜索响应
func MockSymbolsResponse() []byte {
	symbols := []types.Symbol{
		{Name: "main", Kind: "function", Path: "main.go", Line: 10, Column: 1},
		{Name: "Config", Kind: "struct", Path: "config.go", Line: 5, Column: 1},
	}
	return MockResponse(symbols)
}

// MockFindMatchesResponse 模拟查找匹配响应
func MockFindMatchesResponse() []byte {
	matches := []types.FindMatch{
		{Path: "main.go", LineNumber: 10, Lines: "func main() {}", AbsoluteOffset: 100, Submatches: []types.Submatch{{Start: 0, End: 4}}},
	}
	return MockResponse(matches)
}

// MockBoolResponse 创建布尔响应
func MockBoolResponse(b bool) []byte {
	return MockResponse(b)
}

// ErrorHandler 创建错误响应处理器
func ErrorHandler(statusCode int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(message))
	}
}

// JSONHandler 创建 JSON 响应处理器
func JSONHandler(v interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	}
}

// HandlerFuncForPath 创建特定路径的处理器
func HandlerFuncForPath(path string, fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return fn
}

// NewTestContext 创建测试用的 context
func NewTestContext() context.Context {
	return context.Background()
}
