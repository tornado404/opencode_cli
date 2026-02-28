package types

// Session 会话类型
type Session struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ParentID  string `json:"parentId,omitempty"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	Model     string `json:"model"`
	Agent     string `json:"agent"`
}

// SessionStatus 会话状态
type SessionStatus struct {
	Status    string `json:"status"`
	IsReady   bool   `json:"isReady"`
	IsWorking bool   `json:"isWorking"`
	MessageID string `json:"messageId,omitempty"`
}

// Message 消息类型
type Message struct {
	ID        string `json:"id"`
	SessionID string `json:"sessionId"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"createdAt"`
	Content   string `json:"content,omitempty"`
}

// Part 消息部分
type Part struct {
	Type string      `json:"type"`
	Text interface{} `json:"text"`
}

// MessageWithParts 带部分的消息
type MessageWithParts struct {
	Info  Message `json:"info"`
	Parts []Part  `json:"parts"`
}

// MessageRequest 消息请求
type MessageRequest struct {
	MessageID string   `json:"messageId,omitempty"`
	Model     string   `json:"model,omitempty"`
	Agent     string   `json:"agent,omitempty"`
	NoReply   bool     `json:"noReply,omitempty"`
	System    string   `json:"system,omitempty"`
	Tools     []string `json:"tools,omitempty"`
	Parts     []Part   `json:"parts"`
}

// CommandRequest 命令请求
type CommandRequest struct {
	MessageID string            `json:"messageId,omitempty"`
	Agent     string            `json:"agent,omitempty"`
	Model     string            `json:"model,omitempty"`
	Command   string            `json:"command"`
	Arguments map[string]string `json:"arguments,omitempty"`
}

// ShellRequest Shell 命令请求
type ShellRequest struct {
	Agent   string `json:"agent"`
	Model   string `json:"model,omitempty"`
	Command string `json:"command"`
}

// Config 配置类型
type Config struct {
	Providers    map[string]interface{} `json:"providers"`
	DefaultModel string                 `json:"defaultModel"`
	Theme        string                 `json:"theme"`
	Language     string                 `json:"language"`
	AutoApprove  []string               `json:"autoApprove"`
	MaxTokens    int                    `json:"maxTokens"`
	Temperature  float64                `json:"temperature"`
}

// Provider 提供商类型
type Provider struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	BaseURL  string   `json:"baseURL"`
	Models   []string `json:"models"`
	AuthType string   `json:"authType"`
}

// ProviderAuthMethod 提供商认证方式
type ProviderAuthMethod struct {
	Type        string `json:"type"`
	URL         string `json:"url,omitempty"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// ProviderAuthAuthorization OAuth 授权响应
type ProviderAuthAuthorization struct {
	URL           string `json:"url"`
	State         string `json:"state"`
	CodeChallenge string `json:"codeChallenge"`
}

// Project 项目类型
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Vcs  string `json:"vcs"`
}

// Path 路径类型
type Path struct {
	Current string `json:"current"`
	Home    string `json:"home"`
	IsGit   bool   `json:"isGit"`
}

// VcsInfo VCS 信息
type VcsInfo struct {
	Type    string `json:"type"`
	Branch  string `json:"branch"`
	Commit  string `json:"commit"`
	Remote  string `json:"remote"`
	IsDirty bool   `json:"isDirty"`
}

// FileNode 文件节点
type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Type     string     `json:"type"`
	Children []FileNode `json:"children,omitempty"`
}

// File 文件类型 (用于文件状态)
type File struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

// FileContent 文件内容
type FileContent struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// FileDiff 文件差异
type FileDiff struct {
	Path   string `json:"path"`
	Before string `json:"before"`
	After  string `json:"after"`
	Status string `json:"status"`
}

// Todo 待办事项
type Todo struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	MessageID string `json:"messageId"`
}

// Agent 代理类型
type Agent struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`
}

// Command 命令类型
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

// ToolIDs 工具 ID 列表
type ToolIDs struct {
	IDs []string `json:"ids"`
}

// ToolList 工具列表
type ToolList struct {
	Tools []Tool `json:"tools"`
}

// Tool 工具定义
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

// LSPStatus LSP 服务器状态
type LSPStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Port   int    `json:"port"`
}

// FormatterStatus 格式化器状态
type FormatterStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// MCPStatus MCP 服务器状态
type MCPStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// MCPConfig MCP 服务器配置
type MCPConfig struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Healthy bool   `json:"healthy"`
	Version string `json:"version"`
}

// LogRequest 日志请求
type LogRequest struct {
	Service string                 `json:"service"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

// TUIToastRequest Toast 请求
type TUIToastRequest struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message"`
	Variant string `json:"variant"`
}

// TUICommandRequest TUI 命令请求
type TUICommandRequest struct {
	Command string `json:"command"`
}

// AuthRequest 认证请求
type AuthRequest struct {
	ProviderID  string            `json:"providerId"`
	Credentials map[string]string `json:"credentials"`
}

// Symbol 工作区符号
type Symbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Container string `json:"container,omitempty"`
}

// FindMatch 查找匹配
type FindMatch struct {
	Path           string     `json:"path"`
	LineNumber     int        `json:"line_number"`
	Lines          string     `json:"lines"`
	AbsoluteOffset int        `json:"absolute_offset"`
	Submatches     []Submatch `json:"submatches"`
}

// Submatch 子匹配
type Submatch struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Event 事件类型
type Event struct {
	Type string      `json:"type"`
	Text interface{} `json:"text"`
}
