# oho - OpenCode 命令行工具

> 唯一完全基于 Go 实现的 OpenCode 命令行客户端。

[![GitHub Stars](https://img.shields.io/github/stars/tornado404/opencode_cli?style=flat-square)](https://github.com/tornado404/opencode_cli/stargazers)
[![License](https://img.shields.io/badge/license-GPLv3-blue?style=flat-square)](LICENSE)

oho 是 OpenCode Server 的命令行客户端工具，提供对 OpenCode Server API 的完整访问能力。

## 项目定位

### 独特价值

**oho** 是 [OpenCode 生态系统](https://opencode.ai/docs/ecosystem/) 中**唯一完全使用 Go 实现的命令行客户端**。

> "oho is callable from Bash" 代表了强大的扩展性和兼容性——这是项目的独特定位。

### 设计目标

让 OpenCode 更易于被其他 AI 调用和监控：

- 🤖 原生支持被任何 AI Agent 调用
- 🔄 可集成到自动化工作流中
- 📦 可在 CI/CD 管道中运行
- 🔗 可与 shell 工具无缝组合

### 独特的 Linux 能力

在 Linux 环境下，oho 提供 OpenCode CLI 目前不支持的能力：

| 能力 | 说明 |
|------|------|
| 📁 在指定目录创建会话 | 在任意目录启动 AI 编程会话 |
| 💬 基于会话继续发送消息 | 恢复之前的会话上下文 |
| 🗑️ 销毁会话 | 完整的会话生命周期管理 |
| 🔄 会话 Fork 和 Revert | 轻松切换实验性开发 |

## 功能特性

- ✅ 完整的 API 映射和封装
- ✅ HTTP Basic Auth 认证支持
- ✅ JSON/Text 双输出模式
- ✅ 配置文件和环境变量支持
- ✅ 完整的会话管理操作
- ✅ 消息发送和管理
- ✅ 文件和符号查找
- ✅ TUI 界面控制
- ✅ MCP/LSP/Formatter 状态管理
- 📊 **[API 覆盖率](./COMPLETION.md)** - 查看实现覆盖情况

## 安装

### 快速安装（推荐）

```bash
# 使用 curl
curl -sSL https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.sh | bash

# 或使用 wget
wget -qO- https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.sh | bash
```

### 源码构建

```bash
# 克隆仓库
git clone https://github.com/tornado404/opencode_cli.git
cd opencode_cli/oho

# 构建
make build

# 或构建 Linux 版本
make build-linux
```

### 依赖要求

- Go 1.21+
- Cobra CLI 框架
- 标准库 net/http

## 快速开始

### 1. 配置服务器连接

```bash
# 使用环境变量
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# 或使用命令行参数
oho --host 127.0.0.1 --port 4096 --password your-password session list
```

### 2. 基本用法

```bash
# 检查服务器健康状态
oho global health

# 列出所有会话
oho session list

# 创建新会话
oho session create

# 在指定目录创建会话
oho session create --path /your/project

# 发送消息
oho message add -s <session-id> "你好，请帮我分析这个项目"

# 继续现有会话
oho message add -s <session-id> "继续之前的任务"

# 查看消息列表
oho message list -s <session-id>

# 销毁会话
oho session delete <session-id>

# 获取配置
oho config get

# 列出提供商
oho provider list
```

## 与其他生态项目对比

| 特性 | oho | 其他生态项目 |
|------|-----|-------------|
| 实现语言 | Go | TypeScript/Python/Go |
| AI 可调用性 | ✅ 原生支持 | 需要额外适配 |
| 跨平台 | Linux/Mac/Windows | 运行时依赖 |
| 集成难度 | ⭐⭐⭐⭐⭐ 极低 | ⭐⭐⭐ 中等 |

参考：[OpenCode 生态系统中的其他项目](https://opencode.ai/docs/ecosystem/)

## 命令参考

### 全局命令

```bash
oho global health          # 检查服务器健康状态
oho global event           # 监听全局事件流 (SSE)
```

### 项目管理

```bash
oho project list           # 列出所有项目
oho project current        # 获取当前项目
oho path                   # 获取当前路径
oho vcs                    # 获取 VCS 信息
oho instance dispose       # 销毁当前实例
```

### 会话管理

```bash
oho session list                         # 列出所有会话
oho session list --id ses_abc            # 按 ID 过滤（模糊匹配）
oho session list --title "测试"           # 按标题过滤（模糊）
oho session list --project-id proj1      # 按项目 ID 过滤
oho session list --directory babylon     # 按目录过滤
oho session list --created 1773537883643 # 按创建时间过滤
oho session list --updated 1773538142930  # 按更新时间过滤
oho session list --sort updated --order desc  # 按更新时间排序（倒序）
oho session list --limit 10 --offset 0   # 分页
oho session create                       # 创建新会话
oho session create --title "名称"         # 创建带自定义标题的会话
oho session create --parent ses_xxx      # 创建子会话
oho session create --path /path          # 在指定目录创建会话
oho session status                       # 获取所有会话状态
oho session get <id>                     # 获取会话详情
oho session delete <id>                  # 删除会话
oho session update <id> --title "新标题"  # 更新会话
oho session children <id>                # 获取子会话
oho session todo <id>                    # 获取待办事项
oho session init <id> --provider <provider> --model <model>  # 初始化 AGENTS.md
oho session fork <id>                    # 分叉会话
oho session abort <id>                   # 中止会话
oho session share <id>                   # 分享会话
oho session unshare <id>                # 取消分享会话
oho session diff <id>                    # 获取文件差异
oho session diff <id> --message <msg-id> # 获取指定消息后的差异
oho session summarize <id> --provider <provider> --model <model>  # 总结会话
oho session revert <id> --message <msg-id>  # 回退消息
oho session unrevert <id>                # 恢复所有已回退的消息
oho session permissions <id> <perm-id> --response allow  # 响应权限请求
oho session submit "任务描述"             # 一键提交任务（创建会话+发送消息）
oho session submit "任务" --init-project --provider openai --model gpt-4  # 提交任务并初始化项目
oho session achieve <id>                 # 归档会话
```

**List 命令参数**:

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--id` | string | 按 ID 过滤（模糊查询，大小写不敏感） | - |
| `--title` | string | 按标题过滤（模糊查询，大小写不敏感） | - |
| `--created` | int64 | 按创建时间过滤（精确时间戳） | - |
| `--updated` | int64 | 按更新时间过滤（精确时间戳） | - |
| `--project-id` | string | 按项目 ID 过滤（模糊查询） | - |
| `--directory` | string | 按目录过滤（模糊查询） | - |
| `--status` | string | 按状态过滤（running/completed/error/aborted/idle） | - |
| `--running` | bool | 只显示正在运行的会话 | false |
| `--sort` | string | 排序字段（created/updated） | updated |
| `--order` | string | 排序顺序（asc/desc） | desc |
| `--limit` | int | 限制结果数量 | - |
| `--offset` | int | 分页偏移量 | 0 |
| `-j, --json` | bool | JSON 输出格式 | false |

**`session submit` 命令参数**:

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--init-project` | bool | 初始化项目（需配合 --provider 和 --model） | false |
| `--provider` | string | Provider ID（初始化项目时必需） | - |
| `--model` | string | Model ID（初始化项目时必需） | - |
| `--title` | string | 会话标题 | 自动生成 |
| `--directory` | string | 工作目录 | 当前目录 |
| `--agent` | string | Agent ID | - |
| `--message-model` | string | 消息使用的模型 | - |
| `--no-reply` | bool | 不等待 AI 响应 | false |
| `--system` | string | System prompt | - |
| `--tools` | string[] | 工具列表（可多次指定） | - |
| `--file` | string[] | 文件附件（可多次指定） | - |

### 消息管理

```bash
oho message list -s <session>              # 列出消息
oho message add -s <session> "内容"        # 发送消息
oho message get -s <session> <msg-id>     # 获取消息详情
oho message prompt-async -s <session> "内容"  # 异步发送消息
oho message command -s <session> /help     # 执行斜杠命令
oho message shell -s <session> --agent default "ls -la"  # 运行 shell 命令
```

### 快捷命令 (oho add)

```bash
oho add "帮我分析这个项目"                     # 创建会话并发送消息
oho add "修复登录 bug" --title "Bug 修复"       # 创建带自定义标题的会话
oho add "测试功能" --no-reply --agent default  # 不等待 AI 响应
oho add "分析日志" --file /var/log/app.log    # 附加文件到消息
oho add "任务描述" --directory /path/to/project # 指定工作目录
oho add "消息内容" --json                      # JSON 格式输出
```

### ⚠️ 超时注意事项

`oho add` 命令默认会等待 AI 响应后返回。对于复杂任务，AI 可能需要较长时间思考，可能导致超时。

**避免超时的方法**:

1. **使用 `--no-reply` 参数**（推荐）:
   ```bash
   # 发送消息后立即返回，不等待 AI 响应
   oho add "分析项目结构" --no-reply
   
   # 稍后检查结果
   oho message list -s <session-id>
   ```

2. **增加超时时间**:
   ```bash
   # 设置超时为 10 分钟（600 秒）
   export OPENCODE_CLIENT_TIMEOUT=600
   oho add "复杂任务"
   
   # 或临时设置
   OPENCODE_CLIENT_TIMEOUT=600 oho add "复杂任务"
   ```

3. **使用异步命令**:
   ```bash
   # 先创建会话
   oho session create --title "任务"
   
   # 异步发送消息
   oho message prompt-async -s <session-id> "任务描述"
   ```

**超时配置**:

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `OPENCODE_CLIENT_TIMEOUT` | 300 秒 | HTTP 请求超时时间（秒） |

**`oho add` 参数**:

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--title` | string | 会话标题（不提供则自动生成） | 自动生成 |
| `--parent` | string | 父会话 ID（用于创建子会话） | - |
| `--directory` | string | 会话的工作目录 | 当前目录 |
| `--agent` | string | Agent ID | - |
| `--model` | string | 模型 ID（如 `provider:model`） | 默认模型 |
| `--no-reply` | bool | 不等待 AI 响应 | false |
| `--system` | string | System prompt | - |
| `--tools` | string[] | 工具列表（可多次指定） | - |
| `--file` | string[] | 文件附件（可多次指定） | - |
| `-j, --json` | bool | JSON 格式输出 | false |

### 配置管理

```bash
oho config get                      # 获取配置
oho config set --theme dark         # 更新配置
oho config providers                # 列出提供商和默认模型
```

### 提供商管理

```bash
oho provider list                   # 列出所有提供商
oho provider auth                  # 获取提供商认证方式
oho provider oauth authorize <id>  # OAuth 授权
oho provider oauth callback <id>   # 处理 OAuth 回调
```

### 文件操作

```bash
oho file list [path]                # 列出文件
oho file content <path>             # 读取文件内容
oho file status                     # 获取文件状态
```

### 查找功能

```bash
oho find text "pattern"             # 搜索文本
oho find file "query"              # 查找文件
oho find symbol "query"            # 查找符号
```

### 其他命令

```bash
oho agent list                      # 列出 agents
oho command list                    # 列出命令
oho tool ids                        # 列出工具 ID
oho tool list --provider xxx --model xxx  # 列出工具
oho lsp status                      # LSP 状态
oho formatter status                # Formatter 状态
oho mcp list                        # 列出 MCP 服务器
oho mcp add <name> --config '{}'    # 添加 MCP 服务器
oho tui open-help                   # 打开帮助
oho tui show-toast --message "消息" # 显示提示
oho auth set <provider> --credentials key=value  # 设置认证
```

## 输出格式

使用 `-j` 或 `--json` 参数获取 JSON 输出：

```bash
oho session list -j
oho config get --json
```

## 配置文件

配置文件位于 `~/.config/oho/config.json`：

```json
{
  "host": "127.0.0.1",
  "port": 4096,
  "username": "opencode",
  "password": "",
  "json": false
}
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `OPENCODE_SERVER_HOST` | 服务器主机 | `127.0.0.1` |
| `OPENCODE_SERVER_PORT` | 服务器端口 | `4096` |
| `OPENCODE_SERVER_USERNAME` | 用户名 | `opencode` |
| `OPENCODE_SERVER_PASSWORD` | 密码 | 空 |

## 开发

```bash
# 运行
go run ./cmd/oho --help

# 测试
go test ./...

# 格式化
go fmt ./...

# 清理
make clean
```

## 项目结构

```
oho/
├── cmd/
│   └── oho/
│       ├── main.go               # 程序入口
│       ├── root.go               # 根命令
│       ├── cmd/                  # 子命令
│       │   ├── global/           # 全局命令
│       │   ├── project/          # 项目管理
│       │   ├── session/          # 会话管理
│       │   ├── message/          # 消息管理
│       │   ├── configcmd/        # 配置管理
│       │   ├── provider/         # 提供商管理
│       │   ├── file/             # 文件操作
│       │   ├── find/             # 查找功能
│       │   ├── tool/             # 工具命令
│       │   ├── agent/            # Agent 命令
│       │   ├── command/          # 命令管理
│       │   ├── lsp/              # LSP 状态
│       │   ├── formatter/        # Formatter 状态
│       │   ├── mcp/              # MCP 服务器管理
│       │   ├── mcpserver/        # MCP 服务器启动
│       │   ├── tui/              # TUI 控制
│       │   └── auth/             # 认证管理
│       └── internal/
│           ├── client/           # HTTP 客户端
│           ├── config/           # 配置管理
│           ├── types/            # 类型定义
│           └── util/             # 工具函数
├── Makefile
└── README.md
```

## MCP 服务器

oho 可以作为 MCP (Model Context Protocol) 服务器使用，允许外部 MCP 客户端（如 Claude Desktop、Cursor 等）通过 MCP 协议调用 OpenCode API。

### 启动 MCP 服务器

```bash
# 启动 MCP 服务器（stdio 模式）
oho mcpserver
```

MCP 服务器使用 stdio 传输，这是本地 MCP 客户端的标准模式。

### 可用的 MCP 工具

| 工具 | 说明 |
|------|------|
| `session_list` | 列出所有会话 |
| `session_create` | 创建新会话 |
| `session_get` | 获取会话详情 |
| `session_delete` | 删除会话 |
| `session_status` | 获取所有会话状态 |
| `message_list` | 列出会话中的消息 |
| `message_add` | 向会话发送消息 |
| `config_get` | 获取 OpenCode 配置 |
| `project_list` | 列出所有项目 |
| `project_current` | 获取当前项目 |
| `provider_list` | 列出 AI 提供商 |
| `file_list` | 列出目录中的文件 |
| `file_content` | 读取文件内容 |
| `find_text` | 在项目中搜索文本 |
| `find_file` | 按名称查找文件 |
| `global_health` | 检查服务器健康状态 |

### MCP 客户端配置

#### Claude Desktop (macOS/Windows)

添加到 `claude_desktop_config.json`：

```json
{
  "mcpServers": {
    "oho": {
      "command": "oho",
      "args": ["mcpserver"],
      "env": {
        "OPENCODE_SERVER_HOST": "127.0.0.1",
        "OPENCODE_SERVER_PORT": "4096",
        "OPENCODE_SERVER_PASSWORD": "your-password"
      }
    }
  }
}
```

#### Cursor

添加到 Cursor 设置（JSON 模式）：

```json
{
  "mcpServers": {
    "oho": {
      "command": "oho",
      "args": ["mcpserver"],
      "env": {
        "OPENCODE_SERVER_HOST": "127.0.0.1",
        "OPENCODE_SERVER_PORT": "4096",
        "OPENCODE_SERVER_PASSWORD": "your-password"
      }
    }
  }
}
```

#### VS Code (配合 Copilot Free)

VS Code 没有原生 MCP 支持。使用 [MCP VS Code 扩展](https://github.com/modelcontextprotocol/servers)：

```json
{
  "mcpServers": {
    "oho": {
      "command": "oho",
      "args": ["mcpserver"]
    }
  }
}
```

## 许可协议

GPL v3 许可 - 参见项目根目录 [LICENSE](../LICENSE)

## 参考链接

- [OpenCode 官方文档](https://opencode.ai/docs/)
- [OpenCode 生态系统](https://opencode.ai/docs/ecosystem/)
- [OpenCode GitHub](https://github.com/anomalyco/opencode)

## 贡献

欢迎提交 Issues 和 Pull Requests！
