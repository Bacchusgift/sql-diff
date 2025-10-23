#!/bin/bash

# SQL-Diff 快速配置脚本
# 使用此脚本快速设置环境变量

echo "======================================"
echo "SQL-Diff 快速配置向导"
echo "======================================"
echo ""

# 检测是否已有配置
if [ -n "$SQL_DIFF_AI_API_KEY" ]; then
    echo "✓ 检测到已配置的环境变量"
    echo ""
    echo "当前配置:"
    echo "  AI 启用: $SQL_DIFF_AI_ENABLED"
    echo "  提供商: $SQL_DIFF_AI_PROVIDER"
    echo "  API Key: ${SQL_DIFF_AI_API_KEY:0:6}...${SQL_DIFF_AI_API_KEY: -4}"
    echo ""
    read -p "是否重新配置？(y/N) " reconfigure
    if [[ ! "$reconfigure" =~ ^[Yy]$ ]]; then
        exit 0
    fi
    echo ""
fi

# 选择 AI 提供商
echo "选择 AI 提供商:"
echo "  1) DeepSeek (推荐)"
echo "  2) OpenAI"
echo ""
read -p "请选择 [1]: " provider_choice

case $provider_choice in
    2)
        PROVIDER="openai"
        ENDPOINT="https://api.openai.com/v1"
        MODEL="gpt-4"
        ;;
    *)
        PROVIDER="deepseek"
        ENDPOINT="https://api.deepseek.com/v1"
        MODEL="deepseek-chat"
        ;;
esac

echo ""
echo "已选择: $PROVIDER"
echo ""

# 输入 API Key
read -p "请输入 API Key: " API_KEY
if [ -z "$API_KEY" ]; then
    echo "❌ API Key 不能为空"
    exit 1
fi

echo ""
echo "✓ 配置完成！"
echo ""

# 生成配置
cat << EOF
# 将以下内容添加到 ~/.bashrc 或 ~/.zshrc

export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_PROVIDER=$PROVIDER
export SQL_DIFF_AI_API_KEY=$API_KEY
export SQL_DIFF_AI_ENDPOINT=$ENDPOINT
export SQL_DIFF_AI_MODEL=$MODEL
export SQL_DIFF_AI_TIMEOUT=30
EOF

echo ""
echo "======================================"
echo "使用方法:"
echo "======================================"
echo ""
echo "1. 临时使用（仅当前终端）:"
echo "   eval \"\$(./setup-env.sh)\""
echo ""
echo "2. 永久保存（推荐）:"
echo "   ./setup-env.sh >> ~/.bashrc"
echo "   source ~/.bashrc"
echo ""
echo "3. 或者使用工具命令:"
echo "   sql-diff config --ai-enabled --provider $PROVIDER --api-key YOUR_KEY >> ~/.bashrc"
echo ""
echo "4. 验证配置:"
echo "   sql-diff config --show"
echo ""
