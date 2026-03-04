# oho CLI 操作指南 - 模块 8: 查询任务执行情况及状态

> **适用版本**: oho CLI v1.0+  
> **最后更新**: 2026-03-03  
> **作者**: nanobot 🐈  
> **前置模块**: [模块 7: 中断任务](./07-interrupt-task.md)

---

## 📋 目录

1. [任务状态查询概述](#1-任务状态查询概述)
2. [查看会话详情](#2-查看会话详情)
3. [消息历史查询](#3-消息历史查询)
4. [待办事项追踪](#4-待办事项追踪)
5. [任务进度监控](#5-任务进度监控)
6. [会话总结与导出](#6-会话总结与导出)
7. [高级查询技巧](#7-高级查询技巧)

---

## 1. 任务状态查询概述

### 1.1 查询维度

| 维度 | 命令 | 用途 |
|------|------|------|
| 会话状态 | `oho session status` | 查看运行/完成/错误状态 |
| 会话详情 | `oho session get` | 获取完整会话信息 |
| 消息列表 | `oho message list` | 查看对话历史 |
| 消息详情 | `oho message get` | 获取单条消息内容 |
| 待办事项 | `oho session todo` | 追踪任务进度 |
| 会话总结 | `oho session summarize` | 生成会话摘要 |

---

### 1.2 状态类型详解

```json
{
  "status": "running",    // 正在执行
  "status": "completed",  // 已完成
  "status": "aborted",    // 已中止
  "status": "error",      // 发生错误
  "status": "idle"        // 空闲
}
```

**状态转换**:
```
新建 → running → completed
            ↓
        aborted
            ↓
         error
```

---

### 1.3 查询频率建议

| 场景 | 建议频率 |
|------|----------|
| 长时间任务监控 | 每 30-60 秒 |
| 日常检查 | 每 5-10 分钟 |
| 批量任务 | 每 2-5 分钟 |
| 错误排查 | 立即查询 |

---

## 2. 查看会话详情

### 2.1 基本用法

```bash
# 获取会话详情
oho session get -s tidy-panda

# JSON 格式
oho session get -s tidy-panda --json
```

**预期输出**:
```json
{
  "id": "ses_352a39c7bffe7RQv3VaA7Kypgs",
  "slug": "tidy-panda",
  "title": "oho CLI 文档完善任务",
  "status": "completed",
  "createdAt": "2026-03-02T18:00:00Z",
  "updatedAt": "2026-03-03T11:00:00Z",
  "messageCount": 45,
  "model": "alibaba-cn/qwen3.5-plus",
  "project": {
    "path": "/mnt/d/fe/opencode_cli",
    "name": "opencode_cli"
  }
}
```

---

### 2.2 提取关键信息

```bash
# 提取状态
oho session get -s tidy-panda --json | jq -r '.status'

# 提取消息数
oho session get -s tidy-panda --json | jq -r '.messageCount'

# 提取最后更新时间
oho session get -s tidy-panda --json | jq -r '.updatedAt'

# 提取模型
oho session get -s tidy-panda --json | jq -r '.model'

# 提取项目路径
oho session get -s tidy-panda --json | jq -r '.project.path'
```

---

### 2.3 批量获取会话信息

```bash
#!/bin/bash
# 批量获取所有会话信息

oho session list --json | jq -r '.[] | @base64' | \
while read -r line; do
    decoded=$(echo "$line" | base64 -d)
    slug=$(echo "$decoded" | jq -r '.slug')
    status=$(echo "$decoded" | jq -r '.status')
    msg_count=$(echo "$decoded" | jq -r '.messageCount')
    updated=$(echo "$decoded" | jq -r '.updated')
    
    printf "%-20s %-10s %3d 条消息  %s\n" "$slug" "$status" "$msg_count" "$updated"
done
```

---

### 2.4 按条件过滤

```bash
# 只查看特定项目的会话
oho session list --project /mnt/d/fe/opencode_cli --json

# 只查看运行中的会话
oho session status --json | jq '.[] | select(.status == "running")'

# 只查看今天创建的会话
oho session list --json | jq '.[] | select(.createdAt | startswith("2026-03-03"))'

# 查看消息数超过 10 条的会话
oho session list --json | jq '.[] | select(.messageCount > 10)'
```

---

## 3. 消息历史查询

### 3.1 列出消息

```bash
# 列出所有消息
oho message list -s tidy-panda

# 限制数量
oho message list -s tidy-panda --limit 10

# JSON 格式
oho message list -s tidy-panda --json
```

**预期输出**:
```
消息列表 (tidy-panda):
  [user]      10:00  开始文档编写
  [assistant] 10:01  好的，我来帮你编写文档...
  [user]      10:05  继续模块 2
  [assistant] 10:08  模块 2 已完成
  ...
```

---

### 3.2 查看最新消息

```bash
# 查看最后 5 条消息
oho message list -s tidy-panda --limit 5

# 查看最后一条 AI 响应
oho message list -s tidy-panda --limit 1 --json | jq '.[] | select(.role == "assistant")'
```

---

### 3.3 获取单条消息详情

```bash
# 获取消息详情
oho message get msg_xxxxx -s tidy-panda

# JSON 格式
oho message get msg_xxxxx -s tidy-panda --json
```

**消息详情包含**:
```json
{
  "id": "msg_xxxxx",
  "role": "assistant",
  "content": "消息内容...",
  "createdAt": "2026-03-03T10:00:00Z",
  "tokens": 1500,
  "model": "alibaba-cn/qwen3.5-plus"
}
```

---

### 3.4 消息内容提取

```bash
# 提取所有用户消息
oho message list -s tidy-panda --json | \
  jq -r '.[] | select(.role == "user") | .content'

# 提取所有 AI 响应
oho message list -s tidy-panda --json | \
  jq -r '.[] | select(.role == "assistant") | .content'

# 提取代码块
oho message list -s tidy-panda --json | \
  jq -r '.[] | select(.role == "assistant") | .content' | \
  grep -A 100 '```go' | head -50

# 统计 Token 使用
oho message list -s tidy-panda --json | \
  jq '[.[].tokens] | add'
```

---

### 3.5 消息搜索

```bash
#!/bin/bash
# 搜索消息内容

search_messages() {
    local session=$1
    local keyword=$2
    
    oho message list -s "$session" --json | \
      jq -r --arg kw "$keyword" '.[] | select(.content | ascii_downcase | contains($kw)) | "\(.role): \(.createdAt)\n\(.content)[0:200]...\n---"'
}

# 使用
search_messages "tidy-panda" "文档"
search_messages "tidy-panda" "TODO"
```

---

## 4. 待办事项追踪

### 4.1 查看待办事项

```bash
# 获取会话待办事项
oho session todo -s tidy-panda

# JSON 格式
oho session todo -s tidy-panda --json
```

**预期输出**:
```json
{
  "todos": [
    {
      "id": "todo_1",
      "content": "编写模块 1 文档",
      "status": "completed"
    },
    {
      "id": "todo_2",
      "content": "编写模块 2 文档",
      "status": "in_progress"
    },
    {
      "id": "todo_3",
      "content": "编写模块 3 文档",
      "status": "pending"
    }
  ]
}
```

---

### 4.2 待办状态统计

```bash
# 统计各状态数量
oho session todo -s tidy-panda --json | \
  jq 'group_by(.status) | map({status: .[0].status, count: length})'

# 查看未完成的待办
oho session todo -s tidy-panda --json | \
  jq '.[] | select(.status != "completed")'

# 计算完成百分比
oho session todo -s tidy-panda --json | \
  jq '([.[] | select(.status == "completed")] | length) / length * 100'
```

---

### 4.3 待办进度可视化

```bash
#!/bin/bash
# 待办进度条

todo_progress() {
    local session=$1
    
    todos=$(oho session todo -s "$session" --json)
    total=$(echo "$todos" | jq 'length')
    completed=$(echo "$todos" | jq '[.[] | select(.status == "completed")] | length')
    
    if [ $total -eq 0 ]; then
        echo "无待办事项"
        return
    fi
    
    percent=$((completed * 100 / total))
    filled=$((percent / 5))
    empty=$((20 - filled))
    
    printf "进度：["
    printf "%${filled}s" | tr ' ' '█'
    printf "%${empty}s" | tr ' ' '░'
    printf "] %d%% (%d/%d)\n" "$percent" "$completed" "$total"
}

# 使用
todo_progress "tidy-panda"
```

---

## 5. 任务进度监控

### 5.1 实时监控

```bash
#!/bin/bash
# 实时监控任务进度

monitor_task() {
    local session=$1
    local interval=${2:-5}
    
    echo "=== 监控会话：$session ==="
    echo "刷新间隔：${interval}秒"
    echo ""
    
    last_msg_count=0
    
    while true; do
        clear
        echo "=== $(date) ==="
        echo "会话：$session"
        echo ""
        
        # 获取状态
        status=$(oho session get -s "$session" --json | jq -r '.status')
        msg_count=$(oho session get -s "$session" --json | jq -r '.messageCount')
        updated=$(oho session get -s "$session" --json | jq -r '.updatedAt')
        
        echo "状态：$status"
        echo "消息数：$msg_count ($(($msg_count - $last_msg_count)) 新增)"
        echo "最后更新：$updated"
        echo ""
        
        # 显示最后一条消息
        echo "=== 最新消息 ==="
        oho message list -s "$session" --limit 1 --json | \
          jq -r '.[0].content' | head -10
        
        last_msg_count=$msg_count
        
        if [ "$status" != "running" ]; then
            echo ""
            echo "✓ 任务已结束"
            break
        fi
        
        sleep $interval
    done
}

# 使用
monitor_task "tidy-panda" 10
```

---

### 5.2 进度日志

```bash
#!/bin/bash
# 记录任务进度日志

log_progress() {
    local session=$1
    local log_file="/tmp/oho_progress_${session}.log"
    
    while true; do
        timestamp=$(date -Iseconds)
        status=$(oho session get -s "$session" --json | jq -r '.status')
        msg_count=$(oho session get -s "$session" --json | jq -r '.messageCount')
        
        echo "[$timestamp] status=$status messages=$msg_count" >> "$log_file"
        
        if [ "$status" != "running" ]; then
            echo "[$timestamp] 任务结束" >> "$log_file"
            break
        fi
        
        sleep 30
    done
    
    echo "进度日志已保存到：$log_file"
}

# 后台运行
log_progress "tidy-panda" &
```

---

### 5.3 多任务并行监控

```bash
#!/bin/bash
# 监控多个任务

monitor_multiple() {
    sessions=("tidy-panda" "hidden-sailor" "shiny-squid")
    
    while true; do
        clear
        echo "=== 多任务监控 ==="
        echo "时间：$(date)"
        echo ""
        
        for session in "${sessions[@]}"; do
            status=$(oho session get -s "$session" --json 2>/dev/null | jq -r '.status')
            msg_count=$(oho session get -s "$session" --json 2>/dev/null | jq -r '.messageCount')
            
            case $status in
                "running")
                    icon="🔄"
                    ;;
                "completed")
                    icon="✅"
                    ;;
                "aborted")
                    icon="⛔"
                    ;;
                "error")
                    icon="❌"
                    ;;
                *)
                    icon="⏸️"
                    ;;
            esac
            
            printf "%s %-15s %-10s %3d 条消息\n" "$icon" "$session" "$status" "$msg_count"
        done
        
        sleep 10
    done
}

# 使用
monitor_multiple
```

---

### 5.4 完成通知

```bash
#!/bin/bash
# 任务完成通知

wait_for_completion() {
    local session=$1
    
    echo "等待任务完成：$session"
    
    while true; do
        status=$(oho session get -s "$session" --json | jq -r '.status')
        
        if [ "$status" != "running" ]; then
            echo ""
            echo "✓ 任务完成！状态：$status"
            
            # 发送通知
            notify-send "OpenCode" "会话 $session 已完成 (状态：$status)"
            
            # 显示总结
            echo ""
            echo "=== 任务总结 ==="
            oho session get -s "$session" --json | \
              jq '{status, messageCount, updatedAt, model}'
            
            break
        fi
        
        printf "."
        sleep 10
    done
}

# 使用
wait_for_completion "tidy-panda"
```

---

## 6. 会话总结与导出

### 6.1 生成会话总结

```bash
# 生成会话总结
oho session summarize -s tidy-panda

# JSON 格式
oho session summarize -s tidy-panda --json
```

**总结内容**:
```json
{
  "summary": "本次会话主要完成了 oho CLI 操作文档的模块 1-3 编写...",
  "keyPoints": [
    "完成了客户端初始化文档",
    "完成了验证模块文档",
    "完成了 Session 检查模块文档"
  ],
  "filesModified": [
    "docs/oho-cli-usage/01-client-initialization.md",
    "docs/oho-cli-usage/02-validation.md",
    "docs/oho-cli-usage/03-check-session.md"
  ],
  "nextSteps": [
    "继续编写模块 4",
    "完善示例代码"
  ]
}
```

---

### 6.2 导出会话数据

```bash
#!/bin/bash
# 导出完整会话数据

export_session() {
    local session=$1
    local output_dir="./export_${session}_$(date +%Y%m%d_%H%M%S)"
    
    mkdir -p "$output_dir"
    
    # 导出会话信息
    oho session get -s "$session" --json > "$output_dir/session_info.json"
    
    # 导出消息历史
    oho message list -s "$session" --json > "$output_dir/messages.json"
    
    # 导出待办事项
    oho session todo -s "$session" --json > "$output_dir/todos.json"
    
    # 生成总结
    oho session summarize -s "$session" --json > "$output_dir/summary.json"
    
    # 生成可读报告
    cat > "$output_dir/README.md" << EOF
# 会话导出报告

**会话**: $session
**导出时间**: $(date)

## 文件列表
- session_info.json - 会话基本信息
- messages.json - 完整消息历史
- todos.json - 待办事项
- summary.json - AI 生成的总结

## 统计
$(oho session get -s "$session" --json | jq '{
  状态：.status,
  消息数：.messageCount,
  模型：.model,
  创建时间：.createdAt
}')
EOF
    
    echo "会话已导出到：$output_dir"
}

# 使用
export_session "tidy-panda"
```

---

### 6.3 导出为 Markdown

```bash
#!/bin/bash
# 导出为 Markdown 格式

export_to_markdown() {
    local session=$1
    local output_file="${session}_conversation.md"
    
    {
        echo "# 会话记录：$session"
        echo ""
        echo "导出时间：$(date)"
        echo ""
        echo "---"
        echo ""
        
        oho message list -s "$session" --json | jq -r '
            .[] | 
            "**[\(.role)]** \(.createdAt)\n\n\(.content)\n\n---\n"
        '
    } > "$output_file"
    
    echo "已导出到：$output_file"
}

# 使用
export_to_markdown "tidy-panda"
```

---

### 6.4 导出为 HTML

```bash
#!/bin/bash
# 导出为 HTML 格式

export_to_html() {
    local session=$1
    local output_file="${session}_conversation.html"
    
    {
        cat << 'HTML_HEAD'
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>会话记录</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .message { margin: 20px 0; padding: 15px; border-radius: 8px; }
        .user { background: #e3f2fd; }
        .assistant { background: #f5f5f5; }
        .timestamp { color: #666; font-size: 0.9em; }
        pre { background: #f8f8f8; padding: 10px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>会话记录</h1>
HTML_HEAD
        
        oho message list -s "$session" --json | jq -r '
            .[] | 
            "<div class=\"message \(.role)\"><strong class=\"timestamp\">\(.createdAt)</strong><p>\(.content | @html)</p></div>"
        '
        
        echo "</body></html>"
    } > "$output_file"
    
    echo "已导出到：$output_file"
}
```

---

## 7. 高级查询技巧

### 7.1 自定义查询脚本

```bash
#!/bin/bash
# 综合查询仪表板

oho_dashboard() {
    clear
    echo "╔════════════════════════════════════════════════════════╗"
    echo "║           OpenCode 会话监控仪表板                      ║"
    echo "╚════════════════════════════════════════════════════════╝"
    echo ""
    echo "时间：$(date)"
    echo ""
    
    # 总会话统计
    total=$(oho session list --json | jq 'length')
    running=$(oho session status --json | jq '[.[] | select(.status == "running")] | length')
    completed=$(oho session status --json | jq '[.[] | select(.status == "completed")] | length')
    
    echo "📊 总体统计"
    echo "   总会话数：$total"
    echo "   运行中：$running"
    echo "   已完成：$completed"
    echo ""
    
    # 运行中的会话
    if [ $running -gt 0 ]; then
        echo "🔄 运行中的会话"
        oho session status --json | jq -r '.[] | select(.status == "running") | "   • \(.slug) - \(.messageCount) 条消息"'
        echo ""
    fi
    
    # 今日活动
    echo "📈 今日活动"
    today=$(date +%Y-%m-%d)
    today_sessions=$(oho session list --json | jq --arg t "$today" '[.[] | select(.createdAt | startswith($t))] | length')
    today_messages=$(oho session list --json | jq --arg t "$today" '[.[] | select(.createdAt | startswith($t)) | .messageCount] | add // 0')
    echo "   新会话：$today_sessions"
    echo "   消息数：$today_messages"
    echo ""
    
    # 资源使用
    echo "💾 存储使用"
    du -sh ~/.local/share/opencode/ 2>/dev/null | awk '{print "   OpenCode 数据：" $1}'
    echo ""
}

# 使用
oho_dashboard
```

---

### 7.2 性能分析

```bash
#!/bin/bash
# 分析会话性能

analyze_performance() {
    local session=$1
    
    echo "=== 会话性能分析：$session ==="
    echo ""
    
    # 获取消息列表
    messages=$(oho message list -s "$session" --json)
    
    # 统计 Token
    total_tokens=$(echo "$messages" | jq '[.[].tokens] | add // 0')
    avg_tokens=$(echo "$messages" | jq '[.[].tokens] | add / length // 0')
    
    echo "Token 统计"
    echo "   总 Token: $total_tokens"
    echo "   平均每条：$(printf "%.0f" $avg_tokens)"
    echo ""
    
    # 响应时间分析（如果有）
    echo "响应分析"
    msg_count=$(echo "$messages" | jq 'length')
    echo "   消息数：$msg_count"
    
    # 模型使用
    models=$(echo "$messages" | jq '[.[].model] | unique')
    echo "   使用模型：$models"
}

# 使用
analyze_performance "tidy-panda"
```

---

### 7.3 趋势分析

```bash
#!/bin/bash
# 分析使用趋势

analyze_trend() {
    echo "=== 使用趋势分析 ==="
    echo ""
    
    # 按天统计会话数
    for day in $(seq 7 -1 1); do
        date_str=$(date -d "$day days ago" +%Y-%m-%d)
        count=$(oho session list --json | jq --arg d "$date_str" '[.[] | select(.createdAt | startswith($d))] | length')
        printf "%s: %2d 会话\n" "$date_str" "$count"
    done
}

# 使用
analyze_trend
```

---

### 7.4 自动化报告

```bash
#!/bin/bash
# 生成日报

generate_daily_report() {
    local date_str=$(date +%Y-%m-%d)
    local report_file="opencode_daily_${date_str}.md"
    
    {
        echo "# OpenCode 日报"
        echo ""
        echo "日期：$date_str"
        echo "生成时间：$(date)"
        echo ""
        
        echo "## 总体统计"
        echo ""
        total=$(oho session list --json | jq --arg d "$date_str" '[.[] | select(.createdAt | startswith($d))] | length')
        messages=$(oho session list --json | jq --arg d "$date_str" '[.[] | select(.createdAt | startswith($d)) | .messageCount] | add // 0')
        echo "- 新会话数：$total"
        echo "- 总消息数：$messages"
        echo ""
        
        echo "## 会话详情"
        echo ""
        oho session list --json | jq --arg d "$date_str" -r '
            .[] | select(.createdAt | startswith($d)) | 
            "- **\(.slug)**: \(.messageCount) 条消息，状态：\(.status)"
        '
        echo ""
        
        echo "## 活跃项目"
        echo ""
        oho session list --json | jq --arg d "$date_str" -r '
            [.[] | select(.createdAt | startswith($d))] | 
            group_by(.project.path) | 
            map({project: .[0].project.path, sessions: length}) | 
            .[] | "- \(.project): \(.sessions) 会话"
        '
    } > "$report_file"
    
    echo "日报已生成：$report_file"
}

# 使用
generate_daily_report
```

---

## 🔧 实用技巧

### 技巧 1: 快速状态检查

```bash
# 一行命令检查所有会话
oho session status --json | jq -r '.[] | "\(.slug): \(.status)"'
```

---

### 技巧 2: 消息计数

```bash
# 统计每个会话的消息数
oho session list --json | jq -r '.[] | "\(.slug): \(.messageCount) 条消息"'
```

---

### 技巧 3: 错误会话快速定位

```bash
# 找出所有错误会话
oho session status --json | jq -r '.[] | select(.status == "error") | .slug'
```

---

### 技巧 4: 会话活动热力图

```bash
# 显示 24 小时活动分布
oho session list --json | jq -r '.[].createdAt' | \
  cut -dT -f2 | cut -d: -f1 | sort | uniq -c | \
  awk '{printf "%02d:00 - %d 会话\n", $2, $1}'
```

---

## 📝 检查清单

在查询任务状态时，请确认：

- [ ] 使用正确的会话 ID 或 Slug
- [ ] 选择合适的查询粒度（详情/列表/总结）
- [ ] 必要时使用 JSON 格式便于处理
- [ ] 定期保存重要查询结果
- [ ] 监控长时间运行的任务

---

## 🔗 相关文档

- [模块 3: 检查 Session](./03-check-session.md) - 会话管理基础
- [模块 7: 中断任务](./07-interrupt-task.md) - 任务控制
- [模块 1: 客户端初始化](./01-client-initialization.md) - 连接配置

---

## 🆘 常见问题

### Q1: 如何查看任务执行进度？

**A**:
```bash
# 查看会话状态
oho session get -s slug --json | jq '.status'

# 查看消息增长
oho message list -s slug --json | jq 'length'
```

---

### Q2: 如何知道 AI 正在做什么？

**A**:
```bash
# 查看最新消息
oho message list -s slug --limit 5

# 查看待办事项
oho session todo -s slug
```

---

### Q3: 如何导出完整的对话记录？

**A**:
```bash
# 导出为 JSON
oho message list -s slug --json > conversation.json

# 导出为 Markdown
# 使用 6.3 节的脚本
```

---

### Q4: 如何监控多个任务？

**A**:
```bash
# 使用 5.3 节的多任务监控脚本
# 或简单轮询
while true; do
    oho session status
    sleep 30
done
```

---

*文档生成时间：2026-03-03 11:17 CST*
