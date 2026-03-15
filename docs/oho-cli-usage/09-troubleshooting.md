# oho CLI 问题排查指南

> **适用版本**: oho CLI v1.0+  
> **最后更新**: 2026-03-15  
> **状态**: 🟢 维护中

---

## 📋 目录

1. [快速诊断流程](#1-快速诊断流程)
2. [连接问题](#2-连接问题)
3. [认证问题](#3-认证问题)
4. [消息无响应问题](#4-消息无响应问题)
5. [会话状态问题](#5-会话状态问题)
6. [模型/提供商问题](#6-模型提供商问题)
7. [文件附件问题](#7-文件附件问题)
8. [性能问题](#8-性能问题)
9. [诊断脚本](#9-诊断脚本)
10. [获取帮助](#10-获取帮助)

---

## 1. 快速诊断流程

### 1.1 5 分钟快速检查

按顺序执行以下命令，快速定位问题：

```bash
# Step 1: 检查服务器是否运行
curl -s http://localhost:4096/global/health || echo "❌ 服务器未运行"

# Step 2: 检查环境变量
echo "密码配置：${OPENCODE_SERVER_PASSWORD:+已设置}" || echo "❌ 未设置密码"

# Step 3: 检查认证
oho config get || echo "❌ 认证失败"

# Step 4: 检查会话
oho session list || echo "❌ 无法获取会话"

# Step 5: 检查可用模型
oho config providers || echo "❌ 无法获取模型列表"
```

---

### 1.2 问题定位流程图

```
┌─────────────────────────────────────────────────────────────┐
│                    消息提交后无响应                           │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
            ┌─────────────────────┐
            │ 1. 服务器是否运行？   │
            │ curl /global/health │
            └──────────┬──────────┘
                       │
           ┌───────────┴───────────┐
           │ NO                    │ YES
           ▼                       ▼
    ┌─────────────┐         ┌─────────────────────┐
    │ 启动服务器   │         │ 2. 认证是否正确？    │
    │ opencode    │         │ oho config get      │
    │ serve       │         └──────────┬──────────┘
    └─────────────┘                    │
                              ┌────────┴────────┐
                              │ NO              │ YES
                              ▼                 ▼
                       ┌─────────────┐   ┌─────────────────────┐
                       │ 检查密码配置 │   │ 3. 会话是否忙碌？    │
                       │ 检查用户名   │   │ oho session status  │
                       └─────────────┘   └──────────┬──────────┘
                                                    │
                                          ┌─────────┴─────────┐
                                          │ YES (忙碌)        │ NO
                                          ▼                   ▼
                                   ┌─────────────┐     ┌─────────────────────┐
                                   │ 等待或中止   │     │ 4. 模型是否可用？    │
                                   │ 当前任务     │     │ oho config providers│
                                   └─────────────┘     └──────────┬──────────┘
                                                                   │
                                                         ┌─────────┴─────────┐
                                                         │ NO                │ YES
                                                         ▼                   ▼
                                                  ┌─────────────┐     ┌─────────────────────┐
                                                  │ 配置提供商   │     │ 5. 检查请求格式     │
                                                  │ 检查 API Key │     │ parts 数组是否正确   │
                                                  └─────────────┘     └─────────────────────┘
```

---

## 2. 连接问题

### 2.1 服务器无法连接

**症状**:
```bash
$ oho session list
Error: 请求失败：dial tcp 127.0.0.1:4096: connect: connection refused
```

**可能原因**:
1. OpenCode 服务器未启动
2. 服务器运行在非默认端口
3. 防火墙阻止连接

**解决方案**:

```bash
# 1. 检查服务器进程
ps aux | grep opencode

# 2. 检查端口监听
netstat -tlnp | grep 4096
# 或
lsof -i :4096

# 3. 启动服务器（如果未运行）
opencode serve --port 4096 --hostname 127.0.0.1

# 4. 测试连接
curl http://localhost:4096/global/health

# 5. 如果是端口问题，指定正确端口
export OPENCODE_SERVER_PORT=4096
oho session list
```

---

### 2.2 服务器运行在远程主机

**配置远程连接**:

```bash
# 设置环境变量
export OPENCODE_SERVER_HOST=192.168.1.100
export OPENCODE_SERVER_PORT=4096
export OPENCODE_SERVER_PASSWORD=your-password

# 或使用配置文件 ~/.config/oho/config.json
{
  "host": "192.168.1.100",
  "port": 4096,
  "password": "your-password"
}

# 测试连接
oho session list
```

---

### 2.3 CORS 错误（Web 客户端）

**症状**: 浏览器控制台显示 CORS 错误

**解决方案**:

```bash
# 启动服务器时启用 CORS
opencode serve --cors http://localhost:5173 --cors https://app.example.com

# 或配置允许所有来源（开发环境）
opencode serve --cors "*"
```

---

## 3. 认证问题

### 3.1 401 认证失败

**症状**:
```bash
$ oho session list
Error: API 错误 [401]: 认证失败：用户名或密码错误
```

**解决方案**:

```bash
# 方法 1: 环境变量（推荐）
export OPENCODE_SERVER_PASSWORD=your-password
oho session list

# 方法 2: 命令行参数
oho --password your-password session list

# 方法 3: 配置文件
cat > ~/.config/oho/config.json <<EOF
{
  "password": "your-password"
}
EOF

# 方法 4: 检查用户名（如果自定义过）
export OPENCODE_SERVER_USERNAME=opencode  # 默认用户名

# 验证认证
curl -u "opencode:your-password" http://localhost:4096/config
```

---

### 3.2 密码包含特殊字符

**问题**: 密码包含 `$`, `!`, `\` 等字符时被 shell 解释

**解决方案**:

```bash
# 使用单引号
export OPENCODE_SERVER_PASSWORD='my$p@ssw0rd!'

# 或使用双引号并转义
export OPENCODE_SERVER_PASSWORD="my\$p@ssw0rd\!"

# 或写入配置文件（推荐）
cat > ~/.config/oho/config.json <<EOF
{
  "password": "my$p@ssw0rd!"
}
EOF
```

---

### 3.3 认证配置检查清单

```bash
#!/bin/bash
# 认证配置检查脚本

echo "=== 认证配置检查 ==="

# 检查环境变量
echo -n "OPENCODE_SERVER_HOST: "
echo "${OPENCODE_SERVER_HOST:-未设置 (默认：127.0.0.1)}"

echo -n "OPENCODE_SERVER_PORT: "
echo "${OPENCODE_SERVER_PORT:-未设置 (默认：4096)}"

echo -n "OPENCODE_SERVER_USERNAME: "
echo "${OPENCODE_SERVER_USERNAME:-未设置 (默认：opencode)}"

echo -n "OPENCODE_SERVER_PASSWORD: "
if [ -n "$OPENCODE_SERVER_PASSWORD" ]; then
    echo "已设置 (${#OPENCODE_SERVER_PASSWORD} 字符)"
else
    echo "未设置 ❌"
fi

# 检查配置文件
if [ -f ~/.config/oho/config.json ]; then
    echo -e "\n配置文件：~/.config/oho/config.json"
    cat ~/.config/oho/config.json | grep -v password
else
    echo -e "\n配置文件：不存在"
fi

# 测试认证
echo -e "\n测试认证..."
curl -s -w "\nHTTP 状态码：%{http_code}\n" \
    -u "opencode:${OPENCODE_SERVER_PASSWORD}" \
    http://localhost:4096/config
```

---

## 4. 消息无响应问题

### 4.0 超时错误（最新）

**症状**:
```bash
Error: 请求失败：Post "http://127.0.0.1:4096/session/.../message": 
       context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

**原因分析**:
- ✅ 请求已成功发送到服务器
- ✅ 服务器正在处理（AI 正在思考）
- ❌ 响应时间超过了客户端超时限制（默认 30 秒）

**这不是 Bug**：消息实际上已经提交成功，只是 CLI 等不及服务器响应就断开了。

**解决方案**:

```bash
# 方法 1: 设置环境变量增加超时时间（推荐）
# 单位：秒，默认 300 秒（5 分钟）
export OPENCODE_CLIENT_TIMEOUT=600  # 10 分钟
oho message add -s ses_xxx "复杂任务"

# 方法 2: 使用 --no-reply 不等待响应
oho message add -s ses_xxx "任务" --no-reply

# 方法 3: 使用异步提交
oho message prompt-async -s ses_xxx "任务"

# 方法 4: 等待后查看结果
sleep 30
oho message list -s ses_xxx --limit 3
```

**超时配置说明**:
| 配置 | 默认值 | 说明 |
|------|--------|------|
| `OPENCODE_CLIENT_TIMEOUT` | 300 秒（5 分钟） | 客户端 HTTP 请求超时 |
| 最小值 | 1 秒 | - |
| 推荐值 | 300-600 秒 | 根据任务复杂度调整 |

---

### 4.1 症状识别

**症状 A: 请求卡住不返回**
```bash
$ oho message add -s ses_xxx "Hello"
[长时间等待，无输出]
```

**症状 B: 返回成功但无 AI 响应**
```bash
$ oho message add -s ses_xxx "Hello" --no-reply
消息已发送
[但后续查看消息历史，没有 assistant 响应]
```

**症状 C: 返回错误**
```bash
Error: API 错误 [500]: Internal server error
```

---

### 4.2 原因 A: 会话忙碌

**检查方法**:
```bash
# 检查会话状态
oho session status

# 或检查特定会话
oho session get ses_xxx --json | jq '.status'
```

**预期输出**:
```json
{
  "ses_xxx": {
    "type": "idle",      // 空闲 - 正常
    "isWorking": false
  }
}
```

```json
{
  "ses_xxx": {
    "type": "busy",      // 忙碌 - 需要等待
    "isWorking": true
  }
}
```

**解决方案**:
```bash
# 方法 1: 等待当前任务完成
sleep 30
oho session status

# 方法 2: 中止当前任务
oho session abort ses_xxx

# 方法 3: 创建新会话
oho session create
```

---

### 4.3 原因 B: noReply 参数

**问题**: 使用 `--no-reply` 或 `noReply: true` 时，消息不会触发 AI 响应

**检查方法**:
```bash
# 查看发送的消息
oho message list -s ses_xxx --limit 3

# 如果最后一条是 user 消息且没有 assistant 响应
# 可能是因为使用了 --no-reply
```

**解决方案**:
```bash
# 正确：等待 AI 响应（默认行为）
oho message add -s ses_xxx "Hello"

# 错误：不等待响应
oho message add -s ses_xxx "Hello" --no-reply  # 不会触发 AI
```

---

### 4.4 原因 C: 模型/提供商问题

**检查方法**:
```bash
# 检查可用模型
oho config providers

# 检查默认模型
oho config get --json | jq '.model'
```

**常见问题**:

| 问题 | 症状 | 解决方案 |
|------|------|----------|
| 未配置提供商 | `config providers` 返回空 | 在 OpenCode 中配置提供商 |
| API Key 过期 | 请求超时或返回 401 | 更新 API Key |
| 模型不可用 | 返回 `model not found` | 使用 `config providers` 查看可用模型 |
| 配额耗尽 | 返回 `quota exceeded` | 等待配额重置或升级 |

**解决方案**:
```bash
# 1. 列出可用提供商
oho config providers

# 2. 设置默认模型（通过配置文件）
cat > ~/.config/oho/config.json <<EOF
{
  "model": {
    "providerID": "alibaba-cn",
    "modelID": "qwen3.5-plus"
  }
}
EOF

# 3. 测试模型可用性
oho message add -s ses_xxx "测试" --no-reply
```

---

### 4.5 原因 D: 权限请求等待确认

**症状**: 服务器等待用户确认权限（如文件写入、命令执行）

**检查方法**:
```bash
# 检查是否有待处理的权限请求
oho session get ses_xxx --json | jq '.permissions'

# 查看服务器日志
# （在运行 opencode serve 的终端查看）
```

**解决方案**:
```bash
# 方法 1: 在 TUI 中确认权限
# OpenCode TUI 会弹出权限确认对话框

# 方法 2: 通过 API 响应权限请求
oho session permissions ses_xxx perm_id --response allow

# 方法 3: 配置自动允许（开发环境）
# 在 opencode.json 中配置：
{
  "permissions": {
    "write": "allow",
    "bash": "allow"
  }
}
```

---

### 4.6 原因 E: 请求格式错误

**症状**: 返回 400 错误

**检查方法**:
```bash
# 使用 debug 模式查看详细请求
oho message add -s ses_xxx "Hello" 2>&1 | grep -A 20 "DEBUG"
```

**正确格式**:
```json
{
  "parts": [
    {
      "type": "text",
      "text": "Hello"
    }
  ]
}
```

**错误格式示例**:
```json
// 错误 1: 缺少 parts 数组
{
  "message": "Hello"  // ❌ 应该是 parts: [{type: "text", text: "..."}]
}

// 错误 2: model 格式错误
{
  "model": "alibaba-cn/qwen3.5-plus",  // ❌ 应该是对象
  "parts": [...]
}

// 正确的 model 格式
{
  "model": {
    "providerID": "alibaba-cn",
    "modelID": "qwen3.5-plus"
  },
  "parts": [...]
}
```

---

## 5. 会话状态问题

### 5.1 会话 ID 无效

**症状**:
```bash
Error: API 错误 [400]: Invalid string: must start with "ses"
```

**原因**: 使用了 Slug（如 `tidy-panda`）而不是完整会话 ID

**解决方案**:
```bash
# 错误：使用 Slug
oho message list -s tidy-panda  # ❌

# 正确：使用完整会话 ID
oho session list  # 找到 ses_xxx
oho message list -s ses_34dbffe0dffe8SfdMTbL53MWFP  # ✅
```

**获取会话 ID**:
```bash
# 列出所有会话
oho session list

# JSON 格式（便于脚本处理）
oho session list --json | jq '.[].id'

# 获取第一个会话 ID
oho session list --json | jq -r '.[0].id'

# 按项目过滤
oho session list --project /mnt/d/fe/babylon3DWorld --json | jq -r '.[0].id'
```

---

### 5.2 会话被中止

**症状**:
```bash
Error: 会话已被中止
```

**检查方法**:
```bash
oho session get ses_xxx --json | jq '.status'
```

**解决方案**:
```bash
# 方法 1: 创建新会话
oho session create

# 方法 2: 分叉会话（保留历史）
oho session fork ses_xxx

# 方法 3: 继续未中止的会话
oho session list --status idle
```

---

### 5.3 会话状态监控脚本

```bash
#!/bin/bash
# 会话状态监控脚本

SESSION_ID="${1:-}"

if [ -z "$SESSION_ID" ]; then
    echo "用法：$0 <session_id>"
    echo "示例：$0 ses_34dbffe0dffe8SfdMTbL53MWFP"
    exit 1
fi

echo "监控会话：$SESSION_ID"
echo "按 Ctrl+C 停止"
echo ""

while true; do
    TIMESTAMP=$(date '+%H:%M:%S')
    
    # 获取会话状态
    STATUS=$(oho session get "$SESSION_ID" --json 2>/dev/null | jq -r '.status // "unknown"')
    
    # 获取最新消息
    LAST_MSG=$(oho message list -s "$SESSION_ID" --limit 1 2>/dev/null | head -3)
    
    echo "[$TIMESTAMP] 状态：$STATUS"
    echo "$LAST_MSG"
    echo "---"
    
    sleep 5
done
```

**使用方法**:
```bash
chmod +x monitor_session.sh
./monitor_session.sh ses_34dbffe0dffe8SfdMTbL53MWFP
```

---

## 6. 模型/提供商问题

### 6.1 检查可用模型

```bash
# 列出所有提供商和默认模型
oho config providers

# JSON 格式
oho config providers --json | jq '.'

# 输出示例：
# 可用提供商:
# 
# 默认模型:
#   minimax-cn: MiniMax-M2.5-highspeed
#   openrouter: google/gemini-3-pro-preview
#   alibaba-cn: tongyi-intent-detect-v3
#   opencode: big-pickle
#   google: gemini-3-pro-preview
#   minimax: MiniMax-M2.5-highspeed
#   deepseek: deepseek-reasoner
```

---

### 6.2 模型配置问题

**症状**: 消息提交后返回 `model not found` 或长时间无响应

**解决方案**:

```bash
# 1. 检查当前配置
oho config get --json

# 2. 如果没有配置模型，在 OpenCode 中配置
# 编辑 ~/.opencode/config.json 或通过 TUI 配置

# 3. 测试模型可用性
oho message add -s ses_xxx "测试" --no-reply

# 4. 如果仍然失败，尝试其他模型
# 在 OpenCode TUI 中切换模型
```

---

### 6.3 API Key 问题

**症状**: 返回 401 或 403 错误

**检查方法**:
```bash
# 查看 OpenCode 日志
# 在运行 opencode serve 的终端查看错误信息
```

**解决方案**:
```bash
# 1. 在 OpenCode 配置中更新 API Key
# 编辑 ~/.opencode/config.json

# 2. 重新认证提供商
oho provider auth

# 3. 对于 OAuth 提供商
oho provider oauth authorize <provider_id>
```

---

## 7. 文件附件问题

### 7.1 文件不存在

**症状**:
```bash
Error: 文件不存在：/path/to/file.txt
```

**解决方案**:
```bash
# 1. 检查文件路径
ls -la /path/to/file.txt

# 2. 使用绝对路径
oho message add -s ses_xxx "分析文件" \
    --file /absolute/path/to/file.txt

# 3. 检查文件权限
chmod +r /path/to/file.txt
```

---

### 7.2 文件过大

**症状**: 请求超时或返回 413 错误

**限制**:
- 单文件最大：10MB（取决于服务器配置）
- 总请求大小：取决于 HTTP 服务器配置

**解决方案**:
```bash
# 方法 1: 压缩文件
tar -czf archive.tar.gz ./large_directory/
oho message add -s ses_xxx "分析代码" --file archive.tar.gz

# 方法 2: 分批发送
find ./src -name "*.ts" | head -5 | while read file; do
    oho message add -s ses_xxx "分析 $file" --file "$file"
done

# 方法 3: 使用文件 URL（如果服务器支持）
# 在消息中引用文件路径而不是附件
```

---

### 7.3 文件类型不支持

**支持的 MIME 类型**:

| 类型 | 扩展名 |
|------|--------|
| 图片 | `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`, `.bmp`, `.svg` |
| 文档 | `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx` |
| 文本 | `.txt`, `.md`, `.html`, `.css`, `.js`, `.json`, `.xml`, `.yaml`, `.yml` |
| 代码 | `.py`, `.go`, `.java`, `.c`, `.cpp`, `.h`, `.rs`, `.ts`, `.tsx` |
| 其他 | `.zip`, `.tar`, `.gz`, `.mp3`, `.mp4`, `.wav` |

**检测 MIME 类型**:
```bash
file --mime-type /path/to/file
```

---

## 8. 性能问题

### 8.1 响应时间长

**正常响应时间**:
- 简单查询：1-5 秒
- 复杂任务：10-60 秒
- 代码分析：30-120 秒

**优化方法**:

```bash
# 方法 1: 使用更快的模型
oho message add -s ses_xxx "任务" \
    --model alibaba-cn/qwen3.5-plus  # 快速模型

# 方法 2: 异步提交
oho message prompt-async -s ses_xxx "长时间任务"

# 方法 3: 简化任务描述
# 避免过长的上下文
```

---

### 8.2 内存占用高

**检查方法**:
```bash
# 检查 OpenCode 进程内存
ps aux | grep opencode | awk '{print $2, $4}'

# 或使用 htop
htop -p $(pgrep opencode)
```

**解决方案**:
```bash
# 方法 1: 限制会话历史长度
# 在 opencode.json 中配置：
{
  "maxTokens": 4096
}

# 方法 2: 定期清理旧会话
oho session list --json | jq -r '.[].id' | head -5 | xargs -I {} oho session delete {}

# 方法 3: 重启服务器
pkill opencode
opencode serve
```

---

### 8.3 性能监控脚本

```bash
#!/bin/bash
# 性能监控脚本

echo "=== OpenCode 性能监控 ==="
echo ""

# 服务器响应时间
echo -n "服务器响应时间："
START=$(date +%s%N)
curl -s http://localhost:4096/global/health > /dev/null
END=$(date +%s%N)
ELAPSED=$(( (END - START) / 1000000 ))
echo "${ELAPSED}ms"

# 进程信息
echo ""
echo "进程信息:"
ps aux | grep "[o]pencode" | awk '{printf "  PID: %s, CPU: %s%%, MEM: %s%%\n", $2, $3, $4}'

# 会话数量
echo ""
echo -n "活跃会话数："
oho session list --json 2>/dev/null | jq 'length' || echo "N/A"

# 端口监听
echo ""
echo "端口监听:"
netstat -tlnp 2>/dev/null | grep 4096 || echo "未监听 4096 端口"
```

---

## 9. 诊断脚本

### 9.1 完整诊断脚本

```bash
#!/bin/bash
# oho CLI 完整诊断脚本
# 用于诊断消息提交后无响应的问题

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
SERVER_HOST="${OPENCODE_SERVER_HOST:-127.0.0.1}"
SERVER_PORT="${OPENCODE_SERVER_PORT:-4096}"
SERVER_PASSWORD="${OPENCODE_SERVER_PASSWORD:-}"
BASE_URL="http://${SERVER_HOST}:${SERVER_PORT}"

echo "========================================="
echo "oho CLI 诊断脚本"
echo "========================================="
echo "服务器地址：${BASE_URL}"
echo ""

# 1. 检查服务器健康状态
echo -e "${YELLOW}[1/7] 检查服务器健康状态...${NC}"
if curl -s -f "${BASE_URL}/global/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 服务器正常运行${NC}"
else
    echo -e "${RED}✗ 服务器无法连接${NC}"
    echo "解决方案：运行 'opencode serve --port 4096'"
    exit 1
fi
echo ""

# 2. 检查认证配置
echo -e "${YELLOW}[2/7] 检查认证配置...${NC}"
if [ -z "$SERVER_PASSWORD" ]; then
    echo -e "${RED}✗ 未设置 OPENCODE_SERVER_PASSWORD${NC}"
    echo "解决方案：export OPENCODE_SERVER_PASSWORD=your-password"
    exit 1
else
    echo -e "${GREEN}✓ 密码已配置${NC}"
fi

AUTH_RESULT=$(curl -s -w "%{http_code}" -u "opencode:${SERVER_PASSWORD}" "${BASE_URL}/config" 2>/dev/null)
HTTP_CODE="${AUTH_RESULT: -3}"
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ 认证成功${NC}"
else
    echo -e "${RED}✗ 认证失败 (HTTP ${HTTP_CODE})${NC}"
    echo "解决方案：检查密码是否正确"
    exit 1
fi
echo ""

# 3. 检查会话列表
echo -e "${YELLOW}[3/7] 获取会话列表...${NC}"
SESSIONS=$(curl -s -u "opencode:${SERVER_PASSWORD}" "${BASE_URL}/session")
SESSION_COUNT=$(echo "$SESSIONS" | grep -o '"id"' | wc -l)
echo "找到 ${SESSION_COUNT} 个会话"

if [ "$SESSION_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}! 创建新会话用于测试...${NC}"
    CREATE_RESP=$(curl -s -X POST -u "opencode:${SERVER_PASSWORD}" \
        -H "Content-Type: application/json" \
        -d '{"title":"diagnostic-session"}' \
        "${BASE_URL}/session")
    SESSION_ID=$(echo "$CREATE_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "创建会话：${SESSION_ID}"
else
    SESSION_ID=$(echo "$SESSIONS" | grep -o '"id":"ses_[^"]*"' | head -1 | cut -d'"' -f4)
    echo "使用会话：${SESSION_ID}"
fi
echo ""

# 4. 检查会话状态
echo -e "${YELLOW}[4/7] 检查会话状态...${NC}"
STATUS_RESP=$(curl -s -u "opencode:${SERVER_PASSWORD}" "${BASE_URL}/session/status")
if echo "$STATUS_RESP" | grep -q "\"${SESSION_ID}\""; then
    echo -e "${GREEN}✓ 会话状态正常${NC}"
else
    echo -e "${YELLOW}! 会话状态无法获取${NC}"
fi
echo ""

# 5. 检查可用模型
echo -e "${YELLOW}[5/7] 检查可用模型...${NC}"
PROVIDERS=$(curl -s -u "opencode:${SERVER_PASSWORD}" "${BASE_URL}/config/providers")
if echo "$PROVIDERS" | grep -q "default"; then
    echo -e "${GREEN}✓ 模型配置正常${NC}"
    echo "$PROVIDERS" | grep -A 5 "default"
else
    echo -e "${YELLOW}! 模型配置可能有问题${NC}"
fi
echo ""

# 6. 测试消息提交
echo -e "${YELLOW}[6/7] 测试消息提交...${NC}"
MESSAGE_REQ='{"parts": [{"type": "text", "text": "诊断测试，请回复 OK"}], "noReply": false}'

MSG_RESP=$(curl -s -X POST -u "opencode:${SERVER_PASSWORD}" \
    -H "Content-Type: application/json" \
    -d "$MESSAGE_REQ" \
    "${BASE_URL}/session/${SESSION_ID}/message")

if echo "$MSG_RESP" | grep -q '"id"'; then
    MSG_ID=$(echo "$MSG_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo -e "${GREEN}✓ 消息提交成功，ID: ${MSG_ID}${NC}"
else
    echo -e "${RED}✗ 消息提交失败${NC}"
    echo "响应：$MSG_RESP"
    exit 1
fi
echo ""

# 7. 检查 AI 响应
echo -e "${YELLOW}[7/7] 检查 AI 响应...${NC}"
sleep 5

MSG_HISTORY=$(curl -s -u "opencode:${SERVER_PASSWORD}" \
    "${BASE_URL}/session/${SESSION_ID}/message?limit=5")

if echo "$MSG_HISTORY" | grep -q '"role":"assistant"'; then
    echo -e "${GREEN}✓ 检测到 AI 响应${NC}"
else
    echo -e "${YELLOW}! 未检测到 AI 响应${NC}"
    echo ""
    echo "可能原因:"
    echo "  1. AI 模型配置问题"
    echo "  2. 会话被中止"
    echo "  3. 权限请求等待确认"
    echo "  4. 服务器日志错误"
fi
echo ""

echo "========================================="
echo "诊断完成"
echo "========================================="
```

**使用方法**:
```bash
chmod +x diagnose.sh
export OPENCODE_SERVER_PASSWORD=your-password
./diagnose.sh
```

---

### 9.2 快速检查脚本

```bash
#!/bin/bash
# 快速检查脚本（30 秒内完成）

echo "=== oho 快速检查 ==="

# 1. 服务器
curl -s http://localhost:4096/global/health > /dev/null && \
    echo "✓ 服务器运行" || echo "✗ 服务器未运行"

# 2. 认证
[ -n "$OPENCODE_SERVER_PASSWORD" ] && \
    echo "✓ 密码已配置" || echo "✗ 密码未配置"

# 3. 会话
oho session list > /dev/null 2>&1 && \
    echo "✓ 会话可访问" || echo "✗ 会话不可访问"

# 4. 模型
oho config providers > /dev/null 2>&1 && \
    echo "✓ 模型可用" || echo "✗ 模型不可用"

echo "=== 检查完成 ==="
```

---

## 10. 获取帮助

### 10.1 内置帮助

```bash
# 查看 oho 命令帮助
oho --help
oho message --help
oho session --help

# 查看特定命令帮助
oho message add --help
```

### 10.2 日志位置

```bash
# OpenCode 服务器日志
# 在运行 opencode serve 的终端查看

# oho CLI 调试输出
oho message add -s ses_xxx "test" 2>&1 | grep -i debug

# 配置文件位置
~/.config/oho/config.json
~/.opencode/config.json
```

### 10.3 反馈渠道

| 渠道 | 用途 |
|------|------|
| GitHub Issues | Bug 报告和功能请求 |
| 项目文档 | https://opencode.ai/docs/ |
| oho CLI 文档 | 本目录下各模块文档 |

---

## 附录：常见问题速查表

| 问题 | 快速命令 | 解决方案 |
|------|----------|----------|
| 服务器未运行 | `curl localhost:4096/global/health` | `opencode serve` |
| 认证失败 | `oho config get` | 设置 `OPENCODE_SERVER_PASSWORD` |
| 会话忙碌 | `oho session status` | 等待或 `oho session abort` |
| 模型不可用 | `oho config providers` | 配置提供商 API Key |
| 文件不存在 | `ls -la <path>` | 使用绝对路径 |
| 响应超时 | - | 使用 `prompt-async` |
| 会话 ID 无效 | `oho session list` | 使用 `ses_` 开头的 ID |

---

*文档维护：nanobot 🐈*  
*最后更新：2026-03-15*
