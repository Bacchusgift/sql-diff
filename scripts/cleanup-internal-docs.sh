#!/bin/bash
set -e

echo "🧹 清理内部开发文档"
echo "===================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查当前分支
CURRENT_BRANCH=$(git branch --show-current)
echo -e "${YELLOW}当前分支: $CURRENT_BRANCH${NC}"
echo ""

# 列出要删除的文件
echo "📋 将从 Git 历史中删除以下文件："
echo ""

FILES_TO_REMOVE=(
    ".AI_SQL_GENERATION_DOC.md"
)

# 检查这些文件是否存在于 Git 中
EXISTING_FILES=()
for file in "${FILES_TO_REMOVE[@]}"; do
    if git ls-files | grep -q "^${file}$"; then
        echo "  ❌ $file"
        EXISTING_FILES+=("$file")
    fi
done

echo ""

if [ ${#EXISTING_FILES[@]} -eq 0 ]; then
    echo -e "${GREEN}✅ 没有需要清理的文件${NC}"
    exit 0
fi

# 确认操作
echo -e "${RED}⚠️  警告：此操作将从 Git 历史中永久删除这些文件！${NC}"
echo ""
read -p "是否继续？(yes/no): " -r
echo ""

if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "操作已取消"
    exit 0
fi

# 从 Git 中删除文件（保留在本地，只从 Git 追踪中移除）
echo "🗑️  从 Git 追踪中移除文件..."
for file in "${EXISTING_FILES[@]}"; do
    git rm --cached "$file" 2>/dev/null || true
    echo "  ✓ Removed from tracking: $file"
done

echo ""
echo -e "${GREEN}✅ 清理完成！${NC}"
echo ""
echo "📝 下一步操作："
echo "  1. 提交变更:"
echo "     git commit -m 'chore: remove internal development docs from Git tracking'"
echo ""
echo "  2. 推送到远程:"
echo "     git push origin $CURRENT_BRANCH"
echo ""
echo "  3. 清理本地文件（如果需要）:"
echo "     rm -f .AI_SQL_GENERATION_DOC.md"
echo ""
