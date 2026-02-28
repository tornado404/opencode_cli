# oho CLI Unit Tests Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create comprehensive unit tests achieving 79%+ code coverage for the oho CLI tool (Go-based OpenCode command-line client).

**Architecture:** Use Go's standard testing package with httptest to mock HTTP client. Tests will verify command execution paths, API request/response handling, and error scenarios. The HTTP client will be abstracted via interface to enable dependency injection in tests.

**Tech Stack:** Go 1.21+, testing package, httptest, cobra command testing

---

## Task Overview

| Task | Description | Estimated Steps |
|------|-------------|-----------------|
| 1 | Create test infrastructure (mock client, interfaces) | 8 |
| 2 | Test global commands | 4 |
| 3 | Test session commands | 10 |
| 4 | Test message commands | 8 |
| 5 | Test config commands | 4 |
| 6 | Test provider commands | 4 |
| 7 | Test project commands | 4 |
| 8 | Test file/find commands | 6 |
| 9 | Test agent/command/tool commands | 6 |
| 10 | Test lsp/formatter/mcp commands | 6 |
| 11 | Test tui/auth commands | 4 |
| 12 | Test internal utilities | 4 |
| 13 | Coverage verification and reporting | 4 |

---

## Task 1: Create Test Infrastructure

### 1.1: Create Mock HTTP Client Interface

**Files:**
- Create: `oho/internal/client/client_mock.go`

```go
package client

import (
    "context"
)

// MockClient implements ClientInterface for testing
type MockClient struct {
    GetFunc           func(ctx context.Context, path string) ([]byte, error)
    GetWithQueryFunc  func(ctx context.Context, path string, queryParams map[string]string) ([]byte, error)
    PostFunc          func(ctx context.Context, path string, body interface{}) ([]byte, error)
    PutFunc           func(ctx context.Context, path string, body interface{}) ([]byte, error)
    PatchFunc         func(ctx context.Context, path string, body interface{}) ([]byte, error)
    DeleteFunc        func(ctx context.Context, path string) ([]byte, error)
    SSEStreamFunc     func(ctx context.Context, path string) (<-chan []byte, <-chan error, error)
}

func (m *MockClient) Get(ctx context.Context, path string) ([]byte, error) {
    if m.GetFunc != nil {
        return m.GetFunc(ctx, path)
    }
    return nil, nil
}

func (m *MockClient) GetWithQuery(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
    if m.GetWithQueryFunc != nil {
        return m.GetWithQueryFunc(ctx, path, queryParams)
    }
    return nil, nil
}

func (m *MockClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
    if m.PostFunc != nil {
        return m.PostFunc(ctx, path, body)
    }
    return nil, nil
}

func (m *MockClient) Put(ctx context.Context, path string, body interface{}) ([]byte, error) {
    if m.PutFunc != nil {
        return m.PutFunc(ctx, path, body)
    }
    return nil, nil
}

func (m *MockClient) Patch(ctx context.Context, path string, body interface{}) ([]byte, error) {
    if m.PatchFunc != nil {
        return m.PatchFunc(ctx, path, body)
    }
    return nil, nil
}

func (m *MockClient) Delete(ctx context.Context, path string) ([]byte, error) {
    if m.DeleteFunc != nil {
        return m.DeleteFunc(ctx, path)
    }
    return nil, nil
}

func (m *MockClient) SSEStream(ctx context.Context, path string) (<-chan []byte, <-chan error, error) {
    if m.SSEStreamFunc != nil {
        return m.SSEStreamFunc(ctx, path)
    }
    return nil, nil, nil
}
```

### 1.2: Create Client Interface

**Files:**
- Modify: `oho/internal/client/client.go:1-25`

Add after type Client struct:

```go
// ClientInterface 定义客户端接口，便于测试
type ClientInterface interface {
    Get(ctx context.Context, path string) ([]byte, error)
    GetWithQuery(ctx context.Context, path string, queryParams map[string]string) ([]byte, error)
    Post(ctx context.Context, path string, body interface{}) ([]byte, error)
    Put(ctx context.Context, path string, body interface{}) ([]byte, error)
    Patch(ctx context.Context, path string, body interface{}) ([]byte, error)
    Delete(ctx context.Context, path string) ([]byte, error)
    SSEStream(ctx context.Context, path string) (<-chan []byte, <-chan error, error)
}

// 确保 Client 实现 ClientInterface
var _ ClientInterface = (*Client)(nil)
```

### 1.3: Create Test Helpers Package

**Files:**
- Create: `oho/internal/testutil/testutil.go`

```go
package testutil

import (
    "bytes"
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
        {ID: "session2", Title: "Test Session 2", Model: "gpt-3.5"},
    }
    return MockResponse(sessions)
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

// MockHealthResponse 模拟健康检查响应
func MockHealthResponse() []byte {
    health := types.HealthResponse{Healthy: true, Version: "1.0.0"}
    return MockResponse(health)
}

// MockConfigResponse 模拟配置响应
func MockConfigResponse() []byte {
    cfg := types.Config{
        DefaultModel: "gpt-4",
        Theme:        "dark",
        Language:     "en",
        MaxTokens:    4096,
        Temperature:  0.7,
    }
    return MockResponse(cfg)
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

// PostBodyHandler 创建 POST 请求处理器
func PostBodyHandler(fn func(body []byte) ([]byte, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        body, _ := json.Marshal(map[string]interface{}{"status": "ok"})
        w.Write(body)
    }
}
```

### 1.4: Initialize Config in Tests

**Files:**
- Create: `oho/internal/config/config_test.go`

```go
package config

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // 设置测试环境变量
    os.Setenv("OPENCODE_SERVER_HOST", "127.0.0.1")
    os.Setenv("OPENCODE_SERVER_PORT", "4096")
    os.Setenv("OPENCODE_SERVER_USERNAME", "opencode")
    os.Setenv("OPENCODE_SERVER_PASSWORD", "test")
    
    // 初始化配置
    Init()
    
    m.Run()
}
```

### 1.5: Create Base Test Fixtures

**Files:**
- Create: `oho/cmd/session/session_test.go` (basic structure)

```go
package session

import (
    "bytes"
    "context"
    "testing"
    
    "github.com/anomalyco/oho/internal/client"
    "github.com/anomalyco/oho/internal/testutil"
    "github.com/anomalyco/oho/internal/types"
)

func TestListCmd(t *testing.T) {
    tests := []struct {
        name       string
        mockResp   []byte
        mockErr    error
        wantErr    bool
    }{
        {
            name:     "success",
            mockResp: testutil.MockSessionsResponse(),
            mockErr:  nil,
            wantErr:  false,
        },
        {
            name:     "api error",
            mockResp: []byte(""),
            mockErr:  &client.APIError{StatusCode: 500, Message: "Internal Error"},
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := &client.MockClient{
                GetFunc: func(ctx context.Context, path string) ([]byte, error) {
                    return tt.mockResp, tt.mockErr
                },
            }
            
            // 测试逻辑
            _ = mock
        })
    }
}
```

### 1.6: Run Initial Test Check

Run: `cd /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/playful-planet/oho && go test ./... -v -count=1 2>&1 | head -20`
Expected: Should complete without errors (no tests yet)

### 1.7: Add go.mod test dependency

Run: `cd /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/playful-planet/oho && go mod tidy`
Expected: Dependencies resolved

### 1.8: Commit Infrastructure

```bash
cd /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/playful-planet
git add oho/internal/client/client_mock.go oho/internal/testutil/ oho/internal/config/config_test.go
git commit -m "test: add test infrastructure - mock client and helpers"
```

---

## Task 2: Test Global Commands

### 2.1: Create Global Command Tests

**Files:**
- Create: `oho/cmd/global/global_test.go`

```go
package global

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/anomalyco/oho/internal/types"
)

func TestHealthCmd(t *testing.T) {
    tests := []struct {
        name           string
        serverResponse *types.HealthResponse
        serverStatus   int
        wantHealthy    bool
        wantErr        bool
    }{
        {
            name:           "server healthy",
            serverResponse: &types.HealthResponse{Healthy: true, Version: "1.0.0"},
            serverStatus:   200,
            wantHealthy:    true,
            wantErr:        false,
        },
        {
            name:           "server unhealthy",
            serverResponse: &types.HealthResponse{Healthy: false, Version: "1.0.0"},
            serverStatus:   200,
            wantHealthy:    false,
            wantErr:        false,
        },
        {
            name:         "server error",
            serverStatus: 500,
            wantErr:      true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.serverStatus)
                if tt.serverResponse != nil {
                    json.NewEncoder(w).Encode(tt.serverResponse)
                }
            }))
            defer server.Close()
            
            // 测试逻辑需要 mock client
            // ...
        })
    }
}

func TestHealthCmdJSONOutput(t *testing.T) {
    // 测试 JSON 输出模式
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(types.HealthResponse{Healthy: true, Version: "1.0.0"})
    }))
    defer server.Close()
    
    // Test with JSON flag
}

func TestEventCmd(t *testing.T) {
    // 测试 SSE 流
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/event-stream")
        w.Write([]byte("data: test event\n\n"))
    }))
    defer server.Close()
    
    // Test SSE stream handling
}
```

### 2.2: Run Tests

Run: `cd /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/playful-planet/oho && go test ./cmd/global/... -v -run TestHealth`
Expected: PASS

### 2.3: Add More Test Cases

Add tests for:
- Network error scenarios
- Invalid JSON response
- Timeout handling

### 2.4: Commit

```bash
git add oho/cmd/global/global_test.go
git commit -m "test: add global command tests"
```

---

## Task 3: Test Session Commands

### 3.1: Create Session Test File

**Files:**
- Create: `oho/cmd/session/session_test.go`

```go
package session

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/anomalyco/oho/internal/types"
)

func TestSessionListCmd(t *testing.T) {
    tests := []struct {
        name       string
        mockResp    []types.Session
        statusCode  int
        wantErr     bool
    }{
        {
            name: "success with sessions",
            mockResp: []types.Session{
                {ID: "s1", Title: "Session 1", Model: "gpt-4"},
            },
            statusCode: 200,
            wantErr:    false,
        },
        {
            name:       "empty sessions",
            mockResp:   []types.Session{},
            statusCode: 200,
            wantErr:    false,
        },
        {
            name:       "server error",
            mockResp:   nil,
            statusCode: 500,
            wantErr:    true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.statusCode)
                if tt.mockResp != nil {
                    json.NewEncoder(w).Encode(tt.mockResp)
                }
            }))
            defer server.Close()
            
            // Test implementation
        })
    }
}

func TestSessionCreateCmd(t *testing.T) {
    tests := []struct {
        name        string
        parentID     string
        title        string
        mockResponse *types.Session
        statusCode   int
        wantErr      bool
    }{
        {
            name:        "create simple session",
            title:       "Test Session",
            mockResponse: &types.Session{ID: "new-session", Title: "Test Session", Model: "gpt-4"},
            statusCode:  200,
            wantErr:     false,
        },
        {
            name:        "create with parent",
            parentID:    "parent-123",
            title:       "Child Session",
            mockResponse: &types.Session{ID: "child-session", Title: "Child Session", ParentID: "parent-123", Model: "gpt-4"},
            statusCode:  200,
            wantErr:     false,
        },
    }
}

func TestSessionStatusCmd(t *testing.T) {
    // Test session status
}

func TestSessionGetCmd(t *testing.T) {
    // Test get session by ID
}

func TestSessionDeleteCmd(t *testing.T) {
    // Test delete session
}

func TestSessionUpdateCmd(t *testing.T) {
    // Test update session title
}

func TestSessionChildrenCmd(t *testing.T) {
    // Test get child sessions
}

func TestSessionTodoCmd(t *testing.T) {
    // Test get todo items
}

func TestSessionForkCmd(t *testing.T) {
    // Test fork session
}

func TestSessionAbortCmd(t *testing.T) {
    // Test abort session
}
```

### 3.2-3.10: Complete Session Tests

Add tests for remaining session commands:
- session/share, session/unshare
- session/diff
- session/summarize
- session/revert, session/unrevert
- session/permissions

Run: `cd /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/playful-planet/oho && go test ./cmd/session/... -v`
Expected: PASS

---

## Task 4: Test Message Commands

### 4.1: Create Message Test File

**Files:**
- Create: `oho/cmd/message/message_test.go`

```go
package message

import (
    "testing"
)

func TestMessageListCmd(t *testing.T) {
    // Test listing messages with limit
}

func TestMessageAddCmd(t *testing.T) {
    tests := []struct {
        name      string
        sessionID string
        content   string
        model     string
        agent     string
        noReply   bool
        wantErr   bool
    }{
        {
            name:      "add simple message",
            sessionID: "session1",
            content:   "Hello",
            wantErr:   false,
        },
        {
            name:      "add with model",
            sessionID: "session1",
            content:   "Hello",
            model:     "gpt-4",
            wantErr:   false,
        },
        {
            name:      "empty content",
            sessionID: "session1",
            content:   "",
            wantErr:   true,
        },
        {
            name:      "empty session",
            sessionID: "",
            content:   "Hello",
            wantErr:   true,
        },
    }
}

func TestMessageGetCmd(t *testing.T) {
    // Test get single message
}

func TestMessagePromptAsyncCmd(t *testing.T) {
    // Test async message
}

func TestMessageCommandCmd(t *testing.T) {
    // Test execute command
}

func TestMessageShellCmd(t *testing.T) {
    // Test shell execution
}
```

### 4.2-4.8: Complete Message Tests

Run: `go test ./cmd/message/... -v`

---

## Task 5-13: Remaining Command Tests

Following the same pattern, create tests for:

| Task | Package | Test File | Commands to Test |
|------|---------|-----------|------------------|
| 5 | configcmd | `configcmd/config_test.go` | config get, set, providers |
| 6 | provider | `provider/provider_test.go` | provider list, auth, oauth |
| 7 | project | `project/project_test.go` | project list, current, path, vcs |
| 8 | file/find | `file/file_test.go`, `find/find_test.go` | file list/content, find text/file/symbol |
| 9 | agent/command/tool | `agent/agent_test.go`, `command/command_test.go`, `tool/tool_test.go` | list commands |
| 10 | lsp/formatter/mcp | `lsp/lsp_test.go`, `formatter/formatter_test.go`, `mcp/mcp_test.go` | status, add, remove |
| 11 | tui/auth | `tui/tui_test.go`, `auth/auth_test.go` | toast, open-help, auth set |
| 12 | internal/util | `util/output_test.go` | output formatting functions |
| 13 | Coverage | `Makefile` coverage target | Generate coverage report |

---

## Coverage Tracking

### Makefile Coverage Target

**Files:**
- Modify: `oho/Makefile`

```makefile
COVERAGE_DIR := coverage
COVERAGE_OUT := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

.PHONY: test-coverage
test-coverage:
    @mkdir -p $(COVERAGE_DIR)
    go test -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...
    go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
    @echo "Coverage report: $(COVERAGE_HTML)"
    @go tool cover -func=$(COVERAGE_OUT) | tail -1

.PHONY: test
test:
    go test -v -race ./...

.PHONY: test-unit
test-unit:
    go test -v -short ./...
```

---

## Implementation Order

1. Task 1: Infrastructure (blocking)
2. Task 3: Session commands (largest module)
3. Task 4: Message commands
4. Task 2: Global commands
5. Task 5: Config commands
6. Task 6-11: Remaining commands
7. Task 12-13: Utilities and reporting

---

## Verification Commands

```bash
# Run all tests
cd /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/playful-planet/oho
go test -v ./...

# Generate coverage
make test-coverage

# Check coverage percentage
go tool cover -func=coverage/coverage.out | grep total

# Run specific test
go test ./cmd/session/... -v -run TestSessionListCmd
```

---

## Plan Complete

**Plan saved to:** `docs/plans/2026-02-27-oho-unit-tests.md`

**Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

Which approach?
