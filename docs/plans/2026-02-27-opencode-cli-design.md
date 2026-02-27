# OpenCode CLI 工具架构设计
## 2026-02-27

## 项目概述

**目标**：构建一个生产级的 Go CLI 工具，用于管理 OpenCode 会话和消息，支持子命令、REST API 包装、认证、配置管理和错误处理。

**核心需求**：
1. 支持子命令如 `oho session create`、`oho message add -s sessionxxx`
2. 封装 REST API 调用
3. 处理认证、配置和错误处理
4. 遵循生产级 Go CLI 工具最佳实践

## 架构设计

### 1. 项目结构

```
oho/                           # 主包
├── cmd/                       # 命令定义 (Cobra)
│   ├── root.go               # 根命令和全局配置
│   ├── session/              # 会话子命令
│   │   ├── create.go
│   │   ├── list.go
│   │   └── delete.go
│   └── message/              # 消息子命令
│       ├── add.go
│       └── list.go
├── internal/
│   ├── client/               # REST API 客户端
│   │   ├── client.go         # 主客户端接口
│   │   ├── config.go         # 客户端配置
│   │   └── transport.go      # HTTP 传输层
│   ├── config/               # 配置管理
│   │   ├── manager.go
│   │   └── types.go
│   ├── auth/                 # 认证管理
│   │   ├── token.go
│   │   └── credentials.go
│   └── output/               # 输出格式化
│       ├── formatter.go
│       └── printer.go
├── pkg/
│   └── oho/                  # 公共API（如需要）
└── main.go
```

### 2. 命令框架选择

**选择：Cobra（推荐方案）**

**理由**：
- **生产验证**：Kubernetes (kubectl)、Docker、Hugo、GitHub CLI等使用
- **强大功能**：自动生成帮助、bash/zsh补全、子命令嵌套、标志解析
- **生态完整**：与Viper（配置管理）完美集成
- **模式成熟**：清晰的代码组织模式

**权衡**：
- ✅ 企业级功能完善
- ✅ 社区支持强大
- ✅ 代码组织清晰
- ⚠️ 学习曲线稍陡（但文档丰富）
- ⚠️ 依赖较多

**备选方案**：
- **urfave/cli**：更轻量，API更简洁
- **标准库flag**：最小依赖，但需要更多模板代码

### 3. 子命令组织模式

**生产级模式**（参考Kubernetes/kubectl）：

```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "oho",
    Short: "OpenCode CLI tool",
    Long:  "CLI for managing OpenCode sessions and messages",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        return initializeApp() // 初始化配置、客户端等
    },
}

// cmd/session/create.go
var createCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new session",
    Args:  cobra.ExactArgs(1),
    RunE:  runCreateSession,
}

func runCreateSession(cmd *cobra.Command, args []string) error {
    timeout, _ := cmd.Flags().GetDuration("timeout")
    return client.Session().Create(args[0], timeout)
}
```

### 4. 标志处理模式

**短标志（-s）和长标志（--session）支持**：

```go
// cmd/message/add.go
func init() {
    addCmd.Flags().StringP("session", "s", "", "Session ID (required)")
    addCmd.MarkFlagRequired("session")
    addCmd.Flags().StringP("content", "c", "", "Message content")
    addCmd.Flags().Bool("json", false, "Output as JSON")
}

// 标志值获取
sessionID, _ := cmd.Flags().GetString("session")
outputJSON, _ := cmd.Flags().GetBool("json")
```

**最佳实践**：
- 全局标志使用 `PersistentFlags()`
- 使用 `MarkFlagRequired()` 标记必填标志
- 支持环境变量覆盖：`cobra.OnInitialize(initConfig)`

### 5. REST API 客户端包装模式

**结构化客户端模式**（参考Kubernetes client-go）：

```go
// internal/client/client.go
type Client struct {
    baseURL    *url.URL
    httpClient *http.Client
    token      string
    
    Session *SessionService
    Message *MessageService
}

type SessionService struct {
    client *Client
}

func (s *SessionService) Create(name string, timeout time.Duration) (*Session, error) {
    // REST调用实现
}

// 工厂函数
func NewClient(baseURL string, token string) (*Client, error) {
    c := &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
    c.Session = &SessionService{client: c}
    c.Message = &MessageService{client: c}
    return c, nil
}
```

**高级特性**：
- 请求/响应拦截器（日志、重试、认证刷新）
- 连接池配置
- 超时控制（请求/连接/空闲超时）
- 错误类型定义

### 6. 配置管理模式

**Viper + 环境变量 + 配置文件**：

```go
// internal/config/manager.go
type Config struct {
    ServerURL    string        `mapstructure:"server_url" yaml:"server_url"`
    APIToken     string        `mapstructure:"api_token" yaml:"api_token"`
    Timeout      time.Duration `mapstructure:"timeout" yaml:"timeout"`
    OutputFormat string        `mapstructure:"output_format" yaml:"output_format"`
}

func Load() (*Config, error) {
    viper.SetDefault("server_url", "http://localhost:4096")
    viper.SetDefault("timeout", "30s")
    viper.SetDefault("output_format", "text")
    
    // 配置文件搜索
    viper.SetConfigName("oho")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("$HOME/.config/oho/")
    viper.AddConfigPath(".")
    
    // 环境变量
    viper.SetEnvPrefix("OHO")
    viper.AutomaticEnv()
    
    // 命令行标志绑定
    viper.BindPFlag("server_url", rootCmd.PersistentFlags().Lookup("server"))
    
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

**配置优先级链**（从高到低）：
1. 命令行标志（`--server http://example.com`）
2. 环境变量（`OHO_SERVER_URL=http://example.com`）
3. 配置文件（`~/.config/oho/config.yaml`）
4. 默认值

### 7. 错误处理和输出格式化

**分层错误处理**：

```go
// 自定义错误类型
type APIError struct {
    StatusCode int
    Message    string
    Details    map[string]interface{}
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// 错误包装和分类
func classifyError(err error) error {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        switch apiErr.StatusCode {
        case 401:
            return fmt.Errorf("authentication failed: %w", err)
        case 404:
            return fmt.Errorf("resource not found: %w", err)
        default:
            return fmt.Errorf("API request failed: %w", err)
        }
    }
    return err
}
```

**输出格式化**：

```go
// internal/output/printer.go
type Printer interface {
    PrintSession(*Session) error
    PrintMessage(*Message) error
    PrintError(error) error
}

type TextPrinter struct{ Out io.Writer }
type JSONPrinter struct{ Out io.Writer }
type TablePrinter struct{ Out io.Writer }

// 根据--output标志选择打印机
func NewPrinter(format string, w io.Writer) Printer {
    switch format {
    case "json":
        return &JSONPrinter{Out: w}
    case "table":
        return &TablePrinter{Out: w}
    default:
        return &TextPrinter{Out: w}
    }
}
```

### 8. 认证管理

**认证流程**：
1. **API Token 认证**：支持配置文件、环境变量、命令行标志
2. **Token 刷新**：自动刷新过期令牌（如支持）
3. **凭据存储**：安全存储认证信息（使用系统密钥链或加密文件）

**实现模式**：
```go
// internal/auth/token.go
type TokenManager struct {
    token    string
    expiry   time.Time
    refresh  func() (string, time.Time, error)
}

func (tm *TokenManager) GetToken() (string, error) {
    if tm.token == "" || time.Now().After(tm.expiry) {
        token, expiry, err := tm.refresh()
        if err != nil {
            return "", err
        }
        tm.token = token
        tm.expiry = expiry
    }
    return tm.token, nil
}
```

### 9. 测试策略

**生产级测试模式**：
- **单元测试**：每个包独立的测试
- **集成测试**：使用 `httptest.Server` 模拟 API
- **E2E 测试**：实际的 CLI 命令执行测试
- **Golden 测试**：输出格式的回归测试

**测试工具**：
- `testify/assert`：断言库
- `httptest`：HTTP测试服务器
- `io/ioutil`：捕获输出
- `exec.Command`：执行 CLI 命令

### 10. 推荐的依赖

```toml
# go.mod
github.com/spf13/cobra v1.8.0      # CLI框架
github.com/spf13/viper v1.18.0     # 配置管理
github.com/stretchr/testify v1.9.0 # 测试
golang.org/x/term v0.18.0          # 终端处理
github.com/fatih/color v1.16.0     # 彩色输出（可选）
github.com/olekukonko/tablewriter v0.0.5  # 表格输出（可选）
```

## 实施优先级

### 第一阶段：基础框架
1. 初始化项目结构
2. 设置 Cobra 根命令
3. 实现配置管理（Viper）
4. 创建基础的 REST 客户端

### 第二阶段：核心功能
1. 实现 session 子命令（create、list、delete）
2. 实现 message 子命令（add、list）
3. 完善错误处理和输出格式化

### 第三阶段：高级功能
1. 添加认证管理
2. 实现请求重试和超时控制
3. 添加测试套件
4. 完善文档和帮助信息

## 成功标准

1. **功能完整性**：支持所有指定的子命令和标志
2. **错误处理**：清晰的错误消息和适当的退出码
3. **配置管理**：多源配置支持（文件、环境变量、命令行）
4. **代码质量**：遵循 Go 最佳实践，完善的测试覆盖
5. **用户体验**：直观的帮助信息，清晰的输出格式

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| API 变更 | 高 | 使用接口抽象，提供适配器模式 |
| 认证复杂度 | 中 | 模块化认证管理，支持多种认证方式 |
| 性能问题 | 低 | 实现连接池，优化 HTTP 客户端配置 |
| 依赖过多 | 低 | 严格控制依赖，使用最小化依赖集 |

## 验收标准

1. ✅ `oho --help` 显示完整的帮助信息
2. ✅ `oho session create mysession` 成功创建会话
3. ✅ `oho message add -s sessionid "message content"` 成功添加消息
4. ✅ 配置文件和环境变量正常工作
5. ✅ 错误情况有清晰的错误消息和适当的退出码
6. ✅ 支持 JSON 和表格输出格式

---
*设计文档版本：1.0*
*创建日期：2026-02-27*
*最后更新：2026-02-27*