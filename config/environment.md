# 环境变量配置

SQL-Diff 支持通过环境变量配置 AI 功能，这是**推荐的配置方式**。

## 为什么使用环境变量？

相比配置文件，环境变量有以下优势：

- ✅ **更安全** - API Key 不会被提交到 Git
- ✅ **更灵活** - 轻松切换不同账号
- ✅ **更方便** - 一次配置全局使用
- ✅ **CI/CD 友好** - 集成简单

## 支持的环境变量

| 环境变量 | 说明 | 示例值 |
|---------|------|--------|
| `SQL_DIFF_AI_ENABLED` | 启用/禁用 AI | `true` / `false` |
| `SQL_DIFF_AI_PROVIDER` | AI 提供商 | `deepseek` / `openai` |
| `SQL_DIFF_AI_API_KEY` | API 密钥 | `sk-xxx...` |
| `SQL_DIFF_AI_ENDPOINT` | API 端点 | `https://api.deepseek.com/v1` |
| `SQL_DIFF_AI_MODEL` | 模型名称 | `deepseek-chat` |
| `SQL_DIFF_AI_TIMEOUT` | 超时时间（秒） | `30` |

## 快速配置

### 方法 1: 使用配置命令（推荐）

```bash
# 一条命令生成配置并保存到 ~/.bashrc
sql-diff config \
  --ai-enabled \
  --provider deepseek \
  --api-key YOUR_API_KEY \
  --endpoint https://api.deepseek.com/v1 \
  --model deepseek-chat \
  --timeout 30 \
  >> ~/.bashrc

# 立即生效
source ~/.bashrc

# 验证配置
sql-diff config --show
```

### 方法 2: 手动配置

编辑 `~/.bashrc` 或 `~/.zshrc`：

```bash
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_PROVIDER=deepseek
export SQL_DIFF_AI_API_KEY=sk-your-api-key-here
export SQL_DIFF_AI_ENDPOINT=https://api.deepseek.com/v1
export SQL_DIFF_AI_MODEL=deepseek-chat
export SQL_DIFF_AI_TIMEOUT=30
```

然后生效：

```bash
source ~/.bashrc
```

### 方法 3: 临时使用

仅在当前终端会话生效：

```bash
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_API_KEY=sk-xxx

sql-diff -s "..." -t "..." --ai
```

## 查看配置

使用配置命令查看当前配置：

```bash
sql-diff config --show
```

输出示例：

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

## CI/CD 集成

### GitHub Actions

```yaml
name: Schema Check
on: [pull_request]

jobs:
  check:
    runs-on: ubuntu-latest
    env:
      SQL_DIFF_AI_ENABLED: true
      SQL_DIFF_AI_PROVIDER: deepseek
      SQL_DIFF_AI_API_KEY: ${{ secrets.DEEPSEEK_API_KEY }}
      SQL_DIFF_AI_ENDPOINT: https://api.deepseek.com/v1
      SQL_DIFF_AI_MODEL: deepseek-chat
    
    steps:
      - uses: actions/checkout@v3
      - name: Analyze Schema Changes
        run: |
          sql-diff -s "$(cat old.sql)" -t "$(cat new.sql)" --ai
```

### GitLab CI

```yaml
schema-check:
  variables:
    SQL_DIFF_AI_ENABLED: "true"
    SQL_DIFF_AI_PROVIDER: "deepseek"
    SQL_DIFF_AI_API_KEY: $DEEPSEEK_API_KEY
  script:
    - sql-diff -s "$(cat old.sql)" -t "$(cat new.sql)" --ai
```

## 配置优先级

配置加载顺序（后者覆盖前者）：

1. **默认配置** - 内置默认值
2. **配置文件** - `.sql-diff-config.yaml`
3. **环境变量** - `SQL_DIFF_AI_*` ⭐ 优先级最高

## 安全建议

### 保护 API Key

```bash
# 设置文件权限
chmod 600 ~/.bashrc

# 或使用单独的环境文件
echo "export SQL_DIFF_AI_API_KEY=sk-xxx" > ~/.sql-diff.env
chmod 600 ~/.sql-diff.env
source ~/.sql-diff.env
```

### 不要硬编码

❌ **错误示例**：

```bash
# 不要在脚本中硬编码 API Key
./deploy.sh --api-key sk-xxx
```

✅ **正确示例**：

```bash
# 使用环境变量
export SQL_DIFF_AI_API_KEY=sk-xxx
./deploy.sh
```

## 故障排查

### 环境变量未生效

```bash
# 检查环境变量
echo $SQL_DIFF_AI_ENABLED
echo $SQL_DIFF_AI_API_KEY

# 确保 source 了配置文件
source ~/.bashrc

# 查看完整配置
sql-diff config --show
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
sql-diff config --ai-enabled --provider deepseek --api-key sk-xxx >> ~/.bashrc
source ~/.bashrc
```

## 下一步

- 📝 [配置文件方式](/config/file.md)
- 🔧 [配置命令详解](/config/command.md)
- 🤖 [AI 功能指南](/ai/guide.md)
