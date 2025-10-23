# 快速开始

欢迎使用 SQL-Diff！本指南将帮助您在 5 分钟内开始使用。

## 安装

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff

# 下载依赖
go mod download

# 编译
make build

# 验证安装
./bin/sql-diff --help
```

### 使用 Go Install

```bash
go install github.com/Bacchusgift/sql-diff@latest
```

## 第一次使用

### 1. 基础比对

比对两个简单的表结构：

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))"
```

输出：

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
新增列: 1 个
  + email VARCHAR

✓ 生成的 DDL 语句:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

➕ 新增列 (1):
  1. ALTER TABLE users ADD COLUMN email VARCHAR(255);

📋 完整执行脚本:
ALTER TABLE users ADD COLUMN email VARCHAR(255);
```

### 2. 输出到文件

将生成的 DDL 保存到文件：

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))" \
  -o migration.sql
```

### 3. 配置 AI 功能（可选）

```bash
# 方法 1: 使用环境变量（推荐）
sql-diff config \
  --ai-enabled \
  --provider deepseek \
  --api-key YOUR_API_KEY \
  >> ~/.bashrc

source ~/.bashrc

# 方法 2: 使用配置文件
cp .sql-diff-config.example.yaml .sql-diff-config.yaml
# 编辑文件填入 API Key
```

### 4. 使用 AI 分析

```bash
./bin/sql-diff \
  -s "CREATE TABLE users (id INT)" \
  -t "CREATE TABLE users (id INT, email VARCHAR(255))" \
  --ai
```

## 常用命令

```bash
# 查看帮助
sql-diff --help

# 查看配置
sql-diff config --show

# 运行演示
make run-demo

# 运行测试
make test
```

## 下一步

- 📖 阅读[完整文档](./introduction.md)
- 🔧 了解[配置选项](/config/environment.md)
- 🤖 探索 [AI 功能](/ai/guide.md)
- 💡 查看[使用示例](/examples/basic.md)
