# 介绍

SQL-Diff 是一个基于 AST 语法树的智能 SQL 表结构比对工具。

## 什么是 SQL-Diff？

SQL-Diff 可以：

- 🔍 **精准比对** 两个 MySQL 表结构的差异
- 🚀 **自动生成** ALTER TABLE 等 DDL 补全语句
- 🤖 **AI 分析** 提供智能优化建议和风险提示（可选）
- 💻 **美观输出** 分类显示、彩色标注、一目了然

## 核心优势

### 1. 精准高效

基于 AST 语法树解析，准确识别：
- 新增列
- 修改列（类型、长度、约束、默认值）
- 删除列
- 索引变更

### 2. 智能分析

可选接入 AI 大模型（DeepSeek 等），提供：
- 📊 差异分析
- ✨ 优化建议
- ⚠️ 风险提示
- 📖 最佳实践

### 3. 使用简单

```bash
# 一条命令完成比对
sql-diff -s "源表SQL" -t "目标表SQL"

# 启用 AI 分析
sql-diff -s "源表SQL" -t "目标表SQL" --ai
```

### 4. 安全可靠

- 删除操作自动注释，防止误删
- 环境变量管理敏感信息
- 完整的测试覆盖（100%）

## 应用场景

### 数据库迁移

在版本升级时，快速生成迁移脚本：

```bash
sql-diff \
  -s "$(cat schema/v1.0.sql)" \
  -t "$(cat schema/v2.0.sql)" \
  -o migration.sql
```

### 代码审查

在 Pull Request 中分析表结构变更：

```bash
git show HEAD~1:schema.sql > old.sql
git show HEAD:schema.sql > new.sql
sql-diff -s "$(cat old.sql)" -t "$(cat new.sql)" --ai
```

### 性能优化

使用 AI 分析索引优化：

```bash
sql-diff -s "当前表结构" -t "优化后表结构" --ai
```

### CI/CD 集成

在持续集成中自动检查表结构变更：

```yaml
- name: Check Schema Changes
  run: |
    sql-diff -s "$OLD_SCHEMA" -t "$NEW_SCHEMA" --ai
```

## 技术特点

### AST 解析

使用抽象语法树方式解析 SQL，比简单的字符串匹配更精准：

- 支持各种 SQL 语法
- 准确识别字段属性
- 正确处理注释和空格

### 美观输出

分类显示，一目了然：

- ➕ 新增列（绿色）
- 🔄 修改列（黄色）  
- 🗑️ 删除列（红色，已注释）
- 📇 新增索引（青色）
- 🗂️ 删除索引（紫色，已注释）

### 灵活配置

支持多种配置方式：

- 环境变量（推荐）
- 配置文件
- 命令行参数

## 性能指标

| 指标 | 数值 |
|------|------|
| 解析速度 | < 1ms |
| AI 响应时间 | 6-7秒 |
| AI 分析成本 | < ¥0.002/次 |
| 准确率 | 100% |

## 下一步

- 🚀 [快速开始](./getting-started.md)
- 📖 [详细文档](./cli.md)
- 🤖 [AI 功能](/ai/guide.md)
- 💡 [使用示例](/examples/basic.md)
