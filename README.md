# SQL-Diff

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Status-Active-success?style=for-the-badge" alt="Status">
</p>

<p align="center">
  一个基于 AST 语法树的 SQL 表结构比对工具，支持交互式多行输入，自动生成 DDL 差异语句，可选接入 AI 大模型进行智能比对和优化建议。
</p>

<p align="center">
  <a href="https://bacchusgift.github.io/sql-diff/">📚 完整文档</a> ·
  <a href="https://bacchusgift.github.io/sql-diff/guide/getting-started">🚀 快速开始</a> ·
  <a href="https://bacchusgift.github.io/sql-diff/examples/basic">💡 示例</a> ·
  <a href="https://github.com/Bacchusgift/sql-diff/issues">💬 问题反馈</a>
</p>

---

## ✨ 特性

### 🎯 交互式输入（新）
支持多行 SQL 直接粘贴，完美解决换行符问题：
- ✅ 从 Navicat、MySQL Workbench 等工具直接复制
- ✅ 支持包含注释的复杂 SQL
- ✅ 无需处理换行符和转义字符
- ✅ 实时字符统计和友好提示

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

**🌐 完整文档网站：[https://bacchusgift.github.io/sql-diff/](https://bacchusgift.github.io/sql-diff/)**

- [🚀 快速开始](https://bacchusgift.github.io/sql-diff/guide/getting-started) - 5 分钟快速上手指南
- [💻 命令行工具](https://bacchusgift.github.io/sql-diff/guide/cli) - 详细的使用说明
- [🤖 AI 功能指南](https://bacchusgift.github.io/sql-diff/ai/guide) - AI 智能分析配置和使用
- [💡 使用示例](https://bacchusgift.github.io/sql-diff/examples/basic) - 实际应用场景
- [🏛️ 架构设计](https://bacchusgift.github.io/sql-diff/architecture) - 项目架构和设计思想
- [🤝 贡献指南](https://bacchusgift.github.io/sql-diff/CONTRIBUTING) - 如何为项目贡献

## 📦 安装

### macOS (Homebrew)

```bash
# 添加 tap
brew tap Bacchusgift/tap

# 安装
brew install sql-diff

# 或者一条命令
brew install Bacchusgift/tap/sql-diff
```

### Go Install

```bash
go install github.com/Bacchusgift/sql-diff@latest
```

### 从源码构建

```bash
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff
go build -o sql-diff cmd/sql-diff/main.go
```

## 🚀 快速开始

### 交互式模式（推荐）

对于多行 SQL 或从数据库工具复制的语句，交互式模式是最佳选择：

```bash
# 启动交互式模式
sql-diff -i

# 按提示操作：
# 1. 粘贴源表 SQL（支持多行）
# 2. 按 Ctrl+D（macOS/Linux）或 Ctrl+Z（Windows）结束输入
# 3. 粘贴目标表 SQL
# 4. 再次按 Ctrl+D
# 5. 自动生成 DDL！
```

**交互式 + AI 分析：**
```bash
sql-diff -i --ai
```

**交互式 + 输出到文件：**
```bash
sql-diff -i -o migration.sql
```

### 命令行参数模式

对于简单的单行 SQL，也可以使用命令行参数：

### 基础比对

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

### 本地开发

```bash
# 克隆项目
git clone https://github.com/Bacchusgift/sql-diff.git
cd sql-diff

# 安装依赖
go mod download

# 运行测试
make test

# 本地构建
make build

# 运行
./bin/sql-diff -i
```

### 可用的 Make 命令

```bash
make help          # 显示所有可用命令
make build         # 编译当前平台
make test          # 运行测试
make build-all     # 跨平台编译（所有平台）
make release       # 打包发布版本
make clean         # 清理构建产物
```

### 发布新版本

本项目使用 GitHub Actions 自动化发布流程，只需推送版本标签即可：

```bash
# 1. 确保代码已提交
git add .
git commit -m "feat: 新功能描述"
git push origin main

# 2. 创建并推送版本标签（推荐使用 v 前缀）
git tag v1.0.2
git push origin v1.0.2

# 或者不带 v 前缀也支持
git tag 1.0.2
git push origin 1.0.2
```

**自动化流程会完成：**

1. ✅ **跨平台编译** - 自动编译 6 个平台的二进制文件：
   - Linux (AMD64, ARM64)
   - macOS (Intel, Apple Silicon)
   - Windows (AMD64, ARM64)

2. ✅ **创建 GitHub Release** - 自动创建发布页面并上传：
   - 所有平台的二进制文件
   - SHA256 校验和文件
   - 自动生成的更新日志

3. ✅ **生成 Homebrew 更新信息** - 在 Release 评论中提供：
   - 更新后的 Formula 代码
   - SHA256 校验和
   - 详细的更新步骤

4. 🚧 **Homebrew Tap 更新** - 目前需要手动更新（后续将自动化）

### CI/CD 工作流

项目配置了两个主要的 GitHub Actions 工作流：

#### 1. CI 工作流 (`.github/workflows/ci.yml`)

每次推送到 `main` 或 `develop` 分支，或创建 PR 时触发：

- ✅ 代码格式检查 (`go fmt`)
- ✅ 代码质量检查 (`go vet`)
- ✅ 运行所有单元测试
- ✅ 6 平台编译验证

#### 2. Release 工作流 (`.github/workflows/release.yml`)

推送版本标签时触发（如 `v1.0.2`）：

- 🏗️ 跨平台编译
- 📦 创建 GitHub Release
- 📝 生成更新日志和安装说明
- 🔐 计算 SHA256 校验和
- 🍺 提供 Homebrew Formula 更新信息

### 手动更新 Homebrew Tap

发布新版本后，需要手动更新 Homebrew Tap 仓库（未来将自动化）：

```bash
# 1. 进入 homebrew-tap 仓库
cd ../homebrew-tap

# 2. 从 GitHub Release 评论中复制新的 Formula 代码
# 3. 编辑 Formula/sql-diff.rb，更新 url 和 sha256

# 4. 提交并推送
git add Formula/sql-diff.rb
git commit -m "chore: update sql-diff to v1.0.2"
git push origin main
```

详细说明请参考 [HOMEBREW.md](./HOMEBREW.md)

## 📝 License

MIT License
