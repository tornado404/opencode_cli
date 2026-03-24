# 共享技能库 opencode 技能更新

**更新时间**: 2026-03-24 10:50 CST  
**更新人**: nanobot 🐈  
**技能版本**: v4.5.0 → v4.6.0

---

## 📋 更新内容

### 1. 更新 oho CLI 版本
- **v1.1.0** → **v1.1.1**
- 新增 `--timeout` 参数支持

### 2. 添加超时配置说明

#### 新增参数
在关键参数表中添加：
| `--timeout` | int | 请求超时时间（秒） | 300 |

#### 新增超时处理方法

**方法 1: 使用 `--no-reply` 参数** (推荐)
```bash
oho add "分析项目结构" --directory /mnt/d/fe/project --no-reply
```

**方法 2: 增加超时时间**
```bash
export OPENCODE_CLIENT_TIMEOUT=600
oho add "复杂任务" --directory /mnt/d/fe/project
```

**方法 3: 使用 `--timeout` 参数** (最方便)
```bash
oho add "复杂任务" --directory /mnt/d/fe/project --timeout 600
```

**方法 4: 使用异步命令**
```bash
oho session create --title "任务" --path /mnt/d/fe/project
oho message prompt-async -s <session-id> "任务描述"
```

#### 新增超时配置优先级表

| 配置方式 | 优先级 | 说明 |
|----------|--------|------|
| `--timeout` 参数 | 最高 | 临时覆盖，仅对当前命令有效 |
| `OPENCODE_CLIENT_TIMEOUT` 环境变量 | 中 | 对当前 shell 会话有效 |
| 默认值 | 最低 | 300 秒（5 分钟） |

#### 新增超时错误提示示例
```
请求超时（300 秒）

建议:
  1. 使用 --no-reply 参数避免等待
  2. 设置环境变量增加超时：export OPENCODE_CLIENT_TIMEOUT=600
  3. 使用异步命令：oho message prompt-async -s <session-id> "任务"
```

---

## 📁 修改文件

| 文件 | 修改内容 | 状态 |
|------|---------|------|
| `/mnt/d/fe/shared-skills/opencode/SKILL.md` | 添加超时配置完整说明 | ✅ |

---

## 🔄 同步状态

### 已同步
- ✅ `oho/README.md` (英文)
- ✅ `oho/README_zh.md` (中文)
- ✅ `oho/internal/client/client.go` (错误提示)
- ✅ `oho/cmd/add/add.go` (--timeout 参数)
- ✅ `/mnt/d/fe/shared-skills/opencode/SKILL.md` (技能文档)

### 技能文档版本
- **opencode/SKILL.md**: v4.6.0 (2026-03-24 10:50)

---

## 📚 使用示例

### babylon3DWorld 项目任务（更新版）

```bash
#!/bin/bash
# 提交编码任务（带超时配置）

oho add "@hephaestus ulw 优化编辑器与 world 页面的导航逻辑

**编码目标**:
1. 编辑器返回 world 页面时：直接刷新页面，不再判断是否 editor 造成了改动
2. world 页进入 editor 时：不再缓存 world 本身

**关键词**: ulw" \
  --directory /mnt/d/fe/babylon3DWorld \
  --title "ulw - 优化编辑器导航逻辑" \
  --timeout 600 \
  --no-reply

echo "✅ 任务已提交（超时 600 秒）"
```

---

## ✅ 验收清单

- [x] 参数表添加 `--timeout` 参数
- [x] 添加 4 种避免超时的方法
- [x] 添加超时配置优先级表
- [x] 添加超时错误提示示例
- [x] 更新 oho CLI 版本号
- [x] 更新最后更新时间

---

**更新完成时间**: 2026-03-24 10:50 CST
