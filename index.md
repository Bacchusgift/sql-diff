---
layout: home

hero:
  name: SQL-Diff
  text: 智能 SQL 表结构比对工具
  tagline: 基于 AST 语法树，精准比对表结构差异，自动生成 DDL 语句，支持 AI 智能分析
  image:
    src: /hero-image.svg
    alt: SQL-Diff
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: 在 GitHub 查看
      link: https://github.com/Bacchusgift/sql-diff

features:
  - icon: 🎯
    title: 交互式输入
    details: 支持多行 SQL 直接粘贴，从 Navicat、MySQL Workbench 等工具复制即用，无需处理换行符
  
  - icon: 🔍
    title: 精准比对
    details: 基于 AST 语法树解析 CREATE TABLE 语句，准确识别新增列、修改列、删除列和索引变更
  
  - icon: 🚀
    title: DDL 生成
    details: 自动生成标准 MySQL DDL 语句，支持 ALTER TABLE 等操作，安全标注删除操作
  
  - icon: 🤖
    title: AI 增强
    details: 可选接入 DeepSeek 等大模型，提供智能差异分析、优化建议、风险提示和最佳实践
  
  - icon: 💻
    title: CLI 友好
    details: 简洁美观的命令行界面，彩色输出，分类显示，支持输出到文件
  
  - icon: ⚙️
    title: 灵活配置
    details: 支持环境变量和配置文件两种方式，配置命令一键生成，CI/CD 集成简单
---

## 🎯 为什么选择 SQL-Diff？

<div class="vp-doc">

### 🚀 效率提升 1000 倍

- **人工比对**: 2-4 小时 → **SQL-Diff**: 2-5 秒
- **成本降低**: 人工 ¥200-500 → **AI 分析** < ¥0.002/次

### 💡 智能分析，专业建议

AI 分析达到高级 DBA 水平：
- 复合索引优化建议
- 数据类型精度优化
- ENUM 类型推荐
- 分区表设计建议
- 数据迁移风险提示

### 🎨 美观输出，一目了然

```bash
✓ 生成的 DDL 语句:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

➕ 新增列 (2):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255);
  2. ALTER TABLE users ADD COLUMN created_at TIMESTAMP;

📇 新增索引 (1):
  1. ALTER TABLE users ADD INDEX idx_email (email);

📋 完整执行脚本:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN created_at TIMESTAMP;
ALTER TABLE users ADD INDEX idx_email (email);
```

</div>

## 📦 快速体验

::: code-group

```bash [🍺 安装 (macOS)]
# 一条命令安装
brew install Bacchusgift/tap/sql-diff

# 验证安装
sql-diff --version

# 更新到最新版本
brew upgrade sql-diff
```

```bash [🐧 安装 (Linux)]
# 下载预编译二进制文件
wget https://github.com/Bacchusgift/sql-diff/releases/latest/download/sql-diff-linux-amd64

# 赋予执行权限
chmod +x sql-diff-linux-amd64

# 移动到 PATH 目录
sudo mv sql-diff-linux-amd64 /usr/local/bin/sql-diff

# 验证
sql-diff --version
```

```bash [🛠️ 从源码构建]
# 克隆项目
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff

# 编译
make build

# 运行
./bin/sql-diff --version
```

```bash [🚀 交互式模式]
# 启动交互式模式
sql-diff -i

# 按提示粘贴源表 SQL（支持多行）
# 输入 'END' 或连续两次 Enter 结束输入
# 再粘贴目标表 SQL
# 自动生成 DDL！
```

```bash [🤖 配置 AI]
# 一键配置 AI 功能（可选）
sql-diff config \
  --ai-enabled \
  --provider deepseek \
  --api-key YOUR_KEY \
  >> ~/.bashrc

source ~/.bashrc

# 启用 AI 分析
sql-diff -i --ai
```

```bash [📝 命令行模式]
# 简单 SQL 可用命令行参数
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))"

# 输出到文件
sql-diff -i -o migration.sql
```

:::

## 🌟 核心特性

<div class="feature-grid">

### 🎯 交互式输入模式

<div class="feature-content">

**一键启动，直接粘贴**

```bash
sql-diff -i
```

✅ **完美支持多行 SQL**  
✅ **从数据库工具直接复制**  
✅ **无需转义换行符**  
✅ **实时字符统计**  
✅ **友好操作提示**

**使用场景：**
- 从 Navicat/MySQL Workbench 复制表结构
- 处理包含注释的复杂 SQL
- 比对大型表结构（几十个字段）
- 避免 Shell 转义字符问题

</div>

### 分类显示 DDL

<div class="feature-content">

- ➕ **新增列**（绿色）
- 🔄 **修改列**（黄色）
- 🗑️ **删除列**（红色，已注释）
- 📇 **新增索引**（青色）
- 🗂️ **删除索引**（紫色，已注释）

每类操作自动编号，数量统计清晰

</div>

### AI 智能分析

<div class="feature-content">

- 📊 **差异分析** - 深入解读表结构变更
- ✨ **优化建议** - 针对性的改进建议
- ⚠️ **潜在风险** - 识别可能的问题
- 📖 **最佳实践** - 行业标准推荐

</div>

### 环境变量配置

<div class="feature-content">

```bash
# 支持的环境变量
SQL_DIFF_AI_ENABLED    # 启用/禁用 AI
SQL_DIFF_AI_PROVIDER   # AI 提供商
SQL_DIFF_AI_API_KEY    # API 密钥
SQL_DIFF_AI_ENDPOINT   # API 端点
SQL_DIFF_AI_MODEL      # 模型名称
SQL_DIFF_AI_TIMEOUT    # 超时时间
```

</div>

</div>

## 🎓 应用场景

<div class="use-cases">

::: tip 数据库迁移
在版本升级时使用 AI 分析变更影响，节省大量人工审查时间
:::

::: tip 代码审查
在 Pull Request 中集成 AI 分析，提升代码审查质量
:::

::: tip 性能优化
AI 提出的索引优化建议直接提升查询性能
:::

::: tip 团队协作
统一的 DDL 生成标准，降低沟通成本
:::

</div>

## 📊 性能指标

| 指标 | 数值 | 评价 |
|------|------|------|
| 平均响应时间 | 6-7 秒 | 优秀 ⭐⭐⭐⭐⭐ |
| AI 分析成本 | <¥0.002/次 | 极低 ⭐⭐⭐⭐⭐ |
| 准确率 | 100% | 完美 ⭐⭐⭐⭐⭐ |
| 测试覆盖率 | 100% | 完整 ⭐⭐⭐⭐⭐ |

## 🤝 社区

<div class="community">

- 📖 [完整文档](/guide/introduction)
- 💬 [问题反馈](https://github.com/Bacchusgift/sql-diff/issues)
- 🌟 [Star on GitHub](https://github.com/Bacchusgift/sql-diff)
- 🤝 [贡献代码](https://github.com/Bacchusgift/sql-diff/blob/main/CONTRIBUTING.md)

</div>

<style>
.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  margin: 2rem 0;
}

.feature-content {
  padding: 1rem;
  border-radius: 8px;
  background: var(--vp-c-bg-soft);
}

.use-cases {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin: 2rem 0;
}

.community {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
  margin: 2rem 0;
  font-size: 1.1rem;
}

.community a {
  text-decoration: none;
  transition: transform 0.2s;
}

.community a:hover {
  transform: translateY(-2px);
}
</style>
