# 快速开始指南

这是一个 5 分钟的快速入门教程，帮助你快速上手 SQL-Diff 工具。

## 步骤 1: 构建项目

```bash
# 克隆项目（如果还没有）
cd /Users/youzi/CascadeProjects/sql-diff

# 下载依赖
go mod download

# 编译
make build

# 或者直接运行
go build -o bin/sql-diff cmd/sql-diff/main.go
```

## 步骤 2: 第一次使用

比对两个简单的表结构：

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))"
```

你会看到类似这样的输出：

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       SQL 表结构比对工具
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📖 正在解析源表结构...
✓ 源表: users (2 列)

📖 正在解析目标表结构...
✓ 目标表: users (3 列)

🔍 正在比对表结构...

📊 差异摘要:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
新增列: 1 个
  + email VARCHAR

✓ 生成的 DDL 语句:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ALTER TABLE users ADD COLUMN email VARCHAR(255);
```

## 步骤 3: 保存到文件

将生成的 DDL 保存到文件：

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))" \
  -o migration.sql
```

现在你可以在 `migration.sql` 文件中查看 DDL 语句，并在数据库中执行。

## 步骤 4: 运行演示

查看更多示例：

```bash
# 运行完整演示
make run-demo

# 或者
chmod +x examples/demo.sh
./examples/demo.sh
```

## 步骤 5: 配置 AI（可选）

如果你想使用 AI 增强功能：

1. 复制配置文件模板：

```bash
cp .sql-diff-config.example.yaml .sql-diff-config.yaml
```

2. 编辑配置文件，填入你的 API Key：

```yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-your-api-key-here  # 替换为你的 API Key
  api_endpoint: https://api.deepseek.com/v1
  model: deepseek-chat
```

3. 使用 AI 功能：

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (...)" \
  -t "CREATE TABLE users (...)" \
  --ai
```

## 常用命令

```bash
# 构建项目
make build

# 运行测试
make test

# 清理构建产物
make clean

# 安装到系统
make install

# 查看帮助
./bin/sql-diff --help
```

## 实际使用场景

### 场景 1: 开发环境同步到生产

你在开发环境修改了表结构，现在需要生成迁移脚本：

```bash
# 开发环境的表结构
DEV_SQL="CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    user_id INT NOT NULL,
    total DECIMAL(10,2),
    status VARCHAR(20),
    created_at TIMESTAMP,
    INDEX idx_user (user_id)
)"

# 生产环境的表结构
PROD_SQL="CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    user_id INT NOT NULL,
    total DECIMAL(10,2),
    created_at TIMESTAMP
)"

# 生成迁移脚本
./bin/sql-diff -s "$PROD_SQL" -t "$DEV_SQL" -o prod_migration.sql
```

### 场景 2: 代码审查

在代码审查时，快速查看表结构变更：

```bash
# 从 Git 获取变更前后的 SQL
git show HEAD~1:schema/users.sql > /tmp/old.sql
git show HEAD:schema/users.sql > /tmp/new.sql

# 比对差异
./bin/sql-diff -s "$(cat /tmp/old.sql)" -t "$(cat /tmp/new.sql)"
```

### 场景 3: 数据库重构

使用 AI 帮助你优化表结构：

```bash
./bin/sql-diff \
  -s "CREATE TABLE products (id INT, name TEXT, price FLOAT)" \
  -t "CREATE TABLE products (id INT PRIMARY KEY, name VARCHAR(200) NOT NULL, price DECIMAL(10,2) NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)" \
  --ai
```

AI 会分析变更并提供：
- 数据类型选择建议
- 索引优化建议
- 潜在性能影响
- 最佳实践建议

## 下一步

- 📖 阅读完整文档
- 🏗️ 了解架构设计
- 🐛 [报告问题](https://github.com/Bacchusgift/sql-diff/issues)
- 💡 [功能建议](https://github.com/Bacchusgift/sql-diff/issues)

## 常见问题

**Q: 支持哪些数据库？**
A: 目前主要支持 MySQL 语法，未来会支持 PostgreSQL、SQLite 等。

**Q: 生成的 DDL 可以直接执行吗？**
A: 建议先在测试环境验证。删除操作会被注释掉以避免误删。

**Q: AI 功能需要付费吗？**
A: 取决于你选择的 AI 提供商。DeepSeek 和 OpenAI 都有免费额度。

**Q: 如何报告 Bug？**
A: 请在 GitHub Issues 中提交，附上 SQL 语句和错误信息。

**Q: 可以在 CI/CD 中使用吗？**
A: 可以！工具支持命令行模式，适合集成到 CI/CD 流程中。

---

🎉 恭喜！你已经掌握了 SQL-Diff 的基本用法。开始使用吧！
