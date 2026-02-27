# OpenCode CLI 工具实施计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 实现一个生产级的 Go CLI 工具，支持子命令、REST API 包装、认证、配置管理和错误处理。

**Architecture:** 使用 Cobra 作为 CLI 框架，Viper 进行配置管理，结构化 HTTP 客户端封装 REST API，分层错误处理和可配置的输出格式化。

**Tech Stack:** Go 1.21+, Cobra, Viper, testify (测试), net/http

---

## 第一阶段：基础框架

### Task 1: 初始化 Go 模块和基础结构

**文件:**
- 创建: `go.mod`
- 创建: `main.go`
- 创建: `cmd/root.go`

**Step 1: 初始化 Go 模块**

```bash
go mod init opencode-cli
```

**Step 2: 添加依赖**

```bash
go get github.com/spf13/cobra@v1.8.0
go get github.com/spf13/viper@v1.18.0
go get github.com/stretchr/testify@v1.9.0
```

**Step 3: 验证依赖安装**

```bash
go mod tidy
go list -m all | grep -E "(cobra|viper|testify)"
```

预期: 显示三个依赖包及其版本

**Step 4: 创建 main.go**

```go
// main.go
package main

import (
	"os"
	"github.com/spf13/cobra"
	"opencode-cli/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

**Step 5: 创建基础根命令**

```go
// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oho",
		Short: "OpenCode CLI tool",
		Long:  "CLI for managing OpenCode sessions and messages",
	}
	return cmd
}
```

**Step 6: 测试基础命令**

```bash
go run main.go --help
```

预期: 显示 "OpenCode CLI tool" 帮助信息

**Step 7: 提交**

```bash
git add go.mod go.sum main.go cmd/
git commit -m "feat: initialize Go module and basic command structure"
```

---

### Task 2: 配置管理模块

**文件:**
- 创建: `internal/config/types.go`
- 创建: `internal/config/manager.go`
- 修改: `cmd/root.go`

**Step 1: 定义配置结构**

```go
// internal/config/types.go
package config

import (
	"time"
)

type Config struct {
	ServerURL    string        `mapstructure:"server_url"`
	APIToken     string        `mapstructure:"api_token"`
	Timeout      time.Duration `mapstructure:"timeout"`
	OutputFormat string        `mapstructure:"output_format"`
}
```

**Step 2: 实现配置加载器**

```go
// internal/config/manager.go
package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetDefault("server_url", "http://localhost:4096")
	viper.SetDefault("timeout", "30s")
	viper.SetDefault("output_format", "text")
	
	viper.SetConfigName("oho")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/oho/")
	viper.AddConfigPath(".")
	
	viper.SetEnvPrefix("OHO")
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}
	
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &cfg, nil
}
```

**Step 3: 添加配置测试**

```go
// internal/config/manager_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:4096", cfg.ServerURL)
	assert.Equal(t, "30s", cfg.Timeout.String())
	assert.Equal(t, "text", cfg.OutputFormat)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("OHO_SERVER_URL", "http://test.example.com")
	defer os.Unsetenv("OHO_SERVER_URL")
	
	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, "http://test.example.com", cfg.ServerURL)
}
```

**Step 4: 运行配置测试**

```bash
go test ./internal/config/...
```

预期: 测试通过

**Step 5: 提交**

```bash
git add internal/config/
git commit -m "feat: add configuration management with Viper"
```

---

### Task 3: REST API 客户端基础

**文件:**
- 创建: `internal/client/client.go`
- 创建: `internal/client/config.go`
- 创建: `internal/client/transport.go`
- 创建: `internal/client/errors.go`

**Step 1: 定义客户端接口**

```go
// internal/client/client.go
package client

import (
	"net/http"
	"net/url"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	token      string
}

func NewClient(baseURL, token string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	
	return &Client{
		baseURL:    parsedURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
	}, nil
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	u := c.baseURL.ResolveReference(&url.URL{Path: path})
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "oho-cli/1.0")
	
	return req, nil
}
```

**Step 2: 实现错误类型**

```go
// internal/client/errors.go
package client

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
	Details    map[string]interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func classifyError(statusCode int, body []byte) error {
	switch statusCode {
	case 401:
		return fmt.Errorf("authentication failed (status %d)", statusCode)
	case 404:
		return fmt.Errorf("resource not found (status %d)", statusCode)
	case 500:
		return fmt.Errorf("server error (status %d)", statusCode)
	default:
		return &APIError{
			StatusCode: statusCode,
			Message:    string(body),
		}
	}
}
```

**Step 3: 添加客户端测试**

```go
// internal/client/client_test.go
package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("http://example.com", "test-token")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "http", c.baseURL.Scheme)
	assert.Equal(t, "example.com", c.baseURL.Host)
}
```

**Step 4: 运行客户端测试**

```bash
go test ./internal/client/...
```

预期: 测试通过

**Step 5: 提交**

```bash
git add internal/client/
git commit -m "feat: add REST API client with error handling"
```

---

## 第二阶段：核心功能

### Task 4: Session 子命令

**文件:**
- 创建: `cmd/session/create.go`
- 创建: `cmd/session/list.go`
- 创建: `cmd/session/delete.go`
- 修改: `cmd/root.go`

**Step 1: 创建 session 命令组**

```go
// cmd/session.go
package cmd

import (
	"github.com/spf13/cobra"
)

func NewSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage OpenCode sessions",
	}
	return cmd
}
```

**Step 2: 实现 session create 命令**

```go
// cmd/session/create.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new session",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreateSession,
}

func runCreateSession(cmd *cobra.Command, args []string) error {
	name := args[0]
	fmt.Printf("Creating session: %s\n", name)
	// TODO: Implement actual API call
	return nil
}

func init() {
	createCmd.Flags().DurationP("timeout", "t", 30*time.Second, "Operation timeout")
}
```

**Step 3: 实现 session list 命令**

```go
// cmd/session/list.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sessions",
	RunE:  runListSessions,
}

func runListSessions(cmd *cobra.Command, args []string) error {
	fmt.Println("Listing sessions...")
	// TODO: Implement actual API call
	return nil
}
```

**Step 4: 注册 session 命令**

```go
// 在 cmd/root.go 中添加
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oho",
		Short: "OpenCode CLI tool",
		Long:  "CLI for managing OpenCode sessions and messages",
	}
	
	// Add subcommands
	cmd.AddCommand(NewSessionCmd())
	
	return cmd
}
```

**Step 5: 测试 session 命令**

```bash
go run main.go session --help
go run main.go session create mysession
go run main.go session list
```

预期: 命令执行成功，显示相应输出

**Step 6: 提交**

```bash
git add cmd/session/
git commit -m "feat: add session subcommands (create, list)"
```

---

### Task 5: Message 子命令

**文件:**
- 创建: `cmd/message/add.go`
- 创建: `cmd/message/list.go`
- 修改: `cmd/root.go`

**Step 1: 创建 message 命令组**

```go
// cmd/message.go
package cmd

import (
	"github.com/spf13/cobra"
)

func NewMessageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message",
		Short: "Manage session messages",
	}
	return cmd
}
```

**Step 2: 实现 message add 命令**

```go
// cmd/message/add.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a message to a session",
	Args:  cobra.ExactArgs(1),
	RunE:  runAddMessage,
}

func runAddMessage(cmd *cobra.Command, args []string) error {
	content := args[0]
	sessionID, _ := cmd.Flags().GetString("session")
	
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	
	fmt.Printf("Adding message to session %s: %s\n", sessionID, content)
	// TODO: Implement actual API call
	return nil
}

func init() {
	addCmd.Flags().StringP("session", "s", "", "Session ID (required)")
	addCmd.MarkFlagRequired("session")
}
```

**Step 3: 实现 message list 命令**

```go
// cmd/message/list.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List messages in a session",
	RunE:  runListMessages,
}

func runListMessages(cmd *cobra.Command, args []string) error {
	sessionID, _ := cmd.Flags().GetString("session")
	
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	
	fmt.Printf("Listing messages for session %s...\n", sessionID)
	// TODO: Implement actual API call
	return nil
}

func init() {
	listCmd.Flags().StringP("session", "s", "", "Session ID (required)")
	listCmd.MarkFlagRequired("session")
}
```

**Step 4: 注册 message 命令**

```go
// 在 cmd/root.go 中添加
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oho",
		Short: "OpenCode CLI tool",
		Long:  "CLI for managing OpenCode sessions and messages",
	}
	
	cmd.AddCommand(NewSessionCmd())
	cmd.AddCommand(NewMessageCmd())
	
	return cmd
}
```

**Step 5: 测试 message 命令**

```bash
go run main.go message --help
go run main.go message add -s session123 "Hello world"
go run main.go message list -s session123
```

预期: 命令执行成功，显示相应输出

**Step 6: 提交**

```bash
git add cmd/message/
git commit -m "feat: add message subcommands (add, list)"
```

---

### Task 6: 输出格式化模块

**文件:**
- 创建: `internal/output/printer.go`
- 创建: `internal/output/formatter.go`
- 修改: `cmd/root.go`

**Step 1: 定义打印机接口**

```go
// internal/output/printer.go
package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type Printer interface {
	Print(data interface{}) error
	PrintError(err error) error
}

type TextPrinter struct {
	Out io.Writer
}

func (p *TextPrinter) Print(data interface{}) error {
	_, err := fmt.Fprintf(p.Out, "%v\n", data)
	return err
}

func (p *TextPrinter) PrintError(err error) error {
	_, err := fmt.Fprintf(p.Out, "Error: %v\n", err)
	return err
}

type JSONPrinter struct {
	Out io.Writer
}

func (p *JSONPrinter) Print(data interface{}) error {
	encoder := json.NewEncoder(p.Out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (p *JSONPrinter) PrintError(err error) error {
	return p.Print(map[string]string{"error": err.Error()})
}
```

**Step 2: 实现打印机工厂**

```go
// internal/output/formatter.go
package output

import (
	"io"
)

func NewPrinter(format string, w io.Writer) Printer {
	switch format {
	case "json":
		return &JSONPrinter{Out: w}
	default:
		return &TextPrinter{Out: w}
	}
}
```

**Step 3: 集成打印机到命令**

```go
// 在 cmd/root.go 中添加全局标志
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oho",
		Short: "OpenCode CLI tool",
		Long:  "CLI for managing OpenCode sessions and messages",
	}
	
	cmd.PersistentFlags().StringP("output", "o", "text", "Output format (text, json)")
	
	return cmd
}
```

**Step 4: 测试输出格式化**

```go
// internal/output/printer_test.go
package output

import (
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTextPrinter(t *testing.T) {
	var buf bytes.Buffer
	p := &TextPrinter{Out: &buf}
	
	err := p.Print("test message")
	assert.NoError(t, err)
	assert.Equal(t, "test message\n", buf.String())
}

func TestJSONPrinter(t *testing.T) {
	var buf bytes.Buffer
	p := &JSONPrinter{Out: &buf}
	
	data := map[string]string{"key": "value"}
	err := p.Print(data)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"key": "value"`)
}
```

**Step 5: 运行输出测试**

```bash
go test ./internal/output/...
```

预期: 测试通过

**Step 6: 提交**

```bash
git add internal/output/
git commit -m "feat: add output formatting with text and JSON support"
```

---

## 第三阶段：高级功能

### Task 7: API 客户端实现

**文件:**
- 修改: `internal/client/client.go`
- 创建: `internal/client/session.go`
- 创建: `internal/client/message.go`

**Step 1: 实现 Session 服务**

```go
// internal/client/session.go
package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Session struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) CreateSession(name string) (*Session, error) {
	req, err := c.newRequest("POST", "/api/sessions")
	if err != nil {
		return nil, err
	}
	
	body := map[string]string{"name": name}
	req.Header.Set("Content-Type", "application/json")
	json.NewEncoder(req.Body).Encode(body)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		return nil, classifyError(resp.StatusCode, resp.Body)
	}
	
	var session Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}
	
	return &session, nil
}

func (c *Client) ListSessions() ([]Session, error) {
	req, err := c.newRequest("GET", "/api/sessions")
	if err != nil {
		return nil, err
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, classifyError(resp.StatusCode, resp.Body)
	}
	
	var sessions []Session
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return nil, err
	}
	
	return sessions, nil
}
```

**Step 2: 实现 Message 服务**

```go
// internal/client/message.go
package client

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (c *Client) AddMessage(sessionID, content string) (*Message, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("/api/sessions/%s/messages", sessionID))
	if err != nil {
		return nil, err
	}
	
	body := map[string]string{"content": content}
	req.Header.Set("Content-Type", "application/json")
	json.NewEncoder(req.Body).Encode(body)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		return nil, classifyError(resp.StatusCode, resp.Body)
	}
	
	var message Message
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, err
	}
	
	return &message, nil
}

func (c *Client) ListMessages(sessionID string) ([]Message, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/sessions/%s/messages", sessionID))
	if err != nil {
		return nil, err
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, classifyError(resp.StatusCode, resp.Body)
	}
	
	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}
	
	return messages, nil
}
```

**Step 3: 集成 API 客户端到命令**

```go
// 更新 cmd/session/create.go
func runCreateSession(cmd *cobra.Command, args []string) error {
	name := args[0]
	
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	client, err := client.NewClient(cfg.ServerURL, cfg.APIToken)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	
	session, err := client.CreateSession(name)
	if err != nil {
		return err
	}
	
	outputFormat, _ := cmd.Flags().GetString("output")
	printer := output.NewPrinter(outputFormat, cmd.OutOrStdout())
	return printer.Print(session)
}
```

**Step 4: 测试 API 集成**

```bash
# 使用模拟服务器测试
go test ./internal/client/... -v
```

**Step 5: 提交**

```bash
git add internal/client/
git commit -m "feat: implement REST API client for sessions and messages"
```

---

### Task 8: 集成测试和验证

**文件:**
- 创建: `tests/integration/cli_test.go`
- 修改: `Makefile` (或创建)

**Step 1: 创建集成测试**

```go
// tests/integration/cli_test.go
package integration

import (
	"bytes"
	"os/exec"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCLI_HelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	assert.NoError(t, err)
	assert.Contains(t, stdout.String(), "OpenCode CLI tool")
}

func TestCLI_SessionCreate(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "session", "create", "test-session")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// 设置环境变量以避免实际 API 调用
	cmd.Env = append(os.Environ(), "OHO_SERVER_URL=http://localhost:9999")
	
	err := cmd.Run()
	// 预期会失败（因为服务器不存在），但验证命令结构正确
	assert.Error(t, err)
}
```

**Step 2: 创建 Makefile**

```makefile
# Makefile
.PHONY: build test lint clean

build:
	go build -o oho main.go

test:
	go test ./... -v

test-unit:
	go test ./internal/... -v

test-integration:
	go test ./tests/integration/... -v

lint:
	golangci-lint run ./...

clean:
	rm -f oho
```

**Step 3: 运行完整测试套件**

```bash
make test
```

预期: 所有测试通过

**Step 4: 构建二进制文件**

```bash
make build
./oho --help
```

预期: 成功构建并显示帮助信息

**Step 5: 提交**

```bash
git add tests/ Makefile
git commit -m "feat: add integration tests and build system"
```

---

## 完成标准验证

1. ✅ `oho --help` 显示完整的帮助信息
2. ✅ `oho session create mysession` 成功创建会话（模拟）
3. ✅ `oho message add -s sessionid "message content"` 成功添加消息（模拟）
4. ✅ 配置文件和环境变量正常工作
5. ✅ 错误情况有清晰的错误消息和适当的退出码
6. ✅ 支持 JSON 和表格输出格式

**计划完成并保存到 `docs/plans/2026-02-28-opencode-cli-implementation.md`。**

## 执行选项

**两个执行选项：**

**1. Subagent-Driven (本会话)** - 我按任务分派新的子代理，任务间进行代码审查，快速迭代

**2. Parallel Session (独立会话)** - 在新的工作树中打开独立会话使用 executing-plans，批量执行并设置检查点

**您希望选择哪种方式？**