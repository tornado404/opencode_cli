# OpenCode CLI

> 让 OpenCode 成为可被其他 AI 调用和监督的命令行工具

本项目是 OpenCode 生态系统的扩展，提供基于 Bash 的完整命令行客户端，让 OpenCode 能够被其他 AI 系统调用和监督。

## 项目愿景

**oho** 的设计目标是让 OpenCode 更好地被其他 AI 调用和监督。在 [OpenCode 生态系统](https://opencode.ai/docs/zh-cn/ecosystem/) 中，存在许多类似的应用，但 **oho 是唯一一个完全基于 Bash 实现的方案**。

> "oho 在 Bash 中可调用" 代表着强大的扩展性和兼容性 —— 这是本项目独一无二的定位。

## 核心特性

### 完整的 API 覆盖

oho 基于 OpenCode REST API 构建，提供完整的命令行接口：

- ✅ 会话管理（创建、删除、继续、终止）
- ✅ 消息发送与接收
- ✅ 项目与文件操作
- ✅ 配置与提供商管理
- ✅ MCP/LSP/格式化器状态管理

### 独特的 Linux 能力

在 Linux 环境中，oho 可以做到 OpenCode CLI 暂时不具备的功能：

- 📁 **指定目录创建 Session**：在任意目录启动 AI 编程会话
- 💬 **基于 Session 继续发消息**：接续之前的会话上下文
- 🗑️ **销毁 Session**：完整管理会话生命周期
- 🔄 **会话分叉与回退**：实验性开发轻松切换

### Bash 可调用性

作为纯 Bash 实现，oho 可以：

- 被任何 AI Agent 调用
- 集成到自动化工作流
- 在 CI/CD 管道中运行
- 与其他 shell 工具无缝组合

## 快速开始

### 安装

```bash
cd oho
make build
```

### 基本用法

```bash
# 配置服务器连接
export OPENCODE_SERVER_HOST=127.0.0.1
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# 列出所有会话
oho session list

# 创建新会话
oho session create

# 在指定目录创建会话
oho session create --path /your/project

# 发送消息
oho message add -s <session-id> "帮我分析这个代码"

# 继续已有会话
oho message add -s <session-id> "继续刚才的任务"

# 销毁会话
oho session delete <session-id>
```

## 与其他生态项目的对比

| 特性 | oho | 其他生态项目 |
|------|-----|-------------|
| 实现语言 | Bash | TypeScript/Python/Go |
| AI 可调用 | ✅ 天然支持 | 需要额外适配 |
| 集成难度 | ⭐⭐⭐⭐⭐ 极低 | ⭐⭐⭐ 中等 |

## 项目结构

```
.
├── oho/                    # OpenCode Bash 客户端
│   ├── cmd/                # 命令行入口
│   ├── internal/           # 内部包
│   ├── go.mod              # Go 模块定义
│   └── README.md           # 客户端详细文档
├── docs/                   # 项目文档
│   └── plans/              # 设计计划
├── assets/                 # 资源文件
│   └── oho_cli.png         # 命令行截图
├── AGENTS.md               # AI 开发指南
└── LICENSE                 # GPL v3 许可证
```

## 命令参考

完整命令列表请参考 [oho/README.md](oho/README.md)

## 许可证

本项目基于 GPL v3 许可证开源，详见 [LICENSE](LICENSE) 文件。

## 参考资源

- [OpenCode 官方文档](https://opencode.ai/docs/zh-cn/)
- [OpenCode 生态系统](https://opencode.ai/docs/zh-cn/ecosystem/)
- [OpenCode GitHub](https://github.com/anomalyco/opencode)
