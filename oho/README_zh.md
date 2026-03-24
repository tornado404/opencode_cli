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

# 发送消息并附带文件
oho message add -s <session-id> "请分析这个文件" --file /path/to/file.txt

# 发送消息并附带多个文件
oho message add -s <session-id> "请比较这些文件" --file file1.txt --file file2.txt

# 发送图片进行分析
oho message add -s <session-id> "请描述这张图片" --file image.jpg

# 销毁会话
oho session delete <session-id>
```

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

### 支持的文件类型

- **图片**: `.jpg`, `.png`, `.gif`, `.webp`, `.bmp`, `.svg`
- **文档**: `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`
- **文本**: `.txt`, `.md`, `.html`, `.css`, `.js`, `.json`, `.yaml`
- **代码**: `.py`, `.go`, `.java`, `.c`, `.cpp`, `.rs`, `.ts`, `.tsx`
- **其他**: `.zip`, `.tar`, `.gz`, `.mp3`, `.mp4`, `.wav`

## 完整文档

完整的命令参考、安装说明和详细文档，请查看 [主项目 README](../README_zh.md)。
