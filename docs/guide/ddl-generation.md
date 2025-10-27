# DDL 生成

SQL-Diff 会根据表结构差异自动生成标准的 MySQL DDL (Data Definition Language) 语句。

## DDL 语句分类

生成的 DDL 语句按类型分为五个类别,使用不同的颜色和图标显示:

### ✅ 新增列 (ADD COLUMN)

添加新列到表中:

```sql
ALTER TABLE table_name ADD COLUMN column_name datatype [constraints];
```

**示例**:
```sql
ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true;
```

### 🔄 修改列 (MODIFY COLUMN)

修改现有列的定义 (类型、约束等):

```sql
ALTER TABLE table_name MODIFY COLUMN column_name new_datatype [new_constraints];
```

**示例**:
```sql
ALTER TABLE products MODIFY COLUMN name VARCHAR(255) NOT NULL;
ALTER TABLE products MODIFY COLUMN price DECIMAL(12,4);
ALTER TABLE users MODIFY COLUMN status ENUM('active', 'inactive', 'suspended');
```

::: tip MODIFY vs CHANGE
SQL-Diff 使用 `MODIFY COLUMN` 来修改列定义,它保持列名不变。如果需要重命名列,需要手动使用 `CHANGE COLUMN`。
:::

### ❌ 删除列 (DROP COLUMN)

从表中删除列:

```sql
ALTER TABLE table_name DROP COLUMN column_name;
```

**示例**:
```sql
ALTER TABLE users DROP COLUMN old_field;
ALTER TABLE products DROP COLUMN deprecated_column;
```

::: danger 数据丢失警告
删除列会永久删除该列的所有数据,此操作不可撤销!执行前请务必备份数据。
:::

### ✅ 新增索引 (ADD INDEX)

添加新索引:

```sql
ALTER TABLE table_name ADD [INDEX_TYPE] KEY index_name (column_list);
```

**示例**:
```sql
ALTER TABLE users ADD KEY idx_email (email);
ALTER TABLE users ADD UNIQUE KEY uk_username (username);
ALTER TABLE posts ADD FULLTEXT KEY idx_content (title, content);
ALTER TABLE locations ADD SPATIAL KEY idx_coordinates (coordinates);
```

**支持的索引类型**:
- `KEY` - 普通索引
- `UNIQUE KEY` - 唯一索引
- `FULLTEXT KEY` - 全文索引
- `SPATIAL KEY` - 空间索引

### ❌ 删除索引 (DROP INDEX)

删除现有索引:

```sql
ALTER TABLE table_name DROP INDEX index_name;
```

**示例**:
```sql
ALTER TABLE users DROP INDEX idx_old_field;
ALTER TABLE posts DROP INDEX idx_deprecated;
```

## 输出格式

### 命令行输出

SQL-Diff 提供美观的分类输出:

```
🔍 开始比对表结构...

表名: users
源表列数: 3, 目标表列数: 5
差异统计: 新增 2 列, 修改 1 列, 删除 0 列

📋 生成的 DDL 语句:

➕ 新增列 (2):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;
  2. ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

🔄 修改列 (1):
  1. ALTER TABLE users MODIFY COLUMN name VARCHAR(255) NOT NULL;

📊 索引变更:

➕ 新增索引 (1):
  1. ALTER TABLE users ADD UNIQUE KEY uk_email (email);

✅ 比对完成! 共生成 4 条 DDL 语句
```

### 文件输出

使用重定向可以将 DDL 保存到文件:

```bash
sql-diff -s "..." -t "..." > migration.sql
```

生成的文件内容:

```sql
-- SQL-Diff Generated Migration
-- Table: users
-- Generated at: 2025-10-22 10:30:00

-- Add Columns
ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Modify Columns
ALTER TABLE users MODIFY COLUMN name VARCHAR(255) NOT NULL;

-- Add Indexes
ALTER TABLE users ADD UNIQUE KEY uk_email (email);
```

## 高级特性

### 1. 保持列顺序

对于需要在特定位置添加列的场景:

```sql
-- 在第一列之前添加
ALTER TABLE users ADD COLUMN id INT FIRST;

-- 在指定列之后添加
ALTER TABLE users ADD COLUMN middle_name VARCHAR(100) AFTER first_name;
```

::: tip
SQL-Diff 当前版本不会自动生成 `FIRST` 或 `AFTER` 子句,如需精确控制列顺序,请手动调整生成的 DDL。
:::

### 2. 组合多个变更

可以将多个变更合并到一个 ALTER TABLE 语句中:

```sql
-- 原始输出
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users DROP COLUMN old_field;

-- 手动优化为单个语句
ALTER TABLE users
  ADD COLUMN email VARCHAR(255),
  ADD COLUMN phone VARCHAR(20),
  DROP COLUMN old_field;
```

### 3. 在线 DDL 优化

对于大表,建议使用 MySQL 5.6+ 的在线 DDL 功能:

```sql
-- 使用 ALGORITHM 和 LOCK 子句
ALTER TABLE large_table 
  ADD COLUMN new_field VARCHAR(100),
  ALGORITHM=INPLACE, 
  LOCK=NONE;
```

## 执行 DDL 的最佳实践

### 1. 测试环境验证

```bash
# 1. 生成 DDL
sql-diff -s source.sql -t target.sql > migration.sql

# 2. 在测试环境执行
mysql -h testdb.example.com -u user -p database < migration.sql

# 3. 验证结果
mysql -h testdb.example.com -u user -p -e "DESCRIBE users"
```

### 2. 备份数据

```bash
# 执行变更前备份
mysqldump -h prod.example.com -u user -p database > backup_$(date +%Y%m%d).sql

# 执行变更
mysql -h prod.example.com -u user -p database < migration.sql
```

### 3. 使用事务 (如果支持)

```sql
START TRANSACTION;

-- 执行 DDL 语句
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- 验证结果
SELECT COUNT(*) FROM users;

-- 如果正确则提交,否则回滚
COMMIT;
-- ROLLBACK;
```

::: warning DDL 事务限制
MySQL 的大部分 DDL 语句会导致隐式提交,无法回滚。只有在某些存储引擎和特定条件下才支持 DDL 事务。
:::

### 4. 监控执行时间

对于大表的结构变更:

```sql
-- 查看当前正在执行的 DDL
SHOW PROCESSLIST;

-- 使用 pt-online-schema-change (Percona Toolkit)
pt-online-schema-change \
  --alter "ADD COLUMN email VARCHAR(255)" \
  D=database,t=users \
  --execute
```

## 性能考虑

### DDL 操作的性能影响

不同类型的 DDL 操作对性能的影响不同:

| 操作类型 | 是否锁表 | 性能影响 | 建议 |
|---------|---------|---------|------|
| ADD COLUMN (末尾) | 否 (5.6+) | 低 | 可在线执行 |
| ADD COLUMN (中间) | 是 | 高 | 低峰期执行 |
| MODIFY COLUMN | 取决于类型 | 中-高 | 谨慎执行 |
| DROP COLUMN | 否 (5.6+) | 低 | 可在线执行 |
| ADD INDEX | 否 (5.6+) | 中 | 可在线执行 |
| DROP INDEX | 否 | 低 | 可在线执行 |

### 大表优化策略

对于百万级以上数据的表:

1. **使用 pt-online-schema-change**
2. **分批执行变更**
3. **选择低峰期时段**
4. **监控主从延迟**
5. **预留足够磁盘空间**

## 常见问题

### Q: 如何撤销 DDL 操作?

A: 大部分 DDL 操作不可直接撤销,建议:
- 执行前备份数据
- 准备回滚 DDL (例如 ADD 对应 DROP)
- 在测试环境充分测试

### Q: DDL 执行失败怎么办?

A: 检查以下几点:
- 数据类型兼容性
- 磁盘空间是否充足
- 是否有权限限制
- 是否有外键约束冲突

### Q: 如何优化 DDL 执行速度?

A: 可以:
- 使用 `ALGORITHM=INPLACE`
- 对大表使用在线 DDL 工具
- 在从库先执行,然后主从切换
- 考虑使用 Ghost 等工具

## 下一步

- [命令行工具](/guide/cli) - 了解所有 CLI 选项
- [AI 功能](/ai/guide) - 使用 AI 优化 DDL
- [示例](/examples/advanced) - 查看复杂场景示例
