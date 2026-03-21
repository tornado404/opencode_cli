# OpenCode Server API 技能手册

> **版本**: v4.5.0 | **Server**: v1.2.15 | **oho CLI**: v1.1.0
> 
> **⚠️ 重要**: 所有 OpenCode 操作必须通过 **oho CLI** 发送，不再直接使用 Python SDK 或 curl。

---

## 🚫 Telegram 消息规范（重要）

**禁止发送的内容**:
- ❌ 不要发送"有更新"、"技能已更新"等通知类消息
- ❌ 不要发送格式化的版本更新日志
- ❌ 不要在任务提交后发送冗长的确认报告

**正确的做法**:
- ✅ 任务提交后只返回简洁的会话 ID 和状态
- ✅ 有实质性进展时才通知用户
- ✅ 保持消息简短、直接

**示例**:
```
✅ 任务已提交
会话：ses_xxx
项目：babylon3DWorld
```

---

## 🚀 快速使用

### 提交任务（异步模式，避免超时）

```bash
# 最常用：一条命令完成会话创建 + 消息发送
oho add "任务描述" --directory /mnt/d/fe/project --no-reply

# 指定 Agent
oho add "@hephaestus 分析性能瓶颈" --directory /mnt/d/fe/project --no-reply

# 指定标题
oho add "修复登录 bug" --title "Bug 修复" --directory /mnt/d/fe/project --no-reply

# 附加文件
oho add "分析日志" --file /var/log/app.log --directory /mnt/d/fe/project --no-reply
```

### 关键参数

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `--title` | string | 会话标题（不提供则自动生成） | 自动生成 |
| `--parent` | string | 父会话 ID（用于创建子会话） | - |
| `--directory` | string | 会话工作目录 | 当前目录 |
| `--agent` | string | 消息的 Agent ID | - |
| `--model` | string | 消息的模型 ID（如 `provider:model`） | 默认模型 |
| `--no-reply` | bool | 不等待 AI 响应 | false |
| `--system` | string | 系统提示词 | - |
| `--tools` | string[] | 工具列表（可多次指定） | - |
| `--file` | string[] | 文件附件（可多次指定） | - |
| `-j, --json` | bool | JSON 格式输出 | false |

---

## 📋 常用命令

### 会话管理

```bash
# 列出所有会话
oho session list

# 按 ID 过滤（模糊匹配）
oho session list --id ses_abc

# 按标题过滤（模糊）
oho session list --title "测试"

# 按项目 ID 过滤
oho session list --project-id proj1

# 按目录过滤
oho session list --directory babylon

# 按时间戳过滤
oho session list --created 1773537883643
oho session list --updated 1773538142930

# 按状态过滤
oho session list --status running    # running/completed/error/aborted/idle
oho session list --running           # 仅显示运行中的会话

# 排序和分页
oho session list --sort updated --order desc  # 按 updated 降序
oho session list --limit 10 --offset 0        # 分页

# JSON 输出
oho session list -j

# 创建会话
oho session create
oho session create --title "名称"
oho session create --parent ses_xxx    # 创建子会话
oho session create --path /path        # 在指定目录创建

# 获取会话详情
oho session get <id>

# 更新会话
oho session update <id> --title "New Title"

# 获取子会话
oho session children <id>

# 获取待办事项
oho session todo <id>

# 会话分支
oho session fork <id>

# 中止会话
oho session abort <id>

# 分享/取消分享会话
oho session share <id>
oho session unshare <id>

# 获取文件差异
oho session diff <id>

# 会话摘要
oho session summarize <id>

# 回滚消息
oho session revert <id> --message <msg-id>
oho session unrevert <id>

# 响应权限请求
oho session permissions <id> <perm-id> --response allow

# 删除会话
oho session delete ses_xxx
```

### 消息管理

```bash
# 查看消息列表
oho message list -s ses_xxx

# 获取消息详情
oho message get -s ses_xxx <msg-id>

# 发送消息（同步）
oho message add -s ses_xxx "继续任务"

# 异步发送（不等待响应）
oho message prompt-async -s ses_xxx "任务内容"

# 执行命令
oho message command -s ses_xxx /help

# 运行 shell 命令
oho message shell -s ses_xxx --agent default "ls -la"
```

### 项目管理

```bash
# 列出所有项目
oho project list

# 获取当前项目
oho project current

# 获取当前路径
oho path

# 获取 VCS 信息
oho vcs

# 处置当前实例
oho instance dispose
```

### 全局命令

```bash
# 检查服务器健康状态
oho global health

# 监听全局事件流（SSE）
oho global event
```

### 配置管理

```bash
# 获取配置
oho config get

# 更新配置
oho config set --theme dark

# 列出提供者
oho config providers
```

### 提供者管理

```bash
# 列出所有提供者
oho provider list

# 获取认证方法
oho provider auth

# OAuth 授权
oho provider oauth authorize <id>

# 处理回调
oho provider oauth callback <id>
```

### 文件操作

```bash
# 列出文件
oho file list [path]

# 读取文件内容
oho file content <path>

# 获取文件状态
oho file status
```

### 查找功能

```bash
# 搜索文本
oho find text "pattern"

# 查找文件
oho find file "query"

# 查找符号
oho find symbol "query"
```

### 其他命令

```bash
# 列出 Agents
oho agent list

# 列出命令
oho command list

# 列出工具 ID
oho tool ids

# 列出工具
oho tool list --provider xxx --model xxx

# LSP 状态
oho lsp status

# Formatter 状态
oho formatter status

# MCP 服务器
oho mcp list
oho mcp add <name> --config '{}'

# TUI
oho tui open-help
oho tui show-toast --message "message"

# 认证设置
oho auth set <provider> --credentials key=value
```

---

## ⚠️ 超时处理（重要）

**必须使用 `--no-reply` 参数**，避免 AI 调用超时：

```bash
# ✅ 正确：异步提交
oho add "任务" --directory /mnt/d/fe/project --no-reply

# ❌ 错误：同步等待（可能超时）
oho add "任务" --directory /mnt/d/fe/project
```

**后台轮询等待完成**（可选）:

```bash
#!/bin/bash
session_id=$(oho add "任务" --json | jq -r '.sessionId')

while true; do
    count=$(oho message list -s "$session_id" -j | jq 'length')
    [ "$count" -ge 2 ] && echo "✅ 完成" && break
    echo "⏳ 执行中... ($count 条消息)"
    sleep 10
done
```

---

## 🤖 Agent 系统

| Agent | 角色 | 适用场景 |
|-------|------|---------|
| **@sisyphus** | 主协调器 | 大型功能开发、并行执行 |
| **@hephaestus** | 深度工作者 | 代码探索、性能优化 |
| **@prometheus** | 战略规划师 | 需求澄清、架构设计 |

---

## 📝 实战案例

### babylon3DWorld 项目任务

```bash
#!/bin/bash
# 提交编码任务

oho add "@hephaestus ulw 优化编辑器与 world 页面的导航逻辑

**编码目标**:
1. 编辑器返回 world 页面时：直接刷新页面，不再判断是否 editor 造成了改动
2. world 页进入 editor 时：不再缓存 world 本身

**关键词**: ulw" \
  --directory /mnt/d/fe/babylon3DWorld \
  --title "ulw - 优化编辑器导航逻辑" \
  --no-reply

echo "✅ 任务已提交"
```

### 多项目批量任务

```bash
#!/bin/bash
# 批量提交任务

oho add "任务 1" --directory /mnt/d/fe/babylon3DWorld --no-reply
oho add "任务 2" --directory /mnt/d/fe/wujimanager --no-reply
oho add "任务 3" --directory /mnt/d/fe/armdraw --no-reply

echo "✅ 所有任务已提交"
```

---

## 🔧 故障排除

### 401 Unauthorized
```bash
# 检查密码
echo $OPENCODE_SERVER_PASSWORD

# 或命令行指定
oho --password cs516123456 session list
```

### Session not found
```bash
# 重新创建会话
oho session create --path /mnt/d/fe/babylon3DWorld
```

### 任务超时
```bash
# 使用 --no-reply 异步提交
oho add "任务" --directory /mnt/d/fe/project --no-reply
```

---

## 🔗 MCP Server

oho 可以作为 MCP (Model Context Protocol) 服务器，允许外部 MCP 客户端（如 Claude Desktop、Cursor 等）通过 MCP 协议调用 OpenCode API。

### 启动 MCP 服务器

```bash
# 启动 MCP 服务器（stdio 模式）
oho mcpserver
```

### 可用的 MCP 工具

| 工具 | 说明 |
|------|------|
| `session_list` | 列出所有会话 |
| `session_create` | 创建新会话 |
| `session_get` | 获取会话详情 |
| `session_delete` | 删除会话 |
| `session_status` | 获取所有会话状态 |
| `message_list` | 列出会话中的消息 |
| `message_add` | 发送消息到会话 |
| `config_get` | 获取 OpenCode 配置 |
| `project_list` | 列出所有项目 |
| `project_current` | 获取当前项目 |
| `provider_list` | 列出 AI 提供者 |
| `file_list` | 列出目录中的文件 |
| `file_content` | 读取文件内容 |
| `find_text` | 在项目中搜索文本 |
| `find_file` | 按名称查找文件 |
| `global_health` | 检查服务器健康状态 |

### MCP 客户端配置

#### Claude Desktop (macOS/Windows)

在 `claude_desktop_config.json` 中添加：

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

在 Cursor 设置（JSON 模式）中添加：

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

---

## 🔗 相关资源

- **OpenCode 官方文档**: https://opencode.ai/docs/
- **oho CLI 仓库**: https://github.com/tornado404/opencode_cli
- **OpenAPI Spec**: http://localhost:4096/doc

---

*创建时间：2026-02-27 13:46 CST*  
*最后更新：2026-03-21 09:28 CST*
