# oho - OpenCode CLI

> 让 OpenCode 成为可被其他 AI 调用和监督的命令行工具

[![GitHub Stars](https://img.shields.io/github/stars/tornado404/opencode_cli?style=flat-square)](https://github.com/tornado404/opencode_cli/stargazers)
[![License](https://img.shields.io/badge/license-GPLv3-blue?style=flat-square)](LICENSE)

oho 是 OpenCode Server 的命令行客户端工具，提供对 OpenCode Server API 的完整访问。

## 项目定位

### 独特价值

**oho** 是 [OpenCode 生态系统](https://opencode.ai/docs/zh-cn/ecosystem/) 中 **唯一一个完全基于 Bash 实现的命令行客户端**。

> "oho 在 Bash 中可调用" 代表着强大的扩展性和兼容性 —— 这是本项目独一无二的定位。

### 设计目标

让 OpenCode 更好地被其他 AI 调用和监督：

- 🤖 被任何 AI Agent 天然调用
- 🔄 集成到自动化工作流
- 📦 在 CI/CD 管道中运行
- 🔗 与其他 shell 工具无缝组合

### 独特的 Linux 能力

在 Linux 环境中，oho 可以做到 OpenCode CLI 暂时不具备的功能：

| 功能 | 说明 |
|------|------|
| 📁 指定目录创建 Session | 在任意目录启动 AI 编程会话 |
| 💬 基于 Session 继续发消息 | 接续之前的会话上下文 |
| 🗑️ 销毁 Session | 完整管理会话生命周期 |
| 🔄 会话分叉与回退 | 实验性开发轻松切换 |

## 界面预览

![oho CLI](assets/oho_cli.png)

## 功能特性

- ✅ 完整的 API 映射封装
- ✅ 支持 HTTP Basic Auth 认证
- ✅ JSON/文本双输出模式
- ✅ 配置文件和环境变量支持
- ✅ 所有会话管理操作
- ✅ 消息发送和管理
- ✅ 文件和符号查找
- ✅ TUI 界面控制
- ✅ MCP/LSP/格式化器状态管理
- 📊 **[API 完成情况](./COMPLETION.md)** - 查看实现覆盖率

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/tornado404/opencode_cli.git
cd opencode_cli/oho

# 编译
make build

# 或编译 Linux 版本
make build-linux
```

### 依赖

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

# 或使用命令行标志
oho --host 127.0.0.1 --port 4096 --password your-password session list
```

### 2. 基本用法

```bash
# 检查服务器状态
oho global health

# 列出所有会话
oho session list

# 创建新会话
oho session create

# 在指定目录创建会话
oho session create --path /your/project

# 发送消息
oho message add -s <session-id> "你好，请帮我分析这个项目"

# 继续已有会话
oho message add -s <session-id> "继续刚才的任务"

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
| 实现语言 | Bash | TypeScript/Python/Go |
| AI 可调用 | ✅ 天然支持 | 需要额外适配 |
| 跨平台 | Linux/Mac/Windows | 依赖运行时 |
| 集成难度 | ⭐⭐⭐⭐⭐ 极低 | ⭐⭐⭐ 中等 |

参考：[OpenCode 生态系统中的其他项目](https://opencode.ai/docs/zh-cn/ecosystem/)

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
oho session list                      # 列出所有会话
oho session create                    # 创建新会话
oho session create --path /path        # 在指定目录创建会话
oho session status                    # 获取所有会话状态
oho session get <id>                  # 获取会话详情
oho session delete <id>               # 删除会话
oho session update <id> --title "新标题"  # 更新会话
oho session children <id>             # 获取子会话
oho session todo <id>                 # 获取待办事项
oho session fork <id>                 # 分叉会话
oho session abort <id>                # 中止会话
oho session share <id>                # 分享会话
oho session unshare <id>              # 取消分享
oho session diff <id>                 # 获取文件差异
oho session summarize <id>            # 总结会话
oho session revert <id> --message <msg-id>  # 回退消息
oho session unrevert <id>             # 恢复回退
oho session permissions <id> <perm-id> --response allow  # 响应权限
```

### 消息管理

```bash
oho message list -s <session>         # 列出消息
oho message add -s <session> "内容"   # 发送消息
oho message get -s <session> <msg-id> # 获取消息详情
oho message prompt-async -s <session> "内容"  # 异步发送
oho message command -s <session> /help  # 执行命令
oho message shell -s <session> --agent default "ls -la"  # 运行 shell
```

### 配置管理

```bash
oho config get                      # 获取配置
oho config set --theme dark         # 更新配置
oho config providers                # 列出提供商和默认模型
```

### 提供商管理

```bash
oho provider list                   # 列出所有提供商
oho provider auth                   # 获取认证方式
oho provider oauth authorize <id>   # OAuth 授权
oho provider oauth callback <id>    # 处理回调
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
oho find file "query"               # 查找文件
oho find symbol "query"             # 查找符号
```

### 其他命令

```bash
oho agent list                      # 列出代理
oho command list                    # 列出命令
oho tool ids                        # 列出工具 ID
oho tool list --provider xxx --model xxx  # 列出工具
oho lsp status                      # LSP 状态
oho formatter status                # 格式化器状态
oho mcp list                        # MCP 服务器列表
oho mcp add <name> --config '{}'    # 添加 MCP 服务器
oho tui open-help                   # 打开帮助
oho tui show-toast --message "提示"  # 显示提示
oho auth set <provider> --credentials key=value  # 设置认证
```

## 输出格式

使用 `-j` 或 `--json` 标志以 JSON 格式输出：

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

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
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
│       ├── main.go           # 入口文件
│       ├── root.go           # 根命令
│       ├── cmd/              # 子命令
│       │   ├── global/
│       │   ├── project/
│       │   ├── session/
│       │   ├── message/
│       │   ├── configcmd/
│       │   ├── provider/
│       │   ├── file/
│       │   ├── find/
│       │   ├── tool/
│       │   ├── agent/
│       │   ├── command/
│       │   ├── lsp/
│       │   ├── formatter/
│       │   ├── mcp/
│       │   ├── tui/
│       │   └── auth/
│       └── internal/
│           ├── client/       # HTTP 客户端
│           ├── config/       # 配置管理
│           ├── types/        # 类型定义
│           └── util/         # 工具函数
├── Makefile
├── build.sh
└── README.md
```

## 许可证

GPL v3 License - 详见项目根目录 [LICENSE](../LICENSE)

## 参考资源

- [OpenCode 官方文档](https://opencode.ai/docs/zh-cn/)
- [OpenCode 生态系统](https://opencode.ai/docs/zh-cn/ecosystem/)
- [OpenCode GitHub](https://github.com/anomalyco/opencode)

## 贡献

欢迎提交 Issue 和 Pull Request！
