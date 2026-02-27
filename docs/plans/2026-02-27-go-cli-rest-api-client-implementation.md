# Go CLI REST API 客户端实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标:** 实现一个生产级的 Go CLI REST API 客户端，包含可靠的 HTTP 通信、错误处理、重试机制和直观的命令行界面。

**架构:** 采用混合架构，核心使用 retryablehttp 提供重试功能，服务接口按资源分离，使用 Cobra 构建 CLI。三层错误处理：网络错误、HTTP 状态码错误、业务逻辑错误。

**技术栈:**
- Go 1.21+
- github.com/hashicorp/go-retryablehttp (重试逻辑)
- github.com/spf13/cobra (CLI 框架)
- 标准库 net/http, encoding/json

---

## 阶段 1: 项目初始化与核心客户端

### 任务 1: 项目结构与模块初始化

**文件:**
- 创建: `go.mod`
- 创建: `cmd/cli/main.go`
- 创建: `internal/api/client.go`
- 创建: `internal/api/config.go`
- 创建: `internal/api/errors.go`

**步骤 1: 初始化 Go 模块**

```bash
go mod init github.com/username/opencode_cli
```

**步骤 2: 添加依赖**

```bash
go get github.com/hashicorp/go-retryablehttp
go get github.com/spf13/cobra
```

**步骤 3: 创建主程序入口**

创建 `cmd/cli/main.go`:

```go
package main

import "github.com/username/opencode_cli/cmd"

func main() {
    cmd.Execute()
}
```

**步骤 4: 创建配置结构体**

创建 `internal/api/config.go`:

```go
package api

import "time"

// Config 客户端配置
type Config struct {
    BaseURL    string        `json:"base_url"`
    Timeout    time.Duration `json:"timeout"`
    MaxRetries int           `json:"max_retries"`
    AuthToken  string        `json:"auth_token"`
    Headers    map[string]string `json:"headers"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
    return Config{
        BaseURL:    "https://api.example.com",
        Timeout:    30 * time.Second,
        MaxRetries: 3,
        Headers:    make(map[string]string),
    }
}
```

**步骤 5: 创建错误类型**

创建 `internal/api/errors.go`:

```go
package api

import "fmt"

// APIError 表示 API 返回的错误
type APIError struct {
    StatusCode int    `json:"-"`
    Code       string `json:"code"`
    Message    string `json:"message"`
    Details    any    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
    return fmt.Sprintf("%s: %s (status: %d)", e.Code, e.Message, e.StatusCode)
}

// HTTPError 表示 HTTP 层面的错误
type HTTPError struct {
    StatusCode int
    Status     string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Status)
}

// ConnectionError 表示网络连接错误
type ConnectionError struct {
    Err error
}

func (e *ConnectionError) Error() string {
    return fmt.Sprintf("connection error: %v", e.Err)
}

func (e *ConnectionError) Unwrap() error {
    return e.Err
}
```

**步骤 6: 提交基础结构**

```bash
git add go.mod go.sum cmd/cli/main.go internal/api/
git commit -m "chore: initialize project structure and dependencies"
```

---

### 任务 2: 核心客户端实现

**文件:**
- 修改: `internal/api/client.go`
- 创建: `internal/api/retry.go`
- 创建: `internal/api/request.go`

**步骤 1: 实现客户端结构体**

修改 `internal/api/client.go`:

```go
package api

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    
    retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// Client 是 REST API 客户端
type Client struct {
    baseURL    string
    httpClient *retryablehttp.Client
    config     Config
}

// ClientOption 函数式选项
type ClientOption func(*Client)

// NewClient 创建新的 API 客户端
func NewClient(config Config, opts ...ClientOption) (*Client, error) {
    client := &Client{
        baseURL: config.BaseURL,
        config:  config,
    }
    
    // 创建带重试的 HTTP 客户端
    client.httpClient = newRetryableClient(config)
    
    // 应用选项
    for _, opt := range opts {
        opt(client)
    }
    
    // 设置认证头
    if config.AuthToken != "" {
        client.httpClient.HTTPClient.Transport = &transportWithToken{
            token:     config.AuthToken,
            transport: client.httpClient.HTTPClient.Transport,
        }
    }
    
    return client, nil
}

// WithBaseURL 设置基础 URL
func WithBaseURL(url string) ClientOption {
    return func(c *Client) {
        c.baseURL = url
    }
}

// WithAuthToken 设置认证令牌
func WithAuthToken(token string) ClientOption {
    return func(c *Client) {
        c.config.AuthToken = token
    }
}

// transportWithToken 添加 Bearer Token 的传输层
type transportWithToken struct {
    token     string
    transport http.RoundTripper
}

func (t *transportWithToken) RoundTrip(req *http.Request) (*http.Response, error) {
    req.Header.Set("Authorization", "Bearer "+t.token)
    return t.transport.RoundTrip(req)
}
```

**步骤 2: 实现重试逻辑**

创建 `internal/api/retry.go`:

```go
package api

import (
    "context"
    "net/http"
    
    retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// newRetryableClient 创建带重试的 HTTP 客户端
func newRetryableClient(config Config) *retryablehttp.Client {
    client := retryablehttp.NewClient()
    client.RetryMax = config.MaxRetries
    client.RetryWaitMin = 100 * time.Millisecond
    client.RetryWaitMax = 30 * time.Second
    client.HTTPClient.Timeout = config.Timeout
    
    // 自定义重试条件
    client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
        if resp != nil {
            // 只在特定状态码下重试
            switch resp.StatusCode {
            case 429: // 速率限制
                return true, nil
            case 502, 503, 504: // 服务器错误
                return true, nil
            }
        }
        return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
    }
    
    return client
}
```

**步骤 3: 实现请求/响应处理**

创建 `internal/api/request.go`:

```go
package api

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

// marshalBody 序列化请求体
func (c *Client) marshalBody(v any) (io.Reader, error) {
    if v == nil {
        return nil, nil
    }
    
    data, err := json.Marshal(v)
    if err != nil {
        return nil, fmt.Errorf("marshal request body: %w", err)
    }
    
    return bytes.NewReader(data), nil
}

// unmarshalResponse 反序列化响应体
func (c *Client) unmarshalResponse(resp *http.Response, v any) error {
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("read response body: %w", err)
    }
    
    if err := json.Unmarshal(body, v); err != nil {
        return fmt.Errorf("unmarshal response: %w", err)
    }
    
    return nil
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
    // 构建完整 URL
    u, err := url.JoinPath(c.baseURL, path)
    if err != nil {
        return nil, fmt.Errorf("build url: %w", err)
    }
    
    // 创建请求
    req, err := retryablehttp.NewRequestWithContext(ctx, method, u, body)
    if err != nil {
        return nil, &ConnectionError{Err: err}
    }
    
    // 添加自定义头部
    for k, v := range c.config.Headers {
        req.Header.Set(k, v)
    }
    
    // 如果是 JSON 请求体，设置 Content-Type
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    req.Header.Set("Accept", "application/json")
    
    // 执行请求
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, &ConnectionError{Err: err}
    }
    
    // 检查 HTTP 状态码
    if resp.StatusCode >= 400 {
        defer resp.Body.Close()
        var apiErr APIError
        if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil {
            apiErr.StatusCode = resp.StatusCode
            return nil, &apiErr
        }
        return nil, &HTTPError{
            StatusCode: resp.StatusCode,
            Status:     resp.Status,
        }
    }
    
    return resp, nil
}
```

**步骤 4: 运行基本测试**

创建测试文件 `internal/api/client_test.go`:

```go
package api

import (
    "context"
    "testing"
)

func TestNewClient(t *testing.T) {
    config := DefaultConfig()
    client, err := NewClient(config)
    if err != nil {
        t.Fatalf("NewClient failed: %v", err)
    }
    
    if client == nil {
        t.Fatal("client is nil")
    }
}
```

**步骤 5: 运行测试**

```bash
go test ./internal/api -v
```

**步骤 6: 提交客户端实现**

```bash
git add internal/api/client.go internal/api/retry.go internal/api/request.go internal/api/client_test.go
git commit -m "feat: implement core HTTP client with retry logic"
```

---

## 阶段 2: 服务接口与用户服务实现

### 任务 3: 用户服务接口

**文件:**
- 创建: `internal/api/services/users.go`
- 创建: `internal/api/services/users_test.go`
- 创建: `internal/api/types.go`

**步骤 1: 定义通用类型**

创建 `internal/api/types.go`:

```go
package api

import "time"

// User 用户模型
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role,omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
    Name  string `json:"name,omitempty"`
    Email string `json:"email,omitempty"`
    Role  string `json:"role,omitempty"`
}

// ListOptions 列表查询选项
type ListOptions struct {
    Page  int `json:"page,omitempty"`
    Limit int `json:"limit,omitempty"`
}

// PaginatedResponse 分页响应
type PaginatedResponse[T any] struct {
    Items []T `json:"items"`
    Total int `json:"total"`
    Page  int `json:"page"`
    Limit int `json:"limit"`
}
```

**步骤 2: 实现用户服务接口**

创建 `internal/api/services/users.go`:

```go
package services

import (
    "context"
    "fmt"
    "net/http"
    
    "github.com/username/opencode_cli/internal/api"
)

// UserService 用户服务接口
type UserService interface {
    Get(ctx context.Context, id string) (*api.User, error)
    List(ctx context.Context, opts api.ListOptions) (*api.PaginatedResponse[api.User], error)
    Create(ctx context.Context, req api.CreateUserRequest) (*api.User, error)
    Update(ctx context.Context, id string, req api.UpdateUserRequest) (*api.User, error)
    Delete(ctx context.Context, id string) error
}

// userService 用户服务实现
type userService struct {
    client *api.Client
}

// NewUserService 创建用户服务
func NewUserService(client *api.Client) UserService {
    return &userService{client: client}
}

// Get 获取单个用户
func (s *userService) Get(ctx context.Context, id string) (*api.User, error) {
    resp, err := s.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/users/%s", id), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var user api.User
    if err := s.client.unmarshalResponse(resp, &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

// List 列出用户
func (s *userService) List(ctx context.Context, opts api.ListOptions) (*api.PaginatedResponse[api.User], error) {
    // 构建查询参数
    path := "/users"
    if opts.Page > 0 || opts.Limit > 0 {
        path = fmt.Sprintf("%s?page=%d&limit=%d", path, opts.Page, opts.Limit)
    }
    
    resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result api.PaginatedResponse[api.User]
    if err := s.client.unmarshalResponse(resp, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req api.CreateUserRequest) (*api.User, error) {
    body, err := s.client.marshalBody(req)
    if err != nil {
        return nil, err
    }
    
    resp, err := s.client.doRequest(ctx, http.MethodPost, "/users", body)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var user api.User
    if err := s.client.unmarshalResponse(resp, &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, id string, req api.UpdateUserRequest) (*api.User, error) {
    body, err := s.client.marshalBody(req)
    if err != nil {
        return nil, err
    }
    
    resp, err := s.client.doRequest(ctx, http.MethodPut, fmt.Sprintf("/users/%s", id), body)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var user api.User
    if err := s.client.unmarshalResponse(resp, &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id string) error {
    resp, err := s.client.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/users/%s", id), nil)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNoContent {
        return &api.HTTPError{
            StatusCode: resp.StatusCode,
            Status:     resp.Status,
        }
    }
    
    return nil
}
```

**步骤 3: 编写用户服务测试**

创建 `internal/api/services/users_test.go`:

```go
package services

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/username/opencode_cli/internal/api"
)

func TestUserService_Get(t *testing.T) {
    // 创建测试服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id":"123","name":"Test User","email":"test@example.com","created_at":"2024-01-01T00:00:00Z"}`))
    }))
    defer server.Close()
    
    // 创建客户端
    config := api.DefaultConfig()
    config.BaseURL = server.URL
    client, err := api.NewClient(config)
    if err != nil {
        t.Fatalf("NewClient failed: %v", err)
    }
    
    // 创建用户服务
    userService := NewUserService(client)
    
    // 测试获取用户
    user, err := userService.Get(context.Background(), "123")
    if err != nil {
        t.Fatalf("Get failed: %v", err)
    }
    
    if user.ID != "123" {
        t.Errorf("expected user ID '123', got %s", user.ID)
    }
    if user.Name != "Test User" {
        t.Errorf("expected user name 'Test User', got %s", user.Name)
    }
}
```

**步骤 4: 运行测试**

```bash
go test ./internal/api/services -v
```

**步骤 5: 提交用户服务实现**

```bash
git add internal/api/types.go internal/api/services/
git commit -m "feat: implement user service with CRUD operations"
```

---

## 阶段 3: CLI 集成与命令实现

### 任务 4: CLI 框架与根命令

**文件:**
- 创建: `cmd/root.go`
- 创建: `cmd/context.go`
- 创建: `cmd/output.go`

**步骤 1: 实现根命令**

创建 `cmd/root.go`:

```go
package cmd

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "opencode",
    Short: "OpenCode CLI - REST API client",
    Long: `OpenCode CLI is a command-line tool for interacting with REST APIs.
    
This tool provides a production-grade HTTP client with retry logic,
error handling, and intuitive commands.`,
    PersistentPreRunE: setupContext,
}

// Execute 执行根命令
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    // 全局标志
    rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.opencode.yaml)")
    rootCmd.PersistentFlags().StringP("base-url", "u", "https://api.example.com", "API base URL")
    rootCmd.PersistentFlags().StringP("token", "t", "", "authentication token")
    rootCmd.PersistentFlags().Duration("timeout", 30, "request timeout in seconds")
    rootCmd.PersistentFlags().Int("retries", 3, "maximum retry attempts")
}
```

**步骤 2: 实现上下文管理**

创建 `cmd/context.go`:

```go
package cmd

import (
    "context"
    
    "github.com/spf13/cobra"
    "github.com/username/opencode_cli/internal/api"
    "github.com/username/opencode_cli/internal/api/services"
)

type contextKey string

const (
    clientKey contextKey = "client"
)

// setupContext 设置命令上下文
func setupContext(cmd *cobra.Command, args []string) error {
    // 获取配置值
    baseURL, _ := cmd.Flags().GetString("base-url")
    token, _ := cmd.Flags().GetString("token")
    timeout, _ := cmd.Flags().GetDuration("timeout")
    retries, _ := cmd.Flags().GetInt("retries")
    
    // 创建配置
    config := api.Config{
        BaseURL:    baseURL,
        AuthToken:  token,
        Timeout:    timeout,
        MaxRetries: retries,
        Headers:    make(map[string]string),
    }
    
    // 创建客户端
    client, err := api.NewClient(config)
    if err != nil {
        return err
    }
    
    // 设置到上下文
    ctx := context.WithValue(cmd.Context(), clientKey, client)
    cmd.SetContext(ctx)
    
    return nil
}

// getClientFromContext 从上下文获取客户端
func getClientFromContext(ctx context.Context) *api.Client {
    client, ok := ctx.Value(clientKey).(*api.Client)
    if !ok {
        panic("client not found in context")
    }
    return client
}

// getUserService 获取用户服务
func getUserService(ctx context.Context) services.UserService {
    client := getClientFromContext(ctx)
    return services.NewUserService(client)
}
```

**步骤 3: 实现输出格式化**

创建 `cmd/output.go`:

```go
package cmd

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
    
    "github.com/username/opencode_cli/internal/api"
    "gopkg.in/yaml.v3"
)

// OutputFormat 输出格式
type OutputFormat string

const (
    FormatTable OutputFormat = "table"
    FormatJSON  OutputFormat = "json"
    FormatYAML  OutputFormat = "yaml"
)

// printOutput 打印输出
func printOutput(format OutputFormat, data any) error {
    switch format {
    case FormatJSON:
        return printJSON(data)
    case FormatYAML:
        return printYAML(data)
    case FormatTable:
        return printTable(data)
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}

// printJSON 以 JSON 格式打印
func printJSON(data any) error {
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    return encoder.Encode(data)
}

// printYAML 以 YAML 格式打印
func printYAML(data any) error {
    encoder := yaml.NewEncoder(os.Stdout)
    defer encoder.Close()
    return encoder.Encode(data)
}

// printTable 以表格格式打印
func printTable(data any) error {
    switch v := data.(type) {
    case *api.User:
        return printUserTable(v)
    case *api.PaginatedResponse[api.User]:
        return printUsersTable(v)
    default:
        // 默认回退到 JSON
        return printJSON(data)
    }
}

// printUserTable 打印单个用户表格
func printUserTable(user *api.User) error {
    fmt.Println("ID:", user.ID)
    fmt.Println("Name:", user.Name)
    fmt.Println("Email:", user.Email)
    if user.Role != "" {
        fmt.Println("Role:", user.Role)
    }
    fmt.Println("Created:", user.CreatedAt.Format("2006-01-02 15:04:05"))
    return nil
}

// printUsersTable 打印用户列表表格
func printUsersTable(users *api.PaginatedResponse[api.User]) error {
    fmt.Printf("Showing %d of %d users (page %d)\n\n", len(users.Items), users.Total, users.Page)
    fmt.Println("ID\tName\tEmail\tRole\tCreated")
    fmt.Println("--\t----\t-----\t----\t-------")
    
    for _, user := range users.Items {
        fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
            user.ID,
            user.Name,
            user.Email,
            user.Role,
            user.CreatedAt.Format("2006-01-02"),
        )
    }
    
    return nil
}
```

**步骤 4: 运行 CLI 编译测试**

```bash
go build -o opencode ./cmd/cli
./opencode --help
```

**步骤 5: 提交 CLI 框架**

```bash
git add cmd/root.go cmd/context.go cmd/output.go
git commit -m "feat: implement CLI framework with context and output formatting"
```

---

### 任务 5: 用户命令实现

**文件:**
- 创建: `cmd/users.go`
- 创建: `cmd/users_list.go`
- 创建: `cmd/users_get.go`
- 创建: `cmd/users_create.go`

**步骤 1: 实现用户根命令**

创建 `cmd/users.go`:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(NewCmdUsers())
}

// NewCmdUsers 创建用户命令
func NewCmdUsers() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "users",
        Short: "Manage users",
    }
    
    cmd.AddCommand(
        NewCmdUsersList(),
        NewCmdUsersGet(),
        NewCmdUsersCreate(),
        NewCmdUsersUpdate(),
        NewCmdUsersDelete(),
    )
    
    return cmd
}
```

**步骤 2: 实现用户列表命令**

创建 `cmd/users_list.go`:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

// NewCmdUsersList 创建用户列表命令
func NewCmdUsersList() *cobra.Command {
    var (
        format string
        page   int
        limit  int
    )
    
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List users",
        RunE: func(cmd *cobra.Command, args []string) error {
            userService := getUserService(cmd.Context())
            
            users, err := userService.List(cmd.Context(), api.ListOptions{
                Page:  page,
                Limit: limit,
            })
            if err != nil {
                return fmt.Errorf("list users: %w", err)
            }
            
            return printOutput(OutputFormat(format), users)
        },
    }
    
    cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")
    cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
    cmd.Flags().IntVarP(&limit, "limit", "l", 100, "Maximum number of users to list")
    
    return cmd
}
```

**步骤 3: 实现用户获取命令**

创建 `cmd/users_get.go`:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

// NewCmdUsersGet 创建用户获取命令
func NewCmdUsersGet() *cobra.Command {
    var format string
    
    cmd := &cobra.Command{
        Use:   "get <id>",
        Short: "Get a user by ID",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            userService := getUserService(cmd.Context())
            
            user, err := userService.Get(cmd.Context(), args[0])
            if err != nil {
                return fmt.Errorf("get user: %w", err)
            }
            
            return printOutput(OutputFormat(format), user)
        },
    }
    
    cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")
    
    return cmd
}
```

**步骤 4: 实现用户创建命令**

创建 `cmd/users_create.go`:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

// NewCmdUsersCreate 创建用户创建命令
func NewCmdUsersCreate() *cobra.Command {
    var (
        name  string
        email string
        role  string
        format string
    )
    
    cmd := &cobra.Command{
        Use:   "create",
        Short: "Create a new user",
        RunE: func(cmd *cobra.Command, args []string) error {
            userService := getUserService(cmd.Context())
            
            user, err := userService.Create(cmd.Context(), api.CreateUserRequest{
                Name:  name,
                Email: email,
                Role:  role,
            })
            if err != nil {
                return fmt.Errorf("create user: %w", err)
            }
            
            return printOutput(OutputFormat(format), user)
        },
    }
    
    cmd.Flags().StringVarP(&name, "name", "n", "", "User name (required)")
    cmd.Flags().StringVarP(&email, "email", "e", "", "User email (required)")
    cmd.Flags().StringVarP(&role, "role", "r", "", "User role")
    cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")
    
    cmd.MarkFlagRequired("name")
    cmd.MarkFlagRequired("email")
    
    return cmd
}
```

**步骤 5: 添加缺失的更新和删除命令（占位符）**

创建 `cmd/users_update.go`:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

// NewCmdUsersUpdate 创建用户更新命令
func NewCmdUsersUpdate() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "update <id>",
        Short: "Update a user",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return fmt.Errorf("not implemented yet")
        },
    }
    
    return cmd
}

// NewCmdUsersDelete 创建用户删除命令
func NewCmdUsersDelete() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "delete <id>",
        Short: "Delete a user",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return fmt.Errorf("not implemented yet")
        },
    }
    
    return cmd
}
```

**步骤 6: 测试用户命令**

```bash
go build -o opencode ./cmd/cli
./opencode users --help
./opencode users list --help
./opencode users get --help
./opencode users create --help
```

**步骤 7: 提交用户命令**

```bash
git add cmd/users*.go
git commit -m "feat: implement user commands (list, get, create)"
```

---

## 阶段 4: 高级功能与完善

### 任务 6: 配置管理与环境变量支持

**文件:**
- 创建: `config/config.go`
- 创建: `config/loader.go`
- 修改: `cmd/root.go`

**步骤 1: 实现配置管理**

创建 `config/config.go`:

```go
package config

import (
    "os"
    "path/filepath"
    
    "github.com/spf13/viper"
)

// LoadConfig 加载配置
func LoadConfig(configFile string) error {
    if configFile != "" {
        viper.SetConfigFile(configFile)
    } else {
        // 默认配置路径
        home, err := os.UserHomeDir()
        if err != nil {
            return err
        }
        
        viper.AddConfigPath(".")
        viper.AddConfigPath(filepath.Join(home, ".config"))
        viper.SetConfigName(".opencode")
        viper.SetConfigType("yaml")
    }
    
    // 环境变量支持
    viper.SetEnvPrefix("OPENCODE")
    viper.AutomaticEnv()
    
    // 读取配置
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return err
        }
    }
    
    return nil
}
```

**步骤 2: 更新根命令支持配置**

修改 `cmd/root.go` 中的 `setupContext` 函数:

```go
import (
    // 添加
    "github.com/spf13/viper"
    "github.com/username/opencode_cli/config"
)

func setupContext(cmd *cobra.Command, args []string) error {
    // 加载配置
    configFile, _ := cmd.Flags().GetString("config")
    if err := config.LoadConfig(configFile); err != nil {
        return err
    }
    
    // 合并配置：命令行标志 > 配置文件 > 环境变量 > 默认值
    viper.BindPFlag("base_url", cmd.Flags().Lookup("base-url"))
    viper.BindPFlag("token", cmd.Flags().Lookup("token"))
    viper.BindPFlag("timeout", cmd.Flags().Lookup("timeout"))
    viper.BindPFlag("retries", cmd.Flags().Lookup("retries"))
    
    // 获取配置值
    baseURL := viper.GetString("base_url")
    token := viper.GetString("token")
    timeout := viper.GetDuration("timeout")
    retries := viper.GetInt("retries")
    
    // ... 剩余代码不变
}
```

**步骤 3: 添加 viper 依赖**

```bash
go get github.com/spf13/viper
```

**步骤 4: 创建示例配置文件**

创建 `example.config.yaml`:

```yaml
base_url: "https://api.example.com"
timeout: 30
retries: 3
# token: "your-token-here"  # 可以通过环境变量 OPENCODE_TOKEN 设置
```

**步骤 5: 测试配置加载**

```bash
export OPENCODE_TOKEN="test-token"
go build -o opencode ./cmd/cli
./opencode users list --base-url "https://test.api.com"
```

**步骤 6: 提交配置管理**

```bash
git add config/ example.config.yaml
git commit -m "feat: add config management with viper support"
```

---

### 任务 7: 日志记录与调试支持

**文件:**
- 创建: `internal/logging/logger.go`
- 修改: `internal/api/client.go`

**步骤 1: 实现结构化日志**

创建 `internal/logging/logger.go`:

```go
package logging

import (
    "context"
    "io"
    "os"
    
    "github.com/rs/zerolog"
)

type contextKey string

const (
    loggerKey contextKey = "logger"
)

// NewLogger 创建新的日志记录器
func NewLogger(w io.Writer, level zerolog.Level) zerolog.Logger {
    logger := zerolog.New(w).
        Level(level).
        With().
        Timestamp().
        Logger()
    
    return logger
}

// DefaultLogger 创建默认日志记录器
func DefaultLogger() zerolog.Logger {
    return NewLogger(os.Stderr, zerolog.InfoLevel)
}

// WithLogger 将日志记录器添加到上下文
func WithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

// FromContext 从上下文获取日志记录器
func FromContext(ctx context.Context) zerolog.Logger {
    logger, ok := ctx.Value(loggerKey).(zerolog.Logger)
    if !ok {
        return DefaultLogger()
    }
    return logger
}
```

**步骤 2: 更新客户端支持日志**

修改 `internal/api/client.go`:

```go
import (
    // 添加
    "github.com/rs/zerolog"
    "github.com/username/opencode_cli/internal/logging"
)

type Client struct {
    baseURL    string
    httpClient *retryablehttp.Client
    config     Config
    logger     zerolog.Logger
}

// NewClient 更新支持日志记录器
func NewClient(config Config, opts ...ClientOption) (*Client, error) {
    client := &Client{
        baseURL: config.BaseURL,
        config:  config,
        logger:  logging.DefaultLogger(),
    }
    
    // ... 剩余代码
}

// WithLogger 设置日志记录器选项
func WithLogger(logger zerolog.Logger) ClientOption {
    return func(c *Client) {
        c.logger = logger
    }
}
```

**步骤 3: 添加日志记录到请求处理**

更新 `internal/api/request.go`:

```go
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
    logger := c.logger.With().
        Str("method", method).
        Str("path", path).
        Logger()
    
    logger.Debug().Msg("sending request")
    
    // ... 在错误处添加日志
    if err != nil {
        logger.Error().Err(err).Msg("request failed")
        return nil, &ConnectionError{Err: err}
    }
    
    logger.Debug().
        Int("status", resp.StatusCode).
        Msg("received response")
    
    // ... 剩余代码
}
```

**步骤 4: 添加 zerolog 依赖**

```bash
go get github.com/rs/zerolog
```

**步骤 5: 测试日志记录**

```bash
export OPENCODE_LOG_LEVEL=debug
go build -o opencode ./cmd/cli
./opencode users list --base-url "http://localhost:8080" 2>&1 | head -20
```

**步骤 6: 提交日志支持**

```bash
git add internal/logging/ internal/api/client.go internal/api/request.go
git commit -m "feat: add structured logging with zerolog"
```

---

## 完成与交付

### 任务 8: 最终集成测试与文档

**文件:**
- 创建: `README.md`
- 创建: `CONTRIBUTING.md`
- 创建: `examples/example.go`

**步骤 1: 创建 README**

创建 `README.md`:

```markdown
# OpenCode CLI

A production-grade Go CLI for interacting with REST APIs.

## Features

- ✅ Reliable HTTP client with automatic retry logic
- ✅ Comprehensive error handling (network, HTTP, API errors)
- ✅ Multiple authentication methods (Bearer token, API key, Basic auth)
- ✅ Configurable timeout and retry policies
- ✅ Structured logging with zerolog
- ✅ Multiple output formats (table, JSON, YAML)
- ✅ Configuration file and environment variable support
- ✅ Intuitive CLI with Cobra framework

## Installation

```bash
# From source
go install github.com/username/opencode_cli/cmd/cli@latest

# Or build locally
git clone https://github.com/username/opencode_cli
cd opencode_cli
go build -o opencode ./cmd/cli
```

## Configuration

1. Config file: `~/.config/.opencode.yaml`
2. Environment variables: `OPENCODE_*`
3. Command line flags

Example config:
```yaml
base_url: "https://api.example.com"
timeout: 30
retries: 3
```

## Usage

```bash
# List users
opencode users list --format table

# Get a user
opencode users get 123 --format json

# Create a user
opencode users create --name "John Doe" --email "john@example.com"

# With custom config
opencode --config ~/.myconfig.yaml users list
```

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.
```

**步骤 2: 创建贡献指南**

创建 `CONTRIBUTING.md`:

```markdown
# Contributing to OpenCode CLI

## Development Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Build: `go build ./cmd/cli`
4. Run tests: `go test ./...`

## Project Structure

- `cmd/` - CLI commands and entry points
- `internal/api/` - Core HTTP client and services
- `internal/logging/` - Structured logging
- `config/` - Configuration management
- `tests/` - Test files

## Code Style

- Use `gofmt` for formatting
- Follow Go naming conventions
- Write comprehensive tests
- Document public APIs

## Testing

- Unit tests: `go test ./internal/...`
- Integration tests: `go test ./tests/...`
- Coverage: `go test -cover ./...`
```

**步骤 3: 创建使用示例**

创建 `examples/example.go`:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/username/opencode_cli/internal/api"
    "github.com/username/opencode_cli/internal/api/services"
)

func main() {
    // Create client
    config := api.DefaultConfig()
    config.BaseURL = "https://api.example.com"
    config.AuthToken = "your-token-here"
    
    client, err := api.NewClient(config)
    if err != nil {
        panic(err)
    }
    
    // Create user service
    userService := services.NewUserService(client)
    
    // List users
    ctx := context.Background()
    users, err := userService.List(ctx, api.ListOptions{
        Page:  1,
        Limit: 10,
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Found %d users\n", users.Total)
    for _, user := range users.Items {
        fmt.Printf("- %s (%s)\n", user.Name, user.Email)
    }
}
```

**步骤 4: 运行最终测试**

```bash
go test ./...
go build ./cmd/cli
./opencode --help
```

**步骤 5: 提交最终文档**

```bash
git add README.md CONTRIBUTING.md examples/
git commit -m "docs: add comprehensive documentation and examples"
```

---

## 执行选项

**计划已完成并保存到 `docs/plans/2026-02-27-go-cli-rest-api-client-implementation.md`。两个执行选项：**

**1. 子代理驱动（当前会话）** - 我为每个任务分发新的子代理，在任务间进行代码审查，快速迭代

**2. 并行会话（分离）** - 使用 executing-plans 打开新会话，批量执行并设置检查点

**选择哪种方法？**

**如果选择子代理驱动：**
- **必需子技能：** 使用 superpowers:subagent-driven-development
- 保持当前会话
- 每个任务使用新的子代理 + 代码审查

**如果选择并行会话：**
- 引导用户在工作树中打开新会话
- **必需子技能：** 新会话使用 superpowers:executing-plans