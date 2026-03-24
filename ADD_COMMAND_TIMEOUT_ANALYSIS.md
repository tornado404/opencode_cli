# oho add 命令超时问题分析报告

**分析时间**: 2026-03-24 10:30 (CST)  
**分析人**: nanobot 🐈

---

## 🔍 问题现象

用户反馈使用 `oho add` 命令时会超时。

---

## 📋 分析结果

### 1. 代码实现分析

#### add.go 执行流程

```go
func runAdd(cmd *cobra.Command, args []string) error {
    // Step 1: 获取当前工作目录
    sessionDir := addDirectory
    if sessionDir == "" {
        sessionDir, _ = os.Getwd()
    }

    // Step 2: 生成标题
    sessionTitle := addTitle
    if sessionTitle == "" {
        sessionTitle = fmt.Sprintf("New session - %s", time.Now().Format("2006-01-02T15:04:05"))
    }

    // Step 3: 创建会话 ✅ (快速，通常 < 1 秒)
    sessionID, err := createSession(c, ctx, sessionTitle, addParent, sessionDir)

    // Step 4: 发送消息 ⚠️ (可能超时)
    messageID, err := sendMessage(c, ctx, sessionID, message, addAgent, addModel, addNoReply, ...)
    
    // Step 5: 输出结果
    ...
}
```

#### sendMessage 函数

```go
func sendMessage(...) (string, error) {
    // 构建消息请求
    msgReq := types.MessageRequest{
        Model:   convertModel(model),
        Agent:   agent,
        NoReply: noReply,  // ⚠️ 关键：默认 false，会等待 AI 响应
        System:  system,
        Tools:   tools,
        Parts:   parts,
    }

    // 发送 POST 请求
    resp, err := c.Post(ctx, fmt.Sprintf("/session/%s/message", sessionID), msgReq)
    // ⚠️ 这里会等待 AI 响应，可能耗时很长
    ...
}
```

#### client.go 超时配置

```go
func NewClient() *Client {
    timeoutSec := 300 // 5 分钟
    if envTimeout := os.Getenv("OPENCODE_CLIENT_TIMEOUT"); envTimeout != "" {
        if parsed, err := strconv.Atoi(envTimeout); err == nil && parsed > 0 {
            timeoutSec = parsed
        }
    }

    return &Client{
        httpClient: &http.Client{
            Timeout: time.Duration(timeoutSec) * time.Second,
        },
    }
}
```

---

## 🐛 问题根源

### 问题 1: 默认等待 AI 响应

**代码位置**: `add.go:sendMessage()`

```go
NoReply: noReply,  // 默认值：false
```

**问题**:
- `addNoReply` 默认值为 `false`
- 这意味着 `oho add` 默认会**等待 AI 响应后才返回**
- 对于复杂的 AI 任务（分析代码、重构等），5 分钟可能不够

**对比 `oho message add`**:
```bash
# oho message add 也需要 --no-reply 才不等待
oho message add -s <session> "任务"           # 等待响应
oho message add -s <session> "任务" --no-reply  # 不等待
```

---

### 问题 2: 文档未说明超时配置

**文档位置**: `README.md`, `oho/README.md`

**缺失内容**:
1. ❌ 未说明 `add` 命令默认会等待 AI 响应
2. ❌ 未说明超时时间为 5 分钟
3. ❌ 未说明如何通过环境变量调整超时
4. ❌ 未建议使用 `--no-reply` 避免超时

**现有文档**:
```markdown
### Quick Start (Session + Message)

```bash
oho add "帮我分析这个项目"                    # Create session and send message
oho add "修复登录 bug" --title "Bug 修复"       # Create session with custom title
oho add "测试功能" --no-reply --agent default  # Don't wait for AI response
```
```

**问题**: 虽然示例中有 `--no-reply`，但**没有解释为什么需要这个参数**。

---

### 问题 3: 超时错误提示不友好

**代码位置**: `client.go:Request()`

```go
resp, err := c.httpClient.Do(req)
if err != nil {
    return nil, fmt.Errorf("请求失败：%w", err)
}
```

**问题**:
- 超时错误只显示 "请求失败：context deadline exceeded"
- 用户不知道可以调整超时时间
- 没有建议解决方案

---

## ✅ 解决方案

### 方案 1: 修改文档（推荐，立即生效）

**修改文件**: `oho/README.md`, `README_zh.md`

**添加内容**:

```markdown
### ⚠️ 超时注意事项

`oho add` 命令默认会等待 AI 响应后返回。对于复杂任务，AI 可能需要较长时间思考，可能导致超时。

**避免超时的方法**:

1. **使用 `--no-reply` 参数** (推荐):
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

---

### 方案 2: 修改代码默认行为

**修改文件**: `oho/cmd/add/add.go`

**选项 A: 默认使用 `--no-reply`**

```go
var (
    addNoReply bool = true  // 改为默认 true
)
```

**优点**: 
- 避免大多数超时问题
- 符合 CLI 工具的最佳实践（快速返回）

**缺点**:
- 破坏向后兼容性
- 用户可能期望立即看到 AI 响应

**选项 B: 添加超时提示**

```go
func runAdd(cmd *cobra.Command, args []string) error {
    ...
    
    // 发送消息
    if !addNoReply {
        fmt.Println("⏳ 等待 AI 响应（按 Ctrl+C 中断，或使用 --no-reply 避免等待）...")
    }
    
    messageID, err := sendMessage(...)
    ...
}
```

---

### 方案 3: 改进错误提示

**修改文件**: `oho/internal/client/client.go`

```go
resp, err := c.httpClient.Do(req)
if err != nil {
    if strings.Contains(err.Error(), "context deadline exceeded") {
        return nil, fmt.Errorf("请求超时（%d 秒）\n\n建议:\n  1. 使用 --no-reply 参数避免等待\n  2. 设置环境变量增加超时：export OPENCODE_CLIENT_TIMEOUT=600\n  3. 使用异步命令：oho message prompt-async", timeoutSec)
    }
    return nil, fmt.Errorf("请求失败：%w", err)
}
```

---

### 方案 4: 添加超时配置标志

**修改文件**: `oho/cmd/add/add.go`

```go
var (
    addTimeout int  // 新增超时标志
)

func init() {
    ...
    Cmd.Flags().IntVar(&addTimeout, "timeout", 0, "请求超时时间（秒），0 使用默认值")
}

func runAdd(cmd *cobra.Command, args []string) error {
    ...
    
    // 如果指定了超时，临时覆盖
    if addTimeout > 0 {
        os.Setenv("OPENCODE_CLIENT_TIMEOUT", strconv.Itoa(addTimeout))
        c = client.NewClient()  // 重新创建客户端
    }
    
    ...
}
```

**使用方式**:
```bash
oho add "复杂任务" --timeout 600
```

---

## 📝 推荐修复顺序

### 立即修复（今天）
1. ✅ **修改文档** - 添加超时说明和解决方案
2. ✅ **改进错误提示** - 超时错误给出明确建议

### 短期修复（本周）
3. ✅ **添加超时标志** - `--timeout` 参数
4. ✅ **添加等待提示** - 显示 "等待 AI 响应..."

### 长期考虑（下个版本）
5. ⏳ **讨论默认行为** - 是否默认 `--no-reply`
6. ⏳ **添加进度显示** - 显示 AI 思考进度

---

## 🔧 快速修复脚本

### 修复文档

```bash
# 1. 更新 oho/README.md
# 在 "Quick Start (Session + Message)" 部分后添加超时说明

# 2. 更新 README_zh.md
# 同步中文文档
```

### 修复错误提示

```bash
# 修改 oho/internal/client/client.go
# 在 Request() 函数中添加超时错误处理
```

---

## 📊 影响评估

| 修复方案 | 影响范围 | 风险 | 优先级 |
|----------|---------|------|--------|
| 修改文档 | 用户文档 | 低 | P0 |
| 改进错误提示 | 用户体验 | 低 | P0 |
| 添加超时标志 | 新增功能 | 低 | P1 |
| 修改默认行为 | 破坏性变更 | 中 | P2 |

---

## 🎯 结论

**主要原因**: 
1. ✅ **文档未说明** - 用户不知道 `add` 默认会等待 AI 响应
2. ✅ **文档未说明超时配置** - 用户不知道可以调整超时时间
3. ⚠️ 代码设计合理，但缺少用户友好的提示

**建议修复**:
1. **立即更新文档**，说明超时配置和 `--no-reply` 用法
2. **改进错误提示**，超时给出明确解决方案
3. **考虑添加 `--timeout` 参数**，方便临时调整

**代码本身没有问题**，主要是**文档和用户引导不足**。

---

*报告生成时间：2026-03-24 10:30 CST*
