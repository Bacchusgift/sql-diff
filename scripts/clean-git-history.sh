#!/bin/bash
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}       清除 Git 历史中的内部文档${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 要从 Git 历史中彻底删除的文件和目录
# 注意：只删除真正的内部文档！
FILES_TO_REMOVE=(
    ".UPDATE_SUMMARY.md"              # 项目更新摘要（内部跟踪）
    ".AI_SQL_GENERATION_DOC.md"       # AI 相关内部文档
    "docs/AI_GUIDE.md"                # 早期 AI 指南（已迁移到 VitePress）
    "docs/ARCHITECTURE.md"            # 早期架构文档（已迁移）
    "docs/ENV_CONFIG.md"              # 早期环境配置（已迁移）
    "docs/EXAMPLES.md"                # 早期示例（已迁移）
    "docs/QUICKSTART.md"              # 早期快速开始（已迁移）
    "docs/"                           # 早期文档目录
    "Formula/"                        # Homebrew Formula（已迁移到独立仓库）
)

echo -e "${YELLOW}⚠️  警告：此操作将重写 Git 历史！${NC}"
echo ""
echo "将从 Git 历史中彻底删除以下文件/目录："
echo ""
for file in "${FILES_TO_REMOVE[@]}"; do
    echo "  ❌ $file"
done
echo ""
echo -e "${RED}注意事项：${NC}"
echo "  1. 这将重写整个 Git 历史"
echo "  2. 所有分支和标签都会受影响"
echo "  3. 需要强制推送到远程仓库 (git push --force)"
echo "  4. 其他协作者需要重新克隆仓库"
echo ""
echo -e "${YELLOW}建议：${NC}"
echo "  1. 确保已备份重要数据"
echo "  2. 通知所有协作者"
echo "  3. 在私人项目或独自维护的项目上执行"
echo ""

read -p "确认要继续吗？输入 'yes' 继续，其他任意键取消: " -r
echo ""

if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${GREEN}✅ 操作已取消${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}🔍 检查工具...${NC}"

# 检查是否安装了 git-filter-repo
if ! command -v git-filter-repo &> /dev/null; then
    echo -e "${YELLOW}⚠️  未安装 git-filter-repo，正在安装...${NC}"
    
    if command -v brew &> /dev/null; then
        echo "使用 Homebrew 安装..."
        brew install git-filter-repo
    elif command -v pip3 &> /dev/null; then
        echo "使用 pip3 安装..."
        pip3 install git-filter-repo
    elif command -v pip &> /dev/null; then
        echo "使用 pip 安装..."
        pip install git-filter-repo
    else
        echo -e "${RED}❌ 无法自动安装 git-filter-repo${NC}"
        echo "请手动安装："
        echo "  macOS: brew install git-filter-repo"
        echo "  Linux: pip3 install git-filter-repo"
        exit 1
    fi
fi

echo -e "${GREEN}✅ git-filter-repo 已就绪${NC}"
echo ""

# 创建备份
BACKUP_DIR="../sql-diff-backup-$(date +%Y%m%d_%H%M%S)"
echo -e "${BLUE}📦 创建备份到: $BACKUP_DIR${NC}"
cp -r . "$BACKUP_DIR"
echo -e "${GREEN}✅ 备份完成${NC}"
echo ""

# 确保在 Git 仓库中
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ 错误：当前目录不是 Git 仓库${NC}"
    exit 1
fi

# 构建 git-filter-repo 命令参数
FILTER_ARGS=""
for file in "${FILES_TO_REMOVE[@]}"; do
    FILTER_ARGS="$FILTER_ARGS --path $file --invert-paths"
done

echo -e "${BLUE}🗑️  开始清除 Git 历史...${NC}"
echo ""

# 执行过滤
git-filter-repo $FILTER_ARGS --force

echo ""
echo -e "${GREEN}✅ Git 历史清理完成！${NC}"
echo ""

# 显示清理结果
echo -e "${BLUE}📊 清理结果：${NC}"
echo ""
git log --oneline --graph --all -10
echo ""

echo -e "${BLUE}📦 仓库大小变化：${NC}"
du -sh .git
echo ""

echo -e "${YELLOW}📝 下一步操作：${NC}"
echo ""
echo "1. 检查本地仓库状态："
echo "   git log --all --oneline"
echo ""
echo "2. 添加远程仓库（filter-repo 会删除远程）："
echo "   git remote add origin https://github.com/Bacchusgift/sql-diff.git"
echo ""
echo "3. 强制推送到远程（⚠️  这将重写远程历史）："
echo "   git push origin main --force"
echo ""
echo "4. 推送所有分支和标签："
echo "   git push origin --all --force"
echo "   git push origin --tags --force"
echo ""
echo -e "${RED}⚠️  警告：${NC}"
echo "   - 其他协作者需要重新克隆仓库"
echo "   - 已有的 fork 需要重新同步"
echo "   - CI/CD 可能需要重新触发"
echo ""
echo -e "${GREEN}✅ 如果遇到问题，可以从备份恢复：${NC}"
echo "   cd .."
echo "   rm -rf sql-diff"
echo "   mv $BACKUP_DIR sql-diff"
echo ""
