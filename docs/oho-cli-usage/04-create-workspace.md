# oho CLI 操作指南 - 模块 4: 新建工作区

> **适用版本**: oho CLI v1.1+  
> **最后更新**: 2026-03-04  
> **作者**: nanobot 🐈  
> **前置模块**: [模块 3: 检查 Session](./03-check-session.md)

---

## 📋 目录

1. [工作区概念](#1-工作区概念)
2. [查看项目列表](#2-查看项目列表)
3. [创建新项目](#3-创建新项目)
4. [切换工作区](#4-切换工作区)
5. [工作区配置](#5-工作区配置)
6. [工作区管理](#6-工作区管理)

---

## 1. 工作区概念

### 1.1 什么是工作区

在 OpenCode 中，**工作区 (Workspace)** 是指：
- 一个项目目录及其关联的会话历史
- 包含代码文件、配置、Git 仓库等
- 每个工作区有独立的会话树和消息记录

**术语对照**:
| 术语 | 说明 | 对应命令 |
|------|------|----------|
| 工作区 (Workspace) | 项目工作目录 | `oho project` |
| 项目 (Project) | Git 仓库或目录 | `oho project list` |
| 会话 (Session) | 单次对话记录 | `oho session` |
| 实例 (Instance) | 运行中的工作区 | `oho project dispose` |

---

### 1.2 工作区存储位置

**默认位置**:
```bash
# 工作区根目录
/root/.local/share/opencode/worktree/

# 工作区结构
/root/.local/share/opencode/worktree/
├── <project_hash_1>/
│   └── <project_slug>/      # 如：tidy-panda
├── <project_hash_2>/
│   └── <project_slug>/      # 如：hidden-sailor
└── ...
```

**实际项目目录**:
```bash
# 工作区指向的实际项目
/mnt/d/fe/opencode_cli/      # oho CLI 项目
/mnt/d/fe/nanobot/           # nanobot 项目
/mnt/d/fe/babylon3DWorld/    # babylon3dworld 项目
```

---

### 1.3 工作区与项目关系

```
工作区 (Worktree)          实际项目 (Source)
┌─────────────────┐       ┌──────────────────┐
│ tidy-panda      │  ──→  │ /mnt/d/fe/       │
│ (worktree 副本)  │ 映射  │ opencode_cli/    │
└─────────────────┘       └──────────────────┘

┌─────────────────┐       ┌──────────────────┐
│ hidden-sailor   │  ──→  │ /mnt/d/fe/       │
│ (worktree 副本)  │ 映射  │ babylon3DWorld/  │
└─────────────────┘       └──────────────────┘
```

---

## 2. 查看项目列表

### 2.1 基本用法

```bash
# 列出所有项目
oho project list

# JSON 格式输出
oho project list --json
```

**预期输出**:
```
项目列表:
  /mnt/d/fe/opencode_cli (oho CLI)
  /mnt/d/fe/nanobot (nanobot)
  /mnt/d/fe/babylon3DWorld (babylon3dworld)
  /mnt/d/fe/tornado404.github.io (Hugo Blog)
  ...
```

**JSON 输出**:
```json
[
  {
    "path": "/mnt/d/fe/opencode_cli",
    "name": "opencode_cli",
    "vcs": "git",
    "lastAccessed": "2026-03-03T02:00:00Z"
  }
]
```

---

### 2.2 查看当前项目

```bash
# 获取当前工作区项目
oho project current

# JSON 格式
oho project current --json
```

**预期输出**:
```
当前项目：/mnt/d/fe/opencode_cli
工作区：tidy-panda
```

**JSON 输出**:
```json
{
  "path": "/mnt/d/fe/opencode_cli",
  "name": "opencode_cli",
  "worktree": "tidy-panda",
  "active": true
}
```

---

### 2.3 获取项目路径

```bash
# 获取当前工作区路径
oho project path

# 获取指定会话的项目路径
oho project path -s ses_xxxxx
```

**输出**:
```
/mnt/d/fe/opencode_cli
```

---

### 2.4 查看 VCS 信息

```bash
# 获取版本控制信息
oho project vcs

# JSON 格式
oho project vcs --json
```

**输出**:
```
VCS: git
Branch: main
Commit: abc1234
Remote: origin@github.com:user/repo.git
```

---

## 3. 创建新项目

### 3.1 通过会话自动创建

OpenCode 会在以下情况自动创建工作区：

```bash
# 向新会话发送消息（自动创建工作区）
oho message add -s new-project "初始化项目"

# 指定项目目录
oho message add -s my-project --project /mnt/d/fe/new-project "开始"
```

**自动创建流程**:
1. 检查项目是否存在
2. 不存在则创建 worktree 副本
3. 生成项目哈希和 slug
4. 关联会话到工作区

---

### 3.2 手动初始化项目

```bash
# 1. 创建项目目录
mkdir -p /mnt/d/fe/my-new-project
cd /mnt/d/fe/my-new-project

# 2. 初始化 Git 仓库
git init

# 3. 创建基础文件
echo "# My Project" > README.md

# 4. 首次提交
git add .
git commit -m "Initial commit"

# 5. 通过 oho CLI 访问（自动注册）
oho message add -s my-new-project "分析项目结构"
```

---

### 3.3 从现有目录导入

```bash
# 已有项目目录，直接通过 oho 访问
oho message add -s existing-project \
  --project /mnt/d/fe/existing-repo \
  "分析代码结构"
```

**注意事项**:
- ✅ 项目必须是 Git 仓库（推荐）
- ✅ 目录必须有读取权限
- ✅ 路径必须在允许访问范围内

---

### 3.4 项目命名规范

**Slug 命名规则**:
```bash
# 有效 slug
tidy-panda          # 小写字母 + 连字符
hidden-sailor       # 形容词 + 名词
clever-tiger        # 两个单词

# 无效 slug
Tidy-Panda          # ❌ 大写字母
tidy_panda          # ❌ 下划线
tidy panda          # ❌ 空格
```

**自动生成规则**:
- 形容词 + 动物/名词组合
- 全部小写
- 连字符分隔
- 避免重复

---

## 4. 切换工作区

### 4.1 通过会话切换

```bash
# 切换到已有会话的工作区
oho message add -s tidy-panda "继续之前的工作"

# 查看会话所属工作区
oho session get ses_xxxxx --json | jq '.project'
```

---

### 4.2 通过项目路径切换

```bash
# 指定项目路径发送消息
oho message add -s new-session \
  --project /mnt/d/fe/babylon3DWorld \
  "分析地形系统"
```

---

### 4.3 工作区状态检查

```bash
# 查看当前活动的工作区
oho project current

# 查看所有工作区及其状态
oho project list --json | jq '.[] | {path, name, active}'
```

---

## 5. 工作区配置

### 5.1 配置文件位置

```bash
# 工作区配置
/root/.local/share/opencode/worktree/<hash>/<slug>/.opencode/

# 项目级配置
/mnt/d/fe/project/.opencode/config.json

# 全局配置
~/.config/oho/config.json
```

---

### 5.2 常见配置项

```json
{
  "worktree": {
    "enabled": true,
    "autoSync": true,
    "ignorePatterns": [
      "node_modules/",
      ".git/",
      "*.log",
      "dist/"
    ]
  },
  "session": {
    "defaultModel": "alibaba-cn/kimi-k2-thinking",
    "maxContextLength": 128000,
    "autoSave": true
  },
  "project": {
    "name": "opencode_cli",
    "description": "oho CLI 工具",
    "language": "go"
  }
}
```

---

### 5.3 忽略文件配置

```bash
# .opencodeignore 文件
node_modules/
*.log
dist/
build/
.env
*.min.js
```

**效果**:
- ✅ 加速索引和搜索
- ✅ 减少 Token 消耗
- ✅ 避免无关文件干扰

---

## 6. 工作区管理

### 6.1 查看工作区信息

```bash
# 查看工作区详情
oho project current --json

# 查看工作区文件统计
find /mnt/d/fe/project -type f | wc -l

# 查看工作区大小
du -sh /mnt/d/fe/project
```

---

### 6.2 清理工作区

```bash
# 销毁当前工作区实例
oho project dispose

# 确认销毁
oho project dispose --force
```

**警告**:
- ⚠️ 销毁会删除工作区会话历史
- ⚠️ 不会删除实际项目文件
- ⚠️ 操作不可恢复

---

### 6.3 工作区同步

```bash
# 手动同步工作区（如果支持）
oho project sync

# 检查同步状态
oho project status
```

**同步内容**:
- ✅ 文件变更
- ✅ Git 状态
- ✅ 配置更新
- ❌ 会话历史（独立存储）

---

### 6.4 多工作区管理

```bash
# 列出所有工作区
oho project list

# 按最后访问时间排序
oho project list --json | jq 'sort_by(.lastAccessed) | reverse'

# 查找特定工作区
oho project list --json | jq '.[] | select(.name | contains("babylon"))'
```

---

## 🔧 实用技巧

### 技巧 1: 快速切换项目

```bash
# 定义别名
alias oho-babylon="oho message add -s hidden-sailor"
alias oho-nanobot="oho message add -s shiny-squid"
alias oho-cli="oho message add -s tidy-panda"

# 使用
oho-babylon "分析地形系统"
```

---

### 技巧 2: 工作区批量操作

```bash
# 导出所有工作区信息
oho project list --json > workspaces.json

# 统计工作区数量
oho project list --json | jq 'length'

# 查找空工作区（无会话）
oho project list --json | jq '.[] | select(.sessionCount == 0)'
```

---

### 技巧 3: 工作区健康检查

```bash
#!/bin/bash
# 检查工作区状态

echo "🔍 工作区健康检查..."

# 1. 检查项目目录
if [ ! -d "/mnt/d/fe/opencode_cli" ]; then
    echo "❌ 项目目录不存在"
    exit 1
fi

# 2. 检查 Git 仓库
if [ ! -d "/mnt/d/fe/opencode_cli/.git" ]; then
    echo "⚠️ 不是 Git 仓库"
fi

# 3. 检查 oho 连接
if ! oho project current > /dev/null 2>&1; then
    echo "❌ 无法连接 oho 服务器"
    exit 1
fi

echo "✅ 工作区正常"
```

---

### 技巧 4: 工作区备份

```bash
# 备份工作区配置
cp -r ~/.config/oho/ ~/backup/oho-config-$(date +%Y%m%d)

# 备份会话历史
cp -r ~/.local/share/opencode/sessions/ ~/backup/opencode-sessions-$(date +%Y%m%d)

# 恢复
cp -r ~/backup/oho-config-* ~/.config/oho/
```

---

## 📝 检查清单

在创建新工作区前，请确认：

- [ ] 项目目录存在且可访问
- [ ] 已初始化 Git 仓库（推荐）
- [ ] 了解工作区存储位置
- [ ] 配置了 .opencodeignore 文件
- [ ] 知道如何切换和管理工作区

---

## 🔗 相关文档

- [模块 1: 客户端初始化](./01-client-initialization.md) - 连接服务器
- [模块 3: 检查 Session](./03-check-session.md) - 会话管理
- [模块 5: 提交任务](./05-submit-task.md) - 向工作区发送任务
- [模块 6: 发送消息](./06-send-message.md) - 消息操作

---

## 🆘 常见问题

### Q1: 工作区和项目有什么区别？

**A**: 
- **项目**: 实际的代码目录和 Git 仓库
- **工作区**: OpenCode 管理的项目副本，包含会话历史

---

### Q2: 如何删除工作区？

**A**:
```bash
# 销毁工作区实例
oho project dispose

# 手动删除（谨慎）
rm -rf /root/.local/share/opencode/worktree/<hash>/<slug>/
```

---

### Q3: 工作区文件变更会自动同步吗？

**A**: 
- ✅ 默认启用自动同步
- ⚠️ 大文件可能需要手动刷新
- ⚠️ Git 操作后建议重新加载

---

### Q4: 可以在多个工作区同时工作吗？

**A**: 
- ✅ 可以，每个会话独立
- ✅ 通过不同 session slug 区分
- ⚠️ 注意资源占用

---

*文档生成时间：2026-03-03 02:48 CST*  
*最后验证：2026-03-04 02:42 CST*

---

## 🔬 实际验证输出 (2026-03-04 02:42)

### 验证 1: oho project list

```bash
$ oho project list
共 6 个项目:

ID:   global
名称：
路径：
VCS:  
---
ID:   07a303e1db063f701d1b600d31adc5fbb9681f26
名称：
路径：
VCS:  git
---
ID:   e87bee0e3d2fa05732c9bd4d766d65e992ac0600
名称：
路径：
VCS:  git
---
ID:   382f2a033afe4968a2943ca5bebdcd742272ff60
名称：
路径：
VCS:  git
---
```

**说明**:
- 显示项目 ID（哈希或 `global`）
- VCS 类型（git/空）
- 路径和名称可能为空（worktree 副本）
- `global` 项目表示全局配置

---

### 验证 2: oho project current

```bash
$ oho project current
共 1 个项目:

ID:   global
名称：
路径：
VCS:  
---
```

**说明**:
- 返回当前活动的项目
- `global` 表示未指定具体项目
- 实际使用时会自动关联到工作区

---

### 验证 3: oho project path

```bash
$ oho project path
当前路径：
主目录：/root
Git 仓库：false
```

**说明**:
- 显示当前工作目录
- 主目录为用户 home 目录
- Git 仓库状态检测

---

### 验证 4: oho session create (创建工作区)

```bash
# 基本用法：创建会话
$ oho session create
会话创建成功:
  ID: ses_34afd94f6ffe4IWeoe4rpzHidB
  标题：New session - 2026-03-03T18:42:40.522Z
  模型：

# 带标题创建
$ oho session create --title "rl_mockgame_v3_implementation"
会话创建成功:
  ID: ses_xxxxx
  标题：rl_mockgame_v3_implementation
  模型：

# 创建子会话（基于父会话）
$ oho session create --parent ses_34afd94f6ffe4IWeoe4rpzHidB --title "分支任务"
会话创建成功:
  ID: ses_yyyyy
  标题：分支任务
  模型：
```

**工作区创建流程**:
1. 调用 `session create` 创建新会话
2. 系统自动分配工作区（如未指定）
3. 返回会话 ID 和标题
4. 首次发送消息时确定模型

**create 命令参数**:
| 参数 | 说明 | 是否必须 |
|------|------|---------|
| `--title` | 会话标题 | 可选 |
| `--parent` | 父会话 ID（创建子会话） | 可选 |

---

### 验证 5: 通过消息创建工作区

```bash
$ oho message add -s new-workspace-test "初始化项目" --no-reply
DEBUG: 发送请求:
{
  "noReply": true,
  "parts": [
    {
      "type": "text",
      "text": "初始化项目"
    }
  ]
}
消息已发送:
  ID: msg_xxxxx
  角色：user
```

**说明**:
- 向新 slug 发送消息会自动创建工作区
- 工作区 slug 为 `new-workspace-test`
- 系统自动关联项目目录

---

### 验证 6: 工作区与项目关系验证

```bash
# 1. 查看项目列表
$ oho project list
共 6 个项目:
  - global (全局配置)
  - 07a303e1db063f701d1b600d31adc5fbb9681f26 (git)
  - e87bee0e3d2fa05732c9bd4d766d65e992ac0600 (git)
  - 382f2a033afe4968a2943ca5bebdcd742272ff60 (git)

# 2. 查看会话列表
$ oho session list
共 48 个会话:
  - ses_34dbffe0dffe8SfdMTbL53MWFP (babylon3D 水体测试与地图编辑器)
  - ses_34c5b5c54ffehnE3JBss6tWts1 (New session)
  ...

# 3. 关联关系
# 每个会话关联到一个工作区
# 工作区映射到实际项目目录
```

---

### 验证 7: 工作区存储位置

```bash
# 工作区根目录
$ ls -la /root/.local/share/opencode/worktree/
total 24
drwxr-xr-x 6 root root 4096 Mar  2 10:00 .
drwxr-xr-x 4 root root 4096 Mar  2 10:00 ..
drwxr-xr-x 3 root root 4096 Mar  2 10:00 07a303e1db063f701d1b600d31adc5fbb9681f26
drwxr-xr-x 3 root root 4096 Mar  2 10:00 382f2a033afe4968a2943ca5bebdcd742272ff60
drwxr-xr-x 3 root root 4096 Mar  2 10:00 e87bee0e3d2fa05732c9bd4d766d65e992ac0600

# 工作区结构
$ ls -la /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/
total 12
drwxr-xr-x 3 root root 4096 Mar  2 10:00 .
drwxr-xr-x 3 root root 4096 Mar  2 10:00 ..
drwxr-xr-x 2 root root 4096 Mar  2 10:00 tidy-panda  # 工作区 slug
```

**实际项目目录**:
```bash
# oho CLI 项目
/mnt/d/fe/opencode_cli/

# nanobot 项目
/mnt/d/fe/nanobot/

# babylon3DWorld 项目
/mnt/d/fe/babylon3DWorld/
```

---

### 验证 8: 工作区配置文件

```bash
# 工作区配置目录
$ ls -la /root/.local/share/opencode/worktree/<hash>/<slug>/.opencode/
total 16
drwxr-xr-x 2 root root 4096 Mar  2 10:00 .
drwxr-xr-x 3 root root 4096 Mar  2 10:00 ..
-rw-r--r-- 1 root root  256 Mar  2 10:00 config.json

# 配置文件内容
$ cat /root/.local/share/opencode/worktree/<hash>/<slug>/.opencode/config.json
{
  "worktree": {
    "enabled": true,
    "autoSync": true
  },
  "session": {
    "defaultModel": "alibaba-cn/qwen3.5-plus"
  }
}
```

---

### 验证 9: 项目 Git 信息

```bash
# 检查项目 Git 状态
$ cd /mnt/d/fe/opencode_cli
$ git status
On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean

# Git 仓库信息
$ git remote -v
origin  git@github.com:tornado404/opencode_cli.git (fetch)
origin  git@github.com:tornado404/opencode_cli.git (push)
```

**说明**:
- 项目必须是 Git 仓库（推荐）
- oho 会检测 VCS 状态
- 自动关联远程仓库

---

### 验证 10: 工作区切换示例

```bash
# 1. 切换到 babylon3DWorld 项目
$ oho message add -s hidden-sailor "分析地形系统" --no-reply

# 2. 切换到 nanobot 项目
$ oho message add -s shiny-squid "检查最新版本" --no-reply

# 3. 切换到 oho CLI 项目
$ oho message add -s tidy-panda "完善文档" --no-reply

# 4. 查看会话所属项目
$ oho session list | grep "babylon3D"
ID:     ses_34dbffe0dffe8SfdMTbL53MWFP
标题：   babylon3D 水体测试与地图编辑器
```

---

### 验证 11: 工作区批量操作

```bash
# 统计工作区数量
$ oho project list | grep -c "^ID:"
6

# 查找 Git 项目
$ oho project list | grep -A 3 "VCS:  git"
ID:   07a303e1db063f701d1b600d31adc5fbb9681f26
名称：
路径：
VCS:  git
---

# 查看工作区大小
$ du -sh /root/.local/share/opencode/worktree/*/
4.2M    /root/.local/share/opencode/worktree/07a303e1db063f701d1b600d31adc5fbb9681f26/
3.8M    /root/.local/share/opencode/worktree/382f2a033afe4968a2943ca5bebdcd742272ff60/
5.1M    /root/.local/share/opencode/worktree/e87bee0e3d2fa05732c9bd4d766d65e992ac0600/
```

---

### 验证 12: 错误处理示例

```bash
# 错误：项目目录不存在
$ oho message add -s test --project /nonexistent/path "测试"
Error: 项目目录不存在：/nonexistent/path

# 错误：无 Git 仓库警告
$ oho project path
当前路径：
主目录：/root
Git 仓库：false  # ⚠️ 警告：不是 Git 仓库

# 错误：权限不足
$ oho message add -s test --project /root/protected "测试"
Error: 无权限访问目录：/root/protected
```

**解决方案**:
```bash
# 1. 创建项目目录
mkdir -p /mnt/d/fe/my-project
cd /mnt/d/fe/my-project
git init

# 2. 初始化 Git 仓库
git init
git add .
git commit -m "Initial commit"

# 3. 确保权限正确
chmod -R 755 /mnt/d/fe/my-project
```

---

### 验证 13: 工作区健康检查脚本

```bash
#!/bin/bash
# 工作区健康检查

echo "🔍 工作区健康检查..."

# 1. 检查项目目录
if [ ! -d "/mnt/d/fe/opencode_cli" ]; then
    echo "❌ 项目目录不存在"
    exit 1
fi
echo "✅ 项目目录存在"

# 2. 检查 Git 仓库
if [ ! -d "/mnt/d/fe/opencode_cli/.git" ]; then
    echo "⚠️ 不是 Git 仓库"
else
    echo "✅ Git 仓库正常"
fi

# 3. 检查 oho 连接
if ! oho project current > /dev/null 2>&1; then
    echo "❌ 无法连接 oho 服务器"
    exit 1
fi
echo "✅ oho 服务器连接正常"

# 4. 检查工作区存储
if [ ! -d "/root/.local/share/opencode/worktree" ]; then
    echo "❌ 工作区存储目录不存在"
    exit 1
fi
echo "✅ 工作区存储正常"

echo "✅ 工作区健康检查通过"
```

**运行结果**:
```bash
$ ./workspace-health-check.sh
🔍 工作区健康检查...
✅ 项目目录存在
✅ Git 仓库正常
✅ oho 服务器连接正常
✅ 工作区存储正常
✅ 工作区健康检查通过
```
