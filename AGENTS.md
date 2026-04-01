# AGENTS.md - OpenCode CLI 开发指南

本文档为智能编码代理提供开发 OpenCode CLI 工具的指南。

## 项目概述

- **项目名称**: opencode_cli
- **技术栈**: Shell / Python
- **用途**: OpenCode 命令行工具开发

---

## Build/Lint/Test Commands

### 测试命令

```bash
# 运行所有 Python 测试
python -m pytest scripts/

# 运行单个测试文件
python -m pytest scripts/opencode-test.py -v

# 运行特定测试
python -m pytest scripts/opencode-test.py -k "test_name"
```

### 代码检查

```bash
# Python 语法检查
python -m py_compile scripts/*.py

# Shell 脚本语法检查
bash -n scripts/*.sh
```

---

## 代码风格指南

### 通用原则

- **语言**: 中文回答，中文注释
- **简洁**: 文字简练，代码详细
- **可维护**: 代码结构清晰，易于阅读

### 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| Shell 脚本 | kebab-case | `opencode-submit.sh` |
| Python 文件 | snake_case | `opencode_submit.py` |
| Python 函数 | snake_case | `def make_request()` |
| 常量 | UPPER_SNAKE_CASE | `MAX_RETRY_COUNT` |

### Shell 脚本规范

```bash
#!/bin/bash
# 错误时立即退出
set -e

# 变量使用 ${VAR} 引用
PROJECT_DIR="${HOME}/projects/myproject"

# 函数定义使用 snake_case
function check_prerequisites() {
    if ! command -v git &> /dev/null; then
        echo "Error: git is required"
        exit 1
    fi
}
```

### Python 规范

```python
#!/usr/bin/env python3
"""模块文档字符串"""

import json
import sys
from typing import Optional

# 常量定义
DEFAULT_TIMEOUT = 30

def make_request(path: str, method: str = 'GET') -> dict:
    """发送 HTTP 请求"""
    pass

class OpenCodeClient:
    """OpenCode 客户端类"""
    
    def __init__(self, host: str = 'localhost', port: int = 4096):
        self.host = host
        self.port = port
```

### 错误处理

```python
# 使用 try-except 包装可能失败的操作
try:
    result = make_request('/api/endpoint')
except Exception as e:
    print(f"Error: {e}")
    sys.exit(1)

# Shell 脚本使用 set -e
set -e
```

### Import 顺序

1. 标准库 (import json)
2. 第三方库 (import requests)
3. 本地模块 (from . import utils)

---

## oho 项目构建指南

oho 是用 Go 编写的 CLI 工具，位于 `oho/` 子目录。

### ⚠️ 重要：构建命令必须包含 ldflags

oho 使用 ldflags 注入版本信息。**必须使用以下格式**，否则版本信息会显示为 `commit: none, built: unknown`：

```bash
# 正确：包含完整 ldflags
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
go build -ldflags "-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${DATE}" -o bin/oho ./cmd

# 错误：缺少 ldflags（版本信息会丢失）
go build -o bin/oho ./cmd
```

### 推荐使用 Makefile

```bash
cd oho

# 本地构建（推荐）
make build

# Linux 交叉编译
make build-linux

# 开发模式（无版本信息）
make build-dev
```

### 验证构建

```bash
./bin/oho --version
# 应该显示: oho version v1.x.x (commit: xxxxx, built: YYYY-MM-DDThh:mm:ssZ)
# 而不是:    oho version dev (commit: none, built: unknown)
```

---

## 项目结构

```
opencode_cli/
├── oho/                    # Go CLI 项目 (主项目)
│   ├── cmd/                # 命令入口
│   ├── internal/           # 内部包
│   ├── Makefile           # 构建目标
│   └── README.md           # oho 详细文档
├── scripts/               # Python/Shell 脚本
│   ├── opencode-submit.py
│   └── opencode-test.py
└── AGENTS.md              # 本文件
```

---

## 开发流程

1. 创建新脚本时使用合适的命名 (kebab-case/shnake_case)
2. 添加可执行权限: `chmod +x script.sh`
3. 运行测试验证: `python -m pytest scripts/`
4. 代码审查前确保无语法错误
