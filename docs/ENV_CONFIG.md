# 环境变量配置指南

SQL-Diff 支持通过环境变量配置 AI 功能，无需管理配置文件。

## 📋 支持的环境变量

| 环境变量 | 说明 | 示例值 |
|---------|------|--------|
| `SQL_DIFF_AI_ENABLED` | 是否启用 AI | `true` / `false` |
| `SQL_DIFF_AI_PROVIDER` | AI 提供商 | `deepseek` / `openai` |
| `SQL_DIFF_AI_API_KEY` | API 密钥 | `sk-xxx...` |
| `SQL_DIFF_AI_ENDPOINT` | API 端点 | `https://api.deepseek.com/v1` |
| `SQL_DIFF_AI_MODEL` | 使用的模型 | `deepseek-chat` |
| `SQL_DIFF_AI_TIMEOUT` | 超时时间（秒） | `30` |

## 🚀 快速配置

### 方法 1: 使用配置命令（推荐）

```bash
# 生成配置并保存到 ~/.bashrc
sql-diff config \
  --ai-enabled \
  --provider deepseek \
  --api-key sk-your-api-key-here \
  >> ~/.bashrc

# 立即生效
source ~/.bashrc

# 验证配置
sql-diff config --show
```

### 方法 2: 使用配置向导

```bash
# 运行交互式配置向导
chmod +x setup-env.sh
./setup-env.sh >> ~/.bashrc
source ~/.bashrc
```

### 方法 3: 手动配置

```bash
# 编辑 ~/.bashrc 或 ~/.zshrc
vim ~/.bashrc

# 添加以下内容
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_PROVIDER=deepseek
export SQL_DIFF_AI_API_KEY=sk-your-api-key-here
export SQL_DIFF_AI_ENDPOINT=https://api.deepseek.com/v1
export SQL_DIFF_AI_MODEL=deepseek-chat
export SQL_DIFF_AI_TIMEOUT=30

# 生效
source ~/.bashrc
```

### 方法 4: 临时使用（不保存）

```bash
# 仅在当前终端会话生效
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_PROVIDER=deepseek
export SQL_DIFF_AI_API_KEY=sk-xxx
export SQL_DIFF_AI_ENDPOINT=https://api.deepseek.com/v1
export SQL_DIFF_AI_MODEL=deepseek-chat

# 或使用 eval
eval "$(sql-diff config --ai-enabled --provider deepseek --api-key sk-xxx)"
```

## 🔍 查看配置

```bash
# 查看当前环境变量配置
sql-diff config --show
```

输出示例:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       当前环境变量配置
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ SQL_DIFF_AI_ENABLED = true
✓ SQL_DIFF_AI_PROVIDER = deepseek
✓ SQL_DIFF_AI_API_KEY = sk-b50...89c6
✓ SQL_DIFF_AI_ENDPOINT = https://api.deepseek.com/v1
✓ SQL_DIFF_AI_MODEL = deepseek-chat
✓ SQL_DIFF_AI_TIMEOUT = 30

✓ 已检测到环境变量配置

📋 最终生效的配置:
  AI 启用状态: true
  AI 提供商:   deepseek
  API Key:     sk-b50...89c6
  API 端点:    https://api.deepseek.com/v1
  模型:        deepseek-chat
  超时时间:    30 秒
```

## ⚙️ 配置优先级

配置加载顺序（后者覆盖前者）:

1. **默认配置** - 内置的默认值
2. **配置文件** - `.sql-diff-config.yaml`
3. **环境变量** - `SQL_DIFF_AI_*` (优先级最高)

这意味着:
- 如果设置了环境变量，将覆盖配置文件
- 环境变量优先级最高，最灵活
- 配置文件作为备选方案

## 🎯 实际使用

### 使用环境变量运行

```bash
# 使用 AI 分析
sql-diff -s "CREATE TABLE users (...)" -t "CREATE TABLE users (...)" --ai

# 不使用 AI（即使环境变量已配置）
sql-diff -s "CREATE TABLE users (...)" -t "CREATE TABLE users (...)"
```

### CI/CD 集成

在 GitHub Actions / GitLab CI 中:

```yaml
env:
  SQL_DIFF_AI_ENABLED: true
  SQL_DIFF_AI_PROVIDER: deepseek
  SQL_DIFF_AI_API_KEY: ${{ secrets.DEEPSEEK_API_KEY }}
  SQL_DIFF_AI_ENDPOINT: https://api.deepseek.com/v1
  SQL_DIFF_AI_MODEL: deepseek-chat

steps:
  - name: Analyze Schema Changes
    run: |
      sql-diff -s "$(cat old.sql)" -t "$(cat new.sql)" --ai
```

### Docker 环境

```bash
docker run -e SQL_DIFF_AI_ENABLED=true \
           -e SQL_DIFF_AI_PROVIDER=deepseek \
           -e SQL_DIFF_AI_API_KEY=sk-xxx \
           sql-diff:latest \
           -s "..." -t "..." --ai
```

## 🔐 安全建议

### 保护 API Key

1. **不要硬编码** API Key 到脚本中
2. **使用环境变量** 存储敏感信息
3. **权限控制** ~/.bashrc 文件权限设为 600
4. **定期轮换** API Key

```bash
# 设置文件权限
chmod 600 ~/.bashrc

# 使用 .env 文件（不要提交到 git）
echo "export SQL_DIFF_AI_API_KEY=sk-xxx" > ~/.sql-diff.env
chmod 600 ~/.sql-diff.env
source ~/.sql-diff.env
```

### 团队使用

```bash
# 每个开发者使用自己的 API Key
# ~/.bashrc
export SQL_DIFF_AI_API_KEY=${MY_DEEPSEEK_KEY}

# 在项目中提供配置模板
# .env.example
SQL_DIFF_AI_ENABLED=true
SQL_DIFF_AI_PROVIDER=deepseek
SQL_DIFF_AI_API_KEY=your-api-key-here
```

## 💡 常见场景

### 场景 1: 开发环境

```bash
# 开发时启用 AI
export SQL_DIFF_AI_ENABLED=true
sql-diff -s "..." -t "..." --ai
```

### 场景 2: 生产环境

```bash
# 生产环境禁用 AI，使用配置文件
export SQL_DIFF_AI_ENABLED=false
sql-diff -s "..." -t "..."
```

### 场景 3: 多账号切换

```bash
# 使用不同的 API Key
export SQL_DIFF_AI_API_KEY=sk-project-a
sql-diff -s "..." -t "..." --ai

export SQL_DIFF_AI_API_KEY=sk-project-b
sql-diff -s "..." -t "..." --ai
```

## 🐛 故障排查

### 环境变量未生效

```bash
# 检查环境变量是否正确设置
echo $SQL_DIFF_AI_ENABLED
echo $SQL_DIFF_AI_API_KEY

# 查看完整配置
sql-diff config --show

# 确保 source 了配置文件
source ~/.bashrc
```

### 配置冲突

```bash
# 清除所有环境变量
unset SQL_DIFF_AI_ENABLED
unset SQL_DIFF_AI_PROVIDER
unset SQL_DIFF_AI_API_KEY
unset SQL_DIFF_AI_ENDPOINT
unset SQL_DIFF_AI_MODEL
unset SQL_DIFF_AI_TIMEOUT

# 重新配置
sql-diff config --ai-enabled --provider deepseek --api-key sk-xxx
```

## 📚 相关文档

- [AI 功能使用指南](AI_GUIDE.md)
- [快速开始](QUICKSTART.md)
- [配置文件说明](../README.md#配置文件)

---

**推荐方式**: 使用 `sql-diff config` 命令管理配置，简单直观！
