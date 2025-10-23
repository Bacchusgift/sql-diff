# 命令行工具

SQL-Diff 提供了强大的命令行界面,支持多种使用场景。

## 基本命令

### 比对表结构

使用 `-s` (source) 和 `-t` (target) 参数:

```bash
sql-diff -s "CREATE TABLE users (id INT);" -t "CREATE TABLE users (id INT, name VARCHAR(100));"
```

### 从文件读取

```bash
sql-diff -s "$(cat source.sql)" -t "$(cat target.sql)"
```

或使用子命令:

```bash
sql-diff -s @source.sql -t @target.sql
```

## 命令选项

### 主要选项

| 选项 | 简写 | 说明 | 示例 |
|------|------|------|------|
| `--source` | `-s` | 源表 SQL 语句 | `-s "CREATE TABLE..."` |
| `--target` | `-t` | 目标表 SQL 语句 | `-t "CREATE TABLE..."` |
| `--ai` | | 启用 AI 分析 | `--ai` |
| `--help` | `-h` | 显示帮助信息 | `-h` |
| `--version` | `-v` | 显示版本号 | `-v` |

### 输出选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `--output` | 输出到文件 | `--output migration.sql` |
| `--format` | 输出格式 (text/json) | `--format json` |
| `--quiet` | 静默模式 | `--quiet` |
| `--verbose` | 详细输出 | `--verbose` |

### AI 选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `--ai` | 启用 AI 分析 | `--ai` |
| `--ai-provider` | AI 提供商 | `--ai-provider deepseek` |

## 配置命令

### 设置配置

```bash
# 配置 AI 功能
sql-diff config \
  --ai-enabled=true \
  --api-key="sk-xxx" \
  --provider=deepseek \
  --model="deepseek-chat"
```

### 查看配置

```bash
# 查看当前配置
sql-diff config --show

# 输出示例:
# AI 配置:
# ✓ AI 功能: 已启用
# ✓ API Key: sk-xxx***
# ✓ 提供商: deepseek
# ✓ 模型: deepseek-chat
```

### 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `--ai-enabled` | 启用 AI | false |
| `--api-key` | API 密钥 | - |
| `--provider` | AI 提供商 | deepseek |
| `--model` | 模型名称 | deepseek-chat |
| `--api-url` | API 地址 | - |
| `--show` | 显示当前配置 | - |

## 使用示例

### 1. 基础比对

```bash
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(255), email VARCHAR(255));"
```

输出:
```
🔍 开始比对表结构...

表名: users
源表列数: 2, 目标表列数: 3

📋 生成的 DDL 语句:

🔄 修改列 (1):
  1. ALTER TABLE users MODIFY COLUMN name VARCHAR(255);

➕ 新增列 (1):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255);

✅ 比对完成! 共生成 2 条 DDL 语句
```

### 2. 使用 AI 分析

首先配置 AI:

```bash
sql-diff config --ai-enabled=true --api-key="sk-xxx"
```

然后启用 AI 分析:

```bash
sql-diff -s "..." -t "..." --ai
```

输出会包含 AI 建议:

```
🤖 AI 分析结果:

📊 差异分析:
检测到 1 处列修改和 1 处新增列...

💡 优化建议:
- 建议为 email 字段添加唯一索引
- 考虑添加 created_at 和 updated_at 时间戳字段

⚠️  潜在风险:
- name 字段扩容可能导致索引重建
- 建议在低峰期执行

✅ 最佳实践:
- 为新字段设置默认值
- 添加适当的注释
```

### 3. 从文件读取并输出到文件

```bash
sql-diff \
  -s "$(cat tables/source/users.sql)" \
  -t "$(cat tables/target/users.sql)" \
  --output migrations/users_001.sql
```

### 4. JSON 格式输出

```bash
sql-diff -s "..." -t "..." --format json
```

输出:
```json
{
  "table_name": "users",
  "source_columns": 2,
  "target_columns": 3,
  "ddl_statements": [
    "ALTER TABLE users MODIFY COLUMN name VARCHAR(255);",
    "ALTER TABLE users ADD COLUMN email VARCHAR(255);"
  ],
  "statistics": {
    "added_columns": 1,
    "modified_columns": 1,
    "dropped_columns": 0,
    "added_indexes": 0,
    "dropped_indexes": 0
  }
}
```

### 5. 静默模式

只输出 DDL 语句,不显示额外信息:

```bash
sql-diff -s "..." -t "..." --quiet
```

输出:
```sql
ALTER TABLE users MODIFY COLUMN name VARCHAR(255);
ALTER TABLE users ADD COLUMN email VARCHAR(255);
```

### 6. 批量处理

使用 shell 脚本批量处理多个表:

```bash
#!/bin/bash

TABLES=(users products orders)

for table in "${TABLES[@]}"; do
  echo "Processing $table..."
  
  sql-diff \
    -s "$(cat source/${table}.sql)" \
    -t "$(cat target/${table}.sql)" \
    --output "migrations/${table}_$(date +%Y%m%d).sql" \
    --ai
  
  echo "✓ Done: migrations/${table}_$(date +%Y%m%d).sql"
done
```

## 高级用法

### 1. 管道操作

```bash
# 直接应用到数据库
sql-diff -s "..." -t "..." --quiet | mysql -h localhost -u user -p database

# 与其他工具结合
sql-diff -s "..." -t "..." --format json | jq '.ddl_statements[]'
```

### 2. 环境变量

使用环境变量简化命令:

```bash
# 设置环境变量
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_API_KEY=sk-xxx

# 命令会自动使用环境变量
sql-diff -s "..." -t "..." --ai
```

### 3. 配置文件

创建 `.sql-diff-config.yaml`:

```yaml
ai:
  enabled: true
  provider: deepseek
  api_key: sk-xxx
  model: deepseek-chat
```

SQL-Diff 会自动读取该文件。

### 4. 调试模式

```bash
# 显示详细日志
sql-diff -s "..." -t "..." --verbose --ai

# 输出包含:
# - SQL 解析过程
# - 差异检测详情
# - AI API 调用日志
# - 完整的错误堆栈
```

## 退出码

SQL-Diff 使用标准退出码:

| 退出码 | 含义 |
|--------|------|
| 0 | 成功 |
| 1 | 一般错误 |
| 2 | 配置错误 |
| 3 | SQL 解析错误 |
| 4 | AI API 错误 |

使用示例:

```bash
sql-diff -s "..." -t "..."
if [ $? -eq 0 ]; then
  echo "比对成功"
else
  echo "比对失败"
fi
```

## 常见用例

### CI/CD 集成

```yaml
# .github/workflows/schema-check.yml
name: Schema Check

on: [pull_request]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Install SQL-Diff
        run: go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
      
      - name: Check Schema Changes
        run: |
          sql-diff \
            -s "$(cat db/schema/current.sql)" \
            -t "$(cat db/schema/new.sql)" \
            --ai \
            --output migration.sql
        env:
          SQL_DIFF_AI_API_KEY: ${{ secrets.DEEPSEEK_API_KEY }}
```

### 数据库迁移

```bash
# 生成迁移文件
sql-diff \
  -s "$(mysqldump --no-data -h prod -u user -p db table)" \
  -t "$(cat new_schema.sql)" \
  --output "migrations/$(date +%Y%m%d_%H%M%S)_update_table.sql"
```

### 代码审查

```bash
# 比对并生成报告
sql-diff -s old.sql -t new.sql --ai --format json > review.json

# 用 jq 提取关键信息
cat review.json | jq '{
  table: .table_name,
  changes: .statistics,
  risks: .ai_analysis.risks
}'
```

## 故障排查

### 命令找不到

```bash
# 检查安装路径
which sql-diff

# 添加到 PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### 权限错误

```bash
# 检查文件权限
ls -l $(which sql-diff)

# 添加执行权限
chmod +x $(which sql-diff)
```

### AI 功能不工作

```bash
# 检查配置
sql-diff config --show

# 测试 API 连接
curl -X POST https://api.deepseek.com/v1/chat/completions \
  -H "Authorization: Bearer $SQL_DIFF_AI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"deepseek-chat","messages":[{"role":"user","content":"test"}]}'
```

## 下一步

- [配置](/config/environment) - 详细配置说明
- [AI 功能](/ai/guide) - AI 功能使用指南
- [示例](/examples/basic) - 更多实际示例
