# Go CLI REST API 客户端设计文档

## 概述
本项目旨在构建一个生产级的 Go CLI 工具，用于包装 REST API。该客户端需要提供可靠的 HTTP 通信、完善的错误处理、重试机制以及清晰的用户界面。

## 设计目标
1. **可靠性**：生产级重试和错误处理
2. **可维护性**：清晰的架构和接口设计
3. **用户体验**：直观的 CLI 界面和丰富的输出格式
4. **可测试性**：易于测试和模拟的组件

## 架构设计

### 推荐的混合架构
基于对生产级 Go 项目的分析，推荐采用混合架构：

```go
// 主客户端结构
type APIClient struct {
    baseURL    string
    httpClient *retryablehttp.Client // 核心：重试功能
    token      string
    
    // 服务接口（面向资源）
    Users    UserService
    Projects ProjectService
    // ... 其他资源
}
```

### 三种生产级模式分析

#### 1. 标准库基础模式（Kubernetes 风格）
- **特点**：纯标准库，无外部依赖，接口清晰
- **适用场景**：企业级应用，需要严格控制依赖
- **真实案例**：Kubernetes client-go 库

#### 2. Resty 增强模式
- **特点**：丰富的中间件，自动序列化，内置功能完善
- **适用场景**：快速开发，需要丰富功能的项目
- **真实案例**：Apache Airflow Go SDK, Teleport 集成组件

#### 3. RetryableHTTP 模式
- **特点**：内置指数退避重试，高可靠性
- **适用场景**：网络不稳定或 API 有速率限制的场景
- **真实案例**：TruffleHog, Ory 生态项目

## 核心组件设计

### 配置管理
```go
// 配置结构体
type Config struct {
    BaseURL    string
    Timeout    time.Duration
    MaxRetries int
    AuthToken  string
    Headers    map[string]string
}

// 函数式选项模式
type ClientOption func(*Client)
func WithBaseURL(url string) ClientOption
func WithAuthToken(token string) ClientOption
func WithTimeout(d time.Duration) ClientOption
func WithRetries(max int) ClientOption
```

### 认证处理
支持多种认证方式：
1. **Bearer Token**：`Authorization: Bearer <token>`
2. **API Key**：自定义头部或查询参数
3. **Basic Auth**：标准 HTTP 基本认证
4. **OAuth2**：支持令牌刷新

### 请求/响应类型定义
```go
// 强类型请求/响应
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role,omitempty"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// 分页支持
type PaginatedResponse[T any] struct {
    Items []T `json:"items"`
    Total int `json:"total"`
    Page  int `json:"page"`
    Limit int `json:"limit"`
}
```

### 错误处理策略
三层错误处理体系：

1. **网络/连接错误**：`ConnectionError`
2. **HTTP 状态码错误**：`HTTPError` 和 `APIError`
3. **业务逻辑错误**：在具体服务方法中处理

```go
type APIError struct {
    StatusCode int    `json:"-"`
    Code       string `json:"code"`
    Message    string `json:"message"`
    Details    any    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
    return fmt.Sprintf("%s: %s (status: %d)", e.Code, e.Message, e.StatusCode)
}
```

### 重试和超时机制
```go
func NewClientWithRetry(config Config) *retryablehttp.Client {
    client := retryablehttp.NewClient()
    client.RetryMax = config.MaxRetries
    client.RetryWaitMin = 100 * time.Millisecond
    client.RetryWaitMax = 30 * time.Second
    client.HTTPClient.Timeout = config.Timeout
    
    // 自定义重试条件
    client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
        // 只在特定状态码下重试
        if resp != nil {
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

### JSON 序列化最佳实践
```go
// 统一的 JSON 处理
func (c *Client) marshalBody(v any) (io.Reader, error) {
    if v == nil {
        return nil, nil
    }
    
    // 使用标准库的 Marshal，支持自定义标签
    data, err := json.Marshal(v)
    if err != nil {
        return nil, fmt.Errorf("marshal request body: %w", err)
    }
    
    return bytes.NewReader(data), nil
}

func (c *Client) unmarshalResponse(resp *http.Response, v any) error {
    defer resp.Body.Close()
    
    // 读取响应体
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("read response body: %w", err)
    }
    
    // 调试日志
    if c.debug {
        c.logger.Debug("raw response", "status", resp.StatusCode, "body", string(body))
    }
    
    // 解组响应
    if err := json.Unmarshal(body, v); err != nil {
        return fmt.Errorf("unmarshal response: %w", err)
    }
    
    return nil
}
```

## CLI 设计

### 命令结构
```bash
# 命令层次结构
<cli> users list [--format table|json|yaml] [--limit 100]
<cli> users get <id>
<cli> users create --name <name> --email <email>
<cli> users update <id> --role <role>
<cli> users delete <id>

<cli> projects list
<cli> projects create --name <name>
# ... 其他资源
```

### Cobra 集成示例
```go
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

### 输出格式支持
1. **Table**：适合终端查看
2. **JSON**：适合脚本处理
3. **YAML**：适合配置文件
4. **CSV**：适合数据导出

## 目录结构
```
cmd/
  cli/
    main.go           # CLI 入口
    root.go          # Root 命令
    commands/        # 子命令
      users.go
      projects.go
internal/
  api/
    client.go        # 主客户端
    config.go        # 配置
    errors.go        # 错误类型
    retry.go         # 重试逻辑
    services/        # 服务接口
      users.go
      projects.go
  cmd/
    context.go       # CLI 上下文
    output.go        # 输出格式化
pkg/
  utils/
    json.go          # JSON 工具
    validation.go    # 验证工具
docs/
  plans/
    2026-02-27-go-cli-rest-api-client-design.md
tests/
  integration/
    api_test.go
  unit/
    client_test.go
go.mod
go.sum
```

## 测试策略

### 单元测试
- 测试客户端配置和初始化
- 测试错误处理逻辑
- 测试 JSON 序列化/反序列化
- 测试重试逻辑

### 集成测试
- 使用测试服务器（httptest.Server）
- 模拟 API 响应
- 测试完整的请求/响应周期

### 端到端测试
- 针对真实 API 的测试（可选）
- CLI 命令执行测试

## 依赖管理

### 核心依赖
- `github.com/hashicorp/go-retryablehttp`：重试逻辑
- `github.com/spf13/cobra`：CLI 框架
- `github.com/spf13/viper`：配置管理（可选）

### 可选依赖
- `github.com/go-resty/resty/v2`：如果需要 Resty 功能
- `github.com/stretchr/testify`：测试断言
- `github.com/rs/zerolog`：结构化日志

## 开发计划

### 阶段 1：核心客户端
1. 实现基础客户端结构
2. 实现配置管理
3. 实现 HTTP 请求封装
4. 实现错误处理

### 阶段 2：服务接口
1. 定义服务接口
2. 实现用户服务示例
3. 实现项目服务示例
4. 添加测试

### 阶段 3：CLI 集成
1. 集成 Cobra
2. 实现基本命令
3. 添加输出格式化
4. 添加帮助文档

### 阶段 4：高级功能
1. 添加重试逻辑
2. 添加日志记录
3. 添加配置持久化
4. 添加插件支持

## 成功标准
1. 客户端能够稳定地调用目标 REST API
2. 完善的错误处理和用户反馈
3. 直观的 CLI 界面和帮助文档
4. 完整的测试覆盖
5. 符合 Go 项目的最佳实践

---
*设计文档创建时间：2026-02-27*
*基于生产级 Go 项目的最佳实践分析*