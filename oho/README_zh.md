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

### 支持的文件类型

- **图片**: `.jpg`, `.png`, `.gif`, `.webp`, `.bmp`, `.svg`
- **文档**: `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`
- **文本**: `.txt`, `.md`, `.html`, `.css`, `.js`, `.json`, `.yaml`
- **代码**: `.py`, `.go`, `.java`, `.c`, `.cpp`, `.rs`, `.ts`, `.tsx`
- **其他**: `.zip`, `.tar`, `.gz`, `.mp3`, `.mp4`, `.wav`

## 完整文档

完整的命令参考、安装说明和详细文档，请查看 [主项目 README](../README_zh.md)。
