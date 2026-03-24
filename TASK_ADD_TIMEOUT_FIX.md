# oho add 超时问题修复任务

## 任务目标
修复 oho add 命令超时问题，主要是文档和用户引导不足。

## 需要完成的工作

### 1. 更新文档 (P0)
**文件**: `oho/README.md` 和 `README_zh.md`

在 Quick Start 部分后添加超时说明：

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

### 2. 改进错误提示 (P0)
**文件**: `oho/internal/client/client.go`

在 `Request()` 函数中改进超时错误提示：

```go
resp, err := c.httpClient.Do(req)
if err != nil {
    if strings.Contains(err.Error(), "context deadline exceeded") {
        return nil, fmt.Errorf("请求超时（%d 秒）\n\n建议:\n  1. 使用 --no-reply 参数避免等待\n  2. 设置环境变量增加超时：export OPENCODE_CLIENT_TIMEOUT=600\n  3. 使用异步命令：oho message prompt-async", timeoutSec)
    }
    return nil, fmt.Errorf("请求失败：%w", err)
}
```

### 3. 添加 --timeout 参数 (P1)
**文件**: `oho/cmd/add/add.go`

添加超时标志：

```go
var (
    addTimeout int  // 新增
)

func init() {
    Cmd.Flags().IntVar(&addTimeout, "timeout", 0, "请求超时时间（秒），0 使用默认值")
}

func runAdd(cmd *cobra.Command, args []string) error {
    // 如果指定了超时，临时覆盖
    if addTimeout > 0 {
        os.Setenv("OPENCODE_CLIENT_TIMEOUT", strconv.Itoa(addTimeout))
        c = client.NewClient()  // 重新创建客户端
    }
    ...
}
```

## 验收标准
- [ ] README.md 和 README_zh.md 都有超时说明
- [ ] 超时错误提示包含解决方案
- [ ] --timeout 参数可用
- [ ] 所有测试通过

## 参考资料
- 分析报告：`/mnt/d/fe/opencode_cli/ADD_COMMAND_TIMEOUT_ANALYSIS.md`

任务完成后中止会话。
