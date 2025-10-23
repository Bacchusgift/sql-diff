# 快速开始

欢迎使用 SQL-Diff！本指南将帮助您在 5 分钟内开始使用。

## 📦 安装

### 🍺 macOS (Homebrew) - 推荐

最简单、最快捷的安装方式：

```bash
# 一条命令安装
brew install Bacchusgift/tap/sql-diff

# 验证安装
sql-diff --version
```

**更新到最新版本：**
```bash
brew upgrade sql-diff
```

### 🐧 Linux / Windows

#### 方式 1: 下载预编译二进制文件（推荐）

从 [GitHub Releases](https://github.com/Bacchusgift/sql-diff/releases/latest) 下载对应平台：

**Linux:**
```bash
# AMD64
wget https://github.com/Bacchusgift/sql-diff/releases/latest/download/sql-diff-linux-amd64
chmod +x sql-diff-linux-amd64
sudo mv sql-diff-linux-amd64 /usr/local/bin/sql-diff

# ARM64
wget https://github.com/Bacchusgift/sql-diff/releases/latest/download/sql-diff-linux-arm64
chmod +x sql-diff-linux-arm64
sudo mv sql-diff-linux-arm64 /usr/local/bin/sql-diff
```

**Windows:**
1. 下载 [sql-diff-windows-amd64.exe](https://github.com/Bacchusgift/sql-diff/releases/latest/download/sql-diff-windows-amd64.exe)
2. 重命名为 `sql-diff.exe`
3. 添加到 PATH 环境变量

#### 方式 2: 使用 Go Install

如果已安装 Go 1.21+：

```bash
go install github.com/Bacchusgift/sql-diff/cmd/sql-diff@latest
```

#### 方式 3: 从源码构建

```bash
# 克隆仓库
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff

# 下载依赖
go mod download

# 编译
make build

# 安装（可选）
sudo mv bin/sql-diff /usr/local/bin/

# 验证安装
sql-diff --version
```

## 🚀 第一次使用

### 1. 交互式模式（推荐）

对于多行 SQL 或从数据库工具复制的语句，使用交互式模式最方便：

```bash
sql-diff -i
```

按照提示操作：
1. 粘贴源表的 CREATE TABLE 语句（可以是多行）
2. 按 **Ctrl+D**（macOS/Linux）或 **Ctrl+Z 然后 Enter**（Windows）结束输入
3. 粘贴目标表的 CREATE TABLE 语句
4. 再次按 **Ctrl+D** 完成

示例输出：
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       SQL 表结构比对工具 - 交互式模式
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 请粘贴源表的 CREATE TABLE 语句：
（粘贴完成后按 Ctrl+D 结束输入，macOS/Linux）

[粘贴你的 SQL...]
✓ 已读取 245 个字符

📋 请粘贴目标表的 CREATE TABLE 语句：
[粘贴你的 SQL...]
✓ 已读取 312 个字符

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       开始比对
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
...
```

### 2. 命令行参数模式

对于简单的单行 SQL，可以直接使用命令行参数：

```bash
sql-diff \
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

### 3. 输出到文件

交互式模式也支持输出到文件：

```bash
sql-diff -i -o migration.sql
```

或命令行模式：

```bash
sql-diff \
  -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
  -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255))" \
  -o migration.sql
```

### 4. 配置 AI 功能（可选）

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

### 5. 使用 AI 分析

交互式模式 + AI：

```bash
sql-diff -i --ai
```

或命令行模式：

```bash
sql-diff \
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
