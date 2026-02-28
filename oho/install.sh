#!/usr/bin/env bash
#
# oho 一键安装脚本
# 支持从源码编译或从 GitHub Releases 下载预编译二进制
#

set -e

# 配置
REPO_OWNER="tornado404"
REPO_NAME="opencode_cli"
BINARY_NAME="oho"
INSTALL_DIR="${HOME}/.local/bin"
CONFIG_DIR="${HOME}/.config/oho"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取最新版本
get_latest_version() {
    # 尝试从 GitHub API 获取最新 release
    if command -v curl &> /dev/null; then
        VERSION=$(curl -sSL https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest 2>/dev/null | grep -o '"tag_name":.*' | cut -d'"' -f4 || echo "")
    elif command -v wget &> /dev/null; then
        VERSION=$(wget -qO- https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest 2>/dev/null | grep -o '"tag_name":.*' | cut -d'"' -f4 || echo "")
    fi
    
    if [ -z "$VERSION" ]; then
        echo "${YELLOW}无法获取最新版本，将从源码编译${NC}"
        VERSION="build-from-source"
    fi
    
    echo "$VERSION"
}

# 检测操作系统和架构
detect_os_arch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7)
            ARCH="armv7"
            ;;
        *)
            echo "${RED}不支持的架构: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    # 标准化 OS 名称
    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            echo "${RED}不支持的操作系统: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo "${OS}-${ARCH}"
}

# 下载预编译二进制
download_binary() {
    local version="$1"
    local os_arch="$2"
    local download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${version}/${BINARY_NAME}-${os_arch}"
    local temp_file="/tmp/${BINARY_NAME}-${os_arch}"
    
    echo "${YELLOW}下载中: ${download_url}${NC}"
    
    if command -v curl &> /dev/null; then
        curl -L -o "$temp_file" "$download_url"
    elif command -v wget &> /dev/null; then
        wget -O "$temp_file" "$download_url"
    else
        echo "${RED}错误: 需要 curl 或 wget${NC}"
        return 1
    fi
    
    if [ $? -ne 0 ]; then
        echo "${RED}下载失败${NC}"
        rm -f "$temp_file"
        return 1
    fi
    
    # 移动到安装目录
    mkdir -p "$INSTALL_DIR"
    mv "$temp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    echo "${GREEN}安装成功!${NC}"
    echo "二进制文件: ${INSTALL_DIR}/${BINARY_NAME}"
}

# 从源码编译
build_from_source() {
    echo "${YELLOW}从源码编译...${NC}"
    
    # 检查 Go 是否安装
    if ! command -v go &> /dev/null; then
        echo "${RED}错误: 需要 Go 1.21+ 来编译${NC}"
        echo "请访问 https://golang.org/dl/ 安装 Go"
        exit 1
    fi
    
    # 检查 Go 版本
    GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
    GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
    GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)
    
    if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 21 ]); then
        echo "${RED}错误: 需要 Go 1.21+${NC}"
        exit 1
    fi
    
    # 创建临时目录
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # 克隆仓库
    echo "克隆仓库..."
    if command -v git &> /dev/null; then
        git clone --depth 1 "https://github.com/${REPO_OWNER}/${REPO_NAME}.git"
    else
        echo "${RED}错误: 需要 git${NC}"
        exit 1
    fi
    
    cd "${REPO_NAME}/oho"
    
    # 编译
    echo "编译中..."
    go build -o "${INSTALL_DIR}/${BINARY_NAME}" ./cmd
    
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    # 清理
    cd /
    rm -rf "$TEMP_DIR"
    
    echo "${GREEN}编译安装成功!${NC}"
    echo "二进制文件: ${INSTALL_DIR}/${BINARY_NAME}"
}

# 创建默认配置
create_config() {
    mkdir -p "$CONFIG_DIR"
    
    if [ ! -f "${CONFIG_DIR}/config.json" ]; then
        cat > "${CONFIG_DIR}/config.json" << EOF
{
  "host": "127.0.0.1",
  "port": 4096,
  "username": "opencode",
  "password": "",
  "json": false
}
EOF
        echo "${GREEN}配置文件已创建: ${CONFIG_DIR}/config.json${NC}"
    fi
}

# 添加到 PATH（如果需要）
add_to_path() {
    # 检查是否已经在 PATH 中
    if [[ ":$PATH:" == *":${INSTALL_DIR}:"* ]]; then
        return
    fi
    
    # 检测 shell 配置文件
    local shell_rc=""
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="${HOME}/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="${HOME}/.bashrc"
    fi
    
    if [ -n "$shell_rc" ] && [ -f "$shell_rc" ]; then
        # 检查是否已经添加过
        if ! grep -q "oho.*local/bin" "$shell_rc"; then
            echo "" >> "$shell_rc"
            echo "# oho CLI" >> "$shell_rc"
            echo "export PATH=\"\${HOME}/.local/bin:\$PATH\"" >> "$shell_rc"
            echo "${YELLOW}已将 ${INSTALL_DIR} 添加到 PATH${NC}"
            echo "请运行: source $shell_rc"
        fi
    fi
}

# 验证安装
verify_install() {
    if "${INSTALL_DIR}/${BINARY_NAME}" --version &> /dev/null; then
        echo "${GREEN}验证成功!${NC}"
        "${INSTALL_DIR}/${BINARY_NAME}" --version
        return 0
    else
        echo "${RED}验证失败${NC}"
        return 1
    fi
}

# 显示帮助
show_help() {
    cat << EOF
oho 一键安装脚本

用法: $0 [选项]

选项:
    -h, --help              显示帮助信息
    -v, --version VERSION   指定版本 (默认: latest)
    -d, --dir DIRECTORY     指定安装目录 (默认: ~/.local/bin)
    -s, --source            强制从源码编译
    -c, --check             仅验证安装

示例:
    $0                      # 安装最新版本
    $0 -s                   # 从源码编译
    $0 -d /usr/local/bin    # 安装到指定目录

EOF
}

# 主函数
main() {
    local install_dir="$INSTALL_DIR"
    local version=""
    local force_source=false
    local check_only=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                version="$2"
                shift 2
                ;;
            -d|--dir)
                install_dir="$2"
                shift 2
                ;;
            -s|--source)
                force_source=true
                shift
                ;;
            -c|--check)
                check_only=true
                shift
                ;;
            *)
                echo "${RED}未知选项: $1${NC}"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 验证模式
    if [ "$check_only" = true ]; then
        if [ -x "${install_dir}/${BINARY_NAME}" ]; then
            echo "${GREEN}oho 已安装${NC}"
            "${install_dir}/${BINARY_NAME}" --version
            exit 0
        else
            echo "${RED}oho 未安装${NC}"
            exit 1
        fi
    fi
    
    echo "=========================================="
    echo "         oho 安装脚本"
    echo "=========================================="
    echo ""
    
    # 设置安装目录
    INSTALL_DIR="$install_dir"
    
    # 获取版本
    if [ -z "$version" ]; then
        version=$(get_latest_version)
    fi
    
    echo "版本: $version"
    echo "安装目录: $INSTALL_DIR"
    echo ""
    
    # 尝试下载或编译
    if [ "$force_source" = true ]; then
        build_from_source
    else
        os_arch=$(detect_os_arch)
        
        # 尝试下载
        if ! download_binary "$version" "$os_arch"; then
            echo "${YELLOW}下载失败，尝试从源码编译...${NC}"
            build_from_source
        fi
    fi
    
    # 创建配置
    create_config
    
    # 添加到 PATH
    add_to_path
    
    # 验证安装
    echo ""
    echo "=========================================="
    verify_install
    echo "=========================================="
    echo ""
    echo "使用说明:"
    echo "  1. 确保 ${INSTALL_DIR} 在 PATH 中"
    echo "  2. 配置 OpenCode Server 连接信息:"
    echo "       export OPENCODE_SERVER_HOST=127.0.0.1"
    echo "       export OPENCODE_SERVER_PORT=4096"
    echo "       export OPENCODE_SERVER_PASSWORD=your-password"
    echo "  3. 运行: oho --help"
    echo ""
    echo "快速开始:"
    echo "  oho global health          # 检查服务器状态"
    echo "  oho session list          # 列出所有会话"
    echo "  oho mcpserver             # 启动 MCP 服务器"
    echo ""
}

# 运行主函数
main "$@"
