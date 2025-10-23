# SQL-Diff

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Status-Active-success?style=for-the-badge" alt="Status">
</p>

<p align="center">
  一个基于 AST 语法树的 SQL 表结构比对工具，支持生成 DDL 差异语句，可选接入 AI 大模型进行智能比对和优化建议。
</p>

---

## ✨ 特性

### 🔍 精准比对
基于 AST 语法树解析，准确识别：
- ✅ 新增列
- ✅ 修改列（类型、长度、约束、默认值）
- ✅ 删除列（安全注释）
- ✅ 索引变更

### 🛠️ DDL 生成
自动生成标准 MySQL DDL 语句：
```sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users MODIFY COLUMN name VARCHAR(200) NOT NULL;
ALTER TABLE users ADD INDEX idx_email (email);
```

### 🤖 AI 增强
可选接入 DeepSeek 等大模型，提供：
- 💡 智能分析表结构差异
- 📊 SQL 优化建议
- ⚠️ 潜在风险提示
- ✨ 最佳实践建议

### 💻 CLI 友好
简洁美观的命令行界面：
- 🎨 彩色输出
- 📦 清晰的进度提示
- 📝 结构化的结果展示
- 💾 支持输出到文件

## 📚 文档

- [**快速开始**](docs/QUICKSTART.md) - 5 分钟快速上手指南
- [**使用示例**](docs/EXAMPLES.md) - 详细的使用示例
- [**环境变量配置**](docs/ENV_CONFIG.md) - 使用环境变量配置 AI 功能
- [**AI 功能指南**](docs/AI_GUIDE.md) - AI 智能分析配置和使用
- [**架构设计**](docs/ARCHITECTURE.md) - 项目架构和设计思想
- [**贡献指南**](CONTRIBUTING.md) - 如何为项目贡献

## 📦 安装

```bash
go install github.com/Bacchusgift/sql-diff@latest
```

或从源码构建：

```bash
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff
go build -o sql-diff cmd/sql-diff/main.go
```

## 🚀 快速开始

### 交互式模式（推荐）

对于包含多行的复杂 SQL 语句，推荐使用交互式模式：

```bash
# 启动交互式模式
sql-diff -i

# 交互式模式 + AI 分析
sql-diff -i --ai

# 交互式模式 + 输出到文件
sql-diff -i -o output.sql
```

操作流程：
1. 运行命令后，程序提示粘贴源表 SQL
2. 直接粘贴（支持多行），粘贴完成后按 **Ctrl+D**（macOS/Linux）或 **Ctrl+Z 然后 Enter**（Windows）
3. 接着粘贴目标表 SQL，再次按 **Ctrl+D**
4. 程序自动比对并显示结果

### 基础用法

比对两个表结构并生成 DDL 语句：

```bash
sql-diff -s "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))" \
         -t "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255), created_at TIMESTAMP)"
```

输出：

```sql
-- 需要执行的 DDL 语句
ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN created_at TIMESTAMP;
```

### 使用 AI 增强功能

**方法 1: 使用环境变量（推荐）**

```bash
# 配置 AI 功能
sql-diff config \
  --ai-enabled \
  --provider deepseek \
  --api-key your-api-key-here \
  >> ~/.bashrc

# 生效
source ~/.bashrc

# 验证配置
sql-diff config --show

# 使用 AI 分析
sql-diff -s "CREATE TABLE users (...)" -t "CREATE TABLE users (...)" --ai
```

**方法 2: 使用配置文件**

1. 创建配置文件 `.sql-diff-config.yaml`：

```yaml
ai:
  enabled: true
  provider: deepseek
  api_key: your-api-key-here
  api_endpoint: https://api.deepseek.com/v1
```

2. 运行带 AI 分析的比对：

```bash
sql-diff -s "CREATE TABLE users (...)" -t "CREATE TABLE users (...)" --ai
```

## 📖 使用说明

### 命令行参数

- `-i, --interactive`: 交互式模式（推荐，支持多行粘贴）
- `-s, --source`: 源表的 CREATE TABLE 语句
- `-t, --target`: 目标表的 CREATE TABLE 语句
- `--ai`: 启用 AI 智能分析（需要配置）
- `--config`: 指定配置文件路径（默认：`.sql-diff-config.yaml`）
- `-o, --output`: 输出文件路径（默认：输出到控制台）

### 配置文件

在项目根目录或用户目录创建 `.sql-diff-config.yaml`：

```yaml
ai:
  enabled: true
  provider: deepseek  # 支持 deepseek, openai 等
  api_key: sk-** (替换成你的）
  api_endpoint: https://api.deepseek.com/v1
  model: deepseek-chat
```

## 🔧 开发

```bash
# 安装依赖
go mod download

# 运行测试
go test ./...

# 构建
go build -o bin/sql-diff cmd/sql-diff/main.go
```

## 📝 License

MIT License
