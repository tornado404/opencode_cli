# oho CLI 操作指南 - 模块 7: 中断任务

> **适用版本**: oho CLI v1.0+  
> **最后更新**: 2026-03-03  
> **作者**: nanobot 🐈  
> **前置模块**: [模块 6: 指定 session_id 和模型发消息](./06-send-message.md)

---

## 📋 目录

1. [任务中断概述](#1-任务中断概述)
2. [中止正在运行的会话](#2-中止正在运行的会话)
3. [检查会话状态](#3-检查会话状态)
4. [中断后恢复](#4-中断后恢复)
5. [批量中断管理](#5-批量中断管理)
6. [中断策略与最佳实践](#6-中断策略与最佳实践)

---

## 1. 任务中断概述

### 1.1 为什么需要中断任务

**常见场景**:
- ❌ AI 陷入死循环或重复输出
- ❌ 任务方向错误，需要重新指导
- ❌ 发现更好的解决方案
- ❌ 资源占用过高，需要释放
- ❌ 紧急任务需要优先处理
- ❌ API 错误导致卡住

---

### 1.2 中断方式对比

| 方式 | 命令 | 效果 | 适用场景 |
|------|------|------|----------|
| 中止会话 | `oho session abort` | 立即停止 | 紧急中断 |
| 不等待响应 | `--no-reply` | 后台执行 | 长时间任务 |
| 异步提交 | `prompt-async` | 非阻塞 | 批量任务 |
| 回退消息 | `oho session revert` | 撤销操作 | 错误修改 |

---

### 1.3 中断的影响

**中断后会话状态**:
```
运行中 (running) → 已中止 (aborted)
```

**数据保留**:
- ✅ 已生成的消息保留
- ✅ 会话历史完整
- ✅ 文件修改保留（已写入的）
- ❌ 当前正在执行的操作丢失
- ❌ 未完成的思考过程丢失

---

## 2. 中止正在运行的会话

### 2.1 基本用法

```bash
# 中止指定会话
oho session abort -s tidy-panda

# 使用会话 ID
oho session abort -s ses_352a39c7bffe7RQv3VaA7Kypgs

# JSON 格式输出
oho session abort -s tidy-panda --json
```

**预期输出**:
```
✓ 会话 tidy-panda 已中止
状态：aborted
```

---

### 2.2 中止所有运行中的会话

```bash
#!/bin/bash
# 批量中止所有运行中的会话

running_sessions=$(oho session status --json | jq -r '.[] | select(.status == "running") | .id')

for session_id in $running_sessions; do
    echo "中止会话：$session_id"
    oho session abort -s "$session_id"
done

echo "所有运行中的会话已中止"
```

---

### 2.3 条件中断

```bash
#!/bin/bash
# 根据条件中断会话

# 中断运行超过 30 分钟的会话
timeout_minutes=30

oho session status --json | jq -r '.[] | select(.status == "running") | @base64' | \
while read -r line; do
    decoded=$(echo "$line" | base64 -d)
    session_id=$(echo "$decoded" | jq -r '.id')
    updated=$(echo "$decoded" | jq -r '.updated')
    
    # 计算运行时间
    now=$(date -u +%s)
    last_update=$(date -d "$updated" +%s 2>/dev/null || echo 0)
    running_minutes=$(( (now - last_update) / 60 ))
    
    if [ $running_minutes -gt $timeout_minutes ]; then
        echo "⚠️  会话 $session_id 运行 ${running_minutes} 分钟，执行中断"
        oho session abort -s "$session_id"
    else
        echo "✓ 会话 $session_id 运行 ${running_minutes} 分钟，正常"
    fi
done
```

---

### 2.4 中断前保存状态

```bash
#!/bin/bash
# 中断前保存会话状态

save_session_state() {
    local session=$1
    
    # 导出会话消息
    oho message list -s "$session" --json > "/tmp/session_${session}_messages.json"
    
    # 导出会话详情
    oho session get -s "$session" --json > "/tmp/session_${session}_info.json"
    
    echo "会话状态已保存到 /tmp/session_${session}_*.json"
}

# 使用
save_session_state "tidy-panda"

# 然后中断
oho session abort -s "tidy-panda"
```

---

### 2.5 优雅中断

```bash
#!/bin/bash
# 优雅中断：等待当前操作完成后不再继续

graceful_abort() {
    local session=$1
    local wait_seconds=${2:-30}
    
    echo "等待 ${wait_seconds} 秒后中断会话 $session..."
    sleep $wait_seconds
    
    # 检查会话是否仍在运行
    status=$(oho session get -s "$session" --json 2>/dev/null | jq -r '.status')
    
    if [ "$status" == "running" ]; then
        echo "会话仍在运行，执行中断"
        oho session abort -s "$session"
    else
        echo "会话已结束，状态：$status"
    fi
}

# 使用
graceful_abort "tidy-panda" 60
```

---

## 3. 检查会话状态

### 3.1 查看所有会话状态

```bash
# 列出所有会话状态
oho session status

# JSON 格式
oho session status --json
```

**预期输出**:
```
会话状态:
  tidy-panda       running    (2 分钟前)
  hidden-sailor    completed  (1 小时前)
  shiny-squid      aborted    (30 分钟前)
```

**JSON 输出**:
```json
[
  {
    "id": "ses_xxxxx",
    "slug": "tidy-panda",
    "status": "running",
    "updated": "2026-03-03T10:00:00Z",
    "messageCount": 15
  }
]
```

---

### 3.2 查看单个会话状态

```bash
# 获取会话详情
oho session get -s tidy-panda

# JSON 格式
oho session get -s tidy-panda --json | jq '.status'
```

**状态类型**:
| 状态 | 说明 | 可操作 |
|------|------|--------|
| `running` | 正在执行 | 可中断 |
| `completed` | 已完成 | 可继续发消息 |
| `aborted` | 已中止 | 可恢复或新建 |
| `error` | 发生错误 | 需检查错误信息 |
| `idle` | 空闲 | 可发送新消息 |

---

### 3.3 监控会话状态

```bash
#!/bin/bash
# 实时监控会话状态

monitor_session() {
    local session=$1
    
    while true; do
        clear
        echo "=== 会话监控：$session ==="
        echo "时间：$(date)"
        echo ""
        
        status=$(oho session get -s "$session" --json 2>/dev/null | jq -r '.status')
        updated=$(oho session get -s "$session" --json 2>/dev/null | jq -r '.updated')
        msg_count=$(oho session get -s "$session" --json 2>/dev/null | jq -r '.messageCount')
        
        echo "状态：$status"
        echo "最后更新：$updated"
        echo "消息数：$msg_count"
        
        if [ "$status" != "running" ]; then
            echo ""
            echo "会话已结束"
            break
        fi
        
        sleep 5
    done
}

# 使用
monitor_session "tidy-panda"
```

---

### 3.4 状态过滤

```bash
# 只查看运行中的会话
oho session status --json | jq '.[] | select(.status == "running")'

# 只查看已中止的会话
oho session status --json | jq '.[] | select(.status == "aborted")'

# 只查看错误的会话
oho session status --json | jq '.[] | select(.status == "error")'

# 统计各状态数量
oho session status --json | jq 'group_by(.status) | map({status: .[0].status, count: length})'
```

---

## 4. 中断后恢复

### 4.1 从中断点继续

```bash
# 查看中止前的最后消息
oho message list -s tidy-panda --tail 5

# 继续任务
oho message add -s tidy-panda "继续之前的任务，但请简化方案"
```

---

### 4.2 恢复会话的最佳实践

```bash
#!/bin/bash
# 中断后恢复会话

recover_session() {
    local session=$1
    local context=$2
    
    # 1. 检查会话状态
    status=$(oho session get -s "$session" --json | jq -r '.status')
    
    if [ "$status" != "aborted" ] && [ "$status" != "error" ]; then
        echo "会话状态正常：$status，无需恢复"
        return 0
    fi
    
    # 2. 查看最后消息
    echo "=== 最后 3 条消息 ==="
    oho message list -s "$session" --tail 3
    
    # 3. 发送恢复消息
    echo ""
    echo "发送恢复消息..."
    oho message add -s "$session" \
      --model alibaba-cn/qwen3.5-plus \
      "我们继续之前的任务。$context"
}

# 使用
recover_session "tidy-panda" "请从代码分析部分继续"
```

---

### 4.3 创建新会话替代

```bash
# 如果原会话问题太多，创建新会话
oho message add -s new-task-slug \
  "重新开始任务：分析项目结构。之前在中断的会话 tidy-panda 中已经讨论了部分，这里是新的开始。"
```

---

### 4.4 回退错误操作

```bash
# 回退最后一条消息（如果 AI 做了错误的文件修改）
oho session revert -s tidy-panda

# 回退指定数量的消息
oho session revert -s tidy-panda --count 3

# 恢复已回退的消息
oho session unrevert -s tidy-panda
```

**注意**:
- ⚠️ `revert` 只回退消息，不撤销文件修改
- ⚠️ 文件修改需要手动恢复（使用 Git）
- ✅ 适合撤销错误的指令

---

## 5. 批量中断管理

### 5.1 清理所有运行中的会话

```bash
#!/bin/bash
# 清理所有运行中的会话

echo "=== 清理运行中的会话 ==="

# 获取所有运行中的会话
running=$(oho session status --json | jq -r '.[] | select(.status == "running") | .slug')

if [ -z "$running" ]; then
    echo "✓ 没有运行中的会话"
    exit 0
fi

count=0
for session in $running; do
    echo "中止：$session"
    oho session abort -s "$session"
    ((count++))
    sleep 1  # 避免限流
done

echo "✓ 已中止 $count 个会话"
```

---

### 5.2 按项目清理

```bash
#!/bin/bash
# 按项目清理会话

cleanup_by_project() {
    local project_path=$1
    
    # 获取项目相关的会话
    sessions=$(oho session list --project "$project_path" --json | jq -r '.[] | select(.status == "running") | .id')
    
    for session_id in $sessions; do
        echo "中止会话：$session_id (项目：$project_path)"
        oho session abort -s "$session_id"
    done
}

# 使用
cleanup_by_project "/mnt/d/fe/opencode_cli"
cleanup_by_project "/mnt/d/fe/babylon3DWorld"
```

---

### 5.3 定时清理

```bash
#!/bin/bash
# 定时清理长时间运行的会话

# 添加到 crontab
# */30 * * * * /path/to/cleanup_script.sh

MAX_RUNNING_MINUTES=60

oho session status --json | jq -r '.[] | select(.status == "running") | @base64' | \
while read -r line; do
    decoded=$(echo "$line" | base64 -d)
    session_id=$(echo "$decoded" | jq -r '.id')
    slug=$(echo "$decoded" | jq -r '.slug')
    updated=$(echo "$decoded" | jq -r '.updated')
    
    now=$(date -u +%s)
    last_update=$(date -d "$updated" +%s 2>/dev/null || echo 0)
    running_minutes=$(( (now - last_update) / 60 ))
    
    if [ $running_minutes -gt $MAX_RUNNING_MINUTES ]; then
        echo "[$(date)] 中止超时会话：$slug (${running_minutes}分钟)"
        oho session abort -s "$session_id"
    fi
done
```

---

### 5.4 会话健康检查脚本

```bash
#!/bin/bash
# 完整的会话健康检查

echo "=== OpenCode 会话健康检查 ==="
echo "时间：$(date)"
echo ""

# 统计
total=$(oho session status --json | jq 'length')
running=$(oho session status --json | jq '[.[] | select(.status == "running")] | length')
aborted=$(oho session status --json | jq '[.[] | select(.status == "aborted")] | length')
error=$(oho session status --json | jq '[.[] | select(.status == "error")] | length')

echo "总会话数：$total"
echo "运行中：$running"
echo "已中止：$aborted"
echo "错误：$error"
echo ""

# 显示运行中的会话
if [ $running -gt 0 ]; then
    echo "=== 运行中的会话 ==="
    oho session status --json | jq -r '.[] | select(.status == "running") | "\(.slug) - 更新于 \(.updated)"'
fi

# 显示错误的会话
if [ $error -gt 0 ]; then
    echo ""
    echo "=== 错误的会话 ==="
    oho session status --json | jq -r '.[] | select(.status == "error") | "\(.slug) - \(.error)"'
fi

# 建议
if [ $running -gt 5 ]; then
    echo ""
    echo "⚠️  警告：运行中的会话过多，建议清理"
    echo "运行：oho session abort -s <session>"
fi
```

---

## 6. 中断策略与最佳实践

### 6.1 何时应该中断

**立即中断**:
- ❌ AI 重复输出相同内容
- ❌ 明显偏离任务目标
- ❌ 产生大量无意义输出
- ❌ API 错误导致卡住

**等待完成**:
- ✅ 长时间但正常的分析
- ✅ 大文件处理
- ✅ 复杂代码生成
- ✅ 多步骤任务执行中

---

### 6.2 中断前的检查清单

```bash
# 中断前检查清单
pre_abort_checklist() {
    local session=$1
    
    echo "=== 中断前检查：$session ==="
    
    # 1. 查看当前状态
    status=$(oho session get -s "$session" --json | jq -r '.status')
    echo "1. 当前状态：$status"
    
    # 2. 查看最后消息
    echo "2. 最后消息:"
    oho message list -s "$session" --tail 2 --json | jq -r '.[].content' | head -50
    
    # 3. 检查是否有文件修改
    echo "3. 检查 Git 状态..."
    git status --short
    
    # 4. 确认是否真的需要中断
    echo ""
    read -p "确认中断？(y/N): " confirm
    if [ "$confirm" != "y" ]; then
        echo "取消中断"
        return 1
    fi
    
    return 0
}
```

---

### 6.3 中断后的跟进

```bash
#!/bin/bash
# 中断后跟进流程

post_abort_followup() {
    local session=$1
    local reason=$2
    
    echo "=== 中断后跟进 ==="
    echo "会话：$session"
    echo "原因：$reason"
    echo ""
    
    # 1. 记录中断原因
    echo "[$(date)] 中断会话 $session: $reason" >> ~/.oho/interrupt_log.md
    
    # 2. 分析是否需要调整策略
    echo "建议:"
    case $reason in
        "loop")
            echo "- 下次使用更明确的指令"
            echo "- 考虑分步骤执行"
            ;;
        "wrong_direction")
            echo "- 提供更详细的上下文"
            echo "- 使用系统提示约束方向"
            ;;
        "api_error")
            echo "- 检查网络连接"
            echo "- 尝试更换模型"
            ;;
        *)
            echo "- 审查任务描述是否清晰"
            ;;
    esac
    
    # 3. 决定是否继续
    echo ""
    read -p "是否创建新会话继续？(y/N): " continue_task
    if [ "$continue_task" == "y" ]; then
        oho message add -s "${session}-retry" "重新开始：继续之前的任务，但改进方法"
    fi
}
```

---

### 6.4 避免频繁中断

**最佳实践**:
1. ✅ 任务开始前明确目标
2. ✅ 使用系统提示约束行为
3. ✅ 分步骤执行复杂任务
4. ✅ 设置合理的期望
5. ✅ 给 AI 足够时间思考

**避免**:
- ❌ 稍有不满意就中断
- ❌ 频繁切换任务方向
- ❌ 不清晰的指令导致反复

---

### 6.5 中断日志

```bash
# 记录中断历史
log_interrupt() {
    local session=$1
    local reason=$2
    local timestamp=$(date -Iseconds)
    
    echo "- [$timestamp] **$session**: $reason" >> ~/.oho/interrupt_log.md
}

# 查看中断历史
view_interrupt_log() {
    tail -20 ~/.oho/interrupt_log.md
}

# 使用
log_interrupt "tidy-panda" "AI 陷入死循环，重复输出相同代码"
```

---

## 🔧 实用技巧

### 技巧 1: 一键清理

```bash
# 添加到 ~/.bashrc
alias oho-cleanup="oho session status --json | jq -r '.[] | select(.status == \"running\") | .slug' | xargs -I {} oho session abort -s {}"
```

---

### 技巧 2: 中断通知

```bash
# 中断后发送通知
oho session abort -s tidy-panda && \
  notify-send "OpenCode" "会话 tidy-panda 已中断"
```

---

### 技巧 3: 自动重试

```bash
# 中断后自动重试（最多 3 次）
retry_with_abort() {
    local session=$1
    local message=$2
    local max_retries=3
    local retry=0
    
    while [ $retry -lt $max_retries ]; do
        echo "尝试 $((retry+1))/$max_retries"
        
        # 发送消息并设置超时
        timeout 300 oho message add -s "$session" "$message"
        
        if [ $? -eq 0 ]; then
            echo "✓ 成功"
            return 0
        else
            echo "⚠️  超时或错误，中断并重试"
            oho session abort -s "$session"
            ((retry++))
            sleep 5
        fi
    done
    
    echo "❌ 所有重试失败"
    return 1
}
```

---

### 技巧 4: 会话超时自动中断

```bash
# 设置会话超时
SESSION_TIMEOUT=1800  # 30 分钟

check_session_timeout() {
    oho session status --json | jq -r '.[] | select(.status == "running") | @base64' | \
    while read -r line; do
        decoded=$(echo "$line" | base64 -d)
        session_id=$(echo "$decoded" | jq -r '.id')
        updated=$(echo "$decoded" | jq -r '.updated')
        
        now=$(date -u +%s)
        last_update=$(date -d "$updated" +%s 2>/dev/null || echo 0)
        
        if [ $((now - last_update)) -gt $SESSION_TIMEOUT ]; then
            echo "会话 $session_id 超时，自动中断"
            oho session abort -s "$session_id"
        fi
    done
}
```

---

## 📝 检查清单

在中断任务前，请确认：

- [ ] 会话确实需要中断（不是正常执行中）
- [ ] 已保存重要进度
- [ ] 记录了中断原因
- [ ] 考虑了替代方案（如等待完成）
- [ ] 准备好恢复或重试计划

---

## 🔗 相关文档

- [模块 5: 指定工作区提交任务](./05-submit-task.md) - 任务提交
- [模块 6: 指定 session_id 和模型发消息](./06-send-message.md) - 消息控制
- [模块 8: 查询状态](./08-query-status.md) - 状态监控

---

## 🆘 常见问题

### Q1: 中断后会话还能继续使用吗？

**A**: 可以，中断后会话状态变为 `aborted`，可以发送新消息继续：
```bash
oho session abort -s slug
oho message add -s slug "继续任务"
```

---

### Q2: 中断会撤销文件修改吗？

**A**: 不会，中断只停止 AI 执行，已写入的文件修改保留。使用 Git 恢复：
```bash
git checkout -- path/to/file
```

---

### Q3: 如何防止 AI 卡住？

**A**:
```bash
# 使用异步提交
oho message prompt-async -s slug "任务"

# 设置超时
timeout 300 oho message add -s slug "任务"

# 使用 --no-reply
oho message add -s slug "任务" --no-reply
```

---

### Q4: 中断多个会话有影响吗？

**A**: 可以同时中断多个会话，但建议：
- 逐个中断，避免 API 限流
- 记录每个会话的中断原因
- 清理后检查系统资源

---

*文档生成时间：2026-03-03 11:17 CST*
