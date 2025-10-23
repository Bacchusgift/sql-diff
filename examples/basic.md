# 基础示例

本文档提供 SQL-Diff 的基础使用示例。

## 推荐：使用交互式模式

对于多行 SQL 或从数据库工具复制的语句，**强烈推荐使用交互式模式**：

```bash
# 启动交互式模式
sql-diff -i

# 按提示粘贴源表 SQL（支持多行）
# 粘贴完成后按 Ctrl+D（Mac/Linux）或 Ctrl+Z（Windows）
# 再粘贴目标表 SQL
# 自动生成 DDL！
```

### 交互式模式优势

✅ **支持多行 SQL** - 从 Navicat/MySQL Workbench 直接复制  
✅ **无需转义** - 不用处理换行符和引号  
✅ **实时反馈** - 字符统计和友好提示  
✅ **操作简单** - 一键启动，直接粘贴

---

## 简单的列添加

### 场景

需要在 `users` 表中添加 `email` 字段。

### 源表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 目标表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 执行命令

**方式 1：交互式模式（推荐）**

```bash
sql-diff -i
# 粘贴源表 SQL，按 Ctrl+D
# 粘贴目标表 SQL，再次 Ctrl+D
```

**方式 2：命令行模式**

```bash
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(100) NOT NULL);" \
  -t "CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(100) NOT NULL, email VARCHAR(255));"
```

### 输出

```
🔍 开始比对表结构...

表名: users
源表列数: 2, 目标表列数: 3

📋 生成的 DDL 语句:

➕ 新增列 (1):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255);

✅ 比对完成! 共生成 1 条 DDL 语句
```

### 生成的 DDL

```sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);
```

## 修改列定义

### 场景

扩展 `name` 字段长度,并添加 NOT NULL 约束。

### 源表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(50)
);
```

### 目标表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);
```

### 执行命令

```bash
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50));" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(255) NOT NULL);"
```

### 生成的 DDL

```sql
ALTER TABLE users MODIFY COLUMN name VARCHAR(255) NOT NULL;
```

## 添加索引

### 场景

为 `email` 字段添加唯一索引。

### 源表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  email VARCHAR(255)
);
```

### 目标表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  email VARCHAR(255),
  UNIQUE KEY uk_email (email)
);
```

### 执行命令

```bash
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, email VARCHAR(255));" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, email VARCHAR(255), UNIQUE KEY uk_email (email));"
```

### 输出

```
📋 生成的 DDL 语句:

➕ 新增索引 (1):
  1. ALTER TABLE users ADD UNIQUE KEY uk_email (email);
```

### 生成的 DDL

```sql
ALTER TABLE users ADD UNIQUE KEY uk_email (email);
```

## 删除列

### 场景

移除不再使用的 `old_field` 字段。

### 源表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100),
  old_field VARCHAR(50)
);
```

### 目标表

```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
);
```

### 执行命令

```bash
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), old_field VARCHAR(50));" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));"
```

### 输出

```
📋 生成的 DDL 语句:

❌ 删除列 (1):
  1. ALTER TABLE users DROP COLUMN old_field;

⚠️  警告: 删除列会导致数据永久丢失,请务必备份数据!
```

### 生成的 DDL

```sql
ALTER TABLE users DROP COLUMN old_field;
```

::: danger 数据丢失风险
删除列会永久删除该列的所有数据,执行前请务必备份!
:::

## 组合变更

### 场景

同时进行多种类型的变更。

### 源表

```sql
CREATE TABLE products (
  id INT PRIMARY KEY,
  name VARCHAR(100),
  price DECIMAL(10,2),
  old_status VARCHAR(20),
  KEY idx_old (old_status)
);
```

### 目标表

```sql
CREATE TABLE products (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  price DECIMAL(12,4),
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY idx_name (name)
);
```

### 执行命令

```bash
sql-diff \
  -s "$(cat source.sql)" \
  -t "$(cat target.sql)"
```

### 输出

```
🔍 开始比对表结构...

表名: products
源表列数: 4, 目标表列数: 5
差异统计: 新增 2 列, 修改 2 列, 删除 1 列

📋 生成的 DDL 语句:

🔄 修改列 (2):
  1. ALTER TABLE products MODIFY COLUMN id INT PRIMARY KEY AUTO_INCREMENT;
  2. ALTER TABLE products MODIFY COLUMN name VARCHAR(255) NOT NULL;
  3. ALTER TABLE products MODIFY COLUMN price DECIMAL(12,4);

➕ 新增列 (2):
  1. ALTER TABLE products ADD COLUMN description TEXT;
  2. ALTER TABLE products ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

❌ 删除列 (1):
  1. ALTER TABLE products DROP COLUMN old_status;

📊 索引变更:

❌ 删除索引 (1):
  1. ALTER TABLE products DROP INDEX idx_old;

➕ 新增索引 (1):
  1. ALTER TABLE products ADD KEY idx_name (name);

✅ 比对完成! 共生成 8 条 DDL 语句
```

## 从文件读取

### 准备文件

**source.sql**:
```sql
CREATE TABLE users (
  id INT PRIMARY KEY,
  name VARCHAR(100)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**target.sql**:
```sql
CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 执行命令

```bash
sql-diff -s "$(cat source.sql)" -t "$(cat target.sql)"
```

### 保存到文件

```bash
sql-diff -s "$(cat source.sql)" -t "$(cat target.sql)" > migration.sql
```

**migration.sql**:
```sql
ALTER TABLE users MODIFY COLUMN id INT PRIMARY KEY AUTO_INCREMENT;
ALTER TABLE users MODIFY COLUMN name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE users ADD UNIQUE KEY uk_email (email);
```

## 使用 AI 分析

### 执行命令

```bash
sql-diff \
  -s "CREATE TABLE users (id INT, name VARCHAR(50));" \
  -t "CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255) NOT NULL, email VARCHAR(255));" \
  --ai
```

### 输出 (包含 AI 分析)

```
🔍 开始比对表结构...

表名: users
源表列数: 2, 目标表列数: 3

📋 生成的 DDL 语句:

🔄 修改列 (2):
  1. ALTER TABLE users MODIFY COLUMN id INT PRIMARY KEY AUTO_INCREMENT;
  2. ALTER TABLE users MODIFY COLUMN name VARCHAR(255) NOT NULL;

➕ 新增列 (1):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255);

✅ 比对完成! 共生成 3 条 DDL 语句

🤖 AI 分析结果:

📊 差异分析:
检测到以下重要变更:
1. id 字段添加了主键和自增属性,这是数据库设计的最佳实践
2. name 字段扩展到 255 字符并添加 NOT NULL 约束,提高数据完整性
3. 新增 email 字段用于存储用户邮箱

💡 优化建议:
- 为 email 字段添加唯一索引: ALTER TABLE users ADD UNIQUE KEY uk_email (email)
- 考虑为 email 添加格式验证
- 建议添加 created_at 和 updated_at 时间戳字段
- name 字段可能不需要 255 字符,根据实际需求调整

⚠️  潜在风险:
- name 字段添加 NOT NULL 约束,确保现有数据没有空值
- 如果表已有数据,修改 id 为自增可能需要特殊处理
- email 字段允许 NULL,建议根据业务需求决定是否允许

✅ 最佳实践:
- 主键使用自增 ID 是常见做法
- 为邮箱等唯一字段添加索引
- 为字段添加注释说明用途
- 在测试环境充分验证后再应用到生产环境
```

## JSON 格式输出

### 执行命令

```bash
sql-diff \
  -s "CREATE TABLE users (id INT);" \
  -t "CREATE TABLE users (id INT, name VARCHAR(100));" \
  --format json
```

### 输出

```json
{
  "table_name": "users",
  "source_columns": 1,
  "target_columns": 2,
  "ddl_statements": [
    "ALTER TABLE users ADD COLUMN name VARCHAR(100);"
  ],
  "statistics": {
    "added_columns": 1,
    "modified_columns": 0,
    "dropped_columns": 0,
    "added_indexes": 0,
    "dropped_indexes": 0
  }
}
```

## 静默模式

### 执行命令

```bash
sql-diff \
  -s "CREATE TABLE users (id INT);" \
  -t "CREATE TABLE users (id INT, name VARCHAR(100));" \
  --quiet
```

### 输出 (仅 DDL)

```sql
ALTER TABLE users ADD COLUMN name VARCHAR(100);
```

适合在脚本中使用:

```bash
# 直接应用到数据库
sql-diff -s "..." -t "..." --quiet | mysql -h localhost -u user -p database
```

## 批量处理示例

### 脚本

```bash
#!/bin/bash
# batch_compare.sh

TABLES=(users products orders)

for table in "${TABLES[@]}"; do
  echo "=== Processing $table ==="
  
  sql-diff \
    -s "$(cat old/${table}.sql)" \
    -t "$(cat new/${table}.sql)" \
    --output "migrations/${table}_$(date +%Y%m%d).sql"
  
  echo "✓ Generated: migrations/${table}_$(date +%Y%m%d).sql"
  echo ""
done

echo "All migrations generated!"
```

### 执行

```bash
chmod +x batch_compare.sh
./batch_compare.sh
```

### 输出

```
=== Processing users ===
✓ Generated: migrations/users_20251022.sql

=== Processing products ===
✓ Generated: migrations/products_20251022.sql

=== Processing orders ===
✓ Generated: migrations/orders_20251022.sql

All migrations generated!
```

## 实际应用场景

### 场景 1: 开发环境同步

```bash
# 从开发数据库导出当前结构
mysqldump --no-data -h dev-db -u user -p myapp users > current.sql

# 比对新设计
sql-diff -s current.sql -t new_design.sql

# 生成迁移脚本
sql-diff -s current.sql -t new_design.sql > apply_to_dev.sql

# 应用到开发库
mysql -h dev-db -u user -p myapp < apply_to_dev.sql
```

### 场景 2: 代码审查

```bash
# 在 PR 中比对变更
git show main:db/schema/users.sql > old_schema.sql
git show HEAD:db/schema/users.sql > new_schema.sql

sql-diff -s old_schema.sql -t new_schema.sql --ai > schema_review.txt

# 审查变更
cat schema_review.txt
```

### 场景 3: 生产部署准备

```bash
#!/bin/bash
# prepare_prod_migration.sh

# 1. 导出生产当前结构
echo "Exporting production schema..."
mysqldump --no-data -h prod-db users > prod_current.sql

# 2. 生成迁移脚本和分析
echo "Generating migration..."
sql-diff \
  -s prod_current.sql \
  -t new_users_schema.sql \
  --ai > migration_analysis.txt

sql-diff \
  -s prod_current.sql \
  -t new_users_schema.sql \
  > migration.sql

# 3. 创建完整的部署包
mkdir -p deploy_$(date +%Y%m%d)
cp prod_current.sql deploy_$(date +%Y%m%d)/backup_schema.sql
cp migration.sql deploy_$(date +%Y%m%d)/
cp migration_analysis.txt deploy_$(date +%Y%m%d)/

echo "Deployment package created: deploy_$(date +%Y%m%d)"
echo "Please review migration_analysis.txt before deployment"
```

## 下一步

- [复杂场景示例](/examples/advanced) - 学习更高级的用法
- [CLI 工具](/guide/cli) - 了解所有命令选项
- [AI 功能](/ai/guide) - 使用 AI 增强分析
