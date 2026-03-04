# oho CLI 操作指南 - 模块 1: 客户端初始化

> **适用版本**: oho CLI v1.0+  
> **最后更新**: 2026-03-02  
> **作者**: nanobot 🐈

---

## 📋 目录

1. [认证配置](#1-认证配置)
2. [服务器连接](#2-服务器连接)
3. [配置管理](#3-配置管理)
4. [连接验证](#4-连接验证)
5. [常见问题](#5-常见问题)

---

## 1. 认证配置

### 1.1 设置认证凭据

**命令**: `oho auth set`

```bash
# 设置服务器密码（交互式输入）
oho auth set

# 或通过环境变量设置
export OPENCODE_SERVER_PASSWORD=your_password
oho auth set
```

**说明**:
- 认证凭据存储在本地配置文件中
- 密码用于访问 OpenCode Server API
- 建议通过环境变量管理敏感信息

**预期输出**:
```
✓ 认证凭据已设置
```

---

### 1.2 环境变量配置

**推荐方式**: 在 `~/.bashrc` 或 `~/.zshrc` 中添加:

```bash
# OpenCode Server 配置
export OPENCODE_SERVER_PASSWORD="your_password"
export OPENCODE_HOST="127.0.0.1"
export OPENCODE_PORT="4096"
```

**临时设置** (当前会话有效):
```bash
export OPENCODE_SERVER_PASSWORD="cs516123456"
```

---

## 2. 服务器连接

### 2.1 默认连接参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--host` | `127.0.0.1` | 服务器主机地址 |
| `--port` | `4096` | 服务器端口 |
| `--password` | 环境变量 | 服务器密码 (可覆盖环境变量) |

### 2.2 连接服务器

```bash
# 使用默认配置连接
oho config get

# 指定主机和端口
oho config get --host 192.168.1.100 --port 4096

# 临时覆盖密码
oho config get --password "new_password"
```

**预期输出**:
```json
{
  "model": "alibaba-cn/qwen3.5-plus",
  "host": "127.0.0.1",
  "port": 4096
}
```

---

### 2.3 多服务器配置

```bash
# 开发环境
oho config get --host localhost --port 4096

# 生产环境
oho config get --host api.example.com --port 443

# 测试环境
oho config get --host test-server.local --port 4097
```

---

## 3. 配置管理

### 3.1 查看当前配置

```bash
# 获取完整配置
oho config get

# JSON 格式输出
oho config get --json
```

**示例输出**:
```json
{
  "model": "alibaba-cn/qwen3.5-plus",
  "host": "127.0.0.1",
  "port": 4096,
  "providers": [...]
}
```

---

### 3.2 设置默认模型

```bash
# 设置默认模型
oho config set --model alibaba-cn/qwen3.5-plus

# 验证设置
oho config get --json | jq '.model'
```

**可用模型**:
```bash
# 列出所有可用提供商和模型
oho config providers
```

**常见模型**:
- `alibaba-cn/qwen3.5-plus` - 通义千问 3.5 Plus (推荐)
- `alibaba-cn/kimi-k2-thinking` - Kimi K2 Thinking
- `openai/gpt-4o` - GPT-4o
- `anthropic/claude-opus-4-5` - Claude Opus

---

### 3.3 配置项说明

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| `model` | string | 默认 LLM 模型 | `alibaba-cn/qwen3.5-plus` |
| `host` | string | 服务器主机 | `127.0.0.1` |
| `port` | int | 服务器端口 | `4096` |
| `providers` | array | 可用模型列表 | 自动获取 |

---

## 4. 连接验证

### 4.1 测试连接

```bash
# 检查配置
oho config get

# 列出会话 (验证 API 访问)
oho session list

# 查看服务器状态
oho global health
```

**成功标志**:
- 返回 JSON 配置信息
- 会话列表正常显示
- 健康检查返回 `200 OK`

---

### 4.2 诊断命令

```bash
# 详细输出模式
oho session list --json

# 检查认证
oho auth set  # 重新设置凭据

# 网络连通性测试
curl -v http://127.0.0.1:4096/global/health
```

---

### 4.3 连接问题排查

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| `connection refused` | 服务器未启动 | `opencode serve --port 4096` |
| `401 Unauthorized` | 密码错误 | `oho auth set` 重新设置 |
| `timeout` | 网络问题 | 检查 `--host` 和 `--port` |
| `model not found` | 模型配置错误 | `oho config providers` 查看可用模型 |

---

## 5. 常见问题

### Q1: 如何切换不同环境的服务器？

**A**: 使用 `--host` 和 `--port` 参数:

```bash
# 本地开发
oho message add -s ses_xxx "Hello" --host localhost --port 4096

# 远程服务器
oho message add -s ses_xxx "Hello" --host api.example.com --port 443
```

---

### Q2: 密码保存在哪里？

**A**: 密码存储在:
- 环境变量: `OPENCODE_SERVER_PASSWORD` (推荐)
- 本地配置文件: `~/.config/oho/config.json` (通过 `oho auth set` 设置)

**安全建议**:
- ✅ 使用环境变量管理生产密码
- ✅ 使用密钥管理工具 (如 1Password, pass)
- ❌ 不要在代码中硬编码密码
- ❌ 不要提交密码到版本控制

---

### Q3: 如何验证模型是否可用？

**A**: 使用 `oho config providers` 命令:

```bash
oho config providers --json
```

**输出示例**:
```json
{
  "providers": [
    {
      "id": "alibaba-cn",
      "models": [
        {"id": "qwen3.5-plus", "status": "available"},
        {"id": "kimi-k2-thinking", "status": "available"}
      ]
    }
  ]
}
```

---

### Q4: 连接超时怎么办？

**A**: 检查以下步骤:

```bash
# 1. 确认服务器运行
ps aux | grep opencode

# 2. 检查端口监听
netstat -tlnp | grep 4096

# 3. 测试本地连接
curl http://127.0.0.1:4096/global/health

# 4. 重启服务器
opencode serve --port 4096
```

---

### Q5: 如何在脚本中使用 oho CLI？

**A**: 示例脚本:

```bash
#!/bin/bash

# 设置环境变量
export OPENCODE_SERVER_PASSWORD="your_password"

# 创建会话
SESSION=$(oho session create --json | jq -r '.id')

# 发送消息
oho message add -s "$SESSION" "Hello, World!"

# 获取响应
oho message list -s "$SESSION" --json | jq '.[-1].content'
```

---

## 📚 相关文档

- [模块 2: 验证](./02-validation.md) - 身份验证、权限检查
- [模块 3: 检查 Session](./03-check-session.md) - 查询、列表、详情
- [模块 4: 新建工作区](./04-create-workspace.md) - 创建和管理
- [模块 5: 指定工作区提交任务](./05-submit-task.md) - 任务提交
- [模块 6: 指定 session_id 和模型发消息](./06-send-message.md) - 消息发送
- [模块 7: 中断任务](./07-interrupt-task.md) - 任务中断
- [模块 8: 查询任务执行情况及状态](./08-query-status.md) - 状态查询

---

## 🔗 参考链接

- [OpenCode Server 文档](https://github.com/opencode-ai/opencode)
- [oho CLI GitHub](https://github.com/opencode-ai/opencode_cli)
- [API 参考](https://coding.dashscope.aliyuncs.com/v1)

---

*文档生成时间：2026-03-02 18:14 CST*  
*最后验证：2026-03-03 23:28 CST*

---

## 🔬 实际验证输出 (2026-03-03 23:28)

### 验证 1: oho config get

```bash
$ oho config get
当前配置:
  默认模型：
  主题：
  语言：
  最大 Token：0
  温度：0.00
```

---

### 验证 2: oho config providers

```bash
$ oho config providers
可用提供商:

默认模型:
  google: gemini-3-pro-preview
  minimax: MiniMax-M2.5-highspeed
  deepseek: deepseek-reasoner
  minimax-cn: MiniMax-M2.5-highspeed
  openrouter: google/gemini-3-pro-preview
  alibaba-cn: tongyi-intent-detect-v3
  opencode: big-pickle
```

---

### 验证 3: oho session list

```bash
$ oho session list
共 48 个会话:

ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
模型：   
---
ID:     ses_34c5b5c54ffehnE3JBss6tWts1
标题：   New session - 2026-03-03T12:20:37.425Z
模型：   
---
ID:     ses_35725f2eeffecp7ZPxdGfCnPkO
标题：   New session - 2026-03-01T10:03:08.433Z
模型：   
---
... (共 48 个会话)
```

---

### 验证 4: oho --help

```bash
$ oho --help
oho 是 OpenCode Server 的命令行客户端工具。
	
它提供了对 OpenCode Server API 的完整访问，允许你通过命令行管理会话、消息、配置等。

示例:
  oho session create              # 创建新会话
  oho message add -s session123   # 添加消息到会话
  oho session list                # 列出所有会话
  oho config get                  # 获取配置
  oho provider list               # 列出所有提供商

Available Commands:
  agent       代理命令
  auth        认证管理
  command     命令管理
  config      配置管理命令
  message     消息管理命令
  project     项目管理命令
  provider    提供商管理命令
  session     会话管理命令
  ...
```

---

### 验证 5: oho session get

```bash
$ oho session get ses_34dbffe0dffe8SfdMTbL53MWFP
共 1 个会话:

ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
模型：   
---
```

---

### 验证 6: oho message add (带超时错误示例)

```bash
$ oho message add -s ses_34c5b5c54ffehnE3JBss6tWts1 "测试文档完善"
DEBUG: 发送请求:
{
  "parts": [
    {
      "type": "text",
      "text": "测试文档完善"
    }
  ]
}
Error: 请求失败：Post "http://127.0.0.1:4096/session/ses_34c5b5c54ffehnE3JBss6tWts1/message": 
       context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

**说明**: 此错误表示服务器响应超时，可能原因：
- 服务器处理时间过长
- 网络连接问题
- 服务器负载过高

**解决方案**:
```bash
# 使用 --no-reply 不等待响应
oho message add -s ses_xxx "内容" --no-reply

# 检查服务器状态
ps aux | grep opencode

# 重启服务器
opencode serve --port 4096
```
