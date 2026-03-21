# oho (OpenCode CLI) 项目完整摘要文档

## 1. 项目概述

### 1.1 项目目标和定位

**oho** 是 OpenCode Server 的命令行客户端工具，为 OpenCode Server API 提供完整的命令行访问能力。

**核心定位**：
- 🤖 **AI 可调用**：原生支持被任何 AI Agent 调用
- 🔄 **工作流集成**：可集成到自动化工作流中
- 📦 **CI/CD 支持**：可在 CI/CD 管道中运行
- 🔗 **工具链兼容**：可与 shell 工具无缝组合

**独特价值**：oho 是 OpenCode 生态系统中**唯一完全使用 Go 实现的命令行客户端**，专注于提供 Linux/Unix 环境下的高效命令行体验。

### 1.2 核心功能列表

| 功能类别 | 具体功能 |
|----------|----------|
| **会话管理** | 创建/删除/列出/更新会话，会话 fork/revert，会话状态查询 |
| **消息管理** | 发送消息/命令/shell 命令，异步消息，消息历史查询 |
| **项目管理** | 项目列表，当前项目，VCS 信息，路径管理 |
| **配置管理** | 配置获取/设置，提供商配置，默认模型设置 |
| **文件操作** | 文件列表，文件内容读取，文件状态查询 |
| **查找功能** | 文本搜索，文件查找，符号查找 |
| **工具管理** | 工具 ID 列表，工具能力查询，MCP 服务器管理 |
| **状态查询** | LSP 状态，格式化器状态，全局健康检查 |
| **TUI 控制** | Toast 通知，帮助界面，命令执行 |
| **认证管理** | OAuth 授权，凭证设置 |

### 1.3 技术栈说明

| 技术组件 | 版本/说明 |
|----------|-----------|
| **编程语言** | Go 1.21+ |
| **CLI 框架** | Cobra v1.8.0 |
| **标志处理** | pflag v1.0.5 |
| **HTTP 客户端** | 标准库 net/http |
| **JSON 处理** | 标准库 encoding/json |
| **认证方式** | HTTP Basic Auth |
| **构建工具** | Make / go build |
| **测试框架** | Go testing 包 |

---

## 2. 架构设计

### 2.1 目录结构说明

```
oho/
├── cmd/                          # 命令行子命令目录
│   ├── main.go                   # 程序入口，注册所有子命令
│   ├── add/                      # oho add 命令（快捷创建会话 + 发消息）
│   ├── agent/                    # agent 相关命令
│   ├── auth/                     # 认证相关命令
│   ├── command/                  # 命令管理
│   ├── configcmd/                # 配置管理
│   ├── file/                     # 文件操作
│   ├── find/                     # 查找功能
│   ├── formatter/                # 格式化器状态
│   ├── global/                   # 全局命令（健康检查、事件流）
│   ├── lsp/                      # LSP 状态
│   ├── mcp/                      # MCP 服务器管理
│   ├── mcpserver/                # MCP 服务器启动
│   ├── message/                  # 消息管理
│   ├── project/                  # 项目管理
│   ├── provider/                 # 提供商管理
│   ├── session/                  # 会话管理
│   ├── tool/                     # 工具命令
│   └── tui/                      # TUI 控制
├── internal/                     # 内部包（不对外暴露）
│   ├── client/                   # HTTP 客户端封装
│   │   ├── client.go             # 客户端实现
│   │   ├── client_interface.go   # 客户端接口定义
│   │   ├── client_mock.go        # Mock 客户端（测试用）
│   │   └── client_test.go        # 客户端测试
│   ├── config/                   # 配置管理
│   │   ├── config.go             # 配置加载和优先级
│   │   └── config_test.go        # 配置测试
│   ├── types/                    # 数据类型定义
│   │   ├── types.go              # 所有 API 类型定义
│   │   └── types_test.go         # 类型测试
│   ├── util/                     # 工具函数
│   │   ├── output.go             # 输出格式化工具
│   │   └── output_test.go        # 输出测试
│   └── testutil/                 # 测试工具
│       └── testutil.go           # Mock 数据生成
├── bin/                          # 编译输出目录
├── tasks/                        # 任务定义（可选）
├── Makefile                      # 构建配置
├── go.mod                        # Go 模块定义
├── coverage.out                  # 测试覆盖率报告
├── README.md                     # 英文文档
└── README_zh.md                  # 中文文档
```

### 2.2 模块划分和职责

| 模块 | 职责 | 关键文件 |
|------|------|----------|
| **cmd/** | 命令行接口层，解析用户输入并调用内部 API | `main.go`, `*/<command>.go` |
| **internal/client/** | HTTP 通信层，封装所有 API 请求 | `client.go` |
| **internal/config/** | 配置管理，处理优先级（命令行 > 配置文件 > 环境变量） | `config.go` |
| **internal/types/** | 数据类型定义，与 API 响应结构对应 | `types.go` |
| **internal/util/** | 工具函数，如输出格式化 | `output.go` |
| **internal/testutil/** | 测试辅助，提供 Mock 数据 | `testutil.go` |

### 2.3 数据流和调用关系

```
用户输入 (oho add "消息")
    │
    ▼
┌─────────────────────────────────────┐
│  cmd/main.go                        │
│  - 初始化配置 (config.Init())       │
│  - 注册子命令                        │
│  - 绑定全局标志                      │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  cmd/add/add.go                     │
│  - 解析 add 命令参数                 │
│  - 调用 createSession()             │
│  - 调用 sendMessage()               │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  internal/client/client.go          │
│  - NewClient() 创建 HTTP 客户端      │
│  - Post/PostWithQuery 发送请求      │
│  - 处理认证和错误                    │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  OpenCode Server API                │
│  - POST /session                    │
│  - POST /session/{id}/message       │
└─────────────────────────────────────┘
    │
    ▼
JSON 响应 → types 解析 → 输出结果
```

---

## 3. 核心实现分析

### 3.1 `oho add` 命令实现流程

`oho add` 是一个复合命令，将"创建会话"和"发送消息"两个操作合并为一步。

#### 完整流程（5 步）

```
Step 1: 获取工作目录
   │
   ▼
Step 2: 生成会话标题（如未提供）
   │
   ▼
Step 3: 调用 createSession() → POST /session?directory={dir}
   │
   ▼
Step 4: 调用 sendMessage() → POST /session/{id}/message
   │
   ▼
Step 5: 输出结果（JSON 或文本格式）
```

#### 关键代码位置

| 函数 | 文件 | 行号 | 功能 |
|------|------|------|------|
| `runAdd()` | `cmd/add/add.go` | 69-136 | add 命令主执行逻辑 |
| `createSession()` | `cmd/add/add.go` | 140-162 | 创建会话，返回 session ID |
| `sendMessage()` | `cmd/add/add.go` | 187-251 | 发送消息，支持文件附件 |
| `convertModel()` | `cmd/add/add.go` | 165-183 | 模型格式转换（provider:model） |
| `detectMimeType()` | `cmd/add/add.go` | 254-317 | MIME 类型检测（文件附件） |

#### API 调用详情

```go
// 1. 创建会话
POST /session?directory={directory}
Body: {"title": "会话标题", "parentID": "父会话 ID(可选)"}
Response: {"id": "ses_xxx", ...}

// 2. 发送消息
POST /session/{sessionID}/message
Body: {
  "model": "provider:model" 或 "model-string",
  "agent": "agent-id",
  "noReply": false,
  "system": "system-prompt",
  "tools": ["tool1", "tool2"],
  "parts": [
    {"type": "text", "text": "消息内容"},
    {"type": "file", "url": "data:...;base64,xxx", "mime": "image/png"}
  ]
}
Response: {"info": {"id": "msg_xxx", ...}, "parts": [...]}
```

### 3.2 HTTP 客户端通信层设计

#### Client 结构

```go
type Client struct {
    baseURL    string         // API 基础 URL (http://host:port)
    httpClient *http.Client   // HTTP 客户端（带超时配置）
    username   string         // Basic Auth 用户名
    password   string         // Basic Auth 密码
}
```

#### 核心方法

| 方法 | 签名 | 功能 |
|------|------|------|
| `NewClient()` | `func NewClient() *Client` | 创建客户端，加载配置，设置 5 分钟超时 |
| `Request()` | `func (c *Client) Request(ctx, method, path, body)` | 发送 HTTP 请求，处理认证和错误 |
| `RequestWithQuery()` | `func (c *Client) RequestWithQuery(ctx, method, path, queryParams, body)` | 发送带 query 参数的请求 |
| `Get/Post/Put/Delete()` | 各种 HTTP 方法封装 | 便捷方法 |
| `SSEStream()` | `func (c *Client) SSEStream(ctx, path)` | 服务器发送事件流（SSE） |

#### 认证处理

```go
// 在 Request() 方法中（第 67-69 行）
if c.username != "" && c.password != "" {
    req.SetBasicAuth(c.username, c.password)
}
```

#### 错误处理

```go
// 状态码检查（第 91-96 行）
if resp.StatusCode >= 400 {
    if resp.StatusCode == 401 {
        return nil, fmt.Errorf("认证失败 [401]: 用户名或密码错误\n\n请配置认证信息...")
    }
    return nil, fmt.Errorf("API 错误 [%d]: %s", resp.StatusCode, string(respBody))
}
```

### 3.3 错误处理和超时机制

#### 超时配置

| 配置项 | 默认值 | 环境变量 | 位置 |
|--------|--------|----------|------|
| HTTP 客户端超时 | 300 秒 (5 分钟) | `OPENCODE_CLIENT_TIMEOUT` | `client.go:32-37` |

```go
// client.go:32-37
timeoutSec := 300 // 5 分钟
if envTimeout := os.Getenv("OPENCODE_CLIENT_TIMEOUT"); envTimeout != "" {
    if parsed, err := strconv.Atoi(envTimeout); err == nil && parsed > 0 {
        timeoutSec = parsed
    }
}
```

#### 错误处理模式

1. **错误包装**：使用 `%w` 包装底层错误
   ```go
   return fmt.Errorf("failed to create session: %w", err)
   ```

2. **错误传播**：错误逐层向上传递
   ```go
   if err != nil {
       return fmt.Errorf("API request failed: %w", err)
   }
   ```

3. **特殊错误处理**：401 认证错误提供详细配置指南
   ```go
   if resp.StatusCode == 401 {
       return nil, fmt.Errorf("认证失败 [401]: ...\n\n请配置认证信息...")
   }
   ```

#### 部分失败处理

`oho add` 命令实现了部分失败处理：
- 会话创建成功但消息发送失败时，返回 session ID 并提示警告
- 用户可以选择继续使用已创建的会话

```go
// add.go:99-112
if err != nil {
    if addJSONOutput {
        output := map[string]interface{}{
            "sessionId": sessionID,
            "status":    "partial",
            "error":     fmt.Sprintf("failed to send message: %v", err),
        }
        // ...
    } else {
        fmt.Printf("Session created: %s\n", sessionID)
        fmt.Printf("Warning: Message send failed: %v\n", err)
    }
    return nil  // 不返回错误，让用户知道会话已创建
}
```

---

## 4. 测试覆盖

### 4.1 现有测试结构

| 测试目录 | 测试文件数 | 测试内容 |
|----------|------------|----------|
| `cmd/*/` | 15 个 `*_test.go` | 各子命令的单元测试 |
| `internal/client/` | 1 个 `client_test.go` | HTTP 客户端测试 |
| `internal/config/` | 1 个 `config_test.go` | 配置管理测试 |
| `internal/types/` | 1 个 `types_test.go` | 类型定义测试 |
| `internal/util/` | 2 个 `*_test.go` | 工具函数测试 |

**总计**：21 个测试文件

### 4.2 新增测试用例说明

#### `cmd/add/add_test.go` 测试用例（769 行）

| 测试函数 | 测试内容 | 覆盖场景 |
|----------|----------|----------|
| `TestConvertModel` | 模型格式转换 | 空模型、简单字符串、provider:model 格式 |
| `TestDetectMimeType` | MIME 类型检测 | 18 种文件扩展名 |
| `TestCreateSession` | 会话创建 | 成功/失败/API 错误/JSON 解析错误 |
| `TestSendMessage` | 消息发送 | 简单消息、带附件、no-reply 模式、错误处理 |
| `TestRunAddSuccess` | add 命令完整流程 | 基本使用、自定义标题、no-reply、指定 agent/model |
| `TestRaceConditionScenarios` | 竞态条件 | 会话就绪延迟场景 |
| `TestTimeoutScenarios` | 超时场景 | 请求在超时内完成/超时 |
| `TestErrorPropagation` | 错误传播 | 会话创建失败、消息发送失败 |
| `TestPartialFailureHandling` | 部分失败处理 | 会话成功但消息失败 |
| `TestJSONOutputFormat` | 输出格式 | JSON 输出/文本输出 |

#### `internal/client/client_test.go` 测试用例（122 行）

| 测试函数 | 测试内容 |
|----------|----------|
| `TestNewClient` | 客户端创建 |
| `TestClientGetSuccess` | GET 请求成功 |
| `TestClientPostSuccess` | POST 请求成功 |
| `TestClientEmptyResponse` | 空响应处理 |
| `TestClientErrorResponse` | 错误响应处理 |

### 4.3 测试覆盖率分析

根据 `coverage.out` 报告：

| 模块 | 覆盖状态 | 说明 |
|------|----------|------|
| `cmd/add/` | ✅ 高覆盖 | 核心功能有完整单元测试 |
| `cmd/message/` | ✅ 已覆盖 | 消息命令测试 |
| `cmd/session/` | ✅ 已覆盖 | 会话命令测试 |
| `internal/client/` | ✅ 高覆盖 | HTTP 客户端核心逻辑已测试 |
| `internal/config/` | ✅ 已覆盖 | 配置加载测试 |
| `internal/types/` | ✅ 已覆盖 | 类型定义测试 |
| `internal/util/` | ✅ 已覆盖 | 工具函数测试 |

**测试特点**：
- 使用 Mock 客户端 (`client.MockClient`) 隔离外部依赖
- 使用 `httptest.NewServer` 模拟 HTTP 服务器
- 覆盖正常流程和错误流程
- 包含边界条件测试（超时、空响应、JSON 解析错误）

---

## 5. 已知问题和修复

### 5.1 间歇性失败原因分析

基于测试代码中的场景分析：

| 问题类型 | 原因 | 解决方案 |
|----------|------|----------|
| **会话就绪延迟** | 创建会话后服务器需要时间初始化 | 测试中模拟了 50-200ms 延迟场景 |
| **请求超时** | AI 响应时间超过默认 5 分钟 | 可通过 `OPENCODE_CLIENT_TIMEOUT` 调整 |
| **认证失败** | 密码变更或配置未同步 | 错误信息提供详细配置指南 |
| **部分失败** | 会话创建成功但消息发送失败 | 实现部分失败处理，返回 session ID |

### 5.2 已修复的 Bug 列表

根据版本信息 `v1.1.0-17-g512e296-dirty`：

| Bug | 描述 | 修复状态 |
|-----|------|----------|
| 模型格式兼容 | `provider:model` 格式解析 | ✅ 已修复 (`convertModel()`) |
| 文件附件 MIME 检测 | 文件扩展名大小写处理 | ✅ 已修复 (`detectMimeType()` 使用 `ToLower()`) |
| 部分失败处理 | 会话创建成功但消息失败时返回空 | ✅ 已修复 (返回 session ID + 警告) |
| JSON 输出格式 | 部分失败时 JSON 格式不完整 | ✅ 已修复 (添加 `status: partial`) |

### 5.3 待改进项

| 改进项 | 优先级 | 说明 |
|--------|--------|------|
| **重试机制** | 中 | 当前无自动重试逻辑，可添加指数退避重试 |
| **连接池** | 低 | 当前每次命令创建新客户端，可复用连接 |
| **日志系统** | 中 | 可添加 debug 日志输出，便于排查问题 |
| **进度指示** | 低 | 长时间等待 AI 响应时可显示进度 |
| **命令补全** | 低 | 可生成 shell 补全脚本（已在 README 中提及） |

---

## 6. 使用指南

### 6.1 安装步骤

#### 方法一：快速安装（推荐）

```bash
# 使用 curl
curl -sSL https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.sh | bash

# 或使用 wget
wget -qO- https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.sh | bash
```

#### 方法二：源码编译

```bash
# 克隆仓库
git clone https://github.com/tornado404/opencode_cli.git
cd opencode_cli/oho

# 构建
make build

# 或构建 Linux 版本
make build-linux
```

#### 依赖要求

- Go 1.21+
- Cobra CLI 框架（通过 `go mod` 自动安装）

### 6.2 常用命令示例

#### 配置连接

```bash
# 使用环境变量
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# 或使用命令行标志
oho --host 127.0.0.1 --port 4096 --password your-password <command>
```

#### 会话管理

```bash
# 列出所有会话
oho session list

# 创建新会话
oho session create

# 在指定目录创建会话
oho session create --path /your/project

# 创建子会话
oho session create --parent ses_xxx

# 获取会话详情
oho session get ses_xxx

# 删除会话
oho session delete ses_xxx

# Fork 会话
oho session fork ses_xxx

# Revert 到指定消息
oho session revert ses_xxx --message msg_yyy
```

#### 消息管理

```bash
# 发送消息
oho message add -s ses_xxx "帮我分析这个项目"

# 使用指定 agent 和模型
oho message add -s ses_xxx "代码审查" --agent default --model anthropic:claude-3

# 不等待 AI 响应（异步）
oho message add -s ses_xxx "后台运行测试" --no-reply

# 附加文件
oho message add -s ses_xxx "分析这个日志" --file /var/log/app.log

# 执行 shell 命令
oho message shell -s ses_xxx --agent default "ls -la"
```

#### 快捷命令 (oho add)

```bash
# 创建会话并发送消息（一步完成）
oho add "帮我分析这个项目"

# 自定义会话标题
oho add "修复登录 bug" --title "Bug 修复"

# 指定工作目录
oho add "任务描述" --directory /path/to/project

# JSON 格式输出
oho add "消息内容" --json
```

#### 其他常用命令

```bash
# 健康检查
oho global health

# 获取配置
oho config get

# 列出提供商
oho provider list

# 文件列表
oho file list

# 文本搜索
oho find text "pattern"

# 查找符号
oho find symbol "FunctionName"

# LSP 状态
oho lsp status

# MCP 服务器列表
oho mcp list
```

### 6.3 故障排查方法

#### 问题 1：认证失败 (401)

```
错误：认证失败 [401]: 用户名或密码错误

解决方案：
1. 检查密码是否正确
2. 确认服务器是否启用 Basic Auth
3. 尝试使用环境变量：
   export OPENCODE_SERVER_PASSWORD=your-password
```

#### 问题 2：连接超时

```
错误：请求失败：context deadline exceeded

解决方案：
1. 确认服务器地址和端口正确
2. 检查服务器是否运行
3. 增加超时时间：
   export OPENCODE_CLIENT_TIMEOUT=600  # 10 分钟
```

#### 问题 3：会话创建失败

```
错误：API 请求失败：...

解决方案：
1. 检查服务器健康状态：oho global health
2. 确认目录路径有效
3. 查看详细错误信息（使用 --json 输出）
```

#### 问题 4：消息发送失败但会话已创建

```
警告：消息发送失败：...
会话已创建：ses_xxx

解决方案：
1. 会话仍然可用，可继续使用 oho message add -s ses_xxx 发送消息
2. 检查消息内容是否包含非法字符
3. 确认附件文件存在且可读
```

#### 调试技巧

```bash
# 使用 JSON 输出查看详细响应
oho session list --json

# 检查服务器配置
oho config get

# 查看当前项目信息
oho project current

# 检查 VCS 状态
oho vcs
```

---

## 附录

### A. 配置文件位置

```
~/.config/oho/config.json
```

配置文件格式：
```json
{
  "host": "127.0.0.1",
  "port": 4096,
  "username": "opencode",
  "password": "",
  "json": false
}
```

### B. 环境变量优先级

配置优先级：**命令行标志 > 配置文件 > 环境变量 > 默认值**

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `OPENCODE_SERVER_HOST` | `127.0.0.1` | 服务器主机地址 |
| `OPENCODE_SERVER_PORT` | `4096` | 服务器端口 |
| `OPENCODE_SERVER_USERNAME` | `opencode` | 用户名 |
| `OPENCODE_SERVER_PASSWORD` | 空 | 密码 |
| `OPENCODE_CLIENT_TIMEOUT` | `300` | HTTP 超时（秒） |

### C. 项目信息

- **版本**：v1.1.0-17-g512e296-dirty
- **提交**：512e296
- **构建时间**：2026-03-21
- **仓库**：https://github.com/tornado404/opencode_cli
- **许可证**：GPL v3

---

*文档生成时间：2026-03-21*
