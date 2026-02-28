# oho - OpenCode Bash 命令行工具

> 唯一完全基于 Bash 实现的 OpenCode 命令行客户端。

## oho 是什么？

**oho** 是一个纯 Bash 实现的命令行工具，提供对 OpenCode Server API 的完整访问，专为 AI Agent 和自动化工作流调用而设计。

## 快速开始

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

# 发送消息
oho message add -s <session-id> "帮我分析这个项目"

# 销毁会话
oho session delete <session-id>
```

## 完整文档

完整的命令参考、安装说明和详细文档，请查看 [主项目 README](../README_zh.md)。
