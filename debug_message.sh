#!/bin/bash
# =============================================================================
# OpenCode 消息提交调试脚本
# =============================================================================
# 用途：诊断为什么消息提交后没有 AI 响应
# 使用：export OPENCODE_SERVER_PASSWORD=your-password && ./debug_message.sh
# 参考：docs/oho-cli-usage/09-troubleshooting.md
# =============================================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置（支持环境变量覆盖）
SERVER_HOST="${OPENCODE_SERVER_HOST:-127.0.0.1}"
SERVER_PORT="${OPENCODE_SERVER_PORT:-4096}"
SERVER_PASSWORD="${OPENCODE_SERVER_PASSWORD:-}"
SERVER_USERNAME="${OPENCODE_SERVER_USERNAME:-opencode}"
BASE_URL="http://${SERVER_HOST}:${SERVER_PORT}"

echo "========================================="
echo "OpenCode 消息提交调试脚本"
echo "========================================="
echo "服务器地址：${BASE_URL}"
echo "用户名：${SERVER_USERNAME}"
echo ""

# 1. 检查服务器健康状态
echo -e "${YELLOW}[1/6] 检查服务器健康状态...${NC}"
if curl -s -f "${BASE_URL}/global/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 服务器正常运行${NC}"
else
    echo -e "${RED}✗ 服务器无法连接${NC}"
    echo "请确保 OpenCode 服务器正在运行：opencode serve"
    exit 1
fi
echo ""

# 2. 检查认证
echo -e "${YELLOW}[2/6] 检查认证配置...${NC}"
if [ -z "$SERVER_PASSWORD" ]; then
    echo -e "${RED}✗ 未设置 OPENCODE_SERVER_PASSWORD 环境变量${NC}"
    echo "请设置：export OPENCODE_SERVER_PASSWORD=your-password"
    exit 1
else
    echo -e "${GREEN}✓ 密码已配置${NC}"
fi

# 测试认证
AUTH_RESULT=$(curl -s -w "%{http_code}" -u "opencode:${SERVER_PASSWORD}" "${BASE_URL}/config" 2>/dev/null)
HTTP_CODE="${AUTH_RESULT: -3}"
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ 认证成功${NC}"
else
    echo -e "${RED}✗ 认证失败 (HTTP ${HTTP_CODE})${NC}"
    echo "请检查密码是否正确"
    exit 1
fi
echo ""

# 3. 列出会话
echo -e "${YELLOW}[3/6] 获取会话列表...${NC}"
SESSIONS=$(curl -s -u "opencode:${SERVER_PASSWORD}" "${BASE_URL}/session")
SESSION_COUNT=$(echo "$SESSIONS" | grep -o '"id"' | wc -l)
echo "找到 ${SESSION_COUNT} 个会话"

if [ "$SESSION_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}! 没有现有会话，将创建一个新会话${NC}"
    CREATE_RESP=$(curl -s -X POST -u "opencode:${SERVER_PASSWORD}" \
        -H "Content-Type: application/json" \
        -d '{"title":"debug-session"}' \
        "${BASE_URL}/session")
    SESSION_ID=$(echo "$CREATE_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "创建会话：${SESSION_ID}"
else
    SESSION_ID=$(echo "$SESSIONS" | grep -o '"id":"ses_[^"]*"' | head -1 | cut -d'"' -f4)
    echo "使用会话：${SESSION_ID}"
fi
echo ""

# 4. 检查会话状态
echo -e "${YELLOW}[4/6] 检查会话状态...${NC}"
STATUS_RESP=$(curl -s -u "opencode:${SERVER_PASSWORD}" "${BASE_URL}/session/status")
echo "会话状态响应："
echo "$STATUS_RESP" | head -c 500
echo ""

# 5. 测试消息提交（无响应模式）
echo -e "${YELLOW}[5/6] 测试消息提交（no-reply 模式）...${NC}"
MESSAGE_REQ='{
    "parts": [
        {
            "type": "text",
            "text": "这是一个测试消息，请回复 OK"'
        }
    ],
    "noReply": true
}'

MSG_RESP=$(curl -s -X POST -u "opencode:${SERVER_PASSWORD}" \
    -H "Content-Type: application/json" \
    -d "$MESSAGE_REQ" \
    "${BASE_URL}/session/${SESSION_ID}/message")

echo "消息提交响应："
echo "$MSG_RESP"

if echo "$MSG_RESP" | grep -q '"id"'; then
    MSG_ID=$(echo "$MSG_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo -e "${GREEN}✓ 消息提交成功，ID: ${MSG_ID}${NC}"
else
    echo -e "${RED}✗ 消息提交失败${NC}"
    echo "响应内容：$MSG_RESP"
    exit 1
fi
echo ""

# 6. 检查消息历史
echo -e "${YELLOW}[6/6] 检查消息历史...${NC}"
sleep 2  # 等待 2 秒让服务器处理

MSG_HISTORY=$(curl -s -u "opencode:${SERVER_PASSWORD}" \
    "${BASE_URL}/session/${SESSION_ID}/message?limit=5")

echo "最新消息历史："
echo "$MSG_HISTORY" | head -c 1000
echo ""

# 检查是否有 AI 响应
if echo "$MSG_HISTORY" | grep -q '"role":"assistant"'; then
    echo -e "${GREEN}✓ 检测到 AI 响应${NC}"
else
    echo -e "${YELLOW}! 未检测到 AI 响应${NC}"
    echo ""
    echo "可能的原因："
    echo "  1. AI 模型配置问题 - 检查 provider 配置"
    echo "  2. 会话被中止 - 检查会话状态"
    echo "  3. 权限请求等待确认 - 检查是否有权限弹窗"
    echo "  4. 服务器日志 - 查看 opencode serve 输出"
fi
echo ""

echo "========================================="
echo "调试完成"
echo "========================================="
echo ""
echo "建议的下一步："
echo "1. 检查 OpenCode 服务器日志：查看 opencode serve 的输出"
echo "2. 检查提供商配置：oho config providers"
echo "3. 检查会话权限：oho session permissions <session-id>"
echo "4. 尝试使用 oho CLI: oho message add -s ${SESSION_ID} '测试' --no-reply"
