# 配置命令

SQL-Diff 提供了便捷的配置命令,用于管理 AI 功能的配置。

## 基础用法

```bash
sql-diff config [options]
```

## 配置选项

### AI 功能配置

| 选项 | 说明 | 示例 |
|------|------|------|
| `--ai-enabled` | 启用/禁用 AI | `--ai-enabled=true` |
| `--api-key` | API 密钥 | `--api-key="sk-xxx"` |
| `--provider` | AI 提供商 | `--provider=deepseek` |
| `--model` | 模型名称 | `--model=deepseek-chat` |
| `--api-url` | API 地址 | `--api-url="https://..."` |

### 查看配置

| 选项 | 说明 |
|------|------|
| `--show` | 显示当前配置 |

## 使用示例

### 1. 首次配置

配置 AI 功能:

```bash
sql-diff config \
  --ai-enabled=true \
  --api-key="sk-your-api-key" \
  --provider=deepseek \
  --model=deepseek-chat
```

输出:

```bash
export SQL_DIFF_AI_ENABLED=true
export SQL_DIFF_AI_API_KEY=sk-your-api-key
export SQL_DIFF_AI_PROVIDER=deepseek
export SQL_DIFF_AI_MODEL=deepseek-chat
```

### 2. 保存到环境变量

将配置保存到 shell 配置文件:

```bash
# 生成环境变量并保存到文件
sql-diff config \
  --ai-enabled=true \
  --api-key="sk-xxx" > ~/.sql-diff-env

# 在 shell 配置中引用
echo "source ~/.sql-diff-env" >> ~/.zshrc

# 立即生效
source ~/.zshrc
```

### 3. 查看当前配置

```bash
sql-diff config --show
```

输出:

```
📋 当前配置:

AI 配置:
✓ AI 功能: 已启用
✓ API Key: sk-xxx******
✓ 提供商: deepseek
✓ 模型: deepseek-chat

配置来源:
环境变量:
  - SQL_DIFF_AI_ENABLED=true
  - SQL_DIFF_AI_API_KEY=sk-***
  - SQL_DIFF_AI_PROVIDER=deepseek
  - SQL_DIFF_AI_MODEL=deepseek-chat
```

### 4. 更新配置

更新 API Key:

```bash
sql-diff config --api-key="new-api-key"
```

禁用 AI:

```bash
sql-diff config --ai-enabled=false
```

### 5. 切换提供商

切换到其他 AI 提供商:

```bash
# 使用 OpenAI
sql-diff config \
  --provider=openai \
  --api-key="sk-openai-key" \
  --model=gpt-4

# 使用自定义提供商
sql-diff config \
  --provider=custom \
  --api-url="https://your-api.com/v1/chat/completions" \
  --api-key="your-key" \
  --model="your-model"
```

## 配置工作流

### 个人开发环境

```bash
#!/bin/bash
# setup-dev.sh

echo "Setting up SQL-Diff for development..."

# 配置 AI
sql-diff config \
  --ai-enabled=true \
  --api-key="sk-dev-key" \
  --provider=deepseek > ~/.sql-diff-env

# 添加到 shell 配置
if ! grep -q "sql-diff-env" ~/.zshrc; then
  echo "source ~/.sql-diff-env" >> ~/.zshrc
fi

echo "✅ Configuration saved to ~/.sql-diff-env"
echo "Run 'source ~/.zshrc' to apply changes"
```

### 团队环境

```bash
#!/bin/bash
# team-setup.sh

echo "=== SQL-Diff Team Setup ==="
echo ""
echo "This script will help you configure SQL-Diff for your team."
echo ""

# 提示输入 API Key
read -p "Enter your DeepSeek API Key: " API_KEY

# 配置
sql-diff config \
  --ai-enabled=true \
  --api-key="$API_KEY" \
  --provider=deepseek \
  --model=deepseek-chat > ~/.sql-diff-env

# 设置权限
chmod 600 ~/.sql-diff-env

# 添加到 shell 配置
SHELL_RC="${HOME}/.zshrc"
if [ ! -f "$SHELL_RC" ]; then
  SHELL_RC="${HOME}/.bashrc"
fi

if ! grep -q "sql-diff-env" "$SHELL_RC"; then
  echo "" >> "$SHELL_RC"
  echo "# SQL-Diff Configuration" >> "$SHELL_RC"
  echo "[ -f ~/.sql-diff-env ] && source ~/.sql-diff-env" >> "$SHELL_RC"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Run: source $SHELL_RC"
echo "2. Test: sql-diff config --show"
echo "3. Try:  sql-diff -s old.sql -t new.sql --ai"
```

### CI/CD 环境

```yaml
# .github/workflows/schema-check.yml
steps:
  - name: Configure SQL-Diff
    run: |
      sql-diff config \
        --ai-enabled=true \
        --api-key="${{ secrets.DEEPSEEK_API_KEY }}" \
        --provider=deepseek \
        --model=deepseek-chat > .sql-diff-env
      
      source .sql-diff-env
  
  - name: Run Schema Check
    run: sql-diff -s old.sql -t new.sql --ai
```

## 高级用法

### 1. 条件配置

根据环境自动配置:

```bash
#!/bin/bash
# auto-config.sh

ENV=${1:-dev}

case $ENV in
  dev)
    API_KEY=$DEV_API_KEY
    MODEL="deepseek-chat"
    ;;
  staging)
    API_KEY=$STAGING_API_KEY
    MODEL="deepseek-chat"
    ;;
  prod)
    API_KEY=$PROD_API_KEY
    MODEL="deepseek-chat"
    ;;
esac

sql-diff config \
  --ai-enabled=true \
  --api-key="$API_KEY" \
  --model="$MODEL"
```

### 2. 配置验证

验证配置是否正确:

```bash
#!/bin/bash
# verify-config.sh

echo "Verifying SQL-Diff configuration..."

# 检查配置
sql-diff config --show

# 测试 AI 功能
echo ""
echo "Testing AI functionality..."
RESULT=$(sql-diff \
  -s "CREATE TABLE test (id INT);" \
  -t "CREATE TABLE test (id INT, name VARCHAR(100));" \
  --ai 2>&1)

if echo "$RESULT" | grep -q "AI 分析"; then
  echo "✅ AI functionality working!"
else
  echo "❌ AI functionality failed!"
  echo "$RESULT"
  exit 1
fi
```

### 3. 多配置管理

管理多个配置文件:

```bash
#!/bin/bash
# config-manager.sh

PROFILE=${1:-default}
CONFIG_DIR=~/.sql-diff-profiles

mkdir -p $CONFIG_DIR

case $1 in
  save)
    # 保存当前配置
    sql-diff config \
      --ai-enabled=true \
      --api-key="$2" > $CONFIG_DIR/$PROFILE
    echo "✅ Profile '$PROFILE' saved"
    ;;
  
  load)
    # 加载配置
    if [ -f "$CONFIG_DIR/$PROFILE" ]; then
      source $CONFIG_DIR/$PROFILE
      echo "✅ Profile '$PROFILE' loaded"
    else
      echo "❌ Profile '$PROFILE' not found"
    fi
    ;;
  
  list)
    # 列出所有配置
    echo "Available profiles:"
    ls -1 $CONFIG_DIR
    ;;
esac
```

使用:

```bash
# 保存配置
./config-manager.sh save dev sk-dev-key

# 加载配置
./config-manager.sh load dev

# 列出所有配置
./config-manager.sh list
```

## 故障排查

### 配置不生效

```bash
# 1. 检查环境变量
env | grep SQL_DIFF

# 2. 查看当前配置
sql-diff config --show

# 3. 重新配置
sql-diff config --ai-enabled=true --api-key="sk-xxx"

# 4. 刷新环境
source ~/.zshrc
```

### 权限问题

```bash
# 检查配置文件权限
ls -la ~/.sql-diff-env

# 修复权限
chmod 600 ~/.sql-diff-env

# 重新生成
sql-diff config --ai-enabled=true --api-key="sk-xxx" > ~/.sql-diff-env
```

### API Key 无效

```bash
# 测试 API Key
curl -X POST https://api.deepseek.com/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"deepseek-chat","messages":[{"role":"user","content":"test"}]}'

# 如果失败,重新配置
sql-diff config --api-key="new-correct-key"
```

## 最佳实践

1. **使用独立的配置文件**: 不要直接修改 `.zshrc`,使用独立的 `.sql-diff-env`
2. **设置正确的权限**: 配置文件应该是 `600` 权限
3. **不要提交到版本控制**: 将配置文件添加到 `.gitignore`
4. **定期轮换 API Key**: 每 90 天更新一次
5. **使用环境特定的配置**: 开发、测试、生产使用不同的配置

## 下一步

- [环境变量配置](/config/environment) - 详细的环境变量说明
- [配置文件](/config/file) - 使用配置文件
- [快速开始](/guide/getting-started) - 开始使用 SQL-Diff
