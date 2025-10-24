# 帮助信息优化说明

## 🎯 优化目标

让用户更容易发现和使用新增的 AI 自然语言生成 SQL 功能。

## 📊 改进对比

### 主命令帮助信息（`sql-diff --help`）

#### ✅ 改进后

```
sql-diff 是一个基于 AST 的 SQL 表结构比对工具。

核心功能：
  • 精准比对表结构差异，自动生成 DDL 语句
  • AI 自然语言生成 SQL（需配置 AI）
  • 交互式模式，支持多行粘贴
  • 可选 AI 智能分析，提供优化建议

Examples:
  # 1️⃣  表结构比对（交互式，推荐）
  sql-diff -i
  
  # 2️⃣  AI 生成 CREATE TABLE
  sql-diff generate -d "创建用户表：ID、用户名、邮箱、密码"
  
  # 3️⃣  AI 生成 ALTER TABLE（交互式）
  sql-diff alter -i -d "添加手机号字段"
  
  # 4️⃣  命令行模式比对
  sql-diff -s "CREATE TABLE users (id INT)" -t "CREATE TABLE users (id INT, name VARCHAR(100))"
  
  # 5️⃣  启用 AI 分析
  sql-diff -i --ai
  
  # 查看详细帮助
  sql-diff generate --help
  sql-diff alter --help
```

**改进点**：
1. ✅ 在 Long 描述中列出了 4 个核心功能
2. ✅ 使用数字标记（1️⃣2️⃣3️⃣）组织示例
3. ✅ **新增**：AI 生成 CREATE TABLE 示例
4. ✅ **新增**：AI 生成 ALTER TABLE 示例
5. ✅ 提示用户查看详细帮助

#### ❌ 改进前

```
sql-diff 是一个基于 AST 的 SQL 表结构比对工具。

可以比对两个表结构的差异，并自动生成 DDL 补全语句。
支持可选的 AI 智能分析功能，提供优化建议。

Examples:
  # 交互式模式（推荐，支持多行粘贴）
  sql-diff -i
  
  # 基础用法
  sql-diff -s "CREATE TABLE users (id INT)" -t "CREATE TABLE users (id INT, name VARCHAR(100))"
  
  # 启用 AI 分析
  sql-diff -s "..." -t "..." --ai
  
  # 交互式 + AI
  sql-diff -i --ai
  
  # 输出到文件
  sql-diff -i -o output.sql
```

**问题**：
- ❌ 没有提到 AI 生成 SQL 功能
- ❌ 示例中没有展示 generate 和 alter 命令
- ❌ 功能描述过于简单

---

### generate 命令帮助（`sql-diff generate --help`）

#### ✅ 改进后

```
使用 AI 根据自然语言描述生成标准的 MySQL CREATE TABLE 语句。

功能特点：
  ✓ 自动推断字段类型（VARCHAR、INT、DECIMAL、DATETIME 等）
  ✓ 自动添加主键、索引、唯一约束
  ✓ 应用 MySQL 最佳实践（InnoDB、UTF8MB4、注释等）
  ✓ 使用标准命名规范（snake_case）
  ✓ 支持输出到文件

注意：此功能需要启用 AI（配置 .sql-diff-config.yaml 或使用 --ai 参数）

Examples:
  # 基础用法
  sql-diff generate -d "创建用户表，包含 ID、用户名、邮箱、密码、创建时间"
  
  # 复杂示例
  sql-diff generate -d "创建订单表：订单号（唯一）、用户ID（外键）、订单金额（精确到分）、订单状态、下单时间"
  
  # 输出到文件
  sql-diff generate -d "创建商品表：商品ID、名称、价格、库存" -o product.sql
  
  # 启用 AI（如果配置文件中未启用）
  sql-diff generate --ai -d "创建博客文章表"
```

**改进点**：
1. ✅ 列出 5 个功能特点
2. ✅ 明确说明需要 AI 配置
3. ✅ 提供从简单到复杂的示例
4. ✅ 展示输出到文件的用法

#### ❌ 改进前

```
使用 AI 根据自然语言描述生成标准的 MySQL CREATE TABLE 语句。

示例：
  sql-diff generate -d "创建一个用户表，包含 ID、用户名、邮箱、密码、创建时间"
  sql-diff generate -d "创建商品表：商品ID、名称、价格、库存、分类、状态" -o product.sql
```

**问题**：
- ❌ 没有说明功能特点
- ❌ 没有提示 AI 配置要求
- ❌ 示例较少

---

### alter 命令帮助（`sql-diff alter --help`）

#### ✅ 改进后

```
使用 AI 根据现有表结构和自然语言描述生成 MySQL ALTER TABLE 语句。

功能特点：
  ✓ 基于现有表结构生成变更语句
  ✓ 支持交互式输入表结构（推荐）
  ✓ 可生成多条 ALTER 语句
  ✓ 考虑数据迁移安全性
  ✓ 保持与现有表的一致性

注意：此功能需要启用 AI（配置 .sql-diff-config.yaml 或使用 --ai 参数）

Examples:
  # 交互式模式（推荐，支持粘贴大型表结构）
  sql-diff alter -i -d "添加手机号字段、邮箱改为唯一索引"
  
  # 命令行模式
  sql-diff alter -t "CREATE TABLE users (id INT, name VARCHAR(100));" -d "添加邮箱字段"
  
  # 复杂修改
  sql-diff alter -i -d "添加实名认证字段：真实姓名、身份证号、认证状态、认证时间"
  
  # 输出到文件
  sql-diff alter -i -d "添加商品状态字段" -o alter_product.sql
```

**改进点**：
1. ✅ 列出 5 个功能特点
2. ✅ 明确说明需要 AI 配置
3. ✅ 强调交互式模式（推荐）
4. ✅ 提供多种使用场景示例

#### ❌ 改进前

```
使用 AI 根据现有表结构和自然语言描述生成 MySQL ALTER TABLE 语句。

示例：
  # 命令行模式
  sql-diff alter -t "CREATE TABLE users ..." -d "添加手机号字段、邮箱改为唯一索引"
  
  # 交互式模式
  sql-diff alter -i -d "添加商品状态字段，默认值为上架"
```

**问题**：
- ❌ 没有说明功能特点
- ❌ 没有提示 AI 配置要求
- ❌ 示例较少
- ❌ 没有强调交互式模式的优势

---

## 📝 改进总结

### 主要改进

1. **突出新功能**：在主帮助中显著展示 AI 生成 SQL 功能
2. **结构化描述**：使用 ✓ 和 • 符号列出功能特点
3. **分类示例**：使用数字标记组织不同使用场景
4. **明确要求**：说明需要配置 AI 才能使用
5. **引导探索**：提示用户查看子命令的详细帮助

### 用户体验提升

| 改进项 | 改进前 | 改进后 |
|--------|--------|--------|
| 发现新功能 | 需要滚动到 Available Commands | 首屏可见，示例 2️⃣3️⃣ |
| 了解功能特点 | 无 | 列出 5 个特点 |
| 知道如何配置 | 无提示 | 明确说明 AI 配置要求 |
| 示例丰富度 | 2 个示例 | 4-5 个示例 |
| 使用引导 | 无 | 标注推荐方式、难度 |

### 关键改进点

1. **主帮助信息**：
   - ✅ 在核心功能中列出"AI 自然语言生成 SQL"
   - ✅ 示例 2️⃣ 展示 generate 命令
   - ✅ 示例 3️⃣ 展示 alter 命令

2. **generate 命令**：
   - ✅ 功能特点清单（5 项）
   - ✅ AI 配置提示
   - ✅ 从简单到复杂的示例

3. **alter 命令**：
   - ✅ 功能特点清单（5 项）
   - ✅ 强调交互式模式
   - ✅ 多场景示例

## 🎯 使用建议

### 用户发现流程

1. **新用户**：
   ```bash
   sql-diff --help
   # 看到示例 2️⃣ 和 3️⃣，了解到 AI 生成功能
   ```

2. **进一步了解**：
   ```bash
   sql-diff generate --help
   # 查看详细功能特点和示例
   ```

3. **开始使用**：
   ```bash
   sql-diff generate -d "创建用户表：ID、用户名、邮箱"
   ```

### 快速参考

```bash
# 查看所有命令
sql-diff --help

# 查看子命令详细帮助
sql-diff generate --help
sql-diff alter --help
sql-diff config --help
sql-diff version --help

# 快速体验
sql-diff -i                          # 表结构比对
sql-diff generate -d "创建用户表"     # AI 生成表
sql-diff alter -i -d "添加字段"       # AI 生成 DDL
```

## 📚 相关文档

- [README.md](./README.md) - 项目主文档
- [AI_SQL_GENERATOR.md](./AI_SQL_GENERATOR.md) - AI 功能详细文档
- [VERSION_INFO.md](./VERSION_INFO.md) - 版本信息功能
- [快速开始指南](./guide/getting-started.md) - 新手入门

## 🎉 总结

通过这次优化，用户可以：

1. ✅ 在主帮助信息中直接看到 AI 生成功能
2. ✅ 通过清晰的示例了解如何使用
3. ✅ 获得功能特点和配置要求的明确说明
4. ✅ 找到从简单到复杂的使用示例

这大大提升了新功能的可发现性和易用性！
