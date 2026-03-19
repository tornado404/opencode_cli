# 任务：no-reply 模式下返回 sessionID 和 messageID

## 问题描述

当前使用 `oho message add --no-reply` 发送消息时：
1. 服务器不等待 AI 响应，返回空响应
2. CLI 只打印 "消息已发送"，**不输出 sessionID 和 messageID**
3. 用户无法追踪任务状态，交互体验差

## 需求

优化 `--no-reply` 模式的交互体验，确保执行后至少能获得：
- ✅ **sessionID**（从命令行参数已知）
- ✅ **messageID**（用户提供的或服务器生成的）
- ✅ **执行状态确认**

## 实现方案

### 方案 A：CLI 端优化（优先实施）

修改 `/mnt/d/fe/opencode_cli/oho/cmd/message/message.go` 的 `addCmd`：

```go
// 当前代码（问题）
if len(resp) == 0 {
    fmt.Println("消息已发送")
    return nil
}

// 优化后
if len(resp) == 0 {
    // no-reply 模式，服务器返回空响应
    fmt.Println("消息已发送 (no-reply 模式)")
    fmt.Printf("  Session ID: %s\n", sessionID)
    if messageID != "" {
        fmt.Printf("  Message ID: %s\n", messageID)
    } else {
        fmt.Println("  Message ID: (服务器生成，可通过 session 消息列表查询)")
    }
    fmt.Println("\n提示：使用 oho message list -s <session> 查看消息状态")
    return nil
}
```

### 方案 B：服务器端优化（可选，增强版）

修改 OpenCode Server 的 `/session/:id/message` 接口：
- 当 `noReply=true` 时，立即返回已创建的消息对象
- 返回字段：`{ id, sessionId, role: "user", createdAt, parts: [...] }`
- 不等待 AI 响应，但确认消息已入队

## 验收标准

1. ✅ 执行 `oho message add -s xxx "消息" --no-reply` 后，输出包含：
   - Session ID
   - Message ID（如果提供了 `--message` 参数）
   - 状态确认信息

2. ✅ 输出格式清晰，便于脚本解析（支持 `--json` 模式）

3. ✅ 不影响正常模式（无 `--no-reply`）的行为

## 相关文件

- `/mnt/d/fe/opencode_cli/oho/cmd/message/message.go` - 主要修改位置
- `/mnt/d/fe/opencode_cli/oho/internal/types/types.go` - 类型定义（可能需要新增响应类型）

## 测试用例

```bash
# 测试 1: 无 messageID
oho message add -s test123 "测试消息" --no-reply
# 期望输出：Session ID: test123

# 测试 2: 有 messageID
oho message add -s test123 "测试消息" --no-reply --message my-task-001
# 期望输出：Session ID: test123, Message ID: my-task-001

# 测试 3: JSON 模式
oho message add -s test123 "测试消息" --no-reply --json
# 期望输出：{"sessionId":"test123","messageId":"...","status":"queued"}
```

## 优先级

- P0: CLI 端优化（方案 A）- 立即实施
- P1: 服务器端优化（方案 B）- 后续增强
