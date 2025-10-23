# SQL-Diff 使用示例

本文档包含 sql-diff 工具的各种使用示例。

## 1. 基础用法

### 1.1 比对两个简单表结构

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))"
```

输出：
```sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);
```

### 1.2 比对复杂表结构

源表：
```sql
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

目标表（新增字段和索引）：
```sql
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    remark TEXT,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

运行命令：
```bash
./bin/sql-diff \
  -s "CREATE TABLE orders (...)" \
  -t "CREATE TABLE orders (...)"
```

## 2. 输出到文件

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (...)" \
  -t "CREATE TABLE users (...)" \
  -o migration.sql
```

生成的 `migration.sql` 文件可以直接在数据库中执行。

## 3. 使用 AI 增强功能

### 3.1 配置 AI

创建配置文件 `.sql-diff-config.yaml`：

```yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-your-api-key-here
  api_endpoint: https://api.deepseek.com/v1
  model: deepseek-chat
  timeout: 30
```

### 3.2 运行带 AI 分析的比对

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (...)" \
  -t "CREATE TABLE users (...)" \
  --ai
```

AI 会提供：
- 表结构变更的智能分析
- SQL 优化建议
- 潜在风险提示
- 最佳实践建议

## 4. 常见场景

### 4.1 新增字段

源表：
```sql
CREATE TABLE products (
    id INT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10,2)
)
```

目标表：
```sql
CREATE TABLE products (
    id INT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10,2),
    description TEXT,
    stock INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

生成的 DDL：
```sql
ALTER TABLE products ADD COLUMN description TEXT;
ALTER TABLE products ADD COLUMN stock INT DEFAULT 0;
ALTER TABLE products ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
```

### 4.2 修改字段属性

源表：
```sql
CREATE TABLE users (
    id INT,
    name VARCHAR(100),
    email VARCHAR(200)
)
```

目标表：
```sql
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL
)
```

生成的 DDL：
```sql
ALTER TABLE users MODIFY COLUMN name VARCHAR(100) NOT NULL;
ALTER TABLE users MODIFY COLUMN email VARCHAR(255) NOT NULL;
```

### 4.3 新增索引

源表：
```sql
CREATE TABLE posts (
    id INT PRIMARY KEY,
    title VARCHAR(200),
    content TEXT,
    author_id INT,
    created_at TIMESTAMP
)
```

目标表：
```sql
CREATE TABLE posts (
    id INT PRIMARY KEY,
    title VARCHAR(200),
    content TEXT,
    author_id INT,
    created_at TIMESTAMP,
    INDEX idx_author (author_id),
    INDEX idx_created (created_at)
)
```

生成的 DDL：
```sql
ALTER TABLE posts ADD INDEX idx_author (author_id);
ALTER TABLE posts ADD INDEX idx_created (created_at);
```

## 5. 高级用法

### 5.1 从文件读取 SQL

```bash
# 将 SQL 保存到文件
echo "CREATE TABLE ..." > source.sql
echo "CREATE TABLE ..." > target.sql

# 使用命令替换读取文件内容
./bin/sql-diff -s "$(cat source.sql)" -t "$(cat target.sql)"
```

### 5.2 指定自定义配置文件

```bash
./bin/sql-diff \
  -s "CREATE TABLE ..." \
  -t "CREATE TABLE ..." \
  --config ./custom-config.yaml \
  --ai
```

### 5.3 管道使用

```bash
# 生成 DDL 并直接执行
./bin/sql-diff -s "..." -t "..." | mysql -u root -p database_name
```

## 6. 注意事项

1. **删除操作**：工具生成的删除列和删除索引的 DDL 会被注释掉，避免误删数据
2. **数据类型**：确保 SQL 语句格式正确，支持标准 MySQL DDL 语法
3. **AI 功能**：需要有效的 API Key，并且网络可以访问 API 端点
4. **备份**：在执行任何 DDL 之前，请务必备份数据库

## 7. 故障排查

### 解析失败

如果遇到解析失败，请检查：
- SQL 语句是否完整且格式正确
- 是否使用了标准的 CREATE TABLE 语法
- 特殊字符是否正确转义

### AI 功能不工作

如果 AI 功能无法使用：
- 检查配置文件路径是否正确
- 验证 API Key 是否有效
- 确认网络可以访问 API 端点
- 查看错误信息获取详细原因
