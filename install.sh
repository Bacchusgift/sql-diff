#!/bin/bash

# SQL-Diff 智能安装脚本
# 自动检测网络环境，选择最佳安装方式

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 版本号
VERSION="v1.0.2"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}       SQL-Diff 智能安装脚本${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 检测操作系统和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}✗ 不支持的架构: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}✓ 检测到系统: $OS-$ARCH${NC}"
echo ""

# 检测网络环境
echo -e "${YELLOW}⏳ 检测网络环境...${NC}"

if curl -s --max-time 3 https://api.github.com > /dev/null 2>&1; then
    echo -e "${GREEN}✓ GitHub 访问正常，使用国际源${NC}"
    USE_MIRROR=false
else
    echo -e "${YELLOW}⚠ GitHub 访问较慢，建议使用镜像源${NC}"
    USE_MIRROR=true
fi

echo ""

# 构建下载 URL
BINARY_NAME="sql-diff-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

if [ "$USE_MIRROR" = true ]; then
    # 使用 Gitee 镜像（如果已配置）或 GitHub 代理
    echo -e "${BLUE}📦 选择下载源...${NC}"
    echo "1) GitHub 官方源（可能较慢）"
    echo "2) ghproxy.com 代理加速"
    echo "3) 手动指定 URL"
    
    read -p "请选择 (1/2/3，默认 2): " choice
    choice=${choice:-2}
    
    case $choice in
        1)
            DOWNLOAD_URL="https://github.com/Bacchusgift/sql-diff/releases/download/${VERSION}/${BINARY_NAME}"
            ;;
        2)
            DOWNLOAD_URL="https://ghproxy.com/https://github.com/Bacchusgift/sql-diff/releases/download/${VERSION}/${BINARY_NAME}"
            ;;
        3)
            read -p "请输入下载 URL: " DOWNLOAD_URL
            ;;
        *)
            echo -e "${RED}✗ 无效选择${NC}"
            exit 1
            ;;
    esac
else
    DOWNLOAD_URL="https://github.com/Bacchusgift/sql-diff/releases/download/${VERSION}/${BINARY_NAME}"
fi

echo ""
echo -e "${BLUE}📥 开始下载...${NC}"
echo -e "${YELLOW}URL: ${DOWNLOAD_URL}${NC}"
echo ""

# 下载文件
TMP_FILE="/tmp/sql-diff-${BINARY_NAME}"

if command -v wget > /dev/null 2>&1; then
    wget -O "$TMP_FILE" "$DOWNLOAD_URL" || {
        echo -e "${RED}✗ 下载失败${NC}"
        exit 1
    }
elif command -v curl > /dev/null 2>&1; then
    curl -L -o "$TMP_FILE" "$DOWNLOAD_URL" || {
        echo -e "${RED}✗ 下载失败${NC}"
        exit 1
    }
else
    echo -e "${RED}✗ 需要 wget 或 curl 命令${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 下载完成${NC}"
echo ""

# 赋予执行权限
chmod +x "$TMP_FILE"

# 安装到系统路径
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="${INSTALL_DIR}/sql-diff"

echo -e "${BLUE}📦 安装到 ${INSTALL_PATH}...${NC}"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$INSTALL_PATH"
else
    echo -e "${YELLOW}⚠ 需要 sudo 权限${NC}"
    sudo mv "$TMP_FILE" "$INSTALL_PATH"
fi

echo -e "${GREEN}✓ 安装成功！${NC}"
echo ""

# 验证安装
if command -v sql-diff > /dev/null 2>&1; then
    VERSION_OUTPUT=$(sql-diff --version 2>&1 || echo "unknown")
    echo -e "${GREEN}✓ 验证安装:${NC}"
    echo -e "${BLUE}   $VERSION_OUTPUT${NC}"
else
    echo -e "${RED}✗ 安装验证失败${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 安装完成！${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}快速开始:${NC}"
echo -e "  ${BLUE}sql-diff -i${NC}             # 交互式模式"
echo -e "  ${BLUE}sql-diff --help${NC}         # 查看帮助"
echo ""
echo -e "${YELLOW}完整文档:${NC}"
echo -e "  ${BLUE}https://bacchusgift.github.io/sql-diff/${NC}"
echo ""
