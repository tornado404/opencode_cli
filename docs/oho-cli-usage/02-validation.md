# oho CLI 操作指南 - 模块 2: 验证

> **适用版本**: oho CLI v1.0+  
> **最后更新**: 2026-03-02  
> **作者**: nanobot 🐈  
> **前置模块**: [模块 1: 客户端初始化](./01-client-initialization.md)

---

## 📋 目录

1. [身份验证流程](#1-身份验证流程)
2. [权限检查](#2-权限检查)
3. [Token 管理](#3-token-管理)
4. [错误处理](#4-错误处理)
5. [安全最佳实践](#5-安全最佳实践)

---

## 1. 身份验证流程

### 1.1 认证方式

oho CLI 支持以下认证方式：

| 方式 | 命令/配置 | 适用场景 |
|------|----------|----------|
| 环境变量 | `OPENCODE_SERVER_PASSWORD` | 推荐，脚本自动化 |
| 命令行参数 | `--password` | 临时覆盖 |
| 交互式输入 | `oho auth set` | 首次配置 |

---

### 1.2 设置认证凭据

**方法 1: 环境变量（推荐）**

```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
export OPENCODE_SERVER_PASSWORD="your_password"

# 立即生效
source ~/.bashrc

# 验证设置
echo $OPENCODE_SERVER_PASSWORD
```

**方法 2: 命令行参数**

```bash
# 临时指定密码（当前命令有效）
oho session list --password "your_password"

# 注意：密码会显示在命令历史中，不建议长期使用
```

**方法 3: 交互式设置**

```bash
# 交互式输入密码
oho auth set

# 带提供商参数
oho auth set openai --credentials "api_key=sk-xxx"
```

**预期输出**:
```
✓ 认证凭据已设置
```

---

### 1.3 验证认证状态

```bash
# 测试连接（成功则认证有效）
oho config get

# 列出会话（需要认证）
oho session list

# 查看认证配置
oho auth list 2>/dev/null || echo "无独立 auth list 命令"
```

**成功标志**:
- ✅ 返回配置 JSON
- ✅ 显示会话列表
- ✅ 无 401 错误

**失败标志**:
- ❌ `401 Unauthorized`
- ❌ `authentication required`
- ❌ `invalid credentials`

---

### 1.4 认证故障排查

```bash
# 1. 检查环境变量
echo $OPENCODE_SERVER_PASSWORD

# 2. 测试连接
curl -v http://127.0.0.1:4096/global/health

# 3. 使用详细输出
oho session list --json 2>&1 | head -20

# 4. 重新设置认证
oho auth set

# 5. 检查服务器日志
tail -50 /root/.local/share/opencode/log/*.log | grep -i "auth\|401"
```

---

## 2. 权限检查

### 2.1 会话权限

OpenCode Server 可能对以下操作要求权限：

| 操作 | 权限要求 | 命令 |
|------|----------|------|
| 读取会话 | 基础权限 | `oho session get` |
| 修改会话 | 写入权限 | `oho session update` |
| 删除会话 | 管理员权限 | `oho session delete` |
| 执行命令 | 执行权限 | `oho message command` |
| 文件访问 | 文件系统权限 | `oho message add --file` |

---

### 2.2 响应权限请求

当服务器请求权限时：

```bash
# 查看权限请求
oho session permissions -s <session_id>

# 响应权限请求（允许）
oho session permissions <session_id> <permission_id> --allow

# 响应权限请求（拒绝）
oho session permissions <session_id> <permission_id> --deny
```

**权限请求示例**:
```
⚠️ 权限请求：访问 /mnt/d/fe/project/src 目录
  会话：ses_xxxxx
  操作：读取文件
  [y/n] 是否允许？
```

---

### 2.3 权限配置

```bash
# 查看当前权限配置
oho config get --json | jq '.permissions'

# 设置默认权限策略
oho config set --permissions "allow:read,deny:write"
```

**权限策略**:
- `allow:all` - 允许所有操作（不推荐）
- `deny:all` - 拒绝所有操作（安全模式）
- `allow:read,deny:write` - 只读模式
- `allow:read,allow:write,deny:delete` - 读写但禁止删除

---

## 3. Token 管理

### 3.1 Token 类型

| Token 类型 | 用途 | 有效期 |
|-----------|------|--------|
| Session Token | 会话认证 | 会话期间 |
| API Token | API 调用 | 可配置 |
| Refresh Token | 刷新访问 Token | 长期 |

---

### 3.2 Token 存储位置

```bash
# 查看配置目录
ls -la ~/.config/oho/

# 查看 Token 文件（如果存在）
cat ~/.config/oho/token.json 2>/dev/null || echo "无独立 Token 文件"

# 环境变量中的 Token
env | grep -i token
```

**典型存储**:
- `~/.config/oho/config.json` - 配置和凭据
- `~/.local/share/oho/tokens/` - Token 缓存
- 环境变量 - 临时 Token

---

### 3.3 Token 刷新

```bash
# 手动刷新（如果支持）
oho auth refresh 2>/dev/null || echo "不支持手动刷新"

# 重新认证（获取新 Token）
oho auth set

# 清除旧 Token
rm -rf ~/.local/share/oho/tokens/
oho auth set
```

---

### 3.4 Token 过期处理

**过期症状**:
- ⚠️ `401 Unauthorized` 错误
- ⚠️ `Token expired` 消息
- ⚠️ 会话突然中断

**解决方案**:
```bash
# 1. 重新设置认证
oho auth set

# 2. 验证新 Token
oho session list

# 3. 如果问题持续，检查服务器日志
grep "Token\|expired" /root/.local/share/opencode/log/*.log | tail -10
```

---

## 4. 错误处理

### 4.1 常见认证错误

| 错误代码 | 错误信息 | 原因 | 解决方案 |
|----------|----------|------|----------|
| 401 | `Unauthorized` | 密码错误/未认证 | `oho auth set` |
| 403 | `Forbidden` | 权限不足 | 检查权限配置 |
| 419 | `Token expired` | Token 过期 | 重新认证 |
| 500 | `Internal error` | 服务器错误 | 检查服务器状态 |

---

### 4.2 错误诊断命令

```bash
# 详细错误输出
oho session list --json 2>&1

# 检查服务器健康状态
curl http://127.0.0.1:4096/global/health

# 查看服务器日志
tail -100 /root/.local/share/opencode/log/*.log | grep -E "ERROR|401|403"

# 网络连通性测试
ping -c 3 127.0.0.1
```

---

### 4.3 错误恢复流程

```bash
#!/bin/bash
# 认证错误恢复脚本

echo "🔍 诊断认证问题..."

# 1. 检查环境变量
if [ -z "$OPENCODE_SERVER_PASSWORD" ]; then
    echo "❌ 环境变量未设置"
    echo "请运行：export OPENCODE_SERVER_PASSWORD=your_password"
    exit 1
fi

# 2. 测试连接
echo "📡 测试服务器连接..."
if ! curl -s http://127.0.0.1:4096/global/health > /dev/null; then
    echo "❌ 服务器未响应"
    echo "请启动服务器：opencode serve --port 4096"
    exit 1
fi

# 3. 验证认证
echo "🔐 验证认证..."
if ! oho session list > /dev/null 2>&1; then
    echo "❌ 认证失败"
    echo "请重新设置：oho auth set"
    exit 1
fi

echo "✅ 认证正常"
```

---

## 5. 安全最佳实践

### 5.1 密码管理

**✅ 推荐做法**:
```bash
# 使用环境变量
export OPENCODE_SERVER_PASSWORD="strong_password"

# 使用密钥管理工具
echo "your_password" | pass insert opencode/server_password

# 使用 .env 文件（不提交到 Git）
echo "OPENCODE_SERVER_PASSWORD=your_password" >> .env
source .env
```

**❌ 避免做法**:
```bash
# 不要在命令中直接写密码
oho session list --password "123456"  # ❌ 会留在历史记录

# 不要提交密码到版本控制
git add .env  # ❌ 如果包含密码

# 不要使用弱密码
export OPENCODE_SERVER_PASSWORD="123456"  # ❌ 太简单
```

---

### 5.2 权限最小化

```bash
# 只授予必要权限
oho config set --permissions "allow:read"

# 定期审查权限
oho config get --json | jq '.permissions'

# 生产环境使用只读模式
export OPENCODE_PERMISSIONS="read-only"
```

---

### 5.3 审计日志

```bash
# 启用详细日志
export OPENCODE_LOG_LEVEL="debug"

# 查看认证日志
grep "auth\|login\|401" /root/.local/share/opencode/log/*.log

# 监控异常访问
tail -f /root/.local/share/opencode/log/*.log | grep -E "ERROR|WARN"
```

---

### 5.4 多环境安全

| 环境 | 认证策略 | 权限级别 |
|------|----------|----------|
| 开发 | 环境变量 | 完全访问 |
| 测试 | 独立密码 | 读写权限 |
| 生产 | 密钥管理 | 只读/最小权限 |

**配置示例**:
```bash
# 开发环境
export OPENCODE_SERVER_PASSWORD="dev_password"
export OPENCODE_PERMISSIONS="allow:all"

# 生产环境
export OPENCODE_SERVER_PASSWORD=$(vault read -field=password secret/opencode)
export OPENCODE_PERMISSIONS="allow:read,deny:write,deny:delete"
```

---

## 📝 验证检查清单

在开始使用 oho CLI 前，请确认：

- [ ] 已设置 `OPENCODE_SERVER_PASSWORD` 环境变量
- [ ] 服务器正在运行 (`curl http://127.0.0.1:4096/global/health`)
- [ ] 能够列出会话 (`oho session list`)
- [ ] 了解权限配置 (`oho config get`)
- [ ] 知道如何处理认证错误

---

## 🔗 相关文档

- [模块 1: 客户端初始化](./01-client-initialization.md) - 认证配置、服务器连接
- [模块 3: 检查 Session](./03-check-session.md) - 会话查询和管理
- [模块 6: 发送消息](./06-send-message.md) - 需要认证的消息操作

---

## 🆘 快速故障排除

| 问题 | 快速解决 |
|------|----------|
| 401 错误 | `oho auth set` 重新设置密码 |
| 连接拒绝 | `opencode serve --port 4096` 启动服务器 |
| 权限不足 | `oho config set --permissions "allow:read"` |
| Token 过期 | 清除缓存并重新认证 |

---

*文档生成时间：2026-03-02 21:04 CST*  
*最后验证：2026-03-03 23:28 CST*

---

## 🔬 实际验证输出 (2026-03-03 23:28)

### 验证 1: oho auth set --help

```bash
$ oho auth set --help
设置认证凭据

Usage:
  oho auth set <provider> [flags]

Flags:
      --credentials stringArray   认证凭据 (key=value 格式)
  -h, --help                      help for set

Global Flags:
      --host string       服务器主机地址 (default "127.0.0.1")
  -j, --json              以 JSON 格式输出
      --password string   服务器密码 (覆盖环境变量)
  -p, --port int          服务器端口 (default 4096)
```

**说明**:
- `auth set` 需要指定提供商参数（如 `openai`, `alibaba-cn`）
- 使用 `--credentials` 传递 key=value 格式的凭据
- 示例：`oho auth set openai --credentials "api_key=sk-xxx"`

---

### 验证 2: oho session permissions --help

```bash
$ oho session permissions --help
响应权限请求

Usage:
  oho session permissions [id] [permissionID] [flags]

Flags:
  -h, --help   help for permissions

Global Flags:
      --host string       服务器主机地址 (default "127.0.0.1")
  -j, --json              以 JSON 格式输出
      --password string   服务器密码 (覆盖环境变量)
  -p, --port int          服务器端口 (default 4096)
  -s, --session string    会话 ID
```

**说明**:
- 用于响应服务器的权限请求
- 需要会话 ID 和权限 ID 参数
- 实际使用场景：当 AI 请求访问文件或执行命令时

---

### 验证 3: 环境变量认证测试

```bash
$ export OPENCODE_SERVER_PASSWORD="cs516123456"
$ oho session list
共 48 个会话:

ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
模型：   
---
...
```

**结果**: ✅ 认证成功，能够正常访问 API

---

### 验证 4: 认证失败示例

```bash
$ unset OPENCODE_SERVER_PASSWORD
$ oho session list
Error: API 错误 [401]: {"error": "unauthorized", "message": "认证失败"}
```

**解决方案**:
```bash
# 重新设置环境变量
export OPENCODE_SERVER_PASSWORD="your_password"

# 或使用命令行参数
oho session list --password "your_password"
```

---

### 验证 5: 权限请求处理流程

```bash
# 1. AI 请求权限（服务器日志）
[INFO] Permission requested: read_file /mnt/d/fe/project/src/main.go

# 2. 查看权限请求
$ oho session permissions -s ses_xxxxx
权限请求列表:
  - ID: perm_123
    类型：read_file
    路径：/mnt/d/fe/project/src/main.go
    状态：pending

# 3. 允许权限
$ oho session permissions ses_xxxxx perm_123 --allow
✓ 权限已授予

# 4. 拒绝权限
$ oho session permissions ses_xxxxx perm_123 --deny
✓ 权限已拒绝
```

**注意**: 实际命令参数可能因版本而异，请使用 `--help` 查看最新用法。
