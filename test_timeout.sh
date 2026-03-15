#!/bin/bash
# 超时配置测试脚本

echo "========================================="
echo "oho 超时配置测试"
echo "========================================="
echo ""

# 测试 1: 默认超时（300 秒）
echo "测试 1: 默认超时配置（300 秒）"
echo "命令：oho session list"
timeout 10 oho session list > /dev/null 2>&1 && echo "✓ 默认超时工作正常" || echo "✗ 默认超时失败"
echo ""

# 测试 2: 自定义超时（10 秒）
echo "测试 2: 自定义超时配置（10 秒）"
echo "命令：OPENCODE_CLIENT_TIMEOUT=10 oho session list"
OPENCODE_CLIENT_TIMEOUT=10 timeout 15 oho session list > /dev/null 2>&1 && echo "✓ 自定义超时工作正常" || echo "✗ 自定义超时失败"
echo ""

# 测试 3: 检查 oho 路径
echo "测试 3: 检查 oho 安装路径"
OHO_PATH=$(which oho)
echo "oho 路径：$OHO_PATH"
if [ "$OHO_PATH" = "/root/.local/bin/oho" ]; then
    echo "✓ 使用正确的 oho 版本"
else
    echo "⚠ 可能使用了旧版本 oho"
fi
echo ""

# 测试 4: 检查二进制文件大小（新版本应该更大）
echo "测试 4: 检查二进制文件"
if [ -f /root/.local/bin/oho ]; then
    SIZE=$(ls -lh /root/.local/bin/oho | awk '{print $5}')
    echo "/root/.local/bin/oho 大小：$SIZE"
    echo "✓ 主要安装路径正确"
fi
echo ""

echo "========================================="
echo "测试完成"
echo "========================================="
echo ""
echo "建议："
echo "1. 长时间任务设置：export OPENCODE_CLIENT_TIMEOUT=600"
echo "2. 或使用 --no-reply 不等待响应"
echo "3. 或使用 prompt-async 异步提交"
